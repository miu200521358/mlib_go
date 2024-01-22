package main

import (
	"embed"
	"fmt"
	"path/filepath"
	"runtime"

	"github.com/miu200521358/walk/pkg/declarative"
	"github.com/miu200521358/walk/pkg/walk"

	"github.com/miu200521358/mlib_go/pkg/utils/config"
	"github.com/miu200521358/mlib_go/pkg/widget/file_picker"
)

//go:embed resources/app_config.json
var appConfig embed.FS

func init() {
	runtime.LockOSThread()

	walk.AppendToWalkInit(func() {
		walk.MustRegisterWindowClass(file_picker.FilePickerClass)
	})
}

type Foo struct {
	Bar string
	Baz int
}

func main() {
	appConfig := config.ReadAppConfig(appConfig)

	var mw *walk.MainWindow

	if err := (declarative.MainWindow{
		AssignTo: &mw,
		Title:    fmt.Sprintf("%s %s", appConfig.AppName, appConfig.AppVersion),
		Size:     declarative.Size{Width: 1024, Height: 768},
		Layout:   declarative.VBox{Alignment: declarative.AlignHNearVNear},
	}).Create(); err != nil {
		panic(err)
	}

	if err := (file_picker.NewPmxReadFilePicker(
		mw,
		"PmxPath",
		"Pmxファイル",
		"Pmxファイルを選択してください",
		func(path string) {
			dir, file := filepath.Split(path)
			ext := filepath.Ext(file)
			outputPath := filepath.Join(dir, file[:len(file)-len(ext)]+"_out"+ext)
			println(outputPath)
			// pmxOutputFilePicker.PathEntry.SetText(outputPath)
		})); err != nil {
		panic(err)
	}

	mw.Run()
}
