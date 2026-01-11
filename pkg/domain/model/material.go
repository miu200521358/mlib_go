package model

import "github.com/miu200521358/mlib_go/pkg/domain/mmath"

// SphereMode はスフィアテクスチャの合成モードを表す。
type SphereMode int

const (
	// SPHERE_MODE_INVALID は無効モード。
	SPHERE_MODE_INVALID SphereMode = iota
	// SPHERE_MODE_MULTIPLICATION は乗算モード。
	SPHERE_MODE_MULTIPLICATION
	// SPHERE_MODE_ADDITION は加算モード。
	SPHERE_MODE_ADDITION
	// SPHERE_MODE_SUBTEXTURE はサブテクスチャモード。
	SPHERE_MODE_SUBTEXTURE
)

// DrawFlag は材質の描画フラグを表す。
type DrawFlag int

const (
	// DRAW_FLAG_NONE はフラグなし。
	DRAW_FLAG_NONE DrawFlag = 0
	// DRAW_FLAG_DOUBLE_SIDED_DRAWING は両面描画。
	DRAW_FLAG_DOUBLE_SIDED_DRAWING DrawFlag = 0x0001
	// DRAW_FLAG_GROUND_SHADOW は地面影。
	DRAW_FLAG_GROUND_SHADOW DrawFlag = 0x0002
	// DRAW_FLAG_DRAWING_ON_SELF_SHADOW_MAPS はセルフシャドウマップ描画。
	DRAW_FLAG_DRAWING_ON_SELF_SHADOW_MAPS DrawFlag = 0x0004
	// DRAW_FLAG_DRAWING_SELF_SHADOWS はセルフシャドウ描画。
	DRAW_FLAG_DRAWING_SELF_SHADOWS DrawFlag = 0x0008
	// DRAW_FLAG_DRAWING_EDGE はエッジ描画。
	DRAW_FLAG_DRAWING_EDGE DrawFlag = 0x0010
)

// ToonSharingFlag はトゥーン共有モードを表す。
type ToonSharingFlag int

const (
	// TOON_SHARING_INDIVIDUAL は個別トゥーン。
	TOON_SHARING_INDIVIDUAL ToonSharingFlag = iota
	// TOON_SHARING_SHARING は共有トゥーン。
	TOON_SHARING_SHARING
)

// Material は材質要素を表す。
type Material struct {
	index              int
	name               string
	EnglishName        string
	Memo               string
	Diffuse            mmath.Vec4
	Specular           mmath.Vec4
	Ambient            mmath.Vec3
	DrawFlag           DrawFlag
	Edge               mmath.Vec4
	EdgeSize           float64
	TextureIndex       int
	SphereTextureIndex int
	SphereMode         SphereMode
	ToonSharingFlag    ToonSharingFlag
	ToonTextureIndex   int
	VerticesCount      int
}

// NewMaterial は既定値で Material を生成する。
func NewMaterial() *Material {
	return &Material{
		index:              -1,
		name:               "",
		EnglishName:        "",
		Memo:               "",
		DrawFlag:           DRAW_FLAG_NONE,
		EdgeSize:           0.0,
		TextureIndex:       -1,
		SphereTextureIndex: -1,
		SphereMode:         SPHERE_MODE_INVALID,
		ToonSharingFlag:    TOON_SHARING_INDIVIDUAL,
		ToonTextureIndex:   -1,
		VerticesCount:      0,
	}
}

// Index は材質 index を返す。
func (m *Material) Index() int {
	return m.index
}

// SetIndex は材質 index を設定する。
func (m *Material) SetIndex(index int) {
	m.index = index
}

// Name は材質名を返す。
func (m *Material) Name() string {
	return m.name
}

// SetName は材質名を設定する。
func (m *Material) SetName(name string) {
	m.name = name
}

// IsValid は材質が有効か判定する。
func (m *Material) IsValid() bool {
	return m != nil && m.index >= 0
}
