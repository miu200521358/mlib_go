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

func (r *VpdRepository) Save(overridePath string, data core.IHashModel, includeSystem bool) error {
	return nil
}

// 指定されたパスのファイルからデータを読み込む
func (r *VpdRepository) Load(path string) (core.IHashModel, error) {
	// モデルを新規作成
	motion := r.newFunc(path)

	hash, err := r.LoadHash(path)
	if err != nil {
		mlog.E("Load.LoadHash error: %v", err)
		return motion, err
	}
	motion.SetHash(hash)

	// ファイルを開く
	err = r.open(path)
	if err != nil {
		mlog.E("Load.Open error: %v", err)
		return motion, err
	}

	err = r.readLines()
	if err != nil {
		mlog.E("Load.readLines error: %v", err)
		return motion, err
	}

	err = r.loadHeader(motion)
	if err != nil {
		mlog.E("Load.readHeader error: %v", err)
		return motion, err
	}

	err = r.loadModel(motion)
	if err != nil {
		mlog.E("Load.readData error: %v", err)
		return motion, err
	}

	r.close()

	return motion, nil
}

func (r *VpdRepository) LoadName(path string) (string, error) {
	// モデルを新規作成
	motion := r.newFunc(path)

	// ファイルを開く
	err := r.open(path)
	if err != nil {
		mlog.E("LoadName.Open error: %v", err)
		return "", err
	}

	err = r.readLines()
	if err != nil {
		mlog.E("Load.readLines error: %v", err)
		return "", err
	}

	err = r.loadHeader(motion)
	if err != nil {
		mlog.E("LoadName.readHeader error: %v", err)
		return "", err
	}

	r.close()

	return motion.ModelName, nil
}

func (r *VpdRepository) readLines() error {
	var lines []string

	sjisReader := transform.NewReader(r.file, japanese.ShiftJIS.NewDecoder())
	scanner := bufio.NewScanner(sjisReader)
	for scanner.Scan() {
		txt := scanner.Text()
		txt = strings.ReplaceAll(txt, "\t", "    ")
		lines = append(lines, txt)
	}
	r.lines = lines
	return scanner.Err()
}

func (r *VpdRepository) readText(line string, pattern *regexp.Regexp) ([]string, error) {
	matches := pattern.FindStringSubmatch(line)
	if len(matches) > 0 {
		return matches, nil
	}
	return nil, nil
}

func (r *VpdRepository) loadHeader(motion *vmd.VmdMotion) error {
	signaturePattern := regexp.MustCompile(`Vocaloid Pose Data file`)
	modelNamePattern := regexp.MustCompile(`(.*)(\.osm;.*// 親ファイル名)`)

	// signature
	{
		matches, err := r.readText(r.lines[0], signaturePattern)
		if err != nil || len(matches) == 0 {
			mlog.E("readHeader.readText error: %v", err)
			return err
		}
	}

	// モデル名
	matches, err := r.readText(r.lines[2], modelNamePattern)
	if err != nil {
		mlog.E("readHeader.readText error: %v", err)
		return err
	}

	if len(matches) > 0 {
		motion.ModelName = matches[1]
	}

	return nil
}

func (r *VpdRepository) loadModel(motion *vmd.VmdMotion) error {
	boneStartPattern := regexp.MustCompile(`(?:.*)(?:{)(.*)`)
	bonePosPattern := regexp.MustCompile(`([+-]?\d+(?:\.\d+))(?:,)([+-]?\d+(?:\.\d+))(?:,)([+-]?\d+(?:\.\d+))(?:;)(?:.*trans.*)`)
	boneRotPattern := regexp.MustCompile(`([+-]?\d+(?:\.\d+))(?:,)([+-]?\d+(?:\.\d+))(?:,)([+-]?\d+(?:\.\d+))(?:,)([+-]?\d+(?:\.\d+))(?:;)(?:.*Quaternion.*)`)

	var bf *vmd.BoneFrame
	var boneName string
	for _, line := range r.lines {
		{
			// 括弧開始: ボーン名
			matches, err := r.readText(line, boneStartPattern)
			if err == nil && len(matches) > 0 {
				boneName = matches[1]
				bf = vmd.NewBoneFrame(0)
				continue
			}
		}
		{
			// ボーン位置
			matches, err := r.readText(line, bonePosPattern)
			if err == nil && len(matches) > 0 {
				x, _ := strconv.ParseFloat(matches[1], 64)
				y, _ := strconv.ParseFloat(matches[2], 64)
				z, _ := strconv.ParseFloat(matches[3], 64)
				bf.Position = &mmath.MVec3{x, y, z}
				continue
			}
		}
		{
			// ボーン角度
			matches, err := r.readText(line, boneRotPattern)
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
