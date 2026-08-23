//go:build windows
// +build windows

// 指示: miu200521358
package mgl

import (
	"math"

	"github.com/go-gl/gl/v4.3-core/gl"
	"github.com/go-gl/mathgl/mgl32"

	"github.com/miu200521358/mlib_go/pkg/adapter/graphics_api"
	"github.com/miu200521358/mlib_go/pkg/domain/mmath"
	"github.com/miu200521358/mlib_go/pkg/shared/state"
)

const operationPointSize = float32(17)

var (
	operationPointColor       = state.TrajectoryColor{R: 1.0, G: 0.72, B: 0.15, A: 0.95}
	operationPointActiveColor = state.TrajectoryColor{R: 0.20, G: 0.90, B: 1.0, A: 1.0}
)

// OperationPointRenderPoint は操作点の描画位置とドラッグ状態を表す。
type OperationPointRenderPoint struct {
	Position mmath.Vec3
	Active   bool
}

// OperationPointRenderer は画面ピクセル幅を保つ操作点ハンドルを描画する。
type OperationPointRenderer struct {
	buffer *VertexBufferHandle
}

// NewOperationPointRenderer は空の操作点レンダラーを生成する。
func NewOperationPointRenderer() *OperationPointRenderer {
	return &OperationPointRenderer{}
}

// Render は操作点をモデルより前面に円形ハンドルとして描画する。
func (r *OperationPointRenderer) Render(
	shader graphics_api.IShader,
	points []OperationPointRenderPoint,
	viewportWidth, viewportHeight int,
) {
	if r == nil || shader == nil || len(points) == 0 || viewportWidth <= 0 || viewportHeight <= 0 {
		return
	}
	vertices := buildOperationPointVertices(points)
	if len(vertices) == 0 {
		return
	}
	r.buffer = updateTrajectoryBuffer(r.buffer, vertices)
	if r.buffer == nil {
		return
	}

	shader.UseProgram(graphics_api.ProgramTypeTrajectory)
	program := shader.Program(graphics_api.ProgramTypeTrajectory)
	viewportUniform := GetUniformLocation(program, "viewportSize\x00")
	gl.Uniform2f(viewportUniform, float32(viewportWidth), float32(viewportHeight))
	// 操作点はボーン位置に常時見える必要があるため深度判定だけ一時停止し、後続描画へ状態を残さない。
	depthEnabled := gl.IsEnabled(gl.DEPTH_TEST)
	gl.Disable(gl.DEPTH_TEST)
	r.buffer.Bind()
	gl.DrawArrays(gl.TRIANGLES, 0, int32(len(vertices)/trajectoryVertexFloatCount))
	r.buffer.Unbind()
	if depthEnabled {
		gl.Enable(gl.DEPTH_TEST)
	}
	shader.ResetProgram()
}

// Delete は操作点用 OpenGL リソースを解放する。
func (r *OperationPointRenderer) Delete() {
	if r == nil || r.buffer == nil {
		return
	}
	r.buffer.Delete()
	r.buffer = nil
}

// ScreenToCameraFacingPlane は画面位置のレイと、指定点を通るカメラ正対平面の交点を返す。
// planeNormal はカメラの視線方向を渡し、ドラッグ中の奥行きを固定して 2D 入力を一意な 3D 位置へ写す。
func ScreenToCameraFacingPlane(
	screenX, screenY float64,
	width, height int,
	view, projection mgl32.Mat4,
	planePoint, planeNormal mmath.Vec3,
) (mmath.Vec3, bool) {
	if width <= 0 || height <= 0 || !finiteOperationPointValues(screenX, screenY) || planeNormal.LengthSqr() <= 1e-12 {
		return mmath.Vec3{}, false
	}
	nearWorld, nearErr := mgl32.UnProject(
		mgl32.Vec3{float32(screenX), float32(height) - float32(screenY), 0},
		view, projection, 0, 0, width, height,
	)
	farWorld, farErr := mgl32.UnProject(
		mgl32.Vec3{float32(screenX), float32(height) - float32(screenY), 1},
		view, projection, 0, 0, width, height,
	)
	if nearErr != nil || farErr != nil {
		return mmath.Vec3{}, false
	}
	glPoint := NewGlVec3(&planePoint)
	glNormal := NewGlVec3(&planeNormal)
	direction := farWorld.Sub(nearWorld)
	denominator := glNormal.Dot(direction)
	if math.Abs(float64(denominator)) <= 1e-7 {
		return mmath.Vec3{}, false
	}
	ratio := glNormal.Dot(glPoint.Sub(nearWorld)) / denominator
	intersection := nearWorld.Add(direction.Mul(ratio))
	out := mmath.Vec3{}
	out.X = -float64(intersection.X())
	out.Y = float64(intersection.Y())
	out.Z = float64(intersection.Z())
	if !finiteOperationPointValues(out.X, out.Y, out.Z) {
		return mmath.Vec3{}, false
	}
	return out, true
}

// buildOperationPointVertices は操作点を軌跡 marker と同じ頂点形式へ変換する。
func buildOperationPointVertices(points []OperationPointRenderPoint) []float32 {
	vertices := make([]float32, 0, len(points)*len(trajectoryMarkerCorners)*trajectoryVertexFloatCount)
	for _, point := range points {
		color := operationPointColor
		if point.Active {
			color = operationPointActiveColor
		}
		marker := state.TrajectoryPoint{Position: [3]float32{
			float32(point.Position.X), float32(point.Position.Y), float32(point.Position.Z),
		}}
		vertices = appendTrajectoryQuad(vertices, marker, marker, color, color, operationPointSize, true)
	}
	return vertices
}

// finiteOperationPointValues は値が有限か判定する。
func finiteOperationPointValues(values ...float64) bool {
	for _, value := range values {
		if math.IsNaN(value) || math.IsInf(value, 0) {
			return false
		}
	}
	return true
}
