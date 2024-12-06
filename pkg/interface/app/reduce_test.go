package app

import (
	"fmt"
	"testing"

	"github.com/miu200521358/mlib_go/pkg/domain/vmd"
	"github.com/miu200521358/mlib_go/pkg/infrastructure/repository"
)

func TestReduce_Gimme(t *testing.T) {
	tests := []struct {
		name     string
		minFrame float32
		maxFrame float32
	}{
		{
			name:     "センター",
			minFrame: 500,
			maxFrame: 517,
		},
	}

	rep := repository.NewVmdRepository()
	// data, err := rep.Load("../../../test_resources/ぎみぎみっちゃん_sam式燭台切光忠（ベスト）Ver1.51_LUSAFW.vmd")
	data, err := rep.Load("D:/MMD/MikuMikuDance_v926x64/UserFile/Motion/ダンス_1人/GimmeGimme シガー/ぎみぎみっちゃん.vmd")
	if err != nil {
		t.Errorf("Load error: %v", err)
		return
	}
	motion := data.(*vmd.VmdMotion)

	reduceMotion := vmd.NewVmdMotion("")
	reduceMotion.SetName("reduced")
	reduceMotion.BoneFrames = motion.BoneFrames.Reduce()
	rep.Save("../../../test_resources/test.vmd", reduceMotion, false)

	data, _ = rep.Load("../../../test_resources/test.vmd")
	reduceMotion = data.(*vmd.VmdMotion)

	for n, tt := range tests {
		t.Run(fmt.Sprintf("%s[%04.0f-%04.0f]", tt.name, tt.minFrame, tt.maxFrame), func(t *testing.T) {
			for f := tt.minFrame; f <= tt.maxFrame; f++ {
				reduceBf := reduceMotion.BoneFrames.Get(tt.name).Get(f)
				bf := motion.BoneFrames.Get(tt.name).Get(f)

				if bf.Position.NearEquals(reduceBf.Position, 1e-1) {
					t.Errorf("%d:%s [%04.0f]: BoneFrame.Position want %v, expect %v, diff %v",
						n, tt.name, f, bf.Position, reduceBf.Position, bf.Position.Subed(reduceBf.Position))
				}
				if !bf.Rotation.NearEquals(reduceBf.Rotation, 1e-1) {
					t.Errorf("%d:%s [%04.0f]: BoneFrame.Rotation want %v(%v), expect %v(%v), diff %v(%v)",
						n, tt.name, f, bf.Rotation.ToDegrees(), bf.Rotation, reduceBf.Rotation.ToDegrees(),
						reduceBf.Rotation, bf.Rotation.ToDegrees().Subed(reduceBf.Rotation.ToDegrees()),
						bf.Rotation.Vec4().Subed(reduceBf.Rotation.Vec4()))
				}
			}
		})
	}
}
