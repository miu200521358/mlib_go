//go:build windows
// +build windows

// 指示: miu200521358
package app

import (
	"github.com/miu200521358/mlib_go/pkg/domain/mmath"
	"github.com/miu200521358/mlib_go/pkg/shared/base/config"
	"golang.org/x/sys/windows"
)

var (
	user32               = windows.NewLazySystemDLL("user32.dll")
	procGetSystemMetrics = user32.NewProc("GetSystemMetrics")
)

const (
	smCxScreen          = 0
	smCyScreen          = 1
	defaultWindowWidth  = 512
	defaultWindowHeight = 768
)

// getSystemMetrics はWin32のシステムメトリクスを取得する。
func getSystemMetrics(index int) int {
	ret, _, _ := procGetSystemMetrics.Call(uintptr(index))
	return int(ret)
}

// GetCenterSizeAndWidth は画面中央に配置するサイズと位置を算出する。
func GetCenterSizeAndWidth(appConfig *config.AppConfig, viewerCount int) (
	widths []int, heights []int, positionXs []int, positionYs []int,
) {
	if viewerCount < 0 {
		viewerCount = 0
	}

	screenWidth := getSystemMetrics(smCxScreen)
	screenHeight := getSystemMetrics(smCyScreen)
	if screenWidth <= 0 {
		screenWidth = defaultWindowWidth
	}
	if screenHeight <= 0 {
		screenHeight = defaultWindowHeight
	}

	controlWidth := defaultWindowWidth
	controlHeight := defaultWindowHeight
	viewerWidth := defaultWindowWidth
	viewerHeight := defaultWindowHeight
	isHorizontal := true
	if appConfig != nil {
		if appConfig.ControlWindowSize.Width > 0 {
			controlWidth = appConfig.ControlWindowSize.Width
		}
		if appConfig.ControlWindowSize.Height > 0 {
			controlHeight = appConfig.ControlWindowSize.Height
		}
		if appConfig.ViewerWindowSize.Width > 0 {
			viewerWidth = appConfig.ViewerWindowSize.Width
		}
		if appConfig.ViewerWindowSize.Height > 0 {
			viewerHeight = appConfig.ViewerWindowSize.Height
		}
		isHorizontal = appConfig.Horizontal
	}

	widths = make([]int, 0, viewerCount+1)
	heights = make([]int, 0, viewerCount+1)
	positionXs = make([]int, viewerCount+1)
	positionYs = make([]int, viewerCount+1)

	widths = append(widths, controlWidth)
	heights = append(heights, controlHeight)
	for i := 0; i < viewerCount; i++ {
		widths = append(widths, viewerWidth)
		heights = append(heights, viewerHeight)
	}

	appWidth := 0
	appHeight := 0
	widthRatio := 1.0
	heightRatio := 1.0
	if isHorizontal {
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

	if widthRatio < 1 || heightRatio < 1 {
		for i := range widths {
			widths[i] = int(float64(widths[i]) * widthRatio)
			heights[i] = int(float64(heights[i]) * heightRatio)
		}
		if isHorizontal {
			appWidth = mmath.Sum(widths)
			appHeight = mmath.Max(heights)
		} else {
			appWidth = mmath.Max(widths)
			appHeight = mmath.Sum(heights)
		}
	}

	if isHorizontal {
		centerX := (screenWidth - appWidth) / 2
		centerY := (screenHeight - appHeight) / 2

		centerX += mmath.Sum(widths[1:])
		positionXs[0] = centerX
		positionYs[0] = centerY
		for i := 0; i < viewerCount; i++ {
			centerX -= widths[i+1]
			positionXs[i+1] = centerX
			positionYs[i+1] = centerY
		}
		return
	}

	centerX := (screenWidth - appWidth) / 2
	centerY := (screenHeight - appHeight) / 2
	positionXs[0] = centerX
	positionYs[0] = centerY
	for i := 0; i < viewerCount; i++ {
		centerY += heights[i+1]
		positionXs[i+1] = centerX
		positionYs[i+1] = centerY
	}
	return
}
