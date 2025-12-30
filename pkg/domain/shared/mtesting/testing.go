// Package mtesting はテスト/再現性の共通契約を定義する。
// テストは数値/文字列/インデックスの一致を中心に検証する。
package mtesting

// TestResourcesDir はテスト資産の基準ディレクトリ名。
const TestResourcesDir = "test_resources"

// AllowAbsolutePathForNonRedistributable は配布できないデータの絶対パス参照を許容する規約。
const AllowAbsolutePathForNonRedistributable = true

// GoldenInputUsesDirectFiles はゴールデン入力を直接参照する規約。
const GoldenInputUsesDirectFiles = true

// ImageComparisonEnabled は画像比較テストの有無（行わない）。
const ImageComparisonEnabled = false

// DefaultMmdReproTolerance はMMD再現テストの許容差。
const DefaultMmdReproTolerance = 0.03

// 許容差の目安:
// - NearEquals は 1e-10 ～ 1e-6 程度を使い分ける。
