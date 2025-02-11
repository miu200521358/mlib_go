package state

import "sync/atomic"

type SharedState struct {
	flags              uint32       // 32ビット分のフラグを格納
	frameValue         atomic.Value // 現在フレーム
	maxFrameValue      atomic.Value // 最大フレーム
	frameIntervalValue atomic.Value // FPS制限
}

type frameState struct {
	frame float32
}

type maxFrameState struct {
	maxFrame float32
}

type frameIntervalState struct {
	frameInterval float32
}

const (
	flagEnabledFrameDrop   = 1 << iota // フレームドロップON/OFF
	flagEnabledPhysics     = 1 << iota // 物理ON/OFF
	flagPhysicsReset       = 1 << iota // 物理リセット
	flagShowNormal         = 1 << iota // ボーンデバッグ表示
	flagShowWire           = 1 << iota // ワイヤーフレームデバッグ表示
	flagShowOverride       = 1 << iota // オーバーライドデバッグ表示
	flagShowSelectedVertex = 1 << iota // 選択頂点デバッグ表示
	flagShowBoneAll        = 1 << iota // 全ボーンデバッグ表示
	flagShowBoneIk         = 1 << iota // IKボーンデバッグ表示
	flagShowBoneEffector   = 1 << iota // 付与親ボーンデバッグ表示
	flagShowBoneFixed      = 1 << iota // 軸制限ボーンデバッグ表示
	flagShowBoneRotate     = 1 << iota // 回転ボーンデバッグ表示
	flagShowBoneTranslate  = 1 << iota // 移動ボーンデバッグ表示
	flagShowBoneVisible    = 1 << iota // 表示ボーンデバッグ表示
	flagShowRigidBodyFront = 1 << iota // 剛体デバッグ表示(前面)
	flagShowRigidBodyBack  = 1 << iota // 剛体デバッグ表示(埋め込み)
	flagShowJoint          = 1 << iota // ジョイントデバッグ表示
	flagShowInfo           = 1 << iota // 情報デバッグ表示
	flagLimitFps30         = 1 << iota // 30FPS制限
	flagLimitFps60         = 1 << iota // 60FPS制限
	flagUnLimitFps         = 1 << iota // FPS無制限
	flagCameraSync         = 1 << iota // カメラ同期
	flagPlaying            = 1 << iota // 再生中フラグ
	flagClosed             = 1 << iota // 描画ウィンドウクローズ
)

// NewSharedState は2つのStateを注入して生成するコンストラクタ
func NewSharedState() *SharedState {
	return &SharedState{
		flags: 0,
	}
}

// アトミックに取得
func (ss *SharedState) Load() uint32 {
	return atomic.LoadUint32(&ss.flags)
}

// アトミックに書き込み (全体をswap)
func (ss *SharedState) Store(newVal uint32) {
	atomic.StoreUint32(&ss.flags, newVal)
}

// 特定ビットをON/OFFにする
func (ss *SharedState) setBit(bitMask uint32, enable bool) {
	for {
		oldVal := ss.Load()
		newVal := oldVal
		if enable {
			newVal |= bitMask
		} else {
			newVal &= ^bitMask
		}
		if atomic.CompareAndSwapUint32(&ss.flags, oldVal, newVal) {
			return
		}
	}
}

// ビットが立っているかどうか
func (ss *SharedState) isBitSet(bitMask uint32) bool {
	return (ss.Load() & bitMask) != 0
}

func (ss *SharedState) IsEnabledFrameDrop() bool {
	return ss.isBitSet(flagEnabledFrameDrop)
}

func (ss *SharedState) SetEnabledFrameDrop(enabled bool) {
	ss.setBit(flagEnabledFrameDrop, enabled)
}

func (ss *SharedState) IsEnabledPhysics() bool {
	return ss.isBitSet(flagEnabledPhysics)
}

func (ss *SharedState) SetEnabledPhysics(enabled bool) {
	ss.setBit(flagEnabledPhysics, enabled)
}

func (ss *SharedState) IsPhysicsReset() bool {
	return ss.isBitSet(flagPhysicsReset)
}

func (ss *SharedState) SetPhysicsReset(reset bool) {
	ss.setBit(flagPhysicsReset, reset)
}

func (ss *SharedState) IsShowNormal() bool {
	return ss.isBitSet(flagShowNormal)
}

func (ss *SharedState) SetShowNormal(show bool) {
	ss.setBit(flagShowNormal, show)
}

func (ss *SharedState) IsShowWire() bool {
	return ss.isBitSet(flagShowWire)
}

func (ss *SharedState) SetShowWire(show bool) {
	ss.setBit(flagShowWire, show)
}

func (ss *SharedState) IsShowOverride() bool {
	return ss.isBitSet(flagShowOverride)
}

func (ss *SharedState) SetShowOverride(show bool) {
	ss.setBit(flagShowOverride, show)
}

func (ss *SharedState) IsShowSelectedVertex() bool {
	return ss.isBitSet(flagShowSelectedVertex)
}

func (ss *SharedState) SetShowSelectedVertex(show bool) {
	ss.setBit(flagShowSelectedVertex, show)
}

func (ss *SharedState) IsShowBoneAll() bool {
	return ss.isBitSet(flagShowBoneAll)
}

func (ss *SharedState) SetShowBoneAll(show bool) {
	ss.setBit(flagShowBoneAll, show)
}

func (ss *SharedState) IsShowBoneIk() bool {
	return ss.isBitSet(flagShowBoneIk)
}

func (ss *SharedState) SetShowBoneIk(show bool) {
	ss.setBit(flagShowBoneIk, show)
}

func (ss *SharedState) IsShowBoneEffector() bool {
	return ss.isBitSet(flagShowBoneEffector)
}

func (ss *SharedState) SetShowBoneEffector(show bool) {
	ss.setBit(flagShowBoneEffector, show)
}

func (ss *SharedState) IsShowBoneFixed() bool {
	return ss.isBitSet(flagShowBoneFixed)
}

func (ss *SharedState) SetShowBoneFixed(show bool) {
	ss.setBit(flagShowBoneFixed, show)
}

func (ss *SharedState) IsShowBoneRotate() bool {
	return ss.isBitSet(flagShowBoneRotate)
}

func (ss *SharedState) SetShowBoneRotate(show bool) {
	ss.setBit(flagShowBoneRotate, show)
}

func (ss *SharedState) IsShowBoneTranslate() bool {
	return ss.isBitSet(flagShowBoneTranslate)
}

func (ss *SharedState) SetShowBoneTranslate(show bool) {
	ss.setBit(flagShowBoneTranslate, show)
}

func (ss *SharedState) IsShowBoneVisible() bool {
	return ss.isBitSet(flagShowBoneVisible)
}

func (ss *SharedState) SetShowBoneVisible(show bool) {
	ss.setBit(flagShowBoneVisible, show)
}

func (ss *SharedState) IsShowRigidBodyFront() bool {
	return ss.isBitSet(flagShowRigidBodyFront)
}

func (ss *SharedState) SetShowRigidBodyFront(show bool) {
	ss.setBit(flagShowRigidBodyFront, show)
}

func (ss *SharedState) IsShowRigidBodyBack() bool {
	return ss.isBitSet(flagShowRigidBodyBack)
}

func (ss *SharedState) SetShowRigidBodyBack(show bool) {
	ss.setBit(flagShowRigidBodyBack, show)
}

func (ss *SharedState) IsShowJoint() bool {
	return ss.isBitSet(flagShowJoint)
}

func (ss *SharedState) SetShowJoint(show bool) {
	ss.setBit(flagShowJoint, show)
}

func (ss *SharedState) IsShowInfo() bool {
	return ss.isBitSet(flagShowInfo)
}

func (ss *SharedState) SetShowInfo(show bool) {
	ss.setBit(flagShowInfo, show)
}

func (ss *SharedState) IsLimitFps30() bool {
	return ss.isBitSet(flagLimitFps30)
}

func (ss *SharedState) SetLimitFps30(limit bool) {
	ss.setBit(flagLimitFps30, limit)
}

func (ss *SharedState) IsLimitFps60() bool {
	return ss.isBitSet(flagLimitFps60)
}

func (ss *SharedState) SetLimitFps60(limit bool) {
	ss.setBit(flagLimitFps60, limit)
}

func (ss *SharedState) IsUnLimitFps() bool {
	return ss.isBitSet(flagUnLimitFps)
}

func (ss *SharedState) SetUnLimitFps(limit bool) {
	ss.setBit(flagUnLimitFps, limit)
}

func (ss *SharedState) IsCameraSync() bool {
	return ss.isBitSet(flagCameraSync)
}

func (ss *SharedState) SetCameraSync(sync bool) {
	ss.setBit(flagCameraSync, sync)
}

func (ss *SharedState) Playing() bool {
	return ss.isBitSet(flagPlaying)
}

func (ss *SharedState) SetPlaying(p bool) {
	ss.setBit(flagPlaying, p)
}

func (ss *SharedState) IsClosed() bool {
	return ss.isBitSet(flagClosed)
}

func (ss *SharedState) SetClosed(closed bool) {
	ss.setBit(flagClosed, closed)
}

func (ss *SharedState) Frame() float32 {
	return ss.frameValue.Load().(frameState).frame
}

func (ss *SharedState) SetFrame(frame float32) {
	ss.frameValue.Store(frameState{frame: frame})
}

func (ss *SharedState) MaxFrame() float32 {
	return ss.maxFrameValue.Load().(maxFrameState).maxFrame
}

func (ss *SharedState) UpdateMaxFrame(maxFrame float32) {
	if ss.MaxFrame() < maxFrame {
		ss.SetMaxFrame(maxFrame)
	}
}

func (ss *SharedState) SetMaxFrame(maxFrame float32) {
	ss.maxFrameValue.Store(maxFrameState{maxFrame: maxFrame})
}

func (ss *SharedState) FrameInterval() float32 {
	return ss.frameIntervalValue.Load().(frameIntervalState).frameInterval
}

func (ss *SharedState) SetFrameInterval(spf float32) {
	ss.frameIntervalValue.Store(frameIntervalState{frameInterval: spf})
}
