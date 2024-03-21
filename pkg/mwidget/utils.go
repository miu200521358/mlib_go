package mwidget

import "github.com/miu200521358/walk/pkg/walk"

func CheckError(err error, w *MWindow, message string) {
	if err != nil {
		if w != nil {
			walk.MsgBox(w.MainWindow, message, err.Error(), walk.MsgBoxIconError)
			w.Close()
		} else {
			walk.MsgBox(nil, message, err.Error(), walk.MsgBoxIconError)
		}
	}
}
