package repository

import (
	"strings"
	"testing"

	"github.com/miu200521358/mlib_go/pkg/config/mlog"
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

	model.Vertices.ForEach(func(index int, vertex *pmx.Vertex) {
		expectedV, _ := expectedModel.Vertices.Get(vertex.Index())
		if !vertex.Position.NearEquals(expectedV.Position, 1e-5) {
			t.Errorf("Expected Position to be %v, got %v", expectedV.Position, vertex.Position)
		}
	})

	model.Faces.ForEach(func(index int, face *pmx.Face) {
		expectedF, _ := expectedModel.Faces.Get(face.Index())
		if face.VertexIndexes[0] != expectedF.VertexIndexes[0] || face.VertexIndexes[1] != expectedF.VertexIndexes[1] || face.VertexIndexes[2] != expectedF.VertexIndexes[2] {
			t.Errorf("Expected VertexIndexes to be %v, got %v", expectedF.VertexIndexes, face.VertexIndexes)
		}
	})

	model.Materials.ForEach(func(index int, material *pmx.Material) {
		expectedT, _ := expectedModel.Materials.Get(material.Index())
		if !material.Diffuse.NearEquals(expectedT.Diffuse, 1e-5) {
			t.Errorf("Expected Diffuse to be %v, got %v", expectedT.Diffuse, material.Diffuse)
		}
		if !material.Ambient.NearEquals(expectedT.Ambient, 1e-5) {
			t.Errorf("Expected Ambient to be %v, got %v", expectedT.Ambient, material.Ambient)
		}
		if !material.Specular.NearEquals(expectedT.Specular, 1e-5) {
			t.Errorf("Expected Specular to be %v, got %v", expectedT.Specular, material.Specular)
		}
		if !material.Edge.NearEquals(expectedT.Edge, 1e-5) {
			t.Errorf("Expected EdgeColor to be %v, got %v", expectedT.Edge, material.Edge)
		}
		if material.DrawFlag != expectedT.DrawFlag {
			t.Errorf("Expected DrawFlag to be %v, got %v", expectedT.DrawFlag, material.DrawFlag)
		}
		if material.EdgeSize != expectedT.EdgeSize {
			t.Errorf("Expected EdgeSize to be %v, got %v", expectedT.EdgeSize, material.EdgeSize)
		}
		if material.TextureIndex != expectedT.TextureIndex {
			t.Errorf("Expected TextureIndex to be %v, got %v", expectedT.TextureIndex, material.TextureIndex)
		}
		if material.SphereTextureIndex != expectedT.SphereTextureIndex {
			t.Errorf("Expected SphereTextureIndex to be %v, got %v", expectedT.SphereTextureIndex, material.SphereTextureIndex)
		}
		if material.ToonTextureIndex != expectedT.ToonTextureIndex {
			t.Errorf("Expected ToonTextureIndex to be %v, got %v", expectedT.ToonTextureIndex, material.ToonTextureIndex)
		}
		if material.SphereMode != expectedT.SphereMode {
			t.Errorf("Expected SphereMode to be %v, got %v", expectedT.SphereMode, material.SphereMode)
		}
		if material.ToonSharingFlag != expectedT.ToonSharingFlag {
			t.Errorf("Expected ToonSharingFlag to be %v, got %v", expectedT.ToonSharingFlag, material.ToonSharingFlag)
		}
		if material.VerticesCount != expectedT.VerticesCount {
			t.Errorf("Expected VerticesCount to be %v, got %v", expectedT.VerticesCount, material.VerticesCount)
		}
	})

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

	model.Vertices.ForEach(func(index int, vertex *pmx.Vertex) {
		expectedV, _ := expectedModel.Vertices.Get(vertex.Index())
		if !vertex.Position.NearEquals(expectedV.Position, 1e-5) {
			t.Errorf("Expected Position to be %v, got %v", expectedV.Position, vertex.Position)
		}
	})

	model.Faces.ForEach(func(index int, face *pmx.Face) {
		expectedF, _ := expectedModel.Faces.Get(face.Index())
		if face.VertexIndexes[0] != expectedF.VertexIndexes[0] || face.VertexIndexes[1] != expectedF.VertexIndexes[1] || face.VertexIndexes[2] != expectedF.VertexIndexes[2] {
			t.Errorf("Expected Face[%d] VertexIndexes to be %v, got %v", face.Index(), expectedF.VertexIndexes, face.VertexIndexes)
		}
	})

	for _, is := range [][]int{{0, 0}, {1, 1}, {2, 2}, {3, 4}} {
		expectIndex, materialIndex := is[0], is[1]
		m, _ := model.Materials.Get(materialIndex)
		expectedT, _ := expectedModel.Materials.Get(expectIndex)
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

	model.Vertices.ForEach(func(index int, vertex *pmx.Vertex) {
		expectedV, _ := expectedModel.Vertices.Get(vertex.Index())
		if !vertex.Position.NearEquals(expectedV.Position, 1e-5) {
			t.Errorf("Expected Position to be %v, got %v", expectedV.Position, vertex.Position)
		}
	})

	model.Faces.ForEach(func(index int, face *pmx.Face) {
		expectedF, _ := expectedModel.Faces.Get(face.Index())
		if face.VertexIndexes[0] != expectedF.VertexIndexes[0] || face.VertexIndexes[1] != expectedF.VertexIndexes[1] || face.VertexIndexes[2] != expectedF.VertexIndexes[2] {
			t.Errorf("Expected Face[%d] VertexIndexes to be %v, got %v", face.Index(), expectedF.VertexIndexes, face.VertexIndexes)
		}
	})

	model.Materials.ForEach(func(index int, material *pmx.Material) {
		expectedT, _ := expectedModel.Materials.Get(material.Index())
		if !material.Diffuse.NearEquals(expectedT.Diffuse, 1e-5) {
			t.Errorf("Expected Diffuse to be %v, got %v", expectedT.Diffuse, material.Diffuse)
		}
		if !material.Ambient.NearEquals(expectedT.Ambient, 1e-5) {
			t.Errorf("Expected Ambient to be %v, got %v", expectedT.Ambient, material.Ambient)
		}
		if !material.Specular.NearEquals(expectedT.Specular, 1e-5) {
			t.Errorf("Expected Specular to be %v, got %v", expectedT.Specular, material.Specular)
		}
		if !material.Edge.NearEquals(expectedT.Edge, 1e-5) {
			t.Errorf("Expected EdgeColor to be %v, got %v", expectedT.Edge, material.Edge)
		}
		if material.DrawFlag != expectedT.DrawFlag {
			t.Errorf("Expected DrawFlag to be %v, got %v", expectedT.DrawFlag, material.DrawFlag)
		}
		if material.EdgeSize != expectedT.EdgeSize {
			t.Errorf("Expected EdgeSize to be %v, got %v", expectedT.EdgeSize, material.EdgeSize)
		}
		if material.TextureIndex != expectedT.TextureIndex {
			t.Errorf("Expected TextureIndex to be %v, got %v", expectedT.TextureIndex, material.TextureIndex)
		}
		if material.SphereTextureIndex != expectedT.SphereTextureIndex {
			t.Errorf("Expected SphereTextureIndex to be %v, got %v", expectedT.SphereTextureIndex, material.SphereTextureIndex)
		}
		if material.ToonTextureIndex != expectedT.ToonTextureIndex {
			t.Errorf("Expected ToonTextureIndex to be %v, got %v", expectedT.ToonTextureIndex, material.ToonTextureIndex)
		}
		if material.SphereMode != expectedT.SphereMode {
			t.Errorf("Expected SphereMode to be %v, got %v", expectedT.SphereMode, material.SphereMode)
		}
		if material.ToonSharingFlag != expectedT.ToonSharingFlag {
			t.Errorf("Expected ToonSharingFlag to be %v, got %v", expectedT.ToonSharingFlag, material.ToonSharingFlag)
		}
		if material.VerticesCount != expectedT.VerticesCount {
			t.Errorf("Expected VerticesCount to be %v, got %v", expectedT.VerticesCount, material.VerticesCount)
		}
	})
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

	model.Vertices.ForEach(func(index int, vertex *pmx.Vertex) {
		expectedV, _ := expectedModel.Vertices.Get(vertex.Index())
		if !vertex.Position.NearEquals(expectedV.Position, 1e-5) {
			t.Errorf("Expected Position to be %v, got %v", expectedV.Position, vertex.Position)
		}
	})

	model.Faces.ForEach(func(index int, face *pmx.Face) {
		expectedF, _ := expectedModel.Faces.Get(face.Index())
		if face.VertexIndexes[0] != expectedF.VertexIndexes[0] || face.VertexIndexes[1] != expectedF.VertexIndexes[1] || face.VertexIndexes[2] != expectedF.VertexIndexes[2] {
			t.Errorf("Expected Face[%d] VertexIndexes to be %v, got %v", face.Index(), expectedF.VertexIndexes, face.VertexIndexes)
		}
	})

	model.Materials.ForEach(func(index int, material *pmx.Material) {
		expectedT, _ := expectedModel.Materials.Get(material.Index())
		if !material.Diffuse.NearEquals(expectedT.Diffuse, 1e-5) {
			t.Errorf("Expected Diffuse to be %v, got %v", expectedT.Diffuse, material.Diffuse)
		}
		if !material.Ambient.NearEquals(expectedT.Ambient, 1e-5) {
			t.Errorf("Expected Ambient to be %v, got %v", expectedT.Ambient, material.Ambient)
		}
		if !material.Specular.NearEquals(expectedT.Specular, 1e-5) {
			t.Errorf("Expected Specular to be %v, got %v", expectedT.Specular, material.Specular)
		}
		if !material.Edge.NearEquals(expectedT.Edge, 1e-5) {
			t.Errorf("Expected EdgeColor to be %v, got %v", expectedT.Edge, material.Edge)
		}
		if material.DrawFlag != expectedT.DrawFlag {
			t.Errorf("Expected DrawFlag to be %v, got %v", expectedT.DrawFlag, material.DrawFlag)
		}
		if material.EdgeSize != expectedT.EdgeSize {
			t.Errorf("Expected EdgeSize to be %v, got %v", expectedT.EdgeSize, material.EdgeSize)
		}
		if material.TextureIndex != expectedT.TextureIndex {
			t.Errorf("Expected TextureIndex to be %v, got %v", expectedT.TextureIndex, material.TextureIndex)
		}
		if material.SphereTextureIndex != expectedT.SphereTextureIndex {
			t.Errorf("Expected SphereTextureIndex to be %v, got %v", expectedT.SphereTextureIndex, material.SphereTextureIndex)
		}
		if material.ToonTextureIndex != expectedT.ToonTextureIndex {
			t.Errorf("Expected ToonTextureIndex to be %v, got %v", expectedT.ToonTextureIndex, material.ToonTextureIndex)
		}
		if material.SphereMode != expectedT.SphereMode {
			t.Errorf("Expected SphereMode to be %v, got %v", expectedT.SphereMode, material.SphereMode)
		}
		if material.ToonSharingFlag != expectedT.ToonSharingFlag {
			t.Errorf("Expected ToonSharingFlag to be %v, got %v", expectedT.ToonSharingFlag, material.ToonSharingFlag)
		}
		if material.VerticesCount != expectedT.VerticesCount {
			t.Errorf("Expected VerticesCount to be %v, got %v", expectedT.VerticesCount, material.VerticesCount)
		}
	})

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

	model.Vertices.ForEach(func(index int, vertex *pmx.Vertex) {
		expectedV, _ := expectedModel.Vertices.Get(vertex.Index())
		if !vertex.Position.NearEquals(expectedV.Position, 1e-5) {
			t.Errorf("Expected Position to be %v, got %v", expectedV.Position, vertex.Position)
		}
	})

	model.Faces.ForEach(func(index int, face *pmx.Face) {
		expectedF, _ := expectedModel.Faces.Get(face.Index())
		if face.VertexIndexes[0] != expectedF.VertexIndexes[0] || face.VertexIndexes[1] != expectedF.VertexIndexes[1] || face.VertexIndexes[2] != expectedF.VertexIndexes[2] {
			t.Errorf("Expected Face[%d] VertexIndexes to be %v, got %v", face.Index(), expectedF.VertexIndexes, face.VertexIndexes)
		}
	})

	model.Materials.ForEach(func(index int, material *pmx.Material) {
		expectedT, _ := expectedModel.Materials.Get(material.Index())
		if !material.Diffuse.NearEquals(expectedT.Diffuse, 1e-5) {
			t.Errorf("Expected Diffuse to be %v, got %v", expectedT.Diffuse, material.Diffuse)
		}
		if !material.Ambient.NearEquals(expectedT.Ambient, 1e-5) {
			t.Errorf("Expected Ambient to be %v, got %v", expectedT.Ambient, material.Ambient)
		}
		if !material.Specular.NearEquals(expectedT.Specular, 1e-5) {
			t.Errorf("Expected Specular to be %v, got %v", expectedT.Specular, material.Specular)
		}
		if !material.Edge.NearEquals(expectedT.Edge, 1e-5) {
			t.Errorf("Expected EdgeColor to be %v, got %v", expectedT.Edge, material.Edge)
		}
		if material.DrawFlag != expectedT.DrawFlag {
			t.Errorf("Expected DrawFlag to be %v, got %v", expectedT.DrawFlag, material.DrawFlag)
		}
		if material.EdgeSize != expectedT.EdgeSize {
			t.Errorf("Expected EdgeSize to be %v, got %v", expectedT.EdgeSize, material.EdgeSize)
		}
		if material.TextureIndex != expectedT.TextureIndex {
			t.Errorf("Expected TextureIndex to be %v, got %v", expectedT.TextureIndex, material.TextureIndex)
		}
		if material.SphereTextureIndex != expectedT.SphereTextureIndex {
			t.Errorf("Expected SphereTextureIndex to be %v, got %v", expectedT.SphereTextureIndex, material.SphereTextureIndex)
		}
		if material.ToonTextureIndex != expectedT.ToonTextureIndex {
			t.Errorf("Expected ToonTextureIndex to be %v, got %v", expectedT.ToonTextureIndex, material.ToonTextureIndex)
		}
		if material.SphereMode != expectedT.SphereMode {
			t.Errorf("Expected SphereMode to be %v, got %v", expectedT.SphereMode, material.SphereMode)
		}
		if material.ToonSharingFlag != expectedT.ToonSharingFlag {
			t.Errorf("Expected ToonSharingFlag to be %v, got %v", expectedT.ToonSharingFlag, material.ToonSharingFlag)
		}
		if material.VerticesCount != expectedT.VerticesCount {
			t.Errorf("Expected VerticesCount to be %v, got %v", expectedT.VerticesCount, material.VerticesCount)
		}
	})

}

func TestXRepository_Load7(t *testing.T) {
	mlog.SetLevel(mlog.DEBUG)
	rep := NewXRepository()

	// Test case 1: Successful read
	path := "D:/MMD/MikuMikuDance_v926x64/UserFile/BackGround/ゆづき/桜切紙ステージ/桜切紙ステージ/桜切紙ステージ.x"
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

	if model.Vertices.Length() != expectedModel.Vertices.Length() {
		t.Errorf("Expected Vertices Count to be %v, got %v", expectedModel.Vertices.Length(), model.Vertices.Length())
	}

	model.Vertices.ForEach(func(index int, vertex *pmx.Vertex) {
		expectedV, _ := expectedModel.Vertices.Get(vertex.Index())
		if !vertex.Position.NearEquals(expectedV.Position, 1e-5) {
			t.Errorf("Expected Position to be %v, got %v", expectedV.Position, vertex.Position)
		}
	})

	if model.Faces.Length() != expectedModel.Faces.Length() {
		t.Errorf("Expected Faces Count to be %v, got %v", expectedModel.Faces.Length(), model.Faces.Length())
	}

	model.Faces.ForEach(func(index int, face *pmx.Face) {
		expectedF, _ := expectedModel.Faces.Get(face.Index())
		if face.VertexIndexes[0] != expectedF.VertexIndexes[0] || face.VertexIndexes[1] != expectedF.VertexIndexes[1] || face.VertexIndexes[2] != expectedF.VertexIndexes[2] {
			t.Errorf("Expected Face[%d] VertexIndexes to be %v, got %v", face.Index(), expectedF.VertexIndexes, face.VertexIndexes)
		}
	})

	if model.Materials.Length() != expectedModel.Materials.Length() {
		t.Errorf("Expected Materials Count to be %v, got %v", expectedModel.Materials.Length(), model.Materials.Length())
	}

	model.Materials.ForEach(func(index int, material *pmx.Material) {
		expectedT, _ := expectedModel.Materials.Get(material.Index())
		if !material.Diffuse.NearEquals(expectedT.Diffuse, 1e-5) {
			t.Errorf("Expected Diffuse to be %v, got %v", expectedT.Diffuse, material.Diffuse)
		}
		if !material.Ambient.NearEquals(expectedT.Ambient, 1e-5) {
			t.Errorf("Expected Ambient to be %v, got %v", expectedT.Ambient, material.Ambient)
		}
		if !material.Specular.NearEquals(expectedT.Specular, 1e-5) {
			t.Errorf("Expected Specular to be %v, got %v", expectedT.Specular, material.Specular)
		}
		if !material.Edge.NearEquals(expectedT.Edge, 1e-5) {
			t.Errorf("Expected EdgeColor to be %v, got %v", expectedT.Edge, material.Edge)
		}
		if material.DrawFlag != expectedT.DrawFlag {
			t.Errorf("Expected DrawFlag to be %v, got %v", expectedT.DrawFlag, material.DrawFlag)
		}
		if material.EdgeSize != expectedT.EdgeSize {
			t.Errorf("Expected EdgeSize to be %v, got %v", expectedT.EdgeSize, material.EdgeSize)
		}
		if material.TextureIndex != expectedT.TextureIndex {
			t.Errorf("Expected TextureIndex to be %v, got %v", expectedT.TextureIndex, material.TextureIndex)
		}
		if material.SphereTextureIndex != expectedT.SphereTextureIndex {
			t.Errorf("Expected SphereTextureIndex to be %v, got %v", expectedT.SphereTextureIndex, material.SphereTextureIndex)
		}
		if material.ToonTextureIndex != expectedT.ToonTextureIndex {
			t.Errorf("Expected ToonTextureIndex to be %v, got %v", expectedT.ToonTextureIndex, material.ToonTextureIndex)
		}
		if material.SphereMode != expectedT.SphereMode {
			t.Errorf("Expected SphereMode to be %v, got %v", expectedT.SphereMode, material.SphereMode)
		}
		if material.ToonSharingFlag != expectedT.ToonSharingFlag {
			t.Errorf("Expected ToonSharingFlag to be %v, got %v", expectedT.ToonSharingFlag, material.ToonSharingFlag)
		}
		if material.VerticesCount != expectedT.VerticesCount {
			t.Errorf("Expected VerticesCount to be %v, got %v", expectedT.VerticesCount, material.VerticesCount)
		}
	})

}
