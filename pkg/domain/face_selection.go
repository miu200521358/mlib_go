// 指示: miu200521358
package domain

import (
	"fmt"

	"github.com/miu200521358/mlib_go/pkg/domain/model"
)

// faceEdge は三角面の辺を頂点 index の昇順で表す。
// 頂点順が異なる隣接面でも同じキーになるよう正規化する。
type faceEdge struct {
	first  int
	second int
}

// ExpandConnectedFaceIndexes は選択済み面から辺共有で連続する面を拡張する。
//
// 拡張は同じ材質に属する面だけを対象とする。入力に複数材質の面が含まれる
// 場合も、各面を起点として材質ごとに独立した連結成分を探索する。戻り値は
// 重複を除いた面 index の昇順であり、seedFaceIndexes が空なら空スライスを返す。
// 材質面数の不整合、seed の範囲外、nil 面、負の頂点 index はエラーとする。
func ExpandConnectedFaceIndexes(modelData *model.PmxModel, seedFaceIndexes []int) ([]int, error) {
	if modelData == nil {
		return nil, fmt.Errorf("モデルが未設定です")
	}
	if modelData.Faces == nil {
		return nil, fmt.Errorf("面コレクションが未設定です")
	}
	faceMaterialIndex, err := model.NewFaceMaterialIndex(modelData)
	if err != nil {
		return nil, err
	}
	faceCount := modelData.Faces.Len()
	if len(seedFaceIndexes) == 0 {
		return []int{}, nil
	}

	selected := make(map[int]struct{}, len(seedFaceIndexes))
	queue := make([]int, 0, len(seedFaceIndexes))
	for _, faceIndex := range seedFaceIndexes {
		if faceIndex < 0 || faceIndex >= faceCount {
			return nil, fmt.Errorf("選択面 index が範囲外です: index=%d faces=%d", faceIndex, faceCount)
		}
		if _, exists := selected[faceIndex]; exists {
			continue
		}
		selected[faceIndex] = struct{}{}
		queue = append(queue, faceIndex)
	}

	// 辺から面を逆引きしておくことで、面数に対して線形の BFS になる。
	edges := make(map[faceEdge][]int, faceCount*3)
	for faceIndex, faceData := range modelData.Faces.Values() {
		if faceData == nil {
			return nil, fmt.Errorf("面が未設定です: index=%d", faceIndex)
		}
		vertices := faceData.VertexIndexes
		for _, vertexIndex := range vertices {
			if vertexIndex < 0 {
				return nil, fmt.Errorf("面の頂点 index が不正です: face=%d vertex=%d", faceIndex, vertexIndex)
			}
		}
		for edgeIndex := 0; edgeIndex < len(vertices); edgeIndex++ {
			first := vertices[edgeIndex]
			second := vertices[(edgeIndex+1)%len(vertices)]
			// 退化三角形の自己辺は隣接判定に使わない。
			if first == second {
				continue
			}
			edge := normalizeFaceEdge(first, second)
			edges[edge] = append(edges[edge], faceIndex)
		}
	}

	for head := 0; head < len(queue); head++ {
		faceIndex := queue[head]
		materialIndex, ok := faceMaterialIndex.MaterialIndex(faceIndex)
		if !ok {
			return nil, fmt.Errorf("面の材質 index を解決できません: face=%d", faceIndex)
		}
		vertices := modelData.Faces.Values()[faceIndex].VertexIndexes
		for edgeIndex := 0; edgeIndex < len(vertices); edgeIndex++ {
			first := vertices[edgeIndex]
			second := vertices[(edgeIndex+1)%len(vertices)]
			if first == second {
				continue
			}
			for _, neighborIndex := range edges[normalizeFaceEdge(first, second)] {
				if _, exists := selected[neighborIndex]; exists {
					continue
				}
				neighborMaterialIndex, ok := faceMaterialIndex.MaterialIndex(neighborIndex)
				if !ok || neighborMaterialIndex != materialIndex {
					continue
				}
				selected[neighborIndex] = struct{}{}
				queue = append(queue, neighborIndex)
			}
		}
	}

	result := make([]int, 0, len(selected))
	for faceIndex := 0; faceIndex < faceCount; faceIndex++ {
		if _, exists := selected[faceIndex]; exists {
			result = append(result, faceIndex)
		}
	}
	return result, nil
}

// ExpandContinuousFaceIndexes は ExpandConnectedFaceIndexes の別名である。
// 仕様書の「連続面」表現に合わせた呼び出し名を提供する。
func ExpandContinuousFaceIndexes(modelData *model.PmxModel, seedFaceIndexes []int) ([]int, error) {
	return ExpandConnectedFaceIndexes(modelData, seedFaceIndexes)
}

// normalizeFaceEdge は辺の端点順を正規化して返す。
func normalizeFaceEdge(first, second int) faceEdge {
	if first > second {
		first, second = second, first
	}
	return faceEdge{first: first, second: second}
}
