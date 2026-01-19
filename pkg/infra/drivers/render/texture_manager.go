//go:build windows
// +build windows

// 指示: miu200521358
package render

import (
	"bytes"
	"embed"
	"errors"
	"fmt"
	"image"
	"image/draw"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/ftrvxmtrx/tga"
	"github.com/go-gl/gl/v4.3-core/gl"
	"github.com/miu200521358/dds/pkg/dds"
	"github.com/miu200521358/mlib_go/pkg/domain/mmath"
	"github.com/miu200521358/mlib_go/pkg/domain/model"
	"github.com/miu200521358/mlib_go/pkg/infra/base/logging"
	baseerr "github.com/miu200521358/mlib_go/pkg/shared/base/err"
	"golang.org/x/image/bmp"
	_ "golang.org/x/image/riff"
	_ "golang.org/x/image/tiff"
)

// ----------------------------------------------------------------------------
// 埋め込みリソース：toonフォルダ内のファイルをすべて埋め込む
// ----------------------------------------------------------------------------

//go:embed toon/*
var toonFiles embed.FS

// ----------------------------------------------------------------------------
// テクスチャ管理用：OpenGLのテクスチャ情報
// ----------------------------------------------------------------------------

type textureGl struct {
	// Texture は元のテクスチャ情報。
	*model.Texture
	// Id はOpenGLのテクスチャID。
	Id uint32
	// TextureType はテクスチャ種別。
	TextureType model.TextureType
	// TextureUnitId はOpenGLのテクスチャユニット。
	TextureUnitId uint32
	// TextureUnitNo はテクスチャユニット番号。
	TextureUnitNo uint32
	// IsGeneratedMipmap はミップマップ生成済みフラグ。
	IsGeneratedMipmap bool
	// Initialized は初期化済みフラグ。
	Initialized bool
}

// bind はテクスチャをバインドする。
func (texGl *textureGl) bind() {
	if texGl == nil || texGl.Id == 0 {
		return
	}
	gl.ActiveTexture(texGl.TextureUnitId)
	gl.BindTexture(gl.TEXTURE_2D, texGl.Id)

	// Toonテクスチャのみ wrap を REPEAT にする
	if texGl.TextureType == model.TEXTURE_TYPE_TOON {
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.REPEAT)
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.REPEAT)
	}

	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)

	if texGl.IsGeneratedMipmap {
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR_MIPMAP_LINEAR)
	} else {
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAX_LEVEL, 0)
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	}
}

// unbind はテクスチャのバインドを解除する。
func (texGl *textureGl) unbind() {
	gl.BindTexture(gl.TEXTURE_2D, 0)
}

// delete はテクスチャを削除する。
func (texGl *textureGl) delete() {
	if texGl == nil || texGl.Id == 0 {
		return
	}
	gl.DeleteTextures(1, &texGl.Id)
}

// ----------------------------------------------------------------------------
// TextureManager : 通常テクスチャ・Toonテクスチャ・スフィアテクスチャなどを一括管理
// ----------------------------------------------------------------------------

type TextureManager struct {
	textures     []*textureGl
	toonTextures []*textureGl
}

// NewTextureManager はTextureManagerを生成する。
func NewTextureManager() *TextureManager {
	return &TextureManager{
		textures:     nil,
		toonTextures: make([]*textureGl, 10),
	}
}

// LoadAllTextures はモデルに紐づくテクスチャをロードする。
func (tm *TextureManager) LoadAllTextures(windowIndex int, textures *model.TextureCollection, modelPath string) error {
	if tm == nil || textures == nil {
		return nil
	}
	var loadErr error
	// インデックスで参照できるようにスライスを確保する
	tm.textures = make([]*textureGl, textures.Len())

	for _, texture := range textures.Values() {
		if texture == nil || !texture.IsValid() {
			continue
		}
		texGl, err := tm.loadTextureGl(windowIndex, texture, modelPath)
		if err != nil {
			if loadErr == nil {
				loadErr = err
			}
			logging.DefaultLogger().Warn("テクスチャ読み込みに失敗しました: %v", err)
			continue
		}
		idx := texture.Index()
		if idx < 0 || idx >= len(tm.textures) {
			logging.DefaultLogger().Warn("テクスチャインデックスが範囲外です: %d", idx)
			continue
		}
		tm.textures[idx] = texGl
	}

	return loadErr
}

// Texture はインデックスに対応するテクスチャを返す。
func (tm *TextureManager) Texture(textureIndex int) *textureGl {
	if tm == nil || textureIndex < 0 || textureIndex >= len(tm.textures) {
		return nil
	}
	return tm.textures[textureIndex]
}

// LoadToonTextures は埋め込みトゥーンテクスチャをロードする。
func (tm *TextureManager) LoadToonTextures(windowIndex int) error {
	if tm == nil {
		return nil
	}
	if tm.toonTextures == nil || len(tm.toonTextures) != 10 {
		tm.toonTextures = make([]*textureGl, 10)
	}

	for i := 0; i < 10; i++ {
		filePath := fmt.Sprintf("toon/toon%02d.bmp", i+1)

		tex := model.NewTexture()
		tex.SetIndex(i)
		tex.SetName(filePath)
		tex.EnglishName = filePath
		tex.TextureType = model.TEXTURE_TYPE_TOON
		tex.SetValid(true)

		toonGl := &textureGl{
			Texture:     tex,
			TextureType: model.TEXTURE_TYPE_TOON,
		}

		gl.GenTextures(1, &toonGl.Id)

		// Toon用テクスチャユニットを設定
		switch windowIndex {
		case 0:
			toonGl.TextureUnitId = gl.TEXTURE10
			toonGl.TextureUnitNo = 10
		case 1:
			toonGl.TextureUnitId = gl.TEXTURE11
			toonGl.TextureUnitNo = 11
		case 2:
			toonGl.TextureUnitId = gl.TEXTURE12
			toonGl.TextureUnitNo = 12
		}

		img, err := loadImageFromResources(toonFiles, filePath)
		if err != nil {
			return err
		}
		image := convertToNRGBA(img)
		if image == nil {
			return baseerr.NewImagePackageError("トゥーンテクスチャの変換に失敗しました: "+filePath, nil)
		}

		toonGl.bind()

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

		toonGl.unbind()
		toonGl.Initialized = true
		tm.toonTextures[i] = toonGl
	}
	return nil
}

// ToonTexture は指定インデックスのトゥーンテクスチャを返す。
func (tm *TextureManager) ToonTexture(index int) *textureGl {
	if tm == nil || index < 0 || index >= len(tm.toonTextures) {
		return nil
	}
	return tm.toonTextures[index]
}

// Delete は管理しているテクスチャを削除する。
func (tm *TextureManager) Delete() {
	if tm == nil {
		return
	}
	for _, texGl := range tm.textures {
		if texGl != nil {
			texGl.delete()
		}
	}
	for _, toonGl := range tm.toonTextures {
		if toonGl != nil {
			toonGl.delete()
		}
	}
}

// loadTextureGl は単一テクスチャをOpenGL向けにロードする。
func (tm *TextureManager) loadTextureGl(windowIndex int, texture *model.Texture, modelPath string) (*textureGl, error) {
	if texture == nil || texture.Name() == "" {
		return nil, baseerr.NewImagePackageError("テクスチャ名が不正です", nil)
	}

	texGl := &textureGl{
		Texture:     texture,
		TextureType: texture.TextureType,
	}

	// モデルパス + テクスチャ相対パス
	texPath := filepath.Join(filepath.Dir(modelPath), texture.Name())

	img, err := loadImageFromFile(texPath)
	if err != nil {
		texGl.Initialized = false
		return texGl, err
	}
	image := convertToNRGBA(img)
	if image == nil {
		return texGl, baseerr.NewImagePackageError("テクスチャの変換に失敗しました: "+texture.Name(), nil)
	}

	texGl.IsGeneratedMipmap =
		mmath.IsPowerOfTwo(image.Rect.Size().X) && mmath.IsPowerOfTwo(image.Rect.Size().Y)

	gl.GenTextures(1, &texGl.Id)
	texGl.setTextureUnit(windowIndex)

	texGl.bind()

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
		logging.DefaultLogger().Debug("ミップマップ生成エラー: %s", texture.Name())
	}

	texGl.unbind()
	texGl.Initialized = true

	return texGl, nil
}

// setTextureUnit はwindowIndexとTextureTypeに応じてテクスチャユニットを決める。
func (texGl *textureGl) setTextureUnit(windowIndex int) {
	if texGl == nil {
		return
	}
	var baseUnit uint32
	switch windowIndex {
	case 0:
		baseUnit = gl.TEXTURE0
	case 1:
		baseUnit = gl.TEXTURE3
	case 2:
		baseUnit = gl.TEXTURE6
	}

	var offsetUnit uint32
	switch texGl.TextureType {
	case model.TEXTURE_TYPE_TEXTURE:
		offsetUnit = 0
	case model.TEXTURE_TYPE_TOON:
		offsetUnit = 1
	case model.TEXTURE_TYPE_SPHERE:
		offsetUnit = 2
	}

	texGl.TextureUnitId = baseUnit + offsetUnit

	baseVal := uint32(gl.TEXTURE0)
	texGl.TextureUnitNo = (texGl.TextureUnitId - baseVal)
}

// ----------------------------------------------------------------------------
// 画像読み込み補助
// ----------------------------------------------------------------------------

// imageOpenFunc は画像データを取得する関数型。
type imageOpenFunc func() (io.ReadCloser, error)

var errUnsupportedImageFormat = errors.New("unsupported image format")

// loadImageFromResources は埋め込みFSから画像を読み込む。
func loadImageFromResources(resources embed.FS, fileName string) (image.Image, error) {
	data, err := fs.ReadFile(resources, fileName)
	if err != nil {
		return nil, baseerr.NewFsPackageError("トゥーンテクスチャの読み込みに失敗しました: "+fileName, err)
	}
	open := func() (io.ReadCloser, error) {
		return io.NopCloser(bytes.NewReader(data)), nil
	}
	return loadImage(fileName, fileName, open)
}

// loadImageFromFile はファイルパスから画像を読み込む。
func loadImageFromFile(path string) (image.Image, error) {
	baseName := filepath.Base(path)
	open := func() (io.ReadCloser, error) {
		file, err := os.Open(path)
		if err != nil {
			return nil, baseerr.NewOsPackageError("テクスチャファイルの読み込みに失敗しました: "+baseName, err)
		}
		return file, nil
	}
	return loadImage(path, baseName, open)
}

// loadImage は拡張子に応じて画像を読み込む。
func loadImage(path string, displayName string, open imageOpenFunc) (image.Image, error) {
	candidates := imageExtensionCandidates(path)
	if len(candidates) == 0 {
		return nil, baseerr.NewImagePackageError("未対応の画像形式です: "+displayName, nil)
	}
	var lastErr error
	for _, ext := range candidates {
		reader, err := open()
		if err != nil {
			return nil, err
		}
		img, err := loadImageByExtension(reader, ext)
		_ = reader.Close()
		if err == nil {
			return img, nil
		}
		lastErr = err
	}

	reader, err := open()
	if err != nil {
		return nil, err
	}
	defer reader.Close()
	img, _, err := image.Decode(reader)
	if err != nil {
		if lastErr == nil {
			lastErr = err
		}
		return nil, baseerr.NewImagePackageError("テクスチャのデコードに失敗しました: "+displayName, lastErr)
	}
	return img, nil
}

// imageExtensionCandidates は拡張子に応じたデコード順を返す。
func imageExtensionCandidates(path string) []string {
	ext := strings.ToLower(strings.TrimPrefix(filepath.Ext(path), "."))
	switch ext {
	case "png":
		return []string{"png", "gif", "jpg", "bmp", "tga", "dds"}
	case "tga":
		return []string{"tga", "png", "gif", "jpg", "bmp", "dds"}
	case "gif":
		return []string{"gif", "png", "jpg", "bmp", "tga", "dds"}
	case "dds":
		return []string{"dds", "png", "gif", "jpg", "bmp", "tga"}
	case "jpg", "jpeg":
		return []string{"jpg", "png", "gif", "bmp", "tga", "dds"}
	case "bmp":
		return []string{"bmp", "png", "gif", "jpg", "tga", "dds"}
	case "spa", "sph":
		return []string{"bmp", "png", "gif", "jpg", "tga", "dds"}
	default:
		return nil
	}
}

// loadImageByExtension は指定形式でデコードを試みる。
func loadImageByExtension(reader io.Reader, ext string) (image.Image, error) {
	switch ext {
	case "png":
		return png.Decode(reader)
	case "tga":
		return tga.Decode(reader)
	case "gif":
		return gif.Decode(reader)
	case "dds":
		return dds.Decode(reader)
	case "jpg", "jpeg":
		return jpeg.Decode(reader)
	case "bmp":
		return bmp.Decode(reader)
	default:
		return nil, errUnsupportedImageFormat
	}
}

// convertToNRGBA は画像を NRGBA に変換する。
func convertToNRGBA(img image.Image) *image.NRGBA {
	if img == nil {
		return nil
	}
	if rgba, ok := img.(*image.NRGBA); ok {
		return rgba
	}
	bounds := img.Bounds()
	rgba := image.NewNRGBA(bounds)
	draw.Draw(rgba, bounds, img, bounds.Min, draw.Src)
	return rgba
}
