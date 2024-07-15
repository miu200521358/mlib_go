package repository

import (
	"bytes"

	"github.com/miu200521358/mlib_go/pkg/domain/core"
	"github.com/miu200521358/mlib_go/pkg/domain/mmath"
	"github.com/miu200521358/mlib_go/pkg/domain/vmd"
	"github.com/miu200521358/mlib_go/pkg/mutils/mlog"
	"golang.org/x/text/encoding/japanese"
)

// 指定されたパスのファイルからデータを読み込む
func (r *VmdRepository) Load(path string) (core.IHashModel, error) {
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

	err = r.readHeader(motion)
	if err != nil {
		mlog.E("Load.readHeader error: %v", err)
		return motion, err
	}

	err = r.loadModel(motion)
	if err != nil {
		mlog.E("Load.loadModel error: %v", err)
		return motion, err
	}

	r.close()

	return motion, nil
}

func (r *VmdRepository) LoadName(path string) (string, error) {
	// モデルを新規作成
	motion := r.newFunc(path)

	// ファイルを開く
	err := r.open(path)
	if err != nil {
		mlog.E("LoadName.Open error: %v", err)
		return "", err
	}

	err = r.readHeader(motion)
	if err != nil {
		mlog.E("LoadName.readHeader error: %v", err)
		return "", err
	}

	r.close()

	return motion.ModelName, nil
}

func (r *VmdRepository) decodeShiftJIS(fbytes []byte) (string, error) {
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

func (r *VmdRepository) readText(size int) (string, error) {
	fbytes, err := r.unpackBytes(size)
	if err != nil {
		mlog.E("ReadText error: %v", err)
		return "", err
	}
	return r.decodeShiftJIS(fbytes)
}

func (r *VmdRepository) readHeader(motion *vmd.VmdMotion) error {
	r.defineEncoding(japanese.ShiftJIS)

	// vmdバージョン
	signature, err := r.readText(30)
	if err != nil {
		mlog.E("readHeader.ReadText error: %v", err)
		return err
	}
	motion.Signature = signature

	// モデル名
	motion.ModelName, err = r.readText(20)
	if err != nil {
		mlog.E("readHeader.ReadText error: %v", err)
		return err
	}

	return nil
}

func (r *VmdRepository) loadModel(motion *vmd.VmdMotion) error {
	err := r.loadBones(motion)
	if err != nil {
		mlog.E("loadModel.readBones error: %v", err)
		return err
	}

	err = r.loadMorphs(motion)
	if err != nil {
		mlog.E("loadModel.readMorphs error: %v", err)
		return err
	}

	err = r.loadCameras(motion)
	if err != nil {
		// カメラがなくてもエラーにしないが、後続は読まない
		mlog.E("loadModel.readCameras error: %v", err)
		return nil
	}

	err = r.loadLights(motion)
	if err != nil {
		// ライトがなくてもエラーにしないが、後続は読まない
		mlog.E("loadModel.readLights error: %v", err)
		return nil
	}

	err = r.loadShadows(motion)
	if err != nil {
		// シャドウがなくてもエラーにしないが、後続は読まない
		mlog.E("loadModel.readShadows error: %v", err)
		return nil
	}

	err = r.loadIks(motion)
	if err != nil {
		// IKがなくてもエラーにしないが、後続は読まない
		mlog.E("loadModel.readIks error: %v", err)
		return nil
	}

	return nil
}

func (r *VmdRepository) loadBones(motion *vmd.VmdMotion) error {
	totalCount, err := r.unpackUInt()
	if err != nil {
		mlog.E("readBones.totalCount error: %v", err)
		return err
	}

	for i := 0; i < int(totalCount); i++ {
		v := &vmd.BoneFrame{
			BaseFrame: vmd.NewFrame(i).(*vmd.BaseFrame),
		}
		v.Read = true

		// ボーン名
		boneName, err := r.readText(15)
		if err != nil {
			mlog.E("[%d] readBones.boneName error: %v", i, err)
			return err
		}

		// キーフレ番号
		index, err := r.unpackUInt()
		if err != nil {
			mlog.E("[%d] readBones.index error: %v", i, err)
			return err
		}
		v.SetIndex(int(index))

		// 位置X,Y,Z
		position, err := r.unpackVec3()
		if err != nil {
			mlog.E("[%d] readBones.Position error: %v", i, err)
			return err
		}
		v.Position = &position

		// 回転X,Y,Z,W
		qq, err := r.unpackQuaternion()
		if err != nil {
			mlog.E("[%d] readBones.Quaternion error: %v", i, err)
			return err
		}
		v.Rotation = &qq

		// 補間曲線
		curves, err := r.unpackBytes(64)
		if err != nil {
			mlog.E("[%d] readBones.Curves error: %v", i, err)
			return err
		}
		v.Curves = vmd.NewBoneCurvesByValues(curves)

		motion.AppendRegisteredBoneFrame(boneName, v)
	}

	return nil
}

func (r *VmdRepository) loadMorphs(motion *vmd.VmdMotion) error {
	totalCount, err := r.unpackUInt()
	if err != nil {
		mlog.E("readMorphs.totalCount error: %v", err)
		return err
	}

	for i := 0; i < int(totalCount); i++ {
		v := vmd.NewMorphFrame(0)
		v.Registered = true
		v.Read = true

		// モーフ名
		morphName, err := r.readText(15)
		if err != nil {
			mlog.E("[%d] readMorphs.morphName error: %v", i, err)
			return err
		}

		// キーフレ番号
		index, err := r.unpackUInt()
		if err != nil {
			mlog.E("[%d] readMorphs.index error: %v", i, err)
			return err
		}
		v.SetIndex(int(index))

		// ratio
		v.Ratio, err = r.unpackFloat()
		if err != nil {
			mlog.E("[%d] readMorphs.Ratio error: %v", i, err)
			return err
		}

		motion.AppendMorphFrame(morphName, v)
	}

	return nil
}

func (r *VmdRepository) loadCameras(motion *vmd.VmdMotion) error {
	totalCount, err := r.unpackUInt()
	if err != nil {
		mlog.E("readCameras.totalCount error: %v", err)
		return err
	}

	for i := 0; i < int(totalCount); i++ {
		v := vmd.NewCameraFrame(0)
		v.Registered = true
		v.Read = true

		// キーフレ番号
		index, err := r.unpackUInt()
		if err != nil {
			mlog.E("[%d] readCameras.index error: %v", i, err)
			return err
		}
		v.SetIndex(int(index))

		// 距離
		v.Distance, err = r.unpackFloat()
		if err != nil {
			mlog.E("[%d] readCameras.Distance error: %v", i, err)
			return err
		}

		// 位置X,Y,Z
		position, err := r.unpackVec3()
		if err != nil {
			mlog.E("[%d] readCameras.Position error: %v", i, err)
			return err
		}
		v.Position = &position

		// 回転(オイラー角度)
		degrees, err := r.unpackVec3()
		if err != nil {
			mlog.E("[%d] readCameras.Degrees error: %v", i, err)
			return err
		}
		v.Rotation = mmath.NewMRotationFromDegrees(&degrees)

		// 補間曲線
		curves, err := r.unpackBytes(24)
		if err != nil {
			mlog.E("[%d] readCameras.Curves error: %v", i, err)
			return err
		}
		v.Curves = vmd.NewCameraCurvesByValues(curves)

		// 視野角
		viewOfAngle, err := r.unpackUInt()
		if err != nil {
			mlog.E("[%d] readCameras.ViewOfAngle error: %v", i, err)
			return err
		}
		v.ViewOfAngle = int(viewOfAngle)

		// パースOFF
		perspective, err := r.unpackByte()
		if err != nil {
			mlog.E("[%d] readCameras.Perspective error: %v", i, err)
			return err
		}
		v.IsPerspectiveOff = perspective == 1

		motion.AppendCameraFrame(v)
	}

	return nil
}

func (r *VmdRepository) loadLights(motion *vmd.VmdMotion) error {
	totalCount, err := r.unpackUInt()
	if err != nil {
		mlog.E("readLights.totalCount error: %v", err)
		return err
	}

	for i := 0; i < int(totalCount); i++ {
		v := vmd.NewLightFrame(0)
		v.Registered = true
		v.Read = true

		// キーフレ番号
		index, err := r.unpackUInt()
		if err != nil {
			mlog.E("[%d] readLights.index error: %v", i, err)
			return err
		}
		v.SetIndex(int(index))

		// 照明色
		color, err := r.unpackVec3()
		if err != nil {
			mlog.E("[%d] readLights.Color error: %v", i, err)
			return err
		}
		v.Color = &color

		// 位置X,Y,Z
		position, err := r.unpackVec3()
		if err != nil {
			mlog.E("[%d] readLights.Position error: %v", i, err)
			return err
		}
		v.Position = &position

		motion.AppendLightFrame(v)
	}

	return nil
}

func (r *VmdRepository) loadShadows(motion *vmd.VmdMotion) error {
	totalCount, err := r.unpackUInt()
	if err != nil {
		mlog.E("readShadows.totalCount error: %v", err)
		return err
	}

	for i := 0; i < int(totalCount); i++ {
		v := vmd.NewShadowFrame(0)
		v.Registered = true
		v.Read = true

		// キーフレ番号
		index, err := r.unpackUInt()
		if err != nil {
			mlog.E("[%d] readShadows.index error: %v", i, err)
			return err
		}
		v.SetIndex(int(index))

		// セルフ影タイプ
		shadowMode, err := r.unpackByte()
		if err != nil {
			mlog.E("[%d] readShadows.ShadowMode error: %v", i, err)
			return err
		}
		v.ShadowMode = int(shadowMode)

		// 距離
		v.Distance, err = r.unpackFloat()
		if err != nil {
			mlog.E("[%d] readShadows.Distance error: %v", i, err)
			return err
		}

		motion.AppendShadowFrame(v)
	}

	return nil
}

func (r *VmdRepository) loadIks(motion *vmd.VmdMotion) error {
	totalCount, err := r.unpackUInt()
	if err != nil {
		mlog.E("readIks.totalCount error: %v", err)
		return err
	}

	for i := 0; i < int(totalCount); i++ {
		v := vmd.NewIkFrame(0)
		v.Registered = true
		v.Read = true

		// キーフレ番号
		index, err := r.unpackUInt()
		if err != nil {
			mlog.E("[%d] readIks.index error: %v", i, err)
			return err
		}
		v.SetIndex(int(index))

		// モデル表示
		visible, err := r.unpackByte()
		if err != nil {
			mlog.E("[%d] readIks.Visible error: %v", i, err)
			return err
		}
		v.Visible = visible == 1

		// IKリストの数
		ikCount, err := r.unpackUInt()
		if err != nil {
			mlog.E("[%d] readIks.IkCount error: %v", i, err)
			return err
		}
		for j := 0; j < int(ikCount); j++ {
			ik := vmd.NewIkEnableFrame(v.GetIndex())

			// IKボーン名
			ik.BoneName, err = r.readText(20)
			if err != nil {
				mlog.E("[%d][%d] readIks.Ik.BoneName error: %v", i, j, err)
				return err
			}

			// IK有効無効
			enabled, err := r.unpackByte()
			if err != nil {
				mlog.E("[%d][%d] readIks.Ik.Enabled error: %v", i, j, err)
				return err
			}
			ik.Enabled = enabled == 1

			v.IkList = append(v.IkList, ik)
		}

		motion.AppendIkFrame(v)
	}

	return nil
}
