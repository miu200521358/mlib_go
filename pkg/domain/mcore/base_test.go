package mcore

import "testing"

func TestIndexModel_New(t *testing.T) {
	t.Run("正常: インデックス0", func(t *testing.T) {
		m := NewIndexModel(0)
		if m.Index() != 0 {
			t.Errorf("Index() = %v, want 0", m.Index())
		}
		if !m.IsValid() {
			t.Errorf("IsValid() = false, want true")
		}
	})

	t.Run("正常: インデックス5", func(t *testing.T) {
		m := NewIndexModel(5)
		if m.Index() != 5 {
			t.Errorf("Index() = %v, want 5", m.Index())
		}
	})

	t.Run("無効: インデックス-1", func(t *testing.T) {
		m := NewIndexModel(-1)
		if m.IsValid() {
			t.Errorf("IsValid() = true, want false")
		}
	})
}

func TestIndexModel_SetIndex(t *testing.T) {
	m := NewIndexModel(0)
	m.SetIndex(10)
	if m.Index() != 10 {
		t.Errorf("Index() = %v, want 10", m.Index())
	}
}

func TestIndexNameModel_New(t *testing.T) {
	t.Run("正常: 名前付き", func(t *testing.T) {
		m := NewIndexNameModel(0, "センター", "center")
		if m.Index() != 0 {
			t.Errorf("Index() = %v, want 0", m.Index())
		}
		if m.Name() != "センター" {
			t.Errorf("Name() = %v, want センター", m.Name())
		}
		if m.EnglishName() != "center" {
			t.Errorf("EnglishName() = %v, want center", m.EnglishName())
		}
		if !m.IsValid() {
			t.Errorf("IsValid() = false, want true")
		}
	})

	t.Run("無効: インデックス-1", func(t *testing.T) {
		m := NewIndexNameModel(-1, "test", "test")
		if m.IsValid() {
			t.Errorf("IsValid() = true, want false")
		}
	})
}

func TestIndexNameModel_SetName(t *testing.T) {
	m := NewIndexNameModel(0, "old", "old_en")
	m.SetName("new")
	m.SetEnglishName("new_en")

	if m.Name() != "new" {
		t.Errorf("Name() = %v, want new", m.Name())
	}
	if m.EnglishName() != "new_en" {
		t.Errorf("EnglishName() = %v, want new_en", m.EnglishName())
	}
}

func TestIndexModel_ImplementsInterface(t *testing.T) {
	// IIndexModelインターフェースを満たすことを確認
	var _ IIndexModel = (*IndexModel)(nil)
}

func TestIndexNameModel_ImplementsInterface(t *testing.T) {
	// IIndexNameModelインターフェースを満たすことを確認
	var _ IIndexNameModel = (*IndexNameModel)(nil)
}

// 埋め込みテスト: 派生structでのオーバーライド
type testDerivedModel struct {
	IndexNameModel
	data string
}

func (m *testDerivedModel) IsValid() bool {
	return m.IndexNameModel.IsValid() && m.data != ""
}

func TestDerivedModel_Override(t *testing.T) {
	t.Run("派生: 有効", func(t *testing.T) {
		m := &testDerivedModel{
			IndexNameModel: *NewIndexNameModel(0, "test", "test"),
			data:           "some data",
		}
		if !m.IsValid() {
			t.Errorf("IsValid() = false, want true")
		}
	})

	t.Run("派生: dataが空で無効", func(t *testing.T) {
		m := &testDerivedModel{
			IndexNameModel: *NewIndexNameModel(0, "test", "test"),
			data:           "",
		}
		if m.IsValid() {
			t.Errorf("IsValid() = true, want false")
		}
	})

	t.Run("派生: インデックスが負で無効", func(t *testing.T) {
		m := &testDerivedModel{
			IndexNameModel: *NewIndexNameModel(-1, "test", "test"),
			data:           "some data",
		}
		if m.IsValid() {
			t.Errorf("IsValid() = true, want false")
		}
	})
}
