package vmd

import (
	"math"

	"github.com/miu200521358/mlib_go/pkg/mmath"
)

type VertexDelta struct {
	Position *mmath.MVec3
}

func NewVertexDelta(pos *mmath.MVec3) *VertexDelta {
	return &VertexDelta{
		Position: pos,
	}
}

type VertexDeltas struct {
	Data     map[int]*VertexDelta
	IndexMap map[mmath.MVec3]map[int]*VertexDelta
}

func NewVertexDeltas() *VertexDeltas {
	return &VertexDeltas{
		Data: make(map[int]*VertexDelta),
	}
}

func (vds *VertexDeltas) SetupMapKeys() {
	vds.IndexMap = make(map[mmath.MVec3]map[int]*VertexDelta)
	for k, v := range vds.Data {
		baseKey := v.GetMapKey()
		// 前後のオフセット込みでマッピング
		for _, offset := range []*mmath.MVec3{
			{0, 0, 0}, {1, 0, 0}, {0, 1, 0}, {0, 0, 1},
			{0, 0, 0}, {-1, 0, 0}, {0, -1, 0}, {0, 0, -1},
			{1, 1, 0}, {1, 0, 1}, {0, 1, 1}, {1, 1, 1},
			{-1, -1, 0}, {-1, 0, -1}, {0, -1, -1}, {-1, -1, -1},
			{1, -1, 0}, {1, 0, -1}, {0, 1, -1}, {1, -1, 1},
			{-1, 1, 0}, {-1, 0, 1}, {0, -1, 1}, {-1, 1, -1},
		} {
			key := *baseKey.Added(offset)
			if _, ok := vds.IndexMap[key]; !ok {
				vds.IndexMap[key] = make(map[int]*VertexDelta)
			}
			vds.IndexMap[key][k] = v
		}
	}
}

func (vds *VertexDeltas) GetMapValues(v *VertexDelta) ([]int, []*mmath.MVec3) {
	if vds.Data == nil {
		return nil, nil
	}
	key := v.GetMapKey()
	indexes := make([]int, 0)
	values := make([]*mmath.MVec3, 0)
	if mapIndexes, ok := vds.IndexMap[key]; ok {
		for i, iv := range mapIndexes {
			indexes = append(indexes, i)
			values = append(values, iv.Position)
		}
		return indexes, values
	}
	return nil, nil
}

func (vd *VertexDelta) GetMapKey() mmath.MVec3 {
	return mmath.MVec3{math.Round(vd.Position.GetX()), math.Round(vd.Position.GetY()), math.Round(vd.Position.GetZ())}
}
