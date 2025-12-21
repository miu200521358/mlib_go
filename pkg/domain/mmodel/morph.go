package mmodel

import (
	"github.com/miu200521358/mlib_go/pkg/domain/mcore"
	"github.com/miu200521358/mlib_go/pkg/domain/mmath"
	"github.com/tiendc/go-deepcopy"
)

// MorphPanel は操作パネルを表します。
type MorphPanel byte

const (
	MORPH_PANEL_SYSTEM             MorphPanel = 0 // システム予約
	MORPH_PANEL_EYEBROW_LOWER_LEFT MorphPanel = 1 // 眉(左下)
	MORPH_PANEL_EYE_UPPER_LEFT     MorphPanel = 2 // 目(左上)
	MORPH_PANEL_LIP_UPPER_RIGHT    MorphPanel = 3 // 口(右上)
	MORPH_PANEL_OTHER_LOWER_RIGHT  MorphPanel = 4 // その他(右下)
)

// PanelName は操作パネルの名前を返します。
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

// MorphType はモーフ種類を表します。
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

// --------------------------------------------
// IMorphOffset インターフェース
// --------------------------------------------

// IMorphOffset はモーフオフセットのインターフェースです。
type IMorphOffset interface {
	// Type はモーフタイプを返します。
	Type() MorphType
	// Copy は深いコピーを作成します。
	Copy() (IMorphOffset, error)
}

// --------------------------------------------
// VertexMorphOffset
// --------------------------------------------

// VertexMorphOffset は頂点モーフオフセットを表します。
type VertexMorphOffset struct {
	VertexIndex int         // 頂点Index
	Position    *mmath.Vec3 // 座標オフセット量
}

// NewVertexMorphOffset は新しいVertexMorphOffsetを生成します。
func NewVertexMorphOffset(vertexIndex int, position *mmath.Vec3) *VertexMorphOffset {
	return &VertexMorphOffset{
		VertexIndex: vertexIndex,
		Position:    position,
	}
}

// Type はモーフタイプを返します。
func (o *VertexMorphOffset) Type() MorphType {
	return MORPH_TYPE_VERTEX
}

// Copy は深いコピーを作成します。
func (o *VertexMorphOffset) Copy() (IMorphOffset, error) {
	cp := &VertexMorphOffset{}
	if err := deepcopy.Copy(cp, o); err != nil {
		return nil, err
	}
	return cp, nil
}

// --------------------------------------------
// UvMorphOffset
// --------------------------------------------

// UvMorphOffset はUVモーフオフセットを表します。
type UvMorphOffset struct {
	VertexIndex int         // 頂点Index
	Uv          *mmath.Vec4 // UVオフセット量
}

// NewUvMorphOffset は新しいUvMorphOffsetを生成します。
func NewUvMorphOffset(vertexIndex int, uv *mmath.Vec4) *UvMorphOffset {
	return &UvMorphOffset{
		VertexIndex: vertexIndex,
		Uv:          uv,
	}
}

// Type はモーフタイプを返します。
func (o *UvMorphOffset) Type() MorphType {
	return MORPH_TYPE_UV
}

// Copy は深いコピーを作成します。
func (o *UvMorphOffset) Copy() (IMorphOffset, error) {
	cp := &UvMorphOffset{}
	if err := deepcopy.Copy(cp, o); err != nil {
		return nil, err
	}
	return cp, nil
}

// --------------------------------------------
// BoneMorphOffset
// --------------------------------------------

// BoneMorphOffset はボーンモーフオフセットを表します。
type BoneMorphOffset struct {
	BoneIndex          int               // ボーンIndex
	Position           *mmath.Vec3       // グローバル移動量
	CancelablePosition *mmath.Vec3       // 親キャンセル位置
	Rotation           *mmath.Quaternion // グローバル回転量
	CancelableRotation *mmath.Quaternion // 親キャンセル回転
	Scale              *mmath.Vec3       // グローバル縮尺量
	CancelableScale    *mmath.Vec3       // 親キャンセルスケール
	LocalMat           *mmath.Mat4       // ローカル変換行列
}

// NewBoneMorphOffset は新しいBoneMorphOffsetを生成します。
func NewBoneMorphOffset(boneIndex int) *BoneMorphOffset {
	return &BoneMorphOffset{
		BoneIndex: boneIndex,
		Position:  mmath.NewVec3(),
		Rotation:  mmath.NewQuaternion(),
	}
}

// Type はモーフタイプを返します。
func (o *BoneMorphOffset) Type() MorphType {
	return MORPH_TYPE_BONE
}

// Copy は深いコピーを作成します。
func (o *BoneMorphOffset) Copy() (IMorphOffset, error) {
	cp := &BoneMorphOffset{}
	if err := deepcopy.Copy(cp, o); err != nil {
		return nil, err
	}
	return cp, nil
}

// --------------------------------------------
// GroupMorphOffset
// --------------------------------------------

// GroupMorphOffset はグループモーフオフセットを表します。
type GroupMorphOffset struct {
	MorphIndex  int     // モーフIndex
	MorphFactor float64 // モーフ変動量
}

// NewGroupMorphOffset は新しいGroupMorphOffsetを生成します。
func NewGroupMorphOffset(morphIndex int, morphFactor float64) *GroupMorphOffset {
	return &GroupMorphOffset{
		MorphIndex:  morphIndex,
		MorphFactor: morphFactor,
	}
}

// Type はモーフタイプを返します。
func (o *GroupMorphOffset) Type() MorphType {
	return MORPH_TYPE_GROUP
}

// Copy は深いコピーを作成します。
func (o *GroupMorphOffset) Copy() (IMorphOffset, error) {
	return &GroupMorphOffset{
		MorphIndex:  o.MorphIndex,
		MorphFactor: o.MorphFactor,
	}, nil
}

// --------------------------------------------
// MaterialMorphOffset
// --------------------------------------------

// MaterialMorphCalcMode は材質モーフの計算モードを表します。
type MaterialMorphCalcMode int

const (
	CALC_MODE_MULTIPLICATION MaterialMorphCalcMode = 0 // 乗算
	CALC_MODE_ADDITION       MaterialMorphCalcMode = 1 // 加算
)

// MaterialMorphOffset は材質モーフオフセットを表します。
type MaterialMorphOffset struct {
	MaterialIndex       int                   // 材質Index（-1:全材質対象）
	CalcMode            MaterialMorphCalcMode // 計算モード（0:乗算, 1:加算）
	Diffuse             *mmath.Vec4           // Diffuse (R,G,B,A)
	Specular            *mmath.Vec4           // SpecularColor (R,G,B, 係数)
	Ambient             *mmath.Vec3           // AmbientColor (R,G,B)
	Edge                *mmath.Vec4           // エッジ色 (R,G,B,A)
	EdgeSize            float64               // エッジサイズ
	TextureFactor       *mmath.Vec4           // テクスチャ係数 (R,G,B,A)
	SphereTextureFactor *mmath.Vec4           // スフィアテクスチャ係数 (R,G,B,A)
	ToonTextureFactor   *mmath.Vec4           // Toonテクスチャ係数 (R,G,B,A)
}

// NewMaterialMorphOffset は新しいMaterialMorphOffsetを生成します。
func NewMaterialMorphOffset(
	materialIndex int,
	calcMode MaterialMorphCalcMode,
	diffuse, specular *mmath.Vec4,
	ambient *mmath.Vec3,
	edge *mmath.Vec4,
	edgeSize float64,
	textureFactor, sphereTextureFactor, toonTextureFactor *mmath.Vec4,
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

// Type はモーフタイプを返します。
func (o *MaterialMorphOffset) Type() MorphType {
	return MORPH_TYPE_MATERIAL
}

// Copy は深いコピーを作成します。
func (o *MaterialMorphOffset) Copy() (IMorphOffset, error) {
	cp := &MaterialMorphOffset{}
	if err := deepcopy.Copy(cp, o); err != nil {
		return nil, err
	}
	return cp, nil
}

// --------------------------------------------
// Morph
// --------------------------------------------

// Morph はモーフを表します。
type Morph struct {
	mcore.IndexNameModel
	Panel       MorphPanel     // モーフパネル
	MorphType   MorphType      // モーフ種類
	Offsets     []IMorphOffset // モーフオフセット
	DisplaySlot int            // 表示枠
	IsSystem    bool           // ツール側で追加したモーフ
}

// NewMorph は新しいMorphを生成します。
func NewMorph() *Morph {
	return &Morph{
		IndexNameModel: *mcore.NewIndexNameModel(-1, "", ""),
		Panel:          MORPH_PANEL_SYSTEM,
		MorphType:      MORPH_TYPE_VERTEX,
		Offsets:        make([]IMorphOffset, 0),
		DisplaySlot:    -1,
		IsSystem:       false,
	}
}

// IsValid はMorphが有効かどうかを返します。
func (m *Morph) IsValid() bool {
	return m != nil && m.Index() >= 0
}

// Copy は深いコピーを作成します。
func (m *Morph) Copy() (*Morph, error) {
	cp := &Morph{
		IndexNameModel: *mcore.NewIndexNameModel(m.Index(), m.Name(), m.EnglishName()),
		Panel:          m.Panel,
		MorphType:      m.MorphType,
		Offsets:        make([]IMorphOffset, len(m.Offsets)),
		DisplaySlot:    m.DisplaySlot,
		IsSystem:       m.IsSystem,
	}

	for i, offset := range m.Offsets {
		copiedOffset, err := offset.Copy()
		if err != nil {
			return nil, err
		}
		cp.Offsets[i] = copiedOffset
	}

	return cp, nil
}
