package main

import (
	"embed"
	"fmt"
	"path/filepath"
	"runtime"

	"github.com/miu200521358/walk/pkg/walk"

	"github.com/miu200521358/mlib_go/pkg/mutil"
	"github.com/miu200521358/mlib_go/pkg/mwidget"
	"github.com/miu200521358/mlib_go/pkg/pmx"
)

//go:embed resources/app_config.json
var appConfigFile embed.FS

func init() {
	runtime.LockOSThread()

	walk.AppendToWalkInit(func() {
		walk.MustRegisterWindowClass(mwidget.FilePickerClass)
	})
}

func main() {
	mWindow, err := mwidget.NewMWindow(appConfigFile)
	if err != nil {
		panic(err)
	}

	pmxReadPicker, err := (mwidget.NewPmxReadFilePicker(
		mWindow,
		"PmxPath",
		"Pmxファイル",
		"Pmxファイルを選択してください",
		func(path string) {}))
	if err != nil {
		panic(err)
	}

	pmxSavePicker, err := (mwidget.NewPmxSaveFilePicker(
		mWindow,
		"出力Pmxファイル",
		"出力Pmxファイルパスを入力もしくは選択してください",
		func(path string) {}))
	if err != nil {
		panic(err)
	}

	execButton, err := walk.NewPushButton(&mWindow.MainWindow)
	if err != nil {
		panic(err)
	}
	execButton.SetText("モデル描画")
	execButton.SetEnabled(false)

	console, err := (mwidget.NewConsoleView(mWindow))
	if err != nil {
		panic(err)
	}

	execButton.Clicked().Attach(func() {
		data, err := pmxReadPicker.GetData()
		if err != nil {
			panic(err)
		}
		model := data.(*pmx.PmxModel)

		console.AppendText(fmt.Sprintf("モデル名: %s", model.Name))
		console.AppendText(fmt.Sprintf("頂点数: %d", len(model.Vertices.Indexes)))
		console.AppendText(fmt.Sprintf("面数: %d", len(model.Faces.Indexes)))
		console.AppendText(fmt.Sprintf("材質数: %d", len(model.Materials.Indexes)))
		console.AppendText(fmt.Sprintf("ボーン数: %d", len(model.Bones.Indexes)))
		console.AppendText(fmt.Sprintf("表情数: %d", len(model.Morphs.Indexes)))

		glWindow, err := mwidget.NewGlWindow("モデル描画")
		if err != nil {
			panic(err)
		}
		glWindow.AddData()
	})

	pmxReadPicker.OnPathChanged = func(path string) {
		isExist, err := mutil.ExistsFile(path)
		if !isExist || err != nil {
			pmxSavePicker.PathLineEdit.SetText("")
			execButton.SetEnabled(false)
			return
		}

		dir, file := filepath.Split(path)
		ext := filepath.Ext(file)
		outputPath := filepath.Join(dir, file[:len(file)-len(ext)]+"_out"+ext)
		pmxSavePicker.PathLineEdit.SetText(outputPath)
		execButton.SetEnabled(true)
	}

	mWindow.Center()
	mWindow.Run()
}
