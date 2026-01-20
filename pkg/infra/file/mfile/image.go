// 指示: miu200521358
package mfile

import (
	"image"
	"image/draw"
	"io/fs"
	"os"

	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	"github.com/miu200521358/mlib_go/pkg/shared/base/merr"
)

// ReadImage は画像ファイルを読み込む。
func ReadImage(path string) (image.Image, error) {
	if path == "" {
		return nil, newFileNotFound("パスが空です", nil)
	}
	file, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, newFileNotFound("画像が存在しません: "+path, nil)
		}
		return nil, newFileReadFailed("画像ファイルを開けません: "+path, merr.NewOsPackageError("os.Openに失敗しました", err))
	}
	defer file.Close()
	img, _, err := image.Decode(file)
	if err != nil {
		return nil, newImageDecodeFailed("画像のデコードに失敗しました: "+path, merr.NewImagePackageError("image.Decodeに失敗しました", err))
	}
	return img, nil
}

// ReadImageFromFS は埋め込みFSから画像を読み込む。
func ReadImageFromFS(resources fs.FS, fileName string) (image.Image, error) {
	if fileName == "" {
		return nil, newFileNotFound("パスが空です", nil)
	}
	file, err := resources.Open(fileName)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, newFileNotFound("画像が存在しません: "+fileName, nil)
		}
		return nil, newFileReadFailed("画像ファイルを開けません: "+fileName, merr.NewFsPackageError("fs.Openに失敗しました", err))
	}
	defer file.Close()
	img, _, err := image.Decode(file)
	if err != nil {
		return nil, newImageDecodeFailed("画像のデコードに失敗しました: "+fileName, merr.NewImagePackageError("image.Decodeに失敗しました", err))
	}
	return img, nil
}

// ToNRGBA はNRGBA形式へ変換する。
func ToNRGBA(img image.Image) *image.NRGBA {
	if img == nil {
		return nil
	}
	if converted, ok := img.(*image.NRGBA); ok {
		return converted
	}
	bounds := img.Bounds()
	nrgba := image.NewNRGBA(bounds)
	draw.Draw(nrgba, bounds, img, bounds.Min, draw.Src)
	return nrgba
}

// FlipImageY は縦方向に反転する。
func FlipImageY(img *image.RGBA) *image.RGBA {
	if img == nil {
		return nil
	}
	bounds := img.Bounds()
	flipped := image.NewRGBA(bounds)
	width := bounds.Dx()
	height := bounds.Dy()
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			flipped.Set(x, height-1-y, img.At(x, y))
		}
	}
	return flipped
}
