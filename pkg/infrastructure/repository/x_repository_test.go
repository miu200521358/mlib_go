package repository

import (
	"strings"
	"testing"

	"github.com/miu200521358/mlib_go/pkg/domain/pmx"
)

func TestXRepository_LoadName(t *testing.T) {
	rep := NewXRepository()

	// Test case 1: Successful read
	path := "D:/MMD/MikuMikuDance_v926x64/UserFile/Accessory/咲音マイク.x"
	modelName := rep.LoadName(path)

	expectedModelName := ""
	if modelName != expectedModelName {
		t.Errorf("Expected modelName to be %q, got %q", expectedModelName, modelName)
	}
}

func TestXRepository_Load1(t *testing.T) {
	rep := NewXRepository()

	// Test case 1: Successful read
	path := "D:/MMD/MikuMikuDance_v926x64/UserFile/Accessory/咲音マイク.x"
	data, err := rep.Load(path)

	if err != nil {
		t.Errorf("Expected err to be nil, got %v", err)
	}
	if data == nil {
		t.Errorf("Expected model to be not nil, got nil")
	}
	model := data.(*pmx.PmxModel)

	pmxRep := NewPmxRepository()
	pmxRep.Save("../../../test_resources/test.pmx", model, false)

	pmxPath := strings.Replace(path, ".x", ".pmx", -1)
	expectedData, _ := pmxRep.Load(pmxPath)
	expectedModel := expectedData.(*pmx.PmxModel)

	for i, v := range model.Vertices.Data {
		expectedV := expectedModel.Vertices.Get(i)
		if !v.Position.NearEquals(expectedV.Position, 1e-5) {
			t.Errorf("Expected Position to be %v, got %v", expectedV.Position, v.Position)
		}
	}

	for i, f := range model.Faces.Data {
		expectedF := expectedModel.Faces.Get(i)
		if f.VertexIndexes[0] != expectedF.VertexIndexes[0] || f.VertexIndexes[1] != expectedF.VertexIndexes[1] || f.VertexIndexes[2] != expectedF.VertexIndexes[2] {
			t.Errorf("Expected VertexIndexes to be %v, got %v", expectedF.VertexIndexes, f.VertexIndexes)
		}
	}

	for i, m := range model.Materials.Data {
		expectedT := expectedModel.Materials.Get(i)
		if !m.Diffuse.NearEquals(expectedT.Diffuse, 1e-5) {
			t.Errorf("Expected Diffuse to be %v, got %v", expectedT.Diffuse, m.Diffuse)
		}
		if !m.Ambient.NearEquals(expectedT.Ambient, 1e-5) {
			t.Errorf("Expected Ambient to be %v, got %v", expectedT.Ambient, m.Ambient)
		}
		if !m.Specular.NearEquals(expectedT.Specular, 1e-5) {
			t.Errorf("Expected Specular to be %v, got %v", expectedT.Specular, m.Specular)
		}
		if !m.Edge.NearEquals(expectedT.Edge, 1e-5) {
			t.Errorf("Expected EdgeColor to be %v, got %v", expectedT.Edge, m.Edge)
		}
		if m.DrawFlag != expectedT.DrawFlag {
			t.Errorf("Expected DrawFlag to be %v, got %v", expectedT.DrawFlag, m.DrawFlag)
		}
		if m.EdgeSize != expectedT.EdgeSize {
			t.Errorf("Expected EdgeSize to be %v, got %v", expectedT.EdgeSize, m.EdgeSize)
		}
		if m.TextureIndex != expectedT.TextureIndex {
			t.Errorf("Expected TextureIndex to be %v, got %v", expectedT.TextureIndex, m.TextureIndex)
		}
		if m.SphereTextureIndex != expectedT.SphereTextureIndex {
			t.Errorf("Expected SphereTextureIndex to be %v, got %v", expectedT.SphereTextureIndex, m.SphereTextureIndex)
		}
		if m.ToonTextureIndex != expectedT.ToonTextureIndex {
			t.Errorf("Expected ToonTextureIndex to be %v, got %v", expectedT.ToonTextureIndex, m.ToonTextureIndex)
		}
		if m.SphereMode != expectedT.SphereMode {
			t.Errorf("Expected SphereMode to be %v, got %v", expectedT.SphereMode, m.SphereMode)
		}
		if m.ToonSharingFlag != expectedT.ToonSharingFlag {
			t.Errorf("Expected ToonSharingFlag to be %v, got %v", expectedT.ToonSharingFlag, m.ToonSharingFlag)
		}
		if m.VerticesCount != expectedT.VerticesCount {
			t.Errorf("Expected VerticesCount to be %v, got %v", expectedT.VerticesCount, m.VerticesCount)
		}
	}

}

func TestXRepository_Load2(t *testing.T) {
	rep := NewXRepository()

	// Test case 1: Successful read
	path := "D:/MMD/MikuMikuDance_v926x64/UserFile/Accessory/食べ物/ファミレスメニューセットver1.0 キャベツ鉢/オムライス/Xファイル/オムライス20151226.x"
	data, err := rep.Load(path)

	if err != nil {
		t.Errorf("Expected err to be nil, got %v", err)
	}
	if data == nil {
		t.Errorf("Expected model to be not nil, got nil")
	}
	model := data.(*pmx.PmxModel)

	pmxRep := NewPmxRepository()
	pmxRep.Save("../../../test_resources/test.pmx", model, false)

	pmxPath := strings.Replace(path, ".x", ".pmx", -1)
	expectedData, _ := pmxRep.Load(pmxPath)
	expectedModel := expectedData.(*pmx.PmxModel)

	for i, v := range model.Vertices.Data {
		expectedV := expectedModel.Vertices.Get(i)
		if !v.Position.NearEquals(expectedV.Position, 1e-5) {
			t.Errorf("Expected Position to be %v, got %v", expectedV.Position, v.Position)
		}
	}

	for i, f := range model.Faces.Data {
		expectedF := expectedModel.Faces.Get(i)
		if f.VertexIndexes[0] != expectedF.VertexIndexes[0] || f.VertexIndexes[1] != expectedF.VertexIndexes[1] || f.VertexIndexes[2] != expectedF.VertexIndexes[2] {
			t.Errorf("Expected Face[%d] VertexIndexes to be %v, got %v", f.Index(), expectedF.VertexIndexes, f.VertexIndexes)
		}
	}

	for _, is := range [][]int{{0, 0}, {1, 1}, {2, 2}, {3, 4}} {
		expectIndex, materialIndex := is[0], is[1]
		m := model.Materials.Get(materialIndex)
		expectedT := expectedModel.Materials.Get(expectIndex)
		if !m.Diffuse.NearEquals(expectedT.Diffuse, 1e-5) {
			t.Errorf("Expected Diffuse to be %v, got %v", expectedT.Diffuse, m.Diffuse)
		}
		if !m.Ambient.NearEquals(expectedT.Ambient, 1e-5) {
			t.Errorf("Expected Ambient to be %v, got %v", expectedT.Ambient, m.Ambient)
		}
		if !m.Specular.NearEquals(expectedT.Specular, 1e-5) {
			t.Errorf("Expected Specular to be %v, got %v", expectedT.Specular, m.Specular)
		}
		if !m.Edge.NearEquals(expectedT.Edge, 1e-5) {
			t.Errorf("Expected EdgeColor to be %v, got %v", expectedT.Edge, m.Edge)
		}
		if m.DrawFlag != expectedT.DrawFlag {
			t.Errorf("Expected DrawFlag to be %v, got %v", expectedT.DrawFlag, m.DrawFlag)
		}
		if m.EdgeSize != expectedT.EdgeSize {
			t.Errorf("Expected EdgeSize to be %v, got %v", expectedT.EdgeSize, m.EdgeSize)
		}
		if m.TextureIndex != expectedT.TextureIndex {
			t.Errorf("Expected TextureIndex to be %v, got %v", expectedT.TextureIndex, m.TextureIndex)
		}
		if m.SphereTextureIndex != expectedT.SphereTextureIndex {
			t.Errorf("Expected SphereTextureIndex to be %v, got %v", expectedT.SphereTextureIndex, m.SphereTextureIndex)
		}
		if m.ToonTextureIndex != expectedT.ToonTextureIndex {
			t.Errorf("Expected ToonTextureIndex to be %v, got %v", expectedT.ToonTextureIndex, m.ToonTextureIndex)
		}
		if m.SphereMode != expectedT.SphereMode {
			t.Errorf("Expected SphereMode to be %v, got %v", expectedT.SphereMode, m.SphereMode)
		}
		if m.ToonSharingFlag != expectedT.ToonSharingFlag {
			t.Errorf("Expected ToonSharingFlag to be %v, got %v", expectedT.ToonSharingFlag, m.ToonSharingFlag)
		}
		if m.VerticesCount != expectedT.VerticesCount {
			t.Errorf("Expected VerticesCount to be %v, got %v", expectedT.VerticesCount, m.VerticesCount)
		}
	}

}

func TestXRepository_Load3(t *testing.T) {
	rep := NewXRepository()

	// Test case 1: Successful read
	path := "D:/MMD/MikuMikuDance_v926x64/UserFile/Effect/_色調補正/ikClut ikeno/ikClut.x"
	_, err := rep.Load(path)

	if err == nil {
		t.Errorf("Expected err to be nil, got %v", err)
	}
}

func TestXRepository_Load4(t *testing.T) {
	rep := NewXRepository()

	// Test case 1: Successful read
	path := "D:/MMD/MikuMikuDance_v926x64/UserFile/Accessory/食べ物/お箸セット1 モノゾフ/みやこ箸.x"
	data, err := rep.Load(path)

	if err != nil {
		t.Errorf("Expected err to be nil, got %v", err)
	}
	if data == nil {
		t.Errorf("Expected model to be not nil, got nil")
	}
	model := data.(*pmx.PmxModel)

	pmxRep := NewPmxRepository()
	pmxRep.Save("../../../test_resources/test.pmx", model, false)

	pmxPath := strings.Replace(path, ".x", ".pmx", -1)
	expectedData, _ := pmxRep.Load(pmxPath)
	expectedModel := expectedData.(*pmx.PmxModel)

	for i, v := range model.Vertices.Data {
		expectedV := expectedModel.Vertices.Get(i)
		if !v.Position.NearEquals(expectedV.Position, 1e-5) {
			t.Errorf("Expected Position to be %v, got %v", expectedV.Position, v.Position)
		}
	}

	for i, f := range model.Faces.Data {
		expectedF := expectedModel.Faces.Get(i)
		if f.VertexIndexes[0] != expectedF.VertexIndexes[0] || f.VertexIndexes[1] != expectedF.VertexIndexes[1] || f.VertexIndexes[2] != expectedF.VertexIndexes[2] {
			t.Errorf("Expected Face[%d] VertexIndexes to be %v, got %v", f.Index(), expectedF.VertexIndexes, f.VertexIndexes)
		}
	}

	for i, m := range model.Materials.Data {
		expectedT := expectedModel.Materials.Get(i)
		if !m.Diffuse.NearEquals(expectedT.Diffuse, 1e-5) {
			t.Errorf("Expected Diffuse to be %v, got %v", expectedT.Diffuse, m.Diffuse)
		}
		if !m.Ambient.NearEquals(expectedT.Ambient, 1e-5) {
			t.Errorf("Expected Ambient to be %v, got %v", expectedT.Ambient, m.Ambient)
		}
		if !m.Specular.NearEquals(expectedT.Specular, 1e-5) {
			t.Errorf("Expected Specular to be %v, got %v", expectedT.Specular, m.Specular)
		}
		if !m.Edge.NearEquals(expectedT.Edge, 1e-5) {
			t.Errorf("Expected EdgeColor to be %v, got %v", expectedT.Edge, m.Edge)
		}
		if m.DrawFlag != expectedT.DrawFlag {
			t.Errorf("Expected DrawFlag to be %v, got %v", expectedT.DrawFlag, m.DrawFlag)
		}
		if m.EdgeSize != expectedT.EdgeSize {
			t.Errorf("Expected EdgeSize to be %v, got %v", expectedT.EdgeSize, m.EdgeSize)
		}
		if m.TextureIndex != expectedT.TextureIndex {
			t.Errorf("Expected TextureIndex to be %v, got %v", expectedT.TextureIndex, m.TextureIndex)
		}
		if m.SphereTextureIndex != expectedT.SphereTextureIndex {
			t.Errorf("Expected SphereTextureIndex to be %v, got %v", expectedT.SphereTextureIndex, m.SphereTextureIndex)
		}
		if m.ToonTextureIndex != expectedT.ToonTextureIndex {
			t.Errorf("Expected ToonTextureIndex to be %v, got %v", expectedT.ToonTextureIndex, m.ToonTextureIndex)
		}
		if m.SphereMode != expectedT.SphereMode {
			t.Errorf("Expected SphereMode to be %v, got %v", expectedT.SphereMode, m.SphereMode)
		}
		if m.ToonSharingFlag != expectedT.ToonSharingFlag {
			t.Errorf("Expected ToonSharingFlag to be %v, got %v", expectedT.ToonSharingFlag, m.ToonSharingFlag)
		}
		if m.VerticesCount != expectedT.VerticesCount {
			t.Errorf("Expected VerticesCount to be %v, got %v", expectedT.VerticesCount, m.VerticesCount)
		}
	}

}

func TestXRepository_Load5(t *testing.T) {
	rep := NewXRepository()

	// Test case 1: Successful read
	path := "D:/MMD/MikuMikuDance_v926x64/UserFile/Accessory/食べ物/お箸セット1 モノゾフ/丸金箔箸.x"
	data, err := rep.Load(path)

	if err != nil {
		t.Errorf("Expected err to be nil, got %v", err)
	}
	if data == nil {
		t.Errorf("Expected model to be not nil, got nil")
	}
	model := data.(*pmx.PmxModel)

	pmxRep := NewPmxRepository()
	pmxRep.Save("../../../test_resources/test.pmx", model, false)

	pmxPath := strings.Replace(path, ".x", ".pmx", -1)
	expectedData, _ := pmxRep.Load(pmxPath)
	expectedModel := expectedData.(*pmx.PmxModel)

	for i, v := range model.Vertices.Data {
		expectedV := expectedModel.Vertices.Get(i)
		if !v.Position.NearEquals(expectedV.Position, 1e-5) {
			t.Errorf("Expected Position to be %v, got %v", expectedV.Position, v.Position)
		}
	}

	for i, f := range model.Faces.Data {
		expectedF := expectedModel.Faces.Get(i)
		if f.VertexIndexes[0] != expectedF.VertexIndexes[0] || f.VertexIndexes[1] != expectedF.VertexIndexes[1] || f.VertexIndexes[2] != expectedF.VertexIndexes[2] {
			t.Errorf("Expected Face[%d] VertexIndexes to be %v, got %v", f.Index(), expectedF.VertexIndexes, f.VertexIndexes)
		}
	}

	for i, m := range model.Materials.Data {
		expectedT := expectedModel.Materials.Get(i)
		if !m.Diffuse.NearEquals(expectedT.Diffuse, 1e-5) {
			t.Errorf("Expected Diffuse to be %v, got %v", expectedT.Diffuse, m.Diffuse)
		}
		if !m.Ambient.NearEquals(expectedT.Ambient, 1e-5) {
			t.Errorf("Expected Ambient to be %v, got %v", expectedT.Ambient, m.Ambient)
		}
		if !m.Specular.NearEquals(expectedT.Specular, 1e-5) {
			t.Errorf("Expected Specular to be %v, got %v", expectedT.Specular, m.Specular)
		}
		if !m.Edge.NearEquals(expectedT.Edge, 1e-5) {
			t.Errorf("Expected EdgeColor to be %v, got %v", expectedT.Edge, m.Edge)
		}
		if m.DrawFlag != expectedT.DrawFlag {
			t.Errorf("Expected DrawFlag to be %v, got %v", expectedT.DrawFlag, m.DrawFlag)
		}
		if m.EdgeSize != expectedT.EdgeSize {
			t.Errorf("Expected EdgeSize to be %v, got %v", expectedT.EdgeSize, m.EdgeSize)
		}
		if m.TextureIndex != expectedT.TextureIndex {
			t.Errorf("Expected TextureIndex to be %v, got %v", expectedT.TextureIndex, m.TextureIndex)
		}
		if m.SphereTextureIndex != expectedT.SphereTextureIndex {
			t.Errorf("Expected SphereTextureIndex to be %v, got %v", expectedT.SphereTextureIndex, m.SphereTextureIndex)
		}
		if m.ToonTextureIndex != expectedT.ToonTextureIndex {
			t.Errorf("Expected ToonTextureIndex to be %v, got %v", expectedT.ToonTextureIndex, m.ToonTextureIndex)
		}
		if m.SphereMode != expectedT.SphereMode {
			t.Errorf("Expected SphereMode to be %v, got %v", expectedT.SphereMode, m.SphereMode)
		}
		if m.ToonSharingFlag != expectedT.ToonSharingFlag {
			t.Errorf("Expected ToonSharingFlag to be %v, got %v", expectedT.ToonSharingFlag, m.ToonSharingFlag)
		}
		if m.VerticesCount != expectedT.VerticesCount {
			t.Errorf("Expected VerticesCount to be %v, got %v", expectedT.VerticesCount, m.VerticesCount)
		}
	}

}

func TestXRepository_Load6(t *testing.T) {
	rep := NewXRepository()

	// Test case 1: Successful read
	path := "D:/MMD/MikuMikuDance_v926x64/UserFile/Accessory/食べ物/カレーライス アサシンP/カツカレー.x"
	data, err := rep.Load(path)

	if err != nil {
		t.Errorf("Expected err to be nil, got %v", err)
	}
	if data == nil {
		t.Errorf("Expected model to be not nil, got nil")
	}
	model := data.(*pmx.PmxModel)

	pmxRep := NewPmxRepository()
	pmxRep.Save("../../../test_resources/test.pmx", model, false)

	pmxPath := strings.Replace(path, ".x", ".pmx", -1)
	expectedData, _ := pmxRep.Load(pmxPath)
	expectedModel := expectedData.(*pmx.PmxModel)

	for i, v := range model.Vertices.Data {
		expectedV := expectedModel.Vertices.Get(i)
		if !v.Position.NearEquals(expectedV.Position, 1e-5) {
			t.Errorf("Expected Position to be %v, got %v", expectedV.Position, v.Position)
		}
	}

	for i, f := range model.Faces.Data {
		expectedF := expectedModel.Faces.Get(i)
		if f.VertexIndexes[0] != expectedF.VertexIndexes[0] || f.VertexIndexes[1] != expectedF.VertexIndexes[1] || f.VertexIndexes[2] != expectedF.VertexIndexes[2] {
			t.Errorf("Expected Face[%d] VertexIndexes to be %v, got %v", f.Index(), expectedF.VertexIndexes, f.VertexIndexes)
		}
	}

	for i, m := range model.Materials.Data {
		expectedT := expectedModel.Materials.Get(i)
		if !m.Diffuse.NearEquals(expectedT.Diffuse, 1e-5) {
			t.Errorf("Expected Diffuse to be %v, got %v", expectedT.Diffuse, m.Diffuse)
		}
		if !m.Ambient.NearEquals(expectedT.Ambient, 1e-5) {
			t.Errorf("Expected Ambient to be %v, got %v", expectedT.Ambient, m.Ambient)
		}
		if !m.Specular.NearEquals(expectedT.Specular, 1e-5) {
			t.Errorf("Expected Specular to be %v, got %v", expectedT.Specular, m.Specular)
		}
		if !m.Edge.NearEquals(expectedT.Edge, 1e-5) {
			t.Errorf("Expected EdgeColor to be %v, got %v", expectedT.Edge, m.Edge)
		}
		if m.DrawFlag != expectedT.DrawFlag {
			t.Errorf("Expected DrawFlag to be %v, got %v", expectedT.DrawFlag, m.DrawFlag)
		}
		if m.EdgeSize != expectedT.EdgeSize {
			t.Errorf("Expected EdgeSize to be %v, got %v", expectedT.EdgeSize, m.EdgeSize)
		}
		if m.TextureIndex != expectedT.TextureIndex {
			t.Errorf("Expected TextureIndex to be %v, got %v", expectedT.TextureIndex, m.TextureIndex)
		}
		if m.SphereTextureIndex != expectedT.SphereTextureIndex {
			t.Errorf("Expected SphereTextureIndex to be %v, got %v", expectedT.SphereTextureIndex, m.SphereTextureIndex)
		}
		if m.ToonTextureIndex != expectedT.ToonTextureIndex {
			t.Errorf("Expected ToonTextureIndex to be %v, got %v", expectedT.ToonTextureIndex, m.ToonTextureIndex)
		}
		if m.SphereMode != expectedT.SphereMode {
			t.Errorf("Expected SphereMode to be %v, got %v", expectedT.SphereMode, m.SphereMode)
		}
		if m.ToonSharingFlag != expectedT.ToonSharingFlag {
			t.Errorf("Expected ToonSharingFlag to be %v, got %v", expectedT.ToonSharingFlag, m.ToonSharingFlag)
		}
		if m.VerticesCount != expectedT.VerticesCount {
			t.Errorf("Expected VerticesCount to be %v, got %v", expectedT.VerticesCount, m.VerticesCount)
		}
	}

}
