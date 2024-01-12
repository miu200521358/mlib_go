package material

import (
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
