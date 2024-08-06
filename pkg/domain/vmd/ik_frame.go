package vmd

type IkEnabledFrame struct {
	*BaseFrame        // キーフレ
	BoneName   string // IKボーン名
	Enabled    bool   // IKON/OFF
}

func NewIkEnableFrame(index float64) *IkEnabledFrame {
	return &IkEnabledFrame{
		BaseFrame: NewFrame(index).(*BaseFrame),
		BoneName:  "",
		Enabled:   true,
	}
}

func (kf *IkEnabledFrame) Copy() *IkEnabledFrame {
	vv := &IkEnabledFrame{
		BoneName: kf.BoneName,
		Enabled:  kf.Enabled,
	}
	return vv
}

func (nextKf *IkEnabledFrame) LerpFrame(prevFrame IBaseFrame, index float64) IBaseFrame {
	// 補間なしで前のキーフレを引き継ぐ
	return prevFrame
}

type IkFrame struct {
	*BaseFrame                   // キーフレ
	Visible    bool              // 表示ON/OFF
	IkList     []*IkEnabledFrame // IKリスト
}

func NewIkFrame(index float64) *IkFrame {
	return &IkFrame{
		BaseFrame: NewFrame(index).(*BaseFrame),
		Visible:   true,
		IkList:    make([]*IkEnabledFrame, 0),
	}
}

func NullNewIkFrame() *IkFrame {
	return nil
}

func (ikf *IkFrame) Copy() IBaseFrame {
	vv := &IkFrame{
		BaseFrame: ikf.BaseFrame.Copy().(*BaseFrame),
		Visible:   ikf.Visible,
		IkList:    make([]*IkEnabledFrame, len(ikf.IkList)),
	}
	for i, v := range ikf.IkList {
		vv.IkList[i] = v.Copy()
	}
	return vv
}

func (nextIkf *IkFrame) lerpFrame(prevFrame IBaseFrame, index float64) IBaseFrame {
	prevIkf := prevFrame.(*IkFrame)
	// 補間なしで前のキーフレを引き継ぐ
	vv := &IkFrame{
		Visible: prevIkf.Visible,
		IkList:  make([]*IkEnabledFrame, 0, len(prevIkf.IkList)),
	}
	for _, v := range prevIkf.IkList {
		vv.IkList = append(vv.IkList, v.Copy())
	}
	return vv
}

func (ikf *IkFrame) splitCurve(prevFrame IBaseFrame, nextFrame IBaseFrame, index float64) {
}

func (ikf *IkFrame) IsEnable(boneName string) bool {
	for _, ik := range ikf.IkList {
		if ik.BoneName == boneName {
			// IKボーン名がある場合、その結果を取得する
			return ik.Enabled
		}
	}
	// キーが無い場合、IK有効がデフォルト
	return true
}
