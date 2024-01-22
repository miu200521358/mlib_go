package main

import (
	"embed"
	"fmt"
	"path/filepath"
	"runtime"

	"github.com/miu200521358/walk/pkg/walk"

	"github.com/miu200521358/mlib_go/pkg/pmx/pmx_model"
	"github.com/miu200521358/mlib_go/pkg/widget/console_view"
	"github.com/miu200521358/mlib_go/pkg/widget/file_picker"
	"github.com/miu200521358/mlib_go/pkg/widget/gl_window"
	"github.com/miu200521358/mlib_go/pkg/widget/m_window"

)

//go:embed resources/app_config.json
var appConfigFile embed.FS

func init() {
	runtime.LockOSThread()

	walk.AppendToWalkInit(func() {
		walk.MustRegisterWindowClass(file_picker.FilePickerClass)
	})
}

func main() {
	mWindow, err := m_window.NewMWindow(appConfigFile)
	if err != nil {
		panic(err)
	}

	pmxReadPicker, err := (file_picker.NewPmxReadFilePicker(
		mWindow,
		"PmxPath",
		"Pmxファイル",
		"Pmxファイルを選択してください",
		func(path string) {}))
	if err != nil {
		panic(err)
	}

	pmxSavePicker, err := (file_picker.NewPmxSaveFilePicker(
		mWindow,
		"出力Pmxファイル",
		"出力Pmxファイルパスを入力もしくは選択してください",
		func(path string) {}))
	if err != nil {
		panic(err)
	}

	pmxReadPicker.OnPathChanged = func(path string) {
		dir, file := filepath.Split(path)
		ext := filepath.Ext(file)
		outputPath := filepath.Join(dir, file[:len(file)-len(ext)]+"_out"+ext)
		pmxSavePicker.PathLineEdit.SetText(outputPath)
	}

	execButton, err := walk.NewPushButton(&mWindow.MainWindow)
	if err != nil {
		panic(err)
	}
	execButton.SetText("グリッド描画")

	console, err := (console_view.NewConsoleView(mWindow))
	if err != nil {
		panic(err)
	}

	execButton.Clicked().Attach(func() {
		if pmxReadPicker.PathLineEdit.Text() != "" {
			data, err := pmxReadPicker.GetData()
			if err != nil {
				panic(err)
			}
			model := data.(*pmx_model.PmxModel)

			console.AppendText(fmt.Sprintf("モデル名: %s", model.Name))
			console.AppendText(fmt.Sprintf("頂点数: %d", len(model.Vertices.Indexes)))
			console.AppendText(fmt.Sprintf("面数: %d", len(model.Faces.Indexes)))
			console.AppendText(fmt.Sprintf("材質数: %d", len(model.Materials.Indexes)))
			console.AppendText(fmt.Sprintf("ボーン数: %d", len(model.Bones.Indexes)))
			console.AppendText(fmt.Sprintf("表情数: %d", len(model.Morphs.Indexes)))
		}

		glWindow, err := gl_window.NewGlWindow("モデル描画")
		if err != nil {
			panic(err)
		}
		glWindow.AddData()

		// glWindow.Run()
		// glWindow.Draw()
	})

	mWindow.Center()
	mWindow.Run()
}
