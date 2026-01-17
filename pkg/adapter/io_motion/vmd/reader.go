// 指示: miu200521358
package vmd

import (
	"io"

	"github.com/miu200521358/mlib_go/pkg/adapter/io_common"
	"github.com/miu200521358/mlib_go/pkg/domain/mmath"
	"github.com/miu200521358/mlib_go/pkg/domain/motion"
)

// vmdReader はVMD読み込み処理を表す。
type vmdReader struct {
	reader *io_common.BinaryReader
}

// newVmdReader はvmdReaderを生成する。
func newVmdReader(r io.Reader) *vmdReader {
	return &vmdReader{reader: io_common.NewBinaryReader(r)}
}

// Read はVMDを読み込む。
func (v *vmdReader) Read(motionData *motion.VmdMotion) error {
	if motionData == nil {
		return io_common.NewIoParseFailed("VMDモーションがnilです", nil)
	}
	if err := v.readHeader(motionData); err != nil {
		return err
	}
	if err := v.readBoneFrames(motionData); err != nil {
		return err
	}
	if err := v.readMorphFrames(motionData); err != nil {
		return err
	}
	if err := v.readCameraFrames(motionData); err != nil {
		return nil
	}
	if err := v.readLightFrames(motionData); err != nil {
		return nil
	}
	if err := v.readShadowFrames(motionData); err != nil {
		return nil
	}
	if err := v.readIkFrames(motionData); err != nil {
		return nil
	}
	return nil
}

// readHeader はヘッダを読み込む。
func (v *vmdReader) readHeader(motionData *motion.VmdMotion) error {
	signatureRaw, err := v.reader.ReadBytes(30)
	if err != nil {
		return io_common.NewIoParseFailed("VMD署名の読み込みに失敗しました", err)
	}
	signature, err := io_common.DecodeShiftJISFixed(signatureRaw)
	if err != nil {
		return io_common.NewIoParseFailed("VMD署名のデコードに失敗しました", err)
	}
	modelNameRaw, err := v.reader.ReadBytes(20)
	if err != nil {
		return io_common.NewIoParseFailed("VMDモデル名の読み込みに失敗しました", err)
	}
	modelName, err := io_common.DecodeShiftJISFixed(modelNameRaw)
	if err != nil {
		return io_common.NewIoParseFailed("VMDモデル名のデコードに失敗しました", err)
	}
	motionData.Signature = signature
	motionData.SetName(modelName)
	return nil
}

// readBoneFrames はボーンフレームを読み込む。
func (v *vmdReader) readBoneFrames(motionData *motion.VmdMotion) error {
	count, err := v.reader.ReadUint32()
	if err != nil {
		return io_common.NewIoParseFailed("VMDボーンフレーム数の読み込みに失敗しました", err)
	}
	for i := 0; i < int(count); i++ {
		nameRaw, err := v.reader.ReadBytes(15)
		if err != nil {
			return io_common.NewIoParseFailed("VMDボーン名の読み込みに失敗しました", err)
		}
		name, err := io_common.DecodeShiftJISFixed(nameRaw)
		if err != nil {
			return io_common.NewIoParseFailed("VMDボーン名のデコードに失敗しました", err)
		}
		frameIndex, err := v.reader.ReadUint32()
		if err != nil {
			return io_common.NewIoParseFailed("VMDフレーム番号の読み込みに失敗しました", err)
		}
		pos, err := v.reader.ReadVec3()
		if err != nil {
			return io_common.NewIoParseFailed("VMDボーン位置の読み込みに失敗しました", err)
		}
		rotRaw, err := v.reader.ReadVec4()
		if err != nil {
			return io_common.NewIoParseFailed("VMDボーン回転の読み込みに失敗しました", err)
		}
		curveRaw, err := v.reader.ReadBytes(64)
		if err != nil {
			return io_common.NewIoParseFailed("VMDボーン補間曲線の読み込みに失敗しました", err)
		}

		frame := motion.NewBoneFrame(motion.Frame(float32(frameIndex)))
		frame.Read = true
		frame.Position = &pos
		q := mmath.NewQuaternionByValues(rotRaw.X, rotRaw.Y, rotRaw.Z, rotRaw.W)
		frame.Rotation = &q
		frame.Curves = motion.NewBoneCurvesByValues(curveRaw)
		motionData.AppendBoneFrame(name, frame)
	}
	return nil
}

// readMorphFrames はモーフフレームを読み込む。
func (v *vmdReader) readMorphFrames(motionData *motion.VmdMotion) error {
	count, err := v.reader.ReadUint32()
	if err != nil {
		return io_common.NewIoParseFailed("VMDモーフフレーム数の読み込みに失敗しました", err)
	}
	for i := 0; i < int(count); i++ {
		nameRaw, err := v.reader.ReadBytes(15)
		if err != nil {
			return io_common.NewIoParseFailed("VMDモーフ名の読み込みに失敗しました", err)
		}
		name, err := io_common.DecodeShiftJISFixed(nameRaw)
		if err != nil {
			return io_common.NewIoParseFailed("VMDモーフ名のデコードに失敗しました", err)
		}
		frameIndex, err := v.reader.ReadUint32()
		if err != nil {
			return io_common.NewIoParseFailed("VMDモーフフレーム番号の読み込みに失敗しました", err)
		}
		ratio, err := v.reader.ReadFloat32()
		if err != nil {
			return io_common.NewIoParseFailed("VMDモーフ比率の読み込みに失敗しました", err)
		}
		frame := motion.NewMorphFrame(motion.Frame(float32(frameIndex)))
		frame.Read = true
		frame.Ratio = ratio
		motionData.AppendMorphFrame(name, frame)
	}
	return nil
}

// readCameraFrames はカメラフレームを読み込む。
func (v *vmdReader) readCameraFrames(motionData *motion.VmdMotion) error {
	count, err := v.reader.ReadUint32()
	if err != nil {
		return io_common.NewIoParseFailed("VMDカメラフレーム数の読み込みに失敗しました", err)
	}
	for i := 0; i < int(count); i++ {
		frameIndex, err := v.reader.ReadUint32()
		if err != nil {
			return io_common.NewIoParseFailed("VMDカメラフレーム番号の読み込みに失敗しました", err)
		}
		frame := motion.NewCameraFrame(motion.Frame(float32(frameIndex)))
		frame.Read = true
		frame.Distance, err = v.reader.ReadFloat32()
		if err != nil {
			return io_common.NewIoParseFailed("VMDカメラ距離の読み込みに失敗しました", err)
		}
		pos, err := v.reader.ReadVec3()
		if err != nil {
			return io_common.NewIoParseFailed("VMDカメラ位置の読み込みに失敗しました", err)
		}
		frame.Position = &pos
		deg, err := v.reader.ReadVec3()
		if err != nil {
			return io_common.NewIoParseFailed("VMDカメラ回転の読み込みに失敗しました", err)
		}
		frame.Degrees = &deg
		curveRaw, err := v.reader.ReadBytes(24)
		if err != nil {
			return io_common.NewIoParseFailed("VMDカメラ補間曲線の読み込みに失敗しました", err)
		}
		frame.Curves = motion.NewCameraCurvesByValues(curveRaw)
		viewOfAngle, err := v.reader.ReadUint32()
		if err != nil {
			return io_common.NewIoParseFailed("VMDカメラ視野角の読み込みに失敗しました", err)
		}
		frame.ViewOfAngle = int(viewOfAngle)
		perspective, err := v.reader.ReadUint8()
		if err != nil {
			return io_common.NewIoParseFailed("VMDカメラパース設定の読み込みに失敗しました", err)
		}
		frame.IsPerspectiveOff = perspective == 1
		motionData.AppendCameraFrame(frame)
	}
	return nil
}

// readLightFrames はライトフレームを読み込む。
func (v *vmdReader) readLightFrames(motionData *motion.VmdMotion) error {
	count, err := v.reader.ReadUint32()
	if err != nil {
		return io_common.NewIoParseFailed("VMDライトフレーム数の読み込みに失敗しました", err)
	}
	for i := 0; i < int(count); i++ {
		frameIndex, err := v.reader.ReadUint32()
		if err != nil {
			return io_common.NewIoParseFailed("VMDライトフレーム番号の読み込みに失敗しました", err)
		}
		frame := motion.NewLightFrame(motion.Frame(float32(frameIndex)))
		frame.Read = true
		color, err := v.reader.ReadVec3()
		if err != nil {
			return io_common.NewIoParseFailed("VMDライト色の読み込みに失敗しました", err)
		}
		pos, err := v.reader.ReadVec3()
		if err != nil {
			return io_common.NewIoParseFailed("VMDライト位置の読み込みに失敗しました", err)
		}
		frame.Color = color
		frame.Position = pos
		motionData.AppendLightFrame(frame)
	}
	return nil
}

// readShadowFrames はシャドウフレームを読み込む。
func (v *vmdReader) readShadowFrames(motionData *motion.VmdMotion) error {
	count, err := v.reader.ReadUint32()
	if err != nil {
		return io_common.NewIoParseFailed("VMDシャドウフレーム数の読み込みに失敗しました", err)
	}
	for i := 0; i < int(count); i++ {
		frameIndex, err := v.reader.ReadUint32()
		if err != nil {
			return io_common.NewIoParseFailed("VMDシャドウフレーム番号の読み込みに失敗しました", err)
		}
		shadowMode, err := v.reader.ReadUint8()
		if err != nil {
			return io_common.NewIoParseFailed("VMDシャドウ種別の読み込みに失敗しました", err)
		}
		distance, err := v.reader.ReadFloat32()
		if err != nil {
			return io_common.NewIoParseFailed("VMDシャドウ距離の読み込みに失敗しました", err)
		}
		frame := motion.NewShadowFrame(motion.Frame(float32(frameIndex)))
		frame.Read = true
		frame.ShadowMode = int(shadowMode)
		frame.Distance = distance
		motionData.AppendShadowFrame(frame)
	}
	return nil
}

// readIkFrames はIKフレームを読み込む。
func (v *vmdReader) readIkFrames(motionData *motion.VmdMotion) error {
	count, err := v.reader.ReadUint32()
	if err != nil {
		return io_common.NewIoParseFailed("VMDIKフレーム数の読み込みに失敗しました", err)
	}
	for i := 0; i < int(count); i++ {
		frameIndex, err := v.reader.ReadUint32()
		if err != nil {
			return io_common.NewIoParseFailed("VMDIKフレーム番号の読み込みに失敗しました", err)
		}
		frame := motion.NewIkFrame(motion.Frame(float32(frameIndex)))
		frame.Read = true
		visible, err := v.reader.ReadUint8()
		if err != nil {
			return io_common.NewIoParseFailed("VMDIK表示フラグの読み込みに失敗しました", err)
		}
		frame.Visible = visible == 1
		ikCount, err := v.reader.ReadUint32()
		if err != nil {
			return io_common.NewIoParseFailed("VMDIK数の読み込みに失敗しました", err)
		}
		frame.IkList = make([]*motion.IkEnabledFrame, 0, ikCount)
		for j := 0; j < int(ikCount); j++ {
			nameRaw, err := v.reader.ReadBytes(20)
			if err != nil {
				return io_common.NewIoParseFailed("VMDIKボーン名の読み込みに失敗しました", err)
			}
			name, err := io_common.DecodeShiftJISFixed(nameRaw)
			if err != nil {
				return io_common.NewIoParseFailed("VMDIKボーン名のデコードに失敗しました", err)
			}
			enabled, err := v.reader.ReadUint8()
			if err != nil {
				return io_common.NewIoParseFailed("VMDIK有効フラグの読み込みに失敗しました", err)
			}
			item := motion.NewIkEnabledFrame(frame.Index(), name)
			item.Enabled = enabled == 1
			frame.IkList = append(frame.IkList, item)
		}
		motionData.AppendIkFrame(frame)
	}
	return nil
}
