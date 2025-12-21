package mmodel

import (
	"github.com/miu200521358/mlib_go/pkg/domain/mcore"
	"github.com/miu200521358/mlib_go/pkg/domain/mmath"
	"github.com/tiendc/go-deepcopy"
)

// SphereMode はスフィアモードを表します。
type SphereMode byte

const (
	SPHERE_MODE_INVALID        SphereMode = 0 // 無効
	SPHERE_MODE_MULTIPLICATION SphereMode = 1 // 乗算(sph)
	SPHERE_MODE_ADDITION       SphereMode = 2 // 加算(spa)
	SPHERE_MODE_SUBTEXTURE     SphereMode = 3 // サブテクスチャ
)

// DrawFlag は描画フラグを表します。
type DrawFlag byte

const (
	DRAW_FLAG_NONE                        DrawFlag = 0x0000 // 初期値
	DRAW_FLAG_DOUBLE_SIDED_DRAWING        DrawFlag = 0x0001 // 両面描画
	DRAW_FLAG_GROUND_SHADOW               DrawFlag = 0x0002 // 地面影
	DRAW_FLAG_DRAWING_ON_SELF_SHADOW_MAPS DrawFlag = 0x0004 // セルフシャドウマップへの描画
	DRAW_FLAG_DRAWING_SELF_SHADOWS        DrawFlag = 0x0008 // セルフシャドウの描画
	DRAW_FLAG_DRAWING_EDGE                DrawFlag = 0x0010 // エッジ描画
)

// IsDoubleSidedDrawing は両面描画かどうかを返します。
func (f DrawFlag) IsDoubleSidedDrawing() bool {
	return f&DRAW_FLAG_DOUBLE_SIDED_DRAWING != 0
}

// IsGroundShadow は地面影かどうかを返します。
func (f DrawFlag) IsGroundShadow() bool {
	return f&DRAW_FLAG_GROUND_SHADOW != 0
}

// IsDrawingOnSelfShadowMaps はセルフシャドウマップへの描画かどうかを返します。
func (f DrawFlag) IsDrawingOnSelfShadowMaps() bool {
	return f&DRAW_FLAG_DRAWING_ON_SELF_SHADOW_MAPS != 0
}

// IsDrawingSelfShadows はセルフシャドウの描画かどうかを返します。
func (f DrawFlag) IsDrawingSelfShadows() bool {
	return f&DRAW_FLAG_DRAWING_SELF_SHADOWS != 0
}

// IsDrawingEdge はエッジ描画かどうかを返します。
func (f DrawFlag) IsDrawingEdge() bool {
	return f&DRAW_FLAG_DRAWING_EDGE != 0
}

// SetDoubleSidedDrawing は両面描画フラグを設定します。
func (f DrawFlag) SetDoubleSidedDrawing(on bool) DrawFlag {
	if on {
		return f | DRAW_FLAG_DOUBLE_SIDED_DRAWING
	}
	return f &^ DRAW_FLAG_DOUBLE_SIDED_DRAWING
}

// SetGroundShadow は地面影フラグを設定します。
func (f DrawFlag) SetGroundShadow(on bool) DrawFlag {
	if on {
		return f | DRAW_FLAG_GROUND_SHADOW
	}
	return f &^ DRAW_FLAG_GROUND_SHADOW
}

// SetDrawingOnSelfShadowMaps はセルフシャドウマップへの描画フラグを設定します。
func (f DrawFlag) SetDrawingOnSelfShadowMaps(on bool) DrawFlag {
	if on {
		return f | DRAW_FLAG_DRAWING_ON_SELF_SHADOW_MAPS
	}
	return f &^ DRAW_FLAG_DRAWING_ON_SELF_SHADOW_MAPS
}

// SetDrawingSelfShadows はセルフシャドウの描画フラグを設定します。
func (f DrawFlag) SetDrawingSelfShadows(on bool) DrawFlag {
	if on {
		return f | DRAW_FLAG_DRAWING_SELF_SHADOWS
	}
	return f &^ DRAW_FLAG_DRAWING_SELF_SHADOWS
}

// SetDrawingEdge はエッジ描画フラグを設定します。
func (f DrawFlag) SetDrawingEdge(on bool) DrawFlag {
	if on {
		return f | DRAW_FLAG_DRAWING_EDGE
	}
	return f &^ DRAW_FLAG_DRAWING_EDGE
}

// ToonSharing は共有Toonフラグを表します。
type ToonSharing byte

const (
	TOON_SHARING_INDIVIDUAL ToonSharing = 0 // 個別Toon
	TOON_SHARING_SHARED     ToonSharing = 1 // 共有Toon
)

// Material は材質を表します。
type Material struct {
	mcore.IndexNameModel             // インデックス・名前
	Diffuse              *mmath.Vec4 // 拡散色＋非透過度
	Specular             *mmath.Vec4 // 反射色＋反射強度
	Ambient              *mmath.Vec3 // 環境色
	DrawFlag             DrawFlag    // 描画フラグ
	Edge                 *mmath.Vec4 // エッジ色
	EdgeSize             float64     // エッジサイズ
	TextureIndex         int         // 通常テクスチャINDEX
	SphereTextureIndex   int         // スフィアテクスチャINDEX
	SphereMode           SphereMode  // スフィアモード
	ToonSharingFlag      ToonSharing // 共有Toonフラグ
	ToonTextureIndex     int         // ToonテクスチャINDEX
	Memo                 string      // メモ
	VerticesCount        int         // 材質に対応する面(頂点)数
}

// NewMaterial は新しいMaterialを生成します。
func NewMaterial() *Material {
	return &Material{
		IndexNameModel:     *mcore.NewIndexNameModel(-1, "", ""),
		Diffuse:            mmath.NewVec4(),
		Specular:           mmath.NewVec4(),
		Ambient:            mmath.NewVec3(),
		DrawFlag:           DRAW_FLAG_NONE,
		Edge:               mmath.NewVec4(),
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

// IsValid はMaterialが有効かどうかを返します。
func (m *Material) IsValid() bool {
	return m != nil && m.IndexNameModel.IsValid()
}

// Copy は深いコピーを作成します。
func (m *Material) Copy() (*Material, error) {
	cp := &Material{}
	if err := deepcopy.Copy(cp, m); err != nil {
		return nil, err
	}
	return cp, nil
}
