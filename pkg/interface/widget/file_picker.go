//go:build windows
// +build windows

package widget

import (
	"fmt"
	"path/filepath"
	"runtime"

	"github.com/miu200521358/walk/pkg/walk"
	"github.com/miu200521358/win"

	"github.com/miu200521358/mlib_go/pkg/domain/core"
	"github.com/miu200521358/mlib_go/pkg/infrastructure/reader"
	"github.com/miu200521358/mlib_go/pkg/mutils"
	"github.com/miu200521358/mlib_go/pkg/mutils/mconfig"
	"github.com/miu200521358/mlib_go/pkg/mutils/mi18n"
	"github.com/miu200521358/mlib_go/pkg/mutils/mlog"
)

const FilePickerClass = "FilePicker Class"

type FilePicker struct {
	walk.WidgetBase
	title             string                    // ファイルタイトル
	historyKey        string                    // 履歴用キー(空欄の場合、保存用ファイルと見なす)
	filterExtension   map[int]map[string]string // フィルタ拡張子
	PathLineEdit      *walk.LineEdit            // ファイルパス入力欄
	nameLineEdit      *walk.LineEdit            // ファイル名入力欄
	openPushButton    *walk.PushButton          // 開くボタン
	historyPushButton *walk.PushButton          // 履歴ボタン
	OnPathChanged     func(string)              // パス変更時のコールバック
	limitHistory      int                       // 履歴リスト
	modelReader       core.IReader              // mcore
	initialDirPath    string                    // 初期ディレクトリ
	cacheData         core.IHashModel           // キャッシュデータ
	window            *walk.MainWindow          // MainWindow
	historyDialog     *walk.Dialog              // 履歴ダイアログ
}

func NewPmxReadFilePicker(
	window *walk.MainWindow,
	parent walk.Container,
	historyKey string,
	title string,
	tooltip string,
	description string,
	OnPathChanged func(string),
) (*FilePicker, error) {
	return NewFilePicker(
		window,
		parent,
		historyKey,
		title,
		tooltip,
		description,
		map[int]map[string]string{
			0: {"*.pmx": "Pmx Files (*.pmx)"}, 2: {"*.*": "All Files (*.*)"}},
		50,
		&reader.PmxReader{},
		OnPathChanged)
}

func NewVpdReadFilePicker(
	window *walk.MainWindow,
	parent walk.Container,
	historyKey string,
	title string,
	tooltip string,
	description string,
	OnPathChanged func(string),
) (*FilePicker, error) {
	return NewFilePicker(
		window,
		parent,
		historyKey,
		title,
		tooltip,
		description,
		map[int]map[string]string{
			0: {"*.vpd": "Vpd Files (*.vpd)"}, 1: {"*.*": "All Files (*.*)"}},
		50,
		&reader.VpdMotionReader{},
		OnPathChanged)
}

func NewVmdVpdReadFilePicker(
	window *walk.MainWindow,
	parent walk.Container,
	historyKey string,
	title string,
	tooltip string,
	description string,
	OnPathChanged func(string),
) (*FilePicker, error) {
	return NewFilePicker(
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
		reader.NewVmdVpdMotionReader(),
		OnPathChanged)
}

func NewVmdReadFilePicker(
	window *walk.MainWindow,
	parent walk.Container,
	historyKey string,
	title string,
	tooltip string,
	description string,
	OnPathChanged func(string),
) (*FilePicker, error) {
	return NewFilePicker(
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
		&reader.VmdMotionReader{},
		OnPathChanged)
}

func NewVmdSaveFilePicker(
	window *walk.MainWindow,
	parent walk.Container,
	title string,
	tooltip string,
	description string,
	OnPathChanged func(string),
) (*FilePicker, error) {
	return NewFilePicker(
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
		nil,
		OnPathChanged)
}

func NewPmxSaveFilePicker(
	window *walk.MainWindow,
	parent walk.Container,
	title string,
	tooltip string,
	description string,
	onPathChanged func(string),
) (*FilePicker, error) {
	return NewFilePicker(
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
		nil,
		onPathChanged)
}

func NewFilePicker(
	window *walk.MainWindow,
	parent walk.Container,
	historyKey string,
	title string,
	tooltip string,
	description string,
	filterExtension map[int]map[string]string,
	limitHistory int,
	modelReader core.IReader,
	onPathChanged func(string),
) (*FilePicker, error) {
	picker := new(FilePicker)
	picker.title = title
	picker.historyKey = historyKey
	picker.filterExtension = filterExtension
	picker.OnPathChanged = onPathChanged
	picker.limitHistory = limitHistory
	picker.modelReader = modelReader
	picker.window = window

	if err := walk.InitWidget(
		picker,
		parent,
		FilePickerClass,
		win.WS_DISABLED,
		0); err != nil {

		return nil, err
	}

	// タイトル
	titleComposite, err := walk.NewComposite(parent)
	if err != nil {
		return nil, err
	}
	titleComposite.SetLayout(walk.NewHBoxLayout())

	titleLabel, err := walk.NewTextLabel(titleComposite)
	if err != nil {
		return nil, err
	}
	titleLabel.SetText(title)
	titleLabel.SetToolTipText(tooltip)
	titleLabel.MouseDown().Attach(func(x, y int, button walk.MouseButton) {
		mlog.IL(description)
	})

	if historyKey != "" {
		startBracketLabel, err := walk.NewTextLabel(titleComposite)
		if err != nil {
			return nil, err
		}
		startBracketLabel.SetText("  (")
		startBracketLabel.SetToolTipText(tooltip)

		nameLineEdit, err := walk.NewLineEditStaticEdge(titleComposite)
		if err != nil {
			return nil, err
		}
		nameLineEdit.SetText(mi18n.T("未設定"))
		nameLineEdit.SetToolTipText(tooltip)
		nameLineEdit.SetReadOnly(true)
		picker.nameLineEdit = nameLineEdit

		endBracketLabel, err := walk.NewTextLabel(titleComposite)
		if err != nil {
			return nil, err
		}
		endBracketLabel.SetText(")")
	}

	// パス入力欄
	inputComposite, err := walk.NewComposite(parent)
	if err != nil {
		return nil, err
	}
	inputComposite.SetLayout(walk.NewHBoxLayout())

	picker.PathLineEdit, err = walk.NewLineEdit(inputComposite)
	if err != nil {
		return nil, err
	}
	picker.PathLineEdit.SetToolTipText(tooltip)
	picker.PathLineEdit.DropFiles().Attach(func(files []string) {
		if len(files) > 0 {
			path := files[0]
			// パスを入力欄に設定
			picker.PathLineEdit.SetText(path)
			// コールバックを呼び出し
			picker.OnChanged(path)
		}
	})

	picker.openPushButton, err = walk.NewPushButton(inputComposite)
	if err != nil {
		return nil, err
	}
	picker.openPushButton.SetToolTipText(tooltip)
	picker.openPushButton.SetText(mi18n.T("開く"))
	picker.openPushButton.Clicked().Attach(picker.onClickOpenButton())

	if historyKey != "" {
		picker.historyPushButton, err = walk.NewPushButton(inputComposite)
		if err != nil {
			return nil, err
		}
		picker.historyPushButton.SetToolTipText(tooltip)
		picker.historyPushButton.SetText(mi18n.T("履歴"))
		picker.historyDialog, err = picker.createHistoryDialog()
		if err != nil {
			return nil, err
		}

		picker.historyPushButton.Clicked().Attach(picker.onClickHistoryButton())
	}

	return picker, nil
}

func (picker *FilePicker) GetData() (core.IHashModel, error) {
	if picker.PathLineEdit.Text() == "" || picker.modelReader == nil {
		return nil, nil
	}

	if isExist, err := mutils.ExistsFile(picker.PathLineEdit.Text()); err != nil || !isExist {
		return nil, fmt.Errorf(mi18n.T("ファイルが存在しません"))
	}

	// キャッシュの有無は見ずに、必ず取得し直す
	data, err := picker.modelReader.ReadByFilepath(picker.PathLineEdit.Text())
	defer runtime.GC() // 読み込み時のメモリ解放

	if err != nil {
		return nil, err
	}
	picker.cacheData = data

	return data, nil
}

// パスが正しいことが分かっている上でデータだけ取り直したい場合
func (picker *FilePicker) GetDataForce() core.IHashModel {
	data, err := picker.modelReader.ReadByFilepath(picker.PathLineEdit.Text())
	defer runtime.GC() // 読み込み時のメモリ解放

	if err != nil {
		return nil
	}

	return data
}

func (picker *FilePicker) SetCache(data core.IHashModel) {
	if data == nil {
		picker.PathLineEdit.SetText("")
		return
	}

	picker.cacheData = data
	picker.PathLineEdit.SetText(data.GetPath())
}

func (picker *FilePicker) IsCached() bool {
	if isExist, err := mutils.ExistsFile(picker.PathLineEdit.Text()); err != nil || !isExist {
		return false
	}

	hash, err := picker.modelReader.ReadHashByFilePath(picker.PathLineEdit.Text())
	if err != nil {
		return false
	}

	return picker.cacheData != nil && picker.cacheData.GetHash() == hash
}

func (picker *FilePicker) ClearCache() {
	picker.cacheData = nil
}

func (picker *FilePicker) GetCache() core.IHashModel {
	return picker.cacheData
}

func (picker *FilePicker) SetPath(path string) {
	picker.PathLineEdit.SetText(path)
}

func (picker *FilePicker) GetPath() string {
	return picker.PathLineEdit.Text()
}

func (picker *FilePicker) OnChanged(path string) {
	picker.PathLineEdit.SetText(path)

	if picker.modelReader != nil && picker.historyKey != "" {
		if path == "" {
			picker.nameLineEdit.SetText(mi18n.T("未設定"))
		} else {
			modelName, err := picker.modelReader.ReadNameByFilepath(path)
			if err != nil {
				picker.nameLineEdit.SetText(mi18n.T("読み込み失敗"))
			} else {
				picker.nameLineEdit.SetText(modelName)
			}
		}
	}

	if picker.historyKey != "" {
		// 履歴用キーを指定して履歴リストを保存
		mconfig.SaveUserConfig(picker.historyKey, path, picker.limitHistory)
	}

	if picker.OnPathChanged != nil {
		picker.OnPathChanged(path)
	}
}

func (picker *FilePicker) onClickHistoryButton() walk.EventHandler {
	return func() {
		if dlg, err := picker.createHistoryDialog(); dlg != nil && err == nil {
			if ok := dlg.Run(); ok == walk.DlgCmdOK {
				// コールバックを呼び出し
				picker.OnChanged(picker.PathLineEdit.Text())
			}
			dlg.Dispose()
		}
	}
}

func (picker *FilePicker) createHistoryDialog() (*walk.Dialog, error) {
	// 履歴リストを取得
	choices := mconfig.LoadUserConfig(picker.historyKey)

	// 履歴ダイアログを開く
	dlg, err := walk.NewDialog(picker.Form())
	if err != nil {
		walk.MsgBox(nil, mi18n.T("履歴ダイアログ生成エラー"), err.Error(), walk.MsgBoxIconError)
		return nil, err
	}
	dlg.SetTitle(mi18n.T("履歴ダイアログタイトル", map[string]interface{}{"Title": picker.title}))
	dlg.SetLayout(walk.NewVBoxLayout())
	dlg.SetSize(walk.Size{Width: 800, Height: 400})

	// 履歴リストを表示
	historyListBox, err := walk.NewListBox(dlg)
	if err != nil {
		walk.MsgBox(nil, mi18n.T("履歴リスト生成エラー"), err.Error(), walk.MsgBoxIconError)
		return nil, err
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
		picker.PathLineEdit.SetText(item)
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
		return nil, err
	}
	buttonComposite.SetLayout(walk.NewHBoxLayout())

	// OKボタン
	okButton, err := walk.NewPushButton(buttonComposite)
	if err != nil {
		walk.MsgBox(nil, mi18n.T("OKボタン生成エラー"), err.Error(), walk.MsgBoxIconError)
		return nil, err
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
		return nil, err
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
		if picker.PathLineEdit.Text() != "" {
			// ファイルパスからディレクトリパスを取得
			dirPath := filepath.Dir(picker.PathLineEdit.Text())
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
			picker.PathLineEdit.SetText(dlg.FilePath)
			// コールバックを呼び出し
			picker.OnChanged(dlg.FilePath)
		}
	}
}

func (f *FilePicker) convertFilterExtension(filterExtension map[int]map[string]string) string {
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

func (*FilePicker) CreateLayoutItem(ctx *walk.LayoutContext) walk.LayoutItem {
	return &filePickerLayoutItem{idealSize: walk.SizeFrom96DPI(walk.Size{Width: 50, Height: 50}, ctx.DPI())}
}

func (f *FilePicker) Exists() bool {
	if f.PathLineEdit.Text() == "" {
		return false
	}
	isExist, err := mutils.ExistsFile(f.PathLineEdit.Text())
	if err != nil {
		return false
	}
	return isExist
}

func (f *FilePicker) ExistsOrEmpty() bool {
	if f.PathLineEdit.Text() == "" {
		return true
	}
	isExist, err := mutils.ExistsFile(f.PathLineEdit.Text())
	if err != nil {
		return false
	}
	return isExist
}

func (f *FilePicker) SetEnabled(enable bool) {
	f.PathLineEdit.SetEnabled(enable)
	f.openPushButton.SetEnabled(enable)
	if f.historyPushButton != nil {
		f.historyPushButton.SetEnabled(enable)
	}
}

func (f *FilePicker) Enabled() bool {
	return f.PathLineEdit.Enabled()
}

func (f *FilePicker) SetVisible(visible bool) {
	f.PathLineEdit.SetVisible(visible)
	f.openPushButton.SetVisible(visible)
	if f.historyPushButton != nil {
		f.historyPushButton.SetVisible(visible)
	}
}

func (f *FilePicker) Dispose() {
	f.PathLineEdit.Dispose()
	f.openPushButton.Dispose()
	if f.historyPushButton != nil {
		f.historyPushButton.Dispose()
	}
	f.WidgetBase.Dispose()
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
