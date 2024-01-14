package material

import (
	"github.com/miu200521358/mlib_go/pkg/core/index_model"
	"github.com/miu200521358/mlib_go/pkg/math/mvec3"
)

type SphereMode int

const (
	INVALID        SphereMode = 0
	MULTIPLICATION SphereMode = 1
	ADDITION       SphereMode = 2
	SUBTEXTURE     SphereMode = 3
)

type DrawFlag int

const (
	NONE                        DrawFlag = 0x0000
	DOUBLE_SIDED_DRAWING        DrawFlag = 0x0001
	GROUND_SHADOW               DrawFlag = 0x0002
	DRAWING_ON_SELF_SHADOW_MAPS DrawFlag = 0x0004
	DRAWING_SELF_SHADOWS        DrawFlag = 0x0008
	DRAWING_EDGE                DrawFlag = 0x0010
)

type ToonSharing int

const (
	INDIVIDUAL ToonSharing = 0
	SHARING    ToonSharing = 1
)

type Material struct {
	index_model.IndexModel
	Index              int
	Name               string
	EnglishName        string
	DiffuseColor       *mvec3.T
	DiffuseAlpha       float64
	SpecularColor      *mvec3.T
	SpecularFactor     float64
	AmbientColor       *mvec3.T
	DrawFlag           DrawFlag
	EdgeColor          *mvec3.T
	EdgeAlpha          float64
	EdgeSize           float64
	TextureIndex       int
	SphereTextureIndex int
	SphereMode         SphereMode
	ToonSharingFlag    ToonSharing
	ToonTextureIndex   int
	Comment            string
	VerticesCount      int
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
		Index:              index,
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

// 材質リスト
type Materials struct {
	index_model.IndexModelCorrection[*Material]
}

func NewMaterials(name string) *Materials {
	return &Materials{
		IndexModelCorrection: *index_model.NewIndexModelCorrection[*Material](),
	}
}
