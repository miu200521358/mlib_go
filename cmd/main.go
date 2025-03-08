//go:build windows
// +build windows

package main

import (
	"embed"
	"fmt"
	"runtime"

	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/miu200521358/walk/pkg/declarative"
	"github.com/miu200521358/walk/pkg/walk"

	"github.com/miu200521358/mlib_go/pkg/config/mconfig"
	"github.com/miu200521358/mlib_go/pkg/config/mi18n"
	"github.com/miu200521358/mlib_go/pkg/config/mlog"
	"github.com/miu200521358/mlib_go/pkg/domain/pmx"
	"github.com/miu200521358/mlib_go/pkg/domain/vmd"
	"github.com/miu200521358/mlib_go/pkg/infrastructure/repository"
	"github.com/miu200521358/mlib_go/pkg/interface/app"
	"github.com/miu200521358/mlib_go/pkg/interface/controller"
	"github.com/miu200521358/mlib_go/pkg/interface/controller/widget"
	"github.com/miu200521358/mlib_go/pkg/interface/state"
	"github.com/miu200521358/mlib_go/pkg/interface/viewer"
)

var env string

func init() {
	runtime.LockOSThread()

	// システム上の半分の論理プロセッサを使用させる
	runtime.GOMAXPROCS(max(1, int(runtime.NumCPU()/4)))

	walk.AppendToWalkInit(func() {
		walk.MustRegisterWindowClass(controller.ConsoleViewClass)
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

	viewerCount := 2

	appConfig := mconfig.LoadAppConfig(appFiles)
	appConfig.Env = env
	mi18n.Initialize(appI18nFiles)
	shared := state.NewSharedState(viewerCount)

	widths, heights, positionXs, positionYs := app.GetCenterSizeAndWidth(appConfig, viewerCount)

	var controlWindow *controller.ControlWindow
	viewerWindowList := viewer.NewViewerList(shared, appConfig)
	var err error

	go func() {
		// 操作ウィンドウは別スレッドで起動
		defer app.SafeExecute(appConfig.IsSetEnv(), func() {
			widgets := &controller.MWidgets{}

			controlWindow, err = controller.NewControlWindow(shared, appConfig,
				newMenuItems(), newTabPages(widgets), widgets.EnabledInPlaying,
				widths[0], heights[0], positionXs[0], positionYs[0])
			if err != nil {
				app.ShowErrorDialog(appConfig.IsSetEnv(), err)
				return
			}

			widgets.SetWindow(controlWindow)
			controlWindow.Run()
		})
	}()

	// GL初期化
	if err := glfw.Init(); err != nil {
		app.ShowErrorDialog(appConfig.IsSetEnv(), fmt.Errorf("failed to initialize GLFW: %v", err))
		return
	}

	// 描画ウィンドウはメインスレッドで起動
	defer app.SafeExecute(appConfig.IsSetEnv(), func() {
		for n := range viewerCount {
			nIdx := n + 1
			if err := viewerWindowList.Add(fmt.Sprintf("Viewer%d", nIdx),
				widths[nIdx], heights[nIdx], positionXs[nIdx], positionYs[nIdx]); err != nil {
				app.ShowErrorDialog(appConfig.IsSetEnv(), err)
				return
			}
		}

		viewerWindowList.Run()
	})
}

func newMenuItems() []declarative.MenuItem {
	return []declarative.MenuItem{
		declarative.Action{
			Text:        mi18n.T("&サンプルメニュー"),
			OnTriggered: func() { mlog.IL("%s", mi18n.T("サンプルヘルプ")) },
		},
	}
}

func newTabPages(mWidgets *controller.MWidgets) []declarative.TabPage {
	var fileTab *walk.TabPage

	player := widget.NewMotionPlayer()

	pmxLoad11Picker := widget.NewPmxLoadFilePicker(
		"pmx",
		"モデルファイル1-1",
		"モデルファイルを選択してください",
		func(cw *controller.ControlWindow, rep repository.IRepository, path string) {
			if data, err := rep.Load(path); err == nil {
				cw.StoreModel(0, 0, data.(*pmx.PmxModel))
			} else {
				mlog.ET(mi18n.T("読み込み失敗"), err.Error())
			}
		},
	)

	vmdLoader11Picker := widget.NewVmdVpdLoadFilePicker(
		"vmd",
		"モーションファイル1-1",
		"モーションファイルを選択してください",
		func(cw *controller.ControlWindow, rep repository.IRepository, path string) {
			if data, err := rep.Load(path); err == nil {
				motion := data.(*vmd.VmdMotion)
				player.Reset(motion.MaxFrame())
				cw.StoreMotion(0, 0, motion)
			} else {
				mlog.ET(mi18n.T("読み込み失敗"), err.Error())
			}
		},
	)

	pmxLoad21Picker := widget.NewPmxLoadFilePicker(
		"pmx",
		"モデルファイル2-1",
		"モデルファイルを選択してください",
		func(cw *controller.ControlWindow, rep repository.IRepository, path string) {
			if data, err := rep.Load(path); err == nil {
				cw.StoreModel(1, 0, data.(*pmx.PmxModel))
			} else {
				mlog.ET(mi18n.T("読み込み失敗"), err.Error())
			}
		},
	)

	vmdLoader21Picker := widget.NewVmdVpdLoadFilePicker(
		"vmd",
		"モーションファイル2-1",
		"モーションファイルを選択してください",
		func(cw *controller.ControlWindow, rep repository.IRepository, path string) {
			if data, err := rep.Load(path); err == nil {
				motion := data.(*vmd.VmdMotion)
				player.Reset(motion.MaxFrame())
				cw.StoreMotion(1, 0, motion)
			} else {
				mlog.ET(mi18n.T("読み込み失敗"), err.Error())
			}
		},
	)

	mWidgets.Widgets = append(mWidgets.Widgets, player, pmxLoad11Picker, vmdLoader11Picker,
		pmxLoad21Picker, vmdLoader21Picker)

	return []declarative.TabPage{
		{
			Title:    "ファイル",
			AssignTo: &fileTab,
			Layout:   declarative.VBox{},
			Background: declarative.SystemColorBrush{
				Color: walk.SysColorInactiveCaption,
			},
			Children: []declarative.Widget{
				declarative.Composite{
					Layout: declarative.VBox{},
					Children: []declarative.Widget{
						declarative.TextLabel{
							Text: "表示用モデル設定説明",
						},
						pmxLoad11Picker.Widgets(),
						vmdLoader11Picker.Widgets(),
						declarative.VSeparator{},
						pmxLoad21Picker.Widgets(),
						vmdLoader21Picker.Widgets(),
						declarative.VSeparator{},
						player.Widgets(),
						declarative.VSpacer{},
					},
				},
			},
		},
	}
}
