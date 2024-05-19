package vmd

import (
	"testing"

	"github.com/miu200521358/mlib_go/pkg/mmath"
	"github.com/miu200521358/mlib_go/pkg/pmx"
)

func TestVmdMotionReader_ReadNameByFilepath(t *testing.T) {
	r := &VmdMotionReader{}

	// Test case 1: Successful read
	path := "../../test_resources/サンプルモーション_0046.vmd"
	modelName, err := r.ReadNameByFilepath(path)

	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err)
	}

	expectedModelName := "初音ミク準標準"
	if modelName != expectedModelName {
		t.Errorf("Expected modelName to be %q, got %q", expectedModelName, modelName)
	}
}

func TestVmdMotionReader_ReadByFilepath(t *testing.T) {
	r := &VmdMotionReader{}

	// Test case 2: File not found
	invalidPath := "../../test_resources/nonexistent.vmd"
	_, err := r.ReadByFilepath(invalidPath)

	if err == nil {
		t.Errorf("Expected error to be not nil, got nil")
	}

	// Test case 1: Successful read
	path := "../../test_resources/サンプルモーション.vmd"
	model, err := r.ReadByFilepath(path)

	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err)
	}

	// Verify the model properties
	expectedPath := path
	if model.GetPath() != expectedPath {
		t.Errorf("Expected Path to be %q, got %q", expectedPath, model.GetPath())
	}

	// モデル名
	expectedModelName := "日本 roco式 トレス用"
	if model.GetName() != expectedModelName {
		t.Errorf("Expected modelName to be %q, got %q", expectedModelName, model.GetName())
	}

	motion := model.(*VmdMotion)

	// キーフレがある
	{
		bf := motion.BoneFrames.Get(pmx.CENTER.String()).Get(358)

		// フレーム番号
		expectedFrameNo := int(358)
		if bf.GetIndex() != expectedFrameNo {
			t.Errorf("Expected FrameNo to be %d, got %d", expectedFrameNo, bf.GetIndex())
		}

		// 位置
		expectedPosition := &mmath.MVec3{1.094920158, 0, 0.100637913}
		if !bf.Position.MMD().PracticallyEquals(expectedPosition, 1e-8) {
			t.Errorf("Expected Position to be %v, got %v", expectedPosition, bf.Position.MMD())
		}

		// 回転
		expectedRotation := mmath.NewMQuaternionByValues(0, 0, 0, 1)
		if 1-bf.Rotation.GetQuaternion().MMD().Dot(expectedRotation) > 1e-8 {
			t.Errorf("Expected Rotation to be %v, got %v", expectedRotation, bf.Rotation)
		}

		// 補間曲線
		expectedTranslateXStart := &mmath.MVec2{64, 0}
		if !bf.Curves.TranslateX.Start.PracticallyEquals(expectedTranslateXStart, 1e-5) {
			t.Errorf("Expected TranslateX.Start to be %v, got %v", expectedTranslateXStart, bf.Curves.TranslateX.Start)
		}

		expectedTranslateXEnd := &mmath.MVec2{87, 87}
		if !bf.Curves.TranslateX.End.PracticallyEquals(expectedTranslateXEnd, 1e-5) {
			t.Errorf("Expected TranslateX.End to be %v, got %v", expectedTranslateXEnd, bf.Curves.TranslateX.End)
		}

		expectedTranslateYStart := &mmath.MVec2{20, 20}
		if !bf.Curves.TranslateY.Start.PracticallyEquals(expectedTranslateYStart, 1e-5) {
			t.Errorf("Expected TranslateY.Start to be %v, got %v", expectedTranslateYStart, bf.Curves.TranslateY.Start)
		}

		expectedTranslateYEnd := &mmath.MVec2{107, 107}
		if !bf.Curves.TranslateY.End.PracticallyEquals(expectedTranslateYEnd, 1e-5) {
			t.Errorf("Expected TranslateY.End to be %v, got %v", expectedTranslateYEnd, bf.Curves.TranslateY.End)
		}

		expectedTranslateZStart := &mmath.MVec2{64, 0}
		if !bf.Curves.TranslateZ.Start.PracticallyEquals(expectedTranslateZStart, 1e-5) {
			t.Errorf("Expected TranslateZ.Start to be %v, got %v", expectedTranslateZStart, bf.Curves.TranslateZ.Start)
		}

		expectedTranslateZEnd := &mmath.MVec2{87, 87}
		if !bf.Curves.TranslateZ.End.PracticallyEquals(expectedTranslateZEnd, 1e-5) {
			t.Errorf("Expected TranslateZ.End to be %v, got %v", expectedTranslateZEnd, bf.Curves.TranslateZ.End)
		}

		expectedRotateStart := &mmath.MVec2{20, 20}
		if !bf.Curves.Rotate.Start.PracticallyEquals(expectedRotateStart, 1e-5) {
			t.Errorf("Expected Rotate.Start to be %v, got %v", expectedRotateStart, bf.Curves.Rotate.Start)
		}

		expectedRotateEnd := &mmath.MVec2{107, 107}
		if !bf.Curves.Rotate.End.PracticallyEquals(expectedRotateEnd, 1e-5) {
			t.Errorf("Expected Rotate.End to be %v, got %v", expectedRotateEnd, bf.Curves.Rotate.End)
		}
	}

	{
		bf := motion.BoneFrames.Get(pmx.UPPER.String()).Get(689)

		// フレーム番号
		expectedFrameNo := int(689)
		if bf.GetIndex() != expectedFrameNo {
			t.Errorf("Expected FrameNo to be %d, got %d", expectedFrameNo, bf.GetIndex())
		}

		// 位置
		expectedPosition := &mmath.MVec3{0, 0, 0}
		if !bf.Position.MMD().PracticallyEquals(expectedPosition, 1e-8) {
			t.Errorf("Expected Position to be %v, got %v", expectedPosition, bf.Position.MMD())
		}

		// 回転
		expectedDegrees := &mmath.MVec3{-6.270921156, -26.96361355, 0.63172903}
		if bf.Rotation.GetDegreesMMD().PracticallyEquals(expectedDegrees, 1e-8) {
			t.Errorf("Expected Rotation to be %v, got %v", expectedDegrees, bf.Rotation)
		}

		// 補間曲線
		expectedTranslateXStart := &mmath.MVec2{20, 20}
		if !bf.Curves.TranslateX.Start.PracticallyEquals(expectedTranslateXStart, 1e-5) {
			t.Errorf("Expected TranslateX.Start to be %v, got %v", expectedTranslateXStart, bf.Curves.TranslateX.Start)
		}

		expectedTranslateXEnd := &mmath.MVec2{107, 107}
		if !bf.Curves.TranslateX.End.PracticallyEquals(expectedTranslateXEnd, 1e-5) {
			t.Errorf("Expected TranslateX.End to be %v, got %v", expectedTranslateXEnd, bf.Curves.TranslateX.End)
		}

		expectedTranslateYStart := &mmath.MVec2{20, 20}
		if !bf.Curves.TranslateY.Start.PracticallyEquals(expectedTranslateYStart, 1e-5) {
			t.Errorf("Expected TranslateY.Start to be %v, got %v", expectedTranslateYStart, bf.Curves.TranslateY.Start)
		}

		expectedTranslateYEnd := &mmath.MVec2{107, 107}
		if !bf.Curves.TranslateY.End.PracticallyEquals(expectedTranslateYEnd, 1e-5) {
			t.Errorf("Expected TranslateY.End to be %v, got %v", expectedTranslateYEnd, bf.Curves.TranslateY.End)
		}

		expectedTranslateZStart := &mmath.MVec2{20, 20}
		if !bf.Curves.TranslateZ.Start.PracticallyEquals(expectedTranslateZStart, 1e-5) {
			t.Errorf("Expected TranslateZ.Start to be %v, got %v", expectedTranslateZStart, bf.Curves.TranslateZ.Start)
		}

		expectedTranslateZEnd := &mmath.MVec2{107, 107}
		if !bf.Curves.TranslateZ.End.PracticallyEquals(expectedTranslateZEnd, 1e-5) {
			t.Errorf("Expected TranslateZ.End to be %v, got %v", expectedTranslateZEnd, bf.Curves.TranslateZ.End)
		}

		expectedRotateStart := &mmath.MVec2{20, 20}
		if !bf.Curves.Rotate.Start.PracticallyEquals(expectedRotateStart, 1e-5) {
			t.Errorf("Expected Rotate.Start to be %v, got %v", expectedRotateStart, bf.Curves.Rotate.Start)
		}

		expectedRotateEnd := &mmath.MVec2{107, 107}
		if !bf.Curves.Rotate.End.PracticallyEquals(expectedRotateEnd, 1e-5) {
			t.Errorf("Expected Rotate.End to be %v, got %v", expectedRotateEnd, bf.Curves.Rotate.End)
		}
	}

	{
		bf := motion.BoneFrames.Get(pmx.LEG_IK.Right()).Get(384)

		// フレーム番号
		expectedFrameNo := int(384)
		if bf.GetIndex() != expectedFrameNo {
			t.Errorf("Expected FrameNo to be %d, got %d", expectedFrameNo, bf.GetIndex())
		}

		// 位置
		expectedPosition := &mmath.MVec3{0.548680067, 0.134522215, -2.504074097}
		if !bf.Position.MMD().PracticallyEquals(expectedPosition, 1e-8) {
			t.Errorf("Expected Position to be %v, got %v", expectedPosition, bf.Position.MMD())
		}

		// 回転
		expectedDegrees := &mmath.MVec3{22.20309405, 6.80959631, 2.583712695}
		if bf.Rotation.GetDegreesMMD().PracticallyEquals(expectedDegrees, 1e-8) {
			t.Errorf("Expected Rotation to be %v, got %v", expectedDegrees, bf.Rotation)
		}

		// 補間曲線
		expectedTranslateXStart := &mmath.MVec2{64, 0}
		if !bf.Curves.TranslateX.Start.PracticallyEquals(expectedTranslateXStart, 1e-5) {
			t.Errorf("Expected TranslateX.Start to be %v, got %v", expectedTranslateXStart, bf.Curves.TranslateX.Start)
		}

		expectedTranslateXEnd := &mmath.MVec2{64, 127}
		if !bf.Curves.TranslateX.End.PracticallyEquals(expectedTranslateXEnd, 1e-5) {
			t.Errorf("Expected TranslateX.End to be %v, got %v", expectedTranslateXEnd, bf.Curves.TranslateX.End)
		}

		expectedTranslateYStart := &mmath.MVec2{64, 0}
		if !bf.Curves.TranslateY.Start.PracticallyEquals(expectedTranslateYStart, 1e-5) {
			t.Errorf("Expected TranslateY.Start to be %v, got %v", expectedTranslateYStart, bf.Curves.TranslateY.Start)
		}

		expectedTranslateYEnd := &mmath.MVec2{87, 87}
		if !bf.Curves.TranslateY.End.PracticallyEquals(expectedTranslateYEnd, 1e-5) {
			t.Errorf("Expected TranslateY.End to be %v, got %v", expectedTranslateYEnd, bf.Curves.TranslateY.End)
		}

		expectedTranslateZStart := &mmath.MVec2{64, 0}
		if !bf.Curves.TranslateZ.Start.PracticallyEquals(expectedTranslateZStart, 1e-5) {
			t.Errorf("Expected TranslateZ.Start to be %v, got %v", expectedTranslateZStart, bf.Curves.TranslateZ.Start)
		}

		expectedTranslateZEnd := &mmath.MVec2{64, 127}
		if !bf.Curves.TranslateZ.End.PracticallyEquals(expectedTranslateZEnd, 1e-5) {
			t.Errorf("Expected TranslateZ.End to be %v, got %v", expectedTranslateZEnd, bf.Curves.TranslateZ.End)
		}

		expectedRotateStart := &mmath.MVec2{64, 0}
		if !bf.Curves.Rotate.Start.PracticallyEquals(expectedRotateStart, 1e-5) {
			t.Errorf("Expected Rotate.Start to be %v, got %v", expectedRotateStart, bf.Curves.Rotate.Start)
		}

		expectedRotateEnd := &mmath.MVec2{87, 87}
		if !bf.Curves.Rotate.End.PracticallyEquals(expectedRotateEnd, 1e-5) {
			t.Errorf("Expected Rotate.End to be %v, got %v", expectedRotateEnd, bf.Curves.Rotate.End)
		}
	}

	{
		// キーがないフレーム
		bf := motion.BoneFrames.Get(pmx.LEG_IK.Left()).Get(384)

		// フレーム番号
		expectedFrameNo := int(384)
		if bf.GetIndex() != expectedFrameNo {
			t.Errorf("Expected FrameNo to be %d, got %d", expectedFrameNo, bf.GetIndex())
		}

		// 位置
		expectedPosition := &mmath.MVec3{-1.63, 0.05, 2.58}
		if !bf.Position.MMD().PracticallyEquals(expectedPosition, 1e-2) {
			t.Errorf("Expected Position to be %v, got %v", expectedPosition, bf.Position.MMD())
		}

		// 回転
		expectedDegrees := &mmath.MVec3{-1.4, 6.7, -5.2}
		if bf.Rotation.GetDegreesMMD().PracticallyEquals(expectedDegrees, 1e-2) {
			t.Errorf("Expected Rotation to be %v, got %v", expectedDegrees, bf.Rotation)
		}
	}

	{
		// キーがないフレーム
		bf := motion.BoneFrames.Get(pmx.LEG_IK.Left()).Get(394)

		// フレーム番号
		expectedFrameNo := int(394)
		if bf.GetIndex() != expectedFrameNo {
			t.Errorf("Expected FrameNo to be %d, got %d", expectedFrameNo, bf.GetIndex())
		}

		// 位置
		expectedPosition := &mmath.MVec3{0.76, 1.17, 1.34}
		if !bf.Position.MMD().PracticallyEquals(expectedPosition, 1e-2) {
			t.Errorf("Expected Position to be %v, got %v", expectedPosition, bf.Position.MMD())
		}

		// 回転
		expectedDegrees := &mmath.MVec3{-41.9, -1.6, 1.0}
		if bf.Rotation.GetDegreesMMD().PracticallyEquals(expectedDegrees, 1e-2) {
			t.Errorf("Expected Rotation to be %v, got %v", expectedDegrees, bf.Rotation)
		}
	}

	{
		// キーがないフレーム
		bf := motion.BoneFrames.Get(pmx.LEG_IK.Left()).Get(412)

		// フレーム番号
		expectedFrameNo := int(412)
		if bf.GetIndex() != expectedFrameNo {
			t.Errorf("Expected FrameNo to be %d, got %d", expectedFrameNo, bf.GetIndex())
		}

		// 位置
		expectedPosition := &mmath.MVec3{-0.76, -0.61, -1.76}
		if !bf.Position.MMD().PracticallyEquals(expectedPosition, 1e-2) {
			t.Errorf("Expected Position to be %v, got %v", expectedPosition, bf.Position.MMD())
		}

		// 回転
		expectedDegrees := &mmath.MVec3{43.1, 0.0, 0.0}
		if bf.Rotation.GetDegreesMMD().PracticallyEquals(expectedDegrees, 1e-2) {
			t.Errorf("Expected Rotation to be %v, got %v", expectedDegrees, bf.Rotation)
		}
	}

	{
		// キーがないフレーム
		bf := motion.BoneFrames.Get(pmx.ARM.Right()).Get(384)

		// フレーム番号
		expectedFrameNo := int(384)
		if bf.GetIndex() != expectedFrameNo {
			t.Errorf("Expected FrameNo to be %d, got %d", expectedFrameNo, bf.GetIndex())
		}

		// 位置
		expectedPosition := &mmath.MVec3{0.0, 0.0, 0.0}
		if !bf.Position.MMD().PracticallyEquals(expectedPosition, 1e-2) {
			t.Errorf("Expected Position to be %v, got %v", expectedPosition, bf.Position.MMD())
		}

		// 回転
		expectedDegrees := &mmath.MVec3{13.5, -4.3, 27.0}
		if bf.Rotation.GetDegreesMMD().PracticallyEquals(expectedDegrees, 1e-2) {
			t.Errorf("Expected Rotation to be %v, got %v", expectedDegrees, bf.Rotation)
		}
	}
}
