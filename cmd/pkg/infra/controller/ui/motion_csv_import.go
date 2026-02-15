// 指示: miu200521358
package ui

import (
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"

	"github.com/miu200521358/mlib_go/cmd/pkg/adapter/mpresenter/messages"
	"github.com/miu200521358/mlib_go/pkg/adapter/io_common"
	"github.com/miu200521358/mlib_go/pkg/adapter/io_csv"
	"github.com/miu200521358/mlib_go/pkg/adapter/io_motion"
	"github.com/miu200521358/mlib_go/pkg/domain/mmath"
	"github.com/miu200521358/mlib_go/pkg/domain/motion"
	"github.com/miu200521358/mlib_go/pkg/shared/base/merr"
)

const (
	motionCsvBoneColumnsMin  = 72
	motionCsvMorphColumnsMin = 3
)

const (
	motionCsvNoDataErrorID              = "15601"
	motionCsvHeaderNotFoundErrorID      = "15602"
	motionCsvFormatUnknownErrorID       = "15603"
	motionCsvColumnsInsufficientErrorID = "15604"
	motionCsvParseFailedErrorID         = "15605"
	motionCsvVmdSaveFailedErrorID       = "15606"
	motionCsvModelInvalidErrorID        = "95601"
)

var (
	motionCsvBoneFilePattern  = regexp.MustCompile(`(?i)_bone(_\d{8}_\d{6})?\.csv$`)
	motionCsvMorphFilePattern = regexp.MustCompile(`(?i)_morph(_\d{8}_\d{6})?\.csv$`)
	errMotionCsvNoData        = merr.NewCommonError(
		motionCsvNoDataErrorID,
		merr.ErrorKindValidate,
		messages.MessageMotionCsvImportNoData,
		nil,
	)
)

type motionCsvKind int

const (
	motionCsvKindUnknown motionCsvKind = iota
	motionCsvKindBone
	motionCsvKindMorph
)

// buildMotionVmdDefaultOutputPath はCSV入力パスからVMD出力先の既定値を生成する。
func buildMotionVmdDefaultOutputPath(inputPath string) string {
	if inputPath == "" {
		return ""
	}
	if motionCsvBoneFilePattern.MatchString(inputPath) {
		return motionCsvBoneFilePattern.ReplaceAllString(inputPath, motionVmdExt)
	}
	if motionCsvMorphFilePattern.MatchString(inputPath) {
		return motionCsvMorphFilePattern.ReplaceAllString(inputPath, motionVmdExt)
	}
	return strings.TrimSuffix(inputPath, filepath.Ext(inputPath)) + motionVmdExt
}

// importMotionCsvByInputPath はCSVファイルを読み込みVMDとして保存する。
func importMotionCsvByInputPath(inputPath string, outputPath string) error {
	motionData := motion.NewVmdMotion(outputPath)
	motionData.SetName(strings.TrimSuffix(filepath.Base(outputPath), filepath.Ext(outputPath)))

	kind, _, err := appendMotionCsvByPath(motionData, inputPath)
	if err != nil {
		return err
	}

	pairPath := inferMotionCsvPairPath(inputPath, kind)
	if pairPath != "" && !sameFilePath(inputPath, pairPath) && existsFile(pairPath) {
		if _, _, err := appendMotionCsvByPath(motionData, pairPath); err != nil {
			return err
		}
	}

	if motionData.BoneFrames.Len() == 0 && motionData.MorphFrames.Len() == 0 {
		return errMotionCsvNoData
	}

	motionData.UpdateHash()
	if err := io_motion.NewVmdVpdRepository().Save(outputPath, motionData, io_common.SaveOptions{}); err != nil {
		return merr.NewCommonError(
			motionCsvVmdSaveFailedErrorID,
			merr.ErrorKindValidate,
			messages.MessageMotionCsvVmdSaveFailed,
			err,
			filepath.Base(outputPath),
		)
	}
	return nil
}

// appendMotionCsvByPath はCSV1ファイルを読み込み、モーションへ追記する。
func appendMotionCsvByPath(motionData *motion.VmdMotion, path string) (motionCsvKind, int, error) {
	if motionData == nil {
		return motionCsvKindUnknown, 0, merr.NewCommonError(
			motionCsvModelInvalidErrorID,
			merr.ErrorKindInternal,
			messages.MessageMotionCsvMotionDataNil,
			nil,
		)
	}

	data, err := io_csv.NewCsvRepository().Load(path)
	if err != nil {
		return motionCsvKindUnknown, 0, err
	}
	model, ok := data.(*io_csv.CsvModel)
	if !ok || model == nil {
		return motionCsvKindUnknown, 0, merr.NewCommonError(
			motionCsvModelInvalidErrorID,
			merr.ErrorKindInternal,
			messages.MessageMotionCsvModelConvertFailed,
			nil,
			filepath.Base(path),
		)
	}

	records := model.Records()
	if len(records) == 0 || len(records[0]) == 0 {
		return motionCsvKindUnknown, 0, merr.NewCommonError(
			motionCsvHeaderNotFoundErrorID,
			merr.ErrorKindValidate,
			messages.MessageMotionCsvHeaderNotFound,
			nil,
			filepath.Base(path),
		)
	}

	kind := detectMotionCsvKind(records[0])
	switch kind {
	case motionCsvKindBone:
		if err := validateMotionCsvColumns(records, motionCsvBoneColumnsMin, "ボーンCSV", path); err != nil {
			return kind, 0, err
		}
		count, err := appendMotionBoneCsvRecords(motionData, model, path)
		return kind, count, err
	case motionCsvKindMorph:
		if err := validateMotionCsvColumns(records, motionCsvMorphColumnsMin, "モーフCSV", path); err != nil {
			return kind, 0, err
		}
		count, err := appendMotionMorphCsvRecords(motionData, model, path)
		return kind, count, err
	default:
		return motionCsvKindUnknown, 0, merr.NewCommonError(
			motionCsvFormatUnknownErrorID,
			merr.ErrorKindValidate,
			messages.MessageMotionCsvFormatUnknown,
			nil,
			filepath.Base(path),
		)
	}
}

// detectMotionCsvKind はヘッダからCSV種別を判定する。
func detectMotionCsvKind(header []string) motionCsvKind {
	if len(header) == 0 {
		return motionCsvKindUnknown
	}
	head := strings.TrimSpace(strings.TrimPrefix(header[0], "\ufeff"))
	switch head {
	case "ボーン名":
		return motionCsvKindBone
	case "モーフ名":
		return motionCsvKindMorph
	default:
		return motionCsvKindUnknown
	}
}

// validateMotionCsvColumns はCSV列数を検証する。
func validateMotionCsvColumns(records [][]string, minColumns int, label string, path string) error {
	if len(records) == 0 {
		return nil
	}
	if len(records[0]) < minColumns {
		return merr.NewCommonError(
			motionCsvColumnsInsufficientErrorID,
			merr.ErrorKindValidate,
			messages.MessageMotionCsvColumnsInsufficient,
			nil,
			label,
			filepath.Base(path),
		)
	}

	for rowIndex := 1; rowIndex < len(records); rowIndex++ {
		row := records[rowIndex]
		if isMotionCsvEmptyRow(row) {
			continue
		}
		if len(row) < minColumns {
			return merr.NewCommonError(
				motionCsvColumnsInsufficientErrorID,
				merr.ErrorKindValidate,
				messages.MessageMotionCsvColumnsInsufficientWithRow,
				nil,
				label,
				rowIndex+1,
				filepath.Base(path),
			)
		}
	}
	return nil
}

// appendMotionBoneCsvRecords はボーンCSVレコードをモーションへ追記する。
func appendMotionBoneCsvRecords(motionData *motion.VmdMotion, model *io_csv.CsvModel, path string) (int, error) {
	if model == nil {
		return 0, merr.NewCommonError(
			motionCsvModelInvalidErrorID,
			merr.ErrorKindInternal,
			messages.MessageMotionCsvModelNil,
			nil,
		)
	}

	rows := make([]motionBoneCsvRow, 0)
	if err := io_csv.UnmarshalWithOptions(model, &rows, io_csv.CsvUnmarshalOptions{
		ColumnMapping: io_csv.CsvColumnMappingOrder,
	}); err != nil {
		return 0, merr.NewCommonError(
			motionCsvParseFailedErrorID,
			merr.ErrorKindValidate,
			messages.MessageMotionCsvBoneParseFailed,
			err,
			filepath.Base(path),
		)
	}

	count := 0
	for _, row := range rows {
		boneName := strings.TrimSpace(row.BoneName)
		if boneName == "" {
			continue
		}

		frameIndex := row.Frame
		if frameIndex < 0 {
			frameIndex = 0
		}
		frame := motion.NewBoneFrame(motion.Frame(frameIndex))
		position := mmath.NewVec3()
		position.X = row.PositionX
		position.Y = row.PositionY
		position.Z = row.PositionZ
		frame.Position = &position
		rotation := mmath.NewQuaternionFromDegrees(row.RotationX, -row.RotationY, -row.RotationZ)
		frame.Rotation = &rotation
		interpolation := buildMotionBoneInterpolationValues(row)
		frame.Curves = motion.NewBoneCurvesByValues(interpolation[:])
		motionData.AppendBoneFrame(boneName, frame)
		count++
	}

	return count, nil
}

// appendMotionMorphCsvRecords はモーフCSVレコードをモーションへ追記する。
func appendMotionMorphCsvRecords(motionData *motion.VmdMotion, model *io_csv.CsvModel, path string) (int, error) {
	if model == nil {
		return 0, merr.NewCommonError(
			motionCsvModelInvalidErrorID,
			merr.ErrorKindInternal,
			messages.MessageMotionCsvModelNil,
			nil,
		)
	}

	rows := make([]motionMorphCsvRow, 0)
	if err := io_csv.UnmarshalWithOptions(model, &rows, io_csv.CsvUnmarshalOptions{
		ColumnMapping: io_csv.CsvColumnMappingHeader,
	}); err != nil {
		return 0, merr.NewCommonError(
			motionCsvParseFailedErrorID,
			merr.ErrorKindValidate,
			messages.MessageMotionCsvMorphParseFailed,
			err,
			filepath.Base(path),
		)
	}

	count := 0
	for _, row := range rows {
		morphName := strings.TrimSpace(row.MorphName)
		if morphName == "" {
			continue
		}
		frameIndex := row.Frame
		if frameIndex < 0 {
			frameIndex = 0
		}
		frame := motion.NewMorphFrame(motion.Frame(frameIndex))
		frame.Ratio = row.Size
		motionData.AppendMorphFrame(morphName, frame)
		count++
	}

	return count, nil
}

// buildMotionBoneInterpolationValues はCSV行からボーン補間値64要素を抽出する。
func buildMotionBoneInterpolationValues(row motionBoneCsvRow) [64]byte {
	values := [64]int{
		row.Interp00, row.Interp01, row.Interp02, row.Interp03,
		row.Interp04, row.Interp05, row.Interp06, row.Interp07,
		row.Interp08, row.Interp09, row.Interp10, row.Interp11,
		row.Interp12, row.Interp13, row.Interp14, row.Interp15,
		row.Interp16, row.Interp17, row.Interp18, row.Interp19,
		row.Interp20, row.Interp21, row.Interp22, row.Interp23,
		row.Interp24, row.Interp25, row.Interp26, row.Interp27,
		row.Interp28, row.Interp29, row.Interp30, row.Interp31,
		row.Interp32, row.Interp33, row.Interp34, row.Interp35,
		row.Interp36, row.Interp37, row.Interp38, row.Interp39,
		row.Interp40, row.Interp41, row.Interp42, row.Interp43,
		row.Interp44, row.Interp45, row.Interp46, row.Interp47,
		row.Interp48, row.Interp49, row.Interp50, row.Interp51,
		row.Interp52, row.Interp53, row.Interp54, row.Interp55,
		row.Interp56, row.Interp57, row.Interp58, row.Interp59,
		row.Interp60, row.Interp61, row.Interp62, row.Interp63,
	}
	var interpolation [64]byte
	for i, value := range values {
		interpolation[i] = clampMotionCsvInterpolation(value)
	}
	return interpolation
}

// clampMotionCsvInterpolation は補間値をbyte範囲へ丸める。
func clampMotionCsvInterpolation(value int) byte {
	if value < 0 {
		return 0
	}
	if value > 255 {
		return 255
	}
	return byte(value)
}

// inferMotionCsvPairPath は入力CSVに対応するペアCSVパスを返す。
func inferMotionCsvPairPath(path string, kind motionCsvKind) string {
	switch kind {
	case motionCsvKindBone:
		return motionCsvBoneFilePattern.ReplaceAllString(path, "_morph$1.csv")
	case motionCsvKindMorph:
		return motionCsvMorphFilePattern.ReplaceAllString(path, "_bone$1.csv")
	default:
		return ""
	}
}

// isMotionCsvEmptyRow は空行か判定する。
func isMotionCsvEmptyRow(row []string) bool {
	if len(row) == 0 {
		return true
	}
	for _, cell := range row {
		if strings.TrimSpace(cell) != "" {
			return false
		}
	}
	return true
}

// existsFile は通常ファイルの存在を判定する。
func existsFile(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return !info.IsDir()
}

// sameFilePath は同一パスか判定する。
func sameFilePath(left string, right string) bool {
	return sameFilePathByOS(runtime.GOOS, left, right)
}

// sameFilePathByOS はOSごとの比較方式で同一パスか判定する。
func sameFilePathByOS(osName string, left string, right string) bool {
	leftPath := filepath.Clean(left)
	rightPath := filepath.Clean(right)
	if osName == "windows" {
		return strings.EqualFold(leftPath, rightPath)
	}
	return leftPath == rightPath
}
