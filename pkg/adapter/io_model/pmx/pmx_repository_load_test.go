// 指示: miu200521358
package pmx

import (
	"math"
	"path/filepath"
	"testing"

	"github.com/miu200521358/mlib_go/pkg/domain/mmath"
	"github.com/miu200521358/mlib_go/pkg/domain/model"
	"gonum.org/v1/gonum/spatial/r3"
)

func TestPmxRepository_Load(t *testing.T) {
	r := NewPmxRepository()

	data, err := r.Load(testResourcePath("サンプルモデル_PMX読み取り確認用.pmx"))
	if err != nil {
		t.Fatalf("Expected error to be nil, got %q", err)
	}
	modelData, ok := data.(*model.PmxModel)
	if !ok {
		t.Fatalf("Expected model type to be *PmxModel, got %T", data)
	}

	assertSampleModel(t, modelData)
}

func TestPmxRepository_Load_NotExist(t *testing.T) {
	r := NewPmxRepository()

	_, err := r.Load(testResourcePath("サンプルモデル_Nothing.pmx"))
	if err == nil {
		t.Errorf("Expected error to be not nil, got nil")
	}
}

func TestPmxRepository_Load_2_1(t *testing.T) {
	r := NewPmxRepository()

	data, err := r.Load(testResourcePath("サンプルモデル_PMX2.1_UTF-8.pmx"))
	if err != nil {
		t.Fatalf("Expected error to be nil, got %q", err)
	}
	modelData, ok := data.(*model.PmxModel)
	if !ok {
		t.Fatalf("Expected model type to be *PmxModel, got %T", data)
	}

	expectedName := "サンプルモデル迪卢克"
	if modelData.Name() != expectedName {
		t.Errorf("Expected Name to be %q, got %q", expectedName, modelData.Name())
	}
}

// assertSampleModel はサンプルモデルの読み込み結果を検証する。
func assertSampleModel(t *testing.T, modelData *model.PmxModel) {
	t.Helper()
	if modelData == nil {
		t.Fatalf("Expected modelData to be not nil")
	}
	pmxModel := modelData

	expectedName := "v2配布用素体03"
	if pmxModel.Name() != expectedName {
		t.Errorf("Expected Name to be %q, got %q", expectedName, pmxModel.Name())
	}

	{
		v, _ := pmxModel.Vertices.Get(13)
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
		v, _ := pmxModel.Vertices.Get(120)
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
		f, _ := pmxModel.Faces.Get(19117)
		expectedFaceVertexIndexes := []int{8857, 8893, 8871}
		if f.VertexIndexes[0] != expectedFaceVertexIndexes[0] {
			t.Errorf("Expected Deform to be %v, got %v", expectedFaceVertexIndexes[0], f.VertexIndexes[0])
		}
		if f.VertexIndexes[1] != expectedFaceVertexIndexes[1] {
			t.Errorf("Expected Deform to be %v, got %v", expectedFaceVertexIndexes[1], f.VertexIndexes[1])
		}
		if f.VertexIndexes[2] != expectedFaceVertexIndexes[2] {
			t.Errorf("Expected Deform to be %v, got %v", expectedFaceVertexIndexes[2], f.VertexIndexes[2])
		}
	}

	{
		tex, _ := pmxModel.Textures.Get(10)
		expectedName := "tex\\_13_Toon.bmp"
		if tex.Name() != expectedName {
			t.Errorf("Expected Path to be %q, got %q", expectedName, tex.Name())
		}
	}

	{
		m, _ := pmxModel.Materials.Get(10)
		expectedName := "00_EyeWhite_はぅ"
		if m.Name() != expectedName {
			t.Errorf("Expected Path to be %q, got %q", expectedName, m.Name())
		}
		expectedEnglishName := "N00_000_00_EyeWhite_00_EYE (Instance)_Hau"
		if m.EnglishName != expectedEnglishName {
			t.Errorf("Expected Path to be %q, got %q", expectedEnglishName, m.EnglishName)
		}
		expectedDiffuse := mmath.Vec4{X: 1.0, Y: 1.0, Z: 1.0, W: 0.0}
		if !m.Diffuse.NearEquals(expectedDiffuse, 1e-5) {
			t.Errorf("Expected Diffuse to be %v, got %v", expectedDiffuse, m.Diffuse)
		}
		expectedSpecular := mmath.Vec4{X: 0.0, Y: 0.0, Z: 0.0, W: 0.0}
		if !m.Specular.NearEquals(expectedSpecular, 1e-5) {
			t.Errorf("Expected Specular to be %v, got %v", expectedSpecular, m.Specular)
		}
		expectedAmbient := vec3(0.5, 0.5, 0.5)
		if !m.Ambient.NearEquals(expectedAmbient, 1e-5) {
			t.Errorf("Expected Ambient to be %v, got %v", expectedAmbient, m.Ambient)
		}
		expectedDrawFlag := model.DRAW_FLAG_GROUND_SHADOW | model.DRAW_FLAG_DRAWING_ON_SELF_SHADOW_MAPS | model.DRAW_FLAG_DRAWING_SELF_SHADOWS
		if m.DrawFlag != expectedDrawFlag {
			t.Errorf("Expected DrawFlag to be %v, got %v", expectedDrawFlag, m.DrawFlag)
		}
		expectedEdge := mmath.Vec4{X: 0.2745098, Y: 0.09019607, Z: 0.1254902, W: 1.0}
		if !m.Edge.NearEquals(expectedEdge, 1e-5) {
			t.Errorf("Expected Edge to be %v, got %v", expectedEdge, m.Edge)
		}
		expectedEdgeSize := 1.0
		if math.Abs(m.EdgeSize-expectedEdgeSize) > 1e-5 {
			t.Errorf("Expected EdgeSize to be %v, got %v", expectedEdgeSize, m.EdgeSize)
		}
		expectedTextureIndex := 22
		if m.TextureIndex != expectedTextureIndex {
			t.Errorf("Expected TextureIndex to be %v, got %v", expectedTextureIndex, m.TextureIndex)
		}
		expectedSphereTextureIndex := 4
		if m.SphereTextureIndex != expectedSphereTextureIndex {
			t.Errorf("Expected SphereTextureIndex to be %v, got %v", expectedSphereTextureIndex, m.SphereTextureIndex)
		}
		expectedSphereMode := model.SPHERE_MODE_ADDITION
		if m.SphereMode != expectedSphereMode {
			t.Errorf("Expected SphereMode to be %v, got %v", expectedSphereMode, m.SphereMode)
		}
		expectedToonSharingFlag := model.TOON_SHARING_INDIVIDUAL
		if m.ToonSharingFlag != expectedToonSharingFlag {
			t.Errorf("Expected ToonSharingFlag to be %v, got %v", expectedToonSharingFlag, m.ToonSharingFlag)
		}
		expectedToonTextureIndex := 21
		if m.ToonTextureIndex != expectedToonTextureIndex {
			t.Errorf("Expected ToonTextureIndex to be %v, got %v", expectedToonTextureIndex, m.ToonTextureIndex)
		}
		expectedMemo := ""
		if m.Memo != expectedMemo {
			t.Errorf("Expected Memo to be %v, got %v", expectedMemo, m.Memo)
		}
		expectedVerticesCount := 1584
		if m.VerticesCount != expectedVerticesCount {
			t.Errorf("Expected VerticesCount to be %v, got %v", expectedVerticesCount, m.VerticesCount)
		}
	}

	{
		b, _ := pmxModel.Bones.Get(5)
		expectedName := "上半身"
		if b.Name() != expectedName {
			t.Errorf("Expected Name to be %q, got %q", expectedName, b.Name())
		}
		expectedEnglishName := "J_Bip_C_Spine2"
		if b.EnglishName != expectedEnglishName {
			t.Errorf("Expected EnglishName to be %q, got %q", expectedEnglishName, b.EnglishName)
		}
		expectedPosition := vec3(0.0, 12.39097, -0.2011687)
		if !b.Position.MMD().NearEquals(expectedPosition, 1e-5) {
			t.Errorf("Expected Position to be %v, got %v", expectedPosition, b.Position)
		}
		expectedParentBoneIndex := 3
		if b.ParentIndex != expectedParentBoneIndex {
			t.Errorf("Expected ParentBoneIndex to be %v, got %v", expectedParentBoneIndex, b.ParentIndex)
		}
		expectedLayer := 0
		if b.Layer != expectedLayer {
			t.Errorf("Expected Layer to be %v, got %v", expectedLayer, b.Layer)
		}
		expectedBoneFlag := model.BONE_FLAG_CAN_ROTATE | model.BONE_FLAG_IS_VISIBLE | model.BONE_FLAG_CAN_MANIPULATE | model.BONE_FLAG_TAIL_IS_BONE
		if b.BoneFlag != expectedBoneFlag {
			t.Errorf("Expected BoneFlag to be %v, got %v", expectedBoneFlag, b.BoneFlag)
		}
		expectedTailPosition := vec3(0.0, 0.0, 0.0)
		if !b.TailPosition.MMD().NearEquals(expectedTailPosition, 1e-5) {
			t.Errorf("Expected TailPosition to be %v, got %v", expectedTailPosition, b.TailPosition)
		}
		expectedTailIndex := 6
		if b.TailIndex != expectedTailIndex {
			t.Errorf("Expected TailIndex to be %v, got %v", expectedTailIndex, b.TailIndex)
		}
	}

	{
		b, _ := pmxModel.Bones.Get(12)
		expectedName := "右目"
		if b.Name() != expectedName {
			t.Errorf("Expected Name to be %q, got %q", expectedName, b.Name())
		}
		expectedEnglishName := "J_Adj_R_FaceEye"
		if b.EnglishName != expectedEnglishName {
			t.Errorf("Expected EnglishName to be %q, got %q", expectedEnglishName, b.EnglishName)
		}
		expectedPosition := vec3(-0.1984593, 18.47478, 0.04549573)
		if !b.Position.MMD().NearEquals(expectedPosition, 1e-5) {
			t.Errorf("Expected Position to be %v, got %v", expectedPosition, b.Position)
		}
		expectedParentBoneIndex := 9
		if b.ParentIndex != expectedParentBoneIndex {
			t.Errorf("Expected ParentBoneIndex to be %v, got %v", expectedParentBoneIndex, b.ParentIndex)
		}
		expectedLayer := 0
		if b.Layer != expectedLayer {
			t.Errorf("Expected Layer to be %v, got %v", expectedLayer, b.Layer)
		}
		expectedBoneFlag := model.BONE_FLAG_CAN_ROTATE | model.BONE_FLAG_IS_VISIBLE | model.BONE_FLAG_CAN_MANIPULATE | model.BONE_FLAG_TAIL_IS_BONE | model.BONE_FLAG_IS_EXTERNAL_ROTATION
		if b.BoneFlag != expectedBoneFlag {
			t.Errorf("Expected BoneFlag to be %v, got %v", expectedBoneFlag, b.BoneFlag)
		}
		expectedTailPosition := vec3(0.0, 0.0, 0.0)
		if !b.TailPosition.MMD().NearEquals(expectedTailPosition, 1e-5) {
			t.Errorf("Expected TailPosition to be %v, got %v", expectedTailPosition, b.TailPosition)
		}
		expectedTailIndex := -1
		if b.TailIndex != expectedTailIndex {
			t.Errorf("Expected TailIndex to be %v, got %v", expectedTailIndex, b.TailIndex)
		}
		expectedEffectBoneIndex := 10
		if b.EffectIndex != expectedEffectBoneIndex {
			t.Errorf("Expected EffectorBoneIndex to be %v, got %v", expectedEffectBoneIndex, b.EffectIndex)
		}
		expectedEffectFactor := 0.3
		if math.Abs(b.EffectFactor-expectedEffectFactor) > 1e-5 {
			t.Errorf("Expected EffectorBoneIndex to be %v, got %v", expectedEffectFactor, b.EffectFactor)
		}
	}

	{
		b, _ := pmxModel.Bones.Get(28)
		expectedName := "左腕捩"
		if b.Name() != expectedName {
			t.Errorf("Expected Name to be %q, got %q", expectedName, b.Name())
		}
		expectedEnglishName := "arm_twist_L"
		if b.EnglishName != expectedEnglishName {
			t.Errorf("Expected EnglishName to be %q, got %q", expectedEnglishName, b.EnglishName)
		}
		expectedPosition := vec3(2.482529, 15.52887, 0.3184027)
		if !b.Position.MMD().NearEquals(expectedPosition, 1e-5) {
			t.Errorf("Expected Position to be %v, got %v", expectedPosition, b.Position)
		}
		expectedParentBoneIndex := 27
		if b.ParentIndex != expectedParentBoneIndex {
			t.Errorf("Expected ParentBoneIndex to be %v, got %v", expectedParentBoneIndex, b.ParentIndex)
		}
		expectedLayer := 0
		if b.Layer != expectedLayer {
			t.Errorf("Expected Layer to be %v, got %v", expectedLayer, b.Layer)
		}
		expectedBoneFlag := model.BONE_FLAG_CAN_ROTATE | model.BONE_FLAG_IS_VISIBLE | model.BONE_FLAG_CAN_MANIPULATE | model.BONE_FLAG_TAIL_IS_BONE | model.BONE_FLAG_HAS_FIXED_AXIS | model.BONE_FLAG_HAS_LOCAL_AXIS
		if b.BoneFlag != expectedBoneFlag {
			t.Errorf("Expected BoneFlag to be %v, got %v", expectedBoneFlag, b.BoneFlag)
		}
		expectedTailPosition := vec3(0.0, 0.0, 0.0)
		if !b.TailPosition.MMD().NearEquals(expectedTailPosition, 1e-5) {
			t.Errorf("Expected TailPosition to be %v, got %v", expectedTailPosition, b.TailPosition)
		}
		expectedTailIndex := -1
		if b.TailIndex != expectedTailIndex {
			t.Errorf("Expected TailIndex to be %v, got %v", expectedTailIndex, b.TailIndex)
		}
		expectedFixedAxis := vec3(0.819152, -0.5735764, -4.355049e-15)
		if !b.FixedAxis.MMD().NearEquals(expectedFixedAxis, 1e-5) {
			t.Errorf("Expected FixedAxis to be %v, got %v", expectedFixedAxis, b.FixedAxis)
		}
		expectedLocalAxisX := vec3(0.8191521, -0.5735765, -4.35505e-15)
		if !b.LocalAxisX.MMD().NearEquals(expectedLocalAxisX, 1e-5) {
			t.Errorf("Expected LocalAxisX to be %v, got %v", expectedLocalAxisX, b.LocalAxisX)
		}
		expectedLocalAxisZ := vec3(-3.567448e-15, 2.497953e-15, -1)
		if !b.LocalAxisZ.MMD().NearEquals(expectedLocalAxisZ, 1e-5) {
			t.Errorf("Expected LocalAxisZ to be %v, got %v", expectedLocalAxisZ, b.LocalAxisZ)
		}
	}

	{
		b, _ := pmxModel.Bones.Get(98)
		expectedName := "左足ＩＫ"
		if b.Name() != expectedName {
			t.Errorf("Expected Name to be %q, got %q", expectedName, b.Name())
		}
		expectedEnglishName := "leg_IK_L"
		if b.EnglishName != expectedEnglishName {
			t.Errorf("Expected EnglishName to be %q, got %q", expectedEnglishName, b.EnglishName)
		}
		expectedPosition := vec3(0.9644502, 1.647273, 0.4050385)
		if !b.Position.MMD().NearEquals(expectedPosition, 1e-5) {
			t.Errorf("Expected Position to be %v, got %v", expectedPosition, b.Position)
		}
		expectedParentBoneIndex := 97
		if b.ParentIndex != expectedParentBoneIndex {
			t.Errorf("Expected ParentBoneIndex to be %v, got %v", expectedParentBoneIndex, b.ParentIndex)
		}
		expectedLayer := 0
		if b.Layer != expectedLayer {
			t.Errorf("Expected Layer to be %v, got %v", expectedLayer, b.Layer)
		}
		expectedBoneFlag := model.BONE_FLAG_CAN_ROTATE | model.BONE_FLAG_IS_VISIBLE | model.BONE_FLAG_CAN_MANIPULATE | model.BONE_FLAG_IS_IK | model.BONE_FLAG_CAN_TRANSLATE
		if b.BoneFlag != expectedBoneFlag {
			t.Errorf("Expected BoneFlag to be %v, got %v", expectedBoneFlag, b.BoneFlag)
		}
		expectedTailPosition := vec3(0.0, 0.0, 1.0)
		if !b.TailPosition.MMD().NearEquals(expectedTailPosition, 1e-5) {
			t.Errorf("Expected TailPosition to be %v, got %v", expectedTailPosition, b.TailPosition)
		}
		expectedTailIndex := -1
		if b.TailIndex != expectedTailIndex {
			t.Errorf("Expected TailIndex to be %v, got %v", expectedTailIndex, b.TailIndex)
		}
		expectedIkTargetBoneIndex := 95
		if b.Ik == nil || b.Ik.BoneIndex != expectedIkTargetBoneIndex {
			t.Errorf("Expected IkTargetBoneIndex to be %v, got %v", expectedIkTargetBoneIndex, b.Ik.BoneIndex)
		}
		expectedIkLoopCount := 40
		if b.Ik == nil || b.Ik.LoopCount != expectedIkLoopCount {
			t.Errorf("Expected IkLoopCount to be %v, got %v", expectedIkLoopCount, b.Ik.LoopCount)
		}
		expectedIkLimitedDegree := 57.29578
		if b.Ik == nil || math.Abs(mmath.RadToDeg(b.Ik.UnitRotation.X)-expectedIkLimitedDegree) > 1e-5 {
			t.Errorf("Expected IkLimitedRadian to be %v, got %v", expectedIkLimitedDegree, mmath.RadToDeg(b.Ik.UnitRotation.X))
		}
		if b.Ik != nil {
			il := b.Ik.Links[0]
			expectedLinkBoneIndex := 94
			if il.BoneIndex != expectedLinkBoneIndex {
				t.Errorf("Expected LinkBoneIndex to be %v, got %v", expectedLinkBoneIndex, il.BoneIndex)
			}
			expectedAngleLimit := true
			if il.AngleLimit != expectedAngleLimit {
				t.Errorf("Expected AngleLimit to be %v, got %v", expectedAngleLimit, il.AngleLimit)
			}
			expectedMinAngleLimit := vec3(-180.0, 0.0, 0.0)
			if !il.MinAngleLimit.RadToDeg().NearEquals(expectedMinAngleLimit, 1e-5) {
				t.Errorf("Expected MinAngleLimit to be %v, got %v", expectedMinAngleLimit, il.MinAngleLimit.RadToDeg())
			}
			expectedMaxAngleLimit := vec3(-0.5, 0.0, 0.0)
			if !il.MaxAngleLimit.RadToDeg().NearEquals(expectedMaxAngleLimit, 1e-5) {
				t.Errorf("Expected MaxAngleLimit to be %v, got %v", expectedMaxAngleLimit, il.MaxAngleLimit.RadToDeg())
			}
		}
		if b.Ik != nil {
			il := b.Ik.Links[1]
			expectedLinkBoneIndex := 93
			if il.BoneIndex != expectedLinkBoneIndex {
				t.Errorf("Expected LinkBoneIndex to be %v, got %v", expectedLinkBoneIndex, il.BoneIndex)
			}
			expectedAngleLimit := false
			if il.AngleLimit != expectedAngleLimit {
				t.Errorf("Expected AngleLimit to be %v, got %v", expectedAngleLimit, il.AngleLimit)
			}
		}
	}

	{
		m, _ := pmxModel.Morphs.Get(2)
		expectedName := "にこり"
		if m.Name() != expectedName {
			t.Errorf("Expected Name to be %q, got %q", expectedName, m.Name())
		}
		expectedEnglishName := "Fcl_BRW_Fun"
		if m.EnglishName != expectedEnglishName {
			t.Errorf("Expected EnglishName to be %q, got %q", expectedEnglishName, m.EnglishName)
		}
		expectedMorphType := model.MORPH_TYPE_VERTEX
		if m.MorphType != expectedMorphType {
			t.Errorf("Expected MorphType to be %v, got %v", expectedMorphType, m.MorphType)
		}
		expectedOffsetCount := 68
		if len(m.Offsets) != expectedOffsetCount {
			t.Errorf("Expected OffsetCount to be %v, got %v", expectedOffsetCount, len(m.Offsets))
		}
		{
			o := m.Offsets[30].(*model.VertexMorphOffset)
			expectedVertexIndex := 19329
			if o.VertexIndex != expectedVertexIndex {
				t.Errorf("Expected VertexIndex to be %v, got %v", expectedVertexIndex, o.VertexIndex)
			}
			expectedPosition := vec3(-0.01209146, 0.1320038, -0.0121327)
			if !o.Position.MMD().NearEquals(expectedPosition, 1e-5) {
				t.Errorf("Expected Position to be %v, got %v", expectedPosition, o.Position)
			}
		}
	}

	{
		m, _ := pmxModel.Morphs.Get(111)
		expectedName := "いボーン"
		if m.Name() != expectedName {
			t.Errorf("Expected Name to be %q, got %q", expectedName, m.Name())
		}
		expectedEnglishName := "Fcl_MTH_I_Bone"
		if m.EnglishName != expectedEnglishName {
			t.Errorf("Expected EnglishName to be %q, got %q", expectedEnglishName, m.EnglishName)
		}
		expectedMorphType := model.MORPH_TYPE_BONE
		if m.MorphType != expectedMorphType {
			t.Errorf("Expected MorphType to be %v, got %v", expectedMorphType, m.MorphType)
		}
		expectedOffsetCount := 3
		if len(m.Offsets) != expectedOffsetCount {
			t.Errorf("Expected OffsetCount to be %v, got %v", expectedOffsetCount, len(m.Offsets))
		}
		{
			o := m.Offsets[1].(*model.BoneMorphOffset)
			expectedBoneIndex := 17
			if o.BoneIndex != expectedBoneIndex {
				t.Errorf("Expected BoneIndex to be %v, got %v", expectedBoneIndex, o.BoneIndex)
			}
			expectedPosition := vec3(0.0, 0.0, 0.0)
			if !o.Position.MMD().NearEquals(expectedPosition, 1e-5) {
				t.Errorf("Expected BoneIndex to be %v, got %v", expectedBoneIndex, o.BoneIndex)
			}
			expectedRotation := vec3(-6.000048, 3.952523e-14, -3.308919e-14)
			if !o.Rotation.ToDegrees().NearEquals(expectedRotation, 1e-5) {
				t.Errorf("Expected Rotation to be %v, got %v", expectedRotation, o.Rotation.ToDegrees())
			}
		}
	}

	{
		m, _ := pmxModel.Morphs.Get(122)
		expectedName := "なごみ材質"
		if m.Name() != expectedName {
			t.Errorf("Expected Name to be %q, got %q", expectedName, m.Name())
		}
		expectedEnglishName := "eye_Nagomi_Material"
		if m.EnglishName != expectedEnglishName {
			t.Errorf("Expected EnglishName to be %q, got %q", expectedEnglishName, m.EnglishName)
		}
		expectedMorphType := model.MORPH_TYPE_MATERIAL
		if m.MorphType != expectedMorphType {
			t.Errorf("Expected MorphType to be %v, got %v", expectedMorphType, m.MorphType)
		}
		expectedOffsetCount := 4
		if len(m.Offsets) != expectedOffsetCount {
			t.Errorf("Expected OffsetCount to be %v, got %v", expectedOffsetCount, len(m.Offsets))
		}
		{
			o := m.Offsets[3].(*model.MaterialMorphOffset)
			expectedMaterialIndex := 12
			if o.MaterialIndex != expectedMaterialIndex {
				t.Errorf("Expected MaterialIndex to be %v, got %v", expectedMaterialIndex, o.MaterialIndex)
			}
			expectedCalcMode := model.CALC_MODE_ADDITION
			if o.CalcMode != expectedCalcMode {
				t.Errorf("Expected CalcMode to be %v, got %v", expectedCalcMode, o.CalcMode)
			}
			expectedDiffuse := mmath.Vec4{X: 0.0, Y: 0.0, Z: 0.0, W: 1.0}
			if !o.Diffuse.NearEquals(expectedDiffuse, 1e-5) {
				t.Errorf("Expected Diffuse to be %v, got %v", expectedDiffuse, o.Diffuse)
			}
			expectedSpecular := mmath.Vec4{X: 0.0, Y: 0.0, Z: 0.0, W: 0.0}
			if !o.Specular.NearEquals(expectedSpecular, 1e-5) {
				t.Errorf("Expected Specular to be %v, got %v", expectedSpecular, o.Specular)
			}
			expectedAmbient := vec3(0.0, 0.0, 0.0)
			if !o.Ambient.NearEquals(expectedAmbient, 1e-5) {
				t.Errorf("Expected Ambient to be %v, got %v", expectedAmbient, o.Ambient)
			}
			expectedEdge := mmath.Vec4{X: 0.0, Y: 0.0, Z: 0.0, W: 0.0}
			if !o.Edge.NearEquals(expectedEdge, 1e-5) {
				t.Errorf("Expected Edge to be %v, got %v", expectedEdge, o.Edge)
			}
			expectedEdgeSize := 0.0
			if math.Abs(o.EdgeSize-expectedEdgeSize) > 1e-5 {
				t.Errorf("Expected EdgeSize to be %v, got %v", expectedEdgeSize, o.EdgeSize)
			}
			expectedTextureFactor := mmath.Vec4{X: 0.0, Y: 0.0, Z: 0.0, W: 0.0}
			if !o.TextureFactor.NearEquals(expectedTextureFactor, 1e-5) {
				t.Errorf("Expected TextureFactor to be %v, got %v", expectedTextureFactor, o.TextureFactor)
			}
			expectedSphereTextureFactor := mmath.Vec4{X: 0.0, Y: 0.0, Z: 0.0, W: 0.0}
			if !o.SphereTextureFactor.NearEquals(expectedSphereTextureFactor, 1e-5) {
				t.Errorf("Expected SphereTextureFactor to be %v, got %v", expectedSphereTextureFactor, o.SphereTextureFactor)
			}
			expectedToonTextureFactor := mmath.Vec4{X: 0.0, Y: 0.0, Z: 0.0, W: 0.0}
			if !o.ToonTextureFactor.NearEquals(expectedToonTextureFactor, 1e-5) {
				t.Errorf("Expected ToonTextureFactor to be %v, got %v", expectedToonTextureFactor, o.ToonTextureFactor)
			}
		}
	}

	{
		m, _ := pmxModel.Morphs.Get(137)
		expectedName := "ひそめ"
		if m.Name() != expectedName {
			t.Errorf("Expected Name to be %q, got %q", expectedName, m.Name())
		}
		expectedEnglishName := "brow_Frown"
		if m.EnglishName != expectedEnglishName {
			t.Errorf("Expected EnglishName to be %q, got %q", expectedEnglishName, m.EnglishName)
		}
		expectedMorphType := model.MORPH_TYPE_GROUP
		if m.MorphType != expectedMorphType {
			t.Errorf("Expected MorphType to be %v, got %v", expectedMorphType, m.MorphType)
		}
		expectedOffsetCount := 6
		if len(m.Offsets) != expectedOffsetCount {
			t.Errorf("Expected OffsetCount to be %v, got %v", expectedOffsetCount, len(m.Offsets))
		}
		{
			o := m.Offsets[2].(*model.GroupMorphOffset)
			expectedMorphIndex := 21
			if o.MorphIndex != expectedMorphIndex {
				t.Errorf("Expected MorphIndex to be %v, got %v", expectedMorphIndex, o.MorphIndex)
			}
			expectedFactor := 0.3
			if math.Abs(o.MorphFactor-expectedFactor) > 1e-5 {
				t.Errorf("Expected Factor to be %v, got %v", expectedFactor, o.MorphFactor)
			}
		}
	}

	{
		d, _ := pmxModel.DisplaySlots.Get(0)
		expectedName := "Root"
		if d.Name() != expectedName {
			t.Errorf("Expected Name to be %q, got %q", expectedName, d.Name())
		}
		expectedEnglishName := "Root"
		if d.EnglishName != expectedEnglishName {
			t.Errorf("Expected EnglishName to be %q, got %q", expectedEnglishName, d.EnglishName)
		}
		expectedSpecialFlag := model.SPECIAL_FLAG_ON
		if d.SpecialFlag != expectedSpecialFlag {
			t.Errorf("Expected SpecialFlag to be %v, got %v", expectedSpecialFlag, d.SpecialFlag)
		}
		expectedReferenceCount := 1
		if len(d.References) != expectedReferenceCount {
			t.Errorf("Expected ReferenceCount to be %v, got %v", expectedReferenceCount, len(d.References))
		}
		{
			r := d.References[0]
			expectedDisplayType := model.DISPLAY_TYPE_BONE
			if r.DisplayType != expectedDisplayType {
				t.Errorf("Expected DisplayType to be %v, got %v", expectedDisplayType, r.DisplayType)
			}
			expectedIndex := 0
			if r.DisplayIndex != expectedIndex {
				t.Errorf("Expected DisplayIndex to be %v, got %v", expectedIndex, r.DisplayIndex)
			}
		}
	}

	{
		d, _ := pmxModel.DisplaySlots.Get(1)
		expectedName := "表情"
		if d.Name() != expectedName {
			t.Errorf("Expected Name to be %q, got %q", expectedName, d.Name())
		}
		expectedEnglishName := "Exp"
		if d.EnglishName != expectedEnglishName {
			t.Errorf("Expected EnglishName to be %q, got %q", expectedEnglishName, d.EnglishName)
		}
		expectedSpecialFlag := model.SPECIAL_FLAG_ON
		if d.SpecialFlag != expectedSpecialFlag {
			t.Errorf("Expected SpecialFlag to be %v, got %v", expectedSpecialFlag, d.SpecialFlag)
		}
		expectedReferenceCount := 143
		if len(d.References) != expectedReferenceCount {
			t.Errorf("Expected ReferenceCount to be %v, got %v", expectedReferenceCount, len(d.References))
		}
		{
			r := d.References[50]
			expectedDisplayType := model.DISPLAY_TYPE_MORPH
			if r.DisplayType != expectedDisplayType {
				t.Errorf("Expected DisplayType to be %v, got %v", expectedDisplayType, r.DisplayType)
			}
			expectedIndex := 142
			if r.DisplayIndex != expectedIndex {
				t.Errorf("Expected DisplayIndex to be %v, got %v", expectedIndex, r.DisplayIndex)
			}
		}
	}

	{
		d, _ := pmxModel.DisplaySlots.Get(9)
		expectedName := "右指"
		if d.Name() != expectedName {
			t.Errorf("Expected Name to be %q, got %q", expectedName, d.Name())
		}
		expectedEnglishName := "right hand"
		if d.EnglishName != expectedEnglishName {
			t.Errorf("Expected EnglishName to be %q, got %q", expectedEnglishName, d.EnglishName)
		}
		expectedSpecialFlag := model.SPECIAL_FLAG_OFF
		if d.SpecialFlag != expectedSpecialFlag {
			t.Errorf("Expected SpecialFlag to be %v, got %v", expectedSpecialFlag, d.SpecialFlag)
		}
		expectedReferenceCount := 15
		if len(d.References) != expectedReferenceCount {
			t.Errorf("Expected ReferenceCount to be %v, got %v", expectedReferenceCount, len(d.References))
		}
		{
			r := d.References[7]
			expectedDisplayType := model.DISPLAY_TYPE_BONE
			if r.DisplayType != expectedDisplayType {
				t.Errorf("Expected DisplayType to be %v, got %v", expectedDisplayType, r.DisplayType)
			}
			expectedIndex := 81
			if r.DisplayIndex != expectedIndex {
				t.Errorf("Expected DisplayIndex to be %v, got %v", expectedIndex, r.DisplayIndex)
			}
		}
	}

	{
		r, _ := pmxModel.RigidBodies.Get(14)
		expectedName := "右腕"
		if r.Name() != expectedName {
			t.Errorf("Expected Name to be %q, got %q", expectedName, r.Name())
		}
		expectedEnglishName := "J_Bip_R_UpperArm"
		if r.EnglishName != expectedEnglishName {
			t.Errorf("Expected EnglishName to be %q, got %q", expectedEnglishName, r.EnglishName)
		}
		expectedBoneIndex := 61
		if r.BoneIndex != expectedBoneIndex {
			t.Errorf("Expected BoneIndex to be %v, got %v", expectedBoneIndex, r.BoneIndex)
		}
		expectedGroupIndex := byte(2)
		if r.CollisionGroup.Group != expectedGroupIndex {
			t.Errorf("Expected GroupIndex to be %v, got %v", expectedGroupIndex, r.CollisionGroup.Group)
		}
		expectedCollisionGroupMasks := []uint16{1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
		actualCollisionMasks := collisionMaskToSlice(r.CollisionGroup.Mask)
		for i := 0; i < 16; i++ {
			if actualCollisionMasks[i] != expectedCollisionGroupMasks[i] {
				t.Errorf("Expected CollisionGroupMask to be %v, got %v (%v)", expectedCollisionGroupMasks[i], actualCollisionMasks[i], i)
			}
		}
		expectedShapeType := model.SHAPE_CAPSULE
		if r.Shape != expectedShapeType {
			t.Errorf("Expected ShapeType to be %v, got %v", expectedShapeType, r.Shape)
		}
		expectedSize := vec3(0.5398922, 2.553789, 0.0)
		if !r.Size.NearEquals(expectedSize, 1e-5) {
			t.Errorf("Expected Size to be %v, got %v", expectedSize, r.Size)
		}
		expectedPosition := vec3(-2.52586, 15.45157, 0.241455)
		if !r.Position.MMD().NearEquals(expectedPosition, 1e-5) {
			t.Errorf("Expected Position to be %v, got %v", expectedPosition, r.Position)
		}
		expectedRotation := vec3(0.0, 0.0, 125.00)
		if !r.Rotation.RadToDeg().NearEquals(expectedRotation, 1e-5) {
			t.Errorf("Expected Rotation to be %v, got %v", expectedRotation, r.Rotation.RadToDeg())
		}
		expectedMass := 1.0
		if math.Abs(r.Param.Mass-expectedMass) > 1e-5 {
			t.Errorf("Expected Mass to be %v, got %v", expectedMass, r.Param.Mass)
		}
		expectedLinearDamping := 0.5
		if math.Abs(r.Param.LinearDamping-expectedLinearDamping) > 1e-5 {
			t.Errorf("Expected LinearDamping to be %v, got %v", expectedLinearDamping, r.Param.LinearDamping)
		}
		expectedAngularDamping := 0.5
		if math.Abs(r.Param.AngularDamping-expectedAngularDamping) > 1e-5 {
			t.Errorf("Expected AngularDamping to be %v, got %v", expectedAngularDamping, r.Param.AngularDamping)
		}
		expectedRestitution := 0.0
		if math.Abs(r.Param.Restitution-expectedRestitution) > 1e-5 {
			t.Errorf("Expected Restitution to be %v, got %v", expectedRestitution, r.Param.Restitution)
		}
		expectedFriction := 0.0
		if math.Abs(r.Param.Friction-expectedFriction) > 1e-5 {
			t.Errorf("Expected Friction to be %v, got %v", expectedFriction, r.Param.Friction)
		}
		expectedPhysicsType := model.PHYSICS_TYPE_STATIC
		if r.PhysicsType != expectedPhysicsType {
			t.Errorf("Expected PhysicsType to be %v, got %v", expectedPhysicsType, r.PhysicsType)
		}
	}

	{
		j, _ := pmxModel.Joints.Get(13)
		expectedName := "↓|頭|髪_06-01"
		if j.Name() != expectedName {
			t.Errorf("Expected Name to be %q, got %q", expectedName, j.Name())
		}
		expectedEnglishName := "↓|頭|髪_06-01"
		if j.EnglishName != expectedEnglishName {
			t.Errorf("Expected EnglishName to be %q, got %q", expectedEnglishName, j.EnglishName)
		}
		expectedRigidBodyIndexA := 5
		if j.RigidBodyIndexA != expectedRigidBodyIndexA {
			t.Errorf("Expected RigidBodyIndexA to be %v, got %v", expectedRigidBodyIndexA, j.RigidBodyIndexA)
		}
		expectedRigidBodyIndexB := 45
		if j.RigidBodyIndexB != expectedRigidBodyIndexB {
			t.Errorf("Expected RigidBodyIndexB to be %v, got %v", expectedRigidBodyIndexB, j.RigidBodyIndexB)
		}
		expectedPosition := vec3(-1.189768, 18.56266, 0.07277258)
		if !j.Param.Position.MMD().NearEquals(expectedPosition, 1e-5) {
			t.Errorf("Expected Position to be %v, got %v", expectedPosition, j.Param.Position)
		}
		expectedRotation := vec3(-15.10554, 91.26718, -4.187446)
		if !j.Param.Rotation.RadToDeg().NearEquals(expectedRotation, 1e-5) {
			t.Errorf("Expected Rotation to be %v, got %v", expectedRotation, j.Param.Rotation.RadToDeg())
		}
		expectedTranslationLimitMin := vec3(0.0, 0.0, 0.0)
		if !j.Param.TranslationLimitMin.NearEquals(expectedTranslationLimitMin, 1e-5) {
			t.Errorf("Expected TranslationLimitation1 to be %v, got %v", expectedTranslationLimitMin, j.Param.TranslationLimitMin)
		}
		expectedTranslationLimitMax := vec3(0.0, 0.0, 0.0)
		if !j.Param.TranslationLimitMax.NearEquals(expectedTranslationLimitMax, 1e-5) {
			t.Errorf("Expected TranslationLimitation2 to be %v, got %v", expectedTranslationLimitMax, j.Param.TranslationLimitMax)
		}
		expectedRotationLimitMin := vec3(-29.04, -14.3587, -29.04)
		if !j.Param.RotationLimitMin.RadToDeg().NearEquals(expectedRotationLimitMin, 1e-5) {
			t.Errorf("Expected RotationLimitation1 to be %v, got %v", expectedRotationLimitMin, j.Param.RotationLimitMin.RadToDeg())
		}
		expectedRotationLimitMax := vec3(29.04, 14.3587, 29.04)
		if !j.Param.RotationLimitMax.RadToDeg().NearEquals(expectedRotationLimitMax, 1e-5) {
			t.Errorf("Expected RotationLimitation2 to be %v, got %v", expectedRotationLimitMax, j.Param.RotationLimitMax.RadToDeg())
		}
		expectedSpringConstantTranslation := vec3(0.0, 0.0, 0.0)
		if !j.Param.SpringConstantTranslation.NearEquals(expectedSpringConstantTranslation, 1e-5) {
			t.Errorf("Expected SpringConstantTranslation to be %v, got %v", expectedSpringConstantTranslation, j.Param.SpringConstantTranslation)
		}
		expectedSpringConstantRotation := vec3(60.0, 29.6667, 60.0)
		if !j.Param.SpringConstantRotation.NearEquals(expectedSpringConstantRotation, 1e-5) {
			t.Errorf("Expected SpringConstantRotation to be %v, got %v", expectedSpringConstantRotation, j.Param.SpringConstantRotation)
		}
	}
}

// testResourcePath はテストリソースのパスを組み立てる。
func testResourcePath(name string) string {
	return filepath.Join("..", "..", "..", "..", "internal", "test_resources", name)
}

// vec3 は3次元ベクトルを生成する。
func vec3(x, y, z float64) mmath.Vec3 {
	return mmath.Vec3{r3.Vec{X: x, Y: y, Z: z}}
}

// collisionMaskToSlice は衝突マスクを衝突可否の配列へ展開する。
func collisionMaskToSlice(mask uint16) []uint16 {
	flags := []uint16{
		0x0001,
		0x0002,
		0x0004,
		0x0008,
		0x0010,
		0x0020,
		0x0040,
		0x0080,
		0x0100,
		0x0200,
		0x0400,
		0x0800,
		0x1000,
		0x2000,
		0x4000,
		0x8000,
	}
	result := make([]uint16, 16)
	for i, flag := range flags {
		if mask&flag == flag {
			result[i] = 0
		} else {
			result[i] = 1
		}
	}
	return result
}
