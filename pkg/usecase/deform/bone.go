package deform

import (
	"fmt"
	"log"
	"math"
	"os"
	"path/filepath"
	"slices"
	"time"

	"github.com/miu200521358/mlib_go/pkg/config/mlog"
	"github.com/miu200521358/mlib_go/pkg/domain/delta"
	"github.com/miu200521358/mlib_go/pkg/domain/mmath"
	"github.com/miu200521358/mlib_go/pkg/domain/pmx"
	"github.com/miu200521358/mlib_go/pkg/domain/vmd"
	"github.com/miu200521358/mlib_go/pkg/infrastructure/mfile"
	"github.com/miu200521358/mlib_go/pkg/infrastructure/repository"
)

func getLinkAxis(
	minAngleLimitRadians *mmath.MVec3,
	maxAngleLimitRadians *mmath.MVec3,
	ikTargetLocalPosition, ikLocalPosition *mmath.MVec3,
	frame float32,
	count int,
	loop int,
	linkBoneName string,
	ikMotion *vmd.VmdMotion,
	ikFile *os.File,
) (*mmath.MVec3, *mmath.MVec3) {
	// デバッグ出力条件を1回だけ評価
	isDebug := mlog.IsIkVerbose() && ikMotion != nil && ikFile != nil

	// 回転軸
	linkAxis := ikTargetLocalPosition.Cross(ikLocalPosition).Normalize()

	if isDebug {
		fmt.Fprintf(ikFile,
			"[%.3f][%03d][%s][%05d][linkAxis] %s\n",
			frame, loop, linkBoneName, count-1, linkAxis.MMD().String(),
		)
	}

	// linkMat := linkQuat.ToMat4()
	// if isDebug {
	// 	fmt.Fprintf(ikFile,
	// 		"[%.3f][%03d][%s][%05d][linkMat] %s (x: %s, y: %s, z: %s)\n",
	// 		frame, loop, linkBoneName, count-1, linkMat.String(), linkMat.AxisX().String(), linkMat.AxisY().String(), linkMat.AxisZ().String())
	// }

	if minAngleLimitRadians.IsOnlyX() || maxAngleLimitRadians.IsOnlyX() {
		// X軸のみの制限の場合
		vv := linkAxis.X

		if isDebug {
			fmt.Fprintf(ikFile,
				"[%.3f][%03d][%s][%05d][linkAxis(X軸制限)] vv: %.8f\n",
				frame, loop, linkBoneName, count-1, vv)
		}

		if vv < 0 {
			return mmath.MVec3UnitXNeg, linkAxis
		}
		return mmath.MVec3UnitX, linkAxis
	} else if minAngleLimitRadians.IsOnlyY() || maxAngleLimitRadians.IsOnlyY() {
		// Y軸のみの制限の場合
		vv := linkAxis.Y

		if isDebug {
			fmt.Fprintf(ikFile,
				"[%.3f][%03d][%s][%05d][linkAxis(Y軸制限)] vv: %.8f\n",
				frame, loop, linkBoneName, count-1, vv)
		}

		if vv < 0 {
			return mmath.MVec3UnitYNeg, linkAxis
		}
		return mmath.MVec3UnitY, linkAxis
	} else if minAngleLimitRadians.IsOnlyZ() || maxAngleLimitRadians.IsOnlyZ() {
		// Z軸のみの制限の場合
		vv := linkAxis.Z

		if isDebug {
			fmt.Fprintf(ikFile,
				"[%.3f][%03d][%s][%05d][linkAxis(Z軸制限)] vv: %.8f\n",
				frame, loop, linkBoneName, count-1, vv)
		}

		if vv < 0 {
			return mmath.MVec3UnitZNeg, linkAxis
		}
		return mmath.MVec3UnitZ, linkAxis
	}

	return linkAxis, linkAxis
}

func calcIkLimitQuaternion(
	totalIkQuat *mmath.MQuaternion, // リンクボーンの全体回転量
	minAngleLimitRadians *mmath.MVec3, // 最小軸制限（ラジアン）
	maxAngleLimitRadians *mmath.MVec3, // 最大軸制限（ラジアン）
	xAxisVector *mmath.MVec3, // X軸ベクトル
	yAxisVector *mmath.MVec3, // Y軸ベクトル
	zAxisVector *mmath.MVec3, // Z軸ベクトル
	loop int, // ループ回数
	loopCount int, // ループ総回数
	frame float32, // キーフレーム
	count int, // デバッグ用: キーフレ位置
	linkBoneName string, // デバッグ用: リンクボーン名
	ikMotion *vmd.VmdMotion, // デバッグ用: IKモーション
	ikFile *os.File, // デバッグ用: IKファイル
) (*mmath.MQuaternion, int) {
	// デバッグ出力条件を1回だけ評価
	isDebug := mlog.IsIkVerbose() && ikMotion != nil && ikFile != nil

	ikMat := totalIkQuat.ToMat4()
	if isDebug {
		fmt.Fprintf(ikFile,
			"[%.3f][%03d][%s][%05d][ikMat] %s (x: %s, y: %s, z: %s)\n",
			frame, loop, linkBoneName, count-1, ikMat.String(), ikMat.AxisX().String(), ikMat.AxisY().String(), ikMat.AxisZ().String())
	}

	// 軸回転角度を算出
	if minAngleLimitRadians.X > -mmath.HALF_RAD && maxAngleLimitRadians.X < mmath.HALF_RAD {
		// Z*X*Y順
		// X軸回り
		fSX := -ikMat.AxisZ().Y // sin(θx) = -m32
		fX := math.Asin(fSX)    // X軸回り決定
		fCX := math.Cos(fX)     // cos(θx)

		if isDebug {
			fmt.Fprintf(ikFile, "[%.3f][%03d][%s][%05d][X軸制限] fSX: %f, fX: %f, fCX: %f\n",
				frame, loop, linkBoneName, count-1, fSX, fX, fCX)
		}

		// ジンバルロック回避
		if math.Abs(fX) > mmath.GIMBAL1_RAD {
			if fX < 0 {
				fX = -mmath.GIMBAL1_RAD
			} else {
				fX = mmath.GIMBAL1_RAD
			}
			fCX = math.Cos(fX)

			if isDebug {
				fmt.Fprintf(ikFile, "[%.3f][%03d][%s][%05d][X軸制限-ジンバル] fSX: %f, fX: %f, fCX: %f\n",
					frame, loop, linkBoneName, count-1, fSX, fX, fCX)
			}
		}

		// Y軸回り - fCXで割る計算が2回あるので、逆数を一度だけ計算して再利用
		fCXInv := 1.0 / fCX
		fSY := ikMat.AxisZ().X * fCXInv // sin(θy) = m31 / cos(θx)
		fCY := ikMat.AxisZ().Z * fCXInv // cos(θy) = m33 / cos(θx)
		fY := math.Atan2(fSY, fCY)      // Y軸回り決定

		if isDebug {
			fmt.Fprintf(ikFile, "[%.3f][%03d][%s][%05d][X軸制限-Y軸回り] fSY: %f, fCY: %f, fY: %f\n",
				frame, loop, linkBoneName, count-1, fSY, fCY, fY)
		}

		// Z軸周り - 同様にfCXの逆数を再利用
		fSZ := ikMat.AxisX().Y * fCXInv // sin(θz) = m12 / cos(θx)
		fCZ := ikMat.AxisY().Y * fCXInv // cos(θz) = m22 / cos(θx)
		fZ := math.Atan2(fSZ, fCZ)      // Z軸回り決定

		if isDebug {
			fmt.Fprintf(ikFile, "[%.3f][%03d][%s][%05d][X軸制限-Z軸回り] fSZ: %f, fCZ: %f, fZ: %f\n",
				frame, loop, linkBoneName, count-1, fSZ, fCZ, fZ)
		}

		// 角度の制限
		fX = getIkAxisValue(fX, minAngleLimitRadians.X, maxAngleLimitRadians.X, loop, loopCount,
			frame, count, "X軸制限-X", linkBoneName, ikMotion, ikFile)
		fY = getIkAxisValue(fY, minAngleLimitRadians.Y, maxAngleLimitRadians.Y, loop, loopCount,
			frame, count, "X軸制限-Y", linkBoneName, ikMotion, ikFile)
		fZ = getIkAxisValue(fZ, minAngleLimitRadians.Z, maxAngleLimitRadians.Z, loop, loopCount,
			frame, count, "X軸制限-Z", linkBoneName, ikMotion, ikFile)

		// 決定した角度でベクトルを回転
		xQuat := mmath.NewMQuaternionFromAxisAnglesRotate(xAxisVector, fX)
		if isDebug {
			fmt.Fprintf(ikFile, "[%.3f][%03d][%s][%05d][X軸制限-xQuat] xAxisVector: %s, fX: %f, xQuat: %s(%s)\n",
				frame, loop, linkBoneName, count-1, xAxisVector.String(), fX, xQuat.String(), xQuat.ToMMDDegrees().String())
		}

		yQuat := mmath.NewMQuaternionFromAxisAnglesRotate(yAxisVector, fY)
		if isDebug {
			fmt.Fprintf(ikFile, "[%.3f][%03d][%s][%05d][X軸制限-yQuat] yAxisVector: %s, fY: %f, yQuat: %s(%s)\n",
				frame, loop, linkBoneName, count-1, yAxisVector.String(), fY, yQuat.String(), yQuat.ToMMDDegrees().String())
		}

		zQuat := mmath.NewMQuaternionFromAxisAnglesRotate(zAxisVector, fZ)
		if isDebug {
			fmt.Fprintf(ikFile, "[%.3f][%03d][%s][%05d][X軸制限-zQuat] zAxisVector: %s, fZ: %f, zQuat: %s(%s)\n",
				frame, loop, linkBoneName, count-1, zAxisVector.String(), fZ, zQuat.String(), zQuat.ToMMDDegrees().String())
		}

		return yQuat.Muled(xQuat).Muled(zQuat), count
	} else if minAngleLimitRadians.Y > -mmath.HALF_RAD && maxAngleLimitRadians.Y < mmath.HALF_RAD {
		// X*Y*Z順
		// Y軸回り
		fSY := -ikMat.AxisX().Z // sin(θy) = m13
		fY := math.Asin(fSY)    // Y軸回り決定
		fCY := math.Cos(fY)     // cos(θy)

		if isDebug {
			fmt.Fprintf(ikFile, "[%.3f][%03d][%s][%05d][Y軸制限] fSY: %f, fY: %f, fCY: %f\n",
				frame, loop, linkBoneName, count-1, fSY, fY, fCY)
		}

		// ジンバルロック回避
		if math.Abs(fY) > mmath.GIMBAL1_RAD {
			if fY < 0 {
				fY = -mmath.GIMBAL1_RAD
			} else {
				fY = mmath.GIMBAL1_RAD
			}
			fCY = math.Cos(fY)

			if isDebug {
				fmt.Fprintf(ikFile, "[%.3f][%03d][%s][%05d][Y軸制限-ジンバル] fSY: %f, fY: %f, fCY: %f\n",
					frame, loop, linkBoneName, count-1, fSY, fY, fCY)
			}
		}

		// X軸回り - fCYで割る計算が2回あるので、逆数を一度だけ計算して再利用
		fCYInv := 1.0 / fCY
		fSX := ikMat.AxisY().Z * fCYInv // sin(θx) = m23 / cos(θy)
		fCX := ikMat.AxisZ().Z * fCYInv // cos(θx) = m33 / cos(θy)
		fX := math.Atan2(fSX, fCX)      // X軸回り決定

		if isDebug {
			fmt.Fprintf(ikFile, "[%.3f][%03d][%s][%05d][Y軸制限-X軸回り] fSX: %f, fCX: %f, fX: %f\n",
				frame, loop, linkBoneName, count-1, fSX, fCX, fX)
		}

		// Z軸周り - 同様にfCYの逆数を再利用
		fSZ := ikMat.AxisX().Y * fCYInv // sin(θz) = m12 / cos(θy)
		fCZ := ikMat.AxisX().X * fCYInv // cos(θz) = m11 / cos(θy)
		fZ := math.Atan2(fSZ, fCZ)      // Z軸回り決定

		if isDebug {
			fmt.Fprintf(ikFile, "[%.3f][%03d][%s][%05d][Y軸制限-Z軸回り] fSZ: %f, fCZ: %f, fZ: %f\n",
				frame, loop, linkBoneName, count-1, fSZ, fCZ, fZ)
		}

		// 角度の制限
		fX = getIkAxisValue(fX, minAngleLimitRadians.X, maxAngleLimitRadians.X, loop, loopCount,
			frame, count, "Y軸制限-X", linkBoneName, ikMotion, ikFile)
		fY = getIkAxisValue(fY, minAngleLimitRadians.Y, maxAngleLimitRadians.Y, loop, loopCount,
			frame, count, "Y軸制限-Y", linkBoneName, ikMotion, ikFile)
		fZ = getIkAxisValue(fZ, minAngleLimitRadians.Z, maxAngleLimitRadians.Z, loop, loopCount,
			frame, count, "Y軸制限-Z", linkBoneName, ikMotion, ikFile)

		// 決定した角度でベクトルを回転
		xQuat := mmath.NewMQuaternionFromAxisAnglesRotate(xAxisVector, fX)
		if isDebug {
			fmt.Fprintf(ikFile, "[%.3f][%03d][%s][%05d][Y軸制限-xQuat] xAxisVector: %s, fX: %f, xQuat: %s(%s)\n",
				frame, loop, linkBoneName, count-1, xAxisVector.String(), fX, xQuat.String(), xQuat.ToMMDDegrees().String())
		}

		yQuat := mmath.NewMQuaternionFromAxisAnglesRotate(yAxisVector, fY)
		if isDebug {
			fmt.Fprintf(ikFile, "[%.3f][%03d][%s][%05d][Y軸制限-yQuat] yAxisVector: %s, fY: %f, yQuat: %s(%s)\n",
				frame, loop, linkBoneName, count-1, yAxisVector.String(), fY, yQuat.String(), yQuat.ToMMDDegrees().String())
		}

		zQuat := mmath.NewMQuaternionFromAxisAnglesRotate(zAxisVector, fZ)
		if isDebug {
			fmt.Fprintf(ikFile, "[%.3f][%03d][%s][%05d][Y軸制限-zQuat] zAxisVector: %s, fZ: %f, zQuat: %s(%s)\n",
				frame, loop, linkBoneName, count-1, zAxisVector.String(), fZ, zQuat.String(), zQuat.ToMMDDegrees().String())
		}

		return zQuat.Muled(yQuat).Muled(xQuat), count
	}

	// Y*Z*X順
	// Z軸回り
	fSZ := -ikMat.AxisY().X // sin(θz) = m21
	fZ := math.Asin(fSZ)    // Z軸回り決定
	fCZ := math.Cos(fZ)     // cos(θz)

	if isDebug {
		fmt.Fprintf(ikFile, "[%.3f][%03d][%s][%05d][Z軸制限] fSZ: %f, fZ: %f, fCZ: %f\n",
			frame, loop, linkBoneName, count-1, fSZ, fZ, fCZ)
	}

	// ジンバルロック回避
	if math.Abs(fZ) > mmath.GIMBAL1_RAD {
		if fZ < 0 {
			fZ = -mmath.GIMBAL1_RAD
		} else {
			fZ = mmath.GIMBAL1_RAD
		}
		fCZ = math.Cos(fZ)

		if isDebug {
			fmt.Fprintf(ikFile, "[%.3f][%03d][%s][%05d][Z軸制限-ジンバル] fSZ: %f, fZ: %f, fCZ: %f\n",
				frame, loop, linkBoneName, count-1, fSZ, fZ, fCZ)
		}
	}

	// X軸回り
	fSX := ikMat.AxisY().Z / fCZ // sin(θx) = m23 / cos(θz)
	fCX := ikMat.AxisY().Y / fCZ // cos(θx) = m22 / cos(θz)
	fX := math.Atan2(fSX, fCX)   // X軸回り決定

	if isDebug {
		fmt.Fprintf(ikFile, "[%.3f][%03d][%s][%05d][Z軸制限-X軸回り] fSX: %f, fCX: %f, fX: %f\n",
			frame, loop, linkBoneName, count-1, fSX, fCX, fX)
	}

	// Y軸周り
	fSY := ikMat.AxisZ().X / fCZ // sin(θy) = m31 / cos(θz)
	fCY := ikMat.AxisX().X / fCZ // cos(θy) = m11 / cos(θz)
	fY := math.Atan2(fSY, fCY)   // Y軸回り決定

	if isDebug {
		fmt.Fprintf(ikFile, "[%.3f][%03d][%s][%05d][Z軸制限-Y軸回り] fSY: %f, fCY: %f, fY: %f\n",
			frame, loop, linkBoneName, count-1, fSY, fCY, fY)
	}

	// 角度の制限
	fX = getIkAxisValue(fX, minAngleLimitRadians.X, maxAngleLimitRadians.X, loop, loopCount,
		frame, count, "Z軸制限-X", linkBoneName, ikMotion, ikFile)
	fY = getIkAxisValue(fY, minAngleLimitRadians.Y, maxAngleLimitRadians.Y, loop, loopCount,
		frame, count, "Z軸制限-Y", linkBoneName, ikMotion, ikFile)
	fZ = getIkAxisValue(fZ, minAngleLimitRadians.Z, maxAngleLimitRadians.Z, loop, loopCount,
		frame, count, "Z軸制限-Z", linkBoneName, ikMotion, ikFile)

	// 決定した角度でベクトルを回転
	xQuat := mmath.NewMQuaternionFromAxisAnglesRotate(xAxisVector, fX)
	if isDebug {
		fmt.Fprintf(ikFile, "[%.3f][%03d][%s][%05d][Z軸制限-xQuat] xAxisVector: %s, fX: %f, xQuat: %s(%s)\n",
			frame, loop, linkBoneName, count-1, xAxisVector.String(), fX, xQuat.String(), xQuat.ToMMDDegrees().String())
	}

	yQuat := mmath.NewMQuaternionFromAxisAnglesRotate(yAxisVector, fY)
	if isDebug {
		fmt.Fprintf(ikFile, "[%.3f][%03d][%s][%05d][Z軸制限-yQuat] yAxisVector: %s, fY: %f, yQuat: %s(%s)\n",
			frame, loop, linkBoneName, count-1, yAxisVector.String(), fY, yQuat.String(), yQuat.ToMMDDegrees().String())
	}

	zQuat := mmath.NewMQuaternionFromAxisAnglesRotate(zAxisVector, fZ)
	if isDebug {
		fmt.Fprintf(ikFile, "[%.3f][%03d][%s][%05d][Z軸制限-zQuat] zAxisVector: %s, fZ: %f, zQuat: %s(%s)\n",
			frame, loop, linkBoneName, count-1, zAxisVector.String(), fZ, zQuat.String(), zQuat.ToMMDDegrees().String())
	}

	return xQuat.Muled(zQuat).Muled(yQuat), count
}

func getIkAxisValue(
	fV, minAngleLimit, maxAngleLimit float64,
	loop, loopCount int,
	frame float32,
	count int,
	axisName, linkBoneName string,
	ikMotion *vmd.VmdMotion,
	ikFile *os.File,
) float64 {
	// デバッグ出力条件を1回だけ評価
	isDebug := mlog.IsIkVerbose() && ikMotion != nil && ikFile != nil

	isInLoop := float64(loop) < float64(loopCount)/2.0

	if isDebug {
		fmt.Fprintf(ikFile, "[%.3f][%03d][%s][%05d][%s-getIkAxisValue] loop: %d, isInLoop: %v\n",
			frame, loop, linkBoneName, count-1, axisName, loop, isInLoop)
	}

	if fV < minAngleLimit {
		tf := 2*minAngleLimit - fV
		if tf <= maxAngleLimit && isInLoop {
			if isDebug {
				fmt.Fprintf(ikFile, "[%.3f][%03d][%s][%05d][%s-最小角度(loop内)] minAngleLimit: %f, fV: %f, tf: %f\n",
					frame, loop, linkBoneName, count-1, axisName, minAngleLimit, fV, tf)
			}

			fV = tf
		} else {
			if isDebug {
				fmt.Fprintf(ikFile, "[%.3f][%03d][%s][%05d][%s-最小角度(loop外)] minAngleLimit: %f, fV: %f, tf: %f\n",
					frame, loop, linkBoneName, count-1, axisName, minAngleLimit, fV, tf)
			}

			fV = minAngleLimit
		}
	}

	if fV > maxAngleLimit {
		tf := 2*maxAngleLimit - fV
		if tf >= minAngleLimit && isInLoop {
			if isDebug {
				fmt.Fprintf(ikFile, "[%.3f][%03d][%s][%05d][%s-最大角度(loop内)] maxAngleLimit: %f, fV: %f, tf: %f\n",
					frame, loop, linkBoneName, count-1, axisName, maxAngleLimit, fV, tf)
			}

			fV = tf
		} else {
			if isDebug {
				fmt.Fprintf(ikFile, "[%.3f][%03d][%s][%05d][%s-最大角度(loop外)] maxAngleLimit: %f, fV: %f, tf: %f\n",
					frame, loop, linkBoneName, count-1, axisName, maxAngleLimit, fV, tf)
			}

			fV = maxAngleLimit
		}
	}

	return fV
}

// デフォーム対象ボーン情報一覧生成
func newVmdDeltas(
	model *pmx.PmxModel,
	motion *vmd.VmdMotion,
	deltas *delta.VmdDeltas,
	frame float32,
	boneNames []string,
	isAfterPhysics bool,
) ([]int, *delta.VmdDeltas) {
	// ボーン名の存在チェック用マップ
	targetSortedBones := model.Bones.LayerSortedBones[isAfterPhysics]

	if deltas == nil {
		deltas = delta.NewVmdDeltas(frame, model.Bones, model.Hash(), motion.Hash())
	}

	if len(boneNames) == 1 && model.Bones.ContainsByName(boneNames[0]) {
		// 1ボーン指定の場合
		bone, _ := model.Bones.GetByName(boneNames[0])
		return model.Bones.DeformBoneIndexes[bone.Index()], deltas
	}

	// 変形階層順ボーンIndexリスト
	deformBoneIndexes := make([]int, 0, len(targetSortedBones))

	// 関連ボーンINDEXリスト（順不同）
	relativeBoneIndexes := make(map[int]struct{})

	if len(boneNames) > 0 {
		// 指定ボーンに関連するボーンのみ対象とする
		for _, boneName := range boneNames {
			if !model.Bones.ContainsByName(boneName) {
				continue
			}

			// ボーン
			bone, _ := model.Bones.GetByName(boneName)

			// 対象のボーンは常に追加
			if _, ok := relativeBoneIndexes[bone.Index()]; !ok {
				relativeBoneIndexes[bone.Index()] = struct{}{}
			}

			// 関連するボーンの追加
			for _, index := range bone.RelativeBoneIndexes {
				if _, ok := relativeBoneIndexes[index]; !ok {
					relativeBoneIndexes[index] = struct{}{}
				}
			}
		}
	} else {
		// 物理前かつボーン名の指定が無い場合、物理前全ボーンを対象とする
		for _, bone := range model.Bones.LayerSortedBones[isAfterPhysics] {
			deformBoneIndexes = append(deformBoneIndexes, bone.Index())
			if !deltas.Bones.Contains(bone.Index()) {
				deltas.Bones.Update(&delta.BoneDelta{Bone: bone, Frame: frame})
			}
		}

		return deformBoneIndexes, deltas
	}

	// 変形階層・ボーンINDEXでソート
	for _, boneIndex := range model.Bones.LayerSortedIndexes {
		bone, _ := model.Bones.Get(boneIndex)
		if _, ok := relativeBoneIndexes[bone.Index()]; ok {
			deformBoneIndexes = append(deformBoneIndexes, bone.Index())
			if !deltas.Bones.Contains(bone.Index()) {
				deltas.Bones.Update(&delta.BoneDelta{Bone: bone, Frame: frame})
			}
		}
	}

	return deformBoneIndexes, deltas
}

// デフォーム情報を求めて設定
func fillBoneDeform(
	model *pmx.PmxModel,
	motion *vmd.VmdMotion,
	deltas *delta.VmdDeltas,
	frame float32,
	deformBoneIndexes []int,
	isCalcIk bool,
	isAfterPhysics bool,
) *delta.BoneDeltas {

	// IKのON/OFFフレームを取得
	ikFrame := motion.IkFrames.Get(frame)

	// deformBoneIndexes はすでに「変形対象ボーンのインデックス」のリスト
	// ループで対象ボーンごとにデフォーム情報を更新する
	for i, boneIndex := range deformBoneIndexes {
		// ボーンと対応する BoneDelta をローカル変数にキャッシュ
		bone, _ := model.Bones.Get(boneIndex)
		d := deltas.Bones.Get(boneIndex)
		if d == nil {
			d = &delta.BoneDelta{Bone: bone, Frame: frame}
		}

		// 1. ボーンがEffector系 (子がIKエフェクタなど) でなければ、または計算が必要な場合のみ処理
		//    ただし、通常のローカル行列を更新する必要があるのでそこは常に更新。
		isEffectorRotation := bone.IsEffectorRotation()
		isEffectorTranslation := bone.IsEffectorTranslation()

		// 2. ボーンのUnitMatrix が未初期化（nil）か、あるいはEffector指定のときはフレーム値を再取得する
		if d.UnitMatrix == nil || d.UnitMatrix.IsIdent() || isEffectorRotation || isEffectorTranslation {
			// 該当するキー情報を取得
			bf := motion.BoneFrames.Get(bone.Name()).Get(frame)
			if bf.Position != nil {
				d.FramePosition = bf.Position.Copy()
			}
			if bf.Rotation != nil {
				d.FrameRotation = bf.Rotation.Copy()
			}

			// ボーンの移動位置・回転角度・スケールなどを取得
			// (getLocalMat / getPosition / getRotation / getScale の呼び出し)
			d.FrameLocalMat, d.FrameLocalMorphMat = getLocalMat(deltas, bone)
			d.FrameCancelablePosition, d.FrameMorphPosition, d.FrameMorphCancelablePosition =
				getPosition(deltas, bone, bf)
			d.FrameCancelableRotation, d.FrameMorphRotation, d.FrameMorphCancelableRotation =
				getRotation(deltas, bone, bf)
			d.FrameScale, d.FrameCancelableScale, d.FrameMorphScale, d.FrameMorphCancelableScale =
				getScale(deltas, bone, bf)

			// BoneDelta更新（ローカル行列や拡大行列などを組み合わせる）
			updateBoneDelta(deltas.Bones, d, bone)
		}

		// 3. IK対象ボーンならIK計算を実行
		//    IKが有効であり、かつ isCalcIk が true のときだけ
		if isCalcIk && bone.IsIK() && ikFrame.IsEnable(bone.Name()) {
			// IK変形リスト（IKボーンの子孫にあたるボーンインデックス一覧）
			ikTargetDeformBoneIndexes := model.Bones.DeformBoneIndexes[bone.Index()]

			// 変形リストを再帰的に更新 (IKの前に対象ボーンを先に最新化)
			// IK対象ボーンの子階層がまだ最新でない場合、先に更新する
			deltas.Bones = fillBoneDeform(
				model,
				motion,
				deltas,
				frame,
				ikTargetDeformBoneIndexes,
				false, // IK再帰呼び出ししない
				isAfterPhysics,
			)

			// 親→子の順にグローバル行列を再更新
			updateGlobalMatrix(deltas.Bones, ikTargetDeformBoneIndexes)

			// IK適用前のグローバル行列を保存
			for _, idx := range ikTargetDeformBoneIndexes {
				linkD := deltas.Bones.Get(idx)
				if linkD != nil {
					linkD.GlobalIkOffMatrix = linkD.GlobalMatrix.Copy()
					deltas.Bones.Update(linkD)
				}
			}

			// IKボーンのグローバル位置
			ikBoneDelta := deltas.Bones.Get(bone.Index())
			if ikBoneDelta == nil {
				continue
			}
			ikGlobalPosition := ikBoneDelta.FilledGlobalPosition()

			// IK計算実行
			deformIk(
				model,
				motion,
				deltas,
				frame,
				isAfterPhysics,
				bone,
				ikGlobalPosition,
				ikTargetDeformBoneIndexes,
				i, // deformIndex
			)
		}
	}

	return deltas.Bones
}

func updateBoneDelta(
	boneDeltas *delta.BoneDeltas,
	d *delta.BoneDelta,
	bone *pmx.Bone,
) {
	d.UnitMatrix = mmath.NewMMat4()
	d.GlobalMatrix = nil
	d.LocalMatrix = nil
	d.GlobalPosition = nil

	boneDeltas.Update(d)

	// ローカル行列
	localMat := calculateTotalLocalMat(boneDeltas, bone.Index())
	if localMat != nil && !localMat.IsIdent() {
		d.UnitMatrix.Mul(localMat)
	}

	// スケール(回転が変わるため、先にスケールを計算する)
	scaleMat := calculateTotalScaleMat(boneDeltas, bone.Index())
	if scaleMat != nil && !scaleMat.IsIdent() {
		d.UnitMatrix.Mul(scaleMat)
	}

	// 移動
	posMat := calculateTotalPositionMat(boneDeltas, bone.Index())
	if posMat != nil && !posMat.IsIdent() {
		d.UnitMatrix.Mul(posMat)
	}

	// 回転
	rotMat := calculateTotalRotationMat(boneDeltas, bone.Index())
	if rotMat != nil && !rotMat.IsIdent() {
		d.UnitMatrix.Mul(rotMat)
	}

	// 逆BOf行列(初期姿勢行列)
	d.UnitMatrix = d.Bone.RevertOffsetMatrix.Muled(d.UnitMatrix)

	boneDeltas.Update(d)
}

func updateGlobalMatrix(
	boneDeltas *delta.BoneDeltas,
	deformBoneIndexes []int,
) {
	for _, boneIndex := range deformBoneIndexes {
		delta := boneDeltas.Get(boneIndex)
		if delta == nil {
			break
		}
		if delta.UnitMatrix == nil {
			updateBoneDelta(boneDeltas, delta, delta.Bone)
		}
		delta.GlobalMatrix = nil
		delta.LocalMatrix = nil
		delta.GlobalPosition = nil

		parentDelta := boneDeltas.Get(delta.Bone.ParentIndex)
		if parentDelta != nil && parentDelta.GlobalIkOffMatrix != nil && parentDelta.Bone.IsIK() {
			delta.GlobalMatrix = parentDelta.GlobalIkOffMatrix.Muled(delta.UnitMatrix)
		} else if parentDelta != nil && parentDelta.GlobalMatrix != nil {
			delta.GlobalMatrix = parentDelta.GlobalMatrix.Muled(delta.UnitMatrix)
		} else {
			// 対象ボーン自身の行列をかける
			delta.GlobalMatrix = delta.UnitMatrix.Copy()
		}

		boneDeltas.Update(delta)
	}
}

// deformIk IK計算
func deformIk(
	model *pmx.PmxModel,
	motion *vmd.VmdMotion,
	deltas *delta.VmdDeltas,
	frame float32,
	isAfterPhysics bool,
	ikBone *pmx.Bone,
	ikGlobalPosition *mmath.MVec3,
	ikTargetDeformBoneIndexes []int,
	deformIndex int,
) *delta.BoneDeltas {

	// IKリンクが無ければスルー
	if len(ikBone.Ik.Links) < 1 {
		return deltas.Bones
	}

	// デバッグ関連設定
	isDebug := mlog.IsIkVerbose()
	var prefixPath string
	var ikFile *os.File
	var ikMotion, globalMotion *vmd.VmdMotion
	var err error
	count := 1

	if isDebug {
		// IK計算用出力パスの作成
		dirPath := fmt.Sprintf("%s/IK_step", filepath.Dir(model.Path()))
		if mkdirErr := os.MkdirAll(dirPath, 0755); mkdirErr != nil {
			log.Fatal(mkdirErr)
		}
		_, motionFileName, _ := mfile.SplitPath(motion.Path())
		date := time.Now().Format("20060102_150405")
		prefixPath = fmt.Sprintf("%s/%s_%.2f_%s_%03d", dirPath, motionFileName, frame, date, deformIndex)

		// IKモーション作成
		ikMotionPath := fmt.Sprintf("%s_%s.vmd", prefixPath, ikBone.Name())
		ikMotion = vmd.NewVmdMotion(ikMotionPath)

		// グローバルデバッグ用モーション
		globalMotionPath := fmt.Sprintf("%s_%s_global.vmd", prefixPath, ikBone.Name())
		globalMotion = vmd.NewVmdMotion(globalMotionPath)

		// IKログファイル
		ikLogPath := fmt.Sprintf("%s_%s.log", prefixPath, ikBone.Name())
		ikFile, err = os.OpenFile(ikLogPath, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Fprintln(ikFile, "----------------------------------------")
		fmt.Fprintf(ikFile, "[IK計算出力先][%.3f][%s] %s\n", frame, ikBone.Name(), ikMotionPath)
	}

	isFullDebug := isDebug && ikMotion != nil && ikFile != nil
	isGlobalDebug := isDebug && globalMotion != nil && ikFile != nil

	// 終了時に IKモーション等を保存・後処理
	defer func() {
		mlog.IV("[IK計算終了][%.3f][%s]", frame, ikBone.Name())
		if ikMotion != nil {
			r := repository.NewVmdRepository()
			if saveErr := r.Save("", ikMotion, true); saveErr != nil {
				mlog.E("[IK計算出力失敗][%.3f][%s] %s", frame, ikBone.Name(), saveErr)
			}
		}
		if globalMotion != nil {
			r := repository.NewVmdRepository()
			if saveErr := r.Save("", globalMotion, true); saveErr != nil {
				mlog.E("[IK計算出力失敗][%.3f][%s] %s", frame, ikBone.Name(), saveErr)
			}
		}
		if ikFile != nil {
			ikFile.Close()
		}
	}()

	// 一段IKであるか
	isSingleIk := (len(ikBone.Ik.Links) == 1)

	// ループ回数
	loopCount := max(ikBone.Ik.LoopCount, 1)

	// IKターゲットボーン
	ikTargetBone, _ := model.Bones.Get(ikBone.Ik.BoneIndex)

	// IKターゲットボーンが先に変形されているかによってループ回数を多少調整
	ikTargetDeformIndex := slices.Index(ikTargetDeformBoneIndexes, ikTargetBone.Index())
	ikDeformIndex := slices.Index(ikTargetDeformBoneIndexes, ikBone.Index())
	if ikTargetDeformIndex < ikDeformIndex {
		// 初回に余分に1回まわす
		loopCount++
	}

	// つま先IK (左右) 対策
	var ikOnGlobalPosition *mmath.MVec3
	if ikTargetDeformIndex < ikDeformIndex && isSingleIk &&
		(ikBone.Name() == pmx.TOE_IK.Left() || ikBone.Name() == pmx.TOE_IK.Right()) {

		// IK ON時の位置を退避
		ikOnGlobalPosition = ikGlobalPosition.Copy()

		// IK OFF時のグローバル位置を取得して一旦IK目標にする
		ikOffDeltas := DeformBoneByPhysicsFlag(
			model, motion, nil,
			false, frame,
			[]string{ikTargetBone.Name()},
			isAfterPhysics,
		)
		ikGlobalPosition = ikOffDeltas.Bones.Get(ikTargetBone.Index()).FilledGlobalPosition()

		if isGlobalDebug {
			bf := vmd.NewBoneFrame(float32(count))
			bf.Position = ikGlobalPosition
			globalMotion.AppendRegisteredBoneFrame(ikBone.Name(), bf)
			count++
		}
	}

	// IK 計算本体
ikLoop:
	for loop := range loopCount {
		for linkIndex, ikLink := range ikBone.Ik.Links {
			// リンクボーンを取得
			if !model.Bones.Contains(ikLink.BoneIndex) {
				continue
			}
			linkBone, _ := model.Bones.Get(ikLink.BoneIndex)

			// 角度制限が完全0の場合はスキップ
			if (linkBone.AngleLimit &&
				linkBone.MinAngleLimit.IsZero() &&
				linkBone.MaxAngleLimit.IsZero()) ||
				(linkBone.LocalAngleLimit &&
					linkBone.LocalMinAngleLimit.IsZero() &&
					linkBone.LocalMaxAngleLimit.IsZero()) {
				continue
			}

			if isFullDebug {
				fmt.Fprintf(ikFile,
					"[%.3f][%03d][%s][%05d] -------------------------------------------- \n",
					frame, loop, linkBone.Name(), count-1)
			}

			// リンクボーンの現在回転を取得
			linkDelta := deltas.Bones.Get(linkBone.Index())
			if linkDelta == nil {
				linkDelta = &delta.BoneDelta{Bone: linkBone, Frame: frame}
			}
			linkQuat := linkDelta.FilledTotalRotation()

			if isFullDebug {
				bf := vmd.NewBoneFrame(float32(count))
				bf.Rotation = linkQuat.Copy()
				ikMotion.AppendRegisteredBoneFrame(linkBone.Name(), bf)
				count++
				fmt.Fprintf(ikFile,
					"[%.3f][%03d][%s][%05d][linkQuat] %s(%s)\n",
					frame, loop, linkBone.Name(), count-1,
					linkQuat.String(), linkQuat.ToMMDDegrees().String())
			}

			// 現在のIKターゲットボーンのグローバル位置
			ikTargetGlobalPos := deltas.Bones.Get(ikTargetBone.Index()).FilledGlobalPosition()

			// リンクボーンのグローバル行列を逆行列化
			linkGlobalMat := deltas.Bones.Get(linkBone.Index()).FilledGlobalMatrix()
			linkInvMatrix := linkGlobalMat.Inverted()

			if isGlobalDebug {
				// 位置をそれぞれ出力
				bfIk := vmd.NewBoneFrame(float32(count))
				bfIk.Position = ikGlobalPosition
				globalMotion.AppendRegisteredBoneFrame(ikBone.Name(), bfIk)

				bfTgt := vmd.NewBoneFrame(float32(count))
				bfTgt.Position = ikTargetGlobalPos
				globalMotion.AppendRegisteredBoneFrame(ikTargetBone.Name(), bfTgt)

				count++
			}

			// つま先IKで初回ループが終わった後、IK ON位置へ戻す
			if loop == 1 && linkIndex == 0 && ikTargetDeformIndex < ikDeformIndex && ikOnGlobalPosition != nil {
				ikGlobalPosition = ikOnGlobalPosition
				if isGlobalDebug {
					bfLink := vmd.NewBoneFrame(float32(count))
					bfLink.Position = deltas.Bones.Get(linkBone.Index()).FilledGlobalPosition()
					globalMotion.AppendRegisteredBoneFrame(linkBone.Name(), bfLink)

					bfIk := vmd.NewBoneFrame(float32(count))
					bfIk.Position = ikGlobalPosition
					globalMotion.AppendRegisteredBoneFrame(ikBone.Name(), bfIk)

					bfTgt := vmd.NewBoneFrame(float32(count))
					bfTgt.Position = ikTargetGlobalPos
					globalMotion.AppendRegisteredBoneFrame(ikTargetBone.Name(), bfTgt)
					count++
				}
			}

			if isFullDebug {
				fmt.Fprintf(ikFile,
					"[%.3f][%03d][%s][%05d][Global] [%s]ikGlobalPosition: %s, [%s]ikTargetGlobalPosition: %s, [%s]linkGlobalPosition: %s\n",
					frame, loop, linkBone.Name(), count-1,
					ikBone.Name(), ikGlobalPosition.MMD().String(),
					ikTargetBone.Name(), ikTargetGlobalPos.MMD().String(),
					linkBone.Name(), deltas.Bones.Get(linkBone.Index()).FilledGlobalPosition().MMD().String())
			}

			// 注目ノードを起点としたローカル座標系へ変換
			ikTargetLocalPos := linkInvMatrix.MulVec3(ikTargetGlobalPos).Normalize()
			ikLocalPos := linkInvMatrix.MulVec3(ikGlobalPosition).Normalize()

			if isFullDebug {
				dist := ikTargetLocalPos.Distance(ikLocalPos)
				fmt.Fprintf(ikFile,
					"[%.3f][%03d][%s][%05d][Local] ikTargetLocalPos: %s, ikLocalPos: %s (%f)\n",
					frame, loop, linkBone.Name(), count-1,
					ikTargetLocalPos.MMD().String(), ikLocalPos.MMD().String(),
					dist)
			}

			// ターゲットとIK位置がほぼ同じなら抜ける
			distanceThreshold := ikTargetLocalPos.Distance(ikLocalPos)
			if distanceThreshold < 1e-5 {
				if isFullDebug {
					fmt.Fprintf(ikFile,
						"[%.3f][%03d][%s][%05d][Local] ***BREAK*** distanceThreshold: %f\n",
						frame, loop, linkBone.Name(), count-1, distanceThreshold)
				}
				break ikLoop
			}

			// IK回転角度計算 (unitRad はリンクによって重み付け)
			unitRad := ikBone.Ik.UnitRotation.X * float64(linkIndex+1)
			linkDot := ikLocalPos.Dot(ikTargetLocalPos)

			if isFullDebug {
				fmt.Fprintf(ikFile,
					"[%.3f][%03d][%s][%05d][回転角度] unitRad: %.8f (%.5f), linkDot: %.8f\n",
					frame, loop, linkBone.Name(), count-1, unitRad,
					mmath.RadToDeg(unitRad), linkDot)
			}

			originalLinkAngle := math.Acos(mmath.Clamped(linkDot, -1, 1))
			linkAngle := mmath.Clamped(originalLinkAngle, -unitRad, unitRad)

			if isFullDebug {
				fmt.Fprintf(ikFile,
					"[%.3f][%03d][%s][%05d][単位角制限] linkAngle: %.8f(%.5f), originalLinkAngle: %.8f(%.5f)\n",
					frame, loop, linkBone.Name(), count-1,
					linkAngle, mmath.RadToDeg(linkAngle),
					originalLinkAngle, mmath.RadToDeg(originalLinkAngle))
			}

			// 回転軸の算出
			var axisLimited, axisOriginal *mmath.MVec3
			noAxisLimit := (!isSingleIk || linkAngle > mmath.GIMBAL1_RAD) && (ikLink.AngleLimit || ikLink.LocalAngleLimit)
			if noAxisLimit {
				// グローバル or ローカル角度制限
				if ikLink.AngleLimit {
					axisLimited, axisOriginal = getLinkAxis(
						ikLink.MinAngleLimit, ikLink.MaxAngleLimit,
						ikTargetLocalPos, ikLocalPos,
						frame, count, loop, linkBone.Name(), ikMotion, ikFile,
					)
				} else {
					axisLimited, axisOriginal = getLinkAxis(
						ikLink.LocalMinAngleLimit, ikLink.LocalMaxAngleLimit,
						ikTargetLocalPos, ikLocalPos,
						frame, count, loop, linkBone.Name(), ikMotion, ikFile,
					)
				}
			} else {
				axisLimited, axisOriginal = getLinkAxis(
					mmath.MVec3MinVal, mmath.MVec3MaxVal,
					ikTargetLocalPos, ikLocalPos,
					frame, count, loop, linkBone.Name(), ikMotion, ikFile,
				)
			}

			if isFullDebug {
				fmt.Fprintf(ikFile,
					"[%.3f][%03d][%s][%05d][回転軸] linkAxis: %s, originalLinkAxis: %s\n",
					frame, loop, linkBone.Name(), count-1,
					axisLimited.String(), axisOriginal.String())
			}

			// 回転量算出
			originalIkQuat := mmath.NewMQuaternionFromAxisAnglesRotate(axisOriginal, originalLinkAngle)
			var ikQuat *mmath.MQuaternion

			// Fixed軸を持つボーンか否か
			if linkBone.HasFixedAxis() {
				if !(ikLink.AngleLimit || ikLink.LocalAngleLimit) {
					// 角度制限なし => 軸に沿った理想回転
					tmpQuat := mmath.NewMQuaternionFromAxisAnglesRotate(axisLimited, linkAngle)
					ikQuat, _ = tmpQuat.SeparateTwistByAxis(linkBone.NormalizedFixedAxis)
				} else {
					// 角度制限あり => 軸に沿った回転をそのまま適用
					// (軸の向きが反転していれば角度を反転)
					if axisLimited.Dot(linkBone.NormalizedFixedAxis) < 0 {
						linkAngle = -linkAngle
					}
					ikQuat = mmath.NewMQuaternionFromAxisAnglesRotate(linkBone.NormalizedFixedAxis, linkAngle)
				}
			} else {
				ikQuat = mmath.NewMQuaternionFromAxisAnglesRotate(axisLimited, linkAngle)
			}

			originalTotalIkQuat := linkQuat.Muled(originalIkQuat)
			totalIkQuat := linkQuat.Muled(ikQuat)

			if isFullDebug {
				// originalTotalIkQuat
				bfOrig := vmd.NewBoneFrame(float32(count))
				bfOrig.Rotation = originalTotalIkQuat.Copy()
				ikMotion.AppendRegisteredBoneFrame(linkBone.Name(), bfOrig)
				count++
				fmt.Fprintf(ikFile,
					"[%.3f][%03d][%s][%05d][originalTotalIkQuat] %s(%s)\n",
					frame, loop, linkBone.Name(), count-1,
					originalTotalIkQuat.String(),
					originalTotalIkQuat.ToMMDDegrees().String())

				fmt.Fprintf(ikFile,
					"[%.3f][%03d][%s][%05d][originalIkQuat] %s(%s)\n",
					frame, loop, linkBone.Name(), count-1,
					originalIkQuat.String(), originalIkQuat.ToMMDDegrees().String())

				// totalIkQuat
				bfTotal := vmd.NewBoneFrame(float32(count))
				bfTotal.Rotation = totalIkQuat.Copy()
				ikMotion.AppendRegisteredBoneFrame(linkBone.Name(), bfTotal)
				count++
				fmt.Fprintf(ikFile,
					"[%.3f][%03d][%s][%05d][totalIkQuat] %s(%s)\n",
					frame, loop, linkBone.Name(), count-1,
					totalIkQuat.String(), totalIkQuat.ToMMDDegrees().String())

				fmt.Fprintf(ikFile,
					"[%.3f][%03d][%s][%05d][ikQuat] %s(%s)\n",
					frame, loop, linkBone.Name(), count-1,
					ikQuat.String(), ikQuat.ToMMDDegrees().String())
			}

			var resultIkQuat *mmath.MQuaternion
			// グローバル角度制限
			if ikLink.AngleLimit {
				resultIkQuat, count = calcIkLimitQuaternion(
					totalIkQuat,
					ikLink.MinAngleLimit,
					ikLink.MaxAngleLimit,
					mmath.MVec3UnitX,
					mmath.MVec3UnitY,
					mmath.MVec3UnitZ,
					loop, loopCount,
					frame, count,
					linkBone.Name(),
					ikMotion, ikFile)
			} else if ikLink.LocalAngleLimit {
				// ローカル角度制限
				resultIkQuat, count = calcIkLimitQuaternion(
					totalIkQuat,
					ikLink.LocalMinAngleLimit,
					ikLink.LocalMaxAngleLimit,
					linkBone.NormalizedLocalAxisX,
					linkBone.NormalizedLocalAxisY,
					linkBone.NormalizedLocalAxisZ,
					loop, loopCount,
					frame, count,
					linkBone.Name(),
					ikMotion, ikFile)
			} else {
				// 角度制限なし
				resultIkQuat = totalIkQuat
			}

			// モーフによる回転がある場合、最初のループ時に合成
			if loop == 0 && deltas.Morphs != nil && deltas.Morphs.Bones != nil {
				morphBoneDelta := deltas.Morphs.Bones.Get(linkBone.Index())
				if morphBoneDelta != nil && morphBoneDelta.FrameRotation != nil {
					resultIkQuat = resultIkQuat.Muled(morphBoneDelta.FrameRotation)
				}
			}

			// Fixed軸対応
			if linkBone.HasFixedAxis() {
				resultIkQuat = resultIkQuat.ToFixedAxisRotation(linkBone.NormalizedFixedAxis)
				if isFullDebug {
					bf := vmd.NewBoneFrame(float32(count))
					bf.Rotation = resultIkQuat.Copy()
					ikMotion.AppendRegisteredBoneFrame(linkBone.Name(), bf)
					count++
					fmt.Fprintf(ikFile,
						"[%.3f][%03d][%s][%05d][軸制限後] resultIkQuat: %s(%s)\n",
						frame, loop, linkBone.Name(), count-1,
						resultIkQuat.String(), resultIkQuat.ToMMDDegrees().String())
				}
			}

			// IK 結果反映
			linkDelta.FrameRotation = resultIkQuat
			updateBoneDelta(deltas.Bones, linkDelta, linkBone)
			updateGlobalMatrix(deltas.Bones, ikTargetDeformBoneIndexes)

			if isFullDebug {
				linkBf := vmd.NewBoneFrame(float32(count))
				linkBf.Rotation = linkDelta.FilledTotalRotation().Copy()
				ikMotion.AppendRegisteredBoneFrame(linkBone.Name(), linkBf)
				count++
				fmt.Fprintf(ikFile,
					"[%.3f][%03d][%s][%05d][結果] linkBf.Rotation: %s(%s)\n",
					frame, loop, linkBone.Name(), count-1,
					linkBf.Rotation.String(), linkBf.Rotation.ToMMDDegrees().String())
			}
		}
	}

	return deltas.Bones
}

func getLocalMat(
	deltas *delta.VmdDeltas,
	bone *pmx.Bone,
) (*mmath.MMat4, *mmath.MMat4) {
	var localMat *mmath.MMat4
	if deltas.Bones != nil && deltas.Bones.Get(bone.Index()) != nil && deltas.Bones.Get(bone.Index()).FrameLocalMat != nil {
		localMat = deltas.Bones.Get(bone.Index()).FrameLocalMat.Copy()
	}

	var morphLocalMat *mmath.MMat4
	if deltas.Morphs != nil && deltas.Morphs.Bones.Get(bone.Index()) != nil &&
		deltas.Morphs.Bones.Get(bone.Index()).FrameLocalMat != nil {
		morphLocalMat = deltas.Morphs.Bones.Get(bone.Index()).FrameLocalMat.Copy()
	}

	return localMat, morphLocalMat
}

// 該当キーフレにおけるボーンの移動位置
func getPosition(
	deltas *delta.VmdDeltas,
	bone *pmx.Bone,
	bf *vmd.BoneFrame,
) (*mmath.MVec3, *mmath.MVec3, *mmath.MVec3) {
	var cancelablePos *mmath.MVec3
	if bf != nil && bf.CancelablePosition != nil {
		cancelablePos = bf.CancelablePosition.Copy()
	} else {
		cancelablePos = mmath.NewMVec3()
	}

	var morphPos *mmath.MVec3
	if deltas.Morphs != nil && deltas.Morphs.Bones.Get(bone.Index()) != nil &&
		deltas.Morphs.Bones.Get(bone.Index()).FramePosition != nil {
		morphPos = deltas.Morphs.Bones.Get(bone.Index()).FramePosition.Copy()
	}

	var morphCancelablePos *mmath.MVec3
	if deltas.Morphs != nil && deltas.Morphs.Bones.Get(bone.Index()) != nil &&
		deltas.Morphs.Bones.Get(bone.Index()).FrameCancelablePosition != nil {
		morphCancelablePos = deltas.Morphs.Bones.Get(bone.Index()).FrameCancelablePosition.Copy()
	}

	return cancelablePos, morphPos, morphCancelablePos
}

// 該当キーフレにおけるボーンの回転角度
func getRotation(
	deltas *delta.VmdDeltas,
	bone *pmx.Bone,
	bf *vmd.BoneFrame,
) (*mmath.MQuaternion, *mmath.MQuaternion, *mmath.MQuaternion) {
	var cancelableRot *mmath.MQuaternion
	var morphRot *mmath.MQuaternion
	var morphCancelableRot *mmath.MQuaternion

	if bf != nil && bf.CancelableRotation != nil && !bf.CancelableRotation.IsIdent() {
		cancelableRot = bf.CancelableRotation.Copy()
	} else {
		cancelableRot = mmath.NewMQuaternion()
	}

	if deltas.Morphs != nil && deltas.Morphs.Bones.Get(bone.Index()) != nil &&
		deltas.Morphs.Bones.Get(bone.Index()).FrameRotation != nil {
		// IKの場合はIK計算時に組み込まれているので、まだframeRotationが無い場合のみ加味
		morphRot = deltas.Morphs.Bones.Get(bone.Index()).FrameRotation.Copy()
		// mlog.I("[%s][%.3f][%d]: rot: %s(%s), morphRot: %s(%s)\n", bone.Name(), frame, loop,
		// 	rot.String(), rot.ToMMDDegrees().String(), morphRot.String(), morphRot.ToMMDDegrees().String())
	}

	if deltas.Morphs != nil && deltas.Morphs.Bones.Get(bone.Index()) != nil &&
		deltas.Morphs.Bones.Get(bone.Index()).FrameCancelableRotation != nil {
		morphCancelableRot = deltas.Morphs.Bones.Get(bone.Index()).FrameCancelableRotation.Copy()
	}

	// if bone.HasFixedAxis() {
	// 	rot = rot.ToFixedAxisRotation(bone.NormalizedFixedAxis)
	// }

	return cancelableRot, morphRot, morphCancelableRot
}

// 該当キーフレにおけるボーンの拡大率
func getScale(
	deltas *delta.VmdDeltas,
	bone *pmx.Bone,
	bf *vmd.BoneFrame,
) (*mmath.MVec3, *mmath.MVec3, *mmath.MVec3, *mmath.MVec3) {

	var scale *mmath.MVec3
	if deltas.Bones != nil && deltas.Bones.Get(bone.Index()) != nil &&
		deltas.Bones.Get(bone.Index()).FrameScale != nil {
		scale = deltas.Bones.Get(bone.Index()).FrameScale
	} else if bf != nil && bf.Scale != nil && !bf.Scale.IsOne() {
		scale = bf.Scale
	}

	var cancelableScale *mmath.MVec3
	if deltas.Bones != nil && deltas.Bones.Get(bone.Index()) != nil &&
		deltas.Bones.Get(bone.Index()).FrameCancelableScale != nil {
		cancelableScale = deltas.Bones.Get(bone.Index()).FrameCancelableScale
	} else if bf != nil && bf.CancelableScale != nil && !bf.CancelableScale.IsOne() {
		cancelableScale = bf.CancelableScale
	}

	var morphScale *mmath.MVec3
	if deltas.Morphs != nil && deltas.Morphs.Bones.Get(bone.Index()) != nil &&
		deltas.Morphs.Bones.Get(bone.Index()).FrameScale != nil {
		morphScale = deltas.Morphs.Bones.Get(bone.Index()).FrameScale.Copy()
	}

	var morphCancelableScale *mmath.MVec3
	if deltas.Morphs != nil && deltas.Morphs.Bones.Get(bone.Index()) != nil &&
		deltas.Morphs.Bones.Get(bone.Index()).FrameCancelableScale != nil {
		morphCancelableScale = deltas.Morphs.Bones.Get(bone.Index()).FrameCancelableScale.Copy()
	}

	return scale, cancelableScale, morphScale, morphCancelableScale
}
