package texture

type TextureType int

const (
	TEXTURE TextureType = 0
	TOON    TextureType = 1
	SPHERE  TextureType = 2
)

type Texture struct {
	Index       int
	Name        string
	TextureType TextureType
	Path        string
	Valid       bool
}

func NewTexture(index int, name string, textureType TextureType) *Texture {
	return &Texture{
		Index:       index,
		Name:        name,
		TextureType: textureType,
		Path:        "",
		Valid:       false,
	}
}
