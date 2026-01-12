// 指示: miu200521358
package motion

import (
	"fmt"
	"strings"

	"github.com/miu200521358/mlib_go/pkg/domain/mmath"
	"github.com/miu200521358/mlib_go/pkg/shared/hashable"
)

// VmdMotion はモーション全体を表す。
type VmdMotion struct {
	*hashable.HashableBase
	Signature                  string
	BoneFrames                 *BoneFrames
	MorphFrames                *MorphFrames
	CameraFrames               *CameraFrames
	LightFrames                *LightFrames
	ShadowFrames               *ShadowFrames
	IkFrames                   *IkFrames
	MaxSubStepsFrames          *MaxSubStepsFrames
	FixedTimeStepFrames        *FixedTimeStepFrames
	GravityFrames              *GravityFrames
	PhysicsResetFrames         *PhysicsResetFrames
	RigidBodyFrames            *RigidBodyFrames
	JointFrames                *JointFrames
	WindEnabledFrames          *WindEnabledFrames
	WindDirectionFrames        *WindDirectionFrames
	WindLiftCoeffFrames        *WindLiftCoeffFrames
	WindDragCoeffFrames        *WindDragCoeffFrames
	WindRandomnessFrames       *WindRandomnessFrames
	WindSpeedFrames            *WindSpeedFrames
	WindTurbulenceFreqHzFrames *WindTurbulenceFreqHzFrames
}

// NewVmdMotion はVmdMotionを生成する。
func NewVmdMotion(path string) *VmdMotion {
	motion := &VmdMotion{
		HashableBase:               hashable.NewHashableBase("", path),
		Signature:                  "",
		BoneFrames:                 NewBoneFrames(),
		MorphFrames:                NewMorphFrames(),
		CameraFrames:               NewCameraFrames(),
		LightFrames:                NewLightFrames(),
		ShadowFrames:               NewShadowFrames(),
		IkFrames:                   NewIkFrames(),
		MaxSubStepsFrames:          NewMaxSubStepsFrames(),
		FixedTimeStepFrames:        NewFixedTimeStepFrames(),
		GravityFrames:              NewGravityFrames(),
		PhysicsResetFrames:         NewPhysicsResetFrames(),
		RigidBodyFrames:            NewRigidBodyFrames(),
		JointFrames:                NewJointFrames(),
		WindEnabledFrames:          NewWindEnabledFrames(),
		WindDirectionFrames:        NewWindDirectionFrames(),
		WindLiftCoeffFrames:        NewWindLiftCoeffFrames(),
		WindDragCoeffFrames:        NewWindDragCoeffFrames(),
		WindRandomnessFrames:       NewWindRandomnessFrames(),
		WindSpeedFrames:            NewWindSpeedFrames(),
		WindTurbulenceFreqHzFrames: NewWindTurbulenceFreqHzFrames(),
	}
	motion.SetHashPartsFunc(motion.GetHashParts)
	return motion
}

// IsVpd はVPDか判定する。
func (m *VmdMotion) IsVpd() bool {
	if m == nil {
		return false
	}
	return strings.Contains(strings.ToLower(m.Path()), ".vpd")
}

// UpdateHash はName/Path/FileModTime/GetHashPartsでハッシュを更新する。
func (m *VmdMotion) UpdateHash() {
	if m == nil {
		return
	}
	m.HashableBase.UpdateHash()
}

// UpdateRandomHash はランダムハッシュを設定する。
func (m *VmdMotion) UpdateRandomHash() {
	if m == nil {
		return
	}
	m.HashableBase.UpdateRandomHash()
}

// GetHashParts はハッシュ用の追加要素を返す。
func (m *VmdMotion) GetHashParts() string {
	if m == nil {
		return ""
	}
	boneLen := 0
	if m.BoneFrames != nil {
		boneLen = m.BoneFrames.Len()
	}
	morphLen := 0
	if m.MorphFrames != nil {
		morphLen = m.MorphFrames.Len()
	}
	cameraLen := 0
	if m.CameraFrames != nil {
		cameraLen = m.CameraFrames.Len()
	}
	lightLen := 0
	if m.LightFrames != nil {
		lightLen = m.LightFrames.Len()
	}
	shadowLen := 0
	if m.ShadowFrames != nil {
		shadowLen = m.ShadowFrames.Len()
	}
	ikLen := 0
	if m.IkFrames != nil {
		ikLen = m.IkFrames.Len()
	}
	return fmt.Sprintf("%08d%08d%08d%08d%08d%08d",
		boneLen,
		morphLen,
		cameraLen,
		lightLen,
		shadowLen,
		ikLen,
	)
}

// MaxFrame は最大フレーム番号を返す。
func (m *VmdMotion) MaxFrame() Frame {
	if m == nil {
		return 0
	}
	boneMax := Frame(0)
	morphMax := Frame(0)
	if m.BoneFrames != nil {
		boneMax = m.BoneFrames.MaxFrame()
	}
	if m.MorphFrames != nil {
		morphMax = m.MorphFrames.MaxFrame()
	}
	if boneMax == 0 {
		return morphMax
	}
	if morphMax == 0 {
		return boneMax
	}
	if boneMax > morphMax {
		return boneMax
	}
	return morphMax
}

// MinFrame は最小フレーム番号を返す。
func (m *VmdMotion) MinFrame() Frame {
	if m == nil {
		return 0
	}
	hasBone := m.BoneFrames != nil && m.BoneFrames.Len() > 0
	hasMorph := m.MorphFrames != nil && m.MorphFrames.Len() > 0
	if !hasBone && !hasMorph {
		return 0
	}
	if !hasBone {
		return m.MorphFrames.MinFrame()
	}
	if !hasMorph {
		return m.BoneFrames.MinFrame()
	}
	boneMin := m.BoneFrames.MinFrame()
	morphMin := m.MorphFrames.MinFrame()
	if boneMin < morphMin {
		return boneMin
	}
	return morphMin
}

// Indexes はボーン/モーフの全フレーム番号を返す。
func (m *VmdMotion) Indexes() []int {
	if m == nil {
		return nil
	}
	indexes := make([]int, 0)
	if m.BoneFrames != nil {
		indexes = append(indexes, m.BoneFrames.Indexes()...)
	}
	if m.MorphFrames != nil {
		indexes = append(indexes, m.MorphFrames.Indexes()...)
	}
	indexes = uniqueSortedInts(indexes)
	return indexes
}

// AppendBoneFrame はボーンフレームを追加する。
func (m *VmdMotion) AppendBoneFrame(name string, frame *BoneFrame) {
	if m == nil || m.BoneFrames == nil {
		return
	}
	m.BoneFrames.Get(name).Append(frame)
}

// InsertBoneFrame はボーンフレームを挿入する。
func (m *VmdMotion) InsertBoneFrame(name string, frame *BoneFrame) {
	if m == nil || m.BoneFrames == nil {
		return
	}
	m.BoneFrames.Get(name).Insert(frame)
}

// AppendMorphFrame はモーフフレームを追加する。
func (m *VmdMotion) AppendMorphFrame(name string, frame *MorphFrame) {
	if m == nil || m.MorphFrames == nil {
		return
	}
	m.MorphFrames.Get(name).Append(frame)
}

// InsertMorphFrame はモーフフレームを挿入する。
func (m *VmdMotion) InsertMorphFrame(name string, frame *MorphFrame) {
	if m == nil || m.MorphFrames == nil {
		return
	}
	m.MorphFrames.Get(name).Insert(frame)
}

// AppendCameraFrame はカメラフレームを追加する。
func (m *VmdMotion) AppendCameraFrame(frame *CameraFrame) {
	if m == nil || m.CameraFrames == nil {
		return
	}
	m.CameraFrames.Append(frame)
}

// InsertCameraFrame はカメラフレームを挿入する。
func (m *VmdMotion) InsertCameraFrame(frame *CameraFrame) {
	if m == nil || m.CameraFrames == nil {
		return
	}
	m.CameraFrames.Insert(frame)
}

// AppendLightFrame はライトフレームを追加する。
func (m *VmdMotion) AppendLightFrame(frame *LightFrame) {
	if m == nil || m.LightFrames == nil {
		return
	}
	m.LightFrames.Append(frame)
}

// InsertLightFrame はライトフレームを挿入する。
func (m *VmdMotion) InsertLightFrame(frame *LightFrame) {
	if m == nil || m.LightFrames == nil {
		return
	}
	m.LightFrames.Insert(frame)
}

// AppendShadowFrame はシャドウフレームを追加する。
func (m *VmdMotion) AppendShadowFrame(frame *ShadowFrame) {
	if m == nil || m.ShadowFrames == nil {
		return
	}
	m.ShadowFrames.Append(frame)
}

// InsertShadowFrame はシャドウフレームを挿入する。
func (m *VmdMotion) InsertShadowFrame(frame *ShadowFrame) {
	if m == nil || m.ShadowFrames == nil {
		return
	}
	m.ShadowFrames.Insert(frame)
}

// AppendIkFrame はIKフレームを追加する。
func (m *VmdMotion) AppendIkFrame(frame *IkFrame) {
	if m == nil || m.IkFrames == nil {
		return
	}
	m.IkFrames.Append(frame)
}

// InsertIkFrame はIKフレームを挿入する。
func (m *VmdMotion) InsertIkFrame(frame *IkFrame) {
	if m == nil || m.IkFrames == nil {
		return
	}
	m.IkFrames.Insert(frame)
}

// AppendMaxSubStepsFrame は最大演算回数フレームを追加する。
func (m *VmdMotion) AppendMaxSubStepsFrame(frame *MaxSubStepsFrame) {
	if m == nil || m.MaxSubStepsFrames == nil {
		return
	}
	m.MaxSubStepsFrames.Append(frame)
}

// AppendFixedTimeStepFrame は演算頻度フレームを追加する。
func (m *VmdMotion) AppendFixedTimeStepFrame(frame *FixedTimeStepFrame) {
	if m == nil || m.FixedTimeStepFrames == nil {
		return
	}
	m.FixedTimeStepFrames.Append(frame)
}

// AppendGravityFrame は重力フレームを追加する。
func (m *VmdMotion) AppendGravityFrame(frame *GravityFrame) {
	if m == nil || m.GravityFrames == nil {
		return
	}
	m.GravityFrames.Append(frame)
}

// AppendPhysicsResetFrame は物理リセットフレームを追加する。
func (m *VmdMotion) AppendPhysicsResetFrame(frame *PhysicsResetFrame) {
	if m == nil || m.PhysicsResetFrames == nil {
		return
	}
	m.PhysicsResetFrames.Append(frame)
}

// AppendRigidBodyFrame は剛体フレームを追加する。
func (m *VmdMotion) AppendRigidBodyFrame(name string, frame *RigidBodyFrame) {
	if m == nil || m.RigidBodyFrames == nil {
		return
	}
	frames := m.RigidBodyFrames.Get(name)
	if frames == nil {
		frames = NewRigidBodyNameFrames(name)
		m.RigidBodyFrames.Update(frames)
	}
	frames.Append(frame)
}

// InsertRigidBodyFrame は剛体フレームを挿入する。
func (m *VmdMotion) InsertRigidBodyFrame(name string, frame *RigidBodyFrame) {
	if m == nil || m.RigidBodyFrames == nil {
		return
	}
	frames := m.RigidBodyFrames.Get(name)
	if frames == nil {
		frames = NewRigidBodyNameFrames(name)
		m.RigidBodyFrames.Update(frames)
	}
	frames.Insert(frame)
}

// AppendJointFrame はジョイントフレームを追加する。
func (m *VmdMotion) AppendJointFrame(name string, frame *JointFrame) {
	if m == nil || m.JointFrames == nil {
		return
	}
	frames := m.JointFrames.Get(name)
	if frames == nil {
		frames = NewJointNameFrames(name)
		m.JointFrames.Update(frames)
	}
	frames.Append(frame)
}

// InsertJointFrame はジョイントフレームを挿入する。
func (m *VmdMotion) InsertJointFrame(name string, frame *JointFrame) {
	if m == nil || m.JointFrames == nil {
		return
	}
	frames := m.JointFrames.Get(name)
	if frames == nil {
		frames = NewJointNameFrames(name)
		m.JointFrames.Update(frames)
	}
	frames.Insert(frame)
}

// AppendWindEnabledFrame は風有効フレームを追加する。
func (m *VmdMotion) AppendWindEnabledFrame(frame *WindEnabledFrame) {
	if m == nil || m.WindEnabledFrames == nil {
		return
	}
	m.WindEnabledFrames.Append(frame)
}

// AppendWindDirectionFrame は風向きフレームを追加する。
func (m *VmdMotion) AppendWindDirectionFrame(frame *WindDirectionFrame) {
	if m == nil || m.WindDirectionFrames == nil {
		return
	}
	m.WindDirectionFrames.Append(frame)
}

// AppendWindLiftCoeffFrame は風揚力係数フレームを追加する。
func (m *VmdMotion) AppendWindLiftCoeffFrame(frame *WindLiftCoeffFrame) {
	if m == nil || m.WindLiftCoeffFrames == nil {
		return
	}
	m.WindLiftCoeffFrames.Append(frame)
}

// AppendWindDragCoeffFrame は風抗力係数フレームを追加する。
func (m *VmdMotion) AppendWindDragCoeffFrame(frame *WindDragCoeffFrame) {
	if m == nil || m.WindDragCoeffFrames == nil {
		return
	}
	m.WindDragCoeffFrames.Append(frame)
}

// AppendWindRandomnessFrame は風乱流係数フレームを追加する。
func (m *VmdMotion) AppendWindRandomnessFrame(frame *WindRandomnessFrame) {
	if m == nil || m.WindRandomnessFrames == nil {
		return
	}
	m.WindRandomnessFrames.Append(frame)
}

// AppendWindSpeedFrame は風速フレームを追加する。
func (m *VmdMotion) AppendWindSpeedFrame(frame *WindSpeedFrame) {
	if m == nil || m.WindSpeedFrames == nil {
		return
	}
	m.WindSpeedFrames.Append(frame)
}

// AppendWindTurbulenceFreqHzFrame は風乱流周波数フレームを追加する。
func (m *VmdMotion) AppendWindTurbulenceFreqHzFrame(frame *WindTurbulenceFreqHzFrame) {
	if m == nil || m.WindTurbulenceFreqHzFrames == nil {
		return
	}
	m.WindTurbulenceFreqHzFrames.Append(frame)
}

// Clean は不要なフレームを削除する。
func (m *VmdMotion) Clean() {
	if m == nil {
		return
	}
	if m.BoneFrames != nil {
		m.BoneFrames.Clean()
	}
	if m.MorphFrames != nil {
		m.MorphFrames.Clean()
	}
	if m.CameraFrames != nil {
		m.CameraFrames.Clean()
	}
	if m.LightFrames != nil {
		m.LightFrames.Clean()
	}
	if m.ShadowFrames != nil {
		m.ShadowFrames.Clean()
	}
	if m.IkFrames != nil {
		m.IkFrames.Clean()
	}
}

// Copy はモーションを複製しランダムハッシュに更新する。
func (m *VmdMotion) Copy() (VmdMotion, error) {
	if m == nil {
		return VmdMotion{}, nil
	}
	copied, err := deepCopy(m)
	if err != nil {
		return VmdMotion{}, err
	}
	copied.SetHashPartsFunc(copied.GetHashParts)
	copied.UpdateRandomHash()
	return *copied, nil
}

func uniqueSortedInts(values []int) []int {
	if len(values) == 0 {
		return nil
	}
	values = mmath.Unique(values)
	mmath.Sort(values)
	return values
}
