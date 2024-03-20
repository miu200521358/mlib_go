package main

import (
	"embed"
	"fmt"
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

	vmdReadPicker, err := (mwidget.NewVmdReadFilePicker(
		mWindow,
		"VmdPath",
		"Vmdファイル",
		"Vmdファイルを選択してください",
		func(path string) {}))
	if err != nil {
		walk.MsgBox(&mWindow.MainWindow, "Vmdファイル選択エラー", err.Error(), walk.MsgBoxIconError)
	}

	pmxSavePicker, err := (mwidget.NewPmxSaveFilePicker(
		mWindow,
		"出力Pmxファイル",
		"出力Pmxファイルパスを入力もしくは選択してください",
		func(path string) {}))
	if err != nil {
		walk.MsgBox(&mWindow.MainWindow, "出力Pmxファイル選択エラー", err.Error(), walk.MsgBoxIconError)
	}

	motionFrameEdit, err := walk.NewNumberEdit(&mWindow.MainWindow)
	if err != nil {
		walk.MsgBox(&mWindow.MainWindow, "モーションフレーム生成エラー", err.Error(), walk.MsgBoxIconError)
	}
	motionFrameEdit.SetDecimals(0)
	motionFrameEdit.SetRange(0, 1)
	motionFrameEdit.SetValue(0)
	motionFrameEdit.SetEnabled(false)

	motionSlider, err := walk.NewSlider(&mWindow.MainWindow)
	if err != nil {
		walk.MsgBox(&mWindow.MainWindow, "モーションスライダー生成エラー", err.Error(), walk.MsgBoxIconError)
	}
	motionSlider.SetRange(0, 1)
	motionSlider.SetValue(0)
	motionSlider.SetEnabled(false)

	boneDebugCheckBox, err := walk.NewCheckBox(&mWindow.MainWindow)
	if err != nil {
		walk.MsgBox(&mWindow.MainWindow, "ボーンデバッグチェックボックス生成エラー", err.Error(), walk.MsgBoxIconError)
	}
	boneDebugCheckBox.SetText("ボーンデバッグ表示")
	boneDebugCheckBox.SetChecked(false)

	boneDebugCheckBox.Clicked().Attach(func() {
		glWindow.EnableBoneDebug = boneDebugCheckBox.Checked()
	})

	physicsDebugCheckBox, err := walk.NewCheckBox(&mWindow.MainWindow)
	if err != nil {
		walk.MsgBox(&mWindow.MainWindow, "物理デバッグチェックボックス生成エラー", err.Error(), walk.MsgBoxIconError)
	}
	physicsDebugCheckBox.SetText("物理デバッグ表示")
	physicsDebugCheckBox.SetChecked(false)

	physicsDebugCheckBox.Clicked().Attach(func() {
		glWindow.Physics.EnableDebug(physicsDebugCheckBox.Checked())
	})

	physicsCheckBox, err := walk.NewCheckBox(&mWindow.MainWindow)
	if err != nil {
		walk.MsgBox(&mWindow.MainWindow, "静的物理チェックボックス生成エラー", err.Error(), walk.MsgBoxIconError)
	}
	physicsCheckBox.SetText("物理ON/OFF")
	physicsCheckBox.SetChecked(true)
	physicsCheckBox.Clicked().Attach(func() {
		glWindow.EnablePhysics = physicsCheckBox.Checked()
	})

	execButton, err := walk.NewPushButton(&mWindow.MainWindow)
	if err != nil {
		walk.MsgBox(&mWindow.MainWindow, "モデル描画ボタン生成エラー", err.Error(), walk.MsgBoxIconError)
	}
	execButton.SetText("モデル描画")
	execButton.SetEnabled(false)

	paused := false
	pauseButton, err := walk.NewPushButton(&mWindow.MainWindow)
	if err != nil {
		walk.MsgBox(&mWindow.MainWindow, "一時停止ボタン生成エラー", err.Error(), walk.MsgBoxIconError)
	}
	pauseButton.SetText("一時停止")
	pauseButton.SetEnabled(false)
	pauseButton.Clicked().Attach(func() {
		paused = !paused
		glWindow.Pause(paused)
		if paused {
			pauseButton.SetText("再開")
		} else {
			pauseButton.SetText("一時停止")
		}
	})

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
		model.SetUp()

		println(fmt.Sprintf("頂点数: %d", len(model.Vertices.Indexes)))
		println(fmt.Sprintf("面数: %d", len(model.Faces.Indexes)))

		console.AppendText(fmt.Sprintf("モデル名: %s", model.Name))
		console.AppendText(fmt.Sprintf("頂点数: %d", len(model.Vertices.Indexes)))
		console.AppendText(fmt.Sprintf("面数: %d", len(model.Faces.Indexes)))
		console.AppendText(fmt.Sprintf("材質数: %d", len(model.Materials.GetIndexes())))
		console.AppendText(fmt.Sprintf("ボーン数: %d", len(model.Bones.GetIndexes())))
		console.AppendText(fmt.Sprintf("表情数: %d", len(model.Morphs.GetIndexes())))

		var motion *vmd.VmdMotion
		if vmdReadPicker.Exists() {
			motionData, err := vmdReadPicker.GetData()
			if err != nil {
				walk.MsgBox(&mWindow.MainWindow, "Vmdファイル読み込みエラー", err.Error(), walk.MsgBoxIconError)
				return
			}
			motion = motionData.(*vmd.VmdMotion)
		}

		if motion == nil {
			motion = vmd.NewVmdMotion("")
		}

		motionFrameEdit.SetRange(0, float64(motion.GetMaxFrame()+1))
		motionFrameEdit.SetValue(0)

		motionFrameEdit.ValueChanged().Attach(func() {
			if paused {
				mWindow.GetMainGlWindow().SetValue(
					float32(motionFrameEdit.Value()) / mWindow.GetMainGlWindow().Physics.Fps)
				motionSlider.SetValue(int(motionFrameEdit.Value()))
			}
		})

		motionSlider.SetRange(0, int(motion.GetMaxFrame()+1))
		motionSlider.SetValue(0)
		motionSlider.ValueChanged().Attach(func() {
			if paused {
				mWindow.GetMainGlWindow().SetValue(
					float32(motionSlider.Value()) / mWindow.GetMainGlWindow().Physics.Fps)
				motionFrameEdit.SetValue(float64(motionSlider.Value()))
			}
		})

		mWindow.GetMainGlWindow().ClearData()
		mWindow.GetMainGlWindow().AddData(model, motion)
		mWindow.GetMainGlWindow().Run(motionFrameEdit, motionSlider)
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
			physicsDebugCheckBox.SetEnabled(false)
			execButton.SetEnabled(false)
			subExecButton.SetEnabled(false)
			return
		}

		dir, file := filepath.Split(path)
		ext := filepath.Ext(file)
		outputPath := filepath.Join(dir, file[:len(file)-len(ext)]+"_out"+ext)
		pmxSavePicker.PathLineEdit.SetText(outputPath)
		physicsDebugCheckBox.SetEnabled(true)
		execButton.SetEnabled(true)
		subExecButton.SetEnabled(true)
		motionFrameEdit.SetEnabled(true)
		motionSlider.SetEnabled(true)
		pauseButton.SetEnabled(true)
	}

	mWindow.Center()
	mWindow.Run()
}
