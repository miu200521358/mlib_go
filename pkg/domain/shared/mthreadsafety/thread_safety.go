// Package mthreadsafety はスレッド安全性の共通契約を定義する。
// 共有状態の最小化を前提にし、明示コピーで並列安全性を担保する。
package mthreadsafety

// SharedStateUsesAtomic は共有状態に atomic を使う前提。
const SharedStateUsesAtomic = true

// MotionCopyUsesMutex はモーションの Copy が mutex で排他される前提。
const MotionCopyUsesMutex = true

// FrameNameCacheUsesAtomicValue はフレーム名キャッシュが atomic.Value を使う前提。
const FrameNameCacheUsesAtomicValue = true

// RepositoryIsThreadSafe はリポジトリのスレッド安全性（安全ではない）。
const RepositoryIsThreadSafe = false

// ModelsAreMutable はモデル/モーションが可変である前提。
const ModelsAreMutable = true

// RandomSeedManaged は乱数 seed 管理の有無（管理しない）。
const RandomSeedManaged = false
