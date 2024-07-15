package renderer

import "github.com/go-gl/mathgl/mgl32"

type MeshDelta struct {
	Diffuse          mgl32.Vec4
	Specular         mgl32.Vec4
	Ambient          mgl32.Vec3
	Edge             mgl32.Vec4
	EdgeSize         float32
	TextureMulFactor mgl32.Vec4
	TextureAddFactor mgl32.Vec4
	SphereMulFactor  mgl32.Vec4
	SphereAddFactor  mgl32.Vec4
	ToonMulFactor    mgl32.Vec4
	ToonAddFactor    mgl32.Vec4
}
