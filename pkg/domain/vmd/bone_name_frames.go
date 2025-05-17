package vmd

import (
	"slices"

	"github.com/miu200521358/mlib_go/pkg/domain/mmath"
	"github.com/tiendc/go-deepcopy"
)

type BoneNameFrames struct {
	*BaseFrames[*BoneFrame]
	Name string // ボーン名
}

func NewBoneNameFrames(name string) *BoneNameFrames {
	return &BoneNameFrames{
		BaseFrames: NewBaseFrames(NewBoneFrame, nilBoneFrame),
		Name:       name,
	}
}

func nilBoneFrame() *BoneFrame {
	return nil
}

func (boneNameFrames *BoneNameFrames) Reduce() *BoneNameFrames {
	maxFrame := boneNameFrames.values.Max()
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

	inflectionFrames := make([]float32, 0, boneNameFrames.Length())
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
	}

	inflectionFrames = mmath.Unique(inflectionFrames)
	mmath.Sort(inflectionFrames)

	if len(inflectionFrames) <= 2 {
		// 変曲点がない場合、そのまま終了
		return boneNameFrames
	}

	reduceBfs := NewBoneNameFrames(boneNameFrames.Name)
	{
		// 最初のフレームを登録
		bf := boneNameFrames.Get(inflectionFrames[0])
		reduceBf := NewBoneFrame(inflectionFrames[0])
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
			// 全ての曲線が正常に生成された場合、検算する
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
			// 検算が成功した場合、最終フレームを登録して終了
			bf := boneNameFrames.Get(endFrame)

			reduceBf := NewBoneFrame(endFrame)
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

	checkNowX := mmath.Lerp(startX, endX, xy)
	if !mmath.NearEquals(checkNowX, nowX, 1e-1) {
		return false
	}

	checkNowY := mmath.Lerp(startY, endY, yy)
	if !mmath.NearEquals(checkNowY, nowY, 1e-1) {
		return false
	}

	checkNowZ := mmath.Lerp(startZ, endZ, zy)
	return mmath.NearEquals(checkNowZ, nowZ, 1e-1)
}

// ContainsActive 有効なキーフレが存在するか
func (boneNameFrames *BoneNameFrames) ContainsActive() bool {
	if boneNameFrames.Length() == 0 {
		return false
	}

	isActive := false
	boneNameFrames.ForEach(func(index float32, bf *BoneFrame) bool {
		if bf == nil {
			return true
		}

		if (bf.Position != nil && !bf.Position.NearEquals(mmath.MVec3Zero, 1e-2)) ||
			(bf.Rotation != nil && !bf.Rotation.NearEquals(mmath.MQuaternionIdent, 1e-2)) {
			isActive = true
			return false
		}

		nextBf := boneNameFrames.Get(boneNameFrames.NextFrame(bf.Index()))

		if nextBf == nil {
			return true
		}

		if bf.Position != nil && nextBf.Position != nil && !bf.Position.NearEquals(nextBf.Position, 1e-2) {
			isActive = true
			return false
		}

		if bf.Rotation != nil && nextBf.Rotation != nil && !bf.Rotation.NearEquals(nextBf.Rotation, 1e-2) {
			isActive = true
			return false
		}

		return true
	})

	return isActive
}

func (boneNameFrames *BoneNameFrames) Copy() (*BoneNameFrames, error) {
	copied := new(BoneNameFrames)
	err := deepcopy.Copy(copied, boneNameFrames)
	return copied, err
}
