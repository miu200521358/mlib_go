package mfile

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
	"golang.org/x/image/bmp"
	_ "golang.org/x/image/riff"
	_ "golang.org/x/image/tiff"
)

// 指定されたパスから画像を読み込む
func LoadImage(path string) (image.Image, error) {
	return loadImage(path, func(path string) (io.Reader, error) {
		file, err := os.Open(path)
		if err != nil {
			return nil, err
		}
		return file, nil
	}, func(file io.Reader) {
		file.(io.Closer).Close()
	})
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
func LoadImageFromResources(resources embed.FS, fileName string) (image.Image, error) {
	return loadImage(fileName, func(path string) (io.Reader, error) {
		fileData, err := fs.ReadFile(resources, fileName)
		if err != nil {
			return nil, fmt.Errorf("image not found: %v", err)
		}
		return bytes.NewReader(fileData), nil
	}, func(file io.Reader) {})
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

func loadImage(path string, loadFunc func(path string) (io.Reader, error), closeFunc func(file io.Reader)) (image.Image, error) {
	paths := strings.Split(path, ".")
	if len(paths) < 2 {
		return nil, fmt.Errorf("invalid file path: %s", path)
	}

	extensions := []string{}

	switch strings.ToLower(paths[len(paths)-1]) {
	case "png":
		extensions = append(extensions, "png", "gif", "jpg", "bmp", "tga", "dds")
	case "tga":
		extensions = append(extensions, "tga", "png", "gif", "jpg", "bmp", "dds")
	case "gif":
		extensions = append(extensions, "gif", "png", "jpg", "bmp", "tga", "dds")
	case "dds":
		extensions = append(extensions, "dds", "png", "gif", "jpg", "bmp", "tga")
	case "jpg", "jpeg":
		extensions = append(extensions, "jpg", "png", "gif", "bmp", "tga", "dds")
	case "bmp":
		extensions = append(extensions, "bmp", "png", "gif", "jpg", "tga", "dds")
	case "spa", "sph":
		// スフィアファイルはまずbmpとして読み込む
		extensions = append(extensions, "bmp", "png", "gif", "jpg", "tga", "dds")
	}

	// 拡張子に合わせて画像を読み込む
	for _, extension := range extensions {
		img, err := loadImageByExtension(path, extension, loadFunc, closeFunc)
		if err == nil {
			return img, nil
		}
	}

	// どの拡張子にも合致しなかった場合はとりあえずそのまま呼んでみる
	file, err := loadFunc(path)
	if err != nil {
		return nil, err
	}
	defer closeFunc(file)

	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}

	return img, nil
}

// loadImageByExtension 拡張子に合わせて画像を読み込む
func loadImageByExtension(path, extension string, loadFunc func(path string) (io.Reader, error), closeFunc func(file io.Reader)) (image.Image, error) {
	file, err := loadFunc(path)
	if err != nil {
		return nil, err
	}
	defer closeFunc(file)

	switch extension {
	case "png":
		if img, err := png.Decode(file); err == nil {
			return img, nil
		} else {
			return nil, err
		}
	case "tga":
		if img, err := tga.Decode(file); err == nil {
			return img, nil
		} else {
			return nil, err
		}
	case "gif":
		if img, err := gif.Decode(file); err == nil {
			return img, nil
		} else {
			return nil, err
		}
	case "dds":
		if img, err := dds.Decode(file); err == nil {
			return img, nil
		} else {
			return nil, err
		}
	case "jpg":
		if img, err := jpeg.Decode(file); err == nil {
			return img, nil
		} else {
			return nil, err
		}
	case "bmp":
		if img, err := bmp.Decode(file); err == nil {
			return img, nil
		} else {
			return nil, err
		}
	}

	return nil, fmt.Errorf("unsupported image format: %s", extension)
}
