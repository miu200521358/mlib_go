// 指示: miu200521358
package vrm

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/miu200521358/mlib_go/pkg/domain/model"
	modelvrm "github.com/miu200521358/mlib_go/pkg/domain/model/vrm"
	"github.com/miu200521358/mlib_go/pkg/shared/base/merr"
)

func TestVrmRepositoryCanLoad(t *testing.T) {
	repository := NewVrmRepository()

	if !repository.CanLoad("sample.vrm") {
		t.Fatalf("expected sample.vrm to be loadable")
	}
	if !repository.CanLoad("sample.VRM") {
		t.Fatalf("expected sample.VRM to be loadable")
	}
	if repository.CanLoad("sample.pmx") {
		t.Fatalf("expected sample.pmx to be not loadable")
	}
}

func TestVrmRepositoryInferName(t *testing.T) {
	repository := NewVrmRepository()

	got := repository.InferName("C:/work/avatar.vrm")
	if got != "avatar" {
		t.Fatalf("expected avatar, got %s", got)
	}
}

func TestVrmRepositoryLoadReturnsExtInvalid(t *testing.T) {
	repository := NewVrmRepository()

	_, err := repository.Load("sample.pmx")
	if err == nil {
		t.Fatalf("expected error to be not nil")
	}
	if merr.ExtractErrorID(err) != "14102" {
		t.Fatalf("expected error id 14102, got %s", merr.ExtractErrorID(err))
	}
}

func TestVrmRepositoryLoadReturnsFileNotFound(t *testing.T) {
	repository := NewVrmRepository()

	_, err := repository.Load(filepath.Join(t.TempDir(), "missing.vrm"))
	if err == nil {
		t.Fatalf("expected error to be not nil")
	}
	if merr.ExtractErrorID(err) != "14101" {
		t.Fatalf("expected error id 14101, got %s", merr.ExtractErrorID(err))
	}
}

func TestVrmRepositoryLoadVrm1PreferredAndBoneMapping(t *testing.T) {
	repository := NewVrmRepository()
	tempDir := t.TempDir()
	path := filepath.Join(tempDir, "avatar.vrm")

	doc := map[string]any{
		"asset": map[string]any{
			"version":   "2.0",
			"generator": "VRoid Studio v1.0.0",
		},
		"extensionsUsed": []string{"VRM", "VRMC_vrm"},
		"nodes": []any{
			map[string]any{
				"name":        "hips_node",
				"translation": []float64{0, 0.9, 0},
				"children":    []int{1},
			},
			map[string]any{
				"name":        "spine_node",
				"translation": []float64{0, 0.2, 0},
				"children":    []int{2},
			},
			map[string]any{
				"name":        "chest_node",
				"translation": []float64{0, 0.2, 0},
			},
			map[string]any{
				"name":        "extra_node",
				"translation": []float64{0.1, 0.3, 0.2},
			},
		},
		"extensions": map[string]any{
			"VRM": map[string]any{
				"exporterVersion": "VRoid Studio v0.14.0",
				"humanoid": map[string]any{
					"humanBones": []any{
						map[string]any{"bone": "hips", "node": 0},
						map[string]any{"bone": "spine", "node": 1},
						map[string]any{"bone": "chest", "node": 2},
					},
				},
			},
			"VRMC_vrm": map[string]any{
				"specVersion": "1.0",
				"humanoid": map[string]any{
					"humanBones": map[string]any{
						"hips":       map[string]any{"node": 0},
						"spine":      map[string]any{"node": 1},
						"upperChest": map[string]any{"node": 2},
					},
				},
			},
		},
	}
	writeGLBFileForTest(t, path, doc)

	hashableModel, err := repository.Load(path)
	if err != nil {
		t.Fatalf("load failed: %v", err)
	}
	pmxModel, ok := hashableModel.(*model.PmxModel)
	if !ok {
		t.Fatalf("expected *model.PmxModel, got %T", hashableModel)
	}
	if pmxModel.VrmData == nil {
		t.Fatalf("expected vrm data")
	}
	if pmxModel.VrmData.Version != modelvrm.VRM_VERSION_1 {
		t.Fatalf("expected VRM_VERSION_1, got %s", pmxModel.VrmData.Version)
	}
	if pmxModel.VrmData.Profile != modelvrm.VRM_PROFILE_VROID {
		t.Fatalf("expected VRM_PROFILE_VROID, got %s", pmxModel.VrmData.Profile)
	}
	if pmxModel.VrmData.Vrm1 == nil {
		t.Fatalf("expected Vrm1 to be not nil")
	}
	if pmxModel.VrmData.Vrm0 != nil {
		t.Fatalf("expected Vrm0 to be nil when VRM1 is selected")
	}

	hips, err := pmxModel.Bones.GetByName("下半身")
	if err != nil || hips == nil {
		t.Fatalf("expected 下半身 bone: %v", err)
	}
	spine, err := pmxModel.Bones.GetByName("上半身")
	if err != nil || spine == nil {
		t.Fatalf("expected 上半身 bone: %v", err)
	}
	upperBody2, err := pmxModel.Bones.GetByName("上半身2")
	if err != nil || upperBody2 == nil {
		t.Fatalf("expected 上半身2 bone: %v", err)
	}
	if pmxModel.Bones.ContainsByName("上半身3") {
		t.Fatalf("unexpected 上半身3 bone")
	}
	if spine.ParentIndex != hips.Index() {
		t.Fatalf("expected 上半身 parent to be 下半身")
	}
	if upperBody2.ParentIndex != spine.Index() {
		t.Fatalf("expected 上半身2 parent to be 上半身")
	}

	extra, err := pmxModel.Bones.GetByName("extra_node")
	if err != nil || extra == nil {
		t.Fatalf("expected extra_node bone: %v", err)
	}
	if extra.Position.Z >= 0 {
		t.Fatalf("expected z to be converted to minus, got %f", extra.Position.Z)
	}
}

// writeGLBFileForTest はテスト用のJSONをGLBとして書き込む。
func writeGLBFileForTest(t *testing.T, path string, doc map[string]any) {
	t.Helper()
	jsonBytes, err := json.Marshal(doc)
	if err != nil {
		t.Fatalf("json marshal failed: %v", err)
	}
	padSize := (4 - (len(jsonBytes) % 4)) % 4
	if padSize > 0 {
		jsonBytes = append(jsonBytes, bytes.Repeat([]byte(" "), padSize)...)
	}
	totalLength := uint32(12 + 8 + len(jsonBytes))

	var buf bytes.Buffer
	if err := binary.Write(&buf, binary.LittleEndian, uint32(0x46546C67)); err != nil {
		t.Fatalf("write magic failed: %v", err)
	}
	if err := binary.Write(&buf, binary.LittleEndian, uint32(2)); err != nil {
		t.Fatalf("write version failed: %v", err)
	}
	if err := binary.Write(&buf, binary.LittleEndian, totalLength); err != nil {
		t.Fatalf("write length failed: %v", err)
	}
	if err := binary.Write(&buf, binary.LittleEndian, uint32(len(jsonBytes))); err != nil {
		t.Fatalf("write chunk length failed: %v", err)
	}
	if err := binary.Write(&buf, binary.LittleEndian, uint32(0x4E4F534A)); err != nil {
		t.Fatalf("write chunk type failed: %v", err)
	}
	if _, err := buf.Write(jsonBytes); err != nil {
		t.Fatalf("write chunk body failed: %v", err)
	}
	if err := os.WriteFile(path, buf.Bytes(), 0o644); err != nil {
		t.Fatalf("write file failed: %v", err)
	}
}
