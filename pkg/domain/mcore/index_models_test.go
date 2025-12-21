package mcore

import (
	"errors"
	"testing"

	"github.com/miu200521358/mlib_go/pkg/domain/merr"
)

// テスト用のモック実装
type mockIndexModel struct {
	index int
	valid bool
}

func (m *mockIndexModel) Index() int     { return m.index }
func (m *mockIndexModel) SetIndex(i int) { m.index = i }
func (m *mockIndexModel) IsValid() bool  { return m.valid }

func newMockIndexModel(index int, valid bool) *mockIndexModel {
	return &mockIndexModel{index: index, valid: valid}
}

func TestNewIndexModels(t *testing.T) {
	t.Run("正常: 長さ0", func(t *testing.T) {
		im, err := NewIndexModels[*mockIndexModel](0)
		if err != nil {
			t.Errorf("NewIndexModels() error = %v, want nil", err)
		}
		if im.Length() != 0 {
			t.Errorf("NewIndexModels() length = %v, want 0", im.Length())
		}
	})

	t.Run("正常: 長さ5", func(t *testing.T) {
		im, err := NewIndexModels[*mockIndexModel](5)
		if err != nil {
			t.Errorf("NewIndexModels() error = %v, want nil", err)
		}
		if im.Length() != 5 {
			t.Errorf("NewIndexModels() length = %v, want 5", im.Length())
		}
	})

	t.Run("エラー: 負の長さ", func(t *testing.T) {
		_, err := NewIndexModels[*mockIndexModel](-1)
		if err == nil {
			t.Errorf("NewIndexModels() error = nil, want error")
		}
		var target *merr.InvalidArgumentError
		if !errors.As(err, &target) {
			t.Errorf("NewIndexModels() error type mismatch")
		}
	})
}

func TestIndexModels_Get(t *testing.T) {
	im, _ := NewIndexModels[*mockIndexModel](3)
	model := newMockIndexModel(0, true)
	_ = im.Update(model)

	t.Run("正常: インデックス0", func(t *testing.T) {
		_, err := im.Get(0)
		if err != nil {
			t.Errorf("Get() error = %v, want nil", err)
		}
	})

	t.Run("エラー: 負のインデックス", func(t *testing.T) {
		_, err := im.Get(-1)
		if err == nil {
			t.Errorf("Get() error = nil, want error")
		}
		var target *merr.IndexOutOfRangeError
		if !errors.As(err, &target) {
			t.Errorf("Get() error type mismatch")
		}
	})

	t.Run("エラー: 範囲外のインデックス", func(t *testing.T) {
		_, err := im.Get(10)
		if err == nil {
			t.Errorf("Get() error = nil, want error")
		}
		var target *merr.IndexOutOfRangeError
		if !errors.As(err, &target) {
			t.Errorf("Get() error type mismatch")
		}
	})
}

func TestIndexModels_Update(t *testing.T) {
	t.Run("正常: インデックス0", func(t *testing.T) {
		im, _ := NewIndexModels[*mockIndexModel](3)
		err := im.Update(newMockIndexModel(0, true))
		if err != nil {
			t.Errorf("Update() error = %v, want nil", err)
		}
	})

	t.Run("エラー: 負のインデックス", func(t *testing.T) {
		im, _ := NewIndexModels[*mockIndexModel](3)
		err := im.Update(newMockIndexModel(-1, true))
		if err == nil {
			t.Errorf("Update() error = nil, want error")
		}
		var target *merr.IndexOutOfRangeError
		if !errors.As(err, &target) {
			t.Errorf("Update() error type mismatch")
		}
	})

	t.Run("エラー: 範囲外のインデックス", func(t *testing.T) {
		im, _ := NewIndexModels[*mockIndexModel](3)
		err := im.Update(newMockIndexModel(10, true))
		if err == nil {
			t.Errorf("Update() error = nil, want error")
		}
		var target *merr.IndexOutOfRangeError
		if !errors.As(err, &target) {
			t.Errorf("Update() error type mismatch")
		}
	})
}

func TestIndexModels_Append(t *testing.T) {
	t.Run("正常: 自動インデックス設定", func(t *testing.T) {
		im, _ := NewIndexModels[*mockIndexModel](0)
		err := im.Append(newMockIndexModel(-1, true))
		if err != nil {
			t.Errorf("Append() error = %v, want nil", err)
		}
	})

	t.Run("エラー: 既存インデックスへの追加", func(t *testing.T) {
		im, _ := NewIndexModels[*mockIndexModel](3)
		err := im.Append(newMockIndexModel(0, true))
		if err == nil {
			t.Errorf("Append() error = nil, want error")
		}
		var target *merr.InvalidOperationError
		if !errors.As(err, &target) {
			t.Errorf("Append() error type mismatch")
		}
	})
}

func TestIndexModels_Contains(t *testing.T) {
	im, _ := NewIndexModels[*mockIndexModel](3)
	_ = im.Update(newMockIndexModel(0, true))
	_ = im.Update(newMockIndexModel(1, false))

	t.Run("有効なモデル", func(t *testing.T) {
		if !im.Contains(0) {
			t.Errorf("Contains() = false, want true")
		}
	})

	t.Run("無効なモデル", func(t *testing.T) {
		if im.Contains(1) {
			t.Errorf("Contains() = true, want false")
		}
	})

	t.Run("範囲外", func(t *testing.T) {
		if im.Contains(10) {
			t.Errorf("Contains() = true, want false")
		}
	})
}

func TestIndexModels_ForEach(t *testing.T) {
	im, _ := NewIndexModels[*mockIndexModel](3)
	for i := 0; i < 3; i++ {
		_ = im.Update(newMockIndexModel(i, true))
	}

	t.Run("全要素をイテレート", func(t *testing.T) {
		count := 0
		im.ForEach(func(index int, value *mockIndexModel) bool {
			count++
			return true
		})
		if count != 3 {
			t.Errorf("ForEach() count = %v, want 3", count)
		}
	})

	t.Run("途中で中断", func(t *testing.T) {
		count := 0
		im.ForEach(func(index int, value *mockIndexModel) bool {
			count++
			return count < 2
		})
		if count != 2 {
			t.Errorf("ForEach() count = %v, want 2", count)
		}
	})
}

func TestIndexModels_FirstLast(t *testing.T) {
	t.Run("空のコレクション", func(t *testing.T) {
		im, _ := NewIndexModels[*mockIndexModel](0)
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
		im, _ := NewIndexModels[*mockIndexModel](3)
		first := newMockIndexModel(0, true)
		last := newMockIndexModel(2, true)
		_ = im.Update(first)
		_ = im.Update(last)

		result, err := im.First()
		if err != nil || result != first {
			t.Errorf("First() failed")
		}

		result, err = im.Last()
		if err != nil || result != last {
			t.Errorf("Last() failed")
		}
	})
}

func TestIndexModels_Clear(t *testing.T) {
	im, _ := NewIndexModels[*mockIndexModel](3)
	_ = im.Update(newMockIndexModel(0, true))

	im.Clear()

	if im.Length() != 0 || !im.IsEmpty() {
		t.Errorf("Clear() failed")
	}
}

func TestNewIndexModelsWithCapacity(t *testing.T) {
	t.Run("正常", func(t *testing.T) {
		im, err := NewIndexModelsWithCapacity[*mockIndexModel](3, 10)
		if err != nil {
			t.Errorf("NewIndexModelsWithCapacity() error = %v", err)
		}
		if im.Length() != 3 {
			t.Errorf("Length() = %v, want 3", im.Length())
		}
	})

	t.Run("エラー: 負の長さ", func(t *testing.T) {
		_, err := NewIndexModelsWithCapacity[*mockIndexModel](-1, 10)
		if err == nil {
			t.Errorf("NewIndexModelsWithCapacity() error = nil, want error")
		}
	})

	t.Run("エラー: capacity < length", func(t *testing.T) {
		_, err := NewIndexModelsWithCapacity[*mockIndexModel](10, 5)
		if err == nil {
			t.Errorf("NewIndexModelsWithCapacity() error = nil, want error")
		}
	})
}

func TestIndexModels_Remove(t *testing.T) {
	t.Run("正常", func(t *testing.T) {
		im, _ := NewIndexModels[*mockIndexModel](3)
		_ = im.Update(newMockIndexModel(0, true))
		_ = im.Update(newMockIndexModel(1, true))
		_ = im.Update(newMockIndexModel(2, true))

		err := im.Remove(1)
		if err != nil {
			t.Errorf("Remove() error = %v", err)
		}
		if im.Length() != 2 {
			t.Errorf("Length() = %v, want 2", im.Length())
		}
	})

	t.Run("エラー: 範囲外", func(t *testing.T) {
		im, _ := NewIndexModels[*mockIndexModel](3)
		err := im.Remove(10)
		if err == nil {
			t.Errorf("Remove() error = nil, want error")
		}
	})
}

func TestIndexModels_Values(t *testing.T) {
	im, _ := NewIndexModels[*mockIndexModel](3)
	_ = im.Update(newMockIndexModel(0, true))

	values := im.Values()
	if len(values) != 3 {
		t.Errorf("Values() length = %v, want 3", len(values))
	}
}

func TestIndexModels_IsNotEmpty(t *testing.T) {
	im, _ := NewIndexModels[*mockIndexModel](3)
	if !im.IsNotEmpty() {
		t.Errorf("IsNotEmpty() = false, want true")
	}

	im.Clear()
	if im.IsNotEmpty() {
		t.Errorf("IsNotEmpty() = true, want false")
	}
}
