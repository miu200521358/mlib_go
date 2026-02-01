//go:build windows
// +build windows

// 指示: miu200521358
package mgl

import (
	"fmt"

	"github.com/go-gl/gl/v4.3-core/gl"
	"github.com/miu200521358/mlib_go/pkg/adapter/graphics_api"
)

// getFrameBufferStatusString はフレームバッファステータスコードを文字列に変換する。
func getFrameBufferStatusString(status uint32) string {
	switch status {
	case gl.FRAMEBUFFER_COMPLETE:
		return "FRAMEBUFFER_COMPLETE"
	case gl.FRAMEBUFFER_INCOMPLETE_ATTACHMENT:
		return "FRAMEBUFFER_INCOMPLETE_ATTACHMENT"
	case gl.FRAMEBUFFER_INCOMPLETE_MISSING_ATTACHMENT:
		return "FRAMEBUFFER_INCOMPLETE_MISSING_ATTACHMENT"
	case gl.FRAMEBUFFER_INCOMPLETE_DRAW_BUFFER:
		return "FRAMEBUFFER_INCOMPLETE_DRAW_BUFFER"
	case gl.FRAMEBUFFER_INCOMPLETE_READ_BUFFER:
		return "FRAMEBUFFER_INCOMPLETE_READ_BUFFER"
	case gl.FRAMEBUFFER_UNSUPPORTED:
		return "FRAMEBUFFER_UNSUPPORTED"
	case gl.FRAMEBUFFER_INCOMPLETE_MULTISAMPLE:
		return "FRAMEBUFFER_INCOMPLETE_MULTISAMPLE"
	case gl.FRAMEBUFFER_INCOMPLETE_LAYER_TARGETS:
		return "FRAMEBUFFER_INCOMPLETE_LAYER_TARGETS"
	default:
		return "UNKNOWN_ERROR"
	}
}

// getOpenGLErrorString はOpenGLエラーコードを文字列に変換する。
func getOpenGLErrorString(errCode uint32) string {
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

// CheckGLError はOpenGLのエラー状態を検査してエラーを返す。
func CheckGLError() error {
	if errCode := gl.GetError(); errCode != gl.NO_ERROR {
		return graphics_api.NewOpenGLError(
			"OpenGLエラーが発生しました: code=%d detail=%s",
			nil,
			errCode,
			getOpenGLErrorString(errCode),
		)
	}
	return nil
}

// HasGLError はOpenGLエラーの有無を返す。
func HasGLError() bool {
	return gl.GetError() != gl.NO_ERROR
}
