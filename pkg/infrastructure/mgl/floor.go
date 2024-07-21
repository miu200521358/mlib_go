//go:build windows
// +build windows

package mgl

import (
	"github.com/go-gl/gl/v4.4-core/gl"
	"github.com/miu200521358/mlib_go/pkg/infrastructure/mgl/buffer"
)

type MFloor struct {
	vao   *buffer.VAO
	vbo   *buffer.VBO
	count int32
}

func newMFloor() *MFloor {
	mf := &MFloor{}

	mf.vao = buffer.NewVAO()
	mf.vao.Bind()
	mf.vbo, mf.count = buffer.NewVBOForFloor()
	mf.vbo.Unbind()
	mf.vao.Unbind()

	return mf
}

func (shader *MShader) DrawFloor() {
	// mlog.D("MFloor.DrawLine")
	program := shader.Program(PROGRAM_TYPE_FLOOR)
	gl.UseProgram(program)

	windowOpacityUniform := gl.GetUniformLocation(program, gl.Str(SHADER_WINDOW_OPACITY))
	gl.Uniform1f(windowOpacityUniform, shader.WindowOpacity())

	// 平面を引く
	shader.floor.vao.Bind()
	shader.floor.vbo.BindFloor()

	gl.DrawArrays(gl.LINES, 0, shader.floor.count)

	shader.floor.vbo.Unbind()
	shader.floor.vao.Unbind()

	gl.UseProgram(0)
}
