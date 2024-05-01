//go:build !for_linux
// +build !for_linux

package mgl

import (
	"fmt"

	"github.com/go-gl/gl/v4.4-core/gl"
)

func GetOpenGLErrorString(errCode uint32) string {
	switch errCode {
	case gl.NO_ERROR:
		return "No error"
	case gl.INVALID_ENUM:
		return "Invalid enum"
	case gl.INVALID_VALUE:
		return "Invalid value"
	case gl.INVALID_OPERATION:
		return "Invalid operation"
	case gl.STACK_OVERFLOW:
		return "Stack overflow"
	case gl.STACK_UNDERFLOW:
		return "Stack underflow"
	case gl.OUT_OF_MEMORY:
		return "Out of memory"
	case gl.INVALID_FRAMEBUFFER_OPERATION:
		return "Invalid framebuffer operation"
	default:
		return fmt.Sprintf("Unknown error code: %v", errCode)
	}
}

func CheckGLError() error {
	if errCode := gl.GetError(); errCode != gl.NO_ERROR {
		return fmt.Errorf("OpenGL error: %v - %s", errCode, GetOpenGLErrorString(errCode))
	}
	return nil
}
