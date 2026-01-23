// 指示: miu200521358
package mgl

import (
	"math"
	"time"

	"github.com/go-gl/gl/v4.3-core/gl"

	"github.com/miu200521358/mlib_go/pkg/adapter/graphics_api"
	"github.com/miu200521358/mlib_go/pkg/domain/mmath"
	"github.com/miu200521358/mlib_go/pkg/domain/model"
	"gonum.org/v1/gonum/spatial/r3"
)

const (
	rigidBodyHighlightColorR   = float32(1.0)
	rigidBodyHighlightColorG   = float32(0.8)
	rigidBodyHighlightColorB   = float32(0.2)
	rigidBodyHighlightColorA   = float32(0.75)
	rigidBodyHighlightSegments = 24
)

// DebugRigidBodyHover はデバッグカーソル下の剛体情報を保持する。
type DebugRigidBodyHover struct {
	RigidBody *model.RigidBody
	HitPoint  *mmath.Vec3
}

// RigidBodyHighlighter は剛体デバッグハイライトを管理する。
type RigidBodyHighlighter struct {
	highlightVertices       []float32
	highlightBuffer         *VertexBufferHandle
	debugHover              *DebugRigidBodyHover
	debugHoverMatrix        mmath.Mat4
	hasDebugHoverMatrix     bool
	debugHoverStartTime     time.Time
	prevRigidBodyDebugState bool
}

// NewRigidBodyHighlighter はRigidBodyHighlighterを生成する。
func NewRigidBodyHighlighter() *RigidBodyHighlighter {
	return &RigidBodyHighlighter{
		highlightVertices: make([]float32, 0),
	}
}

// DebugHoverInfo は剛体デバッグホバー情報を返す。
func (mp *RigidBodyHighlighter) DebugHoverInfo() *DebugRigidBodyHover {
	return mp.debugHover
}

// UpdateDebugHoverByRigidBody は剛体ハイライト情報を更新する。
func (mp *RigidBodyHighlighter) UpdateDebugHoverByRigidBody(modelIndex int, rb *model.RigidBody, enable bool) {
	mp.UpdateDebugHoverByRigidBodyWithMatrix(modelIndex, rb, nil, enable)
}

// UpdateDebugHoverByRigidBodyWithMatrix は剛体ハイライト情報を行列付きで更新する。
func (mp *RigidBodyHighlighter) UpdateDebugHoverByRigidBodyWithMatrix(
	modelIndex int,
	rb *model.RigidBody,
	worldMatrix *mmath.Mat4,
	enable bool,
) {
	if !enable || rb == nil {
		mp.clearDebugHover()
		return
	}
	mp.debugHover = &DebugRigidBodyHover{RigidBody: rb}
	if worldMatrix != nil {
		mp.debugHoverMatrix = *worldMatrix
		mp.hasDebugHoverMatrix = true
		mp.rebuildHighlightVertices(rb, &mp.debugHoverMatrix)
	} else {
		mp.hasDebugHoverMatrix = false
		mp.highlightVertices = mp.highlightVertices[:0]
	}
	mp.debugHoverStartTime = time.Now()
}

// DrawDebugHighlight は剛体デバッグハイライトを描画する。
func (mp *RigidBodyHighlighter) DrawDebugHighlight(shader graphics_api.IShader, isDrawRigidBodyFront bool) {
	if mp.debugHover == nil || !mp.hasDebugHoverMatrix || len(mp.highlightVertices) == 0 {
		return
	}
	program := shader.Program(graphics_api.ProgramTypePhysics)
	if program == 0 {
		return
	}
	gl.UseProgram(program)

	if isDrawRigidBodyFront {
		gl.Enable(gl.DEPTH_TEST)
		gl.DepthFunc(gl.ALWAYS)
	}

	if mp.highlightBuffer == nil {
		mp.highlightBuffer = NewBufferFactory().NewDebugBuffer(gl.Ptr(&mp.highlightVertices[0]), len(mp.highlightVertices))
	}

	mp.highlightBuffer.Bind()
	mp.highlightBuffer.UpdateDebugBuffer(mp.highlightVertices)
	gl.DrawArrays(gl.LINES, 0, int32(len(mp.highlightVertices)/7))
	mp.highlightBuffer.Unbind()

	if isDrawRigidBodyFront {
		gl.Enable(gl.DEPTH_TEST)
		gl.DepthFunc(gl.LEQUAL)
	}
	gl.UseProgram(0)
}

// CheckAndClearHighlightOnDebugChange は剛体デバッグ状態変更時にハイライトをクリアする。
func (mp *RigidBodyHighlighter) CheckAndClearHighlightOnDebugChange(currentDebugState bool) {
	if mp.prevRigidBodyDebugState != currentDebugState {
		if !currentDebugState {
			mp.clearDebugHover()
		}
		mp.prevRigidBodyDebugState = currentDebugState
	}
}

// CheckAndClearExpiredHighlight は2秒経過したハイライトを自動的にクリアする。
func (mp *RigidBodyHighlighter) CheckAndClearExpiredHighlight() {
	if mp.debugHover == nil {
		return
	}
	if time.Since(mp.debugHoverStartTime) >= 2*time.Second {
		mp.clearDebugHover()
	}
}

// clearDebugHover は剛体デバッグホバー情報をクリアする。
func (mp *RigidBodyHighlighter) clearDebugHover() {
	mp.debugHover = nil
	mp.hasDebugHoverMatrix = false
	if mp.highlightVertices != nil {
		mp.highlightVertices = mp.highlightVertices[:0]
	}
}

// rebuildHighlightVertices は剛体形状のワイヤーフレームを構築する。
func (mp *RigidBodyHighlighter) rebuildHighlightVertices(rb *model.RigidBody, world *mmath.Mat4) {
	mp.highlightVertices = mp.highlightVertices[:0]
	if rb == nil || world == nil {
		return
	}

	size := rb.Size.Absed()
	switch rb.Shape {
	case model.SHAPE_SPHERE:
		mp.appendSphereLines(world, size.X)
	case model.SHAPE_CAPSULE:
		mp.appendCapsuleLines(world, size.X, size.Y*0.5)
	case model.SHAPE_BOX:
		mp.appendBoxLines(world, size)
	default:
		mp.appendBoxLines(world, size)
	}
}

// appendBoxLines は箱剛体のワイヤーフレームを追加する。
func (mp *RigidBodyHighlighter) appendBoxLines(world *mmath.Mat4, size mmath.Vec3) {
	x := size.X
	y := size.Y
	z := size.Z
	corners := []mmath.Vec3{
		{Vec: r3.Vec{X: -x, Y: -y, Z: -z}},
		{Vec: r3.Vec{X: -x, Y: -y, Z: z}},
		{Vec: r3.Vec{X: -x, Y: y, Z: -z}},
		{Vec: r3.Vec{X: -x, Y: y, Z: z}},
		{Vec: r3.Vec{X: x, Y: -y, Z: -z}},
		{Vec: r3.Vec{X: x, Y: -y, Z: z}},
		{Vec: r3.Vec{X: x, Y: y, Z: -z}},
		{Vec: r3.Vec{X: x, Y: y, Z: z}},
	}
	edges := [][2]int{
		{0, 1}, {0, 2}, {1, 3}, {2, 3},
		{4, 5}, {4, 6}, {5, 7}, {6, 7},
		{0, 4}, {1, 5}, {2, 6}, {3, 7},
	}
	for _, e := range edges {
		mp.appendLine(world.MulVec3(corners[e[0]]), world.MulVec3(corners[e[1]]))
	}
}

// appendSphereLines は球剛体のワイヤーフレームを追加する。
func (mp *RigidBodyHighlighter) appendSphereLines(world *mmath.Mat4, radius float64) {
	if radius <= 0 {
		return
	}
	mp.appendCircleLines(world, radius, mmath.Vec3{Vec: r3.Vec{X: 1}}, mmath.Vec3{Vec: r3.Vec{Y: 1}})
	mp.appendCircleLines(world, radius, mmath.Vec3{Vec: r3.Vec{Y: 1}}, mmath.Vec3{Vec: r3.Vec{Z: 1}})
	mp.appendCircleLines(world, radius, mmath.Vec3{Vec: r3.Vec{X: 1}}, mmath.Vec3{Vec: r3.Vec{Z: 1}})
}

// appendCapsuleLines はカプセル剛体のワイヤーフレームを追加する。
func (mp *RigidBodyHighlighter) appendCapsuleLines(world *mmath.Mat4, radius, halfHeight float64) {
	if radius <= 0 {
		return
	}
	if halfHeight < 0 {
		halfHeight = 0
	}
	top := mmath.Vec3{Vec: r3.Vec{Y: halfHeight}}
	bottom := mmath.Vec3{Vec: r3.Vec{Y: -halfHeight}}

	mp.appendCircleLinesAt(world, radius, top, mmath.Vec3{Vec: r3.Vec{X: 1}}, mmath.Vec3{Vec: r3.Vec{Z: 1}})
	mp.appendCircleLinesAt(world, radius, bottom, mmath.Vec3{Vec: r3.Vec{X: 1}}, mmath.Vec3{Vec: r3.Vec{Z: 1}})

	axisOffsets := []mmath.Vec3{
		{Vec: r3.Vec{X: radius}},
		{Vec: r3.Vec{X: -radius}},
		{Vec: r3.Vec{Z: radius}},
		{Vec: r3.Vec{Z: -radius}},
	}
	for _, offset := range axisOffsets {
		mp.appendLine(world.MulVec3(offset.Added(bottom)), world.MulVec3(offset.Added(top)))
	}
}

// appendCircleLines は任意平面の円周ラインを追加する。
func (mp *RigidBodyHighlighter) appendCircleLines(
	world *mmath.Mat4,
	radius float64,
	axisA, axisB mmath.Vec3,
) {
	mp.appendCircleLinesAt(world, radius, mmath.NewVec3(), axisA, axisB)
}

// appendCircleLinesAt は中心指定で円周ラインを追加する。
func (mp *RigidBodyHighlighter) appendCircleLinesAt(
	world *mmath.Mat4,
	radius float64,
	center mmath.Vec3,
	axisA, axisB mmath.Vec3,
) {
	segments := rigidBodyHighlightSegments
	for i := 0; i < segments; i++ {
		theta0 := 2 * math.Pi * float64(i) / float64(segments)
		theta1 := 2 * math.Pi * float64(i+1) / float64(segments)
		p0 := axisA.MuledScalar(math.Cos(theta0) * radius).Added(axisB.MuledScalar(math.Sin(theta0) * radius)).Added(center)
		p1 := axisA.MuledScalar(math.Cos(theta1) * radius).Added(axisB.MuledScalar(math.Sin(theta1) * radius)).Added(center)
		mp.appendLine(world.MulVec3(p0), world.MulVec3(p1))
	}
}

// appendLine はハイライト用ラインを追加する。
func (mp *RigidBodyHighlighter) appendLine(from, to mmath.Vec3) {
	// OpenGL座標系に合わせてX軸を反転する。
	from.X = -from.X
	to.X = -to.X
	mp.highlightVertices = append(mp.highlightVertices,
		float32(from.X), float32(from.Y), float32(from.Z),
		rigidBodyHighlightColorR, rigidBodyHighlightColorG, rigidBodyHighlightColorB, rigidBodyHighlightColorA,
		float32(to.X), float32(to.Y), float32(to.Z),
		rigidBodyHighlightColorR, rigidBodyHighlightColorG, rigidBodyHighlightColorB, rigidBodyHighlightColorA,
	)
}
