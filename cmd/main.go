package main

import (
	"embed"
	"fmt"
	"runtime"

	"github.com/miu200521358/walk/pkg/declarative"

	"github.com/miu200521358/mlib_go/pkg/utils/config"

)

//go:embed resources/app_config.json
var appConfig embed.FS

func init() {
	runtime.LockOSThread()
}

type Foo struct {
	Bar string
	Baz int
}

func main() {
	appConfig := config.ReadAppConfig(appConfig)

	// window, err := declarative.NewMainWindowWithName(fmt.Sprintf("%s %s", appConfig.AppName, appConfig.AppVersion))
	// if err != nil {
	// 	panic(err)
	// }
	// window.Show()
	// window.Run()
	foo := &Foo{"b", 0}

	declarative.MainWindow{
		Title:   fmt.Sprintf("%s %s", appConfig.AppName, appConfig.AppVersion),
		MinSize: declarative.Size{Width: 320, Height: 240},
		Layout:  declarative.VBox{},
		DataBinder: declarative.DataBinder{
			DataSource: foo,
			AutoSubmit: true,
			OnSubmitted: func() {
				fmt.Println(foo)
			},
		},
		Children: []declarative.Widget{
			// RadioButtonGroup is needed for data binding only.
			declarative.RadioButtonGroup{
				DataMember: "Bar",
				Buttons: []declarative.RadioButton{
					{
						Name:  "aRB",
						Text:  "A",
						Value: "a",
					},
					{
						Name:  "bRB",
						Text:  "B",
						Value: "b",
					},
					{
						Name:  "cRB",
						Text:  "C",
						Value: "c",
					},
				},
			},
			declarative.Label{
				Text:    "A",
				Enabled: declarative.Bind("aRB.Checked"),
			},
			declarative.Label{
				Text:    "B",
				Enabled: declarative.Bind("bRB.Checked"),
			},
			declarative.Label{
				Text:    "C",
				Enabled: declarative.Bind("cRB.Checked"),
			},
			declarative.RadioButtonGroup{
				DataMember: "Baz",
				Buttons: []declarative.RadioButton{
					{
						Name:  "oneRB",
						Text:  "1",
						Value: 1,
					},
					{
						Name:  "twoRB",
						Text:  "2",
						Value: 2,
					},
					{
						Name:  "threeRB",
						Text:  "3",
						Value: 3,
					},
				},
			},
			declarative.Label{
				Text:    "1",
				Enabled: declarative.Bind("oneRB.Checked"),
			},
			declarative.Label{
				Text:    "2",
				Enabled: declarative.Bind("twoRB.Checked"),
			},
			declarative.Label{
				Text:    "3",
				Enabled: declarative.Bind("threeRB.Checked"),
			},
		},
	}.Run()
}
