// 指示: miu200521358
package pmx

import (
	"math"
	"path/filepath"
	"testing"

	"github.com/miu200521358/mlib_go/pkg/adapter/io_common"
	"github.com/miu200521358/mlib_go/pkg/domain/mmath"
	"github.com/miu200521358/mlib_go/pkg/domain/model"
)

func TestPmxRepository_Save(t *testing.T) {
	r := NewPmxRepository()

	data, err := r.Load(testResourcePath("サンプルモデル_PMX読み取り確認用.pmx"))
	if err != nil {
		t.Fatalf("Expected error to be nil, got %q", err)
	}
	modelData, ok := data.(*model.PmxModel)
	if !ok {
		t.Fatalf("Expected model type to be *PmxModel, got %T", data)
	}

	savePath := filepath.Join(t.TempDir(), "サンプルモデル_PMX読み取り確認用_output.pmx")
	if err := r.Save(savePath, modelData, io_common.SaveOptions{IncludeSystem: false}); err != nil {
		t.Fatalf("Expected error to be nil, got %q", err)
	}

	savedData, err := r.Load(savePath)
	if err != nil {
		t.Fatalf("Expected error to be nil, got %q", err)
	}
	savedModel, ok := savedData.(*model.PmxModel)
	if !ok {
		t.Fatalf("Expected model type to be *PmxModel, got %T", savedData)
	}

	expectedName := "v2配布用素体03"
	if savedModel.Name() != expectedName {
		t.Errorf("Expected Name to be %q, got %q", expectedName, savedModel.Name())
	}

	{
		v, _ := savedModel.Vertices.Get(13)
		expectedPosition := vec3(0.1565633, 16.62944, -0.2150156)
		if !v.Position.MMD().NearEquals(expectedPosition, 1e-5) {
			t.Errorf("Expected Position to be %v, got %v", expectedPosition, v.Position)
		}
		expectedNormal := vec3(0.2274586, 0.6613649, -0.714744)
		if !v.Normal.MMD().NearEquals(expectedNormal, 1e-5) {
			t.Errorf("Expected Normal to be %v, got %v", expectedNormal, v.Normal)
		}
		expectedUV := mmath.Vec2{X: 0.5112334, Y: 0.1250942}
		if !v.Uv.NearEquals(expectedUV, 1e-5) {
			t.Errorf("Expected UV to be %v, got %v", expectedUV, v.Uv)
		}
		expectedDeformType := model.BDEF4
		if v.DeformType != expectedDeformType {
			t.Errorf("Expected DeformType to be %d, got %d", expectedDeformType, v.DeformType)
		}
		expectedDeform := model.NewBdef4(
			[4]int{7, 8, 25, 9},
			[4]float64{0.6375693, 0.2368899, 0.06898639, 0.05655446},
		)
		if v.Deform.Indexes()[0] != expectedDeform.Indexes()[0] {
			t.Errorf("Expected Deform to be %v, got %v", expectedDeform.Indexes()[0], v.Deform.Indexes()[0])
		}
		if v.Deform.Indexes()[1] != expectedDeform.Indexes()[1] {
			t.Errorf("Expected Deform to be %v, got %v", expectedDeform.Indexes()[1], v.Deform.Indexes()[1])
		}
		if v.Deform.Indexes()[2] != expectedDeform.Indexes()[2] {
			t.Errorf("Expected Deform to be %v, got %v", expectedDeform.Indexes()[2], v.Deform.Indexes()[2])
		}
		if v.Deform.Indexes()[3] != expectedDeform.Indexes()[3] {
			t.Errorf("Expected Deform to be %v, got %v", expectedDeform.Indexes()[3], v.Deform.Indexes()[3])
		}
		if math.Abs(v.Deform.Weights()[0]-expectedDeform.Weights()[0]) > 1e-5 {
			t.Errorf("Expected Deform to be %v, got %v", expectedDeform.Weights()[0], v.Deform.Weights()[0])
		}
		if math.Abs(v.Deform.Weights()[1]-expectedDeform.Weights()[1]) > 1e-5 {
			t.Errorf("Expected Deform to be %v, got %v", expectedDeform.Weights()[1], v.Deform.Weights()[1])
		}
		if math.Abs(v.Deform.Weights()[2]-expectedDeform.Weights()[2]) > 1e-5 {
			t.Errorf("Expected Deform to be %v, got %v", expectedDeform.Weights()[2], v.Deform.Weights()[2])
		}
		if math.Abs(v.Deform.Weights()[3]-expectedDeform.Weights()[3]) > 1e-5 {
			t.Errorf("Expected Deform to be %v, got %v", expectedDeform.Weights()[3], v.Deform.Weights()[3])
		}
		expectedEdgeFactor := 1.0
		if math.Abs(v.EdgeFactor-expectedEdgeFactor) > 1e-5 {
			t.Errorf("Expected EdgeFactor to be %v, got %v", expectedEdgeFactor, v.EdgeFactor)
		}
	}

	{
		v, _ := savedModel.Vertices.Get(120)
		expectedPosition := vec3(1.529492, 5.757646, 0.4527041)
		if !v.Position.MMD().NearEquals(expectedPosition, 1e-5) {
			t.Errorf("Expected Position to be %v, got %v", expectedPosition, v.Position)
		}
		expectedNormal := vec3(0.9943396, 0.1054612, -0.0129031)
		if !v.Normal.MMD().NearEquals(expectedNormal, 1e-5) {
			t.Errorf("Expected Normal to be %v, got %v", expectedNormal, v.Normal)
		}
		expectedUV := mmath.Vec2{X: 0.8788766, Y: 0.7697825}
		if !v.Uv.NearEquals(expectedUV, 1e-5) {
			t.Errorf("Expected UV to be %v, got %v", expectedUV, v.Uv)
		}
		expectedDeformType := model.BDEF2
		if v.DeformType != expectedDeformType {
			t.Errorf("Expected DeformType to be %d, got %d", expectedDeformType, v.DeformType)
		}
		expectedDeform := model.NewBdef2(109, 108, 0.9845969)
		if v.Deform.Indexes()[0] != expectedDeform.Indexes()[0] {
			t.Errorf("Expected Deform to be %v, got %v", expectedDeform.Indexes()[0], v.Deform.Indexes()[0])
		}
		if v.Deform.Indexes()[1] != expectedDeform.Indexes()[1] {
			t.Errorf("Expected Deform to be %v, got %v", expectedDeform.Indexes()[1], v.Deform.Indexes()[1])
		}
		if math.Abs(v.Deform.Weights()[0]-expectedDeform.Weights()[0]) > 1e-5 {
			t.Errorf("Expected Deform to be %v, got %v", expectedDeform.Weights()[0], v.Deform.Weights()[0])
		}
		if math.Abs(v.Deform.Weights()[1]-expectedDeform.Weights()[1]) > 1e-5 {
			t.Errorf("Expected Deform to be %v, got %v", expectedDeform.Weights()[1], v.Deform.Weights()[1])
		}
		expectedEdgeFactor := 1.0
		if math.Abs(v.EdgeFactor-expectedEdgeFactor) > 1e-5 {
			t.Errorf("Expected EdgeFactor to be %v, got %v", expectedEdgeFactor, v.EdgeFactor)
		}
	}

	{
		b, _ := savedModel.Bones.Get(12)
		expectedBoneFlag := model.BONE_FLAG_CAN_ROTATE | model.BONE_FLAG_IS_VISIBLE | model.BONE_FLAG_CAN_MANIPULATE | model.BONE_FLAG_TAIL_IS_BONE | model.BONE_FLAG_IS_EXTERNAL_ROTATION
		if b.BoneFlag != expectedBoneFlag {
			t.Errorf("Expected BoneFlag to be %v, got %v", expectedBoneFlag, b.BoneFlag)
		}
		expectedTailIndex := -1
		if b.TailIndex != expectedTailIndex {
			t.Errorf("Expected TailIndex to be %v, got %v", expectedTailIndex, b.TailIndex)
		}
	}

	{
		b, _ := savedModel.Bones.Get(28)
		expectedBoneFlag := model.BONE_FLAG_CAN_ROTATE | model.BONE_FLAG_IS_VISIBLE | model.BONE_FLAG_CAN_MANIPULATE | model.BONE_FLAG_TAIL_IS_BONE | model.BONE_FLAG_HAS_FIXED_AXIS | model.BONE_FLAG_HAS_LOCAL_AXIS
		if b.BoneFlag != expectedBoneFlag {
			t.Errorf("Expected BoneFlag to be %v, got %v", expectedBoneFlag, b.BoneFlag)
		}
		expectedTailIndex := -1
		if b.TailIndex != expectedTailIndex {
			t.Errorf("Expected TailIndex to be %v, got %v", expectedTailIndex, b.TailIndex)
		}
	}
}
