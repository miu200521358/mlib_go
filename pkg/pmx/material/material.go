package material

import (
	"github.com/miu200521358/mlib_go/pkg/core/index_model"
	"github.com/miu200521358/mlib_go/pkg/math/mvec3"
)

// スフィアモード
type SphereMode int

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

type DrawFlag int

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
type ToonSharing int

const (
	// 0:継続値は個別Toon
	TOON_SHARING_INDIVIDUAL ToonSharing = 0
	// 1:継続値は共有Toon
	TOON_SHARING_SHARING ToonSharing = 1
)

type Material struct {
	index_model.IndexModel
	// 材質名
	Name string
	// 材質名英
	EnglishName string
	// Diffuse (R,G,B,A)(拡散色＋非透過度)
	DiffuseColor *mvec3.T
	DiffuseAlpha float64
	// Specular (R,G,B)(反射色)
	SpecularColor *mvec3.T
	// Specular係数(反射強度)
	SpecularFactor float64
	// Ambient (R,G,B)(環境色)
	AmbientColor *mvec3.T
	// 描画フラグ(8bit) - 各bit 0:OFF 1:ON
	DrawFlag DrawFlag
	// エッジ色 (R,G,B,A)
	EdgeColor *mvec3.T
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
	Comment string
	// 材質に対応する面(頂点)数 (必ず3の倍数になる)
	VerticesCount int
}

func NewMaterial(
	index int,
	name string,
	englishName string,
	diffuseColor *mvec3.T,
	diffuseAlpha float64,
	specularColor *mvec3.T,
	specularFactor float64,
	ambientColor *mvec3.T,
	drawFlag DrawFlag,
	edgeColor *mvec3.T,
	edgeAlpha float64,
	edgeSize float64,
	textureIndex int,
	sphereTextureIndex int,
	sphereMode SphereMode,
	toonSharingFlag ToonSharing,
	toonTextureIndex int,
	comment string,
	verticesCount int,
) *Material {
	return &Material{
		IndexModel:         index_model.IndexModel{Index: index},
		Name:               name,
		EnglishName:        englishName,
		DiffuseColor:       diffuseColor,
		DiffuseAlpha:       diffuseAlpha,
		SpecularColor:      specularColor,
		SpecularFactor:     specularFactor,
		AmbientColor:       ambientColor,
		DrawFlag:           drawFlag,
		EdgeColor:          edgeColor,
		EdgeAlpha:          edgeAlpha,
		EdgeSize:           edgeSize,
		TextureIndex:       textureIndex,
		SphereTextureIndex: sphereTextureIndex,
		SphereMode:         sphereMode,
		ToonSharingFlag:    toonSharingFlag,
		ToonTextureIndex:   toonTextureIndex,
		Comment:            comment,
		VerticesCount:      verticesCount,
	}
}

func (m *Material) Copy() index_model.IndexModelInterface {
	copied := *m
	return &copied
}

// 材質リスト
type Materials struct {
	index_model.IndexModelCorrection[*Material]
}

func NewMaterials(name string) *Materials {
	return &Materials{
		IndexModelCorrection: *index_model.NewIndexModelCorrection[*Material](),
	}
}

// シェーダー用材質
type ShaderMaterial struct {
	LightAmbient4             *mvec3.T
	Material                  *Material
	ShaderTextureFactor       *mvec3.T
	SphereShaderTextureFactor *mvec3.T
	ToonShaderTextureFactor   *mvec3.T
}

func NewShaderMaterial(
	material *Material,
	lightAmbient4 *mvec3.T,
	textureFactor *mvec3.T,
	toonTextureFactor *mvec3.T,
	sphereTextureFactor *mvec3.T,
) *ShaderMaterial {
	return &ShaderMaterial{
		LightAmbient4:             lightAmbient4,
		Material:                  material.Copy().(*Material),
		ShaderTextureFactor:       textureFactor,
		SphereShaderTextureFactor: toonTextureFactor,
		ToonShaderTextureFactor:   sphereTextureFactor,
	}
}

func (sm *ShaderMaterial) Diffuse() []float64 {
	diffuse := make([]float64, 3)
	diffuse[0] = sm.Material.DiffuseColor.GetX()*sm.LightAmbient4.GetX() + sm.Material.AmbientColor.GetX()
	diffuse[1] = sm.Material.DiffuseColor.GetY()*sm.LightAmbient4.GetY() + sm.Material.AmbientColor.GetY()
	diffuse[2] = sm.Material.DiffuseColor.GetZ()*sm.LightAmbient4.GetZ() + sm.Material.AmbientColor.GetZ()
	return diffuse
}

func (sm *ShaderMaterial) DiffuseAlpha() float64 {
	return sm.Material.DiffuseAlpha
}

func (sm *ShaderMaterial) Ambient() []float64 {
	ambient := make([]float64, 3)
	ambient[0] = sm.Material.DiffuseColor.GetX() * sm.LightAmbient4.GetX()
	ambient[1] = sm.Material.DiffuseColor.GetY() * sm.LightAmbient4.GetY()
	ambient[2] = sm.Material.DiffuseColor.GetZ() * sm.LightAmbient4.GetZ()
	return ambient
}

func (sm *ShaderMaterial) Specular() []float64 {
	specular := make([]float64, 3)
	specular[0] = sm.Material.SpecularColor.GetX() * sm.LightAmbient4.GetX()
	specular[1] = sm.Material.SpecularColor.GetY() * sm.LightAmbient4.GetY()
	specular[2] = sm.Material.SpecularColor.GetZ() * sm.LightAmbient4.GetZ()
	return specular
}

func (sm *ShaderMaterial) SpecularFactor() float64 {
	return sm.Material.SpecularFactor
}

func (sm *ShaderMaterial) EdgeColor() []float64 {
	edgeColor := make([]float64, 3)
	edgeColor[0] = sm.Material.EdgeColor.GetX() * sm.Material.DiffuseAlpha
	edgeColor[1] = sm.Material.EdgeColor.GetY() * sm.Material.DiffuseAlpha
	edgeColor[2] = sm.Material.EdgeColor.GetZ() * sm.Material.DiffuseAlpha
	return edgeColor
}

func (sm *ShaderMaterial) EdgeAlpha() float64 {
	return sm.Material.EdgeAlpha
}

func (sm *ShaderMaterial) EdgeSize() float64 {
	return sm.Material.EdgeSize
}

func (sm *ShaderMaterial) TextureFactor() []float64 {
	return *sm.ShaderTextureFactor.Vector()
}

func (sm *ShaderMaterial) SphereTextureFactor() []float64 {
	return *sm.SphereShaderTextureFactor.Vector()
}

func (sm *ShaderMaterial) ToonTextureFactor() []float64 {
	return *sm.ToonShaderTextureFactor.Vector()
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
		sm.IMul(float64(v))
	case *ShaderMaterial:
		sm.Material.DiffuseColor.Mul(v.Material.DiffuseColor)
		sm.Material.DiffuseAlpha *= v.Material.DiffuseAlpha
		sm.Material.AmbientColor.Mul(v.Material.AmbientColor)
		sm.Material.SpecularColor.Mul(v.Material.SpecularColor)
		sm.Material.EdgeColor.Mul(v.Material.EdgeColor)
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
		sm.IAdd(float64(v))
	case *ShaderMaterial:
		sm.Material.DiffuseColor.Add(v.Material.DiffuseColor)
		sm.Material.DiffuseAlpha += v.Material.DiffuseAlpha
		sm.Material.AmbientColor.Add(v.Material.AmbientColor)
		sm.Material.SpecularColor.Add(v.Material.SpecularColor)
		sm.Material.EdgeColor.Add(v.Material.EdgeColor)
		sm.Material.EdgeSize += v.Material.EdgeSize
		sm.Material.EdgeAlpha += v.Material.EdgeAlpha
		sm.ShaderTextureFactor.Add(v.ShaderTextureFactor)
		sm.SphereShaderTextureFactor.Add(v.SphereShaderTextureFactor)
		sm.ToonShaderTextureFactor.Add(v.ToonShaderTextureFactor)
	}
}
