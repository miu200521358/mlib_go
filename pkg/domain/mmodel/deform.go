// Package mmodel はPMXモデルのエンティティを定義します。
package mmodel

import (
	"slices"

	"github.com/miu200521358/mlib_go/pkg/domain/mmath"
)

// DeformType ウェイト変形方式
type DeformType byte

const (
	DEFORM_BDEF1 DeformType = 0 // BDEF1: 単一ボーンウェイト
	DEFORM_BDEF2 DeformType = 1 // BDEF2: 2ボーンウェイト
	DEFORM_BDEF4 DeformType = 2 // BDEF4: 4ボーンウェイト
	DEFORM_SDEF  DeformType = 3 // SDEF: 球面変形
)

// IDeform はデフォームのインターフェースです。
type IDeform interface {
	// Type はデフォームタイプを返します。
	Type() DeformType
	// Indexes はボーンインデックスリストを返します。
	Indexes() []int
	// SetIndexes はボーンインデックスリストを設定します。
	SetIndexes(indexes []int)
	// Weights はウェイトリストを返します。
	Weights() []float64
	// IndexesByWeight はウェイト閾値以上のボーンインデックスを返します。
	IndexesByWeight(weightThreshold float64) []int
	// WeightsByWeight はウェイト閾値以上のウェイトを返します。
	WeightsByWeight(weightThreshold float64) []float64
	// Packed はGPU用の8要素配列を返します。
	Packed() [8]float32
	// Normalize はウェイトを正規化します。
	Normalize(align bool)
	// Index はボーンインデックスの位置を返します。
	Index(boneIndex int) int
	// IndexWeight はボーンのウェイトを返します。
	IndexWeight(boneIndex int) float64
	// SplitWeight はウェイトを分割して追加します。
	SplitWeight(fromIdx, toIdx int, ratio float64)
}

// Deform はデフォームの基底構造体です。
type Deform struct {
	indexes []int
	weights []float64
	typ     DeformType
}

// Type はデフォームタイプを返します。
func (d *Deform) Type() DeformType {
	return d.typ
}

// Index はボーンインデックスの位置を返します。見つからない場合は-1を返します。
func (d *Deform) Index(boneIdx int) int {
	for i, idx := range d.indexes {
		if idx == boneIdx {
			return i
		}
	}
	return -1
}

// IndexWeight はボーンのウェイトを返します。見つからない場合は0を返します。
func (d *Deform) IndexWeight(boneIdx int) float64 {
	i := d.Index(boneIdx)
	if i == -1 {
		return 0
	}
	return d.weights[i]
}

// Indexes はボーンインデックスリストを返します。
func (d *Deform) Indexes() []int {
	return d.indexes
}

// SetIndexes はボーンインデックスリストを設定します。
func (d *Deform) SetIndexes(indexes []int) {
	d.indexes = indexes
}

// Weights はウェイトリストを返します。
func (d *Deform) Weights() []float64 {
	return d.weights
}

// SetWeights はウェイトリストを設定します。
func (d *Deform) SetWeights(weights []float64) {
	d.weights = weights
}

// IndexesByWeight はウェイト閾値以上のボーンインデックスを返します。
func (d *Deform) IndexesByWeight(threshold float64) []int {
	result := make([]int, 0, len(d.indexes))
	for i, w := range d.weights {
		if w >= threshold {
			result = append(result, d.indexes[i])
		}
	}
	return result
}

// WeightsByWeight はウェイト閾値以上のウェイトを返します。
func (d *Deform) WeightsByWeight(threshold float64) []float64 {
	result := make([]float64, 0, len(d.weights))
	for _, w := range d.weights {
		if w >= threshold {
			result = append(result, w)
		}
	}
	return result
}

// Packed はGPU用の8要素配列（ボーンインデックス4つ + ウェイト4つ）を返します。
func (d *Deform) Packed() [8]float32 {
	out := [8]float32{0, 0, 0, 0, 0, 0, 0, 0}
	for i, idx := range d.indexes {
		if i < 4 {
			out[i] = float32(idx)
		}
	}
	for i, w := range d.weights {
		if i < 4 {
			out[i+4] = float32(w)
		}
	}
	return out
}

// Normalize はウェイトを正規化します。
// align=trueの場合、ウェイトを統合して1,2,4個に揃えます。
func (d *Deform) Normalize(align bool) {
	if align {
		// ウェイトを統合する
		merged := make(map[int]float64)
		for i, idx := range d.indexes {
			merged[idx] += d.weights[i]
		}

		// 数が足りるようかさ増しする
		idxs := make([]int, 0, len(merged)+4)
		wgts := make([]float64, 0, len(merged)+4)
		for idx, w := range merged {
			idxs = append(idxs, idx)
			wgts = append(wgts, w)
		}
		for i := len(merged); i < 8; i++ {
			idxs = append(idxs, 0)
			wgts = append(wgts, 0)
		}

		// 正規化
		sum := 0.0
		for _, w := range wgts {
			sum += w
		}
		if sum > 0 {
			for i := range wgts {
				wgts[i] /= sum
			}
		}

		// ウェイトの大きい順にソート
		d.indexes, d.weights = sortByWeight(idxs, wgts)
	}

	// ウェイト正規化
	sum := 0.0
	for _, w := range d.weights {
		sum += w
	}
	if sum > 0 {
		for i := range d.weights {
			d.weights[i] /= sum
		}
	}
}

// SplitWeight はウェイトを分割して追加します。
// fromIdxのウェイトをratioで分割してtoIdxに追加します。
func (d *Deform) SplitWeight(fromIdx, toIdx int, ratio float64) {
	for i, idx := range d.indexes {
		if idx == fromIdx {
			d.indexes = append(d.indexes, toIdx)
			d.weights = append(d.weights, d.weights[i]*ratio)
			d.weights[i] *= 1 - ratio
			break
		}
	}
	d.Normalize(true)
}

// sortByWeight はウェイトの大きい順にソートし、有効な個数(1,2,4)を返します。
func sortByWeight(idxs []int, wgts []float64) ([]int, []float64) {
	order := argsort(wgts)
	slices.Reverse(order) // 降順

	outIdxs := make([]int, 0)
	outWgts := make([]float64, 0)

	for i, pos := range order {
		if wgts[pos] == 0 && (i >= 4 || i == 2 || i == 1) {
			break
		}
		outIdxs = append(outIdxs, idxs[pos])
		outWgts = append(outWgts, wgts[pos])
	}

	return outIdxs, outWgts
}

// argsort はスライスを昇順にソートした場合のインデックス配列を返します。
func argsort(vals []float64) []int {
	idxs := make([]int, len(vals))
	for i := range idxs {
		idxs[i] = i
	}
	slices.SortFunc(idxs, func(a, b int) int {
		if vals[a] < vals[b] {
			return -1
		}
		if vals[a] > vals[b] {
			return 1
		}
		return 0
	})
	return idxs
}

// --------------------------------------------
// Bdef1
// --------------------------------------------

// Bdef1 は単一ボーンウェイトのデフォームです。
type Bdef1 struct {
	Deform
}

// NewBdef1 は新しいBdef1を生成します。
func NewBdef1(idx int) *Bdef1 {
	return &Bdef1{
		Deform: Deform{
			indexes: []int{idx},
			weights: []float64{1.0},
			typ:     DEFORM_BDEF1,
		},
	}
}

// Packed はGPU用の8要素配列を返します。
func (b *Bdef1) Packed() [8]float32 {
	return [8]float32{float32(b.indexes[0]), 0, 0, 0, 1.0, 0, 0, 0}
}

// --------------------------------------------
// Bdef2
// --------------------------------------------

// Bdef2 は2ボーンウェイトのデフォームです。
type Bdef2 struct {
	Deform
}

// NewBdef2 は新しいBdef2を生成します。
func NewBdef2(idx0, idx1 int, weight float64) *Bdef2 {
	return &Bdef2{
		Deform: Deform{
			indexes: []int{idx0, idx1},
			weights: []float64{weight, 1 - weight},
			typ:     DEFORM_BDEF2,
		},
	}
}

// Packed はGPU用の8要素配列を返します。
func (b *Bdef2) Packed() [8]float32 {
	return [8]float32{
		float32(b.indexes[0]), float32(b.indexes[1]), 0, 0,
		float32(b.weights[0]), float32(1 - b.weights[0]), 0, 0,
	}
}

// --------------------------------------------
// Bdef4
// --------------------------------------------

// Bdef4 は4ボーンウェイトのデフォームです。
type Bdef4 struct {
	Deform
}

// NewBdef4 は新しいBdef4を生成します。
func NewBdef4(idx0, idx1, idx2, idx3 int, w0, w1, w2, w3 float64) *Bdef4 {
	return &Bdef4{
		Deform: Deform{
			indexes: []int{idx0, idx1, idx2, idx3},
			weights: []float64{w0, w1, w2, w3},
			typ:     DEFORM_BDEF4,
		},
	}
}

// Packed はGPU用の8要素配列を返します。
func (b *Bdef4) Packed() [8]float32 {
	return [8]float32{
		float32(b.indexes[0]), float32(b.indexes[1]), float32(b.indexes[2]), float32(b.indexes[3]),
		float32(b.weights[0]), float32(b.weights[1]), float32(b.weights[2]), float32(b.weights[3]),
	}
}

// --------------------------------------------
// Sdef
// --------------------------------------------

// Sdef は球面変形のデフォームです。
type Sdef struct {
	Deform
	C  *mmath.Vec3 // SDEF-C値
	R0 *mmath.Vec3 // SDEF-R0値
	R1 *mmath.Vec3 // SDEF-R1値
}

// NewSdef は新しいSdefを生成します。
func NewSdef(idx0, idx1 int, weight float64, c, r0, r1 *mmath.Vec3) *Sdef {
	return &Sdef{
		Deform: Deform{
			indexes: []int{idx0, idx1},
			weights: []float64{weight, 1 - weight},
			typ:     DEFORM_SDEF,
		},
		C:  c,
		R0: r0,
		R1: r1,
	}
}

// Packed はGPU用の8要素配列を返します。
// TODO: SDEFパラメーターの正規化
func (s *Sdef) Packed() [8]float32 {
	return [8]float32{
		float32(s.indexes[0]), float32(s.indexes[1]), 0, 0,
		float32(s.weights[0]), float32(1 - s.weights[0]), 0, 0,
	}
}
