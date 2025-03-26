package widget

import (
	"fmt"
	"slices"
	"sort"

	"github.com/miu200521358/mlib_go/pkg/config/mi18n"
	"github.com/miu200521358/mlib_go/pkg/domain/pmx"
	"github.com/miu200521358/mlib_go/pkg/interface/controller"
	"github.com/miu200521358/walk/pkg/declarative"
	"github.com/miu200521358/walk/pkg/walk"
)

type MaterialTableView struct {
	*walk.TableView
	window              *controller.ControlWindow // メインウィンドウ
	Materials           *pmx.Materials
	MaterialModel       *MaterialModel
	tooltip             string
	prevSelectedIndexes []int
	changeFunc          func(cw *controller.ControlWindow, indexes []int)
}

func NewMaterialTableView(tooltip string, changeFunc func(cw *controller.ControlWindow, indexes []int)) *MaterialTableView {
	m := new(MaterialTableView)
	m.tooltip = tooltip
	m.changeFunc = changeFunc
	m.MaterialModel = new(MaterialModel)

	return m
}

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
					{Title: mi18n.T("日本語名称"), Width: 150},
					{Title: mi18n.T("英語名称"), Width: 150},
					{Title: mi18n.T("有効テクスチャ"), Width: 50},
					{Title: mi18n.T("テクスチャ"), Width: 150},
					{Title: mi18n.T("有効Toon"), Width: 50},
					{Title: "Toon", Width: 150},
					{Title: mi18n.T("有効スフィア"), Width: 50},
					{Title: mi18n.T("スフィア"), Width: 150},
				},
				StyleCell: func(style *walk.CellStyle) {
					m := lb.MaterialModel.Records[style.Row()]
					if (!m.TextureValid && m.TextureNameText != "") ||
						(!m.ToonValid && m.ToonNameText != "") ||
						(!m.SphereValid && m.SphereNameText != "") {
						// テクスチャ、Toon、スフィアが無効の場合は赤文字
						style.TextColor = walk.RGB(255, 0, 0)
					}
					if lb.MaterialModel.Checked(style.Row()) {
						style.BackgroundColor = walk.RGB(159, 255, 243)
					} else {
						style.BackgroundColor = walk.RGB(255, 255, 255)
					}
				},
				OnSelectedIndexesChanged: func() {
					for _, i := range lb.SelectedIndexes() {
						if slices.Equal(lb.prevSelectedIndexes, lb.SelectedIndexes()) || !slices.Contains(lb.prevSelectedIndexes, i) {
							lb.MaterialModel.SetChecked(i, !lb.MaterialModel.Checked(i))
						}
					}
					lb.prevSelectedIndexes = lb.SelectedIndexes()

					if lb.changeFunc != nil {
						lb.changeFunc(lb.window, lb.MaterialModel.CheckedIndexes())
					}
				},
			},
		},
	}
}

func (lb *MaterialTableView) EnabledInPlaying(playing bool) {
	lb.TableView.SetEnabled(!playing)
}

func (lb *MaterialTableView) SetWindow(window *controller.ControlWindow) {
	lb.window = window
}

// ---------------------------------------

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

type MaterialModel struct {
	walk.TableModelBase
	walk.SorterBase
	sortColumn int
	sortOrder  walk.SortOrder
	Records    []*MaterialItem
}

func (m *MaterialModel) RowCount() int {
	return len(m.Records)
}

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

func (m *MaterialModel) Checked(row int) bool {
	return m.Records[row].Checked
}

func (m *MaterialModel) SetChecked(row int, checked bool) error {
	m.Records[row].Checked = checked

	return nil
}

func (m *MaterialModel) Sort(col int, order walk.SortOrder) error {
	m.sortColumn, m.sortOrder = col, order

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
			return c(a.SphereNameText < b.ToonNameText)
		}

		panic("unreachable")
	})

	return m.SorterBase.Sort(col, order)
}

func (m *MaterialModel) AddRecord(
	material *pmx.Material, texture *pmx.Texture, toon *pmx.Texture, sphere *pmx.Texture,
) {
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
		EnglishNameText:  material.EnglishName(),
		TextureValid:     textureValid,
		TextureNameText:  textureName,
		ToonValid:        toonValid,
		ToonNameText:     toonName,
		SphereValid:      sphereValid,
		SphereNameText:   sphereName,
	}
	m.Records = append(m.Records, item)
}

func (m *MaterialModel) ResetRows(model *pmx.PmxModel) {
	m.Records = make([]*MaterialItem, 0)

	m.PublishRowsReset()

	if model == nil {
		return
	}

	model.Materials.ForEach(func(i int, mat *pmx.Material) {
		texture, _ := model.Textures.Get(mat.TextureIndex)
		toon, _ := model.Textures.Get(mat.ToonTextureIndex)
		if mat.ToonSharingFlag == pmx.TOON_SHARING_SHARING {
			toon = pmx.NewTexture()
			toon.SetName(fmt.Sprintf("toon/toon%02d.bmp", mat.ToonTextureIndex))
			toon.SetValid(true)
		}
		sphere, _ := model.Textures.Get(mat.SphereTextureIndex)

		m.AddRecord(mat, texture, toon, sphere)
	})

	m.PublishRowsReset()
}

func (m *MaterialModel) CheckedIndexes() []int {
	indexes := make([]int, 0, len(m.Records))
	for i, r := range m.Records {
		if r.Checked {
			indexes = append(indexes, i)
		}
	}
	return indexes
}
