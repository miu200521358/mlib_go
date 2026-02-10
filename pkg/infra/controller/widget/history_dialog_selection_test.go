// 指示: miu200521358
package widget

import "testing"

func TestResolveHistoryDialogSelectionConfirmResolvesPath(t *testing.T) {
	called := false
	gotIndex := -1
	pathResolver := func(idx int) (string, bool) {
		called = true
		gotIndex = idx
		if idx != 1 {
			return "", false
		}
		return "selected-path", true
	}

	path, ok := resolveHistoryDialogSelection(historyDialogActionConfirm, 1, pathResolver)
	if !ok {
		t.Fatalf("確認操作でパス解決に失敗しました")
	}
	if path != "selected-path" {
		t.Fatalf("解決パスが不正です: got=%s", path)
	}
	if !called {
		t.Fatalf("パス解決関数が呼び出されていません")
	}
	if gotIndex != 1 {
		t.Fatalf("選択インデックスが不正です: got=%d", gotIndex)
	}
}

func TestResolveHistoryDialogSelectionCancelDoesNotResolvePath(t *testing.T) {
	called := false
	pathResolver := func(idx int) (string, bool) {
		called = true
		return "selected-path", true
	}

	path, ok := resolveHistoryDialogSelection(historyDialogActionCancel, 0, pathResolver)
	if ok {
		t.Fatalf("キャンセル操作ではパスを解決してはいけません")
	}
	if path != "" {
		t.Fatalf("キャンセル操作時は空パスであるべきです: got=%s", path)
	}
	if called {
		t.Fatalf("キャンセル操作でパス解決関数が呼び出されています")
	}
}

func TestResolveHistoryDialogSelectionNilResolverReturnsFalse(t *testing.T) {
	path, ok := resolveHistoryDialogSelection(historyDialogActionConfirm, 0, nil)
	if ok {
		t.Fatalf("resolver が nil の場合は失敗であるべきです")
	}
	if path != "" {
		t.Fatalf("resolver が nil の場合は空パスであるべきです: got=%s", path)
	}
}
