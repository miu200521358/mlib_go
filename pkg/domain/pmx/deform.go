package pmx

import (
	"slices"

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
	AllIndexes() []int
	AllWeights() []float64
	Indexes(weightThreshold float64) []int
	Weights(weightThreshold float64) []float64
	NormalizedDeform() [8]float32
	Normalize(align bool)
	Index(boneIndex int) int
	IndexWeight(boneIndex int) float64
	Add(boneIndex, separateIndex int, ratio float64)
}

// Deform デフォーム既定構造体
type Deform struct {
	indexes    []int      // ボーンINDEXリスト
	weights    []float64  // ウェイトリスト
	deformType DeformType // デフォームタイプ
}

func (deform *Deform) Index(boneIndex int) int {
	for i, index := range deform.indexes {
		if index == boneIndex {
			return i
		}
	}
	return -1
}

func (deform *Deform) IndexWeight(boneIndex int) float64 {
	i := deform.Index(boneIndex)
	if i == -1 {
		return 0
	}
	return deform.weights[i]
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
		for i := len(indexWeights); i < 8; i++ {
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

// Add ウェイトを分割して追加する(separateIndexのウェイトをratioで分割して追加)
func (deform *Deform) Add(boneIndex, separateIndex int, ratio float64) {
	for i, index := range deform.indexes {
		if index == separateIndex {
			deform.indexes = append(deform.indexes, boneIndex)
			deform.weights = append(deform.weights, deform.weights[i]*ratio)
			deform.weights[i] *= 1 - ratio
			break
		}
	}

	deform.Normalize(true)
}

// sortIndexesByWeight ウェイトの大きい順に指定個数(1,2,4)までを対象とする
func sortIndexesByWeight(indexes []int, weights []float64) ([]int, []float64) {
	sortableWeights := mmath.Float64Slice(weights)
	weightIndexes := mmath.ArgSort(sortableWeights)
	// 降順にする
	slices.Reverse(weightIndexes)

	// ウェイトの大きい順に指定個数までを対象とする
	sortedIndexes := make([]int, 0)
	sortedWeights := make([]float64, 0)

	for i, weightIndex := range weightIndexes {
		if weights[weightIndex] == 0 && ((i >= 4) || (i == 2) || (i == 1)) {
			// 1, 2, 4個溜まったら抜ける
			break
		}
		sortedIndexes = append(sortedIndexes, indexes[weightIndex])
		sortedWeights = append(sortedWeights, weights[weightIndex])
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
			indexes:    []int{index0},
			weights:    []float64{1.0},
			deformType: BDEF1,
		},
	}
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
			indexes:    []int{index0, index1},
			weights:    []float64{weight0, 1 - weight0},
			deformType: BDEF2,
		},
	}
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
			indexes:    []int{index0, index1, index2, index3},
			weights:    []float64{weight0, weight1, weight2, weight3},
			deformType: BDEF4,
		},
	}
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
			indexes:    []int{index0, index1},
			weights:    []float64{weight0, 1 - weight0},
			deformType: SDEF,
		},
		SdefC:  sdefC,
		SdefR0: sdefR0,
		SdefR1: sdefR1,
	}
}

// NormalizedDeform 4つのボーンINDEXとウェイトを返す（合計8個）
// TODO: SDEFパラメーターの正規化
func (sdef *Sdef) NormalizedDeform() [8]float32 {
	return [8]float32{
		float32(sdef.indexes[0]), float32(sdef.indexes[1]), 0, 0,
		float32(sdef.weights[0]), float32(1 - sdef.weights[0]), 0, 0}
}
