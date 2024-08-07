//go:build windows
// +build windows

package mgl

import (
	"github.com/go-gl/gl/v4.4-core/gl"
)

type MFloor struct {
	vao   *VAO
	vbo   *VBO
	count int32
}

func newMFloor() *MFloor {
	mf := &MFloor{}

	mf.vao = NewVAO()
	mf.vao.Bind()
	mf.vbo, mf.count = NewVBOForFloor()
	mf.vbo.Unbind()
	mf.vao.Unbind()

	return mf
}

func (shader *MShader) DrawFloor() {
	// mlog.D("MFloor.DrawLine")
	program := shader.Program(PROGRAM_TYPE_FLOOR)
	gl.UseProgram(program)

	// 平面を引く
	shader.floor.vao.Bind()
	shader.floor.vbo.BindFloor()

	gl.DrawArrays(gl.LINES, 0, shader.floor.count)

	shader.floor.vbo.Unbind()
	shader.floor.vao.Unbind()

	gl.UseProgram(0)
}
