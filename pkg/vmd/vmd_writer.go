package vmd

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"os"

	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"

	"github.com/miu200521358/mlib_go/pkg/mmath"
	"github.com/miu200521358/mlib_go/pkg/mutils"
	"github.com/miu200521358/mlib_go/pkg/mutils/mi18n"
	"github.com/miu200521358/mlib_go/pkg/mutils/mlog"
)

func (motion *VmdMotion) Save(overrideModelName, overridePath string) error {
	path := motion.GetPath()
	// 保存可能なパスである場合、上書き
	if mutils.CanSave(overridePath) {
		path = overridePath
	}

	modelName := motion.GetName()
	// モデル名が指定されている場合、上書き
	if overrideModelName != "" {
		modelName = overrideModelName
	}

	// Open the output file
	fout, err := os.Create(path)
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
	modelBName, err := encodeName(modelName, 20)
	if err != nil {
		mlog.W(mi18n.T("モデル名エンコードエラー", map[string]interface{}{"Name": modelName}))
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

	// Write the camera frames
	err = writeCameraFrames(fout, motion)
	if err != nil {
		mlog.E(mi18n.T("カメラフレーム書き込みエラー"))
		return err
	}

	// Write the Light frames
	err = writeLightFrames(fout, motion)
	if err != nil {
		mlog.E(mi18n.T("照明フレーム書き込みエラー"))
		return err
	}

	// Write the Shadow frames
	err = writeShadowFrames(fout, motion)
	if err != nil {
		mlog.E(mi18n.T("照明フレーム書き込みエラー"))
		return err
	}

	// Write the IK frames
	err = writeIkFrames(fout, motion)
	if err != nil {
		mlog.E(mi18n.T("IKフレーム書き込みエラー"))
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

	binary.Write(fout, binary.LittleEndian, uint32(motion.BoneFrames.Len()))
	for _, name := range names {
		fs := motion.BoneFrames.Data[name]

		if fs.Len() > 0 {
			// 各ボーンの最大キーフレを先に出力する
			bf := motion.BoneFrames.Data[name].Get(fs.RegisteredIndexes.Max())
			err := writeBoneFrame(fout, name, bf)
			if err != nil {
				return err
			}
		}
	}

	for _, name := range names {
		fs := motion.BoneFrames.Data[name]
		if fs.Len() > 1 {
			// 普通のキーフレをそのまま出力する
			for _, fno := range fs.RegisteredIndexes.List()[:fs.Len()-1] {
				bf := fs.Get(fno)
				err := writeBoneFrame(fout, name, bf)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func writeBoneFrame(fout *os.File, name string, bf *BoneFrame) error {
	if bf == nil {
		return fmt.Errorf("BoneFrame is nil")
	}

	encodedName, err := encodeName(name, 15)
	if err != nil {
		mlog.W(mi18n.T("ボーン名エンコードエラー", map[string]interface{}{"Name": name}))
		return err
	}

	var posMMD *mmath.MVec3
	if bf.Position != nil {
		posMMD = bf.Position
	} else {
		posMMD = mmath.MVec3Zero
	}
	binary.Write(fout, binary.LittleEndian, encodedName)
	binary.Write(fout, binary.LittleEndian, uint32(bf.Index))
	binary.Write(fout, binary.LittleEndian, float32(posMMD.GetX()))
	binary.Write(fout, binary.LittleEndian, float32(posMMD.GetY()))
	binary.Write(fout, binary.LittleEndian, float32(posMMD.GetZ()))

	var quatMMD *mmath.MQuaternion
	if bf.Rotation != nil {
		quatMMD = bf.Rotation.Normalized()
	} else {
		quatMMD = &mmath.MQuaternionIdent
	}
	binary.Write(fout, binary.LittleEndian, float32(quatMMD.GetX()))
	binary.Write(fout, binary.LittleEndian, float32(quatMMD.GetY()))
	binary.Write(fout, binary.LittleEndian, float32(quatMMD.GetZ()))
	binary.Write(fout, binary.LittleEndian, float32(quatMMD.GetW()))

	var curves []byte
	if bf.Curves == nil {
		curves = InitialBoneCurves
	} else {
		curves = make([]byte, len(bf.Curves.Values))
		for i, x := range bf.Curves.Merge() {
			curves[i] = byte(math.Min(255, math.Max(0, float64(x))))
		}
	}
	binary.Write(fout, binary.LittleEndian, curves)

	return nil
}

func writeMorphFrames(fout *os.File, motion *VmdMotion) error {
	binary.Write(fout, binary.LittleEndian, uint32(motion.MorphFrames.Len()))

	names := motion.MorphFrames.GetNames()

	for _, name := range names {
		fs := motion.MorphFrames.Data[name]
		if fs.RegisteredIndexes.Len() > 0 {
			// 普通のキーフレをそのまま出力する
			for _, fno := range fs.RegisteredIndexes.List() {
				mf := fs.Get(fno)
				err := writeMorphFrame(fout, name, mf)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func writeMorphFrame(fout *os.File, name string, mf *MorphFrame) error {
	if mf == nil {
		return fmt.Errorf("MorphFrame is nil")
	}

	encodedName, err := encodeName(name, 15)
	if err != nil {
		mlog.W(mi18n.T("ボーン名エンコードエラー", map[string]interface{}{"Name": name}))
		return err
	}

	binary.Write(fout, binary.LittleEndian, encodedName)
	binary.Write(fout, binary.LittleEndian, uint32(mf.Index))
	binary.Write(fout, binary.LittleEndian, float32(mf.Ratio))

	return nil
}

func writeCameraFrames(fout *os.File, motion *VmdMotion) error {
	binary.Write(fout, binary.LittleEndian, uint32(motion.CameraFrames.Len()))

	fs := motion.CameraFrames
	if fs.RegisteredIndexes.Len() > 0 {
		// 普通のキーフレをそのまま出力する
		for _, fno := range fs.RegisteredIndexes.List() {
			cf := fs.Get(fno)
			err := writeCameraFrame(fout, cf)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func writeCameraFrame(fout *os.File, cf *CameraFrame) error {
	if cf == nil {
		return fmt.Errorf("CameraFrame is nil")
	}

	binary.Write(fout, binary.LittleEndian, uint32(cf.Index))
	binary.Write(fout, binary.LittleEndian, float32(cf.Distance))

	var posMMD *mmath.MVec3
	if cf.Position != nil {
		posMMD = cf.Position
	} else {
		posMMD = mmath.MVec3Zero
	}

	binary.Write(fout, binary.LittleEndian, float32(posMMD.GetX()))
	binary.Write(fout, binary.LittleEndian, float32(posMMD.GetY()))
	binary.Write(fout, binary.LittleEndian, float32(posMMD.GetZ()))

	var degreeMMD *mmath.MVec3
	if cf.Rotation != nil {
		degreeMMD = cf.Rotation.GetDegrees()
	} else {
		degreeMMD = mmath.MVec3Zero
	}
	binary.Write(fout, binary.LittleEndian, float32(degreeMMD.GetX()))
	binary.Write(fout, binary.LittleEndian, float32(degreeMMD.GetY()))
	binary.Write(fout, binary.LittleEndian, float32(degreeMMD.GetZ()))

	var curves []byte
	if cf.Curves == nil {
		curves = InitialCameraCurves
	} else {
		curves = make([]byte, len(cf.Curves.Values))
		for i, x := range cf.Curves.Merge() {
			curves[i] = byte(math.Min(255, math.Max(0, float64(x))))
		}
	}
	binary.Write(fout, binary.LittleEndian, curves)

	binary.Write(fout, binary.LittleEndian, uint32(cf.ViewOfAngle))
	binary.Write(fout, binary.LittleEndian, byte(mmath.BoolToInt(cf.IsPerspectiveOff)))

	return nil
}

func writeLightFrames(fout *os.File, motion *VmdMotion) error {
	binary.Write(fout, binary.LittleEndian, uint32(motion.LightFrames.Len()))

	fs := motion.LightFrames
	if fs.RegisteredIndexes.Len() > 0 {
		// 普通のキーフレをそのまま出力する
		for _, fno := range fs.RegisteredIndexes.List() {
			cf := fs.Get(fno)
			err := writeLightFrame(fout, cf)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func writeLightFrame(fout *os.File, cf *LightFrame) error {
	if cf == nil {
		return fmt.Errorf("LightFrame is nil")
	}

	binary.Write(fout, binary.LittleEndian, uint32(cf.Index))

	var colorMMD *mmath.MVec3
	if cf.Color != nil {
		colorMMD = cf.Color
	} else {
		colorMMD = mmath.MVec3Zero
	}

	binary.Write(fout, binary.LittleEndian, float32(colorMMD.GetX()))
	binary.Write(fout, binary.LittleEndian, float32(colorMMD.GetY()))
	binary.Write(fout, binary.LittleEndian, float32(colorMMD.GetZ()))

	var posMMD *mmath.MVec3
	if cf.Position != nil {
		posMMD = cf.Position
	} else {
		posMMD = mmath.MVec3Zero
	}

	binary.Write(fout, binary.LittleEndian, float32(posMMD.GetX()))
	binary.Write(fout, binary.LittleEndian, float32(posMMD.GetY()))
	binary.Write(fout, binary.LittleEndian, float32(posMMD.GetZ()))

	return nil
}

func writeShadowFrames(fout *os.File, motion *VmdMotion) error {
	binary.Write(fout, binary.LittleEndian, uint32(motion.ShadowFrames.Len()))

	fs := motion.ShadowFrames
	if fs.RegisteredIndexes.Len() > 0 {
		// 普通のキーフレをそのまま出力する
		for _, fno := range fs.RegisteredIndexes.List() {
			cf := fs.Get(fno)
			err := writeShadowFrame(fout, cf)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func writeShadowFrame(fout *os.File, cf *ShadowFrame) error {
	if cf == nil {
		return fmt.Errorf("ShadowFrame is nil")
	}

	binary.Write(fout, binary.LittleEndian, uint32(cf.Index))

	binary.Write(fout, binary.LittleEndian, float32(cf.ShadowMode))
	binary.Write(fout, binary.LittleEndian, float32(cf.Distance))

	return nil
}

func writeIkFrames(fout *os.File, motion *VmdMotion) error {
	binary.Write(fout, binary.LittleEndian, uint32(motion.IkFrames.Len()))

	fs := motion.IkFrames
	if fs.RegisteredIndexes.Len() > 0 {
		// 普通のキーフレをそのまま出力する
		for _, fno := range fs.RegisteredIndexes.List() {
			cf := fs.Get(fno)
			err := writeIkFrame(fout, cf)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func writeIkFrame(fout *os.File, cf *IkFrame) error {
	if cf == nil {
		return fmt.Errorf("IkFrame is nil")
	}

	binary.Write(fout, binary.LittleEndian, uint32(cf.Index))
	binary.Write(fout, binary.LittleEndian, byte(mmath.BoolToInt(cf.Visible)))
	binary.Write(fout, binary.LittleEndian, uint32(len(cf.IkList)))

	fs := cf.IkList
	if len(fs) > 0 {
		// 普通のキーフレをそのまま出力する
		for _, ik := range fs {
			encodedName, err := encodeName(ik.BoneName, 20)
			if err != nil {
				mlog.W(mi18n.T("ボーン名エンコードエラー", map[string]interface{}{"Name": ik.BoneName}))
				return err
			}

			binary.Write(fout, binary.LittleEndian, encodedName)
			binary.Write(fout, binary.LittleEndian, byte(mmath.BoolToInt(ik.Enabled)))
		}
	}

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
