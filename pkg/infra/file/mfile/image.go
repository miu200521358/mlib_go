// 指示: miu200521358
package mfile

import (
	"errors"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/ftrvxmtrx/tga"
	"github.com/miu200521358/dds/pkg/dds"
	"github.com/miu200521358/mlib_go/pkg/shared/base/merr"
	"golang.org/x/image/bmp"
	_ "golang.org/x/image/riff"
	_ "golang.org/x/image/tiff"
)

// LoadImage は指定パスの画像を読み込む。
func LoadImage(path string) (image.Image, error) {
	baseName := filepath.Base(path)
	open := func() (io.ReadCloser, error) {
		file, err := os.Open(path)
		if err != nil {
			return nil, merr.NewOsPackageError("画像ファイルの読み込みに失敗しました: "+baseName, err)
		}
		return file, nil
	}
	return loadImage(path, baseName, open)
}

// imageOpenFunc は画像データを取得する関数型。
type imageOpenFunc func() (io.ReadCloser, error)

var errUnsupportedImageFormat = errors.New("unsupported image format")

// loadImage は拡張子候補に従って画像を読み込む。
func loadImage(path string, displayName string, open imageOpenFunc) (image.Image, error) {
	candidates := imageExtensionCandidates(path)
	if len(candidates) == 0 {
		return nil, merr.NewImagePackageError("未対応の画像形式です: "+displayName, nil)
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
		return nil, merr.NewImagePackageError("画像のデコードに失敗しました: "+displayName, lastErr)
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
