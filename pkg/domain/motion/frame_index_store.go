// 指示: miu200521358
package motion

import "sort"

// IFrameIndexStore はフレーム番号の索引管理を表す。
type IFrameIndexStore interface {
	Upsert(frame Frame)
	Delete(frame Frame)
	Has(frame Frame) bool
	Prev(frame Frame) (Frame, bool)
	Next(frame Frame) (Frame, bool)
	Max() (Frame, bool)
	Min() (Frame, bool)
	Len() int
	ForEach(fn func(Frame) bool)
	IsDirty() bool
	Finalize()
}

// SortedFrameIndexStore はソート配列でフレーム番号を管理する。
type SortedFrameIndexStore struct {
	indexes  []Frame
	indexSet map[Frame]struct{}
	dirty    bool
}

// NewSortedFrameIndexStore はSortedFrameIndexStoreを生成する。
func NewSortedFrameIndexStore() *SortedFrameIndexStore {
	return &SortedFrameIndexStore{
		indexes:  make([]Frame, 0),
		indexSet: make(map[Frame]struct{}),
	}
}

// Upsert はフレーム番号を追加または更新する。
func (s *SortedFrameIndexStore) Upsert(frame Frame) {
	if s == nil {
		return
	}
	if _, exists := s.indexSet[frame]; exists {
		return
	}
	s.indexSet[frame] = struct{}{}
	s.indexes = append(s.indexes, frame)
	s.dirty = true
}

// Delete はフレーム番号を削除する。
func (s *SortedFrameIndexStore) Delete(frame Frame) {
	if s == nil {
		return
	}
	if _, exists := s.indexSet[frame]; !exists {
		return
	}
	delete(s.indexSet, frame)
	for i, v := range s.indexes {
		if v == frame {
			s.indexes = append(s.indexes[:i], s.indexes[i+1:]...)
			s.dirty = true
			return
		}
	}
}

// Has はフレーム番号の存在を判定する。
func (s *SortedFrameIndexStore) Has(frame Frame) bool {
	if s == nil {
		return false
	}
	_, exists := s.indexSet[frame]
	return exists
}

// Prev は指定フレーム未満の最大値を返す。
func (s *SortedFrameIndexStore) Prev(frame Frame) (Frame, bool) {
	if s == nil || len(s.indexes) == 0 {
		return 0, false
	}
	s.Finalize()
	pos := sort.Search(len(s.indexes), func(i int) bool { return s.indexes[i] >= frame })
	if pos == 0 {
		return s.indexes[0], false
	}
	return s.indexes[pos-1], true
}

// Next は指定フレームより大きい最小値を返す。
func (s *SortedFrameIndexStore) Next(frame Frame) (Frame, bool) {
	if s == nil || len(s.indexes) == 0 {
		return frame, false
	}
	s.Finalize()
	pos := sort.Search(len(s.indexes), func(i int) bool { return s.indexes[i] > frame })
	if pos >= len(s.indexes) {
		return frame, false
	}
	return s.indexes[pos], true
}

// Max は最大のフレーム番号を返す。
func (s *SortedFrameIndexStore) Max() (Frame, bool) {
	if s == nil || len(s.indexes) == 0 {
		return 0, false
	}
	s.Finalize()
	return s.indexes[len(s.indexes)-1], true
}

// Min は最小のフレーム番号を返す。
func (s *SortedFrameIndexStore) Min() (Frame, bool) {
	if s == nil || len(s.indexes) == 0 {
		return 0, false
	}
	s.Finalize()
	return s.indexes[0], true
}

// Len は登録数を返す。
func (s *SortedFrameIndexStore) Len() int {
	if s == nil {
		return 0
	}
	return len(s.indexes)
}

// ForEach は昇順で走査する。
func (s *SortedFrameIndexStore) ForEach(fn func(Frame) bool) {
	if s == nil || fn == nil || len(s.indexes) == 0 {
		return
	}
	s.Finalize()
	for _, v := range s.indexes {
		if !fn(v) {
			return
		}
	}
}

// IsDirty は並び替えが必要か返す。
func (s *SortedFrameIndexStore) IsDirty() bool {
	if s == nil {
		return false
	}
	return s.dirty
}

// Finalize は昇順化と重複排除を確定する。
func (s *SortedFrameIndexStore) Finalize() {
	if s == nil || !s.dirty {
		return
	}
	if len(s.indexes) == 0 {
		s.dirty = false
		return
	}
	sort.Slice(s.indexes, func(i, j int) bool { return s.indexes[i] < s.indexes[j] })
	out := s.indexes[:0]
	var prev Frame
	for i, v := range s.indexes {
		if i == 0 || v != prev {
			out = append(out, v)
			prev = v
		}
	}
	s.indexes = out
	s.dirty = false
}
