package main

import (
	"embed"
	"fmt"
	"path/filepath"
	"runtime"

	"github.com/miu200521358/walk/pkg/walk"

	"github.com/miu200521358/mlib_go/pkg/mutils"
	"github.com/miu200521358/mlib_go/pkg/mwidget"
	"github.com/miu200521358/mlib_go/pkg/pmx"

)

func init() {
	runtime.LockOSThread()

	walk.AppendToWalkInit(func() {
		walk.MustRegisterWindowClass(mwidget.FilePickerClass)
	})
}

//go:embed resources/*
var resourceFiles embed.FS

func main() {
	mWindow, err := mwidget.NewMWindow(resourceFiles, true)
	if err != nil {
		walk.MsgBox(nil, "メインウィンドウ生成エラー", err.Error(), walk.MsgBoxIconError)
	}

	glWindow, err := mwidget.NewGlWindow("モデル描画", 512, 768, 0, resourceFiles, nil)
	if err != nil {
		walk.MsgBox(&mWindow.MainWindow, "モデル描画ウィンドウ生成エラー", err.Error(), walk.MsgBoxIconError)
	}
	mWindow.AddGlWindow(glWindow)

	pmxReadPicker, err := (mwidget.NewPmxReadFilePicker(
		mWindow,
		"PmxPath",
		"Pmxファイル",
		"Pmxファイルを選択してください",
		func(path string) {}))
	if err != nil {
		walk.MsgBox(&mWindow.MainWindow, "Pmxファイル選択エラー", err.Error(), walk.MsgBoxIconError)
	}

	pmxSavePicker, err := (mwidget.NewPmxSaveFilePicker(
		mWindow,
		"出力Pmxファイル",
		"出力Pmxファイルパスを入力もしくは選択してください",
		func(path string) {}))
	if err != nil {
		walk.MsgBox(&mWindow.MainWindow, "出力Pmxファイル選択エラー", err.Error(), walk.MsgBoxIconError)
	}

	execButton, err := walk.NewPushButton(&mWindow.MainWindow)
	if err != nil {
		walk.MsgBox(&mWindow.MainWindow, "モデル描画ボタン生成エラー", err.Error(), walk.MsgBoxIconError)
	}
	execButton.SetText("モデル描画")
	execButton.SetEnabled(false)

	subExecButton, err := walk.NewPushButton(&mWindow.MainWindow)
	if err != nil {
		walk.MsgBox(&mWindow.MainWindow, "サブウィンドウ描画ボタン生成エラー", err.Error(), walk.MsgBoxIconError)
	}
	subExecButton.SetText("サブウィンドウ描画")
	subExecButton.SetEnabled(false)

	console, err := (mwidget.NewConsoleView(mWindow))
	if err != nil {
		walk.MsgBox(&mWindow.MainWindow, "コンソール生成エラー", err.Error(), walk.MsgBoxIconError)
	}

	execButton.Clicked().Attach(func() {
		data, err := pmxReadPicker.GetData()
		if err != nil {
			walk.MsgBox(&mWindow.MainWindow, "Pmxファイル読み込みエラー", err.Error(), walk.MsgBoxIconError)
			return
		}
		model := data.(*pmx.PmxModel)

		console.AppendText(fmt.Sprintf("モデル名: %s", model.Name))
		console.AppendText(fmt.Sprintf("頂点数: %d", len(model.Vertices.Indexes)))
		console.AppendText(fmt.Sprintf("面数: %d", len(model.Faces.Indexes)))
		console.AppendText(fmt.Sprintf("材質数: %d", len(model.Materials.Indexes)))
		console.AppendText(fmt.Sprintf("ボーン数: %d", len(model.Bones.Indexes)))
		console.AppendText(fmt.Sprintf("表情数: %d", len(model.Morphs.Indexes)))

		mWindow.GetMainGlWindow().AddData(model)
		mWindow.GetMainGlWindow().Run()
	})

	subExecButton.Clicked().Attach(func() {
		subGlWindow, err := mwidget.NewGlWindow("サブウィンドウ", 300, 300, 1, resourceFiles, glWindow)
		if err != nil {
			walk.MsgBox(&mWindow.MainWindow, "サブウィンドウ生成エラー", err.Error(), walk.MsgBoxIconError)
			return
		}
		mWindow.AddGlWindow(subGlWindow)
	})

	pmxReadPicker.OnPathChanged = func(path string) {
		isExist, err := mutils.ExistsFile(path)
		if !isExist || err != nil {
			pmxSavePicker.PathLineEdit.SetText("")
			execButton.SetEnabled(false)
			subExecButton.SetEnabled(false)
			return
		}

		dir, file := filepath.Split(path)
		ext := filepath.Ext(file)
		outputPath := filepath.Join(dir, file[:len(file)-len(ext)]+"_out"+ext)
		pmxSavePicker.PathLineEdit.SetText(outputPath)
		execButton.SetEnabled(true)
		subExecButton.SetEnabled(true)
	}

	mWindow.Center()
	mWindow.Run()
}
