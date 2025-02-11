package state

type IViewerState interface {
	Frame() float32
	SetFrame(frame float32)
	MaxFrame() float32
	SetMaxFrame(maxFrame float32)
	UpdateMaxFrame(maxFrame float32)
	Playing() bool
	SetPlaying(p bool)
	FrameInterval() float64
	SetFrameInterval(spf float64)
	IsClosed() bool
	SetClosed(closed bool)
}

type viewerState struct {
	frame         float32 // フレーム
	maxFrame      float32 // 最大フレーム
	playing       bool    // 再生中フラグ
	frameInterval float64 // FPS制限
	isClosed      bool    // 描画ウィンドウクローズ
}

func NewViewerState() IViewerState {
	return &viewerState{
		frame:         0.0,
		maxFrame:      1,
		frameInterval: 1.0 / 30.0, // 30fps
	}
}

func (viewerState *viewerState) Frame() float32 {
	return viewerState.frame
}

func (viewerState *viewerState) SetFrame(frame float32) {
	viewerState.frame = frame
}

func (viewerState *viewerState) MaxFrame() float32 {
	return viewerState.maxFrame
}

func (viewerState *viewerState) UpdateMaxFrame(maxFrame float32) {
	if viewerState.maxFrame < maxFrame {
		viewerState.maxFrame = maxFrame
	}
}

func (viewerState *viewerState) SetMaxFrame(maxFrame float32) {
	viewerState.maxFrame = maxFrame
}

func (viewerState *viewerState) Playing() bool {
	return viewerState.playing
}

func (viewerState *viewerState) SetPlaying(p bool) {
	viewerState.playing = p
}

func (viewerState *viewerState) FrameInterval() float64 {
	return viewerState.frameInterval
}

func (viewerState *viewerState) SetFrameInterval(spf float64) {
	viewerState.frameInterval = spf
}

func (viewerState *viewerState) IsClosed() bool {
	return viewerState.isClosed
}

func (viewerState *viewerState) SetClosed(closed bool) {
	viewerState.isClosed = closed
}
