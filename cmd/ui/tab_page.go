package ui

import (
	"github.com/miu200521358/mlib_go/pkg/config/mi18n"
	"github.com/miu200521358/mlib_go/pkg/config/mlog"
	"github.com/miu200521358/mlib_go/pkg/domain/pmx"
	"github.com/miu200521358/mlib_go/pkg/domain/vmd"
	"github.com/miu200521358/mlib_go/pkg/infrastructure/repository"
	"github.com/miu200521358/mlib_go/pkg/interface/controller"
	"github.com/miu200521358/mlib_go/pkg/interface/controller/widget"
	"github.com/miu200521358/walk/pkg/declarative"
	"github.com/miu200521358/walk/pkg/walk"
)

func NewTabPage(mWidgets *controller.MWidgets) declarative.TabPage {
	var fileTab *walk.TabPage

	player := widget.NewMotionPlayer()

	pmxLoad11Picker := widget.NewPmxXLoadFilePicker(
		"pmx",
		"モデルファイル1-1",
		"モデルファイルを選択してください",
		func(cw *controller.ControlWindow, rep repository.IRepository, path string) {
			if data, err := rep.Load(path); err == nil {
				model := data.(*pmx.PmxModel)
				if err := model.Bones.InsertShortageBones(); err != nil {
					mlog.ET(mi18n.T("システム用ボーン追加失敗"), err.Error())
				} else {
					cw.StoreModel(0, 0, model)
				}
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

	pmxLoad21Picker := widget.NewPmxXLoadFilePicker(
		"pmx",
		"モデルファイル2-1",
		"モデルファイルを選択してください",
		func(cw *controller.ControlWindow, rep repository.IRepository, path string) {
			if data, err := rep.Load(path); err == nil {
				model := data.(*pmx.PmxModel)
				if err := model.Bones.InsertShortageBones(); err != nil {
					mlog.ET(mi18n.T("システム用ボーン追加失敗"), err.Error())
				} else {
					cw.StoreModel(1, 0, model)
				}
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

	return declarative.TabPage{
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
	}
}
