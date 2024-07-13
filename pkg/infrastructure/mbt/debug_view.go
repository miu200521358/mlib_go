//go:build windows
// +build windows

package mbt

import (
	"github.com/go-gl/gl/v4.4-core/gl"

	"github.com/miu200521358/mlib_go/pkg/domain/buffer"
	"github.com/miu200521358/mlib_go/pkg/infrastructure/bt"
	"github.com/miu200521358/mlib_go/pkg/infrastructure/mgl"
)

func NewConstBtMDefaultColors() bt.BtMDefaultColors {
	return bt.NewBtMDefaultColors(
		bt.NewBtVector3(float32(1.0), float32(0.0), float32(0.0)), // activeObject	(物理剛体)
		bt.NewBtVector3(float32(0.5), float32(0.5), float32(0.0)), // deactivatedObject
		bt.NewBtVector3(float32(0.5), float32(0.0), float32(0.5)), // wantsDeactivationObject
		bt.NewBtVector3(float32(0.0), float32(0.5), float32(0.5)), // disabledDeactivationObject
		bt.NewBtVector3(float32(0.0), float32(1.0), float32(0.0)), // disabledSimulationObject	(ボーン追従剛体)
		bt.NewBtVector3(float32(1.0), float32(1.0), float32(0.0)), // aabb
		bt.NewBtVector3(float32(0.0), float32(0.0), float32(1.0)), // contactPoint
	)
}

type MDebugDrawLiner struct {
	bt.BtMDebugDrawLiner
	shader   *mgl.MShader
	vao      *buffer.VAO
	vbo      *buffer.VBO
	vertices []float32
	count    int
}

func NewMDebugDrawLiner(shader *mgl.MShader) *MDebugDrawLiner {
	ddl := &MDebugDrawLiner{
		shader:   shader,
		vertices: make([]float32, 0),
	}

	ddl.vao = buffer.NewVAO()
	ddl.vao.Bind()
	ddl.vbo = buffer.NewVBOForDebug()
	ddl.vbo.Unbind()
	ddl.vao.Unbind()

	ddl.BtMDebugDrawLiner = bt.NewDirectorBtMDebugDrawLiner(ddl)

	return ddl
}

func (ddl *MDebugDrawLiner) DrawLine(from bt.BtVector3, to bt.BtVector3, color bt.BtVector3) {
	ddl.vertices = append(ddl.vertices, from.GetX(), from.GetY(), from.GetZ())
	ddl.vertices = append(ddl.vertices, color.GetX(), color.GetY(), color.GetZ(), 0.6)
	ddl.count++

	ddl.vertices = append(ddl.vertices, to.GetX(), to.GetY(), to.GetZ())
	ddl.vertices = append(ddl.vertices, color.GetX(), color.GetY(), color.GetZ(), 0.6)
	ddl.count++
}

func (ddl *MDebugDrawLiner) DrawDebugLines() {
	if ddl.shader.IsDrawRigidBodyFront {
		// モデルメッシュの前面に描画するために深度テストを無効化
		gl.Enable(gl.DEPTH_TEST)
		gl.DepthFunc(gl.ALWAYS)
	}

	ddl.shader.Use(mgl.PROGRAM_TYPE_PHYSICS)

	// 線を引く
	ddl.vao.Bind()
	ddl.vbo.BindDebug(ddl.vertices)

	// ライン描画
	gl.DrawArrays(gl.LINES, 0, int32(ddl.count))

	ddl.vbo.Unbind()
	ddl.vao.Unbind()

	ddl.shader.Unuse()

	// 深度テストを有効に戻す
	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LEQUAL)
}
