// 指示: miu200521358
package vrm

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/miu200521358/mlib_go/pkg/domain/model"
	"github.com/miu200521358/mlib_go/pkg/domain/model/vrm"
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
	if pmxModel.VrmData.Version != vrm.VRM_VERSION_1 {
		t.Fatalf("expected VRM_VERSION_1, got %s", pmxModel.VrmData.Version)
	}
	if pmxModel.VrmData.Profile != vrm.VRM_PROFILE_VROID {
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
	if extra.Position.X <= 0 {
		t.Fatalf("expected x to keep plus direction for VRM1 VRoid profile, got %f", extra.Position.X)
	}
	if extra.Position.Z >= 0 {
		t.Fatalf("expected z to be converted to minus for VRM1 VRoid profile, got %f", extra.Position.Z)
	}
	if extra.Position.X < 1.2 || extra.Position.X > 1.3 {
		t.Fatalf("expected x to be scaled by 12.5, got %f", extra.Position.X)
	}
}

func TestVrmRepositoryLoadVrm0VroidKeepsLegacyAxisConversion(t *testing.T) {
	repository := NewVrmRepository()
	tempDir := t.TempDir()
	path := filepath.Join(tempDir, "avatar.vrm")

	doc := map[string]any{
		"asset": map[string]any{
			"version":   "2.0",
			"generator": "VRoid Studio v0.14.0",
		},
		"extensionsUsed": []string{"VRM"},
		"nodes": []any{
			map[string]any{
				"name":        "hips_node",
				"translation": []float64{0, 0.9, 0},
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
	if pmxModel.VrmData.Version != vrm.VRM_VERSION_0 {
		t.Fatalf("expected VRM_VERSION_0, got %s", pmxModel.VrmData.Version)
	}
	if pmxModel.VrmData.Profile != vrm.VRM_PROFILE_VROID {
		t.Fatalf("expected VRM_PROFILE_VROID, got %s", pmxModel.VrmData.Profile)
	}

	extra, err := pmxModel.Bones.GetByName("extra_node")
	if err != nil || extra == nil {
		t.Fatalf("expected extra_node bone: %v", err)
	}
	if extra.Position.X >= 0 {
		t.Fatalf("expected x to be converted to minus for VRM0 VRoid profile, got %f", extra.Position.X)
	}
	if extra.Position.Z <= 0 {
		t.Fatalf("expected z to keep plus direction for VRM0 VRoid profile, got %f", extra.Position.Z)
	}
	if extra.Position.X > -1.2 || extra.Position.X < -1.3 {
		t.Fatalf("expected x to be scaled by 12.5, got %f", extra.Position.X)
	}
}

func TestVrmRepositoryLoadVrm0UniVrmUsesMmdScaleConversion(t *testing.T) {
	repository := NewVrmRepository()
	tempDir := t.TempDir()
	path := filepath.Join(tempDir, "avatar.vrm")

	doc := map[string]any{
		"asset": map[string]any{
			"version":   "2.0",
			"generator": "UniGLTF-1.28",
		},
		"extensionsUsed": []string{"VRM"},
		"nodes": []any{
			map[string]any{
				"name":        "hips_node",
				"translation": []float64{0, 0.9, 0},
			},
			map[string]any{
				"name":        "extra_node",
				"translation": []float64{0.1, 0.3, 0.2},
			},
		},
		"extensions": map[string]any{
			"VRM": map[string]any{
				"exporterVersion": "UniVRM-0.51.0",
				"humanoid": map[string]any{
					"humanBones": []any{
						map[string]any{"bone": "hips", "node": 0},
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
	if pmxModel.VrmData.Version != vrm.VRM_VERSION_0 {
		t.Fatalf("expected VRM_VERSION_0, got %s", pmxModel.VrmData.Version)
	}
	if pmxModel.VrmData.Profile != vrm.VRM_PROFILE_STANDARD {
		t.Fatalf("expected VRM_PROFILE_STANDARD, got %s", pmxModel.VrmData.Profile)
	}

	extra, err := pmxModel.Bones.GetByName("extra_node")
	if err != nil || extra == nil {
		t.Fatalf("expected extra_node bone: %v", err)
	}
	if extra.Position.X >= 0 {
		t.Fatalf("expected x to be converted to minus for UniVRM conversion, got %f", extra.Position.X)
	}
	if extra.Position.Z <= 0 {
		t.Fatalf("expected z to keep plus direction for UniVRM conversion, got %f", extra.Position.Z)
	}
	if extra.Position.X > -1.2 || extra.Position.X < -1.3 {
		t.Fatalf("expected x to be scaled by 12.5, got %f", extra.Position.X)
	}
}

func TestExportArtifactsWritesGltfAndTextures(t *testing.T) {
	tempDir := t.TempDir()
	vrmPath := filepath.Join(tempDir, "avatar.vrm")
	gltfDir := filepath.Join(tempDir, "out", "glTF")
	textureDir := filepath.Join(tempDir, "out", "tex")
	pngBytes := []byte{
		0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A,
		0x00, 0x00, 0x00, 0x0D, 0x49, 0x48, 0x44, 0x52,
		0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x01,
		0x08, 0x06, 0x00, 0x00, 0x00, 0x1F, 0x15, 0xC4,
		0x89, 0x00, 0x00, 0x00, 0x0A, 0x49, 0x44, 0x41,
		0x54, 0x78, 0x9C, 0x63, 0x60, 0x00, 0x00, 0x00,
		0x02, 0x00, 0x01, 0xE5, 0x27, 0xD4, 0xA2, 0x00,
		0x00, 0x00, 0x00, 0x49, 0x45, 0x4E, 0x44, 0xAE,
		0x42, 0x60, 0x82,
	}

	doc := map[string]any{
		"asset": map[string]any{
			"version": "2.0",
		},
		"buffers": []any{
			map[string]any{
				"byteLength": len(pngBytes),
			},
		},
		"bufferViews": []any{
			map[string]any{
				"buffer":     0,
				"byteOffset": 0,
				"byteLength": len(pngBytes),
			},
		},
		"images": []any{
			map[string]any{
				"name":       "face",
				"bufferView": 0,
				"mimeType":   "image/png",
			},
		},
	}
	writeGLBFileForTestWithBin(t, vrmPath, doc, pngBytes)

	result, err := ExportArtifacts(vrmPath, gltfDir, textureDir)
	if err != nil {
		t.Fatalf("export failed: %v", err)
	}
	if result == nil {
		t.Fatalf("result is nil")
	}
	if _, err := os.Stat(result.GltfPath); err != nil {
		t.Fatalf("gltf output not found: %v", err)
	}
	if _, err := os.Stat(result.BinPath); err != nil {
		t.Fatalf("bin output not found: %v", err)
	}
	if len(result.TextureNames) != 1 {
		t.Fatalf("texture name length mismatch: %d", len(result.TextureNames))
	}
	if strings.TrimSpace(result.TextureNames[0]) == "" {
		t.Fatalf("texture name is empty")
	}
	if _, err := os.Stat(filepath.Join(textureDir, result.TextureNames[0])); err != nil {
		t.Fatalf("texture output not found: %v", err)
	}
}

func TestVrmRepositoryLoadCreatesMeshFromPrimitive(t *testing.T) {
	repository := NewVrmRepository()
	tempDir := t.TempDir()
	path := filepath.Join(tempDir, "mesh.vrm")

	positions := []float32{
		0.0, 0.0, 0.0,
		0.0, 1.0, 0.0,
		1.0, 0.0, 0.0,
	}
	normals := []float32{
		0.0, 0.0, 1.0,
		0.0, 0.0, 1.0,
		0.0, 0.0, 1.0,
	}
	uvs := []float32{
		0.0, 0.0,
		0.0, 1.0,
		1.0, 0.0,
	}
	indices := []uint16{0, 1, 2}

	binChunk := buildInterleavedBinForMeshTest(t, positions, normals, uvs, indices)
	doc := map[string]any{
		"asset": map[string]any{
			"version":   "2.0",
			"generator": "VRM Test",
		},
		"extensionsUsed": []string{"VRMC_vrm"},
		"nodes": []any{
			map[string]any{
				"name": "hips_node",
			},
			map[string]any{
				"name": "mesh_node",
				"mesh": 0,
				"skin": 0,
			},
		},
		"skins": []any{
			map[string]any{
				"joints": []int{0},
			},
		},
		"meshes": []any{
			map[string]any{
				"name": "mesh0",
				"primitives": []any{
					map[string]any{
						"attributes": map[string]any{
							"POSITION":   0,
							"NORMAL":     1,
							"TEXCOORD_0": 2,
						},
						"indices":  3,
						"material": 0,
						"mode":     4,
					},
				},
			},
		},
		"materials": []any{
			map[string]any{
				"name": "body",
				"pbrMetallicRoughness": map[string]any{
					"baseColorFactor": []float64{1.0, 1.0, 1.0, 1.0},
				},
			},
		},
		"buffers": []any{
			map[string]any{
				"byteLength": len(binChunk),
			},
		},
		"bufferViews": []any{
			map[string]any{
				"buffer":     0,
				"byteOffset": 0,
				"byteLength": len(positions) * 4,
			},
			map[string]any{
				"buffer":     0,
				"byteOffset": len(positions) * 4,
				"byteLength": len(normals) * 4,
			},
			map[string]any{
				"buffer":     0,
				"byteOffset": (len(positions) + len(normals)) * 4,
				"byteLength": len(uvs) * 4,
			},
			map[string]any{
				"buffer":     0,
				"byteOffset": (len(positions) + len(normals) + len(uvs)) * 4,
				"byteLength": len(indices) * 2,
			},
		},
		"accessors": []any{
			map[string]any{
				"bufferView":    0,
				"componentType": 5126,
				"count":         3,
				"type":          "VEC3",
			},
			map[string]any{
				"bufferView":    1,
				"componentType": 5126,
				"count":         3,
				"type":          "VEC3",
			},
			map[string]any{
				"bufferView":    2,
				"componentType": 5126,
				"count":         3,
				"type":          "VEC2",
			},
			map[string]any{
				"bufferView":    3,
				"componentType": 5123,
				"count":         3,
				"type":          "SCALAR",
			},
		},
		"extensions": map[string]any{
			"VRMC_vrm": map[string]any{
				"specVersion": "1.0",
				"humanoid": map[string]any{
					"humanBones": map[string]any{
						"hips": map[string]any{"node": 0},
					},
				},
			},
		},
	}
	writeGLBFileForUsecaseMeshTest(t, path, doc, binChunk)

	hashableModel, err := repository.Load(path)
	if err != nil {
		t.Fatalf("load failed: %v", err)
	}
	pmxModel, ok := hashableModel.(*model.PmxModel)
	if !ok {
		t.Fatalf("expected *model.PmxModel, got %T", hashableModel)
	}
	if pmxModel.Vertices.Len() != 3 {
		t.Fatalf("expected 3 vertices, got %d", pmxModel.Vertices.Len())
	}
	if pmxModel.Faces.Len() != 1 {
		t.Fatalf("expected 1 face, got %d", pmxModel.Faces.Len())
	}
	if pmxModel.Materials.Len() != 1 {
		t.Fatalf("expected 1 material, got %d", pmxModel.Materials.Len())
	}
}

func TestVrmRepositoryLoadContinuesWhenNormalAccessorIsInvalid(t *testing.T) {
	repository := NewVrmRepository()
	tempDir := t.TempDir()
	path := filepath.Join(tempDir, "mesh_invalid_normal.vrm")

	positions := []float32{
		0.0, 0.0, 0.0,
		0.0, 1.0, 0.0,
		1.0, 0.0, 0.0,
	}
	uvs := []float32{
		0.0, 0.0,
		0.0, 1.0,
		1.0, 0.0,
	}
	indices := []uint16{0, 1, 2}

	binChunk := buildInterleavedBinForMeshTest(t, positions, nil, uvs, indices)
	doc := map[string]any{
		"asset": map[string]any{
			"version":   "2.0",
			"generator": "VRM Test",
		},
		"extensionsUsed": []string{"VRMC_vrm"},
		"nodes": []any{
			map[string]any{
				"name": "hips_node",
			},
			map[string]any{
				"name": "mesh_node",
				"mesh": 0,
				"skin": 0,
			},
		},
		"skins": []any{
			map[string]any{
				"joints": []int{0},
			},
		},
		"meshes": []any{
			map[string]any{
				"name": "mesh0",
				"primitives": []any{
					map[string]any{
						"attributes": map[string]any{
							"POSITION":   0,
							"NORMAL":     99,
							"TEXCOORD_0": 1,
						},
						"indices":  2,
						"material": 0,
						"mode":     4,
					},
				},
			},
		},
		"materials": []any{
			map[string]any{
				"name": "body",
				"pbrMetallicRoughness": map[string]any{
					"baseColorFactor": []float64{1.0, 1.0, 1.0, 1.0},
				},
			},
		},
		"buffers": []any{
			map[string]any{
				"byteLength": len(binChunk),
			},
		},
		"bufferViews": []any{
			map[string]any{
				"buffer":     0,
				"byteOffset": 0,
				"byteLength": len(positions) * 4,
			},
			map[string]any{
				"buffer":     0,
				"byteOffset": len(positions) * 4,
				"byteLength": len(uvs) * 4,
			},
			map[string]any{
				"buffer":     0,
				"byteOffset": (len(positions) + len(uvs)) * 4,
				"byteLength": len(indices) * 2,
			},
		},
		"accessors": []any{
			map[string]any{
				"bufferView":    0,
				"componentType": 5126,
				"count":         3,
				"type":          "VEC3",
			},
			map[string]any{
				"bufferView":    1,
				"componentType": 5126,
				"count":         3,
				"type":          "VEC2",
			},
			map[string]any{
				"bufferView":    2,
				"componentType": 5123,
				"count":         3,
				"type":          "SCALAR",
			},
		},
		"extensions": map[string]any{
			"VRMC_vrm": map[string]any{
				"specVersion": "1.0",
				"humanoid": map[string]any{
					"humanBones": map[string]any{
						"hips": map[string]any{"node": 0},
					},
				},
			},
		},
	}
	writeGLBFileForUsecaseMeshTest(t, path, doc, binChunk)

	hashableModel, err := repository.Load(path)
	if err != nil {
		t.Fatalf("load failed: %v", err)
	}
	pmxModel, ok := hashableModel.(*model.PmxModel)
	if !ok {
		t.Fatalf("expected *model.PmxModel, got %T", hashableModel)
	}
	if pmxModel.Vertices.Len() != 3 {
		t.Fatalf("expected 3 vertices, got %d", pmxModel.Vertices.Len())
	}
	vertex, err := pmxModel.Vertices.Get(0)
	if err != nil {
		t.Fatalf("get vertex failed: %v", err)
	}
	if vertex == nil {
		t.Fatalf("vertex is nil")
	}
	if vertex.Normal.Y < 0.99 {
		t.Fatalf("expected fallback normal Y close to 1.0, got %f", vertex.Normal.Y)
	}
}

func TestShouldSkipPrimitiveForUnsupportedTargets(t *testing.T) {
	indices := 0
	material := 0
	mode := 4
	seen := map[string]int{}
	withTargets := gltfPrimitive{
		Attributes: map[string]int{"POSITION": 1},
		Indices:    &indices,
		Material:   &material,
		Mode:       &mode,
		Targets: []map[string]int{
			{"POSITION": 2},
		},
	}
	withoutTargets := gltfPrimitive{
		Attributes: map[string]int{"POSITION": 1},
		Indices:    &indices,
		Material:   &material,
		Mode:       &mode,
	}

	if shouldSkipPrimitiveForUnsupportedTargets(withTargets, 0, seen) {
		t.Fatalf("first target primitive should not be skipped")
	}
	if !shouldSkipPrimitiveForUnsupportedTargets(withTargets, 1, seen) {
		t.Fatalf("duplicated target primitive should be skipped")
	}
	if shouldSkipPrimitiveForUnsupportedTargets(withoutTargets, 2, seen) {
		t.Fatalf("primitive without targets should not be skipped")
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

// writeGLBFileForTestWithBin はテスト用のJSON/BINをGLBとして書き込む。
func writeGLBFileForTestWithBin(t *testing.T, path string, doc map[string]any, binChunk []byte) {
	t.Helper()
	jsonBytes, err := json.Marshal(doc)
	if err != nil {
		t.Fatalf("json marshal failed: %v", err)
	}
	jsonPadSize := (4 - (len(jsonBytes) % 4)) % 4
	if jsonPadSize > 0 {
		jsonBytes = append(jsonBytes, bytes.Repeat([]byte(" "), jsonPadSize)...)
	}
	binBytes := append([]byte(nil), binChunk...)
	if len(binBytes) > 0 {
		binPadSize := (4 - (len(binBytes) % 4)) % 4
		if binPadSize > 0 {
			binBytes = append(binBytes, bytes.Repeat([]byte{0x00}, binPadSize)...)
		}
	}

	totalLength := uint32(12 + 8 + len(jsonBytes))
	if len(binBytes) > 0 {
		totalLength += uint32(8 + len(binBytes))
	}
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
		t.Fatalf("write json chunk type failed: %v", err)
	}
	if _, err := buf.Write(jsonBytes); err != nil {
		t.Fatalf("write json chunk body failed: %v", err)
	}
	if len(binBytes) > 0 {
		if err := binary.Write(&buf, binary.LittleEndian, uint32(len(binBytes))); err != nil {
			t.Fatalf("write bin chunk length failed: %v", err)
		}
		if err := binary.Write(&buf, binary.LittleEndian, uint32(0x004E4942)); err != nil {
			t.Fatalf("write bin chunk type failed: %v", err)
		}
		if _, err := buf.Write(binBytes); err != nil {
			t.Fatalf("write bin chunk body failed: %v", err)
		}
	}
	if err := os.WriteFile(path, buf.Bytes(), 0o644); err != nil {
		t.Fatalf("write file failed: %v", err)
	}
}

// buildInterleavedBinForMeshTest はメッシュ検証用のBINチャンクを構築する。
func buildInterleavedBinForMeshTest(t *testing.T, positions []float32, normals []float32, uvs []float32, indices []uint16) []byte {
	t.Helper()
	var buf bytes.Buffer
	for _, value := range positions {
		if err := binary.Write(&buf, binary.LittleEndian, value); err != nil {
			t.Fatalf("write position failed: %v", err)
		}
	}
	for _, value := range normals {
		if err := binary.Write(&buf, binary.LittleEndian, value); err != nil {
			t.Fatalf("write normal failed: %v", err)
		}
	}
	for _, value := range uvs {
		if err := binary.Write(&buf, binary.LittleEndian, value); err != nil {
			t.Fatalf("write uv failed: %v", err)
		}
	}
	for _, value := range indices {
		if err := binary.Write(&buf, binary.LittleEndian, value); err != nil {
			t.Fatalf("write index failed: %v", err)
		}
	}
	return buf.Bytes()
}

// writeGLBFileForUsecaseMeshTest はテスト用JSON/BINをGLBとして書き込む。
func writeGLBFileForUsecaseMeshTest(t *testing.T, path string, doc map[string]any, binChunk []byte) {
	t.Helper()
	writeGLBFileForTestWithBin(t, path, doc, binChunk)
}
