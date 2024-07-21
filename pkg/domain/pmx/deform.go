package pmx

import (
	"sort"

	"github.com/miu200521358/mlib_go/pkg/domain/mmath"
)

// DeformType ウェイト変形方式
type DeformType byte

const (
	BDEF1 DeformType = 0
	BDEF2 DeformType = 1
	BDEF4 DeformType = 2
	SDEF  DeformType = 3
)

type IDeform interface {
	GetType() DeformType
	AllIndexes() []int
	AllWeights() []float64
	Indexes(weightThreshold float64) []int
	Weights(weightThreshold float64) []float64
	NormalizedDeform() [8]float32
}

// Deform デフォーム既定構造体
type Deform struct {
	indexes []int     // ボーンINDEXリスト
	weights []float64 // ウェイトリスト
	Count   int       // デフォームボーン個数
}

// NewDeform creates a new Deform instance.
func NewDeform(indexes []int, weights []float64, count int) Deform {
	return Deform{
		indexes: indexes,
		weights: weights,
		Count:   count,
	}
}

func (deform *Deform) AllIndexes() []int {
	return deform.indexes
}

func (deform *Deform) AllWeights() []float64 {
	return deform.weights
}

// Indexes ウェイト閾値以上のウェイトを持っているINDEXのみを取得する
func (deform *Deform) Indexes(weightThreshold float64) []int {
	var indexes []int
	for i, weight := range deform.weights {
		if weight >= weightThreshold {
			indexes = append(indexes, deform.indexes[i])
		}
	}
	return indexes
}

// Weights ウェイト閾値以上のウェイトを持っているウェイトのみを取得する
func (deform *Deform) Weights(weightThreshold float64) []float64 {
	var weights []float64
	for _, weight := range deform.weights {
		if weight >= weightThreshold {
			weights = append(weights, weight)
		}
	}
	return weights
}

// Normalize ウェイト正規化
func (deform *Deform) Normalize(align bool) {
	if align {
		// ウェイトを統合する
		indexWeights := make(map[int]float64)
		for i, index := range deform.indexes {
			if _, ok := indexWeights[index]; !ok {
				indexWeights[index] = 0.0
			}
			indexWeights[index] += deform.weights[i]
		}

		// 揃える必要がある場合、数が足りるよう、かさ増しする
		ilist := make([]int, 0, len(indexWeights)+4)
		wlist := make([]float64, 0, len(indexWeights)+4)
		for index, weight := range indexWeights {
			ilist = append(ilist, index)
			wlist = append(wlist, weight)
		}
		for i := len(indexWeights); i < deform.Count; i++ {
			ilist = append(ilist, 0)
			wlist = append(wlist, 0)
		}

		// 正規化
		sum := 0.0
		for _, weight := range wlist {
			sum += weight
		}
		for i := range wlist {
			wlist[i] /= sum
		}

		// ウェイトの大きい順に指定個数までを対象とする
		deform.indexes, deform.weights = sortIndexesByWeight(ilist, wlist)
	}

	// ウェイト正規化
	sum := 0.0
	for _, weight := range deform.weights {
		sum += weight
	}
	for i := range deform.weights {
		deform.weights[i] /= sum
	}
}

// NormalizedDeform 4つのボーンINDEXとウェイトを返す（合計8個）
func (deform *Deform) NormalizedDeform() [8]float32 {
	normalizedDeform := [8]float32{0, 0, 0, 0, 0, 0, 0, 0}
	for i, index := range deform.indexes {
		normalizedDeform[i] = float32(index)
	}
	for i, weight := range deform.weights {
		normalizedDeform[i+4] = float32(weight)
	}

	return normalizedDeform
}

// sortIndexesByWeight ウェイトの大きい順に指定個数までを対象とする
func sortIndexesByWeight(indexes []int, weights []float64) ([]int, []float64) {
	sort.SliceStable(weights, func(i, j int) bool {
		return weights[i] > weights[j]
	})

	sortedIndexes := make([]int, len(indexes))
	sortedWeights := make([]float64, len(weights))

	for i, weight := range weights {
		for j, w := range weights {
			if weight == w {
				sortedIndexes[i] = indexes[j]
				sortedWeights[i] = w
				break
			}
		}
	}

	return sortedIndexes, sortedWeights
}

// Bdef1 represents the BDEF1 deformation.
type Bdef1 struct {
	Deform
}

// NewBdef1 creates a new Bdef1 instance.
func NewBdef1(index0 int) *Bdef1 {
	return &Bdef1{
		Deform: Deform{
			indexes: []int{index0},
			weights: []float64{1.0},
			Count:   1,
		},
	}
}

// GetType returns the deformation type.
func (bdef1 *Bdef1) GetType() DeformType {
	return BDEF1
}

// NormalizedDeform 4つのボーンINDEXとウェイトを返す（合計8個）
func (bdef1 *Bdef1) NormalizedDeform() [8]float32 {
	return [8]float32{float32(bdef1.indexes[0]), 0, 0, 0, 1.0, 0, 0, 0}
}

// Bdef2 represents the BDEF2 deformation.
type Bdef2 struct {
	Deform
}

// NewBdef2 creates a new Bdef2 instance.
func NewBdef2(index0, index1 int, weight0 float64) *Bdef2 {
	return &Bdef2{
		Deform: Deform{
			indexes: []int{index0, index1},
			weights: []float64{weight0, 1 - weight0},
			Count:   2,
		},
	}
}

// GetType returns the deformation type.
func (bdef2 *Bdef2) GetType() DeformType {
	return BDEF2
}

// NormalizedDeform 4つのボーンINDEXとウェイトを返す（合計8個）
func (bdef2 *Bdef2) NormalizedDeform() [8]float32 {
	return [8]float32{
		float32(bdef2.indexes[0]), float32(bdef2.indexes[1]), 0, 0,
		float32(bdef2.weights[0]), float32(1 - bdef2.weights[0]), 0, 0}
}

// Bdef4 represents the BDEF4 deformation.
type Bdef4 struct {
	Deform
}

// NewBdef4 creates a new Bdef4 instance.
func NewBdef4(index0, index1, index2, index3 int, weight0, weight1, weight2, weight3 float64) *Bdef4 {
	return &Bdef4{
		Deform: Deform{
			indexes: []int{index0, index1, index2, index3},
			weights: []float64{weight0, weight1, weight2, weight3},
			Count:   4,
		},
	}
}

// GetType returns the deformation type.
func (bdef4 *Bdef4) GetType() DeformType {
	return BDEF4
}

// NormalizedDeform 4つのボーンINDEXとウェイトを返す（合計8個）
func (bdef4 *Bdef4) NormalizedDeform() [8]float32 {
	return [8]float32{
		float32(bdef4.indexes[0]), float32(bdef4.indexes[1]), float32(bdef4.indexes[2]), float32(bdef4.indexes[3]),
		float32(bdef4.weights[0]), float32(bdef4.weights[1]), float32(bdef4.weights[2]), float32(bdef4.weights[3])}
}

// Sdef represents the SDEF deformation.
type Sdef struct {
	Deform
	SdefC  *mmath.MVec3
	SdefR0 *mmath.MVec3
	SdefR1 *mmath.MVec3
}

// NewSdef creates a new Sdef instance.
func NewSdef(index0, index1 int, weight0 float64, sdefC, sdefR0, sdefR1 *mmath.MVec3) *Sdef {
	return &Sdef{
		Deform: Deform{
			indexes: []int{index0, index1},
			weights: []float64{weight0, 1 - weight0},
			Count:   2,
		},
		SdefC:  sdefC,
		SdefR0: sdefR0,
		SdefR1: sdefR1,
	}
}

// GetType returns the deformation type.
func (sdef *Sdef) GetType() DeformType {
	return SDEF
}

// NormalizedDeform 4つのボーンINDEXとウェイトを返す（合計8個）
// TODO: SDEFパラメーターの正規化
func (sdef *Sdef) NormalizedDeform() [8]float32 {
	return [8]float32{
		float32(sdef.indexes[0]), float32(sdef.indexes[1]), 0, 0,
		float32(sdef.weights[0]), float32(1 - sdef.weights[0]), 0, 0}
}
