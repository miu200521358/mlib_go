// 指示: miu200521358
package widget

// historyDialogAction は履歴ダイアログで実行された操作種別を表す。
type historyDialogAction int

const (
	historyDialogActionConfirm historyDialogAction = iota
	historyDialogActionCancel
)

// resolveHistoryDialogSelection は操作種別に応じて選択中の履歴パスを解決する。
func resolveHistoryDialogSelection(action historyDialogAction, selectedIndex int, pathResolver func(int) (string, bool)) (string, bool) {
	if action != historyDialogActionConfirm || pathResolver == nil {
		return "", false
	}
	return pathResolver(selectedIndex)
}
