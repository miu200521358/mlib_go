package mutils

import (
	"bytes"
	"embed"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	"image/png"
	"io/fs"
	"os"
	"strings"

	"github.com/ftrvxmtrx/tga"
	"golang.org/x/image/bmp"
	_ "golang.org/x/image/riff"
	_ "golang.org/x/image/tiff"
)

// 指定されたパスから画像を読み込む
func LoadImage(path string) (image.Image, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	if strings.ToLower(path[len(path)-4:]) == ".png" {
		img, err := png.Decode(file)
		if err != nil {
			return nil, err
		}

		return img, nil
	} else if strings.ToLower(path[len(path)-4:]) == ".tga" {
		img, err := tga.Decode(file)
		if err != nil {
			return nil, err
		}

		return img, nil
	} else if strings.ToLower(path[len(path)-4:]) == ".spa" || strings.ToLower(path[len(path)-4:]) == ".bmp" {
		// スフィアファイルはbmpとして読み込む
		img, err := bmp.Decode(file)
		if err != nil {
			return nil, err
		}

		return img, nil
	}
	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}

	return img, nil
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
func LoadImageFromResources(resourceFiles embed.FS, fileName string) (image.Image, error) {
	fileData, err := fs.ReadFile(resourceFiles, fileName)
	if err != nil {
		return nil, fmt.Errorf("image not found: %v", err)
	}
	file := bytes.NewReader(fileData)

	if strings.ToLower(fileName[len(fileName)-4:]) == ".png" {
		img, err := png.Decode(file)
		if err != nil {
			return nil, err
		}

		return img, nil
	} else if strings.ToLower(fileName[len(fileName)-4:]) == ".tga" {
		img, err := tga.Decode(file)
		if err != nil {
			return nil, err
		}

		return img, nil
	} else if strings.ToLower(fileName[len(fileName)-4:]) == ".spa" ||
		strings.ToLower(fileName[len(fileName)-4:]) == ".bmp" {
		// スフィアファイルはbmpとして読み込む
		img, err := bmp.Decode(file)
		if err != nil {
			return nil, err
		}

		return img, nil
	}
	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}

	return img, nil
}
