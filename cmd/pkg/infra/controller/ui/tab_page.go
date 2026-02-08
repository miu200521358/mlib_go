//go:build windows
// +build windows

// 指示: miu200521358
package ui

import (
	"errors"
	"fmt"

	"github.com/miu200521358/mlib_go/cmd/pkg/adapter/mpresenter/messages"
	"github.com/miu200521358/mlib_go/pkg/adapter/audio_api"
	"github.com/miu200521358/mlib_go/pkg/adapter/io_common"
	"github.com/miu200521358/mlib_go/pkg/adapter/io_motion"
	"github.com/miu200521358/mlib_go/pkg/domain/motion"
	"github.com/miu200521358/mlib_go/pkg/infra/controller"
	"github.com/miu200521358/mlib_go/pkg/infra/controller/widget"
	"github.com/miu200521358/mlib_go/pkg/shared/base"
	"github.com/miu200521358/mlib_go/pkg/shared/base/config"
	"github.com/miu200521358/mlib_go/pkg/shared/base/i18n"
	"github.com/miu200521358/mlib_go/pkg/shared/base/logging"
	"github.com/miu200521358/mlib_go/pkg/shared/base/merr"
	"github.com/miu200521358/mlib_go/pkg/usecase"
	"github.com/miu200521358/walk/pkg/declarative"
	"github.com/miu200521358/walk/pkg/walk"
)

// NewTabPages はサンプル用のタブページ群を生成する。
func NewTabPages(mWidgets *controller.MWidgets, baseServices base.IBaseServices, audioPlayer audio_api.IAudioPlayer) []declarative.TabPage {
	var fileTab *walk.TabPage
	var csvTab *walk.TabPage
	var materialTab *walk.TabPage
	var vertexTab *walk.TabPage

	var translator i18n.II18n
	var logger logging.ILogger
	var userConfig config.IUserConfig
	if baseServices != nil {
		translator = baseServices.I18n()
		logger = baseServices.Logger()
		if cfg := baseServices.Config(); cfg != nil {
			userConfig = cfg.UserConfig()
		}
	}
	if logger == nil {
		logger = logging.DefaultLogger()
	}

	player := widget.NewMotionPlayer(translator)
	player.SetAudioPlayer(audioPlayer, userConfig)

	materialView := widget.NewMaterialTableView(
		translator,
		i18n.TranslateOrMark(translator, messages.HelpMaterialView),
		func(cw *controller.ControlWindow, indexes []int) {
			if cw == nil {
				return
			}
			cw.SetSelectedMaterialIndexes(0, 0, indexes)
		},
	)
	vertexView := widget.NewVertexTableView(
		translator,
		i18n.TranslateOrMark(translator, messages.HelpVertexView),
	)

	allMaterialButton := widget.NewMPushButton()
	allMaterialButton.SetLabel(i18n.TranslateOrMark(translator, messages.LabelAll))
	allMaterialButton.SetMinSize(declarative.Size{Width: 50})
	allMaterialButton.SetOnClicked(func(cw *controller.ControlWindow) {
		if materialView == nil {
			return
		}
		materialView.SetAllChecked(true)
	})

	invertMaterialButton := widget.NewMPushButton()
	invertMaterialButton.SetLabel(i18n.TranslateOrMark(translator, messages.LabelInvert))
	invertMaterialButton.SetMinSize(declarative.Size{Width: 50})
	invertMaterialButton.SetOnClicked(func(cw *controller.ControlWindow) {
		if materialView == nil {
			return
		}
		materialView.InvertChecked()
	})

	pmxLoad11Picker := widget.NewPmxPmdXLoadFilePicker(
		userConfig,
		translator,
		config.UserConfigKeyPmxHistory,
		i18n.TranslateOrMark(translator, messages.LabelModelFile11),
		i18n.TranslateOrMark(translator, messages.LabelModelFileTip),
		func(cw *controller.ControlWindow, rep io_common.IFileReader, path string) {
			loadModel(logger, translator, cw, rep, path, materialView, vertexView, 0, 0)
		},
	)

	vmdLoad11Picker := widget.NewVmdVpdLoadFilePicker(
		userConfig,
		translator,
		config.UserConfigKeyVmdHistory,
		i18n.TranslateOrMark(translator, messages.LabelMotionFile11),
		i18n.TranslateOrMark(translator, messages.LabelMotionFileTip),
		func(cw *controller.ControlWindow, rep io_common.IFileReader, path string) {
			loadMotion(logger, translator, cw, rep, player, path, 0, 0)
		},
	)

	cameraVmdLoad11Picker := widget.NewVmdLoadFilePicker(
		userConfig,
		translator,
		config.UserConfigKeyCameraVmdHistory,
		i18n.TranslateOrMark(translator, messages.LabelCameraMotionFile11),
		i18n.TranslateOrMark(translator, messages.LabelCameraMotionFileTip),
		func(cw *controller.ControlWindow, rep io_common.IFileReader, path string) {
			loadCameraMotion(logger, translator, cw, rep, player, path, 0)
		},
	)

	pmxLoad21Picker := widget.NewPmxPmdXLoadFilePicker(
		userConfig,
		translator,
		config.UserConfigKeyPmxHistory,
		i18n.TranslateOrMark(translator, messages.LabelModelFile21),
		i18n.TranslateOrMark(translator, messages.LabelModelFileTip),
		func(cw *controller.ControlWindow, rep io_common.IFileReader, path string) {
			loadModel(logger, translator, cw, rep, path, nil, nil, 1, 0)
		},
	)

	vmdLoad21Picker := widget.NewVmdVpdLoadFilePicker(
		userConfig,
		translator,
		config.UserConfigKeyVmdHistory,
		i18n.TranslateOrMark(translator, messages.LabelMotionFile21),
		i18n.TranslateOrMark(translator, messages.LabelMotionFileTip),
		func(cw *controller.ControlWindow, rep io_common.IFileReader, path string) {
			loadMotion(logger, translator, cw, rep, player, path, 1, 0)
		},
	)

	cameraVmdLoad21Picker := widget.NewVmdLoadFilePicker(
		userConfig,
		translator,
		config.UserConfigKeyCameraVmdHistory,
		i18n.TranslateOrMark(translator, messages.LabelCameraMotionFile21),
		i18n.TranslateOrMark(translator, messages.LabelCameraMotionFileTip),
		func(cw *controller.ControlWindow, rep io_common.IFileReader, path string) {
			loadCameraMotion(logger, translator, cw, rep, player, path, 1)
		},
	)

	csvMotionPath := ""
	csvOutputPath := ""
	csvMotionRepository := io_common.IFileReader(nil)
	var csvSavePicker *widget.FilePicker

	csvMotionPicker := widget.NewVmdVpdLoadFilePicker(
		userConfig,
		translator,
		config.UserConfigKeyVmdHistory,
		i18n.TranslateOrMark(translator, messages.LabelCsvMotionFile),
		i18n.TranslateOrMark(translator, messages.LabelCsvMotionFileTip),
		func(cw *controller.ControlWindow, rep io_common.IFileReader, path string) {
			csvMotionRepository = rep
			csvMotionPath = path
			if csvSavePicker != nil && path != "" {
				csvSavePicker.SetPath(buildMotionCsvDefaultOutputPath(path))
			}
		},
	)

	csvSavePicker = widget.NewCsvSaveFilePicker(
		userConfig,
		translator,
		i18n.TranslateOrMark(translator, messages.LabelCsvOutputFile),
		i18n.TranslateOrMark(translator, messages.LabelCsvOutputFileTip),
		func(cw *controller.ControlWindow, rep io_common.IFileReader, path string) {
			csvOutputPath = path
		},
	)

	csvSaveButton := widget.NewMPushButton()
	csvSaveButton.SetLabel(i18n.TranslateOrMark(translator, messages.LabelCsvSave))
	csvSaveButton.SetMinSize(declarative.Size{Width: 90})
	csvSaveButton.SetOnClicked(func(cw *controller.ControlWindow) {
		saveMotionCsv(logger, translator, csvMotionRepository, csvMotionPath, csvOutputPath, cw)
	})

	mWidgets.Widgets = append(mWidgets.Widgets, player, pmxLoad11Picker, vmdLoad11Picker, cameraVmdLoad11Picker,
		pmxLoad21Picker, vmdLoad21Picker, cameraVmdLoad21Picker, csvMotionPicker, csvSavePicker, csvSaveButton,
		materialView, allMaterialButton, invertMaterialButton, vertexView)

	mWidgets.SetOnLoaded(func() {
		if mWidgets == nil || mWidgets.Window() == nil {
			return
		}
		mWidgets.Window().SetOnEnabledInPlaying(func(playing bool) {
			for _, w := range mWidgets.Widgets {
				w.SetEnabledInPlaying(playing)
			}
		})
		if vertexView != nil {
			vertexView.StartSelectionSync(0, 0)
		}
	})

	fileTabPage := declarative.TabPage{
		Title:    i18n.TranslateOrMark(translator, messages.LabelFile),
		AssignTo: &fileTab,
		Layout:   declarative.VBox{},
		Background: declarative.SolidColorBrush{
			Color: controller.ColorTabBackground,
		},
		Children: []declarative.Widget{
			declarative.Composite{
				Layout: declarative.VBox{},
				Children: []declarative.Widget{
					declarative.TextLabel{Text: i18n.TranslateOrMark(translator, messages.HelpDisplayModelSetting)},
					pmxLoad11Picker.Widgets(),
					vmdLoad11Picker.Widgets(),
					cameraVmdLoad11Picker.Widgets(),
					declarative.VSeparator{},
					pmxLoad21Picker.Widgets(),
					vmdLoad21Picker.Widgets(),
					cameraVmdLoad21Picker.Widgets(),
					declarative.VSeparator{},
					player.Widgets(),
					declarative.VSpacer{},
				},
			},
		},
	}

	csvTabPage := declarative.TabPage{
		Title:    i18n.TranslateOrMark(translator, messages.LabelCsv),
		AssignTo: &csvTab,
		Layout:   declarative.VBox{},
		Background: declarative.SolidColorBrush{
			Color: controller.ColorTabBackground,
		},
		Children: []declarative.Widget{
			declarative.Composite{
				Layout: declarative.VBox{},
				Children: []declarative.Widget{
					csvMotionPicker.Widgets(),
					csvSavePicker.Widgets(),
					declarative.VSeparator{},
					csvSaveButton.Widgets(),
					declarative.VSpacer{},
				},
			},
		},
	}

	materialTabPage := declarative.TabPage{
		Title:    i18n.TranslateOrMark(translator, messages.LabelMaterialView),
		AssignTo: &materialTab,
		Layout:   declarative.VBox{},
		Background: declarative.SolidColorBrush{
			Color: controller.ColorTabBackground,
		},
		Children: []declarative.Widget{
			declarative.Composite{
				Layout: declarative.HBox{},
				Children: []declarative.Widget{
					declarative.TextLabel{Text: i18n.TranslateOrMark(translator, messages.LabelMaterialView)},
					declarative.HSpacer{},
					allMaterialButton.Widgets(),
					invertMaterialButton.Widgets(),
				},
			},
			materialView.Widgets(),
		},
	}

	vertexTabPage := declarative.TabPage{
		Title:    i18n.TranslateOrMark(translator, messages.LabelVertexView),
		AssignTo: &vertexTab,
		Layout:   declarative.VBox{},
		Background: declarative.SolidColorBrush{
			Color: controller.ColorTabBackground,
		},
		Children: []declarative.Widget{
			declarative.Composite{
				Layout: declarative.HBox{},
				Children: []declarative.Widget{
					declarative.TextLabel{Text: i18n.TranslateOrMark(translator, messages.LabelVertexView)},
					declarative.HSpacer{},
				},
			},
			vertexView.Widgets(),
		},
	}

	return []declarative.TabPage{fileTabPage, csvTabPage, materialTabPage, vertexTabPage}
}

// NewTabPage はサンプル用のタブページを生成する。
func NewTabPage(mWidgets *controller.MWidgets, baseServices base.IBaseServices, audioPlayer audio_api.IAudioPlayer) declarative.TabPage {
	return NewTabPages(mWidgets, baseServices, audioPlayer)[0]
}

// loadModel はモデル読み込み結果をControlWindowへ反映する。
func loadModel(logger logging.ILogger, translator i18n.II18n, cw *controller.ControlWindow, rep io_common.IFileReader, path string, materialView *widget.MaterialTableView, vertexView *widget.VertexTableView, windowIndex, modelIndex int) {
	if cw == nil {
		return
	}
	if path == "" {
		if materialView != nil {
			materialView.ResetRows(nil)
		}
		if vertexView != nil {
			vertexView.ResetRows(nil)
		}
		cw.SetModel(windowIndex, modelIndex, nil)
		return
	}
	modelData, err := usecase.LoadModel(rep, path)
	if err != nil {
		logLoadFailed(logger, translator, err)
		if materialView != nil {
			materialView.ResetRows(nil)
		}
		if vertexView != nil {
			vertexView.ResetRows(nil)
		}
		cw.SetModel(windowIndex, modelIndex, nil)
		return
	}
	if modelData == nil {
		if materialView != nil {
			materialView.ResetRows(nil)
		}
		if vertexView != nil {
			vertexView.ResetRows(nil)
		}
		cw.SetModel(windowIndex, modelIndex, nil)
		return
	}
	if materialView != nil {
		materialView.ResetRows(modelData)
	}
	cw.SetModel(windowIndex, modelIndex, modelData)
	if vertexView != nil {
		vertexView.ResetRows(modelData)
	}
}

// loadMotion はモーション読み込み結果をControlWindowへ反映する。
func loadMotion(logger logging.ILogger, translator i18n.II18n, cw *controller.ControlWindow, rep io_common.IFileReader, player *widget.MotionPlayer, path string, windowIndex, modelIndex int) {
	if cw == nil {
		return
	}
	if path == "" {
		cw.SetMotion(windowIndex, modelIndex, nil)
		return
	}
	motionResult, err := usecase.LoadMotionWithMeta(rep, path)
	if err != nil {
		logLoadFailed(logger, translator, err)
		cw.SetMotion(windowIndex, modelIndex, nil)
		return
	}
	motionData := (*motion.VmdMotion)(nil)
	maxFrame := motion.Frame(0)
	if motionResult != nil {
		motionData = motionResult.Motion
		maxFrame = motionResult.MaxFrame
	}
	if motionData == nil {
		cw.SetMotion(windowIndex, modelIndex, nil)
		return
	}
	if player != nil {
		player.Reset(maxFrame)
	}
	cw.SetMotion(windowIndex, modelIndex, motionData)
}

// saveMotionCsv は指定モーションをCSVへ保存する。
func saveMotionCsv(logger logging.ILogger, translator i18n.II18n, rep io_common.IFileReader, sourcePath string, outputPath string, cw *controller.ControlWindow) {
	if sourcePath == "" {
		logErrorTitle(
			logger,
			i18n.TranslateOrMark(translator, messages.MessageMotionCsvExportFailed),
			errors.New(i18n.TranslateOrMark(translator, messages.MessageMotionCsvSourcePathRequired)),
		)
		return
	}
	if outputPath == "" {
		logErrorTitle(
			logger,
			i18n.TranslateOrMark(translator, messages.MessageMotionCsvExportFailed),
			errors.New(i18n.TranslateOrMark(translator, messages.MessageMotionCsvOutputPathRequired)),
		)
		return
	}
	if rep == nil {
		rep = io_motion.NewVmdVpdRepository()
	}

	motionResult, err := usecase.LoadMotionWithMeta(rep, sourcePath)
	if err != nil {
		logLoadFailed(logger, translator, err)
		return
	}
	if motionResult == nil || motionResult.Motion == nil {
		logLoadFailed(
			logger,
			translator,
			errors.New(i18n.TranslateOrMark(translator, messages.MessageMotionCsvMotionNotFound)),
		)
		return
	}
	if err := exportMotionCsvByOutputPath(outputPath, motionResult.Motion); err != nil {
		logErrorTitle(logger, i18n.TranslateOrMark(translator, messages.MessageMotionCsvExportFailed), err)
		return
	}

	completedTitle := i18n.TranslateOrMark(translator, messages.MessageMotionCsvExportCompleted)
	completedMessage := fmt.Sprintf(
		i18n.TranslateOrMark(translator, messages.MessageMotionCsvExportCompletedDetail),
		outputPath,
	)
	if logger != nil {
		logger.Info("%s", completedMessage)
	}
	controller.Beep()
	if cw != nil {
		walk.MsgBox(cw, completedTitle, completedMessage, walk.MsgBoxIconInformation)
	}
}

// loadCameraMotion はカメラモーション読み込み結果をControlWindowへ反映する。
func loadCameraMotion(logger logging.ILogger, translator i18n.II18n, cw *controller.ControlWindow, rep io_common.IFileReader, player *widget.MotionPlayer, path string, windowIndex int) {
	if cw == nil {
		return
	}
	if path == "" {
		cw.SetCameraMotion(windowIndex, nil)
		return
	}
	motionResult, err := usecase.LoadCameraMotionWithMeta(rep, path)
	if err != nil {
		logLoadFailed(logger, translator, err)
		cw.SetCameraMotion(windowIndex, nil)
		return
	}
	motionData := (*motion.VmdMotion)(nil)
	maxFrame := motion.Frame(0)
	if motionResult != nil {
		motionData = motionResult.Motion
		maxFrame = motionResult.MaxFrame
	}
	if motionData == nil {
		cw.SetCameraMotion(windowIndex, nil)
		return
	}
	if player != nil {
		player.Reset(maxFrame)
	}
	cw.SetCameraMotion(windowIndex, motionData)
}

// logLoadFailed は読み込み失敗ログを出力する。
func logLoadFailed(logger logging.ILogger, translator i18n.II18n, err error) {
	if logger == nil {
		logger = logging.DefaultLogger()
	}
	logErrorTitle(logger, i18n.TranslateOrMark(translator, messages.MessageLoadFailed), err)
}

// logErrorTitle はタイトル付きエラーを出力する。
func logErrorTitle(logger logging.ILogger, title string, err error) {
	if logger == nil {
		return
	}
	if titled, ok := logger.(interface {
		ErrorTitle(title string, err error, msg string, params ...any)
	}); ok {
		titled.ErrorTitle(title, err, "")
		return
	}
	errText := ""
	if err != nil {
		errText = err.Error()
		if errID := merr.ExtractErrorID(err); errID != "" {
			errText = "エラーID: " + errID + "\n" + errText
		}
	}
	logger.Error("%s: %s", title, errText)
}
