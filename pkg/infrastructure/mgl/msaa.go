//go:build windows
// +build windows

package mgl

import (
	"github.com/go-gl/gl/v4.4-core/gl"
)

// Msaa はマルチサンプルアンチエイリアス（MSAA）を管理
type Msaa struct {
	fbo     uint32
	texture uint32
	rbo     uint32
	width   int
	height  int
	samples int32
}

// NewMsaa はMSAAを初期化
func NewMsaa(width, height int) *Msaa {
	msaa := &Msaa{
		width:   width,
		height:  height,
		samples: 4,
	}

	// MSAA用のフレームバッファを作成
	gl.GenFramebuffers(1, &msaa.fbo)
	gl.BindFramebuffer(gl.FRAMEBUFFER, msaa.fbo)

	// マルチサンプルカラーテクスチャを作成
	gl.GenTextures(1, &msaa.texture)
	gl.BindTexture(gl.TEXTURE_2D_MULTISAMPLE, msaa.texture)
	gl.TexImage2DMultisample(gl.TEXTURE_2D_MULTISAMPLE, msaa.samples, gl.RGBA8, int32(msaa.width), int32(msaa.height), true)
	gl.FramebufferTexture2D(gl.FRAMEBUFFER, gl.COLOR_ATTACHMENT0, gl.TEXTURE_2D_MULTISAMPLE, msaa.texture, 0)

	// マルチサンプルレンダーバッファ（深度・ステンシル用）を作成
	gl.GenRenderbuffers(1, &msaa.rbo)
	gl.BindRenderbuffer(gl.RENDERBUFFER, msaa.rbo)
	gl.RenderbufferStorageMultisample(gl.RENDERBUFFER, msaa.samples, gl.DEPTH24_STENCIL8, int32(msaa.width), int32(msaa.height))
	gl.FramebufferRenderbuffer(gl.FRAMEBUFFER, gl.DEPTH_STENCIL_ATTACHMENT, gl.RENDERBUFFER, msaa.rbo)

	// フレームバッファの整合性チェック
	if gl.CheckFramebufferStatus(gl.FRAMEBUFFER) != gl.FRAMEBUFFER_COMPLETE {
		panic("MSAA Framebuffer is not complete")
	}

	gl.BindFramebuffer(gl.FRAMEBUFFER, 0)

	return msaa
}

// Bind はMSAAフレームバッファをバインド
func (msaa *Msaa) Bind() {
	gl.BindFramebuffer(gl.FRAMEBUFFER, msaa.fbo)
}

// Unbind はデフォルトフレームバッファに戻す
func (msaa *Msaa) Unbind() {
	gl.BindFramebuffer(gl.FRAMEBUFFER, 0)
}

// Resolve はMSAAバッファの内容を通常のフレームバッファへ転送
func (msaa *Msaa) Resolve() {
	gl.BindFramebuffer(gl.READ_FRAMEBUFFER, msaa.fbo)
	gl.BindFramebuffer(gl.DRAW_FRAMEBUFFER, 0)
	gl.BlitFramebuffer(0, 0, int32(msaa.width), int32(msaa.height), 0, 0, int32(msaa.width), int32(msaa.height), gl.COLOR_BUFFER_BIT, gl.NEAREST)
	gl.BindFramebuffer(gl.FRAMEBUFFER, 0)
}

// Delete はリソースを解放
func (msaa *Msaa) Delete() {
	gl.DeleteFramebuffers(1, &msaa.fbo)
	gl.DeleteTextures(1, &msaa.texture)
	gl.DeleteRenderbuffers(1, &msaa.rbo)
}

func (msaa *Msaa) Resize(width, height int) {
	if msaa.width == width && msaa.height == height {
		return
	}

	msaa.width = width
	msaa.height = height

	// マルチサンプルカラーテクスチャを再作成
	gl.BindTexture(gl.TEXTURE_2D_MULTISAMPLE, msaa.texture)
	gl.TexImage2DMultisample(gl.TEXTURE_2D_MULTISAMPLE, msaa.samples, gl.RGBA8, int32(msaa.width), int32(msaa.height), true)

	// マルチサンプルレンダーバッファを再作成
	gl.BindRenderbuffer(gl.RENDERBUFFER, msaa.rbo)
	gl.RenderbufferStorageMultisample(gl.RENDERBUFFER, msaa.samples, gl.DEPTH24_STENCIL8, int32(msaa.width), int32(msaa.height))
}
