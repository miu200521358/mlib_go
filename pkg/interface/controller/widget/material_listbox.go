package widget

import (
	"github.com/miu200521358/mlib_go/pkg/domain/pmx"
	"github.com/miu200521358/mlib_go/pkg/interface/controller"
	"github.com/miu200521358/walk/pkg/declarative"
	"github.com/miu200521358/walk/pkg/walk"
)

type MaterialListBox struct {
	*walk.ListBox
	window            *controller.ControlWindow // メインウィンドウ
	materials         *pmx.Materials
	MaterialListModel *MaterialListModel
	tooltip           string
	changeFunc        func(cw *controller.ControlWindow, indexes []int)
}

func NewMaterialListbox(tooltip string, changeFunc func(cw *controller.ControlWindow, indexes []int)) *MaterialListBox {
	// 複数選択リストボックス
	m := &MaterialListModel{
		items:                  make([]string, 0),
		itemsResetPublisher:    new(walk.EventPublisher),
		itemChangedPublisher:   new(walk.IntEventPublisher),
		itemsInsertedPublisher: new(walk.IntRangeEventPublisher),
		itemsRemovedPublisher:  new(walk.IntRangeEventPublisher),
	}

	return &MaterialListBox{MaterialListModel: m, changeFunc: changeFunc, tooltip: tooltip}
}

func (lb *MaterialListBox) Widgets() declarative.Composite {
	return declarative.Composite{
		Layout: declarative.VBox{},
		Children: []declarative.Widget{
			declarative.ListBox{
				AssignTo:       &lb.ListBox,
				Model:          lb.MaterialListModel,
				MultiSelection: true,
				MinSize:        declarative.Size{Width: 0, Height: 100},
				MaxSize:        declarative.Size{Width: 1024, Height: 200},
				ToolTipText:    lb.tooltip,
			},
		},
	}
}

func (lb *MaterialListBox) SetMaterials(
	materials *pmx.Materials,
) {
	lb.MaterialListModel.items = make([]string, 0)
	materials.ForEach(func(i int, material *pmx.Material) {
		lb.MaterialListModel.items = append(lb.MaterialListModel.items, material.Name())
	})
	lb.MaterialListModel.PublishItemsReset()
	lb.SetSelectedIndexes(materials.Indexes())
}

func (lb *MaterialListBox) EnabledInPlaying(enable bool) {
	lb.ListBox.SetEnabled(enable)
}

func (lb *MaterialListBox) SetWindow(window *controller.ControlWindow) {
	lb.window = window
	lb.SelectedIndexesChanged().Attach(func() {
		if lb.changeFunc != nil {
			lb.changeFunc(lb.window, lb.SelectedIndexes())
		}
	})
}

// ---------------------------------------

type MaterialListModel struct {
	*walk.ReflectListModelBase
	itemsResetPublisher    *walk.EventPublisher
	itemChangedPublisher   *walk.IntEventPublisher
	itemsInsertedPublisher *walk.IntRangeEventPublisher
	itemsRemovedPublisher  *walk.IntRangeEventPublisher
	items                  []string
}

func (m *MaterialListModel) ItemCount() int {
	return len(m.items)
}

func (m *MaterialListModel) Value(index int) interface{} {
	return m.items[index]
}

func (m *MaterialListModel) Items() interface{} {
	return m.items
}

func (m *MaterialListModel) ItemsReset() *walk.Event {
	return m.itemsResetPublisher.Event()
}

func (m *MaterialListModel) ItemChanged() *walk.IntEvent {
	return m.itemChangedPublisher.Event()
}

func (m *MaterialListModel) ItemsInserted() *walk.IntRangeEvent {
	return m.itemsInsertedPublisher.Event()
}

func (m *MaterialListModel) ItemsRemoved() *walk.IntRangeEvent {
	return m.itemsRemovedPublisher.Event()
}

func (m *MaterialListModel) PublishItemsReset() {
	m.itemsResetPublisher.Publish()
}

func (m *MaterialListModel) PublishItemChanged(index int) {
	m.itemChangedPublisher.Publish(index)
}

func (m *MaterialListModel) PublishItemsInserted(from, to int) {
	m.itemsInsertedPublisher.Publish(from, to)
}

func (m *MaterialListModel) PublishItemsRemoved(from, to int) {
	m.itemsRemovedPublisher.Publish(from, to)
}
