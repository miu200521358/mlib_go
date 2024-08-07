package mutils

import (
	"bytes"
	"embed"
	"fmt"
	"image"
	"image/draw"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"io/fs"
	"os"
	"strings"

	"github.com/ftrvxmtrx/tga"
	"github.com/miu200521358/dds/pkg/dds"
	"github.com/miu200521358/walk/pkg/walk"
	"golang.org/x/image/bmp"
	_ "golang.org/x/image/riff"
	_ "golang.org/x/image/tiff"
)

// LoadWalkIcon 画像ファイルの読み込み
func LoadWalkIcon(resources embed.FS, imagePath string, dpi int) (*walk.Icon, error) {
	image, err := LoadImageFromResources(resources, imagePath)
	if err != nil {
		return nil, err
	}
	img, err := walk.NewIconFromImageForDPI(*image, dpi)
	if err != nil {
		return nil, err
	}
	return img, nil
}

// 指定されたパスから画像を読み込む
func LoadImage(path string) (*image.Image, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return loadImage(path, file)
}

// 指定された画像を反転させる
func FlipImage(img *image.RGBA) *image.RGBA {
	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y

	flipped := image.NewRGBA(image.Rect(0, 0, width, height))

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			srcX := width - x - 1
			srcY := y
			srcColor := img.At(srcX, srcY)
			flipped.Set(x, y, srcColor)
		}
	}

	return flipped
}

// ReadIconFile アイコンファイルの読み込み
func LoadImageFromResources(resources embed.FS, fileName string) (*image.Image, error) {
	fileData, err := fs.ReadFile(resources, fileName)
	if err != nil {
		return nil, fmt.Errorf("image not found: %v", err)
	}
	file := bytes.NewReader(fileData)

	return loadImage(fileName, file)
}

func ConvertToNRGBA(img image.Image) *image.NRGBA {
	// 画像がすでに*image.NRGBA型の場合はそのまま返す
	if rgba, ok := img.(*image.NRGBA); ok {
		return rgba
	}

	// それ以外の場合は、新しい*image.NRGBAイメージに描画して変換する
	bounds := img.Bounds()
	rgba := image.NewNRGBA(bounds)
	draw.Draw(rgba, rgba.Bounds(), img, bounds.Min, draw.Src)

	return rgba
}

func loadImage(path string, file io.Reader) (*image.Image, error) {
	paths := strings.Split(path, ".")
	if len(paths) < 2 {
		return nil, fmt.Errorf("invalid file path: %s", path)
	}

	switch paths[len(paths)-1] {
	case "png":
		img, err := png.Decode(file)
		if err != nil {
			return nil, err
		}

		return &img, nil
	case "tga":
		img, err := tga.Decode(file)
		if err != nil {
			return nil, err
		}

		return &img, nil
	case "gif":
		img, err := gif.Decode(file)
		if err != nil {
			return nil, err
		}

		return &img, nil
	case "dds":
		img, err := dds.Decode(file)
		if err != nil {
			return nil, err
		}

		return &img, nil
	case "jpg", "jpeg":
		img, err := jpeg.Decode(file)
		if err != nil {
			return nil, err
		}

		return &img, nil
	case "spa", "bmp":
		// スフィアファイルはbmpとして読み込む
		img, err := bmp.Decode(file)
		if err != nil {
			return nil, err
		}

		return &img, nil
	}

	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}

	return &img, nil
}
