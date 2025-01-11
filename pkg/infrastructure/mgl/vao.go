//go:build windows
// +build windows

package mgl

import (
	"github.com/go-gl/gl/v4.4-core/gl"
)

// Vertex Array Object.
type VAO struct {
	id uint32 // ID
}

// Creates a new VAO.
// Bind and use VBOs for rendering later.
func NewVAO() *VAO {
	var vaoId uint32
	gl.GenVertexArrays(1, &vaoId)
	vao := &VAO{id: vaoId}

	return vao
}

// Delete this VAO.
func (v *VAO) Delete() {
	gl.DeleteVertexArrays(1, &v.id)
}

// Binds VAO for rendering.
func (v *VAO) Bind() {
	gl.BindVertexArray(v.id)
}

// Unbinds.
func (v *VAO) Unbind() {
	gl.BindVertexArray(0)
}

// Returns the GL ID.
func (v *VAO) GetId() uint32 {
	return v.id
}
