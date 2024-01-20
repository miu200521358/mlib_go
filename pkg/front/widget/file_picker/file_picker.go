package file_picker

import (
	"path/filepath"
	"syscall"
	"unicode/utf16"
	"unsafe"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"

	"github.com/miu200521358/mlib_go/pkg/core/reader"
	"github.com/miu200521358/mlib_go/pkg/pmx/pmx_reader"
	"github.com/miu200521358/mlib_go/pkg/utils/config"
)

const (
	// 既存ファイルのみ選択可能とする（開く用）
	OFN_FILEMUSTEXIST = 0x1000
	// 既存ファイルが選択された時に上書き確認プロンプトを表示する（保存用）
	OFN_OVERWRITEPROMPT = 0x2
)

var (
	comdlg32            = syscall.NewLazyDLL("comdlg32.dll")
	procGetOpenFileName = comdlg32.NewProc("GetOpenFileNameW")
)

type OpenFileName struct {
	lStructSize       uint32
	hwndOwner         syscall.Handle
	hInstance         syscall.Handle
	lpstrFilter       *uint16
	lpstrCustomFilter *uint16
	nMaxCustFilter    uint32
	nFilterIndex      uint32
	lpstrFile         *uint16
	nMaxFile          uint32
	lpstrFileTitle    *uint16
	nMaxFileTitle     uint32
	lpstrInitialDir   *uint16
	lpstrTitle        *uint16
	Flags             uint32
	nFileOffset       uint16
	nFileExtension    uint16
	lpstrDefExt       *uint16
	lCustData         uintptr
	lpfnHook          uintptr
	lpTemplateName    *uint16
	pvReserved        unsafe.Pointer
	dwReserved        uint32
	FlagsEx           uint32
}

type FilePicker struct {
	ofn OpenFileName
	// ウィンドウ
	appWindow *fyne.Window
	// コンテナ
	Container *fyne.Container
	// 履歴用キー(空欄の場合、保存用ファイルと見なす)
	historyKey string
	// フィルタ拡張子
	filterExtension map[string]string
	// ファイルパス入力欄
	PathEntry *widget.Entry
	// ファイル名入力欄
	EntryName *widget.Label
	// パス変更時のコールバック
	onPathChanged func(string)
	// 履歴リスト
	limitHistory int
	// reader
	modelReader reader.ReaderInterface
	// 初期ディレクトリのUTF-16エンコードされたバージョンを保持
	initialDirUtf16 []uint16
	// 履歴選択ダイアログ
	historyDialog *dialog.CustomDialog
}

func NewPmxReadFilePicker(
	window *fyne.Window,
	historyKey string,
	title string,
	placeHolder string,
	onPathChanged func(string),
) (*FilePicker, error) {
	picker, err := NewFilePicker(
		window,
		historyKey,
		title,
		placeHolder,
		map[string]string{"*.pmx": "Pmx Files (*.pmx)", "*.*": "All Files (*.*)"},
		50,
		&pmx_reader.PmxReader{},
		onPathChanged)
	return picker, err
}

func NewPmxSaveFilePicker(
	window *fyne.Window,
	historyKey string,
	title string,
	placeHolder string,
	onPathChanged func(string),
) (*FilePicker, error) {
	picker, err := NewFilePicker(
		window,
		"",
		title,
		placeHolder,
		map[string]string{"*.pmx": "Pmx Files (*.pmx)", "*.*": "All Files (*.*)"},
		0,
		nil,
		onPathChanged)
	return picker, err
}

func NewFilePicker(
	window *fyne.Window,
	historyKey string,
	title string,
	placeHolder string,
	filterExtension map[string]string,
	limitHistory int,
	modelReader reader.ReaderInterface,
	onPathChanged func(string),
) (*FilePicker, error) {
	picker := FilePicker{}
	picker.historyKey = historyKey
	picker.filterExtension = filterExtension
	picker.limitHistory = limitHistory
	picker.modelReader = modelReader
	picker.onPathChanged = onPathChanged
	picker.appWindow = window
	if historyKey != "" {
		picker.historyDialog = picker.NewHistoryDialog()
	}

	picker.ofn = OpenFileName{}
	picker.ofn.lStructSize = uint32(unsafe.Sizeof(picker.ofn))
	buf := make([]uint16, 260)
	picker.ofn.lpstrFile = &buf[0]
	picker.ofn.nMaxFile = uint32(len(buf))
	titlePtr, err := syscall.UTF16PtrFromString(title + "の選択")
	if err != nil {
		return nil, err
	}
	picker.ofn.lpstrTitle = titlePtr

	if historyKey != "" {
		// 履歴用キーが指定されている場合、開く専用
		picker.ofn.Flags = OFN_FILEMUSTEXIST
	} else {
		// 履歴用キーが指定されていない場合、保存専用
		picker.ofn.Flags = OFN_OVERWRITEPROMPT
	}

	// 許可拡張子
	filterStringUTF16 := convertFilterExtension(filterExtension)
	filterStringPtr := &filterStringUTF16[0]
	picker.ofn.lpstrFilter = filterStringPtr
	picker.ofn.nFilterIndex = 0

	picker.PathEntry = widget.NewEntry()
	picker.PathEntry.OnChanged = picker.OnChanged
	picker.PathEntry.SetPlaceHolder(placeHolder)
	picker.EntryName = widget.NewLabel("")

	// タイトルコンテナ
	titleContainer := container.New(layout.NewHBoxLayout(),
		widget.NewLabel(title),
		layout.NewSpacer())
	// ボタン入力欄コンテナ
	buttonsContainer := container.New(layout.NewHBoxLayout(),
		widget.NewButton("開く", picker.ShowFileDialog))
	if historyKey != "" {
		// 履歴用キーが指定されている場合
		// モデル名を表示
		titleContainer.Add(widget.NewLabel("("))
		titleContainer.Add(picker.EntryName)
		titleContainer.Add(widget.NewLabel(")"))
		// 履歴ボタンを追加
		buttonsContainer.Add(widget.NewButton("履歴", picker.ShowHistoryDialog))
	}
	buttonsContainer.Resize(fyne.NewSize(200, buttonsContainer.MinSize().Height))

	// ファイルパス入力コンテナ
	pathContainer := container.New(
		layout.NewBorderLayout(nil, nil, nil, buttonsContainer),
		picker.PathEntry,
		buttonsContainer)

	// ピッカーコンテナ
	picker.Container = container.New(layout.NewVBoxLayout(), titleContainer, pathContainer)

	return &picker, nil
}

func convertFilterExtension(filterExtension map[string]string) []uint16 {
	var filterString []uint16
	for ext, desc := range filterExtension {
		filterString = append(filterString, utf16.Encode([]rune(desc+"\x00"+ext+"\x00"))...)
	}
	// 最後にUTF-16のヌル終端を追加
	filterString = append(filterString, 0)
	return filterString
}

func (picker *FilePicker) OnChanged(path string) {
	picker.PathEntry.SetText(path)

	if picker.modelReader != nil && picker.historyKey != "" {
		modelName, err := picker.modelReader.ReadNameByFilepath(path)
		if err != nil {
			picker.EntryName.SetText("読み込み失敗")
		} else {
			picker.EntryName.SetText(modelName)
		}
	}

	if picker.historyKey != "" {
		// 履歴用キーを指定して履歴リストを保存
		config.SaveUserConfig(picker.historyKey, path, picker.limitHistory)
	}

	if picker.onPathChanged != nil {
		picker.onPathChanged(path)
	}
}

// ファイル選択ダイアログを表示する
func (picker *FilePicker) ShowFileDialog() {
	var fileBuf [4096]uint16
	picker.ofn.lpstrFile = &fileBuf[0]
	picker.ofn.nMaxFile = uint32(len(fileBuf))

	if picker.historyKey != "" {
		// 履歴用キーを指定して履歴リストを取得
		choices, _ := config.LoadUserConfig(picker.historyKey)
		if len(choices) > 0 {
			// ファイルパスからディレクトリパスを取得
			dirPath := filepath.Dir(choices[0])
			// 履歴リストの先頭を初期パスとして設定
			picker.initialDirUtf16 = utf16.Encode([]rune(dirPath))
			picker.ofn.lpstrInitialDir = &picker.initialDirUtf16[0]
		}
	} else if picker.PathEntry.Text != "" {
		// ファイルパスからディレクトリパスを取得
		dirPath := filepath.Dir(picker.PathEntry.Text)
		// ファイルパスのディレクトリを初期パスとして設定
		picker.initialDirUtf16 = utf16.Encode([]rune(dirPath))
		picker.ofn.lpstrInitialDir = &picker.initialDirUtf16[0]
	}

	ret, _, _ := procGetOpenFileName.Call(uintptr(unsafe.Pointer(&picker.ofn)))
	if ret != 0 {
		path := syscall.UTF16ToString(fileBuf[:])
		picker.PathEntry.SetText(path)
	}
}

func (picker *FilePicker) NewHistoryDialog() *dialog.CustomDialog {
	// 履歴用キーから履歴リストを取得
	choices, _ := config.LoadUserConfig(picker.historyKey)

	listWidget := widget.NewList(
		func() int {
			return len(choices)
		},
		func() fyne.CanvasObject {
			return picker.NewHistoryLabel()
		},
		func(id widget.ListItemID, item fyne.CanvasObject) {
			item.(*HistoryLabel).SetText(choices[id])
		},
	)

	// フォントサイズ
	fontSize := int(listWidget.MinSize().Height)
	// 表示したい行数
	linesToShow := 15
	// リスト全体の高さを計算
	listHeight := linesToShow * fontSize

	// スクロール可能なコンテナにリストを追加し、高さを設定
	scrollContainer := container.NewVScroll(listWidget)
	scrollContainer.SetMinSize(fyne.NewSize(500, float32(listHeight)))

	d := dialog.NewCustom("ファイル選択ダイアログ", "Close", container.NewVBox(
		widget.NewLabel("ファイルを選択すると、ファイルパスが入力欄に入力されます。\n"+
			"ダブルクリックすると、ファイルパスが入力欄に入力され、ダイアログを閉じます。"),
		scrollContainer,
	), *picker.appWindow)
	d.Resize(fyne.NewSize(600, 700))
	return d
}

// 履歴リストを文字列選択ダイアログで表示する
func (picker *FilePicker) ShowHistoryDialog() {
	picker.historyDialog.Show()
}

func (picker *FilePicker) GetPath() string {
	return picker.PathEntry.Text
}

func (picker *FilePicker) SetPath(path string) {
	picker.PathEntry.SetText(path)
}

type HistoryLabel struct {
	widget.Label
	picker *FilePicker
}

func (picker *FilePicker) NewHistoryLabel() *HistoryLabel {
	label := HistoryLabel{}
	label.Label = *widget.NewLabel("")
	label.picker = picker
	return &label
}

func (l *HistoryLabel) Tapped(_ *fyne.PointEvent) {
	// シングルは履歴を入れるだけ
	l.picker.PathEntry.SetText(l.Text)
}

func (l *HistoryLabel) DoubleTapped(_ *fyne.PointEvent) {
	// ダブルタップで履歴リストをクローズ
	l.picker.PathEntry.SetText(l.Text)
	l.picker.historyDialog.Hide()
}

func (l *HistoryLabel) MinSize() fyne.Size {
	return l.BaseWidget.MinSize()
}
