// 指示: miu200521358
package collection

// ReindexResult は再インデックス処理の結果を表す。
type ReindexResult struct {
	Changed  bool
	OldToNew []int
	NewToOld []int
	Removed  []int
	Added    []int
}
