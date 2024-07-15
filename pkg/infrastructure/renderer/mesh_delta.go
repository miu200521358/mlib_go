package renderer

import "github.com/go-gl/mathgl/mgl32"

type MeshDelta struct {
	Diffuse  mgl32.Vec4
	Specular mgl32.Vec4
	Ambient  mgl32.Vec3
	Edge     mgl32.Vec4
	EdgeSize float32
}
