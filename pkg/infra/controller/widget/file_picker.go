//go:build windows
// +build windows

// 指示: miu200521358
package widget

import (
	"path/filepath"
	"strings"

	"github.com/miu200521358/mlib_go/pkg/adapter/io_common"
	"github.com/miu200521358/mlib_go/pkg/adapter/io_model"
	"github.com/miu200521358/mlib_go/pkg/adapter/io_model/pmx"
	"github.com/miu200521358/mlib_go/pkg/adapter/io_motion"
	"github.com/miu200521358/mlib_go/pkg/adapter/io_motion/vmd"
	infraconfig "github.com/miu200521358/mlib_go/pkg/infra/base/config"
	"github.com/miu200521358/mlib_go/pkg/infra/base/i18n"
	"github.com/miu200521358/mlib_go/pkg/infra/base/logging"
	"github.com/miu200521358/mlib_go/pkg/infra/controller"
	"github.com/miu200521358/walk/pkg/declarative"
	"github.com/miu200521358/walk/pkg/walk"
)

type filterExtension struct {
	extension   string
	description string
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
	userConfig        io_commonUserConfig
	pathEdit          *walk.LineEdit
	nameEdit          *walk.LineEdit
	openPushButton    *walk.PushButton
	historyPushButton *walk.PushButton
	historyDialog     *walk.Dialog
	historyListBox    *walk.ListBox
	prevPath          string
	onPathChanged     func(*controller.ControlWindow, io_common.IFileReader, string)
}

type io_commonUserConfig interface {
	GetStringSlice(key string) ([]string, error)
	SetStringSlice(key string, values []string, limit int) error
}

// NewPmxLoadFilePicker はPMX読み込み用のFilePickerを生成する。
func NewPmxLoadFilePicker(historyKey string, title string, tooltip string, onPathChanged func(*controller.ControlWindow, io_common.IFileReader, string)) *FilePicker {
	return newFilePicker(
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

// NewPmxXLoadFilePicker はPMX/X読み込み用のFilePickerを生成する。
func NewPmxXLoadFilePicker(historyKey string, title string, tooltip string, onPathChanged func(*controller.ControlWindow, io_common.IFileReader, string)) *FilePicker {
	return newFilePicker(
		historyKey,
		title,
		tooltip,
		onPathChanged,
		[]filterExtension{
			{extension: "*.pmx;*.x", description: "Pmx/X Files (*.pmx,*.x)"},
			{extension: "*.*", description: "All Files (*.*)"},
		},
		io_model.NewModelRepository(),
	)
}

// NewVmdVpdLoadFilePicker はVMD/VPD読み込み用のFilePickerを生成する。
func NewVmdVpdLoadFilePicker(historyKey string, title string, tooltip string, onPathChanged func(*controller.ControlWindow, io_common.IFileReader, string)) *FilePicker {
	return newFilePicker(
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
func NewVmdLoadFilePicker(historyKey string, title string, tooltip string, onPathChanged func(*controller.ControlWindow, io_common.IFileReader, string)) *FilePicker {
	return newFilePicker(
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

// NewPmxSaveFilePicker はPMX保存用のFilePickerを生成する。
func NewPmxSaveFilePicker(title string, tooltip string, onPathChanged func(*controller.ControlWindow, io_common.IFileReader, string)) *FilePicker {
	return newFilePicker(
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

// NewVmdSaveFilePicker はVMD保存用のFilePickerを生成する。
func NewVmdSaveFilePicker(title string, tooltip string, onPathChanged func(*controller.ControlWindow, io_common.IFileReader, string)) *FilePicker {
	return newFilePicker(
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

// newFilePicker はFilePickerを生成する。
func newFilePicker(historyKey string, title string, tooltip string, onPathChanged func(*controller.ControlWindow, io_common.IFileReader, string),
	filterExtensions []filterExtension, repository io_common.IFileReader) *FilePicker {
	picker := &FilePicker{
		title:            title,
		tooltip:          tooltip,
		historyKey:       historyKey,
		filterExtensions: filterExtensions,
		repository:       repository,
		onPathChanged:    onPathChanged,
		userConfig:       infraconfig.NewUserConfigStore(),
	}
	return picker
}

// SetWindow はウィンドウ参照を設定する。
func (fp *FilePicker) SetWindow(window *controller.ControlWindow) {
	fp.window = window
}

// SetEnabledInPlaying は再生中の有効状態を設定する。
func (fp *FilePicker) SetEnabledInPlaying(playing bool) {
	fp.SetEnabled(!playing)
}

// SetEnabled はウィジェットの有効状態を設定する。
func (fp *FilePicker) SetEnabled(enabled bool) {
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
					Text:        i18n.T("未設定"),
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
		},
		declarative.PushButton{
			AssignTo:    &fp.openPushButton,
			Text:        i18n.T("開く"),
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
			Text:        i18n.T("履歴"),
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
	ok, err := fd.ShowOpen(fp.window)
	if err != nil {
		walk.MsgBox(fp.window, i18n.T("読み込み失敗"), err.Error(), walk.MsgBoxIconError)
		return
	}
	if !ok {
		return
	}
	fp.handlePathChanged(fd.FilePath)
}

// showSaveDialog は保存用ダイアログを表示する。
func (fp *FilePicker) showSaveDialog() {
	fd := new(walk.FileDialog)
	fd.Title = fp.title
	fd.FilePath = fp.pathEdit.Text()
	fd.Filter = fp.buildFilterString()
	fd.InitialDirPath = fp.resolveInitialDir()
	ok, err := fd.ShowSave(fp.window)
	if err != nil {
		walk.MsgBox(fp.window, i18n.T("保存失敗"), err.Error(), walk.MsgBoxIconError)
		return
	}
	if !ok {
		return
	}
	fp.handlePathChanged(fd.FilePath)
}

// handlePathChanged はパス変更時の処理を行う。
func (fp *FilePicker) handlePathChanged(path string) {
	cleaned := fp.cleanPath(path)
	if cleaned == "" || cleaned == fp.prevPath {
		return
	}
	if fp.historyKey != "" && fp.repository != nil && !fp.repository.CanLoad(cleaned) {
		return
	}

	fp.prevPath = cleaned
	if fp.pathEdit != nil {
		fp.pathEdit.SetText(cleaned)
	}
	if fp.nameEdit != nil {
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
	if fp.historyKey == "" || fp.userConfig == nil {
		return
	}
	values, err := fp.userConfig.GetStringSlice(fp.historyKey)
	if err != nil || len(values) == 0 {
		return
	}

	if fp.historyDialog != nil {
		fp.historyListBox.SetModel(values)
		fp.historyDialog.Show()
		return
	}

	dlg := new(walk.Dialog)
	lb := new(walk.ListBox)
	push := new(walk.PushButton)
	if err := (declarative.Dialog{
		AssignTo: &dlg,
		Title:    i18n.T("履歴"),
		MinSize:  declarative.Size{Width: 800, Height: 400},
		Layout:   declarative.VBox{},
		Children: []declarative.Widget{
			declarative.ListBox{
				AssignTo: &lb,
				Model:    values,
				MinSize:  declarative.Size{Width: 800, Height: 400},
				OnItemActivated: func() {
					idx := lb.CurrentIndex()
					if idx < 0 || idx >= len(values) {
						return
					}
					push.SetEnabled(true)
					fp.handlePathChanged(values[idx])
					dlg.Accept()
				},
			}, declarative.Composite{
				Layout: declarative.HBox{},
				Children: []declarative.Widget{

					declarative.PushButton{
						AssignTo: &push,
						Text:     "OK",
						Enabled:  true,
						OnClicked: func() {
							dlg.Accept()
						},
					},
					declarative.PushButton{
						Text: "Cancel",
						OnClicked: func() {
							dlg.Cancel()
						},
					},
				},
			},
		},
	}).Create(fp.window); err != nil {
		return
	}

	fp.historyDialog = dlg
	fp.historyListBox = lb
	push.SetEnabled(true)
	fp.historyDialog.Show()
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
