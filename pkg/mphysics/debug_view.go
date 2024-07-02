//go:build windows
// +build windows

package mphysics

import (
	"github.com/go-gl/gl/v4.4-core/gl"

	"github.com/miu200521358/mlib_go/pkg/mphysics/mbt"
	"github.com/miu200521358/mlib_go/pkg/mview"
)

func NewConstBtMDefaultColors() mbt.BtMDefaultColors {
	return mbt.NewBtMDefaultColors(
		mbt.NewBtVector3(float32(1.0), float32(0.0), float32(0.0)), // activeObject	(物理剛体)
		mbt.NewBtVector3(float32(0.5), float32(0.5), float32(0.0)), // deactivatedObject
		mbt.NewBtVector3(float32(0.5), float32(0.0), float32(0.5)), // wantsDeactivationObject
		mbt.NewBtVector3(float32(0.0), float32(0.5), float32(0.5)), // disabledDeactivationObject
		mbt.NewBtVector3(float32(0.0), float32(1.0), float32(0.0)), // disabledSimulationObject	(ボーン追従剛体)
		mbt.NewBtVector3(float32(1.0), float32(1.0), float32(0.0)), // aabb
		mbt.NewBtVector3(float32(0.0), float32(0.0), float32(1.0)), // contactPoint
	)
}

type MDebugDrawLiner struct {
	mbt.BtMDebugDrawLiner
	shader *mview.MShader
	vao    *mview.VAO
	vbo    *mview.VBO
}

func NewMDebugDrawLiner(shader *mview.MShader) *MDebugDrawLiner {
	ddl := &MDebugDrawLiner{
		shader: shader,
	}
	ddl.BtMDebugDrawLiner = mbt.NewDirectorBtMDebugDrawLiner(ddl)

	ddl.vao = mview.NewVAO()
	ddl.vao.Bind()
	ddl.vbo = mview.NewVBOForDebug()
	ddl.vbo.BindDebug(mbt.NewBtVector3(), mbt.NewBtVector3())
	ddl.vbo.Unbind()
	ddl.vao.Unbind()

	return ddl
}

func (ddl MDebugDrawLiner) DrawLine(from mbt.BtVector3, to mbt.BtVector3, color mbt.BtVector3) {
	if ddl.shader.IsDrawRigidBodyFront {
		// モデルメッシュの前面に描画するために深度テストを無効化
		gl.Enable(gl.DEPTH_TEST)
		gl.DepthFunc(gl.ALWAYS)
	}

	// mlog.D("MDebugDrawLiner.DrawLine")
	ddl.shader.Use(mview.PROGRAM_TYPE_PHYSICS)

	// 色を設定
	colorUniform := gl.GetUniformLocation(ddl.shader.PhysicsProgram, gl.Str(mview.SHADER_COLOR))
	gl.Uniform3f(colorUniform, color.GetX(), color.GetY(), color.GetZ())

	alphaUniform := gl.GetUniformLocation(ddl.shader.PhysicsProgram, gl.Str(mview.SHADER_ALPHA))
	gl.Uniform1f(alphaUniform, 0.6)

	// 線を引く
	ddl.vao.Bind()
	ddl.vbo.BindDebug(from, to)

	gl.DrawArrays(gl.LINES, 0, int32(2))

	ddl.vbo.Unbind()
	ddl.vao.Unbind()

	ddl.shader.Unuse()

	// 深度テストを有効に戻す
	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LEQUAL)
}
