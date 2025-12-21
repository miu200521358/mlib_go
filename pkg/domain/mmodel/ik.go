package mmodel

import (
	"github.com/miu200521358/mlib_go/pkg/domain/mmath"
	"github.com/tiendc/go-deepcopy"
)

// IkLink はIKリンクを表します。
type IkLink struct {
	BoneIndex          int         // リンクボーンのボーンIndex
	AngleLimit         bool        // 角度制限有無
	MinAngleLimit      *mmath.Vec3 // 角度制限下限
	MaxAngleLimit      *mmath.Vec3 // 角度制限上限
	LocalAngleLimit    bool        // ローカル軸角度制限有無
	LocalMinAngleLimit *mmath.Vec3 // ローカル軸制限下限
	LocalMaxAngleLimit *mmath.Vec3 // ローカル軸制限上限
}

// NewIkLink は新しいIkLinkを生成します。
func NewIkLink() *IkLink {
	return &IkLink{
		BoneIndex:          -1,
		AngleLimit:         false,
		MinAngleLimit:      mmath.NewVec3(),
		MaxAngleLimit:      mmath.NewVec3(),
		LocalAngleLimit:    false,
		LocalMinAngleLimit: mmath.NewVec3(),
		LocalMaxAngleLimit: mmath.NewVec3(),
	}
}

// Copy は深いコピーを作成します。
func (l *IkLink) Copy() (*IkLink, error) {
	cp := &IkLink{}
	if err := deepcopy.Copy(cp, l); err != nil {
		return nil, err
	}
	return cp, nil
}

// Ik はIKデータを表します。
type Ik struct {
	BoneIndex    int         // IKターゲットボーンのボーンIndex
	LoopCount    int         // IKループ回数 (最大255)
	UnitRotation *mmath.Vec3 // IKループ計算時の1回あたりの制限角度
	Links        []*IkLink   // IKリンクリスト
}

// NewIk は新しいIkを生成します。
func NewIk() *Ik {
	return &Ik{
		BoneIndex:    -1,
		LoopCount:    0,
		UnitRotation: mmath.NewVec3(),
		Links:        make([]*IkLink, 0),
	}
}

// Copy は深いコピーを作成します。
func (ik *Ik) Copy() (*Ik, error) {
	cp := &Ik{}
	if err := deepcopy.Copy(cp, ik); err != nil {
		return nil, err
	}
	return cp, nil
}
