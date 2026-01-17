// 指示: miu200521358
package pmx

import (
	"sort"

	"github.com/miu200521358/mlib_go/pkg/domain/model"
)

// indexMapping は旧新インデックスの対応を表す。
type indexMapping struct {
	oldToNew []int
	newToOld []int
}

// buildIndexMapping は対象のインデックス対応を生成する。
func buildIndexMapping(total int, include func(index int) bool) indexMapping {
	oldToNew := make([]int, total)
	newToOld := make([]int, 0, total)
	for i := 0; i < total; i++ {
		oldToNew[i] = -1
	}
	for i := 0; i < total; i++ {
		if include == nil || include(i) {
			oldToNew[i] = len(newToOld)
			newToOld = append(newToOld, i)
		}
	}
	return indexMapping{oldToNew: oldToNew, newToOld: newToOld}
}

// mapIndex は旧インデックスを新インデックスへ変換する。
func (m indexMapping) mapIndex(index int) int {
	if index < 0 || index >= len(m.oldToNew) {
		return -1
	}
	return m.oldToNew[index]
}

// compressLayers はレイヤー値を0..N-1へ圧縮する。
func compressLayers(bones []*model.Bone, mapping indexMapping) map[int]int {
	layerSet := make(map[int]struct{})
	for _, oldIndex := range mapping.newToOld {
		if oldIndex < 0 || oldIndex >= len(bones) {
			continue
		}
		bone := bones[oldIndex]
		if bone == nil {
			continue
		}
		layerSet[bone.Layer] = struct{}{}
	}
	layers := make([]int, 0, len(layerSet))
	for layer := range layerSet {
		layers = append(layers, layer)
	}
	sort.Ints(layers)
	mapped := make(map[int]int, len(layers))
	for i, layer := range layers {
		mapped[layer] = i
	}
	return mapped
}

// defineVertexIndexSize は頂点インデックスのサイズを返す。
func defineVertexIndexSize(count int) byte {
	switch {
	case count < 256:
		return 1
	case count <= 65535:
		return 2
	default:
		return 4
	}
}

// defineOtherIndexSize は頂点以外のインデックスサイズを返す。
func defineOtherIndexSize(count int) byte {
	switch {
	case count < 128:
		return 1
	case count <= 32767:
		return 2
	default:
		return 4
	}
}
