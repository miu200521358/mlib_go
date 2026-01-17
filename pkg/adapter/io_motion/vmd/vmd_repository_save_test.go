// 指示: miu200521358
package vmd

import (
	"path/filepath"
	"testing"

	"github.com/miu200521358/mlib_go/pkg/adapter/io_common"
	"github.com/miu200521358/mlib_go/pkg/domain/mmath"
	"github.com/miu200521358/mlib_go/pkg/domain/motion"
)

func TestVmdRepository_Save1(t *testing.T) {
	path := filepath.Join(t.TempDir(), "test_output.vmd")

	motionData := motion.NewVmdMotion(path)
	motionData.SetName("Null_00")

	bf := motion.NewBoneFrame(motion.Frame(0))
	pos := vec3(1, 2, 3)
	bf.Position = &pos
	rot := mmath.NewQuaternionFromDegrees(10, 20, 30)
	bf.Rotation = &rot
	motionData.AppendBoneFrame(boneCenter, bf)

	r := NewVmdRepository()

	if err := r.Save("", motionData, io_common.SaveOptions{}); err != nil {
		t.Fatalf("Expected error to be nil, got %q", err)
	}

	reloadData, err := r.Load(path)
	if err != nil {
		t.Fatalf("Expected error to be nil, got %q", err)
	}
	reloadMotion, ok := reloadData.(*motion.VmdMotion)
	if !ok {
		t.Fatalf("Expected motion type to be *VmdMotion, got %T", reloadData)
	}

	if reloadMotion.Name() != motionData.Name() {
		t.Errorf("Expected model name to be '%s', got %q", motionData.Name(), reloadMotion.Name())
	}

	frames := reloadMotion.BoneFrames.Get(boneCenter)
	if !frames.Has(motion.Frame(0)) {
		t.Errorf("Expected %s to contain frame 0", boneCenter)
	}

	reloadBf := frames.Get(motion.Frame(0))
	if reloadBf.Position == nil {
		t.Fatalf("Expected position to be not nil")
	}
	if !reloadBf.Position.NearEquals(pos, 1e-8) {
		t.Errorf("Expected position to be %v, got %v", pos, reloadBf.Position.MMD())
	}
}

func TestVmdRepository_Save2(t *testing.T) {
	readPath := testResourcePath("サンプルモーション.vmd")

	r := NewVmdRepository()
	data, err := r.Load(readPath)
	if err != nil {
		t.Fatalf("Expected error to be nil, got %q", err)
	}

	outputPath := filepath.Join(t.TempDir(), "test_output.vmd")
	motionData, ok := data.(*motion.VmdMotion)
	if !ok {
		t.Fatalf("Expected motion type to be *VmdMotion, got %T", data)
	}

	if err := r.Save(outputPath, motionData, io_common.SaveOptions{}); err != nil {
		t.Fatalf("Expected error to be nil, got %q", err)
	}

	reloadData, err := r.Load(outputPath)
	if err != nil {
		t.Fatalf("Expected error to be nil, got %q", err)
	}
	reloadMotion, ok := reloadData.(*motion.VmdMotion)
	if !ok {
		t.Fatalf("Expected motion type to be *VmdMotion, got %T", reloadData)
	}

	if reloadMotion.Name() != motionData.Name() {
		t.Errorf("Expected model name to be '%s', got %q", motionData.Name(), reloadMotion.Name())
	}
}

func TestVmdRepository_Save3(t *testing.T) {
	readPath := testResourcePath("ドクヘビ_178cmカメラ.vmd")

	r := NewVmdRepository()
	data, err := r.Load(readPath)
	if err != nil {
		t.Fatalf("Expected error to be nil, got %q", err)
	}

	outputPath := filepath.Join(t.TempDir(), "test_output.vmd")
	motionData, ok := data.(*motion.VmdMotion)
	if !ok {
		t.Fatalf("Expected motion type to be *VmdMotion, got %T", data)
	}

	if err := r.Save(outputPath, motionData, io_common.SaveOptions{}); err != nil {
		t.Fatalf("Expected error to be nil, got %q", err)
	}

	reloadData, err := r.Load(outputPath)
	if err != nil {
		t.Fatalf("Expected error to be nil, got %q", err)
	}
	reloadMotion, ok := reloadData.(*motion.VmdMotion)
	if !ok {
		t.Fatalf("Expected motion type to be *VmdMotion, got %T", reloadData)
	}

	if reloadMotion.Name() != motionData.Name() {
		t.Errorf("Expected model name to be '%s', got %q", motionData.Name(), reloadMotion.Name())
	}
}

func TestVmdRepository_Save4(t *testing.T) {
	readPath := testResourcePath("モーフ_まばたき.vmd")

	r := NewVmdRepository()
	data, err := r.Load(readPath)
	if err != nil {
		t.Fatalf("Expected error to be nil, got %q", err)
	}

	outputPath := filepath.Join(t.TempDir(), "test_output.vmd")
	motionData, ok := data.(*motion.VmdMotion)
	if !ok {
		t.Fatalf("Expected motion type to be *VmdMotion, got %T", data)
	}

	if err := r.Save(outputPath, motionData, io_common.SaveOptions{}); err != nil {
		t.Fatalf("Expected error to be nil, got %q", err)
	}

	reloadData, err := r.Load(outputPath)
	if err != nil {
		t.Fatalf("Expected error to be nil, got %q", err)
	}
	reloadMotion, ok := reloadData.(*motion.VmdMotion)
	if !ok {
		t.Fatalf("Expected motion type to be *VmdMotion, got %T", reloadData)
	}

	if reloadMotion.Name() != motionData.Name() {
		t.Errorf("Expected model name to be '%s', got %q", motionData.Name(), reloadMotion.Name())
	}
}
