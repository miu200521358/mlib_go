package file_picker

import (
	"syscall"
	"unicode/utf16"
	"unsafe"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"

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
	HistoryKey string
	// フィルタ拡張子
	FilterExtension map[string]string
	// ピッカータイトル
	Title *widget.Label
	// ファイルパス入力欄
	PathEntry *widget.Entry
	// ファイル名入力欄
	EntryName *widget.Label
	// ツールチップ用テキスト
	TooltipText string
}

func NewFilePicker(window fyne.Window,
	historyKey string,
	title string,
	tooltip string,
	filterExtension map[string]string,
	onChanged func(string),
) (*FilePicker, error) {
	picker := FilePicker{}
	picker.HistoryKey = historyKey
	picker.FilterExtension = filterExtension

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
	picker.PathEntry.OnChanged = onChanged
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

// ファイル選択ダイアログを表示する
func (picker *FilePicker) ShowFileDialog() {
	var fileBuf [4096]uint16
	picker.ofn.lpstrFile = &fileBuf[0]
	picker.ofn.nMaxFile = uint32(len(fileBuf))

	ret, _, _ := procGetOpenFileName.Call(uintptr(unsafe.Pointer(&picker.ofn)))
	if ret != 0 {
		picker.PathEntry.SetText(syscall.UTF16ToString(fileBuf[:]))
	}
}

// 履歴リストを文字列選択ダイアログで表示する
func (picker *FilePicker) ShowHistoryDialog() {
	if picker.HistoryKey == "" {
		return
	}

	// 履歴用キーから履歴リストを取得
	choices := config.LoadUserConfig(picker.HistoryKey)
	selectEntry := widget.NewSelect(choices, func(value string) {
		picker.PathEntry.SetText(value)
	})

	dialog.NewCustom("ファイル選択ダイアログ", "OK", container.NewVBox(
		widget.NewLabel("選択してください:"),
		selectEntry,
	), *picker.Window).Show()
}
