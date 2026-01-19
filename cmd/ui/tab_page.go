//go:build windows
// +build windows

// 指示: miu200521358
package ui

import (
	"errors"

	"github.com/miu200521358/mlib_go/pkg/adapter/io_common"
	"github.com/miu200521358/mlib_go/pkg/domain/model"
	"github.com/miu200521358/mlib_go/pkg/domain/motion"
	"github.com/miu200521358/mlib_go/pkg/infra/base/i18n"
	"github.com/miu200521358/mlib_go/pkg/infra/base/logging"
	"github.com/miu200521358/mlib_go/pkg/infra/controller"
	"github.com/miu200521358/mlib_go/pkg/infra/controller/widget"
	"github.com/miu200521358/walk/pkg/declarative"
	"github.com/miu200521358/walk/pkg/walk"
)

// overrideBoneInserter は不足ボーンの補完を行うI/F。
type overrideBoneInserter interface {
	InsertShortageOverrideBones() error
}

// NewTabPage はサンプル用のタブページを生成する。
func NewTabPage(mWidgets *controller.MWidgets) declarative.TabPage {
	var fileTab *walk.TabPage

	player := widget.NewMotionPlayer()

	pmxLoad11Picker := widget.NewPmxXLoadFilePicker(
		"pmx",
		i18n.T("モデルファイル1-1"),
		i18n.T("モデルファイルを選択してください"),
		func(cw *controller.ControlWindow, rep io_common.IFileReader, path string) {
			loadModel(cw, rep, path, 0, 0)
		},
	)

	vmdLoad11Picker := widget.NewVmdVpdLoadFilePicker(
		"vmd",
		i18n.T("モーションファイル1-1"),
		i18n.T("モーションファイルを選択してください"),
		func(cw *controller.ControlWindow, rep io_common.IFileReader, path string) {
			loadMotion(cw, rep, player, path, 0, 0)
		},
	)

	pmxLoad21Picker := widget.NewPmxXLoadFilePicker(
		"pmx",
		i18n.T("モデルファイル2-1"),
		i18n.T("モデルファイルを選択してください"),
		func(cw *controller.ControlWindow, rep io_common.IFileReader, path string) {
			loadModel(cw, rep, path, 1, 0)
		},
	)

	vmdLoad21Picker := widget.NewVmdVpdLoadFilePicker(
		"vmd",
		i18n.T("モーションファイル2-1"),
		i18n.T("モーションファイルを選択してください"),
		func(cw *controller.ControlWindow, rep io_common.IFileReader, path string) {
			loadMotion(cw, rep, player, path, 1, 0)
		},
	)

	mWidgets.Widgets = append(mWidgets.Widgets, player, pmxLoad11Picker, vmdLoad11Picker,
		pmxLoad21Picker, vmdLoad21Picker)

	mWidgets.SetOnLoaded(func() {
		if mWidgets == nil || mWidgets.Window() == nil {
			return
		}
		mWidgets.Window().SetOnEnabledInPlaying(func(playing bool) {
			for _, w := range mWidgets.Widgets {
				w.SetEnabledInPlaying(playing)
			}
		})
	})

	return declarative.TabPage{
		Title:    i18n.T("ファイル"),
		AssignTo: &fileTab,
		Layout:   declarative.VBox{},
		Background: declarative.SolidColorBrush{
			Color: controller.ColorTabBackground,
		},
		Children: []declarative.Widget{
			declarative.Composite{
				Layout: declarative.VBox{},
				Children: []declarative.Widget{
					declarative.TextLabel{Text: i18n.T("表示用モデル設定説明")},
					pmxLoad11Picker.Widgets(),
					vmdLoad11Picker.Widgets(),
					declarative.VSeparator{},
					pmxLoad21Picker.Widgets(),
					vmdLoad21Picker.Widgets(),
					declarative.VSeparator{},
					player.Widgets(),
					declarative.VSpacer{},
				},
			},
		},
	}
}

// loadModel はモデル読み込み結果をControlWindowへ反映する。
func loadModel(cw *controller.ControlWindow, rep io_common.IFileReader, path string, windowIndex, modelIndex int) {
	if cw == nil {
		return
	}
	if path == "" {
		cw.SetModel(windowIndex, modelIndex, nil)
		return
	}
	if rep == nil {
		logLoadFailed(errors.New("モデル読み込みリポジトリがありません"))
		cw.SetModel(windowIndex, modelIndex, nil)
		return
	}
	data, err := rep.Load(path)
	if err != nil {
		logLoadFailed(err)
		cw.SetModel(windowIndex, modelIndex, nil)
		return
	}
	modelData, ok := data.(*model.PmxModel)
	if !ok {
		logLoadFailed(errors.New("モデル形式が不正です"))
		cw.SetModel(windowIndex, modelIndex, nil)
		return
	}
	if modelData.Bones != nil {
		if inserter, ok := any(modelData.Bones).(overrideBoneInserter); ok {
			if err := inserter.InsertShortageOverrideBones(); err != nil {
				logging.DefaultLogger().ErrorTitle(i18n.T("システム用ボーン追加失敗"), err, "")
			}
		}
	}
	cw.SetModel(windowIndex, modelIndex, modelData)
}

// loadMotion はモーション読み込み結果をControlWindowへ反映する。
func loadMotion(cw *controller.ControlWindow, rep io_common.IFileReader, player *widget.MotionPlayer, path string, windowIndex, modelIndex int) {
	if cw == nil {
		return
	}
	if path == "" {
		cw.SetMotion(windowIndex, modelIndex, nil)
		return
	}
	if rep == nil {
		logLoadFailed(errors.New("モーション読み込みリポジトリがありません"))
		cw.SetMotion(windowIndex, modelIndex, nil)
		return
	}
	data, err := rep.Load(path)
	if err != nil {
		logLoadFailed(err)
		cw.SetMotion(windowIndex, modelIndex, nil)
		return
	}
	motionData, ok := data.(*motion.VmdMotion)
	if !ok {
		logLoadFailed(errors.New("モーション形式が不正です"))
		cw.SetMotion(windowIndex, modelIndex, nil)
		return
	}
	if player != nil {
		player.Reset(motionData.MaxFrame())
	}
	cw.SetMotion(windowIndex, modelIndex, motionData)
}

// logLoadFailed は読み込み失敗ログを出力する。
func logLoadFailed(err error) {
	logging.DefaultLogger().ErrorTitle(i18n.T("読み込み失敗"), err, "")
}
