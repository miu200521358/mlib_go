package deform

import (
	"testing"

	"github.com/miu200521358/mlib_go/pkg/domain/mmath"
	"github.com/miu200521358/mlib_go/pkg/domain/pmx"
	"github.com/miu200521358/mlib_go/pkg/domain/vmd"
	"github.com/miu200521358/mlib_go/pkg/infrastructure/repository"
	"github.com/miu200521358/mlib_go/pkg/mutils/mlog"
)

func TestVmdMotion_DeformArmIk_Mahoujin_01(t *testing.T) {
	mlog.SetLevel(mlog.IK_VERBOSE)

	vr := repository.NewVmdRepository()
	motionData, err := vr.Load("../../../test_resources/arm_ik_mahoujin_001F.vmd")

	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err)
	}

	motion := motionData.(*vmd.VmdMotion)

	pr := repository.NewPmxRepository()
	modelData, err := pr.Load("D:/MMD/MikuMikuDance_v926x64/UserFile/Model/刀剣乱舞/107_髭切/髭切mkmk009c 刀剣乱舞/髭切mkmk009c/髭切上着無mkmk009b_腕ＩＫ2.pmx")

	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err)
	}

	model := modelData.(*pmx.PmxModel)

	boneDeltas := DeformBone(model, motion, true, 0, nil)
	{
		boneName := pmx.ARM.Right()
		expectedPosition := &mmath.MVec3{X: -1.801768, Y: 18.555544, Z: 0.482812}
		if !boneDeltas.GetByName(boneName).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(boneName).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(boneName).FilledGlobalPosition().MMD()))
		}
	}
	{
		boneName := pmx.ELBOW.Right()
		expectedPosition := &mmath.MVec3{X: -4.091116, Y: 18.629446, Z: -1.670793}
		if !boneDeltas.GetByName(boneName).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(boneName).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(boneName).FilledGlobalPosition().MMD()))
		}
	}
	{
		boneName := pmx.WRIST.Right()
		expectedPosition := &mmath.MVec3{X: -6.370411, Y: 18.910606, Z: -4.062796}
		if !boneDeltas.GetByName(boneName).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(boneName).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(boneName).FilledGlobalPosition().MMD()))
		}
	}
	{
		boneName := pmx.INDEX3.Right()
		expectedPosition := &mmath.MVec3{X: -7.256862, Y: 18.269156, Z: -5.428672}
		if !boneDeltas.GetByName(boneName).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(boneName).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(boneName).FilledGlobalPosition().MMD()))
		}
	}
}

func TestVmdMotion_DeformArmIk_Mahoujin_02(t *testing.T) {
	mlog.SetLevel(mlog.IK_VERBOSE)

	vr := repository.NewVmdRepository()
	motionData, err := vr.Load("../../../test_resources/arm_ik_mahoujin_006F.vmd")

	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err)
	}

	motion := motionData.(*vmd.VmdMotion)

	pr := repository.NewPmxRepository()
	modelData, err := pr.Load("D:/MMD/MikuMikuDance_v926x64/UserFile/Model/刀剣乱舞/107_髭切/髭切mkmk009c 刀剣乱舞/髭切mkmk009c/髭切上着無mkmk009b_腕ＩＫ2.pmx")

	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err)
	}

	model := modelData.(*pmx.PmxModel)

	boneDeltas := DeformBone(model, motion, true, 0, nil)
	{
		boneName := pmx.ARM.Right()
		expectedPosition := &mmath.MVec3{X: -1.801768, Y: 18.555544, Z: 0.482812}
		if !boneDeltas.GetByName(boneName).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(boneName).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(boneName).FilledGlobalPosition().MMD()))
		}
	}
	{
		boneName := pmx.ELBOW.Right()
		expectedPosition := &mmath.MVec3{X: -3.273916, Y: 17.405672, Z: -2.046059}
		if !boneDeltas.GetByName(boneName).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(boneName).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(boneName).FilledGlobalPosition().MMD()))
		}
	}
	{
		boneName := pmx.WRIST.Right()
		expectedPosition := &mmath.MVec3{X: -1.240410, Y: 18.910606, Z: -4.062796}
		if !boneDeltas.GetByName(boneName).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(boneName).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(boneName).FilledGlobalPosition().MMD()))
		}
	}
	{
		boneName := pmx.INDEX3.Right()
		expectedPosition := &mmath.MVec3{X: -0.614190, Y: 19.042362, Z: -5.691705}
		if !boneDeltas.GetByName(boneName).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(boneName).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(boneName).FilledGlobalPosition().MMD()))
		}
	}
}
