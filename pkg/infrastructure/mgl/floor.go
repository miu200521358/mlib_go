package mgl

import (
	"github.com/go-gl/gl/v2.1/gl"
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

func (s *MShader) DrawFloor() {
	// mlog.D("MFloor.DrawLine")
	program := s.GetProgram(PROGRAM_TYPE_FLOOR)
	gl.UseProgram(program)

	// 平面を引く
	s.floor.vao.Bind()
	s.floor.vbo.BindFloor()

	gl.DrawArrays(gl.LINES, 0, s.floor.count)

	s.floor.vbo.Unbind()
	s.floor.vao.Unbind()

	gl.UseProgram(0)
}
