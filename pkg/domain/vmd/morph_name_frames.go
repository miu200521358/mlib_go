package vmd

import "github.com/miu200521358/mlib_go/pkg/domain/mmath"

type MorphNameFrames struct {
	*BaseFrames[*MorphFrame]
	Name string // モーフ名
}

func NewMorphNameFrames(name string) *MorphNameFrames {
	return &MorphNameFrames{
		BaseFrames: NewBaseFrames[*MorphFrame](),
		Name:       name,
	}
}

// ContainsActive 有効なキーフレが存在するか
func (morphNameFrames *MorphNameFrames) ContainsActive() bool {
	morphNameFrames.lock.RLock()
	defer morphNameFrames.lock.RUnlock()

	if morphNameFrames.Length() == 0 {
		return false
	}

	for mf := range morphNameFrames.Iterator() {
		if mf == nil {
			return false
		}

		if !mmath.NearEquals(mf.Ratio, 0.0, 1e-2) {
			return true
		}

		nextMf := morphNameFrames.Get(morphNameFrames.NextFrame(mf.Index()))

		if nextMf == nil {
			return false
		}

		if !mmath.NearEquals(nextMf.Ratio, 0.0, 1e-2) {
			return true
		}
	}

	return false
}
