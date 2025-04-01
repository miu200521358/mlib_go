package pmx

import (
	"slices"

	"github.com/miu200521358/mlib_go/pkg/domain/core"
)

// 材質リスト
type Materials struct {
	*core.IndexNameModels[*Material]
	Vertices map[int][]int
	Faces    map[int][]int
}

func NewMaterials(capacity int) *Materials {
	return &Materials{
		IndexNameModels: core.NewIndexNameModels[*Material](capacity),
		Vertices:        make(map[int][]int),
		Faces:           make(map[int][]int),
	}
}

func (materials *Materials) Setup(vertices *Vertices, faces *Faces, textures *Textures) {
	prevVertexCount := 0

	vertices.ForEach(func(index int, vertex *Vertex) bool {
		vertex.MaterialIndexes = make([]int, 0)
		return true
	})

	materials.ForEach(func(index int, material *Material) bool {
		for j := prevVertexCount; j < prevVertexCount+int(material.VerticesCount/3); j++ {
			face, err := faces.Get(j)
			if err != nil {
				continue
			}
			for _, vertexIndex := range face.VertexIndexes {
				vertex, err := vertices.Get(vertexIndex)
				if err != nil {
					continue
				}
				if !slices.Contains(vertex.MaterialIndexes, material.Index()) {
					vertex.MaterialIndexes = append(vertex.MaterialIndexes, material.Index())
				}
			}
		}

		prevVertexCount += int(material.VerticesCount / 3)

		if material.TextureIndex != -1 && textures.Contains(material.TextureIndex) {
			texture, err := textures.Get(material.TextureIndex)
			if err == nil {
				texture.TextureType = TEXTURE_TYPE_TEXTURE
			}
		}
		if material.ToonTextureIndex != -1 && material.ToonSharingFlag == TOON_SHARING_INDIVIDUAL &&
			textures.Contains(material.ToonTextureIndex) {
			texture, err := textures.Get(material.ToonTextureIndex)
			if err == nil {
				texture.TextureType = TEXTURE_TYPE_TOON
			}
		}
		if material.SphereTextureIndex != -1 && textures.Contains(material.SphereTextureIndex) {
			texture, err := textures.Get(material.SphereTextureIndex)
			if err == nil {
				texture.TextureType = TEXTURE_TYPE_SPHERE
			}
		}

		return true
	})
}
