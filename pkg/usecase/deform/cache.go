package deform

import (
	"math"
	"sync"

	"github.com/miu200521358/mlib_go/pkg/domain/delta"
	"github.com/miu200521358/mlib_go/pkg/domain/physics"
	"github.com/miu200521358/mlib_go/pkg/domain/pmx"
	"github.com/miu200521358/mlib_go/pkg/domain/state"
	"github.com/miu200521358/mlib_go/pkg/domain/vmd"
)

// SpeculativeDeformCache は投機的にDeformした結果をキャッシュする
type SpeculativeDeformCache struct {
	mutex          sync.RWMutex
	predictedFrame float32
	modelHash      string
	motionHash     string
	result         *delta.VmdDeltas
	computing      bool
}

// NewSpeculativeDeformCache は新しいキャッシュインスタンスを作成
func NewSpeculativeDeformCache() *SpeculativeDeformCache {
	return &SpeculativeDeformCache{}
}

// StartComputing は計算開始フラグを設定
func (c *SpeculativeDeformCache) StartComputing(frame float32, modelHash, motionHash string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.computing = true
	c.predictedFrame = frame
	c.modelHash = modelHash
	c.motionHash = motionHash
}

// StoreResult は計算結果を保存
func (c *SpeculativeDeformCache) StoreResult(result *delta.VmdDeltas) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.result = result
	c.computing = false
}

// GetResult はキャッシュされた結果を取得（許容範囲内なら一致とみなす）
func (c *SpeculativeDeformCache) GetResult(frame float32, modelHash, motionHash string) *delta.VmdDeltas {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	// キャッシュが存在しない場合はnil
	if c.result == nil {
		return nil
	}

	// モデルやモーションが変わっていたらキャッシュ無効
	if c.modelHash != modelHash || c.motionHash != motionHash {
		return nil
	}

	// フレームが許容範囲内なら結果を返す
	frameDiff := float32(math.Abs(float64(c.predictedFrame - frame)))
	if frameDiff <= 0.3 {
		return c.result
	}

	return nil
}

// IsComputing は計算中かどうかを返す
func (c *SpeculativeDeformCache) IsComputing() bool {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.computing
}

// Reset はキャッシュをリセットする
func (c *SpeculativeDeformCache) Reset() {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.result = nil
	c.computing = false
	c.predictedFrame = 0
}

// --------------------------------------------------

// 投機実行版のDeform関数
func SpeculativeDeform(
	shared *state.SharedState,
	physics physics.IPhysics,
	model *pmx.PmxModel,
	motion *vmd.VmdMotion,
	deltas *delta.VmdDeltas,
	cache *SpeculativeDeformCache,
	timeStep float32,
	nextFrame float32,
) {
	// フレームドロップONの場合は何もしない
	if shared.IsEnabledFrameDrop() {
		return
	}

	// キャッシュがnilの場合は何もしない
	if cache == nil {
		return
	}

	// 現在計算中の場合は何もしない
	if cache.IsComputing() {
		return
	}

	// モデルやモーションがnilの場合は何もしない
	if model == nil || motion == nil {
		return
	}

	modelHash := model.Hash()
	motionHash := motion.Hash()

	// 計算開始フラグを立てる
	cache.StartComputing(nextFrame, modelHash, motionHash)

	// ゴルーチンで投機的に計算
	go func() {
		// 新しいdeltasオブジェクトを作成
		speculativeDeltas := delta.NewVmdDeltas(nextFrame, model.Bones, modelHash, motionHash)

		// 次のフレームの変形を計算
		result := deformBeforePhysics(model, motion, speculativeDeltas, nextFrame)

		// 結果をキャッシュに保存
		cache.StoreResult(result)
	}()
}
