// 指示: miu200521358
package motion

import "testing"

// dummyFrame はBaseFrames検証用のフレーム。
type dummyFrame struct {
	*BaseFrame
	Value      int
	LerpCount  *int
	SplitCount *int
}

// newDummyFrame はdummyFrameを生成する。
func newDummyFrame(index Frame) *dummyFrame {
	return &dummyFrame{BaseFrame: NewBaseFrame(index)}
}

// Copy はフレームを複製する。
func (d *dummyFrame) Copy() (dummyFrame, error) {
	if d == nil {
		return dummyFrame{}, nil
	}
	copied := dummyFrame{
		BaseFrame:  &BaseFrame{index: d.Index(), Read: d.Read},
		Value:      d.Value,
		LerpCount:  d.LerpCount,
		SplitCount: d.SplitCount,
	}
	return copied, nil
}

// copyWithIndex は指定フレーム番号で複製する。
func (d *dummyFrame) copyWithIndex(index Frame) *dummyFrame {
	if d == nil {
		return nil
	}
	return &dummyFrame{
		BaseFrame:  &BaseFrame{index: index, Read: d.Read},
		Value:      d.Value,
		LerpCount:  d.LerpCount,
		SplitCount: d.SplitCount,
	}
}

// lerpFrame は補間呼び出しを記録する。
func (next *dummyFrame) lerpFrame(prev *dummyFrame, index Frame) *dummyFrame {
	if next != nil && next.LerpCount != nil {
		*next.LerpCount = *next.LerpCount + 1
	}
	out := newDummyFrame(index)
	if prev != nil {
		out.Value = prev.Value + 100
	}
	return out
}

// splitCurve は分割呼び出しを記録する。
func (f *dummyFrame) splitCurve(prev *dummyFrame, next *dummyFrame, index Frame) {
	if f != nil && f.SplitCount != nil {
		*f.SplitCount++
	}
}

// TestBaseFramesGetPaths はGetの分岐を確認する。
func TestBaseFramesGetPaths(t *testing.T) {
	lerpCount := 0
	frames := NewBaseFrames(newDummyFrame, func() *dummyFrame { return nil })

	// empty -> newFunc
	got := frames.Get(2)
	if got == nil || got.Index() != 2 {
		t.Fatalf("Get empty: got=%v", got)
	}

	f0 := newDummyFrame(0)
	f0.Value = 1
	f0.LerpCount = &lerpCount
	f10 := newDummyFrame(10)
	f10.Value = 2
	f10.LerpCount = &lerpCount
	frames.Append(f0)
	frames.Append(f10)

	// existing
	got = frames.Get(0)
	if got == nil || got.Index() != 0 || got.Value != 1 {
		t.Fatalf("Get existing: got=%v", got)
	}

	// between -> lerp
	got = frames.Get(5)
	if got == nil || got.Value != 101 {
		t.Fatalf("Get lerp: got=%v", got)
	}
	if lerpCount == 0 {
		t.Fatalf("lerp not called")
	}

	// no next -> copy max
	got = frames.Get(20)
	if got == nil || got.Index() != 20 || got.Value != 2 {
		t.Fatalf("Get no next: got=%v", got)
	}
}

// TestBaseFramesInsertSplit はInsert時の分割を確認する。
func TestBaseFramesInsertSplit(t *testing.T) {
	splitCount := 0
	frames := NewBaseFrames(newDummyFrame, func() *dummyFrame { return nil })

	f0 := newDummyFrame(0)
	f10 := newDummyFrame(10)
	frames.Append(f0)
	frames.Append(f10)

	f5 := newDummyFrame(5)
	f5.SplitCount = &splitCount
	frames.Insert(f5)

	if splitCount == 0 {
		t.Fatalf("split not called")
	}
}

// TestBaseFramesForEachOrder はForEachの順序を確認する。
func TestBaseFramesForEachOrder(t *testing.T) {
	frames := NewBaseFrames(newDummyFrame, func() *dummyFrame { return nil })
	frames.Append(newDummyFrame(5))
	frames.Append(newDummyFrame(1))
	frames.Append(newDummyFrame(3))

	seen := make([]Frame, 0)
	frames.ForEach(func(frame Frame, _ *dummyFrame) bool {
		seen = append(seen, frame)
		return true
	})
	if len(seen) != 3 || seen[0] != 1 || seen[1] != 3 || seen[2] != 5 {
		t.Fatalf("ForEach order: got=%v", seen)
	}
}

// TestBaseFramesDelete はDeleteの挙動を確認する。
func TestBaseFramesDelete(t *testing.T) {
	frames := NewBaseFrames(newDummyFrame, func() *dummyFrame { return nil })
	frames.Append(newDummyFrame(1))
	frames.Delete(1)
	if frames.Has(1) {
		t.Fatalf("Delete failed")
	}
}
