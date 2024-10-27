//go:build windows
// +build windows

package widget

import (
	"path/filepath"
	"runtime"

	"github.com/miu200521358/walk/pkg/walk"

	"github.com/miu200521358/mlib_go/pkg/domain/core"
	"github.com/miu200521358/mlib_go/pkg/infrastructure/repository"
	"github.com/miu200521358/mlib_go/pkg/infrastructure/state"
	"github.com/miu200521358/mlib_go/pkg/mutils"
	"github.com/miu200521358/mlib_go/pkg/mutils/mconfig"
	"github.com/miu200521358/mlib_go/pkg/mutils/mi18n"
	"github.com/miu200521358/mlib_go/pkg/mutils/mlog"
)

type FilePicker struct {
	walk.WidgetBase
	title             string                    // ファイルタイトル
	historyKey        string                    // 履歴用キー(空欄の場合、保存用ファイルと見なす)
	filterExtension   map[int]map[string]string // フィルタ拡張子
	pathEdit          *walk.LineEdit            // ファイルパス入力欄
	nameEdit          *walk.LineEdit            // ファイル名入力欄
	openPushButton    *walk.PushButton          // 開くボタン
	historyPushButton *walk.PushButton          // 履歴ボタン
	onPathChanged     func(string)              // パス変更時のコールバック
	limitHistory      int                       // 履歴リスト
	rep               repository.IRepository    // リポジトリ
	initialDirPath    string                    // 初期ディレクトリ
	cacheData         core.IHashModel           // キャッシュデータ
	window            state.IControlWindow      // コントローラー画面
	historyDialog     *walk.Dialog              // 履歴ダイアログ
}

func NewPmxReadFilePicker(
	window state.IControlWindow,
	parent walk.Container,
	historyKey string,
	title string,
	tooltip string,
	description string,
) *FilePicker {
	return newFilePicker(
		window,
		parent,
		historyKey,
		title,
		tooltip,
		description,
		map[int]map[string]string{
			0: {"*.pmx": "Pmx Files (*.pmx)"}, 1: {"*.*": "All Files (*.*)"}},
		50,
		repository.NewPmxRepository(),
	)
}

func NewPmxJsonReadFilePicker(
	window state.IControlWindow,
	parent walk.Container,
	historyKey string,
	title string,
	tooltip string,
	description string,
) *FilePicker {
	return newFilePicker(
		window,
		parent,
		historyKey,
		title,
		tooltip,
		description,
		map[int]map[string]string{
			0: {"*.pmx": "Pmx Files (*.pmx)"},
			1: {"*.json": "Json Files (*.json)"},
			2: {"*.*": "All Files (*.*)"}},
		50,
		repository.NewPmxPmxJsonRepository(),
	)
}

func NewCsvReadFilePicker(
	window state.IControlWindow,
	parent walk.Container,
	historyKey string,
	title string,
	tooltip string,
	description string,
) *FilePicker {
	return newFilePicker(
		window,
		parent,
		historyKey,
		title,
		tooltip,
		description,
		map[int]map[string]string{
			0: {"*.csv": "Csv Files (*.csv)"},
			2: {"*.*": "All Files (*.*)"}},
		50,
		repository.NewCsvRepository(),
	)
}

func NewVpdReadFilePicker(
	window state.IControlWindow,
	parent walk.Container,
	historyKey string,
	title string,
	tooltip string,
	description string,
) *FilePicker {
	return newFilePicker(
		window,
		parent,
		historyKey,
		title,
		tooltip,
		description,
		map[int]map[string]string{
			0: {"*.vpd": "Vpd Files (*.vpd)"}, 1: {"*.*": "All Files (*.*)"}},
		50,
		repository.NewVpdRepository())
}

func NewVmdVpdReadFilePicker(
	window state.IControlWindow,
	parent walk.Container,
	historyKey string,
	title string,
	tooltip string,
	description string,
) *FilePicker {
	return newFilePicker(
		window,
		parent,
		historyKey,
		title,
		tooltip,
		description,
		map[int]map[string]string{
			0: {"*.vmd": "Vmd Files (*.vmd)"},
			1: {"*.vpd": "Vpd Files (*.vpd)"},
			2: {"*.*": "All Files (*.*)"}},
		50,
		repository.NewVmdVpdRepository())
}

func NewVmdReadFilePicker(
	window state.IControlWindow,
	parent walk.Container,
	historyKey string,
	title string,
	tooltip string,
	description string,
) *FilePicker {
	return newFilePicker(
		window,
		parent,
		historyKey,
		title,
		tooltip,
		description,
		map[int]map[string]string{
			0: {"*.vmd": "Vmd Files (*.vmd)"},
			1: {"*.*": "All Files (*.*)"}},
		50,
		repository.NewVmdRepository())
}

func NewVmdSaveFilePicker(
	window state.IControlWindow,
	parent walk.Container,
	title string,
	tooltip string,
	description string,
) *FilePicker {
	return newFilePicker(
		window,
		parent,
		"",
		title,
		tooltip,
		description,
		map[int]map[string]string{
			0: {"*.vmd": "Vmd Files (*.vmd)"},
			1: {"*.*": "All Files (*.*)"}},
		0,
		repository.NewVmdRepository())
}

func NewPmxSaveFilePicker(
	window state.IControlWindow,
	parent walk.Container,
	title string,
	tooltip string,
	description string,
) *FilePicker {
	return newFilePicker(
		window,
		parent,
		"",
		title,
		tooltip,
		description,
		map[int]map[string]string{
			0: {"*.pmx": "Pmx Files (*.pmx)"},
			1: {"*.*": "All Files (*.*)"}},
		0,
		repository.NewPmxRepository())
}

func NewPmxJsonSaveFilePicker(
	window state.IControlWindow,
	parent walk.Container,
	title string,
	tooltip string,
	description string,
) *FilePicker {
	return newFilePicker(
		window,
		parent,
		"",
		title,
		tooltip,
		description,
		map[int]map[string]string{
			0: {"*.json": "Json Files (*.json)"},
			1: {"*.*": "All Files (*.*)"}},
		0,
		repository.NewPmxJsonRepository())
}

func NewCsvSaveFilePicker(
	window state.IControlWindow,
	parent walk.Container,
	title string,
	tooltip string,
	description string,
) *FilePicker {
	return newFilePicker(
		window,
		parent,
		"",
		title,
		tooltip,
		description,
		map[int]map[string]string{
			0: {"*.csv": "Csv Files (*.csv)"},
			1: {"*.*": "All Files (*.*)"}},
		0,
		repository.NewCsvRepository())
}

func newFilePicker(
	window state.IControlWindow,
	parent walk.Container,
	historyKey string,
	title string,
	tooltip string,
	description string,
	filterExtension map[int]map[string]string,
	limitHistory int,
	rep repository.IRepository,
) *FilePicker {
	picker := new(FilePicker)
	picker.title = title
	picker.historyKey = historyKey
	picker.filterExtension = filterExtension
	picker.limitHistory = limitHistory
	picker.rep = rep
	picker.window = window

	// タイトル
	titleComposite, err := walk.NewComposite(parent)
	if err != nil {
		RaiseError(err)
	}
	titleComposite.SetLayout(walk.NewHBoxLayout())

	titleLabel, err := walk.NewTextLabel(titleComposite)
	if err != nil {
		RaiseError(err)
	}
	titleLabel.SetText(title)
	titleLabel.MouseDown().Attach(func(x, y int, button walk.MouseButton) {
		mlog.IL(description)
	})

	if historyKey != "" {
		startBracketLabel, err := walk.NewTextLabel(titleComposite)
		if err != nil {
			RaiseError(err)
		}
		startBracketLabel.SetText("  (")
		startBracketLabel.SetToolTipText(tooltip)

		nameLineEdit, err := walk.NewLineEditStaticEdge(titleComposite)
		if err != nil {
			RaiseError(err)
		}
		nameLineEdit.SetText(mi18n.T("未設定"))
		nameLineEdit.SetReadOnly(true)
		bg, err := walk.NewSystemColorBrush(walk.SysColorInactiveCaption)
		if err != nil {
			RaiseError(err)
		}
		nameLineEdit.SetBackground(bg)
		picker.nameEdit = nameLineEdit

		endBracketLabel, err := walk.NewTextLabel(titleComposite)
		if err != nil {
			RaiseError(err)
		}
		endBracketLabel.SetText(")")
	}

	// パス入力欄
	inputComposite, err := walk.NewComposite(parent)
	if err != nil {
		RaiseError(err)
	}
	inputComposite.SetLayout(walk.NewHBoxLayout())

	picker.pathEdit, err = walk.NewLineEdit(inputComposite)
	if err != nil {
		RaiseError(err)
	}
	picker.pathEdit.SetToolTipText(tooltip)
	picker.pathEdit.DropFiles().Attach(func(files []string) {
		if len(files) > 0 {
			path := files[0]
			// パスを入力欄に設定
			picker.pathEdit.ChangeText(path)
			// コールバックを呼び出し
			picker.OnChanged(path)
		}
	})
	picker.pathEdit.TextChanged().Attach(func() {
		// コールバックを呼び出し
		picker.OnChanged(picker.pathEdit.Text())
	})

	picker.openPushButton, err = walk.NewPushButton(inputComposite)
	if err != nil {
		RaiseError(err)
	}
	picker.openPushButton.SetToolTipText(tooltip)
	picker.openPushButton.SetText(mi18n.T("開く"))
	picker.openPushButton.SetMinMaxSize(walk.Size{Width: 70, Height: 20}, walk.Size{Width: 70, Height: 20})
	picker.openPushButton.Clicked().Attach(picker.onClickOpenButton())

	if historyKey != "" {
		picker.historyPushButton, err = walk.NewPushButton(inputComposite)
		if err != nil {
			RaiseError(err)
		}
		picker.historyPushButton.SetToolTipText(tooltip)
		picker.historyPushButton.SetText(mi18n.T("履歴"))
		picker.historyPushButton.SetMinMaxSize(walk.Size{Width: 70, Height: 20}, walk.Size{Width: 70, Height: 20})
		picker.historyDialog, err = picker.createHistoryDialog()
		if err != nil {
			RaiseError(err)
		}

		picker.historyPushButton.Clicked().Attach(picker.onClickHistoryButton())
	}

	return picker
}

func (picker *FilePicker) Load() (core.IHashModel, error) {
	if picker.pathEdit.Text() == "" || picker.rep == nil {
		return nil, nil
	}

	if ok, err := picker.rep.CanLoad(picker.pathEdit.Text()); !ok || err != nil {
		return nil, err
	}

	// キャッシュの有無は見ずに、必ず取得し直す
	data, err := picker.rep.Load(picker.pathEdit.Text())
	defer runtime.GC() // 読み込み時のメモリ解放

	if err != nil {
		RaiseError(err)
	}
	picker.cacheData = data

	return data, nil
}

// パスが正しいことが分かっている上でデータだけ取り直したい場合
func (picker *FilePicker) LoadForce() core.IHashModel {
	data, err := picker.rep.Load(picker.pathEdit.Text())
	defer runtime.GC() // 読み込み時のメモリ解放

	if err != nil {
		return nil
	}

	return data
}

func (picker *FilePicker) ClearCache() {
	picker.cacheData = nil
}

func (picker *FilePicker) GetCache() core.IHashModel {
	return picker.cacheData
}

func (picker *FilePicker) SetCache(data core.IHashModel) {
	if data == nil {
		picker.pathEdit.ChangeText("")
		return
	}

	picker.cacheData = data
	picker.pathEdit.ChangeText(data.Path())
}

func (picker *FilePicker) SetPath(path string) {
	// コールバックを呼び出し
	picker.pathEdit.SetText(path)
}

func (picker *FilePicker) ChangePath(path string) {
	// コールバックを呼び出さない
	picker.pathEdit.ChangeText(path)
}

func (picker *FilePicker) GetPath() string {
	return picker.pathEdit.Text()
}

func (picker *FilePicker) SetName(path string) {
	if picker.nameEdit == nil {
		return
	}
	picker.nameEdit.SetText(path)
}

func (picker *FilePicker) GetName() string {
	if picker.nameEdit == nil {
		return ""
	}
	return picker.nameEdit.Text()
}

func (picker *FilePicker) OnChanged(path string) {
	if picker.rep != nil && picker.historyKey != "" {
		if path == "" {
			picker.nameEdit.SetText(mi18n.T("未設定"))
		} else {
			picker.nameEdit.SetText(picker.rep.LoadName(path))

			if picker.onPathChanged != nil {
				picker.onPathChanged(path)
			}

			if ok, err := picker.rep.CanLoad(picker.pathEdit.Text()); ok && err == nil {
				// ロード系のみ履歴用キーを指定して履歴リストを保存
				mconfig.SaveUserConfig(picker.historyKey, path, picker.limitHistory)
			} else {
				// 読み込めない場合、拒否
				picker.pathEdit.ChangeText("")
			}
		}
	}
}

func (picker *FilePicker) SetOnPathChanged(onPathChanged func(string)) {
	picker.onPathChanged = onPathChanged
}

func (picker *FilePicker) onClickHistoryButton() walk.EventHandler {
	return func() {
		if dlg, err := picker.createHistoryDialog(); dlg != nil && err == nil {
			if ok := dlg.Run(); ok == walk.DlgCmdOK {
				// コールバックを呼び出し
				picker.OnChanged(picker.pathEdit.Text())
			}
			dlg.Dispose()
		}
	}
}

func (picker *FilePicker) createHistoryDialog() (*walk.Dialog, error) {
	// 履歴リストを取得
	choices := mconfig.LoadUserConfig(picker.historyKey)

	// 履歴ダイアログを開く
	dlg, err := walk.NewDialog(picker.window)
	if err != nil {
		walk.MsgBox(nil, mi18n.T("履歴ダイアログ生成エラー"), err.Error(), walk.MsgBoxIconError)
		RaiseError(err)
	}
	dlg.SetTitle(mi18n.T("履歴ダイアログタイトル", map[string]interface{}{"Title": picker.title}))
	dlg.SetLayout(walk.NewVBoxLayout())
	dlg.SetSize(walk.Size{Width: 800, Height: 400})

	// 履歴リストを表示
	historyListBox, err := walk.NewListBox(dlg)
	if err != nil {
		walk.MsgBox(nil, mi18n.T("履歴リスト生成エラー"), err.Error(), walk.MsgBoxIconError)
		RaiseError(err)
	}

	// OKボタンの動作を定義する関数
	itemActivated := func() {
		// 選択されたアイテムを取得
		index := historyListBox.CurrentIndex()
		if index < 0 {
			return
		}
		item := choices[index]
		// パスを入力欄に設定
		picker.pathEdit.ChangeText(item)
	}

	historyListBox.SetModel(choices)
	historyListBox.SetMinMaxSize(walk.Size{Width: 800, Height: 400}, walk.Size{Width: 800, Height: 400})
	// 先頭を表示する（選択はできない）
	historyListBox.SetCurrentIndex(-1)

	// ダブルクリック時の動作を定義
	historyListBox.ItemActivated().Attach(func() {
		itemActivated()
		dlg.Accept()
	})

	// ボタンBox
	buttonComposite, err := walk.NewComposite(dlg)
	if err != nil {
		walk.MsgBox(nil, mi18n.T("ボタンBox生成エラー"), err.Error(), walk.MsgBoxIconError)
		RaiseError(err)
	}
	buttonComposite.SetLayout(walk.NewHBoxLayout())

	// OKボタン
	okButton, err := walk.NewPushButton(buttonComposite)
	if err != nil {
		walk.MsgBox(nil, mi18n.T("OKボタン生成エラー"), err.Error(), walk.MsgBoxIconError)
		RaiseError(err)
	}
	okButton.SetText("OK")
	okButton.Clicked().Attach(func() {
		itemActivated()
		dlg.Accept()
	})

	// Cancel ボタン
	cancelButton, err := walk.NewPushButton(buttonComposite)
	if err != nil {
		walk.MsgBox(nil, mi18n.T("Cancelボタン生成エラー"), err.Error(), walk.MsgBoxIconError)
		RaiseError(err)
	}
	cancelButton.SetText("Cancel")
	cancelButton.Clicked().Attach(func() {
		// ダイアログを閉じる
		dlg.Cancel()
	})

	return dlg, nil
}

func (picker *FilePicker) onClickOpenButton() walk.EventHandler {
	return func() {
		if picker.pathEdit.Text() != "" {
			// ファイルパスからディレクトリパスを取得
			dirPath := filepath.Dir(picker.pathEdit.Text())
			// ファイルパスのディレクトリを初期パスとして設定
			picker.initialDirPath = dirPath
		} else if picker.historyKey != "" {
			// 履歴用キーを指定して履歴リストを取得
			choices := mconfig.LoadUserConfig(picker.historyKey)
			if len(choices) > 0 {
				// ファイルパスからディレクトリパスを取得
				dirPath := filepath.Dir(choices[0])
				// 履歴リストの先頭を初期パスとして設定
				picker.initialDirPath = dirPath
			}
		}

		// ファイル選択ダイアログを開く
		dlg := walk.FileDialog{
			Title: mi18n.T(
				"ファイル選択ダイアログタイトル",
				map[string]interface{}{"Title": picker.title}),
			Filter:         picker.convertFilterExtension(picker.filterExtension),
			FilterIndex:    1,
			InitialDirPath: picker.initialDirPath,
		}
		if ok, err := dlg.ShowOpen(nil); err != nil {
			walk.MsgBox(nil, mi18n.T("ファイル選択ダイアログ選択エラー"), err.Error(), walk.MsgBoxIconError)
		} else if ok {
			// パスを入力欄に設定
			picker.pathEdit.ChangeText(dlg.FilePath)
			// コールバックを呼び出し
			picker.OnChanged(dlg.FilePath)
		}
	}
}

func (picker *FilePicker) convertFilterExtension(filterExtension map[int]map[string]string) string {
	var filterString string
	for i := 0; i < len(filterExtension); i++ {
		extData := filterExtension[i]
		for ext, desc := range extData {
			if filterString != "" {
				filterString = filterString + "|"
			}
			filterString = filterString + desc + "|" + ext
		}
	}
	return filterString
}

func (picker *FilePicker) CreateLayoutItem(ctx *walk.LayoutContext) walk.LayoutItem {
	return &filePickerLayoutItem{idealSize: walk.SizeFrom96DPI(walk.Size{Width: 50, Height: 50}, ctx.DPI())}
}

func (picker *FilePicker) Exists() bool {
	if picker.pathEdit.Text() == "" {
		return false
	}
	isExist, err := mutils.ExistsFile(picker.pathEdit.Text())
	if err != nil {
		return false
	}
	return isExist
}

func (picker *FilePicker) ExistsOrEmpty() bool {
	if picker.pathEdit.Text() == "" {
		return true
	}
	isExist, err := mutils.ExistsFile(picker.pathEdit.Text())
	if err != nil {
		return false
	}
	return isExist
}

func (picker *FilePicker) SetEnabled(enable bool) {
	picker.pathEdit.SetEnabled(enable)
	picker.openPushButton.SetEnabled(enable)
	if picker.historyPushButton != nil {
		picker.historyPushButton.SetEnabled(enable)
	}
}

func (picker *FilePicker) Enabled() bool {
	return picker.pathEdit.Enabled()
}

func (picker *FilePicker) SetVisible(visible bool) {
	picker.pathEdit.SetVisible(visible)
	picker.openPushButton.SetVisible(visible)
	if picker.historyPushButton != nil {
		picker.historyPushButton.SetVisible(visible)
	}
}

func (picker *FilePicker) Dispose() {
	picker.pathEdit.Dispose()
	picker.openPushButton.Dispose()
	if picker.historyPushButton != nil {
		picker.historyPushButton.Dispose()
	}
	picker.WidgetBase.Dispose()
}

type filePickerLayoutItem struct {
	walk.LayoutItemBase
	idealSize walk.Size // in native pixels
}

func (li *filePickerLayoutItem) LayoutFlags() walk.LayoutFlags {
	return 0
}

func (li *filePickerLayoutItem) IdealSize() walk.Size {
	return li.idealSize
}
