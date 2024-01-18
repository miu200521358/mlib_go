package pmx_reader

import (
	"math"
	"testing"

	"github.com/miu200521358/mlib_go/pkg/math/mvec2"
	"github.com/miu200521358/mlib_go/pkg/math/mvec3"
	"github.com/miu200521358/mlib_go/pkg/pmx/bone"
	"github.com/miu200521358/mlib_go/pkg/pmx/display_slot"
	"github.com/miu200521358/mlib_go/pkg/pmx/material"
	"github.com/miu200521358/mlib_go/pkg/pmx/morph"
	"github.com/miu200521358/mlib_go/pkg/pmx/rigidbody"
	"github.com/miu200521358/mlib_go/pkg/pmx/vertex/deform"
)

func TestPmxReader_ReadNameByFilepath(t *testing.T) {
	r := &PmxReader{}

	modelName, err := r.ReadNameByFilepath("../../../resources/test/サンプルモデル.pmx")

	expectedName := "v2配布用素体03"
	if modelName != expectedName {
		t.Errorf("Expected Name to be %q, got %q", expectedName, modelName)
	}

	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err)
	}
}

func TestPmxReader_ReadNameByFilepath_2_1(t *testing.T) {
	r := &PmxReader{}

	modelName, err := r.ReadNameByFilepath("../../../resources/test/サンプルモデル_PMX2.1_UTF-8.pmx")

	expectedName := "サンプルモデル迪卢克"
	if modelName != expectedName {
		t.Errorf("Expected Name to be %q, got %q", expectedName, modelName)
	}

	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err)
	}
}

func TestPmxReader_ReadNameByFilepath_NotExist(t *testing.T) {
	r := &PmxReader{}

	modelName, err := r.ReadNameByFilepath("../../../resources/test/サンプルモデル2.pmx")

	expectedName := ""
	if modelName != expectedName {
		t.Errorf("Expected Name to be %q, got %q", expectedName, modelName)
	}

	if err == nil {
		t.Errorf("Expected error to be not nil, got %q", err)
	}
}

func TestPmxReader_ReadByFilepath(t *testing.T) {
	r := &PmxReader{}

	model, err := r.ReadByFilepath("../../../resources/test/サンプルモデル.pmx")

	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err)
	}

	expectedSignature := "PMX "
	if model.Signature != expectedSignature {
		t.Errorf("Expected Signature to be %q, got %q", expectedSignature, model.Signature)
	}

	expectedVersion := 2.0
	if model.Version != expectedVersion {
		t.Errorf("Expected Version to be %f, got %f", expectedVersion, model.Version)
	}

	expectedName := "v2配布用素体03"
	if model.Name != expectedName {
		t.Errorf("Expected Name to be %q, got %q", expectedName, model.Name)
	}

	{
		v := model.Vertices.GetItem(13)
		expectedPosition := &mvec3.T{0.1565633, 16.62944, -0.2150156}
		if !v.Position.PracticallyEquals(expectedPosition, 1e-5) {
			t.Errorf("Expected Position to be %v, got %v", expectedPosition, v.Position)
		}
		expectedNormal := &mvec3.T{0.2274586, 0.6613649, -0.714744}
		if !v.Normal.PracticallyEquals(expectedNormal, 1e-5) {
			t.Errorf("Expected Normal to be %v, got %v", expectedNormal, v.Normal)
		}
		expectedUV := &mvec2.T{0.5112334, 0.1250942}
		if !v.UV.PracticallyEquals(expectedUV, 1e-5) {
			t.Errorf("Expected UV to be %v, got %v", expectedUV, v.UV)
		}
		expectedDeformType := deform.BDEF4
		if v.DeformType != expectedDeformType {
			t.Errorf("Expected DeformType to be %d, got %d", expectedDeformType, v.DeformType)
		}
		expectedDeform := deform.NewBdef4(7, 8, 25, 9, 0.6375693, 0.2368899, 0.06898639, 0.05655446)
		if v.Deform.GetAllIndexes()[0] != expectedDeform.Indexes[0] {
			t.Errorf("Expected Deform to be %v, got %v", expectedDeform.Indexes[0], v.Deform.GetAllIndexes()[0])
		}
		if v.Deform.GetAllIndexes()[1] != expectedDeform.Indexes[1] {
			t.Errorf("Expected Deform to be %v, got %v", expectedDeform.Indexes[1], v.Deform.GetAllIndexes()[1])
		}
		if v.Deform.GetAllIndexes()[2] != expectedDeform.Indexes[2] {
			t.Errorf("Expected Deform to be %v, got %v", expectedDeform.Indexes[2], v.Deform.GetAllIndexes()[2])
		}
		if v.Deform.GetAllIndexes()[3] != expectedDeform.Indexes[3] {
			t.Errorf("Expected Deform to be %v, got %v", expectedDeform.Indexes[3], v.Deform.GetAllIndexes()[3])
		}
		if math.Abs(v.Deform.GetAllWeights()[0]-expectedDeform.Weights[0]) > 1e-5 {
			t.Errorf("Expected Deform to be %v, got %v", expectedDeform.Weights[0], v.Deform.GetAllWeights()[0])
		}
		if math.Abs(v.Deform.GetAllWeights()[1]-expectedDeform.Weights[1]) > 1e-5 {
			t.Errorf("Expected Deform to be %v, got %v", expectedDeform.Weights[1], v.Deform.GetAllWeights()[1])
		}
		if math.Abs(v.Deform.GetAllWeights()[2]-expectedDeform.Weights[2]) > 1e-5 {
			t.Errorf("Expected Deform to be %v, got %v", expectedDeform.Weights[2], v.Deform.GetAllWeights()[2])
		}
		if math.Abs(v.Deform.GetAllWeights()[3]-expectedDeform.Weights[3]) > 1e-5 {
			t.Errorf("Expected Deform to be %v, got %v", expectedDeform.Weights[3], v.Deform.GetAllWeights()[3])
		}
		expectedEdgeFactor := 1.0
		if math.Abs(v.EdgeFactor-expectedEdgeFactor) > 1e-5 {
			t.Errorf("Expected EdgeFactor to be %v, got %v", expectedEdgeFactor, v.EdgeFactor)
		}
	}

	{
		v := model.Vertices.GetItem(120)
		expectedPosition := &mvec3.T{1.529492, 5.757646, 0.4527041}
		if !v.Position.PracticallyEquals(expectedPosition, 1e-5) {
			t.Errorf("Expected Position to be %v, got %v", expectedPosition, v.Position)
		}
		expectedNormal := &mvec3.T{0.9943396, 0.1054612, -0.0129031}
		if !v.Normal.PracticallyEquals(expectedNormal, 1e-5) {
			t.Errorf("Expected Normal to be %v, got %v", expectedNormal, v.Normal)
		}
		expectedUV := &mvec2.T{0.8788766, 0.7697825}
		if !v.UV.PracticallyEquals(expectedUV, 1e-5) {
			t.Errorf("Expected UV to be %v, got %v", expectedUV, v.UV)
		}
		expectedDeformType := deform.BDEF2
		if v.DeformType != expectedDeformType {
			t.Errorf("Expected DeformType to be %d, got %d", expectedDeformType, v.DeformType)
		}
		expectedDeform := deform.NewBdef2(109, 108, 0.9845969)
		if v.Deform.GetAllIndexes()[0] != expectedDeform.Indexes[0] {
			t.Errorf("Expected Deform to be %v, got %v", expectedDeform.Indexes[0], v.Deform.GetAllIndexes()[0])
		}
		if v.Deform.GetAllIndexes()[1] != expectedDeform.Indexes[1] {
			t.Errorf("Expected Deform to be %v, got %v", expectedDeform.Indexes[1], v.Deform.GetAllIndexes()[1])
		}
		if math.Abs(v.Deform.GetAllWeights()[0]-expectedDeform.Weights[0]) > 1e-5 {
			t.Errorf("Expected Deform to be %v, got %v", expectedDeform.Weights[0], v.Deform.GetAllWeights()[0])
		}
		if math.Abs(v.Deform.GetAllWeights()[1]-expectedDeform.Weights[1]) > 1e-5 {
			t.Errorf("Expected Deform to be %v, got %v", expectedDeform.Weights[1], v.Deform.GetAllWeights()[1])
		}
		expectedEdgeFactor := 1.0
		if math.Abs(v.EdgeFactor-expectedEdgeFactor) > 1e-5 {
			t.Errorf("Expected EdgeFactor to be %v, got %v", expectedEdgeFactor, v.EdgeFactor)
		}
	}

	{
		f := model.Faces.GetItem(19117)
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
		tex := model.Textures.GetItem(10)
		expectedName := "tex\\_13_Toon.bmp"
		if tex.Name != expectedName {
			t.Errorf("Expected Path to be %q, got %q", expectedName, tex.Name)
		}
	}

	{
		m := model.Materials.GetItem(10)
		expectedName := "00_EyeWhite_はぅ"
		if m.Name != expectedName {
			t.Errorf("Expected Path to be %q, got %q", expectedName, m.Name)
		}
		expectedEnglishName := "N00_000_00_EyeWhite_00_EYE (Instance)_Hau"
		if m.EnglishName != expectedEnglishName {
			t.Errorf("Expected Path to be %q, got %q", expectedEnglishName, m.EnglishName)
		}
		expectedDiffuseColor := &mvec3.T{1.0, 1.0, 1.0}
		if !m.DiffuseColor.PracticallyEquals(expectedDiffuseColor, 1e-5) {
			t.Errorf("Expected Diffuse to be %v, got %v", expectedDiffuseColor, m.DiffuseColor)
		}
		expectedDiffuseAlpha := 0.0
		if math.Abs(m.DiffuseAlpha-expectedDiffuseAlpha) > 1e-5 {
			t.Errorf("Expected DiffuseAlpha to be %v, got %v", expectedDiffuseAlpha, m.DiffuseAlpha)
		}
		expectedSpecularColor := &mvec3.T{0.0, 0.0, 0.0}
		if !m.SpecularColor.PracticallyEquals(expectedSpecularColor, 1e-5) {
			t.Errorf("Expected Specular to be %v, got %v", expectedSpecularColor, m.SpecularColor)
		}
		expectedSpecularPower := 0.0
		if math.Abs(m.SpecularPower-expectedSpecularPower) > 1e-5 {
			t.Errorf("Expected SpecularPower to be %v, got %v", expectedSpecularPower, m.SpecularPower)
		}
		expectedAmbientColor := &mvec3.T{0.5, 0.5, 0.5}
		if !m.AmbientColor.PracticallyEquals(expectedAmbientColor, 1e-5) {
			t.Errorf("Expected Ambient to be %v, got %v", expectedAmbientColor, m.AmbientColor)
		}
		expectedDrawFlag := material.DRAW_FLAG_GROUND_SHADOW | material.DRAW_FLAG_DRAWING_ON_SELF_SHADOW_MAPS | material.DRAW_FLAG_DRAWING_SELF_SHADOWS
		if m.DrawFlag != expectedDrawFlag {
			t.Errorf("Expected DrawFlag to be %v, got %v", expectedDrawFlag, m.DrawFlag)
		}
		expectedEdgeColor := &mvec3.T{0.2745098, 0.09019607, 0.1254902}
		if !m.EdgeColor.PracticallyEquals(expectedEdgeColor, 1e-5) {
			t.Errorf("Expected Edge to be %v, got %v", expectedEdgeColor, m.EdgeColor)
		}
		expectedEdgeAlpha := 1.0
		if math.Abs(m.EdgeAlpha-expectedEdgeAlpha) > 1e-5 {
			t.Errorf("Expected EdgeAlpha to be %v, got %v", expectedEdgeAlpha, m.EdgeAlpha)
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
		expectedSphereMode := material.SPHERE_MODE_ADDITION
		if m.SphereMode != expectedSphereMode {
			t.Errorf("Expected SphereMode to be %v, got %v", expectedSphereMode, m.SphereMode)
		}
		expectedToonSharingFlag := material.TOON_SHARING_INDIVIDUAL
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
		b := model.Bones.GetItem(5)
		expectedName := "上半身"
		if b.Name != expectedName {
			t.Errorf("Expected Name to be %q, got %q", expectedName, b.Name)
		}
		expectedEnglishName := "J_Bip_C_Spine2"
		if b.EnglishName != expectedEnglishName {
			t.Errorf("Expected EnglishName to be %q, got %q", expectedEnglishName, b.EnglishName)
		}
		expectedPosition := &mvec3.T{0.0, 12.39097, -0.2011687}
		if !b.Position.PracticallyEquals(expectedPosition, 1e-5) {
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
		expectedBoneFlag := bone.BONE_FLAG_CAN_ROTATE | bone.BONE_FLAG_IS_VISIBLE | bone.BONE_FLAG_CAN_MANIPULATE | bone.BONE_FLAG_TAIL_IS_BONE
		if b.BoneFlag != expectedBoneFlag {
			t.Errorf("Expected BoneFlag to be %v, got %v", expectedBoneFlag, b.BoneFlag)
		}
		expectedTailPosition := &mvec3.T{0.0, 0.0, 0.0}
		if !b.TailPosition.PracticallyEquals(expectedTailPosition, 1e-5) {
			t.Errorf("Expected TailPosition to be %v, got %v", expectedTailPosition, b.TailPosition)
		}
		expectedTailIndex := 6
		if b.TailIndex != expectedTailIndex {
			t.Errorf("Expected TailIndex to be %v, got %v", expectedTailIndex, b.TailIndex)
		}
	}

	{
		b := model.Bones.GetItem(12)
		expectedName := "右目"
		if b.Name != expectedName {
			t.Errorf("Expected Name to be %q, got %q", expectedName, b.Name)
		}
		expectedEnglishName := "J_Adj_R_FaceEye"
		if b.EnglishName != expectedEnglishName {
			t.Errorf("Expected EnglishName to be %q, got %q", expectedEnglishName, b.EnglishName)
		}
		expectedPosition := &mvec3.T{-0.1984593, 18.47478, 0.04549573}
		if !b.Position.PracticallyEquals(expectedPosition, 1e-5) {
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
		expectedBoneFlag := bone.BONE_FLAG_CAN_ROTATE | bone.BONE_FLAG_IS_VISIBLE | bone.BONE_FLAG_CAN_MANIPULATE | bone.BONE_FLAG_TAIL_IS_BONE | bone.BONE_FLAG_IS_EXTERNAL_ROTATION
		if b.BoneFlag != expectedBoneFlag {
			t.Errorf("Expected BoneFlag to be %v, got %v", expectedBoneFlag, b.BoneFlag)
		}
		expectedTailPosition := &mvec3.T{0.0, 0.0, 0.0}
		if !b.TailPosition.PracticallyEquals(expectedTailPosition, 1e-5) {
			t.Errorf("Expected TailPosition to be %v, got %v", expectedTailPosition, b.TailPosition)
		}
		expectedTailIndex := -1
		if b.TailIndex != expectedTailIndex {
			t.Errorf("Expected TailIndex to be %v, got %v", expectedTailIndex, b.TailIndex)
		}
		expectedEffectBoneIndex := 10
		if b.EffectIndex != expectedEffectBoneIndex {
			t.Errorf("Expected ExternalBoneIndex to be %v, got %v", expectedEffectBoneIndex, b.EffectIndex)
		}
		expectedEffectFactor := 0.3
		if math.Abs(b.EffectFactor-expectedEffectFactor) > 1e-5 {
			t.Errorf("Expected ExternalBoneIndex to be %v, got %v", expectedEffectFactor, b.EffectFactor)
		}
	}

	{
		b := model.Bones.GetItem(28)
		expectedName := "左腕捩"
		if b.Name != expectedName {
			t.Errorf("Expected Name to be %q, got %q", expectedName, b.Name)
		}
		expectedEnglishName := "arm_twist_L"
		if b.EnglishName != expectedEnglishName {
			t.Errorf("Expected EnglishName to be %q, got %q", expectedEnglishName, b.EnglishName)
		}
		expectedPosition := &mvec3.T{2.482529, 15.52887, 0.3184027}
		if !b.Position.PracticallyEquals(expectedPosition, 1e-5) {
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
		expectedBoneFlag := bone.BONE_FLAG_CAN_ROTATE | bone.BONE_FLAG_IS_VISIBLE | bone.BONE_FLAG_CAN_MANIPULATE | bone.BONE_FLAG_TAIL_IS_BONE | bone.BONE_FLAG_HAS_FIXED_AXIS | bone.BONE_FLAG_HAS_LOCAL_AXIS
		if b.BoneFlag != expectedBoneFlag {
			t.Errorf("Expected BoneFlag to be %v, got %v", expectedBoneFlag, b.BoneFlag)
		}
		expectedTailPosition := &mvec3.T{0.0, 0.0, 0.0}
		if !b.TailPosition.PracticallyEquals(expectedTailPosition, 1e-5) {
			t.Errorf("Expected TailPosition to be %v, got %v", expectedTailPosition, b.TailPosition)
		}
		expectedTailIndex := -1
		if b.TailIndex != expectedTailIndex {
			t.Errorf("Expected TailIndex to be %v, got %v", expectedTailIndex, b.TailIndex)
		}
		expectedFixedAxis := &mvec3.T{0.819152, -0.5735764, -4.355049e-15}
		if !b.FixedAxis.PracticallyEquals(expectedFixedAxis, 1e-5) {
			t.Errorf("Expected FixedAxis to be %v, got %v", expectedFixedAxis, b.FixedAxis)
		}
		expectedLocalAxisX := &mvec3.T{0.8191521, -0.5735765, -4.35505e-15}
		if !b.LocalAxisX.PracticallyEquals(expectedLocalAxisX, 1e-5) {
			t.Errorf("Expected LocalAxisX to be %v, got %v", expectedLocalAxisX, b.LocalAxisX)
		}
		expectedLocalAxisZ := &mvec3.T{-3.567448e-15, 2.497953e-15, -1}
		if !b.LocalAxisZ.PracticallyEquals(expectedLocalAxisZ, 1e-5) {
			t.Errorf("Expected LocalAxisZ to be %v, got %v", expectedLocalAxisZ, b.LocalAxisZ)
		}
	}

	{
		b := model.Bones.GetItem(98)
		expectedName := "左足ＩＫ"
		if b.Name != expectedName {
			t.Errorf("Expected Name to be %q, got %q", expectedName, b.Name)
		}
		expectedEnglishName := "leg_IK_L"
		if b.EnglishName != expectedEnglishName {
			t.Errorf("Expected EnglishName to be %q, got %q", expectedEnglishName, b.EnglishName)
		}
		expectedPosition := &mvec3.T{0.9644502, 1.647273, 0.4050385}
		if !b.Position.PracticallyEquals(expectedPosition, 1e-5) {
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
		expectedBoneFlag := bone.BONE_FLAG_CAN_ROTATE | bone.BONE_FLAG_IS_VISIBLE | bone.BONE_FLAG_CAN_MANIPULATE | bone.BONE_FLAG_IS_IK | bone.BONE_FLAG_CAN_TRANSLATE
		if b.BoneFlag != expectedBoneFlag {
			t.Errorf("Expected BoneFlag to be %v, got %v", expectedBoneFlag, b.BoneFlag)
		}
		expectedTailPosition := &mvec3.T{0.0, 0.0, 1.0}
		if !b.TailPosition.PracticallyEquals(expectedTailPosition, 1e-5) {
			t.Errorf("Expected TailPosition to be %v, got %v", expectedTailPosition, b.TailPosition)
		}
		expectedTailIndex := -1
		if b.TailIndex != expectedTailIndex {
			t.Errorf("Expected TailIndex to be %v, got %v", expectedTailIndex, b.TailIndex)
		}
		expectedIkTargetBoneIndex := 95
		if b.Ik.BoneIndex != expectedIkTargetBoneIndex {
			t.Errorf("Expected IkTargetBoneIndex to be %v, got %v", expectedIkTargetBoneIndex, b.Ik.BoneIndex)
		}
		expectedIkLoopCount := 40
		if b.Ik.LoopCount != expectedIkLoopCount {
			t.Errorf("Expected IkLoopCount to be %v, got %v", expectedIkLoopCount, b.Ik.LoopCount)
		}
		expectedIkLimitedDegree := 57.29578
		if math.Abs(b.Ik.UnitRotation.GetDegrees().GetX()-expectedIkLimitedDegree) > 1e-5 {
			t.Errorf("Expected IkLimitedRadian to be %v, got %v", expectedIkLimitedDegree, b.Ik.UnitRotation.GetDegrees().GetX())
		}
		{
			il := b.Ik.Links[0]
			expectedLinkBoneIndex := 94
			if il.BoneIndex != expectedLinkBoneIndex {
				t.Errorf("Expected LinkBoneIndex to be %v, got %v", expectedLinkBoneIndex, il.BoneIndex)
			}
			expectedAngleLimit := true
			if il.AngleLimit != expectedAngleLimit {
				t.Errorf("Expected AngleLimit to be %v, got %v", expectedAngleLimit, il.AngleLimit)
			}
			expectedMinAngleLimit := &mvec3.T{-180.0, 0.0, 0.0}
			if !il.MinAngleLimit.GetDegrees().PracticallyEquals(expectedMinAngleLimit, 1e-5) {
				t.Errorf("Expected MinAngleLimit to be %v, got %v", expectedMinAngleLimit, il.MinAngleLimit.GetDegrees())
			}
			expectedMaxAngleLimit := &mvec3.T{-0.5, 0.0, 0.0}
			if !il.MaxAngleLimit.GetDegrees().PracticallyEquals(expectedMaxAngleLimit, 1e-5) {
				t.Errorf("Expected MaxAngleLimit to be %v, got %v", expectedMaxAngleLimit, il.MaxAngleLimit.GetDegrees())
			}
		}
		{
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
		m := model.Morphs.GetItem(2)
		expectedName := "にこり"
		if m.Name != expectedName {
			t.Errorf("Expected Name to be %q, got %q", expectedName, m.Name)
		}
		expectedEnglishName := "Fcl_BRW_Fun"
		if m.EnglishName != expectedEnglishName {
			t.Errorf("Expected EnglishName to be %q, got %q", expectedEnglishName, m.EnglishName)
		}
		expectedMorphType := morph.MORPH_TYPE_VERTEX
		if m.MorphType != expectedMorphType {
			t.Errorf("Expected MorphType to be %v, got %v", expectedMorphType, m.MorphType)
		}
		expectedOffsetCount := 68
		if len(m.Offsets) != expectedOffsetCount {
			t.Errorf("Expected OffsetCount to be %v, got %v", expectedOffsetCount, len(m.Offsets))
		}
		{
			o := m.Offsets[30].(*morph.VertexMorphOffset)
			expectedVertexIndex := 19329
			if o.VertexIndex != expectedVertexIndex {
				t.Errorf("Expected VertexIndex to be %v, got %v", expectedVertexIndex, o.VertexIndex)
			}
			expectedPosition := &mvec3.T{-0.01209146, 0.1320038, -0.0121327}
			if !o.Position.PracticallyEquals(expectedPosition, 1e-5) {
				t.Errorf("Expected Position to be %v, got %v", expectedPosition, o.Position)
			}
		}
	}

	{
		m := model.Morphs.GetItem(111)
		expectedName := "いボーン"
		if m.Name != expectedName {
			t.Errorf("Expected Name to be %q, got %q", expectedName, m.Name)
		}
		expectedEnglishName := "Fcl_MTH_I_Bone"
		if m.EnglishName != expectedEnglishName {
			t.Errorf("Expected EnglishName to be %q, got %q", expectedEnglishName, m.EnglishName)
		}
		expectedMorphType := morph.MORPH_TYPE_BONE
		if m.MorphType != expectedMorphType {
			t.Errorf("Expected MorphType to be %v, got %v", expectedMorphType, m.MorphType)
		}
		expectedOffsetCount := 3
		if len(m.Offsets) != expectedOffsetCount {
			t.Errorf("Expected OffsetCount to be %v, got %v", expectedOffsetCount, len(m.Offsets))
		}
		{
			o := m.Offsets[1].(*morph.BoneMorphOffset)
			expectedBoneIndex := 17
			if o.BoneIndex != expectedBoneIndex {
				t.Errorf("Expected BoneIndex to be %v, got %v", expectedBoneIndex, o.BoneIndex)
			}
			expectedPosition := &mvec3.T{0.0, 0.0, 0.0}
			if !o.Position.PracticallyEquals(expectedPosition, 1e-5) {
				t.Errorf("Expected BoneIndex to be %v, got %v", expectedBoneIndex, o.BoneIndex)
			}
			expectedRotation := &mvec3.T{-6.000048, 3.952523e-14, -3.308919e-14}
			if !o.Rotation.GetDegrees().PracticallyEquals(expectedRotation, 1e-5) {
				t.Errorf("Expected Rotation to be %v, got %v", expectedRotation, o.Rotation.GetDegrees())
			}
		}
	}

	{
		m := model.Morphs.GetItem(122)
		expectedName := "なごみ材質"
		if m.Name != expectedName {
			t.Errorf("Expected Name to be %q, got %q", expectedName, m.Name)
		}
		expectedEnglishName := "eye_Nagomi_Material"
		if m.EnglishName != expectedEnglishName {
			t.Errorf("Expected EnglishName to be %q, got %q", expectedEnglishName, m.EnglishName)
		}
		expectedMorphType := morph.MORPH_TYPE_MATERIAL
		if m.MorphType != expectedMorphType {
			t.Errorf("Expected MorphType to be %v, got %v", expectedMorphType, m.MorphType)
		}
		expectedOffsetCount := 4
		if len(m.Offsets) != expectedOffsetCount {
			t.Errorf("Expected OffsetCount to be %v, got %v", expectedOffsetCount, len(m.Offsets))
		}
		{
			o := m.Offsets[3].(*morph.MaterialMorphOffset)
			expectedMaterialIndex := 12
			if o.MaterialIndex != expectedMaterialIndex {
				t.Errorf("Expected MaterialIndex to be %v, got %v", expectedMaterialIndex, o.MaterialIndex)
			}
			expectedCalcMode := morph.CALC_MODE_ADDITION
			if o.CalcMode != expectedCalcMode {
				t.Errorf("Expected CalcMode to be %v, got %v", expectedCalcMode, o.CalcMode)
			}
			expectedDiffuseColor := &mvec3.T{0.0, 0.0, 0.0}
			if !o.DiffuseColor.PracticallyEquals(expectedDiffuseColor, 1e-5) {
				t.Errorf("Expected DiffuseColor to be %v, got %v", expectedDiffuseColor, o.DiffuseColor)
			}
			expectedDiffuseAlpha := 1.0
			if math.Abs(o.DiffuseAlpha-expectedDiffuseAlpha) > 1e-5 {
				t.Errorf("Expected DiffuseAlpha to be %v, got %v", expectedDiffuseAlpha, o.DiffuseAlpha)
			}
			expectedSpecularColor := &mvec3.T{0.0, 0.0, 0.0}
			if !o.SpecularColor.PracticallyEquals(expectedSpecularColor, 1e-5) {
				t.Errorf("Expected SpecularColor to be %v, got %v", expectedSpecularColor, o.SpecularColor)
			}
			expectedSpecularPower := 0.0
			if math.Abs(o.SpecularPower-expectedSpecularPower) > 1e-5 {
				t.Errorf("Expected SpecularPower to be %v, got %v", expectedSpecularPower, o.SpecularPower)
			}
			expectedAmbientColor := &mvec3.T{0.0, 0.0, 0.0}
			if !o.AmbientColor.PracticallyEquals(expectedAmbientColor, 1e-5) {
				t.Errorf("Expected AmbientColor to be %v, got %v", expectedAmbientColor, o.AmbientColor)
			}
			expectedEdgeColor := &mvec3.T{0.0, 0.0, 0.0}
			if !o.EdgeColor.PracticallyEquals(expectedEdgeColor, 1e-5) {
				t.Errorf("Expected EdgeColor to be %v, got %v", expectedEdgeColor, o.EdgeColor)
			}
			expectedEdgeSize := 0.0
			if math.Abs(o.EdgeSize-expectedEdgeSize) > 1e-5 {
				t.Errorf("Expected EdgeSize to be %v, got %v", expectedEdgeSize, o.EdgeSize)
			}
			expectedTextureCoefficient := &mvec3.T{0.0, 0.0, 0.0}
			if !o.TextureCoefficient.PracticallyEquals(expectedTextureCoefficient, 1e-5) {
				t.Errorf("Expected TextureFactor to be %v, got %v", expectedTextureCoefficient, o.TextureCoefficient)
			}
			expectedTextureAlpha := 0.0
			if math.Abs(o.TextureAlpha-expectedTextureAlpha) > 1e-5 {
				t.Errorf("Expected TextureAlpha to be %v, got %v", expectedTextureAlpha, o.TextureAlpha)
			}
			expectedSphereTextureCoefficient := &mvec3.T{0.0, 0.0, 0.0}
			if !o.SphereTextureCoefficient.PracticallyEquals(expectedSphereTextureCoefficient, 1e-5) {
				t.Errorf("Expected SphereTextureFactor to be %v, got %v", expectedSphereTextureCoefficient, o.SphereTextureCoefficient)
			}
			expectedSphereTextureAlpha := 0.0
			if math.Abs(o.SphereTextureAlpha-expectedSphereTextureAlpha) > 1e-5 {
				t.Errorf("Expected SphereTextureAlpha to be %v, got %v", expectedSphereTextureAlpha, o.SphereTextureAlpha)
			}
			expectedToonTextureCoefficient := &mvec3.T{0.0, 0.0, 0.0}
			if !o.ToonTextureCoefficient.PracticallyEquals(expectedToonTextureCoefficient, 1e-5) {
				t.Errorf("Expected ToonTextureFactor to be %v, got %v", expectedToonTextureCoefficient, o.ToonTextureCoefficient)
			}
			expectedToonTextureAlpha := 0.0
			if math.Abs(o.ToonTextureAlpha-expectedToonTextureAlpha) > 1e-5 {
				t.Errorf("Expected ToonTextureAlpha to be %v, got %v", expectedToonTextureAlpha, o.ToonTextureAlpha)
			}
		}
	}

	{
		m := model.Morphs.GetItem(137)
		expectedName := "ひそめ"
		if m.Name != expectedName {
			t.Errorf("Expected Name to be %q, got %q", expectedName, m.Name)
		}
		expectedEnglishName := "brow_Frown"
		if m.EnglishName != expectedEnglishName {
			t.Errorf("Expected EnglishName to be %q, got %q", expectedEnglishName, m.EnglishName)
		}
		expectedMorphType := morph.MORPH_TYPE_GROUP
		if m.MorphType != expectedMorphType {
			t.Errorf("Expected MorphType to be %v, got %v", expectedMorphType, m.MorphType)
		}
		expectedOffsetCount := 6
		if len(m.Offsets) != expectedOffsetCount {
			t.Errorf("Expected OffsetCount to be %v, got %v", expectedOffsetCount, len(m.Offsets))
		}
		{
			o := m.Offsets[2].(*morph.GroupMorphOffset)
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
		d := model.DisplaySlots.GetItem(0)
		expectedName := "Root"
		if d.Name != expectedName {
			t.Errorf("Expected Name to be %q, got %q", expectedName, d.Name)
		}
		expectedEnglishName := "Root"
		if d.EnglishName != expectedEnglishName {
			t.Errorf("Expected EnglishName to be %q, got %q", expectedEnglishName, d.EnglishName)
		}
		expectedSpecialFlag := display_slot.SPECIAL_FLAG_ON
		if d.SpecialFlag != expectedSpecialFlag {
			t.Errorf("Expected SpecialFlag to be %v, got %v", expectedSpecialFlag, d.SpecialFlag)
		}
		expectedReferenceCount := 1
		if len(d.References) != expectedReferenceCount {
			t.Errorf("Expected ReferenceCount to be %v, got %v", expectedReferenceCount, len(d.References))
		}
		{
			r := d.References[0]
			expectedDisplayType := display_slot.DISPLAY_TYPE_BONE
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
		d := model.DisplaySlots.GetItem(1)
		expectedName := "表情"
		if d.Name != expectedName {
			t.Errorf("Expected Name to be %q, got %q", expectedName, d.Name)
		}
		expectedEnglishName := "Exp"
		if d.EnglishName != expectedEnglishName {
			t.Errorf("Expected EnglishName to be %q, got %q", expectedEnglishName, d.EnglishName)
		}
		expectedSpecialFlag := display_slot.SPECIAL_FLAG_ON
		if d.SpecialFlag != expectedSpecialFlag {
			t.Errorf("Expected SpecialFlag to be %v, got %v", expectedSpecialFlag, d.SpecialFlag)
		}
		expectedReferenceCount := 143
		if len(d.References) != expectedReferenceCount {
			t.Errorf("Expected ReferenceCount to be %v, got %v", expectedReferenceCount, len(d.References))
		}
		{
			r := d.References[50]
			expectedDisplayType := display_slot.DISPLAY_TYPE_MORPH
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
		d := model.DisplaySlots.GetItem(9)
		expectedName := "右指"
		if d.Name != expectedName {
			t.Errorf("Expected Name to be %q, got %q", expectedName, d.Name)
		}
		expectedEnglishName := "right hand"
		if d.EnglishName != expectedEnglishName {
			t.Errorf("Expected EnglishName to be %q, got %q", expectedEnglishName, d.EnglishName)
		}
		expectedSpecialFlag := display_slot.SPECIAL_FLAG_OFF
		if d.SpecialFlag != expectedSpecialFlag {
			t.Errorf("Expected SpecialFlag to be %v, got %v", expectedSpecialFlag, d.SpecialFlag)
		}
		expectedReferenceCount := 15
		if len(d.References) != expectedReferenceCount {
			t.Errorf("Expected ReferenceCount to be %v, got %v", expectedReferenceCount, len(d.References))
		}
		{
			r := d.References[7]
			expectedDisplayType := display_slot.DISPLAY_TYPE_BONE
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
		b := model.RigidBodies.GetItem(14)
		expectedName := "右腕"
		if b.Name != expectedName {
			t.Errorf("Expected Name to be %q, got %q", expectedName, b.Name)
		}
		expectedEnglishName := "J_Bip_R_UpperArm"
		if b.EnglishName != expectedEnglishName {
			t.Errorf("Expected EnglishName to be %q, got %q", expectedEnglishName, b.EnglishName)
		}
		expectedBoneIndex := 61
		if b.BoneIndex != expectedBoneIndex {
			t.Errorf("Expected BoneIndex to be %v, got %v", expectedBoneIndex, b.BoneIndex)
		}
		expectedGroupIndex := byte(2)
		if b.CollisionGroup != expectedGroupIndex {
			t.Errorf("Expected GroupIndex to be %v, got %v", expectedGroupIndex, b.CollisionGroup)
		}
		expectedCollisionGroupMasks := []uint16{1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
		for i := 0; i < 16; i++ {
			if b.CollisionGroupMask.IsCollisions[i] != expectedCollisionGroupMasks[i] {
				t.Errorf("Expected CollisionGroupMask to be %v, got %v (%v)", expectedCollisionGroupMasks[i], b.CollisionGroupMask.IsCollisions[i], i)
			}
		}
		expectedShapeType := rigidbody.SHAPE_CAPSULE
		if b.ShapeType != expectedShapeType {
			t.Errorf("Expected ShapeType to be %v, got %v", expectedShapeType, b.ShapeType)
		}
		expectedSize := &mvec3.T{0.5398922, 2.553789, 0.0}
		if !b.Size.PracticallyEquals(expectedSize, 1e-5) {
			t.Errorf("Expected Size to be %v, got %v", expectedSize, b.Size)
		}
		expectedPosition := &mvec3.T{-2.52586, 15.45157, 0.241455}
		if !b.Position.PracticallyEquals(expectedPosition, 1e-5) {
			t.Errorf("Expected Position to be %v, got %v", expectedPosition, b.Position)
		}
		expectedRotation := &mvec3.T{0.0, 0.0, 125.00}
		if !b.Rotation.GetDegrees().PracticallyEquals(expectedRotation, 1e-5) {
			t.Errorf("Expected Rotation to be %v, got %v", expectedRotation, b.Rotation.GetDegrees())
		}
		expectedMass := 1.0
		if math.Abs(b.Param.Mass-expectedMass) > 1e-5 {
			t.Errorf("Expected Mass to be %v, got %v", expectedMass, b.Param.Mass)
		}
		expectedLinearDamping := 0.5
		if math.Abs(b.Param.LinearDamping-expectedLinearDamping) > 1e-5 {
			t.Errorf("Expected LinearDamping to be %v, got %v", expectedLinearDamping, b.Param.LinearDamping)
		}
		expectedAngularDamping := 0.5
		if math.Abs(b.Param.AngularDamping-expectedAngularDamping) > 1e-5 {
			t.Errorf("Expected AngularDamping to be %v, got %v", expectedAngularDamping, b.Param.AngularDamping)
		}
		expectedRestitution := 0.0
		if math.Abs(b.Param.Restitution-expectedRestitution) > 1e-5 {
			t.Errorf("Expected Restitution to be %v, got %v", expectedRestitution, b.Param.Restitution)
		}
		expectedFriction := 0.0
		if math.Abs(b.Param.Friction-expectedFriction) > 1e-5 {
			t.Errorf("Expected Friction to be %v, got %v", expectedFriction, b.Param.Friction)
		}
		expectedPhysicsType := rigidbody.PHYSICS_TYPE_STATIC
		if b.PhysicsType != expectedPhysicsType {
			t.Errorf("Expected PhysicsType to be %v, got %v", expectedPhysicsType, b.PhysicsType)
		}
	}
}

func TestPmxReader_ReadByFilepath_2_1(t *testing.T) {
	r := &PmxReader{}

	model, err := r.ReadByFilepath("../../../resources/test/サンプルモデル_PMX2.1_UTF-8.pmx")

	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err)
	}

	expectedSignature := "PMX "
	if model.Signature != expectedSignature {
		t.Errorf("Expected Signature to be %q, got %q", expectedSignature, model.Signature)
	}

	expectedVersion := 2.1
	if math.Abs(model.Version-expectedVersion) > 1e-5 {
		t.Errorf("Expected Version to be %.8f, got %.8f", expectedVersion, model.Version)
	}

	expectedName := "サンプルモデル迪卢克"
	if model.Name != expectedName {
		t.Errorf("Expected Name to be %q, got %q", expectedName, model.Name)
	}
}
