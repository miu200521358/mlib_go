package repository

import (
	"bufio"
	"fmt"
	"regexp"
	"runtime"
	"strconv"
	"strings"

	"github.com/miu200521358/mlib_go/pkg/config/mfile"
	"github.com/miu200521358/mlib_go/pkg/config/mi18n"
	"github.com/miu200521358/mlib_go/pkg/config/mlog"
	"github.com/miu200521358/mlib_go/pkg/config/mstring"
	"github.com/miu200521358/mlib_go/pkg/domain/core"
	"github.com/miu200521358/mlib_go/pkg/domain/mmath"
	"github.com/miu200521358/mlib_go/pkg/domain/vmd"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

type VpdRepository struct {
	*baseRepository[*vmd.VmdMotion]
	lines []string
}

func NewVpdRepository() *VpdRepository {
	return &VpdRepository{
		baseRepository: &baseRepository[*vmd.VmdMotion]{
			newFunc: func(path string) *vmd.VmdMotion {
				return vmd.NewVmdMotion(path)
			},
		},
	}
}

func (rep *VpdRepository) Save(overridePath string, data core.IHashModel, includeSystem bool) error {
	runtime.GOMAXPROCS(int(runtime.NumCPU()))
	defer runtime.GOMAXPROCS(max(1, int(runtime.NumCPU()/4)))

	mlog.IL(mi18n.T("保存開始", map[string]interface{}{"Type": "Vpd", "Path": overridePath}))
	defer mlog.I(mi18n.T("保存終了", map[string]interface{}{"Type": "Vpd"}))

	return nil
}

func (rep *VpdRepository) CanLoad(path string) (bool, error) {
	if isExist, err := mfile.ExistsFile(path); err != nil || !isExist {
		return false, fmt.Errorf(mi18n.T("ファイル存在エラー", map[string]interface{}{"Path": path}))
	}

	_, _, ext := mfile.SplitPath(path)
	if strings.ToLower(ext) != ".vpd" {
		return false, fmt.Errorf(mi18n.T("拡張子エラー", map[string]interface{}{"Path": path, "Ext": ".vpd"}))
	}

	return true, nil
}

// 指定されたパスのファイルからデータを読み込む
func (rep *VpdRepository) Load(path string) (core.IHashModel, error) {
	runtime.GOMAXPROCS(int(runtime.NumCPU()))
	defer runtime.GOMAXPROCS(max(1, int(runtime.NumCPU()/4)))

	mlog.IL(mi18n.T("読み込み開始", map[string]interface{}{"Type": "Csv", "Path": path}))
	defer mlog.I(mi18n.T("読み込み終了", map[string]interface{}{"Type": "Csv"}))

	// モデルを新規作成
	motion := rep.newFunc(path)

	// ファイルを開く
	err := rep.open(path)
	if err != nil {
		mlog.E("Load.Open error: %v", err)
		return motion, err
	}

	err = rep.readLines()
	if err != nil {
		mlog.E("Load.readLines error: %v", err)
		return motion, err
	}

	err = rep.loadHeader(motion)
	if err != nil {
		mlog.E("Load.readHeader error: %v", err)
		return motion, err
	}

	err = rep.loadModel(motion)
	if err != nil {
		mlog.E("Load.readData error: %v", err)
		return motion, err
	}

	motion.UpdateHash()
	rep.close()

	return motion, nil
}

func (rep *VpdRepository) LoadName(path string) string {
	if ok, err := rep.CanLoad(path); !ok || err != nil {
		return mi18n.T("読み込み失敗")
	}

	// モデルを新規作成
	motion := rep.newFunc(path)

	// ファイルを開く
	err := rep.open(path)
	if err != nil {
		return mi18n.T("読み込み失敗")
	}

	err = rep.readLines()
	if err != nil {
		return mi18n.T("読み込み失敗")
	}

	err = rep.loadHeader(motion)
	if err != nil {
		return mi18n.T("読み込み失敗")
	}

	rep.close()

	return motion.Name()
}

func (rep *VpdRepository) readLines() error {
	var lines []string

	sjisReader := transform.NewReader(rep.file, japanese.ShiftJIS.NewDecoder())
	scanner := bufio.NewScanner(sjisReader)
	for scanner.Scan() {
		txt := scanner.Text()
		txt = strings.ReplaceAll(txt, "\t", "    ")
		lines = append(lines, txt)
	}
	rep.lines = lines
	return scanner.Err()
}

func (rep *VpdRepository) readText(line string, pattern *regexp.Regexp) ([]string, error) {
	matches := pattern.FindStringSubmatch(line)
	if len(matches) > 0 {
		return matches, nil
	}
	return nil, nil
}

func (rep *VpdRepository) loadHeader(motion *vmd.VmdMotion) error {
	signaturePattern := regexp.MustCompile(`Vocaloid Pose Data file`)
	modelNamePattern := regexp.MustCompile(`(.*)(\.osm;.*// 親ファイル名)`)

	// signature
	{
		matches, err := rep.readText(rep.lines[0], signaturePattern)
		if err != nil || len(matches) == 0 {
			return fmt.Errorf("readHeader.readText error: %v\n\n%v", err, mstring.GetStackTrace())
		}
	}

	// モデル名
	matches, err := rep.readText(rep.lines[2], modelNamePattern)
	if err != nil {
		return fmt.Errorf("readHeader.readText error: %v\n\n%v", err, mstring.GetStackTrace())
	}

	if len(matches) > 0 {
		motion.SetName(matches[1])
	}

	return nil
}

func (rep *VpdRepository) loadModel(motion *vmd.VmdMotion) error {
	boneStartPattern := regexp.MustCompile(`(?:.*)(?:{)(.*)`)
	bonePosPattern := regexp.MustCompile(`([+-]?\d+(?:\.\d+))(?:,)([+-]?\d+(?:\.\d+))(?:,)([+-]?\d+(?:\.\d+))(?:;)(?:.*trans.*)`)
	boneRotPattern := regexp.MustCompile(`([+-]?\d+(?:\.\d+))(?:,)([+-]?\d+(?:\.\d+))(?:,)([+-]?\d+(?:\.\d+))(?:,)([+-]?\d+(?:\.\d+))(?:;)(?:.*Quaternion.*)`)

	var bf *vmd.BoneFrame
	var boneName string
	for _, line := range rep.lines {
		{
			// 括弧開始: ボーン名
			matches, err := rep.readText(line, boneStartPattern)
			if err == nil && len(matches) > 0 {
				boneName = matches[1]
				bf = vmd.NewBoneFrame(0)
				continue
			}
		}
		{
			// ボーン位置
			matches, err := rep.readText(line, bonePosPattern)
			if err == nil && len(matches) > 0 {
				x, _ := strconv.ParseFloat(matches[1], 64)
				y, _ := strconv.ParseFloat(matches[2], 64)
				z, _ := strconv.ParseFloat(matches[3], 64)
				bf.Position = &mmath.MVec3{X: x, Y: y, Z: z}
				continue
			}
		}
		{
			// ボーン角度
			matches, err := rep.readText(line, boneRotPattern)
			if err == nil && len(matches) > 0 {
				x, _ := strconv.ParseFloat(matches[1], 64)
				y, _ := strconv.ParseFloat(matches[2], 64)
				z, _ := strconv.ParseFloat(matches[3], 64)
				w, _ := strconv.ParseFloat(matches[4], 64)
				bf.Rotation = mmath.NewMQuaternionByValues(x, y, z, w)

				motion.AppendRegisteredBoneFrame(boneName, bf)
				continue
			}
		}
	}

	return nil
}
