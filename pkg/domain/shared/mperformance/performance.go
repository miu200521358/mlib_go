// Package mperformance は性能/メモリの共通契約を定義する。
// 可能な限りリアルタイム反映を目指し、困難な場合は適用タイミングを選択可能とする。
package mperformance

// MaxBoneFrames はボーンフレーム数の上限。
const MaxBoneFrames = 600000

// CopyUsesDeepCopy は Copy/DeepCopy を多用する前提。
const CopyUsesDeepCopy = true

// TextureCacheEnabled はテクスチャキャッシュの利用有無。
const TextureCacheEnabled = true

// UseSyncPool は sync.Pool 等の再利用機構を使うかの規約（使わない）。
const UseSyncPool = false

// EnableFileCacheByPathMTime はパス/mtimeによるキャッシュを使うかの規約（使わない）。
const EnableFileCacheByPathMTime = false

// AllowUserSelectApplyTiming は適用/反映タイミングをユーザー指定可能にする方針。
const AllowUserSelectApplyTiming = true
