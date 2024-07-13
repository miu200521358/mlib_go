package pmx

import (
	"slices"

	"github.com/jinzhu/copier"

	"github.com/miu200521358/mlib_go/pkg/domain/core"
	"github.com/miu200521358/mlib_go/pkg/mmath"
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
func (d DrawFlag) IsDoubleSidedDrawing() bool {
	return d&DRAW_FLAG_DOUBLE_SIDED_DRAWING == DRAW_FLAG_DOUBLE_SIDED_DRAWING
}

// 0x02:地面影
func (d DrawFlag) IsGroundShadow() bool {
	return d&DRAW_FLAG_GROUND_SHADOW == DRAW_FLAG_GROUND_SHADOW
}

// 0x04:セルフシャドウマップへの描画
func (d DrawFlag) IsDrawingOnSelfShadowMaps() bool {
	return d&DRAW_FLAG_DRAWING_ON_SELF_SHADOW_MAPS == DRAW_FLAG_DRAWING_ON_SELF_SHADOW_MAPS
}

// 0x08:セルフシャドウの描画
func (d DrawFlag) IsDrawingSelfShadows() bool {
	return d&DRAW_FLAG_DRAWING_SELF_SHADOWS == DRAW_FLAG_DRAWING_SELF_SHADOWS
}

// 0x10:エッジ描画
func (d DrawFlag) IsDrawingEdge() bool {
	return d&DRAW_FLAG_DRAWING_EDGE == DRAW_FLAG_DRAWING_EDGE
}

func (d DrawFlag) SetDrawingEdge(flag bool) DrawFlag {
	if flag {
		return d | DRAW_FLAG_DRAWING_EDGE
	}
	return d &^ DRAW_FLAG_DRAWING_EDGE
}

func (d DrawFlag) SetDoubleSidedDrawing(flag bool) DrawFlag {
	if flag {
		return d | DRAW_FLAG_DOUBLE_SIDED_DRAWING
	}
	return d &^ DRAW_FLAG_DOUBLE_SIDED_DRAWING
}

func (d DrawFlag) SetGroundShadow(flag bool) DrawFlag {
	if flag {
		return d | DRAW_FLAG_GROUND_SHADOW
	}
	return d &^ DRAW_FLAG_GROUND_SHADOW
}

func (d DrawFlag) SetDrawingOnSelfShadowMaps(flag bool) DrawFlag {
	if flag {
		return d | DRAW_FLAG_DRAWING_ON_SELF_SHADOW_MAPS
	}
	return d &^ DRAW_FLAG_DRAWING_ON_SELF_SHADOW_MAPS
}

func (d DrawFlag) SetDrawingSelfShadows(flag bool) DrawFlag {
	if flag {
		return d | DRAW_FLAG_DRAWING_SELF_SHADOWS
	}
	return d &^ DRAW_FLAG_DRAWING_SELF_SHADOWS
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
	*core.IndexNameModel
	Diffuse             *mmath.MVec4 // Diffuse (R,G,B,A)(拡散色＋非透過度)
	Specular            *mmath.MVec4 // Specular (R,G,B,A)(反射色 + 反射強度)
	Ambient             *mmath.MVec3 // Ambient (R,G,B)(環境色)
	DrawFlag            DrawFlag     // 描画フラグ(8bit) - 各bit 0:OFF 1:ON
	Edge                *mmath.MVec4 // エッジ色 (R,G,B,A)
	EdgeSize            float64      // エッジサイズ
	TextureIndex        int          // 通常テクスチャINDEX
	SphereTextureIndex  int          // スフィアテクスチャINDEX
	SphereMode          SphereMode   // スフィアモード
	ToonSharingFlag     ToonSharing  // 共有Toonフラグ
	ToonTextureIndex    int          // ToonテクスチャINDEX
	Memo                string       // メモ
	VerticesCount       int          // 材質に対応する面(頂点)数 (必ず3の倍数になる)
	TextureFactor       *mmath.MVec4 // テクスチャ係数
	SphereTextureFactor *mmath.MVec4 // スフィアテクスチャ係数
	ToonTextureFactor   *mmath.MVec4 // トゥーンテクスチャ係数
}

func NewMaterial() *Material {
	return &Material{
		IndexNameModel:      &core.IndexNameModel{Index: -1, Name: "", EnglishName: ""},
		Diffuse:             mmath.NewMVec4(),
		Specular:            mmath.NewMVec4(),
		Ambient:             mmath.NewMVec3(),
		DrawFlag:            DRAW_FLAG_NONE,
		Edge:                mmath.NewMVec4(),
		EdgeSize:            0.0,
		TextureIndex:        -1,
		SphereTextureIndex:  -1,
		SphereMode:          SPHERE_MODE_INVALID,
		ToonSharingFlag:     TOON_SHARING_INDIVIDUAL,
		ToonTextureIndex:    -1,
		Memo:                "",
		VerticesCount:       0,
		TextureFactor:       mmath.NewMVec4(),
		SphereTextureFactor: mmath.NewMVec4(),
		ToonTextureFactor:   mmath.NewMVec4(),
	}
}

func NewMaterialByName(name string) *Material {
	m := NewMaterial()
	m.Name = name
	return m
}

// 材質リスト
type Materials struct {
	*core.IndexNameModels[*Material]
	Vertices map[int][]int
	Faces    map[int][]int
}

func NewMaterials(count int) *Materials {
	return &Materials{
		IndexNameModels: core.NewIndexNameModels[*Material](count, func() *Material { return nil }),
		Vertices:        make(map[int][]int),
		Faces:           make(map[int][]int),
	}
}

func (m *Material) Copy() core.IIndexNameModel {
	copied := NewMaterial()
	copier.CopyWithOption(copied, m, copier.Option{DeepCopy: true})
	return copied
}

func (ms *Materials) setup(vertices *Vertices, faces *Faces) {
	prevVertexCount := 0

	for _, v := range vertices.Data {
		v.MaterialIndexes = make([]int, 0)
	}

	for _, m := range ms.Data {
		for j := prevVertexCount; j < prevVertexCount+int(m.VerticesCount/3); j++ {
			face := faces.Get(j)
			for _, vertexIndexes := range face.VertexIndexes {
				if !slices.Contains(vertices.Get(vertexIndexes).MaterialIndexes, m.Index) {
					vertices.Get(vertexIndexes).MaterialIndexes =
						append(vertices.Get(vertexIndexes).MaterialIndexes, m.Index)
				}
			}
		}

		prevVertexCount += int(m.VerticesCount / 3)
	}
}
