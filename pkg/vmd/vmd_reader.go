package vmd

import (
	"bytes"

	"golang.org/x/text/encoding/japanese"

	"github.com/miu200521358/mlib_go/pkg/mcore"
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
func (r *VmdMotionReader) ReadByFilepath(path string) (mcore.HashModelInterface, error) {
	// モデルを新規作成
	motion := r.createModel(path)

	// ファイルを開く
	err := r.Open(path)
	if err != nil {
		return motion, err
	}

	err = r.readHeader(motion)
	if err != nil {
		return motion, err
	}

	err = r.readData(motion)
	if err != nil {
		return motion, err
	}

	r.Close()

	err = motion.UpdateDigest()
	if err != nil {
		return motion, err
	}

	return motion, nil
}

func (r *VmdMotionReader) ReadNameByFilepath(path string) (string, error) {
	// モデルを新規作成
	motion := r.createModel(path)

	// ファイルを開く
	err := r.Open(path)
	if err != nil {
		return "", err
	}

	err = r.readHeader(motion)
	if err != nil {
		return "", err
	}

	r.Close()

	err = motion.UpdateDigest()
	if err != nil {
		return "", err
	}

	return motion.ModelName, nil
}

func (r *VmdMotionReader) DecodeShiftJIS(fbytes []byte) (string, error) {
	// VMDは空白込みで入っているので、正規表現で空白以降は削除する
	decodedBytes, err := japanese.ShiftJIS.NewDecoder().Bytes(fbytes)
	if err != nil {
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
		return "", err
	}
	return r.DecodeShiftJIS(fbytes)
}

func (r *VmdMotionReader) readHeader(motion *VmdMotion) error {
	r.DefineEncoding(japanese.ShiftJIS)

	// vmdバージョン
	signature, err := r.ReadText(30)
	if err != nil {
		return err
	}
	motion.Signature = signature

	// モデル名
	motion.ModelName, err = r.ReadText(20)
	if err != nil {
		return err
	}

	return nil
}

func (r *VmdMotionReader) readData(motion *VmdMotion) error {
	err := r.readBones(motion)
	if err != nil {
		return err
	}

	err = r.readMorphs(motion)
	if err != nil {
		return err
	}

	err = r.readCameras(motion)
	if err != nil {
		// カメラがなくてもエラーにしないが、後続は読まない
		return nil
	}

	err = r.readLights(motion)
	if err != nil {
		// ライトがなくてもエラーにしないが、後続は読まない
		return nil
	}

	err = r.readShadows(motion)
	if err != nil {
		// シャドウがなくてもエラーにしないが、後続は読まない
		return nil
	}

	err = r.readIks(motion)
	if err != nil {
		// IKがなくてもエラーにしないが、後続は読まない
		return nil
	}

	return nil
}

func (r *VmdMotionReader) readBones(motion *VmdMotion) error {
	totalCount, err := r.UnpackUInt()
	if err != nil {
		return err
	}

	for i := 0; i < int(totalCount); i++ {
		v := NewBoneFrame(0)
		v.Registered = true
		v.Read = true

		// ボーン名
		boneName, err := r.ReadText(15)
		if err != nil {
			return err
		}

		// キーフレ番号
		index, err := r.UnpackUInt()
		if err != nil {
			return err
		}
		v.Index = int(index)

		// 位置X,Y,Z
		v.Position, err = r.UnpackVec3()
		if err != nil {
			return err
		}

		// 回転X,Y,Z,W
		qq, err := r.UnpackQuaternion()
		if err != nil {
			return err
		}
		v.Rotation.SetQuaternion(qq)

		// 補間曲線
		curves, err := r.UnpackBytes(64)
		if err != nil {
			return err
		}
		v.Curves = NewBoneCurvesByValues(curves)

		motion.AppendBoneFrame(boneName, v)
	}

	return nil
}

func (r *VmdMotionReader) readMorphs(motion *VmdMotion) error {
	totalCount, err := r.UnpackUInt()
	if err != nil {
		return err
	}

	for i := 0; i < int(totalCount); i++ {
		v := NewMorphFrame(0)
		v.Registered = true
		v.Read = true

		// モーフ名
		morphName, err := r.ReadText(15)
		if err != nil {
			return err
		}

		// キーフレ番号
		index, err := r.UnpackUInt()
		if err != nil {
			return err
		}
		v.Index = int(index)

		// ratio
		v.Ratio, err = r.UnpackFloat()
		if err != nil {
			return err
		}

		motion.AppendMorphFrame(morphName, v)
	}

	return nil
}

func (r *VmdMotionReader) readCameras(motion *VmdMotion) error {
	totalCount, err := r.UnpackUInt()
	if err != nil {
		return err
	}

	for i := 0; i < int(totalCount); i++ {
		v := NewCameraFrame(0)
		v.Registered = true
		v.Read = true

		// キーフレ番号
		index, err := r.UnpackUInt()
		if err != nil {
			return err
		}
		v.Index = int(index)

		// 距離
		v.Distance, err = r.UnpackFloat()
		if err != nil {
			return err
		}

		// 位置X,Y,Z
		v.Position, err = r.UnpackVec3()
		if err != nil {
			return err
		}

		// 回転(オイラー角度)
		degrees, err := r.UnpackVec3()
		if err != nil {
			return err
		}
		v.Rotation.SetDegrees(degrees)

		// 補間曲線
		curves, err := r.UnpackBytes(24)
		if err != nil {
			return err
		}
		v.Curves = NewCameraCurvesByValues(curves)

		// 視野角
		viewOfAngle, err := r.UnpackUInt()
		if err != nil {
			return err
		}
		v.ViewOfAngle = int(viewOfAngle)

		// パースOFF
		perspective, err := r.UnpackByte()
		if err != nil {
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
		return err
	}

	for i := 0; i < int(totalCount); i++ {
		v := NewLightFrame(0)
		v.Registered = true
		v.Read = true

		// キーフレ番号
		index, err := r.UnpackUInt()
		if err != nil {
			return err
		}
		v.Index = int(index)

		// 照明色
		v.Color, err = r.UnpackVec3()
		if err != nil {
			return err
		}

		// 位置X,Y,Z
		v.Position, err = r.UnpackVec3()
		if err != nil {
			return err
		}

		motion.AppendLightFrame(v)
	}

	return nil
}

func (r *VmdMotionReader) readShadows(motion *VmdMotion) error {
	totalCount, err := r.UnpackUInt()
	if err != nil {
		return err
	}

	for i := 0; i < int(totalCount); i++ {
		v := NewShadowFrame(0)
		v.Registered = true
		v.Read = true

		// キーフレ番号
		index, err := r.UnpackUInt()
		if err != nil {
			return err
		}
		v.Index = int(index)

		// セルフ影タイプ
		shadowMode, err := r.UnpackByte()
		if err != nil {
			return err
		}
		v.ShadowMode = int(shadowMode)

		// 距離
		v.Distance, err = r.UnpackFloat()
		if err != nil {
			return err
		}

		motion.AppendShadowFrame(v)
	}

	return nil
}

func (r *VmdMotionReader) readIks(motion *VmdMotion) error {
	totalCount, err := r.UnpackUInt()
	if err != nil {
		return err
	}

	for i := 0; i < int(totalCount); i++ {
		v := NewIkFrame(0)
		v.Registered = true
		v.Read = true

		// キーフレ番号
		index, err := r.UnpackUInt()
		if err != nil {
			return err
		}
		v.Index = int(index)

		// モデル表示
		visible, err := r.UnpackByte()
		if err != nil {
			return err
		}
		v.Visible = visible == 1

		// IKリストの数
		ikCount, err := r.UnpackUInt()
		if err != nil {
			return err
		}
		for j := 0; j < int(ikCount); j++ {
			ik := NewIkEnableFrame(v.Index)

			// IKボーン名
			ik.BoneName, err = r.ReadText(20)
			if err != nil {
				return err
			}

			// IK有効無効
			enabled, err := r.UnpackByte()
			if err != nil {
				return err
			}
			ik.Enabled = enabled == 1
		}

		motion.AppendIkFrame(v)
	}

	return nil
}