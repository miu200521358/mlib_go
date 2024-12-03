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
			initialT := float64(iF) / float64(maxIFrame)
			quats = append(quats, bf.Rotation)
			fixRs = append(fixRs, mmath.FindSlerpT(mmath.MQuaternionIdent, mmath.MQuaternionUnitX, bf.Rotation, initialT))
		} else {
			quats = append(quats, mmath.MQuaternionIdent)
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
			startQuat := quats[startIFrame]
			endQuat := quats[endIFrame]
			for i := 1; i <= endIFrame-startIFrame; i++ {
				rangeFrames[i] = float32(i + startIFrame)
				initialT := float64(i) / float64(endIFrame-startIFrame)
				quat := quats[i+startIFrame]
				rangeRs[i] = mmath.FindSlerpT(startQuat, endQuat, quat, initialT)
			}

			inflectionQuatRangeFrames := mmath.FindInflectionFrames(rangeFrames, rangeRs)

			// 見つかった最初の変曲点を対象とする(始点と終点は必ずある)
			inflectionQuatFrames = append(inflectionQuatFrames, inflectionQuatRangeFrames...)
		}

		inflectionQuatFrames = mmath.UniqueFloat32s(inflectionQuatFrames)
		mmath.SortFloat32s(inflectionQuatFrames)

		// 変曲点候補の中から、実際に変曲点として採用するフレームを選ぶ
		startFrame := inflectionQuatFrames[0]
		inflectionFrames = append(inflectionFrames, startFrame, inflectionQuatFrames[len(inflectionQuatFrames)-1])

		i := 0
	quatInflection:
		for {
			startIFrame := int(startFrame)
			if i+2 >= len(inflectionQuatFrames) {
				// 最後まで探しても見つからなかった場合、最後のendFrameを変曲点として登録
				inflectionFrames = append(inflectionFrames, inflectionQuatFrames[len(inflectionQuatFrames)-1])
				break quatInflection
			}

			for j := i + 2; j < len(inflectionQuatFrames); j++ {
				endFrame := inflectionQuatFrames[j]

				endIFrame := int(endFrame)
				frameCnt := endIFrame - startIFrame + 1

				rangeFrames := make([]float32, frameCnt)
				rangeFrames[0] = startFrame
				rangeRs := make([]float64, frameCnt)
				startQuat := quats[startIFrame]
				endQuat := quats[endIFrame]
				for k := 1; k <= endIFrame-startIFrame; k++ {
					rangeFrames[k] = float32(k + startIFrame)
					initialT := float64(k) / float64(endIFrame-startIFrame)
					quat := quats[k+startIFrame]
					rangeRs[k] = mmath.FindSlerpT(startQuat, endQuat, quat, initialT)
				}

				inflectionQuatRangeFrames := mmath.FindInflectionFrames(rangeFrames, rangeRs)
				if len(inflectionQuatRangeFrames) > 2 {
					// 中に変曲点がある場合、そのendFrameを変曲点として登録
					inflectionFrames = append(inflectionFrames, inflectionQuatRangeFrames[1])
					startFrame = inflectionQuatRangeFrames[1]
					i = j - 1
					break
				} else if endFrame == inflectionQuatFrames[len(inflectionQuatFrames)-1] {
					// 最後まで探しても見つからなかった場合、最後のendFrameを変曲点として登録
					inflectionFrames = append(inflectionFrames, endFrame)
					break quatInflection
				}
				// 中に変曲点がない場合、次のendFrameを探す
			}
		}
	}

	inflectionFrames = mmath.UniqueFloat32s(inflectionFrames)
	mmath.SortFloat32s(inflectionFrames)

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

		boneNameFrames.reduceRange(startFrame, midFrame, endFrame, xs, ys, zs, quats, reduceBfs)
	}

	return reduceBfs
}

func (boneNameFrames *BoneNameFrames) reduceRange(
	startFrame, midFrame, endFrame float32, xs, ys, zs []float64, quats []*mmath.MQuaternion, reduceBfs *BoneNameFrames,
) {
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
		initialT := float64(i-startFrame) / float64(endFrame-startFrame)
		quat := quats[int(i)]
		rangeRs = append(rangeRs, mmath.FindSlerpT(startQuat, endQuat, quat, initialT))
	}

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

		boneNameFrames.reduceRange(startFrame, float32(int(midFrame+startFrame)/2), midFrame, xs, ys, zs, quats, reduceBfs)
		boneNameFrames.reduceRange(midFrame, float32(int(endFrame+midFrame)/2), endFrame, xs, ys, zs, quats, reduceBfs)
	}
}
