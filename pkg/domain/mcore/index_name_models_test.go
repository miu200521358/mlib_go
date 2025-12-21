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

func TestNewIndexNameModelsWithCapacity(t *testing.T) {
	t.Run("正常", func(t *testing.T) {
		im, err := NewIndexNameModelsWithCapacity[*mockIndexNameModel](3, 10)
		if err != nil {
			t.Errorf("error = %v", err)
		}
		if im.Length() != 3 {
			t.Errorf("Length() = %v, want 3", im.Length())
		}
	})

	t.Run("エラー: 負の長さ", func(t *testing.T) {
		_, err := NewIndexNameModelsWithCapacity[*mockIndexNameModel](-1, 10)
		if err == nil {
			t.Errorf("error = nil, want error")
		}
	})

	t.Run("エラー: capacity < length", func(t *testing.T) {
		_, err := NewIndexNameModelsWithCapacity[*mockIndexNameModel](10, 5)
		if err == nil {
			t.Errorf("error = nil, want error")
		}
	})
}

func TestIndexNameModels_Get(t *testing.T) {
	im, _ := NewIndexNameModels[*mockIndexNameModel](3)
	_ = im.Update(newMockIndexNameModel(0, "センター", true))

	t.Run("正常", func(t *testing.T) {
		_, err := im.Get(0)
		if err != nil {
			t.Errorf("Get() error = %v", err)
		}
	})

	t.Run("エラー: 範囲外", func(t *testing.T) {
		_, err := im.Get(10)
		if err == nil {
			t.Errorf("Get() error = nil, want error")
		}
	})
}

func TestIndexNameModels_Append(t *testing.T) {
	t.Run("正常: 自動インデックス", func(t *testing.T) {
		im, _ := NewIndexNameModels[*mockIndexNameModel](0)
		err := im.Append(newMockIndexNameModel(-1, "test", true))
		if err != nil {
			t.Errorf("Append() error = %v", err)
		}
	})

	t.Run("エラー: 既存インデックス", func(t *testing.T) {
		im, _ := NewIndexNameModels[*mockIndexNameModel](3)
		err := im.Append(newMockIndexNameModel(0, "test", true))
		if err == nil {
			t.Errorf("Append() error = nil, want error")
		}
	})
}

func TestIndexNameModels_Remove(t *testing.T) {
	t.Run("正常", func(t *testing.T) {
		im, _ := NewIndexNameModels[*mockIndexNameModel](3)
		_ = im.Update(newMockIndexNameModel(0, "センター", true))
		err := im.Remove(0)
		if err != nil {
			t.Errorf("Remove() error = %v", err)
		}
	})

	t.Run("エラー: 範囲外", func(t *testing.T) {
		im, _ := NewIndexNameModels[*mockIndexNameModel](3)
		err := im.Remove(10)
		if err == nil {
			t.Errorf("Remove() error = nil, want error")
		}
	})
}

func TestIndexNameModels_ForEach(t *testing.T) {
	im, _ := NewIndexNameModels[*mockIndexNameModel](3)
	_ = im.Update(newMockIndexNameModel(0, "a", true))
	_ = im.Update(newMockIndexNameModel(1, "b", true))

	count := 0
	im.ForEach(func(index int, value *mockIndexNameModel) bool {
		count++
		return true
	})
	if count != 3 {
		t.Errorf("ForEach() count = %v, want 3", count)
	}
}

func TestIndexNameModels_Values(t *testing.T) {
	im, _ := NewIndexNameModels[*mockIndexNameModel](3)
	values := im.Values()
	if len(values) != 3 {
		t.Errorf("Values() length = %v, want 3", len(values))
	}
}

func TestIndexNameModels_Indexes(t *testing.T) {
	im, _ := NewIndexNameModels[*mockIndexNameModel](3)
	indexes := im.Indexes()
	if len(indexes) != 3 {
		t.Errorf("Indexes() length = %v, want 3", len(indexes))
	}
	if indexes[0] != 0 || indexes[1] != 1 || indexes[2] != 2 {
		t.Errorf("Indexes() values mismatch")
	}
}

func TestIndexNameModels_FirstLast(t *testing.T) {
	t.Run("空のコレクション", func(t *testing.T) {
		im, _ := NewIndexNameModels[*mockIndexNameModel](0)
		_, err := im.First()
		if err == nil {
			t.Errorf("First() error = nil, want error")
		}
		_, err = im.Last()
		if err == nil {
			t.Errorf("Last() error = nil, want error")
		}
	})

	t.Run("要素があるコレクション", func(t *testing.T) {
		im, _ := NewIndexNameModels[*mockIndexNameModel](3)
		first := newMockIndexNameModel(0, "first", true)
		last := newMockIndexNameModel(2, "last", true)
		_ = im.Update(first)
		_ = im.Update(last)

		r, _ := im.First()
		if r != first {
			t.Errorf("First() failed")
		}
		r, _ = im.Last()
		if r != last {
			t.Errorf("Last() failed")
		}
	})
}

func TestIndexNameModels_IsEmptyNotEmpty(t *testing.T) {
	im, _ := NewIndexNameModels[*mockIndexNameModel](3)
	if im.IsEmpty() {
		t.Errorf("IsEmpty() = true, want false")
	}
	if !im.IsNotEmpty() {
		t.Errorf("IsNotEmpty() = false, want true")
	}

	im.Clear()
	if !im.IsEmpty() {
		t.Errorf("IsEmpty() = false, want true")
	}
	if im.IsNotEmpty() {
		t.Errorf("IsNotEmpty() = true, want false")
	}
}

func TestIndexNameModels_Contains(t *testing.T) {
	im, _ := NewIndexNameModels[*mockIndexNameModel](3)
	_ = im.Update(newMockIndexNameModel(0, "センター", true))
	_ = im.Update(newMockIndexNameModel(1, "無効", false))

	if !im.Contains(0) {
		t.Errorf("Contains(0) = false, want true")
	}
	if im.Contains(1) {
		t.Errorf("Contains(1) = true, want false (invalid)")
	}
	if im.Contains(10) {
		t.Errorf("Contains(10) = true, want false (out of range)")
	}
}

func TestIndexNameModels_Update_RangeError(t *testing.T) {
	im, _ := NewIndexNameModels[*mockIndexNameModel](3)
	err := im.Update(newMockIndexNameModel(10, "test", true))
	if err == nil {
		t.Errorf("Update() error = nil, want error")
	}
}
