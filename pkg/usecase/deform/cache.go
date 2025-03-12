package deform

import (
	"math"
	"sync"

	"github.com/miu200521358/mlib_go/pkg/config/mlog"
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
		mlog.IS("[%v] SpeculativeDeformCache.GetResult: cache is nil", frame)
		return nil
	}

	// モデルやモーションが変わっていたらキャッシュ無効
	if c.modelHash != modelHash || c.motionHash != motionHash {
		mlog.IS("[%v] SpeculativeDeformCache.GetResult: cache is invalid", frame)
		return nil
	}

	// フレームの整数部分が一致していて、小数部分の差が許容範囲内かチェック
	frameInt := float32(int(frame))
	predictedInt := float32(int(c.predictedFrame))

	// 整数部分が同じで、小数部分の差が0.1以内なら採用
	if frameInt == predictedInt {
		fracDiff := math.Abs(float64(frame-frameInt) - float64(c.predictedFrame-predictedInt))
		if fracDiff <= 0.1 {
			mlog.IS("[%v] SpeculativeDeformCache.GetResult: cache hit, diff=%.4f", frame, fracDiff)
			return c.result
		}
	}

	mlog.IS("[%v] SpeculativeDeformCache.GetResult: cache is invalid", frame)

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
