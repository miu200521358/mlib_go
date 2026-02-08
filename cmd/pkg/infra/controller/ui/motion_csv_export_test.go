// 指示: miu200521358
package ui

import (
	"encoding/csv"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"testing"
	"time"

	"github.com/miu200521358/mlib_go/pkg/domain/mmath"
	"github.com/miu200521358/mlib_go/pkg/domain/motion"
)

func TestExportMotionCsv(t *testing.T) {
	originalNowFunc := motionCsvNowFunc
	motionCsvNowFunc = func() time.Time {
		return time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC)
	}
	defer func() {
		motionCsvNowFunc = originalNowFunc
	}()

	tempDir := t.TempDir()
	inputPath := filepath.Join(tempDir, "sample.vmd")

	motionData := motion.NewVmdMotion(inputPath)

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

	if err := exportMotionCsv(inputPath, motionData); err != nil {
		t.Fatalf("expected export to succeed, got %v", err)
	}

	bonePath := filepath.Join(tempDir, "sample_bone_20260102_030405.csv")
	morphPath := filepath.Join(tempDir, "sample_morph_20260102_030405.csv")
	if _, err := os.Stat(bonePath); err != nil {
		t.Fatalf("expected bone csv to exist, got %v", err)
	}
	if _, err := os.Stat(morphPath); err != nil {
		t.Fatalf("expected morph csv to exist, got %v", err)
	}

	boneRecords := loadCsvRecords(t, bonePath)
	if len(boneRecords) != 2 {
		t.Fatalf("expected 2 bone rows, got %d", len(boneRecords))
	}
	expectedBoneHeader := []string{
		"ボーン名", "フレーム", "位置X", "位置Y", "位置Z", "回転X", "回転Y", "回転Z",
		"【X_x1】", "Y_x1", "Z_x1", "R_x1", "【X_y1】", "Y_y1", "Z_y1", "R_y1",
		"【X_x2】", "Y_x2", "Z_x2", "R_x2", "【X_y2】", "Y_y2", "Z_y2", "R_y2",
		"【Y_x1】", "Z_x1", "R_x1", "X_y1", "【Y_y1】", "Z_y1", "R_y1", "X_x2",
		"【Y_x2】", "Z_x2", "R_x2", "X_y2", "【Y_y2】", "Z_y2", "R_y2", "1",
		"【Z_x1】", "R_x1", "X_y1", "Y_y1", "【Z_y1】", "R_y1", "X_x2", "Y_x2",
		"【Z_x2】", "R_x2", "X_y2", "Y_y2", "【Z_y2】", "R_y2", "1", "0",
		"【R_x1】", "X_y1", "Y_y1", "Z_y1", "【R_y1】", "X_x2", "Y_x2", "Z_x2",
		"【R_x2】", "X_y2", "Y_y2", "Z_y2", "【R_y2】", "01", "00", "00",
	}
	if len(boneRecords[0]) != len(expectedBoneHeader) {
		t.Fatalf("expected bone header columns %d, got %d", len(expectedBoneHeader), len(boneRecords[0]))
	}
	for i := range expectedBoneHeader {
		if boneRecords[0][i] != expectedBoneHeader[i] {
			t.Fatalf("unexpected bone header at %d: expected %q, got %q", i, expectedBoneHeader[i], boneRecords[0][i])
		}
	}

	boneRow := boneRecords[1]
	if boneRow[0] != "センター" {
		t.Fatalf("expected bone name センター, got %q", boneRow[0])
	}
	if boneRow[1] != "12" {
		t.Fatalf("expected frame 12, got %q", boneRow[1])
	}

	assertFloatTextNear(t, boneRow[2], 1.25, 1e-6)
	assertFloatTextNear(t, boneRow[3], -2.5, 1e-6)
	assertFloatTextNear(t, boneRow[4], 3.75, 1e-6)

	degrees := rotation.ToMMDDegrees()
	assertFloatTextNear(t, boneRow[5], degrees.X, 1e-5)
	assertFloatTextNear(t, boneRow[6], degrees.Y, 1e-5)
	assertFloatTextNear(t, boneRow[7], degrees.Z, 1e-5)

	for i := 0; i < 64; i++ {
		expected := strconv.Itoa(i)
		if boneRow[8+i] != expected {
			t.Fatalf("unexpected interpolation value at %d: expected %q, got %q", i, expected, boneRow[8+i])
		}
	}

	morphRecords := loadCsvRecords(t, morphPath)
	expectedMorph := [][]string{
		{"モーフ名", "フレーム", "大きさ"},
		{"笑い", "34", "0.5"},
	}
	if len(morphRecords) != len(expectedMorph) {
		t.Fatalf("expected %d morph rows, got %d", len(expectedMorph), len(morphRecords))
	}
	for rowIndex := range expectedMorph {
		if len(morphRecords[rowIndex]) != len(expectedMorph[rowIndex]) {
			t.Fatalf("expected morph row %d to have %d columns, got %d", rowIndex, len(expectedMorph[rowIndex]), len(morphRecords[rowIndex]))
		}
		for columnIndex := range expectedMorph[rowIndex] {
			if morphRecords[rowIndex][columnIndex] != expectedMorph[rowIndex][columnIndex] {
				t.Fatalf(
					"unexpected morph cell row=%d col=%d: expected %q, got %q",
					rowIndex,
					columnIndex,
					expectedMorph[rowIndex][columnIndex],
					morphRecords[rowIndex][columnIndex],
				)
			}
		}
	}
}

func TestExportMotionCsv_NonVmdPath(t *testing.T) {
	originalNowFunc := motionCsvNowFunc
	motionCsvNowFunc = func() time.Time {
		return time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC)
	}
	defer func() {
		motionCsvNowFunc = originalNowFunc
	}()

	tempDir := t.TempDir()
	inputPath := filepath.Join(tempDir, "sample.vpd")

	motionData := motion.NewVmdMotion(inputPath)
	boneFrame := motion.NewBoneFrame(0)
	motionData.AppendBoneFrame("センター", boneFrame)

	if err := exportMotionCsv(inputPath, motionData); err != nil {
		t.Fatalf("expected export to be skipped, got %v", err)
	}

	entries, err := os.ReadDir(tempDir)
	if err != nil {
		t.Fatalf("expected temp dir to be readable, got %v", err)
	}
	if len(entries) != 0 {
		t.Fatalf("expected no exported files, got %d entries", len(entries))
	}
}

func TestExportMotionCsvByOutputPath(t *testing.T) {
	originalNowFunc := motionCsvNowFunc
	motionCsvNowFunc = func() time.Time {
		return time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC)
	}
	defer func() {
		motionCsvNowFunc = originalNowFunc
	}()

	tempDir := t.TempDir()
	outputPath := filepath.Join(tempDir, "custom_output.csv")
	motionData := motion.NewVmdMotion(filepath.Join(tempDir, "source.vmd"))

	boneFrame := motion.NewBoneFrame(1)
	motionData.AppendBoneFrame("センター", boneFrame)

	morphFrame := motion.NewMorphFrame(2)
	morphFrame.Ratio = 0.75
	motionData.AppendMorphFrame("笑い", morphFrame)

	if err := exportMotionCsvByOutputPath(outputPath, motionData); err != nil {
		t.Fatalf("expected export to succeed, got %v", err)
	}

	bonePath := filepath.Join(tempDir, "custom_output_bone_20260102_030405.csv")
	morphPath := filepath.Join(tempDir, "custom_output_morph_20260102_030405.csv")
	if _, err := os.Stat(bonePath); err != nil {
		t.Fatalf("expected bone csv to exist, got %v", err)
	}
	if _, err := os.Stat(morphPath); err != nil {
		t.Fatalf("expected morph csv to exist, got %v", err)
	}
}

func TestBuildMotionCsvDefaultOutputPath(t *testing.T) {
	path := buildMotionCsvDefaultOutputPath(filepath.Join("C:", "work", "sample.vmd"))
	expected := filepath.Join("C:", "work", "sample.csv")
	if path != expected {
		t.Fatalf("expected %q, got %q", expected, path)
	}
}

func TestBuildMotionBoneInterpolation_Default(t *testing.T) {
	interpolation := buildMotionBoneInterpolation(nil)
	if interpolation[0] != motion.INITIAL_BONE_CURVES[0] {
		t.Fatalf("expected first interpolation to be default %d, got %d", motion.INITIAL_BONE_CURVES[0], interpolation[0])
	}
	if interpolation[63] != motion.INITIAL_BONE_CURVES[63] {
		t.Fatalf("expected last interpolation to be default %d, got %d", motion.INITIAL_BONE_CURVES[63], interpolation[63])
	}
}

func loadCsvRecords(t *testing.T, path string) [][]string {
	t.Helper()
	file, err := os.Open(path)
	if err != nil {
		t.Fatalf("expected csv file to open, got %v", err)
	}
	defer func() {
		_ = file.Close()
	}()

	records, err := csv.NewReader(file).ReadAll()
	if err != nil {
		t.Fatalf("expected csv to parse, got %v", err)
	}
	return records
}

func assertFloatTextNear(t *testing.T, raw string, expected, epsilon float64) {
	t.Helper()
	value, err := strconv.ParseFloat(raw, 64)
	if err != nil {
		t.Fatalf("expected float text %q, got parse error %v", raw, err)
	}
	if math.Abs(value-expected) > epsilon {
		t.Fatalf("expected %f (±%f), got %f", expected, epsilon, value)
	}
}
