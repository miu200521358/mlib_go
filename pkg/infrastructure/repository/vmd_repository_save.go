package repository

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"os"

	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"

	"github.com/miu200521358/mlib_go/pkg/domain/core"
	"github.com/miu200521358/mlib_go/pkg/domain/mmath"
	"github.com/miu200521358/mlib_go/pkg/domain/vmd"
	"github.com/miu200521358/mlib_go/pkg/mutils"
	"github.com/miu200521358/mlib_go/pkg/mutils/mi18n"
	"github.com/miu200521358/mlib_go/pkg/mutils/mlog"
)

func (r *VmdRepository) Save(overridePath string, data core.IHashModel, includeSystem bool) error {
	motion := data.(*vmd.VmdMotion)

	path := motion.Path()
	// 保存可能なパスである場合、上書き
	if mutils.CanSave(overridePath) {
		path = overridePath
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
	modelBName, err := r.encodeName(motion.Name(), 20)
	if err != nil {
		mlog.W(mi18n.T("モデル名エンコードエラー", map[string]interface{}{"Name": motion.Name()}))
		modelBName = []byte("Vmd Model")
	}

	// Write the model name
	err = binary.Write(fout, binary.LittleEndian, modelBName)
	if err != nil {
		return err
	}

	// Write the bone frames
	err = r.saveBoneFrames(fout, motion)
	if err != nil {
		mlog.E(mi18n.T("ボーンフレーム書き込みエラー"))
		return err
	}

	// Write the morph frames
	err = r.saveMorphFrames(fout, motion)
	if err != nil {
		mlog.E(mi18n.T("モーフフレーム書き込みエラー"))
		return err
	}

	// Write the camera frames
	err = r.saveCameraFrames(fout, motion)
	if err != nil {
		mlog.E(mi18n.T("カメラフレーム書き込みエラー"))
		return err
	}

	// Write the Light frames
	err = r.saveLightFrames(fout, motion)
	if err != nil {
		mlog.E(mi18n.T("照明フレーム書き込みエラー"))
		return err
	}

	// Write the Shadow frames
	err = r.saveShadowFrames(fout, motion)
	if err != nil {
		mlog.E(mi18n.T("照明フレーム書き込みエラー"))
		return err
	}

	// Write the IK frames
	err = r.saveIkFrames(fout, motion)
	if err != nil {
		mlog.E(mi18n.T("IKフレーム書き込みエラー"))
		return err
	}

	// foutを書き込んで終了する
	err = fout.Close()
	if err != nil {
		mlog.E(mi18n.T("ファイルクローズエラー", map[string]interface{}{"Path": motion.Path()}))
		return err
	}

	return nil
}

func (r *VmdRepository) saveBoneFrames(fout *os.File, motion *vmd.VmdMotion) error {
	names := motion.BoneFrames.GetNames()

	r.writeNumber(fout, binaryType_unsignedInt, float64(motion.BoneFrames.Len()), 0.0, true)
	for _, name := range names {
		fs := motion.BoneFrames.Data[name]

		if fs.Len() > 0 {
			// 各ボーンの最大キーフレを先に出力する
			bf := motion.BoneFrames.Data[name].Get(fs.RegisteredIndexes.Max())
			err := r.saveBoneFrame(fout, name, bf)
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
				err := r.saveBoneFrame(fout, name, bf)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func (r *VmdRepository) saveBoneFrame(fout *os.File, name string, bf *vmd.BoneFrame) error {
	if bf == nil {
		return fmt.Errorf("BoneFrame is nil")
	}

	encodedName, err := r.encodeName(name, 15)
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
	r.writeNumber(fout, binaryType_unsignedInt, float64(bf.Index()), 0.0, true)
	r.writeNumber(fout, binaryType_float, posMMD.X, 0.0, false)
	r.writeNumber(fout, binaryType_float, posMMD.Y, 0.0, false)
	r.writeNumber(fout, binaryType_float, posMMD.Z, 0.0, false)

	var quatMMD *mmath.MQuaternion
	if bf.Rotation != nil {
		quatMMD = bf.Rotation.Normalized()
	} else {
		quatMMD = &mmath.MQuaternionIdent
	}
	r.writeNumber(fout, binaryType_float, quatMMD.Vec3().X, 0.0, false)
	r.writeNumber(fout, binaryType_float, quatMMD.Vec3().Y, 0.0, false)
	r.writeNumber(fout, binaryType_float, quatMMD.Vec3().Z, 0.0, false)
	r.writeNumber(fout, binaryType_float, quatMMD.W, 0.0, false)

	var curves []byte
	if bf.Curves == nil {
		curves = vmd.InitialBoneCurves
	} else {
		curves = make([]byte, len(bf.Curves.Values))
		for i, x := range bf.Curves.Merge() {
			curves[i] = byte(math.Min(255, math.Max(0, float64(x))))
		}
	}
	binary.Write(fout, binary.LittleEndian, curves)

	return nil
}

func (r *VmdRepository) saveMorphFrames(fout *os.File, motion *vmd.VmdMotion) error {
	r.writeNumber(fout, binaryType_unsignedInt, float64(motion.MorphFrames.Len()), 0.0, true)

	names := motion.MorphFrames.GetNames()

	for _, name := range names {
		fs := motion.MorphFrames.Data[name]
		if fs.RegisteredIndexes.Len() > 0 {
			// 普通のキーフレをそのまま出力する
			for _, fno := range fs.RegisteredIndexes.List() {
				mf := fs.Get(fno)
				err := r.saveMorphFrame(fout, name, mf)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func (r *VmdRepository) saveMorphFrame(fout *os.File, name string, mf *vmd.MorphFrame) error {
	if mf == nil {
		return fmt.Errorf("MorphFrame is nil")
	}

	encodedName, err := r.encodeName(name, 15)
	if err != nil {
		mlog.W(mi18n.T("ボーン名エンコードエラー", map[string]interface{}{"Name": name}))
		return err
	}

	binary.Write(fout, binary.LittleEndian, encodedName)
	r.writeNumber(fout, binaryType_unsignedInt, float64(mf.Index()), 0.0, true)
	r.writeNumber(fout, binaryType_float, mf.Ratio, 0.0, false)

	return nil
}

func (r *VmdRepository) saveCameraFrames(fout *os.File, motion *vmd.VmdMotion) error {
	r.writeNumber(fout, binaryType_unsignedInt, float64(motion.CameraFrames.Len()), 0.0, true)

	fs := motion.CameraFrames
	if fs.RegisteredIndexes.Len() > 0 {
		// 普通のキーフレをそのまま出力する
		for _, fno := range fs.RegisteredIndexes.List() {
			cf := fs.Get(fno)
			err := r.saveCameraFrame(fout, cf)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (r *VmdRepository) saveCameraFrame(fout *os.File, cf *vmd.CameraFrame) error {
	if cf == nil {
		return fmt.Errorf("CameraFrame is nil")
	}

	r.writeNumber(fout, binaryType_unsignedInt, float64(cf.Index()), 0.0, true)
	r.writeNumber(fout, binaryType_float, cf.Distance, 0.0, false)

	var posMMD *mmath.MVec3
	if cf.Position != nil {
		posMMD = cf.Position
	} else {
		posMMD = mmath.MVec3Zero
	}

	r.writeNumber(fout, binaryType_float, posMMD.X, 0.0, false)
	r.writeNumber(fout, binaryType_float, posMMD.Y, 0.0, false)
	r.writeNumber(fout, binaryType_float, posMMD.Z, 0.0, false)

	var degreeMMD *mmath.MVec3
	if cf.Rotation != nil {
		degreeMMD = cf.Rotation.GetDegrees()
	} else {
		degreeMMD = mmath.MVec3Zero
	}
	r.writeNumber(fout, binaryType_float, degreeMMD.X, 0.0, false)
	r.writeNumber(fout, binaryType_float, degreeMMD.Y, 0.0, false)
	r.writeNumber(fout, binaryType_float, degreeMMD.Z, 0.0, false)

	var curves []byte
	if cf.Curves == nil {
		curves = vmd.InitialCameraCurves
	} else {
		curves = make([]byte, len(cf.Curves.Values))
		for i, x := range cf.Curves.Merge() {
			curves[i] = byte(math.Min(255, math.Max(0, float64(x))))
		}
	}
	binary.Write(fout, binary.LittleEndian, curves)

	r.writeNumber(fout, binaryType_unsignedInt, float64(cf.ViewOfAngle), 0.0, true)
	r.writeBool(fout, cf.IsPerspectiveOff)

	return nil
}

func (r *VmdRepository) saveLightFrames(fout *os.File, motion *vmd.VmdMotion) error {
	r.writeNumber(fout, binaryType_unsignedInt, float64(motion.LightFrames.Len()), 0.0, true)

	fs := motion.LightFrames
	if fs.RegisteredIndexes.Len() > 0 {
		// 普通のキーフレをそのまま出力する
		for _, fno := range fs.RegisteredIndexes.List() {
			cf := fs.Get(fno)
			err := r.saveLightFrame(fout, cf)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (r *VmdRepository) saveLightFrame(fout *os.File, cf *vmd.LightFrame) error {
	if cf == nil {
		return fmt.Errorf("LightFrame is nil")
	}

	r.writeNumber(fout, binaryType_unsignedInt, float64(cf.Index()), 0.0, true)

	var colorMMD *mmath.MVec3
	if cf.Color != nil {
		colorMMD = cf.Color
	} else {
		colorMMD = mmath.MVec3Zero
	}

	r.writeNumber(fout, binaryType_float, colorMMD.X, 0.0, false)
	r.writeNumber(fout, binaryType_float, colorMMD.Y, 0.0, false)
	r.writeNumber(fout, binaryType_float, colorMMD.Z, 0.0, false)

	var posMMD *mmath.MVec3
	if cf.Position != nil {
		posMMD = cf.Position
	} else {
		posMMD = mmath.MVec3Zero
	}

	r.writeNumber(fout, binaryType_float, posMMD.X, 0.0, false)
	r.writeNumber(fout, binaryType_float, posMMD.Y, 0.0, false)
	r.writeNumber(fout, binaryType_float, posMMD.Z, 0.0, false)

	return nil
}

func (r *VmdRepository) saveShadowFrames(fout *os.File, motion *vmd.VmdMotion) error {
	r.writeNumber(fout, binaryType_unsignedInt, float64(motion.ShadowFrames.Len()), 0.0, true)

	fs := motion.ShadowFrames
	if fs.RegisteredIndexes.Len() > 0 {
		// 普通のキーフレをそのまま出力する
		for _, fno := range fs.RegisteredIndexes.List() {
			cf := fs.Get(fno)
			err := r.sveShadowFrame(fout, cf)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (r *VmdRepository) sveShadowFrame(fout *os.File, cf *vmd.ShadowFrame) error {
	if cf == nil {
		return fmt.Errorf("ShadowFrame is nil")
	}

	r.writeNumber(fout, binaryType_unsignedInt, float64(cf.Index()), 0.0, true)

	r.writeNumber(fout, binaryType_float, float64(cf.ShadowMode), 0.0, false)
	r.writeNumber(fout, binaryType_float, cf.Distance, 0.0, false)

	return nil
}

func (r *VmdRepository) saveIkFrames(fout *os.File, motion *vmd.VmdMotion) error {
	r.writeNumber(fout, binaryType_unsignedInt, float64(motion.IkFrames.Len()), 0.0, true)

	fs := motion.IkFrames
	if fs.RegisteredIndexes.Len() > 0 {
		// 普通のキーフレをそのまま出力する
		for _, fno := range fs.RegisteredIndexes.List() {
			cf := fs.Get(fno)
			err := r.saveIkFrame(fout, cf)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (r *VmdRepository) saveIkFrame(fout *os.File, cf *vmd.IkFrame) error {
	if cf == nil {
		return fmt.Errorf("IkFrame is nil")
	}

	r.writeNumber(fout, binaryType_unsignedInt, float64(cf.Index()), 0.0, true)
	r.writeBool(fout, cf.Visible)
	r.writeNumber(fout, binaryType_unsignedInt, float64(len(cf.IkList)), 0.0, true)

	fs := cf.IkList
	if len(fs) > 0 {
		// 普通のキーフレをそのまま出力する
		for _, ik := range fs {
			encodedName, err := r.encodeName(ik.BoneName, 20)
			if err != nil {
				mlog.W(mi18n.T("ボーン名エンコードエラー", map[string]interface{}{"Name": ik.BoneName}))
				return err
			}

			binary.Write(fout, binary.LittleEndian, encodedName)
			r.writeBool(fout, ik.Enabled)
		}
	}

	return nil
}

func (r *VmdRepository) encodeName(name string, limit int) ([]byte, error) {
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
