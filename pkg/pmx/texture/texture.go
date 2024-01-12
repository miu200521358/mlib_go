package texture

type TextureType int

const (
	TEXTURE TextureType = 0
	TOON    TextureType = 1
	SPHERE  TextureType = 2
)

type T struct {
	Index       int
	Name        string
	TextureType TextureType
	Path        string
	Valid       bool
}

func NewTexture(index int, name string, textureType TextureType) *T {
	return &T{
		Index:       index,
		Name:        name,
		TextureType: textureType,
		Path:        "",
		Valid:       false,
	}
}

// Copy
func (v *T) Copy() *T {
	copied := *v
	return &copied
}
