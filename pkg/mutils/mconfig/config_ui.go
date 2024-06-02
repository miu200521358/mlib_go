//go:build windows
// +build windows

package mconfig

import (
	"embed"

	"github.com/miu200521358/walk/pkg/walk"

	"github.com/miu200521358/mlib_go/pkg/mutils"
)

// LoadImageFile 画像ファイルの読み込み
func LoadImageFile(resourceFiles embed.FS, imagePath string, dpi int) (walk.Image, error) {
	image, err := mutils.LoadImageFromResources(resourceFiles, imagePath)
	if err != nil {
		return nil, err
	}
	img, err := walk.NewIconFromImageForDPI(*image, dpi)
	if err != nil {
		return nil, err
	}
	return img, nil
}
