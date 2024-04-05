package vmd

import (
	"bytes"
	"encoding/binary"
	"io"
	"math"
	"os"

	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"

	"github.com/miu200521358/mlib_go/pkg/mutils/mi18n"
	"github.com/miu200521358/mlib_go/pkg/mutils/mlog"

)

type VmdWriter struct {
}

func NewVmdWriter() *VmdWriter {
	return &VmdWriter{}
}

func Write(motion *VmdMotion) error {
	// Open the output file
	fout, err := os.Create(motion.GetPath())
	if err != nil {
		return err
	}

	// Write the header
	header := []byte("Vocaloid Motion Data 0002\x00\x00\x00\x00\x00")
	_, err = fout.Write(header)
	if err != nil {
		return err
	}

	// Convert model name to shift_jis encoding
	modelBName, err := encodeName(motion.GetName(), 20)
	if err != nil {
		mlog.W(mi18n.T("モデル名エンコードエラー", map[string]interface{}{"Name": motion.GetName()}))
		modelBName = []byte("Vmd Model")
	}

	// Write the model name
	err = binary.Write(fout, binary.LittleEndian, modelBName)
	if err != nil {
		return err
	}

	// Write the bone frames
	err = writeBoneFrames(fout, motion)
	if err != nil {
		mlog.E(mi18n.T("ボーンフレーム書き込みエラー"))
		return err
	}

	// Write the morph frames
	err = writeMorphFrames(fout, motion)
	if err != nil {
		mlog.E(mi18n.T("モーフフレーム書き込みエラー"))
		return err
	}

	// foutを書き込んで終了する
	err = fout.Close()
	if err != nil {
		mlog.E(mi18n.T("ファイルクローズエラー", map[string]interface{}{"Path": motion.GetPath()}))
		return err
	}

	return nil
}

func writeBoneFrames(fout *os.File, motion *VmdMotion) error {

	names := motion.BoneFrames.GetNames()
	indexes := motion.BoneFrames.GetRegisteredIndexes()

	binary.Write(fout, binary.LittleEndian, uint32(motion.BoneFrames.GetCount()))
	for i := len(indexes) - 1; i >= 0; i-- {
		index := indexes[i]
		for _, name := range names {
			if motion.BoneFrames.Contains(name) && motion.BoneFrames.Data[name].Contains(index) {
				bf := motion.BoneFrames.Data[name].Data[index]
				encodedName, err := encodeName(name, 15)
				if err != nil {
					mlog.W(mi18n.T("ボーン名エンコードエラー", map[string]interface{}{"Name": name}))
					continue
				}

				binary.Write(fout, binary.LittleEndian, encodedName)
				binary.Write(fout, binary.LittleEndian, uint32(bf.Index))
				binary.Write(fout, binary.LittleEndian, float32(bf.Position.GetX()))
				binary.Write(fout, binary.LittleEndian, float32(bf.Position.GetY()))
				binary.Write(fout, binary.LittleEndian, float32(bf.Position.GetZ()))

				v := bf.Rotation.GetQuaternion().Normalized()
				binary.Write(fout, binary.LittleEndian, float32(v.GetX()))
				binary.Write(fout, binary.LittleEndian, float32(v.GetY()))
				binary.Write(fout, binary.LittleEndian, float32(v.GetZ()))
				binary.Write(fout, binary.LittleEndian, float32(v.GetW()))

				curves := make([]byte, len(bf.Curves.values))
				for i, x := range bf.Curves.Merge() {
					curves[i] = byte(math.Min(255, math.Max(0, float64(x))))
				}
				binary.Write(fout, binary.LittleEndian, curves)
			}
		}
	}

	return nil
}

func writeMorphFrames(fout *os.File, motion *VmdMotion) error {
	binary.Write(fout, binary.LittleEndian, uint32(motion.MorphFrames.GetCount()))

	// FIXME カメラ個数
	binary.Write(fout, binary.LittleEndian, uint32(0))
	// 照明 個数
	binary.Write(fout, binary.LittleEndian, uint32(0))
	// 影 個数
	binary.Write(fout, binary.LittleEndian, uint32(0))
	// IK 個数
	binary.Write(fout, binary.LittleEndian, uint32(0))

	return nil
}

func encodeName(name string, limit int) ([]byte, error) {
	// Encode to CP932
	cp932Encoder := japanese.ShiftJIS.NewEncoder()
	cp932Encoded, err := cp932Encoder.String(name)
	if err != nil {
		return []byte(""), err
	}

	// Decode to Shift_JIS
	shiftJISDecoder := japanese.ShiftJIS.NewDecoder()
	reader := transform.NewReader(bytes.NewReader([]byte(cp932Encoded)), shiftJISDecoder)
	shiftJISDecoded, err := io.ReadAll(reader)
	if err != nil {
		return []byte(""), err
	}

	// Encode to Shift_JIS
	shiftJISEncoder := japanese.ShiftJIS.NewEncoder()
	shiftJISEncoded, err := shiftJISEncoder.String(string(shiftJISDecoded))
	if err != nil {
		return []byte(""), err
	}

	encodedName := []byte(shiftJISEncoded)
	if len(encodedName) <= limit {
		// 指定バイト数に足りない場合は b"\x00" で埋める
		encodedName = append(encodedName, make([]byte, limit-len(encodedName))...)
	}

	// 指定バイト数に切り詰め
	return encodedName[:limit], nil
}
