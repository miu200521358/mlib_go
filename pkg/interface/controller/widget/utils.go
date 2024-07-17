//go:build windows
// +build windows

package widget

import (
	"github.com/miu200521358/walk/pkg/walk"
)

func CheckError(err error, w *walk.MainWindow, message string) {
	if err != nil {
		if w != nil {
			walk.MsgBox(w, message, err.Error(), walk.MsgBoxIconError)
			w.Close()
		} else {
			walk.MsgBox(nil, message, err.Error(), walk.MsgBoxIconError)
		}
	}
}

var MARGIN_ZERO = walk.Margins{HNear: 0, VNear: 0, HFar: 0, VFar: 0}
var MARGIN_SMALL = walk.Margins{HNear: 3, VNear: 3, HFar: 3, VFar: 3}
var MARGIN_MEDIUM = walk.Margins{HNear: 6, VNear: 6, HFar: 6, VFar: 6}
