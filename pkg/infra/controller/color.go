//go:build windows
// +build windows

// 指示: miu200521358
package controller

import "github.com/miu200521358/walk/pkg/walk"

var (
	// ColorWindowBackground はウィンドウ背景色。
	ColorWindowBackground = walk.RGB(160, 160, 160)
	// ColorTabBackground はタブ背景色。
	ColorTabBackground = walk.RGB(191, 205, 219)
	// ColorNavBackground はナビ背景色。
	ColorNavBackground = walk.RGB(237, 232, 201)
)
