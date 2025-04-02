package mgl

import (
	"fmt"

	"github.com/go-gl/gl/v4.4-core/gl"
	"github.com/miu200521358/mlib_go/pkg/domain/rendering"
)

// middleMsaa は、中間FBOを利用して解決するMSAAの実装です。
type middleMsaa struct {
	config rendering.MSAAConfig
	// マルチサンプル用FBOとそのレンダーバッファ
	msaaFBO, msaaColorBuffer, msaaDepthBuffer uint32
	// 中間FBO（シングルサンプル）のFBOとレンダーバッファ
	intermediateFBO, intermediateColorBuffer uint32
}

// NewMiddleMsaa は指定されたサイズで middleMsaa を初期化し、IMsaa を返します。
func NewMiddleMsaa(width, height int) rendering.IMsaa {
	config := rendering.MSAAConfig{
		Width:       width,
		Height:      height,
		SampleCount: 4, // 例として4xMSAA
	}
	m := &middleMsaa{
		config: config,
	}
	m.init()
	return m
}

func (m *middleMsaa) init() {
	// --- マルチサンプルFBOの作成（通常のmsaaと同様） ---
	gl.GenFramebuffers(1, &m.msaaFBO)
	gl.BindFramebuffer(gl.FRAMEBUFFER, m.msaaFBO)

	// カラー用マルチサンプルレンダーバッファの生成と設定
	gl.GenRenderbuffers(1, &m.msaaColorBuffer)
	gl.BindRenderbuffer(gl.RENDERBUFFER, m.msaaColorBuffer)
	gl.RenderbufferStorageMultisample(gl.RENDERBUFFER, int32(m.config.SampleCount), gl.RGBA8, int32(m.config.Width), int32(m.config.Height))
	gl.FramebufferRenderbuffer(gl.FRAMEBUFFER, gl.COLOR_ATTACHMENT0, gl.RENDERBUFFER, m.msaaColorBuffer)

	// 深度／ステンシル用マルチサンプルレンダーバッファの生成と設定
	gl.GenRenderbuffers(1, &m.msaaDepthBuffer)
	gl.BindRenderbuffer(gl.RENDERBUFFER, m.msaaDepthBuffer)
	gl.RenderbufferStorageMultisample(gl.RENDERBUFFER, int32(m.config.SampleCount), gl.DEPTH24_STENCIL8, int32(m.config.Width), int32(m.config.Height))
	gl.FramebufferRenderbuffer(gl.FRAMEBUFFER, gl.DEPTH_STENCIL_ATTACHMENT, gl.RENDERBUFFER, m.msaaDepthBuffer)

	if gl.CheckFramebufferStatus(gl.FRAMEBUFFER) != gl.FRAMEBUFFER_COMPLETE {
		fmt.Println("middleMsaa: multisample FBO is not complete")
	}
	gl.BindFramebuffer(gl.FRAMEBUFFER, 0)

	// --- 中間FBOの作成（シングルサンプル） ---
	gl.GenFramebuffers(1, &m.intermediateFBO)
	gl.BindFramebuffer(gl.FRAMEBUFFER, m.intermediateFBO)

	gl.GenRenderbuffers(1, &m.intermediateColorBuffer)
	gl.BindRenderbuffer(gl.RENDERBUFFER, m.intermediateColorBuffer)
	// 中間FBOはシングルサンプルとして生成
	gl.RenderbufferStorage(gl.RENDERBUFFER, gl.RGBA8, int32(m.config.Width), int32(m.config.Height))
	gl.FramebufferRenderbuffer(gl.FRAMEBUFFER, gl.COLOR_ATTACHMENT0, gl.RENDERBUFFER, m.intermediateColorBuffer)

	if gl.CheckFramebufferStatus(gl.FRAMEBUFFER) != gl.FRAMEBUFFER_COMPLETE {
		fmt.Println("middleMsaa: intermediate FBO is not complete")
	}
	gl.BindFramebuffer(gl.FRAMEBUFFER, 0)
}

func (m *middleMsaa) Bind() {
	// 描画はマルチサンプルFBOに対して行う
	gl.BindFramebuffer(gl.FRAMEBUFFER, m.msaaFBO)
}

func (m *middleMsaa) Unbind() {
	gl.BindFramebuffer(gl.FRAMEBUFFER, 0)
}

func (m *middleMsaa) Resolve() {
	// ① マルチサンプルFBOから中間FBOへブリット
	gl.BindFramebuffer(gl.READ_FRAMEBUFFER, m.msaaFBO)
	gl.BindFramebuffer(gl.DRAW_FRAMEBUFFER, m.intermediateFBO)
	gl.BlitFramebuffer(
		0, 0, int32(m.config.Width), int32(m.config.Height),
		0, 0, int32(m.config.Width), int32(m.config.Height),
		gl.COLOR_BUFFER_BIT, gl.NEAREST,
	)
	// ② 中間FBOからデフォルトFBOへブリット
	gl.BindFramebuffer(gl.READ_FRAMEBUFFER, m.intermediateFBO)
	gl.BindFramebuffer(gl.DRAW_FRAMEBUFFER, 0)
	gl.BlitFramebuffer(
		0, 0, int32(m.config.Width), int32(m.config.Height),
		0, 0, int32(m.config.Width), int32(m.config.Height),
		gl.COLOR_BUFFER_BIT, gl.NEAREST,
	)
	gl.BindFramebuffer(gl.READ_FRAMEBUFFER, 0)
}

func (m *middleMsaa) Delete() {
	if m.msaaFBO != 0 {
		gl.DeleteFramebuffers(1, &m.msaaFBO)
		m.msaaFBO = 0
	}
	if m.msaaColorBuffer != 0 {
		gl.DeleteRenderbuffers(1, &m.msaaColorBuffer)
		m.msaaColorBuffer = 0
	}
	if m.msaaDepthBuffer != 0 {
		gl.DeleteRenderbuffers(1, &m.msaaDepthBuffer)
		m.msaaDepthBuffer = 0
	}
	if m.intermediateFBO != 0 {
		gl.DeleteFramebuffers(1, &m.intermediateFBO)
		m.intermediateFBO = 0
	}
	if m.intermediateColorBuffer != 0 {
		gl.DeleteRenderbuffers(1, &m.intermediateColorBuffer)
		m.intermediateColorBuffer = 0
	}
}

func (m *middleMsaa) Resize(width, height int) {
	m.config.Width = width
	m.config.Height = height

	// マルチサンプルFBOの再生成
	if m.msaaColorBuffer != 0 {
		gl.DeleteRenderbuffers(1, &m.msaaColorBuffer)
	}
	if m.msaaDepthBuffer != 0 {
		gl.DeleteRenderbuffers(1, &m.msaaDepthBuffer)
	}
	gl.BindFramebuffer(gl.FRAMEBUFFER, m.msaaFBO)
	gl.GenRenderbuffers(1, &m.msaaColorBuffer)
	gl.BindRenderbuffer(gl.RENDERBUFFER, m.msaaColorBuffer)
	gl.RenderbufferStorageMultisample(gl.RENDERBUFFER, int32(m.config.SampleCount), gl.RGBA8, int32(width), int32(height))
	gl.FramebufferRenderbuffer(gl.FRAMEBUFFER, gl.COLOR_ATTACHMENT0, gl.RENDERBUFFER, m.msaaColorBuffer)

	gl.GenRenderbuffers(1, &m.msaaDepthBuffer)
	gl.BindRenderbuffer(gl.RENDERBUFFER, m.msaaDepthBuffer)
	gl.RenderbufferStorageMultisample(gl.RENDERBUFFER, int32(m.config.SampleCount), gl.DEPTH24_STENCIL8, int32(width), int32(height))
	gl.FramebufferRenderbuffer(gl.FRAMEBUFFER, gl.DEPTH_STENCIL_ATTACHMENT, gl.RENDERBUFFER, m.msaaDepthBuffer)
	gl.BindFramebuffer(gl.FRAMEBUFFER, 0)

	// 中間FBOの再生成
	if m.intermediateColorBuffer != 0 {
		gl.DeleteRenderbuffers(1, &m.intermediateColorBuffer)
	}
	gl.BindFramebuffer(gl.FRAMEBUFFER, m.intermediateFBO)
	gl.GenRenderbuffers(1, &m.intermediateColorBuffer)
	gl.BindRenderbuffer(gl.RENDERBUFFER, m.intermediateColorBuffer)
	gl.RenderbufferStorage(gl.RENDERBUFFER, gl.RGBA8, int32(width), int32(height))
	gl.FramebufferRenderbuffer(gl.FRAMEBUFFER, gl.COLOR_ATTACHMENT0, gl.RENDERBUFFER, m.intermediateColorBuffer)
	gl.BindFramebuffer(gl.FRAMEBUFFER, 0)
}
