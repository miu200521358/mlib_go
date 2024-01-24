package pmx

import (
	"image"
	"image/draw"
	"path/filepath"

	"github.com/go-gl/gl/v4.4-core/gl"

	"github.com/miu200521358/mlib_go/pkg/mcore"
	"github.com/miu200521358/mlib_go/pkg/mutils"

)

// テクスチャ種別
type TextureType int

const (
	TEXTURE_TYPE_TEXTURE TextureType = 0 // テクスチャ
	TEXTURE_TYPE_TOON    TextureType = 1 // Toonテクスチャ
	TEXTURE_TYPE_SPHERE  TextureType = 2 // スフィアテクスチャ
)

type TextureGL struct {
	id          uint32
	typeId      uint32
	Valid       bool
	textureType TextureType
}

type Texture struct {
	*mcore.IndexModel
	Name          string      // テクスチャ名
	TextureType   TextureType // テクスチャ種別
	Path          string      // テクスチャフルパス
	Valid         bool        // テクスチャフルパスが有効であるか否か
	glId          uint32      // OpenGLテクスチャID
	Initialized   bool        // 描画初期化済みフラグ
	Image         image.Image // テクスチャイメージ
	textureTypeId uint32      // テクスチャタイプID
}

func NewTexture() *Texture {
	return &Texture{
		IndexModel:  &mcore.IndexModel{Index: -1},
		Name:        "",
		TextureType: TEXTURE_TYPE_TEXTURE,
		Path:        "",
		Valid:       false,
		glId:        0,
	}
}

func (t *Texture) GL(
	modelPath string,
	textureType TextureType,
	isIndividual bool,
	windowIndex int,
) *TextureGL {
	tGl := &TextureGL{}

	if t.Initialized && t.Valid {
		// 既にフラグが立ってたら描画初期化済み
		tGl.id = t.glId
		tGl.typeId = t.textureTypeId
		tGl.Valid = t.Valid
		tGl.textureType = t.TextureType
		return tGl
	}

	// global texture
	var texPath string
	if isIndividual {
		texPath = filepath.Join(filepath.Dir(modelPath), t.Name)
	} else {
		texPath = t.Name
	}

	// テクスチャがちゃんとある場合のみ初期化処理実施
	valid, err := mutils.ExistsFile(texPath)
	t.Valid = (valid && err == nil)
	if !t.Valid {
		t.Initialized = true
		return nil
	}

	t.Path = texPath
	img, err := mutils.LoadImage(texPath)
	if err != nil {
		t.Valid = false
	} else {
		rgba := image.NewRGBA(img.Bounds())
		if rgba.Stride != rgba.Rect.Size().X*4 {
			return nil
		}
		draw.Draw(rgba, rgba.Bounds(), img, image.Point{0, 0}, draw.Src)

		t.Image = img
		t.Image = mutils.FlipImage(t.Image)
	}

	if !t.Valid {
		t.Initialized = true
		return nil
	}

	t.TextureType = textureType

	// テクスチャオブジェクト生成
	gl.GenTextures(1, &t.glId)

	// テクスチャ種別によってテクスチャユニットを変更
	if windowIndex == 0 {
		t.textureTypeId = gl.TEXTURE0
		if textureType == TEXTURE_TYPE_TOON {
			t.textureTypeId = gl.TEXTURE1
		} else if textureType == TEXTURE_TYPE_SPHERE {
			t.textureTypeId = gl.TEXTURE2
		}
	} else if windowIndex == 1 {
		t.textureTypeId = gl.TEXTURE3
		if textureType == TEXTURE_TYPE_TOON {
			t.textureTypeId = gl.TEXTURE4
		} else if textureType == TEXTURE_TYPE_SPHERE {
			t.textureTypeId = gl.TEXTURE5
		}
	} else if windowIndex == 2 {
		t.textureTypeId = gl.TEXTURE6
		if textureType == TEXTURE_TYPE_TOON {
			t.textureTypeId = gl.TEXTURE7
		} else if textureType == TEXTURE_TYPE_SPHERE {
			t.textureTypeId = gl.TEXTURE8
		}
	}

	tGl.id = t.glId
	tGl.typeId = t.textureTypeId
	tGl.Valid = t.Valid
	tGl.textureType = t.TextureType

	tGl.Bind()
	gl.TexImage2D(
		gl.TEXTURE_2D,
		0,
		gl.RGBA,
		int32(t.Image.Bounds().Dx()),
		int32(t.Image.Bounds().Dy()),
		0,
		gl.RGBA,
		gl.UNSIGNED_BYTE,
		gl.Ptr(t.Image.(*image.RGBA).Pix),
	)
	tGl.Unbind()

	// 描画初期化
	t.Initialized = true
	return tGl
}

func (t *TextureGL) Bind() {
	if !t.Valid {
		return
	}

	gl.ActiveTexture(t.typeId)
	gl.BindTexture(gl.TEXTURE_2D, t.id)

	if t.textureType == TEXTURE_TYPE_TOON {
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	}

	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAX_LEVEL, 0)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
}

func (t *TextureGL) Unbind() {
	if !t.Valid {
		return
	}

	gl.BindTexture(gl.TEXTURE_2D, 0)
}

// テクスチャリスト
type Textures struct {
	*mcore.IndexModelCorrection[*Texture]
}

func NewTextures() *Textures {
	return &Textures{
		IndexModelCorrection: mcore.NewIndexModelCorrection[*Texture](),
	}
}

// 共有テクスチャ辞書
type ToonTextures struct {
	*mcore.IndexModelCorrection[*Texture]
}

func NewToonTextures() *ToonTextures {
	return &ToonTextures{
		IndexModelCorrection: mcore.NewIndexModelCorrection[*Texture](),
	}
}
