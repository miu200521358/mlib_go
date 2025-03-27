//go:build windows
// +build windows

package app

import (
	"github.com/miu200521358/mlib_go/pkg/config/mconfig"
	"github.com/miu200521358/mlib_go/pkg/domain/mmath"
	"golang.org/x/sys/windows"
)

var (
	user32               = windows.NewLazySystemDLL("user32.dll")
	procGetSystemMetrics = user32.NewProc("GetSystemMetrics")
)

const (
	SM_CXSCREEN = 0
	SM_CYSCREEN = 1
)

func getSystemMetrics(nIndex int) int {
	ret, _, _ := procGetSystemMetrics.Call(uintptr(nIndex))
	return int(ret)
}

func GetCenterSizeAndWidth(appConfig *mconfig.AppConfig, viewerCount int) (
	widths []int, heights []int, positionXs []int, positionYs []int,
) {
	// スクリーンの解像度を取得
	screenWidth := getSystemMetrics(SM_CXSCREEN)
	screenHeight := getSystemMetrics(SM_CYSCREEN)

	// ウィンドウのサイズ
	widths = make([]int, 0, viewerCount+1)
	heights = make([]int, 0, viewerCount+1)
	// ウィンドウの位置
	positionXs = make([]int, viewerCount+1)
	positionYs = make([]int, viewerCount+1)

	widths = append(widths, appConfig.ControlWindowSize.Width)
	heights = append(heights, appConfig.ControlWindowSize.Height)

	for range viewerCount {
		widths = append(widths, appConfig.ViewWindowSize.Width)
		heights = append(heights, appConfig.ViewWindowSize.Height)
	}

	var appWidth, appHeight int
	widthRatio := 1.0
	heightRatio := 1.0
	if appConfig.Horizontal {
		appWidth = mmath.Sum(widths)
		appHeight = mmath.Max(heights)
	} else {
		appWidth = mmath.Max(widths)
		appHeight = mmath.Sum(heights)
	}

	if appHeight > screenHeight-50 {
		heightRatio = float64(screenHeight-50) / float64(appHeight)
	}
	if appWidth > screenWidth-50 {
		widthRatio = float64(screenWidth-50) / float64(appWidth)
	}

	// リサイズ
	if widthRatio < 1 || heightRatio < 1 {
		for n := range widths {
			w := int(float64(widths[n]) * widthRatio)
			h := int(float64(heights[n]) * heightRatio)

			widths[n] = w
			heights[n] = h
		}

		if appConfig.Horizontal {
			appWidth = mmath.Sum(widths)
			appHeight = mmath.Max(heights)
		} else {
			appWidth = mmath.Max(widths)
			appHeight = mmath.Sum(heights)
		}
	}

	// ウィンドウを中央に配置
	if appConfig.Horizontal {
		centerX := (screenWidth - appWidth) / 2
		centerY := (screenHeight - appHeight) / 2

		centerX += mmath.Sum(widths[1:])
		positionXs[0] = centerX
		positionYs[0] = centerY

		for n := range viewerCount {
			centerX -= widths[n+1]
			positionXs[n+1] = centerX
			positionYs[n+1] = centerY
		}
	} else {
		centerX := (screenWidth - appWidth) / 2
		centerY := (screenHeight - appHeight) / 2

		positionXs[0] = centerX
		positionYs[0] = centerY

		for n := range viewerCount {
			centerY += heights[n+1]
			positionXs[n+1] = centerX
			positionYs[n+1] = centerY
		}
	}

	return
}
