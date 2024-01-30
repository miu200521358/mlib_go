package vmd

import (
	"testing"

	"github.com/miu200521358/mlib_go/pkg/mmath"
	"github.com/miu200521358/mlib_go/pkg/pmx"
)

func TestVmdMotion_AnimateBone(t *testing.T) {
	vr := &VmdMotionReader{}
	motionData, err := vr.ReadByFilepath("../../test_resources/サンプルモーション.vmd")

	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err)
	}

	motion := motionData.(*VmdMotion)

	pr := &pmx.PmxReader{}
	modelData, err := pr.ReadByFilepath("../../test_resources/サンプルモデル.pmx")

	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err)
	}

	model := modelData.(*pmx.PmxModel)
	model.SetUp()

	matrixes := motion.AnimateBone([]int{10}, model, []string{pmx.INDEX3.Left()}, false, false, "")

	fno := 10
	{
		expectedPosition := &mmath.MVec3{0.0, 0.0, 0.0}
		if !matrixes.GetItem(pmx.ROOT.String(), fno).Position.PracticallyEquals(expectedPosition, 1e-8) {
			t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem(pmx.ROOT.String(), fno).Position)
		}
	}
	{
		expectedPosition := &mmath.MVec3{0.044920, 8.218059, 0.069347}
		if !matrixes.GetItem(pmx.CENTER.String(), fno).Position.PracticallyEquals(expectedPosition, 1e-8) {
			t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem(pmx.CENTER.String(), fno).Position)
		}
	}

}
