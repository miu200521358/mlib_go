package pmx

import (
	"github.com/miu200521358/mlib_go/pkg/domain/core"
	"github.com/miu200521358/mlib_go/pkg/domain/mmath"
)

// スフィアモード
type SphereMode byte

const (
	// 無効
	SPHERE_MODE_INVALID SphereMode = 0
	// 乗算(sph)
	SPHERE_MODE_MULTIPLICATION SphereMode = 1
	// 加算(spa)
	SPHERE_MODE_ADDITION SphereMode = 2
	// サブテクスチャ(追加UV1のx,yをUV参照して通常テクスチャ描画を行う)
	SPHERE_MODE_SUBTEXTURE SphereMode = 3
)

type DrawFlag byte

const (
	// 初期値
	DRAW_FLAG_NONE DrawFlag = 0x0000
	// 0x01:両面描画
	DRAW_FLAG_DOUBLE_SIDED_DRAWING DrawFlag = 0x0001
	// 0x02:地面影
	DRAW_FLAG_GROUND_SHADOW DrawFlag = 0x0002
	// 0x04:セルフシャドウマップへの描画
	DRAW_FLAG_DRAWING_ON_SELF_SHADOW_MAPS DrawFlag = 0x0004
	// 0x08:セルフシャドウの描画
	DRAW_FLAG_DRAWING_SELF_SHADOWS DrawFlag = 0x0008
	// 0x10:エッジ描画
	DRAW_FLAG_DRAWING_EDGE DrawFlag = 0x0010
)

// 0x01:両面描画
func (drawFlg DrawFlag) IsDoubleSidedDrawing() bool {
	return drawFlg&DRAW_FLAG_DOUBLE_SIDED_DRAWING == DRAW_FLAG_DOUBLE_SIDED_DRAWING
}

// 0x02:地面影
func (drawFlg DrawFlag) IsGroundShadow() bool {
	return drawFlg&DRAW_FLAG_GROUND_SHADOW == DRAW_FLAG_GROUND_SHADOW
}

// 0x04:セルフシャドウマップへの描画
func (drawFlg DrawFlag) IsDrawingOnSelfShadowMaps() bool {
	return drawFlg&DRAW_FLAG_DRAWING_ON_SELF_SHADOW_MAPS == DRAW_FLAG_DRAWING_ON_SELF_SHADOW_MAPS
}

// 0x08:セルフシャドウの描画
func (drawFlg DrawFlag) IsDrawingSelfShadows() bool {
	return drawFlg&DRAW_FLAG_DRAWING_SELF_SHADOWS == DRAW_FLAG_DRAWING_SELF_SHADOWS
}

// 0x10:エッジ描画
func (drawFlg DrawFlag) IsDrawingEdge() bool {
	return drawFlg&DRAW_FLAG_DRAWING_EDGE == DRAW_FLAG_DRAWING_EDGE
}

func (drawFlg DrawFlag) SetDrawingEdge(flag bool) DrawFlag {
	if flag {
		return drawFlg | DRAW_FLAG_DRAWING_EDGE
	}
	return drawFlg &^ DRAW_FLAG_DRAWING_EDGE
}

func (drawFlg DrawFlag) SetDoubleSidedDrawing(flag bool) DrawFlag {
	if flag {
		return drawFlg | DRAW_FLAG_DOUBLE_SIDED_DRAWING
	}
	return drawFlg &^ DRAW_FLAG_DOUBLE_SIDED_DRAWING
}

func (drawFlg DrawFlag) SetGroundShadow(flag bool) DrawFlag {
	if flag {
		return drawFlg | DRAW_FLAG_GROUND_SHADOW
	}
	return drawFlg &^ DRAW_FLAG_GROUND_SHADOW
}

func (drawFlg DrawFlag) SetDrawingOnSelfShadowMaps(flag bool) DrawFlag {
	if flag {
		return drawFlg | DRAW_FLAG_DRAWING_ON_SELF_SHADOW_MAPS
	}
	return drawFlg &^ DRAW_FLAG_DRAWING_ON_SELF_SHADOW_MAPS
}

func (drawFlg DrawFlag) SetDrawingSelfShadows(flag bool) DrawFlag {
	if flag {
		return drawFlg | DRAW_FLAG_DRAWING_SELF_SHADOWS
	}
	return drawFlg &^ DRAW_FLAG_DRAWING_SELF_SHADOWS
}

// 共有Toonフラグ
type ToonSharing byte

const (
	// 0:継続値は個別Toon
	TOON_SHARING_INDIVIDUAL ToonSharing = 0
	// 1:継続値は共有Toon
	TOON_SHARING_SHARING ToonSharing = 1
)

type Material struct {
	index              int          // 材質INDEX
	name               string       // 材質名
	englishName        string       // 材質英名
	Diffuse            *mmath.MVec4 // Diffuse (R,G,B,A)(拡散色＋非透過度)
	Specular           *mmath.MVec4 // Specular (R,G,B,A)(反射色 + 反射強度)
	Ambient            *mmath.MVec3 // Ambient (R,G,B)(環境色)
	DrawFlag           DrawFlag     // 描画フラグ(8bit) - 各bit 0:OFF 1:ON
	Edge               *mmath.MVec4 // エッジ色 (R,G,B,A)
	EdgeSize           float64      // エッジサイズ
	TextureIndex       int          // 通常テクスチャINDEX
	SphereTextureIndex int          // スフィアテクスチャINDEX
	SphereMode         SphereMode   // スフィアモード
	ToonSharingFlag    ToonSharing  // 共有Toonフラグ
	ToonTextureIndex   int          // ToonテクスチャINDEX
	Memo               string       // メモ
	VerticesCount      int          // 材質に対応する面(頂点)数 (必ず3の倍数になる)
}

func NewMaterial() *Material {
	return &Material{
		index:              -1,
		name:               "",
		englishName:        "",
		Diffuse:            mmath.NewMVec4(),
		Specular:           mmath.NewMVec4(),
		Ambient:            mmath.NewMVec3(),
		DrawFlag:           DRAW_FLAG_NONE,
		Edge:               mmath.NewMVec4(),
		EdgeSize:           0.0,
		TextureIndex:       -1,
		SphereTextureIndex: -1,
		SphereMode:         SPHERE_MODE_INVALID,
		ToonSharingFlag:    TOON_SHARING_INDIVIDUAL,
		ToonTextureIndex:   -1,
		Memo:               "",
		VerticesCount:      0,
	}
}

func NewMaterialByName(name string) *Material {
	material := NewMaterial()
	material.SetName(name)
	return material
}

func (material *Material) Index() int {
	return material.index
}

func (material *Material) SetIndex(index int) {
	material.index = index
}

func (material *Material) Name() string {
	return material.name
}

func (material *Material) SetName(name string) {
	material.name = name
}

func (material *Material) EnglishName() string {
	return material.englishName
}

func (material *Material) SetEnglishName(englishName string) {
	material.englishName = englishName
}

func (material *Material) IsValid() bool {
	return material != nil && material.index >= 0
}

func (material *Material) Copy() core.IIndexNameModel {
	return &Material{
		index:              material.index,
		name:               material.name,
		englishName:        material.englishName,
		Diffuse:            material.Diffuse.Copy(),
		Specular:           material.Specular.Copy(),
		Ambient:            material.Ambient.Copy(),
		DrawFlag:           material.DrawFlag,
		Edge:               material.Edge.Copy(),
		EdgeSize:           material.EdgeSize,
		TextureIndex:       material.TextureIndex,
		SphereTextureIndex: material.SphereTextureIndex,
		SphereMode:         material.SphereMode,
		ToonSharingFlag:    material.ToonSharingFlag,
		ToonTextureIndex:   material.ToonTextureIndex,
		Memo:               material.Memo,
		VerticesCount:      material.VerticesCount,
	}
}
