// 指示: miu200521358
package domain

import (
	"reflect"
	"testing"

	"github.com/miu200521358/mlib_go/pkg/domain/model"
)

func TestExpandConnectedFaceIndexesRestrictsMaterialAndSharesReversedEdges(t *testing.T) {
	modelData := model.NewPmxModel()
	appendMaterial(modelData, 6) // 面 0, 1
	appendMaterial(modelData, 3) // 面 2
	appendFace(modelData, [3]int{0, 1, 2})
	appendFace(modelData, [3]int{3, 2, 1}) // 面 0 と辺 (1,2) を共有（逆順）
	appendFace(modelData, [3]int{1, 3, 4}) // 面 1 と共有するが材質が異なる

	got, err := ExpandConnectedFaceIndexes(modelData, []int{0})
	if err != nil {
		t.Fatalf("連続面拡張に失敗: %v", err)
	}
	if want := []int{0, 1}; !reflect.DeepEqual(got, want) {
		t.Fatalf("拡張結果=%v, want %v", got, want)
	}
}

func TestExpandConnectedFaceIndexesDeduplicatesAndSortsSeeds(t *testing.T) {
	modelData := model.NewPmxModel()
	appendMaterial(modelData, 9)
	appendFace(modelData, [3]int{0, 1, 2})
	appendFace(modelData, [3]int{2, 1, 3})
	appendFace(modelData, [3]int{10, 11, 12}) // 独立した面

	got, err := ExpandContinuousFaceIndexes(modelData, []int{2, 0, 2})
	if err != nil {
		t.Fatalf("連続面拡張に失敗: %v", err)
	}
	if want := []int{0, 1, 2}; !reflect.DeepEqual(got, want) {
		t.Fatalf("拡張結果=%v, want %v", got, want)
	}
}

func TestExpandConnectedFaceIndexesRejectsInvalidInput(t *testing.T) {
	tests := []struct {
		name string
		make func() *model.PmxModel
		seed []int
	}{
		{
			name: "seed 範囲外",
			make: func() *model.PmxModel {
				modelData := model.NewPmxModel()
				appendMaterial(modelData, 0)
				return modelData
			},
			seed: []int{0},
		},
		{
			name: "nil 面",
			make: func() *model.PmxModel {
				modelData := model.NewPmxModel()
				appendMaterial(modelData, 3)
				appendFace(modelData, [3]int{0, 1, 2})
				// IndexedCollection は nil を AppendRaw できないため、格納後に
				// 値を nil へ置換して破損入力を再現する。
				modelData.Faces.Values()[0] = nil
				return modelData
			},
			seed: []int{0},
		},
		{
			name: "負の頂点 index",
			make: func() *model.PmxModel {
				modelData := model.NewPmxModel()
				appendMaterial(modelData, 3)
				appendFace(modelData, [3]int{-1, 0, 1})
				return modelData
			},
			seed: []int{0},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := ExpandConnectedFaceIndexes(tt.make(), tt.seed); err == nil {
				t.Fatal("不正な入力がエラーになりません")
			}
		})
	}
}

// appendMaterial はテストモデルへ材質を追加する。
func appendMaterial(modelData *model.PmxModel, verticesCount int) {
	materialData := model.NewMaterial()
	materialData.VerticesCount = verticesCount
	modelData.Materials.AppendRaw(materialData)
}

// appendFace はテストモデルへ指定頂点の三角面を追加する。
func appendFace(modelData *model.PmxModel, vertices [3]int) {
	modelData.Faces.AppendRaw(&model.Face{VertexIndexes: vertices})
}
