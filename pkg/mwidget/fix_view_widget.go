package mwidget

import (
	"github.com/miu200521358/mlib_go/pkg/mutils/mi18n"
	"github.com/miu200521358/walk/pkg/walk"
	"github.com/miu200521358/win"
)

const FixViewWidgetClass = "FixViewWidget Class"

type FixViewWidget struct {
	walk.WidgetBase
	mWindow         *MWindow         // メインウィンドウ
	ModelChoice     *walk.ComboBox   // モデル選択
	BoneChoice      *walk.ComboBox   // ボーン選択
	DirectionChoice *walk.ComboBox   // 方向選択
	ShowButton      *walk.PushButton // 表示ボタン
	IsShow          bool             // 表示中かどうか
}

func NewFixViewWidget(parent walk.Container, mWindow *MWindow) (*FixViewWidget, error) {
	fw := new(FixViewWidget)
	fw.mWindow = mWindow

	if err := walk.InitWidget(
		fw,
		parent,
		FixViewWidgetClass,
		win.WS_DISABLED,
		0); err != nil {

		return nil, err
	}

	composite, err := walk.NewComposite(parent)
	if err != nil {
		return nil, err
	}
	layout := walk.NewHBoxLayout()
	composite.SetLayout(layout)

	// モデル選択
	fw.ModelChoice, err = walk.NewDropDownBox(composite)
	if err != nil {
		return nil, err
	}
	fw.ModelChoice.SetModel([]string{"No.1", "No.2", "No.3"})
	fw.ModelChoice.SetCurrentIndex(0)

	// ボーン選択
	fw.BoneChoice, err = walk.NewDropDownBox(composite)
	if err != nil {
		return nil, err
	}
	fw.BoneChoice.SetModel([]string{"頭", "上半身", "右足首"})
	fw.BoneChoice.SetCurrentIndex(0)

	// 方向選択
	fw.DirectionChoice, err = walk.NewDropDownBox(composite)
	if err != nil {
		return nil, err
	}
	fw.DirectionChoice.SetModel([]string{"正面", "左面", "右面", "背面"})
	fw.DirectionChoice.SetCurrentIndex(0)

	fw.ShowButton, err = walk.NewPushButton(composite)
	if err != nil {
		return nil, err
	}
	fw.ShowButton.SetText(mi18n.T("固定ビュー表示"))
	fw.ShowButton.Clicked().Attach(func() {
		if fw.IsShow {
			fw.ShowButton.SetChecked(false)
			fw.ShowButton.SetText(mi18n.T("固定ビュー表示"))
		} else {
			fw.ShowButton.SetChecked(true)
			fw.ShowButton.SetText(mi18n.T("固定ビュー非表示"))
		}
		fw.IsShow = !fw.IsShow
	})

	return fw, nil
}

func (fw *FixViewWidget) Dispose() {
	fw.WidgetBase.Dispose()
	fw.ModelChoice.Dispose()
	fw.BoneChoice.Dispose()
	fw.DirectionChoice.Dispose()
	fw.ShowButton.Dispose()
}

func (fw *FixViewWidget) SetEnabled(enabled bool) {
	fw.WidgetBase.SetEnabled(enabled)
	fw.ModelChoice.SetEnabled(enabled)
	fw.BoneChoice.SetEnabled(enabled)
	fw.DirectionChoice.SetEnabled(enabled)
	fw.ShowButton.SetEnabled(enabled)
}

func (f *FixViewWidget) CreateLayoutItem(ctx *walk.LayoutContext) walk.LayoutItem {
	return &subViewWidgetLayoutItem{idealSize: walk.SizeFrom96DPI(walk.Size{Width: 50, Height: 50}, ctx.DPI())}
}

type subViewWidgetLayoutItem struct {
	walk.LayoutItemBase
	idealSize walk.Size // in native pixels
}

func (li *subViewWidgetLayoutItem) LayoutFlags() walk.LayoutFlags {
	return 0
}

func (li *subViewWidgetLayoutItem) IdealSize() walk.Size {
	return li.idealSize
}
