package pmx

import (
	"github.com/miu200521358/mlib_go/pkg/mcore"
	"github.com/miu200521358/mlib_go/pkg/mmath"
)

// MorphPanel 操作パネル
type MorphPanel byte

const (
	// システム予約
	MORPH_PANEL_SYSTEM MorphPanel = 0
	// 眉(左下)
	MORPH_PANEL_EYEBROW_LOWER_LEFT MorphPanel = 1
	// 目(左上)
	MORPH_PANEL_EYE_UPPER_LEFT MorphPanel = 2
	// 口(右上)
	MORPH_PANEL_LIP_UPPER_RIGHT MorphPanel = 3
	// その他(右下)
	MORPH_PANEL_OTHER_LOWER_RIGHT MorphPanel = 4
)

// PanelName returns the name of the operation panel.
func (p MorphPanel) PanelName() string {
	switch p {
	case MORPH_PANEL_EYEBROW_LOWER_LEFT:
		return "眉"
	case MORPH_PANEL_EYE_UPPER_LEFT:
		return "目"
	case MORPH_PANEL_LIP_UPPER_RIGHT:
		return "口"
	case MORPH_PANEL_OTHER_LOWER_RIGHT:
		return "他"
	default:
		return "システム"
	}
}

// MorphType モーフ種類
type MorphType int

const (
	// グループ
	MORPH_TYPE_GROUP MorphType = 0
	// 頂点
	MORPH_TYPE_VERTEX MorphType = 1
	// ボーン
	MORPH_TYPE_BONE MorphType = 2
	// MORPH_TYPE_UV
	MORPH_TYPE_UV MorphType = 3
	// 追加UV1
	MORPH_TYPE_EXTENDED_UV1 MorphType = 4
	// 追加UV2
	MORPH_TYPE_EXTENDED_UV2 MorphType = 5
	// 追加UV3
	MORPH_TYPE_EXTENDED_UV3 MorphType = 6
	// 追加UV4
	MORPH_TYPE_EXTENDED_UV4 MorphType = 7
	// 材質
	MORPH_TYPE_MATERIAL MorphType = 8
	// ボーン変形後頂点(システム独自)
	MORPH_TYPE_AFTER_VERTEX MorphType = 9
)

// Morph represents a morph.
type Morph struct {
	*mcore.IndexModel
	// モーフ名
	Name string
	// モーフ名英
	EnglishName string
	// モーフパネル
	Panel MorphPanel
	// モーフ種類
	MorphType MorphType
	// モーフオフセット
	Offsets []TMorphOffset
	// 表示枠
	DisplaySlot int
	// ツール側で追加したモーフ
	IsSystem bool
}

// Copy
func (t *Morph) Copy() mcore.IndexModelInterface {
	copied := *t
	copied.Offsets = make([]TMorphOffset, len(t.Offsets))
	copy(copied.Offsets, t.Offsets)
	return &copied
}

// TMorphOffset represents a morph offset.
type TMorphOffset interface {
	GetType() int
}

// VertexMorphOffset represents a vertex morph.
type VertexMorphOffset struct {
	// 頂点INDEX
	VertexIndex int
	// 座標オフセット量(x,y,z)
	Position mmath.MVec3
}

func (v *VertexMorphOffset) GetType() int {
	return int(MORPH_TYPE_VERTEX)
}

func NewVertexMorph(vertexIndex int, position mmath.MVec3) *VertexMorphOffset {
	return &VertexMorphOffset{
		VertexIndex: vertexIndex,
		Position:    position,
	}
}

// UvMorphOffset represents a UV morph.
type UvMorphOffset struct {
	// 頂点INDEX
	VertexIndex int
	// UVオフセット量(x,y,z,w) ※通常UVはz,wが不要項目になるがモーフとしてのデータ値は記録しておく
	Uv mmath.MVec4
}

func (v *UvMorphOffset) GetType() int {
	return int(MORPH_TYPE_UV)
}

func NewUvMorph(vertexIndex int, uv mmath.MVec4) *UvMorphOffset {
	return &UvMorphOffset{
		VertexIndex: vertexIndex,
		Uv:          uv,
	}
}

// BoneMorphOffset represents a bone morph.
type BoneMorphOffset struct {
	// ボーンIndex
	BoneIndex int
	// グローバル移動量(x,y,z)
	Position mmath.MVec3
	// グローバル回転量-クォータニオン(x,y,z,w)
	Rotation mmath.MRotation
	// グローバル縮尺量(x,y,z) ※システム独自
	Scale mmath.MVec3
	// ローカル軸に沿った移動量(x,y,z) ※システム独自
	LocalPosition mmath.MVec3
	// ローカル軸に沿った回転量-クォータニオン(x,y,z,w) ※システム独自
	LocalRotation mmath.MRotation
	// ローカル軸に沿った縮尺量(x,y,z) ※システム独自
	LocalScale mmath.MVec3
}

func (v *BoneMorphOffset) GetType() int {
	return int(MORPH_TYPE_BONE)
}

func NewBoneMorph(boneIndex int, position mmath.MVec3, rotation mmath.MRotation) *BoneMorphOffset {
	return &BoneMorphOffset{
		BoneIndex:     boneIndex,
		Position:      position,
		Rotation:      rotation,
		Scale:         mmath.MVec3{},
		LocalPosition: mmath.MVec3{},
		LocalRotation: mmath.MRotation{},
		LocalScale:    mmath.MVec3{},
	}
}

// GroupMorphOffset represents a group morph.
type GroupMorphOffset struct {
	// モーフINDEX
	MorphIndex int
	// モーフ変動量
	MorphFactor float64
}

func NewGroupMorph(morphIndex int, morphFactor float64) *GroupMorphOffset {
	return &GroupMorphOffset{
		MorphIndex:  morphIndex,
		MorphFactor: morphFactor,
	}
}

func (v *GroupMorphOffset) GetType() int {
	return int(MORPH_TYPE_GROUP)
}

// MaterialMorphCalcMode 材質モーフ：計算モード
type MaterialMorphCalcMode int

const (
	// 乗算
	CALC_MODE_MULTIPLICATION MaterialMorphCalcMode = 0
	// 加算
	CALC_MODE_ADDITION MaterialMorphCalcMode = 1
)

// MaterialMorphOffset represents a material morph.
type MaterialMorphOffset struct {
	// 材質Index -> -1:全材質対象
	MaterialIndex int
	// 0:乗算, 1:加算
	CalcMode MaterialMorphCalcMode
	// DiffuseColor
	DiffuseColor mmath.MVec3
	DiffuseAlpha float64
	// SpecularColor (R,G,B)
	SpecularColor mmath.MVec3
	// Specular係数
	SpecularPower float64
	// AmbientColor (R,G,B)
	AmbientColor mmath.MVec3
	// エッジ色 (R,G,B,A)
	EdgeColor mmath.MVec3
	EdgeAlpha float64
	// エッジサイズ
	EdgeSize float64
	// テクスチャ係数 (R,G,B,A)
	TextureCoefficient mmath.MVec3
	TextureAlpha       float64
	// スフィアテクスチャ係数 (R,G,B,A)
	SphereTextureCoefficient mmath.MVec3
	SphereTextureAlpha       float64
	// Toonテクスチャ係数 (R,G,B,A)
	ToonTextureCoefficient mmath.MVec3
	ToonTextureAlpha       float64
}

func (v *MaterialMorphOffset) GetType() int {
	return int(MORPH_TYPE_MATERIAL)
}

func NewMaterialMorph(materialIndex int, calcMode MaterialMorphCalcMode, diffuseColor mmath.MVec3, diffuseAlpha float64, specularColor mmath.MVec3, specularPower float64, ambientColor mmath.MVec3, edgeColor mmath.MVec3, edgeAlpha float64, edgeSize float64, textureCoefficient mmath.MVec3, textureAlpha float64, sphereTextureCoefficient mmath.MVec3, sphereTextureAlpha float64, toonTextureCoefficient mmath.MVec3, toonTextureAlpha float64) *MaterialMorphOffset {
	return &MaterialMorphOffset{
		MaterialIndex:            materialIndex,
		CalcMode:                 calcMode,
		DiffuseColor:             diffuseColor,
		DiffuseAlpha:             diffuseAlpha,
		SpecularColor:            specularColor,
		SpecularPower:            specularPower,
		AmbientColor:             ambientColor,
		EdgeColor:                edgeColor,
		EdgeAlpha:                edgeAlpha,
		EdgeSize:                 edgeSize,
		TextureCoefficient:       textureCoefficient,
		TextureAlpha:             textureAlpha,
		SphereTextureCoefficient: sphereTextureCoefficient,
		SphereTextureAlpha:       sphereTextureAlpha,
		ToonTextureCoefficient:   toonTextureCoefficient,
		ToonTextureAlpha:         toonTextureAlpha,
	}
}

// NewMorph
func NewMorph() *Morph {
	return &Morph{
		IndexModel:  &mcore.IndexModel{Index: -1},
		Name:        "",
		EnglishName: "",
		Panel:       MORPH_PANEL_SYSTEM,
		MorphType:   MORPH_TYPE_VERTEX,
		Offsets:     make([]TMorphOffset, 0),
		DisplaySlot: -1,
		IsSystem:    false,
	}
}

// モーフリスト
type Morphs struct {
	*mcore.IndexModelCorrection[*Morph]
}

func NewMorphs() *Morphs {
	return &Morphs{
		IndexModelCorrection: mcore.NewIndexModelCorrection[*Morph](),
	}
}
