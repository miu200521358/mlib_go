package vmd

import (
	"testing"

	"github.com/miu200521358/mlib_go/pkg/mmath"
	"github.com/miu200521358/mlib_go/pkg/mutils/mlog"
	"github.com/miu200521358/mlib_go/pkg/pmx"
)

func TestVmdMotion_DeformLegIk25_Addiction_Wa_Left(t *testing.T) {
	mlog.SetLevel(mlog.IK_VERBOSE)

	vr := &VmdMotionReader{}
	motionData, err := vr.ReadByFilepath("../../test_resources/[A]ddiction_和洋_0126F.vmd")

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
		boneDeltas := motion.BoneFrames.Deform(fno, model, []string{"左襟先"}, true, nil, nil)
		{
			expectedPosition := &mmath.MVec3{-0.225006, 9.705784, 2.033072}
			if !boneDeltas.GetByName(pmx.UPPER.String()).GlobalPosition().MMD().PracticallyEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v", expectedPosition, boneDeltas.GetByName(pmx.UPPER.String()).GlobalPosition().MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{-0.237383, 10.769137, 2.039952}
			if !boneDeltas.GetByName(pmx.UPPER2.String()).GlobalPosition().MMD().PracticallyEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v", expectedPosition, boneDeltas.GetByName(pmx.UPPER2.String()).GlobalPosition().MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{0.460140, 13.290816, 2.531440}
			if !boneDeltas.GetByName(pmx.SHOULDER_P.Left()).GlobalPosition().MMD().PracticallyEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v", expectedPosition, boneDeltas.GetByName(pmx.SHOULDER_P.Left()).GlobalPosition().MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{0.460140, 13.290816, 2.531440}
			if !boneDeltas.GetByName(pmx.SHOULDER.Left()).GlobalPosition().MMD().PracticallyEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v", expectedPosition, boneDeltas.GetByName(pmx.SHOULDER.Left()).GlobalPosition().MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{0.784452, 13.728909, 2.608527}
			if !boneDeltas.GetByName(pmx.ARM.Left()).GlobalPosition().MMD().PracticallyEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v", expectedPosition, boneDeltas.GetByName(pmx.ARM.Left()).GlobalPosition().MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{0.272067, 12.381887, 2.182425}
			if !boneDeltas.GetByName("左上半身C-A").GlobalPosition().MMD().PracticallyEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v", expectedPosition, boneDeltas.GetByName("左上半身C-A").GlobalPosition().MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{0.406217, 13.797803, 2.460243}
			if !boneDeltas.GetByName("左鎖骨IK").GlobalPosition().MMD().PracticallyEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v", expectedPosition, boneDeltas.GetByName("左鎖骨IK").GlobalPosition().MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{-0.052427, 13.347448, 2.216718}
			if !boneDeltas.GetByName("左鎖骨").GlobalPosition().MMD().PracticallyEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v", expectedPosition, boneDeltas.GetByName("左鎖骨").GlobalPosition().MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{0.659554, 14.017852, 2.591099}
			if !boneDeltas.GetByName("左鎖骨先").GlobalPosition().MMD().PracticallyEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v", expectedPosition, boneDeltas.GetByName("左鎖骨先").GlobalPosition().MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{0.065147, 13.296564, 2.607907}
			if !boneDeltas.GetByName("左肩Rz検").GlobalPosition().MMD().PracticallyEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v", expectedPosition, boneDeltas.GetByName("左肩Rz検").GlobalPosition().MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{0.517776, 14.134196, 2.645912}
			if !boneDeltas.GetByName("左肩Rz検先").GlobalPosition().MMD().PracticallyEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v", expectedPosition, boneDeltas.GetByName("左肩Rz検先").GlobalPosition().MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{0.065148, 13.296564, 2.607907}
			if !boneDeltas.GetByName("左肩Ry検").GlobalPosition().MMD().PracticallyEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v", expectedPosition, boneDeltas.GetByName("左肩Ry検").GlobalPosition().MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{0.860159, 13.190875, 3.122428}
			if !boneDeltas.GetByName("左肩Ry検先").GlobalPosition().MMD().PracticallyEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v", expectedPosition, boneDeltas.GetByName("左肩Ry検先").GlobalPosition().MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{0.195053, 12.648546, 2.236849}
			if !boneDeltas.GetByName("左上半身C-B").GlobalPosition().MMD().PracticallyEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v", expectedPosition, boneDeltas.GetByName("左上半身C-B").GlobalPosition().MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{0.294257, 12.912640, 2.257159}
			if !boneDeltas.GetByName("左上半身C-C").GlobalPosition().MMD().PracticallyEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v", expectedPosition, boneDeltas.GetByName("左上半身C-C").GlobalPosition().MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{0.210011, 10.897711, 1.973442}
			if !boneDeltas.GetByName("左上半身2").GlobalPosition().MMD().PracticallyEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v", expectedPosition, boneDeltas.GetByName("左上半身2").GlobalPosition().MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{0.320589, 14.049745, 2.637018}
			if !boneDeltas.GetByName("左襟").GlobalPosition().MMD().PracticallyEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v", expectedPosition, boneDeltas.GetByName("左襟").GlobalPosition().MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{0.297636, 14.263302, 2.374467}
			if !boneDeltas.GetByName("左襟先").GlobalPosition().MMD().PracticallyEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v", expectedPosition, boneDeltas.GetByName("左襟先").GlobalPosition().MMD())
			}
		}
	}
}

func TestVmdMotion_DeformLegIk25_Addiction_Wa_Right(t *testing.T) {
	// mlog.SetLevel(mlog.IK_VERBOSE)

	vr := &VmdMotionReader{}
	motionData, err := vr.ReadByFilepath("../../test_resources/[A]ddiction_和洋_0126F.vmd")

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
		boneDeltas := motion.BoneFrames.Deform(fno, model, []string{"右襟先"}, true, nil, nil)
		{
			expectedPosition := &mmath.MVec3{-0.225006, 9.705784, 2.033072}
			if !boneDeltas.GetByName(pmx.UPPER.String()).GlobalPosition().MMD().PracticallyEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v", expectedPosition, boneDeltas.GetByName(pmx.UPPER.String()).GlobalPosition().MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{-0.237383, 10.769137, 2.039952}
			if !boneDeltas.GetByName(pmx.UPPER2.String()).GlobalPosition().MMD().PracticallyEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v", expectedPosition, boneDeltas.GetByName(pmx.UPPER2.String()).GlobalPosition().MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{-0.630130, 13.306682, 2.752505}
			if !boneDeltas.GetByName(pmx.SHOULDER_P.Right()).GlobalPosition().MMD().PracticallyEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v", expectedPosition, boneDeltas.GetByName(pmx.SHOULDER_P.Right()).GlobalPosition().MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{-0.630131, 13.306683, 2.742505}
			if !boneDeltas.GetByName(pmx.SHOULDER.Right()).GlobalPosition().MMD().PracticallyEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v", expectedPosition, boneDeltas.GetByName(pmx.SHOULDER.Right()).GlobalPosition().MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{-0.948004, 13.753115, 2.690539}
			if !boneDeltas.GetByName(pmx.ARM.Right()).GlobalPosition().MMD().PracticallyEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v", expectedPosition, boneDeltas.GetByName(pmx.ARM.Right()).GlobalPosition().MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{-0.611438, 12.394744, 2.353463}
			if !boneDeltas.GetByName("右上半身C-A").GlobalPosition().MMD().PracticallyEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v", expectedPosition, boneDeltas.GetByName("右上半身C-A").GlobalPosition().MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{-0.664344, 13.835273, 2.403165}
			if !boneDeltas.GetByName("右鎖骨IK").GlobalPosition().MMD().PracticallyEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v", expectedPosition, boneDeltas.GetByName("右鎖骨IK").GlobalPosition().MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{-0.270636, 13.350624, 2.258960}
			if !boneDeltas.GetByName("右鎖骨").GlobalPosition().MMD().PracticallyEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v", expectedPosition, boneDeltas.GetByName("右鎖骨").GlobalPosition().MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{-0.963317, 14.098928, 2.497183}
			if !boneDeltas.GetByName("右鎖骨先").GlobalPosition().MMD().PracticallyEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v", expectedPosition, boneDeltas.GetByName("右鎖骨先").GlobalPosition().MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{-0.235138, 13.300934, 2.666039}
			if !boneDeltas.GetByName("右肩Rz検").GlobalPosition().MMD().PracticallyEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v", expectedPosition, boneDeltas.GetByName("右肩Rz検").GlobalPosition().MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{-0.847069, 13.997178, 2.886786}
			if !boneDeltas.GetByName("右肩Rz検先").GlobalPosition().MMD().PracticallyEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v", expectedPosition, boneDeltas.GetByName("右肩Rz検先").GlobalPosition().MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{-0.235138, 13.300934, 2.666039}
			if !boneDeltas.GetByName("右肩Ry検").GlobalPosition().MMD().PracticallyEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v", expectedPosition, boneDeltas.GetByName("右肩Ry検").GlobalPosition().MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{-1.172100, 13.315790, 2.838742}
			if !boneDeltas.GetByName("右肩Ry検先").GlobalPosition().MMD().PracticallyEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v", expectedPosition, boneDeltas.GetByName("右肩Ry検先").GlobalPosition().MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{-0.591152, 12.674325, 2.391185}
			if !boneDeltas.GetByName("右上半身C-B").GlobalPosition().MMD().PracticallyEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v", expectedPosition, boneDeltas.GetByName("右上半身C-B").GlobalPosition().MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{-0.588046, 12.954157, 2.432232}
			if !boneDeltas.GetByName("右上半身C-C").GlobalPosition().MMD().PracticallyEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v", expectedPosition, boneDeltas.GetByName("右上半身C-C").GlobalPosition().MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{-0.672292, 10.939227, 2.148515}
			if !boneDeltas.GetByName("右上半身2").GlobalPosition().MMD().PracticallyEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v", expectedPosition, boneDeltas.GetByName("右上半身2").GlobalPosition().MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{-0.520068, 14.089510, 2.812157}
			if !boneDeltas.GetByName("右襟").GlobalPosition().MMD().PracticallyEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v", expectedPosition, boneDeltas.GetByName("右襟").GlobalPosition().MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{-0.491354, 14.225309, 2.502640}
			if !boneDeltas.GetByName("右襟先").GlobalPosition().MMD().PracticallyEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v", expectedPosition, boneDeltas.GetByName("右襟先").GlobalPosition().MMD())
			}
		}
	}
}

func TestVmdMotion_DeformLegIk25_Ballet(t *testing.T) {
	mlog.SetLevel(mlog.IK_VERBOSE)

	vr := &VmdMotionReader{}
	motionData, err := vr.ReadByFilepath("../../test_resources/青江バレリーコ_1543F.vmd")

	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err)
	}

	motion := motionData.(*VmdMotion)

	pr := &pmx.PmxReader{}
	modelData, err := pr.ReadByFilepath("D:/MMD/MikuMikuDance_v926x64/UserFile/Model/刀剣乱舞/019_にっかり青江/にっかり青江 帽子屋式 ver2.1/帽子屋式にっかり青江（戦装束）.pmx")

	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err)
	}

	model := modelData.(*pmx.PmxModel)

	{

		fno := int(0)
		boneDeltas := motion.BoneFrames.Deform(fno, model, []string{pmx.TOE.Left(), pmx.TOE_EX.Left()}, true, nil, nil)
		{
			expectedPosition := &mmath.MVec3{-4.374956, 13.203792, 1.554190}
			if !boneDeltas.GetByName(pmx.LOWER.String()).GlobalPosition().MMD().PracticallyEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v", expectedPosition, boneDeltas.GetByName(pmx.LOWER.String()).GlobalPosition().MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{-3.481956, 11.214747, 1.127255}
			if !boneDeltas.GetByName(pmx.LEG.Left()).GlobalPosition().MMD().PracticallyEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v", expectedPosition, boneDeltas.GetByName(pmx.LEG.Left()).GlobalPosition().MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{-7.173243, 7.787793, 0.013533}
			if !boneDeltas.GetByName(pmx.KNEE.Left()).GlobalPosition().MMD().PracticallyEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v", expectedPosition, boneDeltas.GetByName(pmx.KNEE.Left()).GlobalPosition().MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{-11.529483, 3.689184, -1.119154}
			if !boneDeltas.GetByName(pmx.ANKLE.Left()).GlobalPosition().MMD().PracticallyEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v", expectedPosition, boneDeltas.GetByName(pmx.ANKLE.Left()).GlobalPosition().MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{-13.408189, 1.877100, -2.183821}
			if !boneDeltas.GetByName(pmx.TOE.Left()).GlobalPosition().MMD().PracticallyEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v", expectedPosition, boneDeltas.GetByName(pmx.TOE.Left()).GlobalPosition().MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{-12.845280, 2.816309, -2.136874}
			if !boneDeltas.GetByName(pmx.TOE_EX.Left()).GlobalPosition().MMD().PracticallyEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v", expectedPosition, boneDeltas.GetByName(pmx.TOE_EX.Left()).GlobalPosition().MMD())
			}
		}
	}
}

func TestVmdMotion_DeformLegIk26_Far(t *testing.T) {
	mlog.SetLevel(mlog.IK_VERBOSE)

	vr := &VmdMotionReader{}
	motionData, err := vr.ReadByFilepath("../../test_resources/足IK乖離.vmd")

	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err)
	}

	motion := motionData.(*VmdMotion)

	pr := &pmx.PmxReader{}
	modelData, err := pr.ReadByFilepath("D:/MMD/MikuMikuDance_v926x64/UserFile/Model/_あにまさ式ミク準標準見せパン/初音ミクVer2 準標準 見せパン 3.pmx")

	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err)
	}

	model := modelData.(*pmx.PmxModel)

	{

		fno := int(0)
		boneDeltas := motion.BoneFrames.Deform(fno, model, []string{pmx.TOE.Right(), pmx.TOE_EX.Right(), pmx.HEEL.Right()}, true, nil, nil)
		{
			expectedPosition := &mmath.MVec3{-0.796811, 10.752734, -0.072743}
			if !boneDeltas.GetByName(pmx.LEG.Right()).GlobalPosition().MMD().PracticallyEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v", expectedPosition, boneDeltas.GetByName(pmx.LEG.Right()).GlobalPosition().MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{-2.202487, 10.921064, -4.695134}
			if !boneDeltas.GetByName(pmx.KNEE.Right()).GlobalPosition().MMD().PracticallyEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v", expectedPosition, boneDeltas.GetByName(pmx.KNEE.Right()).GlobalPosition().MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{-4.193142, 11.026311, -8.844866}
			if !boneDeltas.GetByName(pmx.ANKLE.Right()).GlobalPosition().MMD().PracticallyEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v", expectedPosition, boneDeltas.GetByName(pmx.ANKLE.Right()).GlobalPosition().MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{-5.108798, 10.935530, -11.494570}
			if !boneDeltas.GetByName(pmx.TOE.Right()).GlobalPosition().MMD().PracticallyEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v", expectedPosition, boneDeltas.GetByName(pmx.TOE.Right()).GlobalPosition().MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{-4.800813, 10.964218, -10.612234}
			if !boneDeltas.GetByName(pmx.TOE_EX.Right()).GlobalPosition().MMD().PracticallyEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v", expectedPosition, boneDeltas.GetByName(pmx.TOE_EX.Right()).GlobalPosition().MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{-4.331888, 12.178923, -9.514071}
			if !boneDeltas.GetByName(pmx.HEEL.Right()).GlobalPosition().MMD().PracticallyEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v", expectedPosition, boneDeltas.GetByName(pmx.HEEL.Right()).GlobalPosition().MMD())
			}
		}
	}
}

func TestVmdMotion_DeformLegIk27_Addiction_Shoes(t *testing.T) {
	mlog.SetLevel(mlog.IK_VERBOSE)

	vr := &VmdMotionReader{}
	motionData, err := vr.ReadByFilepath("../../test_resources/[A]ddiction_和洋_1074-1078F.vmd")

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
		fno := int(2)
		boneDeltas := motion.BoneFrames.Deform(fno, model, nil, true, nil, nil)
		{
			expectedPosition := &mmath.MVec3{0, 0, 0}
			if !boneDeltas.GetByName(pmx.ROOT.String()).GlobalPosition().MMD().PracticallyEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v", expectedPosition, boneDeltas.GetByName(pmx.ROOT.String()).GlobalPosition().MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{1.406722, 1.841236, 0.277818}
			if !boneDeltas.GetByName(pmx.LEG_IK.Left()).GlobalPosition().MMD().PracticallyEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v", expectedPosition, boneDeltas.GetByName(pmx.LEG_IK.Right()).GlobalPosition().MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{0.510231, 9.009953, 0.592482}
			if !boneDeltas.GetByName(pmx.LOWER.String()).GlobalPosition().MMD().PracticallyEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v", expectedPosition, boneDeltas.GetByName(pmx.LOWER.String()).GlobalPosition().MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{1.355914, 7.853320, 0.415251}
			if !boneDeltas.GetByName(pmx.LEG.Left()).GlobalPosition().MMD().PracticallyEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v", expectedPosition, boneDeltas.GetByName(pmx.LEG.Right()).GlobalPosition().MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{-0.327781, 5.203806, -1.073718}
			if !boneDeltas.GetByName(pmx.KNEE.Left()).GlobalPosition().MMD().PracticallyEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v", expectedPosition, boneDeltas.GetByName(pmx.KNEE.Right()).GlobalPosition().MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{1.407848, 1.839228, 0.278700}
			if !boneDeltas.GetByName("左脛骨").GlobalPosition().MMD().PracticallyEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v", expectedPosition, boneDeltas.GetByName("左脛骨").GlobalPosition().MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{1.407848, 1.839228, 0.278700}
			if !boneDeltas.GetByName("左脛骨D").GlobalPosition().MMD().PracticallyEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v", expectedPosition, boneDeltas.GetByName("左脛骨D").GlobalPosition().MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{-0.498054, 5.045506, -1.221016}
			if !boneDeltas.GetByName("左脛骨D先").GlobalPosition().MMD().PracticallyEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v", expectedPosition, boneDeltas.GetByName("左脛骨D先").GlobalPosition().MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{1.462306, 7.684025, 0.087026}
			if !boneDeltas.GetByName("左足Dw").GlobalPosition().MMD().PracticallyEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v", expectedPosition, boneDeltas.GetByName("左足Dw").GlobalPosition().MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{1.593721, 0.784840, -0.054141}
			if !boneDeltas.GetByName("左足先EX").GlobalPosition().MMD().PracticallyEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v", expectedPosition, boneDeltas.GetByName("左足先EX").GlobalPosition().MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{1.551940, 1.045847, 0.034003}
			if !boneDeltas.GetByName("左素足先A").GlobalPosition().MMD().PracticallyEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v", expectedPosition, boneDeltas.GetByName("左素足先A").GlobalPosition().MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{1.453982, 0.305976, -0.510022}
			if !boneDeltas.GetByName("左素足先A先").GlobalPosition().MMD().PracticallyEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v", expectedPosition, boneDeltas.GetByName("左素足先A先").GlobalPosition().MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{1.453982, 0.305976, -0.510022}
			if !boneDeltas.GetByName("左素足先AIK").GlobalPosition().MMD().PracticallyEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v", expectedPosition, boneDeltas.GetByName("左素足先AIK").GlobalPosition().MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{0.941880, 2.132958, 0.020403}
			if !boneDeltas.GetByName("左素足先B").GlobalPosition().MMD().PracticallyEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v", expectedPosition, boneDeltas.GetByName("左素足先B").GlobalPosition().MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{1.359364, 0.974298, -0.226041}
			if !boneDeltas.GetByName("左素足先B先").GlobalPosition().MMD().PracticallyEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v", expectedPosition, boneDeltas.GetByName("左素足先B先").GlobalPosition().MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{1.460890, 0.692527, -0.285973}
			if !boneDeltas.GetByName("左素足先BIK").GlobalPosition().MMD().PracticallyEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v", expectedPosition, boneDeltas.GetByName("左素足先BIK").GlobalPosition().MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{1.173929, 2.066327, 0.182685}
			if !boneDeltas.GetByName("左靴調節").GlobalPosition().MMD().PracticallyEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v", expectedPosition, boneDeltas.GetByName("左靴調節").GlobalPosition().MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{1.739235, 1.171441, 0.485052}
			if !boneDeltas.GetByName("左靴追従").GlobalPosition().MMD().PracticallyEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v", expectedPosition, boneDeltas.GetByName("左靴追従").GlobalPosition().MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{1.186359, 2.046771, 0.189367}
			if !boneDeltas.GetByName("左靴追従先").GlobalPosition().MMD().PracticallyEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v", expectedPosition, boneDeltas.GetByName("左靴追従先").GlobalPosition().MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{1.173929, 2.066327, 0.182685}
			if !boneDeltas.GetByName("左靴追従IK").GlobalPosition().MMD().PracticallyEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v", expectedPosition, boneDeltas.GetByName("左靴追従IK").GlobalPosition().MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{1.574899, 6.873434, 0.342768}
			if !boneDeltas.GetByName("左足補D").GlobalPosition().MMD().PracticallyEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v", expectedPosition, boneDeltas.GetByName("左足補D").GlobalPosition().MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{0.150401, 5.170907, -0.712416}
			if !boneDeltas.GetByName("左足補D先").GlobalPosition().MMD().PracticallyEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v", expectedPosition, boneDeltas.GetByName("左足補D先").GlobalPosition().MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{0.150401, 5.170907, -0.712416}
			if !boneDeltas.GetByName("左足補DIK").GlobalPosition().MMD().PracticallyEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v", expectedPosition, boneDeltas.GetByName("左足補DIK").GlobalPosition().MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{1.355915, 7.853319, 0.415251}
			if !boneDeltas.GetByName("左足向検A").GlobalPosition().MMD().PracticallyEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v", expectedPosition, boneDeltas.GetByName("左足向検A").GlobalPosition().MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{-0.327781, 5.203805, -1.073719}
			if !boneDeltas.GetByName("左足向検A先").GlobalPosition().MMD().PracticallyEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v", expectedPosition, boneDeltas.GetByName("左足向検A先").GlobalPosition().MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{-0.327781, 5.203805, -1.073719}
			if !boneDeltas.GetByName("左足向検AIK").GlobalPosition().MMD().PracticallyEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v", expectedPosition, boneDeltas.GetByName("左足向検AIK").GlobalPosition().MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{1.355914, 7.853319, 0.415251}
			if !boneDeltas.GetByName("左足向-").GlobalPosition().MMD().PracticallyEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v", expectedPosition, boneDeltas.GetByName("左足向-").GlobalPosition().MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{1.264808, 7.561551, -0.161703}
			if !boneDeltas.GetByName("左足w").GlobalPosition().MMD().PracticallyEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v", expectedPosition, boneDeltas.GetByName("左足w").GlobalPosition().MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{-0.714029, 3.930234, -1.935889}
			if !boneDeltas.GetByName("左足w先").GlobalPosition().MMD().PracticallyEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v", expectedPosition, boneDeltas.GetByName("左足w先").GlobalPosition().MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{0.016770, 5.319929, -0.781771}
			if !boneDeltas.GetByName("左膝補").GlobalPosition().MMD().PracticallyEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v", expectedPosition, boneDeltas.GetByName("左膝補").GlobalPosition().MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{-0.164672, 4.511360, -0.957886}
			if !boneDeltas.GetByName("左膝補先").GlobalPosition().MMD().PracticallyEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v", expectedPosition, boneDeltas.GetByName("左膝補先").GlobalPosition().MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{-0.099887, 4.800064, -0.895003}
			if !boneDeltas.GetByName("左膝補IK").GlobalPosition().MMD().PracticallyEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v", expectedPosition, boneDeltas.GetByName("左膝補IK").GlobalPosition().MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{-0.327781, 5.203806, -1.073718}
			if !boneDeltas.GetByName("左足捩検B").GlobalPosition().MMD().PracticallyEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v", expectedPosition, boneDeltas.GetByName("左足捩検B").GlobalPosition().MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{-2.392915, 7.450026, -2.735495}
			if !boneDeltas.GetByName("左足捩検B先").GlobalPosition().MMD().PracticallyEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v", expectedPosition, boneDeltas.GetByName("左足捩検B先").GlobalPosition().MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{-2.392915, 7.450026, -2.735495}
			if !boneDeltas.GetByName("左足捩検BIK").GlobalPosition().MMD().PracticallyEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v", expectedPosition, boneDeltas.GetByName("左足捩検BIK").GlobalPosition().MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{0.514067, 6.528563, -0.329234}
			if !boneDeltas.GetByName("左足捩").GlobalPosition().MMD().PracticallyEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v", expectedPosition, boneDeltas.GetByName("左足捩").GlobalPosition().MMD())
			}
		}
		{
			expectedPosition := &mmath.MVec3{0.231636, 6.794109, -0.557747}
			if !boneDeltas.GetByName("左足捩先").GlobalPosition().MMD().PracticallyEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v", expectedPosition, boneDeltas.GetByName("左足捩先").GlobalPosition().MMD())
			}
		}
	}
}
