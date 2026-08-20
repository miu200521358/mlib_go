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

func TestDeselectConnectedFaceIndexesRemovesSeedComponents(t *testing.T) {
	modelData := model.NewPmxModel()
	appendMaterial(modelData, 9) // 面 0, 1, 2
	appendMaterial(modelData, 3) // 面 3
	appendFace(modelData, [3]int{0, 1, 2})
	appendFace(modelData, [3]int{3, 2, 1}) // 面 0 と同一材質で辺共有
	appendFace(modelData, [3]int{10, 11, 12})
	appendFace(modelData, [3]int{2, 1, 13}) // 面 0 と辺共有するが材質が異なる

	selected := []int{0, 0}
	got, err := DeselectConnectedFaceIndexes(modelData, selected)
	if err != nil {
		t.Fatalf("連続選択面解除に失敗: %v", err)
	}
	if want := []int{}; !reflect.DeepEqual(got, want) {
		t.Fatalf("解除結果=%v, want %v", got, want)
	}
	if want := []int{0, 0}; !reflect.DeepEqual(selected, want) {
		t.Fatalf("入力選択面が変更された: got=%v, want=%v", selected, want)
	}
}

func TestDeselectConnectedFaceIndexesRemovesMultipleMaterialComponents(t *testing.T) {
	modelData := model.NewPmxModel()
	appendMaterial(modelData, 9) // 面 0, 1, 2
	appendMaterial(modelData, 3) // 面 3
	appendFace(modelData, [3]int{0, 1, 2})
	appendFace(modelData, [3]int{3, 2, 1})
	appendFace(modelData, [3]int{3, 4, 5})
	appendFace(modelData, [3]int{2, 1, 6}) // 面 0 と辺共有だが別材質

	connected, err := ExpandConnectedFaceIndexes(modelData, []int{0})
	if err != nil {
		t.Fatalf("連続面の確認に失敗: %v", err)
	}
	if want := []int{0, 1}; !reflect.DeepEqual(connected, want) {
		t.Fatalf("連続面=%v, want %v", connected, want)
	}

	got, err := DeselectConnectedFaceIndexes(modelData, []int{0, 1, 2})
	if err != nil {
		t.Fatalf("複数材質の連続選択面解除に失敗: %v", err)
	}
	if want := []int{}; !reflect.DeepEqual(got, want) {
		t.Fatalf("複数材質の解除結果=%v, want %v", got, want)
	}
}

func TestDeselectConnectedFaceIndexesRejectsInvalidInput(t *testing.T) {
	tests := []struct {
		name  string
		model *model.PmxModel
		seed  []int
	}{
		{
			name:  "nil モデル",
			model: nil,
			seed:  []int{0},
		},
		{
			name: "seed 範囲外",
			model: func() *model.PmxModel {
				modelData := model.NewPmxModel()
				appendMaterial(modelData, 0)
				return modelData
			}(),
			seed: []int{0},
		},
		{
			name: "負の頂点 index",
			model: func() *model.PmxModel {
				modelData := model.NewPmxModel()
				appendMaterial(modelData, 3)
				appendFace(modelData, [3]int{-1, 0, 1})
				return modelData
			}(),
			seed: []int{0},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := DeselectConnectedFaceIndexes(tt.model, tt.seed); err == nil {
				t.Fatal("不正な入力がエラーになりません")
			}
		})
	}
}

func TestDeselectConnectedFaceIndexesEmptySelection(t *testing.T) {
	modelData := model.NewPmxModel()
	appendMaterial(modelData, 0)

	got, err := DeselectConnectedFaceIndexes(modelData, nil)
	if err != nil {
		t.Fatalf("空の選択面解除に失敗: %v", err)
	}
	if got == nil || len(got) != 0 {
		t.Fatalf("空の解除結果=%v, want 非 nil の空スライス", got)
	}
}

func TestRemoveConnectedFaceIndexesAlias(t *testing.T) {
	modelData := model.NewPmxModel()
	appendMaterial(modelData, 3)
	appendFace(modelData, [3]int{0, 1, 2})

	got, err := RemoveConnectedFaceIndexes(modelData, []int{0})
	if err != nil {
		t.Fatalf("別名 API の連続選択面解除に失敗: %v", err)
	}
	if len(got) != 0 {
		t.Fatalf("別名 API の解除結果=%v, want []", got)
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
