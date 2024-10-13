package deform

import (
	"testing"

	"github.com/miu200521358/mlib_go/pkg/domain/mmath"
	"github.com/miu200521358/mlib_go/pkg/domain/pmx"
	"github.com/miu200521358/mlib_go/pkg/domain/vmd"
	"github.com/miu200521358/mlib_go/pkg/infrastructure/repository"
	"github.com/miu200521358/mlib_go/pkg/mutils/mlog"
)

func TestVmdMotion_DeformLegIk30_Addiction_Shoes(t *testing.T) {
	// mlog.SetLevel(mlog.IK_VERBOSE)

	vr := repository.NewVmdRepository()
	motionData, err := vr.Load("../../../test_resources/[A]ddiction_和洋_1037F.vmd")

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
		expectedPosition := &mmath.MVec3{X: 1.748025, Y: 1.683590, Z: 0.556993}
		if !boneDeltas.GetByName(pmx.LEG_IK.Left()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.LEG_IK.Right()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.LEG_IK.Right()).FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: 1.111190, Y: 4.955496, Z: 1.070225}
		if !boneDeltas.GetByName(pmx.LOWER.String()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.LOWER.String()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.LOWER.String()).FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: 1.724965, Y: 3.674735, Z: 0.810759}
		if !boneDeltas.GetByName(pmx.LEG.Left()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.LEG.Right()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.LEG.Right()).FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: -0.924419, Y: 3.943549, Z: -1.420897}
		if !boneDeltas.GetByName(pmx.KNEE.Left()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.KNEE.Right()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.KNEE.Right()).FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: 1.743128, Y: 1.672784, Z: 0.551317}
		if !boneDeltas.GetByName("左脛骨").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左脛骨").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左脛骨").FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: 1.743128, Y: 1.672784, Z: 0.551317}
		if !boneDeltas.GetByName("左脛骨D").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左脛骨D").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左脛骨D").FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: -1.118480, Y: 3.432016, Z: -1.657329}
		if !boneDeltas.GetByName("左脛骨D先").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左脛骨D先").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左脛骨D先").FilledGlobalPosition().MMD()))
		}
	}
	{
		// FIXME
		expectedPosition := &mmath.MVec3{X: 1.763123, Y: 3.708842, Z: 0.369619}
		if !boneDeltas.GetByName("左足Dw").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左足Dw").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左足Dw").FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: 1.698137, Y: 0.674393, Z: 0.043128}
		if !boneDeltas.GetByName("左足先EX").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左足先EX").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左足先EX").FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: 1.712729, Y: 0.919505, Z: 0.174835}
		if !boneDeltas.GetByName("左素足先A").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左素足先A").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左素足先A").FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: 1.494684, Y: 0.328695, Z: -0.500715}
		if !boneDeltas.GetByName("左素足先A先").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左素足先A先").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左素足先A先").FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: 1.494684, Y: 0.328695, Z: -0.500715}
		if !boneDeltas.GetByName("左素足先AIK").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左素足先AIK").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左素足先AIK").FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: 1.224127, Y: 1.738605, Z: 0.240794}
		if !boneDeltas.GetByName("左素足先B").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左素足先B").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左素足先B").FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: 1.582481, Y: 0.627567, Z: -0.222556}
		if !boneDeltas.GetByName("左素足先B先").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左素足先B先").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左素足先B先").FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: 1.572457, Y: 0.658645, Z: -0.209595}
		if !boneDeltas.GetByName("左素足先BIK").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左素足先BIK").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左素足先BIK").FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: 1.448783, Y: 1.760351, Z: 0.405702}
		if !boneDeltas.GetByName("左靴調節").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左靴調節").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左靴調節").FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: 1.875421, Y: 0.917357, Z: 0.652144}
		if !boneDeltas.GetByName("左靴追従").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左靴追従").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左靴追従").FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: 1.404960, Y: 1.846940, Z: 0.380388}
		if !boneDeltas.GetByName("左靴追従先").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左靴追従先").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左靴追従先").FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: 1.448783, Y: 1.760351, Z: 0.405702}
		if !boneDeltas.GetByName("左靴追従IK").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左靴追従IK").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左靴追従IK").FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: 1.330700, Y: 3.021906, Z: 0.153679}
		if !boneDeltas.GetByName("左足補D").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左足補D").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左足補D").FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: -0.625984, Y: 3.444340, Z: -1.272553}
		if !boneDeltas.GetByName("左足補D先").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左足補D先").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左足補D先").FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: -0.625984, Y: 3.444340, Z: -1.272553}
		if !boneDeltas.GetByName("左足補DIK").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左足補DIK").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左足補DIK").FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: 1.724964, Y: 3.674735, Z: 0.810759}
		if !boneDeltas.GetByName("左足向検A").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左足向検A").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左足向検A").FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: -0.924420, Y: 3.943550, Z: -1.420897}
		if !boneDeltas.GetByName("左足向検A先").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左足向検A先").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左足向検A先").FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: -0.924419, Y: 3.943550, Z: -1.420896}
		if !boneDeltas.GetByName("左足向検AIK").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左足向検AIK").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左足向検AIK").FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: 1.724965, Y: 3.674735, Z: 0.810760}
		if !boneDeltas.GetByName("左足向-").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左足向-").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左足向-").FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: 1.566813, Y: 3.645544, Z: 0.177956}
		if !boneDeltas.GetByName("左足w").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左足w").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左足w").FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: -1.879506, Y: 3.670895, Z: -2.715526}
		if !boneDeltas.GetByName("左足w先").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左足w先").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左足w先").FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: -0.631964, Y: 3.647012, Z: -1.211210}
		if !boneDeltas.GetByName("左膝補").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左膝補").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左膝補").FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: -0.83923, Y: 2.876222, Z: -1.494900}
		if !boneDeltas.GetByName("左膝補先").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左膝補先").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左膝補先").FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: -0.686400, Y: 3.444961, Z: -1.285575}
		if !boneDeltas.GetByName("左膝補IK").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左膝補IK").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左膝補IK").FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: -0.924420, Y: 3.943550, Z: -1.420896}
		if !boneDeltas.GetByName("左足捩検B").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左足捩検B").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左足捩検B").FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: -0.952927, Y: 7.388766, Z: -0.972059}
		if !boneDeltas.GetByName("左足捩検B先").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左足捩検B先").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左足捩検B先").FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: -0.952927, Y: 7.388766, Z: -0.972060}
		if !boneDeltas.GetByName("左足捩検BIK").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左足捩検BIK").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左足捩検BIK").FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: 0.400272, Y: 3.809143, Z: -0.305068}
		if !boneDeltas.GetByName("左足捩").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左足捩").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左足捩").FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: 0.371963, Y: 4.256704, Z: -0.267830}
		if !boneDeltas.GetByName("左足捩先").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左足捩先").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左足捩先").FilledGlobalPosition().MMD()))
		}
	}
}

func TestVmdMotion_DeformArmIk_Mahoujin_01(t *testing.T) {
	// mlog.SetLevel(mlog.IK_VERBOSE)

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

func TestVmdMotion_DeformArmIk_Mahoujin_04(t *testing.T) {
	// mlog.SetLevel(mlog.IK_VERBOSE)

	vr := repository.NewVmdRepository()
	motionData, err := vr.Load("../../../test_resources/arm_ik_mahoujin_090F.vmd")

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
		boneName := pmx.ARM.Left()
		expectedPosition := &mmath.MVec3{X: 1.830244, Y: 18.596258, Z: 0.482812}
		if !boneDeltas.GetByName(boneName).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(boneName).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(boneName).FilledGlobalPosition().MMD()))
		}
	}
	{
		boneName := pmx.ELBOW.Left()
		expectedPosition := &mmath.MVec3{X: 2.717007, Y: 18.698180, Z: -2.511497}
		if !boneDeltas.GetByName(boneName).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(boneName).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(boneName).FilledGlobalPosition().MMD()))
		}
	}
	{
		boneName := pmx.WRIST.Left()
		expectedPosition := &mmath.MVec3{X: 0.706904, Y: 21.168780, Z: -3.176916}
		if !boneDeltas.GetByName(boneName).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(boneName).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(boneName).FilledGlobalPosition().MMD()))
		}
	}
	{
		boneName := pmx.INDEX3.Left()
		expectedPosition := &mmath.MVec3{X: 0.120014, Y: 22.707282, Z: -3.770402}
		if !boneDeltas.GetByName(boneName).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(boneName).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(boneName).FilledGlobalPosition().MMD()))
		}
	}
}

func TestVmdMotion_DeformArmIk_Choco_01(t *testing.T) {
	mlog.SetLevel(mlog.IK_VERBOSE)

	vr := repository.NewVmdRepository()
	motionData, err := vr.Load("../../../test_resources/ビタチョコ_0676F.vmd")

	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err)
	}

	motion := motionData.(*vmd.VmdMotion)

	pr := repository.NewPmxRepository()
	modelData, err := pr.Load("D:/MMD/MikuMikuDance_v926x64/UserFile/Model/ゲーム/Fate/眞白式ロマニ・アーキマン ver.1.01/眞白式ロマニ・アーキマン_ビタチョコ2.pmx")

	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err)
	}

	model := modelData.(*pmx.PmxModel)

	boneDeltas := DeformBone(model, motion, true, 0, nil)
	{
		boneName := pmx.ARM.Left()
		expectedPosition := &mmath.MVec3{X: 2.260640, Y: 12.404558, Z: -1.519635}
		if !boneDeltas.GetByName(boneName).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(boneName).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(boneName).FilledGlobalPosition().MMD()))
		}
	}
	{
		boneName := pmx.ELBOW.Left()
		expectedPosition := &mmath.MVec3{X: 1.121608, Y: 11.217656, Z: -4.486015}
		if !boneDeltas.GetByName(boneName).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(boneName).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(boneName).FilledGlobalPosition().MMD()))
		}
	}
	{
		boneName := pmx.WRIST.Left()
		expectedPosition := &mmath.MVec3{X: 0.717674, Y: 13.924381, Z: -3.561227}
		if !boneDeltas.GetByName(boneName).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(boneName).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(boneName).FilledGlobalPosition().MMD()))
		}
	}
	{
		boneName := pmx.INDEX3.Left()
		expectedPosition := &mmath.MVec3{X: 1.002670, Y: 15.652058, Z: -3.506799}
		if !boneDeltas.GetByName(boneName).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(boneName).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(boneName).FilledGlobalPosition().MMD()))
		}
	}
	{
		boneName := pmx.ARM.Right()
		expectedPosition := &mmath.MVec3{X: -2.412614, Y: 12.565295, Z: -1.774290}
		if !boneDeltas.GetByName(boneName).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(boneName).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(boneName).FilledGlobalPosition().MMD()))
		}
	}
	{
		boneName := pmx.ELBOW.Right()
		expectedPosition := &mmath.MVec3{X: -1.009609, Y: 11.296631, Z: -4.589892}
		if !boneDeltas.GetByName(boneName).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(boneName).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(boneName).FilledGlobalPosition().MMD()))
		}
	}
	{
		boneName := pmx.WRIST.Right()
		expectedPosition := &mmath.MVec3{X: -0.137049, Y: 14.029240, Z: -4.235312}
		if !boneDeltas.GetByName(boneName).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(boneName).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(boneName).FilledGlobalPosition().MMD()))
		}
	}
	{
		boneName := pmx.INDEX3.Right()
		expectedPosition := &mmath.MVec3{X: -0.395239, Y: 15.750233, Z: -3.984484}
		if !boneDeltas.GetByName(boneName).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(boneName).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(boneName).FilledGlobalPosition().MMD()))
		}
	}
}
