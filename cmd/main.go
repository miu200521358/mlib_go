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

	var mWindow *mwidget.MWindow
	var err error

	appConfig := mconfig.LoadAppConfig(appFiles)
	appConfig.Env = env
	mi18n.Initialize(appI18nFiles)

	if appConfig.IsEnvProd() || appConfig.IsEnvDev() {
		defer mwidget.RecoverFromPanic(mWindow)
	}

	iconImg, err := mconfig.LoadIconFile(appFiles)
	mwidget.CheckError(err, nil, mi18n.T("アイコン生成エラー"))

	glWindow, err := mwidget.NewGlWindow(512, 768, 0, iconImg, appConfig, nil, nil)

	go func() {
		mWindow, err = mwidget.NewMWindow(512, 768, getMenuItems, iconImg, appConfig, true)
		mwidget.CheckError(err, nil, mi18n.T("メインウィンドウ生成エラー"))

		motionPlayer, _, funcWorldPos := NewFileTabPage(mWindow)

		mwidget.CheckError(err, mWindow, mi18n.T("ビューワーウィンドウ生成エラー"))
		mWindow.AddGlWindow(glWindow)
		glWindow.SetFuncWorldPos(funcWorldPos)
		glWindow.SetMotionPlayer(motionPlayer)
		glWindow.SetTitle(fmt.Sprintf("%s %s", mWindow.Title(), mi18n.T("ビューワー")))

		// コンソールはタブ外に表示
		mWindow.ConsoleView, err = mwidget.NewConsoleView(mWindow, 256, 30)
		mwidget.CheckError(err, mWindow, mi18n.T("コンソール生成エラー"))
		log.SetOutput(mWindow.ConsoleView)

		mWindow.AsFormBase().Closing().Attach(func(canceled *bool, reason walk.CloseReason) {
			go func() {
				mWindow.GetMainGlWindow().IsClosedChannel <- true
			}()
			mWindow.Close()
		})

		mWindow.Center()
		mWindow.Run()
	}()

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

func NewFileTabPage(mWindow *mwidget.MWindow) (*mwidget.MotionPlayer, *mwidget.FixViewWidget, func(worldPos *mmath.MVec3, vmdDeltas []*vmd.VmdDeltas, viewMat *mmath.MMat4)) {
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

	motionPlayer, err := mwidget.NewMotionPlayer(page, mWindow)
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
			motionPlayer.SetValue(0)

			go func() {
				mWindow.GetMainGlWindow().FrameChannel <- 0
				mWindow.GetMainGlWindow().IsPlayingChannel <- false
				mWindow.GetMainGlWindow().ReplaceModelSetChannel <- map[int]*mwidget.ModelSet{0: {NextModel: model, NextMotion: motion}}
			}()
		} else {
			go func() {
				mWindow.GetMainGlWindow().RemoveModelSetIndexChannel <- 0
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
					mWindow.GetMainGlWindow().FrameChannel <- 0
					mWindow.GetMainGlWindow().IsPlayingChannel <- false
					mWindow.GetMainGlWindow().ReplaceModelSetChannel <- map[int]*mwidget.ModelSet{0: {NextMotion: motion}}
				}()
			} else {
				go func() {
					mWindow.GetMainGlWindow().RemoveModelSetIndexChannel <- 0
				}()
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
			mWindow.GetMainGlWindow().IsPlayingChannel <- isPlaying
		}()

		return nil
	}

	pmxReadPicker.PathLineEdit.SetFocus()

	funcWorldPos := func(worldPos *mmath.MVec3, vmdDeltas []*vmd.VmdDeltas, viewMat *mmath.MMat4) {
		mlog.L()
		mlog.I("WorldPosResult: x=%.7f, y=%.7f, z=%.7f", worldPos.GetX(), worldPos.GetY(), worldPos.GetZ())

		if pmxReadPicker.Exists() {
			model := pmxReadPicker.GetCache().(*pmx.PmxModel)
			// 直近頂点を取得
			nearestVertexIndexes := vmdDeltas[0].Vertices.GetNearestVertexIndexes(worldPos)
			for _, vertexIndex := range nearestVertexIndexes {
				vertex := model.Vertices.Get(vertexIndex)
				mlog.I("Near Vertex: %d (元: %s)(変形: %s)",
					vertex.Index, vertex.Position.String(),
					vmdDeltas[0].Vertices.Get(vertex.Index).Position.String())
			}
			go func() {
				mWindow.GetMainGlWindow().ReplaceModelSetChannel <- map[int]*mwidget.ModelSet{0: {NextSelectedVertexIndexes: nearestVertexIndexes}}
			}()

			// 大体近いボーンを取得
			nearestBoneIndexes := vmdDeltas[0].Bones.GetNearestBoneIndexes(worldPos)
			for _, boneIndex := range nearestBoneIndexes {
				bone := model.Bones.Get(boneIndex)
				mlog.I("Near Bone: %d, %s (元: %s)(変形: %s)",
					bone.Index, bone.Name, bone.Position.String(),
					vmdDeltas[0].Bones.Get(bone.Index).GlobalPosition().String())
			}
		}
	}

	return motionPlayer, nil, funcWorldPos
}
