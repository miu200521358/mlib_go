package morph

import (
	"github.com/miu200521358/mlib_go/pkg/core/index_model"
	"github.com/miu200521358/mlib_go/pkg/math/mrotation"
	"github.com/miu200521358/mlib_go/pkg/math/mvec3"
	"github.com/miu200521358/mlib_go/pkg/math/mvec4"

)

// MorphPanel 操作パネル
type MorphPanel int

const (
	// システム予約
	SYSTEM MorphPanel = 0
	// 眉(左下)
	EYEBROW_LOWER_LEFT MorphPanel = 1
	// 目(左上)
	EYE_UPPER_LEFT MorphPanel = 2
	// 口(右上)
	LIP_UPPER_RIGHT MorphPanel = 3
	// その他(右下)
	OTHER_LOWER_RIGHT MorphPanel = 4
)

// PanelName returns the name of the operation panel.
func (p MorphPanel) PanelName() string {
	switch p {
	case EYEBROW_LOWER_LEFT:
		return "眉"
	case EYE_UPPER_LEFT:
		return "目"
	case LIP_UPPER_RIGHT:
		return "口"
	case OTHER_LOWER_RIGHT:
		return "他"
	default:
		return "システム"
	}
}

// MorphType モーフ種類
type MorphType int

const (
	// グループ
	GROUP MorphType = 0
	// 頂点
	VERTEX MorphType = 1
	// ボーン
	BONE MorphType = 2
	// UV
	UV MorphType = 3
	// 追加UV1
	EXTENDED_UV1 MorphType = 4
	// 追加UV2
	EXTENDED_UV2 MorphType = 5
	// 追加UV3
	EXTENDED_UV3 MorphType = 6
	// 追加UV4
	EXTENDED_UV4 MorphType = 7
	// 材質
	MATERIAL MorphType = 8
	// ボーン変形後頂点(システム独自)
	AFTER_VERTEX MorphType = 9
)

// Morph represents a morph.
type Morph struct {
	*index_model.IndexModel
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
func (t *Morph) Copy() index_model.IndexModelInterface {
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
	Position mvec3.T
}

func (v *VertexMorphOffset) GetType() int {
	return int(VERTEX)
}

// UvMorphOffset represents a UV morph.
type UvMorphOffset struct {
	// 頂点INDEX
	VertexIndex int
	// UVオフセット量(x,y,z,w) ※通常UVはz,wが不要項目になるがモーフとしてのデータ値は記録しておく
	Uv mvec4.T
}

func (v *UvMorphOffset) GetType() int {
	return int(UV)
}

// BoneMorphOffset represents a bone morph.
type BoneMorphOffset struct {
	// ボーンIndex
	BoneIndex int
	// グローバル移動量(x,y,z)
	Position mvec3.T
	// グローバル回転量-クォータニオン(x,y,z,w)
	Rotation mrotation.T
	// グローバル縮尺量(x,y,z) ※システム独自
	Scale mvec3.T
	// ローカル軸に沿った移動量(x,y,z) ※システム独自
	LocalPosition mvec3.T
	// ローカル軸に沿った回転量-クォータニオン(x,y,z,w) ※システム独自
	LocalRotation mrotation.T
	// ローカル軸に沿った縮尺量(x,y,z) ※システム独自
	LocalScale mvec3.T
}

func (v *BoneMorphOffset) GetType() int {
	return int(BONE)
}

// GroupMorphOffset represents a group morph.
type GroupMorphOffset struct {
	// モーフINDEX
	MorphIndex int
	// モーフ変動量
	MorphFactor float64
}

func (v *GroupMorphOffset) GetType() int {
	return int(GROUP)
}

// MaterialMorphCalcMode 材質モーフ：計算モード
type MaterialMorphCalcMode int

const (
	// 乗算
	MULTIPLICATION MaterialMorphCalcMode = 0
	// 加算
	ADDITION MaterialMorphCalcMode = 1
)

// MaterialMorphOffset represents a material morph.
type MaterialMorphOffset struct {
	// 材質Index -> -1:全材質対象
	MaterialIndex int
	// 0:乗算, 1:加算
	CalcMode MaterialMorphCalcMode
	// Diffuse
	Diffuse      mvec3.T
	DiffuseAlpha float64
	// Specular (R,G,B)
	Specular mvec3.T
	// Specular係数
	SpecularFactor float64
	// Ambient (R,G,B)
	Ambient mvec3.T
	// エッジ色 (R,G,B,A)
	EdgeColor mvec3.T
	EdgeAlpha float64
	// エッジサイズ
	EdgeSize float64
	// テクスチャ係数 (R,G,B,A)
	TextureCoefficient mvec3.T
	TextureAlpha       float64
	// スフィアテクスチャ係数 (R,G,B,A)
	SphereTextureCoefficient mvec3.T
	SphereTextureAlpha       float64
	// Toonテクスチャ係数 (R,G,B,A)
	ToonTextureCoefficient mvec3.T
	ToonTextureAlpha       float64
}

func (v *MaterialMorphOffset) GetType() int {
	return int(MATERIAL)
}

// NewMorph
func NewMorph(
	index int,
	name, englishName string,
	panel MorphPanel,
	morphType MorphType,
	offsets []TMorphOffset,
	displaySlot int,
	isSystem bool,
) *Morph {
	return &Morph{
		IndexModel:  &index_model.IndexModel{Index: index},
		Name:        name,
		EnglishName: englishName,
		Panel:       panel,
		MorphType:   morphType,
		Offsets:     offsets,
		DisplaySlot: displaySlot,
		IsSystem:    isSystem,
	}
}

// モーフリスト
type Morphs struct {
	*index_model.IndexModelCorrection[*Morph]
}

func NewMorphs(name string) *Morphs {
	return &Morphs{
		IndexModelCorrection: index_model.NewIndexModelCorrection[*Morph](),
	}
}
