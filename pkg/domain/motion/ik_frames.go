// 指示: miu200521358
package motion

// IkEnabledFrame はIK有効フレームを表す。
type IkEnabledFrame struct {
	*BaseFrame
	BoneName string
	Enabled  bool
}

// NewIkEnabledFrame はIkEnabledFrameを生成する。
func NewIkEnabledFrame(index Frame, boneName string) *IkEnabledFrame {
	return &IkEnabledFrame{BaseFrame: NewBaseFrame(index), BoneName: boneName, Enabled: true}
}

// Copy はフレームを複製する。
func (f *IkEnabledFrame) Copy() (IkEnabledFrame, error) {
	if f == nil {
		return IkEnabledFrame{}, nil
	}
	copied := IkEnabledFrame{
		BaseFrame: &BaseFrame{index: f.Index(), Read: f.Read},
		BoneName:  f.BoneName,
		Enabled:   f.Enabled,
	}
	return copied, nil
}

// lerpFrame は補間せず前フレームを複製する。
func (next *IkEnabledFrame) lerpFrame(prev *IkEnabledFrame, index Frame) *IkEnabledFrame {
	if prev == nil && next == nil {
		return NewIkEnabledFrame(index, "")
	}
	if prev == nil {
		return next.copyWithIndex(index)
	}
	return prev.copyWithIndex(index)
}

// splitCurve は何もしない。
func (f *IkEnabledFrame) splitCurve(prev *IkEnabledFrame, next *IkEnabledFrame, index Frame) {
}

// copyWithIndex は指定フレーム番号で複製する。
func (f *IkEnabledFrame) copyWithIndex(index Frame) *IkEnabledFrame {
	if f == nil {
		return nil
	}
	return &IkEnabledFrame{
		BaseFrame: &BaseFrame{index: index, Read: f.Read},
		BoneName:  f.BoneName,
		Enabled:   f.Enabled,
	}
}

// IkFrame はIKフレームを表す。
type IkFrame struct {
	*BaseFrame
	Visible bool
	IkList  []*IkEnabledFrame
}

// NewIkFrame はIkFrameを生成する。
func NewIkFrame(index Frame) *IkFrame {
	return &IkFrame{BaseFrame: NewBaseFrame(index), Visible: true, IkList: make([]*IkEnabledFrame, 0)}
}

// Copy はフレームを複製する。
func (f *IkFrame) Copy() (IkFrame, error) {
	if f == nil {
		return IkFrame{}, nil
	}
	copied := IkFrame{
		BaseFrame: &BaseFrame{index: f.Index(), Read: f.Read},
		Visible:   f.Visible,
		IkList:    copyIkList(f.IkList),
	}
	return copied, nil
}

// lerpFrame は補間せず前フレームを複製する。
func (next *IkFrame) lerpFrame(prev *IkFrame, index Frame) *IkFrame {
	if prev == nil && next == nil {
		return NewIkFrame(index)
	}
	if prev == nil {
		return next.copyWithIndex(index)
	}
	return prev.copyWithIndex(index)
}

// splitCurve は何もしない。
func (f *IkFrame) splitCurve(prev *IkFrame, next *IkFrame, index Frame) {
}

// copyWithIndex は指定フレーム番号で複製する。
func (f *IkFrame) copyWithIndex(index Frame) *IkFrame {
	if f == nil {
		return nil
	}
	return &IkFrame{
		BaseFrame: &BaseFrame{index: index, Read: f.Read},
		Visible:   f.Visible,
		IkList:    copyIkList(f.IkList),
	}
}

// IsEnable はIKの有効/無効を返す。
func (f *IkFrame) IsEnable(boneName string) bool {
	if f == nil {
		return true
	}
	for _, item := range f.IkList {
		if item != nil && item.BoneName == boneName {
			return item.Enabled
		}
	}
	return true
}

// IsDefault は既定値か判定する。
func (f *IkFrame) IsDefault() bool {
	if f == nil {
		return true
	}
	return f.Visible && len(f.IkList) == 0
}

// IkFrames はIKフレーム集合を表す。
type IkFrames struct {
	*BaseFrames[*IkFrame]
}

// NewIkFrames はIkFramesを生成する。
func NewIkFrames() *IkFrames {
	return &IkFrames{BaseFrames: NewBaseFrames(NewIkFrame, nilIkFrame)}
}

// Clean は既定値のみのフレームを削除する。
func (i *IkFrames) Clean() {
	if i == nil || i.Len() != 1 {
		return
	}
	var frameIndex Frame
	var frameValue *IkFrame
	i.ForEach(func(idx Frame, value *IkFrame) bool {
		frameIndex = idx
		frameValue = value
		return false
	})
	if frameValue == nil || frameValue.IsDefault() {
		i.Delete(frameIndex)
	}
}

// Copy はフレーム集合を複製する。
func (i *IkFrames) Copy() (IkFrames, error) {
	if i == nil {
		return IkFrames{}, nil
	}
	return deepCopy(*i)
}

// nilIkFrame は既定の空フレームを返す。
func nilIkFrame() *IkFrame {
	return nil
}

// copyIkList はIK有効リストを複製する。
func copyIkList(src []*IkEnabledFrame) []*IkEnabledFrame {
	if len(src) == 0 {
		return nil
	}
	out := make([]*IkEnabledFrame, 0, len(src))
	for _, item := range src {
		if item == nil {
			out = append(out, nil)
			continue
		}
		copied := IkEnabledFrame{
			BaseFrame: &BaseFrame{index: item.Index(), Read: item.Read},
			BoneName:  item.BoneName,
			Enabled:   item.Enabled,
		}
		out = append(out, &copied)
	}
	return out
}
