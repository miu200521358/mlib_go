//go:build windows
// +build windows

package render

import (
	"embed"
	"fmt"
	"path/filepath"

	"github.com/go-gl/gl/v4.4-core/gl"
	"github.com/miu200521358/mlib_go/pkg/config/mi18n"
	"github.com/miu200521358/mlib_go/pkg/config/mlog"
	"github.com/miu200521358/mlib_go/pkg/domain/mmath"
	"github.com/miu200521358/mlib_go/pkg/domain/pmx"
	"github.com/miu200521358/mlib_go/pkg/infrastructure/mfile"
)

// ----------------------------------------------------------------------------
// 埋め込みリソース：toonフォルダ内のファイルをすべて埋め込む
// プロジェクト構成によっては外部ファイルや別の管理方法でもOK
// ----------------------------------------------------------------------------

//go:embed toon/*
var toonFiles embed.FS

// ----------------------------------------------------------------------------
// テクスチャ管理用：OpenGLのテクスチャ情報
// ----------------------------------------------------------------------------

type textureGl struct {
	// pmx.Texture との紐づけ情報
	*pmx.Texture

	// 実際の OpenGL テクスチャID
	Id uint32

	// テクスチャの種類 (通常 or Toon or スフィア)
	TextureType pmx.TextureType

	// OpenGLのテクスチャユニット (gl.TEXTURE0 等)
	TextureUnitId uint32

	// テクスチャユニット番号 (0,1,2,...)
	TextureUnitNo uint32

	// ミップマップが生成済みかどうか
	IsGeneratedMipmap bool

	// 初期化済みフラグ
	Initialized bool
}

// bind / unbind / delete は OpenGL のバインド処理

func (texGl *textureGl) bind() {
	gl.ActiveTexture(texGl.TextureUnitId)
	gl.BindTexture(gl.TEXTURE_2D, texGl.Id)

	// Toonテクスチャの場合はリピートするなど、場合に応じて設定
	if texGl.TextureType == pmx.TEXTURE_TYPE_TOON {
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

func (texGl *textureGl) unbind() {
	gl.BindTexture(gl.TEXTURE_2D, 0)
}

func (texGl *textureGl) delete() {
	gl.DeleteTextures(1, &texGl.Id)
}

// ----------------------------------------------------------------------------
// TextureManager : 通常テクスチャ・Toonテクスチャ・スフィアテクスチャなどを一括管理
// ----------------------------------------------------------------------------

type TextureManager struct {
	// テクスチャを配列 or マップで保持
	// 例: インデックスを pmx.Texture.Index() としてアクセス
	textures     []*textureGl
	toonTextures []*textureGl
}

// NewTextureManager : TextureManagerの生成
func NewTextureManager() *TextureManager {
	return &TextureManager{
		textures:     nil,
		toonTextures: make([]*textureGl, 10), // Toonは10個分(toon01〜toon10)想定
	}
}

// LoadAllTextures : PMXが持つ "通常テクスチャ" / "スフィアテクスチャ" をロードし、
// インデックスを揃えた配列 (tm.textures) に格納する
func (tm *TextureManager) LoadAllTextures(windowIndex int, textures *pmx.Textures, modelPath string) error {
	// まず textures の長さに応じてスライスを確保
	tm.textures = make([]*textureGl, textures.Length())

	textures.ForEach(func(index int, texture *pmx.Texture) bool {
		texGl, err := tm.loadTextureGl(windowIndex, texture, modelPath)
		if err != nil {
			mlog.W(fmt.Sprintf("texture initialize error: %s", err))
			return true // エラーでも次のテクスチャを処理する
		}
		// インデックス位置に格納
		tm.textures[texture.Index()] = texGl
		return true
	})

	return nil
}

// Texture : 引数の pmx.Texture.Index() に対応する textureGl を返す
// 見つからなければ nil
func (tm *TextureManager) Texture(textureIndex int) *textureGl {
	if textureIndex < 0 || textureIndex >= len(tm.textures) {
		return nil
	}
	return tm.textures[textureIndex]
}

// LoadToonTextures : toon01〜toon10 を埋め込みリソースから読み込み
func (tm *TextureManager) LoadToonTextures(windowIndex int) error {
	for i := 0; i < 10; i++ {
		filePath := fmt.Sprintf("toon/toon%02d.bmp", i+1)

		// テクスチャ領域
		toonGl := &textureGl{
			Texture:     pmx.NewTexture(), // 一応 pmx.Texture構造体を割り当て
			TextureType: pmx.TEXTURE_TYPE_TOON,
		}
		toonGl.Texture.SetIndex(i)       // インデックスは 0〜9
		toonGl.Texture.SetName(filePath) // ファイル名
		toonGl.Texture.SetEnglishName(filePath)
		toonGl.Texture.SetValid(true) // 有効

		// OpenGLテクスチャ生成
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

		// ファイルを埋め込みリソースからロード
		img, err := mfile.LoadImageFromResources(toonFiles, filePath)
		if err != nil {
			return err
		}
		image := mfile.ConvertToNRGBA(img)

		toonGl.bind()

		// glTexImage2D
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

// ToonTexture : 引数のインデックス (0〜9) の Toonテクスチャを返す
func (tm *TextureManager) ToonTexture(index int) *textureGl {
	if index < 0 || index >= len(tm.toonTextures) {
		return nil
	}
	return tm.toonTextures[index]
}

// Delete : 管理しているすべてのテクスチャを削除
func (tm *TextureManager) Delete() {
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

// ----------------------------------------------------------------------------
// 内部メソッド : テクスチャのロード処理
// ----------------------------------------------------------------------------

func (tm *TextureManager) loadTextureGl(
	windowIndex int,
	texture *pmx.Texture,
	modelPath string,
) (*textureGl, error) {

	texGl := &textureGl{
		Texture:     texture, // pmx.Texture と紐づけ
		TextureType: texture.TextureType,
	}

	// モデルパス + テクスチャ相対パス
	texPath := filepath.Join(filepath.Dir(modelPath), texture.Name())

	valid, err := mfile.ExistsFile(texPath)
	if !valid || err != nil {
		texGl.Initialized = false
		errMsg := ""
		if err != nil {
			errMsg = err.Error()
		}
		return texGl, fmt.Errorf("not found texture file: %s, error: %s", texPath, errMsg)
	}

	// 画像を読み込み
	img, err := mfile.LoadImage(texPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load image: %s, error: %s", texPath, err.Error())
	}
	image := mfile.ConvertToNRGBA(img)

	// ミップマップが作れるか (2の累乗かどうか)
	texGl.IsGeneratedMipmap =
		mmath.IsPowerOfTwo(image.Rect.Size().X) && mmath.IsPowerOfTwo(image.Rect.Size().Y)

	// OpenGLテクスチャ生成
	gl.GenTextures(1, &texGl.Id)

	// テクスチャユニット決定 (windowIndex によるオフセットを加味)
	texGl.setTextureUnit(windowIndex)

	// バインドしてデータ転送
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

	// ミップマップ生成
	if texGl.IsGeneratedMipmap {
		gl.GenerateMipmap(gl.TEXTURE_2D)
	} else {
		mlog.D(mi18n.T("ミップマップ生成エラー", map[string]interface{}{"Name": texture.Name()}))
	}

	texGl.unbind()
	texGl.Initialized = true

	return texGl, nil
}

// setTextureUnit : windowIndexとTextureTypeに応じてTextureUnitId, TextureUnitNoを決める
func (texGl *textureGl) setTextureUnit(windowIndex int) {
	// windowIndexごとに 0,3,6.. など
	// TextureTypeごとに +0, +1, +2.. のように振り分け
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
	case pmx.TEXTURE_TYPE_TEXTURE:
		offsetUnit = 0
	case pmx.TEXTURE_TYPE_TOON:
		offsetUnit = 1
	case pmx.TEXTURE_TYPE_SPHERE:
		offsetUnit = 2
	}

	// 実際のユニットIDを決定
	texGl.TextureUnitId = baseUnit + offsetUnit

	// ユニット番号も合わせて設定
	// 例: TEXTURE0 → 0, TEXTURE1 → 1, TEXTURE3 → 3, etc.
	// OpenGL定数 TEXTURE0 は 33984 (0x84C0) なので注意
	// ここでは簡易的に (base - TEXTURE0) + offset で計算
	baseVal := uint32(gl.TEXTURE0) // 0x84C0
	texGl.TextureUnitNo = (texGl.TextureUnitId - baseVal)
}
