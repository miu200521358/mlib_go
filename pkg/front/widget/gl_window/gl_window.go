package gl_window

import (
	"fyne.io/fyne/v2"
	"github.com/fyne-io/glfw-js"

	"github.com/miu200521358/mlib_go/pkg/pmx/pmx_model"

)

type ModelData struct {
	Model *pmx_model.PmxModel
}

type GlWindow struct {
	app                    *fyne.App
	window                 *glfw.Window
	icon                   fyne.Resource
	title                  string
	fullScreen             bool
	size                   *fyne.Size
	callbackClose          func()
	callbackCloseIntercept func()
	callbackDrop           func(fyne.Position, []fyne.URI)
	Data                   []ModelData
}

func NewGlWindow(
	app *fyne.App,
	title string,
	icon fyne.Resource,
	callbackClose func(),
	callbackCloseIntercept func(),
	callbackDrop func(fyne.Position, []fyne.URI),
) (*GlWindow, error) {
	w := GlWindow{
		app:                    app,
		icon:                   icon,
		title:                  title,
		fullScreen:             false,
		size:                   &fyne.Size{Width: 1024, Height: 768},
		callbackClose:          callbackClose,
		callbackCloseIntercept: callbackCloseIntercept,
		callbackDrop:           callbackDrop,
	}
	w.Data = make([]ModelData, 0)

	// GLFW の初期化
	if err := glfw.Init(nil); err != nil {
		return nil, err
	}

	// ウィンドウの作成
	window, err := glfw.CreateWindow(800, 600, "OpenGL with GLFW", nil, nil)
	if err != nil {
		return nil, err
	}
	w.window = window

	return &w, nil
}

func (w *GlWindow) AddData(pmxModel *pmx_model.PmxModel) {
	w.Data = append(w.Data, ModelData{Model: pmxModel})
}

func (w *GlWindow) Title() string {
	return w.title
}

func (w *GlWindow) SetTitle(title string) {
	w.title = title
}

func (w *GlWindow) FullScreen() bool {
	return w.fullScreen
}

func (w *GlWindow) SetFullScreen(fullScreen bool) {
	w.fullScreen = fullScreen
}

func (w *GlWindow) Resize(size fyne.Size) {
	w.size = &size
}

func (w *GlWindow) RequestFocus() {
}

func (w *GlWindow) FixedSize() bool {
	return false
}

func (w *GlWindow) SetFixedSize(fixedSize bool) {
}

func (w *GlWindow) CenterOnScreen() {
	// TODO: Implement
}

func (w *GlWindow) Padded() bool {
	return true
}

func (w *GlWindow) SetPadded(padded bool) {
}

func (w *GlWindow) Icon() fyne.Resource {
	return w.icon
}

func (w *GlWindow) SetIcon(icon fyne.Resource) {
	w.icon = icon
}

func (w *GlWindow) SetMaster() {
}

func (w *GlWindow) MainMenu() *fyne.MainMenu {
	return nil
}

func (w *GlWindow) SetMainMenu(menu *fyne.MainMenu) {
}

func (w *GlWindow) SetOnClosed(onClosed func()) {
	w.callbackClose = onClosed
}

func (w *GlWindow) SetCloseIntercept(closeIntercept func()) {
	w.callbackCloseIntercept = closeIntercept
}

func (w *GlWindow) SetOnDropped(onDropped func(fyne.Position, []fyne.URI)) {
	w.callbackDrop = onDropped
}

func (w *GlWindow) Show() {
}

func (w *GlWindow) Hide() {
	w.window.Hide()
}

func (w *GlWindow) Close() {
	w.window.Destroy()
	glfw.Terminate()
}

func (w *GlWindow) ShowAndRun() {
	for !w.window.ShouldClose() {
		glfw.WaitEvents()
	}
}

func (w *GlWindow) Content() fyne.CanvasObject {
	return nil
}

func (w *GlWindow) SetContent(content fyne.CanvasObject) {
}

func (w *GlWindow) Canvas() fyne.Canvas {
	return nil
}

func (w *GlWindow) Clipboard() fyne.Clipboard {
	return nil
}
