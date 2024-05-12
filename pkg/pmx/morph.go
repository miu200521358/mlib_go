package pmx

import (
	"github.com/jinzhu/copier"

	"github.com/miu200521358/mlib_go/pkg/mcore"
	"github.com/miu200521358/mlib_go/pkg/mmath"
)

// MorphPanel 操作パネル
type MorphPanel byte

const (
	MORPH_PANEL_SYSTEM             MorphPanel = 0 // システム予約
	MORPH_PANEL_EYEBROW_LOWER_LEFT MorphPanel = 1 // 眉(左下)
	MORPH_PANEL_EYE_UPPER_LEFT     MorphPanel = 2 // 目(左上)
	MORPH_PANEL_LIP_UPPER_RIGHT    MorphPanel = 3 // 口(右上)
	MORPH_PANEL_OTHER_LOWER_RIGHT  MorphPanel = 4 // その他(右下)
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
	MORPH_TYPE_GROUP        MorphType = 0 // グループ
	MORPH_TYPE_VERTEX       MorphType = 1 // 頂点
	MORPH_TYPE_BONE         MorphType = 2 // ボーン
	MORPH_TYPE_UV           MorphType = 3 // UV
	MORPH_TYPE_EXTENDED_UV1 MorphType = 4 // 追加UV1
	MORPH_TYPE_EXTENDED_UV2 MorphType = 5 // 追加UV2
	MORPH_TYPE_EXTENDED_UV3 MorphType = 6 // 追加UV3
	MORPH_TYPE_EXTENDED_UV4 MorphType = 7 // 追加UV4
	MORPH_TYPE_MATERIAL     MorphType = 8 // 材質
	MORPH_TYPE_AFTER_VERTEX MorphType = 9 // ボーン変形後頂点(システム独自)
)

// Morph represents a morph.
type Morph struct {
	*mcore.IndexNameModel
	Panel       MorphPanel     // モーフパネル
	MorphType   MorphType      // モーフ種類
	Offsets     []TMorphOffset // モーフオフセット
	DisplaySlot int            // 表示枠
	IsSystem    bool           // ツール側で追加したモーフ
}

// TMorphOffset represents a morph offset.
type TMorphOffset interface {
	GetType() int
}

// VertexMorphOffset represents a vertex morph.
type VertexMorphOffset struct {
	VertexIndex int         // 頂点INDEX
	Position    mmath.MVec3 // 座標オフセット量(x,y,z)
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
	VertexIndex int         // 頂点INDEX
	Uv          mmath.MVec4 // UVオフセット量(x,y,z,w)
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
	BoneIndex     int             // ボーンIndex
	Position      mmath.MVec3     // グローバル移動量(x,y,z)
	Rotation      mmath.MRotation // グローバル回転量-クォータニオン(x,y,z,w)
	Scale         mmath.MVec3     // グローバル縮尺量(x,y,z) ※システム独自
	LocalPosition mmath.MVec3     // ローカル軸に沿った移動量(x,y,z) ※システム独自
	LocalRotation mmath.MRotation // ローカル軸に沿った回転量-クォータニオン(x,y,z,w) ※システム独自
	LocalScale    mmath.MVec3     // ローカル軸に沿った縮尺量(x,y,z) ※システム独自
}

func (v *BoneMorphOffset) GetType() int {
	return int(MORPH_TYPE_BONE)
}

func NewBoneMorph(boneIndex int, position mmath.MVec3, rotation mmath.MRotation) *BoneMorphOffset {
	return &BoneMorphOffset{
		BoneIndex:     boneIndex,
		Position:      position,
		Rotation:      rotation,
		Scale:         mmath.NewMVec3(),
		LocalPosition: mmath.NewMVec3(),
		LocalRotation: *mmath.NewRotation(),
		LocalScale:    mmath.NewMVec3(),
	}
}

// GroupMorphOffset represents a group morph.
type GroupMorphOffset struct {
	MorphIndex  int     // モーフINDEX
	MorphFactor float64 // モーフ変動量
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
	CALC_MODE_MULTIPLICATION MaterialMorphCalcMode = 0 // 乗算
	CALC_MODE_ADDITION       MaterialMorphCalcMode = 1 // 加算
)

// MaterialMorphOffset represents a material morph.
type MaterialMorphOffset struct {
	MaterialIndex       int                   // 材質Index -> -1:全材質対象
	CalcMode            MaterialMorphCalcMode // 0:乗算, 1:加算
	Diffuse             mmath.MVec4           // Diffuse (R,G,B,A)
	Specular            mmath.MVec4           // SpecularColor (R,G,B, 係数)
	Ambient             mmath.MVec3           // AmbientColor (R,G,B)
	Edge                mmath.MVec4           // エッジ色 (R,G,B,A)
	EdgeSize            float64               // エッジサイズ
	TextureFactor       mmath.MVec4           // テクスチャ係数 (R,G,B,A)
	SphereTextureFactor mmath.MVec4           // スフィアテクスチャ係数 (R,G,B,A)
	ToonTextureFactor   mmath.MVec4           // Toonテクスチャ係数 (R,G,B,A)
}

func (v *MaterialMorphOffset) GetType() int {
	return int(MORPH_TYPE_MATERIAL)
}

func NewMaterialMorph(
	materialIndex int,
	calcMode MaterialMorphCalcMode,
	diffuse mmath.MVec4,
	specular mmath.MVec4,
	ambient mmath.MVec3,
	edge mmath.MVec4,
	edgeSize float64,
	textureFactor mmath.MVec4,
	sphereTextureFactor mmath.MVec4,
	toonTextureFactor mmath.MVec4,
) *MaterialMorphOffset {
	return &MaterialMorphOffset{
		MaterialIndex:       materialIndex,
		CalcMode:            calcMode,
		Diffuse:             diffuse,
		Specular:            specular,
		Ambient:             ambient,
		Edge:                edge,
		EdgeSize:            edgeSize,
		TextureFactor:       textureFactor,
		SphereTextureFactor: sphereTextureFactor,
		ToonTextureFactor:   toonTextureFactor,
	}
}

// NewMorph
func NewMorph() *Morph {
	return &Morph{
		IndexNameModel: &mcore.IndexNameModel{Index: -1, Name: "", EnglishName: ""},
		Panel:          MORPH_PANEL_SYSTEM,
		MorphType:      MORPH_TYPE_VERTEX,
		Offsets:        make([]TMorphOffset, 0),
		DisplaySlot:    -1,
		IsSystem:       false,
	}
}

func (m *Morph) Copy() mcore.IIndexNameModel {
	copied := NewMorph()
	copier.CopyWithOption(copied, m, copier.Option{DeepCopy: true})
	return copied
}

// モーフリスト
type Morphs struct {
	*mcore.IndexNameModels[*Morph]
}

func NewMorphs() *Morphs {
	return &Morphs{
		IndexNameModels: mcore.NewIndexNameModels[*Morph](),
	}
}
