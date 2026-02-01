//go:build windows
// +build windows

// 指示: miu200521358
package deform

import (
	"testing"

	"github.com/miu200521358/mlib_go/pkg/domain/model"
	"github.com/miu200521358/mlib_go/pkg/domain/motion"
)

func NoTestVmdMotion_DeformLegIk30_Addiction_Shoes(t *testing.T) {
	// mlog.SetLevel(mlog.IK_VERBOSE)

	vmdMotion := loadVmd(t, "../../../internal/test_resources/[A]ddiction_和洋_1037F.vmd")

	pmxModel := loadPmx(t, "D:/MMD/MikuMikuDance_v926x64/UserFile/Model/_VMDサイジング/wa_129cm 20231028/wa_129cm_20240406.pmx")
	boneDeltas, _ := computeBoneDeltas(pmxModel, vmdMotion, motion.Frame(0), nil, true, false, false)
	{
		expectedPosition := vec3(0, 0, 0)
		if !boneDeltas.GetByName(model.ROOT.String()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(model.ROOT.String()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(model.ROOT.String()).FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(1.748025, 1.683590, 0.556993)
		if !boneDeltas.GetByName(model.LEG_IK.Left()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(model.LEG_IK.Right()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(model.LEG_IK.Right()).FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(1.111190, 4.955496, 1.070225)
		if !boneDeltas.GetByName(model.LOWER.String()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(model.LOWER.String()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(model.LOWER.String()).FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(1.724965, 3.674735, 0.810759)
		if !boneDeltas.GetByName(model.LEG.Left()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(model.LEG.Right()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(model.LEG.Right()).FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(-0.924419, 3.943549, -1.420897)
		if !boneDeltas.GetByName(model.KNEE.Left()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(model.KNEE.Right()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(model.KNEE.Right()).FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(1.743128, 1.672784, 0.551317)
		if !boneDeltas.GetByName("左脛骨").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左脛骨").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("左脛骨").FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(1.743128, 1.672784, 0.551317)
		if !boneDeltas.GetByName("左脛骨D").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左脛骨D").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("左脛骨D").FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(-1.118480, 3.432016, -1.657329)
		if !boneDeltas.GetByName("左脛骨D先").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左脛骨D先").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("左脛骨D先").FilledGlobalPosition()))
		}
	}
	{
		// TODO(miu200521358, 2026-01-18): 期待値の再確認
		expectedPosition := vec3(1.763123, 3.708842, 0.369619)
		if !boneDeltas.GetByName("左足Dw").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左足Dw").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("左足Dw").FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(1.698137, 0.674393, 0.043128)
		if !boneDeltas.GetByName("左足先EX").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左足先EX").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("左足先EX").FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(1.712729, 0.919505, 0.174835)
		if !boneDeltas.GetByName("左素足先A").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左素足先A").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("左素足先A").FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(1.494684, 0.328695, -0.500715)
		if !boneDeltas.GetByName("左素足先A先").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左素足先A先").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("左素足先A先").FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(1.494684, 0.328695, -0.500715)
		if !boneDeltas.GetByName("左素足先AIK").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左素足先AIK").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("左素足先AIK").FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(1.224127, 1.738605, 0.240794)
		if !boneDeltas.GetByName("左素足先B").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左素足先B").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("左素足先B").FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(1.582481, 0.627567, -0.222556)
		if !boneDeltas.GetByName("左素足先B先").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左素足先B先").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("左素足先B先").FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(1.572457, 0.658645, -0.209595)
		if !boneDeltas.GetByName("左素足先BIK").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左素足先BIK").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("左素足先BIK").FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(1.448783, 1.760351, 0.405702)
		if !boneDeltas.GetByName("左靴調節").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左靴調節").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("左靴調節").FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(1.875421, 0.917357, 0.652144)
		if !boneDeltas.GetByName("左靴追従").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左靴追従").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("左靴追従").FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(1.404960, 1.846940, 0.380388)
		if !boneDeltas.GetByName("左靴追従先").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左靴追従先").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("左靴追従先").FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(1.448783, 1.760351, 0.405702)
		if !boneDeltas.GetByName("左靴追従IK").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左靴追従IK").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("左靴追従IK").FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(1.330700, 3.021906, 0.153679)
		if !boneDeltas.GetByName("左足補D").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左足補D").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("左足補D").FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(-0.625984, 3.444340, -1.272553)
		if !boneDeltas.GetByName("左足補D先").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左足補D先").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("左足補D先").FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(-0.625984, 3.444340, -1.272553)
		if !boneDeltas.GetByName("左足補DIK").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左足補DIK").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("左足補DIK").FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(1.724964, 3.674735, 0.810759)
		if !boneDeltas.GetByName("左足向検A").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左足向検A").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("左足向検A").FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(-0.924420, 3.943550, -1.420897)
		if !boneDeltas.GetByName("左足向検A先").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左足向検A先").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("左足向検A先").FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(-0.924419, 3.943550, -1.420896)
		if !boneDeltas.GetByName("左足向検AIK").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左足向検AIK").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("左足向検AIK").FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(1.724965, 3.674735, 0.810760)
		if !boneDeltas.GetByName("左足向-").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左足向-").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("左足向-").FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(1.566813, 3.645544, 0.177956)
		if !boneDeltas.GetByName("左足w").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左足w").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("左足w").FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(-1.879506, 3.670895, -2.715526)
		if !boneDeltas.GetByName("左足w先").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左足w先").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("左足w先").FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(-0.631964, 3.647012, -1.211210)
		if !boneDeltas.GetByName("左膝補").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左膝補").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("左膝補").FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(-0.83923, 2.876222, -1.494900)
		if !boneDeltas.GetByName("左膝補先").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左膝補先").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("左膝補先").FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(-0.686400, 3.444961, -1.285575)
		if !boneDeltas.GetByName("左膝補IK").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左膝補IK").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("左膝補IK").FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(-0.924420, 3.943550, -1.420896)
		if !boneDeltas.GetByName("左足捩検B").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左足捩検B").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("左足捩検B").FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(-0.952927, 7.388766, -0.972059)
		if !boneDeltas.GetByName("左足捩検B先").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左足捩検B先").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("左足捩検B先").FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(-0.952927, 7.388766, -0.972060)
		if !boneDeltas.GetByName("左足捩検BIK").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左足捩検BIK").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("左足捩検BIK").FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(0.400272, 3.809143, -0.305068)
		if !boneDeltas.GetByName("左足捩").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左足捩").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("左足捩").FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(0.371963, 4.256704, -0.267830)
		if !boneDeltas.GetByName("左足捩先").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左足捩先").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("左足捩先").FilledGlobalPosition()))
		}
	}
}

func NoTestVmdMotion_DeformArmIk_Mahoujin_01(t *testing.T) {
	// mlog.SetLevel(mlog.IK_VERBOSE)

	vmdMotion := loadVmd(t, "../../../internal/test_resources/arm_ik_mahoujin_001F.vmd")

	pmxModel := loadPmx(t, "D:/MMD/MikuMikuDance_v926x64/UserFile/Model/刀剣乱舞/107_髭切/髭切mkmk009c 刀剣乱舞/髭切mkmk009c/髭切上着無mkmk009b_腕ＩＫ2.pmx")
	boneDeltas, _ := computeBoneDeltas(pmxModel, vmdMotion, motion.Frame(0), nil, true, false, false)
	{
		boneName := model.ARM.Right()
		expectedPosition := vec3(-1.801768, 18.555544, 0.482812)
		if !boneDeltas.GetByName(boneName).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(boneName).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(boneName).FilledGlobalPosition()))
		}
	}
	{
		boneName := model.ELBOW.Right()
		expectedPosition := vec3(-4.091116, 18.629446, -1.670793)
		if !boneDeltas.GetByName(boneName).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(boneName).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(boneName).FilledGlobalPosition()))
		}
	}
	{
		boneName := model.WRIST.Right()
		expectedPosition := vec3(-6.370411, 18.910606, -4.062796)
		if !boneDeltas.GetByName(boneName).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(boneName).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(boneName).FilledGlobalPosition()))
		}
	}
	{
		boneName := model.INDEX3.Right()
		expectedPosition := vec3(-7.256862, 18.269156, -5.428672)
		if !boneDeltas.GetByName(boneName).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(boneName).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(boneName).FilledGlobalPosition()))
		}
	}
}

func NoTestVmdMotion_DeformArmIk_Mahoujin_04(t *testing.T) {
	// mlog.SetLevel(mlog.IK_VERBOSE)

	vmdMotion := loadVmd(t, "../../../internal/test_resources/arm_ik_mahoujin_090F.vmd")

	pmxModel := loadPmx(t, "D:/MMD/MikuMikuDance_v926x64/UserFile/Model/刀剣乱舞/107_髭切/髭切mkmk009c 刀剣乱舞/髭切mkmk009c/髭切上着無mkmk009b_腕ＩＫ2.pmx")
	boneDeltas, _ := computeBoneDeltas(pmxModel, vmdMotion, motion.Frame(0), nil, true, false, false)
	{
		boneName := model.ARM.Left()
		expectedPosition := vec3(1.830244, 18.596258, 0.482812)
		if !boneDeltas.GetByName(boneName).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(boneName).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(boneName).FilledGlobalPosition()))
		}
	}
	{
		boneName := model.ELBOW.Left()
		expectedPosition := vec3(2.717007, 18.698180, -2.511497)
		if !boneDeltas.GetByName(boneName).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(boneName).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(boneName).FilledGlobalPosition()))
		}
	}
	{
		boneName := model.WRIST.Left()
		expectedPosition := vec3(0.706904, 21.168780, -3.176916)
		if !boneDeltas.GetByName(boneName).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(boneName).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(boneName).FilledGlobalPosition()))
		}
	}
	{
		boneName := model.INDEX3.Left()
		expectedPosition := vec3(0.120014, 22.707282, -3.770402)
		if !boneDeltas.GetByName(boneName).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(boneName).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(boneName).FilledGlobalPosition()))
		}
	}
}

func NoTestVmdMotion_DeformLegIk_Up(t *testing.T) {
	// mlog.SetLevel(mlog.IK_VERBOSE)

	vmdMotion := loadVmd(t, "../../../internal/test_resources/左足あげ.vmd")

	pmxModel := loadPmx(t, "D:/MMD/MikuMikuDance_v926x64/UserFile/Costume/モノクロストリート風衣装 夜/ストリート風白_3.pmx")
	// pmxModel := loadPmx(t, "D:/MMD/MikuMikuDance_v926x64/UserFile/Costume/モノクロストリート風衣装 夜/ストリート風白.pmx")
	// pmxModel := loadPmx(t, "D:/MMD/MikuMikuDance_v926x64/UserFile/Model/VOCALOID/紲星☆あかり20180430 お宮/お宮式紲星☆あかりv1.00.pmx")
	_, _ = computeBoneDeltas(pmxModel, vmdMotion, motion.Frame(0), nil, true, false, false)
}
