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

	for v := range vertices.Iterator() {
		v.MaterialIndexes = make([]int, 0)
	}

	for m := range materials.Iterator() {
		for j := prevVertexCount; j < prevVertexCount+int(m.VerticesCount/3); j++ {
			face, err := faces.Get(j)
			if err != nil {
				continue
			}
			for _, vertexIndex := range face.VertexIndexes {
				vertex, err := vertices.Get(vertexIndex)
				if err != nil {
					continue
				}
				if !slices.Contains(vertex.MaterialIndexes, m.Index()) {
					vertex.MaterialIndexes = append(vertex.MaterialIndexes, m.Index())
				}
			}
		}

		prevVertexCount += int(m.VerticesCount / 3)

		if m.TextureIndex != -1 && textures.Contains(m.TextureIndex) {
			texture, err := textures.Get(m.TextureIndex)
			if err == nil {
				texture.TextureType = TEXTURE_TYPE_TEXTURE
			}
		}
		if m.ToonTextureIndex != -1 && m.ToonSharingFlag == TOON_SHARING_INDIVIDUAL &&
			textures.Contains(m.ToonTextureIndex) {
			texture, err := textures.Get(m.ToonTextureIndex)
			if err == nil {
				texture.TextureType = TEXTURE_TYPE_TOON
			}
		}
		if m.SphereTextureIndex != -1 && textures.Contains(m.SphereTextureIndex) {
			texture, err := textures.Get(m.SphereTextureIndex)
			if err == nil {
				texture.TextureType = TEXTURE_TYPE_SPHERE
			}
		}
	}
}
