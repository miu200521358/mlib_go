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

	matrixes := motion.AnimateBone([]int{10, 999}, model, []string{pmx.INDEX3.Left()}, false, false, "")

	{

		fno := 10
		{
			expectedPosition := &mmath.MVec3{0.0, 0.0, 0.0}
			if !matrixes.GetItem(pmx.ROOT.String(), fno).Position.PracticallyEquals(expectedPosition, 1e-5) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem(pmx.ROOT.String(), fno).Position)
			}
		}
		{
			expectedPosition := &mmath.MVec3{0.044920, 8.218059, 0.069347}
			if !matrixes.GetItem(pmx.CENTER.String(), fno).Position.PracticallyEquals(expectedPosition, 1e-5) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem(pmx.CENTER.String(), fno).Position)
			}
		}
		{
			expectedPosition := &mmath.MVec3{0.044920, 9.392067, 0.064877}
			if !matrixes.GetItem(pmx.GROOVE.String(), fno).Position.PracticallyEquals(expectedPosition, 1e-5) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem(pmx.GROOVE.String(), fno).Position)
			}
		}
		{
			expectedPosition := &mmath.MVec3{0.044920, 11.740084, 0.055937}
			if !matrixes.GetItem(pmx.WAIST.String(), fno).Position.PracticallyEquals(expectedPosition, 1e-5) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem(pmx.WAIST.String(), fno).Position)
			}
		}
		{
			expectedPosition := &mmath.MVec3{0.044920, 12.390969, -0.100531}
			if !matrixes.GetItem(pmx.UPPER.String(), fno).Position.PracticallyEquals(expectedPosition, 1e-5) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem(pmx.UPPER.String(), fno).Position)
			}
		}
		{
			expectedPosition := &mmath.MVec3{0.044920, 13.803633, -0.138654}
			if !matrixes.GetItem(pmx.UPPER2.String(), fno).Position.PracticallyEquals(expectedPosition, 1e-5) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem(pmx.UPPER2.String(), fno).Position)
			}
		}
		{
			expectedPosition := &mmath.MVec3{0.044920, 15.149180, 0.044429}
			if !matrixes.GetItem(pmx.UPPER3.String(), fno).Position.PracticallyEquals(expectedPosition, 1e-5) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem(pmx.UPPER3.String(), fno).Position)
			}
		}
		{
			expectedPosition := &mmath.MVec3{0.324862, 16.470263, 0.419041}
			if !matrixes.GetItem(pmx.SHOULDER_P.Left(), fno).Position.PracticallyEquals(expectedPosition, 1e-5) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem(pmx.SHOULDER_P.Left(), fno).Position)
			}
		}
		{
			expectedPosition := &mmath.MVec3{0.324862, 16.470263, 0.419041}
			if !matrixes.GetItem(pmx.SHOULDER.Left(), fno).Position.PracticallyEquals(expectedPosition, 1e-5) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem(pmx.SHOULDER.Left(), fno).Position)
			}
		}
		{
			expectedPosition := &mmath.MVec3{1.369838, 16.312170, 0.676838}
			if !matrixes.GetItem(pmx.ARM.Left(), fno).Position.PracticallyEquals(expectedPosition, 1e-5) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem(pmx.ARM.Left(), fno).Position)
			}
		}
		{
			expectedPosition := &mmath.MVec3{1.845001, 15.024807, 0.747681}
			if !matrixes.GetItem(pmx.ARM_TWIST.Left(), fno).Position.PracticallyEquals(expectedPosition, 1e-5) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem(pmx.ARM_TWIST.Left(), fno).Position)
			}
		}
		{
			expectedPosition := &mmath.MVec3{2.320162, 13.737446, 0.818525}
			if !matrixes.GetItem(pmx.ELBOW.Left(), fno).Position.PracticallyEquals(expectedPosition, 1e-5) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem(pmx.ELBOW.Left(), fno).Position)
			}
		}
		{
			expectedPosition := &mmath.MVec3{2.516700, 12.502447, 0.336127}
			if !matrixes.GetItem(pmx.WRIST_TWIST.Left(), fno).Position.PracticallyEquals(expectedPosition, 1e-2) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem(pmx.WRIST_TWIST.Left(), fno).Position)
			}
		}
		{
			expectedPosition := &mmath.MVec3{2.732219, 11.267447, -0.146273}
			if !matrixes.GetItem(pmx.WRIST.Left(), fno).Position.PracticallyEquals(expectedPosition, 1e-4) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem(pmx.WRIST.Left(), fno).Position)
			}
		}
		{
			expectedPosition := &mmath.MVec3{2.649188, 10.546797, -0.607412}
			if !matrixes.GetItem(pmx.INDEX1.Left(), fno).Position.PracticallyEquals(expectedPosition, 1e-4) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem(pmx.INDEX1.Left(), fno).Position)
			}
		}
		{
			expectedPosition := &mmath.MVec3{2.408238, 10.209290, -0.576288}
			if !matrixes.GetItem(pmx.INDEX2.Left(), fno).Position.PracticallyEquals(expectedPosition, 1e-3) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem(pmx.INDEX2.Left(), fno).Position)
			}
		}
		{
			expectedPosition := &mmath.MVec3{2.360455, 10.422402, -0.442668}
			if !matrixes.GetItem(pmx.INDEX3.Left(), fno).Position.PracticallyEquals(expectedPosition, 1e-3) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem(pmx.INDEX3.Left(), fno).Position)
			}
		}
	}

	{

		fno := 999
		{
			expectedPosition := &mmath.MVec3{0.0, 0.0, 0.0}
			if !matrixes.GetItem(pmx.ROOT.String(), fno).Position.PracticallyEquals(expectedPosition, 1e-3) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem(pmx.ROOT.String(), fno).Position)
			}
		}
		{
			expectedPosition := &mmath.MVec3{-0.508560, 8.218059, 0.791827}
			if !matrixes.GetItem(pmx.CENTER.String(), fno).Position.PracticallyEquals(expectedPosition, 1e-3) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem(pmx.CENTER.String(), fno).Position)
			}
		}
		{
			expectedPosition := &mmath.MVec3{-0.508560, 9.182008, 0.787357}
			if !matrixes.GetItem(pmx.GROOVE.String(), fno).Position.PracticallyEquals(expectedPosition, 1e-3) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem(pmx.GROOVE.String(), fno).Position)
			}
		}
		{
			expectedPosition := &mmath.MVec3{-0.508560, 11.530025, 0.778416}
			if !matrixes.GetItem(pmx.WAIST.String(), fno).Position.PracticallyEquals(expectedPosition, 1e-3) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem(pmx.WAIST.String(), fno).Position)
			}
		}
		{
			expectedPosition := &mmath.MVec3{-0.508560, 12.180910, 0.621949}
			if !matrixes.GetItem(pmx.UPPER.String(), fno).Position.PracticallyEquals(expectedPosition, 1e-3) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem(pmx.UPPER.String(), fno).Position)
			}
		}
		{
			expectedPosition := &mmath.MVec3{-0.437343, 13.588836, 0.523215}
			if !matrixes.GetItem(pmx.UPPER2.String(), fno).Position.PracticallyEquals(expectedPosition, 1e-3) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem(pmx.UPPER2.String(), fno).Position)
			}
		}
		{
			expectedPosition := &mmath.MVec3{-0.552491, 14.941880, 0.528703}
			if !matrixes.GetItem(pmx.UPPER3.String(), fno).Position.PracticallyEquals(expectedPosition, 1e-3) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem(pmx.UPPER3.String(), fno).Position)
			}
		}
		{
			expectedPosition := &mmath.MVec3{-0.590927, 16.312325, 0.819156}
			if !matrixes.GetItem(pmx.SHOULDER_P.Left(), fno).Position.PracticallyEquals(expectedPosition, 1e-3) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem(pmx.SHOULDER_P.Left(), fno).Position)
			}
		}
		{
			expectedPosition := &mmath.MVec3{-0.590927, 16.312325, 0.819156}
			if !matrixes.GetItem(pmx.SHOULDER.Left(), fno).Position.PracticallyEquals(expectedPosition, 1e-3) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem(pmx.SHOULDER.Left(), fno).Position)
			}
		}
		{
			expectedPosition := &mmath.MVec3{0.072990, 16.156742, 1.666761}
			if !matrixes.GetItem(pmx.ARM.Left(), fno).Position.PracticallyEquals(expectedPosition, 1e-3) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem(pmx.ARM.Left(), fno).Position)
			}
		}
		{
			expectedPosition := &mmath.MVec3{0.043336, 15.182318, 2.635117}
			if !matrixes.GetItem(pmx.ARM_TWIST.Left(), fno).Position.PracticallyEquals(expectedPosition, 1e-3) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem(pmx.ARM_TWIST.Left(), fno).Position)
			}
		}
		{
			expectedPosition := &mmath.MVec3{0.013682, 14.207894, 3.603473}
			if !matrixes.GetItem(pmx.ELBOW.Left(), fno).Position.PracticallyEquals(expectedPosition, 1e-3) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem(pmx.ELBOW.Left(), fno).Position)
			}
		}
		{
			expectedPosition := &mmath.MVec3{1.222444, 13.711100, 3.299384}
			if !matrixes.GetItem(pmx.WRIST_TWIST.Left(), fno).Position.PracticallyEquals(expectedPosition, 1e-3) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem(pmx.WRIST_TWIST.Left(), fno).Position)
			}
		}
		{
			expectedPosition := &mmath.MVec3{2.431205, 13.214306, 2.995294}
			if !matrixes.GetItem(pmx.WRIST.Left(), fno).Position.PracticallyEquals(expectedPosition, 1e-3) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem(pmx.WRIST.Left(), fno).Position)
			}
		}
		{
			expectedPosition := &mmath.MVec3{3.283628, 13.209089, 2.884702}
			if !matrixes.GetItem(pmx.INDEX1.Left(), fno).Position.PracticallyEquals(expectedPosition, 1e-3) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem(pmx.INDEX1.Left(), fno).Position)
			}
		}
		{
			expectedPosition := &mmath.MVec3{3.665809, 13.070156, 2.797680}
			if !matrixes.GetItem(pmx.INDEX2.Left(), fno).Position.PracticallyEquals(expectedPosition, 1e-3) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem(pmx.INDEX2.Left(), fno).Position)
			}
		}
		{
			expectedPosition := &mmath.MVec3{3.886795, 12.968100, 2.718276}
			if !matrixes.GetItem(pmx.INDEX3.Left(), fno).Position.PracticallyEquals(expectedPosition, 1e-3) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem(pmx.INDEX3.Left(), fno).Position)
			}
		}
	}

}

func TestVmdMotion_AnimateBoneLegIk1(t *testing.T) {
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

	matrixes := motion.AnimateBone([]int{29}, model, []string{pmx.TOE.Left()}, true, false, "")

	{

		fno := 29
		{
			expectedPosition := &mmath.MVec3{-0.781335, 11.717622, 1.557067}
			if !matrixes.GetItem(pmx.LOWER.String(), fno).Position.PracticallyEquals(expectedPosition, 1e-3) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem(pmx.LOWER.String(), fno).Position)
			}
		}
		{
			expectedPosition := &mmath.MVec3{-0.368843, 10.614175, 2.532657}
			if !matrixes.GetItem(pmx.LEG.Left(), fno).Position.PracticallyEquals(expectedPosition, 1e-3) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem(pmx.LEG.Left(), fno).Position)
			}
		}
		{
			expectedPosition := &mmath.MVec3{0.983212, 6.945313, 0.487476}
			if !matrixes.GetItem(pmx.KNEE.Left(), fno).Position.PracticallyEquals(expectedPosition, 1e-2) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem(pmx.KNEE.Left(), fno).Position)
			}
		}
		{
			expectedPosition := &mmath.MVec3{-0.345842, 2.211842, 2.182894}
			if !matrixes.GetItem(pmx.ANKLE.Left(), fno).Position.PracticallyEquals(expectedPosition, 1e-2) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem(pmx.ANKLE.Left(), fno).Position)
			}
		}
		{
			expectedPosition := &mmath.MVec3{-0.109262, -0.025810, 1.147780}
			if !matrixes.GetItem(pmx.TOE.Left(), fno).Position.PracticallyEquals(expectedPosition, 1e-2) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem(pmx.TOE.Left(), fno).Position)
			}
		}
	}
}
