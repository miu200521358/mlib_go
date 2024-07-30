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
func (rep *VmdRepository) Load(path string) (core.IHashModel, error) {
	// モデルを新規作成
	motion := rep.newFunc(path)

	// ファイルを開く
	err := rep.open(path)
	if err != nil {
		mlog.E("Load.Open error: %v", err)
		return motion, err
	}

	err = rep.readHeader(motion)
	if err != nil {
		mlog.E("Load.readHeader error: %v", err)
		return motion, err
	}

	err = rep.loadModel(motion)
	if err != nil {
		mlog.E("Load.loadModel error: %v", err)
		return motion, err
	}

	motion.UpdateHash()
	rep.close()

	return motion, nil
}

func (rep *VmdRepository) LoadName(path string) (string, error) {
	// モデルを新規作成
	motion := rep.newFunc(path)

	// ファイルを開く
	err := rep.open(path)
	if err != nil {
		mlog.E("LoadName.Open error: %v", err)
		return "", err
	}

	err = rep.readHeader(motion)
	if err != nil {
		mlog.E("LoadName.readHeader error: %v", err)
		return "", err
	}

	rep.close()

	return motion.Name(), nil
}

func (rep *VmdRepository) decodeShiftJIS(fbytes []byte) (string, error) {
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

func (rep *VmdRepository) readText(size int) (string, error) {
	fbytes, err := rep.unpackBytes(size)
	if err != nil {
		mlog.E("ReadText error: %v", err)
		return "", err
	}
	return rep.decodeShiftJIS(fbytes)
}

func (rep *VmdRepository) readHeader(motion *vmd.VmdMotion) error {
	rep.defineEncoding(japanese.ShiftJIS)

	// vmdバージョン
	signature, err := rep.readText(30)
	if err != nil {
		mlog.E("readHeader.ReadText error: %v", err)
		return err
	}
	motion.Signature = signature

	// モデル名
	name, err := rep.readText(20)
	if err != nil {
		mlog.E("readHeader.ReadText error: %v", err)
		return err
	}
	motion.SetName(name)

	return nil
}

func (rep *VmdRepository) loadModel(motion *vmd.VmdMotion) error {
	err := rep.loadBones(motion)
	if err != nil {
		mlog.E("loadModel.readBones error: %v", err)
		return err
	}

	err = rep.loadMorphs(motion)
	if err != nil {
		mlog.E("loadModel.readMorphs error: %v", err)
		return err
	}

	err = rep.loadCameras(motion)
	if err != nil {
		// カメラがなくてもエラーにしないが、後続は読まない
		mlog.D("loadModel.readCameras error: %v", err)
		return nil
	}

	err = rep.loadLights(motion)
	if err != nil {
		// ライトがなくてもエラーにしないが、後続は読まない
		mlog.D("loadModel.readLights error: %v", err)
		return nil
	}

	err = rep.loadShadows(motion)
	if err != nil {
		// シャドウがなくてもエラーにしないが、後続は読まない
		mlog.D("loadModel.readShadows error: %v", err)
		return nil
	}

	err = rep.loadIks(motion)
	if err != nil {
		// IKがなくてもエラーにしないが、後続は読まない
		mlog.D("loadModel.readIks error: %v", err)
		return nil
	}

	return nil
}

func (rep *VmdRepository) loadBones(motion *vmd.VmdMotion) error {
	totalCount, err := rep.unpackUInt()
	if err != nil {
		mlog.E("readBones.totalCount error: %v", err)
		return err
	}

	bfValues := make([]float64, 7)
	for i := 0; i < int(totalCount); i++ {
		bf := &vmd.BoneFrame{
			BaseFrame: vmd.NewFrame(i).(*vmd.BaseFrame),
		}
		bf.Read = true

		// ボーン名
		boneName, err := rep.readText(15)
		if err != nil {
			mlog.E("[%d] readBones.boneName error: %v", i, err)
			return err
		}

		// キーフレ番号
		index, err := rep.unpackUInt()
		if err != nil {
			mlog.E("[%d] readBones.index error: %v", i, err)
			return err
		}
		bf.SetIndex(int(index))

		// 位置X,Y,Z
		// 回転X,Y,Z,W
		bfValues, err = rep.unpackFloats(bfValues, 7)
		if err != nil {
			mlog.E("[%d] readBones.bfValues error: %v", i, err)
			return err
		}

		bf.Position = &mmath.MVec3{X: bfValues[0], Y: bfValues[1], Z: bfValues[2]}
		bf.Rotation = mmath.NewMQuaternionByValues(bfValues[3], bfValues[4], bfValues[5], bfValues[6])

		// 補間曲線
		curves, err := rep.unpackBytes(64)
		if err != nil {
			mlog.E("[%d] readBones.Curves error: %v", i, err)
			return err
		}
		bf.Curves = vmd.NewBoneCurvesByValues(curves)

		motion.AppendRegisteredBoneFrame(boneName, bf)
	}

	return nil
}

func (rep *VmdRepository) loadMorphs(motion *vmd.VmdMotion) error {
	totalCount, err := rep.unpackUInt()
	if err != nil {
		mlog.E("readMorphs.totalCount error: %v", err)
		return err
	}

	for i := 0; i < int(totalCount); i++ {
		mf := vmd.NewMorphFrame(0)
		mf.Registered = true
		mf.Read = true

		// モーフ名
		morphName, err := rep.readText(15)
		if err != nil {
			mlog.E("[%d] readMorphs.morphName error: %v", i, err)
			return err
		}

		// キーフレ番号
		index, err := rep.unpackUInt()
		if err != nil {
			mlog.E("[%d] readMorphs.index error: %v", i, err)
			return err
		}
		mf.SetIndex(int(index))

		// ratio
		mf.Ratio, err = rep.unpackFloat()
		if err != nil {
			mlog.E("[%d] readMorphs.Ratio error: %v", i, err)
			return err
		}

		motion.AppendMorphFrame(morphName, mf)
	}

	return nil
}

func (rep *VmdRepository) loadCameras(motion *vmd.VmdMotion) error {
	totalCount, err := rep.unpackUInt()
	if err != nil {
		mlog.D("readCameras.totalCount error: %v", err)
		return err
	}

	for i := 0; i < int(totalCount); i++ {
		cf := vmd.NewCameraFrame(0)
		cf.Registered = true
		cf.Read = true

		// キーフレ番号
		index, err := rep.unpackUInt()
		if err != nil {
			mlog.E("[%d] readCameras.index error: %v", i, err)
			return err
		}
		cf.SetIndex(int(index))

		// 距離
		cf.Distance, err = rep.unpackFloat()
		if err != nil {
			mlog.E("[%d] readCameras.Distance error: %v", i, err)
			return err
		}

		// 位置X,Y,Z
		cf.Position, err = rep.unpackVec3()
		if err != nil {
			mlog.E("[%d] readCameras.Position error: %v", i, err)
			return err
		}

		// 回転(オイラー角度)
		degrees, err := rep.unpackVec3()
		if err != nil {
			mlog.E("[%d] readCameras.Degrees error: %v", i, err)
			return err
		}
		cf.Rotation = mmath.NewMRotationFromDegrees(degrees)

		// 補間曲線
		curves, err := rep.unpackBytes(24)
		if err != nil {
			mlog.E("[%d] readCameras.Curves error: %v", i, err)
			return err
		}
		cf.Curves = vmd.NewCameraCurvesByValues(curves)

		// 視野角
		viewOfAngle, err := rep.unpackUInt()
		if err != nil {
			mlog.E("[%d] readCameras.ViewOfAngle error: %v", i, err)
			return err
		}
		cf.ViewOfAngle = int(viewOfAngle)

		// パースOFF
		perspective, err := rep.unpackByte()
		if err != nil {
			mlog.E("[%d] readCameras.Perspective error: %v", i, err)
			return err
		}
		cf.IsPerspectiveOff = perspective == 1

		motion.AppendCameraFrame(cf)
	}

	return nil
}

func (rep *VmdRepository) loadLights(motion *vmd.VmdMotion) error {
	totalCount, err := rep.unpackUInt()
	if err != nil {
		mlog.D("readLights.totalCount error: %v", err)
		return err
	}

	for i := 0; i < int(totalCount); i++ {
		lf := vmd.NewLightFrame(0)
		lf.Registered = true
		lf.Read = true

		// キーフレ番号
		index, err := rep.unpackUInt()
		if err != nil {
			mlog.E("[%d] readLights.index error: %v", i, err)
			return err
		}
		lf.SetIndex(int(index))

		// 照明色
		lf.Color, err = rep.unpackVec3()
		if err != nil {
			mlog.E("[%d] readLights.Color error: %v", i, err)
			return err
		}

		// 位置X,Y,Z
		lf.Position, err = rep.unpackVec3()
		if err != nil {
			mlog.E("[%d] readLights.Position error: %v", i, err)
			return err
		}

		motion.AppendLightFrame(lf)
	}

	return nil
}

func (rep *VmdRepository) loadShadows(motion *vmd.VmdMotion) error {
	totalCount, err := rep.unpackUInt()
	if err != nil {
		mlog.D("readShadows.totalCount error: %v", err)
		return err
	}

	for i := 0; i < int(totalCount); i++ {
		sf := vmd.NewShadowFrame(0)
		sf.Registered = true
		sf.Read = true

		// キーフレ番号
		index, err := rep.unpackUInt()
		if err != nil {
			mlog.E("[%d] readShadows.index error: %v", i, err)
			return err
		}
		sf.SetIndex(int(index))

		// セルフ影タイプ
		shadowMode, err := rep.unpackByte()
		if err != nil {
			mlog.E("[%d] readShadows.ShadowMode error: %v", i, err)
			return err
		}
		sf.ShadowMode = int(shadowMode)

		// 距離
		sf.Distance, err = rep.unpackFloat()
		if err != nil {
			mlog.E("[%d] readShadows.Distance error: %v", i, err)
			return err
		}

		motion.AppendShadowFrame(sf)
	}

	return nil
}

func (rep *VmdRepository) loadIks(motion *vmd.VmdMotion) error {
	totalCount, err := rep.unpackUInt()
	if err != nil {
		mlog.D("readIks.totalCount error: %v", err)
		return err
	}

	for i := 0; i < int(totalCount); i++ {
		ikf := vmd.NewIkFrame(0)
		ikf.Registered = true
		ikf.Read = true

		// キーフレ番号
		index, err := rep.unpackUInt()
		if err != nil {
			mlog.E("[%d] readIks.index error: %v", i, err)
			return err
		}
		ikf.SetIndex(int(index))

		// モデル表示
		visible, err := rep.unpackByte()
		if err != nil {
			mlog.E("[%d] readIks.Visible error: %v", i, err)
			return err
		}
		ikf.Visible = visible == 1

		// IKリストの数
		ikCount, err := rep.unpackUInt()
		if err != nil {
			mlog.E("[%d] readIks.IkCount error: %v", i, err)
			return err
		}
		for j := 0; j < int(ikCount); j++ {
			ik := vmd.NewIkEnableFrame(ikf.Index())

			// IKボーン名
			ik.BoneName, err = rep.readText(20)
			if err != nil {
				mlog.E("[%d][%d] readIks.Ik.BoneName error: %v", i, j, err)
				return err
			}

			// IK有効無効
			enabled, err := rep.unpackByte()
			if err != nil {
				mlog.E("[%d][%d] readIks.Ik.Enabled error: %v", i, j, err)
				return err
			}
			ik.Enabled = enabled == 1

			ikf.IkList = append(ikf.IkList, ik)
		}

		motion.AppendIkFrame(ikf)
	}

	return nil
}
