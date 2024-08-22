package pmx

import (
	"github.com/miu200521358/mlib_go/pkg/domain/core"
	"github.com/miu200521358/mlib_go/pkg/domain/mmath"
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
func (morphPanel MorphPanel) PanelName() string {
	switch morphPanel {
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
	index       int            // モーフINDEX
	name        string         // モーフ名
	englishName string         // モーフ英名
	Panel       MorphPanel     // モーフパネル
	MorphType   MorphType      // モーフ種類
	Offsets     []IMorphOffset // モーフオフセット
	DisplaySlot int            // 表示枠
	IsSystem    bool           // ツール側で追加したモーフ
}

func (morph *Morph) Index() int {
	return morph.index
}

func (morph *Morph) SetIndex(index int) {
	morph.index = index
}

func (morph *Morph) Name() string {
	return morph.name
}

func (morph *Morph) SetName(name string) {
	morph.name = name
}

func (morph *Morph) EnglishName() string {
	return morph.englishName
}

func (morph *Morph) SetEnglishName(englishName string) {
	morph.englishName = englishName
}

func (morph *Morph) IsValid() bool {
	return morph != nil && morph.index >= 0
}

// IMorphOffset represents a morph offset.
type IMorphOffset interface {
	Type() int
}

// VertexMorphOffset represents a vertex morph.
type VertexMorphOffset struct {
	VertexIndex int          // 頂点INDEX
	Position    *mmath.MVec3 // 座標オフセット量(x,y,z)
}

func (offset *VertexMorphOffset) Type() int {
	return int(MORPH_TYPE_VERTEX)
}

func NewVertexMorphOffset(vertexIndex int, position *mmath.MVec3) *VertexMorphOffset {
	return &VertexMorphOffset{
		VertexIndex: vertexIndex,
		Position:    position,
	}
}

// UvMorphOffset represents a UV morph.
type UvMorphOffset struct {
	VertexIndex int          // 頂点INDEX
	Uv          *mmath.MVec4 // UVオフセット量(x,y,z,w)
}

func (offset *UvMorphOffset) Type() int {
	return int(MORPH_TYPE_UV)
}

func NewUvMorphOffset(vertexIndex int, uv *mmath.MVec4) *UvMorphOffset {
	return &UvMorphOffset{
		VertexIndex: vertexIndex,
		Uv:          uv,
	}
}

// BoneMorphOffset represents a bone morph.
type BoneMorphOffset struct {
	BoneIndex int                // ボーンIndex
	Position  *mmath.MVec3       // グローバル移動量
	Rotation  *mmath.MQuaternion // グローバル回転量
	Extend    *BoneMorphOffsetExtend
}

type BoneMorphOffsetExtend struct {
	Scale              *mmath.MVec3       // グローバル縮尺量
	CancelablePosition *mmath.MVec3       // 親キャンセル付き移動量
	CancelableRotation *mmath.MQuaternion // 親キャンセル付き回転量
	CancelableScale    *mmath.MVec3       // 親キャンセル付き縮尺量
	LocalMat           *mmath.MMat4       // ローカル変換行列
}

func (offset *BoneMorphOffset) Type() int {
	return int(MORPH_TYPE_BONE)
}

func NewBoneMorphOffset(boneIndex int) *BoneMorphOffset {
	return &BoneMorphOffset{
		BoneIndex: boneIndex,
		Extend:    &BoneMorphOffsetExtend{},
	}
}

// GroupMorphOffset represents a group morph.
type GroupMorphOffset struct {
	MorphIndex  int     // モーフINDEX
	MorphFactor float64 // モーフ変動量
}

func NewGroupMorphOffset(morphIndex int, morphFactor float64) *GroupMorphOffset {
	return &GroupMorphOffset{
		MorphIndex:  morphIndex,
		MorphFactor: morphFactor,
	}
}

func (offset *GroupMorphOffset) Type() int {
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
	Diffuse             *mmath.MVec4          // Diffuse (R,G,B,A)
	Specular            *mmath.MVec4          // SpecularColor (R,G,B, 係数)
	Ambient             *mmath.MVec3          // AmbientColor (R,G,B)
	Edge                *mmath.MVec4          // エッジ色 (R,G,B,A)
	EdgeSize            float64               // エッジサイズ
	TextureFactor       *mmath.MVec4          // テクスチャ係数 (R,G,B,A)
	SphereTextureFactor *mmath.MVec4          // スフィアテクスチャ係数 (R,G,B,A)
	ToonTextureFactor   *mmath.MVec4          // Toonテクスチャ係数 (R,G,B,A)
}

func (offset *MaterialMorphOffset) Type() int {
	return int(MORPH_TYPE_MATERIAL)
}

func NewMaterialMorphOffset(
	materialIndex int,
	calcMode MaterialMorphCalcMode,
	diffuse *mmath.MVec4,
	specular *mmath.MVec4,
	ambient *mmath.MVec3,
	edge *mmath.MVec4,
	edgeSize float64,
	textureFactor *mmath.MVec4,
	sphereTextureFactor *mmath.MVec4,
	toonTextureFactor *mmath.MVec4,
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
		index:       -1,
		name:        "",
		englishName: "",
		Panel:       MORPH_PANEL_SYSTEM,
		MorphType:   MORPH_TYPE_VERTEX,
		Offsets:     make([]IMorphOffset, 0),
		DisplaySlot: -1,
		IsSystem:    false,
	}
}

func (morph *Morph) Copy() core.IIndexNameModel {
	copiedOffsets := make([]IMorphOffset, len(morph.Offsets))
	for i, offset := range morph.Offsets {
		switch offset.Type() {
		case int(MORPH_TYPE_VERTEX):
			copiedOffsets[i] = NewVertexMorphOffset(offset.(*VertexMorphOffset).VertexIndex, offset.(*VertexMorphOffset).Position)
		case int(MORPH_TYPE_UV), int(MORPH_TYPE_EXTENDED_UV1), int(MORPH_TYPE_EXTENDED_UV2), int(MORPH_TYPE_EXTENDED_UV3), int(MORPH_TYPE_EXTENDED_UV4):
			copiedOffsets[i] = NewUvMorphOffset(offset.(*UvMorphOffset).VertexIndex, offset.(*UvMorphOffset).Uv)
		case int(MORPH_TYPE_BONE):
			copiedOffsets[i] = NewBoneMorphOffset(offset.(*BoneMorphOffset).BoneIndex)
		case int(MORPH_TYPE_GROUP):
			copiedOffsets[i] = NewGroupMorphOffset(offset.(*GroupMorphOffset).MorphIndex, offset.(*GroupMorphOffset).MorphFactor)
		case int(MORPH_TYPE_MATERIAL):
			copiedOffsets[i] = NewMaterialMorphOffset(
				offset.(*MaterialMorphOffset).MaterialIndex,
				offset.(*MaterialMorphOffset).CalcMode,
				offset.(*MaterialMorphOffset).Diffuse.Copy(),
				offset.(*MaterialMorphOffset).Specular.Copy(),
				offset.(*MaterialMorphOffset).Ambient.Copy(),
				offset.(*MaterialMorphOffset).Edge.Copy(),
				offset.(*MaterialMorphOffset).EdgeSize,
				offset.(*MaterialMorphOffset).TextureFactor.Copy(),
				offset.(*MaterialMorphOffset).SphereTextureFactor.Copy(),
				offset.(*MaterialMorphOffset).ToonTextureFactor.Copy(),
			)
		}
	}

	return &Morph{
		index:       morph.index,
		name:        morph.name,
		englishName: morph.englishName,
		Panel:       morph.Panel,
		MorphType:   morph.MorphType,
		Offsets:     copiedOffsets,
		DisplaySlot: morph.DisplaySlot,
		IsSystem:    morph.IsSystem,
	}
}

// モーフリスト
type Morphs struct {
	*core.IndexNameModels[*Morph]
}

func NewMorphs(count int) *Morphs {
	return &Morphs{
		IndexNameModels: core.NewIndexNameModels[*Morph](count, func() *Morph { return nil }),
	}
}

func (morphs *Morphs) Copy() *Morphs {
	copied := NewMorphs(len(morphs.Data))
	for i, morph := range morphs.Data {
		copied.Data[i] = morph.Copy().(*Morph)
	}
	return copied
}
