package deform

import (
	"fmt"
	"log"
	"math"
	"os"
	"path/filepath"
	"slices"
	"time"

	"github.com/miu200521358/mlib_go/pkg/domain/delta"
	"github.com/miu200521358/mlib_go/pkg/domain/mmath"
	"github.com/miu200521358/mlib_go/pkg/domain/pmx"
	"github.com/miu200521358/mlib_go/pkg/domain/vmd"
	"github.com/miu200521358/mlib_go/pkg/infrastructure/repository"
	"github.com/miu200521358/mlib_go/pkg/mutils"
	"github.com/miu200521358/mlib_go/pkg/mutils/mlog"
)

// DeformBoneByPhysicsFlag ボーンデフォーム処理を実行する
func DeformBoneByPhysicsFlag(
	model *pmx.PmxModel,
	motion *vmd.VmdMotion,
	deltas *delta.VmdDeltas,
	isCalcIk bool,
	frame float32,
	boneNames []string,
	isAfterPhysics bool,
) *delta.VmdDeltas {
	if model == nil || motion == nil {
		return deltas
	}

	deformBoneIndexes, deltas := newVmdDeltas(model, motion, deltas, frame, boneNames, isAfterPhysics)

	// ボーンデフォーム情報を埋める
	deltas.Bones = fillBoneDeform(model, motion, deltas, frame, deformBoneIndexes, isCalcIk, isAfterPhysics)

	// ボーンデフォーム情報を更新する
	updateGlobalMatrix(deltas.Bones, deformBoneIndexes)

	return deltas
}

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
	// 回転軸
	linkAxis := ikTargetLocalPosition.Cross(ikLocalPosition).Normalize()

	if mlog.IsIkVerbose() && ikMotion != nil && ikFile != nil {
		fmt.Fprintf(ikFile,
			"[%.3f][%03d][%s][%05d][linkAxis] %s\n",
			frame, loop, linkBoneName, count-1, linkAxis.MMD().String(),
		)
	}

	// linkMat := linkQuat.ToMat4()
	// if mlog.IsIkVerbose() && ikMotion != nil && ikFile != nil {
	// 	fmt.Fprintf(ikFile,
	// 		"[%.3f][%03d][%s][%05d][linkMat] %s (x: %s, y: %s, z: %s)\n",
	// 		frame, loop, linkBoneName, count-1, linkMat.String(), linkMat.AxisX().String(), linkMat.AxisY().String(), linkMat.AxisZ().String())
	// }

	if minAngleLimitRadians.IsOnlyX() || maxAngleLimitRadians.IsOnlyX() {
		// X軸のみの制限の場合
		vv := linkAxis.X

		if mlog.IsIkVerbose() && ikMotion != nil && ikFile != nil {
			fmt.Fprintf(ikFile,
				"[%.3f][%03d][%s][%05d][linkAxis(X軸制限)] vv: %.8f\n",
				frame, loop, linkBoneName, count-1, vv)
		}

		if vv < 0 {
			return mmath.MVec3UnitXInv, linkAxis
		}
		return mmath.MVec3UnitX, linkAxis
	} else if minAngleLimitRadians.IsOnlyY() || maxAngleLimitRadians.IsOnlyY() {
		// Y軸のみの制限の場合
		vv := linkAxis.Y

		if mlog.IsIkVerbose() && ikMotion != nil && ikFile != nil {
			fmt.Fprintf(ikFile,
				"[%.3f][%03d][%s][%05d][linkAxis(Y軸制限)] vv: %.8f\n",
				frame, loop, linkBoneName, count-1, vv)
		}

		if vv < 0 {
			return mmath.MVec3UnitYInv, linkAxis
		}
		return mmath.MVec3UnitY, linkAxis
	} else if minAngleLimitRadians.IsOnlyZ() || maxAngleLimitRadians.IsOnlyZ() {
		// Z軸のみの制限の場合
		vv := linkAxis.Z

		if mlog.IsIkVerbose() && ikMotion != nil && ikFile != nil {
			fmt.Fprintf(ikFile,
				"[%.3f][%03d][%s][%05d][linkAxis(Z軸制限)] vv: %.8f\n",
				frame, loop, linkBoneName, count-1, vv)
		}

		if vv < 0 {
			return mmath.MVec3UnitZInv, linkAxis
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
	ikMat := totalIkQuat.ToMat4()
	if mlog.IsIkVerbose() && ikMotion != nil && ikFile != nil {
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

		if mlog.IsIkVerbose() && ikMotion != nil && ikFile != nil {
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

			if mlog.IsIkVerbose() && ikMotion != nil && ikFile != nil {
				fmt.Fprintf(ikFile, "[%.3f][%03d][%s][%05d][X軸制限-ジンバル] fSX: %f, fX: %f, fCX: %f\n",
					frame, loop, linkBoneName, count-1, fSX, fX, fCX)
			}
		}

		// Y軸回り
		fSY := ikMat.AxisZ().X / fCX // sin(θy) = m31 / cos(θx)
		fCY := ikMat.AxisZ().Z / fCX // cos(θy) = m33 / cos(θx)
		fY := math.Atan2(fSY, fCY)   // Y軸回り決定

		if mlog.IsIkVerbose() && ikMotion != nil && ikFile != nil {
			fmt.Fprintf(ikFile, "[%.3f][%03d][%s][%05d][X軸制限-Y軸回り] fSY: %f, fCY: %f, fY: %f\n",
				frame, loop, linkBoneName, count-1, fSY, fCY, fY)
		}

		// Z軸周り
		fSZ := ikMat.AxisX().Y / fCX // sin(θz) = m12 / cos(θx)
		fCZ := ikMat.AxisY().Y / fCX // cos(θz) = m22 / cos(θx)
		fZ := math.Atan2(fSZ, fCZ)   // Z軸回り決定

		if mlog.IsIkVerbose() && ikMotion != nil && ikFile != nil {
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
		if mlog.IsIkVerbose() && ikMotion != nil && ikFile != nil {
			fmt.Fprintf(ikFile, "[%.3f][%03d][%s][%05d][X軸制限-xQuat] xAxisVector: %s, fX: %f, xQuat: %s(%s)\n",
				frame, loop, linkBoneName, count-1, xAxisVector.String(), fX, xQuat.String(), xQuat.ToMMDDegrees().String())
		}

		yQuat := mmath.NewMQuaternionFromAxisAnglesRotate(yAxisVector, fY)
		if mlog.IsIkVerbose() && ikMotion != nil && ikFile != nil {
			fmt.Fprintf(ikFile, "[%.3f][%03d][%s][%05d][X軸制限-yQuat] yAxisVector: %s, fY: %f, yQuat: %s(%s)\n",
				frame, loop, linkBoneName, count-1, yAxisVector.String(), fY, yQuat.String(), yQuat.ToMMDDegrees().String())
		}

		zQuat := mmath.NewMQuaternionFromAxisAnglesRotate(zAxisVector, fZ)
		if mlog.IsIkVerbose() && ikMotion != nil && ikFile != nil {
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

		if mlog.IsIkVerbose() && ikMotion != nil && ikFile != nil {
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

			if mlog.IsIkVerbose() && ikMotion != nil && ikFile != nil {
				fmt.Fprintf(ikFile, "[%.3f][%03d][%s][%05d][Y軸制限-ジンバル] fSY: %f, fY: %f, fCY: %f\n",
					frame, loop, linkBoneName, count-1, fSY, fY, fCY)
			}
		}

		// X軸回り
		fSX := ikMat.AxisY().Z / fCY // sin(θx) = m23 / cos(θy)
		fCX := ikMat.AxisZ().Z / fCY // cos(θx) = m33 / cos(θy)
		fX := math.Atan2(fSX, fCX)   // X軸回り決定

		if mlog.IsIkVerbose() && ikMotion != nil && ikFile != nil {
			fmt.Fprintf(ikFile, "[%.3f][%03d][%s][%05d][Y軸制限-X軸回り] fSX: %f, fCX: %f, fX: %f\n",
				frame, loop, linkBoneName, count-1, fSX, fCX, fX)
		}

		// Z軸周り
		fSZ := ikMat.AxisX().Y / fCY // sin(θz) = m12 / cos(θy)
		fCZ := ikMat.AxisX().X / fCY // cos(θz) = m11 / cos(θy)
		fZ := math.Atan2(fSZ, fCZ)   // Z軸回り決定

		if mlog.IsIkVerbose() && ikMotion != nil && ikFile != nil {
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
		if mlog.IsIkVerbose() && ikMotion != nil && ikFile != nil {
			fmt.Fprintf(ikFile, "[%.3f][%03d][%s][%05d][Y軸制限-xQuat] xAxisVector: %s, fX: %f, xQuat: %s(%s)\n",
				frame, loop, linkBoneName, count-1, xAxisVector.String(), fX, xQuat.String(), xQuat.ToMMDDegrees().String())
		}

		yQuat := mmath.NewMQuaternionFromAxisAnglesRotate(yAxisVector, fY)
		if mlog.IsIkVerbose() && ikMotion != nil && ikFile != nil {
			fmt.Fprintf(ikFile, "[%.3f][%03d][%s][%05d][Y軸制限-yQuat] yAxisVector: %s, fY: %f, yQuat: %s(%s)\n",
				frame, loop, linkBoneName, count-1, yAxisVector.String(), fY, yQuat.String(), yQuat.ToMMDDegrees().String())
		}

		zQuat := mmath.NewMQuaternionFromAxisAnglesRotate(zAxisVector, fZ)
		if mlog.IsIkVerbose() && ikMotion != nil && ikFile != nil {
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

	if mlog.IsIkVerbose() && ikMotion != nil && ikFile != nil {
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

		if mlog.IsIkVerbose() && ikMotion != nil && ikFile != nil {
			fmt.Fprintf(ikFile, "[%.3f][%03d][%s][%05d][Z軸制限-ジンバル] fSZ: %f, fZ: %f, fCZ: %f\n",
				frame, loop, linkBoneName, count-1, fSZ, fZ, fCZ)
		}
	}

	// X軸回り
	fSX := ikMat.AxisY().Z / fCZ // sin(θx) = m23 / cos(θz)
	fCX := ikMat.AxisY().Y / fCZ // cos(θx) = m22 / cos(θz)
	fX := math.Atan2(fSX, fCX)   // X軸回り決定

	if mlog.IsIkVerbose() && ikMotion != nil && ikFile != nil {
		fmt.Fprintf(ikFile, "[%.3f][%03d][%s][%05d][Z軸制限-X軸回り] fSX: %f, fCX: %f, fX: %f\n",
			frame, loop, linkBoneName, count-1, fSX, fCX, fX)
	}

	// Y軸周り
	fSY := ikMat.AxisZ().X / fCZ // sin(θy) = m31 / cos(θz)
	fCY := ikMat.AxisX().X / fCZ // cos(θy) = m11 / cos(θz)
	fY := math.Atan2(fSY, fCY)   // Y軸回り決定

	if mlog.IsIkVerbose() && ikMotion != nil && ikFile != nil {
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
	if mlog.IsIkVerbose() && ikMotion != nil && ikFile != nil {
		fmt.Fprintf(ikFile, "[%.3f][%03d][%s][%05d][Z軸制限-xQuat] xAxisVector: %s, fX: %f, xQuat: %s(%s)\n",
			frame, loop, linkBoneName, count-1, xAxisVector.String(), fX, xQuat.String(), xQuat.ToMMDDegrees().String())
	}

	yQuat := mmath.NewMQuaternionFromAxisAnglesRotate(yAxisVector, fY)
	if mlog.IsIkVerbose() && ikMotion != nil && ikFile != nil {
		fmt.Fprintf(ikFile, "[%.3f][%03d][%s][%05d][Z軸制限-yQuat] yAxisVector: %s, fY: %f, yQuat: %s(%s)\n",
			frame, loop, linkBoneName, count-1, yAxisVector.String(), fY, yQuat.String(), yQuat.ToMMDDegrees().String())
	}

	zQuat := mmath.NewMQuaternionFromAxisAnglesRotate(zAxisVector, fZ)
	if mlog.IsIkVerbose() && ikMotion != nil && ikFile != nil {
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
	isInLoop := float64(loop) < float64(loopCount)/2.0

	if mlog.IsIkVerbose() && ikMotion != nil && ikFile != nil {
		fmt.Fprintf(ikFile, "[%.3f][%03d][%s][%05d][%s-getIkAxisValue] loop: %d, isInLoop: %v\n",
			frame, loop, linkBoneName, count-1, axisName, loop, isInLoop)
	}

	if fV < minAngleLimit {
		tf := 2*minAngleLimit - fV
		if tf <= maxAngleLimit && isInLoop {
			if mlog.IsIkVerbose() && ikMotion != nil && ikFile != nil {
				fmt.Fprintf(ikFile, "[%.3f][%03d][%s][%05d][%s-最小角度(loop内)] minAngleLimit: %f, fV: %f, tf: %f\n",
					frame, loop, linkBoneName, count-1, axisName, minAngleLimit, fV, tf)
			}

			fV = tf
		} else {
			if mlog.IsIkVerbose() && ikMotion != nil && ikFile != nil {
				fmt.Fprintf(ikFile, "[%.3f][%03d][%s][%05d][%s-最小角度(loop外)] minAngleLimit: %f, fV: %f, tf: %f\n",
					frame, loop, linkBoneName, count-1, axisName, minAngleLimit, fV, tf)
			}

			fV = minAngleLimit
		}
	}

	if fV > maxAngleLimit {
		tf := 2*maxAngleLimit - fV
		if tf >= minAngleLimit && isInLoop {
			if mlog.IsIkVerbose() && ikMotion != nil && ikFile != nil {
				fmt.Fprintf(ikFile, "[%.3f][%03d][%s][%05d][%s-最大角度(loop内)] maxAngleLimit: %f, fV: %f, tf: %f\n",
					frame, loop, linkBoneName, count-1, axisName, maxAngleLimit, fV, tf)
			}

			fV = tf
		} else {
			if mlog.IsIkVerbose() && ikMotion != nil && ikFile != nil {
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
		bone := model.Bones.GetByName(boneNames[0])
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
			bone := model.Bones.GetByName(boneName)

			// 対象のボーンは常に追加
			if _, ok := relativeBoneIndexes[bone.Index()]; !ok {
				relativeBoneIndexes[bone.Index()] = struct{}{}
			}

			// 関連するボーンの追加
			for _, index := range bone.Extend.RelativeBoneIndexes {
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
		bone := model.Bones.Get(boneIndex)
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
	// IKのON/OFF
	ikFrame := motion.IkFrames.Get(frame)
	calcIkRelativeBoneIndexes := make([]int, 0, len(deformBoneIndexes))

	for i, boneIndex := range deformBoneIndexes {
		bone := model.Bones.Get(boneIndex)
		d := deltas.Bones.Get(boneIndex)
		if d == nil {
			d = &delta.BoneDelta{Bone: bone, Frame: frame}
		} else if slices.Contains(calcIkRelativeBoneIndexes, boneIndex) {
			// IK計算済みボーンはスキップ
			continue
		}

		bf := motion.BoneFrames.Get(bone.Name()).Get(frame)
		if bf.Position != nil {
			d.FramePosition = bf.Position.Copy()
		}
		if bf.Rotation != nil {
			d.FrameRotation = bf.Rotation.Copy()
		}

		// ボーンの移動位置、回転角度、拡大率を取得
		d.FrameLocalMat, d.FrameLocalMorphMat = getLocalMat(deltas, bone)
		d.FrameCancelablePosition, d.FrameMorphPosition, d.FrameMorphCancelablePosition =
			getPosition(deltas, bone, bf)
		d.FrameCancelableRotation, d.FrameMorphRotation, d.FrameMorphCancelableRotation =
			getRotation(deltas, bone, bf)
		d.FrameScale, d.FrameCancelableScale, d.FrameMorphScale, d.FrameMorphCancelableScale =
			getScale(deltas, bone, bf)

		updateBoneDelta(deltas.Bones, d, bone)

		if isCalcIk && bone.IsIK() && ikFrame.IsEnable(bone.Name()) {
			// IKの変形リスト
			ikTargetDeformBoneIndexes := model.Bones.DeformBoneIndexes[bone.Index()]

			deltas.Bones = fillBoneDeform(model, motion, deltas, frame, ikTargetDeformBoneIndexes, false, isAfterPhysics)
			updateGlobalMatrix(deltas.Bones, ikTargetDeformBoneIndexes)

			for _, boneIndex := range ikTargetDeformBoneIndexes {
				d := deltas.Bones.Get(boneIndex)
				if d == nil {
					continue
				}
				d.GlobalIkOffMatrix = d.GlobalMatrix.Copy()
				deltas.Bones.Update(d)
			}

			// IKボーンのグローバル位置
			ikGlobalPosition := deltas.Bones.Get(bone.Index()).FilledGlobalPosition()
			deformIk(model, motion, deltas, frame, isAfterPhysics, bone, ikGlobalPosition, ikTargetDeformBoneIndexes, i)
			for _, l := range bone.Ik.Links {
				calcIkRelativeBoneIndexes = append(calcIkRelativeBoneIndexes, l.BoneIndex)
			}
			calcIkRelativeBoneIndexes = append(calcIkRelativeBoneIndexes, boneIndex)
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
	localMat := boneDeltas.TotalLocalMat(bone.Index())
	if localMat != nil && !localMat.IsIdent() {
		d.UnitMatrix.Mul(localMat)
	}

	// 移動
	posMat := boneDeltas.TotalPositionMat(bone.Index())
	if posMat != nil && !posMat.IsIdent() {
		d.UnitMatrix.Mul(posMat)
	}

	// 回転
	rotMat := boneDeltas.TotalRotationMat(bone.Index())
	if rotMat != nil && !rotMat.IsIdent() {
		d.UnitMatrix.Mul(rotMat)
	}

	// スケール
	scaleMat := boneDeltas.TotalScaleMat(bone.Index())
	if scaleMat != nil && !scaleMat.IsIdent() {
		d.UnitMatrix.Mul(scaleMat)
	}

	// 逆BOf行列(初期姿勢行列)
	d.UnitMatrix = d.Bone.Extend.RevertOffsetMatrix.Muled(d.UnitMatrix)

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

// IK計算
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
	if len(ikBone.Ik.Links) < 1 {
		// IKリンクが無ければスルー
		return deltas.Bones
	}

	var prefixPath string
	if mlog.IsIkVerbose() {
		// IK計算デバッグ用モーション
		dirPath := fmt.Sprintf("%s/IK_step", filepath.Dir(model.Path()))
		err := os.MkdirAll(dirPath, 0755)
		if err != nil {
			log.Fatal(err)
		}

		_, motionFileName, _ := mutils.SplitPath(motion.Path())
		date := time.Now().Format("20060102_150405")
		prefixPath = fmt.Sprintf("%s/%s_%.2f_%s_%03d", dirPath, motionFileName, frame, date, deformIndex)
	}

	var err error
	var ikFile *os.File
	var ikMotion *vmd.VmdMotion
	var globalMotion *vmd.VmdMotion
	count := 1

	if mlog.IsIkVerbose() {
		// IK計算デバッグ用モーション
		ikMotionPath := fmt.Sprintf("%s_%s.vmd", prefixPath, ikBone.Name())
		ikMotion = vmd.NewVmdMotion(ikMotionPath)

		globalMotionPath := fmt.Sprintf("%s_%s_global.vmd", prefixPath, ikBone.Name())
		globalMotion = vmd.NewVmdMotion(globalMotionPath)

		ikLogPath := fmt.Sprintf("%s_%s.log", prefixPath, ikBone.Name())
		ikFile, err = os.OpenFile(ikLogPath, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		fmt.Println(ikFile, "----------------------------------------")
		fmt.Println(ikFile, "[IK計算出力先][%.3f][%s] %s", frame, ikMotionPath)
	}
	defer func() {
		mlog.IV("[IK計算終了][%.3f][%s]", frame, ikBone.Name())

		if ikMotion != nil {
			r := repository.NewVmdRepository()
			if err := r.Save("", ikMotion, true); err != nil {
				mlog.E("[IK計算出力失敗][%.3f][%s] %s", frame, ikBone.Name(), err)
			}
		}
		if globalMotion != nil {
			r := repository.NewVmdRepository()
			if err := r.Save("", globalMotion, true); err != nil {
				mlog.E("[IK計算出力失敗][%.3f][%s] %s", frame, ikBone.Name(), err)
			}
		}
		if ikFile != nil {
			ikFile.Close()
		}
	}()

	// // つま先ＩＫであるか
	// isToeIk := strings.Contains(ikBone.Name(), "つま先ＩＫ")
	// 一段IKであるか
	isSingleIk := len(ikBone.Ik.Links) == 1

	// ループ回数
	loopCount := max(ikBone.Ik.LoopCount, 1)

	// IKターゲットボーン
	ikTargetBone := model.Bones.Get(ikBone.Ik.BoneIndex)

	if mlog.IsIkVerbose() && ikMotion != nil && ikFile != nil {
		ikOffMotion := vmd.NewVmdMotion(fmt.Sprintf("%s_0_%s.vmd", prefixPath, ikBone.Name()))

		bif := vmd.NewIkFrame(0)
		bif.Registered = true

		for _, bone := range model.Bones.Data {
			if bone.IsIK() {
				ef := vmd.NewIkEnableFrame(0)
				ef.Registered = true
				ef.BoneName = bone.Name()
				ef.Enabled = false

				bif.IkList = append(bif.IkList, ef)
			}
		}

		ikOffMotion.AppendIkFrame(bif)
		boneNames := make(map[string]struct{})

		for _, delta := range deltas.Bones.Data {
			if delta == nil || (delta.FramePosition == nil && delta.FrameRotation == nil) {
				continue
			}
			bf := vmd.NewBoneFrame(0)
			bf.Position = delta.FramePosition
			bf.Rotation = delta.FrameRotation
			ikOffMotion.AppendRegisteredBoneFrame(delta.Bone.Name(), bf)
			boneNames[delta.Bone.Name()] = struct{}{}
		}

		for _, bone := range model.Bones.Data {
			if _, ok := boneNames[bone.Name()]; !ok {
				bf := motion.BoneFrames.Get(bone.Name()).Get(frame)
				ikOffMotion.AppendRegisteredBoneFrame(bone.Name(), bf)
			}
		}

		r := repository.NewVmdRepository()
		if err := r.Save("", ikOffMotion, true); err != nil {
			mlog.E("[IK計算出力失敗][%.3f][%s] %s", frame, ikBone.Name(), err)
		}
	}

	if mlog.IsIkVerbose() && ikMotion != nil && ikFile != nil {
		bf := vmd.NewBoneFrame(float32(count))
		bf.Position = deltas.Bones.Get(ikBone.Index()).FilledFramePosition()
		bf.Rotation = deltas.Bones.Get(ikBone.Index()).FilledTotalRotation()
		ikMotion.AppendRegisteredBoneFrame(ikBone.Name(), bf)
		count++

		// fmt.Fprintf(ikFile, "[%.3f][%03d][%s][%05d][Local] ikGlobalPosition: %s\n",
		// 	frame, loop, linkBone.Name(), count-1, bf.Position.MMD().String())
	}

	ikTargetDeformIndex := slices.Index(ikTargetDeformBoneIndexes, ikTargetBone.Index())
	ikDeformIndex := slices.Index(ikTargetDeformBoneIndexes, ikBone.Index())
	if ikTargetDeformIndex < ikDeformIndex {
		// 初回IK OFF向きの場合、初回に足首位置に向けるのでループを余分に回す
		loopCount += 1
	}

	// 初回にIK事前計算
	var ikOnGlobalPosition *mmath.MVec3
	if ikTargetDeformIndex < ikDeformIndex && isSingleIk && ikBone.IsLegIK() {
		// デフォームの順番がIKの方がターゲットより後の場合、初回はIK OFF向きにする
		ikOffDeltas := DeformBoneByPhysicsFlag(model, motion, nil, false, frame,
			[]string{ikTargetBone.Name()}, isAfterPhysics)

		// IK ON時の IKターゲットボーンのグローバル位置を一旦別保持
		ikOnGlobalPosition = ikGlobalPosition.Copy()

		// IK OFF 時の IKターゲットボーンのグローバル位置を取得
		ikGlobalPosition = ikOffDeltas.Bones.Get(ikTargetBone.Index()).FilledGlobalPosition()

		if mlog.IsIkVerbose() && globalMotion != nil && ikFile != nil {
			{
				bf := vmd.NewBoneFrame(float32(count))
				bf.Position = ikGlobalPosition
				globalMotion.AppendRegisteredBoneFrame(ikBone.Name(), bf)
			}
			count++
		}
	}

	// IK計算
ikLoop:
	for loop := 0; loop < loopCount; loop++ {
		for lidx, ikLink := range ikBone.Ik.Links {
			// ikLink は末端から並んでる
			if !model.Bones.Contains(ikLink.BoneIndex) {
				continue
			}

			// 処理対象IKリンクボーン
			linkBone := model.Bones.Get(ikLink.BoneIndex)

			// 角度制限があってまったく動かさない場合、IK計算しないで次に行く
			if (linkBone.Extend.AngleLimit &&
				linkBone.Extend.MinAngleLimit.Radians().IsZero() &&
				linkBone.Extend.MaxAngleLimit.Radians().IsZero()) ||
				(linkBone.Extend.LocalAngleLimit &&
					linkBone.Extend.LocalMinAngleLimit.Radians().IsZero() &&
					linkBone.Extend.LocalMaxAngleLimit.Radians().IsZero()) {
				continue
			}

			if mlog.IsIkVerbose() && ikMotion != nil && ikFile != nil {
				fmt.Fprintf(ikFile,
					"[%.3f][%03d][%s][%05d] -------------------------------------------- \n",
					frame, loop, linkBone.Name(), count-1)
			}

			// リンクボーンの変形情報を取得
			linkDelta := deltas.Bones.Get(linkBone.Index())
			if linkDelta == nil {
				linkDelta = &delta.BoneDelta{Bone: linkBone, Frame: frame}
			}
			linkQuat := linkDelta.FilledTotalRotation()

			if mlog.IsIkVerbose() && ikMotion != nil && ikFile != nil {
				bf := vmd.NewBoneFrame(float32(count))
				bf.Rotation = linkQuat.Copy()
				ikMotion.AppendRegisteredBoneFrame(linkBone.Name(), bf)
				count++

				fmt.Fprintf(ikFile,
					"[%.3f][%03d][%s][%05d][linkQuat] %s(%s)\n",
					frame, loop, linkBone.Name(), count-1, bf.Rotation.String(), bf.Rotation.ToMMDDegrees().String(),
				)
			}

			// 現在のIKターゲットボーンのグローバル位置を取得
			ikTargetGlobalPosition := deltas.Bones.Get(ikTargetBone.Index()).FilledGlobalPosition()

			// 注目ノード（実際に動かすボーン=リンクボーン）
			// ワールド座標系から注目ノードの局所座標系への変換
			linkInvMatrix := deltas.Bones.Get(linkBone.Index()).FilledGlobalMatrix().Inverted()

			if mlog.IsIkVerbose() && globalMotion != nil && ikFile != nil {
				{
					bf := vmd.NewBoneFrame(float32(count))
					bf.Position = ikGlobalPosition
					globalMotion.AppendRegisteredBoneFrame(ikBone.Name(), bf)
				}
				{
					bf := vmd.NewBoneFrame(float32(count))
					bf.Position = ikTargetGlobalPosition
					globalMotion.AppendRegisteredBoneFrame(ikTargetBone.Name(), bf)
				}
				count++
			}

			// つま先IK(足首INDEXがつま先IKより前に計算される場合)とかは初回ループを超えたらIKグローバル位置を再計算
			if loop == 1 && lidx == 0 && ikTargetDeformIndex < ikDeformIndex && ikOnGlobalPosition != nil {
				ikGlobalPosition = ikOnGlobalPosition

				if mlog.IsIkVerbose() && globalMotion != nil && ikFile != nil {
					{
						bf := vmd.NewBoneFrame(float32(count))
						bf.Position = deltas.Bones.Get(linkBone.Index()).FilledGlobalPosition()
						globalMotion.AppendRegisteredBoneFrame(linkBone.Name(), bf)
					}
					{
						bf := vmd.NewBoneFrame(float32(count))
						bf.Position = ikGlobalPosition
						globalMotion.AppendRegisteredBoneFrame(ikBone.Name(), bf)
					}
					{
						bf := vmd.NewBoneFrame(float32(count))
						bf.Position = ikTargetGlobalPosition
						globalMotion.AppendRegisteredBoneFrame(ikTargetBone.Name(), bf)
					}
					count++
				}
			}

			if mlog.IsIkVerbose() && ikMotion != nil && ikFile != nil {
				fmt.Fprintf(ikFile,
					"[%.3f][%03d][%s][%05d][Global] [%s]ikGlobalPosition: %s, "+
						"[%s]ikTargetGlobalPosition: %s, [%s]linkGlobalPosition: %s\n",
					frame, loop, linkBone.Name(), count-1,
					ikBone.Name(), ikGlobalPosition.MMD().String(),
					ikTargetBone.Name(), ikTargetGlobalPosition.MMD().String(),
					linkBone.Name(), deltas.Bones.Get(linkBone.Index()).FilledGlobalPosition().MMD().String())
			}

			// 注目ノードを起点とした、IKターゲットのローカル位置
			ikTargetLocalPosition := linkInvMatrix.MulVec3(ikTargetGlobalPosition)
			// 注目ノードを起点とした、IK目標のローカル位置
			ikLocalPosition := linkInvMatrix.MulVec3(ikGlobalPosition)

			if mlog.IsIkVerbose() && ikMotion != nil && ikFile != nil {
				fmt.Fprintf(ikFile,
					"[%.3f][%03d][%s][%05d][Local] ikTargetLocalPosition: %s, ikLocalPosition: %s (%f)\n",
					frame, loop, linkBone.Name(), count-1,
					ikTargetLocalPosition.MMD().String(), ikLocalPosition.MMD().String(),
					ikTargetLocalPosition.Distance(ikLocalPosition))
			}

			ikTargetLocalPosition.Normalize()
			ikLocalPosition.Normalize()

			distanceThreshold := ikTargetLocalPosition.Distance(ikLocalPosition)
			if distanceThreshold < 1e-5 {
				if mlog.IsIkVerbose() && ikMotion != nil && ikFile != nil {
					fmt.Fprintf(ikFile,
						"[%.3f][%03d][%s][%05d][Local] ***BREAK*** distanceThreshold: %f\n",
						frame, loop, linkBone.Name(), count-1, distanceThreshold)
				}

				break ikLoop
			}

			if mlog.IsIkVerbose() && ikMotion != nil && ikFile != nil {
				fmt.Fprintf(ikFile,
					"[%.3f][%03d][%s][%05d][Local] ikTargetLocalPositionNorm: %s, ikLocalPositionNorm: %s\n",
					frame, loop, linkBone.Name(), count-1,
					ikTargetLocalPosition.MMD().String(), ikLocalPosition.MMD().String())
			}

			// 単位角
			unitRad := ikBone.Ik.UnitRotation.Radians().X * float64(lidx+1)
			linkDot := ikLocalPosition.Dot(ikTargetLocalPosition)

			if mlog.IsIkVerbose() && ikMotion != nil && ikFile != nil {
				fmt.Fprintf(ikFile,
					"[%.3f][%03d][%s][%05d][回転角度] unitRad: %.8f (%.5f), linkDot: %.8f\n",
					frame, loop, linkBone.Name(), count-1, unitRad, mmath.ToDegree(unitRad), linkDot,
				)
			}

			// 回転角(ラジアン)
			// 単位角を超えないようにする
			originalLinkAngle := math.Acos(mmath.ClampedFloat(linkDot, -1, 1))
			linkAngle := mmath.ClampedFloat(originalLinkAngle, -unitRad, unitRad)

			if mlog.IsIkVerbose() && ikMotion != nil && ikFile != nil {
				fmt.Fprintf(ikFile,
					"[%.3f][%03d][%s][%05d][単位角制限] linkAngle: %.8f(%.5f), originalLinkAngle: %.8f(%.5f)\n",
					frame, loop, linkBone.Name(), count-1, linkAngle, mmath.ToDegree(linkAngle),
					originalLinkAngle, mmath.ToDegree(originalLinkAngle),
				)
			}

			// 回転軸
			var originalLinkAxis, linkAxis *mmath.MVec3
			// 一段IKでない場合、または一段IKでかつ回転角が88度以上の場合
			if (!isSingleIk || (isSingleIk && linkAngle > mmath.GIMBAL1_RAD)) && ikLink.AngleLimit {
				// グローバル軸制限
				linkAxis, originalLinkAxis = getLinkAxis(
					ikLink.MinAngleLimit.Radians(),
					ikLink.MaxAngleLimit.Radians(),
					ikTargetLocalPosition, ikLocalPosition,
					frame, count, loop, linkBone.Name(), ikMotion, ikFile,
				)
			} else if (!isSingleIk || (isSingleIk && linkAngle > mmath.GIMBAL1_RAD)) && ikLink.LocalAngleLimit {
				// ローカル軸制限
				linkAxis, originalLinkAxis = getLinkAxis(
					ikLink.LocalMinAngleLimit.Radians(),
					ikLink.LocalMaxAngleLimit.Radians(),
					ikTargetLocalPosition, ikLocalPosition,
					frame, count, loop, linkBone.Name(), ikMotion, ikFile,
				)
			} else {
				// 軸制限なし or 一段IKでかつ回転角が88度未満の場合
				linkAxis, originalLinkAxis = getLinkAxis(
					mmath.MVec3MinVal,
					mmath.MVec3MaxVal,
					ikTargetLocalPosition, ikLocalPosition,
					frame, count, loop, linkBone.Name(), ikMotion, ikFile,
				)
			}

			if mlog.IsIkVerbose() && ikMotion != nil && ikFile != nil {
				fmt.Fprintf(ikFile,
					"[%.3f][%03d][%s][%05d][回転軸] linkAxis: %s, originalLinkAxis: %s\n",
					frame, loop, linkBone.Name(), count-1, linkAxis.String(), originalLinkAxis.String(),
				)
			}

			var originalIkQuat, ikQuat, linkIkQuat *mmath.MQuaternion

			if linkBone.HasFixedAxis() {
				originalIkQuat = mmath.NewMQuaternionFromAxisAnglesRotate(originalLinkAxis, originalLinkAngle)

				if !(ikLink.AngleLimit || ikLink.LocalAngleLimit) {
					// 軸制限あり＆角度制限なしの場合、軸にそった理想回転量とする
					linkIkQuat = mmath.NewMQuaternionFromAxisAnglesRotate(linkAxis, linkAngle)
					ikQuat, _ = linkIkQuat.SeparateTwistByAxis(linkBone.Extend.NormalizedFixedAxis)
				} else {
					if linkAxis.Dot(linkBone.Extend.NormalizedFixedAxis) < 0 {
						linkAngle = -linkAngle
					}
					// 軸制限あり＆角度制限ありの場合、calcIkLimitQuaternionで処理するのでこっちはそのまま
					ikQuat = mmath.NewMQuaternionFromAxisAnglesRotate(linkBone.Extend.NormalizedFixedAxis, linkAngle)
				}
			} else {
				originalIkQuat = mmath.NewMQuaternionFromAxisAnglesRotate(originalLinkAxis, originalLinkAngle)
				ikQuat = mmath.NewMQuaternionFromAxisAnglesRotate(linkAxis, linkAngle)
			}

			originalTotalIkQuat := linkQuat.Muled(originalIkQuat)
			totalIkQuat := linkQuat.Muled(ikQuat)

			if mlog.IsIkVerbose() && ikMotion != nil && ikFile != nil {
				{
					bf := vmd.NewBoneFrame(float32(count))
					bf.Rotation = originalTotalIkQuat.Copy()
					ikMotion.AppendRegisteredBoneFrame(linkBone.Name(), bf)
					count++

					fmt.Fprintf(ikFile,
						"[%.3f][%03d][%s][%05d][originalTotalIkQuat] %s(%s)\n",
						frame, loop, linkBone.Name(), count-1, originalTotalIkQuat.String(), originalTotalIkQuat.ToMMDDegrees().String())

					fmt.Fprintf(ikFile,
						"[%.3f][%03d][%s][%05d][originalIkQuat] %s(%s)\n",
						frame, loop, linkBone.Name(), count-1, originalIkQuat.String(), originalIkQuat.ToMMDDegrees().String())
				}
				{
					bf := vmd.NewBoneFrame(float32(count))
					bf.Rotation = totalIkQuat.Copy()
					ikMotion.AppendRegisteredBoneFrame(linkBone.Name(), bf)
					count++

					fmt.Fprintf(ikFile,
						"[%.3f][%03d][%s][%05d][totalIkQuat] %s(%s)\n",
						frame, loop, linkBone.Name(), count-1, totalIkQuat.String(), totalIkQuat.ToMMDDegrees().String())

					fmt.Fprintf(ikFile,
						"[%.3f][%03d][%s][%05d][ikQuat] %s(%s)\n",
						frame, loop, linkBone.Name(), count-1, ikQuat.String(), ikQuat.ToMMDDegrees().String())
				}
			}

			var resultIkQuat *mmath.MQuaternion
			if ikLink.AngleLimit {
				// 角度制限が入ってる場合
				resultIkQuat, count = calcIkLimitQuaternion(
					totalIkQuat,
					ikLink.MinAngleLimit.Radians(),
					ikLink.MaxAngleLimit.Radians(),
					mmath.MVec3UnitX, mmath.MVec3UnitY, mmath.MVec3UnitZ,
					loop, loopCount,
					frame, count, linkBone.Name(), ikMotion, ikFile,
				)

				if mlog.IsIkVerbose() && ikMotion != nil && ikFile != nil {
					bf := vmd.NewBoneFrame(float32(count))
					bf.Rotation = resultIkQuat.Copy()
					ikMotion.AppendRegisteredBoneFrame(linkBone.Name(), bf)
					count++

					fmt.Fprintf(ikFile,
						"[%.3f][%03d][%s][%05d][角度制限後] resultIkQuat: %s(%s), totalIkQuat: %s(%s), ikQuat: %s(%s)\n",
						frame, loop, linkBone.Name(), count-1, resultIkQuat.String(), resultIkQuat.ToMMDDegrees().String(),
						totalIkQuat.String(), totalIkQuat.ToMMDDegrees().String(),
						ikQuat.String(), ikQuat.ToMMDDegrees().String())
				}
			} else if ikLink.LocalAngleLimit {
				// ローカル角度制限が入ってる場合
				resultIkQuat, count = calcIkLimitQuaternion(
					totalIkQuat,
					ikLink.LocalMinAngleLimit.Radians(),
					ikLink.LocalMaxAngleLimit.Radians(),
					linkBone.Extend.NormalizedLocalAxisX,
					linkBone.Extend.NormalizedLocalAxisY,
					linkBone.Extend.NormalizedLocalAxisZ,
					loop, loopCount,
					frame, count, linkBone.Name(), ikMotion, ikFile,
				)

				if mlog.IsIkVerbose() && ikMotion != nil && ikFile != nil {
					bf := vmd.NewBoneFrame(float32(count))
					bf.Rotation = resultIkQuat.Copy()
					ikMotion.AppendRegisteredBoneFrame(linkBone.Name(), bf)
					count++

					fmt.Fprintf(ikFile,
						"[%.3f][%03d][%s][%05d][ローカル角度制限後] resultIkQuat: %s(%s), totalIkQuat: %s(%s), ikQuat: %s(%s)\n",
						frame, loop, linkBone.Name(), count-1, resultIkQuat.String(), resultIkQuat.ToMMDDegrees().String(),
						totalIkQuat.String(), totalIkQuat.ToMMDDegrees().String(),
						ikQuat.String(), ikQuat.ToMMDDegrees().String())
				}
			} else {
				// 角度制限なしの場合
				resultIkQuat = totalIkQuat
			}

			if loop == 0 && deltas.Morphs != nil && deltas.Morphs.Bones != nil &&
				deltas.Morphs.Bones.Get(linkBone.Index()) != nil &&
				deltas.Morphs.Bones.Get(linkBone.Index()).FrameRotation != nil {
				// モーフ変形がある場合、モーフ変形を追加適用
				resultIkQuat = resultIkQuat.Muled(deltas.Morphs.Bones.Get(linkBone.Index()).FrameRotation)
			}

			if linkBone.HasFixedAxis() {
				// 軸制限ありの場合、軸にそった理想回転量とする
				resultIkQuat = resultIkQuat.ToFixedAxisRotation(linkBone.Extend.NormalizedFixedAxis)

				if mlog.IsIkVerbose() && ikMotion != nil && ikFile != nil {
					bf := vmd.NewBoneFrame(float32(count))
					bf.Rotation = resultIkQuat.Copy()
					ikMotion.AppendRegisteredBoneFrame(linkBone.Name(), bf)
					count++

					fmt.Fprintf(ikFile,
						"[%.3f][%03d][%s][%05d][軸制限後] resultIkQuat: %s(%s)\n",
						frame, loop, linkBone.Name(), count-1, resultIkQuat.String(), resultIkQuat.ToMMDDegrees().String())
				}
			}

			// IKの結果を更新
			linkDelta.FrameRotation = resultIkQuat

			updateBoneDelta(deltas.Bones, linkDelta, linkBone)
			updateGlobalMatrix(deltas.Bones, ikTargetDeformBoneIndexes)

			if mlog.IsIkVerbose() && ikMotion != nil && ikFile != nil {
				linkBf := vmd.NewBoneFrame(float32(count))
				linkBf.Rotation = linkDelta.FilledTotalRotation().Copy()
				ikMotion.AppendRegisteredBoneFrame(linkBone.Name(), linkBf)
				count++

				fmt.Fprintf(ikFile,
					"[%.3f][%03d][%s][%05d][結果] linkBf.Rotation: %s(%s)\n",
					frame, loop, linkBone.Name(), count-1, linkBf.Rotation.String(), linkBf.Rotation.ToMMDDegrees().String())
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
	// 	rot = rot.ToFixedAxisRotation(bone.Extend.NormalizedFixedAxis)
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
		cancelableScale = deltas.Bones.Get(bone.Index()).FrameScale
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
