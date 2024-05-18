package vmd

import (
	"bytes"

	"golang.org/x/text/encoding/japanese"

	"github.com/miu200521358/mlib_go/pkg/mcore"
	"github.com/miu200521358/mlib_go/pkg/mmath"
	"github.com/miu200521358/mlib_go/pkg/mutils/mlog"
)

// https://hariganep.seesaa.net/article/201103article_1.html
// https://blog.goo.ne.jp/torisu_tetosuki/e/bc9f1c4d597341b394bd02b64597499d
// https://w.atwiki.jp/kumiho_k/pages/15.html

// VMDリーダー
type VmdMotionReader struct {
	mcore.BaseReader[*VmdMotion]
}

func (r *VmdMotionReader) createModel(path string) *VmdMotion {
	model := NewVmdMotion(path)
	return model
}

// 指定されたパスのファイルからデータを読み込む
func (r *VmdMotionReader) ReadByFilepath(path string) (mcore.IHashModel, error) {
	// モデルを新規作成
	motion := r.createModel(path)

	hash, err := r.ReadHashByFilePath(path)
	if err != nil {
		mlog.E("ReadByFilepath.ReadHashByFilePath error: %v", err)
		return motion, err
	}
	motion.Hash = hash

	// ファイルを開く
	err = r.Open(path)
	if err != nil {
		mlog.E("ReadByFilepath.Open error: %v", err)
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

func (r *VmdMotionReader) ReadNameByFilepath(path string) (string, error) {
	// モデルを新規作成
	motion := r.createModel(path)

	// ファイルを開く
	err := r.Open(path)
	if err != nil {
		mlog.E("ReadNameByFilepath.Open error: %v", err)
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

func (r *VmdMotionReader) DecodeShiftJIS(fbytes []byte) (string, error) {
	// VMDは空白込みで入っているので、正規表現で空白以降は削除する
	decodedBytes, err := japanese.ShiftJIS.NewDecoder().Bytes(fbytes)
	if err != nil {
		mlog.E("DecodeShiftJIS error: %v", err)
		return "", err
	}

	trimBytes := bytes.TrimRight(decodedBytes, "\xfd") // PMDで保存したVMDに入ってる
	trimBytes = bytes.TrimRight(trimBytes, "\x00")

	decodedText := string(trimBytes)

	return decodedText, nil
}

func (r *VmdMotionReader) ReadText(size int) (string, error) {
	fbytes, err := r.UnpackBytes(size)
	if err != nil {
		mlog.E("ReadText error: %v", err)
		return "", err
	}
	return r.DecodeShiftJIS(fbytes)
}

func (r *VmdMotionReader) readHeader(motion *VmdMotion) error {
	r.DefineEncoding(japanese.ShiftJIS)

	// vmdバージョン
	signature, err := r.ReadText(30)
	if err != nil {
		mlog.E("readHeader.ReadText error: %v", err)
		return err
	}
	motion.Signature = signature

	// モデル名
	motion.ModelName, err = r.ReadText(20)
	if err != nil {
		mlog.E("readHeader.ReadText error: %v", err)
		return err
	}

	return nil
}

func (r *VmdMotionReader) readData(motion *VmdMotion) error {
	err := r.readBones(motion)
	if err != nil {
		mlog.E("readData.readBones error: %v", err)
		return err
	}

	err = r.readMorphs(motion)
	if err != nil {
		mlog.E("readData.readMorphs error: %v", err)
		return err
	}

	err = r.readCameras(motion)
	if err != nil {
		// カメラがなくてもエラーにしないが、後続は読まない
		mlog.E("readData.readCameras error: %v", err)
		return nil
	}

	err = r.readLights(motion)
	if err != nil {
		// ライトがなくてもエラーにしないが、後続は読まない
		mlog.E("readData.readLights error: %v", err)
		return nil
	}

	err = r.readShadows(motion)
	if err != nil {
		// シャドウがなくてもエラーにしないが、後続は読まない
		mlog.E("readData.readShadows error: %v", err)
		return nil
	}

	err = r.readIks(motion)
	if err != nil {
		// IKがなくてもエラーにしないが、後続は読まない
		mlog.E("readData.readIks error: %v", err)
		return nil
	}

	return nil
}

func (r *VmdMotionReader) readBones(motion *VmdMotion) error {
	totalCount, err := r.UnpackUInt()
	if err != nil {
		mlog.E("readBones.totalCount error: %v", err)
		return err
	}

	for i := 0; i < int(totalCount); i++ {
		v := &BoneFrame{
			BaseFrame: NewFrame(i).(*BaseFrame),
		}
		v.Read = true

		// ボーン名
		boneName, err := r.ReadText(15)
		if err != nil {
			mlog.E("[%d] readBones.boneName error: %v", i, err)
			return err
		}

		// キーフレ番号
		index, err := r.UnpackUInt()
		if err != nil {
			mlog.E("[%d] readBones.index error: %v", i, err)
			return err
		}
		v.SetIndex(int(index))

		// 位置X,Y,Z
		position, err := r.UnpackVec3(true)
		if err != nil {
			mlog.E("[%d] readBones.Position error: %v", i, err)
			return err
		}
		v.Position = &position

		// 回転X,Y,Z,W
		qq, err := r.UnpackQuaternion(true)
		if err != nil {
			mlog.E("[%d] readBones.Quaternion error: %v", i, err)
			return err
		}
		v.Rotation = mmath.NewRotationByQuaternion(&qq)

		// 補間曲線
		curves, err := r.UnpackBytes(64)
		if err != nil {
			mlog.E("[%d] readBones.Curves error: %v", i, err)
			return err
		}
		v.Curves = NewBoneCurvesByValues(curves)

		motion.AppendRegisteredBoneFrame(boneName, v)
	}

	return nil
}

func (r *VmdMotionReader) readMorphs(motion *VmdMotion) error {
	totalCount, err := r.UnpackUInt()
	if err != nil {
		mlog.E("readMorphs.totalCount error: %v", err)
		return err
	}

	for i := 0; i < int(totalCount); i++ {
		v := NewMorphFrame(0)
		v.Registered = true
		v.Read = true

		// モーフ名
		morphName, err := r.ReadText(15)
		if err != nil {
			mlog.E("[%d] readMorphs.morphName error: %v", i, err)
			return err
		}

		// キーフレ番号
		index, err := r.UnpackUInt()
		if err != nil {
			mlog.E("[%d] readMorphs.index error: %v", i, err)
			return err
		}
		v.SetIndex(int(index))

		// ratio
		v.Ratio, err = r.UnpackFloat()
		if err != nil {
			mlog.E("[%d] readMorphs.Ratio error: %v", i, err)
			return err
		}

		motion.AppendMorphFrame(morphName, v)
	}

	return nil
}

func (r *VmdMotionReader) readCameras(motion *VmdMotion) error {
	totalCount, err := r.UnpackUInt()
	if err != nil {
		mlog.E("readCameras.totalCount error: %v", err)
		return err
	}

	for i := 0; i < int(totalCount); i++ {
		v := NewCameraFrame(0)
		v.Registered = true
		v.Read = true

		// キーフレ番号
		index, err := r.UnpackUInt()
		if err != nil {
			mlog.E("[%d] readCameras.index error: %v", i, err)
			return err
		}
		v.SetIndex(int(index))

		// 距離
		v.Distance, err = r.UnpackFloat()
		if err != nil {
			mlog.E("[%d] readCameras.Distance error: %v", i, err)
			return err
		}

		// 位置X,Y,Z
		position, err := r.UnpackVec3(true)
		if err != nil {
			mlog.E("[%d] readCameras.Position error: %v", i, err)
			return err
		}
		v.Position = &position

		// 回転(オイラー角度)
		degrees, err := r.UnpackVec3(true)
		if err != nil {
			mlog.E("[%d] readCameras.Degrees error: %v", i, err)
			return err
		}
		v.Rotation = mmath.NewRotationByDegrees(&degrees)

		// 補間曲線
		curves, err := r.UnpackBytes(24)
		if err != nil {
			mlog.E("[%d] readCameras.Curves error: %v", i, err)
			return err
		}
		v.Curves = NewCameraCurvesByValues(curves)

		// 視野角
		viewOfAngle, err := r.UnpackUInt()
		if err != nil {
			mlog.E("[%d] readCameras.ViewOfAngle error: %v", i, err)
			return err
		}
		v.ViewOfAngle = int(viewOfAngle)

		// パースOFF
		perspective, err := r.UnpackByte()
		if err != nil {
			mlog.E("[%d] readCameras.Perspective error: %v", i, err)
			return err
		}
		v.IsPerspectiveOff = perspective == 1

		motion.AppendCameraFrame(v)
	}

	return nil
}

func (r *VmdMotionReader) readLights(motion *VmdMotion) error {
	totalCount, err := r.UnpackUInt()
	if err != nil {
		mlog.E("readLights.totalCount error: %v", err)
		return err
	}

	for i := 0; i < int(totalCount); i++ {
		v := NewLightFrame(0)
		v.Registered = true
		v.Read = true

		// キーフレ番号
		index, err := r.UnpackUInt()
		if err != nil {
			mlog.E("[%d] readLights.index error: %v", i, err)
			return err
		}
		v.SetIndex(int(index))

		// 照明色
		color, err := r.UnpackVec3(false)
		if err != nil {
			mlog.E("[%d] readLights.Color error: %v", i, err)
			return err
		}
		v.Color = &color

		// 位置X,Y,Z
		position, err := r.UnpackVec3(true)
		if err != nil {
			mlog.E("[%d] readLights.Position error: %v", i, err)
			return err
		}
		v.Position = &position

		motion.AppendLightFrame(v)
	}

	return nil
}

func (r *VmdMotionReader) readShadows(motion *VmdMotion) error {
	totalCount, err := r.UnpackUInt()
	if err != nil {
		mlog.E("readShadows.totalCount error: %v", err)
		return err
	}

	for i := 0; i < int(totalCount); i++ {
		v := NewShadowFrame(0)
		v.Registered = true
		v.Read = true

		// キーフレ番号
		index, err := r.UnpackUInt()
		if err != nil {
			mlog.E("[%d] readShadows.index error: %v", i, err)
			return err
		}
		v.SetIndex(int(index))

		// セルフ影タイプ
		shadowMode, err := r.UnpackByte()
		if err != nil {
			mlog.E("[%d] readShadows.ShadowMode error: %v", i, err)
			return err
		}
		v.ShadowMode = int(shadowMode)

		// 距離
		v.Distance, err = r.UnpackFloat()
		if err != nil {
			mlog.E("[%d] readShadows.Distance error: %v", i, err)
			return err
		}

		motion.AppendShadowFrame(v)
	}

	return nil
}

func (r *VmdMotionReader) readIks(motion *VmdMotion) error {
	totalCount, err := r.UnpackUInt()
	if err != nil {
		mlog.E("readIks.totalCount error: %v", err)
		return err
	}

	for i := 0; i < int(totalCount); i++ {
		v := NewIkFrame(0)
		v.Registered = true
		v.Read = true

		// キーフレ番号
		index, err := r.UnpackUInt()
		if err != nil {
			mlog.E("[%d] readIks.index error: %v", i, err)
			return err
		}
		v.SetIndex(int(index))

		// モデル表示
		visible, err := r.UnpackByte()
		if err != nil {
			mlog.E("[%d] readIks.Visible error: %v", i, err)
			return err
		}
		v.Visible = visible == 1

		// IKリストの数
		ikCount, err := r.UnpackUInt()
		if err != nil {
			mlog.E("[%d] readIks.IkCount error: %v", i, err)
			return err
		}
		for j := 0; j < int(ikCount); j++ {
			ik := NewIkEnableFrame(v.GetIndex())

			// IKボーン名
			ik.BoneName, err = r.ReadText(20)
			if err != nil {
				mlog.E("[%d][%d] readIks.Ik.BoneName error: %v", i, j, err)
				return err
			}

			// IK有効無効
			enabled, err := r.UnpackByte()
			if err != nil {
				mlog.E("[%d][%d] readIks.Ik.Enabled error: %v", i, j, err)
				return err
			}
			ik.Enabled = enabled == 1
		}

		motion.AppendIkFrame(v)
	}

	return nil
}
