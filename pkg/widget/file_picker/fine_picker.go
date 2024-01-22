package file_picker

import (
	"path/filepath"

	"github.com/miu200521358/walk/pkg/walk"
	"github.com/miu200521358/win"

	"github.com/miu200521358/mlib_go/pkg/core/reader"
	"github.com/miu200521358/mlib_go/pkg/pmx/pmx_reader"
	"github.com/miu200521358/mlib_go/pkg/utils/config"
)

const FilePickerClass = "FilePicker Class"

type FilePicker struct {
	walk.WidgetBase
	// ファイルタイトル
	title string
	// 履歴用キー(空欄の場合、保存用ファイルと見なす

	// 履歴用キー(空欄の場合、保存用ファイルと見なす)
	historyKey string
	// フィルタ拡張子
	filterExtension map[string]string
	// ファイルパス入力欄
	PathLineEdit *walk.LineEdit
	// ファイル名入力欄
	NameLineEdit *walk.LineEdit
	// パス変更時のコールバック
	onPathChanged func(string)
	// 履歴リスト
	limitHistory int
	// reader
	modelReader reader.ReaderInterface
	// 初期ディレクトリ
	initialDirPath string
}

func NewPmxReadFilePicker(
	parent walk.Container,
	historyKey string,
	title string,
	tooltip string,
	onPathChanged func(string),
) error {
	return NewFilePicker(
		parent,
		historyKey,
		title,
		tooltip,
		map[string]string{"*.pmx": "Pmx Files (*.pmx)", "*.*": "All Files (*.*)"},
		50,
		&pmx_reader.PmxReader{},
		onPathChanged)
}

func NewFilePicker(
	parent walk.Container,
	historyKey string,
	title string,
	tooltip string,
	filterExtension map[string]string,
	limitHistory int,
	modelReader reader.ReaderInterface,
	onPathChanged func(string),
) error {
	picker := new(FilePicker)
	picker.title = title
	picker.historyKey = historyKey
	picker.filterExtension = filterExtension
	picker.onPathChanged = onPathChanged
	picker.limitHistory = limitHistory
	picker.modelReader = modelReader

	if err := walk.InitWidget(
		picker,
		parent,
		FilePickerClass,
		win.WS_VISIBLE,
		0); err != nil {

		return err
	}

	// ピッカー全体
	pickerComposite, err := walk.NewComposite(parent)
	if err != nil {
		return err
	}
	pickerLayout := walk.NewVBoxLayout()
	pickerComposite.SetLayout(pickerLayout)

	// タイトル
	titleComposite, err := walk.NewComposite(pickerComposite)
	if err != nil {
		return err
	}
	titleLayout := walk.NewHBoxLayout()
	titleComposite.SetLayout(titleLayout)

	startBracket, err := walk.NewTextLabel(titleComposite)
	if err != nil {
		return err
	}
	startBracket.SetText(title + "  (")
	startBracket.SetToolTipText(tooltip)

	nameLineEdit, err := walk.NewLineEdit(titleComposite)
	if err != nil {
		return err
	}
	nameLineEdit.SetText(title)
	nameLineEdit.SetToolTipText(tooltip)
	nameLineEdit.SetReadOnly(true)
	picker.NameLineEdit = nameLineEdit

	endBracket, err := walk.NewTextLabel(titleComposite)
	if err != nil {
		return err
	}
	endBracket.SetText(")")

	// パス入力欄
	inputComposite, err := walk.NewComposite(pickerComposite)
	if err != nil {
		return err
	}
	inputLayout := walk.NewHBoxLayout()
	inputComposite.SetLayout(inputLayout)

	pathTextEdit, err := walk.NewLineEdit(inputComposite)
	if err != nil {
		return err
	}
	pathTextEdit.SetToolTipText(tooltip)
	picker.PathLineEdit = pathTextEdit

	openPushButton, err := walk.NewPushButton(inputComposite)
	if err != nil {
		return err
	}
	openPushButton.SetToolTipText(tooltip)
	openPushButton.SetText("開く")
	openPushButton.Clicked().Attach(picker.onClickOpenButton())

	historyPushButton, err := walk.NewPushButton(inputComposite)
	if err != nil {
		return err
	}
	historyPushButton.SetToolTipText(tooltip)
	historyPushButton.SetText("履歴")
	historyPushButton.Clicked().Attach(picker.onClickHistoryButton())

	return nil
}

func (picker *FilePicker) onClickHistoryButton() walk.EventHandler {
	return func() {
		// 履歴リストを取得
		choices, _ := config.LoadUserConfig(picker.historyKey)

		// 履歴ダイアログを開く
		dlg, err := walk.NewDialog(picker.Form())
		if err != nil {
			walk.MsgBox(nil, "履歴ダイアログ生成エラー", err.Error(), walk.MsgBoxIconError)
			return
		}
		dlg.SetTitle(picker.title + "ファイルの履歴")
		dlg.SetLayout(walk.NewVBoxLayout())
		dlg.SetSize(walk.Size{Width: 500, Height: 400})

		// 履歴リストを表示
		historyListBox, err := walk.NewListBox(dlg)
		if err != nil {
			walk.MsgBox(nil, "履歴リスト生成エラー", err.Error(), walk.MsgBoxIconError)
			return
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
			// コールバックを呼び出し
			picker.onPathChanged(item)
		}

		historyListBox.SetModel(choices)
		historyListBox.SetMinMaxSize(walk.Size{Width: 500, Height: 400}, walk.Size{Width: 500, Height: 400})
		historyListBox.SetCurrentIndex(0)

		// ダブルクリック時の動作を定義
		historyListBox.ItemActivated().Attach(func() {
			itemActivated()
			dlg.Accept()
		})

		// ボタンBox
		buttonComposite, err := walk.NewComposite(dlg)
		if err != nil {
			walk.MsgBox(nil, "ボタンBox生成エラー", err.Error(), walk.MsgBoxIconError)
			return
		}
		buttonComposite.SetLayout(walk.NewHBoxLayout())

		// OKボタン
		okButton, err := walk.NewPushButton(buttonComposite)
		if err != nil {
			walk.MsgBox(nil, "OKボタン生成エラー", err.Error(), walk.MsgBoxIconError)
			return
		}
		okButton.SetText("OK")
		okButton.Clicked().Attach(func() {
			itemActivated()
			dlg.Accept()
		})

		// Cancel ボタン
		cancelButton, err := walk.NewPushButton(buttonComposite)
		if err != nil {
			walk.MsgBox(nil, "Cancelボタン生成エラー", err.Error(), walk.MsgBoxIconError)
			return
		}
		cancelButton.SetText("Cancel")
		cancelButton.Clicked().Attach(func() {
			// ダイアログを閉じる
			dlg.Cancel()
		})

		// ダイアログを表示
		dlg.Run()
	}
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
			choices, _ := config.LoadUserConfig(picker.historyKey)
			if len(choices) > 0 {
				// ファイルパスからディレクトリパスを取得
				dirPath := filepath.Dir(choices[0])
				// 履歴リストの先頭を初期パスとして設定
				picker.initialDirPath = dirPath
			}
		}

		// ファイル選択ダイアログを開く
		dlg := walk.FileDialog{
			Title:          picker.title + "ファイルを選択してください",
			Filter:         picker.convertFilterExtension(picker.filterExtension),
			FilterIndex:    1,
			InitialDirPath: picker.initialDirPath,
		}
		if ok, err := dlg.ShowOpen(nil); err != nil {
			walk.MsgBox(nil, "エラー", err.Error(), walk.MsgBoxIconError)
		} else if ok {
			// パスを入力欄に設定
			picker.PathLineEdit.SetText(dlg.FilePath)
			// コールバックを呼び出し
			picker.onPathChanged(dlg.FilePath)
		}
	}
}

func (f *FilePicker) convertFilterExtension(filterExtension map[string]string) string {
	var filterString string
	for ext, desc := range filterExtension {
		if filterString != "" {
			filterString = filterString + "|"
		}
		filterString = filterString + desc + "|" + ext
	}
	return filterString
}

func (*FilePicker) CreateLayoutItem(ctx *walk.LayoutContext) walk.LayoutItem {
	return &filePickerLayoutItem{idealSize: walk.SizeFrom96DPI(walk.Size{Width: 50, Height: 50}, ctx.DPI())}
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
