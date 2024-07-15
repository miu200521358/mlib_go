//go:build windows
// +build windows

package renderer

import (
	"embed"
	"fmt"
	"path/filepath"

	"github.com/go-gl/gl/v4.4-core/gl"

	"github.com/miu200521358/mlib_go/pkg/domain/mmath"
	"github.com/miu200521358/mlib_go/pkg/domain/pmx"
	"github.com/miu200521358/mlib_go/pkg/mutils"
	"github.com/miu200521358/mlib_go/pkg/mutils/mi18n"
	"github.com/miu200521358/mlib_go/pkg/mutils/mlog"
)

//go:embed toon/*
var toonFiles embed.FS

type textureGl struct {
	*pmx.Texture
	Id                uint32          // OpenGLテクスチャID
	TextureType       pmx.TextureType // テクスチャ種別
	TextureUnitId     uint32          // テクスチャ種類別描画先ユニットID
	TextureUnitNo     uint32          // テクスチャ種類別描画先ユニット番号
	IsGeneratedMipmap bool            // ミップマップが生成されているか否か
	Initialized       bool            // 描画初期化済みフラグ
}

func (t *textureGl) Bind() {
	gl.ActiveTexture(t.TextureUnitId)
	gl.BindTexture(gl.TEXTURE_2D, t.Id)

	if t.TextureType == pmx.TEXTURE_TYPE_TOON {
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.REPEAT)
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.REPEAT)
	}

	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)

	if t.IsGeneratedMipmap {
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR_MIPMAP_LINEAR)
	} else {
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAX_LEVEL, 0)
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	}
}

func (t *textureGl) Unbind() {
	gl.BindTexture(gl.TEXTURE_2D, 0)
}

func (renderModel *RenderModel) initTexturesGl(windowIndex int) error {
	renderModel.textures = make([]*textureGl, len(renderModel.Model.Textures.Data))

	for _, texture := range renderModel.Model.Textures.Data {
		texGl, err := renderModel.initTextureGl(windowIndex, texture)
		if err != nil {
			mlog.D(fmt.Sprintf("texture error: %s", err))
			continue
		}
		renderModel.textures[texture.Index] = texGl
	}

	return nil
}

func (renderModel *RenderModel) initTextureGl(
	windowIndex int, texture *pmx.Texture,
) (*textureGl, error) {
	texGl := &textureGl{
		Texture: texture,
	}

	// 通常テクスチャ
	texPath := filepath.Join(filepath.Dir(renderModel.Model.GetPath()), texture.Name)

	// テクスチャがちゃんとある場合のみ初期化処理実施
	valid, err := mutils.ExistsFile(texPath)
	if !valid || err != nil {
		texGl.Initialized = false
		return texGl, fmt.Errorf("not found texture file: %s", texPath)
	}

	img, err := mutils.LoadImage(texPath)
	if err != nil {
		return nil, err
	}
	image := mutils.ConvertToNRGBA(*img)

	texGl.IsGeneratedMipmap =
		mmath.IsPowerOfTwo(image.Rect.Size().X) && mmath.IsPowerOfTwo(image.Rect.Size().Y)

	// テクスチャオブジェクト生成
	gl.GenTextures(1, &texGl.Id)

	// テクスチャ種別によってテクスチャユニットを変更
	if windowIndex == 0 {
		switch texture.TextureType {
		case pmx.TEXTURE_TYPE_TEXTURE:
			texGl.TextureUnitId = gl.TEXTURE0
			texGl.TextureUnitNo = 0
		case pmx.TEXTURE_TYPE_TOON:
			texGl.TextureUnitId = gl.TEXTURE1
			texGl.TextureUnitNo = 1
		case pmx.TEXTURE_TYPE_SPHERE:
			texGl.TextureUnitId = gl.TEXTURE2
			texGl.TextureUnitNo = 2
		}
	} else if windowIndex == 1 {
		switch texture.TextureType {
		case pmx.TEXTURE_TYPE_TEXTURE:
			texGl.TextureUnitId = gl.TEXTURE3
			texGl.TextureUnitNo = 3
		case pmx.TEXTURE_TYPE_TOON:
			texGl.TextureUnitId = gl.TEXTURE4
			texGl.TextureUnitNo = 4
		case pmx.TEXTURE_TYPE_SPHERE:
			texGl.TextureUnitId = gl.TEXTURE5
			texGl.TextureUnitNo = 5
		}
	} else if windowIndex == 2 {
		switch texture.TextureType {
		case pmx.TEXTURE_TYPE_TEXTURE:
			texGl.TextureUnitId = gl.TEXTURE6
			texGl.TextureUnitNo = 6
		case pmx.TEXTURE_TYPE_TOON:
			texGl.TextureUnitId = gl.TEXTURE7
			texGl.TextureUnitNo = 7
		case pmx.TEXTURE_TYPE_SPHERE:
			texGl.TextureUnitId = gl.TEXTURE8
			texGl.TextureUnitNo = 8
		}
	}

	texGl.Bind()

	gl.TexImage2D(
		gl.TEXTURE_2D,
		0,
		gl.RGBA,
		int32(image.Rect.Size().X),
		int32(image.Rect.Size().Y),
		0,
		gl.RGBA,
		gl.UNSIGNED_BYTE,
		gl.Ptr(image.Pix),
	)

	if texGl.IsGeneratedMipmap {
		gl.GenerateMipmap(gl.TEXTURE_2D)
	} else {
		mlog.D(mi18n.T("ミップマップ生成エラー", map[string]interface{}{"Name": texture.Name}))
	}

	texGl.Unbind()

	// 描画初期化
	texGl.Initialized = true
	return texGl, nil
}

func (renderModel *RenderModel) initToonTexturesGl(windowIndex int) error {
	renderModel.toonTextures = make([]*textureGl, 10)

	for i := 0; i < 10; i++ {
		filePath := fmt.Sprintf("toon/toon%02d.bmp", i+1)

		toon := pmx.NewTexture()
		toon.Index = i
		toon.Name = filePath
		toon.TextureType = pmx.TEXTURE_TYPE_TOON

		img, err := mutils.LoadImageFromResources(toonFiles, filePath)
		if err != nil {
			return err
		}
		image := mutils.ConvertToNRGBA(*img)

		toonGl := &textureGl{
			Texture: toon,
		}

		// テクスチャオブジェクト生成
		gl.GenTextures(1, &toonGl.Id)

		// Toon用テクスチャユニットを設定
		if windowIndex == 0 {
			toonGl.TextureUnitId = gl.TEXTURE10
			toonGl.TextureUnitNo = 10
		} else if windowIndex == 1 {
			toonGl.TextureUnitId = gl.TEXTURE11
			toonGl.TextureUnitNo = 11
		} else if windowIndex == 2 {
			toonGl.TextureUnitId = gl.TEXTURE12
			toonGl.TextureUnitNo = 12
		}

		toonGl.Bind()

		gl.TexImage2D(
			gl.TEXTURE_2D,
			0,
			gl.RGBA,
			int32(image.Rect.Size().X),
			int32(image.Rect.Size().Y),
			0,
			gl.RGBA,
			gl.UNSIGNED_BYTE,
			gl.Ptr(image.Pix),
		)

		toonGl.Unbind()
		toonGl.Initialized = true

		renderModel.toonTextures[i] = toonGl
	}

	return nil
}
