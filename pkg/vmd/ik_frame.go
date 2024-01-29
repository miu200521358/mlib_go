package vmd

type IkOnOffFrame struct {
	*BaseFrame        // キーフレ
	Name       string // IK名
	IkFrag     bool   // IKON/OFF
}

func NewIkOnOffFrame(index int) *IkOnOffFrame {
	return &IkOnOffFrame{
		BaseFrame: NewVmdBaseFrame(index),
		Name:      "",
		IkFrag:    true,
	}
}

func (kf *IkOnOffFrame) Copy() *IkOnOffFrame {
	vv := &IkOnOffFrame{
		Name:   kf.Name,
		IkFrag: kf.IkFrag,
	}
	return vv
}

type IkFrame struct {
	*BaseFrame                 // キーフレ
	Show       bool            // 表示ON/OFF
	IkList     []*IkOnOffFrame // IKリスト
}

func NewIkFrame(index int) *IkFrame {
	return &IkFrame{
		BaseFrame: NewVmdBaseFrame(index),
		Show:      true,
		IkList:    []*IkOnOffFrame{},
	}
}

func (ikf *IkFrame) Copy() *IkFrame {
	vv := &IkFrame{
		Show:   ikf.Show,
		IkList: []*IkOnOffFrame{},
	}
	for _, v := range ikf.IkList {
		vv.IkList = append(vv.IkList, v.Copy())
	}
	return vv
}
