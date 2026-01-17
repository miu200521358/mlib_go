// 指示: miu200521358
package vmd

import (
	"io"
	"math"

	"github.com/miu200521358/mlib_go/pkg/adapter/io_common"
	"github.com/miu200521358/mlib_go/pkg/domain/mmath"
	"github.com/miu200521358/mlib_go/pkg/domain/motion"
)

// vmdWriter はVMD書き込み処理を表す。
type vmdWriter struct {
	writer *io_common.BinaryWriter
}

// newVmdWriter はvmdWriterを生成する。
func newVmdWriter(w io.Writer) *vmdWriter {
	return &vmdWriter{writer: io_common.NewBinaryWriter(w)}
}

// Write はVMDを書き込む。
func (v *vmdWriter) Write(motionData *motion.VmdMotion) error {
	if motionData == nil {
		return io_common.NewIoSaveFailed("VMDモーションがnilです", nil)
	}
	if err := v.writeHeader(motionData); err != nil {
		return err
	}
	if err := v.writeBoneFrames(motionData); err != nil {
		return err
	}
	if err := v.writeMorphFrames(motionData); err != nil {
		return err
	}
	if err := v.writeCameraFrames(motionData); err != nil {
		return err
	}
	if err := v.writeLightFrames(motionData); err != nil {
		return err
	}
	if err := v.writeShadowFrames(motionData); err != nil {
		return err
	}
	if err := v.writeIkFrames(motionData); err != nil {
		return err
	}
	return nil
}

// writeHeader はヘッダを書き込む。
func (v *vmdWriter) writeHeader(motionData *motion.VmdMotion) error {
	signature, err := io_common.EncodeShiftJISFixed("Vocaloid Motion Data 0002", 30)
	if err != nil {
		return io_common.NewIoSaveFailed("VMD署名の書き込みに失敗しました", err)
	}
	if err := v.writer.WriteBytes(signature); err != nil {
		return io_common.NewIoSaveFailed("VMD署名の書き込みに失敗しました", err)
	}
	nameBytes, err := io_common.EncodeShiftJISFixed(motionData.Name(), 20)
	if err != nil {
		return io_common.NewIoNameEncodeFailed("VMDモデル名のエンコードに失敗しました", err)
	}
	if err := v.writer.WriteBytes(nameBytes); err != nil {
		return io_common.NewIoSaveFailed("VMDモデル名の書き込みに失敗しました", err)
	}
	return nil
}

// writeBoneFrames はボーンフレームを書き込む。
func (v *vmdWriter) writeBoneFrames(motionData *motion.VmdMotion) error {
	names := motionData.BoneFrames.Names()
	count := 0
	for _, name := range names {
		frames := motionData.BoneFrames.Get(name)
		if frames != nil {
			count += frames.Len()
		}
	}
	if err := v.writer.WriteUint32(uint32(count)); err != nil {
		return io_common.NewIoSaveFailed("VMDボーンフレーム数の書き込みに失敗しました", err)
	}
	for _, name := range names {
		frames := motionData.BoneFrames.Get(name)
		if frames == nil || frames.Len() == 0 {
			continue
		}
		maxFrame := frames.MaxFrame()
		frame := frames.Get(maxFrame)
		if err := v.writeBoneFrame(name, frame); err != nil {
			return err
		}
	}
	for _, name := range names {
		frames := motionData.BoneFrames.Get(name)
		if frames == nil || frames.Len() <= 1 {
			continue
		}
		maxFrame := frames.MaxFrame()
		frames.ForEach(func(frameIndex motion.Frame, frame *motion.BoneFrame) bool {
			if frameIndex >= maxFrame {
				return true
			}
			if err := v.writeBoneFrame(name, frame); err != nil {
				return false
			}
			return true
		})
	}
	return nil
}

// writeBoneFrame はボーンフレームを書き込む。
func (v *vmdWriter) writeBoneFrame(name string, frame *motion.BoneFrame) error {
	if frame == nil {
		return io_common.NewIoSaveFailed("VMDボーンフレームがnilです", nil)
	}
	nameBytes, err := io_common.EncodeShiftJISFixed(name, 15)
	if err != nil {
		return io_common.NewIoNameEncodeFailed("VMDボーン名のエンコードに失敗しました", err)
	}
	if err := v.writer.WriteBytes(nameBytes); err != nil {
		return io_common.NewIoSaveFailed("VMDボーン名の書き込みに失敗しました", err)
	}
	if err := v.writeFrameIndex(frame.Index()); err != nil {
		return err
	}
	pos := mmath.Vec3{}
	if frame.Position != nil {
		pos = *frame.Position
	}
	if err := v.writeVec3(pos); err != nil {
		return err
	}
	rot := mmath.NewQuaternion()
	if frame.Rotation != nil {
		rot = *frame.Rotation
	}
	if err := v.writeQuaternion(rot); err != nil {
		return err
	}
	curves := buildBoneCurves(frame)
	if err := v.writer.WriteBytes(curves); err != nil {
		return io_common.NewIoSaveFailed("VMDボーン補間曲線の書き込みに失敗しました", err)
	}
	return nil
}

// writeMorphFrames はモーフフレームを書き込む。
func (v *vmdWriter) writeMorphFrames(motionData *motion.VmdMotion) error {
	names := motionData.MorphFrames.Names()
	count := 0
	for _, name := range names {
		frames := motionData.MorphFrames.Get(name)
		if frames != nil {
			count += frames.Len()
		}
	}
	if err := v.writer.WriteUint32(uint32(count)); err != nil {
		return io_common.NewIoSaveFailed("VMDモーフフレーム数の書き込みに失敗しました", err)
	}
	for _, name := range names {
		frames := motionData.MorphFrames.Get(name)
		if frames == nil || frames.Len() == 0 {
			continue
		}
		frames.ForEach(func(frameIndex motion.Frame, frame *motion.MorphFrame) bool {
			if err := v.writeMorphFrame(name, frame); err != nil {
				return false
			}
			return true
		})
	}
	return nil
}

// writeMorphFrame はモーフフレームを書き込む。
func (v *vmdWriter) writeMorphFrame(name string, frame *motion.MorphFrame) error {
	if frame == nil {
		return io_common.NewIoSaveFailed("VMDモーフフレームがnilです", nil)
	}
	nameBytes, err := io_common.EncodeShiftJISFixed(name, 15)
	if err != nil {
		return io_common.NewIoNameEncodeFailed("VMDモーフ名のエンコードに失敗しました", err)
	}
	if err := v.writer.WriteBytes(nameBytes); err != nil {
		return io_common.NewIoSaveFailed("VMDモーフ名の書き込みに失敗しました", err)
	}
	if err := v.writeFrameIndex(frame.Index()); err != nil {
		return err
	}
	if err := v.writer.WriteFloat32(frame.Ratio, 0, false); err != nil {
		return io_common.NewIoSaveFailed("VMDモーフ比率の書き込みに失敗しました", err)
	}
	return nil
}

// writeCameraFrames はカメラフレームを書き込む。
func (v *vmdWriter) writeCameraFrames(motionData *motion.VmdMotion) error {
	if err := v.writer.WriteUint32(uint32(motionData.CameraFrames.Len())); err != nil {
		return io_common.NewIoSaveFailed("VMDカメラフレーム数の書き込みに失敗しました", err)
	}
	motionData.CameraFrames.ForEach(func(frameIndex motion.Frame, frame *motion.CameraFrame) bool {
		if err := v.writeCameraFrame(frame); err != nil {
			return false
		}
		return true
	})
	return nil
}

// writeCameraFrame はカメラフレームを書き込む。
func (v *vmdWriter) writeCameraFrame(frame *motion.CameraFrame) error {
	if frame == nil {
		return io_common.NewIoSaveFailed("VMDカメラフレームがnilです", nil)
	}
	if err := v.writeFrameIndex(frame.Index()); err != nil {
		return err
	}
	if err := v.writer.WriteFloat32(frame.Distance, 0, false); err != nil {
		return io_common.NewIoSaveFailed("VMDカメラ距離の書き込みに失敗しました", err)
	}
	pos := mmath.Vec3{}
	if frame.Position != nil {
		pos = *frame.Position
	}
	if err := v.writeVec3(pos); err != nil {
		return err
	}
	deg := mmath.Vec3{}
	if frame.Degrees != nil {
		deg = *frame.Degrees
	}
	if err := v.writeVec3(deg); err != nil {
		return err
	}
	curves := buildCameraCurves(frame)
	if err := v.writer.WriteBytes(curves); err != nil {
		return io_common.NewIoSaveFailed("VMDカメラ補間曲線の書き込みに失敗しました", err)
	}
	if err := v.writer.WriteUint32(uint32(frame.ViewOfAngle)); err != nil {
		return io_common.NewIoSaveFailed("VMDカメラ視野角の書き込みに失敗しました", err)
	}
	if err := v.writer.WriteUint8(boolToByte(frame.IsPerspectiveOff)); err != nil {
		return io_common.NewIoSaveFailed("VMDカメラパース設定の書き込みに失敗しました", err)
	}
	return nil
}

// writeLightFrames はライトフレームを書き込む。
func (v *vmdWriter) writeLightFrames(motionData *motion.VmdMotion) error {
	if err := v.writer.WriteUint32(uint32(motionData.LightFrames.Len())); err != nil {
		return io_common.NewIoSaveFailed("VMDライトフレーム数の書き込みに失敗しました", err)
	}
	motionData.LightFrames.ForEach(func(frameIndex motion.Frame, frame *motion.LightFrame) bool {
		if err := v.writeLightFrame(frame); err != nil {
			return false
		}
		return true
	})
	return nil
}

// writeLightFrame はライトフレームを書き込む。
func (v *vmdWriter) writeLightFrame(frame *motion.LightFrame) error {
	if frame == nil {
		return io_common.NewIoSaveFailed("VMDライトフレームがnilです", nil)
	}
	if err := v.writeFrameIndex(frame.Index()); err != nil {
		return err
	}
	if err := v.writeVec3(frame.Color); err != nil {
		return err
	}
	if err := v.writeVec3(frame.Position); err != nil {
		return err
	}
	return nil
}

// writeShadowFrames はシャドウフレームを書き込む。
func (v *vmdWriter) writeShadowFrames(motionData *motion.VmdMotion) error {
	if err := v.writer.WriteUint32(uint32(motionData.ShadowFrames.Len())); err != nil {
		return io_common.NewIoSaveFailed("VMDシャドウフレーム数の書き込みに失敗しました", err)
	}
	motionData.ShadowFrames.ForEach(func(frameIndex motion.Frame, frame *motion.ShadowFrame) bool {
		if err := v.writeShadowFrame(frame); err != nil {
			return false
		}
		return true
	})
	return nil
}

// writeShadowFrame はシャドウフレームを書き込む。
func (v *vmdWriter) writeShadowFrame(frame *motion.ShadowFrame) error {
	if frame == nil {
		return io_common.NewIoSaveFailed("VMDシャドウフレームがnilです", nil)
	}
	if err := v.writeFrameIndex(frame.Index()); err != nil {
		return err
	}
	if err := v.writer.WriteUint8(uint8(frame.ShadowMode)); err != nil {
		return io_common.NewIoSaveFailed("VMDシャドウ種別の書き込みに失敗しました", err)
	}
	if err := v.writer.WriteFloat32(frame.Distance, 0, false); err != nil {
		return io_common.NewIoSaveFailed("VMDシャドウ距離の書き込みに失敗しました", err)
	}
	return nil
}

// writeIkFrames はIKフレームを書き込む。
func (v *vmdWriter) writeIkFrames(motionData *motion.VmdMotion) error {
	if err := v.writer.WriteUint32(uint32(motionData.IkFrames.Len())); err != nil {
		return io_common.NewIoSaveFailed("VMDIKフレーム数の書き込みに失敗しました", err)
	}
	motionData.IkFrames.ForEach(func(frameIndex motion.Frame, frame *motion.IkFrame) bool {
		if err := v.writeIkFrame(frame); err != nil {
			return false
		}
		return true
	})
	return nil
}

// writeIkFrame はIKフレームを書き込む。
func (v *vmdWriter) writeIkFrame(frame *motion.IkFrame) error {
	if frame == nil {
		return io_common.NewIoSaveFailed("VMDIKフレームがnilです", nil)
	}
	if err := v.writeFrameIndex(frame.Index()); err != nil {
		return err
	}
	if err := v.writer.WriteUint8(boolToByte(frame.Visible)); err != nil {
		return io_common.NewIoSaveFailed("VMDIK表示フラグの書き込みに失敗しました", err)
	}
	if err := v.writer.WriteUint32(uint32(len(frame.IkList))); err != nil {
		return io_common.NewIoSaveFailed("VMDIK数の書き込みに失敗しました", err)
	}
	for _, item := range frame.IkList {
		if item == nil {
			return io_common.NewIoSaveFailed("VMDIKボーンがnilです", nil)
		}
		nameBytes, err := io_common.EncodeShiftJISFixed(item.BoneName, 20)
		if err != nil {
			return io_common.NewIoNameEncodeFailed("VMDIKボーン名のエンコードに失敗しました", err)
		}
		if err := v.writer.WriteBytes(nameBytes); err != nil {
			return io_common.NewIoSaveFailed("VMDIKボーン名の書き込みに失敗しました", err)
		}
		if err := v.writer.WriteUint8(boolToByte(item.Enabled)); err != nil {
			return io_common.NewIoSaveFailed("VMDIK有効フラグの書き込みに失敗しました", err)
		}
	}
	return nil
}

// writeFrameIndex はフレーム番号を書き込む。
func (v *vmdWriter) writeFrameIndex(index motion.Frame) error {
	value := uint32(math.Round(float64(index)))
	if err := v.writer.WriteUint32(value); err != nil {
		return io_common.NewIoSaveFailed("VMDフレーム番号の書き込みに失敗しました", err)
	}
	return nil
}

// writeVec3 はVec3を書き込む。
func (v *vmdWriter) writeVec3(vec mmath.Vec3) error {
	if err := v.writer.WriteFloat32(vec.X, 0, false); err != nil {
		return io_common.NewIoSaveFailed("VMDVec3の書き込みに失敗しました", err)
	}
	if err := v.writer.WriteFloat32(vec.Y, 0, false); err != nil {
		return io_common.NewIoSaveFailed("VMDVec3の書き込みに失敗しました", err)
	}
	if err := v.writer.WriteFloat32(vec.Z, 0, false); err != nil {
		return io_common.NewIoSaveFailed("VMDVec3の書き込みに失敗しました", err)
	}
	return nil
}

// writeQuaternion はクォータニオンを書き込む。
func (v *vmdWriter) writeQuaternion(q mmath.Quaternion) error {
	if err := v.writer.WriteFloat32(q.X(), 0, false); err != nil {
		return io_common.NewIoSaveFailed("VMDクォータニオンの書き込みに失敗しました", err)
	}
	if err := v.writer.WriteFloat32(q.Y(), 0, false); err != nil {
		return io_common.NewIoSaveFailed("VMDクォータニオンの書き込みに失敗しました", err)
	}
	if err := v.writer.WriteFloat32(q.Z(), 0, false); err != nil {
		return io_common.NewIoSaveFailed("VMDクォータニオンの書き込みに失敗しました", err)
	}
	if err := v.writer.WriteFloat32(q.W(), 0, false); err != nil {
		return io_common.NewIoSaveFailed("VMDクォータニオンの書き込みに失敗しました", err)
	}
	return nil
}

// buildBoneCurves はボーン曲線配列を生成する。
func buildBoneCurves(frame *motion.BoneFrame) []byte {
	if frame == nil {
		return append([]byte(nil), motion.INITIAL_BONE_CURVES...)
	}
	disablePhysics := frame.DisablePhysics != nil && *frame.DisablePhysics
	if frame.Curves == nil {
		curves := append([]byte(nil), motion.INITIAL_BONE_CURVES...)
		if disablePhysics && len(curves) > 3 {
			curves[2] = 99
			curves[3] = 15
		}
		return curves
	}
	values := frame.Curves.Merge(disablePhysics)
	out := make([]byte, len(values))
	for i, value := range values {
		out[i] = byte(minClamp(float64(value)))
	}
	return out
}

// buildCameraCurves はカメラ曲線配列を生成する。
func buildCameraCurves(frame *motion.CameraFrame) []byte {
	if frame == nil || frame.Curves == nil {
		return append([]byte(nil), motion.INITIAL_CAMERA_CURVES...)
	}
	values := frame.Curves.Merge()
	out := make([]byte, len(values))
	for i, value := range values {
		out[i] = byte(minClamp(float64(value)))
	}
	return out
}

// minClamp は0..255に丸める。
func minClamp(value float64) float64 {
	if value < 0 {
		return 0
	}
	if value > 255 {
		return 255
	}
	return value
}

// boolToByte はboolを0/1に変換する。
func boolToByte(value bool) uint8 {
	if value {
		return 1
	}
	return 0
}
