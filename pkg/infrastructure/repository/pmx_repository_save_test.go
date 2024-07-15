package repository

import (
	"math"
	"testing"

	"github.com/miu200521358/mlib_go/pkg/domain/mmath"
	"github.com/miu200521358/mlib_go/pkg/domain/pmx"
)

func TestPmxWriter_Save1(t *testing.T) {
	r := NewPmxRepository()

	data, err := r.Load("../../../test_resources/サンプルモデル_PMX読み取り確認用.pmx")
	originalModel := data.(*pmx.PmxModel)

	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err)
	}

	// ------------------

	overridePath := "../../../test_resources/サンプルモデル_PMX読み取り確認用_output.pmx"
	r.Save(overridePath, originalModel, false)

	// ------------------
	data, err = r.Load(overridePath)
	model := data.(*pmx.PmxModel)

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
		v := model.Vertices.Get(13)
		expectedPosition := &mmath.MVec3{0.1565633, 16.62944, -0.2150156}
		if !v.Position.MMD().NearEquals(expectedPosition, 1e-5) {
			t.Errorf("Expected Position to be %v, got %v", expectedPosition, v.Position)
		}
		expectedNormal := &mmath.MVec3{0.2274586, 0.6613649, -0.714744}
		if !v.Normal.MMD().NearEquals(expectedNormal, 1e-5) {
			t.Errorf("Expected Normal to be %v, got %v", expectedNormal, v.Normal)
		}
		expectedUV := &mmath.MVec2{0.5112334, 0.1250942}
		if !v.Uv.NearEquals(expectedUV, 1e-5) {
			t.Errorf("Expected UV to be %v, got %v", expectedUV, v.Uv)
		}
		expectedDeformType := pmx.BDEF4
		if v.DeformType != expectedDeformType {
			t.Errorf("Expected DeformType to be %d, got %d", expectedDeformType, v.DeformType)
		}
		expectedDeform := pmx.NewBdef4(7, 8, 25, 9, 0.6375693, 0.2368899, 0.06898639, 0.05655446)
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
		v := model.Vertices.Get(120)
		expectedPosition := &mmath.MVec3{1.529492, 5.757646, 0.4527041}
		if !v.Position.MMD().NearEquals(expectedPosition, 1e-5) {
			t.Errorf("Expected Position to be %v, got %v", expectedPosition, v.Position)
		}
		expectedNormal := &mmath.MVec3{0.9943396, 0.1054612, -0.0129031}
		if !v.Normal.MMD().NearEquals(expectedNormal, 1e-5) {
			t.Errorf("Expected Normal to be %v, got %v", expectedNormal, v.Normal)
		}
		expectedUV := &mmath.MVec2{0.8788766, 0.7697825}
		if !v.Uv.NearEquals(expectedUV, 1e-5) {
			t.Errorf("Expected UV to be %v, got %v", expectedUV, v.Uv)
		}
		expectedDeformType := pmx.BDEF2
		if v.DeformType != expectedDeformType {
			t.Errorf("Expected DeformType to be %d, got %d", expectedDeformType, v.DeformType)
		}
		expectedDeform := pmx.NewBdef2(109, 108, 0.9845969)
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
		f := model.Faces.Get(19117)
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
		tex := model.Textures.Get(10)
		expectedName := "tex\\_13_Toon.bmp"
		if tex.Name != expectedName {
			t.Errorf("Expected Path to be %q, got %q", expectedName, tex.Name)
		}
	}

	{
		m := model.Materials.Get(10)
		expectedName := "00_EyeWhite_はぅ"
		if m.Name != expectedName {
			t.Errorf("Expected Path to be %q, got %q", expectedName, m.Name)
		}
		expectedEnglishName := "N00_000_00_EyeWhite_00_EYE (Instance)_Hau"
		if m.EnglishName != expectedEnglishName {
			t.Errorf("Expected Path to be %q, got %q", expectedEnglishName, m.EnglishName)
		}
		expectedDiffuse := &mmath.MVec4{1.0, 1.0, 1.0, 0.0}
		if !m.Diffuse.NearEquals(expectedDiffuse, 1e-5) {
			t.Errorf("Expected Diffuse to be %v, got %v", expectedDiffuse, m.Diffuse)
		}
		expectedSpecular := &mmath.MVec4{0.0, 0.0, 0.0, 0.0}
		if !m.Specular.NearEquals(expectedSpecular, 1e-5) {
			t.Errorf("Expected Specular to be %v, got %v", expectedSpecular, m.Specular)
		}
		expectedAmbient := &mmath.MVec3{0.5, 0.5, 0.5}
		if !m.Ambient.NearEquals(expectedAmbient, 1e-5) {
			t.Errorf("Expected Ambient to be %v, got %v", expectedAmbient, m.Ambient)
		}
		expectedDrawFlag := pmx.DRAW_FLAG_GROUND_SHADOW | pmx.DRAW_FLAG_DRAWING_ON_SELF_SHADOW_MAPS | pmx.DRAW_FLAG_DRAWING_SELF_SHADOWS
		if m.DrawFlag != expectedDrawFlag {
			t.Errorf("Expected DrawFlag to be %v, got %v", expectedDrawFlag, m.DrawFlag)
		}
		expectedEdge := &mmath.MVec4{0.2745098, 0.09019607, 0.1254902, 1.0}
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
		expectedSphereMode := pmx.SPHERE_MODE_ADDITION
		if m.SphereMode != expectedSphereMode {
			t.Errorf("Expected SphereMode to be %v, got %v", expectedSphereMode, m.SphereMode)
		}
		expectedToonSharingFlag := pmx.TOON_SHARING_INDIVIDUAL
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
		b := model.Bones.Get(5)
		expectedName := "上半身"
		if b.Name != expectedName {
			t.Errorf("Expected Name to be %q, got %q", expectedName, b.Name)
		}
		expectedEnglishName := "J_Bip_C_Spine2"
		if b.EnglishName != expectedEnglishName {
			t.Errorf("Expected EnglishName to be %q, got %q", expectedEnglishName, b.EnglishName)
		}
		expectedPosition := &mmath.MVec3{0.0, 12.39097, -0.2011687}
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
		expectedBoneFlag := pmx.BONE_FLAG_CAN_ROTATE | pmx.BONE_FLAG_IS_VISIBLE | pmx.BONE_FLAG_CAN_MANIPULATE | pmx.BONE_FLAG_TAIL_IS_BONE
		if b.BoneFlag != expectedBoneFlag {
			t.Errorf("Expected BoneFlag to be %v, got %v", expectedBoneFlag, b.BoneFlag)
		}
		expectedTailPosition := &mmath.MVec3{0.0, 0.0, 0.0}
		if !b.TailPosition.MMD().NearEquals(expectedTailPosition, 1e-5) {
			t.Errorf("Expected TailPosition to be %v, got %v", expectedTailPosition, b.TailPosition)
		}
		expectedTailIndex := 6
		if b.TailIndex != expectedTailIndex {
			t.Errorf("Expected TailIndex to be %v, got %v", expectedTailIndex, b.TailIndex)
		}
	}

	{
		b := model.Bones.Get(12)
		expectedName := "右目"
		if b.Name != expectedName {
			t.Errorf("Expected Name to be %q, got %q", expectedName, b.Name)
		}
		expectedEnglishName := "J_Adj_R_FaceEye"
		if b.EnglishName != expectedEnglishName {
			t.Errorf("Expected EnglishName to be %q, got %q", expectedEnglishName, b.EnglishName)
		}
		expectedPosition := &mmath.MVec3{-0.1984593, 18.47478, 0.04549573}
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
		expectedBoneFlag := pmx.BONE_FLAG_CAN_ROTATE | pmx.BONE_FLAG_IS_VISIBLE | pmx.BONE_FLAG_CAN_MANIPULATE | pmx.BONE_FLAG_TAIL_IS_BONE | pmx.BONE_FLAG_IS_EXTERNAL_ROTATION
		if b.BoneFlag != expectedBoneFlag {
			t.Errorf("Expected BoneFlag to be %v, got %v", expectedBoneFlag, b.BoneFlag)
		}
		expectedTailPosition := &mmath.MVec3{0.0, 0.0, 0.0}
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
		b := model.Bones.Get(28)
		expectedName := "左腕捩"
		if b.Name != expectedName {
			t.Errorf("Expected Name to be %q, got %q", expectedName, b.Name)
		}
		expectedEnglishName := "arm_twist_L"
		if b.EnglishName != expectedEnglishName {
			t.Errorf("Expected EnglishName to be %q, got %q", expectedEnglishName, b.EnglishName)
		}
		expectedPosition := &mmath.MVec3{2.482529, 15.52887, 0.3184027}
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
		expectedBoneFlag := pmx.BONE_FLAG_CAN_ROTATE | pmx.BONE_FLAG_IS_VISIBLE | pmx.BONE_FLAG_CAN_MANIPULATE | pmx.BONE_FLAG_TAIL_IS_BONE | pmx.BONE_FLAG_HAS_FIXED_AXIS | pmx.BONE_FLAG_HAS_LOCAL_AXIS
		if b.BoneFlag != expectedBoneFlag {
			t.Errorf("Expected BoneFlag to be %v, got %v", expectedBoneFlag, b.BoneFlag)
		}
		expectedTailPosition := &mmath.MVec3{0.0, 0.0, 0.0}
		if !b.TailPosition.MMD().NearEquals(expectedTailPosition, 1e-5) {
			t.Errorf("Expected TailPosition to be %v, got %v", expectedTailPosition, b.TailPosition)
		}
		expectedTailIndex := -1
		if b.TailIndex != expectedTailIndex {
			t.Errorf("Expected TailIndex to be %v, got %v", expectedTailIndex, b.TailIndex)
		}
		expectedFixedAxis := &mmath.MVec3{0.819152, -0.5735764, -4.355049e-15}
		if !b.FixedAxis.MMD().NearEquals(expectedFixedAxis, 1e-5) {
			t.Errorf("Expected FixedAxis to be %v, got %v", expectedFixedAxis, b.FixedAxis)
		}
		expectedLocalAxisX := &mmath.MVec3{0.8191521, -0.5735765, -4.35505e-15}
		if !b.LocalAxisX.MMD().NearEquals(expectedLocalAxisX, 1e-5) {
			t.Errorf("Expected LocalAxisX to be %v, got %v", expectedLocalAxisX, b.LocalAxisX)
		}
		expectedLocalAxisZ := &mmath.MVec3{-3.567448e-15, 2.497953e-15, -1}
		if !b.LocalAxisZ.MMD().NearEquals(expectedLocalAxisZ, 1e-5) {
			t.Errorf("Expected LocalAxisZ to be %v, got %v", expectedLocalAxisZ, b.LocalAxisZ)
		}
	}

	{
		b := model.Bones.Get(98)
		expectedName := "左足ＩＫ"
		if b.Name != expectedName {
			t.Errorf("Expected Name to be %q, got %q", expectedName, b.Name)
		}
		expectedEnglishName := "leg_IK_L"
		if b.EnglishName != expectedEnglishName {
			t.Errorf("Expected EnglishName to be %q, got %q", expectedEnglishName, b.EnglishName)
		}
		expectedPosition := &mmath.MVec3{0.9644502, 1.647273, 0.4050385}
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
		expectedBoneFlag := pmx.BONE_FLAG_CAN_ROTATE | pmx.BONE_FLAG_IS_VISIBLE | pmx.BONE_FLAG_CAN_MANIPULATE | pmx.BONE_FLAG_IS_IK | pmx.BONE_FLAG_CAN_TRANSLATE
		if b.BoneFlag != expectedBoneFlag {
			t.Errorf("Expected BoneFlag to be %v, got %v", expectedBoneFlag, b.BoneFlag)
		}
		expectedTailPosition := &mmath.MVec3{0.0, 0.0, 1.0}
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
		if b.Ik == nil || math.Abs(b.Ik.UnitRotation.GetDegrees().GetX()-expectedIkLimitedDegree) > 1e-5 {
			t.Errorf("Expected IkLimitedRadian to be %v, got %v", expectedIkLimitedDegree, b.Ik.UnitRotation.GetDegrees().GetX())
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
			expectedMinAngleLimit := &mmath.MVec3{-180.0, 0.0, 0.0}
			if !il.MinAngleLimit.GetDegrees().NearEquals(expectedMinAngleLimit, 1e-5) {
				t.Errorf("Expected MinAngleLimit to be %v, got %v", expectedMinAngleLimit, il.MinAngleLimit.GetDegrees())
			}
			expectedMaxAngleLimit := &mmath.MVec3{-0.5, 0.0, 0.0}
			if !il.MaxAngleLimit.GetDegrees().NearEquals(expectedMaxAngleLimit, 1e-5) {
				t.Errorf("Expected MaxAngleLimit to be %v, got %v", expectedMaxAngleLimit, il.MaxAngleLimit.GetDegrees())
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
		m := model.Morphs.Get(2)
		expectedName := "にこり"
		if m.Name != expectedName {
			t.Errorf("Expected Name to be %q, got %q", expectedName, m.Name)
		}
		expectedEnglishName := "Fcl_BRW_Fun"
		if m.EnglishName != expectedEnglishName {
			t.Errorf("Expected EnglishName to be %q, got %q", expectedEnglishName, m.EnglishName)
		}
		expectedMorphType := pmx.MORPH_TYPE_VERTEX
		if m.MorphType != expectedMorphType {
			t.Errorf("Expected MorphType to be %v, got %v", expectedMorphType, m.MorphType)
		}
		expectedOffsetCount := 68
		if len(m.Offsets) != expectedOffsetCount {
			t.Errorf("Expected OffsetCount to be %v, got %v", expectedOffsetCount, len(m.Offsets))
		}
		{
			o := m.Offsets[30].(*pmx.VertexMorphOffset)
			expectedVertexIndex := 19329
			if o.VertexIndex != expectedVertexIndex {
				t.Errorf("Expected VertexIndex to be %v, got %v", expectedVertexIndex, o.VertexIndex)
			}
			expectedPosition := &mmath.MVec3{-0.01209146, 0.1320038, -0.0121327}
			if !o.Position.MMD().NearEquals(expectedPosition, 1e-5) {
				t.Errorf("Expected Position to be %v, got %v", expectedPosition, o.Position)
			}
		}
	}

	{
		m := model.Morphs.Get(111)
		expectedName := "いボーン"
		if m.Name != expectedName {
			t.Errorf("Expected Name to be %q, got %q", expectedName, m.Name)
		}
		expectedEnglishName := "Fcl_MTH_I_Bone"
		if m.EnglishName != expectedEnglishName {
			t.Errorf("Expected EnglishName to be %q, got %q", expectedEnglishName, m.EnglishName)
		}
		expectedMorphType := pmx.MORPH_TYPE_BONE
		if m.MorphType != expectedMorphType {
			t.Errorf("Expected MorphType to be %v, got %v", expectedMorphType, m.MorphType)
		}
		expectedOffsetCount := 3
		if len(m.Offsets) != expectedOffsetCount {
			t.Errorf("Expected OffsetCount to be %v, got %v", expectedOffsetCount, len(m.Offsets))
		}
		{
			o := m.Offsets[1].(*pmx.BoneMorphOffset)
			expectedBoneIndex := 17
			if o.BoneIndex != expectedBoneIndex {
				t.Errorf("Expected BoneIndex to be %v, got %v", expectedBoneIndex, o.BoneIndex)
			}
			expectedPosition := &mmath.MVec3{0.0, 0.0, 0.0}
			if !o.Position.MMD().NearEquals(expectedPosition, 1e-5) {
				t.Errorf("Expected BoneIndex to be %v, got %v", expectedBoneIndex, o.BoneIndex)
			}
			expectedRotation := &mmath.MVec3{-6.000048, 3.952523e-14, -3.308919e-14}
			if !o.Rotation.GetDegreesMMD().NearEquals(expectedRotation, 1e-5) {
				t.Errorf("Expected Rotation to be %v, got %v", expectedRotation, o.Rotation.GetDegrees())
			}
		}
	}

	{
		m := model.Morphs.Get(122)
		expectedName := "なごみ材質"
		if m.Name != expectedName {
			t.Errorf("Expected Name to be %q, got %q", expectedName, m.Name)
		}
		expectedEnglishName := "eye_Nagomi_Material"
		if m.EnglishName != expectedEnglishName {
			t.Errorf("Expected EnglishName to be %q, got %q", expectedEnglishName, m.EnglishName)
		}
		expectedMorphType := pmx.MORPH_TYPE_MATERIAL
		if m.MorphType != expectedMorphType {
			t.Errorf("Expected MorphType to be %v, got %v", expectedMorphType, m.MorphType)
		}
		expectedOffsetCount := 4
		if len(m.Offsets) != expectedOffsetCount {
			t.Errorf("Expected OffsetCount to be %v, got %v", expectedOffsetCount, len(m.Offsets))
		}
		{
			o := m.Offsets[3].(*pmx.MaterialMorphOffset)
			expectedMaterialIndex := 12
			if o.MaterialIndex != expectedMaterialIndex {
				t.Errorf("Expected MaterialIndex to be %v, got %v", expectedMaterialIndex, o.MaterialIndex)
			}
			expectedCalcMode := pmx.CALC_MODE_ADDITION
			if o.CalcMode != expectedCalcMode {
				t.Errorf("Expected CalcMode to be %v, got %v", expectedCalcMode, o.CalcMode)
			}
			expectedDiffuse := &mmath.MVec4{0.0, 0.0, 0.0, 1.0}
			if !o.Diffuse.NearEquals(expectedDiffuse, 1e-5) {
				t.Errorf("Expected Diffuse to be %v, got %v", expectedDiffuse, o.Diffuse)
			}
			expectedSpecular := &mmath.MVec4{0.0, 0.0, 0.0, 0.0}
			if !o.Specular.NearEquals(expectedSpecular, 1e-5) {
				t.Errorf("Expected Specular to be %v, got %v", expectedSpecular, o.Specular)
			}
			expectedAmbient := &mmath.MVec3{0.0, 0.0, 0.0}
			if !o.Ambient.NearEquals(expectedAmbient, 1e-5) {
				t.Errorf("Expected Ambient to be %v, got %v", expectedAmbient, o.Ambient)
			}
			expectedEdge := &mmath.MVec4{0.0, 0.0, 0.0, 0.0}
			if !o.Edge.NearEquals(expectedEdge, 1e-5) {
				t.Errorf("Expected Edge to be %v, got %v", expectedEdge, o.Edge)
			}
			expectedEdgeSize := 0.0
			if math.Abs(o.EdgeSize-expectedEdgeSize) > 1e-5 {
				t.Errorf("Expected EdgeSize to be %v, got %v", expectedEdgeSize, o.EdgeSize)
			}
			expectedTextureFactor := &mmath.MVec4{0.0, 0.0, 0.0, 0.0}
			if !o.TextureFactor.NearEquals(expectedTextureFactor, 1e-5) {
				t.Errorf("Expected TextureFactor to be %v, got %v", expectedTextureFactor, o.TextureFactor)
			}
			expectedSphereTextureFactor := &mmath.MVec4{0.0, 0.0, 0.0, 0.0}
			if !o.SphereTextureFactor.NearEquals(expectedSphereTextureFactor, 1e-5) {
				t.Errorf("Expected SphereTextureFactor to be %v, got %v", expectedSphereTextureFactor, o.SphereTextureFactor)
			}
			expectedToonTextureFactor := &mmath.MVec4{0.0, 0.0, 0.0, 0.0}
			if !o.ToonTextureFactor.NearEquals(expectedToonTextureFactor, 1e-5) {
				t.Errorf("Expected ToonTextureFactor to be %v, got %v", expectedToonTextureFactor, o.ToonTextureFactor)
			}
		}
	}

	{
		m := model.Morphs.Get(137)
		expectedName := "ひそめ"
		if m.Name != expectedName {
			t.Errorf("Expected Name to be %q, got %q", expectedName, m.Name)
		}
		expectedEnglishName := "brow_Frown"
		if m.EnglishName != expectedEnglishName {
			t.Errorf("Expected EnglishName to be %q, got %q", expectedEnglishName, m.EnglishName)
		}
		expectedMorphType := pmx.MORPH_TYPE_GROUP
		if m.MorphType != expectedMorphType {
			t.Errorf("Expected MorphType to be %v, got %v", expectedMorphType, m.MorphType)
		}
		expectedOffsetCount := 6
		if len(m.Offsets) != expectedOffsetCount {
			t.Errorf("Expected OffsetCount to be %v, got %v", expectedOffsetCount, len(m.Offsets))
		}
		{
			o := m.Offsets[2].(*pmx.GroupMorphOffset)
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
		d := model.DisplaySlots.Get(0)
		expectedName := "Root"
		if d.Name != expectedName {
			t.Errorf("Expected Name to be %q, got %q", expectedName, d.Name)
		}
		expectedEnglishName := "Root"
		if d.EnglishName != expectedEnglishName {
			t.Errorf("Expected EnglishName to be %q, got %q", expectedEnglishName, d.EnglishName)
		}
		expectedSpecialFlag := pmx.SPECIAL_FLAG_ON
		if d.SpecialFlag != expectedSpecialFlag {
			t.Errorf("Expected SpecialFlag to be %v, got %v", expectedSpecialFlag, d.SpecialFlag)
		}
		expectedReferenceCount := 1
		if len(d.References) != expectedReferenceCount {
			t.Errorf("Expected ReferenceCount to be %v, got %v", expectedReferenceCount, len(d.References))
		}
		{
			r := d.References[0]
			expectedDisplayType := pmx.DISPLAY_TYPE_BONE
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
		d := model.DisplaySlots.Get(1)
		expectedName := "表情"
		if d.Name != expectedName {
			t.Errorf("Expected Name to be %q, got %q", expectedName, d.Name)
		}
		expectedEnglishName := "Exp"
		if d.EnglishName != expectedEnglishName {
			t.Errorf("Expected EnglishName to be %q, got %q", expectedEnglishName, d.EnglishName)
		}
		expectedSpecialFlag := pmx.SPECIAL_FLAG_ON
		if d.SpecialFlag != expectedSpecialFlag {
			t.Errorf("Expected SpecialFlag to be %v, got %v", expectedSpecialFlag, d.SpecialFlag)
		}
		expectedReferenceCount := 143
		if len(d.References) != expectedReferenceCount {
			t.Errorf("Expected ReferenceCount to be %v, got %v", expectedReferenceCount, len(d.References))
		}
		{
			r := d.References[50]
			expectedDisplayType := pmx.DISPLAY_TYPE_MORPH
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
		d := model.DisplaySlots.Get(9)
		expectedName := "右指"
		if d.Name != expectedName {
			t.Errorf("Expected Name to be %q, got %q", expectedName, d.Name)
		}
		expectedEnglishName := "right hand"
		if d.EnglishName != expectedEnglishName {
			t.Errorf("Expected EnglishName to be %q, got %q", expectedEnglishName, d.EnglishName)
		}
		expectedSpecialFlag := pmx.SPECIAL_FLAG_OFF
		if d.SpecialFlag != expectedSpecialFlag {
			t.Errorf("Expected SpecialFlag to be %v, got %v", expectedSpecialFlag, d.SpecialFlag)
		}
		expectedReferenceCount := 15
		if len(d.References) != expectedReferenceCount {
			t.Errorf("Expected ReferenceCount to be %v, got %v", expectedReferenceCount, len(d.References))
		}
		{
			r := d.References[7]
			expectedDisplayType := pmx.DISPLAY_TYPE_BONE
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
		b := model.RigidBodies.Get(14)
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
		expectedShapeType := pmx.SHAPE_CAPSULE
		if b.ShapeType != expectedShapeType {
			t.Errorf("Expected ShapeType to be %v, got %v", expectedShapeType, b.ShapeType)
		}
		expectedSize := &mmath.MVec3{0.5398922, 2.553789, 0.0}
		if !b.Size.NearEquals(expectedSize, 1e-5) {
			t.Errorf("Expected Size to be %v, got %v", expectedSize, b.Size)
		}
		expectedPosition := &mmath.MVec3{-2.52586, 15.45157, 0.241455}
		if !b.Position.MMD().NearEquals(expectedPosition, 1e-5) {
			t.Errorf("Expected Position to be %v, got %v", expectedPosition, b.Position)
		}
		expectedRotation := &mmath.MVec3{0.0, 0.0, 125.00}
		if !b.Rotation.GetDegrees().NearEquals(expectedRotation, 1e-5) {
			t.Errorf("Expected Rotation to be %v, got %v", expectedRotation, b.Rotation.GetDegrees())
		}
		expectedMass := 1.0
		if math.Abs(b.RigidBodyParam.Mass-expectedMass) > 1e-5 {
			t.Errorf("Expected Mass to be %v, got %v", expectedMass, b.RigidBodyParam.Mass)
		}
		expectedLinearDamping := 0.5
		if math.Abs(b.RigidBodyParam.LinearDamping-expectedLinearDamping) > 1e-5 {
			t.Errorf("Expected LinearDamping to be %v, got %v", expectedLinearDamping, b.RigidBodyParam.LinearDamping)
		}
		expectedAngularDamping := 0.5
		if math.Abs(b.RigidBodyParam.AngularDamping-expectedAngularDamping) > 1e-5 {
			t.Errorf("Expected AngularDamping to be %v, got %v", expectedAngularDamping, b.RigidBodyParam.AngularDamping)
		}
		expectedRestitution := 0.0
		if math.Abs(b.RigidBodyParam.Restitution-expectedRestitution) > 1e-5 {
			t.Errorf("Expected Restitution to be %v, got %v", expectedRestitution, b.RigidBodyParam.Restitution)
		}
		expectedFriction := 0.0
		if math.Abs(b.RigidBodyParam.Friction-expectedFriction) > 1e-5 {
			t.Errorf("Expected Friction to be %v, got %v", expectedFriction, b.RigidBodyParam.Friction)
		}
		expectedPhysicsType := pmx.PHYSICS_TYPE_STATIC
		if b.PhysicsType != expectedPhysicsType {
			t.Errorf("Expected PhysicsType to be %v, got %v", expectedPhysicsType, b.PhysicsType)
		}
	}

	{
		j := model.Joints.Get(13)
		expectedName := "↓|頭|髪_06-01"
		if j.Name != expectedName {
			t.Errorf("Expected Name to be %q, got %q", expectedName, j.Name)
		}
		expectedEnglishName := "↓|頭|髪_06-01"
		if j.EnglishName != expectedEnglishName {
			t.Errorf("Expected EnglishName to be %q, got %q", expectedEnglishName, j.EnglishName)
		}
		expectedRigidBodyIndexA := 5
		if j.RigidbodyIndexA != expectedRigidBodyIndexA {
			t.Errorf("Expected RigidBodyIndexA to be %v, got %v", expectedRigidBodyIndexA, j.RigidbodyIndexA)
		}
		expectedRigidBodyIndexB := 45
		if j.RigidbodyIndexB != expectedRigidBodyIndexB {
			t.Errorf("Expected RigidBodyIndexB to be %v, got %v", expectedRigidBodyIndexB, j.RigidbodyIndexB)
		}
		expectedPosition := &mmath.MVec3{-1.189768, 18.56266, 0.07277258}
		if !j.Position.MMD().NearEquals(expectedPosition, 1e-5) {
			t.Errorf("Expected Position to be %v, got %v", expectedPosition, j.Position)
		}
		expectedRotation := &mmath.MVec3{-15.10554, 91.26718, -4.187446}
		if !j.Rotation.GetDegrees().NearEquals(expectedRotation, 1e-5) {
			t.Errorf("Expected Rotation to be %v, got %v", expectedRotation, j.Rotation.GetDegrees())
		}
		expectedTranslationLimitMin := &mmath.MVec3{0.0, 0.0, 0.0}
		if !j.JointParam.TranslationLimitMin.NearEquals(expectedTranslationLimitMin, 1e-5) {
			t.Errorf("Expected TranslationLimitation1 to be %v, got %v", expectedTranslationLimitMin, j.JointParam.TranslationLimitMin)
		}
		expectedTranslationLimitMax := &mmath.MVec3{0.0, 0.0, 0.0}
		if !j.JointParam.TranslationLimitMax.NearEquals(expectedTranslationLimitMax, 1e-5) {
			t.Errorf("Expected TranslationLimitation2 to be %v, got %v", expectedTranslationLimitMax, j.JointParam.TranslationLimitMax)
		}
		expectedRotationLimitMin := &mmath.MVec3{-29.04, -14.3587, -29.04}
		if !j.JointParam.RotationLimitMin.GetDegrees().NearEquals(expectedRotationLimitMin, 1e-5) {
			t.Errorf("Expected RotationLimitation1 to be %v, got %v", expectedRotationLimitMin, j.JointParam.RotationLimitMin.GetDegrees())
		}
		expectedRotationLimitMax := &mmath.MVec3{29.04, 14.3587, 29.04}
		if !j.JointParam.RotationLimitMax.GetDegrees().NearEquals(expectedRotationLimitMax, 1e-5) {
			t.Errorf("Expected RotationLimitation2 to be %v, got %v", expectedRotationLimitMax, j.JointParam.RotationLimitMax.GetDegrees())
		}
		expectedSpringConstantTranslation := &mmath.MVec3{0.0, 0.0, 0.0}
		if !j.JointParam.SpringConstantTranslation.NearEquals(expectedSpringConstantTranslation, 1e-5) {
			t.Errorf("Expected SpringConstantTranslation to be %v, got %v", expectedSpringConstantTranslation, j.JointParam.SpringConstantTranslation)
		}
		expectedSpringConstantRotation := &mmath.MVec3{60.0, 29.6667, 60.0}
		if !j.JointParam.SpringConstantRotation.NearEquals(expectedSpringConstantRotation, 1e-5) {
			t.Errorf("Expected SpringConstantRotation to be %v, got %v", expectedSpringConstantRotation, j.JointParam.SpringConstantRotation)
		}
	}
}
