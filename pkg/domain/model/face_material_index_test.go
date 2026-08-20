// 指示: miu200521358
package model

import "testing"

func TestNewFaceMaterialIndexBuildsBidirectionalRanges(t *testing.T) {
	modelData := NewPmxModel()
	modelData.Materials.AppendRaw(materialWithVerticesCount(6))
	modelData.Materials.AppendRaw(materialWithVerticesCount(3))
	for i := 0; i < 3; i++ {
		modelData.Faces.AppendRaw(&Face{VertexIndexes: [3]int{i, i + 1, i + 2}})
	}

	index, err := NewFaceMaterialIndex(modelData)
	if err != nil {
		t.Fatalf("索引構築に失敗: %v", err)
	}
	if got, want := index.FaceCount(), 3; got != want {
		t.Fatalf("FaceCount=%d, want %d", got, want)
	}
	if got, want := index.MaterialCount(), 2; got != want {
		t.Fatalf("MaterialCount=%d, want %d", got, want)
	}
	for faceIndex, wantMaterial := range []int{0, 0, 1} {
		if got, ok := index.MaterialIndex(faceIndex); !ok || got != wantMaterial {
			t.Errorf("MaterialIndex(%d)=(%d,%v), want (%d,true)", faceIndex, got, ok, wantMaterial)
		}
	}
	if got, ok := index.FaceRange(0); !ok || got != (FaceIndexRange{Start: 0, End: 2}) {
		t.Errorf("FaceRange(0)=(%+v,%v)", got, ok)
	}
	if got, ok := index.FaceRange(1); !ok || got != (FaceIndexRange{Start: 2, End: 3}) {
		t.Errorf("FaceRange(1)=(%+v,%v)", got, ok)
	}
}

func TestNewFaceMaterialIndexRejectsInconsistentCounts(t *testing.T) {
	tests := []struct {
		name          string
		verticesCount int
		faceCount     int
	}{
		{name: "負数", verticesCount: -3, faceCount: 0},
		{name: "三の倍数でない", verticesCount: 1, faceCount: 0},
		{name: "面数不足", verticesCount: 6, faceCount: 1},
		{name: "面数余り", verticesCount: 0, faceCount: 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			modelData := NewPmxModel()
			modelData.Materials.AppendRaw(materialWithVerticesCount(tt.verticesCount))
			for i := 0; i < tt.faceCount; i++ {
				modelData.Faces.AppendRaw(&Face{})
			}
			if _, err := NewFaceMaterialIndex(modelData); err == nil {
				t.Fatal("不整合した面索引がエラーになりません")
			}
		})
	}
}

func TestFaceMaterialIndexOutOfRange(t *testing.T) {
	modelData := NewPmxModel()
	modelData.Materials.AppendRaw(materialWithVerticesCount(0))
	index, err := NewFaceMaterialIndex(modelData)
	if err != nil {
		t.Fatalf("索引構築に失敗: %v", err)
	}
	if got, ok := index.MaterialIndex(-1); ok || got != -1 {
		t.Errorf("MaterialIndex(-1)=(%d,%v)", got, ok)
	}
	if _, ok := index.MaterialIndex(0); ok {
		t.Error("空の面索引で MaterialIndex(0) が有効です")
	}
	if _, ok := index.FaceRange(-1); ok {
		t.Error("FaceRange(-1) が有効です")
	}
	if _, ok := index.FaceRange(1); ok {
		t.Error("FaceRange(1) が有効です")
	}
}

// materialWithVerticesCount は指定した頂点数のテスト用材質を生成する。
func materialWithVerticesCount(verticesCount int) *Material {
	materialData := NewMaterial()
	materialData.VerticesCount = verticesCount
	return materialData
}
