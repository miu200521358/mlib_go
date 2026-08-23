//go:build windows
// +build windows

// 指示: miu200521358
package mgl

import (
	"math"
	"testing"

	"github.com/go-gl/mathgl/mgl32"

	"github.com/miu200521358/mlib_go/pkg/domain/mmath"
)

func TestScreenToCameraFacingPlane(t *testing.T) {
	view := mgl32.LookAtV(mgl32.Vec3{0, 0, 10}, mgl32.Vec3{}, mgl32.Vec3{0, 1, 0})
	projection := mgl32.Perspective(mgl32.DegToRad(60), 1, 0.1, 100)
	planePoint := mmath.Vec3{}
	planeNormal := mmath.Vec3{}
	planeNormal.Z = -1

	center, ok := ScreenToCameraFacingPlane(50, 50, 100, 100, view, projection, planePoint, planeNormal)
	if !ok || center.Distance(planePoint) > 1e-4 {
		t.Fatalf("中心の交点 = %+v, ok=%t", center, ok)
	}
	hidpi, ok := ScreenToCameraFacingPlane(100, 100, 200, 200, view, projection, planePoint, planeNormal)
	if !ok || hidpi.Distance(center) > 1e-4 {
		t.Fatalf("HiDPI中心の交点 = %+v, ok=%t", hidpi, ok)
	}
	right, ok := ScreenToCameraFacingPlane(75, 50, 100, 100, view, projection, planePoint, planeNormal)
	if !ok || right.X >= 0 || math.Abs(right.Z) > 1e-4 {
		t.Fatalf("右側の交点 = %+v, ok=%t", right, ok)
	}
	if _, ok := ScreenToCameraFacingPlane(50, 50, 0, 100, view, projection, planePoint, planeNormal); ok {
		t.Fatal("幅0を受理しました")
	}
	if _, ok := ScreenToCameraFacingPlane(50, 50, 100, 100, view, projection, planePoint, mmath.Vec3{}); ok {
		t.Fatal("法線0を受理しました")
	}
	parallelNormal := mmath.Vec3{}
	parallelNormal.X = 1
	if _, ok := ScreenToCameraFacingPlane(50, 50, 100, 100, view, projection, planePoint, parallelNormal); ok {
		t.Fatal("平行なレイと平面を受理しました")
	}
	if _, ok := ScreenToCameraFacingPlane(math.NaN(), 50, 100, 100, view, projection, planePoint, planeNormal); ok {
		t.Fatal("NaN座標を受理しました")
	}

	rotatedView := mgl32.LookAtV(mgl32.Vec3{5, 3, 10}, mgl32.Vec3{}, mgl32.Vec3{0, 1, 0})
	rotatedNormal := mmath.Vec3{}
	rotatedNormal.X, rotatedNormal.Y, rotatedNormal.Z = -5, -3, -10
	rotated, ok := ScreenToCameraFacingPlane(100, 50, 200, 100, rotatedView, mgl32.Perspective(mgl32.DegToRad(60), 2, 0.1, 100), planePoint, rotatedNormal)
	if !ok || rotated.Distance(planePoint) > 1e-3 {
		t.Fatalf("回転カメラ中心の交点 = %+v, ok=%t", rotated, ok)
	}
}

func TestBuildOperationPointVertices(t *testing.T) {
	position := mmath.Vec3{}
	position.X, position.Y, position.Z = 1, 2, 3
	vertices := buildOperationPointVertices([]OperationPointRenderPoint{{Position: position, Active: true}})
	if len(vertices) != len(trajectoryMarkerCorners)*trajectoryVertexFloatCount {
		t.Fatalf("頂点数 = %d", len(vertices))
	}
	if vertices[0] != -1 || vertices[1] != 2 || vertices[2] != 3 {
		t.Fatalf("座標 = %v", vertices[:3])
	}
	if vertices[6] != operationPointActiveColor.R || vertices[16] != operationPointSize || vertices[17] != 1 {
		t.Fatalf("属性 = %v", vertices[6:18])
	}
}
