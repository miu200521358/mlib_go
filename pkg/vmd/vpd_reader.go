package vmd

import (
	"bufio"
	"regexp"
	"strconv"
	"strings"

	"github.com/miu200521358/mlib_go/pkg/domain/core"
	"github.com/miu200521358/mlib_go/pkg/mmath"
	"github.com/miu200521358/mlib_go/pkg/mutils/mlog"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

// VMDリーダー
type VpdMotionReader struct {
	core.BaseReader[*VmdMotion]
	lines []string
}

func (r *VpdMotionReader) createModel(path string) *VmdMotion {
	model := NewVmdMotion(path)
	return model
}

// 指定されたパスのファイルからデータを読み込む
func (r *VpdMotionReader) ReadByFilepath(path string) (core.IHashModel, error) {
	// モデルを新規作成
	motion := r.createModel(path)

	hash, err := r.ReadHashByFilePath(path)
	if err != nil {
		mlog.E("ReadByFilepath.ReadHashByFilePath error: %v", err)
		return motion, err
	}
	motion.SetHash(hash)

	// ファイルを開く
	err = r.Open(path)
	if err != nil {
		mlog.E("ReadByFilepath.Open error: %v", err)
		return motion, err
	}

	err = r.readLines()
	if err != nil {
		mlog.E("ReadByFilepath.readLines error: %v", err)
		return motion, err
	}

	err = r.readHeader(motion)
	if err != nil {
		mlog.E("ReadByFilepath.readHeader error: %v", err)
		return motion, err
	}

	err = r.readData(motion)
	if err != nil {
		mlog.E("ReadByFilepath.readData error: %v", err)
		return motion, err
	}

	r.Close()

	return motion, nil
}

func (r *VpdMotionReader) ReadNameByFilepath(path string) (string, error) {
	// モデルを新規作成
	motion := r.createModel(path)

	// ファイルを開く
	err := r.Open(path)
	if err != nil {
		mlog.E("ReadNameByFilepath.Open error: %v", err)
		return "", err
	}

	err = r.readLines()
	if err != nil {
		mlog.E("ReadByFilepath.readLines error: %v", err)
		return "", err
	}

	err = r.readHeader(motion)
	if err != nil {
		mlog.E("ReadNameByFilepath.readHeader error: %v", err)
		return "", err
	}

	r.Close()

	return motion.ModelName, nil
}

func (r *VpdMotionReader) readLines() error {
	var lines []string

	sjisReader := transform.NewReader(r.File, japanese.ShiftJIS.NewDecoder())
	scanner := bufio.NewScanner(sjisReader)
	for scanner.Scan() {
		txt := scanner.Text()
		txt = strings.ReplaceAll(txt, "\t", "    ")
		lines = append(lines, txt)
	}
	r.lines = lines
	return scanner.Err()
}

func (r *VpdMotionReader) ReadText(line string, pattern *regexp.Regexp) ([]string, error) {
	matches := pattern.FindStringSubmatch(line)
	if len(matches) > 0 {
		return matches, nil
	}
	return nil, nil
}

func (r *VpdMotionReader) readHeader(motion *VmdMotion) error {
	signaturePattern := regexp.MustCompile(`Vocaloid Pose Data file`)
	modelNamePattern := regexp.MustCompile(`(.*)(\.osm;.*// 親ファイル名)`)

	// signature
	{
		matches, err := r.ReadText(r.lines[0], signaturePattern)
		if err != nil || len(matches) == 0 {
			mlog.E("readHeader.ReadText error: %v", err)
			return err
		}
	}

	// モデル名
	matches, err := r.ReadText(r.lines[2], modelNamePattern)
	if err != nil {
		mlog.E("readHeader.ReadText error: %v", err)
		return err
	}

	if len(matches) > 0 {
		motion.ModelName = matches[1]
	}

	return nil
}

func (r *VpdMotionReader) readData(motion *VmdMotion) error {
	boneStartPattern := regexp.MustCompile(`(?:.*)(?:{)(.*)`)
	bonePosPattern := regexp.MustCompile(`([+-]?\d+(?:\.\d+))(?:,)([+-]?\d+(?:\.\d+))(?:,)([+-]?\d+(?:\.\d+))(?:;)(?:.*trans.*)`)
	boneRotPattern := regexp.MustCompile(`([+-]?\d+(?:\.\d+))(?:,)([+-]?\d+(?:\.\d+))(?:,)([+-]?\d+(?:\.\d+))(?:,)([+-]?\d+(?:\.\d+))(?:;)(?:.*Quaternion.*)`)

	var bf *BoneFrame
	var boneName string
	for _, line := range r.lines {
		{
			// 括弧開始: ボーン名
			matches, err := r.ReadText(line, boneStartPattern)
			if err == nil && len(matches) > 0 {
				boneName = matches[1]
				bf = NewBoneFrame(0)
				continue
			}
		}
		{
			// ボーン位置
			matches, err := r.ReadText(line, bonePosPattern)
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
			matches, err := r.ReadText(line, boneRotPattern)
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
