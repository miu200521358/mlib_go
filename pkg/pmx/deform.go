package pmx

import (
	"sort"

	"github.com/go-gl/mathgl/mgl32"

	"github.com/miu200521358/mlib_go/pkg/mmath"
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
	GetAllIndexes() []int
	GetAllWeights() []float64
	GetIndexes(weightThreshold float64) []int
	GetWeights(weightThreshold float64) []float64
	NormalizedDeform() [8]float32
	GetSdefParams() (mgl32.Vec3, mgl32.Vec3, mgl32.Vec3)
}

// Deform デフォーム既定構造体
type Deform struct {
	Indexes []int     // ボーンINDEXリスト
	Weights []float64 // ウェイトリスト
	Count   int       // デフォームボーン個数
}

// NewDeform creates a new Deform instance.
func NewDeform(indexes []int, weights []float64, count int) Deform {
	return Deform{
		Indexes: indexes,
		Weights: weights,
		Count:   count,
	}
}

func (d *Deform) GetAllIndexes() []int {
	return d.Indexes
}

func (d *Deform) GetAllWeights() []float64 {
	return d.Weights
}

// GetIndexes ウェイト閾値以上のウェイトを持っているINDEXのみを取得する
func (d *Deform) GetIndexes(weightThreshold float64) []int {
	var indexes []int
	for i, weight := range d.Weights {
		if weight >= weightThreshold {
			indexes = append(indexes, d.Indexes[i])
		}
	}
	return indexes
}

// GetWeights ウェイト閾値以上のウェイトを持っているウェイトのみを取得する
func (d *Deform) GetWeights(weightThreshold float64) []float64 {
	var weights []float64
	for _, weight := range d.Weights {
		if weight >= weightThreshold {
			weights = append(weights, weight)
		}
	}
	return weights
}

// Normalize ウェイト正規化
func (d *Deform) Normalize(align bool) {
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
		d.Indexes, d.Weights = sortIndexesByWeight(ilist, wlist)
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

// NormalizedDeform 4つのボーンINDEXとウェイトを返す（合計8個）
func (d *Deform) NormalizedDeform() [8]float32 {
	normalizedDeform := [8]float32{0, 0, 0, 0, 0, 0, 0, 0}
	for i, index := range d.Indexes {
		normalizedDeform[i] = float32(index)
	}
	for i, weight := range d.Weights {
		normalizedDeform[i+4] = float32(weight)
	}

	return normalizedDeform
}

// SDEF用パラメーターを返す
func (d *Deform) GetSdefParams() (mgl32.Vec3, mgl32.Vec3, mgl32.Vec3) {
	return mgl32.Vec3{0, 0, 0}, mgl32.Vec3{0, 0, 0}, mgl32.Vec3{0, 0, 0}
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
func NewBdef1(index0 int) Bdef1 {
	return Bdef1{
		Deform: Deform{
			Indexes: []int{index0},
			Weights: []float64{1.0},
			Count:   1,
		},
	}
}

// GetType returns the deformation type.
func (b *Bdef1) GetType() DeformType {
	return BDEF1
}

// NormalizedDeform 4つのボーンINDEXとウェイトを返す（合計8個）
func (d *Bdef1) NormalizedDeform() [8]float32 {
	return [8]float32{float32(d.Indexes[0]), 0, 0, 0, 1.0, 0, 0, 0}
}

// Bdef2 represents the BDEF2 deformation.
type Bdef2 struct {
	Deform
}

// NewBdef2 creates a new Bdef2 instance.
func NewBdef2(index0, index1 int, weight0 float64) Bdef2 {
	return Bdef2{
		Deform: Deform{
			Indexes: []int{index0, index1},
			Weights: []float64{weight0, 1 - weight0},
			Count:   2,
		},
	}
}

// GetType returns the deformation type.
func (b *Bdef2) GetType() DeformType {
	return BDEF2
}

// NormalizedDeform 4つのボーンINDEXとウェイトを返す（合計8個）
func (d *Bdef2) NormalizedDeform() [8]float32 {
	return [8]float32{
		float32(d.Indexes[0]), float32(d.Indexes[1]), 0, 0,
		float32(d.Weights[0]), float32(1 - d.Weights[0]), 0, 0}
}

// Bdef4 represents the BDEF4 deformation.
type Bdef4 struct {
	Deform
}

// NewBdef4 creates a new Bdef4 instance.
func NewBdef4(index0, index1, index2, index3 int, weight0, weight1, weight2, weight3 float64) Bdef4 {
	return Bdef4{
		Deform: Deform{
			Indexes: []int{index0, index1, index2, index3},
			Weights: []float64{weight0, weight1, weight2, weight3},
			Count:   4,
		},
	}
}

// GetType returns the deformation type.
func (b *Bdef4) GetType() DeformType {
	return BDEF4
}

// NormalizedDeform 4つのボーンINDEXとウェイトを返す（合計8個）
func (d *Bdef4) NormalizedDeform() [8]float32 {
	return [8]float32{
		float32(d.Indexes[0]), float32(d.Indexes[1]), float32(d.Indexes[2]), float32(d.Indexes[3]),
		float32(d.Weights[0]), float32(d.Weights[1]), float32(d.Weights[2]), float32(d.Weights[3])}
}

// Sdef represents the SDEF deformation.
type Sdef struct {
	Deform
	SdefC  *mmath.MVec3
	SdefR0 *mmath.MVec3
	SdefR1 *mmath.MVec3
}

// NewSdef creates a new Sdef instance.
func NewSdef(index0, index1 int, weight0 float64, sdefC, sdefR0, sdefR1 *mmath.MVec3) Sdef {
	return Sdef{
		Deform: Deform{
			Indexes: []int{index0, index1},
			Weights: []float64{weight0, 1 - weight0},
			Count:   2,
		},
		SdefC:  sdefC,
		SdefR0: sdefR0,
		SdefR1: sdefR1,
	}
}

// GetType returns the deformation type.
func (s *Sdef) GetType() DeformType {
	return SDEF
}

// NormalizedDeform 4つのボーンINDEXとウェイトを返す（合計8個）
func (d *Sdef) NormalizedDeform() [8]float32 {
	return [8]float32{
		float32(d.Indexes[0]), float32(d.Indexes[1]), 0, 0,
		float32(d.Weights[0]), float32(1 - d.Weights[0]), 0, 0}
}
