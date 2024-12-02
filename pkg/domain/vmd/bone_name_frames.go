package vmd

import (
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
	maxFrame := boneNameFrames.Indexes.Max() + 1
	maxIFrame := int(maxFrame)

	frames := make([]float32, 0, maxIFrame)
	xs := make([]float64, 0, maxIFrame)
	ys := make([]float64, 0, maxIFrame)
	zs := make([]float64, 0, maxIFrame)
	rs := make([]float64, 0, maxIFrame)
	fixRs := make([]float64, 0, maxIFrame)

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
			fixRs = append(fixRs, mmath.FindSlerpT(mmath.MQuaternionIdent, mmath.MQuaternionUnitX, bf.Rotation))
		} else {
			fixRs = append(fixRs, 0)
		}
	}

	inflectionFrames := make([]float32, 0, boneNameFrames.Len())
	if !mmath.IsAllSameValues(xs) {
		inflectionFrames = append(inflectionFrames, mmath.FindInflectionFrames(frames, xs)...)
	}
	if !mmath.IsAllSameValues(ys) {
		inflectionFrames = append(inflectionFrames, mmath.FindInflectionFrames(frames, ys)...)
	}
	if !mmath.IsAllSameValues(zs) {
		inflectionFrames = append(inflectionFrames, mmath.FindInflectionFrames(frames, zs)...)
	}
	isAllSameRs := mmath.IsAllSameValues(fixRs)
	if !isAllSameRs {
		inflectionQuatFrames := make([]float32, 0, boneNameFrames.Len())

		// 一旦固定q0, q1を元に算出したtから変曲点を求める
		inflectionFixQuatFrames := mmath.FindInflectionFrames(frames, fixRs)
		inflectionQuatFrames = append(inflectionQuatFrames, inflectionFixQuatFrames...)

		for i, endFrame := range inflectionFixQuatFrames {
			if i == 0 {
				continue
			}

			startFrame := inflectionFixQuatFrames[i-1]
			startIFrame := int(startFrame)
			endIFrame := int(endFrame)
			frameCnt := endIFrame - startIFrame + 1

			rangeFrames := make([]float32, frameCnt)
			rangeFrames[0] = startFrame
			rangeRs := make([]float64, frameCnt)
			startQuat := boneNameFrames.Get(startFrame).Rotation
			endQuat := boneNameFrames.Get(endFrame).Rotation
			for i := 1; i <= endIFrame-startIFrame; i++ {
				rangeFrames[i] = float32(i + startIFrame)
				rangeRs[i] = mmath.FindSlerpT(startQuat, endQuat, boneNameFrames.Get(float32(i+startIFrame)).Rotation)
			}

			inflectionQuatRangeFrames := mmath.FindInflectionFrames(rangeFrames, rangeRs)

			// 見つかった最初の変曲点を対象とする(始点と終点は必ずある)
			inflectionQuatFrames = append(inflectionQuatFrames, inflectionQuatRangeFrames...)
		}

		inflectionQuatFrames = mmath.UniqueFloat32s(inflectionQuatFrames)
		mmath.SortFloat32s(inflectionQuatFrames)

		// 変曲点候補の中から、実際に変曲点として採用するフレームを選ぶ
		startFrame := inflectionQuatFrames[0]
		inflectionFrames = append(inflectionFrames, startFrame, maxFrame)

		i := 0
		for {
			startIFrame := int(startFrame)
			if startIFrame >= int(inflectionQuatFrames[len(inflectionQuatFrames)-3]) {
				inflectionFrames = append(inflectionFrames, inflectionQuatFrames[len(inflectionQuatFrames)-2])
				break
			}

			for j := i + 2; j < len(inflectionQuatFrames); j++ {
				endFrame := inflectionQuatFrames[j]

				endIFrame := int(endFrame)
				frameCnt := endIFrame - startIFrame + 1

				rangeFrames := make([]float32, frameCnt)
				rangeFrames[0] = startFrame
				rangeRs := make([]float64, frameCnt)
				startQuat := boneNameFrames.Get(startFrame).Rotation
				endQuat := boneNameFrames.Get(endFrame).Rotation
				for k := 1; k <= endIFrame-startIFrame; k++ {
					rangeFrames[k] = float32(k + startIFrame)
					rangeRs[k] = mmath.FindSlerpT(startQuat, endQuat, boneNameFrames.Get(float32(k+startIFrame)).Rotation)
				}

				inflectionQuatRangeFrames := mmath.FindInflectionFrames(rangeFrames, rangeRs)
				if len(inflectionQuatRangeFrames) > 2 {
					// 中に変曲点がある場合、そのendFrameを変曲点として登録
					inflectionFrames = append(inflectionFrames, inflectionQuatRangeFrames[1])
					startFrame = inflectionQuatRangeFrames[1]
					i = j - 1
					break
				}
				// 中に変曲点がない場合、次のendFrameを探す
			}
		}
	}

	inflectionFrames = mmath.UniqueFloat32s(inflectionFrames)
	mmath.SortFloat32s(inflectionFrames)

	if !isAllSameRs {
		// 最終的に求まった変曲点リストからtを求める
		for i, endFrame := range inflectionFrames {
			if i == 0 {
				continue
			}

			startFrame := inflectionFrames[i-1]
			startQuat := boneNameFrames.Get(startFrame).Rotation
			endQuat := boneNameFrames.Get(endFrame).Rotation
			for i := startFrame + 1; i <= endFrame; i++ {
				rs = append(rs, mmath.FindSlerpT(startQuat, endQuat, boneNameFrames.Get(i).Rotation))
			}
		}
	} else {
		rs = make([]float64, len(xs))
	}

	reduceBfs := NewBoneNameFrames(boneNameFrames.Name)
	for i := 0; i < len(inflectionFrames); i += 2 {
		if i == 0 {
			bf := boneNameFrames.Get(inflectionFrames[i])

			reduceBf := NewBoneFrame(inflectionFrames[i])
			reduceBf.Registered = true
			reduceBf.Position = bf.Position.Copy()
			reduceBf.Rotation = bf.Rotation.Copy()
			reduceBf.Curves = bf.Curves.Copy()
			reduceBfs.Append(reduceBf)

			continue
		}

		startFrame := inflectionFrames[i-2]
		midFrame := inflectionFrames[i-1]
		endFrame := inflectionFrames[i]

		boneNameFrames.reduceRange(startFrame, midFrame, endFrame, xs, ys, zs, rs, reduceBfs)
	}

	return reduceBfs
}

func (boneNameFrames *BoneNameFrames) reduceRange(
	startFrame, midFrame, endFrame float32, xs, ys, zs, rs []float64, reduceBfs *BoneNameFrames,
) {
	startIFrame := int(startFrame)
	endIFrame := int(endFrame)

	rangeXs := xs[startIFrame:endIFrame]
	rangeYs := ys[startIFrame:endIFrame]
	rangeZs := zs[startIFrame:endIFrame]
	rangeRs := rs[startIFrame:endIFrame]

	xCurve := mmath.NewCurveFromValues(rangeXs)
	yCurve := mmath.NewCurveFromValues(rangeYs)
	zCurve := mmath.NewCurveFromValues(rangeZs)
	rCurve := mmath.NewCurveFromValues(rangeRs)

	if xCurve != nil && yCurve != nil && zCurve != nil && rCurve != nil {
		// 全ての曲線が正常に生成された場合
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

		reduceBfs.Append(reduceBf)
	} else {
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
				} else {
					reduceBf.Curves = NewBoneCurves()
				}
			}
			{
				bf := boneNameFrames.Get(endFrame)

				reduceBf := NewBoneFrame(endFrame)
				reduceBf.Registered = true
				reduceBf.Position = bf.Position.Copy()
				reduceBf.Rotation = bf.Rotation.Copy()
				if bf.Curves != nil {
					reduceBf.Curves = bf.Curves.Copy()
				} else {
					reduceBf.Curves = NewBoneCurves()
				}
			}
		}

		boneNameFrames.reduceRange(startFrame, float32(int(midFrame+startFrame)/2), midFrame, xs, ys, zs, rs, reduceBfs)
		boneNameFrames.reduceRange(midFrame, float32(int(endFrame+midFrame)/2), endFrame, xs, ys, zs, rs, reduceBfs)
	}
}
