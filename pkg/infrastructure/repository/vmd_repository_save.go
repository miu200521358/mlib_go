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

	"github.com/miu200521358/mlib_go/pkg/config/mi18n"
	"github.com/miu200521358/mlib_go/pkg/config/mlog"
	"github.com/miu200521358/mlib_go/pkg/config/mproc"
	"github.com/miu200521358/mlib_go/pkg/domain/core"
	"github.com/miu200521358/mlib_go/pkg/domain/mmath"
	"github.com/miu200521358/mlib_go/pkg/domain/vmd"
	"github.com/miu200521358/mlib_go/pkg/infrastructure/mfile"
)

func (rep *VmdRepository) Save(overridePath string, data core.IHashModel, includeSystem bool) error {
	mproc.SetMaxProcess(true)
	defer mproc.SetMaxProcess(false)

	mlog.IL("%s", mi18n.T("保存開始", map[string]interface{}{"Type": "Vmd", "Path": overridePath}))
	defer mlog.I("%s", mi18n.T("保存終了", map[string]interface{}{"Type": "Vmd"}))

	motion := data.(*vmd.VmdMotion)

	path := motion.Path()
	// 保存可能なパスである場合、上書き
	if mfile.CanSave(overridePath) {
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
	modelBName, err := rep.encodeName(motion.Name(), 20)
	if err != nil {
		mlog.W("%s", mi18n.T("モデル名エンコードエラー", map[string]interface{}{"Name": motion.Name()}))
		modelBName = []byte("Vmd Model")
	}

	// Write the model name
	err = binary.Write(fout, binary.LittleEndian, modelBName)
	if err != nil {
		return err
	}

	// Write the bone frames
	err = rep.saveBoneFrames(fout, motion)
	if err != nil {
		mlog.E("%s", err, mi18n.T("ボーンフレーム書き込みエラー"))
		return err
	}

	// Write the morph frames
	err = rep.saveMorphFrames(fout, motion)
	if err != nil {
		mlog.E("%s", err, mi18n.T("モーフフレーム書き込みエラー"))
		return err
	}

	// Write the camera frames
	err = rep.saveCameraFrames(fout, motion)
	if err != nil {
		mlog.E("%s", err, mi18n.T("カメラフレーム書き込みエラー"))
		return err
	}

	// Write the Light frames
	err = rep.saveLightFrames(fout, motion)
	if err != nil {
		mlog.E("%s", err, mi18n.T("照明フレーム書き込みエラー"))
		return err
	}

	// Write the Shadow frames
	err = rep.saveShadowFrames(fout, motion)
	if err != nil {
		mlog.E("%s", err, mi18n.T("照明フレーム書き込みエラー"))
		return err
	}

	// Write the IK frames
	err = rep.saveIkFrames(fout, motion)
	if err != nil {
		mlog.E("%s", err, mi18n.T("IKフレーム書き込みエラー"))
		return err
	}

	// foutを書き込んで終了する
	err = fout.Close()
	if err != nil {
		mlog.E("%s", err, mi18n.T("ファイルクローズエラー", map[string]interface{}{"Path": motion.Path()}))
		return err
	}

	return nil
}

func (rep *VmdRepository) saveBoneFrames(fout *os.File, motion *vmd.VmdMotion) error {
	names := make([]string, 0, motion.BoneFrames.Length())
	count := 0
	motion.BoneFrames.ForEach(func(name string, boneNameFrames *vmd.BoneNameFrames) {
		if !boneNameFrames.ContainsActive() {
			return
		}
		names = append(names, boneNameFrames.Name)
		count += boneNameFrames.Length()
	})

	rep.writeNumber(fout, binaryType_unsignedInt, float64(count), 0.0, true)
	n := 0
	for _, name := range names {
		boneFrames := motion.BoneFrames.Get(name)

		if boneFrames.Length() > 0 {
			// 各ボーンの最大キーフレを先に出力する
			bf := motion.BoneFrames.Get(name).Get(boneFrames.Indexes.Max())
			err := rep.saveBoneFrame(fout, name, bf)
			if err != nil {
				return err
			}
		}

		if n%10000 == 0 && n > 0 {
			mlog.I("%s", mi18n.T("保存途中", map[string]interface{}{"Type": mi18n.T("ボーン"), "Index": n, "Total": count}))
		}
		n++
	}

	for _, name := range names {
		fs := motion.BoneFrames.Get(name)
		maxFno := fs.Indexes.Max()
		if fs.Length() > 1 {
			// 普通のキーフレをそのまま出力する
			fs.Indexes.ForEach(func(fno float32) bool {
				if fno < maxFno {
					bf := fs.Get(fno)
					err := rep.saveBoneFrame(fout, name, bf)
					if err != nil {
						return false
					}
				}

				return true
			})
		}
	}

	return nil
}

func (rep *VmdRepository) saveBoneFrame(fout *os.File, name string, bf *vmd.BoneFrame) error {
	if bf == nil {
		return fmt.Errorf("BoneFrame is nil")
	}

	encodedName, err := rep.encodeName(name, 15)
	if err != nil {
		mlog.W("%s", mi18n.T("ボーン名エンコードエラー", map[string]interface{}{"Name": name}))
		return err
	}

	var posMMD *mmath.MVec3
	if bf.Position != nil {
		posMMD = bf.Position
	} else {
		posMMD = mmath.MVec3Zero
	}
	binary.Write(fout, binary.LittleEndian, encodedName)
	rep.writeNumber(fout, binaryType_unsignedInt, float64(bf.Index()), 0.0, true)
	rep.writeNumber(fout, binaryType_float, posMMD.X, 0.0, false)
	rep.writeNumber(fout, binaryType_float, posMMD.Y, 0.0, false)
	rep.writeNumber(fout, binaryType_float, posMMD.Z, 0.0, false)

	var quatMMD *mmath.MQuaternion
	if bf.Rotation != nil {
		quatMMD = bf.Rotation.Normalized()
	} else {
		quatMMD = mmath.MQuaternionIdent
	}
	rep.writeNumber(fout, binaryType_float, quatMMD.Vec3().X, 0.0, false)
	rep.writeNumber(fout, binaryType_float, quatMMD.Vec3().Y, 0.0, false)
	rep.writeNumber(fout, binaryType_float, quatMMD.Vec3().Z, 0.0, false)
	rep.writeNumber(fout, binaryType_float, quatMMD.W, 0.0, false)

	var curves []byte
	if bf.Curves == nil {
		curves = vmd.InitialBoneCurves
	} else {
		curves = make([]byte, len(vmd.InitialBoneCurves))
		for i, x := range bf.Curves.Merge() {
			curves[i] = byte(math.Min(255, math.Max(0, float64(x))))
		}
	}
	binary.Write(fout, binary.LittleEndian, curves)

	return nil
}

func (rep *VmdRepository) saveMorphFrames(fout *os.File, motion *vmd.VmdMotion) error {
	names := make([]string, 0, motion.MorphFrames.Length())
	count := 0
	motion.MorphFrames.ForEach(func(name string, morphNameFrames *vmd.MorphNameFrames) {
		if !morphNameFrames.ContainsActive() {
			return
		}
		names = append(names, morphNameFrames.Name)
		count += morphNameFrames.Length()
	})

	rep.writeNumber(fout, binaryType_unsignedInt, float64(count), 0.0, true)
	n := 0
	for _, name := range names {
		morphFrames := motion.MorphFrames.Get(name)
		if morphFrames.Indexes.Len() > 0 {
			// 普通のキーフレをそのまま出力する
			morphFrames.Indexes.ForEach(func(fno float32) bool {
				mf := morphFrames.Get(fno)
				err := rep.saveMorphFrame(fout, name, mf)
				if err != nil {
					return false
				}

				return true
			})
		}

		if n%10000 == 0 && n > 0 {
			mlog.I("%s", mi18n.T("保存途中", map[string]interface{}{"Type": mi18n.T("モーフ"), "Index": n, "Total": count}))
		}
		n++
	}

	return nil
}

func (rep *VmdRepository) saveMorphFrame(fout *os.File, name string, mf *vmd.MorphFrame) error {
	if mf == nil {
		return fmt.Errorf("MorphFrame is nil")
	}

	encodedName, err := rep.encodeName(name, 15)
	if err != nil {
		mlog.W("%s", mi18n.T("モーフ名エンコードエラー", map[string]interface{}{"Name": name}))
		return err
	}

	binary.Write(fout, binary.LittleEndian, encodedName)
	rep.writeNumber(fout, binaryType_unsignedInt, float64(mf.Index()), 0.0, true)
	rep.writeNumber(fout, binaryType_float, mf.Ratio, 0.0, false)

	return nil
}

func (rep *VmdRepository) saveCameraFrames(fout *os.File, motion *vmd.VmdMotion) error {
	rep.writeNumber(fout, binaryType_unsignedInt, float64(motion.CameraFrames.Length()), 0.0, true)

	cameraFrames := motion.CameraFrames
	if cameraFrames.Indexes.Len() > 0 {
		// 普通のキーフレをそのまま出力する
		cameraFrames.Indexes.ForEach(func(fno float32) bool {
			cf := cameraFrames.Get(fno)
			err := rep.saveCameraFrame(fout, cf)
			if err != nil {
				return false
			}
			return true
		})
	}

	return nil
}

func (rep *VmdRepository) saveCameraFrame(fout *os.File, cf *vmd.CameraFrame) error {
	if cf == nil {
		return fmt.Errorf("CameraFrame is nil")
	}

	rep.writeNumber(fout, binaryType_unsignedInt, float64(cf.Index()), 0.0, true)
	rep.writeNumber(fout, binaryType_float, cf.Distance, 0.0, false)

	var posMMD *mmath.MVec3
	if cf.Position != nil {
		posMMD = cf.Position
	} else {
		posMMD = mmath.MVec3Zero
	}

	rep.writeNumber(fout, binaryType_float, posMMD.X, 0.0, false)
	rep.writeNumber(fout, binaryType_float, posMMD.Y, 0.0, false)
	rep.writeNumber(fout, binaryType_float, posMMD.Z, 0.0, false)

	degreeMMD := cf.Degrees.MMD()
	rep.writeNumber(fout, binaryType_float, degreeMMD.X, 0.0, false)
	rep.writeNumber(fout, binaryType_float, degreeMMD.Y, 0.0, false)
	rep.writeNumber(fout, binaryType_float, degreeMMD.Z, 0.0, false)

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

	rep.writeNumber(fout, binaryType_unsignedInt, float64(cf.ViewOfAngle), 0.0, true)
	rep.writeBool(fout, cf.IsPerspectiveOff)

	return nil
}

func (rep *VmdRepository) saveLightFrames(fout *os.File, motion *vmd.VmdMotion) error {
	rep.writeNumber(fout, binaryType_unsignedInt, float64(motion.LightFrames.Length()), 0.0, true)

	lightFrames := motion.LightFrames
	if lightFrames.Indexes.Len() > 0 {
		// 普通のキーフレをそのまま出力する
		lightFrames.Indexes.ForEach(func(fno float32) bool {
			lf := lightFrames.Get(fno)
			err := rep.saveLightFrame(fout, lf)
			if err != nil {
				return false
			}
			return true
		})
	}

	return nil
}

func (rep *VmdRepository) saveLightFrame(fout *os.File, lf *vmd.LightFrame) error {
	if lf == nil {
		return fmt.Errorf("LightFrame is nil")
	}

	rep.writeNumber(fout, binaryType_unsignedInt, float64(lf.Index()), 0.0, true)

	var colorMMD *mmath.MVec3
	if lf.Color != nil {
		colorMMD = lf.Color
	} else {
		colorMMD = mmath.MVec3Zero
	}

	rep.writeNumber(fout, binaryType_float, colorMMD.X, 0.0, false)
	rep.writeNumber(fout, binaryType_float, colorMMD.Y, 0.0, false)
	rep.writeNumber(fout, binaryType_float, colorMMD.Z, 0.0, false)

	var posMMD *mmath.MVec3
	if lf.Position != nil {
		posMMD = lf.Position
	} else {
		posMMD = mmath.MVec3Zero
	}

	rep.writeNumber(fout, binaryType_float, posMMD.X, 0.0, false)
	rep.writeNumber(fout, binaryType_float, posMMD.Y, 0.0, false)
	rep.writeNumber(fout, binaryType_float, posMMD.Z, 0.0, false)

	return nil
}

func (rep *VmdRepository) saveShadowFrames(fout *os.File, motion *vmd.VmdMotion) error {
	rep.writeNumber(fout, binaryType_unsignedInt, float64(motion.ShadowFrames.Length()), 0.0, true)

	shadowFrames := motion.ShadowFrames
	if shadowFrames.Indexes.Len() > 0 {
		// 普通のキーフレをそのまま出力する
		shadowFrames.Indexes.ForEach(func(fno float32) bool {
			sf := shadowFrames.Get(fno)
			err := rep.sveShadowFrame(fout, sf)
			if err != nil {
				return false
			}
			return true
		})
	}

	return nil
}

func (rep *VmdRepository) sveShadowFrame(fout *os.File, sf *vmd.ShadowFrame) error {
	if sf == nil {
		return fmt.Errorf("ShadowFrame is nil")
	}

	rep.writeNumber(fout, binaryType_unsignedInt, float64(sf.Index()), 0.0, true)

	rep.writeNumber(fout, binaryType_float, float64(sf.ShadowMode), 0.0, false)
	rep.writeNumber(fout, binaryType_float, sf.Distance, 0.0, false)

	return nil
}

func (rep *VmdRepository) saveIkFrames(fout *os.File, motion *vmd.VmdMotion) error {
	rep.writeNumber(fout, binaryType_unsignedInt, float64(motion.IkFrames.Length()), 0.0, true)

	ikFrames := motion.IkFrames
	if ikFrames.Indexes.Len() > 0 {
		// 普通のキーフレをそのまま出力する
		ikFrames.Indexes.ForEach(func(fno float32) bool {
			ikf := ikFrames.Get(fno)
			err := rep.saveIkFrame(fout, ikf)
			if err != nil {
				return false
			}
			return true
		})
	}

	return nil
}

func (rep *VmdRepository) saveIkFrame(fout *os.File, ikf *vmd.IkFrame) error {
	if ikf == nil {
		return fmt.Errorf("IkFrame is nil")
	}

	rep.writeNumber(fout, binaryType_unsignedInt, float64(ikf.Index()), 0.0, true)
	rep.writeBool(fout, ikf.Visible)
	rep.writeNumber(fout, binaryType_unsignedInt, float64(len(ikf.IkList)), 0.0, true)

	fs := ikf.IkList
	if len(fs) > 0 {
		// 普通のキーフレをそのまま出力する
		for _, ik := range fs {
			encodedName, err := rep.encodeName(ik.BoneName, 20)
			if err != nil {
				mlog.W("%s", mi18n.T("ボーン名エンコードエラー", map[string]interface{}{"Name": ik.BoneName}))
				return err
			}

			binary.Write(fout, binary.LittleEndian, encodedName)
			rep.writeBool(fout, ik.Enabled)
		}
	}

	return nil
}

func (rep *VmdRepository) encodeName(name string, limit int) ([]byte, error) {
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
