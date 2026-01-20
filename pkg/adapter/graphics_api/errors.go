// 指示: miu200521358
package graphics_api

import "github.com/miu200521358/mlib_go/pkg/shared/base/merr"

const (
	shaderSourceLoadFailedErrorID    = "94501"
	shaderCompileFailedErrorID       = "94502"
	shaderLinkFailedErrorID          = "94503"
	framebufferIncompleteErrorID     = "94504"
	openGlErrorID                    = "94505"
	graphicsContextInitFailedErrorID = "94506"
)

// NewShaderSourceLoadFailed はシェーダソース読込失敗エラーを生成する。
func NewShaderSourceLoadFailed(message string, cause error) error {
	return merr.NewCommonError(shaderSourceLoadFailedErrorID, merr.ErrorKindInternal, message, cause)
}

// NewShaderCompileFailed はシェーダコンパイル失敗エラーを生成する。
func NewShaderCompileFailed(message string, cause error) error {
	return merr.NewCommonError(shaderCompileFailedErrorID, merr.ErrorKindInternal, message, cause)
}

// NewShaderLinkFailed はシェーダリンク失敗エラーを生成する。
func NewShaderLinkFailed(message string, cause error) error {
	return merr.NewCommonError(shaderLinkFailedErrorID, merr.ErrorKindInternal, message, cause)
}

// NewFramebufferIncomplete はフレームバッファ不完全エラーを生成する。
func NewFramebufferIncomplete(message string, cause error) error {
	return merr.NewCommonError(framebufferIncompleteErrorID, merr.ErrorKindInternal, message, cause)
}

// NewOpenGLError はOpenGLエラーを生成する。
func NewOpenGLError(message string, cause error) error {
	return merr.NewCommonError(openGlErrorID, merr.ErrorKindInternal, message, cause)
}

// NewGraphicsContextInitFailed はOpenGLコンテキスト初期化失敗エラーを生成する。
func NewGraphicsContextInitFailed(message string, cause error) error {
	return merr.NewCommonError(graphicsContextInitFailedErrorID, merr.ErrorKindInternal, message, cause)
}
