package vmd

import (
	"testing"

	"github.com/miu200521358/mlib_go/pkg/mmath"
	"github.com/miu200521358/mlib_go/pkg/mutils/mlog"
	"github.com/miu200521358/mlib_go/pkg/pmx"
)

func TestVmdMotion_AnimateBone_Exists(t *testing.T) {
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

	{

		fno := int(10.0)
		matrixes := motion.AnimateBone(fno, model, []string{pmx.INDEX3.Left()}, false)
		{
			expectedPosition := &mmath.MVec3{0.0, 0.0, 0.0}
			if !matrixes.Get(pmx.ROOT.String()).Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get(pmx.ROOT.String()).Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{0.044920, 8.218059, 0.069347}
			if !matrixes.Get(pmx.CENTER.String()).Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get(pmx.CENTER.String()).Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{0.044920, 9.392067, 0.064877}
			if !matrixes.Get(pmx.GROOVE.String()).Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get(pmx.GROOVE.String()).Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{0.044920, 11.740084, 0.055937}
			if !matrixes.Get(pmx.WAIST.String()).Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get(pmx.WAIST.String()).Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{0.044920, 12.390969, -0.100531}
			if !matrixes.Get(pmx.UPPER.String()).Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get(pmx.UPPER.String()).Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{0.044920, 13.803633, -0.138654}
			if !matrixes.Get(pmx.UPPER2.String()).Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get(pmx.UPPER2.String()).Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{0.044920, 15.149180, 0.044429}
			if !matrixes.Get(pmx.UPPER3.String()).Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get(pmx.UPPER3.String()).Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{0.324862, 16.470263, 0.419041}
			if !matrixes.Get(pmx.SHOULDER_P.Left()).Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get(pmx.SHOULDER_P.Left()).Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{0.324862, 16.470263, 0.419041}
			if !matrixes.Get(pmx.SHOULDER.Left()).Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get(pmx.SHOULDER.Left()).Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{1.369838, 16.312170, 0.676838}
			if !matrixes.Get(pmx.ARM.Left()).Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get(pmx.ARM.Left()).Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{1.845001, 15.024807, 0.747681}
			if !matrixes.Get(pmx.ARM_TWIST.Left()).Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get(pmx.ARM_TWIST.Left()).Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{2.320162, 13.737446, 0.818525}
			if !matrixes.Get(pmx.ELBOW.Left()).Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get(pmx.ELBOW.Left()).Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{2.526190, 12.502445, 0.336127}
			if !matrixes.Get(pmx.WRIST_TWIST.Left()).Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get(pmx.WRIST_TWIST.Left()).Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{2.732219, 11.267447, -0.146273}
			if !matrixes.Get(pmx.WRIST.Left()).Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get(pmx.WRIST.Left()).Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{2.649188, 10.546797, -0.607412}
			if !matrixes.Get(pmx.INDEX1.Left()).Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get(pmx.INDEX1.Left()).Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{2.408238, 10.209290, -0.576288}
			if !matrixes.Get(pmx.INDEX2.Left()).Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get(pmx.INDEX2.Left()).Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{2.360455, 10.422402, -0.442668}
			if !matrixes.Get(pmx.INDEX3.Left()).Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get(pmx.INDEX3.Left()).Position.MMD())
			}
		}
	}
}

func TestVmdMotion_AnimateBone_Lerp(t *testing.T) {
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

	{
		fno := int(999)
		matrixes := motion.AnimateBone(fno, model, []string{pmx.INDEX3.Left()}, false)
		{
			expectedPosition := &mmath.MVec3{0.0, 0.0, 0.0}
			if !matrixes.Get(pmx.ROOT.String()).Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get(pmx.ROOT.String()).Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{-0.508560, 8.218059, 0.791827}
			if !matrixes.Get(pmx.CENTER.String()).Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get(pmx.CENTER.String()).Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{-0.508560, 9.182008, 0.787357}
			if !matrixes.Get(pmx.GROOVE.String()).Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get(pmx.GROOVE.String()).Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{-0.508560, 11.530025, 0.778416}
			if !matrixes.Get(pmx.WAIST.String()).Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get(pmx.WAIST.String()).Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{-0.508560, 12.180910, 0.621949}
			if !matrixes.Get(pmx.UPPER.String()).Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get(pmx.UPPER.String()).Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{-0.437343, 13.588836, 0.523215}
			if !matrixes.Get(pmx.UPPER2.String()).Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get(pmx.UPPER2.String()).Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{-0.552491, 14.941880, 0.528703}
			if !matrixes.Get(pmx.UPPER3.String()).Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get(pmx.UPPER3.String()).Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{-0.590927, 16.312325, 0.819156}
			if !matrixes.Get(pmx.SHOULDER_P.Left()).Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get(pmx.SHOULDER_P.Left()).Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{-0.590927, 16.312325, 0.819156}
			if !matrixes.Get(pmx.SHOULDER.Left()).Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get(pmx.SHOULDER.Left()).Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{0.072990, 16.156742, 1.666761}
			if !matrixes.Get(pmx.ARM.Left()).Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get(pmx.ARM.Left()).Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{0.043336, 15.182318, 2.635117}
			if !matrixes.Get(pmx.ARM_TWIST.Left()).Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get(pmx.ARM_TWIST.Left()).Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{0.013682, 14.207894, 3.603473}
			if !matrixes.Get(pmx.ELBOW.Left()).Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get(pmx.ELBOW.Left()).Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{1.222444, 13.711100, 3.299384}
			if !matrixes.Get(pmx.WRIST_TWIST.Left()).Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get(pmx.WRIST_TWIST.Left()).Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{2.431205, 13.214306, 2.995294}
			if !matrixes.Get(pmx.WRIST.Left()).Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get(pmx.WRIST.Left()).Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{3.283628, 13.209089, 2.884702}
			if !matrixes.Get(pmx.INDEX1.Left()).Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get(pmx.INDEX1.Left()).Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{3.665809, 13.070156, 2.797680}
			if !matrixes.Get(pmx.INDEX2.Left()).Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get(pmx.INDEX2.Left()).Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{3.886795, 12.968100, 2.718276}
			if !matrixes.Get(pmx.INDEX3.Left()).Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get(pmx.INDEX3.Left()).Position.MMD())
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

	{
		fno := int(29)
		matrixes := motion.AnimateBone(fno, model, []string{pmx.TOE.Left()}, true)
		{
			expectedPosition := &mmath.MVec3{-0.781335, 11.717622, 1.557067}
			if !matrixes.Get(pmx.LOWER.String()).Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get(pmx.LOWER.String()).Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{-0.368843, 10.614175, 2.532657}
			if !matrixes.Get(pmx.LEG.Left()).Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get(pmx.LEG.Left()).Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{0.983212, 6.945313, 0.487476}
			if !matrixes.Get(pmx.KNEE.Left()).Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get(pmx.KNEE.Left()).Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{-0.345842, 2.211842, 2.182894}
			if !matrixes.Get(pmx.ANKLE.Left()).Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get(pmx.ANKLE.Left()).Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{-0.109262, -0.025810, 1.147780}
			if !matrixes.Get(pmx.TOE.Left()).Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get(pmx.TOE.Left()).Position.MMD())
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

	{

		fno := int(3152)
		matrixes := motion.AnimateBone(fno, model, []string{pmx.TOE.Left()}, true)
		{
			expectedPosition := &mmath.MVec3{7.928583, 11.713336, 1.998830}
			if !matrixes.Get(pmx.LOWER.String()).Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get(pmx.LOWER.String()).Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{7.370017, 10.665785, 2.963280}
			if !matrixes.Get(pmx.LEG.Left()).Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get(pmx.LEG.Left()).Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{9.282883, 6.689319, 2.96825}
			if !matrixes.Get(pmx.KNEE.Left()).Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get(pmx.KNEE.Left()).Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{4.115521, 7.276527, 2.980609}
			if !matrixes.Get(pmx.ANKLE.Left()).Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get(pmx.ANKLE.Left()).Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{1.931355, 6.108739, 2.994883}
			if !matrixes.Get(pmx.TOE.Left()).Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get(pmx.TOE.Left()).Position.MMD())
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

	{

		fno := int(60)
		matrixes := motion.AnimateBone(fno, model, nil, true)
		{
			expectedPosition := &mmath.MVec3{1.931959, 11.695199, -1.411883}
			if !matrixes.Get(pmx.LOWER.String()).Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get(pmx.LOWER.String()).Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{2.927524, 10.550287, -1.218106}
			if !matrixes.Get(pmx.LEG.Left()).Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get(pmx.LEG.Left()).Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{2.263363, 7.061642, -3.837192}
			if !matrixes.Get(pmx.KNEE.Left()).Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get(pmx.KNEE.Left()).Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{2.747242, 2.529942, -1.331971}
			if !matrixes.Get(pmx.ANKLE.Left()).Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get(pmx.ANKLE.Left()).Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{1.809291, 0.242514, -1.182168}
			if !matrixes.Get(pmx.TOE.Left()).Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get(pmx.TOE.Left()).Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{2.263363, 7.061642, -3.837192}
			if !matrixes.Get(pmx.KNEE_D.Left()).Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get(pmx.KNEE.Left()).Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{1.916109, 1.177077, -1.452845}
			if !matrixes.Get(pmx.TOE_EX.Left()).Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get(pmx.TOE.Left()).Position.MMD())
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

	{

		fno := int(0)
		matrixes := motion.AnimateBone(fno, model, []string{pmx.TOE.Right()}, true)
		{
			expectedPosition := &mmath.MVec3{1.316121, 11.687257, 2.263307}
			if !matrixes.Get(pmx.LOWER.String()).Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get(pmx.LOWER.String()).Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{0.175478, 10.780540, 2.728409}
			if !matrixes.Get(pmx.LEG.Right()).Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get(pmx.LEG.Right()).Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{0.950410, 11.256771, -1.589462}
			if !matrixes.Get(pmx.KNEE.Right()).Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get(pmx.KNEE.Right()).Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{-1.025194, 7.871110, 1.828258}
			if !matrixes.Get(pmx.ANKLE.Right()).Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get(pmx.ANKLE.Right()).Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{-1.701147, 6.066556, 3.384271}
			if !matrixes.Get(pmx.TOE.Right()).Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get(pmx.TOE.Right()).Position.MMD())
			}
		}
	}
}

func TestVmdMotion_AnimateBoneLegIk5_Snow(t *testing.T) {
	mlog.SetLevel(mlog.IK_VERBOSE)

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

	{

		fno := int(7409)
		matrixes := motion.AnimateBone(fno, model, []string{pmx.TOE.Right()}, true)
		{
			expectedPosition := &mmath.MVec3{-7.652257, 11.990970, -4.511993}
			if !matrixes.Get(pmx.LOWER.String()).Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get(pmx.LOWER.String()).Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{-8.637265, 10.835548, -4.326830}
			if !matrixes.Get(pmx.LEG.Right()).Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get(pmx.LEG.Right()).Position.MMD())
			}
		}
		{
			// FIXME
			expectedPosition := &mmath.MVec3{-8.693436, 7.595280, -7.321638}
			if !matrixes.Get(pmx.KNEE.Right()).Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get(pmx.KNEE.Right()).Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{-7.521027, 2.827226, -9.035607}
			if !matrixes.Get(pmx.ANKLE.Right()).Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get(pmx.ANKLE.Right()).Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{-7.453236, 0.356456, -8.876783}
			if !matrixes.Get(pmx.TOE.Right()).Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get(pmx.TOE.Right()).Position.MMD())
			}
		}
	}
}

func TestVmdMotion_AnimateBoneLegIk6_SnowOff(t *testing.T) {
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

	{
		// IK OFF
		{

			fno := int(0)
			matrixes := motion.AnimateBone(fno, model, nil, false)
			{
				expectedPosition := &mmath.MVec3{1.622245, 6.632885, 0.713205}
				if !matrixes.Get(pmx.KNEE.Left()).Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
					t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get(pmx.KNEE.Left()).Position.MMD())
				}
			}
			{
				expectedPosition := &mmath.MVec3{1.003185, 1.474691, 0.475763}
				if !matrixes.Get(pmx.ANKLE.Left()).Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
					t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get(pmx.ANKLE.Left()).Position.MMD())
				}
			}
		}
	}
}

func TestVmdMotion_AnimateBoneLegIk6_SnowOn(t *testing.T) {
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

	{
		// IK ON
		{

			fno := int(0)
			matrixes := motion.AnimateBone(fno, model, nil, true)
			{
				expectedPosition := &mmath.MVec3{2.143878, 6.558880, 1.121747}
				if !matrixes.Get(pmx.KNEE.Left()).Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
					t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get(pmx.KNEE.Left()).Position.MMD())
				}
			}
			{
				expectedPosition := &mmath.MVec3{2.214143, 1.689811, 2.947619}
				if !matrixes.Get(pmx.ANKLE.Left()).Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
					t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get(pmx.ANKLE.Left()).Position.MMD())
				}
			}
		}
	}

}

func TestVmdMotion_AnimateBoneLegIk7_Syou(t *testing.T) {
	mlog.SetLevel(mlog.IK_VERBOSE)

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

	// 残存回転判定用
	{

		fno := int(0)
		matrixes := motion.AnimateBone(fno, model, []string{pmx.TOE.Right()}, true)
		{
			expectedPosition := &mmath.MVec3{0.721499, 11.767294, 1.638818}
			if !matrixes.Get(pmx.LOWER.String()).Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get(pmx.LOWER.String()).Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{-0.133304, 10.693992, 2.314730}
			if !matrixes.Get(pmx.LEG.Right()).Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get(pmx.LEG.Right()).Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{-2.833401, 8.174604, -0.100545}
			if !matrixes.Get(pmx.KNEE.Right()).Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get(pmx.KNEE.Right()).Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{-0.409387, 5.341005, 3.524572}
			if !matrixes.Get(pmx.ANKLE.Right()).Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get(pmx.ANKLE.Right()).Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{-0.578271, 2.874233, 3.669599}
			if !matrixes.Get(pmx.TOE.Right()).Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get(pmx.TOE.Right()).Position.MMD())
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

	{

		fno := int(278)
		matrixes := motion.AnimateBone(fno, model, []string{pmx.TOE.Right()}, true)
		{
			expectedPosition := &mmath.MVec3{0.721499, 11.767294, 1.638818}
			if !matrixes.Get(pmx.LOWER.String()).Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get(pmx.LOWER.String()).Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{-0.133304, 10.693992, 2.314730}
			if !matrixes.Get(pmx.LEG.Right()).Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get(pmx.LEG.Right()).Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{-2.833401, 8.174604, -0.100545}
			if !matrixes.Get(pmx.KNEE.Right()).Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get(pmx.KNEE.Right()).Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{-0.409387, 5.341005, 3.524572}
			if !matrixes.Get(pmx.ANKLE.Right()).Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get(pmx.ANKLE.Right()).Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{-0.578271, 2.874233, 3.669599}
			if !matrixes.Get(pmx.TOE.Right()).Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get(pmx.TOE.Right()).Position.MMD())
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

	{

		fno := int(100)
		matrixes := motion.AnimateBone(fno, model, []string{pmx.TOE.Right()}, true)
		{
			expectedPosition := &mmath.MVec3{0.365000, 11.411437, 1.963828}
			if !matrixes.Get(pmx.LOWER.String()).Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get(pmx.LOWER.String()).Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{-0.513678, 10.280550, 2.500991}
			if !matrixes.Get(pmx.LEG.Right()).Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get(pmx.LEG.Right()).Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{-2.891708, 8.162312, -0.553409}
			if !matrixes.Get(pmx.KNEE.Right()).Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get(pmx.KNEE.Right()).Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{-0.826174, 4.330670, 2.292396}
			if !matrixes.Get(pmx.ANKLE.Right()).Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get(pmx.ANKLE.Right()).Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{-1.063101, 1.865613, 2.335564}
			if !matrixes.Get(pmx.TOE.Right()).Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get(pmx.TOE.Right()).Position.MMD())
			}
		}
	}
	{

		fno := int(107)
		matrixes := motion.AnimateBone(fno, model, []string{pmx.TOE.Right()}, true)
		{
			expectedPosition := &mmath.MVec3{0.365000, 12.042871, 2.034023}
			if !matrixes.Get(pmx.LOWER.String()).Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get(pmx.LOWER.String()).Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{-0.488466, 10.920292, 2.626419}
			if !matrixes.Get(pmx.LEG.Right()).Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get(pmx.LEG.Right()).Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{-1.607765, 6.763937, 1.653586}
			if !matrixes.Get(pmx.KNEE.Right()).Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get(pmx.KNEE.Right()).Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{-1.110289, 1.718307, 2.809817}
			if !matrixes.Get(pmx.ANKLE.Right()).Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get(pmx.ANKLE.Right()).Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{-1.753089, -0.026766, 1.173958}
			if !matrixes.Get(pmx.TOE.Right()).Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get(pmx.TOE.Right()).Position.MMD())
			}
		}
	}
	{

		fno := int(272)
		matrixes := motion.AnimateBone(fno, model, []string{pmx.TOE.Right()}, true)
		{
			expectedPosition := &mmath.MVec3{-0.330117, 10.811301, 1.914508}
			if !matrixes.Get(pmx.LOWER.String()).Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get(pmx.LOWER.String()).Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{-1.325985, 9.797281, 2.479780}
			if !matrixes.Get(pmx.LEG.Right()).Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get(pmx.LEG.Right()).Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{-1.394679, 6.299243, -0.209150}
			if !matrixes.Get(pmx.KNEE.Right()).Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get(pmx.KNEE.Right()).Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{-0.865021, 1.642431, 2.044760}
			if !matrixes.Get(pmx.ANKLE.Right()).Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get(pmx.ANKLE.Right()).Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{-1.191817, -0.000789, 0.220605}
			if !matrixes.Get(pmx.TOE.Right()).Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get(pmx.TOE.Right()).Position.MMD())
			}
		}
	}
	{

		fno := int(273)
		matrixes := motion.AnimateBone(fno, model, []string{pmx.TOE.Right()}, true)
		{
			expectedPosition := &mmath.MVec3{-0.154848, 10.862784, 1.868560}
			if !matrixes.Get(pmx.LOWER.String()).Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get(pmx.LOWER.String()).Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{-1.153633, 9.846655, 2.436846}
			if !matrixes.Get(pmx.LEG.Right()).Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get(pmx.LEG.Right()).Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{-1.498977, 6.380789, -0.272370}
			if !matrixes.Get(pmx.KNEE.Right()).Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get(pmx.KNEE.Right()).Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{-0.845777, 1.802650, 2.106815}
			if !matrixes.Get(pmx.ANKLE.Right()).Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get(pmx.ANKLE.Right()).Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{-1.239674, 0.026274, 0.426385}
			if !matrixes.Get(pmx.TOE.Right()).Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get(pmx.TOE.Right()).Position.MMD())
			}
		}
	}
	{

		fno := int(274)
		matrixes := motion.AnimateBone(fno, model, []string{pmx.TOE.Right()}, true)
		{
			expectedPosition := &mmath.MVec3{0.049523, 10.960778, 1.822612}
			if !matrixes.Get(pmx.LOWER.String()).Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get(pmx.LOWER.String()).Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{-0.930675, 9.938401, 2.400088}
			if !matrixes.Get(pmx.LEG.Right()).Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get(pmx.LEG.Right()).Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{-1.710987, 6.669293, -0.459177}
			if !matrixes.Get(pmx.KNEE.Right()).Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get(pmx.KNEE.Right()).Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{-0.773748, 2.387820, 2.340310}
			if !matrixes.Get(pmx.ANKLE.Right()).Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get(pmx.ANKLE.Right()).Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{-1.256876, 0.365575, 0.994345}
			if !matrixes.Get(pmx.TOE.Right()).Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get(pmx.TOE.Right()).Position.MMD())
			}
		}
	}
	{

		fno := int(278)
		matrixes := motion.AnimateBone(fno, model, []string{pmx.TOE.Right()}, true)
		{
			expectedPosition := &mmath.MVec3{0.721499, 11.767294, 1.638818}
			if !matrixes.Get(pmx.LOWER.String()).Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get(pmx.LOWER.String()).Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{-0.133304, 10.693992, 2.314730}
			if !matrixes.Get(pmx.LEG.Right()).Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get(pmx.LEG.Right()).Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{-2.833401, 8.174604, -0.100545}
			if !matrixes.Get(pmx.KNEE.Right()).Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get(pmx.KNEE.Right()).Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{-0.409387, 5.341005, 3.524572}
			if !matrixes.Get(pmx.ANKLE.Right()).Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get(pmx.ANKLE.Right()).Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{-0.578271, 2.874233, 3.669599}
			if !matrixes.Get(pmx.TOE.Right()).Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get(pmx.TOE.Right()).Position.MMD())
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

	{

		fno := int(0)
		matrixes := motion.AnimateBone(fno, model, []string{pmx.LEG_IK.Right(), "足首_R_"}, true)
		{
			expectedPosition := &mmath.MVec3{-1.869911, 2.074591, -0.911531}
			if !matrixes.Get(pmx.LEG_IK.Right()).Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get(pmx.LEG_IK.Right()).Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{0.0, 10.142656, -1.362172}
			if !matrixes.Get(pmx.LOWER.String()).Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get(pmx.LOWER.String()).Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{-0.843381, 8.895412, -0.666409}
			if !matrixes.Get(pmx.LEG.Right()).Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get(pmx.LEG.Right()).Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{-0.274925, 5.679991, -4.384042}
			if !matrixes.Get(pmx.KNEE.Right()).Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get(pmx.KNEE.Right()).Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{-1.870632, 2.072767, -0.910016}
			if !matrixes.Get(pmx.ANKLE.Right()).Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get(pmx.ANKLE.Right()).Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{-1.485913, -0.300011, -1.310446}
			if !matrixes.Get("足首_R_").Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get("足首_R_").Position.MMD())
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

	{

		fno := int(0)
		matrixes := motion.AnimateBone(fno, model, []string{pmx.TOE.Right()}, true)
		{
			expectedPosition := &mmath.MVec3{0.0, 9.379668, -1.051170}
			if !matrixes.Get(pmx.LOWER.String()).Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get(pmx.LOWER.String()).Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{-0.919751, 8.397145, -0.324375}
			if !matrixes.Get(pmx.LEG.Right()).Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get(pmx.LEG.Right()).Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{-0.422861, 6.169319, -4.100779}
			if !matrixes.Get(pmx.KNEE.Right()).Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get(pmx.KNEE.Right()).Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{-1.821804, 2.095607, -1.186269}
			if !matrixes.Get(pmx.ANKLE.Right()).Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get(pmx.ANKLE.Right()).Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{-1.390510, -0.316872, -1.544655}
			if !matrixes.Get(pmx.TOE.Right()).Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get(pmx.TOE.Right()).Position.MMD())
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

	{

		fno := int(0)
		matrixes := motion.AnimateBone(fno, model, []string{pmx.LEG_IK.Right(), "足首_R_"}, true)
		{
			expectedPosition := &mmath.MVec3{-1.012964, 1.623157, 0.680305}
			if !matrixes.Get(pmx.LEG_IK.Right()).Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get(pmx.LEG_IK.Right()).Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{0.0, 5.953951, -0.512170}
			if !matrixes.Get(pmx.LOWER.String()).Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get(pmx.LOWER.String()).Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{-0.896440, 4.569404, -0.337760}
			if !matrixes.Get(pmx.LEG.Right()).Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get(pmx.LEG.Right()).Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{-0.691207, 1.986888, -4.553376}
			if !matrixes.Get(pmx.KNEE.Right()).Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get(pmx.KNEE.Right()).Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{-1.012964, 1.623157, 0.680305}
			if !matrixes.Get(pmx.ANKLE.Right()).Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get(pmx.ANKLE.Right()).Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{-1.013000, 0.002578, -1.146909}
			if !matrixes.Get("足首_R_").Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get("足首_R_").Position.MMD())
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

	matrixes := motion.AnimateBone(0, model,
		[]string{pmx.LEG_IK.Right(), pmx.TOE.Right(), pmx.LEG_IK.Left(), pmx.TOE.Left()}, true)

	{

		{
			expectedPosition := &mmath.MVec3{-1.216134, 1.887670, -10.78867}
			if !matrixes.Get(pmx.LEG_IK.Right()).Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get(pmx.LEG_IK.Right()).Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{0.803149, 6.056844, -10.232766}
			if !matrixes.Get(pmx.LOWER.String()).Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get(pmx.LOWER.String()).Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{0.728442, 4.560226, -11.571869}
			if !matrixes.Get(pmx.LEG.Right()).Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get(pmx.LEG.Right()).Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{4.173470, 0.361388, -11.217197}
			if !matrixes.Get(pmx.KNEE.Right()).Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get(pmx.KNEE.Right()).Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{-1.217569, 1.885731, -10.788104}
			if !matrixes.Get(pmx.ANKLE.Right()).Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get(pmx.ANKLE.Right()).Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{-0.922247, -1.163554, -10.794323}
			if !matrixes.Get(pmx.TOE.Right()).Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get(pmx.TOE.Right()).Position.MMD())
			}
		}
	}
	{

		{
			expectedPosition := &mmath.MVec3{2.322227, 1.150214, -9.644499}
			if !matrixes.Get(pmx.LEG_IK.Left()).Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get(pmx.LEG_IK.Left()).Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{0.803149, 6.056844, -10.232766}
			if !matrixes.Get(pmx.LOWER.String()).Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get(pmx.LOWER.String()).Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{0.720821, 4.639688, -8.810255}
			if !matrixes.Get(pmx.LEG.Left()).Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get(pmx.LEG.Left()).Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{6.126388, 5.074682, -8.346903}
			if !matrixes.Get(pmx.KNEE.Left()).Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get(pmx.KNEE.Left()).Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{2.323599, 1.147291, -9.645196}
			if !matrixes.Get(pmx.ANKLE.Left()).Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get(pmx.ANKLE.Left()).Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{5.163002, -0.000894, -9.714369}
			if !matrixes.Get(pmx.TOE.Left()).Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get(pmx.TOE.Left()).Position.MMD())
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

	{

		fno := int(0)
		matrixes := motion.AnimateBone(fno, model, []string{pmx.LEG_IK.Right(), pmx.TOE.Right()}, true)
		{
			expectedPosition := &mmath.MVec3{11.324574, 10.920002, -7.150005}
			if !matrixes.Get(pmx.LEG_IK.Right()).Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get(pmx.LEG_IK.Right()).Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{2.433170, 13.740387, 0.992719}
			if !matrixes.Get(pmx.LOWER.String()).Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get(pmx.LOWER.String()).Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{1.982654, 11.188538, 0.602013}
			if !matrixes.Get(pmx.LEG.Right()).Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get(pmx.LEG.Right()).Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{5.661557, 11.008962, -2.259013}
			if !matrixes.Get(pmx.KNEE.Right()).Position.MMD().PracticallyEquals(expectedPosition, 0.02) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get(pmx.KNEE.Right()).Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{9.224476, 10.979847, -5.407887}
			if !matrixes.Get(pmx.ANKLE.Right()).Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get(pmx.ANKLE.Right()).Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{11.345482, 10.263426, -7.003638}
			if !matrixes.Get(pmx.TOE.Right()).Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get(pmx.TOE.Right()).Position.MMD())
			}
		}
	}
}

func TestVmdMotion_AnimateBoneLegIk15_Bottom(t *testing.T) {
	mlog.SetLevel(mlog.IK_VERBOSE)
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

	{

		fno := int(218)
		matrixes := motion.AnimateBone(fno, model, []string{pmx.LEG_IK.Right(), "足首_R_"}, true)
		{
			expectedPosition := &mmath.MVec3{-1.358434, 1.913062, 0.611182}
			if !matrixes.Get(pmx.LEG_IK.Right()).Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get(pmx.LEG_IK.Right()).Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{0.150000, 4.253955, 0.237829}
			if !matrixes.Get(pmx.LOWER.String()).Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get(pmx.LOWER.String()).Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{-0.906292, 2.996784, 0.471846}
			if !matrixes.Get(pmx.LEG.Right()).Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get(pmx.LEG.Right()).Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{-2.533418, 3.889916, -4.114837}
			if !matrixes.Get(pmx.KNEE.Right()).Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get(pmx.KNEE.Right()).Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{-1.358807, 1.912181, 0.611265}
			if !matrixes.Get(pmx.ANKLE.Right()).Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get(pmx.ANKLE.Right()).Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{-2.040872, -0.188916, -0.430442}
			if !matrixes.Get("足首_R_").Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get("足首_R_").Position.MMD())
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

	{

		fno := int(0)
		matrixes := motion.AnimateBone(fno, model, []string{pmx.LEG_IK.Right(), pmx.TOE.Right()}, true)
		{
			expectedPosition := &mmath.MVec3{-1.216134, 1.887670, -10.78867}
			if !matrixes.Get(pmx.LEG_IK.Right()).Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get(pmx.LEG_IK.Right()).Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{0.803149, 6.056844, -10.232766}
			if !matrixes.Get(pmx.LOWER.String()).Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get(pmx.LOWER.String()).Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{0.728442, 4.560226, -11.571869}
			if !matrixes.Get(pmx.LEG.Right()).Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get(pmx.LEG.Right()).Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{4.173470, 0.361388, -11.217197}
			if !matrixes.Get(pmx.KNEE.Right()).Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get(pmx.KNEE.Right()).Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{-1.217569, 1.885731, -10.788104}
			if !matrixes.Get(pmx.ANKLE.Right()).Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get(pmx.ANKLE.Right()).Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{-0.922247, -1.163554, -10.794323}
			if !matrixes.Get(pmx.TOE.Right()).Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get(pmx.TOE.Right()).Position.MMD())
			}
		}
	}
}

func TestVmdMotion_AnimateBoneLegIk17_Snow(t *testing.T) {
	mlog.SetLevel(mlog.IK_VERBOSE)

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

	{

		fno := int(0)
		matrixes := motion.AnimateBone(fno, model, []string{pmx.TOE.Right()}, true)
		{
			expectedPosition := &mmath.MVec3{2.049998, 12.957623, 1.477440}
			if !matrixes.Get(pmx.LOWER.String()).Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get(pmx.LOWER.String()).Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{1.201382, 11.353215, 2.266898}
			if !matrixes.Get(pmx.LEG.Right()).Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get(pmx.LEG.Right()).Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{0.443043, 7.640018, -1.308741}
			if !matrixes.Get(pmx.KNEE.Right()).Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get(pmx.KNEE.Right()).Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{-0.574753, 7.943915, 3.279809}
			if !matrixes.Get(pmx.ANKLE.Right()).Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get(pmx.ANKLE.Right()).Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{-1.443098, 6.324932, 4.887177}
			if !matrixes.Get(pmx.TOE.Right()).Position.MMD().PracticallyEquals(expectedPosition, 0.1) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get(pmx.TOE.Right()).Position.MMD())
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

	{

		fno := int(107)
		matrixes := motion.AnimateBone(fno, model, []string{pmx.TOE.Right()}, true)
		{
			expectedPosition := &mmath.MVec3{0.365000, 12.042871, 2.034023}
			if !matrixes.Get(pmx.LOWER.String()).Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get(pmx.LOWER.String()).Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{-0.488466, 10.920292, 2.626419}
			if !matrixes.Get(pmx.LEG.Right()).Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get(pmx.LEG.Right()).Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{-1.607765, 6.763937, 1.653586}
			if !matrixes.Get(pmx.KNEE.Right()).Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get(pmx.KNEE.Right()).Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{-1.110289, 1.718307, 2.809817}
			if !matrixes.Get(pmx.ANKLE.Right()).Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get(pmx.ANKLE.Right()).Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{-1.753089, -0.026766, 1.173958}
			if !matrixes.Get(pmx.TOE.Right()).Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get(pmx.TOE.Right()).Position.MMD())
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

	{

		fno := int(0)
		matrixes := motion.AnimateBone(fno, model, []string{pmx.TOE.Right()}, true)
		{
			expectedPosition := &mmath.MVec3{0.000000, 9.900000, 0.000000}
			if !matrixes.Get(pmx.LOWER.String()).Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get(pmx.LOWER.String()).Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{-0.599319, 8.639606, 0.369618}
			if !matrixes.Get(pmx.LEG.Right()).Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get(pmx.LEG.Right()).Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{-0.486516, 6.323577, -2.217865}
			if !matrixes.Get(pmx.KNEE.Right()).Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get(pmx.KNEE.Right()).Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{-2.501665, 2.859252, -1.902513}
			if !matrixes.Get(pmx.ANKLE.Right()).Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get(pmx.ANKLE.Right()).Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{-3.071062, 0.841962, -2.077063}
			if !matrixes.Get(pmx.TOE.Right()).Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get(pmx.TOE.Right()).Position.MMD())
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

	{

		fno := int(107)
		matrixes := motion.AnimateBone(fno, model, []string{pmx.TOE.Right()}, true)
		{
			expectedPosition := &mmath.MVec3{0.365000, 12.042871, 2.034023}
			if !matrixes.Get(pmx.LOWER.String()).Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get(pmx.LOWER.String()).Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{-0.488466, 10.920292, 2.626419}
			if !matrixes.Get(pmx.LEG.Right()).Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get(pmx.LEG.Right()).Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{-1.607765, 6.763937, 1.653586}
			if !matrixes.Get(pmx.KNEE.Right()).Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get(pmx.KNEE.Right()).Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{-1.110289, 1.718307, 2.809817}
			if !matrixes.Get(pmx.ANKLE.Right()).Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get(pmx.ANKLE.Right()).Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{-1.753089, -0.026766, 1.173958}
			if !matrixes.Get(pmx.TOE.Right()).Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get(pmx.TOE.Right()).Position.MMD())
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

	{

		fno := int(0)
		matrixes := motion.AnimateBone(fno, model, []string{pmx.TOE.Right()}, false)
		{
			expectedPosition := &mmath.MVec3{-0.133305, 10.693993, 2.314730}
			if !matrixes.Get(pmx.LEG.Right()).Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get(pmx.LEG.Right()).Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{2.708069, 9.216356, -0.720822}
			if !matrixes.Get(pmx.KNEE.Right()).Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get(pmx.KNEE.Right()).Position.MMD())
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

	{

		fno := int(0)
		matrixes := motion.AnimateBone(fno, model, []string{pmx.TOE.Right()}, true)
		{
			expectedPosition := &mmath.MVec3{-0.133306, 10.693994, 2.314731}
			if !matrixes.Get(pmx.LEG.Right()).Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get(pmx.LEG.Right()).Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{-3.753989, 8.506582, 1.058842}
			if !matrixes.Get(pmx.KNEE.Right()).Position.MMD().PracticallyEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get(pmx.KNEE.Right()).Position.MMD())
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

	{

		fno := int(0)
		matrixes := motion.AnimateBone(fno, model, []string{pmx.TOE.Right()}, true)
		{
			expectedPosition := &mmath.MVec3{-0.133305, 10.693993, 2.314730}
			if !matrixes.Get(pmx.LEG.Right()).Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get(pmx.LEG.Right()).Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{2.081436, 7.884178, -0.268146}
			if !matrixes.Get(pmx.KNEE.Right()).Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get(pmx.KNEE.Right()).Position.MMD())
			}
		}
	}
}

func TestVmdMotion_AnimateBoneLegIk23_Addiction(t *testing.T) {
	vr := &VmdMotionReader{}
	motionData, err := vr.ReadByFilepath("../../test_resources/[A]ddiction_Lat式_0171F.vmd")

	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err)
	}

	motion := motionData.(*VmdMotion)

	pr := &pmx.PmxReader{}
	modelData, err := pr.ReadByFilepath("D:/MMD/MikuMikuDance_v926x64/UserFile/Model/VOCALOID/初音ミク/Tda式ミクワンピース/Tda式ミクワンピースRSP.pmx")

	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err)
	}

	model := modelData.(*pmx.PmxModel)

	{

		fno := int(0)
		matrixes := motion.AnimateBone(fno, model, []string{pmx.TOE_IK.Right(), "右つま先"}, true)
		{
			expectedPosition := &mmath.MVec3{0, 0.2593031, 0}
			if !matrixes.Get(pmx.ROOT.String()).Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get(pmx.ROOT.String()).Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{-1.528317, 5.033707, 3.125487}
			if !matrixes.Get(pmx.LEG_IK.Right()).Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get(pmx.LEG_IK.Right()).Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{0.609285, 12.001350, 1.555402}
			if !matrixes.Get(pmx.LOWER.String()).Position.MMD().PracticallyEquals(expectedPosition, 0.12) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get(pmx.LOWER.String()).Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{-0.129098, 10.550634, 1.348259}
			if !matrixes.Get(pmx.LEG.Right()).Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get(pmx.LEG.Right()).Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{1.661012, 6.604201, -1.196993}
			if !matrixes.Get(pmx.KNEE.Right()).Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get(pmx.KNEE.Right()).Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{-1.529553, 5.033699, 3.127081}
			if !matrixes.Get("右足首").Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get("右足首").Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{-2.044619, 3.204468, 2.877363}
			if !matrixes.Get("右つま先").Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get("右つま先").Position.MMD())
			}
		}
	}
}

func TestVmdMotion_AnimateBoneLegIk24_Positive(t *testing.T) {
	vr := &VmdMotionReader{}
	motionData, err := vr.ReadByFilepath("../../test_resources/ポジティブパレード_0526.vmd")

	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err)
	}

	motion := motionData.(*VmdMotion)

	pr := &pmx.PmxReader{}
	modelData, err := pr.ReadByFilepath("D:/MMD/MikuMikuDance_v926x64/UserFile/Model/_VMDサイジング/wa_129cm 20231028/wa_129cm_20240406.pmx")

	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err)
	}

	model := modelData.(*pmx.PmxModel)

	{
		fno := int(0)
		matrixes := motion.AnimateBone(fno, model, nil, true)
		{
			expectedPosition := &mmath.MVec3{0, 0, 0}
			if !matrixes.Get(pmx.ROOT.String()).Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get(pmx.ROOT.String()).Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{-3.312041, 6.310613, -1.134230}
			if !matrixes.Get(pmx.LEG_IK.Left()).Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get(pmx.LEG_IK.Right()).Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{-2.754258, 7.935882, -2.298871}
			if !matrixes.Get(pmx.LOWER.String()).Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get(pmx.LOWER.String()).Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{-2.455364, 6.571013, -1.935295}
			if !matrixes.Get(pmx.LEG.Left()).Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get(pmx.LEG.Right()).Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{-2.695464, 4.323516, -4.574024}
			if !matrixes.Get(pmx.KNEE.Left()).Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get(pmx.KNEE.Right()).Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{-3.322137, 6.302598, -1.131305}
			if !matrixes.Get("左脛骨").Position.MMD().PracticallyEquals(expectedPosition, 0.02) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get("左脛骨").Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{-2.575414, 5.447266, -3.254661}
			if !matrixes.Get("左足捩").Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get("左足捩").Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{-2.229677, 5.626327, -3.481028}
			if !matrixes.Get("左足捩先").Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get("左足捩先").Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{-2.455364, 6.571013, -1.935295}
			if !matrixes.Get("左足向検A").Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get("左足向検A").Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{-2.695177, 4.324148, -4.574588}
			if !matrixes.Get("左足向検A先").Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get("左足向検A先").Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{-2.695177, 4.324148, -4.574588}
			if !matrixes.Get("左足捩検B").Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get("左足捩検B").Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{-0.002697, 5.869486, -6.134800}
			if !matrixes.Get("左足捩検B先").Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get("左足捩検B先").Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{-2.877639, 4.4450495, -4.164494}
			if !matrixes.Get("左膝補").Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get("左膝補").Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{-3.523895, 4.135535, -3.716305}
			if !matrixes.Get("左膝補先").Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get("左膝補先").Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{-2.118768, 6.263350, -2.402574}
			if !matrixes.Get("左足w").Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get("左足w").Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{-2.480717, 3.120446, -5.602753}
			if !matrixes.Get("左足w先").Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get("左足w先").Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{-2.455364, 6.571013, -1.935294}
			if !matrixes.Get("左足向-").Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get("左足向-").Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{-3.322137, 6.302598, -1.131305}
			if !matrixes.Get("左脛骨D").Position.MMD().PracticallyEquals(expectedPosition, 0.02) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get("左脛骨D").Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{-3.199167, 3.952319, -4.391296}
			if !matrixes.Get("左脛骨D先").Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get("左脛骨D先").Position.MMD())
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

	{

		fno := int(3182)
		matrixes := motion.AnimateBone(fno, model, nil, true)
		{
			expectedPosition := &mmath.MVec3{0, 0, 0}
			if !matrixes.Get(pmx.ROOT.String()).Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get(pmx.ROOT.String()).Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{12.400011, 9.000000, 1.885650}
			if !matrixes.Get(pmx.CENTER.String()).Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get(pmx.CENTER.String()).Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{12.400011, 8.580067, 1.885650}
			if !matrixes.Get(pmx.GROOVE.String()).Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get(pmx.GROOVE.String()).Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{12.400011, 11.628636, 2.453597}
			if !matrixes.Get(pmx.WAIST.String()).Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get(pmx.WAIST.String()).Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{12.400011, 12.567377, 1.229520}
			if !matrixes.Get(pmx.UPPER.String()).Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get(pmx.UPPER.String()).Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{12.344202, 13.782951, 1.178849}
			if !matrixes.Get(pmx.UPPER2.String()).Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get(pmx.UPPER2.String()).Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{12.425960, 15.893852, 1.481421}
			if !matrixes.Get(pmx.SHOULDER_P.Left()).Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get(pmx.SHOULDER_P.Left()).Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{12.425960, 15.893852, 1.481421}
			if !matrixes.Get(pmx.SHOULDER.Left()).Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get(pmx.SHOULDER.Left()).Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{13.348320, 15.767927, 1.802947}
			if !matrixes.Get(pmx.ARM.Left()).Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get(pmx.ARM.Left()).Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{13.564770, 14.998386, 1.289923}
			if !matrixes.Get(pmx.ARM_TWIST.Left()).Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get(pmx.ARM_TWIST.Left()).Position.MMD())
			}
		}
		{
			// FIXME
			expectedPosition := &mmath.MVec3{14.043257, 13.297290, 0.0155864}
			if !matrixes.Get(pmx.ELBOW.Left()).Position.MMD().PracticallyEquals(expectedPosition, 0.15) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get(pmx.ELBOW.Left()).Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{13.811955, 13.552182, -0.388005}
			if !matrixes.Get(pmx.WRIST_TWIST.Left()).Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get(pmx.WRIST_TWIST.Left()).Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{13.144803, 14.287374, -1.956703}
			if !matrixes.Get(pmx.WRIST.Left()).Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get(pmx.WRIST.Left()).Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{12.813587, 14.873419, -2.570278}
			if !matrixes.Get(pmx.INDEX1.Left()).Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get(pmx.INDEX1.Left()).Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{12.541822, 15.029200, -2.709604}
			if !matrixes.Get(pmx.INDEX2.Left()).Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get(pmx.INDEX2.Left()).Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{12.476499, 14.950351, -2.502167}
			if !matrixes.Get(pmx.INDEX3.Left()).Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get(pmx.INDEX3.Left()).Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{12.620306, 14.795185, -2.295859}
			if !matrixes.Get("左人指先").Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get("左人指先").Position.MMD())
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

	{

		fno := int(0)
		matrixes := motion.AnimateBone(fno, model, nil, true)
		{
			expectedPosition := &mmath.MVec3{0.04952335, 9.0, 1.72378033}
			if !matrixes.Get(pmx.CENTER.String()).Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get(pmx.CENTER.String()).Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{0.04952335, 7.97980869, 1.72378033}
			if !matrixes.Get(pmx.GROOVE.String()).Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get(pmx.GROOVE.String()).Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{0.04952335, 11.02838314, 2.29172656}
			if !matrixes.Get("腰").Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get("腰").Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{0.04952335, 11.9671191, 1.06765032}
			if !matrixes.Get("下半身").Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get("下半身").Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{-0.24102019, 9.79926074, 1.08498769}
			if !matrixes.Get("下半身先").Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get("下半身先").Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{0.90331914, 10.27362702, 1.009499759}
			if !matrixes.Get("腰キャンセル左").Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get("腰キャンセル左").Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{0.90331914, 10.27362702, 1.00949975}
			if !matrixes.Get("左足").Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get("左足").Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{0.08276818, 5.59348757, -1.24981795}
			if !matrixes.Get("左ひざ").Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get("左ひざ").Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{5.63290634e-01, -2.12439821e-04, -3.87768478e-01}
			if !matrixes.Get("左つま先").Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get("左つま先").Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{0.90331914, 10.27362702, 1.00949975}
			if !matrixes.Get("左足D").Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get("左足D").Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{0.23453057, 5.6736954, -0.76228439}
			if !matrixes.Get("左ひざ2").Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get("左ひざ2").Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{0.12060311, 4.95396153, -1.23761938}
			if !matrixes.Get("左ひざ2先").Position.MMD().PracticallyEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get("左ひざ2先").Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{0.90331914, 10.27362702, 1.00949975}
			if !matrixes.Get("左足y+").Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get("左足y+").Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{0.74736036, 9.38409308, 0.58008117}
			if !matrixes.Get("左足yTgt").Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get("左足yTgt").Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{0.74736036, 9.38409308, 0.58008117}
			if !matrixes.Get("左足yIK").Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get("左足yIK").Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{-0.03018836, 10.40081089, 1.26859617}
			if !matrixes.Get("左尻").Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get("左尻").Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{0.08276818, 5.59348757, -1.24981795}
			if !matrixes.Get("左ひざD").Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get("左ひざD").Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{0.08276818, 5.59348757, -1.24981795}
			if !matrixes.Get("左ひざsub").Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get("左ひざsub").Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{-0.09359026, 5.54494997, -1.80895985}
			if !matrixes.Get("左ひざsub先").Position.MMD().PracticallyEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get("左ひざsub先").Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{1.23779916, 1.28891465, 1.65257835}
			if !matrixes.Get("左ひざD2").Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get("左ひざD2").Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{0.1106881, 4.98643066, -1.26321915}
			if !matrixes.Get("左ひざD2先").Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get("左ひざD2先").Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{0.12060311, 4.95396153, -1.23761938}
			if !matrixes.Get("左ひざD2IK").Position.MMD().PracticallyEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get("左ひざD2IK").Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{0.88590917, 0.38407067, 0.56801614}
			if !matrixes.Get("左足ゆび").Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get("左足ゆび").Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{5.63290634e-01, -2.12439821e-04, -3.87768478e-01}
			if !matrixes.Get("左つま先D").Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get("左つま先D").Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{0.90331914, 10.27362702, 1.00949975}
			if !matrixes.Get("左足s").Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get("左足s").Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{0.08276818, 5.59348757, -1.24981795}
			if !matrixes.Get("左ひざs").Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get("左ひざs").Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{1.23779916, 1.28891465, 1.65257835}
			if !matrixes.Get("左足首D").Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get("左足首D").Position.MMD())
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

	{

		fno := int(0)
		matrixes := motion.AnimateBone(fno, model, nil, true)
		{
			expectedPosition := &mmath.MVec3{0.04952335, 9.0, 1.72378033}
			if !matrixes.Get(pmx.CENTER.String()).Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get(pmx.CENTER.String()).Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{0.04952335, 7.97980869, 1.72378033}
			if !matrixes.Get(pmx.GROOVE.String()).Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get(pmx.GROOVE.String()).Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{0.04952335, 11.02838314, 2.29172656}
			if !matrixes.Get("腰").Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get("腰").Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{0.04952335, 11.9671191, 1.06765032}
			if !matrixes.Get("上半身").Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get("上半身").Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{0.26284261, 13.14576297, 0.84720008}
			if !matrixes.Get("上半身2").Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get("上半身2").Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{0.33636433, 15.27729547, 0.77435588}
			if !matrixes.Get("右肩").Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get("右肩").Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{-0.63104276, 15.44542768, 0.8507726}
			if !matrixes.Get("右肩C").Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get("右肩C").Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{-0.63104276, 15.44542768, 0.8507726}
			if !matrixes.Get("右腕").Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get("右腕").Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{-0.90326269, 14.53727204, 0.7925801}
			if !matrixes.Get("右腕捩").Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get("右腕捩").Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{-1.50502977, 12.52976106, 0.66393998}
			if !matrixes.Get("右ひじ").Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get("右ひじ").Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{-1.46843236, 12.88476121, 0.012831076}
			if !matrixes.Get("右手捩").Position.MMD().PracticallyEquals(expectedPosition, 0.12) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get("右手捩").Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{-1.36287259, 13.90869981, -1.41662258}
			if !matrixes.Get("右手首").Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get("右手首").Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{-1.81521586, 14.00661535, -1.55616424}
			if !matrixes.Get("右手先").Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get("右手先").Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{-0.63104276, 15.44542768, 0.8507726}
			if !matrixes.Get("右腕YZ").Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get("右腕YZ").Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{-0.72589296, 15.12898892, 0.83049645}
			if !matrixes.Get("右腕YZ先").Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get("右腕YZ先").Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{-0.72589374, 15.12898632, 0.83049628}
			if !matrixes.Get("右腕YZIK").Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get("右腕YZIK").Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{-0.63104276, 15.44542768, 0.8507726}
			if !matrixes.Get("右腕X").Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get("右腕X").Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{-1.13337279, 15.59959111, 0.79468785}
			if !matrixes.Get("右腕X先").Position.MMD().PracticallyEquals(expectedPosition, 0.1) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get("右腕X先").Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{-1.1253241, 15.60029489, 0.7461294}
			if !matrixes.Get("右腕XIK").Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get("右腕XIK").Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{-0.90325538, 14.53727326, 0.79258165}
			if !matrixes.Get("右腕捩YZ").Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get("右腕捩YZ").Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{-1.01247534, 14.17289417, 0.76923367}
			if !matrixes.Get("右腕捩YZTgt").Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get("右腕捩YZTgt").Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{-1.01248754, 14.17289597, 0.76923112}
			if !matrixes.Get("右腕捩YZIK").Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get("右腕捩YZIK").Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{-0.90325538, 14.53727326, 0.79258165}
			if !matrixes.Get("右腕捩X").Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get("右腕捩X").Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{-1.40656426, 14.68386802, 0.85919594}
			if !matrixes.Get("右腕捩XTgt").Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get("右腕捩XTgt").Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{-1.40657579, 14.68387899, 0.8591982}
			if !matrixes.Get("右腕捩XIK").Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get("右腕捩XIK").Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{-1.50499623, 12.52974836, 0.66394738}
			if !matrixes.Get("右ひじYZ").Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get("右ひじYZ").Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{-1.48334366, 12.74011791, 0.34655051}
			if !matrixes.Get("右ひじYZ先").Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get("右ひじYZ先").Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{-1.48334297, 12.74012453, 0.34654052}
			if !matrixes.Get("右ひじYZIK").Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get("右ひじYZIK").Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{-1.50499623, 12.52974836, 0.66394738}
			if !matrixes.Get("右ひじX").Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get("右ひじX").Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{-2.01179616, 12.66809052, 0.72106658}
			if !matrixes.Get("右ひじX先").Position.MMD().PracticallyEquals(expectedPosition, 0.02) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get("右ひじX先").Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{-2.00760407, 12.67958516, 0.7289003}
			if !matrixes.Get("右ひじXIK").Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get("右ひじXIK").Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{-1.50499623, 12.52974836, 0.66394738}
			if !matrixes.Get("右ひじY").Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get("右ひじY").Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{-1.61429839, 12.16509505, 0.6405818}
			if !matrixes.Get("右ひじY先").Position.MMD().PracticallyEquals(expectedPosition, 0.6) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get("右ひじY先").Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{-1.48334297, 12.74012453, 0.34654052}
			if !matrixes.Get("右ひじYIK").Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get("右ひじYIK").Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{-1.46845628, 12.88475892, 0.012832214}
			if !matrixes.Get("右手捩YZ").Position.MMD().PracticallyEquals(expectedPosition, 0.15) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get("右手捩YZ").Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{-1.41168478, 13.4363328, -0.7038697}
			if !matrixes.Get("右手捩YZTgt").Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get("右手捩YZTgt").Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{-1.41156715, 13.43632015, -0.70389025}
			if !matrixes.Get("右手捩YZIK").Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get("右手捩YZIK").Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{-1.46845628, 12.88475892, 0.012832214}
			if !matrixes.Get("右手捩X").Position.MMD().PracticallyEquals(expectedPosition, 0.15) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get("右手捩X").Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{-1.5965686, 12.06213832, -0.42564769}
			if !matrixes.Get("右手捩XTgt").Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get("右手捩XTgt").Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{-1.5965684, 12.06214091, -0.42565404}
			if !matrixes.Get("右手捩XIK").Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get("右手捩XIK").Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{-1.7198605, 13.98597326, -1.5267472}
			if !matrixes.Get("右手YZ先").Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get("右手YZ先").Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{-1.71969424, 13.98593727, -1.52669587}
			if !matrixes.Get("右手YZIK").Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get("右手YZIK").Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{-1.36306295, 13.90872698, -1.41659848}
			if !matrixes.Get("右手X").Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get("右手X").Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{-1.54727182, 13.56147176, -1.06342964}
			if !matrixes.Get("右手X先").Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get("右手X先").Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{-1.54700171, 13.5614545, -1.0633896}
			if !matrixes.Get("右手XIK").Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get("右手XIK").Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{-0.90581859, 14.5370842, 0.80752276}
			if !matrixes.Get("右腕捩1").Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get("右腕捩1").Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{-0.99954005, 14.2243783, 0.78748743}
			if !matrixes.Get("右腕捩2").Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get("右腕捩2").Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{-1.10880907, 13.85976329, 0.76412793}
			if !matrixes.Get("右腕捩3").Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get("右腕捩3").Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{-1.21298069, 13.51216081, 0.74185819}
			if !matrixes.Get("右腕捩4").Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get("右腕捩4").Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{-1.5074743, 12.52953348, 0.67889319}
			if !matrixes.Get("右ひじsub").Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get("右ひじsub").Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{-1.48606988, 12.49512375, 1.1040098}
			if !matrixes.Get("右ひじsub先").Position.MMD().PracticallyEquals(expectedPosition, 0.5) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get("右ひじsub先").Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{-1.47286246, 12.87280799, 0.012010051}
			if !matrixes.Get("右手捩1").Position.MMD().PracticallyEquals(expectedPosition, 0.2) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get("右手捩1").Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{-1.45874646, 13.00975416, -0.08652926}
			if !matrixes.Get("右手捩2").Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get("右手捩2").Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{-1.44072981, 13.18461598, -0.35036358}
			if !matrixes.Get("右手捩3").Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get("右手捩3").Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{-1.42368773, 13.34980879, -0.59962077}
			if !matrixes.Get("右手捩4").Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get("右手捩4").Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{-1.40457204, 13.511055, -0.84384039}
			if !matrixes.Get("右手捩5").Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get("右手捩5").Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{-1.39275926, 13.62582429, -1.01699954}
			if !matrixes.Get("右手捩6").Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get("右手捩6").Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{-1.36500465, 13.89623575, -1.42501008}
			if !matrixes.Get("右手首R").Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get("右手首R").Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{-1.36500465, 13.89623575, -1.42501008}
			if !matrixes.Get("右手首1").Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get("右手首1").Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{-1.47243053, 13.91720695, -1.52989243}
			if !matrixes.Get("右手首2").Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get("右手首2").Position.MMD())
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

	{

		fno := int(0)
		matrixes := motion.AnimateBone(fno, model, nil, true)
		{
			expectedPosition := &mmath.MVec3{1.018832, 15.840092, 0.532239}
			if !matrixes.Get("左腕").Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get("左腕").Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{1.186002, 14.510550, 0.099023}
			if !matrixes.Get("左腕捩").Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get("左腕捩").Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{1.353175, 13.181011, -0.334196}
			if !matrixes.Get("左ひじ").Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get("左ひじ").Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{1.018832, 15.840092, 0.532239}
			if !matrixes.Get("左腕W").Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get("左腕W").Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{1.353175, 13.181011, -0.334196}
			if !matrixes.Get("左腕W先").Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get("左腕W先").Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{1.353175, 13.181011, -0.334196}
			if !matrixes.Get("左腕WIK").Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get("左腕WIK").Position.MMD())
			}
		}
	}
}

func TestVmdMotion_AnimateBoneArmIk4_DMF(t *testing.T) {
	vr := &VmdMotionReader{}
	motionData, err := vr.ReadByFilepath("../../test_resources/nac_dmf_601.vmd")

	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err)
	}

	motion := motionData.(*VmdMotion)

	pr := &pmx.PmxReader{}
	modelData, err := pr.ReadByFilepath("D:/MMD/MikuMikuDance_v926x64/UserFile/Model/VOCALOID/初音ミク/ISAO式ミク/I_ミクv4/Miku_V4.pmx")

	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err)
	}

	model := modelData.(*pmx.PmxModel)

	{

		fno := int(0)
		matrixes := motion.AnimateBone(fno, model, nil, true)
		{
			expectedPosition := &mmath.MVec3{6.210230, 8.439670, 0.496305}
			if !matrixes.Get(pmx.CENTER.String()).Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get(pmx.CENTER.String()).Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{6.210230, 8.849669, 0.496305}
			if !matrixes.Get(pmx.GROOVE.String()).Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get(pmx.GROOVE.String()).Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{6.210230, 12.836980, -0.159825}
			if !matrixes.Get("上半身").Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get("上半身").Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{6.261481, 13.968025, 0.288966}
			if !matrixes.Get("上半身2").Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get("上半身2").Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{6.541666, 15.754716, 1.421828}
			if !matrixes.Get("左肩").Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get("左肩").Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{7.451898, 16.031992, 1.675949}
			if !matrixes.Get("左腕").Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get("左腕").Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{8.135534, 15.373729, 1.715530}
			if !matrixes.Get("左腕捩").Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get("左腕捩").Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{9.646749, 13.918620, 1.803021}
			if !matrixes.Get("左ひじ").Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get("左ひじ").Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{9.164164, 13.503792, 1.706635}
			if !matrixes.Get("左手捩").Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get("左手捩").Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{7.772219, 12.307291, 1.428628}
			if !matrixes.Get("左手首").Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get("左手首").Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{7.390504, 12.011601, 1.405503}
			if !matrixes.Get("左手先").Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get("左手先").Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{7.451900, 16.031990, 1.675949}
			if !matrixes.Get("左腕YZ").Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get("左腕YZ").Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{7.690105, 15.802624, 1.689741}
			if !matrixes.Get("左腕YZ先").Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get("左腕YZ先").Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{7.690105, 15.802622, 1.689740}
			if !matrixes.Get("左腕YZIK").Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get("左腕YZIK").Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{7.451899, 16.031988, 1.675950}
			if !matrixes.Get("左腕X").Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get("左腕X").Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{7.816861, 16.406412, 1.599419}
			if !matrixes.Get("左腕X先").Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get("左腕X先").Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{7.816858, 16.406418, 1.599418}
			if !matrixes.Get("左腕XIK").Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get("左腕XIK").Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{8.135530, 15.373726, 1.715530}
			if !matrixes.Get("左腕捩YZ").Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get("左腕捩YZ").Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{8.409824, 15.109610, 1.731412}
			if !matrixes.Get("左腕捩YZTgt").Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get("左腕捩YZTgt").Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{8.409830, 15.109617, 1.731411}
			if !matrixes.Get("左腕捩YZIK").Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get("左腕捩YZIK").Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{8.135530, 15.373725, 1.715531}
			if !matrixes.Get("左腕捩X").Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get("左腕捩X").Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{8.500528, 15.748149, 1.639511}
			if !matrixes.Get("左腕捩XTgt").Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get("左腕捩XTgt").Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{8.500531, 15.748233, 1.639508}
			if !matrixes.Get("左腕捩XIK").Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get("左腕捩XIK").Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{9.646743, 13.918595, 1.803029}
			if !matrixes.Get("左ひじYZ").Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get("左ひじYZ").Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{9.360763, 13.672787, 1.745903}
			if !matrixes.Get("左ひじYZ先").Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get("左ひじYZ先").Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{9.360781, 13.672805, 1.745905}
			if !matrixes.Get("左ひじYZIK").Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get("左ひじYZIK").Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{9.646734, 13.918593, 1.803028}
			if !matrixes.Get("左ひじX").Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get("左ひじX").Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{9.944283, 13.652989, 1.456379}
			if !matrixes.Get("左ひじX先").Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get("左ひじX先").Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{9.944304, 13.653007, 1.456381}
			if !matrixes.Get("左ひじXIK").Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get("左ひじXIK").Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{9.646734, 13.918596, 1.803028}
			if !matrixes.Get("左ひじY").Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get("左ひじY").Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{9.560862, 13.926876, 1.431514}
			if !matrixes.Get("左ひじY先").Position.MMD().PracticallyEquals(expectedPosition, 0.1) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get("左ひじY先").Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{9.360781, 13.672805, 1.745905}
			if !matrixes.Get("左ひじYIK").Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get("左ひじYIK").Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{9.164141, 13.503780, 1.706625}
			if !matrixes.Get("左手捩YZ").Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get("左手捩YZ").Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{8.414344, 12.859288, 1.556843}
			if !matrixes.Get("左手捩YZTgt").Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get("左手捩YZTgt").Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{8.414370, 12.859282, 1.556885}
			if !matrixes.Get("左手捩YZIK").Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get("左手捩YZIK").Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{9.164142, 13.503780, 1.706624}
			if !matrixes.Get("左手捩X").Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get("左手捩X").Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{9.511073, 12.928087, 2.447041}
			if !matrixes.Get("左手捩XTgt").Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get("左手捩XTgt").Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{9.511120, 12.928122, 2.447057}
			if !matrixes.Get("左手捩XIK").Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get("左手捩XIK").Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{7.772183, 12.307314, 1.428563}
			if !matrixes.Get("左手YZ先").Position.MMD().PracticallyEquals(expectedPosition, 0.4) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get("左手YZ先").Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{7.471097, 12.074032, 1.410383}
			if !matrixes.Get("左手YZ先").Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get("左手YZ先").Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{7.471111, 12.074042, 1.410384}
			if !matrixes.Get("左手YZIK").Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get("左手YZIK").Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{7.772183, 12.307314, 1.428564}
			if !matrixes.Get("左手X").Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get("左手X").Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{7.802912, 12.308764, 0.901022}
			if !matrixes.Get("左手X先").Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get("左手X先").Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{7.802991, 12.308830, 0.901079}
			if !matrixes.Get("左手XIK").Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get("左手XIK").Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{8.130125, 15.368912, 1.728851}
			if !matrixes.Get("左腕捩1").Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get("左腕捩1").Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{8.365511, 15.142246, 1.742475}
			if !matrixes.Get("左腕捩2").Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get("左腕捩2").Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{8.639965, 14.877952, 1.758356}
			if !matrixes.Get("左腕捩3").Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get("左腕捩3").Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{8.901615, 14.625986, 1.773497}
			if !matrixes.Get("左腕捩4").Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get("左腕捩4").Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{9.641270, 13.913721, 1.816324}
			if !matrixes.Get("左ひじsub").Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get("左ひじsub").Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{9.907782, 13.661371, 2.034630}
			if !matrixes.Get("左ひじsub先").Position.MMD().PracticallyEquals(expectedPosition, 0.1) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get("左ひじsub先").Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{9.165060, 13.499348, 1.721094}
			if !matrixes.Get("左手捩1").Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get("左手捩1").Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{8.978877, 13.339340, 1.683909}
			if !matrixes.Get("左手捩2").Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get("左手捩2").Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{8.741154, 13.135028, 1.636428}
			if !matrixes.Get("左手捩3").Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get("左手捩3").Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{8.516553, 12.942023, 1.591578}
			if !matrixes.Get("左手捩4").Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get("左手捩4").Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{8.301016, 12.748707, 1.544439}
			if !matrixes.Get("左手捩5").Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get("左手捩5").Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{8.145000, 12.614601, 1.513277}
			if !matrixes.Get("左手捩6").Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get("左手捩6").Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{7.777408, 12.298634, 1.439762}
			if !matrixes.Get("左手首R").Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get("左手首R").Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{7.777408, 12.298635, 1.439762}
			if !matrixes.Get("左手首1").Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get("左手首1").Position.MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{7.670320, 12.202144, 1.486689}
			if !matrixes.Get("左手首2").Position.MMD().PracticallyEquals(expectedPosition, 0.01) {
				t.Errorf("Expected %v, got %v", expectedPosition, matrixes.Get("左手首2").Position.MMD())
			}
		}
	}
}
