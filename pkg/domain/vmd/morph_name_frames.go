package vmd

import "github.com/miu200521358/mlib_go/pkg/domain/mmath"

type MorphNameFrames struct {
	*BaseFrames[*MorphFrame]
	Name string // モーフ名
}

func NewMorphNameFrames(name string) *MorphNameFrames {
	return &MorphNameFrames{
		BaseFrames: NewBaseFrames[*MorphFrame](NewMorphFrame, nilMorphFrame),
		Name:       name,
	}
}

func nilMorphFrame() *MorphFrame {
	return nil
}

// ContainsActive 有効なキーフレが存在するか
func (morphNameFrames *MorphNameFrames) ContainsActive() bool {
	if morphNameFrames.Length() == 0 {
		return false
	}

	isActive := false
	morphNameFrames.ForEach(func(index float32, bf *MorphFrame) bool {
		if bf == nil {
			return true
		}

		if !mmath.NearEquals(bf.Ratio, 0.0, 1e-2) {
			isActive = true
			return false
		}

		nextBf := morphNameFrames.Get(morphNameFrames.NextFrame(bf.Index()))

		if nextBf == nil {
			return true
		}

		if !mmath.NearEquals(bf.Ratio, 0.0, 1e-2) {
			isActive = true
			return false
		}

		return true
	})

	return isActive
}
