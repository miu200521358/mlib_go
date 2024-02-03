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

func TestVmdMotion_AnimateBoneLegIk1_Matsu(t *testing.T) {
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

func TestVmdMotion_AnimateBoneLegIk2_Matsu(t *testing.T) {
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

	matrixes := motion.AnimateBone([]int{3152}, model, []string{pmx.TOE.Left()}, true, false, "")

	{

		fno := 3152
		{
			expectedPosition := &mmath.MVec3{7.928583, 11.713336, 1.998830}
			if !matrixes.GetItem(pmx.LOWER.String(), fno).Position.PracticallyEquals(expectedPosition, 1e-3) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem(pmx.LOWER.String(), fno).Position)
			}
		}
		{
			expectedPosition := &mmath.MVec3{7.370017, 10.665785, 2.963280}
			if !matrixes.GetItem(pmx.LEG.Left(), fno).Position.PracticallyEquals(expectedPosition, 1e-3) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem(pmx.LEG.Left(), fno).Position)
			}
		}
		{
			expectedPosition := &mmath.MVec3{9.282883, 6.689319, 2.96825}
			if !matrixes.GetItem(pmx.KNEE.Left(), fno).Position.PracticallyEquals(expectedPosition, 1e-2) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem(pmx.KNEE.Left(), fno).Position)
			}
		}
		{
			expectedPosition := &mmath.MVec3{4.115521, 7.276527, 2.980609}
			if !matrixes.GetItem(pmx.ANKLE.Left(), fno).Position.PracticallyEquals(expectedPosition, 1e-2) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem(pmx.ANKLE.Left(), fno).Position)
			}
		}
		{
			expectedPosition := &mmath.MVec3{1.931355, 6.108739, 2.994883}
			if !matrixes.GetItem(pmx.TOE.Left(), fno).Position.PracticallyEquals(expectedPosition, 1e-2) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem(pmx.TOE.Left(), fno).Position)
			}
		}
	}
}

func TestVmdMotion_AnimateBoneLegIk3_Matsu(t *testing.T) {
	vr := &VmdMotionReader{}
	motionData, err := vr.ReadByFilepath("../../test_resources/好き雪.vmd")

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

	matrixes := motion.AnimateBone([]int{60}, model, nil, true, false, "")

	{

		fno := 60
		{
			expectedPosition := &mmath.MVec3{1.931959, 11.695199, -1.411883}
			if !matrixes.GetItem(pmx.LOWER.String(), fno).Position.PracticallyEquals(expectedPosition, 1e-3) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem(pmx.LOWER.String(), fno).Position)
			}
		}
		{
			expectedPosition := &mmath.MVec3{2.927524, 10.550287, -1.218106}
			if !matrixes.GetItem(pmx.LEG.Left(), fno).Position.PracticallyEquals(expectedPosition, 1e-3) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem(pmx.LEG.Left(), fno).Position)
			}
		}
		{
			expectedPosition := &mmath.MVec3{2.263363, 7.061642, -3.837192}
			if !matrixes.GetItem(pmx.KNEE.Left(), fno).Position.PracticallyEquals(expectedPosition, 1e-2) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem(pmx.KNEE.Left(), fno).Position)
			}
		}
		{
			expectedPosition := &mmath.MVec3{2.747242, 2.529942, -1.331971}
			if !matrixes.GetItem(pmx.ANKLE.Left(), fno).Position.PracticallyEquals(expectedPosition, 1e-2) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem(pmx.ANKLE.Left(), fno).Position)
			}
		}
		{
			expectedPosition := &mmath.MVec3{1.809291, 0.242514, -1.182168}
			if !matrixes.GetItem(pmx.TOE.Left(), fno).Position.PracticallyEquals(expectedPosition, 1e-2) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem(pmx.TOE.Left(), fno).Position)
			}
		}
		{
			expectedPosition := &mmath.MVec3{2.263363, 7.061642, -3.837192}
			if !matrixes.GetItem(pmx.KNEE_D.Left(), fno).Position.PracticallyEquals(expectedPosition, 1e-2) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem(pmx.KNEE.Left(), fno).Position)
			}
		}
		{
			expectedPosition := &mmath.MVec3{1.916109, 1.177077, -1.452845}
			if !matrixes.GetItem(pmx.TOE_EX.Left(), fno).Position.PracticallyEquals(expectedPosition, 1e-2) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem(pmx.TOE.Left(), fno).Position)
			}
		}
	}
}

func TestVmdMotion_AnimateBoneLegIk4_Snow(t *testing.T) {
	vr := &VmdMotionReader{}
	motionData, err := vr.ReadByFilepath("../../test_resources/好き雪_2794.vmd")

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

	matrixes := motion.AnimateBone([]int{0}, model, []string{pmx.TOE.Right()}, true, false, "")

	{

		fno := 0
		{
			expectedPosition := &mmath.MVec3{1.316121, 11.687257, 2.263307}
			if !matrixes.GetItem(pmx.LOWER.String(), fno).Position.PracticallyEquals(expectedPosition, 1e-3) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem(pmx.LOWER.String(), fno).Position)
			}
		}
		{
			expectedPosition := &mmath.MVec3{0.175478, 10.780540, 2.728409}
			if !matrixes.GetItem(pmx.LEG.Right(), fno).Position.PracticallyEquals(expectedPosition, 1e-3) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem(pmx.LEG.Right(), fno).Position)
			}
		}
		{
			expectedPosition := &mmath.MVec3{0.950410, 11.256771, -1.589462}
			if !matrixes.GetItem(pmx.KNEE.Right(), fno).Position.PracticallyEquals(expectedPosition, 1e-2) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem(pmx.KNEE.Right(), fno).Position)
			}
		}
		{
			expectedPosition := &mmath.MVec3{-1.025194, 7.871110, 1.828258}
			if !matrixes.GetItem(pmx.ANKLE.Right(), fno).Position.PracticallyEquals(expectedPosition, 1e-2) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem(pmx.ANKLE.Right(), fno).Position)
			}
		}
		{
			expectedPosition := &mmath.MVec3{-1.701147, 6.066556, 3.384271}
			if !matrixes.GetItem(pmx.TOE.Right(), fno).Position.PracticallyEquals(expectedPosition, 1e-2) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem(pmx.TOE.Right(), fno).Position)
			}
		}
	}
}

func TestVmdMotion_AnimateBoneLegIk5_Snow(t *testing.T) {
	vr := &VmdMotionReader{}
	motionData, err := vr.ReadByFilepath("../../test_resources/好き雪.vmd")

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

	matrixes := motion.AnimateBone([]int{7409}, model, []string{pmx.TOE.Right()}, true, false, "")

	{

		fno := 7409
		{
			expectedPosition := &mmath.MVec3{-7.652257, 11.990970, -4.511993}
			if !matrixes.GetItem(pmx.LOWER.String(), fno).Position.PracticallyEquals(expectedPosition, 1e-3) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem(pmx.LOWER.String(), fno).Position)
			}
		}
		{
			expectedPosition := &mmath.MVec3{-8.637265, 10.835548, -4.326830}
			if !matrixes.GetItem(pmx.LEG.Right(), fno).Position.PracticallyEquals(expectedPosition, 1e-3) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem(pmx.LEG.Right(), fno).Position)
			}
		}
		{
			expectedPosition := &mmath.MVec3{-8.693436, 7.595280, -7.321638}
			if !matrixes.GetItem(pmx.KNEE.Right(), fno).Position.PracticallyEquals(expectedPosition, 1e-2) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem(pmx.KNEE.Right(), fno).Position)
			}
		}
		{
			expectedPosition := &mmath.MVec3{-7.521027, 2.827226, -9.035607}
			if !matrixes.GetItem(pmx.ANKLE.Right(), fno).Position.PracticallyEquals(expectedPosition, 1e-2) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem(pmx.ANKLE.Right(), fno).Position)
			}
		}
		{
			expectedPosition := &mmath.MVec3{-7.453236, 0.356456, -8.876783}
			if !matrixes.GetItem(pmx.TOE.Right(), fno).Position.PracticallyEquals(expectedPosition, 1e-2) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem(pmx.TOE.Right(), fno).Position)
			}
		}
	}
}

func TestVmdMotion_AnimateBoneLegIk6_Snow(t *testing.T) {
	vr := &VmdMotionReader{}
	motionData, err := vr.ReadByFilepath("../../test_resources/好き雪.vmd")

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

	{
		// IK ON
		matrixes := motion.AnimateBone([]int{0}, model, nil, true, false, "")

		{

			fno := 0
			{
				expectedPosition := &mmath.MVec3{2.143878, 6.558880, 1.121747}
				if !matrixes.GetItem(pmx.KNEE.Left(), fno).Position.PracticallyEquals(expectedPosition, 1e-2) {
					t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem(pmx.KNEE.Left(), fno).Position)
				}
			}
			{
				expectedPosition := &mmath.MVec3{2.214143, 1.689811, 2.947619}
				if !matrixes.GetItem(pmx.ANKLE.Left(), fno).Position.PracticallyEquals(expectedPosition, 1e-2) {
					t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem(pmx.ANKLE.Left(), fno).Position)
				}
			}
		}
	}

	{
		// IK OFF
		matrixes := motion.AnimateBone([]int{0}, model, nil, false, false, "")
		{

			fno := 0
			{
				expectedPosition := &mmath.MVec3{1.622245, 6.632885, 0.713205}
				if !matrixes.GetItem(pmx.KNEE.Left(), fno).Position.PracticallyEquals(expectedPosition, 1e-2) {
					t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem(pmx.KNEE.Left(), fno).Position)
				}
			}
			{
				expectedPosition := &mmath.MVec3{1.003185, 1.474691, 0.475763}
				if !matrixes.GetItem(pmx.ANKLE.Left(), fno).Position.PracticallyEquals(expectedPosition, 1e-2) {
					t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem(pmx.ANKLE.Left(), fno).Position)
				}
			}
		}
	}

}

func TestVmdMotion_AnimateBoneLegIk7_Syou(t *testing.T) {
	vr := &VmdMotionReader{}
	motionData, err := vr.ReadByFilepath("../../test_resources/唱(ダンスのみ)_0278F.vmd")

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

	matrixes := motion.AnimateBone([]int{0}, model, []string{pmx.TOE.Right()}, true, false, "")

	// 残存回転判定用
	{

		fno := 0
		{
			expectedPosition := &mmath.MVec3{0.721499, 11.767294, 1.638818}
			if !matrixes.GetItem(pmx.LOWER.String(), fno).Position.PracticallyEquals(expectedPosition, 1e-3) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem(pmx.LOWER.String(), fno).Position)
			}
		}
		{
			expectedPosition := &mmath.MVec3{-0.133304, 10.693992, 2.314730}
			if !matrixes.GetItem(pmx.LEG.Right(), fno).Position.PracticallyEquals(expectedPosition, 1e-3) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem(pmx.LEG.Right(), fno).Position)
			}
		}
		{
			expectedPosition := &mmath.MVec3{-2.833401, 8.174604, -0.100545}
			if !matrixes.GetItem(pmx.KNEE.Right(), fno).Position.PracticallyEquals(expectedPosition, 1e-2) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem(pmx.KNEE.Right(), fno).Position)
			}
		}
		{
			expectedPosition := &mmath.MVec3{-0.409387, 5.341005, 3.524572}
			if !matrixes.GetItem(pmx.ANKLE.Right(), fno).Position.PracticallyEquals(expectedPosition, 1e-2) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem(pmx.ANKLE.Right(), fno).Position)
			}
		}
		{
			expectedPosition := &mmath.MVec3{-0.578271, 2.874233, 3.669599}
			if !matrixes.GetItem(pmx.TOE.Right(), fno).Position.PracticallyEquals(expectedPosition, 1e-2) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem(pmx.TOE.Right(), fno).Position)
			}
		}
	}
}

func TestVmdMotion_AnimateBoneLegIk8_Syou(t *testing.T) {
	vr := &VmdMotionReader{}
	motionData, err := vr.ReadByFilepath("../../test_resources/唱(ダンスのみ)_0-300F.vmd")

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

	matrixes := motion.AnimateBone([]int{278}, model, []string{pmx.TOE.Right()}, true, false, "")

	{

		fno := 278
		{
			expectedPosition := &mmath.MVec3{0.721499, 11.767294, 1.638818}
			if !matrixes.GetItem(pmx.LOWER.String(), fno).Position.PracticallyEquals(expectedPosition, 1e-3) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem(pmx.LOWER.String(), fno).Position)
			}
		}
		{
			expectedPosition := &mmath.MVec3{-0.133304, 10.693992, 2.314730}
			if !matrixes.GetItem(pmx.LEG.Right(), fno).Position.PracticallyEquals(expectedPosition, 1e-3) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem(pmx.LEG.Right(), fno).Position)
			}
		}
		{
			expectedPosition := &mmath.MVec3{-2.833401, 8.174604, -0.100545}
			if !matrixes.GetItem(pmx.KNEE.Right(), fno).Position.PracticallyEquals(expectedPosition, 1e-2) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem(pmx.KNEE.Right(), fno).Position)
			}
		}
		{
			expectedPosition := &mmath.MVec3{-0.409387, 5.341005, 3.524572}
			if !matrixes.GetItem(pmx.ANKLE.Right(), fno).Position.PracticallyEquals(expectedPosition, 1e-2) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem(pmx.ANKLE.Right(), fno).Position)
			}
		}
		{
			expectedPosition := &mmath.MVec3{-0.578271, 2.874233, 3.669599}
			if !matrixes.GetItem(pmx.TOE.Right(), fno).Position.PracticallyEquals(expectedPosition, 1e-2) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem(pmx.TOE.Right(), fno).Position)
			}
		}
	}
}

func TestVmdMotion_AnimateBoneLegIk9_Syou(t *testing.T) {
	vr := &VmdMotionReader{}
	motionData, err := vr.ReadByFilepath("../../test_resources/唱(ダンスのみ)_0-300F.vmd")

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

	matrixes := motion.AnimateBone([]int{274}, model, []string{pmx.TOE.Right()}, true, false, "")

	{

		fno := 274
		{
			expectedPosition := &mmath.MVec3{0.049523, 10.960778, 1.822612}
			if !matrixes.GetItem(pmx.LOWER.String(), fno).Position.PracticallyEquals(expectedPosition, 1e-3) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem(pmx.LOWER.String(), fno).Position)
			}
		}
		{
			expectedPosition := &mmath.MVec3{-0.930675, 9.938401, 2.400088}
			if !matrixes.GetItem(pmx.LEG.Right(), fno).Position.PracticallyEquals(expectedPosition, 1e-3) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem(pmx.LEG.Right(), fno).Position)
			}
		}
		{
			expectedPosition := &mmath.MVec3{-1.710987, 6.669293, -0.459177}
			if !matrixes.GetItem(pmx.KNEE.Right(), fno).Position.PracticallyEquals(expectedPosition, 1e-2) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem(pmx.KNEE.Right(), fno).Position)
			}
		}
		{
			expectedPosition := &mmath.MVec3{-0.773748, 2.387820, 2.340310}
			if !matrixes.GetItem(pmx.ANKLE.Right(), fno).Position.PracticallyEquals(expectedPosition, 1e-2) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem(pmx.ANKLE.Right(), fno).Position)
			}
		}
		{
			expectedPosition := &mmath.MVec3{-1.256876, 0.365575, 0.994345}
			if !matrixes.GetItem(pmx.TOE.Right(), fno).Position.PracticallyEquals(expectedPosition, 1e-2) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem(pmx.TOE.Right(), fno).Position)
			}
		}
	}
}

func TestVmdMotion_AnimateBoneLegIk10_Syou(t *testing.T) {
	vr := &VmdMotionReader{}
	motionData, err := vr.ReadByFilepath("../../test_resources/唱(ダンスのみ)_0-300F.vmd")

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

	matrixes := motion.AnimateBone([]int{100, 107, 272, 273, 274, 275, 278}, model, []string{pmx.TOE.Right()}, true, false, "")

	{

		fno := 100
		{
			expectedPosition := &mmath.MVec3{0.365000, 11.411437, 1.963828}
			if !matrixes.GetItem(pmx.LOWER.String(), fno).Position.PracticallyEquals(expectedPosition, 1e-3) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem(pmx.LOWER.String(), fno).Position)
			}
		}
		{
			expectedPosition := &mmath.MVec3{-0.513678, 10.280550, 2.500991}
			if !matrixes.GetItem(pmx.LEG.Right(), fno).Position.PracticallyEquals(expectedPosition, 1e-3) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem(pmx.LEG.Right(), fno).Position)
			}
		}
		{
			expectedPosition := &mmath.MVec3{-2.891708, 8.162312, -0.553409}
			if !matrixes.GetItem(pmx.KNEE.Right(), fno).Position.PracticallyEquals(expectedPosition, 1e-2) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem(pmx.KNEE.Right(), fno).Position)
			}
		}
		{
			expectedPosition := &mmath.MVec3{-0.826174, 4.330670, 2.292396}
			if !matrixes.GetItem(pmx.ANKLE.Right(), fno).Position.PracticallyEquals(expectedPosition, 1e-2) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem(pmx.ANKLE.Right(), fno).Position)
			}
		}
		{
			expectedPosition := &mmath.MVec3{-1.063101, 1.865613, 2.335564}
			if !matrixes.GetItem(pmx.TOE.Right(), fno).Position.PracticallyEquals(expectedPosition, 1e-2) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem(pmx.TOE.Right(), fno).Position)
			}
		}
	}
	{

		fno := 107
		{
			expectedPosition := &mmath.MVec3{0.365000, 12.042871, 2.034023}
			if !matrixes.GetItem(pmx.LOWER.String(), fno).Position.PracticallyEquals(expectedPosition, 1e-3) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem(pmx.LOWER.String(), fno).Position)
			}
		}
		{
			expectedPosition := &mmath.MVec3{-0.488466, 10.920292, 2.626419}
			if !matrixes.GetItem(pmx.LEG.Right(), fno).Position.PracticallyEquals(expectedPosition, 1e-3) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem(pmx.LEG.Right(), fno).Position)
			}
		}
		{
			expectedPosition := &mmath.MVec3{-1.607765, 6.763937, 1.653586}
			if !matrixes.GetItem(pmx.KNEE.Right(), fno).Position.PracticallyEquals(expectedPosition, 1e-2) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem(pmx.KNEE.Right(), fno).Position)
			}
		}
		{
			expectedPosition := &mmath.MVec3{-1.110289, 1.718307, 2.809817}
			if !matrixes.GetItem(pmx.ANKLE.Right(), fno).Position.PracticallyEquals(expectedPosition, 1e-2) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem(pmx.ANKLE.Right(), fno).Position)
			}
		}
		{
			expectedPosition := &mmath.MVec3{-1.753089, -0.026766, 1.173958}
			if !matrixes.GetItem(pmx.TOE.Right(), fno).Position.PracticallyEquals(expectedPosition, 1e-2) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem(pmx.TOE.Right(), fno).Position)
			}
		}
	}
	{

		fno := 272
		{
			expectedPosition := &mmath.MVec3{-0.330117, 10.811301, 1.914508}
			if !matrixes.GetItem(pmx.LOWER.String(), fno).Position.PracticallyEquals(expectedPosition, 1e-3) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem(pmx.LOWER.String(), fno).Position)
			}
		}
		{
			expectedPosition := &mmath.MVec3{-1.325985, 9.797281, 2.479780}
			if !matrixes.GetItem(pmx.LEG.Right(), fno).Position.PracticallyEquals(expectedPosition, 1e-3) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem(pmx.LEG.Right(), fno).Position)
			}
		}
		{
			expectedPosition := &mmath.MVec3{-1.394679, 6.299243, -0.209150}
			if !matrixes.GetItem(pmx.KNEE.Right(), fno).Position.PracticallyEquals(expectedPosition, 1e-2) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem(pmx.KNEE.Right(), fno).Position)
			}
		}
		{
			expectedPosition := &mmath.MVec3{-0.865021, 1.642431, 2.044760}
			if !matrixes.GetItem(pmx.ANKLE.Right(), fno).Position.PracticallyEquals(expectedPosition, 1e-2) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem(pmx.ANKLE.Right(), fno).Position)
			}
		}
		{
			expectedPosition := &mmath.MVec3{-1.191817, -0.000789, 0.220605}
			if !matrixes.GetItem(pmx.TOE.Right(), fno).Position.PracticallyEquals(expectedPosition, 1e-2) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem(pmx.TOE.Right(), fno).Position)
			}
		}
	}
	{

		fno := 273
		{
			expectedPosition := &mmath.MVec3{-0.154848, 10.862784, 1.868560}
			if !matrixes.GetItem(pmx.LOWER.String(), fno).Position.PracticallyEquals(expectedPosition, 1e-2) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem(pmx.LOWER.String(), fno).Position)
			}
		}
		{
			expectedPosition := &mmath.MVec3{-1.153633, 9.846655, 2.436846}
			if !matrixes.GetItem(pmx.LEG.Right(), fno).Position.PracticallyEquals(expectedPosition, 1e-2) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem(pmx.LEG.Right(), fno).Position)
			}
		}
		{
			expectedPosition := &mmath.MVec3{-1.498977, 6.380789, -0.272370}
			if !matrixes.GetItem(pmx.KNEE.Right(), fno).Position.PracticallyEquals(expectedPosition, 1e-2) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem(pmx.KNEE.Right(), fno).Position)
			}
		}
		{
			expectedPosition := &mmath.MVec3{-0.845777, 1.802650, 2.106815}
			if !matrixes.GetItem(pmx.ANKLE.Right(), fno).Position.PracticallyEquals(expectedPosition, 1e-2) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem(pmx.ANKLE.Right(), fno).Position)
			}
		}
		{
			expectedPosition := &mmath.MVec3{-1.239674, 0.026274, 0.426385}
			if !matrixes.GetItem(pmx.TOE.Right(), fno).Position.PracticallyEquals(expectedPosition, 1e-2) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem(pmx.TOE.Right(), fno).Position)
			}
		}
	}
	{

		fno := 274
		{
			expectedPosition := &mmath.MVec3{0.049523, 10.960778, 1.822612}
			if !matrixes.GetItem(pmx.LOWER.String(), fno).Position.PracticallyEquals(expectedPosition, 1e-3) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem(pmx.LOWER.String(), fno).Position)
			}
		}
		{
			expectedPosition := &mmath.MVec3{-0.930675, 9.938401, 2.400088}
			if !matrixes.GetItem(pmx.LEG.Right(), fno).Position.PracticallyEquals(expectedPosition, 1e-3) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem(pmx.LEG.Right(), fno).Position)
			}
		}
		{
			expectedPosition := &mmath.MVec3{-1.710987, 6.669293, -0.459177}
			if !matrixes.GetItem(pmx.KNEE.Right(), fno).Position.PracticallyEquals(expectedPosition, 1e-2) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem(pmx.KNEE.Right(), fno).Position)
			}
		}
		{
			expectedPosition := &mmath.MVec3{-0.773748, 2.387820, 2.340310}
			if !matrixes.GetItem(pmx.ANKLE.Right(), fno).Position.PracticallyEquals(expectedPosition, 1e-2) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem(pmx.ANKLE.Right(), fno).Position)
			}
		}
		{
			expectedPosition := &mmath.MVec3{-1.256876, 0.365575, 0.994345}
			if !matrixes.GetItem(pmx.TOE.Right(), fno).Position.PracticallyEquals(expectedPosition, 1e-2) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem(pmx.TOE.Right(), fno).Position)
			}
		}
	}
	{

		fno := 278
		{
			expectedPosition := &mmath.MVec3{0.721499, 11.767294, 1.638818}
			if !matrixes.GetItem(pmx.LOWER.String(), fno).Position.PracticallyEquals(expectedPosition, 1e-3) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem(pmx.LOWER.String(), fno).Position)
			}
		}
		{
			expectedPosition := &mmath.MVec3{-0.133304, 10.693992, 2.314730}
			if !matrixes.GetItem(pmx.LEG.Right(), fno).Position.PracticallyEquals(expectedPosition, 1e-3) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem(pmx.LEG.Right(), fno).Position)
			}
		}
		{
			expectedPosition := &mmath.MVec3{-2.833401, 8.174604, -0.100545}
			if !matrixes.GetItem(pmx.KNEE.Right(), fno).Position.PracticallyEquals(expectedPosition, 1e-2) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem(pmx.KNEE.Right(), fno).Position)
			}
		}
		{
			expectedPosition := &mmath.MVec3{-0.409387, 5.341005, 3.524572}
			if !matrixes.GetItem(pmx.ANKLE.Right(), fno).Position.PracticallyEquals(expectedPosition, 1e-2) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem(pmx.ANKLE.Right(), fno).Position)
			}
		}
		{
			expectedPosition := &mmath.MVec3{-0.578271, 2.874233, 3.669599}
			if !matrixes.GetItem(pmx.TOE.Right(), fno).Position.PracticallyEquals(expectedPosition, 1e-2) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem(pmx.TOE.Right(), fno).Position)
			}
		}
	}
	{

		fno := 275
		{
			expectedPosition := &mmath.MVec3{0.271027, 11.113775, 1.776663}
			if !matrixes.GetItem(pmx.LOWER.String(), fno).Position.PracticallyEquals(expectedPosition, 1e-3) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem(pmx.LOWER.String(), fno).Position)
			}
		}
		{
			expectedPosition := &mmath.MVec3{-0.689199, 10.081417, 2.369725}
			if !matrixes.GetItem(pmx.LEG.Right(), fno).Position.PracticallyEquals(expectedPosition, 1e-3) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem(pmx.LEG.Right(), fno).Position)
			}
		}
		{
			expectedPosition := &mmath.MVec3{-1.955139, 7.141531, -0.667679}
			if !matrixes.GetItem(pmx.KNEE.Right(), fno).Position.PracticallyEquals(expectedPosition, 1e-1) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem(pmx.KNEE.Right(), fno).Position)
			}
		}
		{
			expectedPosition := &mmath.MVec3{-0.639503, 3.472883, 2.775674}
			if !matrixes.GetItem(pmx.ANKLE.Right(), fno).Position.PracticallyEquals(expectedPosition, 1e-2) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem(pmx.ANKLE.Right(), fno).Position)
			}
		}
		{
			expectedPosition := &mmath.MVec3{-1.136614, 1.219771, 1.875187}
			if !matrixes.GetItem(pmx.TOE.Right(), fno).Position.PracticallyEquals(expectedPosition, 1e-2) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem(pmx.TOE.Right(), fno).Position)
			}
		}
	}
}

func TestVmdMotion_AnimateBoneLegIk11_Shining_Miku(t *testing.T) {
	vr := &VmdMotionReader{}
	motionData, err := vr.ReadByFilepath("../../test_resources/シャイニングミラクル_50F.vmd")

	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err)
	}

	motion := motionData.(*VmdMotion)

	pr := &pmx.PmxReader{}
	modelData, err := pr.ReadByFilepath("D:/MMD/MikuMikuDance_v926x64/UserFile/Model/VOCALOID/初音ミク/Tda式初音ミク_盗賊つばき流Ｍトレースモデル配布 v1.07/Tda式初音ミク_盗賊つばき流Mトレースモデルv1.07.pmx")

	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err)
	}

	model := modelData.(*pmx.PmxModel)
	model.SetUp()

	matrixes := motion.AnimateBone([]int{0}, model, []string{pmx.LEG_IK.Right(), "足首_R_"}, true, false, "")

	{

		fno := 0
		{
			expectedPosition := &mmath.MVec3{-1.869911, 2.074591, -0.911531}
			if !matrixes.GetItem(pmx.LEG_IK.Right(), fno).Position.PracticallyEquals(expectedPosition, 1e-2) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem(pmx.LEG_IK.Right(), fno).Position)
			}
		}
		{
			expectedPosition := &mmath.MVec3{0.0, 10.142656, -1.362172}
			if !matrixes.GetItem(pmx.LOWER.String(), fno).Position.PracticallyEquals(expectedPosition, 1e-2) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem(pmx.LOWER.String(), fno).Position)
			}
		}
		{
			expectedPosition := &mmath.MVec3{-0.843381, 8.895412, -0.666409}
			if !matrixes.GetItem(pmx.LEG.Right(), fno).Position.PracticallyEquals(expectedPosition, 1e-2) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem(pmx.LEG.Right(), fno).Position)
			}
		}
		{
			expectedPosition := &mmath.MVec3{-0.274925, 5.679991, -4.384042}
			if !matrixes.GetItem(pmx.KNEE.Right(), fno).Position.PracticallyEquals(expectedPosition, 1e-2) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem(pmx.KNEE.Right(), fno).Position)
			}
		}
		{
			expectedPosition := &mmath.MVec3{-1.870632, 2.072767, -0.910016}
			if !matrixes.GetItem(pmx.ANKLE.Right(), fno).Position.PracticallyEquals(expectedPosition, 1e-2) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem(pmx.ANKLE.Right(), fno).Position)
			}
		}
		{
			expectedPosition := &mmath.MVec3{-1.485913, -0.300011, -1.310446}
			if !matrixes.GetItem("足首_R_", fno).Position.PracticallyEquals(expectedPosition, 1e-2) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem("足首_R_", fno).Position)
			}
		}
	}
}

func TestVmdMotion_AnimateBoneLegIk11_Shining_Vroid(t *testing.T) {
	vr := &VmdMotionReader{}
	motionData, err := vr.ReadByFilepath("../../test_resources/シャイニングミラクル_50F.vmd")

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

	matrixes := motion.AnimateBone([]int{0}, model, []string{pmx.TOE.Right()}, true, false, "")

	{

		fno := 0
		{
			expectedPosition := &mmath.MVec3{0.0, 9.379668, -1.051170}
			if !matrixes.GetItem(pmx.LOWER.String(), fno).Position.PracticallyEquals(expectedPosition, 1e-2) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem(pmx.LOWER.String(), fno).Position)
			}
		}
		{
			expectedPosition := &mmath.MVec3{-0.919751, 8.397145, -0.324375}
			if !matrixes.GetItem(pmx.LEG.Right(), fno).Position.PracticallyEquals(expectedPosition, 1e-2) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem(pmx.LEG.Right(), fno).Position)
			}
		}
		{
			expectedPosition := &mmath.MVec3{-0.422861, 6.169319, -4.100779}
			if !matrixes.GetItem(pmx.KNEE.Right(), fno).Position.PracticallyEquals(expectedPosition, 1e-2) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem(pmx.KNEE.Right(), fno).Position)
			}
		}
		{
			expectedPosition := &mmath.MVec3{-1.821804, 2.095607, -1.186269}
			if !matrixes.GetItem(pmx.ANKLE.Right(), fno).Position.PracticallyEquals(expectedPosition, 1e-2) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem(pmx.ANKLE.Right(), fno).Position)
			}
		}
		{
			expectedPosition := &mmath.MVec3{-1.390510, -0.316872, -1.544655}
			if !matrixes.GetItem(pmx.TOE.Right(), fno).Position.PracticallyEquals(expectedPosition, 1e-2) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem(pmx.TOE.Right(), fno).Position)
			}
		}
	}
}

func TestVmdMotion_AnimateBoneLegIk12_Down_Miku(t *testing.T) {
	vr := &VmdMotionReader{}
	motionData, err := vr.ReadByFilepath("../../test_resources/しゃがむ.vmd")

	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err)
	}

	motion := motionData.(*VmdMotion)

	pr := &pmx.PmxReader{}
	modelData, err := pr.ReadByFilepath("D:/MMD/MikuMikuDance_v926x64/UserFile/Model/VOCALOID/初音ミク/Tda式初音ミク_盗賊つばき流Ｍトレースモデル配布 v1.07/Tda式初音ミク_盗賊つばき流Mトレースモデルv1.07.pmx")

	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err)
	}

	model := modelData.(*pmx.PmxModel)
	model.SetUp()

	matrixes := motion.AnimateBone([]int{0}, model, []string{pmx.LEG_IK.Right(), "足首_R_"}, true, false, "")

	{

		fno := 0
		{
			expectedPosition := &mmath.MVec3{-1.012964, 1.623157, 0.680305}
			if !matrixes.GetItem(pmx.LEG_IK.Right(), fno).Position.PracticallyEquals(expectedPosition, 1e-2) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem(pmx.LEG_IK.Right(), fno).Position)
			}
		}
		{
			expectedPosition := &mmath.MVec3{0.0, 5.953951, -0.512170}
			if !matrixes.GetItem(pmx.LOWER.String(), fno).Position.PracticallyEquals(expectedPosition, 1e-2) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem(pmx.LOWER.String(), fno).Position)
			}
		}
		{
			expectedPosition := &mmath.MVec3{-0.896440, 4.569404, -0.337760}
			if !matrixes.GetItem(pmx.LEG.Right(), fno).Position.PracticallyEquals(expectedPosition, 1e-2) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem(pmx.LEG.Right(), fno).Position)
			}
		}
		{
			expectedPosition := &mmath.MVec3{-0.691207, 1.986888, -4.553376}
			if !matrixes.GetItem(pmx.KNEE.Right(), fno).Position.PracticallyEquals(expectedPosition, 1e-2) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem(pmx.KNEE.Right(), fno).Position)
			}
		}
		{
			expectedPosition := &mmath.MVec3{-1.012964, 1.623157, 0.680305}
			if !matrixes.GetItem(pmx.ANKLE.Right(), fno).Position.PracticallyEquals(expectedPosition, 1e-2) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem(pmx.ANKLE.Right(), fno).Position)
			}
		}
		{
			expectedPosition := &mmath.MVec3{-1.013000, 0.002578, -1.146909}
			if !matrixes.GetItem("足首_R_", fno).Position.PracticallyEquals(expectedPosition, 1e-2) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem("足首_R_", fno).Position)
			}
		}
	}
}
