// 指示: miu200521358
package main

import (
	"crypto/sha256"
	"encoding/csv"
	"encoding/hex"
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/miu200521358/mlib_go/pkg/adapter/io_model"
	"github.com/miu200521358/mlib_go/pkg/domain/model"
)

const progressInterval = 100

var generatedTimestampPattern = regexp.MustCompile(`[0-9]{8}_[0-9]{6}`)

type boneRecord struct {
	index       int
	name        string
	englishName string
}

type modelRecord struct {
	path        string
	format      string
	hash        string
	duplicateOf string
	bones       []boneRecord
	loadErr     error
}

type boneStat struct {
	modelCount    int
	englishCounts map[string]int
}

type collector struct {
	records       []*modelRecord
	uniqueRecords []*modelRecord
	failed        []*modelRecord
	byHash        map[string]*modelRecord
	stats         map[string]*boneStat
	targetCount   int
	excludedCount int
}

// main はモデルを収集し、ボーン一覧と使用頻度を CSV に書き出す。
func main() {
	root, output, err := parseArgs()
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "引数エラー: %v\n", err)
		flag.PrintDefaults()
		os.Exit(2)
	}

	result, err := collectModels(root)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "モデル走査に失敗しました: %v\n", err)
		os.Exit(1)
	}
	if err := writeOutputs(output, result); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "CSV出力に失敗しました: %v\n", err)
		os.Exit(1)
	}
	printSummary(result)
}

// parseArgs は入力ルートと出力ディレクトリを CLI から取得する。
func parseArgs() (string, string, error) {
	root := flag.String("root", "", "PMX/PMDを再帰走査する入力ディレクトリ")
	output := flag.String("out", "bone_names_out", "CSV出力先ディレクトリ")
	flag.Parse()
	if strings.TrimSpace(*root) == "" {
		return "", "", errors.New("-root は必須です")
	}
	return *root, *output, nil
}

// collectModels は入力ディレクトリを決定的な順序で走査してモデル情報を集める。
func collectModels(root string) (*collector, error) {
	rootPath, err := filepath.Abs(root)
	if err != nil {
		return nil, fmt.Errorf("入力ルートの絶対パス化に失敗: %w", err)
	}
	info, err := os.Stat(rootPath)
	if err != nil {
		return nil, fmt.Errorf("入力ルートを確認できません: %w", err)
	}
	if !info.IsDir() {
		return nil, fmt.Errorf("入力ルートがディレクトリではありません: %s", root)
	}

	result := &collector{
		byHash: make(map[string]*modelRecord),
		stats:  make(map[string]*boneStat),
	}
	repository := io_model.NewModelRepository()
	err = filepath.Walk(rootPath, func(path string, fileInfo os.FileInfo, walkErr error) error {
		if walkErr != nil {
			relPath := relativePath(rootPath, path)
			_, _ = fmt.Fprintf(os.Stderr, "警告: %s の走査に失敗しました: %v\n", relPath, walkErr)
			return nil
		}
		if fileInfo.IsDir() {
			return nil
		}
		ext := strings.ToLower(filepath.Ext(fileInfo.Name()))
		if ext != ".pmx" && ext != ".pmd" {
			return nil
		}
		result.targetCount++
		if result.targetCount%progressInterval == 0 {
			_, _ = fmt.Fprintf(os.Stderr, "進捗: %d ファイル処理\n", result.targetCount)
		}

		// 日付を配布元が付ける場合もあるため、ツール出力を示す形式だけをファイル名で除外する。
		if generatedTimestampPattern.MatchString(fileInfo.Name()) {
			result.excludedCount++
			return nil
		}
		result.processFile(path, relativePath(rootPath, path), strings.TrimPrefix(ext, "."), repository)
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("入力ルートの走査に失敗: %w", err)
	}
	return result, nil
}

// processFile はファイル内容をハッシュ化し、未登録モデルだけを読み込む。
func (c *collector) processFile(path, relPath, format string, repository *io_model.ModelRepository) {
	contents, err := os.ReadFile(path)
	if err != nil {
		c.recordFailure(relPath, format, "ファイル読み込みに失敗: "+err.Error())
		return
	}
	hashBytes := sha256.Sum256(contents)
	hash := hex.EncodeToString(hashBytes[:])
	if previous, ok := c.byHash[hash]; ok {
		record := &modelRecord{
			path:        relPath,
			format:      format,
			hash:        hash,
			duplicateOf: previous.path,
			bones:       previous.bones,
			loadErr:     previous.loadErr,
		}
		c.records = append(c.records, record)
		if record.loadErr != nil {
			c.failed = append(c.failed, record)
			_, _ = fmt.Fprintf(os.Stderr, "警告: %s の読み込みに失敗しました: %v\n", relPath, record.loadErr)
		}
		return
	}

	loaded, err := repository.Load(path)
	if err != nil {
		c.recordFailureWithHash(relPath, format, hash, err)
		return
	}
	modelData, ok := loaded.(*model.PmxModel)
	if !ok || modelData == nil {
		c.recordFailureWithHash(relPath, format, hash, errors.New("読み込み結果が *model.PmxModel ではありません"))
		return
	}

	record := &modelRecord{
		path:   relPath,
		format: format,
		hash:   hash,
		bones:  extractBones(modelData),
	}
	// 同一内容は一意モデル一件として扱い、統計の分母と各モデル内の名前を一度だけ増やす。
	c.records = append(c.records, record)
	c.uniqueRecords = append(c.uniqueRecords, record)
	c.byHash[hash] = record
	c.addStats(record)
}

// recordFailure は読み込み前の失敗を記録する。
func (c *collector) recordFailure(relPath, format, message string) {
	c.recordFailureWithHash(relPath, format, "", errors.New(message))
}

// recordFailureWithHash は失敗を models.csv と failed.csv の両方へ反映する。
func (c *collector) recordFailureWithHash(relPath, format, hash string, err error) {
	record := &modelRecord{path: relPath, format: format, hash: hash, loadErr: err}
	c.records = append(c.records, record)
	if hash != "" {
		c.byHash[hash] = record
	}
	c.failed = append(c.failed, record)
	_, _ = fmt.Fprintf(os.Stderr, "警告: %s の読み込みに失敗しました: %v\n", relPath, err)
}

// extractBones はモデルのボーンを CSV 用の値へ変換する。
func extractBones(modelData *model.PmxModel) []boneRecord {
	if modelData == nil || modelData.Bones == nil {
		return nil
	}
	bones := make([]boneRecord, 0, modelData.Bones.Len())
	for _, bone := range modelData.Bones.Values() {
		if bone == nil {
			continue
		}
		bones = append(bones, boneRecord{
			index:       bone.Index(),
			name:        bone.Name(),
			englishName: bone.EnglishName,
		})
	}
	return bones
}

// addStats は一意モデル内の同名ボーンを一度だけ使用頻度へ加える。
func (c *collector) addStats(record *modelRecord) {
	seenNames := make(map[string]struct{})
	for _, bone := range record.bones {
		if _, exists := seenNames[bone.name]; exists {
			continue
		}
		seenNames[bone.name] = struct{}{}
		stat, exists := c.stats[bone.name]
		if !exists {
			stat = &boneStat{englishCounts: make(map[string]int)}
			c.stats[bone.name] = stat
		}
		stat.modelCount++
		stat.englishCounts[bone.englishName]++
	}
}

// writeOutputs は指定された四つの BOM 付き CSV を生成する。
func writeOutputs(outputDir string, result *collector) error {
	if err := os.MkdirAll(outputDir, 0o755); err != nil {
		return fmt.Errorf("出力先ディレクトリ作成に失敗: %w", err)
	}
	if err := writeCSV(filepath.Join(outputDir, "models.csv"), []string{"path", "format", "hash", "duplicate_of", "bone_count"}, modelRows(result)); err != nil {
		return err
	}
	if err := writeCSV(filepath.Join(outputDir, "bones.csv"), []string{"path", "index", "name", "english_name"}, boneRows(result)); err != nil {
		return err
	}
	if err := writeCSV(filepath.Join(outputDir, "bone_name_stats.csv"), []string{"rank", "name", "model_count", "ratio", "top_english_name"}, statsRows(result)); err != nil {
		return err
	}
	if err := writeCSV(filepath.Join(outputDir, "failed.csv"), []string{"path", "error"}, failedRows(result)); err != nil {
		return err
	}
	return nil
}

// writeCSV は encoding/csv で一つの CSV を BOM 付きで書き出す。
func writeCSV(path string, header []string, rows [][]string) error {
	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("%s の作成に失敗: %w", path, err)
	}
	if _, err := file.Write([]byte{0xEF, 0xBB, 0xBF}); err != nil {
		_ = file.Close()
		return fmt.Errorf("%s の BOM 書き込みに失敗: %w", path, err)
	}
	writer := csv.NewWriter(file)
	if err := writer.Write(header); err != nil {
		_ = file.Close()
		return fmt.Errorf("%s のヘッダ書き込みに失敗: %w", path, err)
	}
	for _, row := range rows {
		if err := writer.Write(row); err != nil {
			_ = file.Close()
			return fmt.Errorf("%s の行書き込みに失敗: %w", path, err)
		}
	}
	writer.Flush()
	if err := writer.Error(); err != nil {
		_ = file.Close()
		return fmt.Errorf("%s の CSV 書き込みに失敗: %w", path, err)
	}
	if err := file.Close(); err != nil {
		return fmt.Errorf("%s のクローズに失敗: %w", path, err)
	}
	return nil
}

// modelRows は models.csv の行を走査順に作る。
func modelRows(result *collector) [][]string {
	rows := make([][]string, 0, len(result.records))
	for _, record := range result.records {
		rows = append(rows, []string{
			record.path,
			record.format,
			hashPrefix(record.hash),
			record.duplicateOf,
			fmt.Sprintf("%d", len(record.bones)),
		})
	}
	return rows
}

// boneRows は一意モデルごとの bones.csv の行を作る。
func boneRows(result *collector) [][]string {
	rows := make([][]string, 0)
	for _, record := range result.uniqueRecords {
		for _, bone := range record.bones {
			rows = append(rows, []string{
				record.path,
				fmt.Sprintf("%d", bone.index),
				bone.name,
				bone.englishName,
			})
		}
	}
	return rows
}

// statsRows は指定された順位規則で bone_name_stats.csv の行を作る。
func statsRows(result *collector) [][]string {
	stats := sortedStats(result)
	rows := make([][]string, 0, len(stats))
	modelCount := len(result.uniqueRecords)
	for index, item := range stats {
		ratio := 0.0
		if modelCount > 0 {
			ratio = float64(item.stat.modelCount) / float64(modelCount)
		}
		rows = append(rows, []string{
			fmt.Sprintf("%d", index+1),
			item.name,
			fmt.Sprintf("%d", item.stat.modelCount),
			fmt.Sprintf("%.4f", ratio),
			topEnglishName(item.stat),
		})
	}
	return rows
}

// failedRows は読み込み失敗一覧の行を作る。
func failedRows(result *collector) [][]string {
	rows := make([][]string, 0, len(result.failed))
	for _, record := range result.failed {
		rows = append(rows, []string{record.path, record.loadErr.Error()})
	}
	return rows
}

type sortedStat struct {
	name string
	stat *boneStat
}

// sortedStats は model_count 降順、名前昇順の統計一覧を返す。
func sortedStats(result *collector) []sortedStat {
	items := make([]sortedStat, 0, len(result.stats))
	for name, stat := range result.stats {
		items = append(items, sortedStat{name: name, stat: stat})
	}
	sort.Slice(items, func(i, j int) bool {
		if items[i].stat.modelCount != items[j].stat.modelCount {
			return items[i].stat.modelCount > items[j].stat.modelCount
		}
		return items[i].name < items[j].name
	})
	return items
}

// topEnglishName は出現数が最大の英名を決定し、同数は辞書順で固定する。
// 英名未設定のモデルが多数派になると空文字が選ばれて列が役に立たないため、空の英名は候補から外す。
func topEnglishName(stat *boneStat) string {
	bestName := ""
	bestCount := -1
	for name, count := range stat.englishCounts {
		if strings.TrimSpace(name) == "" {
			continue
		}
		if count > bestCount || (count == bestCount && name < bestName) {
			bestName = name
			bestCount = count
		}
	}
	return bestName
}

// printSummary は標準出力へ集計概要と上位 50 件を表示する。
func printSummary(result *collector) {
	_, _ = fmt.Fprintf(os.Stdout, "対象ファイル数: %d\n", result.targetCount)
	_, _ = fmt.Fprintf(os.Stdout, "除外ファイル数(タイムスタンプ): %d\n", result.excludedCount)
	_, _ = fmt.Fprintf(os.Stdout, "読み込み失敗数: %d\n", len(result.failed))
	_, _ = fmt.Fprintf(os.Stdout, "一意モデル数: %d\n", len(result.uniqueRecords))
	_, _ = fmt.Fprintf(os.Stdout, "ボーン名の異なり数: %d\n", len(result.stats))
	_, _ = fmt.Fprintln(os.Stdout, "上位 50 件 (rank\tname\tmodel_count\tratio):")
	stats := sortedStats(result)
	limit := len(stats)
	if limit > 50 {
		limit = 50
	}
	for index := 0; index < limit; index++ {
		item := stats[index]
		ratio := 0.0
		if len(result.uniqueRecords) > 0 {
			ratio = float64(item.stat.modelCount) / float64(len(result.uniqueRecords))
		}
		_, _ = fmt.Fprintf(os.Stdout, "%d\t%s\t%d\t%.4f\n", index+1, item.name, item.stat.modelCount, ratio)
	}
}

// relativePath は走査ルート基準の区切りを CSV 向けに統一する。
func relativePath(root, path string) string {
	relPath, err := filepath.Rel(root, path)
	if err != nil {
		return filepath.ToSlash(path)
	}
	return filepath.ToSlash(relPath)
}

// hashPrefix は CSV に出す SHA-256 の先頭 12 桁を返す。
func hashPrefix(hash string) string {
	if len(hash) <= 12 {
		return hash
	}
	return hash[:12]
}
