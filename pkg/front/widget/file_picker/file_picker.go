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
	"github.com/miu200521358/mlib_go/pkg/utils/config"
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
	Window *fyne.Window
	// コンテナ
	Container *fyne.Container
	// 履歴用キー(空欄の場合、保存用ファイルと見なす)
	historyKey string
	// フィルタ拡張子
	filterExtension map[string]string
	// ピッカータイトル
	Title *widget.Label
	// ファイルパス入力欄
	PathEntry *widget.Entry
	// ファイル名入力欄
	EntryName *widget.Label
	// ツールチップ用テキスト
	TooltipText string
	// パス変更時のコールバック
	onPathChanged func(string)
	// 履歴リスト
	limitHistory int
	// reader
	modelReader reader.ReaderInterface
	// 初期ディレクトリのUTF-16エンコードされたバージョンを保持
	initialDirUtf16 []uint16
}

func NewFilePicker(window fyne.Window,
	historyKey string,
	title string,
	tooltip string,
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

	picker.ofn = OpenFileName{}
	picker.ofn.lStructSize = uint32(unsafe.Sizeof(picker.ofn))
	buf := make([]uint16, 260)
	picker.ofn.lpstrFile = &buf[0]
	picker.ofn.nMaxFile = uint32(len(buf))

	// 許可拡張子
	filterStringUTF16 := convertFilterExtension(filterExtension)
	filterStringPtr := &filterStringUTF16[0]
	picker.ofn.lpstrFilter = filterStringPtr
	picker.ofn.nFilterIndex = 0

	picker.Title = widget.NewLabel(title)
	picker.PathEntry = widget.NewEntry()
	picker.PathEntry.OnChanged = picker.OnChanged
	picker.EntryName = widget.NewLabel("")
	picker.TooltipText = tooltip

	// タイトルコンテナ
	titleContainer := container.New(layout.NewHBoxLayout(),
		picker.Title,
		layout.NewSpacer(),
		widget.NewLabel("("),
		picker.EntryName,
		widget.NewLabel(")"))
	// ボタン入力欄コンテナ
	buttonsContainer := container.New(layout.NewHBoxLayout(),
		widget.NewButton("開く", picker.ShowFileDialog),
		widget.NewButton("履歴", picker.ShowHistoryDialog))
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

	if picker.modelReader != nil {
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
	}

	ret, _, _ := procGetOpenFileName.Call(uintptr(unsafe.Pointer(&picker.ofn)))
	if ret != 0 {
		picker.PathEntry.SetText(syscall.UTF16ToString(fileBuf[:]))
	}
}

// 履歴リストを文字列選択ダイアログで表示する
func (picker *FilePicker) ShowHistoryDialog() {
	if picker.historyKey == "" {
		return
	}

	// 履歴用キーから履歴リストを取得
	choices, _ := config.LoadUserConfig(picker.historyKey)
	selectEntry := widget.NewSelect(choices, func(value string) {
		picker.PathEntry.SetText(value)
	})

	dialog.NewCustom("ファイル選択ダイアログ", "OK", container.NewVBox(
		widget.NewLabel("選択してください:"),
		selectEntry,
	), *picker.Window).Show()
}
