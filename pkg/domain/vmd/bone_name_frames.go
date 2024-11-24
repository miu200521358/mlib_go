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
	frames := make([]float32, 0, boneNameFrames.Len())
	xs := make([]float64, 0, boneNameFrames.Len())
	ys := make([]float64, 0, boneNameFrames.Len())
	zs := make([]float64, 0, boneNameFrames.Len())
	rs := make([]float64, 0, boneNameFrames.Len())

	boneNameFrames.Indexes.LLRB.AscendGreaterOrEqual(boneNameFrames.Indexes.LLRB.Min(), func(item llrb.Item) bool {
		if float32(item.(core.Float)) >= 0 {
			f := float32(item.(core.Float))

			frames = append(frames, f)
			bf := boneNameFrames.data[f]

			xs = append(xs, bf.Position.X)
			ys = append(ys, bf.Position.Y)
			zs = append(zs, bf.Position.Z)
			rs = append(rs, mmath.MQuaternionIdent.Dot(bf.Rotation))
		}
		return true
	})

	inflectionFrames := make([]float32, 0, boneNameFrames.Len())
	inflectionFrames = append(inflectionFrames, mmath.FindInflectionFrames(frames, xs)...)
	inflectionFrames = append(inflectionFrames, mmath.FindInflectionFrames(frames, ys)...)
	inflectionFrames = append(inflectionFrames, mmath.FindInflectionFrames(frames, zs)...)
	inflectionFrames = append(inflectionFrames, mmath.FindInflectionFrames(frames, rs)...)

	inflectionFrames = mmath.UniqueFloat32s(inflectionFrames)
	mmath.SortFloat32s(inflectionFrames)

	reduceBfs := NewBoneNameFrames(boneNameFrames.Name)
	for i, endFrame := range inflectionFrames {
		if i == 0 {
			continue
		}
		startFrame := inflectionFrames[i-1]

		if i == 1 {
			bf := boneNameFrames.data[startFrame]
			reduceBf := NewBoneFrame(startFrame)
			reduceBf.Registered = bf.Registered
			reduceBf.Position = bf.Position.Copy()
			reduceBf.Rotation = bf.Rotation.Copy()
			reduceBf.Curves = bf.Curves.Copy()
			reduceBfs.Append(reduceBf)
		}

		rangeXs := make([]float64, 0, int(endFrame-startFrame))
		rangeYs := make([]float64, 0, int(endFrame-startFrame))
		rangeZs := make([]float64, 0, int(endFrame-startFrame))
		rangeRs := make([]float64, 0, int(endFrame-startFrame))
		for f := startFrame; f <= endFrame; f++ {
			if _, ok := boneNameFrames.data[f]; ok {
				rangeXs = append(rangeXs, boneNameFrames.data[f].Position.X)
				rangeYs = append(rangeYs, boneNameFrames.data[f].Position.Y)
				rangeZs = append(rangeZs, boneNameFrames.data[f].Position.Z)
				rangeRs = append(rangeRs, mmath.MQuaternionIdent.Dot(boneNameFrames.data[f].Rotation))
			}
		}

		bf := boneNameFrames.data[endFrame]
		reduceBf := NewBoneFrame(endFrame)
		reduceBf.Registered = bf.Registered
		reduceBf.Position = bf.Position.Copy()
		reduceBf.Rotation = bf.Rotation.Copy()
		if len(rangeXs) > 2 {
			reduceBf.Curves = &BoneCurves{
				TranslateX: mmath.NewCurveFromValues(rangeXs),
				TranslateY: mmath.NewCurveFromValues(rangeYs),
				TranslateZ: mmath.NewCurveFromValues(rangeZs),
				Rotate:     mmath.NewCurveFromValues(rangeRs),
			}
		} else {
			reduceBf.Curves = bf.Curves.Copy()
		}

		reduceBfs.Append(reduceBf)
	}

	return reduceBfs
}
