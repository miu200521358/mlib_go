//go:build windows
// +build windows

package pmx

import (
	"embed"
	"fmt"
	"path/filepath"

	"github.com/go-gl/gl/v4.4-core/gl"

	"github.com/miu200521358/mlib_go/pkg/mmath"
	"github.com/miu200521358/mlib_go/pkg/mutils"
	"github.com/miu200521358/mlib_go/pkg/mutils/mi18n"
	"github.com/miu200521358/mlib_go/pkg/mutils/mlog"
)

//go:embed toon/*
var toonFiles embed.FS

type TextureGL struct {
	Id                uint32      // OpenGLテクスチャID
	Valid             bool        // テクスチャフルパスが有効であるか否か
	TextureType       TextureType // テクスチャ種別
	TextureUnitId     uint32      // テクスチャ種類別描画先ユニットID
	TextureUnitNo     uint32      // テクスチャ種類別描画先ユニット番号
	IsGeneratedMipmap bool        // ミップマップが生成されているか否か
}

func (t *TextureGL) Bind() {
	if !t.Valid {
		return
	}

	gl.ActiveTexture(t.TextureUnitId)
	gl.BindTexture(gl.TEXTURE_2D, t.Id)

	if t.TextureType == TEXTURE_TYPE_TOON {
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

func (t *TextureGL) Unbind() {
	if !t.Valid {
		return
	}

	gl.BindTexture(gl.TEXTURE_2D, 0)
}

func (t *TextureGL) delete() {
	if t.Valid {
		gl.DeleteTextures(1, &t.Id)
	}
}

func (t *Texture) GL(
	modelPath string,
	textureType TextureType,
	windowIndex int,
) *TextureGL {
	tGl := &TextureGL{}

	if t.Initialized && t.Valid {
		// 既にフラグが立ってたら描画初期化済み
		// 共有Toonテクスチャの場合、既に初期化済み
		tGl.Id = t.glId
		tGl.Valid = t.Valid
		tGl.TextureType = t.TextureType
		tGl.TextureUnitId = t.textureUnitId
		tGl.TextureUnitNo = t.textureUnitNo
		tGl.IsGeneratedMipmap = t.IsGeneratedMipmap
		return tGl
	}

	// 通常テクスチャ
	texPath := filepath.Join(filepath.Dir(modelPath), t.Name)

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
		t.Image = mutils.ConvertToNRGBA(*img)
		*img = nil
	}

	if !t.Valid {
		t.Initialized = true
		return nil
	}

	t.IsGeneratedMipmap =
		mmath.IsPowerOfTwo(t.Image.Rect.Size().X) && mmath.IsPowerOfTwo(t.Image.Rect.Size().Y)

	t.TextureType = textureType

	// テクスチャオブジェクト生成
	gl.GenTextures(1, &tGl.Id)
	t.glId = tGl.Id

	// テクスチャ種別によってテクスチャユニットを変更
	if windowIndex == 0 {
		switch textureType {
		case TEXTURE_TYPE_TEXTURE:
			t.textureUnitId = gl.TEXTURE0
			t.textureUnitNo = 0
		case TEXTURE_TYPE_TOON:
			t.textureUnitId = gl.TEXTURE1
			t.textureUnitNo = 1
		case TEXTURE_TYPE_SPHERE:
			t.textureUnitId = gl.TEXTURE2
			t.textureUnitNo = 2
		}
	} else if windowIndex == 1 {
		switch textureType {
		case TEXTURE_TYPE_TEXTURE:
			t.textureUnitId = gl.TEXTURE3
			t.textureUnitNo = 3
		case TEXTURE_TYPE_TOON:
			t.textureUnitId = gl.TEXTURE4
			t.textureUnitNo = 4
		case TEXTURE_TYPE_SPHERE:
			t.textureUnitId = gl.TEXTURE5
			t.textureUnitNo = 5
		}
	} else if windowIndex == 2 {
		switch textureType {
		case TEXTURE_TYPE_TEXTURE:
			t.textureUnitId = gl.TEXTURE6
			t.textureUnitNo = 6
		case TEXTURE_TYPE_TOON:
			t.textureUnitId = gl.TEXTURE7
			t.textureUnitNo = 7
		case TEXTURE_TYPE_SPHERE:
			t.textureUnitId = gl.TEXTURE8
			t.textureUnitNo = 8
		}
	}

	tGl.Valid = t.Valid
	tGl.TextureType = t.TextureType
	tGl.TextureUnitId = t.textureUnitId
	tGl.TextureUnitNo = t.textureUnitNo
	tGl.IsGeneratedMipmap = t.IsGeneratedMipmap

	tGl.Bind()

	gl.TexImage2D(
		gl.TEXTURE_2D,
		0,
		gl.RGBA,
		int32(t.Image.Rect.Size().X),
		int32(t.Image.Rect.Size().Y),
		0,
		gl.RGBA,
		gl.UNSIGNED_BYTE,
		gl.Ptr(t.Image.Pix),
	)

	if t.IsGeneratedMipmap {
		gl.GenerateMipmap(gl.TEXTURE_2D)
	} else {
		mlog.D(mi18n.T("ミップマップ生成エラー", map[string]interface{}{"Name": t.Name}))
	}

	tGl.Unbind()

	// 描画初期化
	t.Initialized = true
	return tGl
}

func (t *ToonTextures) initGl(
	windowIndex int,
) error {
	for i := 0; i < 10; i++ {
		filePath := fmt.Sprintf("toon/toon%02d.bmp", i+1)

		toon := NewTexture()
		toon.Index = i
		toon.Name = filePath
		toon.TextureType = TEXTURE_TYPE_TOON
		toon.Path = filePath

		img, err := mutils.LoadImageFromResources(toonFiles, filePath)
		if err != nil {
			return err
		}
		toon.Image = mutils.ConvertToNRGBA(*img)
		toon.Valid = true
		*img = nil

		tGl := &TextureGL{}

		// テクスチャオブジェクト生成
		gl.GenTextures(1, &tGl.Id)
		toon.glId = tGl.Id

		// Toon用テクスチャユニットを設定
		if windowIndex == 0 {
			toon.textureUnitId = gl.TEXTURE10
			toon.textureUnitNo = 10
		} else if windowIndex == 1 {
			toon.textureUnitId = gl.TEXTURE11
			toon.textureUnitNo = 11
		} else if windowIndex == 2 {
			toon.textureUnitId = gl.TEXTURE12
			toon.textureUnitNo = 12
		}

		tGl.Valid = toon.Valid
		tGl.TextureType = toon.TextureType
		tGl.TextureUnitId = toon.textureUnitId
		// tGl.IsGeneratedMipmap = toon.IsGeneratedMipmap	// Toonテクスチャはミップマップを生成しない

		tGl.Bind()

		gl.TexImage2D(
			gl.TEXTURE_2D,
			0,
			gl.RGBA,
			int32(toon.Image.Rect.Size().X),
			int32(toon.Image.Rect.Size().Y),
			0,
			gl.RGBA,
			gl.UNSIGNED_BYTE,
			gl.Ptr(toon.Image.Pix),
		)

		tGl.Unbind()
		toon.Initialized = true

		t.Update(toon)
	}

	return nil
}
