//go:build windows
// +build windows

package main

import (
	"embed"
	"log"
	"runtime"

	"github.com/miu200521358/walk/pkg/declarative"
	"github.com/miu200521358/walk/pkg/walk"

	"github.com/miu200521358/mlib_go/pkg/domain/pmx"
	"github.com/miu200521358/mlib_go/pkg/domain/vmd"
	"github.com/miu200521358/mlib_go/pkg/interface/app"
	"github.com/miu200521358/mlib_go/pkg/interface/controller"
	"github.com/miu200521358/mlib_go/pkg/interface/controller/widget"
	"github.com/miu200521358/mlib_go/pkg/interface/viewer"
	"github.com/miu200521358/mlib_go/pkg/mutils/mconfig"
	"github.com/miu200521358/mlib_go/pkg/mutils/mi18n"
	"github.com/miu200521358/mlib_go/pkg/mutils/mlog"
)

var env string

func init() {
	runtime.LockOSThread()

	// システム上の半分の論理プロセッサを使用させる
	runtime.GOMAXPROCS(int(runtime.NumCPU() / 2))

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

	appConfig := mconfig.LoadAppConfig(appFiles)
	appConfig.Env = env
	mi18n.Initialize(appI18nFiles)

	mApp := app.NewMApp(appConfig)
	mApp.RunViewerToControlChannel()
	mApp.RunControlToViewerChannel()

	go func() {
		// 操作ウィンドウは別スレッドで起動
		controlWindow := controller.NewControlWindow(appConfig, mApp.ControlToViewerChannel(), getMenuItems, 2)
		mApp.SetControlWindow(controlWindow)

		controlWindow.InitTabWidget()
		newFilePage(mApp, controlWindow)

		player := widget.NewMotionPlayer(controlWindow.MainWindow, controlWindow)
		player.SetOnTriggerPlay(func(v bool) {
			controlWindow.TabWidget.SetEnabled(!v)
		})
		controlWindow.SetPlayer(player)

		consoleView := widget.NewConsoleView(controlWindow.MainWindow, 256, 50)
		log.SetOutput(consoleView)

		mApp.RunController()
	}()

	mApp.AddViewWindow(viewer.NewViewWindow(
		mApp.ViewerCount(), appConfig, mApp, mApp.ViewerToControlChannel(), "No.1 ビューワー", nil))
	mApp.AddViewWindow(viewer.NewViewWindow(
		mApp.ViewerCount(), appConfig, mApp, mApp.ViewerToControlChannel(), "No.2 ビューワー",
		mApp.MainViewWindow().GetWindow()))

	mApp.Center()
	mApp.RunViewer()
}

func getMenuItems() []declarative.MenuItem {
	return []declarative.MenuItem{
		declarative.Action{
			Text:        mi18n.T("&サンプルメニュー"),
			OnTriggered: func() { mlog.IL(mi18n.T("サンプルヘルプ")) },
		},
	}
}

func newFilePage(mApp *app.MApp, controlWindow *controller.ControlWindow) *widget.MTabPage {
	tabPage := widget.NewMTabPage("ファイル")
	controlWindow.AddTabPage(tabPage.TabPage)

	tabPage.SetLayout(walk.NewVBoxLayout())

	loadedModels := make([][]*pmx.PmxModel, 2)
	loadedModels[0] = make([]*pmx.PmxModel, 2)
	loadedModels[1] = make([]*pmx.PmxModel, 1)

	mApp.SetFuncGetModels(
		func() [][]*pmx.PmxModel {
			return [][]*pmx.PmxModel{
				{loadedModels[0][0], loadedModels[0][1]},
				{loadedModels[1][0]},
			}
		},
	)

	loadedMotions := make([][]*vmd.VmdMotion, 2)
	loadedMotions[0] = make([]*vmd.VmdMotion, 2)
	loadedMotions[1] = make([]*vmd.VmdMotion, 1)

	mApp.SetFuncGetMotions(
		func() [][]*vmd.VmdMotion {
			return [][]*vmd.VmdMotion{
				{loadedMotions[0][0], loadedMotions[0][1]},
				{loadedMotions[1][0]},
			}
		},
	)

	pmx11ReadPicker := widget.NewPmxReadFilePicker(
		controlWindow,
		tabPage,
		"PmxPath",
		"No.1-1 Pmxファイル",
		"Pmxファイルを選択してください",
		"Pmxファイルの使い方")

	pmx11ReadPicker.SetOnPathChanged(func(path string) {
		if data, err := pmx11ReadPicker.Load(); err == nil {
			model := data.(*pmx.PmxModel)
			model.SetIndex(0)
			loadedModels[0][0] = model
			if loadedMotions[0][0] == nil {
				loadedMotions[0][0] = vmd.NewVmdMotion("")
			}
		} else {
			mlog.E(mi18n.T("読み込み失敗"), err)
		}
	})

	vmd11ReadPicker := widget.NewVmdVpdReadFilePicker(
		controlWindow,
		tabPage,
		"VmdPath",
		"No.1-1 Vmdファイル",
		"Vmdファイルを選択してください",
		"Vmdファイルの使い方")

	vmd11ReadPicker.SetOnPathChanged(func(path string) {
		if data, err := vmd11ReadPicker.Load(); err == nil {
			motion := data.(*vmd.VmdMotion)
			controlWindow.UpdateMaxFrame(motion.MaxFrame())
			loadedMotions[0][0] = motion
		} else {
			mlog.ET(mi18n.T("読み込み失敗"), err.Error())
		}
	})

	walk.NewVSeparator(tabPage)

	pmx12ReadPicker := widget.NewPmxReadFilePicker(
		controlWindow,
		tabPage,
		"PmxPath",
		"No.1-2 Pmxファイル",
		"Pmxファイルを選択してください",
		"Pmxファイルの使い方")

	pmx12ReadPicker.SetOnPathChanged(func(path string) {
		if data, err := pmx12ReadPicker.Load(); err == nil {
			model := data.(*pmx.PmxModel)
			model.SetIndex(1)
			loadedModels[0][1] = model
			if loadedMotions[0][1] == nil {
				loadedMotions[0][1] = vmd.NewVmdMotion("")
			}
		} else {
			mlog.E(mi18n.T("読み込み失敗"), err)
		}
	})

	vmd12ReadPicker := widget.NewVmdVpdReadFilePicker(
		controlWindow,
		tabPage,
		"VmdPath",
		"No.1-2 Vmdファイル",
		"Vmdファイルを選択してください",
		"Vmdファイルの使い方")

	vmd12ReadPicker.SetOnPathChanged(func(path string) {
		if data, err := vmd12ReadPicker.Load(); err == nil {
			motion := data.(*vmd.VmdMotion)
			controlWindow.UpdateMaxFrame(motion.MaxFrame())
			loadedMotions[0][1] = motion
		} else {
			mlog.E(mi18n.T("読み込み失敗"), err)
		}
	})

	walk.NewVSeparator(tabPage)

	pmx2ReadPicker := widget.NewPmxReadFilePicker(
		controlWindow,
		tabPage,
		"PmxPath",
		"No.2 Pmxファイル",
		"Pmxファイルを選択してください",
		"Pmxファイルの使い方")

	pmx2ReadPicker.SetOnPathChanged(func(path string) {
		if data, err := pmx2ReadPicker.Load(); err == nil {
			model := data.(*pmx.PmxModel)
			model.SetIndex(0)
			loadedModels[1][0] = model
			if loadedMotions[1][0] == nil {
				loadedMotions[1][0] = vmd.NewVmdMotion("")
			}
		} else {
			mlog.E(mi18n.T("読み込み失敗"), err)
		}
	})

	vmd2ReadPicker := widget.NewVmdVpdReadFilePicker(
		controlWindow,
		tabPage,
		"VmdPath",
		"No.2 Vmdファイル",
		"Vmdファイルを選択してください",
		"Vmdファイルの使い方")

	vmd2ReadPicker.SetOnPathChanged(func(path string) {
		if data, err := vmd2ReadPicker.Load(); err == nil {
			motion := data.(*vmd.VmdMotion)
			controlWindow.UpdateMaxFrame(motion.MaxFrame())
			loadedMotions[1][0] = motion
		} else {
			mlog.E(mi18n.T("読み込み失敗"), err)
		}
	})

	walk.NewVSpacer(tabPage)

	return tabPage
}
