//go:build windows
// +build windows

// 指示: miu200521358
package mgl

import (
	"math"

	"github.com/go-gl/gl/v4.3-core/gl"

	"github.com/miu200521358/mlib_go/pkg/adapter/graphics_api"
	"github.com/miu200521358/mlib_go/pkg/shared/contracts/mtime"
	"github.com/miu200521358/mlib_go/pkg/shared/state"
)

const (
	trajectoryVertexFloatCount = 18
	defaultTrajectoryWidth     = float32(2)
	defaultGroundedWidth       = float32(7)
	defaultCurrentMarkerSize   = float32(11)
)

var trajectoryQuadCorners = [][2]float32{
	{0, -1}, {0, 1}, {1, 1},
	{0, -1}, {1, 1}, {1, -1},
}

var trajectoryMarkerCorners = [][2]float32{
	{-1, -1}, {-1, 1}, {1, 1},
	{-1, -1}, {1, 1}, {1, -1},
}

// TrajectoryRenderer は viewer 固有の軌跡 VBO と描画版を保持する。
type TrajectoryRenderer struct {
	lineBuffer   *VertexBufferHandle
	markerBuffer *VertexBufferHandle
	polylines    []state.TrajectoryPolyline
	version      uint64
	groundCount  int32
	baseCount    int32
}

// NewTrajectoryRenderer は空の軌跡レンダラーを生成する。
func NewTrajectoryRenderer() *TrajectoryRenderer {
	return &TrajectoryRenderer{}
}

// Render は共有状態の軌跡を同期し、接地下地、通常線、現在点の順に描画する。
func (r *TrajectoryRenderer) Render(
	shader graphics_api.IShader,
	shared *state.SharedState,
	viewerIndex, viewportWidth, viewportHeight int,
) {
	if r == nil || shader == nil || shared == nil || viewportWidth <= 0 || viewportHeight <= 0 {
		return
	}
	r.syncLines(shared, viewerIndex)
	markerVertices := buildTrajectoryMarkerVertices(r.polylines, shared.Frame())
	if r.groundCount+r.baseCount == 0 && len(markerVertices) == 0 {
		return
	}
	if len(markerVertices) > 0 {
		r.markerBuffer = updateTrajectoryBuffer(r.markerBuffer, markerVertices)
	}

	shader.UseProgram(graphics_api.ProgramTypeTrajectory)
	program := shader.Program(graphics_api.ProgramTypeTrajectory)
	viewportUniform := GetUniformLocation(program, "viewportSize\x00")
	gl.Uniform2f(viewportUniform, float32(viewportWidth), float32(viewportHeight))
	if r.lineBuffer != nil {
		r.lineBuffer.Bind()
		if r.groundCount > 0 {
			gl.DrawArrays(gl.TRIANGLES, 0, r.groundCount)
		}
		if r.baseCount > 0 {
			gl.DrawArrays(gl.TRIANGLES, r.groundCount, r.baseCount)
		}
		r.lineBuffer.Unbind()
	}
	if len(markerVertices) > 0 && r.markerBuffer != nil {
		r.markerBuffer.Bind()
		gl.DrawArrays(gl.TRIANGLES, 0, int32(len(markerVertices)/trajectoryVertexFloatCount))
		r.markerBuffer.Unbind()
	}
	shader.ResetProgram()
}

// Delete は軌跡用 OpenGL リソースを解放する。
func (r *TrajectoryRenderer) Delete() {
	if r == nil {
		return
	}
	if r.lineBuffer != nil {
		r.lineBuffer.Delete()
		r.lineBuffer = nil
	}
	if r.markerBuffer != nil {
		r.markerBuffer.Delete()
		r.markerBuffer = nil
	}
}

// syncLines は state の版が変わったときだけ線頂点を再生成する。
func (r *TrajectoryRenderer) syncLines(shared *state.SharedState, viewerIndex int) {
	polylines, version := shared.TrajectoryPolylinesWithVersion(viewerIndex)
	if version == r.version {
		return
	}
	r.version = version
	r.polylines = polylines
	vertices, groundCount, baseCount := buildTrajectoryLineVertices(polylines)
	r.groundCount = groundCount
	r.baseCount = baseCount
	if len(vertices) > 0 {
		r.lineBuffer = updateTrajectoryBuffer(r.lineBuffer, vertices)
	}
}

// updateTrajectoryBuffer は同一レイアウトの VBO を作成または置換する。
func updateTrajectoryBuffer(buffer *VertexBufferHandle, vertices []float32) *VertexBufferHandle {
	if len(vertices) == 0 {
		return buffer
	}
	if buffer == nil {
		return NewVertexBufferBuilder().
			AddAttributeWithSize(AttributePosition, 3).
			AddAttributeWithSize(AttributePosition, 3).
			AddAttributeWithSize(AttributeColor, 4).
			AddAttributeWithSize(AttributeColor, 4).
			AddAttributeWithSize(AttributePosition, 2).
			AddAttributeWithSize(AttributeEdgeFactor, 1).
			AddAttributeWithSize(AttributeEdgeFactor, 1).
			SetData(gl.Ptr(vertices), len(vertices)).
			Build()
	}
	buffer.Bind()
	buffer.VBO.BufferDataWithLayout(
		len(vertices), trajectoryVertexFloatCount, buffer.FloatSize, gl.Ptr(vertices), graphics_api.BufferUsageDynamic,
	)
	buffer.Unbind()
	return buffer
}

// buildTrajectoryLineVertices は接地下地を前、通常線を後に並べた三角形頂点を作る。
func buildTrajectoryLineVertices(polylines []state.TrajectoryPolyline) ([]float32, int32, int32) {
	groundVertices := make([]float32, 0)
	baseVertices := make([]float32, 0)
	for _, polyline := range polylines {
		groundedWidth := positiveOrDefault(polyline.GroundedWidth, defaultGroundedWidth)
		baseWidth := positiveOrDefault(polyline.Width, defaultTrajectoryWidth)
		for i := 1; i < len(polyline.Points); i++ {
			start := polyline.Points[i-1]
			end := polyline.Points[i]
			if start.Grounded && end.Grounded {
				groundVertices = appendTrajectoryQuad(
					groundVertices, start, end, start.GroundColor, end.GroundColor, groundedWidth, false,
				)
			}
			baseVertices = appendTrajectoryQuad(
				baseVertices, start, end, start.Color, end.Color, baseWidth, false,
			)
		}
	}
	groundCount := int32(len(groundVertices) / trajectoryVertexFloatCount)
	baseCount := int32(len(baseVertices) / trajectoryVertexFloatCount)
	return append(groundVertices, baseVertices...), groundCount, baseCount
}

// buildTrajectoryMarkerVertices は共有 frame を表示折れ線上へ補間した円形 marker 頂点を作る。
func buildTrajectoryMarkerVertices(polylines []state.TrajectoryPolyline, frame mtime.Frame) []float32 {
	vertices := make([]float32, 0)
	for _, polyline := range polylines {
		point, ok := interpolateTrajectoryPoint(polyline.Points, frame)
		if !ok {
			continue
		}
		markerSize := positiveOrDefault(polyline.CurrentMarkerSize, defaultCurrentMarkerSize)
		markerColor := polyline.CurrentMarkerColor
		vertices = appendTrajectoryQuad(vertices, point, point, markerColor, markerColor, markerSize, true)
	}
	return vertices
}

// interpolateTrajectoryPoint は frame を挟む表示点間を線形補間する。
func interpolateTrajectoryPoint(points []state.TrajectoryPoint, frame mtime.Frame) (state.TrajectoryPoint, bool) {
	if len(points) == 0 || frame < points[0].Frame || frame > points[len(points)-1].Frame {
		return state.TrajectoryPoint{}, false
	}
	for i := range points {
		if points[i].Frame == frame {
			return points[i], true
		}
		if i == 0 || points[i].Frame < frame {
			continue
		}
		start := points[i-1]
		end := points[i]
		span := float64(end.Frame - start.Frame)
		if span <= 0 {
			return state.TrajectoryPoint{}, false
		}
		ratio := float32(float64(frame-start.Frame) / span)
		out := start
		out.Frame = frame
		for axis := range out.Position {
			out.Position[axis] = start.Position[axis] + (end.Position[axis]-start.Position[axis])*ratio
		}
		out.Color = interpolateTrajectoryColor(start.Color, end.Color, ratio)
		return out, true
	}
	return state.TrajectoryPoint{}, false
}

// interpolateTrajectoryColor は RGBA を線形補間する。
func interpolateTrajectoryColor(start, end state.TrajectoryColor, ratio float32) state.TrajectoryColor {
	return state.TrajectoryColor{
		R: start.R + (end.R-start.R)*ratio,
		G: start.G + (end.G-start.G)*ratio,
		B: start.B + (end.B-start.B)*ratio,
		A: start.A + (end.A-start.A)*ratio,
	}
}

// appendTrajectoryQuad は一線分または marker を二三角形として追加する。
func appendTrajectoryQuad(
	vertices []float32,
	start, end state.TrajectoryPoint,
	startColor, endColor state.TrajectoryColor,
	size float32,
	marker bool,
) []float32 {
	corners := trajectoryQuadCorners
	markerValue := float32(0)
	if marker {
		corners = trajectoryMarkerCorners
		markerValue = 1
	}
	for _, corner := range corners {
		vertices = append(vertices,
			-start.Position[0], start.Position[1], start.Position[2],
			-end.Position[0], end.Position[1], end.Position[2],
			startColor.R, startColor.G, startColor.B, startColor.A,
			endColor.R, endColor.G, endColor.B, endColor.A,
			corner[0], corner[1], size, markerValue,
		)
	}
	return vertices
}

// positiveOrDefault は正の有限値だけを採用する。
func positiveOrDefault(value, fallback float32) float32 {
	if value <= 0 || math.IsNaN(float64(value)) || math.IsInf(float64(value), 0) {
		return fallback
	}
	return value
}
