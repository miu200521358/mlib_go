// 指示: miu200521358
package model

import "github.com/miu200521358/mlib_go/pkg/domain/mmath"

// DeformType は頂点デフォーム種別を表す。
type DeformType int

const (
	// BDEF1 は単一ボーンのデフォーム。
	BDEF1 DeformType = iota
	// BDEF2 は2ボーンのデフォーム。
	BDEF2
	// BDEF4 は4ボーンのデフォーム。
	BDEF4
	// SDEF はSDEFデフォーム。
	SDEF
)

// IDeform はボーンデフォームを表す。
type IDeform interface {
	DeformType() DeformType
	Indexes() []int
	Weights() []float64
}

type deformBase struct {
	deformType DeformType
	indexes    []int
	weights    []float64
}

// DeformType はデフォーム種別を返す。
func (d *deformBase) DeformType() DeformType {
	return d.deformType
}

// Indexes はボーンインデックスの一覧を返す。
func (d *deformBase) Indexes() []int {
	return d.indexes
}

// Weights はボーンウェイトの一覧を返す。
func (d *deformBase) Weights() []float64 {
	return d.weights
}

// Bdef1 は単一ボーンのデフォームを表す。
type Bdef1 struct {
	deformBase
}

// NewBdef1 は Bdef1 を生成する。
func NewBdef1(boneIndex int) *Bdef1 {
	return &Bdef1{
		deformBase: deformBase{
			deformType: BDEF1,
			indexes:    []int{boneIndex},
			weights:    []float64{1.0},
		},
	}
}

// Bdef2 は2ボーンのデフォームを表す。
type Bdef2 struct {
	deformBase
}

// NewBdef2 は Bdef2 を生成する。
func NewBdef2(boneIndex0, boneIndex1 int, weight0 float64) *Bdef2 {
	return &Bdef2{
		deformBase: deformBase{
			deformType: BDEF2,
			indexes:    []int{boneIndex0, boneIndex1},
			weights:    []float64{weight0, 1.0 - weight0},
		},
	}
}

// Bdef4 は4ボーンのデフォームを表す。
type Bdef4 struct {
	deformBase
}

// NewBdef4 は Bdef4 を生成する。
func NewBdef4(indexes [4]int, weights [4]float64) *Bdef4 {
	return &Bdef4{
		deformBase: deformBase{
			deformType: BDEF4,
			indexes:    []int{indexes[0], indexes[1], indexes[2], indexes[3]},
			weights:    []float64{weights[0], weights[1], weights[2], weights[3]},
		},
	}
}

// Sdef はSDEFデフォームを表す。
type Sdef struct {
	deformBase
	SdefC  mmath.Vec3
	SdefR0 mmath.Vec3
	SdefR1 mmath.Vec3
}

// NewSdef は Sdef を生成する。
func NewSdef(boneIndex0, boneIndex1 int, weight0 float64) *Sdef {
	return &Sdef{
		deformBase: deformBase{
			deformType: SDEF,
			indexes:    []int{boneIndex0, boneIndex1},
			weights:    []float64{weight0, 1.0 - weight0},
		},
	}
}
