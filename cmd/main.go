package main

import (
	"embed"
	"log"
	"path/filepath"

	"github.com/miu200521358/walk/pkg/walk"

	"github.com/miu200521358/mlib_go/pkg/mutils"
	"github.com/miu200521358/mlib_go/pkg/mwidget"
	"github.com/miu200521358/mlib_go/pkg/pmx"
	"github.com/miu200521358/mlib_go/pkg/vmd"

)

func init() {
	// runtime.LockOSThread()

	walk.AppendToWalkInit(func() {
		walk.MustRegisterWindowClass(mwidget.FilePickerClass)
		walk.MustRegisterWindowClass(mwidget.MotionPlayerClass)
		walk.MustRegisterWindowClass(mwidget.ConsoleViewClass)
	})
}

//go:embed resources/*
var resourceFiles embed.FS

func main() {
	// {
	// 	// CPUプロファイル用のファイルを作成
	// 	f, err := os.Create("cpu.pprof")
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// 	defer f.Close()

	// 	// CPUプロファイリングを開始
	// 	if err := pprof.StartCPUProfile(f); err != nil {
	// 		panic(err)
	// 	}
	// 	defer pprof.StopCPUProfile()
	// }

	// {
	// 	// go tool pprof cmd\main.go cmd\memory.pprof
	// 	// メモリプロファイル用のファイルを作成
	// 	f, err := os.Create("memory.pprof")
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// 	defer f.Close()
	// 	runtime.GC()

	// 	// ヒーププロファイリングを開始
	// 	if err := pprof.WriteHeapProfile(f); err != nil {
	// 		panic(err)
	// 	}
	// }

	mWindow, err := mwidget.NewMWindow(resourceFiles, true, 512, 768)
	mwidget.CheckError(err, nil, "メインウィンドウ生成エラー")

	glWindow, err := mwidget.NewGlWindow("モデル描画", 512, 768, 0, resourceFiles, nil)
	mwidget.CheckError(err, mWindow, "モデル描画ウィンドウ生成エラー")
	mWindow.AddGlWindow(glWindow)

	NewFileTabPage(mWindow)

	// コンソールはタブ外に表示
	mWindow.ConsoleView, err = mwidget.NewConsoleView(mWindow)
	mwidget.CheckError(err, mWindow, "コンソール生成エラー")
	log.SetOutput(mWindow.ConsoleView)

	mWindow.Center()
	mWindow.Run()
}

func NewFileTabPage(mWindow *mwidget.MWindow) *mwidget.MTabPage {
	page := mwidget.NewMTabPage(mWindow, mWindow.TabWidget, "ファイル")

	mainLayout := walk.NewVBoxLayout()
	page.SetLayout(mainLayout)

	pmxReadPicker, err := (mwidget.NewPmxReadFilePicker(
		page,
		"PmxPath",
		"Pmxファイル",
		"Pmxファイルを選択してください",
		func(path string) {}))
	mwidget.CheckError(err, mWindow, "Pmxファイルピッカー生成エラー")

	vmdReadPicker, err := (mwidget.NewVmdReadFilePicker(
		page,
		"VmdPath",
		"Vmdファイル",
		"Vmdファイルを選択してください",
		func(path string) {}))
	mwidget.CheckError(err, mWindow, "Vmdファイルピッカー生成エラー")

	pmxSavePicker, err := (mwidget.NewPmxSaveFilePicker(
		page,
		"出力Pmxファイル",
		"出力Pmxファイルパスを入力もしくは選択してください",
		func(path string) {}))
	mwidget.CheckError(err, mWindow, "出力Pmxファイルピッカー生成エラー")

	_, err = walk.NewVSeparator(page)
	mwidget.CheckError(err, mWindow, "セパレータ生成エラー")

	motionPlayer, err := mwidget.NewMotionPlayer(page, mWindow, resourceFiles)
	mwidget.CheckError(err, mWindow, "モーションプレイヤー生成エラー")
	motionPlayer.SetEnabled(false)
	motionPlayer.PlayButton.SetEnabled(false)

	var onFilePathChanged = func() {
		mWindow.GetMainGlWindow().Play(false)
		motionPlayer.Play(false)
		if pmxReadPicker.Exists() && vmdReadPicker.ExistsOrEmpty() {
			motionPlayer.PlayButton.SetEnabled(true)
		} else {
			motionPlayer.PlayButton.SetEnabled(false)
		}
	}

	pmxReadPicker.OnPathChanged = func(path string) {
		isExist, err := mutils.ExistsFile(path)
		if !isExist || err != nil {
			pmxSavePicker.PathLineEdit.SetText("")
			return
		}

		dir, file := filepath.Split(path)
		ext := filepath.Ext(file)
		outputPath := filepath.Join(dir, file[:len(file)-len(ext)]+"_out"+ext)
		pmxSavePicker.PathLineEdit.SetText(outputPath)

		if pmxReadPicker.Exists() {
			motionPlayer.PlayButton.SetEnabled(true)
		}

		onFilePathChanged()
	}

	vmdReadPicker.OnPathChanged = func(path string) {
		if pmxReadPicker.Exists() {
			motionPlayer.PlayButton.SetEnabled(true)
		}

		onFilePathChanged()
	}

	motionPlayer.OnPlay = func(isPlaying bool) {
		if !pmxReadPicker.Exists() {
			return
		}

		if isPlaying {
			var model *pmx.PmxModel
			modelCached := pmxReadPicker.IsCached()
			motionCached := vmdReadPicker.IsCached()

			if !modelCached {
				data, err := pmxReadPicker.GetData()
				mwidget.CheckError(err, mWindow, "Pmxファイル読み込みエラー")
				model = data.(*pmx.PmxModel)
				model.SetUp()

				log.Printf("モデル名: %s", model.Name)
			} else {
				model = pmxReadPicker.GetCache().(*pmx.PmxModel)
			}

			var motion *vmd.VmdMotion
			if !motionCached {
				if vmdReadPicker.Exists() {
					motionData, err := vmdReadPicker.GetData()
					if err != nil {
						walk.MsgBox(mWindow.MainWindow, "Vmdファイル読み込みエラー", err.Error(), walk.MsgBoxIconError)
						return
					}
					motion = motionData.(*vmd.VmdMotion)
				}

				if motion == nil {
					motion = vmd.NewVmdMotion("")
				}
			} else {
				motion = vmdReadPicker.GetCache().(*vmd.VmdMotion)
			}

			if modelCached && motionCached {
				mWindow.GetMainGlWindow().Play(true)
			} else {
				motionPlayer.SetEnabled(false)
				motionPlayer.SetRange(0, float64(motion.GetMaxFrame()+1))
				motionPlayer.SetValue(0)

				mWindow.GetMainGlWindow().SetFrame(0)
				mWindow.GetMainGlWindow().Play(true)
				mWindow.GetMainGlWindow().ClearData()
				mWindow.GetMainGlWindow().AddData(model, motion)
				mWindow.GetMainGlWindow().Run(motionPlayer)
			}
		} else {
			mWindow.GetMainGlWindow().Play(false)
		}
	}

	return page
}
