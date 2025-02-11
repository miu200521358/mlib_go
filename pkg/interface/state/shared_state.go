package state

import (
	"sync"
)

// SharedState は、コントローラ用Stateとビューワー用Stateをまとめて保持する
// → 各StateのSetterは、それぞれの側(Controller/Viewer)のみが呼ぶ想定。
// → 他方はGetterだけを呼ぶ。
type SharedState struct {
	ctrlState IControllerState
	viewState IViewerState

	// 片方に対する更新/読み取りがもう片方をブロックしないようにMutexを分割する
	ctrlMu sync.RWMutex
	viewMu sync.RWMutex
}

// NewSharedState は2つのStateを注入して生成するコンストラクタ
func NewSharedState(
	c IControllerState,
	v IViewerState,
) *SharedState {
	return &SharedState{
		ctrlState: c,
		viewState: v,
	}
}

// --- コントローラ側Stateの取得/更新 ---

// ReadControllerState はコントローラStateを読み取る (ビューワー/コントローラ両方から呼べる)
func (s *SharedState) ReadControllerState(fn func(cs IControllerState)) {
	s.ctrlMu.RLock()
	defer s.ctrlMu.RUnlock()
	fn(s.ctrlState)
}

// UpdateControllerState は「コントローラが自身の状態を更新する」ためのメソッド
func (s *SharedState) UpdateControllerState(fn func(cs IControllerState)) {
	s.ctrlMu.Lock()
	defer s.ctrlMu.Unlock()
	fn(s.ctrlState)
}

// --- ビューワー側Stateの取得/更新 ---

// ReadViewerState はビューワーStateを読み取る (ビューワー/コントローラ両方から呼べる)
func (s *SharedState) ReadViewerState(fn func(vs IViewerState)) {
	s.viewMu.RLock()
	defer s.viewMu.RUnlock()
	fn(s.viewState)
}

// UpdateViewerState は「ビューワーが自身の状態を更新する」ためのメソッド
func (s *SharedState) UpdateViewerState(fn func(vs IViewerState)) {
	s.viewMu.Lock()
	defer s.viewMu.Unlock()
	fn(s.viewState)
}
