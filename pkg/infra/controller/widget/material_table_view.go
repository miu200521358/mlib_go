//go:build windows
// +build windows

// 指示: miu200521358
package widget

import (
	"fmt"
	"slices"
	"sort"

	"github.com/miu200521358/mlib_go/pkg/domain/model"
	"github.com/miu200521358/mlib_go/pkg/infra/controller"
	"github.com/miu200521358/mlib_go/pkg/shared/base/i18n"
	"github.com/miu200521358/walk/pkg/declarative"
	"github.com/miu200521358/walk/pkg/walk"
)

// MaterialTableView は材質一覧ウィジェットを表す。
type MaterialTableView struct {
	*walk.TableView
	window              *controller.ControlWindow
	Materials           *model.PmxModel
	MaterialModel       *MaterialModel
	tooltip             string
	translator          i18n.II18n
	changeFunc          func(cw *controller.ControlWindow, indexes []int)
}

// NewMaterialTableView はMaterialTableViewを生成する。
func NewMaterialTableView(translator i18n.II18n, tooltip string, changeFunc func(cw *controller.ControlWindow, indexes []int)) *MaterialTableView {
	m := new(MaterialTableView)
	m.tooltip = tooltip
	m.changeFunc = changeFunc
	m.MaterialModel = new(MaterialModel)
	m.MaterialModel.sortColumn = -1
	m.MaterialModel.sortOrder = walk.SortAscending
	m.translator = translator
	return m
}

// t は翻訳済み文言を返す。
func (lb *MaterialTableView) t(key string) string {
	if lb == nil || lb.translator == nil || !lb.translator.IsReady() {
		return "●●" + key + "●●"
	}
	return lb.translator.T(key)
}

// Widgets はUI構成を返す。
func (lb *MaterialTableView) Widgets() declarative.Composite {
	return declarative.Composite{
		Layout: declarative.VBox{},
		Children: []declarative.Widget{
			declarative.TableView{
				AssignTo:         &lb.TableView,
				AlternatingRowBG: true,
				ColumnsOrderable: true,
				MultiSelection:   true,
				Model:            lb.MaterialModel,
				MinSize:          declarative.Size{Width: 400, Height: 250},
				Columns: []declarative.TableViewColumn{
					{Title: "#", Width: 50},
					{Title: "No.", Width: 50},
					{Title: lb.t("日本語名称"), Width: 150},
					{Title: lb.t("英語名称"), Width: 150},
					{Title: lb.t("有効テクスチャ"), Width: 50},
					{Title: lb.t("テクスチャ"), Width: 150},
					{Title: lb.t("有効Toon"), Width: 50},
					{Title: "Toon", Width: 150},
					{Title: lb.t("有効スフィア"), Width: 50},
					{Title: lb.t("スフィア"), Width: 150},
				},
				OnMouseDown: func(x, y int, button walk.MouseButton) {
					if button != walk.LeftButton || lb.TableView == nil || lb.MaterialModel == nil {
						return
					}
					if walk.ModifiersDown()&(walk.ModControl|walk.ModShift) != 0 {
						return
					}
					row := lb.TableView.IndexAt(x, y)
					if row < 0 {
						return
					}
					// 既に選択済みの行を再クリックした場合のみチェックを反転する。
					if !slices.Contains(lb.TableView.SelectedIndexes(), row) {
						return
					}
					_ = lb.MaterialModel.SetChecked(row, !lb.MaterialModel.Checked(row))
					lb.MaterialModel.PublishRowChanged(row)
					if lb.changeFunc != nil {
						lb.changeFunc(lb.window, lb.MaterialModel.CheckedIndexes())
					}
				},
				StyleCell: func(style *walk.CellStyle) {
					m := lb.MaterialModel.Records[style.Row()]
					if (!m.TextureValid && m.TextureNameText != "") ||
						(!m.ToonValid && m.ToonNameText != "") ||
						(!m.SphereValid && m.SphereNameText != "") {
						style.TextColor = walk.RGB(255, 0, 0)
					}
					if lb.MaterialModel.Checked(style.Row()) {
						style.BackgroundColor = walk.RGB(159, 255, 243)
					} else {
						style.BackgroundColor = walk.RGB(255, 255, 255)
					}
				},
				OnSelectedIndexesChanged: func() {
				},
			},
		},
	}
}

// SetEnabledInPlaying は再生中の有効状態を設定する。
func (lb *MaterialTableView) SetEnabledInPlaying(playing bool) {
	lb.TableView.SetEnabled(!playing)
}

// SetWindow はウィンドウ参照を設定する。
func (lb *MaterialTableView) SetWindow(window *controller.ControlWindow) {
	lb.window = window
}

// ResetRows は材質行を再構築する。
func (lb *MaterialTableView) ResetRows(modelData *model.PmxModel) {
	if lb == nil || lb.MaterialModel == nil {
		return
	}
	lb.MaterialModel.ResetRows(modelData)
	if lb.changeFunc != nil {
		lb.changeFunc(lb.window, lb.MaterialModel.CheckedIndexes())
	}
}

// SetAllChecked は全材質のチェック状態を設定する。
func (lb *MaterialTableView) SetAllChecked(checked bool) {
	if lb == nil || lb.MaterialModel == nil {
		return
	}
	for _, record := range lb.MaterialModel.Records {
		if record == nil {
			continue
		}
		record.Checked = checked
	}
	lb.MaterialModel.PublishRowsReset()
	if lb.changeFunc != nil {
		lb.changeFunc(lb.window, lb.MaterialModel.CheckedIndexes())
	}
}

// InvertChecked は全材質のチェック状態を反転する。
func (lb *MaterialTableView) InvertChecked() {
	if lb == nil || lb.MaterialModel == nil {
		return
	}
	for _, record := range lb.MaterialModel.Records {
		if record == nil {
			continue
		}
		record.Checked = !record.Checked
	}
	lb.MaterialModel.PublishRowsReset()
	if lb.changeFunc != nil {
		lb.changeFunc(lb.window, lb.MaterialModel.CheckedIndexes())
	}
}

// MaterialItem は材質の表示行を表す。
type MaterialItem struct {
	Checked          bool
	Index            int
	JapaneseNameText string
	EnglishNameText  string
	TextureValid     bool
	TextureNameText  string
	ToonValid        bool
	ToonNameText     string
	SphereValid      bool
	SphereNameText   string
}

// MaterialModel は材質テーブルのモデルを表す。
type MaterialModel struct {
	walk.TableModelBase
	walk.SorterBase
	sortColumn int
	sortOrder  walk.SortOrder
	Records    []*MaterialItem
}

// RowCount は行数を返す。
func (m *MaterialModel) RowCount() int {
	return len(m.Records)
}

// Value はセルの値を返す。
func (m *MaterialModel) Value(row, col int) interface{} {
	item := m.Records[row]

	switch col {
	case 0:
		return item.Checked
	case 1:
		return item.Index
	case 2:
		return item.JapaneseNameText
	case 3:
		return item.EnglishNameText
	case 4:
		return item.TextureValid
	case 5:
		return item.TextureNameText
	case 6:
		return item.ToonValid
	case 7:
		return item.ToonNameText
	case 8:
		return item.SphereValid
	case 9:
		return item.SphereNameText
	}

	panic("unexpected col")
}

// Checked はチェック状態を返す。
func (m *MaterialModel) Checked(row int) bool {
	return m.Records[row].Checked
}

// SetChecked はチェック状態を設定する。
func (m *MaterialModel) SetChecked(row int, checked bool) error {
	m.Records[row].Checked = checked
	return nil
}

// ColumnSortable はソート可否を返す。
func (m *MaterialModel) ColumnSortable(col int) bool {
	return col >= 0
}

// SortedColumn は現在のソート列を返す。
func (m *MaterialModel) SortedColumn() int {
	return m.sortColumn
}

// SortOrder は現在のソート順を返す。
func (m *MaterialModel) SortOrder() walk.SortOrder {
	return m.sortOrder
}

// Sort はソート条件を設定する。
func (m *MaterialModel) Sort(col int, order walk.SortOrder) error {
	m.sortColumn, m.sortOrder = col, order
	if m.sortColumn < 0 {
		return m.SorterBase.Sort(col, order)
	}

	sort.SliceStable(m.Records, func(i, j int) bool {
		a, b := m.Records[i], m.Records[j]

		c := func(ls bool) bool {
			if m.sortOrder == walk.SortAscending {
				return ls
			}
			return !ls
		}

		switch m.sortColumn {
		case 0:
			av := 0
			if a.Checked {
				av = 1
			}
			bv := 0
			if b.Checked {
				bv = 1
			}
			return c(av < bv)
		case 1:
			return c(a.Index < b.Index)
		case 2:
			return c(a.JapaneseNameText < b.JapaneseNameText)
		case 3:
			return c(a.EnglishNameText < b.EnglishNameText)
		case 4:
			av := 0
			if a.TextureValid {
				av = 1
			}
			bv := 0
			if b.TextureValid {
				bv = 1
			}
			return c(av < bv)
		case 5:
			return c(a.TextureNameText < b.TextureNameText)
		case 6:
			av := 0
			if a.ToonValid {
				av = 1
			}
			bv := 0
			if b.ToonValid {
				bv = 1
			}
			return c(av < bv)
		case 7:
			return c(a.ToonNameText < b.ToonNameText)
		case 8:
			av := 0
			if a.SphereValid {
				av = 1
			}
			bv := 0
			if b.SphereValid {
				bv = 1
			}
			return c(av < bv)
		case 9:
			return c(a.SphereNameText < b.SphereNameText)
		}

		return false
	})

	return m.SorterBase.Sort(col, order)
}

// AddRecord は材質行を追加する。
func (m *MaterialModel) AddRecord(material *model.Material, texture *model.Texture, toon *model.Texture, sphere *model.Texture) {
	var textureValid, toonValid, sphereValid bool
	var textureName, toonName, sphereName string
	if texture != nil {
		textureValid = texture.IsValid()
		textureName = texture.Name()
	}
	if toon != nil {
		toonValid = toon.IsValid()
		toonName = toon.Name()
	}
	if sphere != nil {
		sphereValid = sphere.IsValid()
		sphereName = sphere.Name()
	}

	item := &MaterialItem{
		Checked:          true,
		Index:            material.Index(),
		JapaneseNameText: material.Name(),
		EnglishNameText:  material.EnglishName,
		TextureValid:     textureValid,
		TextureNameText:  textureName,
		ToonValid:        toonValid,
		ToonNameText:     toonName,
		SphereValid:      sphereValid,
		SphereNameText:   sphereName,
	}
	m.Records = append(m.Records, item)
}

// ResetRows は材質行を再構築する。
func (m *MaterialModel) ResetRows(modelData *model.PmxModel) {
	m.Records = make([]*MaterialItem, 0)
	m.PublishRowsReset()

	if modelData == nil || modelData.Materials == nil {
		return
	}

	materials := modelData.Materials.Values()
	for _, mat := range materials {
		if mat == nil {
			continue
		}
		texture, _ := modelData.Textures.Get(mat.TextureIndex)
		toon, _ := modelData.Textures.Get(mat.ToonTextureIndex)
		if mat.ToonSharingFlag == model.TOON_SHARING_SHARING {
			toon = model.NewTexture()
			toon.SetName(fmt.Sprintf("toon/toon%02d.bmp", mat.ToonTextureIndex))
			toon.SetValid(true)
		}
		sphere, _ := modelData.Textures.Get(mat.SphereTextureIndex)

		m.AddRecord(mat, texture, toon, sphere)
	}

	m.PublishRowsReset()
}

// CheckedIndexes はチェック済みの材質インデックスを返す。
func (m *MaterialModel) CheckedIndexes() []int {
	indexes := make([]int, 0, len(m.Records))
	for _, r := range m.Records {
		if r != nil && r.Checked {
			indexes = append(indexes, r.Index)
		}
	}
	return indexes
}
