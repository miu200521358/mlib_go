//go:build windows
// +build windows

// 指示: miu200521358
package widget

import (
	"fmt"
	"github.com/miu200521358/mlib_go/pkg/adapter/mpresenter/messages"
	"slices"
	"strings"
	"time"

	"github.com/miu200521358/mlib_go/pkg/domain/mmath"
	"github.com/miu200521358/mlib_go/pkg/domain/model"
	"github.com/miu200521358/mlib_go/pkg/infra/controller"
	"github.com/miu200521358/mlib_go/pkg/shared/base/i18n"
	"github.com/miu200521358/walk/pkg/declarative"
	"github.com/miu200521358/walk/pkg/walk"
)

const vertexSelectionSyncInterval = 200 * time.Millisecond

// VertexTableView は選択頂点一覧ウィジェットを表す。
type VertexTableView struct {
	*walk.TableView
	window              *controller.ControlWindow
	modelData           *model.PmxModel
	VertexModel         *VertexModel
	translator          i18n.II18n
	tooltip             string
	viewerIndex         int
	modelIndex          int
	lastSelectedIndexes []int
	syncStop            chan struct{}
}

// NewVertexTableView はVertexTableViewを生成する。
func NewVertexTableView(translator i18n.II18n, tooltip string) *VertexTableView {
	v := new(VertexTableView)
	v.tooltip = tooltip
	v.VertexModel = new(VertexModel)
	v.translator = translator
	return v
}

// t は翻訳済み文言を返す。
func (v *VertexTableView) t(key string) string {
	if v == nil || v.translator == nil || !v.translator.IsReady() {
		return "●●" + key + "●●"
	}
	return v.translator.T(key)
}

// Widgets はUI構成を返す。
func (v *VertexTableView) Widgets() declarative.Composite {
	return declarative.Composite{
		Layout: declarative.VBox{},
		Children: []declarative.Widget{
			declarative.TableView{
				AssignTo:         &v.TableView,
				AlternatingRowBG: true,
				ColumnsOrderable: false,
				MultiSelection:   true,
				Model:            v.VertexModel,
				MinSize:          declarative.Size{Width: 500, Height: 250},
				Columns: []declarative.TableViewColumn{
					{Title: v.t(messages.VertexTableViewKey001), Width: 70},
					{Title: v.t(messages.VertexTableViewKey002), Width: 180},
					{Title: v.t(messages.VertexTableViewKey003), Width: 140},
					{Title: v.t(messages.VertexTableViewKey004), Width: 140},
					{Title: v.t(messages.VertexTableViewKey005), Width: 140},
					{Title: v.t(messages.VertexTableViewKey006), Width: 140},
					{Title: v.t(messages.VertexTableViewKey007), Width: 140},
				},
			},
		},
	}
}

// SetWindow はウィンドウ参照を設定する。
func (v *VertexTableView) SetWindow(window *controller.ControlWindow) {
	v.window = window
}

// SetEnabledInPlaying は再生中の有効状態を設定する。
func (v *VertexTableView) SetEnabledInPlaying(playing bool) {
	if v == nil || v.TableView == nil {
		return
	}
	v.TableView.SetEnabled(!playing)
}

// ResetRows はモデル情報を更新し、行を再構築する。
func (v *VertexTableView) ResetRows(modelData *model.PmxModel) {
	if v == nil || v.VertexModel == nil {
		return
	}
	v.modelData = modelData
	v.syncSelection()
}

// StartSelectionSync は選択頂点の同期処理を開始する。
func (v *VertexTableView) StartSelectionSync(viewerIndex, modelIndex int) {
	if v == nil {
		return
	}
	v.viewerIndex = viewerIndex
	v.modelIndex = modelIndex
	v.stopSelectionSync()
	v.syncSelection()
	if v.window == nil || v.window.IsDisposed() {
		return
	}
	v.syncStop = make(chan struct{})
	go v.runSelectionSync(v.syncStop)
}

// stopSelectionSync は選択頂点の同期処理を停止する。
func (v *VertexTableView) stopSelectionSync() {
	if v == nil || v.syncStop == nil {
		return
	}
	close(v.syncStop)
	v.syncStop = nil
}

// runSelectionSync は定期的に選択頂点の同期処理を行う。
func (v *VertexTableView) runSelectionSync(stop <-chan struct{}) {
	ticker := time.NewTicker(vertexSelectionSyncInterval)
	defer ticker.Stop()
	for {
		select {
		case <-stop:
			return
		case <-ticker.C:
			if v.window == nil || v.window.IsDisposed() {
				return
			}
			v.window.Synchronize(func() {
				v.syncSelection()
			})
		}
	}
}

// syncSelection は選択頂点一覧を更新する。
func (v *VertexTableView) syncSelection() {
	if v == nil || v.VertexModel == nil {
		return
	}
	if v.window == nil {
		v.VertexModel.ResetRows(nil, nil)
		v.lastSelectedIndexes = nil
		return
	}
	indexes := v.window.SelectedVertexIndexes(v.viewerIndex, v.modelIndex)
	normalized := normalizeSelectionIndexes(indexes)
	if slices.Equal(normalized, v.lastSelectedIndexes) {
		return
	}
	v.lastSelectedIndexes = normalized
	v.VertexModel.ResetRows(v.modelData, normalized)
}

// normalizeSelectionIndexes は選択頂点インデックスを正規化する。
func normalizeSelectionIndexes(indexes []int) []int {
	if len(indexes) == 0 {
		return nil
	}
	out := slices.Clone(indexes)
	slices.Sort(out)
	out = slices.Compact(out)
	if len(out) == 0 {
		return nil
	}
	return out
}

// VertexItem は頂点の表示行を表す。
type VertexItem struct {
	Index        int
	PositionText string
	MaterialText string
	Weight1Text  string
	Weight2Text  string
	Weight3Text  string
	Weight4Text  string
}

// VertexModel は頂点テーブルのモデルを表す。
type VertexModel struct {
	walk.TableModelBase
	Records []*VertexItem
}

// RowCount は行数を返す。
func (m *VertexModel) RowCount() int {
	return len(m.Records)
}

// Value はセルの値を返す。
func (m *VertexModel) Value(row, col int) interface{} {
	item := m.Records[row]
	switch col {
	case 0:
		return item.Index
	case 1:
		return item.PositionText
	case 2:
		return item.MaterialText
	case 3:
		return item.Weight1Text
	case 4:
		return item.Weight2Text
	case 5:
		return item.Weight3Text
	case 6:
		return item.Weight4Text
	}
	return nil
}

// ResetRows は選択頂点行を再構築する。
func (m *VertexModel) ResetRows(modelData *model.PmxModel, indexes []int) {
	m.Records = make([]*VertexItem, 0)
	m.PublishRowsReset()

	if modelData == nil || modelData.Vertices == nil {
		return
	}
	for _, idx := range indexes {
		vertex, err := modelData.Vertices.Get(idx)
		if err != nil || vertex == nil {
			continue
		}
		m.Records = append(m.Records, buildVertexItem(modelData, vertex))
	}
	m.PublishRowsReset()
}

// buildVertexItem は表示用の頂点情報を生成する。
func buildVertexItem(modelData *model.PmxModel, vertex *model.Vertex) *VertexItem {
	weights := buildVertexWeights(modelData, vertex)
	return &VertexItem{
		Index:        vertex.Index(),
		PositionText: formatVertexPosition(vertex.Position),
		MaterialText: formatMaterialNames(modelData, vertex.MaterialIndexes),
		Weight1Text:  weights[0],
		Weight2Text:  weights[1],
		Weight3Text:  weights[2],
		Weight4Text:  weights[3],
	}
}

// formatVertexPosition は頂点座標の表示文字列を返す。
func formatVertexPosition(pos mmath.Vec3) string {
	return fmt.Sprintf("X=%.2f, Y=%.2f, Z=%.2f", pos.X, pos.Y, pos.Z)
}

// formatMaterialNames は材質名をカンマ区切りで返す。
func formatMaterialNames(modelData *model.PmxModel, materialIndexes []int) string {
	if modelData == nil || modelData.Materials == nil || len(materialIndexes) == 0 {
		return ""
	}
	names := make([]string, 0, len(materialIndexes))
	seen := make(map[int]struct{})
	for _, idx := range materialIndexes {
		if idx < 0 {
			continue
		}
		if _, ok := seen[idx]; ok {
			continue
		}
		seen[idx] = struct{}{}
		material, err := modelData.Materials.Get(idx)
		if err != nil || material == nil {
			continue
		}
		name := material.Name()
		if name == "" {
			name = fmt.Sprintf("#%d", idx)
		}
		names = append(names, name)
	}
	return strings.Join(names, ", ")
}

// buildVertexWeights は頂点ウェイトの表示文字列を構築する。
func buildVertexWeights(modelData *model.PmxModel, vertex *model.Vertex) [4]string {
	var out [4]string
	if vertex == nil || vertex.Deform == nil {
		return out
	}
	indexes := vertex.Deform.Indexes()
	weights := vertex.Deform.Weights()
	for i := 0; i < len(indexes) && i < len(out); i++ {
		boneIndex := indexes[i]
		if boneIndex < 0 {
			continue
		}
		if i >= len(weights) {
			continue
		}
		weight := weights[i]
		if weight <= 0 {
			continue
		}
		boneName := resolveBoneName(modelData, boneIndex)
		if boneName == "" {
			boneName = fmt.Sprintf("#%d", boneIndex)
		}
		out[i] = fmt.Sprintf("%s(%.2f)", boneName, weight)
	}
	return out
}

// resolveBoneName はボーン名を取得する。
func resolveBoneName(modelData *model.PmxModel, boneIndex int) string {
	if modelData == nil || modelData.Bones == nil {
		return ""
	}
	bone, err := modelData.Bones.Get(boneIndex)
	if err != nil || bone == nil {
		return ""
	}
	return bone.Name()
}
