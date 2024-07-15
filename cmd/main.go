//go:build windows
// +build windows

package main

import (
	"embed"
	"fmt"
	"log"
	"path/filepath"
	"runtime"

	"github.com/miu200521358/walk/pkg/declarative"
	"github.com/miu200521358/walk/pkg/walk"

	"github.com/miu200521358/mlib_go/pkg/domain/delta"
	"github.com/miu200521358/mlib_go/pkg/domain/mmath"
	"github.com/miu200521358/mlib_go/pkg/domain/pmx"
	"github.com/miu200521358/mlib_go/pkg/domain/vmd"
	"github.com/miu200521358/mlib_go/pkg/interface/widget"
	"github.com/miu200521358/mlib_go/pkg/interface/window"
	"github.com/miu200521358/mlib_go/pkg/mutils"
	"github.com/miu200521358/mlib_go/pkg/mutils/mconfig"
	"github.com/miu200521358/mlib_go/pkg/mutils/mi18n"
	"github.com/miu200521358/mlib_go/pkg/mutils/mlog"
)

var env string

func init() {
	runtime.LockOSThread()

	// システム上のすべての論理プロセッサを使用させる
	runtime.GOMAXPROCS(runtime.NumCPU())

	walk.AppendToWalkInit(func() {
		walk.MustRegisterWindowClass(widget.FilePickerClass)
		walk.MustRegisterWindowClass(widget.MotionPlayerClass)
		walk.MustRegisterWindowClass(widget.ConsoleViewClass)
	})
}

//go:embed app/*
var appFiles embed.FS

//go:embed i18n/*
var appI18nFiles embed.FS

func main() {
	// defer profile.Start(profile.MemProfileHeap, profile.ProfilePath(time.Now().Format("20060102_150405"))).Stop()
	// defer profile.Start(profile.MemProfile, profile.ProfilePath(time.Now().Format("20060102_150405"))).Stop()
	// defer profile.Start(profile.CPUProfile, profile.ProfilePath(fmt.Sprintf("cpu_%s", time.Now().Format("20060102_150405")))).Stop()
	// defer profile.Start(profile.CPUProfile, profile.ProfilePath(fmt.Sprintf("cpu_%s", time.Now().Format("20060102_150405")))).Stop()

	var mWindow *window.MWindow
	var err error

	appConfig := mconfig.LoadAppConfig(appFiles)
	appConfig.Env = env
	mi18n.Initialize(appI18nFiles)

	if appConfig.IsEnvProd() || appConfig.IsEnvDev() {
		defer widget.RecoverFromPanic(mWindow.MainWindow)
	}

	uiState := widget.NewUiState()

	iconImg, err := mconfig.LoadIconFile(appFiles)
	widget.CheckError(err, nil, mi18n.T("アイコン生成エラー"))

	glWindow, err := window.NewGlWindow(512, 768, 0, iconImg, appConfig, nil, uiState)

	go func() {
		mWindow, err = window.NewMWindow(512, 768, true, getMenuItems, iconImg, appConfig, uiState)
		widget.CheckError(err, nil, mi18n.T("メインウィンドウ生成エラー"))

		motionPlayer, worldPosFunc := NewFileTabPage(mWindow)

		widget.CheckError(err, mWindow.MainWindow, mi18n.T("ビューワーウィンドウ生成エラー"))
		mWindow.AddGlWindow(glWindow)
		mWindow.MotionPlayer = motionPlayer

		glWindow.SetWorldPosFunc(worldPosFunc)
		glWindow.SetTitle(fmt.Sprintf("%s %s", mWindow.Title(), mi18n.T("ビューワー")))

		// コンソールはタブ外に表示
		mWindow.ConsoleView, err = widget.NewConsoleView(mWindow, 256, 30)
		widget.CheckError(err, mWindow.MainWindow, mi18n.T("コンソール生成エラー"))
		log.SetOutput(mWindow.ConsoleView)

		mWindow.AsFormBase().Closing().Attach(func(canceled *bool, reason walk.CloseReason) {
			go func() {
				for _, glWindow := range mWindow.GlWindows {
					glWindow.IsClosedChannel <- true
				}
			}()
			mWindow.Close()
		})

		mWindow.Center()
		mWindow.Run()
	}()

	glWindow.RunChannel()
	glWindow.Run()

	defer glWindow.TriggerClose(glWindow.Window)
}

func getMenuItems() []declarative.MenuItem {
	return []declarative.MenuItem{
		declarative.Action{
			Text:        mi18n.T("&サンプルメニュー"),
			OnTriggered: func() { mlog.IL(mi18n.T("サンプルヘルプ")) },
		},
	}
}

func NewFileTabPage(mWindow *window.MWindow) (*widget.MotionPlayer, func(prevXprevYFrontPos, prevXprevYBackPos, prevXnowYFrontPos, prevXnowYBackPos,
	nowXprevYFrontPos, nowXprevYBackPos, nowXnowYFrontPos, nowXnowYBackPos *mmath.MVec3, vmdDeltas []*delta.VmdDeltas)) {
	page, _ := widget.NewMTabPage(mWindow.MainWindow, mWindow.TabWidget, mi18n.T("ファイル"))

	page.SetLayout(walk.NewVBoxLayout())

	pmxReadPicker, err := (widget.NewPmxReadFilePicker(
		mWindow.MainWindow,
		page,
		"PmxPath",
		mi18n.T("Pmxファイル"),
		mi18n.T("Pmxファイルを選択してください"),
		mi18n.T("Pmxファイルの使い方"),
		func(path string) {}))
	widget.CheckError(err, mWindow.MainWindow, mi18n.T("Pmxファイルピッカー生成エラー"))

	vmdReadPicker, err := (widget.NewVmdVpdReadFilePicker(
		mWindow.MainWindow,
		page,
		"VmdPath",
		mi18n.T("Vmdファイル"),
		mi18n.T("Vmdファイルを選択してください"),
		mi18n.T("Vmdファイルの使い方"),
		func(path string) {}))
	widget.CheckError(err, mWindow.MainWindow, mi18n.T("Vmdファイルピッカー生成エラー"))

	pmxSavePicker, err := (widget.NewPmxSaveFilePicker(
		mWindow.MainWindow,
		page,
		mi18n.T("出力Pmxファイル"),
		mi18n.T("出力Pmxファイルパスを入力もしくは選択してください"),
		mi18n.T("出力Pmxファイルの使い方"),
		func(path string) {}))
	widget.CheckError(err, mWindow.MainWindow, mi18n.T("出力Pmxファイルピッカー生成エラー"))

	_, err = walk.NewVSeparator(page)
	widget.CheckError(err, mWindow.MainWindow, mi18n.T("セパレータ生成エラー"))

	motionPlayer, err := widget.NewMotionPlayer(page, mWindow.MainWindow, mWindow.UiState)
	widget.CheckError(err, mWindow.MainWindow, mi18n.T("モーションプレイヤー生成エラー"))
	motionPlayer.SetEnabled(false)

	// fixViewWidget, err := widget.NewFixViewWidget(page, mWindow)
	// widget.CheckError(err, mWindow.MainWindow, mi18n.T("固定ビューウィジェット生成エラー"))
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

			motionPlayer.SetEnabled(true)
			motionPlayer.SetValue(0)

			go func() {
				for _, glWindow := range mWindow.GlWindows {
					glWindow.FrameChannel <- 0
					glWindow.IsPlayingChannel <- false
					glWindow.ReplaceModelChannel <- map[int]*pmx.PmxModel{0: model}
				}
			}()
		} else {
			go func() {
				for _, glWindow := range mWindow.GlWindows {
					glWindow.RemoveIndexChannel <- 0
				}
			}()
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
				motionPlayer.SetEnabled(true)
				motionPlayer.SetValue(0)

				go func() {
					for _, glWindow := range mWindow.GlWindows {
						glWindow.FrameChannel <- 0
						glWindow.IsPlayingChannel <- false
						glWindow.ReplaceMotionChannel <- map[int]*vmd.VmdMotion{0: motion}
					}
				}()
			} else {
				for _, glWindow := range mWindow.GlWindows {
					glWindow.RemoveIndexChannel <- 0
				}
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
		go func() {
			for _, glWindow := range mWindow.GlWindows {
				glWindow.IsPlayingChannel <- isPlaying
			}
		}()

		return nil
	}

	pmxReadPicker.PathLineEdit.SetFocus()

	// worldPosFunc := func(prevXprevYFrontPos, prevXprevYBackPos, prevXnowYFrontPos, prevXnowYBackPos,
	// 	nowXprevYFrontPos, nowXprevYBackPos, nowXnowYFrontPos, nowXnowYBackPos *mmath.MVec3, vmdDeltas []*delta.VmdDeltas) {
	// 	mlog.L()

	// 	if pmxReadPicker.Exists() {
	// 		model := pmxReadPicker.GetCache().(*pmx.PmxModel)
	// 		var nearestVertexIndexes [][]int
	// 		// 直近頂点を取得
	// 		if prevXnowYFrontPos == nil {
	// 			nearestVertexIndexes = vmdDeltas[0].Vertices.FindNearestVertexIndexes(prevXprevYFrontPos, nil)
	// 		} else {
	// 			nearestVertexIndexes = vmdDeltas[0].Vertices.FindVerticesInBox(prevXprevYFrontPos, prevXprevYBackPos,
	// 				prevXnowYFrontPos, prevXnowYBackPos, nowXprevYFrontPos, nowXprevYBackPos, nowXnowYFrontPos,
	// 				nowXnowYBackPos, nil)
	// 		}

	// 		if len(nearestVertexIndexes) > 0 {
	// 			for _, vertexIndex := range nearestVertexIndexes {
	// 				vertex := model.Vertices.Get(vertexIndex[0])
	// 				mlog.I("In Box Vertex: %d (元: %s)(変形: %s)",
	// 					vertex.Index, vertex.Position.String(),
	// 					vmdDeltas[0].Vertices.Get(vertex.Index).Position.String())
	// 			}
	// 			go func() {
	// 				glWindow.ReplaceModelSetChannel <- map[int]*widget.ModelSet{0: {NextSelectedVertexIndexes: nearestVertexIndexes[0]}}
	// 			}()
	// 		}

	// 		// // 大体近いボーンを取得
	// 		// nearestBoneIndexes := vmdDeltas[0].Bones.GetNearestBoneIndexes(worldPos)
	// 		// for _, boneIndex := range nearestBoneIndexes {
	// 		// 	bone := model.Bones.Get(boneIndex)
	// 		// 	mlog.I("Near Bone: %d, %s (元: %s)(変形: %s)",
	// 		// 		bone.Index, bone.Name, bone.Position.String(),
	// 		// 		vmdDeltas[0].Bones.Get(bone.Index).GlobalPosition().String())
	// 		// }
	// 	}
	// }

	return motionPlayer, nil
}
