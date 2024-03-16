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

	{

		fno := float32(10.0)
		matrixes := motion.AnimateBone(fno, model, []string{pmx.INDEX3.Left()}, false, false, "")
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
		fno := float32(999)
		matrixes := motion.AnimateBone(fno, model, []string{pmx.INDEX3.Left()}, false, false, "")
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

	{
		fno := float32(29)
		matrixes := motion.AnimateBone(fno, model, []string{pmx.TOE.Left()}, true, false, "")
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

	{

		fno := float32(3152)
		matrixes := motion.AnimateBone(fno, model, []string{pmx.TOE.Left()}, true, false, "")
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

	{

		fno := float32(60)
		matrixes := motion.AnimateBone(fno, model, nil, true, false, "")
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

	{

		fno := float32(0)
		matrixes := motion.AnimateBone(fno, model, []string{pmx.TOE.Right()}, true, false, "")
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

	{

		fno := float32(7409)
		matrixes := motion.AnimateBone(fno, model, []string{pmx.TOE.Right()}, true, false, "")
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
		{

			fno := float32(0)
			matrixes := motion.AnimateBone(fno, model, nil, true, false, "")
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
		{

			fno := float32(0)
			matrixes := motion.AnimateBone(fno, model, nil, false, false, "")
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

	// 残存回転判定用
	{

		fno := float32(0)
		matrixes := motion.AnimateBone(fno, model, []string{pmx.TOE.Right()}, true, false, "")
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

	{

		fno := float32(278)
		matrixes := motion.AnimateBone(fno, model, []string{pmx.TOE.Right()}, true, false, "")
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

	{

		fno := float32(100)
		matrixes := motion.AnimateBone(fno, model, []string{pmx.TOE.Right()}, true, false, "")
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

		fno := float32(107)
		matrixes := motion.AnimateBone(fno, model, []string{pmx.TOE.Right()}, true, false, "")
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

		fno := float32(272)
		matrixes := motion.AnimateBone(fno, model, []string{pmx.TOE.Right()}, true, false, "")
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

		fno := float32(273)
		matrixes := motion.AnimateBone(fno, model, []string{pmx.TOE.Right()}, true, false, "")
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

		fno := float32(274)
		matrixes := motion.AnimateBone(fno, model, []string{pmx.TOE.Right()}, true, false, "")
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

		fno := float32(278)
		matrixes := motion.AnimateBone(fno, model, []string{pmx.TOE.Right()}, true, false, "")
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

		fno := float32(275)
		matrixes := motion.AnimateBone(fno, model, []string{pmx.TOE.Right()}, true, false, "")
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

	{

		fno := float32(0)
		matrixes := motion.AnimateBone(fno, model, []string{pmx.LEG_IK.Right(), "足首_R_"}, true, false, "")
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

	{

		fno := float32(0)
		matrixes := motion.AnimateBone(fno, model, []string{pmx.TOE.Right()}, true, false, "")
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

	{

		fno := float32(0)
		matrixes := motion.AnimateBone(fno, model, []string{pmx.LEG_IK.Right(), "足首_R_"}, true, false, "")
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

func TestVmdMotion_AnimateBoneLegIk13_Lamb(t *testing.T) {
	vr := &VmdMotionReader{}
	motionData, err := vr.ReadByFilepath("../../test_resources/Lamb_2689F.vmd")

	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err)
	}

	motion := motionData.(*VmdMotion)

	pr := &pmx.PmxReader{}
	modelData, err := pr.ReadByFilepath("D:/MMD/MikuMikuDance_v926x64/UserFile/Model/ゲーム/戦国BASARA/幸村 たぬき式 ver.1.24/真田幸村没第二衣装1.24軽量版.pmx")

	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err)
	}

	model := modelData.(*pmx.PmxModel)
	model.SetUp()

	matrixes := motion.AnimateBone(0, model,
		[]string{pmx.LEG_IK.Right(), pmx.TOE.Right(), pmx.LEG_IK.Left(), pmx.TOE.Left()}, true, false, "")

	{

		fno := float32(0)
		{
			expectedPosition := &mmath.MVec3{-1.216134, 1.887670, -10.78867}
			if !matrixes.GetItem(pmx.LEG_IK.Right(), fno).Position.PracticallyEquals(expectedPosition, 1e-2) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem(pmx.LEG_IK.Right(), fno).Position)
			}
		}
		{
			expectedPosition := &mmath.MVec3{0.803149, 6.056844, -10.232766}
			if !matrixes.GetItem(pmx.LOWER.String(), fno).Position.PracticallyEquals(expectedPosition, 1e-2) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem(pmx.LOWER.String(), fno).Position)
			}
		}
		{
			expectedPosition := &mmath.MVec3{0.728442, 4.560226, -11.571869}
			if !matrixes.GetItem(pmx.LEG.Right(), fno).Position.PracticallyEquals(expectedPosition, 1e-2) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem(pmx.LEG.Right(), fno).Position)
			}
		}
		{
			expectedPosition := &mmath.MVec3{4.173470, 0.361388, -11.217197}
			if !matrixes.GetItem(pmx.KNEE.Right(), fno).Position.PracticallyEquals(expectedPosition, 1e-2) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem(pmx.KNEE.Right(), fno).Position)
			}
		}
		{
			expectedPosition := &mmath.MVec3{-1.217569, 1.885731, -10.788104}
			if !matrixes.GetItem(pmx.ANKLE.Right(), fno).Position.PracticallyEquals(expectedPosition, 1e-2) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem(pmx.ANKLE.Right(), fno).Position)
			}
		}
		{
			expectedPosition := &mmath.MVec3{-0.922247, -1.163554, -10.794323}
			if !matrixes.GetItem(pmx.TOE.Right(), fno).Position.PracticallyEquals(expectedPosition, 1e-2) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem(pmx.TOE.Right(), fno).Position)
			}
		}
	}
	{

		fno := float32(0)
		{
			expectedPosition := &mmath.MVec3{2.322227, 1.150214, -9.644499}
			if !matrixes.GetItem(pmx.LEG_IK.Left(), fno).Position.PracticallyEquals(expectedPosition, 1e-2) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem(pmx.LEG_IK.Left(), fno).Position)
			}
		}
		{
			expectedPosition := &mmath.MVec3{0.803149, 6.056844, -10.232766}
			if !matrixes.GetItem(pmx.LOWER.String(), fno).Position.PracticallyEquals(expectedPosition, 1e-2) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem(pmx.LOWER.String(), fno).Position)
			}
		}
		{
			expectedPosition := &mmath.MVec3{0.720821, 4.639688, -8.810255}
			if !matrixes.GetItem(pmx.LEG.Left(), fno).Position.PracticallyEquals(expectedPosition, 1e-2) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem(pmx.LEG.Left(), fno).Position)
			}
		}
		{
			expectedPosition := &mmath.MVec3{6.126388, 5.074682, -8.346903}
			if !matrixes.GetItem(pmx.KNEE.Left(), fno).Position.PracticallyEquals(expectedPosition, 1e-2) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem(pmx.KNEE.Left(), fno).Position)
			}
		}
		{
			expectedPosition := &mmath.MVec3{2.323599, 1.147291, -9.645196}
			if !matrixes.GetItem(pmx.ANKLE.Left(), fno).Position.PracticallyEquals(expectedPosition, 1e-2) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem(pmx.ANKLE.Left(), fno).Position)
			}
		}
		{
			expectedPosition := &mmath.MVec3{5.163002, -0.000894, -9.714369}
			if !matrixes.GetItem(pmx.TOE.Left(), fno).Position.PracticallyEquals(expectedPosition, 1e-2) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem(pmx.TOE.Left(), fno).Position)
			}
		}
	}
}

func TestVmdMotion_AnimateBoneLegIk14_Ballet(t *testing.T) {
	vr := &VmdMotionReader{}
	motionData, err := vr.ReadByFilepath("../../test_resources/ミク用バレリーコ_1069.vmd")

	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err)
	}

	motion := motionData.(*VmdMotion)

	pr := &pmx.PmxReader{}
	modelData, err := pr.ReadByFilepath("D:/MMD/MikuMikuDance_v926x64/UserFile/Model/_あにまさ式/初音ミク_準標準.pmx")

	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err)
	}

	model := modelData.(*pmx.PmxModel)
	model.SetUp()

	{

		fno := float32(0)
		matrixes := motion.AnimateBone(fno, model, []string{pmx.LEG_IK.Right(), pmx.TOE.Right()}, true, false, "")
		{
			expectedPosition := &mmath.MVec3{11.324574, 10.920002, -7.150005}
			if !matrixes.GetItem(pmx.LEG_IK.Right(), fno).Position.PracticallyEquals(expectedPosition, 1e-2) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem(pmx.LEG_IK.Right(), fno).Position)
			}
		}
		{
			expectedPosition := &mmath.MVec3{2.433170, 13.740387, 0.992719}
			if !matrixes.GetItem(pmx.LOWER.String(), fno).Position.PracticallyEquals(expectedPosition, 1e-2) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem(pmx.LOWER.String(), fno).Position)
			}
		}
		{
			expectedPosition := &mmath.MVec3{1.982654, 11.188538, 0.602013}
			if !matrixes.GetItem(pmx.LEG.Right(), fno).Position.PracticallyEquals(expectedPosition, 1e-2) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem(pmx.LEG.Right(), fno).Position)
			}
		}
		{
			expectedPosition := &mmath.MVec3{5.661557, 11.008962, -2.259013}
			if !matrixes.GetItem(pmx.KNEE.Right(), fno).Position.PracticallyEquals(expectedPosition, 1e-1) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem(pmx.KNEE.Right(), fno).Position)
			}
		}
		{
			expectedPosition := &mmath.MVec3{9.224476, 10.979847, -5.407887}
			if !matrixes.GetItem(pmx.ANKLE.Right(), fno).Position.PracticallyEquals(expectedPosition, 1e-2) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem(pmx.ANKLE.Right(), fno).Position)
			}
		}
		{
			expectedPosition := &mmath.MVec3{11.345482, 10.263426, -7.003638}
			if !matrixes.GetItem(pmx.TOE.Right(), fno).Position.PracticallyEquals(expectedPosition, 1e-2) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem(pmx.TOE.Right(), fno).Position)
			}
		}
	}
}

func TestVmdMotion_AnimateBoneLegIk15_Bottom(t *testing.T) {
	vr := &VmdMotionReader{}
	motionData, err := vr.ReadByFilepath("../../test_resources/●ボトム_0-300.vmd")

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

	{

		fno := float32(218)
		matrixes := motion.AnimateBone(fno, model, []string{pmx.LEG_IK.Right(), "足首_R_"}, true, false, "")
		{
			expectedPosition := &mmath.MVec3{-1.358434, 1.913062, 0.611182}
			if !matrixes.GetItem(pmx.LEG_IK.Right(), fno).Position.PracticallyEquals(expectedPosition, 1e-2) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem(pmx.LEG_IK.Right(), fno).Position)
			}
		}
		{
			expectedPosition := &mmath.MVec3{0.150000, 4.253955, 0.237829}
			if !matrixes.GetItem(pmx.LOWER.String(), fno).Position.PracticallyEquals(expectedPosition, 1e-2) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem(pmx.LOWER.String(), fno).Position)
			}
		}
		{
			expectedPosition := &mmath.MVec3{-0.906292, 2.996784, 0.471846}
			if !matrixes.GetItem(pmx.LEG.Right(), fno).Position.PracticallyEquals(expectedPosition, 1e-2) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem(pmx.LEG.Right(), fno).Position)
			}
		}
		{
			expectedPosition := &mmath.MVec3{-2.533418, 3.889916, -4.114837}
			if !matrixes.GetItem(pmx.KNEE.Right(), fno).Position.PracticallyEquals(expectedPosition, 1e-2) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem(pmx.KNEE.Right(), fno).Position)
			}
		}
		{
			expectedPosition := &mmath.MVec3{-1.358807, 1.912181, 0.611265}
			if !matrixes.GetItem(pmx.ANKLE.Right(), fno).Position.PracticallyEquals(expectedPosition, 1e-2) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem(pmx.ANKLE.Right(), fno).Position)
			}
		}
		{
			expectedPosition := &mmath.MVec3{-2.040872, -0.188916, -0.430442}
			if !matrixes.GetItem("足首_R_", fno).Position.PracticallyEquals(expectedPosition, 1e-2) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem("足首_R_", fno).Position)
			}
		}
	}
}

func TestVmdMotion_AnimateBoneLegIk16_Lamb(t *testing.T) {
	vr := &VmdMotionReader{}
	motionData, err := vr.ReadByFilepath("../../test_resources/Lamb_2689F.vmd")

	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err)
	}

	motion := motionData.(*VmdMotion)

	pr := &pmx.PmxReader{}
	modelData, err := pr.ReadByFilepath("D:/MMD/MikuMikuDance_v926x64/UserFile/Model/ゲーム/戦国BASARA/幸村 たぬき式 ver.1.24/真田幸村没第二衣装1.24軽量版.pmx")

	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err)
	}

	model := modelData.(*pmx.PmxModel)
	model.SetUp()

	{

		fno := float32(0)
		matrixes := motion.AnimateBone(fno, model, []string{pmx.LEG_IK.Right(), pmx.TOE.Right()}, true, false, "")
		{
			expectedPosition := &mmath.MVec3{-1.216134, 1.887670, -10.78867}
			if !matrixes.GetItem(pmx.LEG_IK.Right(), fno).Position.PracticallyEquals(expectedPosition, 1e-2) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem(pmx.LEG_IK.Right(), fno).Position)
			}
		}
		{
			expectedPosition := &mmath.MVec3{0.803149, 6.056844, -10.232766}
			if !matrixes.GetItem(pmx.LOWER.String(), fno).Position.PracticallyEquals(expectedPosition, 1e-2) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem(pmx.LOWER.String(), fno).Position)
			}
		}
		{
			expectedPosition := &mmath.MVec3{0.728442, 4.560226, -11.571869}
			if !matrixes.GetItem(pmx.LEG.Right(), fno).Position.PracticallyEquals(expectedPosition, 1e-2) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem(pmx.LEG.Right(), fno).Position)
			}
		}
		{
			expectedPosition := &mmath.MVec3{4.173470, 0.361388, -11.217197}
			if !matrixes.GetItem(pmx.KNEE.Right(), fno).Position.PracticallyEquals(expectedPosition, 1e-2) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem(pmx.KNEE.Right(), fno).Position)
			}
		}
		{
			expectedPosition := &mmath.MVec3{-1.217569, 1.885731, -10.788104}
			if !matrixes.GetItem(pmx.ANKLE.Right(), fno).Position.PracticallyEquals(expectedPosition, 1e-2) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem(pmx.ANKLE.Right(), fno).Position)
			}
		}
		{
			expectedPosition := &mmath.MVec3{-0.922247, -1.163554, -10.794323}
			if !matrixes.GetItem(pmx.TOE.Right(), fno).Position.PracticallyEquals(expectedPosition, 1e-2) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem(pmx.TOE.Right(), fno).Position)
			}
		}
	}
}

func TestVmdMotion_AnimateBoneLegIk17_Snow(t *testing.T) {
	vr := &VmdMotionReader{}
	motionData, err := vr.ReadByFilepath("../../test_resources/好き雪_1075.vmd")

	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err)
	}

	motion := motionData.(*VmdMotion)

	pr := &pmx.PmxReader{}
	modelData, err := pr.ReadByFilepath("D:/MMD/MikuMikuDance_v926x64/UserFile/Model/VOCALOID/初音ミク/Lat式ミクVer2.31/Lat式ミクVer2.31_White_準標準.pmx")

	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err)
	}

	model := modelData.(*pmx.PmxModel)
	model.SetUp()

	{

		fno := float32(0)
		matrixes := motion.AnimateBone(fno, model, []string{pmx.TOE.Right()}, true, false, "")
		{
			expectedPosition := &mmath.MVec3{2.049998, 12.957623, 1.477440}
			if !matrixes.GetItem(pmx.LOWER.String(), fno).Position.PracticallyEquals(expectedPosition, 1e-2) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem(pmx.LOWER.String(), fno).Position)
			}
		}
		{
			expectedPosition := &mmath.MVec3{1.201382, 11.353215, 2.266898}
			if !matrixes.GetItem(pmx.LEG.Right(), fno).Position.PracticallyEquals(expectedPosition, 1e-2) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem(pmx.LEG.Right(), fno).Position)
			}
		}
		{
			expectedPosition := &mmath.MVec3{0.443043, 7.640018, -1.308741}
			if !matrixes.GetItem(pmx.KNEE.Right(), fno).Position.PracticallyEquals(expectedPosition, 1e-2) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem(pmx.KNEE.Right(), fno).Position)
			}
		}
		{
			expectedPosition := &mmath.MVec3{-0.574753, 7.943915, 3.279809}
			if !matrixes.GetItem(pmx.ANKLE.Right(), fno).Position.PracticallyEquals(expectedPosition, 1e-2) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem(pmx.ANKLE.Right(), fno).Position)
			}
		}
		{
			expectedPosition := &mmath.MVec3{-1.443098, 6.324932, 4.887177}
			if !matrixes.GetItem(pmx.TOE.Right(), fno).Position.PracticallyEquals(expectedPosition, 1e-1) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem(pmx.TOE.Right(), fno).Position)
			}
		}
	}
}

func TestVmdMotion_AnimateBoneLegIk18_Syou(t *testing.T) {
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

	{

		fno := float32(107)
		matrixes := motion.AnimateBone(fno, model, []string{pmx.TOE.Right()}, true, false, "")
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
}

func TestVmdMotion_AnimateBoneLegIk19_Wa(t *testing.T) {
	vr := &VmdMotionReader{}
	motionData, err := vr.ReadByFilepath("../../test_resources/129cm_001_10F.vmd")

	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err)
	}

	motion := motionData.(*VmdMotion)

	pr := &pmx.PmxReader{}
	modelData, err := pr.ReadByFilepath("D:/MMD/MikuMikuDance_v926x64/UserFile/Model/_VMDサイジング/wa_129cm 20231028/wa_129cm_bone-structure.pmx")

	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err)
	}

	model := modelData.(*pmx.PmxModel)
	model.SetUp()

	{

		fno := float32(0)
		matrixes := motion.AnimateBone(fno, model, []string{pmx.TOE.Right()}, true, false, "")
		{
			expectedPosition := &mmath.MVec3{0.000000, 9.900000, 0.000000}
			if !matrixes.GetItem(pmx.LOWER.String(), fno).Position.PracticallyEquals(expectedPosition, 1e-3) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem(pmx.LOWER.String(), fno).Position)
			}
		}
		{
			expectedPosition := &mmath.MVec3{-0.599319, 8.639606, 0.369618}
			if !matrixes.GetItem(pmx.LEG.Right(), fno).Position.PracticallyEquals(expectedPosition, 1e-3) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem(pmx.LEG.Right(), fno).Position)
			}
		}
		{
			expectedPosition := &mmath.MVec3{-0.486516, 6.323577, -2.217865}
			if !matrixes.GetItem(pmx.KNEE.Right(), fno).Position.PracticallyEquals(expectedPosition, 1e-2) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem(pmx.KNEE.Right(), fno).Position)
			}
		}
		{
			expectedPosition := &mmath.MVec3{-2.501665, 2.859252, -1.902513}
			if !matrixes.GetItem(pmx.ANKLE.Right(), fno).Position.PracticallyEquals(expectedPosition, 1e-2) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem(pmx.ANKLE.Right(), fno).Position)
			}
		}
		{
			expectedPosition := &mmath.MVec3{-3.071062, 0.841962, -2.077063}
			if !matrixes.GetItem(pmx.TOE.Right(), fno).Position.PracticallyEquals(expectedPosition, 1e-2) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem(pmx.TOE.Right(), fno).Position)
			}
		}
	}
}

func TestVmdMotion_AnimateBoneLegIk20_Syou(t *testing.T) {
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

	{

		fno := float32(107)
		matrixes := motion.AnimateBone(fno, model, []string{pmx.TOE.Right()}, true, false, "")
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
}

func TestVmdMotion_AnimateBoneLegIk21_FK(t *testing.T) {
	vr := &VmdMotionReader{}
	motionData, err := vr.ReadByFilepath("../../test_resources/足FK.vmd")

	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err)
	}

	motion := motionData.(*VmdMotion)

	pr := &pmx.PmxReader{}
	modelData, err := pr.ReadByFilepath("../../test_resources/サンプルモデル_ひざ制限なし.pmx")

	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err)
	}

	model := modelData.(*pmx.PmxModel)
	model.SetUp()

	{

		fno := float32(0)
		matrixes := motion.AnimateBone(fno, model, []string{pmx.TOE.Right()}, false, false, "")
		{
			expectedPosition := &mmath.MVec3{-0.133305, 10.693993, 2.314730}
			if !matrixes.GetItem(pmx.LEG.Right(), fno).Position.PracticallyEquals(expectedPosition, 1e-3) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem(pmx.LEG.Right(), fno).Position)
			}
		}
		{
			expectedPosition := &mmath.MVec3{2.708069, 9.216356, -0.720822}
			if !matrixes.GetItem(pmx.KNEE.Right(), fno).Position.PracticallyEquals(expectedPosition, 1e-2) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem(pmx.KNEE.Right(), fno).Position)
			}
		}
	}
}

func TestVmdMotion_AnimateBoneLegIk22_Bake(t *testing.T) {
	vr := &VmdMotionReader{}
	motionData, err := vr.ReadByFilepath("../../test_resources/足FK焼き込み.vmd")

	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err)
	}

	motion := motionData.(*VmdMotion)

	pr := &pmx.PmxReader{}
	modelData, err := pr.ReadByFilepath("../../test_resources/サンプルモデル_ひざ制限なし.pmx")

	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err)
	}

	model := modelData.(*pmx.PmxModel)
	model.SetUp()

	{

		fno := float32(0)
		matrixes := motion.AnimateBone(fno, model, []string{pmx.TOE.Right()}, true, false, "")
		{
			expectedPosition := &mmath.MVec3{-0.133306, 10.693994, 2.314731}
			if !matrixes.GetItem(pmx.LEG.Right(), fno).Position.PracticallyEquals(expectedPosition, 1e-3) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem(pmx.LEG.Right(), fno).Position)
			}
		}
		{
			expectedPosition := &mmath.MVec3{-3.753989, 8.506582, 1.058842}
			if !matrixes.GetItem(pmx.KNEE.Right(), fno).Position.PracticallyEquals(expectedPosition, 1e-1) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem(pmx.KNEE.Right(), fno).Position)
			}
		}
	}
}

func TestVmdMotion_AnimateBoneLegIk22_NoLimit(t *testing.T) {
	vr := &VmdMotionReader{}
	motionData, err := vr.ReadByFilepath("../../test_resources/足FK.vmd")

	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err)
	}

	motion := motionData.(*VmdMotion)

	pr := &pmx.PmxReader{}
	modelData, err := pr.ReadByFilepath("../../test_resources/サンプルモデル_ひざ制限なし.pmx")

	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err)
	}

	model := modelData.(*pmx.PmxModel)
	model.SetUp()

	{

		fno := float32(0)
		matrixes := motion.AnimateBone(fno, model, []string{pmx.TOE.Right()}, true, false, "")
		{
			expectedPosition := &mmath.MVec3{-0.133305, 10.693993, 2.314730}
			if !matrixes.GetItem(pmx.LEG.Right(), fno).Position.PracticallyEquals(expectedPosition, 1e-3) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem(pmx.LEG.Right(), fno).Position)
			}
		}
		{
			expectedPosition := &mmath.MVec3{2.081436, 7.884178, -0.268146}
			if !matrixes.GetItem(pmx.KNEE.Right(), fno).Position.PracticallyEquals(expectedPosition, 1e-2) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem(pmx.KNEE.Right(), fno).Position)
			}
		}
	}
}

func TestVmdMotion_AnimateBoneArmIk(t *testing.T) {
	vr := &VmdMotionReader{}
	motionData, err := vr.ReadByFilepath("../../test_resources/サンプルモーション.vmd")

	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err)
	}

	motion := motionData.(*VmdMotion)

	pr := &pmx.PmxReader{}
	modelData, err := pr.ReadByFilepath("../../test_resources/ボーンツリーテストモデル.pmx")

	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err)
	}

	model := modelData.(*pmx.PmxModel)
	model.SetUp()

	{

		fno := float32(3182)
		matrixes := motion.AnimateBone(fno, model, nil, true, false, "")
		{
			expectedPosition := &mmath.MVec3{0, 0, 0}
			if !matrixes.GetItem(pmx.ROOT.String(), fno).Position.PracticallyEquals(expectedPosition, 1e-3) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem(pmx.ROOT.String(), fno).Position)
			}
		}
		{
			expectedPosition := &mmath.MVec3{12.400011, 9.000000, 1.885650}
			if !matrixes.GetItem(pmx.CENTER.String(), fno).Position.PracticallyEquals(expectedPosition, 1e-2) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem(pmx.CENTER.String(), fno).Position)
			}
		}
		{
			expectedPosition := &mmath.MVec3{12.400011, 8.580067, 1.885650}
			if !matrixes.GetItem(pmx.GROOVE.String(), fno).Position.PracticallyEquals(expectedPosition, 1e-2) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem(pmx.GROOVE.String(), fno).Position)
			}
		}
		{
			expectedPosition := &mmath.MVec3{12.400011, 11.628636, 2.453597}
			if !matrixes.GetItem(pmx.WAIST.String(), fno).Position.PracticallyEquals(expectedPosition, 1e-2) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem(pmx.WAIST.String(), fno).Position)
			}
		}
		{
			expectedPosition := &mmath.MVec3{12.400011, 12.567377, 1.229520}
			if !matrixes.GetItem(pmx.UPPER.String(), fno).Position.PracticallyEquals(expectedPosition, 1e-2) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem(pmx.UPPER.String(), fno).Position)
			}
		}
		{
			expectedPosition := &mmath.MVec3{12.344202, 13.782951, 1.178849}
			if !matrixes.GetItem(pmx.UPPER2.String(), fno).Position.PracticallyEquals(expectedPosition, 1e-2) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem(pmx.UPPER2.String(), fno).Position)
			}
		}
		{
			expectedPosition := &mmath.MVec3{12.425960, 15.893852, 1.481421}
			if !matrixes.GetItem(pmx.SHOULDER_P.Left(), fno).Position.PracticallyEquals(expectedPosition, 1e-2) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem(pmx.SHOULDER_P.Left(), fno).Position)
			}
		}
		{
			expectedPosition := &mmath.MVec3{12.425960, 15.893852, 1.481421}
			if !matrixes.GetItem(pmx.SHOULDER.Left(), fno).Position.PracticallyEquals(expectedPosition, 1e-2) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem(pmx.SHOULDER.Left(), fno).Position)
			}
		}
		{
			expectedPosition := &mmath.MVec3{13.348320, 15.767927, 1.802947}
			if !matrixes.GetItem(pmx.ARM.Left(), fno).Position.PracticallyEquals(expectedPosition, 1e-2) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem(pmx.ARM.Left(), fno).Position)
			}
		}
		{
			expectedPosition := &mmath.MVec3{13.564770, 14.998386, 1.289923}
			if !matrixes.GetItem(pmx.ARM_TWIST.Left(), fno).Position.PracticallyEquals(expectedPosition, 1e-2) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem(pmx.ARM_TWIST.Left(), fno).Position)
			}
		}
		{
			expectedPosition := &mmath.MVec3{14.043257, 13.297290, 0.155864}
			if !matrixes.GetItem(pmx.ELBOW.Left(), fno).Position.PracticallyEquals(expectedPosition, 1e-2) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem(pmx.ELBOW.Left(), fno).Position)
			}
		}
		{
			expectedPosition := &mmath.MVec3{13.811955, 13.552182, -0.388005}
			if !matrixes.GetItem(pmx.WRIST_TWIST.Left(), fno).Position.PracticallyEquals(expectedPosition, 1e-2) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem(pmx.WRIST_TWIST.Left(), fno).Position)
			}
		}
		{
			expectedPosition := &mmath.MVec3{13.144803, 14.287374, -1.956703}
			if !matrixes.GetItem(pmx.WRIST.Left(), fno).Position.PracticallyEquals(expectedPosition, 1e-2) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem(pmx.WRIST.Left(), fno).Position)
			}
		}
		{
			expectedPosition := &mmath.MVec3{12.813587, 14.873419, -2.570278}
			if !matrixes.GetItem(pmx.INDEX1.Left(), fno).Position.PracticallyEquals(expectedPosition, 1e-2) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem(pmx.INDEX1.Left(), fno).Position)
			}
		}
		{
			expectedPosition := &mmath.MVec3{12.541822, 15.029200, -2.709604}
			if !matrixes.GetItem(pmx.INDEX2.Left(), fno).Position.PracticallyEquals(expectedPosition, 1e-2) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem(pmx.INDEX2.Left(), fno).Position)
			}
		}
		{
			expectedPosition := &mmath.MVec3{12.476499, 14.950351, -2.502167}
			if !matrixes.GetItem(pmx.INDEX3.Left(), fno).Position.PracticallyEquals(expectedPosition, 1e-2) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem(pmx.INDEX3.Left(), fno).Position)
			}
		}
		{
			expectedPosition := &mmath.MVec3{12.620306, 14.795185, -2.295859}
			if !matrixes.GetItem("左人指先", fno).Position.PracticallyEquals(expectedPosition, 1e-2) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem("左人指先", fno).Position)
			}
		}
	}
}

func TestVmdMotion_AnimateBoneLegIk2(t *testing.T) {
	vr := &VmdMotionReader{}
	motionData, err := vr.ReadByFilepath("C:/MMD/mmd_base/tests/resources/唱(ダンスのみ)_0274F.vmd")

	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err)
	}

	motion := motionData.(*VmdMotion)

	pr := &pmx.PmxReader{}
	modelData, err := pr.ReadByFilepath("D:/MMD/MikuMikuDance_v926x64/UserFile/Model/VOCALOID/初音ミク/ISAO式ミク/I_ミクv4/Miku_V4_準標準.pmx")

	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err)
	}

	model := modelData.(*pmx.PmxModel)
	model.SetUp()

	{

		fno := float32(0)
		matrixes := motion.AnimateBone(fno, model, nil, true, false, "")
		{
			expectedPosition := &mmath.MVec3{0.04952335, 9.0, 1.72378033}
			if !matrixes.GetItem(pmx.CENTER.String(), fno).Position.PracticallyEquals(expectedPosition, 1e-3) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem(pmx.CENTER.String(), fno).Position)
			}
		}
		{
			expectedPosition := &mmath.MVec3{0.04952335, 7.97980869, 1.72378033}
			if !matrixes.GetItem(pmx.GROOVE.String(), fno).Position.PracticallyEquals(expectedPosition, 1e-1) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem(pmx.GROOVE.String(), fno).Position)
			}
		}
		{
			expectedPosition := &mmath.MVec3{0.04952335, 11.02838314, 2.29172656}
			if !matrixes.GetItem("腰", fno).Position.PracticallyEquals(expectedPosition, 1e-1) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem("腰", fno).Position)
			}
		}
		{
			expectedPosition := &mmath.MVec3{0.04952335, 11.9671191, 1.06765032}
			if !matrixes.GetItem("下半身", fno).Position.PracticallyEquals(expectedPosition, 1e-1) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem("下半身", fno).Position)
			}
		}
		{
			expectedPosition := &mmath.MVec3{-0.24102019, 9.79926074, 1.08498769}
			if !matrixes.GetItem("下半身先", fno).Position.PracticallyEquals(expectedPosition, 1e-1) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem("下半身先", fno).Position)
			}
		}
		{
			expectedPosition := &mmath.MVec3{0.90331914, 10.27362702, 1.009499759}
			if !matrixes.GetItem("腰キャンセル左", fno).Position.PracticallyEquals(expectedPosition, 1e-1) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem("腰キャンセル左", fno).Position)
			}
		}
		{
			expectedPosition := &mmath.MVec3{0.90331914, 10.27362702, 1.00949975}
			if !matrixes.GetItem("左足", fno).Position.PracticallyEquals(expectedPosition, 1e-1) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem("左足", fno).Position)
			}
		}
		{
			expectedPosition := &mmath.MVec3{0.08276818, 5.59348757, -1.24981795}
			if !matrixes.GetItem("左ひざ", fno).Position.PracticallyEquals(expectedPosition, 1e-1) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem("左ひざ", fno).Position)
			}
		}
		{
			expectedPosition := &mmath.MVec3{5.63290634e-01, -2.12439821e-04, -3.87768478e-01}
			if !matrixes.GetItem("左つま先", fno).Position.PracticallyEquals(expectedPosition, 1e-1) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem("左つま先", fno).Position)
			}
		}
		{
			expectedPosition := &mmath.MVec3{0.90331914, 10.27362702, 1.00949975}
			if !matrixes.GetItem("左足D", fno).Position.PracticallyEquals(expectedPosition, 1e-1) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem("左足D", fno).Position)
			}
		}
		{
			expectedPosition := &mmath.MVec3{0.23453057, 5.6736954, -0.76228439}
			if !matrixes.GetItem("左ひざ2", fno).Position.PracticallyEquals(expectedPosition, 1e-1) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem("左ひざ2", fno).Position)
			}
		}
		{
			expectedPosition := &mmath.MVec3{0.12060311, 4.95396153, -1.23761938}
			if !matrixes.GetItem("左ひざ2先", fno).Position.PracticallyEquals(expectedPosition, 1e-1) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem("左ひざ2先", fno).Position)
			}
		}
		{
			expectedPosition := &mmath.MVec3{0.90331914, 10.27362702, 1.00949975}
			if !matrixes.GetItem("左足y+", fno).Position.PracticallyEquals(expectedPosition, 1e-1) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem("左足y+", fno).Position)
			}
		}
		{
			expectedPosition := &mmath.MVec3{0.74736036, 9.38409308, 0.58008117}
			if !matrixes.GetItem("左足yTgt", fno).Position.PracticallyEquals(expectedPosition, 1e-1) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem("左足yTgt", fno).Position)
			}
		}
		{
			expectedPosition := &mmath.MVec3{0.74736036, 9.38409308, 0.58008117}
			if !matrixes.GetItem("左足yIK", fno).Position.PracticallyEquals(expectedPosition, 1e-1) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem("左足yIK", fno).Position)
			}
		}
		{
			expectedPosition := &mmath.MVec3{-0.03018836, 10.40081089, 1.26859617}
			if !matrixes.GetItem("左尻", fno).Position.PracticallyEquals(expectedPosition, 1e-1) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem("左尻", fno).Position)
			}
		}
		{
			expectedPosition := &mmath.MVec3{0.08276818, 5.59348757, -1.24981795}
			if !matrixes.GetItem("左ひざD", fno).Position.PracticallyEquals(expectedPosition, 1e-1) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem("左ひざD", fno).Position)
			}
		}
		{
			expectedPosition := &mmath.MVec3{0.08276818, 5.59348757, -1.24981795}
			if !matrixes.GetItem("左ひざsub", fno).Position.PracticallyEquals(expectedPosition, 1e-1) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem("左ひざsub", fno).Position)
			}
		}
		{
			expectedPosition := &mmath.MVec3{-0.09359026, 5.54494997, -1.80895985}
			if !matrixes.GetItem("左ひざsub先", fno).Position.PracticallyEquals(expectedPosition, 1e-1) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem("左ひざsub先", fno).Position)
			}
		}
		{
			expectedPosition := &mmath.MVec3{1.23779916, 1.28891465, 1.65257835}
			if !matrixes.GetItem("左ひざD2", fno).Position.PracticallyEquals(expectedPosition, 1e-1) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem("左ひざD2", fno).Position)
			}
		}
		{
			expectedPosition := &mmath.MVec3{0.1106881, 4.98643066, -1.26321915}
			if !matrixes.GetItem("左ひざD2先", fno).Position.PracticallyEquals(expectedPosition, 1e-1) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem("左ひざD2先", fno).Position)
			}
		}
		{
			expectedPosition := &mmath.MVec3{0.12060311, 4.95396153, -1.23761938}
			if !matrixes.GetItem("左ひざD2IK", fno).Position.PracticallyEquals(expectedPosition, 1e-1) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem("左ひざD2IK", fno).Position)
			}
		}
		{
			expectedPosition := &mmath.MVec3{0.88590917, 0.38407067, 0.56801614}
			if !matrixes.GetItem("左足ゆび", fno).Position.PracticallyEquals(expectedPosition, 1e-1) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem("左足ゆび", fno).Position)
			}
		}
		{
			expectedPosition := &mmath.MVec3{5.63290634e-01, -2.12439821e-04, -3.87768478e-01}
			if !matrixes.GetItem("左つま先D", fno).Position.PracticallyEquals(expectedPosition, 1e-1) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem("左つま先D", fno).Position)
			}
		}
		{
			expectedPosition := &mmath.MVec3{0.90331914, 10.27362702, 1.00949975}
			if !matrixes.GetItem("左足s", fno).Position.PracticallyEquals(expectedPosition, 1e-1) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem("左足s", fno).Position)
			}
		}
		{
			expectedPosition := &mmath.MVec3{0.08276818, 5.59348757, -1.24981795}
			if !matrixes.GetItem("左ひざs", fno).Position.PracticallyEquals(expectedPosition, 1e-1) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem("左ひざs", fno).Position)
			}
		}
		{
			expectedPosition := &mmath.MVec3{1.23779916, 1.28891465, 1.65257835}
			if !matrixes.GetItem("左足首D", fno).Position.PracticallyEquals(expectedPosition, 1e-1) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem("左足首D", fno).Position)
			}
		}
	}
}

func TestVmdMotion_AnimateBoneArmIk2(t *testing.T) {
	vr := &VmdMotionReader{}
	motionData, err := vr.ReadByFilepath("C:/MMD/mmd_base/tests/resources/唱(ダンスのみ)_0274F.vmd")

	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err)
	}

	motion := motionData.(*VmdMotion)

	pr := &pmx.PmxReader{}
	modelData, err := pr.ReadByFilepath("D:/MMD/MikuMikuDance_v926x64/UserFile/Model/VOCALOID/初音ミク/ISAO式ミク/I_ミクv4/Miku_V4_準標準.pmx")

	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err)
	}

	model := modelData.(*pmx.PmxModel)
	model.SetUp()

	{

		fno := float32(0)
		matrixes := motion.AnimateBone(fno, model, nil, true, false, "")
		{
			expectedPosition := &mmath.MVec3{0.04952335, 9.0, 1.72378033}
			if !matrixes.GetItem(pmx.CENTER.String(), fno).Position.PracticallyEquals(expectedPosition, 1e-3) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem(pmx.CENTER.String(), fno).Position)
			}
		}
		{
			expectedPosition := &mmath.MVec3{0.04952335, 7.97980869, 1.72378033}
			if !matrixes.GetItem(pmx.GROOVE.String(), fno).Position.PracticallyEquals(expectedPosition, 1e-1) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem(pmx.GROOVE.String(), fno).Position)
			}
		}
		{
			expectedPosition := &mmath.MVec3{0.04952335, 11.02838314, 2.29172656}
			if !matrixes.GetItem("腰", fno).Position.PracticallyEquals(expectedPosition, 1e-1) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem("腰", fno).Position)
			}
		}
		{
			expectedPosition := &mmath.MVec3{0.04952335, 11.9671191, 1.06765032}
			if !matrixes.GetItem("上半身", fno).Position.PracticallyEquals(expectedPosition, 1e-1) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem("上半身", fno).Position)
			}
		}
		{
			expectedPosition := &mmath.MVec3{0.26284261, 13.14576297, 0.84720008}
			if !matrixes.GetItem("上半身2", fno).Position.PracticallyEquals(expectedPosition, 1e-1) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem("上半身2", fno).Position)
			}
		}
		{
			expectedPosition := &mmath.MVec3{0.33636433, 15.27729547, 0.77435588}
			if !matrixes.GetItem("右肩", fno).Position.PracticallyEquals(expectedPosition, 1e-1) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem("右肩", fno).Position)
			}
		}
		{
			expectedPosition := &mmath.MVec3{-0.63104276, 15.44542768, 0.8507726}
			if !matrixes.GetItem("右肩C", fno).Position.PracticallyEquals(expectedPosition, 1e-1) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem("右肩C", fno).Position)
			}
		}
		{
			expectedPosition := &mmath.MVec3{-0.63104276, 15.44542768, 0.8507726}
			if !matrixes.GetItem("右腕", fno).Position.PracticallyEquals(expectedPosition, 1e-1) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem("右腕", fno).Position)
			}
		}
		{
			expectedPosition := &mmath.MVec3{-0.90326269, 14.53727204, 0.7925801}
			if !matrixes.GetItem("右腕捩", fno).Position.PracticallyEquals(expectedPosition, 1e-1) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem("右腕捩", fno).Position)
			}
		}
		{
			expectedPosition := &mmath.MVec3{-1.50502977, 12.52976106, 0.66393998}
			if !matrixes.GetItem("右ひじ", fno).Position.PracticallyEquals(expectedPosition, 1e-1) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem("右ひじ", fno).Position)
			}
		}
		{
			expectedPosition := &mmath.MVec3{-1.46843236, 12.88476121, 0.12831076}
			if !matrixes.GetItem("右手捩", fno).Position.PracticallyEquals(expectedPosition, 1e-1) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem("右手捩", fno).Position)
			}
		}
		{
			expectedPosition := &mmath.MVec3{-1.36287259, 13.90869981, -1.41662258}
			if !matrixes.GetItem("右手首", fno).Position.PracticallyEquals(expectedPosition, 1e-1) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem("右手首", fno).Position)
			}
		}
		{
			expectedPosition := &mmath.MVec3{-1.81521586, 14.00661535, -1.55616424}
			if !matrixes.GetItem("右手先", fno).Position.PracticallyEquals(expectedPosition, 1e-1) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem("右手先", fno).Position)
			}
		}
		{
			expectedPosition := &mmath.MVec3{-0.63104276, 15.44542768, 0.8507726}
			if !matrixes.GetItem("右腕YZ", fno).Position.PracticallyEquals(expectedPosition, 1e-1) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem("右腕YZ", fno).Position)
			}
		}
		{
			expectedPosition := &mmath.MVec3{-0.72589296, 15.12898892, 0.83049645}
			if !matrixes.GetItem("右腕YZ先", fno).Position.PracticallyEquals(expectedPosition, 1e-1) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem("右腕YZ先", fno).Position)
			}
		}
		{
			expectedPosition := &mmath.MVec3{-0.72589374, 15.12898632, 0.83049628}
			if !matrixes.GetItem("右腕YZIK", fno).Position.PracticallyEquals(expectedPosition, 1e-1) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem("右腕YZIK", fno).Position)
			}
		}
		{
			expectedPosition := &mmath.MVec3{-0.63104276, 15.44542768, 0.8507726}
			if !matrixes.GetItem("右腕X", fno).Position.PracticallyEquals(expectedPosition, 1e-1) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem("右腕X", fno).Position)
			}
		}
		{
			expectedPosition := &mmath.MVec3{-1.13337279, 15.59959111, 0.79468785}
			if !matrixes.GetItem("右腕X先", fno).Position.PracticallyEquals(expectedPosition, 1e-1) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem("右腕X先", fno).Position)
			}
		}
		{
			expectedPosition := &mmath.MVec3{-1.1253241, 15.60029489, 0.7461294}
			if !matrixes.GetItem("右腕XIK", fno).Position.PracticallyEquals(expectedPosition, 1e-1) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem("右腕XIK", fno).Position)
			}
		}
		{
			expectedPosition := &mmath.MVec3{-0.90325538, 14.53727326, 0.79258165}
			if !matrixes.GetItem("右腕捩YZ", fno).Position.PracticallyEquals(expectedPosition, 1e-1) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem("右腕捩YZ", fno).Position)
			}
		}
		{
			expectedPosition := &mmath.MVec3{-1.01247534, 14.17289417, 0.76923367}
			if !matrixes.GetItem("右腕捩YZTgt", fno).Position.PracticallyEquals(expectedPosition, 1e-1) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem("右腕捩YZTgt", fno).Position)
			}
		}
		{
			expectedPosition := &mmath.MVec3{-1.01248754, 14.17289597, 0.76923112}
			if !matrixes.GetItem("右腕捩YZIK", fno).Position.PracticallyEquals(expectedPosition, 1e-1) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem("右腕捩YZIK", fno).Position)
			}
		}
		{
			expectedPosition := &mmath.MVec3{-0.90325538, 14.53727326, 0.79258165}
			if !matrixes.GetItem("右腕捩X", fno).Position.PracticallyEquals(expectedPosition, 1e-1) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem("右腕捩X", fno).Position)
			}
		}
		{
			expectedPosition := &mmath.MVec3{-1.40656426, 14.68386802, 0.85919594}
			if !matrixes.GetItem("右腕捩XTgt", fno).Position.PracticallyEquals(expectedPosition, 1e-1) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem("右腕捩XTgt", fno).Position)
			}
		}
		{
			expectedPosition := &mmath.MVec3{-1.40657579, 14.68387899, 0.8591982}
			if !matrixes.GetItem("右腕捩XIK", fno).Position.PracticallyEquals(expectedPosition, 1e-1) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem("右腕捩XIK", fno).Position)
			}
		}
		{
			expectedPosition := &mmath.MVec3{-1.50499623, 12.52974836, 0.66394738}
			if !matrixes.GetItem("右ひじYZ", fno).Position.PracticallyEquals(expectedPosition, 1e-1) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem("右ひじYZ", fno).Position)
			}
		}
		{
			expectedPosition := &mmath.MVec3{-1.48334366, 12.74011791, 0.34655051}
			if !matrixes.GetItem("右ひじYZ先", fno).Position.PracticallyEquals(expectedPosition, 1e-1) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem("右ひじYZ先", fno).Position)
			}
		}
		{
			expectedPosition := &mmath.MVec3{-1.48334297, 12.74012453, 0.34654052}
			if !matrixes.GetItem("右ひじYZIK", fno).Position.PracticallyEquals(expectedPosition, 1e-1) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem("右ひじYZIK", fno).Position)
			}
		}
		{
			expectedPosition := &mmath.MVec3{-1.50499623, 12.52974836, 0.66394738}
			if !matrixes.GetItem("右ひじX", fno).Position.PracticallyEquals(expectedPosition, 1e-1) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem("右ひじX", fno).Position)
			}
		}
		{
			expectedPosition := &mmath.MVec3{-2.01179616, 12.66809052, 0.72106658}
			if !matrixes.GetItem("右ひじX先", fno).Position.PracticallyEquals(expectedPosition, 1e-1) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem("右ひじX先", fno).Position)
			}
		}
		{
			expectedPosition := &mmath.MVec3{-2.00760407, 12.67958516, 0.7289003}
			if !matrixes.GetItem("右ひじXIK", fno).Position.PracticallyEquals(expectedPosition, 1e-1) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem("右ひじXIK", fno).Position)
			}
		}
		{
			expectedPosition := &mmath.MVec3{-1.50499623, 12.52974836, 0.66394738}
			if !matrixes.GetItem("右ひじY", fno).Position.PracticallyEquals(expectedPosition, 1e-1) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem("右ひじY", fno).Position)
			}
		}
		{
			expectedPosition := &mmath.MVec3{-1.61429839, 12.16509505, 0.6405818}
			if !matrixes.GetItem("右ひじY先", fno).Position.PracticallyEquals(expectedPosition, 1e-1) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem("右ひじY先", fno).Position)
			}
		}
		{
			expectedPosition := &mmath.MVec3{-1.48334297, 12.74012453, 0.34654052}
			if !matrixes.GetItem("右ひじYIK", fno).Position.PracticallyEquals(expectedPosition, 1e-1) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem("右ひじYIK", fno).Position)
			}
		}
		{
			expectedPosition := &mmath.MVec3{-1.46845628, 12.88475892, 0.12832214}
			if !matrixes.GetItem("右手捩YZ", fno).Position.PracticallyEquals(expectedPosition, 1e-1) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem("右手捩YZ", fno).Position)
			}
		}
		{
			expectedPosition := &mmath.MVec3{-1.41168478, 13.4363328, -0.7038697}
			if !matrixes.GetItem("右手捩YZTgt", fno).Position.PracticallyEquals(expectedPosition, 1e-1) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem("右手捩YZTgt", fno).Position)
			}
		}
		{
			expectedPosition := &mmath.MVec3{-1.41156715, 13.43632015, -0.70389025}
			if !matrixes.GetItem("右手捩YZIK", fno).Position.PracticallyEquals(expectedPosition, 1e-1) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem("右手捩YZIK", fno).Position)
			}
		}
		{
			expectedPosition := &mmath.MVec3{-1.46845628, 12.88475892, 0.12832214}
			if !matrixes.GetItem("右手捩X", fno).Position.PracticallyEquals(expectedPosition, 1e-1) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem("右手捩X", fno).Position)
			}
		}
		{
			expectedPosition := &mmath.MVec3{-1.5965686, 12.06213832, -0.42564769}
			if !matrixes.GetItem("右手捩XTgt", fno).Position.PracticallyEquals(expectedPosition, 1e-1) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem("右手捩XTgt", fno).Position)
			}
		}
		{
			expectedPosition := &mmath.MVec3{-1.5965684, 12.06214091, -0.42565404}
			if !matrixes.GetItem("右手捩XIK", fno).Position.PracticallyEquals(expectedPosition, 1e-1) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem("右手捩XIK", fno).Position)
			}
		}
		{
			expectedPosition := &mmath.MVec3{-1.7198605, 13.98597326, -1.5267472}
			if !matrixes.GetItem("右手YZ先", fno).Position.PracticallyEquals(expectedPosition, 1e-1) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem("右手YZ先", fno).Position)
			}
		}
		{
			expectedPosition := &mmath.MVec3{-1.71969424, 13.98593727, -1.52669587}
			if !matrixes.GetItem("右手YZIK", fno).Position.PracticallyEquals(expectedPosition, 1e-1) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem("右手YZIK", fno).Position)
			}
		}
		{
			expectedPosition := &mmath.MVec3{-1.36306295, 13.90872698, -1.41659848}
			if !matrixes.GetItem("右手X", fno).Position.PracticallyEquals(expectedPosition, 1e-1) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem("右手X", fno).Position)
			}
		}
		{
			expectedPosition := &mmath.MVec3{-1.54727182, 13.56147176, -1.06342964}
			if !matrixes.GetItem("右手X先", fno).Position.PracticallyEquals(expectedPosition, 1e-1) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem("右手X先", fno).Position)
			}
		}
		{
			expectedPosition := &mmath.MVec3{-1.54700171, 13.5614545, -1.0633896}
			if !matrixes.GetItem("右手XIK", fno).Position.PracticallyEquals(expectedPosition, 1e-1) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem("右手XIK", fno).Position)
			}
		}
		{
			expectedPosition := &mmath.MVec3{-0.90581859, 14.5370842, 0.80752276}
			if !matrixes.GetItem("右腕捩1", fno).Position.PracticallyEquals(expectedPosition, 1e-1) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem("右腕捩1", fno).Position)
			}
		}
		{
			expectedPosition := &mmath.MVec3{-0.99954005, 14.2243783, 0.78748743}
			if !matrixes.GetItem("右腕捩2", fno).Position.PracticallyEquals(expectedPosition, 1e-1) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem("右腕捩2", fno).Position)
			}
		}
		{
			expectedPosition := &mmath.MVec3{-1.10880907, 13.85976329, 0.76412793}
			if !matrixes.GetItem("右腕捩3", fno).Position.PracticallyEquals(expectedPosition, 1e-1) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem("右腕捩3", fno).Position)
			}
		}
		{
			expectedPosition := &mmath.MVec3{-1.21298069, 13.51216081, 0.74185819}
			if !matrixes.GetItem("右腕捩4", fno).Position.PracticallyEquals(expectedPosition, 1e-1) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem("右腕捩4", fno).Position)
			}
		}
		{
			expectedPosition := &mmath.MVec3{-1.5074743, 12.52953348, 0.67889319}
			if !matrixes.GetItem("右ひじsub", fno).Position.PracticallyEquals(expectedPosition, 1e-1) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem("右ひじsub", fno).Position)
			}
		}
		{
			expectedPosition := &mmath.MVec3{-1.48606988, 12.49512375, 1.1040098}
			if !matrixes.GetItem("右ひじsub先", fno).Position.PracticallyEquals(expectedPosition, 1e-1) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem("右ひじsub先", fno).Position)
			}
		}
		{
			expectedPosition := &mmath.MVec3{-1.47286246, 12.87280799, 0.12010051}
			if !matrixes.GetItem("右手捩1", fno).Position.PracticallyEquals(expectedPosition, 1e-1) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem("右手捩1", fno).Position)
			}
		}
		{
			expectedPosition := &mmath.MVec3{-1.45874646, 13.00975416, -0.08652926}
			if !matrixes.GetItem("右手捩2", fno).Position.PracticallyEquals(expectedPosition, 1e-1) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem("右手捩2", fno).Position)
			}
		}
		{
			expectedPosition := &mmath.MVec3{-1.44072981, 13.18461598, -0.35036358}
			if !matrixes.GetItem("右手捩3", fno).Position.PracticallyEquals(expectedPosition, 1e-1) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem("右手捩3", fno).Position)
			}
		}
		{
			expectedPosition := &mmath.MVec3{-1.42368773, 13.34980879, -0.59962077}
			if !matrixes.GetItem("右手捩4", fno).Position.PracticallyEquals(expectedPosition, 1e-1) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem("右手捩4", fno).Position)
			}
		}
		{
			expectedPosition := &mmath.MVec3{-1.40457204, 13.511055, -0.84384039}
			if !matrixes.GetItem("右手捩5", fno).Position.PracticallyEquals(expectedPosition, 1e-1) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem("右手捩5", fno).Position)
			}
		}
		{
			expectedPosition := &mmath.MVec3{-1.39275926, 13.62582429, -1.01699954}
			if !matrixes.GetItem("右手捩6", fno).Position.PracticallyEquals(expectedPosition, 1e-1) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem("右手捩6", fno).Position)
			}
		}
		{
			expectedPosition := &mmath.MVec3{-1.36500465, 13.89623575, -1.42501008}
			if !matrixes.GetItem("右手首R", fno).Position.PracticallyEquals(expectedPosition, 1e-1) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem("右手首R", fno).Position)
			}
		}
		{
			expectedPosition := &mmath.MVec3{-1.36500465, 13.89623575, -1.42501008}
			if !matrixes.GetItem("右手首1", fno).Position.PracticallyEquals(expectedPosition, 1e-1) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem("右手首1", fno).Position)
			}
		}
		{
			expectedPosition := &mmath.MVec3{-1.47243053, 13.91720695, -1.52989243}
			if !matrixes.GetItem("右手首2", fno).Position.PracticallyEquals(expectedPosition, 1e-1) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem("右手首2", fno).Position)
			}
		}
	}
}

func TestVmdMotion_AnimateBoneArmIk3(t *testing.T) {
	vr := &VmdMotionReader{}
	motionData, err := vr.ReadByFilepath("C:/MMD/mlib_go/test_resources/Addiction_0F.vmd")

	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err)
	}

	motion := motionData.(*VmdMotion)

	pr := &pmx.PmxReader{}
	modelData, err := pr.ReadByFilepath("D:/MMD/MikuMikuDance_v926x64/UserFile/Model/VOCALOID/初音ミク/Sour式初音ミクVer.1.02/Black_全表示.pmx")

	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err)
	}

	model := modelData.(*pmx.PmxModel)
	model.SetUp()

	{

		fno := float32(0)
		matrixes := motion.AnimateBone(fno, model, nil, true, false, "")
		{
			expectedPosition := &mmath.MVec3{1.018832, 15.840092, 0.532239}
			if !matrixes.GetItem("左腕", fno).Position.PracticallyEquals(expectedPosition, 0.2) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem("左腕", fno).Position)
			}
		}
		{
			expectedPosition := &mmath.MVec3{1.186002, 14.510550, 0.099023}
			if !matrixes.GetItem("左腕捩", fno).Position.PracticallyEquals(expectedPosition, 1e-1) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem("左腕捩", fno).Position)
			}
		}
		{
			expectedPosition := &mmath.MVec3{1.353175, 13.181011, -0.334196}
			if !matrixes.GetItem("左ひじ", fno).Position.PracticallyEquals(expectedPosition, 1e-1) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem("左ひじ", fno).Position)
			}
		}
		{
			expectedPosition := &mmath.MVec3{1.018832, 15.840092, 0.532239}
			if !matrixes.GetItem("左腕W", fno).Position.PracticallyEquals(expectedPosition, 0.2) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem("左腕W", fno).Position)
			}
		}
		{
			expectedPosition := &mmath.MVec3{1.353175, 13.181011, -0.334196}
			if !matrixes.GetItem("左腕W先", fno).Position.PracticallyEquals(expectedPosition, 0.2) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem("左腕W先", fno).Position)
			}
		}
		{
			expectedPosition := &mmath.MVec3{1.353175, 13.181011, -0.334196}
			if !matrixes.GetItem("左腕WIK", fno).Position.PracticallyEquals(expectedPosition, 0.2) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.GetItem("左腕WIK", fno).Position)
			}
		}
	}
}
