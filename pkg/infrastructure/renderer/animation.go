package renderer

import (
	"github.com/miu200521358/mlib_go/pkg/domain/delta"
	"github.com/miu200521358/mlib_go/pkg/domain/vmd"
)

type Animation struct {
	Model                    *RenderModel     // 描画モデル
	Motion                   *vmd.VmdMotion   // モーション
	VmdDeltas                *delta.VmdDeltas // モーション変化量
	InvisibleMaterialIndexes []int            // 非表示材質インデックス
	SelectedVertexIndexes    []int            // 選択頂点インデックス
}
