package mmodel

import (
	"strings"

	"github.com/miu200521358/mlib_go/pkg/domain/mcore"
	"github.com/miu200521358/mlib_go/pkg/domain/mmath"
	"github.com/tiendc/go-deepcopy"
)

// Bone はボーンを表します。
type Bone struct {
	mcore.IndexNameModel

	// 基本データ（PMXから読み込み）
	Position     *mmath.Vec3 // ボーン位置
	ParentIndex  int         // 親ボーンのボーンIndex
	Layer        int         // 変形階層
	Flag         BoneFlag    // ボーンフラグ(16bit)
	TailPosition *mmath.Vec3 // 接続先座標（Flag.IsTailBone()=falseの場合）
	TailIndex    int         // 接続先ボーンIndex（Flag.IsTailBone()=trueの場合）
	EffectIndex  int         // 付与親ボーンのボーンIndex
	EffectFactor float64     // 付与率
	FixedAxis    *mmath.Vec3 // 軸固定ベクトル（Flag.HasFixedAxis()=trueの場合）
	LocalAxisX   *mmath.Vec3 // ローカルX軸（Flag.HasLocalAxis()=trueの場合）
	LocalAxisZ   *mmath.Vec3 // ローカルZ軸（Flag.HasLocalAxis()=trueの場合）
	EffectorKey  int         // 外部親変形Key（Flag.IsExternalParentDeform()=trueの場合）
	Ik           *Ik         // IKデータ（Flag.IsIK()=trueの場合）

	// 表示
	DisplaySlotIndex int  // 該当表示枠Index
	IsSystem         bool // システム計算用追加ボーン

	// 計算済みデータ（Setup後に設定）
	NormalizedLocalAxisX   *mmath.Vec3 // 計算済みX軸方向ベクトル
	NormalizedLocalAxisY   *mmath.Vec3 // 計算済みY軸方向ベクトル
	NormalizedLocalAxisZ   *mmath.Vec3 // 計算済みZ軸方向ベクトル
	NormalizedFixedAxis    *mmath.Vec3 // 計算済み軸制限ベクトル
	LocalAxis              *mmath.Vec3 // ローカル軸方向ベクトル
	ParentRelativePosition *mmath.Vec3 // 親ボーンからの相対位置
	ChildRelativePosition  *mmath.Vec3 // Tailボーンへの相対位置
	RevertOffsetMatrix     *mmath.Mat4 // 逆オフセット行列
	OffsetMatrix           *mmath.Mat4 // オフセット行列（自身の位置を原点に戻す）

	// ボーン関係（Setup後に設定）
	TreeBoneIndexes      []int    // 自分のボーンまでのボーンIndexリスト
	ParentBoneIndexes    []int    // 親ボーンからルートまでのボーンIndexリスト
	ParentBoneNames      []string // 親ボーンからルートまでのボーン名リスト
	RelativeBoneIndexes  []int    // 関連ボーンIndex一覧（付与親, IK等）
	ChildBoneIndexes     []int    // 子ボーンIndex一覧
	EffectiveBoneIndexes []int    // 付与子ボーンIndex一覧
	IkLinkBoneIndexes    []int    // IKリンク登録ボーンIndex一覧
	IkTargetBoneIndexes  []int    // IKターゲット登録ボーンIndex一覧

	// IKリンク時の角度制限（Setup後に設定）
	AngleLimit         bool        // 角度制限有無
	MinAngleLimit      *mmath.Vec3 // 角度制限下限
	MaxAngleLimit      *mmath.Vec3 // 角度制限上限
	LocalAngleLimit    bool        // ローカル軸角度制限有無
	LocalMinAngleLimit *mmath.Vec3 // ローカル軸角度制限下限
	LocalMaxAngleLimit *mmath.Vec3 // ローカル軸角度制限上限

	// その他
	AxisSign      int   // ボーン軸符号（左:-1, 右:1）
	OriginalLayer int   // 元の変形階層
	ParentBone    *Bone // 親ボーン参照（Setup後に設定）
}

// NewBone は新しいBoneを生成します。
func NewBone() *Bone {
	return &Bone{
		IndexNameModel:   *mcore.NewIndexNameModel(-1, "", ""),
		Position:         mmath.NewVec3(),
		ParentIndex:      -1,
		Layer:            -1,
		Flag:             BONE_FLAG_NONE,
		TailPosition:     mmath.NewVec3(),
		TailIndex:        -1,
		EffectIndex:      -1,
		EffectFactor:     0.0,
		FixedAxis:        nil,
		LocalAxisX:       nil,
		LocalAxisZ:       nil,
		EffectorKey:      -1,
		Ik:               nil,
		DisplaySlotIndex: -1,
		IsSystem:         false,
		// 計算済みフィールドはnilで初期化
		NormalizedLocalAxisX:   nil,
		NormalizedLocalAxisY:   nil,
		NormalizedLocalAxisZ:   nil,
		NormalizedFixedAxis:    nil,
		LocalAxis:              nil,
		ParentRelativePosition: nil,
		ChildRelativePosition:  nil,
		RevertOffsetMatrix:     nil,
		OffsetMatrix:           nil,
		TreeBoneIndexes:        nil,
		ParentBoneIndexes:      nil,
		ParentBoneNames:        nil,
		RelativeBoneIndexes:    nil,
		ChildBoneIndexes:       nil,
		EffectiveBoneIndexes:   nil,
		IkLinkBoneIndexes:      nil,
		IkTargetBoneIndexes:    nil,
		AngleLimit:             false,
		MinAngleLimit:          nil,
		MaxAngleLimit:          nil,
		LocalAngleLimit:        false,
		LocalMinAngleLimit:     nil,
		LocalMaxAngleLimit:     nil,
		AxisSign:               1,
		OriginalLayer:          -1,
		ParentBone:             nil,
	}
}

// NewBoneByName は名前を指定して新しいBoneを生成します。
func NewBoneByName(name string) *Bone {
	b := NewBone()
	b.SetName(name)
	return b
}

// IsValid はBoneが有効かどうかを返します。
func (b *Bone) IsValid() bool {
	return b != nil && b.Index() >= 0
}

// Direction はボーンの方向（左/右/体幹）を返します。
func (b *Bone) Direction() BoneDirection {
	if strings.Contains(b.Name(), "左") {
		return BONE_DIRECTION_LEFT
	} else if strings.Contains(b.Name(), "右") {
		return BONE_DIRECTION_RIGHT
	}
	return BONE_DIRECTION_TRUNK
}

// Copy は深いコピーを作成します。
func (b *Bone) Copy() (*Bone, error) {
	cp := &Bone{}
	if err := deepcopy.Copy(cp, b); err != nil {
		return nil, err
	}
	// ParentBoneは参照のためnilにする
	cp.ParentBone = nil
	return cp, nil
}

// NormalizeFixedAxis は軸固定ベクトルを正規化して設定します。
func (b *Bone) NormalizeFixedAxis(fixedAxis *mmath.Vec3) {
	b.NormalizedFixedAxis = fixedAxis.Normalized()
}

// NormalizeLocalAxis はローカル軸を正規化して設定します。
func (b *Bone) NormalizeLocalAxis(localAxisX *mmath.Vec3) {
	b.NormalizedLocalAxisX = localAxisX.Normalized()
	b.NormalizedLocalAxisY = b.NormalizedLocalAxisX.Cross(mmath.VEC3_UNIT_Z_NEG)
	b.NormalizedLocalAxisZ = b.NormalizedLocalAxisX.Cross(b.NormalizedLocalAxisY)
}
