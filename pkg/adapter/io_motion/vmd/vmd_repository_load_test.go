// 指示: miu200521358
package vmd

import (
	"path/filepath"
	"testing"

	"github.com/miu200521358/mlib_go/pkg/domain/mmath"
	"github.com/miu200521358/mlib_go/pkg/domain/model"
	"github.com/miu200521358/mlib_go/pkg/domain/motion"
	"gonum.org/v1/gonum/spatial/r3"
)

func TestVmdRepository_Load(t *testing.T) {
	r := NewVmdRepository()

	// Test case 2: File not found
	invalidPath := testResourcePath("nonexistent.vmd")
	_, err := r.Load(invalidPath)

	if err == nil {
		t.Errorf("Expected error to be not nil, got nil")
	}

	// Test case 1: Successful read
	path := testResourcePath("サンプルモーション.vmd")
	loaded, err := r.Load(path)

	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err)
	}

	// Verify the model properties
	expectedPath := path
	if loaded.Path() != expectedPath {
		t.Errorf("Expected Path to be %q, got %q", expectedPath, loaded.Path())
	}

	// モデル名
	expectedModelName := "日本 roco式 トレス用"
	if loaded.Name() != expectedModelName {
		t.Errorf("Expected modelName to be %q, got %q", expectedModelName, loaded.Name())
	}

	motionData, ok := loaded.(*motion.VmdMotion)
	if !ok {
		t.Fatalf("Expected motion type to be *VmdMotion, got %T", loaded)
	}

	// キーフレがある
	{
		bf := motionData.BoneFrames.Get(model.CENTER.String()).Get(358)

		// フレーム番号
		expectedFrameNo := motion.Frame(358)
		if bf.Index() != expectedFrameNo {
			t.Errorf("Expected FrameNo to be %.4f, got %.4f", expectedFrameNo, bf.Index())
		}

		// 位置
		expectedPosition := vec3(1.094920158, 0, 0.100637913)
		if bf.Position == nil {
			t.Fatalf("Expected Position to be not nil")
		}
		if !bf.Position.MMD().NearEquals(expectedPosition, 1e-8) {
			t.Errorf("Expected Position to be %v, got %v", expectedPosition, bf.Position.MMD())
		}

		// 回転
		expectedRotation := mmath.NewQuaternionByValues(0, 0, 0, 1)
		if bf.Rotation == nil {
			t.Fatalf("Expected Rotation to be not nil")
		}
		if 1-bf.Rotation.Dot(expectedRotation) > 1e-8 {
			t.Errorf("Expected Rotation to be %v, got %v", expectedRotation, bf.Rotation)
		}

		// 補間曲線
		expectedTranslateXStart := mmath.Vec2{X: 64, Y: 0}
		if bf.Curves == nil {
			t.Fatalf("Expected Curves to be not nil")
		}
		if !bf.Curves.TranslateX.Start.NearEquals(expectedTranslateXStart, 1e-5) {
			t.Errorf("Expected TranslateX.Start to be %v, got %v", expectedTranslateXStart, bf.Curves.TranslateX.Start)
		}

		expectedTranslateXEnd := mmath.Vec2{X: 87, Y: 87}
		if !bf.Curves.TranslateX.End.NearEquals(expectedTranslateXEnd, 1e-5) {
			t.Errorf("Expected TranslateX.End to be %v, got %v", expectedTranslateXEnd, bf.Curves.TranslateX.End)
		}

		expectedTranslateYStart := mmath.Vec2{X: 20, Y: 20}
		if !bf.Curves.TranslateY.Start.NearEquals(expectedTranslateYStart, 1e-5) {
			t.Errorf("Expected TranslateY.Start to be %v, got %v", expectedTranslateYStart, bf.Curves.TranslateY.Start)
		}

		expectedTranslateYEnd := mmath.Vec2{X: 107, Y: 107}
		if !bf.Curves.TranslateY.End.NearEquals(expectedTranslateYEnd, 1e-5) {
			t.Errorf("Expected TranslateY.End to be %v, got %v", expectedTranslateYEnd, bf.Curves.TranslateY.End)
		}

		expectedTranslateZStart := mmath.Vec2{X: 64, Y: 0}
		if !bf.Curves.TranslateZ.Start.NearEquals(expectedTranslateZStart, 1e-5) {
			t.Errorf("Expected TranslateZ.Start to be %v, got %v", expectedTranslateZStart, bf.Curves.TranslateZ.Start)
		}

		expectedTranslateZEnd := mmath.Vec2{X: 87, Y: 87}
		if !bf.Curves.TranslateZ.End.NearEquals(expectedTranslateZEnd, 1e-5) {
			t.Errorf("Expected TranslateZ.End to be %v, got %v", expectedTranslateZEnd, bf.Curves.TranslateZ.End)
		}

		expectedRotateStart := mmath.Vec2{X: 20, Y: 20}
		if !bf.Curves.Rotate.Start.NearEquals(expectedRotateStart, 1e-5) {
			t.Errorf("Expected Rotate.Start to be %v, got %v", expectedRotateStart, bf.Curves.Rotate.Start)
		}

		expectedRotateEnd := mmath.Vec2{X: 107, Y: 107}
		if !bf.Curves.Rotate.End.NearEquals(expectedRotateEnd, 1e-5) {
			t.Errorf("Expected Rotate.End to be %v, got %v", expectedRotateEnd, bf.Curves.Rotate.End)
		}
	}

	{
		bf := motionData.BoneFrames.Get(model.UPPER.String()).Get(689)

		// フレーム番号
		expectedFrameNo := motion.Frame(689)
		if bf.Index() != expectedFrameNo {
			t.Errorf("Expected FrameNo to be %.4f, got %.4f", expectedFrameNo, bf.Index())
		}

		// 位置
		expectedPosition := vec3(0, 0, 0)
		if bf.Position == nil {
			t.Fatalf("Expected Position to be not nil")
		}
		if !bf.Position.MMD().NearEquals(expectedPosition, 1e-8) {
			t.Errorf("Expected Position to be %v, got %v", expectedPosition, bf.Position.MMD())
		}

		// 回転
		expectedDegrees := vec3(-6.270921156, -26.96361355, 0.63172903)
		if bf.Rotation == nil {
			t.Fatalf("Expected Rotation to be not nil")
		}
		if bf.Rotation.ToMMDDegrees().NearEquals(expectedDegrees, 1e-8) {
			t.Errorf("Expected Rotation to be %v, got %v", expectedDegrees, bf.Rotation)
		}

		// 補間曲線
		expectedTranslateXStart := mmath.Vec2{X: 20, Y: 20}
		if bf.Curves == nil {
			t.Fatalf("Expected Curves to be not nil")
		}
		if !bf.Curves.TranslateX.Start.NearEquals(expectedTranslateXStart, 1e-5) {
			t.Errorf("Expected TranslateX.Start to be %v, got %v", expectedTranslateXStart, bf.Curves.TranslateX.Start)
		}

		expectedTranslateXEnd := mmath.Vec2{X: 107, Y: 107}
		if !bf.Curves.TranslateX.End.NearEquals(expectedTranslateXEnd, 1e-5) {
			t.Errorf("Expected TranslateX.End to be %v, got %v", expectedTranslateXEnd, bf.Curves.TranslateX.End)
		}

		expectedTranslateYStart := mmath.Vec2{X: 20, Y: 20}
		if !bf.Curves.TranslateY.Start.NearEquals(expectedTranslateYStart, 1e-5) {
			t.Errorf("Expected TranslateY.Start to be %v, got %v", expectedTranslateYStart, bf.Curves.TranslateY.Start)
		}

		expectedTranslateYEnd := mmath.Vec2{X: 107, Y: 107}
		if !bf.Curves.TranslateY.End.NearEquals(expectedTranslateYEnd, 1e-5) {
			t.Errorf("Expected TranslateY.End to be %v, got %v", expectedTranslateYEnd, bf.Curves.TranslateY.End)
		}

		expectedTranslateZStart := mmath.Vec2{X: 20, Y: 20}
		if !bf.Curves.TranslateZ.Start.NearEquals(expectedTranslateZStart, 1e-5) {
			t.Errorf("Expected TranslateZ.Start to be %v, got %v", expectedTranslateZStart, bf.Curves.TranslateZ.Start)
		}

		expectedTranslateZEnd := mmath.Vec2{X: 107, Y: 107}
		if !bf.Curves.TranslateZ.End.NearEquals(expectedTranslateZEnd, 1e-5) {
			t.Errorf("Expected TranslateZ.End to be %v, got %v", expectedTranslateZEnd, bf.Curves.TranslateZ.End)
		}

		expectedRotateStart := mmath.Vec2{X: 20, Y: 20}
		if !bf.Curves.Rotate.Start.NearEquals(expectedRotateStart, 1e-5) {
			t.Errorf("Expected Rotate.Start to be %v, got %v", expectedRotateStart, bf.Curves.Rotate.Start)
		}

		expectedRotateEnd := mmath.Vec2{X: 107, Y: 107}
		if !bf.Curves.Rotate.End.NearEquals(expectedRotateEnd, 1e-5) {
			t.Errorf("Expected Rotate.End to be %v, got %v", expectedRotateEnd, bf.Curves.Rotate.End)
		}
	}

	{
		bf := motionData.BoneFrames.Get(model.LEG_IK.Right()).Get(384)

		// フレーム番号
		expectedFrameNo := motion.Frame(384)
		if bf.Index() != expectedFrameNo {
			t.Errorf("Expected FrameNo to be %.4f, got %.4f", expectedFrameNo, bf.Index())
		}

		// 位置
		expectedPosition := vec3(0.548680067, 0.134522215, -2.504074097)
		if bf.Position == nil {
			t.Fatalf("Expected Position to be not nil")
		}
		if !bf.Position.MMD().NearEquals(expectedPosition, 1e-8) {
			t.Errorf("Expected Position to be %v, got %v", expectedPosition, bf.Position.MMD())
		}

		// 回転
		expectedDegrees := vec3(22.20309405, 6.80959631, 2.583712695)
		if bf.Rotation == nil {
			t.Fatalf("Expected Rotation to be not nil")
		}
		if bf.Rotation.ToMMDDegrees().NearEquals(expectedDegrees, 1e-8) {
			t.Errorf("Expected Rotation to be %v, got %v", expectedDegrees, bf.Rotation)
		}

		// 補間曲線
		expectedTranslateXStart := mmath.Vec2{X: 64, Y: 0}
		if bf.Curves == nil {
			t.Fatalf("Expected Curves to be not nil")
		}
		if !bf.Curves.TranslateX.Start.NearEquals(expectedTranslateXStart, 1e-5) {
			t.Errorf("Expected TranslateX.Start to be %v, got %v", expectedTranslateXStart, bf.Curves.TranslateX.Start)
		}

		expectedTranslateXEnd := mmath.Vec2{X: 64, Y: 127}
		if !bf.Curves.TranslateX.End.NearEquals(expectedTranslateXEnd, 1e-5) {
			t.Errorf("Expected TranslateX.End to be %v, got %v", expectedTranslateXEnd, bf.Curves.TranslateX.End)
		}

		expectedTranslateYStart := mmath.Vec2{X: 64, Y: 0}
		if !bf.Curves.TranslateY.Start.NearEquals(expectedTranslateYStart, 1e-5) {
			t.Errorf("Expected TranslateY.Start to be %v, got %v", expectedTranslateYStart, bf.Curves.TranslateY.Start)
		}

		expectedTranslateYEnd := mmath.Vec2{X: 87, Y: 87}
		if !bf.Curves.TranslateY.End.NearEquals(expectedTranslateYEnd, 1e-5) {
			t.Errorf("Expected TranslateY.End to be %v, got %v", expectedTranslateYEnd, bf.Curves.TranslateY.End)
		}

		expectedTranslateZStart := mmath.Vec2{X: 64, Y: 0}
		if !bf.Curves.TranslateZ.Start.NearEquals(expectedTranslateZStart, 1e-5) {
			t.Errorf("Expected TranslateZ.Start to be %v, got %v", expectedTranslateZStart, bf.Curves.TranslateZ.Start)
		}

		expectedTranslateZEnd := mmath.Vec2{X: 64, Y: 127}
		if !bf.Curves.TranslateZ.End.NearEquals(expectedTranslateZEnd, 1e-5) {
			t.Errorf("Expected TranslateZ.End to be %v, got %v", expectedTranslateZEnd, bf.Curves.TranslateZ.End)
		}

		expectedRotateStart := mmath.Vec2{X: 64, Y: 0}
		if !bf.Curves.Rotate.Start.NearEquals(expectedRotateStart, 1e-5) {
			t.Errorf("Expected Rotate.Start to be %v, got %v", expectedRotateStart, bf.Curves.Rotate.Start)
		}

		expectedRotateEnd := mmath.Vec2{X: 87, Y: 87}
		if !bf.Curves.Rotate.End.NearEquals(expectedRotateEnd, 1e-5) {
			t.Errorf("Expected Rotate.End to be %v, got %v", expectedRotateEnd, bf.Curves.Rotate.End)
		}
	}

	{
		// キーがないフレーム
		bf := motionData.BoneFrames.Get(model.LEG_IK.Left()).Get(384)

		// フレーム番号
		expectedFrameNo := motion.Frame(384)
		if bf.Index() != expectedFrameNo {
			t.Errorf("Expected FrameNo to be %.4f, got %.4f", expectedFrameNo, bf.Index())
		}

		// 位置
		expectedPosition := vec3(-1.63, 0.05, 2.58)
		if bf.Position == nil {
			t.Fatalf("Expected Position to be not nil")
		}
		if !bf.Position.MMD().NearEquals(expectedPosition, 1e-2) {
			t.Errorf("Expected Position to be %v, got %v", expectedPosition, bf.Position.MMD())
		}

		// 回転
		expectedDegrees := vec3(-1.4, 6.7, -5.2)
		if bf.Rotation == nil {
			t.Fatalf("Expected Rotation to be not nil")
		}
		if bf.Rotation.ToMMDDegrees().NearEquals(expectedDegrees, 1e-2) {
			t.Errorf("Expected Rotation to be %v, got %v", expectedDegrees, bf.Rotation)
		}
	}

	{
		// キーがないフレーム
		bf := motionData.BoneFrames.Get(model.LEG_IK.Left()).Get(394)

		// フレーム番号
		expectedFrameNo := motion.Frame(394)
		if bf.Index() != expectedFrameNo {
			t.Errorf("Expected FrameNo to be %.4f, got %.4f", expectedFrameNo, bf.Index())
		}

		// 位置
		expectedPosition := vec3(0.76, 1.17, 1.34)
		if bf.Position == nil {
			t.Fatalf("Expected Position to be not nil")
		}
		if !bf.Position.MMD().NearEquals(expectedPosition, 1e-2) {
			t.Errorf("Expected Position to be %v, got %v", expectedPosition, bf.Position.MMD())
		}

		// 回転
		expectedDegrees := vec3(-41.9, -1.6, 1.0)
		if bf.Rotation == nil {
			t.Fatalf("Expected Rotation to be not nil")
		}
		if bf.Rotation.ToMMDDegrees().NearEquals(expectedDegrees, 1e-2) {
			t.Errorf("Expected Rotation to be %v, got %v", expectedDegrees, bf.Rotation)
		}
	}

	{
		// キーがないフレーム
		bf := motionData.BoneFrames.Get(model.LEG_IK.Left()).Get(412)

		// フレーム番号
		expectedFrameNo := motion.Frame(412)
		if bf.Index() != expectedFrameNo {
			t.Errorf("Expected FrameNo to be %.4f, got %.4f", expectedFrameNo, bf.Index())
		}

		// 位置
		expectedPosition := vec3(-0.76, -0.61, -1.76)
		if bf.Position == nil {
			t.Fatalf("Expected Position to be not nil")
		}
		if !bf.Position.MMD().NearEquals(expectedPosition, 1e-2) {
			t.Errorf("Expected Position to be %v, got %v", expectedPosition, bf.Position.MMD())
		}

		// 回転
		expectedDegrees := vec3(43.1, 0.0, 0.0)
		if bf.Rotation == nil {
			t.Fatalf("Expected Rotation to be not nil")
		}
		if bf.Rotation.ToMMDDegrees().NearEquals(expectedDegrees, 1e-2) {
			t.Errorf("Expected Rotation to be %v, got %v", expectedDegrees, bf.Rotation)
		}
	}

	{
		// キーがないフレーム
		bf := motionData.BoneFrames.Get(model.ARM.Right()).Get(384)

		// フレーム番号
		expectedFrameNo := motion.Frame(384)
		if bf.Index() != expectedFrameNo {
			t.Errorf("Expected FrameNo to be %.4f, got %.4f", expectedFrameNo, bf.Index())
		}

		// 位置
		expectedPosition := vec3(0.0, 0.0, 0.0)
		if bf.Position == nil {
			t.Fatalf("Expected Position to be not nil")
		}
		if !bf.Position.MMD().NearEquals(expectedPosition, 1e-2) {
			t.Errorf("Expected Position to be %v, got %v", expectedPosition, bf.Position.MMD())
		}

		// 回転
		expectedDegrees := vec3(13.5, -4.3, 27.0)
		if bf.Rotation == nil {
			t.Fatalf("Expected Rotation to be not nil")
		}
		if bf.Rotation.ToMMDDegrees().NearEquals(expectedDegrees, 1e-2) {
			t.Errorf("Expected Rotation to be %v, got %v", expectedDegrees, bf.Rotation)
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
