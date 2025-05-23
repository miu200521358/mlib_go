package repository

import (
	"testing"

	"github.com/miu200521358/mlib_go/pkg/domain/mmath"
	"github.com/miu200521358/mlib_go/pkg/domain/vmd"
)

func TestVmdWriter_Write1(t *testing.T) {
	path := "../../../test_resources/test_output.vmd"

	// Create a VmdMotion instance for testing
	motion := vmd.NewVmdMotion(path)
	motion.SetName("Null_00")

	bf := vmd.NewBoneFrame(0)
	bf.Position = &mmath.MVec3{X: 1, Y: 2, Z: 3}
	bf.Rotation = mmath.NewMQuaternionFromDegrees(10, 20, 30)
	motion.AppendBoneFrame("センター", bf)

	r := NewVmdRepository(true)

	// Create a VmdWriter instance
	err := r.Save("", motion, false)
	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err)
	}

	reloadData, err := r.Load(path)
	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err)
	}
	reloadMotion := reloadData.(*vmd.VmdMotion)

	if reloadMotion.Name() != motion.Name() {
		t.Errorf("Expected model name to be '%s', got %q", motion.Name(), reloadMotion.Name())
	}

	if reloadMotion.BoneFrames.Contains("センター") == false {
		t.Errorf("Expected センター to be contained in bone frames")
	}

	if reloadMotion.BoneFrames.Get("センター").Contains(0) == false {
		t.Errorf("Expected センター to contain frame 0")
	}

	reloadBf := reloadMotion.BoneFrames.Get("センター").Get(0)
	if reloadBf.Position.NearEquals(bf.Position, 1e-8) == false {
		t.Errorf("Expected position to be %v, got %v", bf.Position, reloadBf.Position.MMD())
	}

}

func TestVmdWriter_Write2(t *testing.T) {
	// Test case 1: Successful read
	readPath := "../../../test_resources/サンプルモーション.vmd"

	r := NewVmdRepository(true)
	data, err := r.Load(readPath)

	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err)
	}

	outputPath := "../../../test_resources/test_output.vmd"

	motion := data.(*vmd.VmdMotion)

	// Create a VmdWriter instance
	err = r.Save(outputPath, motion, false)
	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err)
	}

	reloadData, err := r.Load(outputPath)
	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err)
	}
	reloadMotion := reloadData.(*vmd.VmdMotion)

	if reloadMotion.Name() != motion.Name() {
		t.Errorf("Expected model name to be '%s', got %q", motion.Name(), reloadMotion.Name())
	}

}

func TestVmdWriter_Write3(t *testing.T) {
	// Test case 1: Successful read
	readPath := "../../../test_resources/ドクヘビ_178cmカメラ.vmd"

	r := NewVmdRepository(true)
	model, err := r.Load(readPath)

	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err)
	}

	outputPath := "../../../test_resources/test_output.vmd"
	motion := model.(*vmd.VmdMotion)

	// Create a VmdWriter instance
	err = r.Save(outputPath, motion, false)
	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err)
	}

	reloadData, err := r.Load(outputPath)
	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err)
	}
	reloadMotion := reloadData.(*vmd.VmdMotion)

	if reloadMotion.Name() != motion.Name() {
		t.Errorf("Expected model name to be '%s', got %q", motion.Name(), reloadMotion.Name())
	}

}

func TestVmdWriter_Write4(t *testing.T) {
	// Test case 1: Successful read
	readPath := "../../../test_resources/モーフ_まばたき.vmd"

	r := NewVmdRepository(true)
	model, err := r.Load(readPath)

	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err)
	}

	outputPath := "../../../test_resources/test_output.vmd"
	motion := model.(*vmd.VmdMotion)

	// Create a VmdWriter instance
	err = r.Save(outputPath, motion, false)
	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err)
	}

	reloadData, err := r.Load(outputPath)
	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err)
	}
	reloadMotion := reloadData.(*vmd.VmdMotion)

	if reloadMotion.Name() != motion.Name() {
		t.Errorf("Expected model name to be '%s', got %q", motion.Name(), reloadMotion.Name())
	}

}
