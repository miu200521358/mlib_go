package vmd

import (
	"slices"

	"github.com/miu200521358/mlib_go/pkg/domain/mmath"
)

type BoneNameFrames struct {
	*BaseFrames[*BoneFrame]
	Name string // ボーン名
}

func NewBoneNameFrames(name string) *BoneNameFrames {
	return &BoneNameFrames{
		BaseFrames: NewBaseFrames[*BoneFrame](NewBoneFrame, NullBoneFrame),
		Name:       name,
	}
}

func (boneNameFrames *BoneNameFrames) Copy() *BoneNameFrames {
	copied := NewBoneNameFrames(boneNameFrames.Name)
	for _, frame := range boneNameFrames.List() {
		copied.Append(frame.Copy().(*BoneFrame))
	}
	return copied
}

func (boneNameFrames *BoneNameFrames) Reduce() *BoneNameFrames {
	maxFrame := boneNameFrames.Indexes.Max()
	maxIFrame := int(maxFrame) + 1

	frames := make([]float32, 0, maxIFrame)
	xs := make([]float64, 0, maxIFrame)
	ys := make([]float64, 0, maxIFrame)
	zs := make([]float64, 0, maxIFrame)
	fixRs := make([]float64, 0, maxIFrame)
	quats := make([]*mmath.MQuaternion, 0, maxIFrame)

	for iF := range maxIFrame {
		f := float32(iF)

		frames = append(frames, f)
		bf := boneNameFrames.Get(f)

		if bf.Position != nil {
			xs = append(xs, bf.Position.X)
			ys = append(ys, bf.Position.Y)
			zs = append(zs, bf.Position.Z)
		} else {
			xs = append(xs, 0)
			ys = append(ys, 0)
			zs = append(zs, 0)
		}

		if bf.Rotation != nil {
			quats = append(quats, bf.Rotation)
			fixRs = append(fixRs, mmath.MQuaternionIdent.Dot(bf.Rotation))
		} else {
			quats = append(quats, mmath.MQuaternionIdent)
			fixRs = append(fixRs, 0)
		}
	}

	inflectionFrames := make([]float32, 0, boneNameFrames.Len())
	if !mmath.IsAllSameValues(xs) {
		inflectionFrames = append(inflectionFrames, mmath.FindInflectionFrames(frames, xs, 1e-4)...)
	}
	if !mmath.IsAllSameValues(ys) {
		inflectionFrames = append(inflectionFrames, mmath.FindInflectionFrames(frames, ys, 1e-4)...)
	}
	if !mmath.IsAllSameValues(zs) {
		inflectionFrames = append(inflectionFrames, mmath.FindInflectionFrames(frames, zs, 1e-4)...)
	}
	if !mmath.IsAllSameValues(fixRs) {
		inflectionFrames = append(inflectionFrames, mmath.FindInflectionFrames(frames, fixRs, 1e-6)...)

		// inflectionFrames = mmath.UniqueFloat32s(inflectionFrames)
		// mmath.SortFloat32s(inflectionFrames)

		// inflectionQuatFrames := make([]float32, 0, len(frames))
		// j := 2
		// for {
		// 	if j >= len(inflectionFrames)-1 {
		// 		break
		// 	}

		// 	startFrame := inflectionFrames[j-2]
		// 	endFrame := inflectionFrames[j]

		// 	startIFrame := int(startFrame)
		// 	endIFrame := int(endFrame)

		// 	startQuat := quats[startIFrame]
		// 	// endQuat := quats[endIFrame]

		// 	rangeFrames := make([]float32, 0, endIFrame-startIFrame+1)
		// 	qs := make([]float64, 0, endIFrame-startIFrame+1)
		// 	for k := startFrame; k <= endFrame; k++ {
		// 		rangeFrames = append(rangeFrames, k)
		// 		qs = append(qs, startQuat.Dot(quats[int(k)]))
		// 		// initialT := float64(j-startFrame) / float64(endFrame-startFrame)
		// 		// qs = append(qs, mmath.FindSlerpT(startQuat, endQuat, quats[int(i)], initialT))
		// 	}

		// 	// 変曲点を追加する
		// 	inflectionQuatFrames = append(inflectionQuatFrames,
		// 		mmath.FindInflectionFrames(rangeFrames, qs, 1e-6)...)

		// 	if j <= len(inflectionFrames)-2 {
		// 		// 残りが2つ以上ある場合、次の次に進む
		// 		j++
		// 	}

		// 	j++
		// }

		// inflectionFrames = append(inflectionFrames, inflectionQuatFrames...)
	}

	if len(inflectionFrames) <= 2 {
		// 変曲点がない場合、そのまま終了
		return boneNameFrames
	}

	inflectionFrames = mmath.UniqueFloat32s(inflectionFrames)
	mmath.SortFloat32s(inflectionFrames)

	reduceBfs := NewBoneNameFrames(boneNameFrames.Name)
	{
		// 最初のフレームを登録
		bf := boneNameFrames.Get(inflectionFrames[0])
		reduceBf := NewBoneFrame(inflectionFrames[0])
		reduceBf.Registered = true
		reduceBf.Position = bf.Position.Copy()
		reduceBf.Rotation = bf.Rotation.Copy()
		if bf.Curves != nil {
			reduceBf.Curves = bf.Curves.Copy()
		}
		reduceBfs.Append(reduceBf)
	}

	startFrame := inflectionFrames[0]
	midFrame := inflectionFrames[1]
	endFrame := inflectionFrames[2]
	exactEndFrame := float32(0)

	i := 2
	for {
		if exactEndFrame >= maxFrame {
			break
		}

		// print(fmt.Sprintf("startFrame: %f, midFrame: %f, endFrame: %f\n", startFrame, midFrame, endFrame))
		exactEndFrame = boneNameFrames.reduceRange(startFrame, midFrame, endFrame, xs, ys, zs, quats, reduceBfs)

		// 実際に繋げた終了フレームまでを繋ぐ
		exactI := slices.Index(inflectionFrames, exactEndFrame)

		if exactI == -1 {
			// 途中で区切った場合、範囲から決める
			if exactEndFrame < midFrame {
				// 前半で区切っている場合
				startFrame = exactEndFrame
				continue
			} else {
				// 後半で区切っている場合
				startFrame = midFrame
				midFrame = exactEndFrame
				continue
			}
		} else {
			i = exactI
		}

		if i >= len(inflectionFrames)-1 {
			break
		}

		i += 2

		if i >= len(inflectionFrames)-1 {
			break
		} else {
			startFrame = exactEndFrame
			midFrame = inflectionFrames[i-1]
			endFrame = inflectionFrames[i]
		}
	}

	// 最後のフレームを登録
	{
		startFrame := exactEndFrame
		endFrame := inflectionFrames[len(inflectionFrames)-1]
		midFrame := float32(int(startFrame+endFrame) / 2)

		exactEndFrame = boneNameFrames.reduceRange(startFrame, midFrame, endFrame, xs, ys, zs, quats, reduceBfs)

		for exactEndFrame < endFrame {
			// 途中までしか繋げなかった場合、そこから次を探す
			startFrame = exactEndFrame
			midFrame = float32(int(exactEndFrame+endFrame) / 2)

			exactEndFrame = boneNameFrames.reduceRange(startFrame, midFrame, endFrame, xs, ys, zs, quats, reduceBfs)
		}
	}

	return reduceBfs
}

func (boneNameFrames *BoneNameFrames) reduceRange(
	startFrame, midFrame, endFrame float32, xs, ys, zs []float64, quats []*mmath.MQuaternion, reduceBfs *BoneNameFrames,
) float32 {
	startIFrame := int(startFrame)
	endIFrame := int(endFrame)

	var rangeXs, rangeYs, rangeZs []float64
	if len(xs) <= endIFrame {
		rangeXs = xs[startIFrame:]
		rangeYs = ys[startIFrame:]
		rangeZs = zs[startIFrame:]
	} else {
		rangeXs = xs[startIFrame : endIFrame+1]
		rangeYs = ys[startIFrame : endIFrame+1]
		rangeZs = zs[startIFrame : endIFrame+1]
	}

	rangeRs := make([]float64, 0, len(rangeXs))
	startQuat := quats[startIFrame]
	endQuat := quats[endIFrame]

	// rangeFrames := make([]float32, 0, len(rangeXs))
	// qs := make([]float64, 0, len(xs))
	// for i := startFrame; i <= endFrame; i++ {
	// 	// initialT := float64(i-startFrame) / float64(endFrame-startFrame)
	// 	// qs = append(qs, mmath.FindSlerpT(startQuat, endQuat, quats[int(i)], initialT))
	// 	qs = append(qs, startQuat.Dot(quats[int(i)]))
	// 	rangeFrames = append(rangeFrames, i)
	// }

	// inflectionQuatFrames := mmath.FindInflectionFrames(rangeFrames, qs, 1e-5)
	// inflectionQuatFrames = mmath.UniqueFloat32s(inflectionQuatFrames)

	// if len(inflectionQuatFrames) > 3 {
	// 	// 変曲点が多すぎる場合、半分に分割する
	// 	for i := 2; i < len(inflectionQuatFrames); i += 2 {
	// 		rangeStartFrame := inflectionQuatFrames[i-2]
	// 		rangeMidFrame := inflectionQuatFrames[i-1]
	// 		rangeEndFrame := inflectionQuatFrames[i]

	// 		boneNameFrames.reduceRange(rangeStartFrame, rangeMidFrame, rangeEndFrame, xs, ys, zs, quats, reduceBfs)
	// 	}
	// 	return false
	// }

	for i := startFrame; i <= endFrame; i++ {
		// initialT := float64(i-startFrame) / float64(endFrame-startFrame)
		quat := quats[int(i)]
		rangeRs = append(rangeRs, mmath.FindSlerpT(startQuat, endQuat, quat, 0))
	}

	xCurve := mmath.NewCurveFromValues(rangeXs, 1e-2)
	yCurve := mmath.NewCurveFromValues(rangeYs, 1e-2)
	zCurve := mmath.NewCurveFromValues(rangeZs, 1e-2)
	rCurve := mmath.NewCurveFromValues(rangeRs, 1e-4)

	if xCurve != nil && yCurve != nil && zCurve != nil && rCurve != nil {
		isSuccess := true
		for i := startIFrame + 1; i < endIFrame; i++ {
			if !boneNameFrames.checkCurve(
				xCurve, yCurve, zCurve, rCurve,
				xs[startIFrame], xs[i], xs[endIFrame],
				ys[startIFrame], ys[i], ys[endIFrame],
				zs[startIFrame], zs[i], zs[endIFrame],
				quats[startIFrame], quats[i], quats[endIFrame],
				startFrame, float32(i), endFrame,
			) {
				isSuccess = false
				break
			}
		}

		if isSuccess {
			// 全ての曲線が正常に生成された場合、検算する
			bf := boneNameFrames.Get(endFrame)

			reduceBf := NewBoneFrame(endFrame)
			reduceBf.Registered = true
			reduceBf.Position = bf.Position.Copy()
			reduceBf.Rotation = bf.Rotation.Copy()
			reduceBf.Curves = &BoneCurves{
				TranslateX: xCurve,
				TranslateY: yCurve,
				TranslateZ: zCurve,
				Rotate:     rCurve,
			}

			// print(fmt.Sprintf("reduceBf: %v\n", endFrame))
			reduceBfs.Append(reduceBf)

			// endまで繋げられた場合
			return endFrame
		}
	}

	// 生成できなかった場合、半分に分割する
	midIFrame := int(midFrame)
	if midIFrame == startIFrame || midIFrame == endIFrame {
		// 半分に出来なかった場合、そのまま全打ち状態で終了
		{
			bf := boneNameFrames.Get(startFrame)

			reduceBf := NewBoneFrame(startFrame)
			reduceBf.Registered = true
			reduceBf.Position = bf.Position.Copy()
			reduceBf.Rotation = bf.Rotation.Copy()
			if bf.Curves != nil {
				reduceBf.Curves = bf.Curves.Copy()
			}

			reduceBfs.Append(reduceBf)
		}

		// endまで繋げられなかった場合
		return midFrame
	}

	return boneNameFrames.reduceRange(startFrame, float32(int(midFrame+startFrame)/2), midFrame, xs, ys, zs, quats, reduceBfs)
}

// 検算
func (boneNameFrames *BoneNameFrames) checkCurve(
	xCurve, yCurve, zCurve, rCurve *mmath.Curve, startX, nowX, endX, startY, nowY, endY, startZ, nowZ, endZ float64,
	startQuat, nowQuat, endQuat *mmath.MQuaternion, startFrame, nowFrame, endFrame float32,
) bool {
	_, xy, _ := mmath.Evaluate(xCurve, startFrame, nowFrame, endFrame)
	_, yy, _ := mmath.Evaluate(yCurve, startFrame, nowFrame, endFrame)
	_, zy, _ := mmath.Evaluate(zCurve, startFrame, nowFrame, endFrame)
	_, ry, _ := mmath.Evaluate(rCurve, startFrame, nowFrame, endFrame)

	checkNowQuat := startQuat.Slerp(endQuat, ry)
	if !checkNowQuat.NearEquals(nowQuat, 1e-1) {
		return false
	}

	checkNowX := mmath.LerpFloat(startX, endX, xy)
	if !mmath.NearEquals(checkNowX, nowX, 1e-1) {
		return false
	}

	checkNowY := mmath.LerpFloat(startY, endY, yy)
	if !mmath.NearEquals(checkNowY, nowY, 1e-1) {
		return false
	}

	checkNowZ := mmath.LerpFloat(startZ, endZ, zy)
	return mmath.NearEquals(checkNowZ, nowZ, 1e-1)
}
