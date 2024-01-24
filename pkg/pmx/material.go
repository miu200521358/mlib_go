package pmx

import (
	"github.com/miu200521358/mlib_go/pkg/mcore"
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
	DiffuseColor  [3]float32
	DiffuseAlpha  float32
	SpecularColor [3]float32
	SpecularPower float32
	AmbientColor  [3]float32
	EdgeColor     [3]float32
	EdgeAlpha     float32
	EdgeSize      float32
	Texture       *TextureGL
	SphereTexture *TextureGL
	ToonTexture   *TextureGL
}

type Material struct {
	*mcore.IndexModel
	// 材質名
	Name string
	// 材質名英
	EnglishName string
	// Diffuse (R,G,B,A)(拡散色＋非透過度)
	DiffuseColor mmath.MVec3
	DiffuseAlpha float64
	// Specular (R,G,B)(反射色)
	SpecularColor mmath.MVec3
	// Specular係数(反射強度)
	SpecularPower float64
	// Ambient (R,G,B)(環境色)
	AmbientColor mmath.MVec3
	// 描画フラグ(8bit) - 各bit 0:OFF 1:ON
	DrawFlag DrawFlag
	// エッジ色 (R,G,B,A)
	EdgeColor mmath.MVec3
	EdgeAlpha float64
	// エッジサイズ
	EdgeSize float64
	// 通常テクスチャINDEX
	TextureIndex int
	// スフィアテクスチャINDEX
	SphereTextureIndex int
	// スフィアモード
	SphereMode SphereMode
	// 共有Toonフラグ
	ToonSharingFlag ToonSharing
	// ToonテクスチャINDEX
	ToonTextureIndex int
	// メモ
	Memo string
	// 材質に対応する面(頂点)数 (必ず3の倍数になる)
	VerticesCount int
}

func NewMaterial() *Material {
	return &Material{
		IndexModel:         &mcore.IndexModel{Index: -1},
		Name:               "",
		EnglishName:        "",
		DiffuseColor:       mmath.MVec3{},
		DiffuseAlpha:       0.0,
		SpecularColor:      mmath.MVec3{},
		SpecularPower:      0.0,
		AmbientColor:       mmath.MVec3{},
		DrawFlag:           DRAW_FLAG_NONE,
		EdgeColor:          mmath.MVec3{},
		EdgeAlpha:          0.0,
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
	sphereTexture *Texture,
	toonTexture *Texture,
	windowIndex int,
) *MaterialGL {
	var textureGL *TextureGL
	if texture != nil {
		textureGL = texture.GL(modelPath, TEXTURE_TYPE_TEXTURE, true, windowIndex)
	}

	var sphereTextureGL *TextureGL
	if sphereTexture != nil {
		sphereTextureGL = sphereTexture.GL(modelPath, TEXTURE_TYPE_SPHERE, true, windowIndex)
	}

	var tooTextureGL *TextureGL
	if toonTexture != nil {
		tooTextureGL = toonTexture.GL(
			modelPath, TEXTURE_TYPE_TOON, m.ToonSharingFlag == TOON_SHARING_INDIVIDUAL, windowIndex)
	}

	return &MaterialGL{
		DiffuseColor:  m.DiffuseColor.GL(),
		DiffuseAlpha:  float32(m.DiffuseAlpha),
		SpecularColor: m.SpecularColor.GL(),
		SpecularPower: float32(m.SpecularPower),
		AmbientColor:  m.AmbientColor.GL(),
		EdgeColor:     m.EdgeColor.GL(),
		EdgeAlpha:     float32(m.EdgeAlpha),
		EdgeSize:      float32(m.EdgeSize),
		Texture:       textureGL,
		SphereTexture: sphereTextureGL,
		ToonTexture:   tooTextureGL,
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

// シェーダー用材質
type ShaderMaterial struct {
	LightAmbient4             *mmath.MVec3
	Material                  *Material
	ShaderTextureFactor       *mmath.MVec4
	SphereShaderTextureFactor *mmath.MVec4
	ToonShaderTextureFactor   *mmath.MVec4
}

func NewShaderMaterial(
	material *Material,
	lightAmbient4 *mmath.MVec3,
	textureFactor *mmath.MVec4,
	toonTextureFactor *mmath.MVec4,
	sphereTextureFactor *mmath.MVec4,
) *ShaderMaterial {
	return &ShaderMaterial{
		LightAmbient4:             lightAmbient4,
		Material:                  material.Copy().(*Material),
		ShaderTextureFactor:       textureFactor,
		SphereShaderTextureFactor: toonTextureFactor,
		ToonShaderTextureFactor:   sphereTextureFactor,
	}
}

func (sm *ShaderMaterial) Diffuse() []float32 {
	diffuse := make([]float32, 3)
	diffuse[0] = float32(sm.Material.DiffuseColor.GetX())*float32(sm.LightAmbient4.GetX()) + float32(sm.Material.AmbientColor.GetX())
	diffuse[1] = float32(sm.Material.DiffuseColor.GetY())*float32(sm.LightAmbient4.GetY()) + float32(sm.Material.AmbientColor.GetY())
	diffuse[2] = float32(sm.Material.DiffuseColor.GetZ())*float32(sm.LightAmbient4.GetZ()) + float32(sm.Material.AmbientColor.GetZ())
	return diffuse
}

func (sm *ShaderMaterial) DiffuseAlpha() float32 {
	return float32(sm.Material.DiffuseAlpha)
}

func (sm *ShaderMaterial) Ambient() []float32 {
	ambient := make([]float32, 3)
	ambient[0] = float32(sm.Material.DiffuseColor.GetX()) * float32(sm.LightAmbient4.GetX())
	ambient[1] = float32(sm.Material.DiffuseColor.GetY()) * float32(sm.LightAmbient4.GetY())
	ambient[2] = float32(sm.Material.DiffuseColor.GetZ()) * float32(sm.LightAmbient4.GetZ())
	return ambient
}

func (sm *ShaderMaterial) Specular() []float32 {
	specular := make([]float32, 3)
	specular[0] = float32(sm.Material.SpecularColor.GetX()) * float32(sm.LightAmbient4.GetX())
	specular[1] = float32(sm.Material.SpecularColor.GetY()) * float32(sm.LightAmbient4.GetY())
	specular[2] = float32(sm.Material.SpecularColor.GetZ()) * float32(sm.LightAmbient4.GetZ())
	return specular
}

func (sm *ShaderMaterial) SpecularPower() float32 {
	return float32(sm.Material.SpecularPower)
}

func (sm *ShaderMaterial) EdgeColor() []float32 {
	edgeColor := make([]float32, 3)
	edgeColor[0] = float32(sm.Material.EdgeColor.GetX()) * float32(sm.Material.DiffuseAlpha)
	edgeColor[1] = float32(sm.Material.EdgeColor.GetY()) * float32(sm.Material.DiffuseAlpha)
	edgeColor[2] = float32(sm.Material.EdgeColor.GetZ()) * float32(sm.Material.DiffuseAlpha)
	return edgeColor
}

func (sm *ShaderMaterial) EdgeAlpha() float32 {
	return float32(sm.Material.EdgeAlpha)
}

func (sm *ShaderMaterial) EdgeSize() float32 {
	return float32(sm.Material.EdgeSize)
}

func (sm *ShaderMaterial) TextureFactor() []float32 {
	textureFactor := make([]float32, 3)
	textureFactor[0] = float32(sm.ShaderTextureFactor.GetX())
	textureFactor[1] = float32(sm.ShaderTextureFactor.GetY())
	textureFactor[2] = float32(sm.ShaderTextureFactor.GetZ())
	textureFactor[3] = float32(sm.ShaderTextureFactor.GetW())
	return textureFactor
}

func (sm *ShaderMaterial) SphereTextureFactor() []float32 {
	sphereTextureFactor := make([]float32, 3)
	sphereTextureFactor[0] = float32(sm.SphereShaderTextureFactor.GetX())
	sphereTextureFactor[1] = float32(sm.SphereShaderTextureFactor.GetY())
	sphereTextureFactor[2] = float32(sm.SphereShaderTextureFactor.GetZ())
	sphereTextureFactor[3] = float32(sm.SphereShaderTextureFactor.GetW())
	return sphereTextureFactor
}

func (sm *ShaderMaterial) ToonTextureFactor() []float32 {
	toonTextureFactor := make([]float32, 3)
	toonTextureFactor[0] = float32(sm.ToonShaderTextureFactor.GetX())
	toonTextureFactor[1] = float32(sm.ToonShaderTextureFactor.GetY())
	toonTextureFactor[2] = float32(sm.ToonShaderTextureFactor.GetZ())
	toonTextureFactor[3] = float32(sm.ToonShaderTextureFactor.GetW())
	return toonTextureFactor
}

func (sm *ShaderMaterial) IMul(v interface{}) {
	switch v := v.(type) {
	case float64:
		sm.Material.DiffuseColor.MulScalar(v)
		sm.Material.DiffuseAlpha *= v
		sm.Material.AmbientColor.MulScalar(v)
		sm.Material.SpecularColor.MulScalar(v)
		sm.Material.EdgeColor.MulScalar(v)
		sm.Material.EdgeSize *= v
		sm.Material.EdgeAlpha *= v
		sm.ShaderTextureFactor.MulScalar(v)
		sm.SphereShaderTextureFactor.MulScalar(v)
		sm.ToonShaderTextureFactor.MulScalar(v)
	case int:
		sm.IMul(float32(v))
	case *ShaderMaterial:
		sm.Material.DiffuseColor.Mul(&v.Material.DiffuseColor)
		sm.Material.DiffuseAlpha *= v.Material.DiffuseAlpha
		sm.Material.AmbientColor.Mul(&v.Material.AmbientColor)
		sm.Material.SpecularColor.Mul(&v.Material.SpecularColor)
		sm.Material.EdgeColor.Mul(&v.Material.EdgeColor)
		sm.Material.EdgeSize *= v.Material.EdgeSize
		sm.Material.EdgeAlpha *= v.Material.EdgeAlpha
		sm.ShaderTextureFactor.Mul(v.ShaderTextureFactor)
		sm.SphereShaderTextureFactor.Mul(v.SphereShaderTextureFactor)
		sm.ToonShaderTextureFactor.Mul(v.ToonShaderTextureFactor)
	}
}

func (sm *ShaderMaterial) IAdd(v interface{}) {
	switch v := v.(type) {
	case float64:
		sm.Material.DiffuseColor.AddScalar(v)
		sm.Material.DiffuseAlpha += v
		sm.Material.AmbientColor.AddScalar(v)
		sm.Material.SpecularColor.AddScalar(v)
		sm.Material.EdgeColor.AddScalar(v)
		sm.Material.EdgeSize += v
		sm.Material.EdgeAlpha += v
		sm.ShaderTextureFactor.AddScalar(v)
		sm.SphereShaderTextureFactor.AddScalar(v)
		sm.ToonShaderTextureFactor.AddScalar(v)
	case int:
		sm.IAdd(float32(v))
	case *ShaderMaterial:
		sm.Material.DiffuseColor.Add(&v.Material.DiffuseColor)
		sm.Material.DiffuseAlpha += v.Material.DiffuseAlpha
		sm.Material.AmbientColor.Add(&v.Material.AmbientColor)
		sm.Material.SpecularColor.Add(&v.Material.SpecularColor)
		sm.Material.EdgeColor.Add(&v.Material.EdgeColor)
		sm.Material.EdgeSize += v.Material.EdgeSize
		sm.Material.EdgeAlpha += v.Material.EdgeAlpha
		sm.ShaderTextureFactor.Add(v.ShaderTextureFactor)
		sm.SphereShaderTextureFactor.Add(v.SphereShaderTextureFactor)
		sm.ToonShaderTextureFactor.Add(v.ToonShaderTextureFactor)
	}
}
