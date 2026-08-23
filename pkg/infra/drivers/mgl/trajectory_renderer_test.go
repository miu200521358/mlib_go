//go:build windows
// +build windows

// 指示: miu200521358
package mgl

import (
	"testing"

	"github.com/miu200521358/mlib_go/pkg/shared/contracts/mtime"
	"github.com/miu200521358/mlib_go/pkg/shared/state"
)

// TestBuildTrajectoryLineVertices は接地線を下地だけに追加することを確認する。
func TestBuildTrajectoryLineVertices(t *testing.T) {
	points := []state.TrajectoryPoint{
		{Frame: 0, Position: [3]float32{1, 0, 0}, Grounded: true},
		{Frame: 1, Position: [3]float32{2, 0, 0}, Grounded: true},
		{Frame: 2, Position: [3]float32{3, 0, 0}, Grounded: false},
	}
	vertices, groundCount, baseCount := buildTrajectoryLineVertices([]state.TrajectoryPolyline{{Points: points}})
	if groundCount != 6 || baseCount != 12 {
		t.Fatalf("vertex counts = ground:%d base:%d, want 6/12", groundCount, baseCount)
	}
	if len(vertices) != int(groundCount+baseCount)*trajectoryVertexFloatCount {
		t.Fatalf("float count = %d", len(vertices))
	}
	if got := vertices[0]; got != -1 {
		t.Fatalf("OpenGL X conversion = %v, want -1", got)
	}
}

// TestBuildTrajectoryMarkerVertices は現在 frame の補間位置と marker 設定を確認する。
func TestBuildTrajectoryMarkerVertices(t *testing.T) {
	color := state.TrajectoryColor{R: 1, G: 0.5, A: 1}
	polyline := state.TrajectoryPolyline{
		CurrentMarkerSize: 9, CurrentMarkerColor: color,
		Points: []state.TrajectoryPoint{
			{Frame: 2, Position: [3]float32{2, 4, 6}},
			{Frame: 6, Position: [3]float32{6, 8, 10}},
		},
	}
	vertices := buildTrajectoryMarkerVertices([]state.TrajectoryPolyline{polyline}, mtime.Frame(4))
	if len(vertices) != 6*trajectoryVertexFloatCount {
		t.Fatalf("marker float count = %d", len(vertices))
	}
	if vertices[0] != -4 || vertices[1] != 6 || vertices[2] != 8 {
		t.Fatalf("marker position = %v", vertices[:3])
	}
	if vertices[6] != color.R || vertices[7] != color.G || vertices[16] != 9 || vertices[17] != 1 {
		t.Fatalf("marker attributes = %v", vertices[6:18])
	}
	if got := buildTrajectoryMarkerVertices([]state.TrajectoryPolyline{polyline}, mtime.Frame(7)); len(got) != 0 {
		t.Fatalf("outside marker = %v, want nil", got)
	}
}
