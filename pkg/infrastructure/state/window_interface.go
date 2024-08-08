//go:build windows
// +build windows

package state

import (
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/miu200521358/mlib_go/pkg/domain/delta"
	"github.com/miu200521358/mlib_go/pkg/domain/mmath"
	"github.com/miu200521358/mlib_go/pkg/domain/pmx"
	"github.com/miu200521358/mlib_go/pkg/infrastructure/mbt"
	"github.com/miu200521358/mlib_go/pkg/infrastructure/mgl"
)

type IControlWindow interface {
	Run()
	Dispose()
	Close()
	Size() (int, int)
	SetPosition(x, y int)
	Frame() float32
	SetFrame(frame float32)
	SetFrameChannel(frame float32)
	UpdateMaxFrame(frame float32)
	UpdateMaxFrameChannel(frame float32)
	UpdateSelectedVertexIndexes(indexes [][][]int)
	SetUpdateSelectedVertexIndexesFunc(f func([][][]int))
	SetPlayingChannel(playing bool)
	SetClosed(closed bool)
}

type IViewWindow interface {
	Dispose()
	Close()
	Size() (int, int)
	SetPosition(x, y int)
	TriggerClose(window *glfw.Window)
	GetWindow() *glfw.Window
	ResetPhysics()
	AppState() IAppState
	Title() string
	OverrideTextureId() uint32
	SetOverrideTextureId(id uint32)
	GetViewerParameter() (float64, float64, *mmath.MVec2, *mmath.MVec3, *mmath.MVec3, *mmath.MVec3)
	UpdateViewerParameter(yaw, pitch float64, size *mmath.MVec2, cameraPos, cameraUp, lookAtCenter *mmath.MVec3)
	Render(models []*pmx.PmxModel, vmdDeltas []*delta.VmdDeltas)
	Physics() mbt.IPhysics
	LoadModels(models []*pmx.PmxModel)
}

type IPlayer interface {
	Playing() bool
	SetPlaying(playing bool)
	Frame() float32
	SetFrame(v float32)
	MaxFrame() float32
	SetMaxFrame(max float32)
	UpdateMaxFrame(max float32)
	SetRange(min, max int)
	SetEnabled(enabled bool)
	SetEnabledPlayButton(enabled bool)
	Enabled() bool
}

type IRenderModel interface {
	Hash() string
	Delete()
	Render(
		shader mgl.IShader, appState IAppState, vmdDeltas *delta.VmdDeltas,
		leftCursorPositions, leftCursorRemovePositions, leftCursorWorldHistoryPositions,
		leftCursorRemoveWorldHistoryPositions []*mgl32.Vec3,
	)
	Model() *pmx.PmxModel
	InvisibleMaterialIndexes() []int
	ExistInvisibleMaterialIndex(index int) bool
	SelectedVertexIndexes() []int
	NoSelectedVertexIndexes() []int
	SetSelectedVertexIndexes(indexes []int)
	SetNoSelectedVertexIndexes(indexes []int)
	ClearSelectedVertexIndexes()
	UpdateSelectedVertexIndexes(indexes []int)
	UpdateNoSelectedVertexIndexes(indexes []int)
}
