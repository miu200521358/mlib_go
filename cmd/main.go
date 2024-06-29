//go:build windows
// +build windows

package main

import (
	"embed"
	"fmt"
	"log"
	"math"
	"path/filepath"
	"runtime"

	"github.com/miu200521358/walk/pkg/declarative"
	"github.com/miu200521358/walk/pkg/walk"

	"github.com/miu200521358/mlib_go/pkg/mmath"
	"github.com/miu200521358/mlib_go/pkg/mutils"
	"github.com/miu200521358/mlib_go/pkg/mutils/mconfig"
	"github.com/miu200521358/mlib_go/pkg/mutils/mi18n"
	"github.com/miu200521358/mlib_go/pkg/mutils/mlog"
	"github.com/miu200521358/mlib_go/pkg/mwidget"
	"github.com/miu200521358/mlib_go/pkg/pmx"
	"github.com/miu200521358/mlib_go/pkg/vmd"
)

var env string

func init() {
	runtime.LockOSThread()

	walk.AppendToWalkInit(func() {
		walk.MustRegisterWindowClass(mwidget.FilePickerClass)
		walk.MustRegisterWindowClass(mwidget.MotionPlayerClass)
		walk.MustRegisterWindowClass(mwidget.ConsoleViewClass)
		walk.MustRegisterWindowClass(mwidget.FixViewWidgetClass)
	})

	// runtime.MemProfileRate = 1
}

//go:embed resources/*
var resourceFiles embed.FS

func main() {
	// defer profile.Start(profile.MemProfileHeap, profile.ProfilePath(time.Now().Format("20060102_150405"))).Stop()
	// defer profile.Start(profile.MemProfile, profile.ProfilePath(time.Now().Format("20060102_150405"))).Stop()
	// defer profile.Start(profile.CPUProfile,
	// 	profile.ProfilePath(fmt.Sprintf("cpu_%s", time.Now().Format("20060102_150405")))).Stop()
	// if err := agent.Listen(agent.Options{}); err != nil {
	// 	log.Fatal(err)
	// }

	// f, err := os.Create("trace.out")
	// if err != nil {
	// 	log.Fatalf("failed to create trace output file: %v", err)
	// }
	// defer f.Close()

	// if err := trace.Start(f); err != nil {
	// 	log.Fatalf("failed to start trace: %v", err)
	// }
	// defer trace.Stop()

	var mWindow *mwidget.MWindow
	var err error

	appConfig := mconfig.LoadAppConfig(resourceFiles)
	appConfig.Env = env

	if appConfig.IsEnvProd() || appConfig.IsEnvDev() {
		defer mwidget.RecoverFromPanic(mWindow)
	}

	mWindow, err = mwidget.NewMWindow(resourceFiles, appConfig, true, 512, 768, getMenuItems)
	mwidget.CheckError(err, nil, mi18n.T("メインウィンドウ生成エラー"))

	motionPlayer, fixViewWidget, funcWorldPos := NewFileTabPage(mWindow)

	glWindow, err := mwidget.NewGlWindow(fmt.Sprintf("%s %s", mWindow.Title(), mi18n.T("ビューワー")),
		512, 768, 0, resourceFiles, nil, motionPlayer, fixViewWidget, funcWorldPos)
	mwidget.CheckError(err, mWindow, mi18n.T("ビューワーウィンドウ生成エラー"))
	mWindow.AddGlWindow(glWindow)
	defer glWindow.Close(glWindow.Window)

	// コンソールはタブ外に表示
	mWindow.ConsoleView, err = mwidget.NewConsoleView(mWindow, 256, 30)
	mwidget.CheckError(err, mWindow, mi18n.T("コンソール生成エラー"))
	log.SetOutput(mWindow.ConsoleView)

	mWindow.Center()
	mWindow.Run()
}

func getMenuItems() []declarative.MenuItem {
	return []declarative.MenuItem{
		declarative.Action{
			Text:        mi18n.T("&サンプルメニュー"),
			OnTriggered: func() { mlog.IL(mi18n.T("サンプルヘルプ")) },
		},
	}
}

func NewFileTabPage(mWindow *mwidget.MWindow) (*mwidget.MotionPlayer, *mwidget.FixViewWidget, func(worldPos *mmath.MVec3, viewMat *mmath.MMat4)) {
	page, _ := mwidget.NewMTabPage(mWindow, mWindow.TabWidget, mi18n.T("ファイル"))

	page.SetLayout(walk.NewVBoxLayout())

	pmxReadPicker, err := (mwidget.NewPmxReadFilePicker(
		mWindow,
		page,
		"PmxPath",
		mi18n.T("Pmxファイル"),
		mi18n.T("Pmxファイルを選択してください"),
		mi18n.T("Pmxファイルの使い方"),
		func(path string) {}))
	mwidget.CheckError(err, mWindow, mi18n.T("Pmxファイルピッカー生成エラー"))

	vmdReadPicker, err := (mwidget.NewVmdVpdReadFilePicker(
		mWindow,
		page,
		"VmdPath",
		mi18n.T("Vmdファイル"),
		mi18n.T("Vmdファイルを選択してください"),
		mi18n.T("Vmdファイルの使い方"),
		func(path string) {}))
	mwidget.CheckError(err, mWindow, mi18n.T("Vmdファイルピッカー生成エラー"))

	pmxSavePicker, err := (mwidget.NewPmxSaveFilePicker(
		mWindow,
		page,
		mi18n.T("出力Pmxファイル"),
		mi18n.T("出力Pmxファイルパスを入力もしくは選択してください"),
		mi18n.T("出力Pmxファイルの使い方"),
		func(path string) {}))
	mwidget.CheckError(err, mWindow, mi18n.T("出力Pmxファイルピッカー生成エラー"))

	_, err = walk.NewVSeparator(page)
	mwidget.CheckError(err, mWindow, mi18n.T("セパレータ生成エラー"))

	motionPlayer, err := mwidget.NewMotionPlayer(page, mWindow, resourceFiles)
	mwidget.CheckError(err, mWindow, mi18n.T("モーションプレイヤー生成エラー"))
	motionPlayer.SetEnabled(false)

	// fixViewWidget, err := mwidget.NewFixViewWidget(page, mWindow)
	// mwidget.CheckError(err, mWindow, mi18n.T("固定ビューウィジェット生成エラー"))
	// fixViewWidget.SetEnabled(false)

	var onFilePathChanged = func() {
		if motionPlayer.Playing() {
			motionPlayer.Play(false)
		}
		enabled := pmxReadPicker.Exists() && vmdReadPicker.ExistsOrEmpty()
		motionPlayer.SetEnabled(enabled)
		// fixViewWidget.SetEnabled(enabled)
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
			data, err := pmxReadPicker.GetData()
			if err != nil {
				mlog.E(mi18n.T("Pmxファイル読み込みエラー"), err.Error())
				return
			}
			model := data.(*pmx.PmxModel)
			var motion *vmd.VmdMotion
			if vmdReadPicker.IsCached() {
				motion = vmdReadPicker.GetCache().(*vmd.VmdMotion)
			} else {
				motion = vmd.NewVmdMotion("")
			}

			motionPlayer.SetEnabled(true)
			// fixViewWidget.SetEnabled(true)
			mWindow.GetMainGlWindow().SetFrame(0)
			motionPlayer.SetValue(0)
			mWindow.GetMainGlWindow().Play(false)
			mWindow.GetMainGlWindow().ClearData()
			mWindow.GetMainGlWindow().AddData(model, motion)
			mWindow.GetMainGlWindow().Run()
		}

		onFilePathChanged()
	}

	vmdReadPicker.OnPathChanged = func(path string) {
		if vmdReadPicker.Exists() {
			motionData, err := vmdReadPicker.GetData()
			if err != nil {
				mlog.E(mi18n.T("Vmdファイル読み込みエラー"), err.Error())
				return
			}
			motion := motionData.(*vmd.VmdMotion)

			motionPlayer.SetRange(0, motion.GetMaxFrame()+1)
			motionPlayer.SetValue(0)

			if pmxReadPicker.Exists() {
				model := pmxReadPicker.GetCache().(*pmx.PmxModel)

				motionPlayer.SetEnabled(true)
				// fixViewWidget.SetEnabled(true)
				mWindow.GetMainGlWindow().SetFrame(0)
				motionPlayer.SetValue(0)
				mWindow.GetMainGlWindow().Play(false)
				mWindow.GetMainGlWindow().ClearData()
				mWindow.GetMainGlWindow().AddData(model, motion)
				mWindow.GetMainGlWindow().Run()
			}
		}

		onFilePathChanged()
	}

	motionPlayer.OnPlay = func(isPlaying bool) error {
		if !isPlaying {
			pmxReadPicker.SetEnabled(true)
			vmdReadPicker.SetEnabled(true)
			pmxSavePicker.SetEnabled(true)
		} else {
			pmxReadPicker.SetEnabled(false)
			vmdReadPicker.SetEnabled(false)
			pmxSavePicker.SetEnabled(false)
		}

		motionPlayer.PlayButton.SetEnabled(true)
		mWindow.GetMainGlWindow().Play(isPlaying)

		return nil
	}

	pmxReadPicker.PathLineEdit.SetFocus()

	funcWorldPos := func(worldPos *mmath.MVec3, viewMat *mmath.MMat4) {
		if pmxReadPicker.Exists() {
			model := pmxReadPicker.GetCache().(*pmx.PmxModel)
			distances := mmath.DistanceLineToPoints(
				worldPos, model.Vertices.Positions)
			sortedVertexIndexes := mutils.ArgSort(mutils.Float64Slice(distances))
			nearestVertexIndex := sortedVertexIndexes[0]
			// 大体近くて最も手前の頂点を選択
			z := math.MaxFloat64
			for i := range min(3, len(sortedVertexIndexes)) {
				sortedVertexIndex := sortedVertexIndexes[i]
				vPos := model.Vertices.Positions[sortedVertexIndex]
				// カメラの向きに応じた奥行きを取得
				cameraVPos := viewMat.MulVec3(vPos)
				if cameraVPos.GetZ() < z {
					nearestVertexIndex = sortedVertexIndex
					z = cameraVPos.GetZ()
				}
				mlog.D("near vertex index: %d worldPos: %s, cameraPos: %s",
					sortedVertexIndex, vPos.String(), cameraVPos.String())
			}
			mlog.D("Nearest Vertex Index: %d (%s)", nearestVertexIndex,
				model.Vertices.Get(nearestVertexIndex).Position.String())
		}
	}

	return motionPlayer, nil, funcWorldPos
}
