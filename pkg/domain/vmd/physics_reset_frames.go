package vmd

import "github.com/tiendc/go-deepcopy"

// 物理リセット方法
type PhysicsResetType int

const (
	PHYSICS_RESET_TYPE_NONE            PhysicsResetType = 0 // 物理リセットなし
	PHYSICS_RESET_TYPE_CONTINUE_FRAME  PhysicsResetType = 1 // 継続フレーム用物理リセット（前フレームの物理を継続）
	PHYSICS_RESET_TYPE_START_FRAME     PhysicsResetType = 2 // 開始フレーム用物理リセット
	PHYSICS_RESET_TYPE_START_FIT_FRAME PhysicsResetType = 3 // 開始フレーム用物理リセット（Yスタンスから開始）
)

type PhysicsResetFrame struct {
	*BaseFrame                        // キーフレ
	PhysicsResetType PhysicsResetType // 物理リセット種別
}

func NewPhysicsResetFrame(index float32) *PhysicsResetFrame {
	return &PhysicsResetFrame{
		BaseFrame:        NewFrame(index).(*BaseFrame),
		PhysicsResetType: PHYSICS_RESET_TYPE_NONE, // デフォルト値
	}
}

func NewPhysicsResetFrameByValue(index float32, physicsReset PhysicsResetType) *PhysicsResetFrame {
	return &PhysicsResetFrame{
		BaseFrame:        NewFrame(index).(*BaseFrame),
		PhysicsResetType: physicsReset,
	}
}

func (nextMf *PhysicsResetFrame) Copy() IBaseFrame {
	vv := &PhysicsResetFrame{
		BaseFrame:        nextMf.BaseFrame.Copy().(*BaseFrame),
		PhysicsResetType: nextMf.PhysicsResetType,
	}
	return vv
}

func (nextMf *PhysicsResetFrame) lerpFrame(prevFrame IBaseFrame, index float32) IBaseFrame {
	prevMf := prevFrame.(*PhysicsResetFrame)
	// 補間なしで前のキーフレを引き継ぐ
	vv := &PhysicsResetFrame{
		BaseFrame:        NewFrame(index).(*BaseFrame),
		PhysicsResetType: prevMf.PhysicsResetType,
	}
	return vv
}

func (nextMf *PhysicsResetFrame) splitCurve(prevFrame IBaseFrame, nextFrame IBaseFrame, index float32) {
}

type PhysicsResetFrames struct {
	*BaseFrames[*PhysicsResetFrame]
}

func NewPhysicsResetFrames() *PhysicsResetFrames {
	return &PhysicsResetFrames{
		BaseFrames: NewBaseFrames(NewPhysicsResetFrame, nilPhysicsResetFrame),
	}
}

func nilPhysicsResetFrame() *PhysicsResetFrame {
	return nil
}

func (mf *PhysicsResetFrames) Copy() (*PhysicsResetFrames, error) {
	copied := new(PhysicsResetFrames)
	err := deepcopy.Copy(copied, mf)
	return copied, err
}

func (mf *PhysicsResetFrames) Get(frame float32) *PhysicsResetFrame {
	if mf.values.Len() == 0 {
		return &PhysicsResetFrame{
			BaseFrame:        NewFrame(frame).(*BaseFrame),
			PhysicsResetType: PHYSICS_RESET_TYPE_NONE, // デフォルト値
		}
	}

	if mf.values.Has(frame) {
		// キーフレが登録済みの場合、値を返す
		return mf.values.Get(NewFrame(frame)).(*PhysicsResetFrame)
	}

	nextFrame := mf.NextFrame(frame)
	if nextFrame > frame {
		// ドンピシャルフレームがなく、次のキーフレがある場合、次のキーフレを取得
		return mf.values.Get(NewFrame(nextFrame)).(*PhysicsResetFrame)
	}

	// 次のキーフレがない場合、前のキーフレを取得
	prevFrame := mf.PrevFrame(frame)
	return mf.values.Get(NewFrame(prevFrame)).(*PhysicsResetFrame)
}
