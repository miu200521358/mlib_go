// 指示: miu200521358
package vpd

import (
	"bufio"
	"io"
	"regexp"
	"strconv"
	"strings"

	"github.com/miu200521358/mlib_go/pkg/adapter/io_common"
	"github.com/miu200521358/mlib_go/pkg/domain/mmath"
	"github.com/miu200521358/mlib_go/pkg/domain/motion"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

var (
	vpdSignaturePattern = regexp.MustCompile(`Vocaloid Pose Data file`)
	vpdModelNamePattern = regexp.MustCompile(`(.*)(\.osm;.*// 親ファイル名)`)
	vpdBonePosPattern   = regexp.MustCompile(`([+-]?\d+(?:\.\d+)?)(?:,)([+-]?\d+(?:\.\d+)?)(?:,)([+-]?\d+(?:\.\d+)?)(?:;)(?:.*trans.*)`)
	vpdBoneRotPattern   = regexp.MustCompile(`([+-]?\d+(?:\.\d+)?)(?:,)([+-]?\d+(?:\.\d+)?)(?:,)([+-]?\d+(?:\.\d+)?)(?:,)([+-]?\d+(?:\.\d+)?)(?:;)(?:.*Quaternion.*)`)
)

// vpdReader はVPD読み取り処理を表す。
type vpdReader struct {
	reader io.Reader
	lines  []string
}

// newVpdReader はvpdReaderを生成する。
func newVpdReader(r io.Reader) *vpdReader {
	return &vpdReader{reader: r}
}

// Read はVPDテキストを読み込む。
func (r *vpdReader) Read(motionData *motion.VmdMotion) error {
	if motionData == nil {
		return io_common.NewIoParseFailed("VPDモーションがnilです", nil)
	}
	if err := r.readLines(); err != nil {
		return err
	}
	if err := r.readHeader(motionData); err != nil {
		return err
	}
	if err := r.readBones(motionData); err != nil {
		return err
	}
	return nil
}

// readLines はVPDを行単位で読み込む。
func (r *vpdReader) readLines() error {
	if r.reader == nil {
		return io_common.NewIoParseFailed("VPD読み取り元がnilです", nil)
	}
	sjisReader := transform.NewReader(r.reader, japanese.ShiftJIS.NewDecoder())
	scanner := bufio.NewScanner(sjisReader)
	lines := make([]string, 0)
	for scanner.Scan() {
		line := strings.ReplaceAll(scanner.Text(), "\t", "    ")
		lines = append(lines, line)
	}
	if err := scanner.Err(); err != nil {
		return io_common.NewIoParseFailed("VPDテキストの読み取りに失敗しました", err)
	}
	r.lines = lines
	return nil
}

// readHeader は署名とモデル名を読み込む。
func (r *vpdReader) readHeader(motionData *motion.VmdMotion) error {
	if len(r.lines) < 3 {
		return io_common.NewIoParseFailed("VPDヘッダが不足しています", nil)
	}
	if !vpdSignaturePattern.MatchString(r.lines[0]) {
		return io_common.NewIoParseFailed("VPD署名が不正です", nil)
	}
	matches := vpdModelNamePattern.FindStringSubmatch(r.lines[2])
	if len(matches) < 2 {
		return io_common.NewIoParseFailed("VPDモデル名が不正です", nil)
	}
	motionData.SetName(strings.TrimSpace(matches[1]))
	return nil
}

// readBones はボーンブロックを読み込む。
func (r *vpdReader) readBones(motionData *motion.VmdMotion) error {
	var (
		boneName string
		frame    *motion.BoneFrame
	)
	for _, line := range r.lines {
		if name, ok := parseBoneStart(line); ok {
			boneName = name
			frame = motion.NewBoneFrame(motion.Frame(0))
			frame.Read = true
			continue
		}
		if frame == nil || boneName == "" {
			continue
		}
		if matches := vpdBonePosPattern.FindStringSubmatch(line); len(matches) >= 4 {
			pos, err := parseVec3(matches[1], matches[2], matches[3])
			if err != nil {
				return err
			}
			frame.Position = &pos
			continue
		}
		if matches := vpdBoneRotPattern.FindStringSubmatch(line); len(matches) >= 5 {
			quat, err := parseQuat(matches[1], matches[2], matches[3], matches[4])
			if err != nil {
				return err
			}
			frame.Rotation = &quat
			motionData.AppendBoneFrame(boneName, frame)
			frame = nil
			boneName = ""
		}
	}
	if frame != nil {
		return io_common.NewIoParseFailed("VPDボーンの回転が不足しています", nil)
	}
	return nil
}

func parseBoneStart(line string) (string, bool) {
	idx := strings.Index(line, "{")
	if idx < 0 {
		return "", false
	}
	name := strings.TrimSpace(line[idx+1:])
	if name == "" {
		return "", false
	}
	return name, true
}

func parseVec3(x, y, z string) (mmath.Vec3, error) {
	fx, err := strconv.ParseFloat(x, 64)
	if err != nil {
		return mmath.Vec3{}, io_common.NewIoParseFailed("VPD位置の読み取りに失敗しました", err)
	}
	fy, err := strconv.ParseFloat(y, 64)
	if err != nil {
		return mmath.Vec3{}, io_common.NewIoParseFailed("VPD位置の読み取りに失敗しました", err)
	}
	fz, err := strconv.ParseFloat(z, 64)
	if err != nil {
		return mmath.Vec3{}, io_common.NewIoParseFailed("VPD位置の読み取りに失敗しました", err)
	}
	vec := mmath.Vec3{}
	vec.X = fx
	vec.Y = fy
	vec.Z = fz
	return vec, nil
}

func parseQuat(x, y, z, w string) (mmath.Quaternion, error) {
	fx, err := strconv.ParseFloat(x, 64)
	if err != nil {
		return mmath.Quaternion{}, io_common.NewIoParseFailed("VPD回転の読み取りに失敗しました", err)
	}
	fy, err := strconv.ParseFloat(y, 64)
	if err != nil {
		return mmath.Quaternion{}, io_common.NewIoParseFailed("VPD回転の読み取りに失敗しました", err)
	}
	fz, err := strconv.ParseFloat(z, 64)
	if err != nil {
		return mmath.Quaternion{}, io_common.NewIoParseFailed("VPD回転の読み取りに失敗しました", err)
	}
	fw, err := strconv.ParseFloat(w, 64)
	if err != nil {
		return mmath.Quaternion{}, io_common.NewIoParseFailed("VPD回転の読み取りに失敗しました", err)
	}
	quat := mmath.NewQuaternionByValues(fx, fy, fz, fw)
	return quat, nil
}
