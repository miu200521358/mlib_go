package vmd

import (
	"github.com/miu200521358/mlib_go/pkg/domain/core"
	"github.com/miu200521358/mlib_go/pkg/domain/mmath"
	"github.com/petar/GoLLRB/llrb"
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
	maxFrame := int(boneNameFrames.Indexes.Max() + 1)

	frames := make([]float32, 0, maxFrame)
	xs := make([]float64, 0, maxFrame)
	ys := make([]float64, 0, maxFrame)
	zs := make([]float64, 0, maxFrame)
	rs := make([]float64, 0, maxFrame)

	for iF := range maxFrame {
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
			rs = append(rs, mmath.FindSlerpT(mmath.MQuaternionIdent, mmath.MQuaternionUnitX, bf.Rotation))
		} else {
			rs = append(rs, 0)
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
	if !mmath.IsAllSameValues(rs) {
		inflectionFrames = append(inflectionFrames, mmath.FindInflectionFrames(frames, rs)...)
	}

	boneNameFrames.RegisteredIndexes.AscendRange(core.Float(0), core.Float(boneNameFrames.RegisteredIndexes.Max()), func(i llrb.Item) bool {
		if v, ok := boneNameFrames.data[float32(i.(core.Float))]; ok && v.Read {
			inflectionFrames = append(inflectionFrames, float32(i.(core.Float)))
		}
		return true
	})

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
