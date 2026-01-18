package graphics_api

import baseerr "github.com/miu200521358/mlib_go/pkg/shared/base/err"

const (
	shaderSourceLoadFailedErrorID   = "94501"
	shaderCompileFailedErrorID      = "94502"
	shaderLinkFailedErrorID         = "94503"
	framebufferIncompleteErrorID    = "94504"
	openGlErrorID                   = "94505"
	graphicsContextInitFailedErrorID = "94506"
)

// NewShaderSourceLoadFailed はシェーダソース読込失敗エラーを生成する。
func NewShaderSourceLoadFailed(message string, cause error) error {
	return baseerr.NewCommonError(shaderSourceLoadFailedErrorID, baseerr.ErrorKindInternal, message, cause)
}

// NewShaderCompileFailed はシェーダコンパイル失敗エラーを生成する。
func NewShaderCompileFailed(message string, cause error) error {
	return baseerr.NewCommonError(shaderCompileFailedErrorID, baseerr.ErrorKindInternal, message, cause)
}

// NewShaderLinkFailed はシェーダリンク失敗エラーを生成する。
func NewShaderLinkFailed(message string, cause error) error {
	return baseerr.NewCommonError(shaderLinkFailedErrorID, baseerr.ErrorKindInternal, message, cause)
}

// NewFramebufferIncomplete はフレームバッファ不完全エラーを生成する。
func NewFramebufferIncomplete(message string, cause error) error {
	return baseerr.NewCommonError(framebufferIncompleteErrorID, baseerr.ErrorKindInternal, message, cause)
}

// NewOpenGLError はOpenGLエラーを生成する。
func NewOpenGLError(message string, cause error) error {
	return baseerr.NewCommonError(openGlErrorID, baseerr.ErrorKindInternal, message, cause)
}

// NewGraphicsContextInitFailed はOpenGLコンテキスト初期化失敗エラーを生成する。
func NewGraphicsContextInitFailed(message string, cause error) error {
	return baseerr.NewCommonError(graphicsContextInitFailedErrorID, baseerr.ErrorKindInternal, message, cause)
}
