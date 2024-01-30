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
		if !matrixes.GetItem(pmx.ROOT.String(), fno).Position.PracticallyEquals(expectedPosition, 1e-6) {
			t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem(pmx.ROOT.String(), fno).Position)
		}
	}
	{
		expectedPosition := &mmath.MVec3{0.044920, 8.218059, 0.069347}
		if !matrixes.GetItem(pmx.CENTER.String(), fno).Position.PracticallyEquals(expectedPosition, 1e-6) {
			t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem(pmx.CENTER.String(), fno).Position)
		}
	}
	{
		expectedPosition := &mmath.MVec3{0.044920, 9.392067, 0.064877}
		if !matrixes.GetItem(pmx.GROOVE.String(), fno).Position.PracticallyEquals(expectedPosition, 1e-6) {
			t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem(pmx.GROOVE.String(), fno).Position)
		}
	}
	{
		expectedPosition := &mmath.MVec3{0.044920, 11.740084, 0.055937}
		if !matrixes.GetItem(pmx.WAIST.String(), fno).Position.PracticallyEquals(expectedPosition, 1e-6) {
			t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem(pmx.WAIST.String(), fno).Position)
		}
	}
	{
		expectedPosition := &mmath.MVec3{0.044920, 12.390969, -0.100531}
		if !matrixes.GetItem(pmx.UPPER.String(), fno).Position.PracticallyEquals(expectedPosition, 1e-6) {
			t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem(pmx.UPPER.String(), fno).Position)
		}
	}
	{
		expectedPosition := &mmath.MVec3{0.044920, 13.803633, -0.138654}
		if !matrixes.GetItem(pmx.UPPER2.String(), fno).Position.PracticallyEquals(expectedPosition, 1e-6) {
			t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem(pmx.UPPER2.String(), fno).Position)
		}
	}
	{
		expectedPosition := &mmath.MVec3{0.044920, 15.149180, 0.044429}
		if !matrixes.GetItem(pmx.UPPER3.String(), fno).Position.PracticallyEquals(expectedPosition, 1e-6) {
			t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem(pmx.UPPER3.String(), fno).Position)
		}
	}
	{
		expectedPosition := &mmath.MVec3{0.324862, 16.470263, 0.419041}
		if !matrixes.GetItem(pmx.SHOULDER_P.Left(), fno).Position.PracticallyEquals(expectedPosition, 1e-6) {
			t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem(pmx.SHOULDER_P.Left(), fno).Position)
		}
	}
	{
		expectedPosition := &mmath.MVec3{0.324862, 16.470263, 0.419041}
		if !matrixes.GetItem(pmx.SHOULDER.Left(), fno).Position.PracticallyEquals(expectedPosition, 1e-6) {
			t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem(pmx.SHOULDER.Left(), fno).Position)
		}
	}
	{
		expectedPosition := &mmath.MVec3{1.369838, 16.312170, 0.676838}
		if !matrixes.GetItem(pmx.ARM.Left(), fno).Position.PracticallyEquals(expectedPosition, 1e-6) {
			t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem(pmx.ARM.Left(), fno).Position)
		}
	}
	{
		expectedPosition := &mmath.MVec3{1.845001, 15.024807, 0.747681}
		if !matrixes.GetItem(pmx.ARM_TWIST.Left(), fno).Position.PracticallyEquals(expectedPosition, 1e-6) {
			t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem(pmx.ARM_TWIST.Left(), fno).Position)
		}
	}
	{
		expectedPosition := &mmath.MVec3{2.320162, 13.737446, 0.818525}
		if !matrixes.GetItem(pmx.ELBOW.Left(), fno).Position.PracticallyEquals(expectedPosition, 1e-6) {
			t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem(pmx.ELBOW.Left(), fno).Position)
		}
	}
	{
		expectedPosition := &mmath.MVec3{2.516700, 12.502447, 0.336127}
		if !matrixes.GetItem(pmx.WRIST_TWIST.Left(), fno).Position.PracticallyEquals(expectedPosition, 1e-6) {
			t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem(pmx.WRIST_TWIST.Left(), fno).Position)
		}
	}
	{
		expectedPosition := &mmath.MVec3{2.732219, 11.267447, -0.146273}
		if !matrixes.GetItem(pmx.WRIST.Left(), fno).Position.PracticallyEquals(expectedPosition, 1e-6) {
			t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem(pmx.WRIST.Left(), fno).Position)
		}
	}
	{
		expectedPosition := &mmath.MVec3{2.649188, 10.546797, -0.607412}
		if !matrixes.GetItem(pmx.INDEX1.Left(), fno).Position.PracticallyEquals(expectedPosition, 1e-6) {
			t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem(pmx.INDEX1.Left(), fno).Position)
		}
	}
	{
		expectedPosition := &mmath.MVec3{2.408238, 10.209290, -0.576288}
		if !matrixes.GetItem(pmx.INDEX2.Left(), fno).Position.PracticallyEquals(expectedPosition, 1e-6) {
			t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem(pmx.INDEX2.Left(), fno).Position)
		}
	}
	{
		expectedPosition := &mmath.MVec3{2.360455, 10.422402, -0.442668}
		if !matrixes.GetItem(pmx.INDEX3.Left(), fno).Position.PracticallyEquals(expectedPosition, 1e-6) {
			t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem(pmx.INDEX3.Left(), fno).Position)
		}
	}
}
