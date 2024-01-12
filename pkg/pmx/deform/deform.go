package deform

import (
	"sort"

	"github.com/miu200521358/mlib_go/pkg/math/mvec3"

)

// DeformType ウェイト変形方式
type DeformType int

const (
	BDEF1 DeformType = 0
	BDEF2 DeformType = 1
	BDEF4 DeformType = 2
	SDEF  DeformType = 3
)

// T デフォーム既定構造体
type T struct {
	Indexes []int
	Weights []float64
	Count   int
}

// NewDeform creates a new Deform instance.
func NewDeform(indexes []int, weights []float64, count int) *T {
	return &T{
		Indexes: indexes,
		Weights: weights,
		Count:   count,
	}
}

// GetIndexes ウェイト閾値以上のウェイトを持っているINDEXのみを取得する
func (d *T) GetIndexes(weightThreshold float64) []int {
	var indexes []int
	for i, weight := range d.Weights {
		if weight >= weightThreshold {
			indexes = append(indexes, d.Indexes[i])
		}
	}
	return indexes
}

// GetWeights ウェイト閾値以上のウェイトを持っているウェイトのみを取得する
func (d *T) GetWeights(weightThreshold float64) []float64 {
	var weights []float64
	for _, weight := range d.Weights {
		if weight >= weightThreshold {
			weights = append(weights, weight)
		}
	}
	return weights
}

// Normalize ウェイト正規化
func (d *T) Normalize(align bool) {
	if align {
		// ウェイトを統合する
		indexWeights := make(map[int]float64)
		for i, index := range d.Indexes {
			if _, ok := indexWeights[index]; !ok {
				indexWeights[index] = 0.0
			}
			indexWeights[index] += d.Weights[i]
		}

		// 揃える必要がある場合、数が足りるよう、かさ増しする
		ilist := make([]int, 0, len(indexWeights)+4)
		wlist := make([]float64, 0, len(indexWeights)+4)
		for index, weight := range indexWeights {
			ilist = append(ilist, index)
			wlist = append(wlist, weight)
		}
		for i := len(indexWeights); i < d.Count; i++ {
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
		sortIndexesByWeight(ilist, wlist)
		d.Indexes = ilist[:d.Count]
		d.Weights = wlist[:d.Count]
	}

	// ウェイト正規化
	sum := 0.0
	for _, weight := range d.Weights {
		sum += weight
	}
	for i := range d.Weights {
		d.Weights[i] /= sum
	}
}

// NormalizedDeform ウェイト正規化して4つのボーンINDEXとウェイトを返す（合計8個）
func (d *T) NormalizedDeform() []float64 {
	// 揃える必要がある場合、ウェイトを統合する
	indexWeights := make(map[int]float64)
	for i, index := range d.Indexes {
		if _, ok := indexWeights[index]; !ok {
			indexWeights[index] = 0.0
		}
		indexWeights[index] += d.Weights[i]
	}

	// 揃える必要がある場合、数が足りるよう、かさ増しする
	ilist := make([]int, 0, len(indexWeights)+4)
	wlist := make([]float64, 0, len(indexWeights)+4)
	for index, weight := range indexWeights {
		ilist = append(ilist, index)
		wlist = append(wlist, weight)
	}
	for i := len(indexWeights); i < 4; i++ {
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
	sortIndexesByWeight(ilist, wlist)
	ilist = ilist[:4]
	wlist = wlist[:4]

	// ウェイト正規化
	sum = 0.0
	for _, weight := range wlist {
		sum += weight
	}
	for i := range wlist {
		wlist[i] /= sum
	}

	normalizedDeform := make([]float64, 0, 8)
	normalizedDeform = append(normalizedDeform, float64(ilist[0]), float64(ilist[1]), float64(ilist[2]), float64(ilist[3]))
	normalizedDeform = append(normalizedDeform, wlist[0], wlist[1], wlist[2], wlist[3])

	return normalizedDeform
}

// sortIndexesByWeight sorts the indexes and weights by weight in descending order.
func sortIndexesByWeight(indexes []int, weights []float64) {
	sortFunc := func(i, j int) bool {
		return weights[i] > weights[j]
	}
	sortIndexes := make([]int, len(indexes))
	for i := range sortIndexes {
		sortIndexes[i] = i
	}
	sort.Slice(sortIndexes, sortFunc)
	for i := range indexes {
		indexes[i] = indexes[sortIndexes[i]]
		weights[i] = weights[sortIndexes[i]]
	}
}

// Bdef1 represents the BDEF1 deformation.
type Bdef1 struct {
	T
}

// NewBdef1 creates a new Bdef1 instance.
func NewBdef1(index0 int) *Bdef1 {
	return &Bdef1{
		T: T{
			Indexes: []int{index0},
			Weights: []float64{1.0},
			Count:   1,
		},
	}
}

// Type returns the deformation type.
func (b *Bdef1) Type() DeformType {
	return BDEF1
}

// Bdef2 represents the BDEF2 deformation.
type Bdef2 struct {
	T
}

// NewBdef2 creates a new Bdef2 instance.
func NewBdef2(index0, index1 int, weight0 float64) *Bdef2 {
	return &Bdef2{
		T: T{
			Indexes: []int{index0, index1},
			Weights: []float64{weight0, 1 - weight0},
			Count:   2,
		},
	}
}

// Type returns the deformation type.
func (b *Bdef2) Type() DeformType {
	return BDEF2
}

// Bdef4 represents the BDEF4 deformation.
type Bdef4 struct {
	T
}

// NewBdef4 creates a new Bdef4 instance.
func NewBdef4(index0, index1, index2, index3 int, weight0, weight1, weight2, weight3 float64) *Bdef4 {
	return &Bdef4{
		T: T{
			Indexes: []int{index0, index1, index2, index3},
			Weights: []float64{weight0, weight1, weight2, weight3},
			Count:   4,
		},
	}
}

// Type returns the deformation type.
func (b *Bdef4) Type() DeformType {
	return BDEF4
}

// Sdef represents the SDEF deformation.
type Sdef struct {
	T
	SdefC  *mvec3.T
	SdefR0 *mvec3.T
	SdefR1 *mvec3.T
}

// NewSdef creates a new Sdef instance.
func NewSdef(index0, index1 int, weight0, sdefCX, sdefCY, sdefCZ, sdefR0X, sdefR0Y, sdefR0Z, sdefR1X, sdefR1Y, sdefR1Z float64) *Sdef {
	return &Sdef{
		T: T{
			Indexes: []int{index0, index1},
			Weights: []float64{weight0, 1 - weight0},
			Count:   2,
		},
		SdefC:  &mvec3.T{sdefCX, sdefCY, sdefCZ},
		SdefR0: &mvec3.T{sdefR0X, sdefR0Y, sdefR0Z},
		SdefR1: &mvec3.T{sdefR1X, sdefR1Y, sdefR1Z},
	}
}

// Type returns the deformation type.
func (s *Sdef) Type() DeformType {
	return SDEF
}
