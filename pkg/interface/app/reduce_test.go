package app

import (
	"fmt"
	"testing"

	"github.com/miu200521358/mlib_go/pkg/domain/mmath"
	"github.com/miu200521358/mlib_go/pkg/domain/vmd"
	"github.com/miu200521358/mlib_go/pkg/infrastructure/repository"
)

func TestReduce_Gimme(t *testing.T) {
	tests := []struct {
		name  string
		frame float32
		curve *vmd.BoneCurves
	}{
		{
			name:  "左肩",
			frame: 200,
			curve: &vmd.BoneCurves{
				TranslateX: mmath.NewCurveByValues(72, 0, 106, 73),
				TranslateY: mmath.NewCurveByValues(72, 0, 106, 73),
				TranslateZ: mmath.NewCurveByValues(72, 0, 106, 73),
				Rotate:     mmath.NewCurveByValues(72, 0, 106, 73),
			},
		},
		{
			name:  "センター",
			frame: 0,
			curve: &vmd.BoneCurves{
				TranslateX: mmath.NewCurveByValues(20, 20, 107, 107),
				TranslateY: mmath.NewCurveByValues(20, 20, 107, 107),
				TranslateZ: mmath.NewCurveByValues(20, 20, 107, 107),
				Rotate:     mmath.NewCurveByValues(20, 20, 107, 107),
			},
		},
		{
			name:  "センター",
			frame: 66,
			curve: &vmd.BoneCurves{
				TranslateX: mmath.NewCurveByValues(72, 0, 106, 73),
				TranslateY: mmath.NewCurveByValues(72, 0, 106, 73),
				TranslateZ: mmath.NewCurveByValues(72, 0, 106, 73),
				Rotate:     mmath.NewCurveByValues(72, 0, 106, 73),
			},
		},
	}

	rep := repository.NewVmdRepository()
	data, err := rep.Load("../../../test_resources/ぎみぎみっちゃん_sam式燭台切光忠（ベスト）Ver1.51_LUSAFW.vmd")
	if err != nil {
		t.Errorf("Load error: %v", err)
		return
	}
	motion := data.(*vmd.VmdMotion)

	for n, tt := range tests {
		t.Run(fmt.Sprintf("%s [%04.0f]", tt.name, tt.frame), func(t *testing.T) {
			bfs := motion.BoneFrames.Get(tt.name).Reduce()

			reduceBf := bfs.Get(tt.frame)
			bf := motion.BoneFrames.Get(tt.name).Get(tt.frame)
			if bf.Curves == nil {
				prevFrame := motion.BoneFrames.Get(tt.name).PrevFrame(tt.frame)
				nextFrame := motion.BoneFrames.Get(tt.name).NextFrame(tt.frame)
				nextBf := motion.BoneFrames.Get(tt.name).Get(nextFrame)
				translateXStart, _ := mmath.SplitCurve(nextBf.Curves.TranslateX, prevFrame, tt.frame, nextFrame)
				translateYStart, _ := mmath.SplitCurve(nextBf.Curves.TranslateY, prevFrame, tt.frame, nextFrame)
				translateZStart, _ := mmath.SplitCurve(nextBf.Curves.TranslateZ, prevFrame, tt.frame, nextFrame)
				rotateStart, _ := mmath.SplitCurve(nextBf.Curves.Rotate, prevFrame, tt.frame, nextFrame)
				bf.Curves = &vmd.BoneCurves{
					TranslateX: translateXStart,
					TranslateY: translateYStart,
					TranslateZ: translateZStart,
					Rotate:     rotateStart,
				}
			}

			if !bf.Position.NearEquals(reduceBf.Position, 1e-6) {
				t.Errorf("%d:%s [%04.0f]: BoneFrame.Position want %v, expect %v",
					n, tt.name, tt.frame, bf.Position, reduceBf.Position)
			}
			if !bf.Rotation.NearEquals(reduceBf.Rotation, 1e-6) {
				t.Errorf("%d:%s [%04.0f]: BoneFrame.Rotation want %v, expect %v",
					n, tt.name, tt.frame, bf.Rotation, reduceBf.Rotation)
			}
			if !bf.Curves.TranslateX.Start.NearEquals(&reduceBf.Curves.TranslateX.Start, 1e-6) {
				t.Errorf("%d:%s [%04.0f]: BoneFrame.Curves.TranslateX.Start want %v, expect %v",
					n, tt.name, tt.frame, bf.Curves.TranslateX.Start, reduceBf.Curves.TranslateX.Start)
			}
			if !bf.Curves.TranslateX.End.NearEquals(&reduceBf.Curves.TranslateX.End, 1e-6) {
				t.Errorf("%d:%s [%04.0f]: BoneFrame.Curves.TranslateX.End want %v, expect %v",
					n, tt.name, tt.frame, bf.Curves.TranslateX.End, reduceBf.Curves.TranslateX.End)
			}
			if !bf.Curves.TranslateY.Start.NearEquals(&reduceBf.Curves.TranslateY.Start, 1e-6) {
				t.Errorf("%d:%s [%04.0f]: BoneFrame.Curves.TranslateY.Start want %v, expect %v",
					n, tt.name, tt.frame, bf.Curves.TranslateY.Start, reduceBf.Curves.TranslateY.Start)
			}
			if !bf.Curves.TranslateY.End.NearEquals(&reduceBf.Curves.TranslateY.End, 1e-6) {
				t.Errorf("%d:%s [%04.0f]: BoneFrame.Curves.TranslateY.End want %v, expect %v",
					n, tt.name, tt.frame, bf.Curves.TranslateY.End, reduceBf.Curves.TranslateY.End)
			}
			if !bf.Curves.TranslateZ.Start.NearEquals(&reduceBf.Curves.TranslateZ.Start, 1e-6) {
				t.Errorf("%d:%s [%04.0f]: BoneFrame.Curves.TranslateZ.Start want %v, expect %v",
					n, tt.name, tt.frame, bf.Curves.TranslateZ.Start, reduceBf.Curves.TranslateZ.Start)
			}
			if !bf.Curves.TranslateZ.End.NearEquals(&reduceBf.Curves.TranslateZ.End, 1e-6) {
				t.Errorf("%d:%s [%04.0f]: BoneFrame.Curves.TranslateZ.End want %v, expect %v",
					n, tt.name, tt.frame, bf.Curves.TranslateZ.End, reduceBf.Curves.TranslateZ.End)
			}
			if !bf.Curves.Rotate.Start.NearEquals(&reduceBf.Curves.Rotate.Start, 1e-6) {
				t.Errorf("%d:%s [%04.0f]: BoneFrame.Curves.Rotate.Start want %v, expect %v",
					n, tt.name, tt.frame, bf.Curves.Rotate.Start, reduceBf.Curves.Rotate.Start)
			}
			if !bf.Curves.Rotate.End.NearEquals(&reduceBf.Curves.Rotate.End, 1e-6) {
				t.Errorf("%d:%s [%04.0f]: BoneFrame.Curves.Rotate.End want %v, expect %v",
					n, tt.name, tt.frame, bf.Curves.Rotate.End, reduceBf.Curves.Rotate.End)
			}
		})
	}

	reduceMotion := vmd.NewVmdMotion("")
	reduceMotion.SetName("reduced")
	reduceMotion.BoneFrames = motion.BoneFrames.Reduce()
	rep.Save("../../../test_resources/test.vmd", reduceMotion, false)
}
