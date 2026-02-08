// 指示: miu200521358
package ui

import (
	"math"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/miu200521358/mlib_go/pkg/adapter/io_motion"
	"github.com/miu200521358/mlib_go/pkg/domain/mmath"
	"github.com/miu200521358/mlib_go/pkg/domain/motion"
)

func TestBuildMotionVmdDefaultOutputPath(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "bone csv with timestamp",
			input: filepath.Join("C:", "work", "sample_bone_20260102_030405.csv"),
			want:  filepath.Join("C:", "work", "sample.vmd"),
		},
		{
			name:  "morph csv with timestamp",
			input: filepath.Join("C:", "work", "sample_morph_20260102_030405.csv"),
			want:  filepath.Join("C:", "work", "sample.vmd"),
		},
		{
			name:  "plain csv",
			input: filepath.Join("C:", "work", "sample.csv"),
			want:  filepath.Join("C:", "work", "sample.vmd"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := buildMotionVmdDefaultOutputPath(tt.input)
			if got != tt.want {
				t.Fatalf("expected %q, got %q", tt.want, got)
			}
		})
	}
}

func TestImportMotionCsvByInputPath_BoneInputLoadsMorphPair(t *testing.T) {
	originalNowFunc := motionCsvNowFunc
	motionCsvNowFunc = func() time.Time {
		return time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC)
	}
	defer func() {
		motionCsvNowFunc = originalNowFunc
	}()

	tempDir := t.TempDir()
	sourceMotion := buildMotionCsvRoundTripSample(filepath.Join(tempDir, "source.vmd"))
	if err := exportMotionCsvByOutputPath(filepath.Join(tempDir, "sample.csv"), sourceMotion); err != nil {
		t.Fatalf("expected export to succeed, got %v", err)
	}

	boneCsvPath := filepath.Join(tempDir, "sample_bone_20260102_030405.csv")
	outputVmdPath := filepath.Join(tempDir, "imported.vmd")
	if err := importMotionCsvByInputPath(boneCsvPath, outputVmdPath); err != nil {
		t.Fatalf("expected import to succeed, got %v", err)
	}

	verifyImportedMotionCsvRoundTrip(t, outputVmdPath)
}

func TestImportMotionCsvByInputPath_MorphInputLoadsBonePair(t *testing.T) {
	originalNowFunc := motionCsvNowFunc
	motionCsvNowFunc = func() time.Time {
		return time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC)
	}
	defer func() {
		motionCsvNowFunc = originalNowFunc
	}()

	tempDir := t.TempDir()
	sourceMotion := buildMotionCsvRoundTripSample(filepath.Join(tempDir, "source.vmd"))
	if err := exportMotionCsvByOutputPath(filepath.Join(tempDir, "sample.csv"), sourceMotion); err != nil {
		t.Fatalf("expected export to succeed, got %v", err)
	}

	morphCsvPath := filepath.Join(tempDir, "sample_morph_20260102_030405.csv")
	outputVmdPath := filepath.Join(tempDir, "imported.vmd")
	if err := importMotionCsvByInputPath(morphCsvPath, outputVmdPath); err != nil {
		t.Fatalf("expected import to succeed, got %v", err)
	}

	verifyImportedMotionCsvRoundTrip(t, outputVmdPath)
}

func TestImportMotionCsvByInputPath_InvalidFormat(t *testing.T) {
	tempDir := t.TempDir()
	inputPath := filepath.Join(tempDir, "invalid.csv")
	outputPath := filepath.Join(tempDir, "invalid.vmd")
	if err := os.WriteFile(inputPath, []byte("hoge,fuga\n1,2\n"), 0o644); err != nil {
		t.Fatalf("expected invalid csv to be written, got %v", err)
	}

	err := importMotionCsvByInputPath(inputPath, outputPath)
	if err == nil {
		t.Fatalf("expected import to fail for invalid format")
	}
}

func buildMotionCsvRoundTripSample(path string) *motion.VmdMotion {
	motionData := motion.NewVmdMotion(path)
	motionData.SetName("sample")

	boneFrame := motion.NewBoneFrame(12)
	position := mmath.NewVec3()
	position.X = 1.25
	position.Y = -2.5
	position.Z = 3.75
	boneFrame.Position = &position
	rotation := mmath.NewQuaternionFromDegrees(10, 20, 30)
	boneFrame.Rotation = &rotation
	interpolation := make([]byte, 64)
	for i := range interpolation {
		interpolation[i] = byte(i)
	}
	boneFrame.Curves = motion.NewBoneCurvesByValues(interpolation)
	motionData.AppendBoneFrame("センター", boneFrame)

	morphFrame := motion.NewMorphFrame(34)
	morphFrame.Ratio = 0.5
	motionData.AppendMorphFrame("笑い", morphFrame)

	return motionData
}

func verifyImportedMotionCsvRoundTrip(t *testing.T, outputVmdPath string) {
	t.Helper()

	loaded, err := io_motion.NewVmdVpdRepository().Load(outputVmdPath)
	if err != nil {
		t.Fatalf("expected imported vmd to load, got %v", err)
	}
	importedMotion, ok := loaded.(*motion.VmdMotion)
	if !ok || importedMotion == nil {
		t.Fatalf("expected loaded data to be VmdMotion")
	}

	boneFrames := importedMotion.BoneFrames.Get("センター")
	if boneFrames == nil || !boneFrames.Has(12) {
		t.Fatalf("expected bone frame [センター:12] to exist")
	}
	boneFrame := boneFrames.Get(12)
	if boneFrame == nil {
		t.Fatalf("expected bone frame to be non-nil")
	}
	if boneFrame.Position == nil {
		t.Fatalf("expected bone position to be non-nil")
	}
	if math.Abs(boneFrame.Position.X-1.25) > 1e-6 {
		t.Fatalf("unexpected position X: %f", boneFrame.Position.X)
	}
	if math.Abs(boneFrame.Position.Y+2.5) > 1e-6 {
		t.Fatalf("unexpected position Y: %f", boneFrame.Position.Y)
	}
	if math.Abs(boneFrame.Position.Z-3.75) > 1e-6 {
		t.Fatalf("unexpected position Z: %f", boneFrame.Position.Z)
	}
	if boneFrame.Rotation == nil {
		t.Fatalf("expected bone rotation to be non-nil")
	}
	expectedDegrees := mmath.NewQuaternionFromDegrees(10, 20, 30).ToMMDDegrees()
	degrees := boneFrame.Rotation.ToMMDDegrees()
	if math.Abs(degrees.X-expectedDegrees.X) > 1e-4 {
		t.Fatalf("unexpected rotation X: %f", degrees.X)
	}
	if math.Abs(degrees.Y-expectedDegrees.Y) > 1e-4 {
		t.Fatalf("unexpected rotation Y: %f", degrees.Y)
	}
	if math.Abs(degrees.Z-expectedDegrees.Z) > 1e-4 {
		t.Fatalf("unexpected rotation Z: %f", degrees.Z)
	}
	if boneFrame.Curves == nil || len(boneFrame.Curves.Values) < 64 {
		t.Fatalf("expected bone curves values to exist")
	}
	sourceInterpolation := make([]byte, 64)
	for i := range sourceInterpolation {
		sourceInterpolation[i] = byte(i)
	}
	expectedInterpolation := motion.NewBoneCurvesByValues(sourceInterpolation).Merge(false)
	for i := 0; i < 64; i++ {
		if boneFrame.Curves.Values[i] != expectedInterpolation[i] {
			t.Fatalf("unexpected interpolation value at %d: got=%d", i, boneFrame.Curves.Values[i])
		}
	}

	morphFrames := importedMotion.MorphFrames.Get("笑い")
	if morphFrames == nil || !morphFrames.Has(34) {
		t.Fatalf("expected morph frame [笑い:34] to exist")
	}
	morphFrame := morphFrames.Get(34)
	if morphFrame == nil {
		t.Fatalf("expected morph frame to be non-nil")
	}
	if math.Abs(morphFrame.Ratio-0.5) > 1e-6 {
		t.Fatalf("unexpected morph ratio: %f", morphFrame.Ratio)
	}
}
