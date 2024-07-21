package repository

import (
	"bufio"
	"regexp"
	"strconv"
	"strings"

	"github.com/miu200521358/mlib_go/pkg/domain/core"
	"github.com/miu200521358/mlib_go/pkg/domain/mmath"
	"github.com/miu200521358/mlib_go/pkg/domain/vmd"
	"github.com/miu200521358/mlib_go/pkg/mutils/mlog"
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
				return &vmd.VmdMotion{
					HashModel: core.NewHashModel(path),
				}
			},
		},
	}
}

func (rep *VpdRepository) Save(overridePath string, data core.IHashModel, includeSystem bool) error {
	return nil
}

// 指定されたパスのファイルからデータを読み込む
func (rep *VpdRepository) Load(path string) (core.IHashModel, error) {
	// モデルを新規作成
	motion := rep.newFunc(path)

	hash, err := rep.LoadHash(path)
	if err != nil {
		mlog.E("Load.LoadHash error: %v", err)
		return motion, err
	}
	motion.SetHash(hash)

	// ファイルを開く
	err = rep.open(path)
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

	rep.close()

	return motion, nil
}

func (rep *VpdRepository) LoadName(path string) (string, error) {
	// モデルを新規作成
	motion := rep.newFunc(path)

	// ファイルを開く
	err := rep.open(path)
	if err != nil {
		mlog.E("LoadName.Open error: %v", err)
		return "", err
	}

	err = rep.readLines()
	if err != nil {
		mlog.E("Load.readLines error: %v", err)
		return "", err
	}

	err = rep.loadHeader(motion)
	if err != nil {
		mlog.E("LoadName.readHeader error: %v", err)
		return "", err
	}

	rep.close()

	return motion.ModelName, nil
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
			mlog.E("readHeader.readText error: %v", err)
			return err
		}
	}

	// モデル名
	matches, err := rep.readText(rep.lines[2], modelNamePattern)
	if err != nil {
		mlog.E("readHeader.readText error: %v", err)
		return err
	}

	if len(matches) > 0 {
		motion.ModelName = matches[1]
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
