package pmx

import (
	"embed"

	"github.com/go-gl/mathgl/mgl32"

	"github.com/miu200521358/mlib_go/pkg/mcore"
	"github.com/miu200521358/mlib_go/pkg/mgl"
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

// 共有Toonフラグ
type ToonSharing byte

const (
	// 0:継続値は個別Toon
	TOON_SHARING_INDIVIDUAL ToonSharing = 0
	// 1:継続値は共有Toon
	TOON_SHARING_SHARING ToonSharing = 1
)

type MaterialGL struct {
	diffuse           *mmath.MVec4 // Diffuse (R,G,B,A)(拡散色＋非透過度)
	specular          *mmath.MVec4 // Specular (R,G,B,A)(反射色 + 反射強度)
	ambient           *mmath.MVec3 // Ambient (R,G,B)(環境色)
	edge              *mmath.MVec4 // エッジ色 (R,G,B,A)
	EdgeSize          float32      // エッジサイズ
	DrawFlag          DrawFlag     // 描画フラグ(8bit) - 各bit 0:OFF 1:ON
	Texture           *TextureGL   // 通常テクスチャ
	TextureFactor     *mmath.MVec4 // テクスチャ係数
	SphereTexture     *TextureGL   // スフィアテクスチャ
	ToonTexture       *TextureGL   // トゥーンテクスチャ
	SphereMode        SphereMode   // スフィアモード
	VerticesCount     int          // 頂点数
	PrevVerticesCount int          // 前の材質までの頂点数
	lightAmbient      *mmath.MVec4
}

func (m *MaterialGL) Diffuse() mgl32.Vec4 {
	diffuse := mgl32.Vec4{
		float32(m.diffuse.GetX())*float32(m.lightAmbient.GetX()) + float32(m.ambient.GetX()),
		float32(m.diffuse.GetY())*float32(m.lightAmbient.GetY()) + float32(m.ambient.GetY()),
		float32(m.diffuse.GetZ())*float32(m.lightAmbient.GetZ()) + float32(m.ambient.GetZ()),
		float32(m.diffuse.GetW()),
	}
	return diffuse
}

func (m *MaterialGL) Ambient() mgl32.Vec3 {
	ambient := mgl32.Vec3{
		float32(m.diffuse.GetX()) * float32(m.lightAmbient.GetX()),
		float32(m.diffuse.GetY()) * float32(m.lightAmbient.GetY()),
		float32(m.diffuse.GetZ()) * float32(m.lightAmbient.GetZ()),
	}
	return ambient
}

func (m *MaterialGL) Specular() mgl32.Vec4 {
	specular := mgl32.Vec4{
		float32(m.specular.GetX()) * float32(m.lightAmbient.GetX()),
		float32(m.specular.GetY()) * float32(m.lightAmbient.GetY()),
		float32(m.specular.GetZ()) * float32(m.lightAmbient.GetZ()),
		float32(m.specular.GetW()),
	}
	return specular
}

func (m *MaterialGL) Edge() [4]float32 {
	edge := [4]float32{
		float32(m.edge.GetX()),
		float32(m.edge.GetY()),
		float32(m.edge.GetZ()),
		float32(m.edge.GetW()) * float32(m.diffuse.GetW()),
	}
	return edge
}

type Material struct {
	*mcore.IndexModel
	Name               string       // 材質名
	EnglishName        string       // 材質名英
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
		IndexModel:         &mcore.IndexModel{Index: -1},
		Name:               "",
		EnglishName:        "",
		Diffuse:            &mmath.MVec4{},
		Specular:           &mmath.MVec4{},
		Ambient:            &mmath.MVec3{},
		DrawFlag:           DRAW_FLAG_NONE,
		Edge:               &mmath.MVec4{},
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

func (m *Material) GL(
	modelPath string,
	texture *Texture,
	toonTexture *Texture,
	sphereTexture *Texture,
	windowIndex int,
	prevVerticesCount int,
	resourceFiles embed.FS,
) *MaterialGL {
	var textureGL *TextureGL
	if texture != nil {
		textureGL = texture.GL(modelPath, TEXTURE_TYPE_TEXTURE, windowIndex, resourceFiles)
	}

	var sphereTextureGL *TextureGL
	if sphereTexture != nil {
		sphereTextureGL = sphereTexture.GL(modelPath, TEXTURE_TYPE_SPHERE, windowIndex, resourceFiles)
	}

	var tooTextureGL *TextureGL
	if toonTexture != nil {
		tooTextureGL = toonTexture.GL(modelPath, TEXTURE_TYPE_TOON, windowIndex, resourceFiles)
	}

	return &MaterialGL{
		diffuse:           m.Diffuse,
		ambient:           m.Ambient,
		specular:          m.Specular,
		edge:              m.Edge,
		EdgeSize:          float32(m.EdgeSize),
		Texture:           textureGL,
		SphereTexture:     sphereTextureGL,
		ToonTexture:       tooTextureGL,
		DrawFlag:          m.DrawFlag,
		SphereMode:        m.SphereMode,
		VerticesCount:     m.VerticesCount,
		PrevVerticesCount: prevVerticesCount * 4,
		lightAmbient:      &mmath.MVec4{mgl.LIGHT_AMBIENT, mgl.LIGHT_AMBIENT, mgl.LIGHT_AMBIENT, 1},
	}
}

func (m *Material) Copy() mcore.IndexModelInterface {
	copied := *m
	return &copied
}

// 材質リスト
type Materials struct {
	*mcore.IndexModelCorrection[*Material]
}

func NewMaterials() *Materials {
	return &Materials{
		IndexModelCorrection: mcore.NewIndexModelCorrection[*Material](),
	}
}

// // シェーダー用材質
// type ShaderMaterial struct {
// 	LightAmbient4             *mmath.MVec4
// 	Material                  *Material
// 	ShaderTextureFactor       *mmath.MVec4
// 	SphereShaderTextureFactor *mmath.MVec4
// 	ToonShaderTextureFactor   *mmath.MVec4
// }

// func NewShaderMaterial(
// 	material *Material,
// 	lightAmbient4 *mmath.MVec4,
// 	textureFactor *mmath.MVec4,
// 	toonTextureFactor *mmath.MVec4,
// 	sphereTextureFactor *mmath.MVec4,
// ) *ShaderMaterial {
// 	return &ShaderMaterial{
// 		LightAmbient4:             lightAmbient4,
// 		Material:                  material.Copy().(*Material),
// 		ShaderTextureFactor:       textureFactor,
// 		SphereShaderTextureFactor: toonTextureFactor,
// 		ToonShaderTextureFactor:   sphereTextureFactor,
// 	}
// }

// func (sm *ShaderMaterial) Diffuse() []float32 {
// 	diffuse := make([]float32, 3)
// 	diffuse[0] = float32(sm.Material.Diffuse.GetX())*float32(sm.LightAmbient4.GetX()) + float32(sm.Material.Ambient.GetX())
// 	diffuse[1] = float32(sm.Material.Diffuse.GetY())*float32(sm.LightAmbient4.GetY()) + float32(sm.Material.Ambient.GetY())
// 	diffuse[2] = float32(sm.Material.Diffuse.GetZ())*float32(sm.LightAmbient4.GetZ()) + float32(sm.Material.Ambient.GetZ())
// 	return diffuse
// }

// func (sm *ShaderMaterial) Ambient() []float32 {
// 	ambient := make([]float32, 3)
// 	ambient[0] = float32(sm.Material.Diffuse.GetX()) * float32(sm.LightAmbient4.GetX())
// 	ambient[1] = float32(sm.Material.Diffuse.GetY()) * float32(sm.LightAmbient4.GetY())
// 	ambient[2] = float32(sm.Material.Diffuse.GetZ()) * float32(sm.LightAmbient4.GetZ())
// 	return ambient
// }

// func (sm *ShaderMaterial) Specular() []float32 {
// 	specular := make([]float32, 4)
// 	specular[0] = float32(sm.Material.Specular.GetX()) * float32(sm.LightAmbient4.GetX())
// 	specular[1] = float32(sm.Material.Specular.GetY()) * float32(sm.LightAmbient4.GetY())
// 	specular[2] = float32(sm.Material.Specular.GetZ()) * float32(sm.LightAmbient4.GetZ())
// 	specular[3] = float32(sm.Material.Specular.GetW()) * float32(sm.LightAmbient4.GetW())
// 	return specular
// }

// func (sm *ShaderMaterial) Edge() []float32 {
// 	edgeColor := make([]float32, 3)
// 	edgeColor[0] = float32(sm.Material.Edge.GetX()) * float32(sm.Material.Diffuse.GetW())
// 	edgeColor[1] = float32(sm.Material.Edge.GetY()) * float32(sm.Material.Diffuse.GetW())
// 	edgeColor[2] = float32(sm.Material.Edge.GetZ()) * float32(sm.Material.Diffuse.GetW())
// 	edgeColor[3] = float32(sm.Material.Edge.GetW()) * float32(sm.Material.Diffuse.GetW())
// 	return edgeColor
// }

// func (sm *ShaderMaterial) EdgeSize() float32 {
// 	return float32(sm.Material.EdgeSize)
// }

// func (sm *ShaderMaterial) TextureFactor() []float32 {
// 	textureFactor := make([]float32, 3)
// 	textureFactor[0] = float32(sm.ShaderTextureFactor.GetX())
// 	textureFactor[1] = float32(sm.ShaderTextureFactor.GetY())
// 	textureFactor[2] = float32(sm.ShaderTextureFactor.GetZ())
// 	textureFactor[3] = float32(sm.ShaderTextureFactor.GetW())
// 	return textureFactor
// }

// func (sm *ShaderMaterial) SphereTextureFactor() []float32 {
// 	sphereTextureFactor := make([]float32, 3)
// 	sphereTextureFactor[0] = float32(sm.SphereShaderTextureFactor.GetX())
// 	sphereTextureFactor[1] = float32(sm.SphereShaderTextureFactor.GetY())
// 	sphereTextureFactor[2] = float32(sm.SphereShaderTextureFactor.GetZ())
// 	sphereTextureFactor[3] = float32(sm.SphereShaderTextureFactor.GetW())
// 	return sphereTextureFactor
// }

// func (sm *ShaderMaterial) ToonTextureFactor() []float32 {
// 	toonTextureFactor := make([]float32, 3)
// 	toonTextureFactor[0] = float32(sm.ToonShaderTextureFactor.GetX())
// 	toonTextureFactor[1] = float32(sm.ToonShaderTextureFactor.GetY())
// 	toonTextureFactor[2] = float32(sm.ToonShaderTextureFactor.GetZ())
// 	toonTextureFactor[3] = float32(sm.ToonShaderTextureFactor.GetW())
// 	return toonTextureFactor
// }

// func (sm *ShaderMaterial) IMul(v interface{}) {
// 	switch v := v.(type) {
// 	case float64:
// 		sm.Material.Diffuse.MulScalar(v)
// 		sm.Material.Ambient.MulScalar(v)
// 		sm.Material.Specular.MulScalar(v)
// 		sm.Material.EdgeSize *= v
// 		sm.ShaderTextureFactor.MulScalar(v)
// 		sm.SphereShaderTextureFactor.MulScalar(v)
// 		sm.ToonShaderTextureFactor.MulScalar(v)
// 	case int:
// 		sm.IMul(float32(v))
// 	case *ShaderMaterial:
// 		sm.Material.Diffuse.Mul(&v.Material.Diffuse)
// 		sm.Material.Ambient.Mul(&v.Material.Ambient)
// 		sm.Material.Specular.Mul(&v.Material.Specular)
// 		sm.Material.EdgeSize *= v.Material.EdgeSize
// 		sm.ShaderTextureFactor.Mul(v.ShaderTextureFactor)
// 		sm.SphereShaderTextureFactor.Mul(v.SphereShaderTextureFactor)
// 		sm.ToonShaderTextureFactor.Mul(v.ToonShaderTextureFactor)
// 	}
// }

// func (sm *ShaderMaterial) IAdd(v interface{}) {
// 	switch v := v.(type) {
// 	case float64:
// 		sm.Material.Diffuse.AddScalar(v)
// 		sm.Material.Ambient.AddScalar(v)
// 		sm.Material.Specular.AddScalar(v)
// 		sm.Material.EdgeSize += v
// 		sm.ShaderTextureFactor.AddScalar(v)
// 		sm.SphereShaderTextureFactor.AddScalar(v)
// 		sm.ToonShaderTextureFactor.AddScalar(v)
// 	case int:
// 		sm.IAdd(float32(v))
// 	case *ShaderMaterial:
// 		sm.Material.Diffuse.Add(&v.Material.Diffuse)
// 		sm.Material.Ambient.Add(&v.Material.Ambient)
// 		sm.Material.Specular.Add(&v.Material.Specular)
// 		sm.Material.EdgeSize += v.Material.EdgeSize
// 		sm.ShaderTextureFactor.Add(v.ShaderTextureFactor)
// 		sm.SphereShaderTextureFactor.Add(v.SphereShaderTextureFactor)
// 		sm.ToonShaderTextureFactor.Add(v.ToonShaderTextureFactor)
// 	}
// }
