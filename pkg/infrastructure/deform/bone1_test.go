package deform

import (
	"testing"

	"github.com/miu200521358/mlib_go/pkg/domain/mmath"
	"github.com/miu200521358/mlib_go/pkg/domain/pmx"
	"github.com/miu200521358/mlib_go/pkg/domain/vmd"
	"github.com/miu200521358/mlib_go/pkg/infrastructure/repository"
)

func TestVmdMotion_Deform_Exists(t *testing.T) {
	vr := repository.NewVmdRepository()
	motionData, err := vr.Load("../../../test_resources/サンプルモーション.vmd")

	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err)
	}

	motion := motionData.(*vmd.VmdMotion)

	pr := repository.NewPmxRepository()
	modelData, err := pr.Load("../../../test_resources/サンプルモデル.pmx")

	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err)
	}

	model := modelData.(*pmx.PmxModel)

	{

		boneDeltas := DeformBone(model, motion, false, 10, []string{pmx.INDEX3.Left()})
		{
			expectedPosition := &mmath.MVec3{X: 0.0, Y: 0.0, Z: 0.0}
			if !boneDeltas.GetByName(pmx.ROOT.String()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.ROOT.String()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.ROOT.String()).FilledGlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{X: 0.044920, Y: 8.218059, Z: 0.069347}
			if !boneDeltas.GetByName(pmx.CENTER.String()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.CENTER.String()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.CENTER.String()).FilledGlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{X: 0.044920, Y: 9.392067, Z: 0.064877}
			if !boneDeltas.GetByName(pmx.GROOVE.String()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.GROOVE.String()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.GROOVE.String()).FilledGlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{X: 0.044920, Y: 11.740084, Z: 0.055937}
			if !boneDeltas.GetByName(pmx.WAIST.String()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.WAIST.String()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.WAIST.String()).FilledGlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{X: 0.044920, Y: 12.390969, Z: -0.100531}
			if !boneDeltas.GetByName(pmx.UPPER.String()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.UPPER.String()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.UPPER.String()).FilledGlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{X: 0.044920, Y: 13.803633, Z: -0.138654}
			if !boneDeltas.GetByName(pmx.UPPER2.String()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.UPPER2.String()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.UPPER2.String()).FilledGlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{X: 0.044920, Y: 15.149180, Z: 0.044429}
			if !boneDeltas.GetByName(pmx.UPPER3.String()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.UPPER3.String()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.UPPER3.String()).FilledGlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{X: 0.324862, Y: 16.470263, Z: 0.419041}
			if !boneDeltas.GetByName(pmx.SHOULDER_P.Left()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.SHOULDER_P.Left()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.SHOULDER_P.Left()).FilledGlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{X: 0.324862, Y: 16.470263, Z: 0.419041}
			if !boneDeltas.GetByName(pmx.SHOULDER.Left()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.SHOULDER.Left()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.SHOULDER.Left()).FilledGlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{X: 1.369838, Y: 16.312170, Z: 0.676838}
			if !boneDeltas.GetByName(pmx.ARM.Left()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.ARM.Left()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.ARM.Left()).FilledGlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{X: 1.845001, Y: 15.024807, Z: 0.747681}
			if !boneDeltas.GetByName(pmx.ARM_TWIST.Left()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.ARM_TWIST.Left()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.ARM_TWIST.Left()).FilledGlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{X: 2.320162, Y: 13.737446, Z: 0.818525}
			if !boneDeltas.GetByName(pmx.ELBOW.Left()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.ELBOW.Left()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.ELBOW.Left()).FilledGlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{X: 2.526190, Y: 12.502445, Z: 0.336127}
			if !boneDeltas.GetByName(pmx.WRIST_TWIST.Left()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.WRIST_TWIST.Left()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.WRIST_TWIST.Left()).FilledGlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{X: 2.732219, Y: 11.267447, Z: -0.146273}
			if !boneDeltas.GetByName(pmx.WRIST.Left()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.WRIST.Left()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.WRIST.Left()).FilledGlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{X: 2.649188, Y: 10.546797, Z: -0.607412}
			if !boneDeltas.GetByName(pmx.INDEX1.Left()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.INDEX1.Left()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.INDEX1.Left()).FilledGlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{X: 2.408238, Y: 10.209290, Z: -0.576288}
			if !boneDeltas.GetByName(pmx.INDEX2.Left()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.INDEX2.Left()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.INDEX2.Left()).FilledGlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{X: 2.360455, Y: 10.422402, Z: -0.442668}
			if !boneDeltas.GetByName(pmx.INDEX3.Left()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.INDEX3.Left()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.INDEX3.Left()).FilledGlobalPosition().MMD()))
			}
		}
	}
}

func TestVmdMotion_Deform_Lerp(t *testing.T) {
	vr := repository.NewVmdRepository()
	motionData, err := vr.Load("../../../test_resources/サンプルモーション.vmd")

	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err)
	}

	motion := motionData.(*vmd.VmdMotion)

	pr := repository.NewPmxRepository()
	modelData, err := pr.Load("../../../test_resources/サンプルモデル.pmx")

	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err)
	}

	model := modelData.(*pmx.PmxModel)

	{
		boneDeltas := DeformBone(model, motion, true, 999, []string{pmx.INDEX3.Left()})
		{
			expectedPosition := &mmath.MVec3{X: 0.0, Y: 0.0, Z: 0.0}
			if !boneDeltas.GetByName(pmx.ROOT.String()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.ROOT.String()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.ROOT.String()).FilledGlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{X: -0.508560, Y: 8.218059, Z: 0.791827}
			if !boneDeltas.GetByName(pmx.CENTER.String()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.CENTER.String()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.CENTER.String()).FilledGlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{X: -0.508560, Y: 9.182008, Z: 0.787357}
			if !boneDeltas.GetByName(pmx.GROOVE.String()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.GROOVE.String()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.GROOVE.String()).FilledGlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{X: -0.508560, Y: 11.530025, Z: 0.778416}
			if !boneDeltas.GetByName(pmx.WAIST.String()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.WAIST.String()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.WAIST.String()).FilledGlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{X: -0.508560, Y: 12.180910, Z: 0.621949}
			if !boneDeltas.GetByName(pmx.UPPER.String()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.UPPER.String()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.UPPER.String()).FilledGlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{X: -0.437343, Y: 13.588836, Z: 0.523215}
			if !boneDeltas.GetByName(pmx.UPPER2.String()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.UPPER2.String()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.UPPER2.String()).FilledGlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{X: -0.552491, Y: 14.941880, Z: 0.528703}
			if !boneDeltas.GetByName(pmx.UPPER3.String()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.UPPER3.String()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.UPPER3.String()).FilledGlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{X: -0.590927, Y: 16.312325, Z: 0.819156}
			if !boneDeltas.GetByName(pmx.SHOULDER_P.Left()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.SHOULDER_P.Left()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.SHOULDER_P.Left()).FilledGlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{X: -0.590927, Y: 16.312325, Z: 0.819156}
			if !boneDeltas.GetByName(pmx.SHOULDER.Left()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.SHOULDER.Left()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.SHOULDER.Left()).FilledGlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{X: 0.072990, Y: 16.156742, Z: 1.666761}
			if !boneDeltas.GetByName(pmx.ARM.Left()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.ARM.Left()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.ARM.Left()).FilledGlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{X: 0.043336, Y: 15.182318, Z: 2.635117}
			if !boneDeltas.GetByName(pmx.ARM_TWIST.Left()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.ARM_TWIST.Left()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.ARM_TWIST.Left()).FilledGlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{X: 0.013682, Y: 14.207894, Z: 3.603473}
			if !boneDeltas.GetByName(pmx.ELBOW.Left()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.ELBOW.Left()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.ELBOW.Left()).FilledGlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{X: 1.222444, Y: 13.711100, Z: 3.299384}
			if !boneDeltas.GetByName(pmx.WRIST_TWIST.Left()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.WRIST_TWIST.Left()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.WRIST_TWIST.Left()).FilledGlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{X: 2.431205, Y: 13.214306, Z: 2.995294}
			if !boneDeltas.GetByName(pmx.WRIST.Left()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.WRIST.Left()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.WRIST.Left()).FilledGlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{X: 3.283628, Y: 13.209089, Z: 2.884702}
			if !boneDeltas.GetByName(pmx.INDEX1.Left()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.INDEX1.Left()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.INDEX1.Left()).FilledGlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{X: 3.665809, Y: 13.070156, Z: 2.797680}
			if !boneDeltas.GetByName(pmx.INDEX2.Left()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.INDEX2.Left()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.INDEX2.Left()).FilledGlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{X: 3.886795, Y: 12.968100, Z: 2.718276}
			if !boneDeltas.GetByName(pmx.INDEX3.Left()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.INDEX3.Left()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.INDEX3.Left()).FilledGlobalPosition().MMD()))
			}
		}
	}

}

func TestVmdMotion_DeformLegIk1_Matsu(t *testing.T) {
	vr := repository.NewVmdRepository()
	motionData, err := vr.Load("../../../test_resources/サンプルモーション.vmd")

	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err)
	}

	motion := motionData.(*vmd.VmdMotion)

	pr := repository.NewPmxRepository()
	modelData, err := pr.Load("../../../test_resources/サンプルモデル.pmx")

	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err)
	}

	model := modelData.(*pmx.PmxModel)

	{
		boneDeltas := DeformBone(model, motion, true, 29, []string{pmx.TOE.Left(), pmx.HEEL.Left()})
		{
			expectedPosition := &mmath.MVec3{X: -0.781335, Y: 11.717622, Z: 1.557067}
			if !boneDeltas.GetByName(pmx.LOWER.String()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.LOWER.String()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.LOWER.String()).FilledGlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{X: -0.368843, Y: 10.614175, Z: 2.532657}
			if !boneDeltas.GetByName(pmx.LEG.Left()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.LEG.Left()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.LEG.Left()).FilledGlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{X: 0.983212, Y: 6.945313, Z: 0.487476}
			if !boneDeltas.GetByName(pmx.KNEE.Left()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.KNEE.Left()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.KNEE.Left()).FilledGlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{X: -0.345842, Y: 2.211842, Z: 2.182894}
			if !boneDeltas.GetByName(pmx.ANKLE.Left()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.ANKLE.Left()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.ANKLE.Left()).FilledGlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{X: -0.109262, Y: -0.025810, Z: 1.147780}
			if !boneDeltas.GetByName(pmx.TOE.Left()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.TOE.Left()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.TOE.Left()).FilledGlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{X: -0.923587, Y: 0.733788, Z: 2.624565}
			if !boneDeltas.GetByName(pmx.HEEL.Left()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.HEEL.Left()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.HEEL.Left()).FilledGlobalPosition().MMD()))
			}
		}
	}
}

func TestVmdMotion_DeformLegIk2_Matsu(t *testing.T) {
	vr := repository.NewVmdRepository()
	motionData, err := vr.Load("../../../test_resources/サンプルモーション.vmd")

	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err)
	}

	motion := motionData.(*vmd.VmdMotion)

	pr := repository.NewPmxRepository()
	modelData, err := pr.Load("../../../test_resources/サンプルモデル.pmx")

	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err)
	}

	model := modelData.(*pmx.PmxModel)

	{
		boneDeltas := DeformBone(model, motion, true, 3152, []string{pmx.TOE.Left(), pmx.HEEL.Left()})
		{
			expectedPosition := &mmath.MVec3{X: 7.928583, Y: 11.713336, Z: 1.998830}
			if !boneDeltas.GetByName(pmx.LOWER.String()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.LOWER.String()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.LOWER.String()).FilledGlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{X: 7.370017, Y: 10.665785, Z: 2.963280}
			if !boneDeltas.GetByName(pmx.LEG.Left()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.LEG.Left()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.LEG.Left()).FilledGlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{X: 9.282883, Y: 6.689319, Z: 2.96825}
			if !boneDeltas.GetByName(pmx.KNEE.Left()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.KNEE.Left()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.KNEE.Left()).FilledGlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{X: 4.115521, Y: 7.276527, Z: 2.980609}
			if !boneDeltas.GetByName(pmx.ANKLE.Left()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.ANKLE.Left()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.ANKLE.Left()).FilledGlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{X: 1.931355, Y: 6.108739, Z: 2.994883}
			if !boneDeltas.GetByName(pmx.TOE.Left()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.TOE.Left()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.TOE.Left()).FilledGlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{X: 2.569512, Y: 7.844740, Z: 3.002920}
			if !boneDeltas.GetByName(pmx.HEEL.Left()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.HEEL.Left()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.HEEL.Left()).FilledGlobalPosition().MMD()))
			}
		}
	}
}

func TestVmdMotion_DeformLegIk3_Matsu(t *testing.T) {
	// mlog.SetLevel(mlog.IK_VERBOSE)
	vr := repository.NewVmdRepository()
	motionData, err := vr.Load("../../../test_resources/腰元.vmd")

	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err)
	}

	motion := motionData.(*vmd.VmdMotion)

	pr := repository.NewPmxRepository()
	modelData, err := pr.Load("../../../test_resources/サンプルモデル.pmx")

	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err)
	}

	model := modelData.(*pmx.PmxModel)

	{
		boneDeltas := DeformBone(model, motion, true, 60, nil)
		{
			expectedPosition := &mmath.MVec3{X: 1.931959, Y: 11.695199, Z: -1.411883}
			if !boneDeltas.GetByName(pmx.LOWER.String()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.LOWER.String()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.LOWER.String()).FilledGlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{X: 2.927524, Y: 10.550287, Z: -1.218106}
			if !boneDeltas.GetByName(pmx.LEG.Left()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.LEG.Left()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.LEG.Left()).FilledGlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{X: 2.263363, Y: 7.061642, Z: -3.837192}
			if !boneDeltas.GetByName(pmx.KNEE.Left()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.KNEE.Left()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.KNEE.Left()).FilledGlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{X: 2.747242, Y: 2.529942, Z: -1.331971}
			if !boneDeltas.GetByName(pmx.ANKLE.Left()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.ANKLE.Left()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.ANKLE.Left()).FilledGlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{X: 2.263363, Y: 7.061642, Z: -3.837192}
			if !boneDeltas.GetByName(pmx.KNEE_D.Left()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.KNEE_D.Left()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.KNEE_D.Left()).FilledGlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{X: 1.916109, Y: 1.177077, Z: -1.452845}
			if !boneDeltas.GetByName(pmx.TOE_EX.Left()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.TOE_EX.Left()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.TOE_EX.Left()).FilledGlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{X: 1.809291, Y: 0.242514, Z: -1.182168}
			if !boneDeltas.GetByName(pmx.TOE.Left()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.TOE.Left()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.TOE.Left()).FilledGlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{X: 3.311764, Y: 1.159233, Z: -0.613653}
			if !boneDeltas.GetByName(pmx.HEEL.Left()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.HEEL.Left()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.HEEL.Left()).FilledGlobalPosition().MMD()))
			}
		}
	}
}

func TestVmdMotion_DeformLegIk4_Snow(t *testing.T) {
	// mlog.SetLevel(mlog.IK_VERBOSE)

	vr := repository.NewVmdRepository()
	motionData, err := vr.Load("../../../test_resources/好き雪_2794.vmd")

	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err)
	}

	motion := motionData.(*vmd.VmdMotion)

	pr := repository.NewPmxRepository()
	modelData, err := pr.Load("../../../test_resources/サンプルモデル.pmx")

	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err)
	}

	model := modelData.(*pmx.PmxModel)

	boneDeltas := DeformBone(model, motion, true, 0, []string{pmx.TOE.Right(), pmx.HEEL.Right()})
	{
		expectedPosition := &mmath.MVec3{X: 1.316121, Y: 11.687257, Z: 2.263307}
		if !boneDeltas.GetByName(pmx.LOWER.String()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.LOWER.String()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.LOWER.String()).FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: 0.175478, Y: 10.780540, Z: 2.728409}
		if !boneDeltas.GetByName(pmx.LEG.Right()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.LEG.Right()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.LEG.Right()).FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: 0.950410, Y: 11.256771, Z: -1.589462}
		if !boneDeltas.GetByName(pmx.KNEE.Right()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.KNEE.Right()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.KNEE.Right()).FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: -1.025194, Y: 7.871110, Z: 1.828258}
		if !boneDeltas.GetByName(pmx.ANKLE.Right()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.ANKLE.Right()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.ANKLE.Right()).FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: -1.701147, Y: 6.066556, Z: 3.384271}
		if !boneDeltas.GetByName(pmx.TOE.Right()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.TOE.Right()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.TOE.Right()).FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: -1.379169, Y: 7.887148, Z: 3.436968}
		if !boneDeltas.GetByName(pmx.HEEL.Right()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.HEEL.Right()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.HEEL.Right()).FilledGlobalPosition().MMD()))
		}
	}
}

func TestVmdMotion_DeformLegIk5_Koshi(t *testing.T) {
	vr := repository.NewVmdRepository()
	motionData, err := vr.Load("../../../test_resources/腰元.vmd")

	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err)
	}

	motion := motionData.(*vmd.VmdMotion)

	pr := repository.NewPmxRepository()
	modelData, err := pr.Load("../../../test_resources/サンプルモデル.pmx")

	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err)
	}

	model := modelData.(*pmx.PmxModel)

	boneDeltas := DeformBone(model, motion, true, 7409, []string{pmx.TOE.Right(), pmx.HEEL.Right()})
	{
		expectedPosition := &mmath.MVec3{X: -7.652257, Y: 11.990970, Z: -4.511993}
		if !boneDeltas.GetByName(pmx.LOWER.String()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.LOWER.String()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.LOWER.String()).FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: -8.637265, Y: 10.835548, Z: -4.326830}
		if !boneDeltas.GetByName(pmx.LEG.Right()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.LEG.Right()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.LEG.Right()).FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: -8.693436, Y: 7.595280, Z: -7.321638}
		if !boneDeltas.GetByName(pmx.KNEE.Right()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.KNEE.Right()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.KNEE.Right()).FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: -7.521027, Y: 2.827226, Z: -9.035607}
		if !boneDeltas.GetByName(pmx.ANKLE.Right()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.ANKLE.Right()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.ANKLE.Right()).FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: -7.453236, Y: 0.356456, Z: -8.876783}
		if !boneDeltas.GetByName(pmx.TOE.Right()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.04) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.TOE.Right()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.TOE.Right()).FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: -7.030497, Y: 1.820072, Z: -7.827912}
		if !boneDeltas.GetByName(pmx.HEEL.Right()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.HEEL.Right()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.HEEL.Right()).FilledGlobalPosition().MMD()))
		}
	}
}

func TestVmdMotion_DeformLegIk6_KoshiOff(t *testing.T) {
	vr := repository.NewVmdRepository()
	motionData, err := vr.Load("../../../test_resources/腰元.vmd")

	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err)
	}

	motion := motionData.(*vmd.VmdMotion)

	pr := repository.NewPmxRepository()
	modelData, err := pr.Load("../../../test_resources/サンプルモデル.pmx")

	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err)
	}

	model := modelData.(*pmx.PmxModel)

	// IK OFF
	boneDeltas := DeformBone(model, motion, false, 0, nil)
	{
		expectedPosition := &mmath.MVec3{X: 1.622245, Y: 6.632885, Z: 0.713205}
		if !boneDeltas.GetByName(pmx.KNEE.Left()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.KNEE.Left()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.KNEE.Left()).FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: 1.003185, Y: 1.474691, Z: 0.475763}
		if !boneDeltas.GetByName(pmx.ANKLE.Left()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.ANKLE.Left()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.ANKLE.Left()).FilledGlobalPosition().MMD()))
		}
	}

}

func TestVmdMotion_DeformLegIk6_KoshiOn(t *testing.T) {
	vr := repository.NewVmdRepository()
	motionData, err := vr.Load("../../../test_resources/腰元.vmd")

	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err)
	}

	motion := motionData.(*vmd.VmdMotion)

	pr := repository.NewPmxRepository()
	modelData, err := pr.Load("../../../test_resources/サンプルモデル.pmx")

	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err)
	}

	model := modelData.(*pmx.PmxModel)

	// IK ON
	boneDeltas := DeformBone(model, motion, true, 0, nil)
	{
		expectedPosition := &mmath.MVec3{X: 2.143878, Y: 6.558880, Z: 1.121747}
		if !boneDeltas.GetByName(pmx.KNEE.Left()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.KNEE.Left()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.KNEE.Left()).FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: 2.214143, Y: 1.689811, Z: 2.947619}
		if !boneDeltas.GetByName(pmx.ANKLE.Left()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.ANKLE.Left()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.ANKLE.Left()).FilledGlobalPosition().MMD()))
		}
	}

}

func TestVmdMotion_DeformLegIk6_KoshiIkOn(t *testing.T) {
	vr := repository.NewVmdRepository()
	motionData, err := vr.Load("../../../test_resources/腰元.vmd")

	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err)
	}

	motion := motionData.(*vmd.VmdMotion)

	pr := repository.NewPmxRepository()
	modelData, err := pr.Load("../../../test_resources/サンプルモデル.pmx")

	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err)
	}

	model := modelData.(*pmx.PmxModel)

	// IK ON
	fno := int(0)

	ikEnabledFrame := vmd.NewIkEnableFrame(float32(fno))
	ikEnabledFrame.Enabled = true
	ikEnabledFrame.BoneName = pmx.LEG_IK.Left()

	ikFrame := vmd.NewIkFrame(float32(fno))
	ikFrame.IkList = append(ikFrame.IkList, ikEnabledFrame)

	boneDeltas := DeformBone(model, motion, true, 0, nil)

	{
		expectedPosition := &mmath.MVec3{X: 2.143878, Y: 6.558880, Z: 1.121747}
		if !boneDeltas.GetByName(pmx.KNEE.Left()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.KNEE.Left()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.KNEE.Left()).FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: 2.214143, Y: 1.689811, Z: 2.947619}
		if !boneDeltas.GetByName(pmx.ANKLE.Left()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.ANKLE.Left()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.ANKLE.Left()).FilledGlobalPosition().MMD()))
		}
	}

}

func TestVmdMotion_DeformLegIk6_KoshiIkOff(t *testing.T) {
	vr := repository.NewVmdRepository()
	motionData, err := vr.Load("../../../test_resources/腰元.vmd")

	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err)
	}

	motion := motionData.(*vmd.VmdMotion)

	pr := repository.NewPmxRepository()
	modelData, err := pr.Load("../../../test_resources/サンプルモデル.pmx")

	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err)
	}

	model := modelData.(*pmx.PmxModel)

	// IK OFF

	fno := int(0)

	ikEnabledFrame := vmd.NewIkEnableFrame(float32(fno))
	ikEnabledFrame.Enabled = false
	ikEnabledFrame.BoneName = pmx.LEG_IK.Left()

	ikFrame := vmd.NewIkFrame(float32(fno))
	ikFrame.IkList = append(ikFrame.IkList, ikEnabledFrame)

	boneDeltas := DeformBone(model, motion, false, 0, nil)
	{
		expectedPosition := &mmath.MVec3{X: 1.622245, Y: 6.632885, Z: 0.713205}
		if !boneDeltas.GetByName(pmx.KNEE.Left()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.KNEE.Left()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.KNEE.Left()).FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: 1.003185, Y: 1.474691, Z: 0.475763}
		if !boneDeltas.GetByName(pmx.ANKLE.Left()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.ANKLE.Left()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.ANKLE.Left()).FilledGlobalPosition().MMD()))
		}
	}
}

func TestVmdMotion_DeformLegIk7_Syou(t *testing.T) {
	// mlog.SetLevel(mlog.IK_VERBOSE)

	vr := repository.NewVmdRepository()
	motionData, err := vr.Load("../../../test_resources/唱(ダンスのみ)_0278F.vmd")

	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err)
	}

	motion := motionData.(*vmd.VmdMotion)

	pr := repository.NewPmxRepository()
	modelData, err := pr.Load("../../../test_resources/サンプルモデル.pmx")

	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err)
	}

	model := modelData.(*pmx.PmxModel)

	// 残存回転判定用
	boneDeltas := DeformBone(model, motion, true, 0, []string{pmx.TOE.Right(), pmx.HEEL.Right()})
	{
		expectedPosition := &mmath.MVec3{X: 0.721499, Y: 11.767294, Z: 1.638818}
		if !boneDeltas.GetByName(pmx.LOWER.String()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.LOWER.String()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.LOWER.String()).FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: -0.133304, Y: 10.693992, Z: 2.314730}
		if !boneDeltas.GetByName(pmx.LEG.Right()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.LEG.Right()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.LEG.Right()).FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: -2.833401, Y: 8.174604, Z: -0.100545}
		if !boneDeltas.GetByName(pmx.KNEE.Right()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.KNEE.Right()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.KNEE.Right()).FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: -0.409387, Y: 5.341005, Z: 3.524572}
		if !boneDeltas.GetByName(pmx.ANKLE.Right()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.ANKLE.Right()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.ANKLE.Right()).FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: -0.578271, Y: 2.874233, Z: 3.669599}
		if !boneDeltas.GetByName(pmx.TOE.Right()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.TOE.Right()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.TOE.Right()).FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: 0.322606, Y: 4.249237, Z: 4.517416}
		if !boneDeltas.GetByName(pmx.HEEL.Right()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.HEEL.Right()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.HEEL.Right()).FilledGlobalPosition().MMD()))
		}
	}
}

func TestVmdMotion_DeformLegIk8_Syou(t *testing.T) {
	vr := repository.NewVmdRepository()
	motionData, err := vr.Load("../../../test_resources/唱(ダンスのみ)_0-300F.vmd")

	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err)
	}

	motion := motionData.(*vmd.VmdMotion)

	pr := repository.NewPmxRepository()
	modelData, err := pr.Load("../../../test_resources/サンプルモデル.pmx")

	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err)
	}

	model := modelData.(*pmx.PmxModel)

	boneDeltas := DeformBone(model, motion, true, 278, []string{pmx.TOE.Right(), pmx.HEEL.Right()})
	{
		expectedPosition := &mmath.MVec3{X: 0.721499, Y: 11.767294, Z: 1.638818}
		if !boneDeltas.GetByName(pmx.LOWER.String()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.LOWER.String()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.LOWER.String()).FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: -0.133304, Y: 10.693992, Z: 2.314730}
		if !boneDeltas.GetByName(pmx.LEG.Right()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.LEG.Right()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.LEG.Right()).FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: -2.833401, Y: 8.174604, Z: -0.100545}
		if !boneDeltas.GetByName(pmx.KNEE.Right()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.KNEE.Right()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.KNEE.Right()).FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: -0.409387, Y: 5.341005, Z: 3.524572}
		if !boneDeltas.GetByName(pmx.ANKLE.Right()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.ANKLE.Right()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.ANKLE.Right()).FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: -0.578271, Y: 2.874233, Z: 3.669599}
		if !boneDeltas.GetByName(pmx.TOE.Right()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.TOE.Right()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.TOE.Right()).FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: 0.322606, Y: 4.249237, Z: 4.517416}
		if !boneDeltas.GetByName(pmx.HEEL.Right()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.HEEL.Right()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.HEEL.Right()).FilledGlobalPosition().MMD()))
		}
	}
}

func TestVmdMotion_DeformLegIk10_Syou1(t *testing.T) {
	vr := repository.NewVmdRepository()
	motionData, err := vr.Load("../../../test_resources/唱(ダンスのみ)_0-300F.vmd")

	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err)
	}

	motion := motionData.(*vmd.VmdMotion)

	pr := repository.NewPmxRepository()
	modelData, err := pr.Load("../../../test_resources/サンプルモデル.pmx")

	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err)
	}

	model := modelData.(*pmx.PmxModel)

	boneDeltas := DeformBone(model, motion, true, 100, []string{pmx.TOE.Right(), pmx.HEEL.Right()})
	{
		expectedPosition := &mmath.MVec3{X: 0.365000, Y: 11.411437, Z: 1.963828}
		if !boneDeltas.GetByName(pmx.LOWER.String()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.LOWER.String()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.LOWER.String()).FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: -0.513678, Y: 10.280550, Z: 2.500991}
		if !boneDeltas.GetByName(pmx.LEG.Right()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.LEG.Right()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.LEG.Right()).FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: -2.891708, Y: 8.162312, Z: -0.553409}
		if !boneDeltas.GetByName(pmx.KNEE.Right()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.KNEE.Right()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.KNEE.Right()).FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: -0.826174, Y: 4.330670, Z: 2.292396}
		if !boneDeltas.GetByName(pmx.ANKLE.Right()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.ANKLE.Right()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.ANKLE.Right()).FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: -1.063101, Y: 1.865613, Z: 2.335564}
		if !boneDeltas.GetByName(pmx.TOE.Right()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.TOE.Right()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.TOE.Right()).FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: -0.178356, Y: 3.184965, Z: 3.282950}
		if !boneDeltas.GetByName(pmx.HEEL.Right()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.HEEL.Right()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.HEEL.Right()).FilledGlobalPosition().MMD()))
		}
	}
}

func TestVmdMotion_DeformLegIk10_Syou2(t *testing.T) {
	// mlog.SetLevel(mlog.IK_VERBOSE)
	vr := repository.NewVmdRepository()
	motionData, err := vr.Load("../../../test_resources/唱(ダンスのみ)_0-300F.vmd")

	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err)
	}

	motion := motionData.(*vmd.VmdMotion)

	pr := repository.NewPmxRepository()
	modelData, err := pr.Load("../../../test_resources/サンプルモデル.pmx")

	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err)
	}

	model := modelData.(*pmx.PmxModel)

	boneDeltas := DeformBone(model, motion, true, 107, []string{pmx.TOE.Right(), pmx.HEEL.Right()})
	{
		expectedPosition := &mmath.MVec3{X: 0.365000, Y: 12.042871, Z: 2.034023}
		if !boneDeltas.GetByName(pmx.LOWER.String()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.LOWER.String()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.LOWER.String()).FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: -0.488466, Y: 10.920292, Z: 2.626419}
		if !boneDeltas.GetByName(pmx.LEG.Right()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.LEG.Right()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.LEG.Right()).FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: -1.607765, Y: 6.763937, Z: 1.653586}
		if !boneDeltas.GetByName(pmx.KNEE.Right()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.KNEE.Right()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.KNEE.Right()).FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: -1.110289, Y: 1.718307, Z: 2.809817}
		if !boneDeltas.GetByName(pmx.ANKLE.Right()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.ANKLE.Right()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.ANKLE.Right()).FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: -1.753089, Y: -0.026766, Z: 1.173958}
		if !boneDeltas.GetByName(pmx.TOE.Right()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.TOE.Right()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.TOE.Right()).FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: -0.952785, Y: 0.078826, Z: 2.838099}
		if !boneDeltas.GetByName(pmx.HEEL.Right()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.HEEL.Right()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.HEEL.Right()).FilledGlobalPosition().MMD()))
		}
	}

}

func TestVmdMotion_DeformLegIk10_Syou3(t *testing.T) {
	// mlog.SetLevel(mlog.IK_VERBOSE)
	vr := repository.NewVmdRepository()
	motionData, err := vr.Load("../../../test_resources/唱(ダンスのみ)_0-300F.vmd")

	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err)
	}

	motion := motionData.(*vmd.VmdMotion)

	pr := repository.NewPmxRepository()
	modelData, err := pr.Load("../../../test_resources/サンプルモデル.pmx")

	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err)
	}

	model := modelData.(*pmx.PmxModel)

	boneDeltas := DeformBone(model, motion, true, 272, []string{pmx.TOE.Right(), pmx.HEEL.Right()})
	{
		expectedPosition := &mmath.MVec3{X: -0.330117, Y: 10.811301, Z: 1.914508}
		if !boneDeltas.GetByName(pmx.LOWER.String()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.LOWER.String()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.LOWER.String()).FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: -1.325985, Y: 9.797281, Z: 2.479780}
		if !boneDeltas.GetByName(pmx.LEG.Right()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.LEG.Right()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.LEG.Right()).FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: -1.394679, Y: 6.299243, Z: -0.209150}
		if !boneDeltas.GetByName(pmx.KNEE.Right()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.KNEE.Right()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.KNEE.Right()).FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: -0.865021, Y: 1.642431, Z: 2.044760}
		if !boneDeltas.GetByName(pmx.ANKLE.Right()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.ANKLE.Right()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.ANKLE.Right()).FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: -1.191817, Y: -0.000789, Z: 0.220605}
		if !boneDeltas.GetByName(pmx.TOE.Right()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.TOE.Right()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.TOE.Right()).FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: -0.958608, Y: -0.002146, Z: 2.055439}
		if !boneDeltas.GetByName(pmx.HEEL.Right()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.HEEL.Right()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.HEEL.Right()).FilledGlobalPosition().MMD()))
		}
	}

}

func TestVmdMotion_DeformLegIk10_Syou4(t *testing.T) {
	vr := repository.NewVmdRepository()
	motionData, err := vr.Load("../../../test_resources/唱(ダンスのみ)_0-300F.vmd")

	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err)
	}

	motion := motionData.(*vmd.VmdMotion)

	pr := repository.NewPmxRepository()
	modelData, err := pr.Load("../../../test_resources/サンプルモデル.pmx")

	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err)
	}

	model := modelData.(*pmx.PmxModel)

	boneDeltas := DeformBone(model, motion, true, 273, []string{pmx.TOE.Right(), pmx.HEEL.Right()})
	{
		expectedPosition := &mmath.MVec3{X: -0.154848, Y: 10.862784, Z: 1.868560}
		if !boneDeltas.GetByName(pmx.LOWER.String()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.LOWER.String()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.LOWER.String()).FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: -1.153633, Y: 9.846655, Z: 2.436846}
		if !boneDeltas.GetByName(pmx.LEG.Right()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.LEG.Right()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.LEG.Right()).FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: -1.498977, Y: 6.380789, Z: -0.272370}
		if !boneDeltas.GetByName(pmx.KNEE.Right()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.KNEE.Right()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.KNEE.Right()).FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: -0.845777, Y: 1.802650, Z: 2.106815}
		if !boneDeltas.GetByName(pmx.ANKLE.Right()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.ANKLE.Right()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.ANKLE.Right()).FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: -1.239674, Y: 0.026274, Z: 0.426385}
		if !boneDeltas.GetByName(pmx.TOE.Right()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.TOE.Right()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.TOE.Right()).FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: -0.797867, Y: 0.159797, Z: 2.217469}
		if !boneDeltas.GetByName(pmx.HEEL.Right()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.HEEL.Right()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.HEEL.Right()).FilledGlobalPosition().MMD()))
		}
	}

}

func TestVmdMotion_DeformLegIk10_Syou5(t *testing.T) {
	vr := repository.NewVmdRepository()
	motionData, err := vr.Load("../../../test_resources/唱(ダンスのみ)_0-300F.vmd")

	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err)
	}

	motion := motionData.(*vmd.VmdMotion)

	pr := repository.NewPmxRepository()
	modelData, err := pr.Load("../../../test_resources/サンプルモデル.pmx")

	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err)
	}

	model := modelData.(*pmx.PmxModel)

	boneDeltas := DeformBone(model, motion, true, 274, []string{pmx.TOE.Right(), pmx.HEEL.Right()})
	{
		expectedPosition := &mmath.MVec3{X: 0.049523, Y: 10.960778, Z: 1.822612}
		if !boneDeltas.GetByName(pmx.LOWER.String()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.LOWER.String()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.LOWER.String()).FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: -0.930675, Y: 9.938401, Z: 2.400088}
		if !boneDeltas.GetByName(pmx.LEG.Right()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.LEG.Right()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.LEG.Right()).FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: -1.710987, Y: 6.669293, Z: -0.459177}
		if !boneDeltas.GetByName(pmx.KNEE.Right()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.KNEE.Right()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.KNEE.Right()).FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: -0.773748, Y: 2.387820, Z: 2.340310}
		if !boneDeltas.GetByName(pmx.ANKLE.Right()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.ANKLE.Right()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.ANKLE.Right()).FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: -1.256876, Y: 0.365575, Z: 0.994345}
		if !boneDeltas.GetByName(pmx.TOE.Right()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.TOE.Right()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.TOE.Right()).FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: -0.556038, Y: 0.785363, Z: 2.653745}
		if !boneDeltas.GetByName(pmx.HEEL.Right()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.HEEL.Right()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.HEEL.Right()).FilledGlobalPosition().MMD()))
		}
	}

}

func TestVmdMotion_DeformLegIk10_Syou6(t *testing.T) {
	vr := repository.NewVmdRepository()
	motionData, err := vr.Load("../../../test_resources/唱(ダンスのみ)_0-300F.vmd")

	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err)
	}

	motion := motionData.(*vmd.VmdMotion)

	pr := repository.NewPmxRepository()
	modelData, err := pr.Load("../../../test_resources/サンプルモデル.pmx")

	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err)
	}

	model := modelData.(*pmx.PmxModel)

	boneDeltas := DeformBone(model, motion, true, 278, []string{pmx.TOE.Right(), pmx.HEEL.Right()})
	{
		expectedPosition := &mmath.MVec3{X: 0.721499, Y: 11.767294, Z: 1.638818}
		if !boneDeltas.GetByName(pmx.LOWER.String()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.LOWER.String()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.LOWER.String()).FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: -0.133304, Y: 10.693992, Z: 2.314730}
		if !boneDeltas.GetByName(pmx.LEG.Right()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.LEG.Right()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.LEG.Right()).FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: -2.833401, Y: 8.174604, Z: -0.100545}
		if !boneDeltas.GetByName(pmx.KNEE.Right()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.KNEE.Right()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.KNEE.Right()).FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: -0.409387, Y: 5.341005, Z: 3.524572}
		if !boneDeltas.GetByName(pmx.ANKLE.Right()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.ANKLE.Right()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.ANKLE.Right()).FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: -0.578271, Y: 2.874233, Z: 3.669599}
		if !boneDeltas.GetByName(pmx.TOE.Right()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.TOE.Right()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.TOE.Right()).FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: 0.322606, Y: 4.249237, Z: 4.517416}
		if !boneDeltas.GetByName(pmx.HEEL.Right()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.HEEL.Right()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.HEEL.Right()).FilledGlobalPosition().MMD()))
		}
	}

}

func TestVmdMotion_DeformLegIk11_Shining_Miku(t *testing.T) {
	// mlog.SetLevel(mlog.IK_VERBOSE)
	vr := repository.NewVmdRepository()
	motionData, err := vr.Load("../../../test_resources/シャイニングミラクル_50F.vmd")

	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err)
	}

	motion := motionData.(*vmd.VmdMotion)

	pr := repository.NewPmxRepository()
	modelData, err := pr.Load("D:/MMD/MikuMikuDance_v926x64/UserFile/Model/VOCALOID/初音ミク/Tda式初音ミク_盗賊つばき流Ｍトレースモデル配布 v1.07/Tda式初音ミク_盗賊つばき流Mトレースモデルv1.07_かかと.pmx")

	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err)
	}

	model := modelData.(*pmx.PmxModel)

	boneDeltas := DeformBone(model, motion, true, 0, []string{pmx.LEG_IK.Right(), pmx.HEEL.Right(), "足首_R_"})
	{
		expectedPosition := &mmath.MVec3{X: -1.869911, Y: 2.074591, Z: -0.911531}
		if !boneDeltas.GetByName(pmx.LEG_IK.Right()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.LEG_IK.Right()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.LEG_IK.Right()).FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: 0.0, Y: 0.002071, Z: 0.0}
		if !boneDeltas.GetByName(pmx.ROOT.String()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.ROOT.String()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.ROOT.String()).FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: 0.0, Y: 8.404771, Z: -0.850001}
		if !boneDeltas.GetByName(pmx.CENTER.String()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.CENTER.String()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.CENTER.String()).FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: 0.0, Y: 5.593470, Z: -0.850001}
		if !boneDeltas.GetByName(pmx.GROOVE.String()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.GROOVE.String()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.GROOVE.String()).FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: 0.0, Y: 9.311928, Z: -0.586922}
		if !boneDeltas.GetByName(pmx.WAIST.String()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.WAIST.String()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.WAIST.String()).FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: 0.0, Y: 10.142656, Z: -1.362172}
		if !boneDeltas.GetByName(pmx.LOWER.String()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.LOWER.String()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.LOWER.String()).FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: -0.843381, Y: 8.895412, Z: -0.666409}
		if !boneDeltas.GetByName(pmx.LEG.Right()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.LEG.Right()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.LEG.Right()).FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: -0.274925, Y: 5.679991, Z: -4.384042}
		if !boneDeltas.GetByName(pmx.KNEE.Right()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.KNEE.Right()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.KNEE.Right()).FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: -1.870632, Y: 2.072767, Z: -0.910016}
		if !boneDeltas.GetByName(pmx.ANKLE.Right()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.ANKLE.Right()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.ANKLE.Right()).FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: -1.485913, Y: -0.300011, Z: -1.310446}
		if !boneDeltas.GetByName("足首_R_").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("足首_R_").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("足首_R_").FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: -1.894769, Y: 0.790468, Z: 0.087442}
		if !boneDeltas.GetByName(pmx.HEEL.Right()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.HEEL.Right()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.HEEL.Right()).FilledGlobalPosition().MMD()))
		}
	}
}

func TestVmdMotion_DeformLegIk11_Shining_Vroid(t *testing.T) {
	vr := repository.NewVmdRepository()
	motionData, err := vr.Load("../../../test_resources/シャイニングミラクル_50F.vmd")

	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err)
	}

	motion := motionData.(*vmd.VmdMotion)

	pr := repository.NewPmxRepository()
	modelData, err := pr.Load("../../../test_resources/サンプルモデル.pmx")

	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err)
	}

	model := modelData.(*pmx.PmxModel)

	boneDeltas := DeformBone(model, motion, true, 0, []string{pmx.TOE.Right(), pmx.HEEL.Right()})
	{
		expectedPosition := &mmath.MVec3{X: 0.0, Y: 9.379668, Z: -1.051170}
		if !boneDeltas.GetByName(pmx.LOWER.String()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.LOWER.String()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.LOWER.String()).FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: -0.919751, Y: 8.397145, Z: -0.324375}
		if !boneDeltas.GetByName(pmx.LEG.Right()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.LEG.Right()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.LEG.Right()).FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: -0.422861, Y: 6.169319, Z: -4.100779}
		if !boneDeltas.GetByName(pmx.KNEE.Right()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.KNEE.Right()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.KNEE.Right()).FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: -1.821804, Y: 2.095607, Z: -1.186269}
		if !boneDeltas.GetByName(pmx.ANKLE.Right()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.ANKLE.Right()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.ANKLE.Right()).FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: -1.390510, Y: -0.316872, Z: -1.544655}
		if !boneDeltas.GetByName(pmx.TOE.Right()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.TOE.Right()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.TOE.Right()).FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: -1.852786, Y: 0.811991, Z: -0.154341}
		if !boneDeltas.GetByName(pmx.HEEL.Right()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.HEEL.Right()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.HEEL.Right()).FilledGlobalPosition().MMD()))
		}
	}

}

func TestVmdMotion_DeformLegIk12_Down_Miku(t *testing.T) {
	// mlog.SetLevel(mlog.IK_VERBOSE)
	vr := repository.NewVmdRepository()
	motionData, err := vr.Load("../../../test_resources/しゃがむ.vmd")

	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err)
	}

	motion := motionData.(*vmd.VmdMotion)

	pr := repository.NewPmxRepository()
	modelData, err := pr.Load("D:/MMD/MikuMikuDance_v926x64/UserFile/Model/VOCALOID/初音ミク/Tda式初音ミク_盗賊つばき流Ｍトレースモデル配布 v1.07/Tda式初音ミク_盗賊つばき流Mトレースモデルv1.07_かかと.pmx")

	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err)
	}

	model := modelData.(*pmx.PmxModel)

	boneDeltas := DeformBone(model, motion, true, 0, []string{pmx.LEG_IK.Right(), pmx.HEEL.Right(), "足首_R_"})
	{
		expectedPosition := &mmath.MVec3{X: -1.012964, Y: 1.623157, Z: 0.680305}
		if !boneDeltas.GetByName(pmx.LEG_IK.Right()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.LEG_IK.Right()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.LEG_IK.Right()).FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: 0.0, Y: 5.953951, Z: -0.512170}
		if !boneDeltas.GetByName(pmx.LOWER.String()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.LOWER.String()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.LOWER.String()).FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: -0.896440, Y: 4.569404, Z: -0.337760}
		if !boneDeltas.GetByName(pmx.LEG.Right()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.LEG.Right()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.LEG.Right()).FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: -0.691207, Y: 1.986888, Z: -4.553376}
		if !boneDeltas.GetByName(pmx.KNEE.Right()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.KNEE.Right()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.KNEE.Right()).FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: -1.012964, Y: 1.623157, Z: 0.680305}
		if !boneDeltas.GetByName(pmx.ANKLE.Right()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.ANKLE.Right()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.ANKLE.Right()).FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: -1.013000, Y: 0.002578, Z: -1.146909}
		if !boneDeltas.GetByName("足首_R_").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("足首_R_").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("足首_R_").FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: -1.056216, Y: -0.001008, Z: 0.676086}
		if !boneDeltas.GetByName(pmx.HEEL.Right()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.HEEL.Right()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.HEEL.Right()).FilledGlobalPosition().MMD()))
		}
	}
}

func TestVmdMotion_DeformLegIk13_Lamb(t *testing.T) {
	vr := repository.NewVmdRepository()
	motionData, err := vr.Load("../../../test_resources/Lamb_2689F.vmd")

	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err)
	}

	motion := motionData.(*vmd.VmdMotion)

	pr := repository.NewPmxRepository()
	modelData, err := pr.Load("D:/MMD/MikuMikuDance_v926x64/UserFile/Model/ゲーム/戦国BASARA/幸村 たぬき式 ver.1.24/真田幸村没第二衣装1.24軽量版.pmx")

	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err)
	}

	model := modelData.(*pmx.PmxModel)
	boneDeltas := DeformBone(model, motion, true, 0, []string{pmx.LEG_IK.Right(), pmx.TOE.Right(), pmx.LEG_IK.Left(), pmx.TOE.Left(), pmx.HEEL.Left()})

	{

		{
			expectedPosition := &mmath.MVec3{X: -1.216134, Y: 1.887670, Z: -10.78867}
			if !boneDeltas.GetByName(pmx.LEG_IK.Right()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.LEG_IK.Right()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.LEG_IK.Right()).FilledGlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{X: 0.803149, Y: 6.056844, Z: -10.232766}
			if !boneDeltas.GetByName(pmx.LOWER.String()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.LOWER.String()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.LOWER.String()).FilledGlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{X: 0.728442, Y: 4.560226, Z: -11.571869}
			if !boneDeltas.GetByName(pmx.LEG.Right()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.LEG.Right()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.LEG.Right()).FilledGlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{X: 4.173470, Y: 0.361388, Z: -11.217197}
			if !boneDeltas.GetByName(pmx.KNEE.Right()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.KNEE.Right()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.KNEE.Right()).FilledGlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{X: -1.217569, Y: 1.885731, Z: -10.788104}
			if !boneDeltas.GetByName(pmx.ANKLE.Right()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.ANKLE.Right()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.ANKLE.Right()).FilledGlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{X: -0.922247, Y: -1.163554, Z: -10.794323}
			if !boneDeltas.GetByName(pmx.TOE.Right()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.TOE.Right()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.TOE.Right()).FilledGlobalPosition().MMD()))
			}
		}
	}
	{

		{
			expectedPosition := &mmath.MVec3{X: 2.322227, Y: 1.150214, Z: -9.644499}
			if !boneDeltas.GetByName(pmx.LEG_IK.Left()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.LEG_IK.Left()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.LEG_IK.Left()).FilledGlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{X: 0.803149, Y: 6.056844, Z: -10.232766}
			if !boneDeltas.GetByName(pmx.LOWER.String()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.LOWER.String()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.LOWER.String()).FilledGlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{X: 0.720821, Y: 4.639688, Z: -8.810255}
			if !boneDeltas.GetByName(pmx.LEG.Left()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.LEG.Left()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.LEG.Left()).FilledGlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{X: 6.126388, Y: 5.074682, Z: -8.346903}
			if !boneDeltas.GetByName(pmx.KNEE.Left()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.KNEE.Left()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.KNEE.Left()).FilledGlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{X: 2.323599, Y: 1.147291, Z: -9.645196}
			if !boneDeltas.GetByName(pmx.ANKLE.Left()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.ANKLE.Left()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.ANKLE.Left()).FilledGlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{X: 5.163002, Y: -0.000894, Z: -9.714369}
			if !boneDeltas.GetByName(pmx.TOE.Left()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.TOE.Left()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.TOE.Left()).FilledGlobalPosition().MMD()))
			}
		}
	}
}

func TestVmdMotion_DeformLegIk14_Ballet(t *testing.T) {
	vr := repository.NewVmdRepository()
	motionData, err := vr.Load("../../../test_resources/ミク用バレリーコ_1069.vmd")

	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err)
	}

	motion := motionData.(*vmd.VmdMotion)

	pr := repository.NewPmxRepository()
	modelData, err := pr.Load("D:/MMD/MikuMikuDance_v926x64/UserFile/Model/_あにまさ式/初音ミク_準標準.pmx")

	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err)
	}

	model := modelData.(*pmx.PmxModel)

	boneDeltas := DeformBone(model, motion, true, 0, []string{pmx.LEG_IK.Right(), pmx.TOE.Right(), pmx.HEEL.Right()})
	{
		expectedPosition := &mmath.MVec3{X: 11.324574, Y: 10.920002, Z: -7.150005}
		if !boneDeltas.GetByName(pmx.LEG_IK.Right()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.LEG_IK.Right()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.LEG_IK.Right()).FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: 2.433170, Y: 13.740387, Z: 0.992719}
		if !boneDeltas.GetByName(pmx.LOWER.String()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.LOWER.String()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.LOWER.String()).FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: 1.982654, Y: 11.188538, Z: 0.602013}
		if !boneDeltas.GetByName(pmx.LEG.Right()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.LEG.Right()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.LEG.Right()).FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: 5.661557, Y: 11.008962, Z: -2.259013}
		if !boneDeltas.GetByName(pmx.KNEE.Right()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.KNEE.Right()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.KNEE.Right()).FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: 9.224476, Y: 10.979847, Z: -5.407887}
		if !boneDeltas.GetByName(pmx.ANKLE.Right()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.ANKLE.Right()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.ANKLE.Right()).FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: 11.345482, Y: 10.263426, Z: -7.003638}
		if !boneDeltas.GetByName(pmx.TOE.Right()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.TOE.Right()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.TOE.Right()).FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: 9.406674, Y: 9.687277, Z: -5.710646}
		if !boneDeltas.GetByName(pmx.HEEL.Right()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.HEEL.Right()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.HEEL.Right()).FilledGlobalPosition().MMD()))
		}
	}
}

func TestVmdMotion_DeformLegIk15_Bottom(t *testing.T) {
	// mlog.SetLevel(mlog.IK_VERBOSE)

	vr := repository.NewVmdRepository()
	motionData, err := vr.Load("../../../test_resources/●ボトム_0-300.vmd")

	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err)
	}

	motion := motionData.(*vmd.VmdMotion)

	pr := repository.NewPmxRepository()
	modelData, err := pr.Load("D:/MMD/MikuMikuDance_v926x64/UserFile/Model/VOCALOID/初音ミク/Tda式初音ミク_盗賊つばき流Ｍトレースモデル配布 v1.07/Tda式初音ミク_盗賊つばき流Mトレースモデルv1.07_かかと.pmx")

	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err)
	}

	model := modelData.(*pmx.PmxModel)

	boneDeltas := DeformBone(model, motion, true, 218, []string{pmx.LEG_IK.Right(), pmx.HEEL.Right(), "足首_R_"})
	{
		expectedPosition := &mmath.MVec3{X: -1.358434, Y: 1.913062, Z: 0.611182}
		if !boneDeltas.GetByName(pmx.LEG_IK.Right()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.LEG_IK.Right()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.LEG_IK.Right()).FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: 0.150000, Y: 4.253955, Z: 0.237829}
		if !boneDeltas.GetByName(pmx.LOWER.String()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.LOWER.String()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.LOWER.String()).FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: -0.906292, Y: 2.996784, Z: 0.471846}
		if !boneDeltas.GetByName(pmx.LEG.Right()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.LEG.Right()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.LEG.Right()).FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: -2.533418, Y: 3.889916, Z: -4.114837}
		if !boneDeltas.GetByName(pmx.KNEE.Right()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.KNEE.Right()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.KNEE.Right()).FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: -1.358807, Y: 1.912181, Z: 0.611265}
		if !boneDeltas.GetByName(pmx.ANKLE.Right()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.ANKLE.Right()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.ANKLE.Right()).FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: -2.040872, Y: -0.188916, Z: -0.430442}
		if !boneDeltas.GetByName("足首_R_").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("足首_R_").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("足首_R_").FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: -1.292688, Y: 0.375211, Z: 1.133899}
		if !boneDeltas.GetByName(pmx.HEEL.Right()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.HEEL.Right()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.HEEL.Right()).FilledGlobalPosition().MMD()))
		}
	}
}

func TestVmdMotion_DeformLegIk16_Lamb(t *testing.T) {
	vr := repository.NewVmdRepository()
	motionData, err := vr.Load("../../../test_resources/Lamb_2689F.vmd")

	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err)
	}

	motion := motionData.(*vmd.VmdMotion)

	pr := repository.NewPmxRepository()
	modelData, err := pr.Load("D:/MMD/MikuMikuDance_v926x64/UserFile/Model/ゲーム/戦国BASARA/幸村 たぬき式 ver.1.24/真田幸村没第二衣装1.24軽量版.pmx")

	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err)
	}

	model := modelData.(*pmx.PmxModel)

	boneDeltas := DeformBone(model, motion, true, 0, []string{pmx.LEG_IK.Right(), pmx.TOE.Right(), pmx.HEEL.Right()})

	{
		expectedPosition := &mmath.MVec3{X: -1.216134, Y: 1.887670, Z: -10.78867}
		if !boneDeltas.GetByName(pmx.LEG_IK.Right()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.LEG_IK.Right()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.LEG_IK.Right()).FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: 0.803149, Y: 6.056844, Z: -10.232766}
		if !boneDeltas.GetByName(pmx.LOWER.String()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.LOWER.String()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.LOWER.String()).FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: 0.728442, Y: 4.560226, Z: -11.571869}
		if !boneDeltas.GetByName(pmx.LEG.Right()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.LEG.Right()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.LEG.Right()).FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: 4.173470, Y: 0.361388, Z: -11.217197}
		if !boneDeltas.GetByName(pmx.KNEE.Right()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.KNEE.Right()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.KNEE.Right()).FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: -1.217569, Y: 1.885731, Z: -10.788104}
		if !boneDeltas.GetByName(pmx.ANKLE.Right()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.ANKLE.Right()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.ANKLE.Right()).FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: -0.922247, Y: -1.163554, Z: -10.794323}
		if !boneDeltas.GetByName(pmx.TOE.Right()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.TOE.Right()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.TOE.Right()).FilledGlobalPosition().MMD()))
		}
	}
}

func TestVmdMotion_DeformLegIk17_Snow(t *testing.T) {
	// mlog.SetLevel(mlog.IK_VERBOSE)

	vr := repository.NewVmdRepository()
	motionData, err := vr.Load("../../../test_resources/好き雪_1075.vmd")

	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err)
	}

	motion := motionData.(*vmd.VmdMotion)

	pr := repository.NewPmxRepository()
	modelData, err := pr.Load("D:/MMD/MikuMikuDance_v926x64/UserFile/Model/VOCALOID/初音ミク/Lat式ミクVer2.31/Lat式ミクVer2.31_White_準標準.pmx")

	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err)
	}

	model := modelData.(*pmx.PmxModel)

	boneDeltas := DeformBone(model, motion, true, 0, []string{pmx.TOE.Right(), pmx.HEEL.Right()})

	{
		expectedPosition := &mmath.MVec3{X: 2.049998, Y: 12.957623, Z: 1.477440}
		if !boneDeltas.GetByName(pmx.LOWER.String()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.LOWER.String()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.LOWER.String()).FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: 1.201382, Y: 11.353215, Z: 2.266898}
		if !boneDeltas.GetByName(pmx.LEG.Right()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.LEG.Right()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.LEG.Right()).FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: 0.443043, Y: 7.640018, Z: -1.308741}
		if !boneDeltas.GetByName(pmx.KNEE.Right()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.KNEE.Right()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.KNEE.Right()).FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: -0.574753, Y: 7.943915, Z: 3.279809}
		if !boneDeltas.GetByName(pmx.ANKLE.Right()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.ANKLE.Right()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.ANKLE.Right()).FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: -1.443098, Y: 6.324932, Z: 4.837177}
		if !boneDeltas.GetByName(pmx.TOE.Right()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.TOE.Right()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.TOE.Right()).FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: -0.701516, Y: 8.181108, Z: 4.687274}
		if !boneDeltas.GetByName(pmx.HEEL.Right()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.HEEL.Right()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.HEEL.Right()).FilledGlobalPosition().MMD()))
		}
	}
}

func TestVmdMotion_DeformLegIk18_Syou(t *testing.T) {
	vr := repository.NewVmdRepository()
	motionData, err := vr.Load("../../../test_resources/唱(ダンスのみ)_0-300F.vmd")

	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err)
	}

	motion := motionData.(*vmd.VmdMotion)

	pr := repository.NewPmxRepository()
	modelData, err := pr.Load("../../../test_resources/サンプルモデル.pmx")

	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err)
	}

	model := modelData.(*pmx.PmxModel)

	boneDeltas := DeformBone(model, motion, true, 107, []string{pmx.TOE.Right(), pmx.HEEL.Right()})

	{
		expectedPosition := &mmath.MVec3{X: 0.365000, Y: 12.042871, Z: 2.034023}
		if !boneDeltas.GetByName(pmx.LOWER.String()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.LOWER.String()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.LOWER.String()).FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: -0.488466, Y: 10.920292, Z: 2.626419}
		if !boneDeltas.GetByName(pmx.LEG.Right()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.LEG.Right()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.LEG.Right()).FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: -1.607765, Y: 6.763937, Z: 1.653586}
		if !boneDeltas.GetByName(pmx.KNEE.Right()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.KNEE.Right()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.KNEE.Right()).FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: -1.110289, Y: 1.718307, Z: 2.809817}
		if !boneDeltas.GetByName(pmx.ANKLE.Right()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.ANKLE.Right()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.ANKLE.Right()).FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: -1.753089, Y: -0.026766, Z: 1.173958}
		if !boneDeltas.GetByName(pmx.TOE.Right()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.TOE.Right()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.TOE.Right()).FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: -0.952785, Y: 0.078826, Z: 2.838099}
		if !boneDeltas.GetByName(pmx.HEEL.Right()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.HEEL.Right()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.HEEL.Right()).FilledGlobalPosition().MMD()))
		}
	}
}

func TestVmdMotion_DeformLegIk19_Wa(t *testing.T) {
	vr := repository.NewVmdRepository()
	motionData, err := vr.Load("../../../test_resources/129cm_001_10F.vmd")

	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err)
	}

	motion := motionData.(*vmd.VmdMotion)

	pr := repository.NewPmxRepository()
	modelData, err := pr.Load("D:/MMD/MikuMikuDance_v926x64/UserFile/Model/_VMDサイジング/wa_129cm 20231028/wa_129cm_bone-structure.pmx")

	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err)
	}

	model := modelData.(*pmx.PmxModel)

	boneDeltas := DeformBone(model, motion, true, 0, []string{pmx.TOE.Right(), pmx.HEEL.Right()})
	{
		expectedPosition := &mmath.MVec3{X: 0.000000, Y: 9.900000, Z: 0.000000}
		if !boneDeltas.GetByName(pmx.LOWER.String()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.LOWER.String()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.LOWER.String()).FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: -0.599319, Y: 8.639606, Z: 0.369618}
		if !boneDeltas.GetByName(pmx.LEG.Right()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.LEG.Right()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.LEG.Right()).FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: -0.486516, Y: 6.323577, Z: -2.217865}
		if !boneDeltas.GetByName(pmx.KNEE.Right()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.KNEE.Right()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.KNEE.Right()).FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: -2.501665, Y: 2.859252, Z: -1.902513}
		if !boneDeltas.GetByName(pmx.ANKLE.Right()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.ANKLE.Right()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.ANKLE.Right()).FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: -3.071062, Y: 0.841962, Z: -2.077063}
		if !boneDeltas.GetByName(pmx.TOE.Right()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.TOE.Right()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.TOE.Right()).FilledGlobalPosition().MMD()))
		}
	}
}

func TestVmdMotion_DeformLegIk20_Syou(t *testing.T) {
	vr := repository.NewVmdRepository()
	motionData, err := vr.Load("../../../test_resources/唱(ダンスのみ)_0-300F.vmd")

	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err)
	}

	motion := motionData.(*vmd.VmdMotion)

	pr := repository.NewPmxRepository()
	modelData, err := pr.Load("../../../test_resources/サンプルモデル.pmx")

	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err)
	}

	model := modelData.(*pmx.PmxModel)

	boneDeltas := DeformBone(model, motion, true, 107, []string{pmx.TOE.Right(), pmx.HEEL.Right()})
	{
		expectedPosition := &mmath.MVec3{X: 0.365000, Y: 12.042871, Z: 2.034023}
		if !boneDeltas.GetByName(pmx.LOWER.String()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.LOWER.String()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.LOWER.String()).FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: -0.488466, Y: 10.920292, Z: 2.626419}
		if !boneDeltas.GetByName(pmx.LEG.Right()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.LEG.Right()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.LEG.Right()).FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: -1.607765, Y: 6.763937, Z: 1.653586}
		if !boneDeltas.GetByName(pmx.KNEE.Right()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.KNEE.Right()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.KNEE.Right()).FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: -1.110289, Y: 1.718307, Z: 2.809817}
		if !boneDeltas.GetByName(pmx.ANKLE.Right()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.ANKLE.Right()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.ANKLE.Right()).FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: -1.753089, Y: -0.026766, Z: 1.173958}
		if !boneDeltas.GetByName(pmx.TOE.Right()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.TOE.Right()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.TOE.Right()).FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: -0.952785, Y: 0.078826, Z: 2.838099}
		if !boneDeltas.GetByName(pmx.HEEL.Right()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.HEEL.Right()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.HEEL.Right()).FilledGlobalPosition().MMD()))
		}
	}
}

func TestVmdMotion_DeformLegIk21_FK(t *testing.T) {
	vr := repository.NewVmdRepository()
	motionData, err := vr.Load("../../../test_resources/足FK.vmd")

	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err)
	}

	motion := motionData.(*vmd.VmdMotion)

	pr := repository.NewPmxRepository()
	modelData, err := pr.Load("../../../test_resources/サンプルモデル_ひざ制限なし.pmx")

	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err)
	}

	model := modelData.(*pmx.PmxModel)

	boneDeltas := DeformBone(model, motion, false, 0, []string{pmx.TOE.Right(), pmx.HEEL.Right()})
	{
		expectedPosition := &mmath.MVec3{X: -0.133305, Y: 10.693993, Z: 2.314730}
		if !boneDeltas.GetByName(pmx.LEG.Right()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.LEG.Right()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.LEG.Right()).FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: 2.708069, Y: 9.216356, Z: -0.720822}
		if !boneDeltas.GetByName(pmx.KNEE.Right()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.KNEE.Right()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.KNEE.Right()).FilledGlobalPosition().MMD()))
		}
	}
}

func TestVmdMotion_DeformLegIk22_Bake(t *testing.T) {
	vr := repository.NewVmdRepository()
	motionData, err := vr.Load("../../../test_resources/足FK焼き込み.vmd")

	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err)
	}

	motion := motionData.(*vmd.VmdMotion)

	pr := repository.NewPmxRepository()
	modelData, err := pr.Load("../../../test_resources/サンプルモデル_ひざ制限なし.pmx")

	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err)
	}

	model := modelData.(*pmx.PmxModel)

	boneDeltas := DeformBone(model, motion, true, 0, []string{pmx.TOE.Right(), pmx.HEEL.Right()})
	{
		expectedPosition := &mmath.MVec3{X: -0.133306, Y: 10.693994, Z: 2.314731}
		if !boneDeltas.GetByName(pmx.LEG.Right()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.LEG.Right()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.LEG.Right()).FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: -3.753989, Y: 8.506582, Z: 1.058842}
		if !boneDeltas.GetByName(pmx.KNEE.Right()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.KNEE.Right()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.KNEE.Right()).FilledGlobalPosition().MMD()))
		}
	}
}

func TestVmdMotion_DeformLegIk22_NoLimit(t *testing.T) {
	vr := repository.NewVmdRepository()
	motionData, err := vr.Load("../../../test_resources/足FK.vmd")

	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err)
	}

	motion := motionData.(*vmd.VmdMotion)

	pr := repository.NewPmxRepository()
	modelData, err := pr.Load("../../../test_resources/サンプルモデル_ひざ制限なし.pmx")

	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err)
	}

	model := modelData.(*pmx.PmxModel)

	boneDeltas := DeformBone(model, motion, true, 0, []string{pmx.TOE.Right(), pmx.HEEL.Right()})
	{
		expectedPosition := &mmath.MVec3{X: -0.133305, Y: 10.693993, Z: 2.314730}
		if !boneDeltas.GetByName(pmx.LEG.Right()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.LEG.Right()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.LEG.Right()).FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: 2.081436, Y: 7.884178, Z: -0.268146}
		if !boneDeltas.GetByName(pmx.KNEE.Right()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.KNEE.Right()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.KNEE.Right()).FilledGlobalPosition().MMD()))
		}
	}
}

func TestVmdMotion_DeformLegIk23_Addiction(t *testing.T) {
	vr := repository.NewVmdRepository()
	motionData, err := vr.Load("../../../test_resources/[A]ddiction_Lat式_0171F.vmd")

	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err)
	}

	motion := motionData.(*vmd.VmdMotion)

	pr := repository.NewPmxRepository()
	modelData, err := pr.Load("D:/MMD/MikuMikuDance_v926x64/UserFile/Model/VOCALOID/初音ミク/Tda式ミクワンピース/Tda式ミクワンピースRSP.pmx")

	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err)
	}

	model := modelData.(*pmx.PmxModel)

	boneDeltas := DeformBone(model, motion, true, 0, []string{pmx.TOE_IK.Right(), pmx.TOE.Right()})

	{
		expectedPosition := &mmath.MVec3{X: 0, Y: 0.2593031, Z: 0}
		if !boneDeltas.GetByName(pmx.ROOT.String()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.ROOT.String()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.ROOT.String()).FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: -1.528317, Y: 5.033707, Z: 3.125487}
		if !boneDeltas.GetByName(pmx.LEG_IK.Right()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.LEG_IK.Right()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.LEG_IK.Right()).FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: 0.609285, Y: 12.001350, Z: 1.666402}
		if !boneDeltas.GetByName(pmx.LOWER.String()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.LOWER.String()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.LOWER.String()).FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: -0.129098, Y: 10.550634, Z: 1.348259}
		if !boneDeltas.GetByName(pmx.LEG.Right()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.LEG.Right()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.LEG.Right()).FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: 1.661012, Y: 6.604201, Z: -1.196993}
		if !boneDeltas.GetByName(pmx.KNEE.Right()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.KNEE.Right()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.KNEE.Right()).FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: -1.529553, Y: 5.033699, Z: 3.127081}
		if !boneDeltas.GetByName(pmx.ANKLE.Right()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.ANKLE.Right()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.ANKLE.Right()).FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: -2.044619, Y: 3.204468, Z: 2.877363}
		if !boneDeltas.GetByName(pmx.TOE.Right()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.TOE.Right()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.TOE.Right()).FilledGlobalPosition().MMD()))
		}
	}
}

func TestVmdMotion_DeformLegIk24_Positive(t *testing.T) {
	vr := repository.NewVmdRepository()
	motionData, err := vr.Load("../../../test_resources/ポジティブパレード_0526.vmd")

	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err)
	}

	motion := motionData.(*vmd.VmdMotion)

	pr := repository.NewPmxRepository()
	modelData, err := pr.Load("D:/MMD/MikuMikuDance_v926x64/UserFile/Model/_VMDサイジング/wa_129cm 20231028/wa_129cm_20240406.pmx")

	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err)
	}

	model := modelData.(*pmx.PmxModel)
	boneDeltas := DeformBone(model, motion, true, 0, nil)
	{
		expectedPosition := &mmath.MVec3{X: 0, Y: 0, Z: 0}
		if !boneDeltas.GetByName(pmx.ROOT.String()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.ROOT.String()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.ROOT.String()).FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: -3.312041, Y: 6.310613, Z: -1.134230}
		if !boneDeltas.GetByName(pmx.LEG_IK.Left()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.LEG_IK.Right()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.LEG_IK.Right()).FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: -2.754258, Y: 7.935882, Z: -2.298871}
		if !boneDeltas.GetByName(pmx.LOWER.String()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.LOWER.String()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.LOWER.String()).FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: -2.455364, Y: 6.571013, Z: -1.935295}
		if !boneDeltas.GetByName(pmx.LEG.Left()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.LEG.Right()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.LEG.Right()).FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: -2.695464, Y: 4.323516, Z: -4.574024}
		if !boneDeltas.GetByName(pmx.KNEE.Left()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.KNEE.Right()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.KNEE.Right()).FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: -3.322137, Y: 6.302598, Z: -1.131305}
		if !boneDeltas.GetByName("左脛骨").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左脛骨").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左脛骨").FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: -2.575414, Y: 5.447266, Z: -3.254661}
		if !boneDeltas.GetByName("左足捩").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左足捩").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左足捩").FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: -2.229677, Y: 5.626327, Z: -3.481028}
		if !boneDeltas.GetByName("左足捩先").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左足捩先").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左足捩先").FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: -2.455364, Y: 6.571013, Z: -1.935295}
		if !boneDeltas.GetByName("左足向検A").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左足向検A").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左足向検A").FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: -2.695177, Y: 4.324148, Z: -4.574588}
		if !boneDeltas.GetByName("左足向検A先").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左足向検A先").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左足向検A先").FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: -2.695177, Y: 4.324148, Z: -4.574588}
		if !boneDeltas.GetByName("左足捩検B").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左足捩検B").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左足捩検B").FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: -0.002697, Y: 5.869486, Z: -6.134800}
		if !boneDeltas.GetByName("左足捩検B先").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左足捩検B先").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左足捩検B先").FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: -2.877639, Y: 4.4450495, Z: -4.164494}
		if !boneDeltas.GetByName("左膝補").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左膝補").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左膝補").FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: -3.523895, Y: 4.135535, Z: -3.716305}
		if !boneDeltas.GetByName("左膝補先").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左膝補先").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左膝補先").FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: -2.118768, Y: 6.263350, Z: -2.402574}
		if !boneDeltas.GetByName("左足w").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左足w").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左足w").FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: -2.480717, Y: 3.120446, Z: -5.602753}
		if !boneDeltas.GetByName("左足w先").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左足w先").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左足w先").FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: -2.455364, Y: 6.571013, Z: -1.935294}
		if !boneDeltas.GetByName("左足向-").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左足向-").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左足向-").FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: -3.322137, Y: 6.302598, Z: -1.131305}
		if !boneDeltas.GetByName("左脛骨D").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左脛骨D").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左脛骨D").FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: -3.199167, Y: 3.952319, Z: -4.391296}
		if !boneDeltas.GetByName("左脛骨D先").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左脛骨D先").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左脛骨D先").FilledGlobalPosition().MMD()))
		}
	}
}

func TestVmdMotion_DeformArmIk(t *testing.T) {
	vr := repository.NewVmdRepository()
	motionData, err := vr.Load("../../../test_resources/サンプルモーション.vmd")

	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err)
	}

	motion := motionData.(*vmd.VmdMotion)

	pr := repository.NewPmxRepository()
	modelData, err := pr.Load("../../../test_resources/ボーンツリーテストモデル.pmx")

	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err)
	}

	model := modelData.(*pmx.PmxModel)

	boneDeltas := DeformBone(model, motion, true, 3182, nil)
	{
		expectedPosition := &mmath.MVec3{X: 0, Y: 0, Z: 0}
		if !boneDeltas.GetByName(pmx.ROOT.String()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.ROOT.String()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.ROOT.String()).FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: 12.400011, Y: 9.000000, Z: 1.885650}
		if !boneDeltas.GetByName(pmx.CENTER.String()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.CENTER.String()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.CENTER.String()).FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: 12.400011, Y: 8.580067, Z: 1.885650}
		if !boneDeltas.GetByName(pmx.GROOVE.String()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.GROOVE.String()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.GROOVE.String()).FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: 12.400011, Y: 11.628636, Z: 2.453597}
		if !boneDeltas.GetByName(pmx.WAIST.String()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.WAIST.String()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.WAIST.String()).FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: 12.400011, Y: 12.567377, Z: 1.229520}
		if !boneDeltas.GetByName(pmx.UPPER.String()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.UPPER.String()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.UPPER.String()).FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: 12.344202, Y: 13.782951, Z: 1.178849}
		if !boneDeltas.GetByName(pmx.UPPER2.String()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.UPPER2.String()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.UPPER2.String()).FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: 12.425960, Y: 15.893852, Z: 1.481421}
		if !boneDeltas.GetByName(pmx.SHOULDER_P.Left()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.SHOULDER_P.Left()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.SHOULDER_P.Left()).FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: 12.425960, Y: 15.893852, Z: 1.481421}
		if !boneDeltas.GetByName(pmx.SHOULDER.Left()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.SHOULDER.Left()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.SHOULDER.Left()).FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: 13.348320, Y: 15.767927, Z: 1.802947}
		if !boneDeltas.GetByName(pmx.ARM.Left()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.ARM.Left()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.ARM.Left()).FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: 13.564770, Y: 14.998386, Z: 1.289923}
		if !boneDeltas.GetByName(pmx.ARM_TWIST.Left()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.ARM_TWIST.Left()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.ARM_TWIST.Left()).FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: 14.043257, Y: 13.297290, Z: 0.155864}
		if !boneDeltas.GetByName(pmx.ELBOW.Left()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.ELBOW.Left()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.ELBOW.Left()).FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: 13.811955, Y: 13.552182, Z: -0.388005}
		if !boneDeltas.GetByName(pmx.WRIST_TWIST.Left()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.WRIST_TWIST.Left()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.WRIST_TWIST.Left()).FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: 13.144803, Y: 14.287374, Z: -1.956703}
		if !boneDeltas.GetByName(pmx.WRIST.Left()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.WRIST.Left()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.WRIST.Left()).FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: 12.813587, Y: 14.873419, Z: -2.570278}
		if !boneDeltas.GetByName(pmx.INDEX1.Left()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.INDEX1.Left()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.INDEX1.Left()).FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: 12.541822, Y: 15.029200, Z: -2.709604}
		if !boneDeltas.GetByName(pmx.INDEX2.Left()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.INDEX2.Left()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.INDEX2.Left()).FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: 12.476499, Y: 14.950351, Z: -2.502167}
		if !boneDeltas.GetByName(pmx.INDEX3.Left()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.INDEX3.Left()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.INDEX3.Left()).FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: 12.620306, Y: 14.795185, Z: -2.295859}
		if !boneDeltas.GetByName("左人指先").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左人指先").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左人指先").FilledGlobalPosition().MMD()))
		}
	}
}

func TestVmdMotion_DeformLegIk2(t *testing.T) {
	vr := repository.NewVmdRepository()
	motionData, err := vr.Load("C:/MMD/mmd_base/tests/resources/唱(ダンスのみ)_0274F.vmd")

	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err)
	}

	motion := motionData.(*vmd.VmdMotion)

	pr := repository.NewPmxRepository()
	modelData, err := pr.Load("D:/MMD/MikuMikuDance_v926x64/UserFile/Model/VOCALOID/初音ミク/ISAO式ミク/I_ミクv4/Miku_V4_準標準.pmx")

	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err)
	}

	model := modelData.(*pmx.PmxModel)

	boneDeltas := DeformBone(model, motion, true, 0, nil)

	{
		expectedPosition := &mmath.MVec3{X: 0.04952335, Y: 9.0, Z: 1.72378033}
		if !boneDeltas.GetByName(pmx.CENTER.String()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.CENTER.String()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.CENTER.String()).FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: 0.04952335, Y: 7.97980869, Z: 1.72378033}
		if !boneDeltas.GetByName(pmx.GROOVE.String()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.GROOVE.String()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.GROOVE.String()).FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: 0.04952335, Y: 11.02838314, Z: 2.29172656}
		if !boneDeltas.GetByName("腰").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("腰").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("腰").FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: 0.04952335, Y: 11.9671191, Z: 1.06765032}
		if !boneDeltas.GetByName("下半身").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("下半身").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("下半身").FilledGlobalPosition().MMD()))
		}
	}
	// FIXME: 物理後なので求められない
	// {
	// 	expectedPosition := &mmath.MVec3{X: -0.24102019, Y:9.79926074,Z:1.08498769}
	// 	if !boneDeltas.GetByName("下半身先").GetGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
	// 		t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("下半身先").GetGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("下半身先").GetGlobalPosition().MMD()))
	// 	}
	// }
	{
		expectedPosition := &mmath.MVec3{X: 0.90331914, Y: 10.27362702, Z: 1.009499759}
		if !boneDeltas.GetByName("腰キャンセル左").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("腰キャンセル左").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("腰キャンセル左").FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: 0.90331914, Y: 10.27362702, Z: 1.00949975}
		if !boneDeltas.GetByName("左足").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左足").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左足").FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: 0.08276818, Y: 5.59348757, Z: -1.24981795}
		if !boneDeltas.GetByName("左ひざ").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左ひざ").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左ひざ").FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: 5.63290634e-01, Y: -2.12439821e-04, Z: -3.87768478e-01}
		if !boneDeltas.GetByName("左つま先").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左つま先").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左つま先").FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: 0.90331914, Y: 10.27362702, Z: 1.00949975}
		if !boneDeltas.GetByName("左足D").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左足D").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左足D").FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: 0.23453057, Y: 5.6736954, Z: -0.76228439}
		if !boneDeltas.GetByName("左ひざ2").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左ひざ2").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左ひざ2").FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: 0.12060311, Y: 4.95396153, Z: -1.23761938}
		if !boneDeltas.GetByName("左ひざ2先").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左ひざ2先").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左ひざ2先").FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: 0.90331914, Y: 10.27362702, Z: 1.00949975}
		if !boneDeltas.GetByName("左足y+").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左足y+").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左足y+").FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: 0.74736036, Y: 9.38409308, Z: 0.58008117}
		if !boneDeltas.GetByName("左足yTgt").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左足yTgt").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左足yTgt").FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: 0.74736036, Y: 9.38409308, Z: 0.58008117}
		if !boneDeltas.GetByName("左足yIK").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左足yIK").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左足yIK").FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: -0.03018836, Y: 10.40081089, Z: 1.26859617}
		if !boneDeltas.GetByName("左尻").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左尻").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左尻").FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: 0.08276818, Y: 5.59348757, Z: -1.24981795}
		if !boneDeltas.GetByName("左ひざsub").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左ひざsub").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左ひざsub").FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: -0.09359026, Y: 5.54494997, Z: -1.80895985}
		if !boneDeltas.GetByName("左ひざsub先").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左ひざsub先").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左ひざsub先").FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: 1.23779916, Y: 1.28891465, Z: 1.65257835}
		if !boneDeltas.GetByName("左ひざD2").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左ひざD2").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左ひざD2").FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: 0.1106881, Y: 4.98643066, Z: -1.26321915}
		if !boneDeltas.GetByName("左ひざD2先").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左ひざD2先").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左ひざD2先").FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: 0.12060311, Y: 4.95396153, Z: -1.23761938}
		if !boneDeltas.GetByName("左ひざD2IK").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左ひざD2IK").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左ひざD2IK").FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: 0.88590917, Y: 0.38407067, Z: 0.56801614}
		if !boneDeltas.GetByName("左足ゆび").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左足ゆび").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左足ゆび").FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: 5.63290634e-01, Y: -2.12439821e-04, Z: -3.87768478e-01}
		if !boneDeltas.GetByName("左つま先D").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左つま先D").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左つま先D").FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: 0.90331914, Y: 10.27362702, Z: 1.00949975}
		if !boneDeltas.GetByName("左足D").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左足D").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左足D").FilledGlobalPosition().MMD()))
		}
	}
	// {
	// 	expectedPosition := &mmath.MVec3{X: 0.08276818, Y:5.59348757,Z:-1.24981795}
	// 	if !boneDeltas.GetByName("左ひざD").GetGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
	// 		t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左ひざD").GetGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左ひざD").GetGlobalPosition().MMD()))
	// 	}
	// }
	// {
	// 	expectedPosition := &mmath.MVec3{X: 1.23779916, Y:1.28891465,Z:1.65257835}
	// 	if !boneDeltas.GetByName("左足首D").GetGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
	// 		t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左足首D").GetGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左足首D").GetGlobalPosition().MMD()))
	// 	}
	// }
}

func TestVmdMotion_DeformArmIk3(t *testing.T) {
	vr := repository.NewVmdRepository()
	motionData, err := vr.Load("C:/MMD/mlib_go/test_resources/Addiction_0F.vmd")

	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err)
	}

	motion := motionData.(*vmd.VmdMotion)

	pr := repository.NewPmxRepository()
	modelData, err := pr.Load("D:/MMD/MikuMikuDance_v926x64/UserFile/Model/VOCALOID/初音ミク/Sour式初音ミクVer.1.02/Black_全表示.pmx")

	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err)
	}

	model := modelData.(*pmx.PmxModel)

	boneDeltas := DeformBone(model, motion, true, 0, nil)
	{
		expectedPosition := &mmath.MVec3{X: 1.018832, Y: 15.840092, Z: 0.532239}
		if !boneDeltas.GetByName("左腕").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左腕").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左腕").FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: 1.186002, Y: 14.510550, Z: 0.099023}
		if !boneDeltas.GetByName("左腕捩").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左腕捩").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左腕捩").FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: 1.353175, Y: 13.181011, Z: -0.334196}
		if !boneDeltas.GetByName("左ひじ").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左ひじ").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左ひじ").FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: 1.018832, Y: 15.840092, Z: 0.532239}
		if !boneDeltas.GetByName("左腕W").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左腕W").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左腕W").FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: 1.353175, Y: 13.181011, Z: -0.334196}
		if !boneDeltas.GetByName("左腕W先").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左腕W先").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左腕W先").FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: 1.353175, Y: 13.181011, Z: -0.334196}
		if !boneDeltas.GetByName("左腕WIK").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左腕WIK").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左腕WIK").FilledGlobalPosition().MMD()))
		}
	}
}

func TestVmdMotion_DeformLegIk25_Ballet(t *testing.T) {
	// mlog.SetLevel(mlog.IK_VERBOSE)

	vr := repository.NewVmdRepository()
	motionData, err := vr.Load("../../../test_resources/青江バレリーコ_1543F.vmd")

	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err)
	}

	motion := motionData.(*vmd.VmdMotion)

	pr := repository.NewPmxRepository()
	modelData, err := pr.Load("D:/MMD/MikuMikuDance_v926x64/UserFile/Model/刀剣乱舞/019_にっかり青江/にっかり青江 帽子屋式 ver2.1/帽子屋式にっかり青江（戦装束）_表示枠.pmx")

	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err)
	}

	model := modelData.(*pmx.PmxModel)

	boneDeltas := DeformBone(model, motion, true, 0, []string{pmx.TOE.Left(), pmx.HEEL.Left(), pmx.TOE_EX.Left()})

	{
		expectedPosition := &mmath.MVec3{X: -4.374956, Y: 13.203792, Z: 1.554190}
		if !boneDeltas.GetByName(pmx.LOWER.String()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.LOWER.String()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.LOWER.String()).FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: -3.481956, Y: 11.214747, Z: 1.127255}
		if !boneDeltas.GetByName(pmx.LEG.Left()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.LEG.Left()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.LEG.Left()).FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: -7.173243, Y: 7.787793, Z: 0.013533}
		if !boneDeltas.GetByName(pmx.KNEE.Left()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.KNEE.Left()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.KNEE.Left()).FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: -11.529483, Y: 3.689184, Z: -1.119154}
		if !boneDeltas.GetByName(pmx.ANKLE.Left()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.ANKLE.Left()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.ANKLE.Left()).FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: -13.408189, Y: 1.877100, Z: -2.183821}
		if !boneDeltas.GetByName(pmx.TOE.Left()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.TOE.Left()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.TOE.Left()).FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: -12.545708, Y: 4.008257, Z: -0.932670}
		if !boneDeltas.GetByName(pmx.HEEL.Left()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.HEEL.Left()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.HEEL.Left()).FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: -3.481956, Y: 11.214747, Z: 1.127255}
		if !boneDeltas.GetByName(pmx.LEG_D.Left()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.LEG_D.Left()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.LEG_D.Left()).FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: -7.173243, Y: 7.787793, Z: 0.013533}
		if !boneDeltas.GetByName(pmx.KNEE_D.Left()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.KNEE_D.Left()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.KNEE_D.Left()).FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: -11.529483, Y: 3.689184, Z: -1.119154}
		if !boneDeltas.GetByName(pmx.ANKLE_D.Left()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.ANKLE_D.Left()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.ANKLE_D.Left()).FilledGlobalPosition().MMD()))
		}
	}
	// {
	// 	expectedPosition := &mmath.MVec3{X: -12.845280, Y:2.816309,Z:-2.136874}
	// 	if !boneDeltas.GetByName(pmx.TOE_EX.Left()).GetGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
	// 		t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.TOE_EX.Left()).GetGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.TOE_EX.Left()).GetGlobalPosition().MMD()))
	// 	}
	// }
}

func TestVmdMotion_DeformLegIk26_Far(t *testing.T) {
	// mlog.SetLevel(mlog.IK_VERBOSE)

	vr := repository.NewVmdRepository()
	motionData, err := vr.Load("../../../test_resources/足IK乖離.vmd")

	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err)
	}

	motion := motionData.(*vmd.VmdMotion)

	pr := repository.NewPmxRepository()
	modelData, err := pr.Load("D:/MMD/MikuMikuDance_v926x64/UserFile/Model/_あにまさ式ミク準標準見せパン/初音ミクVer2 準標準 見せパン 3.pmx")

	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err)
	}

	model := modelData.(*pmx.PmxModel)
	boneDeltas := DeformBone(model, motion, true, 0, []string{pmx.TOE.Right(), pmx.TOE_EX.Right(), pmx.HEEL.Right()})

	{
		expectedPosition := &mmath.MVec3{X: -0.796811, Y: 10.752734, Z: -0.072743}
		if !boneDeltas.GetByName(pmx.LEG.Right()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.LEG.Right()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.LEG.Right()).FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: -2.202487, Y: 10.921064, Z: -4.695134}
		if !boneDeltas.GetByName(pmx.KNEE.Right()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.KNEE.Right()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.KNEE.Right()).FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: -4.193142, Y: 11.026311, Z: -8.844866}
		if !boneDeltas.GetByName(pmx.ANKLE.Right()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.ANKLE.Right()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.ANKLE.Right()).FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: -5.108798, Y: 10.935530, Z: -11.494570}
		if !boneDeltas.GetByName(pmx.TOE.Right()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.TOE.Right()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.TOE.Right()).FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: -4.800813, Y: 10.964218, Z: -10.612234}
		if !boneDeltas.GetByName(pmx.TOE_EX.Right()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.TOE_EX.Right()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.TOE_EX.Right()).FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: -4.331888, Y: 12.178923, Z: -9.514071}
		if !boneDeltas.GetByName(pmx.HEEL.Right()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.HEEL.Right()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.HEEL.Right()).FilledGlobalPosition().MMD()))
		}
	}
}

func TestVmdMotion_DeformLegIk27_Addiction_Shoes(t *testing.T) {
	// mlog.SetLevel(mlog.IK_VERBOSE)

	vr := repository.NewVmdRepository()
	motionData, err := vr.Load("../../../test_resources/[A]ddiction_和洋_1074-1078F.vmd")

	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err)
	}

	motion := motionData.(*vmd.VmdMotion)

	pr := repository.NewPmxRepository()
	modelData, err := pr.Load("D:/MMD/MikuMikuDance_v926x64/UserFile/Model/_VMDサイジング/wa_129cm 20231028/wa_129cm_20240406.pmx")

	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err)
	}

	model := modelData.(*pmx.PmxModel)

	boneDeltas := DeformBone(model, motion, true, 2, nil)
	{
		expectedPosition := &mmath.MVec3{X: 0, Y: 0, Z: 0}
		if !boneDeltas.GetByName(pmx.ROOT.String()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.ROOT.String()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.ROOT.String()).FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: 1.406722, Y: 1.841236, Z: 0.277818}
		if !boneDeltas.GetByName(pmx.LEG_IK.Left()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.LEG_IK.Right()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.LEG_IK.Right()).FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: 0.510231, Y: 9.009953, Z: 0.592482}
		if !boneDeltas.GetByName(pmx.LOWER.String()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.LOWER.String()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.LOWER.String()).FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: 1.355914, Y: 7.853320, Z: 0.415251}
		if !boneDeltas.GetByName(pmx.LEG.Left()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.LEG.Right()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.LEG.Right()).FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: -0.327781, Y: 5.203806, Z: -1.073718}
		if !boneDeltas.GetByName(pmx.KNEE.Left()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.KNEE.Right()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.KNEE.Right()).FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: 1.407848, Y: 1.839228, Z: 0.278700}
		if !boneDeltas.GetByName("左脛骨").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左脛骨").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左脛骨").FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: 1.407848, Y: 1.839228, Z: 0.278700}
		if !boneDeltas.GetByName("左脛骨D").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左脛骨D").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左脛骨D").FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: -0.498054, Y: 5.045506, Z: -1.221016}
		if !boneDeltas.GetByName("左脛骨D先").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左脛骨D先").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左脛骨D先").FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: 1.462306, Y: 7.684025, Z: 0.087026}
		if !boneDeltas.GetByName("左足Dw").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左足Dw").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左足Dw").FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: 1.593721, Y: 0.784840, Z: -0.054141}
		if !boneDeltas.GetByName("左足先EX").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左足先EX").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左足先EX").FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: 1.551940, Y: 1.045847, Z: 0.034003}
		if !boneDeltas.GetByName("左素足先A").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左素足先A").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左素足先A").FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: 1.453982, Y: 0.305976, Z: -0.510022}
		if !boneDeltas.GetByName("左素足先A先").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左素足先A先").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左素足先A先").FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: 1.453982, Y: 0.305976, Z: -0.510022}
		if !boneDeltas.GetByName("左素足先AIK").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左素足先AIK").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左素足先AIK").FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: 0.941880, Y: 2.132958, Z: 0.020403}
		if !boneDeltas.GetByName("左素足先B").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左素足先B").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左素足先B").FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: 1.359364, Y: 0.974298, Z: -0.226041}
		if !boneDeltas.GetByName("左素足先B先").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左素足先B先").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左素足先B先").FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: 1.460890, Y: 0.692527, Z: -0.285973}
		if !boneDeltas.GetByName("左素足先BIK").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左素足先BIK").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左素足先BIK").FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: 1.173929, Y: 2.066327, Z: 0.182685}
		if !boneDeltas.GetByName("左靴調節").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左靴調節").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左靴調節").FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: 1.739235, Y: 1.171441, Z: 0.485052}
		if !boneDeltas.GetByName("左靴追従").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左靴追従").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左靴追従").FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: 1.186359, Y: 2.046771, Z: 0.189367}
		if !boneDeltas.GetByName("左靴追従先").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左靴追従先").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左靴追従先").FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: 1.173929, Y: 2.066327, Z: 0.182685}
		if !boneDeltas.GetByName("左靴追従IK").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左靴追従IK").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左靴追従IK").FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: 1.574899, Y: 6.873434, Z: 0.342768}
		if !boneDeltas.GetByName("左足補D").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左足補D").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左足補D").FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: 0.150401, Y: 5.170907, Z: -0.712416}
		if !boneDeltas.GetByName("左足補D先").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左足補D先").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左足補D先").FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: 0.150401, Y: 5.170907, Z: -0.712416}
		if !boneDeltas.GetByName("左足補DIK").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左足補DIK").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左足補DIK").FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: 1.355915, Y: 7.853319, Z: 0.415251}
		if !boneDeltas.GetByName("左足向検A").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左足向検A").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左足向検A").FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: -0.327781, Y: 5.203805, Z: -1.073719}
		if !boneDeltas.GetByName("左足向検A先").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左足向検A先").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左足向検A先").FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: -0.327781, Y: 5.203805, Z: -1.073719}
		if !boneDeltas.GetByName("左足向検AIK").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左足向検AIK").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左足向検AIK").FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: 1.355914, Y: 7.853319, Z: 0.415251}
		if !boneDeltas.GetByName("左足向-").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左足向-").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左足向-").FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: 1.264808, Y: 7.561551, Z: -0.161703}
		if !boneDeltas.GetByName("左足w").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左足w").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左足w").FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: -0.714029, Y: 3.930234, Z: -1.935889}
		if !boneDeltas.GetByName("左足w先").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左足w先").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左足w先").FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: 0.016770, Y: 5.319929, Z: -0.781771}
		if !boneDeltas.GetByName("左膝補").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左膝補").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左膝補").FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: -0.164672, Y: 4.511360, Z: -0.957886}
		if !boneDeltas.GetByName("左膝補先").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左膝補先").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左膝補先").FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: -0.099887, Y: 4.800064, Z: -0.895003}
		if !boneDeltas.GetByName("左膝補IK").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左膝補IK").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左膝補IK").FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: -0.327781, Y: 5.203806, Z: -1.073718}
		if !boneDeltas.GetByName("左足捩検B").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左足捩検B").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左足捩検B").FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: -2.392915, Y: 7.450026, Z: -2.735495}
		if !boneDeltas.GetByName("左足捩検B先").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左足捩検B先").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左足捩検B先").FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: -2.392915, Y: 7.450026, Z: -2.735495}
		if !boneDeltas.GetByName("左足捩検BIK").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左足捩検BIK").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左足捩検BIK").FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: 0.514067, Y: 6.528563, Z: -0.329234}
		if !boneDeltas.GetByName("左足捩").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左足捩").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左足捩").FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: 0.231636, Y: 6.794109, Z: -0.557747}
		if !boneDeltas.GetByName("左足捩先").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左足捩先").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左足捩先").FilledGlobalPosition().MMD()))
		}
	}
}
