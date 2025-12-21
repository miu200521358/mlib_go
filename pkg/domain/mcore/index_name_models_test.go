package mcore

import (
	"errors"
	"testing"

	"github.com/miu200521358/mlib_go/pkg/domain/merr"
)

// テスト用のモック実装
type mockIndexNameModel struct {
	index       int
	name        string
	englishName string
	valid       bool
}

func (m *mockIndexNameModel) Index() int                        { return m.index }
func (m *mockIndexNameModel) SetIndex(i int)                    { m.index = i }
func (m *mockIndexNameModel) IsValid() bool                     { return m.valid }
func (m *mockIndexNameModel) Name() string                      { return m.name }
func (m *mockIndexNameModel) SetName(name string)               { m.name = name }
func (m *mockIndexNameModel) EnglishName() string               { return m.englishName }
func (m *mockIndexNameModel) SetEnglishName(englishName string) { m.englishName = englishName }

func newMockIndexNameModel(index int, name string, valid bool) *mockIndexNameModel {
	return &mockIndexNameModel{
		index: index,
		name:  name,
		valid: valid,
	}
}

func TestNewIndexNameModels(t *testing.T) {
	t.Run("正常: 長さ0", func(t *testing.T) {
		im, err := NewIndexNameModels[*mockIndexNameModel](0)
		if err != nil {
			t.Errorf("NewIndexNameModels() error = %v, want nil", err)
		}
		if im.Length() != 0 {
			t.Errorf("NewIndexNameModels() length = %v, want 0", im.Length())
		}
	})

	t.Run("エラー: 負の長さ", func(t *testing.T) {
		_, err := NewIndexNameModels[*mockIndexNameModel](-1)
		if err == nil {
			t.Errorf("NewIndexNameModels() error = nil, want error")
		}
		var target *merr.InvalidArgumentError
		if !errors.As(err, &target) {
			t.Errorf("NewIndexNameModels() error type mismatch")
		}
	})
}

func TestIndexNameModels_GetByName(t *testing.T) {
	im, _ := NewIndexNameModels[*mockIndexNameModel](3)
	_ = im.Update(newMockIndexNameModel(0, "センター", true))

	t.Run("正常: 存在する名前", func(t *testing.T) {
		_, err := im.GetByName("センター")
		if err != nil {
			t.Errorf("GetByName() error = %v, want nil", err)
		}
	})

	t.Run("エラー: 存在しない名前", func(t *testing.T) {
		_, err := im.GetByName("不明")
		if err == nil {
			t.Errorf("GetByName() error = nil, want error")
		}
		var target *merr.NameNotFoundError
		if !errors.As(err, &target) {
			t.Errorf("GetByName() error type mismatch")
		}
	})
}

func TestIndexNameModels_ContainsByName(t *testing.T) {
	im, _ := NewIndexNameModels[*mockIndexNameModel](3)
	_ = im.Update(newMockIndexNameModel(0, "センター", true))

	t.Run("存在する名前", func(t *testing.T) {
		if !im.ContainsByName("センター") {
			t.Errorf("ContainsByName() = false, want true")
		}
	})

	t.Run("存在しない名前", func(t *testing.T) {
		if im.ContainsByName("不明") {
			t.Errorf("ContainsByName() = true, want false")
		}
	})
}

func TestIndexNameModels_RemoveByName(t *testing.T) {
	t.Run("正常: 存在する名前を削除", func(t *testing.T) {
		im, _ := NewIndexNameModels[*mockIndexNameModel](3)
		_ = im.Update(newMockIndexNameModel(0, "センター", true))
		err := im.RemoveByName("センター")
		if err != nil {
			t.Errorf("RemoveByName() error = %v, want nil", err)
		}
		if im.ContainsByName("センター") {
			t.Errorf("RemoveByName() name still exists")
		}
	})

	t.Run("エラー: 存在しない名前を削除", func(t *testing.T) {
		im, _ := NewIndexNameModels[*mockIndexNameModel](3)
		err := im.RemoveByName("不明")
		if err == nil {
			t.Errorf("RemoveByName() error = nil, want error")
		}
		var target *merr.NameNotFoundError
		if !errors.As(err, &target) {
			t.Errorf("RemoveByName() error type mismatch")
		}
	})
}

func TestIndexNameModels_Names(t *testing.T) {
	im, _ := NewIndexNameModels[*mockIndexNameModel](3)
	_ = im.Update(newMockIndexNameModel(0, "センター", true))
	_ = im.Update(newMockIndexNameModel(1, "上半身", true))

	names := im.Names()
	if len(names) != 3 {
		t.Errorf("Names() length = %v, want 3", len(names))
	}
	if names[0] != "センター" {
		t.Errorf("Names()[0] = %v, want センター", names[0])
	}
	if names[1] != "上半身" {
		t.Errorf("Names()[1] = %v, want 上半身", names[1])
	}
}

func TestIndexNameModels_UpdateNameIndexes(t *testing.T) {
	im, _ := NewIndexNameModels[*mockIndexNameModel](3)
	_ = im.Update(newMockIndexNameModel(0, "センター", true))
	_ = im.Update(newMockIndexNameModel(1, "上半身", true))

	im.nameIndexes = make(map[string]int)
	im.UpdateNameIndexes()

	if !im.ContainsByName("センター") {
		t.Errorf("UpdateNameIndexes() センター not found")
	}
	if !im.ContainsByName("上半身") {
		t.Errorf("UpdateNameIndexes() 上半身 not found")
	}
}

func TestIndexNameModels_NameFirstWins(t *testing.T) {
	im, _ := NewIndexNameModels[*mockIndexNameModel](3)
	_ = im.Update(newMockIndexNameModel(0, "センター", true))
	_ = im.Update(newMockIndexNameModel(1, "センター", true))

	result, _ := im.GetByName("センター")
	if result.Index() != 0 {
		t.Errorf("GetByName() first wins: index = %v, want 0", result.Index())
	}
}

func TestIndexNameModels_Clear(t *testing.T) {
	im, _ := NewIndexNameModels[*mockIndexNameModel](3)
	_ = im.Update(newMockIndexNameModel(0, "センター", true))

	im.Clear()

	if im.Length() != 0 || im.ContainsByName("センター") {
		t.Errorf("Clear() failed")
	}
}
