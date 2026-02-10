//go:build windows
// +build windows

// 指示: miu200521358
package widget

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/miu200521358/mlib_go/pkg/adapter/mpresenter/messages"

	"github.com/miu200521358/mlib_go/pkg/adapter/io_audio"
	"github.com/miu200521358/mlib_go/pkg/adapter/io_common"
	"github.com/miu200521358/mlib_go/pkg/adapter/io_csv"
	"github.com/miu200521358/mlib_go/pkg/adapter/io_model"
	"github.com/miu200521358/mlib_go/pkg/adapter/io_model/pmd"
	"github.com/miu200521358/mlib_go/pkg/adapter/io_model/pmx"
	"github.com/miu200521358/mlib_go/pkg/adapter/io_motion"
	"github.com/miu200521358/mlib_go/pkg/adapter/io_motion/vmd"
	"github.com/miu200521358/mlib_go/pkg/infra/base/err"
	"github.com/miu200521358/mlib_go/pkg/infra/controller"
	"github.com/miu200521358/mlib_go/pkg/shared/base/i18n"
	"github.com/miu200521358/mlib_go/pkg/shared/base/logging"
	"github.com/miu200521358/walk/pkg/declarative"
	"github.com/miu200521358/walk/pkg/walk"
)

type filterExtension struct {
	extension   string
	description string
}

// FileFilterExtension はFilePickerに渡す拡張子フィルタを表す。
type FileFilterExtension struct {
	Extension   string
	Description string
}

// FilePicker はファイル選択ウィジェットを表す。
type FilePicker struct {
	window            *controller.ControlWindow
	title             string
	tooltip           string
	historyKey        string
	initialDirPath    string
	filterExtensions  []filterExtension
	repository        io_common.IFileReader
	translator        i18n.II18n
	userConfig        iCommonUserConfig
	pathEdit          *walk.LineEdit
	nameEdit          *walk.LineEdit
	openPushButton    *walk.PushButton
	historyPushButton *walk.PushButton
	historyDialog     *walk.Dialog
	historyListBox    *walk.ListBox
	historyValues     []string
	prevPath          string
	prevPathHash      string
	requestedEnabled  bool
	loading           bool
	onPathChanged     func(*controller.ControlWindow, io_common.IFileReader, string)
}

type iCommonUserConfig interface {
	GetStringSlice(key string) ([]string, error)
	SetStringSlice(key string, values []string, limit int) error
}

// NewLoadFilePicker は読み込み用のFilePickerを生成する。
func NewLoadFilePicker(
	userConfig iCommonUserConfig,
	translator i18n.II18n,
	historyKey string,
	title string,
	tooltip string,
	onPathChanged func(*controller.ControlWindow, io_common.IFileReader, string),
	filterExtensions []FileFilterExtension,
	repository io_common.IFileReader,
) *FilePicker {
	converted := make([]filterExtension, 0, len(filterExtensions))
	for _, ext := range filterExtensions {
		converted = append(converted, filterExtension{
			extension:   ext.Extension,
			description: ext.Description,
		})
	}
	return newFilePicker(
		userConfig,
		translator,
		historyKey,
		title,
		tooltip,
		onPathChanged,
		converted,
		repository,
	)
}

// NewPmxLoadFilePicker はPMX読み込み用のFilePickerを生成する。
func NewPmxLoadFilePicker(userConfig iCommonUserConfig, translator i18n.II18n, historyKey string, title string, tooltip string, onPathChanged func(*controller.ControlWindow, io_common.IFileReader, string)) *FilePicker {
	return newFilePicker(
		userConfig,
		translator,
		historyKey,
		title,
		tooltip,
		onPathChanged,
		[]filterExtension{
			{extension: "*.pmx", description: "Pmx Files (*.pmx)"},
			{extension: "*.*", description: "All Files (*.*)"},
		},
		pmx.NewPmxRepository(),
	)
}

// NewPmdLoadFilePicker はPMD読み込み用のFilePickerを生成する。
func NewPmdLoadFilePicker(userConfig iCommonUserConfig, translator i18n.II18n, historyKey string, title string, tooltip string, onPathChanged func(*controller.ControlWindow, io_common.IFileReader, string)) *FilePicker {
	return newFilePicker(
		userConfig,
		translator,
		historyKey,
		title,
		tooltip,
		onPathChanged,
		[]filterExtension{
			{extension: "*.pmd", description: "Pmd Files (*.pmd)"},
			{extension: "*.*", description: "All Files (*.*)"},
		},
		pmd.NewPmdRepository(),
	)
}

// NewPmxPmdXLoadFilePicker はPMX/PMD/X読み込み用のFilePickerを生成する。
func NewPmxPmdXLoadFilePicker(userConfig iCommonUserConfig, translator i18n.II18n, historyKey string, title string, tooltip string, onPathChanged func(*controller.ControlWindow, io_common.IFileReader, string)) *FilePicker {
	return newFilePicker(
		userConfig,
		translator,
		historyKey,
		title,
		tooltip,
		onPathChanged,
		[]filterExtension{
			{extension: "*.pmx;*.pmd;*.x", description: "Pmx/Pmd/X Files (*.pmx;*.pmd;*.x)"},
			{extension: "*.*", description: "All Files (*.*)"},
		},
		io_model.NewModelRepository(),
	)
}

// NewPmxPmdLoadFilePicker はPMX/PMD読み込み用のFilePickerを生成する。
func NewPmxPmdLoadFilePicker(userConfig iCommonUserConfig, translator i18n.II18n, historyKey string, title string, tooltip string, onPathChanged func(*controller.ControlWindow, io_common.IFileReader, string)) *FilePicker {
	return newFilePicker(
		userConfig,
		translator,
		historyKey,
		title,
		tooltip,
		onPathChanged,
		[]filterExtension{
			{extension: "*.pmx;*.pmd", description: "Pmx/Pmd Files (*.pmx;*.pmd)"},
			{extension: "*.*", description: "All Files (*.*)"},
		},
		io_model.NewPmxPmdRepository(),
	)
}

// NewVmdVpdLoadFilePicker はVMD/VPD読み込み用のFilePickerを生成する。
func NewVmdVpdLoadFilePicker(userConfig iCommonUserConfig, translator i18n.II18n, historyKey string, title string, tooltip string, onPathChanged func(*controller.ControlWindow, io_common.IFileReader, string)) *FilePicker {
	return newFilePicker(
		userConfig,
		translator,
		historyKey,
		title,
		tooltip,
		onPathChanged,
		[]filterExtension{
			{extension: "*.vmd;*.vpd", description: "Vmd/Vpd Files (*.vmd;*.vpd)"},
			{extension: "*.*", description: "All Files (*.*)"},
		},
		io_motion.NewVmdVpdRepository(),
	)
}

// NewVmdLoadFilePicker はVMD読み込み用のFilePickerを生成する。
func NewVmdLoadFilePicker(userConfig iCommonUserConfig, translator i18n.II18n, historyKey string, title string, tooltip string, onPathChanged func(*controller.ControlWindow, io_common.IFileReader, string)) *FilePicker {
	return newFilePicker(
		userConfig,
		translator,
		historyKey,
		title,
		tooltip,
		onPathChanged,
		[]filterExtension{
			{extension: "*.vmd", description: "Vmd Files (*.vmd)"},
			{extension: "*.*", description: "All Files (*.*)"},
		},
		vmd.NewVmdRepository(),
	)
}

// NewAudioLoadFilePicker は音楽読み込み用のFilePickerを生成する。
func NewAudioLoadFilePicker(userConfig iCommonUserConfig, translator i18n.II18n, historyKey string, title string, tooltip string, onPathChanged func(*controller.ControlWindow, io_common.IFileReader, string)) *FilePicker {
	return newFilePicker(
		userConfig,
		translator,
		historyKey,
		title,
		tooltip,
		onPathChanged,
		[]filterExtension{
			{extension: "*.wav;*.mp3", description: "Audio Files (*.wav;*.mp3)"},
			{extension: "*.*", description: "All Files (*.*)"},
		},
		io_audio.NewAudioRepository(translator),
	)
}

// NewCsvLoadFilePicker はCSV読み込み用のFilePickerを生成する。
func NewCsvLoadFilePicker(userConfig iCommonUserConfig, translator i18n.II18n, historyKey string, title string, tooltip string, onPathChanged func(*controller.ControlWindow, io_common.IFileReader, string)) *FilePicker {
	return newFilePicker(
		userConfig,
		translator,
		historyKey,
		title,
		tooltip,
		onPathChanged,
		[]filterExtension{
			{extension: "*.csv", description: "Csv Files (*.csv)"},
			{extension: "*.*", description: "All Files (*.*)"},
		},
		io_csv.NewCsvRepository(),
	)
}

// NewPmxSaveFilePicker はPMX保存用のFilePickerを生成する。
func NewPmxSaveFilePicker(userConfig iCommonUserConfig, translator i18n.II18n, title string, tooltip string, onPathChanged func(*controller.ControlWindow, io_common.IFileReader, string)) *FilePicker {
	return newFilePicker(
		userConfig,
		translator,
		"",
		title,
		tooltip,
		onPathChanged,
		[]filterExtension{
			{extension: "*.pmx", description: "Pmx Files (*.pmx)"},
			{extension: "*.*", description: "All Files (*.*)"},
		},
		io_model.NewModelRepository(),
	)
}

// NewPmdSaveFilePicker はPMD保存用のFilePickerを生成する。
func NewPmdSaveFilePicker(userConfig iCommonUserConfig, translator i18n.II18n, title string, tooltip string, onPathChanged func(*controller.ControlWindow, io_common.IFileReader, string)) *FilePicker {
	return newFilePicker(
		userConfig,
		translator,
		"",
		title,
		tooltip,
		onPathChanged,
		[]filterExtension{
			{extension: "*.pmd", description: "Pmd Files (*.pmd)"},
			{extension: "*.*", description: "All Files (*.*)"},
		},
		io_model.NewModelRepository(),
	)
}

// NewVmdSaveFilePicker はVMD保存用のFilePickerを生成する。
func NewVmdSaveFilePicker(userConfig iCommonUserConfig, translator i18n.II18n, title string, tooltip string, onPathChanged func(*controller.ControlWindow, io_common.IFileReader, string)) *FilePicker {
	return newFilePicker(
		userConfig,
		translator,
		"",
		title,
		tooltip,
		onPathChanged,
		[]filterExtension{
			{extension: "*.vmd", description: "Vmd Files (*.vmd)"},
			{extension: "*.*", description: "All Files (*.*)"},
		},
		io_motion.NewVmdVpdRepository(),
	)
}

// NewCsvSaveFilePicker はCSV保存用のFilePickerを生成する。
func NewCsvSaveFilePicker(userConfig iCommonUserConfig, translator i18n.II18n, title string, tooltip string, onPathChanged func(*controller.ControlWindow, io_common.IFileReader, string)) *FilePicker {
	return newFilePicker(
		userConfig,
		translator,
		"",
		title,
		tooltip,
		onPathChanged,
		[]filterExtension{
			{extension: "*.csv", description: "Csv Files (*.csv)"},
			{extension: "*.*", description: "All Files (*.*)"},
		},
		io_csv.NewCsvRepository(),
	)
}

// newFilePicker はFilePickerを生成する。
func newFilePicker(userConfig iCommonUserConfig, translator i18n.II18n, historyKey string, title string, tooltip string,
	onPathChanged func(*controller.ControlWindow, io_common.IFileReader, string), filterExtensions []filterExtension, repository io_common.IFileReader) *FilePicker {
	picker := &FilePicker{
		title:            title,
		tooltip:          tooltip,
		historyKey:       historyKey,
		filterExtensions: filterExtensions,
		repository:       repository,
		onPathChanged:    onPathChanged,
		userConfig:       userConfig,
		translator:       translator,
		requestedEnabled: true,
	}
	return picker
}

// SetWindow はウィンドウ参照を設定する。
func (fp *FilePicker) SetWindow(window *controller.ControlWindow) {
	fp.window = window
}

// t は翻訳済み文言を返す。
func (fp *FilePicker) t(key string) string {
	if fp == nil || fp.translator == nil || !fp.translator.IsReady() {
		return "●●" + key + "●●"
	}
	return fp.translator.T(key)
}

// SetEnabledInPlaying は再生中の有効状態を設定する。
func (fp *FilePicker) SetEnabledInPlaying(playing bool) {
	fp.SetEnabled(!playing)
}

// SetEnabled はウィジェットの有効状態を設定する。
func (fp *FilePicker) SetEnabled(enabled bool) {
	if fp == nil {
		return
	}
	fp.requestedEnabled = enabled
	fp.applyEnabledState()
}

// SetPath は外部からパスを設定し、読み込み処理を実行する。
func (fp *FilePicker) SetPath(path string) {
	if fp == nil || fp.loading {
		return
	}
	fp.handlePathChanged(path)
}

// applyEnabledState は要求状態と読み込み状態からUIの有効状態を反映する。
func (fp *FilePicker) applyEnabledState() {
	if fp == nil {
		return
	}
	enabled := fp.requestedEnabled && !fp.loading
	if fp.pathEdit != nil {
		fp.pathEdit.SetEnabled(enabled)
	}
	if fp.nameEdit != nil {
		fp.nameEdit.SetEnabled(enabled)
	}
	if fp.openPushButton != nil {
		fp.openPushButton.SetEnabled(enabled)
	}
	if fp.historyPushButton != nil {
		fp.historyPushButton.SetEnabled(enabled)
	}
	fp.setHistoryDialogEnabled(enabled)
	fp.applyNameEditBackground()
}

// setHistoryDialogEnabled は履歴ダイアログの有効状態を更新する。
func (fp *FilePicker) setHistoryDialogEnabled(enabled bool) {
	if fp == nil || fp.historyDialog == nil {
		return
	}
	if fp.historyDialog.IsDisposed() {
		fp.historyDialog = nil
		fp.historyListBox = nil
		return
	}
	fp.historyDialog.SetEnabled(enabled)
	if fp.historyListBox != nil {
		fp.historyListBox.SetEnabled(enabled)
	}
}

// beginLoading は読み込み中状態へ遷移する。
func (fp *FilePicker) beginLoading() bool {
	if fp == nil || fp.loading {
		return false
	}
	fp.loading = true
	fp.applyEnabledState()
	return true
}

// endLoading は読み込み中状態を解除する。
func (fp *FilePicker) endLoading() {
	if fp == nil || !fp.loading {
		return
	}
	fp.loading = false
	fp.applyEnabledState()
}

// applyNameEditBackground は表示名欄の背景色を再設定する。
func (fp *FilePicker) applyNameEditBackground() {
	if fp == nil || fp.nameEdit == nil {
		return
	}
	bg, err := walk.NewSolidColorBrush(controller.ColorTabBackground)
	if err != nil {
		return
	}
	fp.nameEdit.SetBackground(bg)
}

// Widgets はUI構成を返す。
func (fp *FilePicker) Widgets() declarative.Composite {
	titleWidgets := []declarative.Widget{
		declarative.TextLabel{
			Text:        fp.title,
			ToolTipText: fp.tooltip,
		},
	}

	if fp.historyKey != "" {
		titleWidgets = append(titleWidgets, declarative.Composite{
			Layout: declarative.HBox{},
			Children: []declarative.Widget{
				declarative.TextLabel{
					Text:        "  (",
					ToolTipText: fp.tooltip,
				},
				declarative.LineEdit{
					ReadOnly: true,
					Background: declarative.SolidColorBrush{
						Color: controller.ColorTabBackground,
					},
					Text:        fp.t(messages.FilePickerKey001),
					ToolTipText: fp.tooltip,
					AssignTo:    &fp.nameEdit,
				},
				declarative.TextLabel{
					Text:        ") ",
					ToolTipText: fp.tooltip,
				},
			},
		})
	}

	inputWidgets := []declarative.Widget{
		declarative.LineEdit{
			AssignTo:    &fp.pathEdit,
			ToolTipText: fp.tooltip,
			OnTextChanged: func() {
				fp.handlePathChanged(fp.pathEdit.Text())
			},
			OnEditingFinished: func() {
				// フォーカス喪失時にも発火するため、同一パスの再適用は行わない。
				fp.handlePathChanged(fp.pathEdit.Text())
			},
			OnDropFiles: func(files []string) {
				fp.handleDropFiles(files)
			},
		},
		declarative.PushButton{
			AssignTo:    &fp.openPushButton,
			Text:        fp.t(messages.FilePickerKey002),
			ToolTipText: fp.tooltip,
			OnClicked: func() {
				fp.onOpenClicked()
			},
			MinSize: declarative.Size{Width: 70, Height: 20},
			MaxSize: declarative.Size{Width: 70, Height: 20},
		},
	}

	if fp.historyKey != "" {
		inputWidgets = append(inputWidgets, declarative.PushButton{
			AssignTo:    &fp.historyPushButton,
			Text:        fp.t(messages.FilePickerKey003),
			ToolTipText: fp.tooltip,
			OnClicked: func() {
				fp.openHistoryDialog()
			},
			MinSize: declarative.Size{Width: 70, Height: 20},
			MaxSize: declarative.Size{Width: 70, Height: 20},
		})
	}

	return declarative.Composite{
		Layout: declarative.VBox{},
		Children: []declarative.Widget{
			declarative.Composite{
				Layout:   declarative.HBox{},
				Children: titleWidgets,
			},
			declarative.Composite{
				Layout:   declarative.HBox{},
				Children: inputWidgets,
			},
		},
	}
}

// onOpenClicked は開くボタンの処理を行う。
func (fp *FilePicker) onOpenClicked() {
	if fp == nil || fp.loading {
		return
	}
	if fp.historyKey == "" {
		fp.showSaveDialog()
		return
	}
	fp.showOpenDialog()
}

// showOpenDialog は読み込み用ダイアログを表示する。
func (fp *FilePicker) showOpenDialog() {
	fd := new(walk.FileDialog)
	fd.Title = fp.title
	fd.FilePath = fp.pathEdit.Text()
	fd.Filter = fp.buildFilterString()
	fd.InitialDirPath = fp.resolveInitialDir()
	ok, dialogErr := fd.ShowOpen(fp.window)
	if dialogErr != nil {
		walk.MsgBox(fp.window, fp.t(messages.FilePickerKey004), err.BuildErrorText(dialogErr), walk.MsgBoxIconError)
		return
	}
	if !ok {
		return
	}
	fp.handlePathConfirmed(fd.FilePath)
}

// showSaveDialog は保存用ダイアログを表示する。
func (fp *FilePicker) showSaveDialog() {
	fd := new(walk.FileDialog)
	fd.Title = fp.title
	fd.FilePath = fp.pathEdit.Text()
	fd.Filter = fp.buildFilterString()
	fd.InitialDirPath = fp.resolveInitialDir()
	ok, dialogErr := fd.ShowSave(fp.window)
	if dialogErr != nil {
		walk.MsgBox(fp.window, fp.t(messages.FilePickerKey005), err.BuildErrorText(dialogErr), walk.MsgBoxIconError)
		return
	}
	if !ok {
		return
	}
	fp.handlePathConfirmed(fd.FilePath)
}

// handlePathChanged はパス変更時の処理を行う。
func (fp *FilePicker) handlePathChanged(path string) {
	// パス変更時は同一パスを再適用しない。
	fp.applyPath(path, false)
}

// handlePathConfirmed は明示確定時にパス変更処理を実行する。
func (fp *FilePicker) handlePathConfirmed(path string) {
	// 明示確定時は同一パスでも再適用する。
	fp.applyPath(path, true)
}

// handleDropFiles はドロップされたファイル一覧から読み込み対象を反映する。
func (fp *FilePicker) handleDropFiles(files []string) {
	if fp == nil || len(files) == 0 {
		return
	}
	// D&Dは最初に有効判定を通過した1件のみ採用する。
	for _, file := range files {
		cleaned := fp.cleanPath(file)
		if cleaned == "" {
			continue
		}
		if fp.repository != nil && !fp.repository.CanLoad(cleaned) {
			continue
		}
		fp.handlePathConfirmed(cleaned)
		return
	}
}

// applyPath はパス更新処理を共通化する。
func (fp *FilePicker) applyPath(path string, allowSame bool) {
	if fp == nil || fp.loading {
		return
	}
	cleaned := fp.cleanPath(path)
	if cleaned == "" {
		return
	}
	// 同一パス時は前回読込時の内容ハッシュ差分がない限り再適用しない。
	currentHash := fp.computeFileHash(cleaned)
	if !allowSame && cleaned == fp.prevPath && currentHash == fp.prevPathHash {
		return
	}
	if fp.historyKey != "" && fp.repository != nil && !fp.repository.CanLoad(cleaned) {
		return
	}
	if !fp.beginLoading() {
		return
	}
	defer fp.endLoading()

	fp.prevPath = cleaned
	fp.prevPathHash = currentHash
	if fp.pathEdit != nil {
		// 同値再設定による不要なTextChanged発火を避ける。
		if fp.pathEdit.Text() != cleaned {
			fp.pathEdit.SetText(cleaned)
		}
	}
	if fp.nameEdit != nil && fp.repository != nil {
		fp.nameEdit.SetText(fp.repository.InferName(cleaned))
	}
	if fp.onPathChanged != nil {
		fp.onPathChanged(fp.window, fp.repository, cleaned)
	}
	fp.saveHistoryIfNeeded(cleaned)
}

// saveHistoryIfNeeded は履歴保存が可能な場合に保存する。
func (fp *FilePicker) saveHistoryIfNeeded(path string) {
	if fp.historyKey == "" || fp.userConfig == nil {
		return
	}
	if fp.repository != nil && !fp.repository.CanLoad(path) {
		return
	}
	values, err := fp.userConfig.GetStringSlice(fp.historyKey)
	if err != nil {
		return
	}
	values = append([]string{path}, values...)
	values = dedupe(values)
	if err := fp.userConfig.SetStringSlice(fp.historyKey, values, 50); err != nil {
		logger := logging.DefaultLogger()
		logger.Warn("履歴保存に失敗しました: %s", err.Error())
	}
}

// buildFilterString は拡張子フィルタ文字列を生成する。
func (fp *FilePicker) buildFilterString() string {
	pairs := make([]string, 0, len(fp.filterExtensions)*2)
	for _, ext := range fp.filterExtensions {
		pairs = append(pairs, ext.description, ext.extension)
	}
	return strings.Join(pairs, "|")
}

// resolveInitialDir は初期ディレクトリを決定する。
func (fp *FilePicker) resolveInitialDir() string {
	if fp.pathEdit != nil {
		current := fp.cleanPath(fp.pathEdit.Text())
		if current != "" {
			return filepath.Dir(current)
		}
	}
	if fp.historyKey == "" || fp.userConfig == nil {
		return fp.initialDirPath
	}
	values, err := fp.userConfig.GetStringSlice(fp.historyKey)
	if err != nil || len(values) == 0 {
		return fp.initialDirPath
	}
	return filepath.Dir(values[0])
}

// openHistoryDialog は履歴ダイアログを表示する。
func (fp *FilePicker) openHistoryDialog() {
	if fp.historyKey == "" || fp.loading || !fp.requestedEnabled {
		return
	}
	values := []string{}
	if fp.userConfig != nil {
		var err error
		values, err = fp.userConfig.GetStringSlice(fp.historyKey)
		if err != nil {
			logger := logging.DefaultLogger()
			logger.Warn("履歴読込に失敗しました")
			values = []string{}
		}
	}
	fp.historyValues = append([]string{}, values...)

	if fp.historyDialog != nil {
		if fp.historyDialog.IsDisposed() {
			fp.historyDialog = nil
			fp.historyListBox = nil
		} else {
			fp.historyListBox.SetModel(fp.historyValues)
			fp.setHistoryDialogEnabled(fp.requestedEnabled && !fp.loading)
			fp.historyDialog.Show()
			return
		}
	}

	dlg := new(walk.Dialog)
	lb := new(walk.ListBox)
	push := new(walk.PushButton)
	var parent walk.Form
	if fp.window != nil {
		parent = fp.window
	} else {
		parent = walk.App().ActiveForm()
	}
	if parent == nil {
		return
	}
	if err := (declarative.Dialog{
		AssignTo: &dlg,
		Title:    fp.t(messages.FilePickerKey003),
		MinSize:  declarative.Size{Width: 800, Height: 400},
		Layout:   declarative.VBox{},
		Children: []declarative.Widget{
			declarative.ListBox{
				AssignTo: &lb,
				Model:    fp.historyValues,
				MinSize:  declarative.Size{Width: 800, Height: 400},
				OnItemActivated: func() {
					if !fp.applyHistoryDialogSelection(historyDialogActionConfirm, lb.CurrentIndex()) {
						return
					}
					push.SetEnabled(true)
					dlg.Accept()
				},
			}, declarative.Composite{
				Layout: declarative.HBox{},
				Children: []declarative.Widget{
					declarative.PushButton{
						AssignTo: &push,
						Text:     fp.t(messages.FilePickerKey006),
						Enabled:  true,
						OnClicked: func() {
							dlg.Accept()
							fp.applyHistoryDialogSelection(historyDialogActionConfirm, lb.CurrentIndex())
						},
					},
					declarative.PushButton{
						Text: fp.t(messages.FilePickerKey007),
						OnClicked: func() {
							fp.applyHistoryDialogSelection(historyDialogActionCancel, lb.CurrentIndex())
							dlg.Cancel()
						},
					},
				},
			},
		},
	}).Create(parent); err != nil {
		return
	}

	fp.historyDialog = dlg
	fp.historyListBox = lb
	fp.historyDialog.Disposing().Attach(func() {
		fp.historyDialog = nil
		fp.historyListBox = nil
		fp.historyValues = nil
	})
	push.SetEnabled(true)
	fp.setHistoryDialogEnabled(fp.requestedEnabled && !fp.loading)
	fp.historyDialog.Show()
}

// applyHistoryDialogSelection は履歴ダイアログの操作に応じて選択中パスを反映する。
func (fp *FilePicker) applyHistoryDialogSelection(action historyDialogAction, index int) bool {
	selectedPath, ok := resolveHistoryDialogSelection(action, index, fp.historyPathByIndex)
	if !ok {
		return false
	}
	fp.handlePathConfirmed(selectedPath)
	return true
}

// historyPathByIndex は履歴インデックスに対応するパスを返す。
func (fp *FilePicker) historyPathByIndex(index int) (string, bool) {
	if fp == nil || index < 0 || index >= len(fp.historyValues) {
		return "", false
	}
	path := fp.cleanPath(fp.historyValues[index])
	if path == "" {
		return "", false
	}
	return path, true
}

// cleanPath は入力パスを正規化する。
func (fp *FilePicker) cleanPath(path string) string {
	if path == "" {
		return ""
	}
	path = filepath.Clean(path)
	path = strings.Trim(path, "\"")
	path = strings.Trim(path, "'")
	path = strings.TrimSpace(path)
	path = strings.Trim(path, ".")
	return path
}

// computeFileHash はファイル内容のSHA-256ハッシュを返す。
func (fp *FilePicker) computeFileHash(path string) string {
	if path == "" {
		return ""
	}
	file, err := os.Open(path)
	if err != nil {
		return ""
	}
	defer file.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return ""
	}
	return hex.EncodeToString(hash.Sum(nil))
}

// dedupe は重複を排除したスライスを返す。
func dedupe(values []string) []string {
	seen := make(map[string]struct{}, len(values))
	result := make([]string, 0, len(values))
	for _, v := range values {
		if v == "" {
			continue
		}
		if _, ok := seen[v]; ok {
			continue
		}
		seen[v] = struct{}{}
		result = append(result, v)
	}
	return result
}
