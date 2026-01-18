// 指示: miu200521358
package vpd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/miu200521358/mlib_go/pkg/domain/mmath"
	"github.com/miu200521358/mlib_go/pkg/domain/motion"
	"golang.org/x/text/encoding/japanese"
	"gonum.org/v1/gonum/spatial/r3"
)

func TestVpdRepository_Load(t *testing.T) {
	path := writeVpdFile(t, strings.Join([]string{
		"Vocaloid Pose Data file",
		"1",
		"Sample.osm; // 親ファイル名",
		"{センター",
		"0.5,1.25,2.75; // trans",
		"0,0,0,1; // Quaternion",
	}, "\n"))

	r := NewVpdRepository()
	data, err := r.Load(path)
	if err != nil {
		t.Fatalf("Expected error to be nil, got %q", err)
	}

	motionData, ok := data.(*motion.VmdMotion)
	if !ok {
		t.Fatalf("Expected motion type to be *VmdMotion, got %T", data)
	}
	if motionData.Name() != "Sample" {
		t.Errorf("Expected model name to be %q, got %q", "Sample", motionData.Name())
	}

	frames := motionData.BoneFrames.Get("センター")
	bf := frames.Get(motion.Frame(0))
	if bf == nil {
		t.Fatalf("Expected bone frame to be not nil")
	}
	if bf.Position == nil {
		t.Fatalf("Expected Position to be not nil")
	}
	if bf.Rotation == nil {
		t.Fatalf("Expected Rotation to be not nil")
	}

	expectedPos := mmath.Vec3{Vec: r3.Vec{X: 0.5, Y: 1.25, Z: 2.75}}
	if !bf.Position.MMD().NearEquals(expectedPos, 1e-8) {
		t.Errorf("Expected Position to be %v, got %v", expectedPos, bf.Position.MMD())
	}
	if 1-bf.Rotation.Dot(mmath.NewQuaternionByValues(0, 0, 0, 1)) > 1e-8 {
		t.Errorf("Expected Rotation to be identity, got %v", bf.Rotation)
	}
}

func TestVpdRepository_Load_InvalidSignature(t *testing.T) {
	path := writeVpdFile(t, strings.Join([]string{
		"Invalid Signature",
		"1",
		"Sample.osm; // 親ファイル名",
	}, "\n"))

	r := NewVpdRepository()
	if _, err := r.Load(path); err == nil {
		t.Fatalf("Expected error to be not nil")
	}
}

func writeVpdFile(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "sample.vpd")
	encoded, err := japanese.ShiftJIS.NewEncoder().Bytes([]byte(content))
	if err != nil {
		t.Fatalf("Expected Shift-JIS encode to succeed, got %q", err)
	}
	if err := os.WriteFile(path, encoded, 0o644); err != nil {
		t.Fatalf("Expected write to succeed, got %q", err)
	}
	return path
}
