package model

import "github.com/miu200521358/mlib_go/pkg/domain/mmath"

// MorphPanel はモーフパネルを表す。
type MorphPanel int

const (
	// MORPH_PANEL_SYSTEM はシステムパネル。
	MORPH_PANEL_SYSTEM MorphPanel = iota
	// MORPH_PANEL_EYEBROW_LOWER_LEFT は眉パネル。
	MORPH_PANEL_EYEBROW_LOWER_LEFT
	// MORPH_PANEL_EYE_UPPER_LEFT は目パネル。
	MORPH_PANEL_EYE_UPPER_LEFT
	// MORPH_PANEL_LIP_UPPER_RIGHT は口パネル。
	MORPH_PANEL_LIP_UPPER_RIGHT
	// MORPH_PANEL_OTHER_LOWER_RIGHT はその他パネル。
	MORPH_PANEL_OTHER_LOWER_RIGHT
)

// MorphType はモーフ種別を表す。
type MorphType int

const (
	// MORPH_TYPE_GROUP はグループモーフ。
	MORPH_TYPE_GROUP MorphType = iota
	// MORPH_TYPE_VERTEX は頂点モーフ。
	MORPH_TYPE_VERTEX
	// MORPH_TYPE_BONE はボーンモーフ。
	MORPH_TYPE_BONE
	// MORPH_TYPE_UV はUVモーフ。
	MORPH_TYPE_UV
	// MORPH_TYPE_EXTENDED_UV1 は追加UV1モーフ。
	MORPH_TYPE_EXTENDED_UV1
	// MORPH_TYPE_EXTENDED_UV2 は追加UV2モーフ。
	MORPH_TYPE_EXTENDED_UV2
	// MORPH_TYPE_EXTENDED_UV3 は追加UV3モーフ。
	MORPH_TYPE_EXTENDED_UV3
	// MORPH_TYPE_EXTENDED_UV4 は追加UV4モーフ。
	MORPH_TYPE_EXTENDED_UV4
	// MORPH_TYPE_MATERIAL は材質モーフ。
	MORPH_TYPE_MATERIAL
	// MORPH_TYPE_AFTER_VERTEX はボーン変形後頂点モーフ。
	MORPH_TYPE_AFTER_VERTEX
)

// MaterialMorphCalcMode は材質モーフ計算モードを表す。
type MaterialMorphCalcMode int

const (
	// CALC_MODE_MULTIPLICATION は乗算。
	CALC_MODE_MULTIPLICATION MaterialMorphCalcMode = iota
	// CALC_MODE_ADDITION は加算。
	CALC_MODE_ADDITION
)

// MorphOffset はモーフオフセットを表す。
type MorphOffset interface {
	MorphType() MorphType
}

// VertexMorphOffset は頂点モーフオフセットを表す。
type VertexMorphOffset struct {
	VertexIndex int
	Position    mmath.Vec3
}

// MorphType はオフセット種別を返す。
func (o *VertexMorphOffset) MorphType() MorphType {
	return MORPH_TYPE_VERTEX
}

// UvMorphOffset はUVモーフオフセットを表す。
type UvMorphOffset struct {
	VertexIndex int
	Uv          mmath.Vec4
	UvType      MorphType
}

// MorphType はオフセット種別を返す。
func (o *UvMorphOffset) MorphType() MorphType {
	return o.UvType
}

// BoneMorphOffset はボーンモーフオフセットを表す。
type BoneMorphOffset struct {
	BoneIndex int
	Position  mmath.Vec3
	Rotation  mmath.Quaternion
}

// MorphType はオフセット種別を返す。
func (o *BoneMorphOffset) MorphType() MorphType {
	return MORPH_TYPE_BONE
}

// GroupMorphOffset はグループモーフオフセットを表す。
type GroupMorphOffset struct {
	MorphIndex  int
	MorphFactor float64
}

// MorphType はオフセット種別を返す。
func (o *GroupMorphOffset) MorphType() MorphType {
	return MORPH_TYPE_GROUP
}

// MaterialMorphOffset は材質モーフオフセットを表す。
type MaterialMorphOffset struct {
	MaterialIndex       int
	CalcMode            MaterialMorphCalcMode
	Diffuse             mmath.Vec4
	Specular            mmath.Vec4
	Ambient             mmath.Vec3
	Edge                mmath.Vec4
	EdgeSize            float64
	TextureFactor       mmath.Vec4
	SphereTextureFactor mmath.Vec4
	ToonTextureFactor   mmath.Vec4
}

// MorphType はオフセット種別を返す。
func (o *MaterialMorphOffset) MorphType() MorphType {
	return MORPH_TYPE_MATERIAL
}

// Morph はモーフ要素を表す。
type Morph struct {
	index       int
	name        string
	EnglishName string
	Panel       MorphPanel
	MorphType   MorphType
	Offsets     []MorphOffset
	DisplaySlot int
}

// Index はモーフ index を返す。
func (m *Morph) Index() int {
	return m.index
}

// SetIndex はモーフ index を設定する。
func (m *Morph) SetIndex(index int) {
	m.index = index
}

// Name はモーフ名を返す。
func (m *Morph) Name() string {
	return m.name
}

// SetName はモーフ名を設定する。
func (m *Morph) SetName(name string) {
	m.name = name
}

// IsValid はモーフが有効か判定する。
func (m *Morph) IsValid() bool {
	return m != nil && m.index >= 0
}
