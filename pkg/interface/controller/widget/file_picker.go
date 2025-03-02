package widget

import (
	"path/filepath"

	"github.com/miu200521358/mlib_go/pkg/config/mconfig"
	"github.com/miu200521358/mlib_go/pkg/config/mi18n"
	"github.com/miu200521358/mlib_go/pkg/config/mlog"
	"github.com/miu200521358/mlib_go/pkg/infrastructure/repository"
	"github.com/miu200521358/mlib_go/pkg/interface/controller"
	"github.com/miu200521358/walk/pkg/declarative"
	"github.com/miu200521358/walk/pkg/walk"
)

type filterExtension struct {
	extension   string
	description string
}

type FilePicker struct {
	window            *controller.ControlWindow // メインウィンドウ
	title             string                    // タイトル
	tooltip           string                    // ツールチップ
	historyKey        string                    // 履歴を保存するキー
	initialDirPath    string                    // 初期ディレクトリパス
	filterExtensions  []filterExtension         // フィルター拡張子
	repository        repository.IRepository    // リポジトリ
	pathEdit          *walk.LineEdit            // パス入力欄
	nameEdit          *walk.LineEdit            // 名前欄(read-only)
	openPushButton    *walk.PushButton          // 開くボタン
	historyPushButton *walk.PushButton          // 履歴ボタン
	historyDialog     *walk.Dialog              // 履歴ダイアログ
	historyListBox    *walk.ListBox             // 履歴リスト
	onPathChanged     func(string)              // パス変更時のコールバック
}

func NewPmxLoadFilePicker(
	window *controller.ControlWindow,
	historyKey string,
	title string,
	tooltip string,
	onPathChanged func(string),
) *FilePicker {
	return newFilePicker(
		window,
		historyKey,
		title,
		tooltip,
		onPathChanged,
		[]filterExtension{
			{extension: "*.pmx", description: "Pmx Files (*.pmx)"},
			{extension: "*.*", description: "All Files (*.*)"},
		},
		repository.NewPmxRepository(),
	)
}

func NewVmdVpdLoadFilePicker(
	window *controller.ControlWindow,
	historyKey string,
	title string,
	tooltip string,
	onPathChanged func(string),
) *FilePicker {
	return newFilePicker(
		window,
		historyKey,
		title,
		tooltip,
		onPathChanged,
		[]filterExtension{
			{extension: "*.vmd;*.vpd", description: "Vmd/Vpd Files (*.vmd;*.vpd)"},
			{extension: "*.*", description: "All Files (*.*)"},
		},
		repository.NewPmxRepository(),
	)
}

func newFilePicker(
	window *controller.ControlWindow,
	historyKey string,
	title string,
	tooltip string,
	onPathChanged func(string),
	filterExtension []filterExtension,
	repository repository.IRepository,
) *FilePicker {
	picker := new(FilePicker)
	picker.title = title
	picker.tooltip = tooltip
	picker.historyKey = historyKey
	picker.initialDirPath = ""
	picker.filterExtensions = filterExtension
	picker.repository = repository
	picker.window = window
	picker.onPathChanged = onPathChanged

	return picker
}

func (fp *FilePicker) Widgets() declarative.Composite {
	titleWidgets := []declarative.Widget{
		declarative.TextLabel{
			Text:        fp.title,
			ToolTipText: fp.tooltip,
			OnMouseDown: func(x, y int, button walk.MouseButton) { mlog.IL("%s", fp.tooltip) },
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
					Background: declarative.SystemColorBrush{
						Color: walk.SysColorInactiveCaption,
					},
					Text:        mi18n.T("未設定"),
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
				fp.onChanged(fp.pathEdit.Text())
			},
			OnDropFiles: func(files []string) {
				if len(files) > 0 {
					path := files[0]
					// パスを入力欄に設定
					fp.pathEdit.ChangeText(path)
					// コールバックを呼び出し
					fp.onChanged(path)
				}
			},
		},
		declarative.PushButton{
			AssignTo:    &fp.openPushButton,
			Text:        mi18n.T("開く"),
			ToolTipText: fp.tooltip,
			OnClicked:   fp.onClickOpenButton(),
			MinSize:     declarative.Size{Width: 70, Height: 20},
			MaxSize:     declarative.Size{Width: 70, Height: 20},
		},
	}

	if fp.historyKey != "" {
		inputWidgets = append(inputWidgets, declarative.PushButton{
			AssignTo:    &fp.historyPushButton,
			Text:        mi18n.T("履歴"),
			ToolTipText: fp.tooltip,
			OnClicked: func() {
				var err error
				if fp.historyDialog, err = fp.newHistoryDialog(); fp.historyDialog != nil && err == nil {
					if ok := fp.historyDialog.Run(); ok == walk.DlgCmdOK {
						// コールバックを呼び出し
						fp.onChanged(fp.pathEdit.Text())
					}
					fp.historyDialog.Dispose()
				}
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

func (fp *FilePicker) onChanged(path string) {
	if fp.repository != nil && fp.historyKey != "" {
		if path == "" {
			fp.nameEdit.SetText(mi18n.T("未設定"))
		} else {
			fp.nameEdit.SetText(fp.repository.LoadName(path))

			if fp.onPathChanged != nil {
				fp.onPathChanged(path)
			}

			if ok, err := fp.repository.CanLoad(fp.pathEdit.Text()); ok && err == nil {
				// ロード系のみ履歴用キーを指定して履歴リストを保存
				mconfig.SaveUserConfig(fp.historyKey, path, 50)
			} else {
				// 読み込めない場合、拒否
				fp.pathEdit.ChangeText("")
			}
		}
	}
}

func (fp *FilePicker) onClickOpenButton() walk.EventHandler {
	return func() {
		if fp.pathEdit.Text() != "" {
			// ファイルパスからディレクトリパスを取得
			dirPath := filepath.Dir(fp.pathEdit.Text())
			// ファイルパスのディレクトリを初期パスとして設定
			fp.initialDirPath = dirPath
		} else if fp.historyKey != "" {
			// 履歴用キーを指定して履歴リストを取得
			choices := mconfig.LoadUserConfig(fp.historyKey)
			if len(choices) > 0 {
				// ファイルパスからディレクトリパスを取得
				dirPath := filepath.Dir(choices[0])
				// 履歴リストの先頭を初期パスとして設定
				fp.initialDirPath = dirPath
			}
		}

		// ファイル選択ダイアログを開く
		dlg := walk.FileDialog{
			Title: mi18n.T(
				"ファイル選択ダイアログタイトル",
				map[string]interface{}{"Title": fp.title}),
			Filter:         fp.convertFilterExtension(),
			FilterIndex:    1,
			InitialDirPath: fp.initialDirPath,
		}
		if ok, err := dlg.ShowOpen(nil); err != nil {
			walk.MsgBox(nil, mi18n.T("ファイル選択ダイアログ選択エラー"), err.Error(), walk.MsgBoxIconError)
		} else if ok {
			// パスを入力欄に設定
			fp.pathEdit.ChangeText(dlg.FilePath)
			// コールバックを呼び出し
			fp.onChanged(dlg.FilePath)
		}
	}
}

func (fp *FilePicker) convertFilterExtension() string {
	var filterString string
	for _, ext := range fp.filterExtensions {
		if filterString != "" {
			filterString = filterString + "|"
		}
		filterString = filterString + ext.description + "|" + ext.extension
	}
	return filterString
}

func (fp *FilePicker) newHistoryDialog() (*walk.Dialog, error) {
	// 履歴リストを取得
	choices := mconfig.LoadUserConfig(fp.historyKey)
	var dlg *walk.Dialog

	err := declarative.Dialog{
		AssignTo: &dlg,
		Title:    mi18n.T("履歴ダイアログタイトル", map[string]interface{}{"Title": fp.title}),
		Layout:   declarative.VBox{},
		Size:     declarative.Size{Width: 800, Height: 400},
		Children: []declarative.Widget{
			declarative.ListBox{
				AssignTo:     &fp.historyListBox,
				Model:        choices,
				MinSize:      declarative.Size{Width: 800, Height: 400},
				CurrentIndex: -1,
				OnItemActivated: func() {
					// 選択されたアイテムを取得
					index := fp.historyListBox.CurrentIndex()
					if index < 0 {
						return
					}
					item := choices[index]
					// パスを入力欄に設定
					fp.pathEdit.SetText(item)
					fp.historyDialog.Accept()
				},
			},
			declarative.Composite{
				Layout: declarative.HBox{},
				Children: []declarative.Widget{
					declarative.PushButton{
						Text: "OK",
						OnClicked: func() {
							// 選択されたアイテムを取得
							index := fp.historyListBox.CurrentIndex()
							if index < 0 {
								return
							}
							item := choices[index]
							// パスを入力欄に設定
							fp.pathEdit.SetText(item)
							fp.historyDialog.Accept()
						},
					},
					declarative.PushButton{
						Text: "Cancel",
						OnClicked: func() {
							// ダイアログを閉じる
							fp.historyDialog.Cancel()
						},
					},
				},
			},
		},
	}.Create(fp.pathEdit.Parent().Form())

	if err != nil {
		return nil, err
	}

	return dlg, nil
}
