package repository

import (
	"testing"

	"github.com/miu200521358/mlib_go/pkg/domain/mmath"
	"github.com/miu200521358/mlib_go/pkg/domain/pmx"
	"github.com/miu200521358/mlib_go/pkg/domain/vmd"
)

func TestVmdMotionReader_LoadName(t *testing.T) {
	r := NewVmdRepository(true)

	// Test case 1: Successful read
	path := "../../../test_resources/サンプルモーション_0046.vmd"
	modelName := r.LoadName(path)

	expectedModelName := "初音ミク準標準"
	if modelName != expectedModelName {
		t.Errorf("Expected modelName to be %q, got %q", expectedModelName, modelName)
	}
}

func TestVmdMotionReader_Load(t *testing.T) {
	r := NewVmdRepository(true)

	// Test case 2: File not found
	invalidPath := "../../../test_resources/nonexistent.vmd"
	_, err := r.Load(invalidPath)

	if err == nil {
		t.Errorf("Expected error to be not nil, got nil")
	}

	// Test case 1: Successful read
	path := "../../../test_resources/サンプルモーション.vmd"
	model, err := r.Load(path)

	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err)
	}

	// Verify the model properties
	expectedPath := path
	if model.Path() != expectedPath {
		t.Errorf("Expected Path to be %q, got %q", expectedPath, model.Path())
	}

	// モデル名
	expectedModelName := "日本 roco式 トレス用"
	if model.Name() != expectedModelName {
		t.Errorf("Expected modelName to be %q, got %q", expectedModelName, model.Name())
	}

	motion := model.(*vmd.VmdMotion)

	// キーフレがある
	{
		bf := motion.BoneFrames.Get(pmx.CENTER.String()).Get(358)

		// フレーム番号
		expectedFrameNo := float32(358)
		if bf.Index() != expectedFrameNo {
			t.Errorf("Expected FrameNo to be %.4f, got %.4f", expectedFrameNo, bf.Index())
		}

		// 位置
		expectedPosition := &mmath.MVec3{X: 1.094920158, Y: 0, Z: 0.100637913}
		if !bf.Position.MMD().NearEquals(expectedPosition, 1e-8) {
			t.Errorf("Expected Position to be %v, got %v", expectedPosition, bf.Position.MMD())
		}

		// 回転
		expectedRotation := mmath.NewMQuaternionByValues(0, 0, 0, 1)
		if 1-bf.Rotation.Dot(expectedRotation) > 1e-8 {
			t.Errorf("Expected Rotation to be %v, got %v", expectedRotation, bf.Rotation)
		}

		// 補間曲線
		expectedTranslateXStart := &mmath.MVec2{X: 64, Y: 0}
		if !bf.Curves.TranslateX.Start.NearEquals(expectedTranslateXStart, 1e-5) {
			t.Errorf("Expected TranslateX.Start to be %v, got %v", expectedTranslateXStart, bf.Curves.TranslateX.Start)
		}

		expectedTranslateXEnd := &mmath.MVec2{X: 87, Y: 87}
		if !bf.Curves.TranslateX.End.NearEquals(expectedTranslateXEnd, 1e-5) {
			t.Errorf("Expected TranslateX.End to be %v, got %v", expectedTranslateXEnd, bf.Curves.TranslateX.End)
		}

		expectedTranslateYStart := &mmath.MVec2{X: 20, Y: 20}
		if !bf.Curves.TranslateY.Start.NearEquals(expectedTranslateYStart, 1e-5) {
			t.Errorf("Expected TranslateY.Start to be %v, got %v", expectedTranslateYStart, bf.Curves.TranslateY.Start)
		}

		expectedTranslateYEnd := &mmath.MVec2{X: 107, Y: 107}
		if !bf.Curves.TranslateY.End.NearEquals(expectedTranslateYEnd, 1e-5) {
			t.Errorf("Expected TranslateY.End to be %v, got %v", expectedTranslateYEnd, bf.Curves.TranslateY.End)
		}

		expectedTranslateZStart := &mmath.MVec2{X: 64, Y: 0}
		if !bf.Curves.TranslateZ.Start.NearEquals(expectedTranslateZStart, 1e-5) {
			t.Errorf("Expected TranslateZ.Start to be %v, got %v", expectedTranslateZStart, bf.Curves.TranslateZ.Start)
		}

		expectedTranslateZEnd := &mmath.MVec2{X: 87, Y: 87}
		if !bf.Curves.TranslateZ.End.NearEquals(expectedTranslateZEnd, 1e-5) {
			t.Errorf("Expected TranslateZ.End to be %v, got %v", expectedTranslateZEnd, bf.Curves.TranslateZ.End)
		}

		expectedRotateStart := &mmath.MVec2{X: 20, Y: 20}
		if !bf.Curves.Rotate.Start.NearEquals(expectedRotateStart, 1e-5) {
			t.Errorf("Expected Rotate.Start to be %v, got %v", expectedRotateStart, bf.Curves.Rotate.Start)
		}

		expectedRotateEnd := &mmath.MVec2{X: 107, Y: 107}
		if !bf.Curves.Rotate.End.NearEquals(expectedRotateEnd, 1e-5) {
			t.Errorf("Expected Rotate.End to be %v, got %v", expectedRotateEnd, bf.Curves.Rotate.End)
		}
	}

	{
		bf := motion.BoneFrames.Get(pmx.UPPER.String()).Get(689)

		// フレーム番号
		expectedFrameNo := float32(689)
		if bf.Index() != expectedFrameNo {
			t.Errorf("Expected FrameNo to be %.4f, got %.4f", expectedFrameNo, bf.Index())
		}

		// 位置
		expectedPosition := &mmath.MVec3{X: 0, Y: 0, Z: 0}
		if !bf.Position.MMD().NearEquals(expectedPosition, 1e-8) {
			t.Errorf("Expected Position to be %v, got %v", expectedPosition, bf.Position.MMD())
		}

		// 回転
		expectedDegrees := &mmath.MVec3{X: -6.270921156, Y: -26.96361355, Z: 0.63172903}
		if bf.Rotation.ToMMDDegrees().NearEquals(expectedDegrees, 1e-8) {
			t.Errorf("Expected Rotation to be %v, got %v", expectedDegrees, bf.Rotation)
		}

		// 補間曲線
		expectedTranslateXStart := &mmath.MVec2{X: 20, Y: 20}
		if !bf.Curves.TranslateX.Start.NearEquals(expectedTranslateXStart, 1e-5) {
			t.Errorf("Expected TranslateX.Start to be %v, got %v", expectedTranslateXStart, bf.Curves.TranslateX.Start)
		}

		expectedTranslateXEnd := &mmath.MVec2{X: 107, Y: 107}
		if !bf.Curves.TranslateX.End.NearEquals(expectedTranslateXEnd, 1e-5) {
			t.Errorf("Expected TranslateX.End to be %v, got %v", expectedTranslateXEnd, bf.Curves.TranslateX.End)
		}

		expectedTranslateYStart := &mmath.MVec2{X: 20, Y: 20}
		if !bf.Curves.TranslateY.Start.NearEquals(expectedTranslateYStart, 1e-5) {
			t.Errorf("Expected TranslateY.Start to be %v, got %v", expectedTranslateYStart, bf.Curves.TranslateY.Start)
		}

		expectedTranslateYEnd := &mmath.MVec2{X: 107, Y: 107}
		if !bf.Curves.TranslateY.End.NearEquals(expectedTranslateYEnd, 1e-5) {
			t.Errorf("Expected TranslateY.End to be %v, got %v", expectedTranslateYEnd, bf.Curves.TranslateY.End)
		}

		expectedTranslateZStart := &mmath.MVec2{X: 20, Y: 20}
		if !bf.Curves.TranslateZ.Start.NearEquals(expectedTranslateZStart, 1e-5) {
			t.Errorf("Expected TranslateZ.Start to be %v, got %v", expectedTranslateZStart, bf.Curves.TranslateZ.Start)
		}

		expectedTranslateZEnd := &mmath.MVec2{X: 107, Y: 107}
		if !bf.Curves.TranslateZ.End.NearEquals(expectedTranslateZEnd, 1e-5) {
			t.Errorf("Expected TranslateZ.End to be %v, got %v", expectedTranslateZEnd, bf.Curves.TranslateZ.End)
		}

		expectedRotateStart := &mmath.MVec2{X: 20, Y: 20}
		if !bf.Curves.Rotate.Start.NearEquals(expectedRotateStart, 1e-5) {
			t.Errorf("Expected Rotate.Start to be %v, got %v", expectedRotateStart, bf.Curves.Rotate.Start)
		}

		expectedRotateEnd := &mmath.MVec2{X: 107, Y: 107}
		if !bf.Curves.Rotate.End.NearEquals(expectedRotateEnd, 1e-5) {
			t.Errorf("Expected Rotate.End to be %v, got %v", expectedRotateEnd, bf.Curves.Rotate.End)
		}
	}

	{
		bf := motion.BoneFrames.Get(pmx.LEG_IK.Right()).Get(384)

		// フレーム番号
		expectedFrameNo := float32(384)
		if bf.Index() != expectedFrameNo {
			t.Errorf("Expected FrameNo to be %.4f, got %.4f", expectedFrameNo, bf.Index())
		}

		// 位置
		expectedPosition := &mmath.MVec3{X: 0.548680067, Y: 0.134522215, Z: -2.504074097}
		if !bf.Position.MMD().NearEquals(expectedPosition, 1e-8) {
			t.Errorf("Expected Position to be %v, got %v", expectedPosition, bf.Position.MMD())
		}

		// 回転
		expectedDegrees := &mmath.MVec3{X: 22.20309405, Y: 6.80959631, Z: 2.583712695}
		if bf.Rotation.ToMMDDegrees().NearEquals(expectedDegrees, 1e-8) {
			t.Errorf("Expected Rotation to be %v, got %v", expectedDegrees, bf.Rotation)
		}

		// 補間曲線
		expectedTranslateXStart := &mmath.MVec2{X: 64, Y: 0}
		if !bf.Curves.TranslateX.Start.NearEquals(expectedTranslateXStart, 1e-5) {
			t.Errorf("Expected TranslateX.Start to be %v, got %v", expectedTranslateXStart, bf.Curves.TranslateX.Start)
		}

		expectedTranslateXEnd := &mmath.MVec2{X: 64, Y: 127}
		if !bf.Curves.TranslateX.End.NearEquals(expectedTranslateXEnd, 1e-5) {
			t.Errorf("Expected TranslateX.End to be %v, got %v", expectedTranslateXEnd, bf.Curves.TranslateX.End)
		}

		expectedTranslateYStart := &mmath.MVec2{X: 64, Y: 0}
		if !bf.Curves.TranslateY.Start.NearEquals(expectedTranslateYStart, 1e-5) {
			t.Errorf("Expected TranslateY.Start to be %v, got %v", expectedTranslateYStart, bf.Curves.TranslateY.Start)
		}

		expectedTranslateYEnd := &mmath.MVec2{X: 87, Y: 87}
		if !bf.Curves.TranslateY.End.NearEquals(expectedTranslateYEnd, 1e-5) {
			t.Errorf("Expected TranslateY.End to be %v, got %v", expectedTranslateYEnd, bf.Curves.TranslateY.End)
		}

		expectedTranslateZStart := &mmath.MVec2{X: 64, Y: 0}
		if !bf.Curves.TranslateZ.Start.NearEquals(expectedTranslateZStart, 1e-5) {
			t.Errorf("Expected TranslateZ.Start to be %v, got %v", expectedTranslateZStart, bf.Curves.TranslateZ.Start)
		}

		expectedTranslateZEnd := &mmath.MVec2{X: 64, Y: 127}
		if !bf.Curves.TranslateZ.End.NearEquals(expectedTranslateZEnd, 1e-5) {
			t.Errorf("Expected TranslateZ.End to be %v, got %v", expectedTranslateZEnd, bf.Curves.TranslateZ.End)
		}

		expectedRotateStart := &mmath.MVec2{X: 64, Y: 0}
		if !bf.Curves.Rotate.Start.NearEquals(expectedRotateStart, 1e-5) {
			t.Errorf("Expected Rotate.Start to be %v, got %v", expectedRotateStart, bf.Curves.Rotate.Start)
		}

		expectedRotateEnd := &mmath.MVec2{X: 87, Y: 87}
		if !bf.Curves.Rotate.End.NearEquals(expectedRotateEnd, 1e-5) {
			t.Errorf("Expected Rotate.End to be %v, got %v", expectedRotateEnd, bf.Curves.Rotate.End)
		}
	}

	{
		// キーがないフレーム
		bf := motion.BoneFrames.Get(pmx.LEG_IK.Left()).Get(384)

		// フレーム番号
		expectedFrameNo := float32(384)
		if bf.Index() != expectedFrameNo {
			t.Errorf("Expected FrameNo to be %.4f, got %.4f", expectedFrameNo, bf.Index())
		}

		// 位置
		expectedPosition := &mmath.MVec3{X: -1.63, Y: 0.05, Z: 2.58}
		if !bf.Position.MMD().NearEquals(expectedPosition, 1e-2) {
			t.Errorf("Expected Position to be %v, got %v", expectedPosition, bf.Position.MMD())
		}

		// 回転
		expectedDegrees := &mmath.MVec3{X: -1.4, Y: 6.7, Z: -5.2}
		if bf.Rotation.ToMMDDegrees().NearEquals(expectedDegrees, 1e-2) {
			t.Errorf("Expected Rotation to be %v, got %v", expectedDegrees, bf.Rotation)
		}
	}

	{
		// キーがないフレーム
		bf := motion.BoneFrames.Get(pmx.LEG_IK.Left()).Get(394)

		// フレーム番号
		expectedFrameNo := float32(394)
		if bf.Index() != expectedFrameNo {
			t.Errorf("Expected FrameNo to be %.4f, got %.4f", expectedFrameNo, bf.Index())
		}

		// 位置
		expectedPosition := &mmath.MVec3{X: 0.76, Y: 1.17, Z: 1.34}
		if !bf.Position.MMD().NearEquals(expectedPosition, 1e-2) {
			t.Errorf("Expected Position to be %v, got %v", expectedPosition, bf.Position.MMD())
		}

		// 回転
		expectedDegrees := &mmath.MVec3{X: -41.9, Y: -1.6, Z: 1.0}
		if bf.Rotation.ToMMDDegrees().NearEquals(expectedDegrees, 1e-2) {
			t.Errorf("Expected Rotation to be %v, got %v", expectedDegrees, bf.Rotation)
		}
	}

	{
		// キーがないフレーム
		bf := motion.BoneFrames.Get(pmx.LEG_IK.Left()).Get(412)

		// フレーム番号
		expectedFrameNo := float32(412)
		if bf.Index() != expectedFrameNo {
			t.Errorf("Expected FrameNo to be %.4f, got %.4f", expectedFrameNo, bf.Index())
		}

		// 位置
		expectedPosition := &mmath.MVec3{X: -0.76, Y: -0.61, Z: -1.76}
		if !bf.Position.MMD().NearEquals(expectedPosition, 1e-2) {
			t.Errorf("Expected Position to be %v, got %v", expectedPosition, bf.Position.MMD())
		}

		// 回転
		expectedDegrees := &mmath.MVec3{X: 43.1, Y: 0.0, Z: 0.0}
		if bf.Rotation.ToMMDDegrees().NearEquals(expectedDegrees, 1e-2) {
			t.Errorf("Expected Rotation to be %v, got %v", expectedDegrees, bf.Rotation)
		}
	}

	{
		// キーがないフレーム
		bf := motion.BoneFrames.Get(pmx.ARM.Right()).Get(384)

		// フレーム番号
		expectedFrameNo := float32(384)
		if bf.Index() != expectedFrameNo {
			t.Errorf("Expected FrameNo to be %.4f, got %.4f", expectedFrameNo, bf.Index())
		}

		// 位置
		expectedPosition := &mmath.MVec3{X: 0.0, Y: 0.0, Z: 0.0}
		if !bf.Position.MMD().NearEquals(expectedPosition, 1e-2) {
			t.Errorf("Expected Position to be %v, got %v", expectedPosition, bf.Position.MMD())
		}

		// 回転
		expectedDegrees := &mmath.MVec3{X: 13.5, Y: -4.3, Z: 27.0}
		if bf.Rotation.ToMMDDegrees().NearEquals(expectedDegrees, 1e-2) {
			t.Errorf("Expected Rotation to be %v, got %v", expectedDegrees, bf.Rotation)
		}
	}
}
