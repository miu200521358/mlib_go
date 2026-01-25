//go:build windows
// +build windows

// 指示: miu200521358
package mgl

import (
	"github.com/go-gl/gl/v4.3-core/gl"
	"github.com/miu200521358/mlib_go/pkg/adapter/graphics_api"
)

type msaaConfig struct {
	width       int
	height      int
	sampleCount int
}

// MsaaBuffer はMSAA用FBOの実装。
type MsaaBuffer struct {
	config              msaaConfig
	msFBO               uint32
	colorBufferMS       uint32
	depthBufferMS       uint32
	colorBuffer         uint32
	depthBuffer         uint32
	resolveFBO          uint32
	initErr             error
}

// NewMsaaBuffer はMSAAバッファを生成する。
func NewMsaaBuffer(width, height int) graphics_api.IMsaa {
	config := msaaConfig{width: width, height: height, sampleCount: 4}
	buf := &MsaaBuffer{config: config}
	buf.init()
	return buf
}

// init はMSAA用FBOを初期化する。
func (m *MsaaBuffer) init() {
	gl.GenFramebuffers(1, &m.msFBO)
	gl.BindFramebuffer(gl.FRAMEBUFFER, m.msFBO)

	gl.GenRenderbuffers(1, &m.colorBufferMS)
	gl.BindRenderbuffer(gl.RENDERBUFFER, m.colorBufferMS)
	gl.RenderbufferStorageMultisample(
		gl.RENDERBUFFER,
		int32(m.config.sampleCount),
		gl.RGBA8,
		int32(m.config.width),
		int32(m.config.height),
	)
	gl.FramebufferRenderbuffer(gl.FRAMEBUFFER, gl.COLOR_ATTACHMENT0, gl.RENDERBUFFER, m.colorBufferMS)

	gl.GenRenderbuffers(1, &m.depthBufferMS)
	gl.BindRenderbuffer(gl.RENDERBUFFER, m.depthBufferMS)
	gl.RenderbufferStorageMultisample(
		gl.RENDERBUFFER,
		int32(m.config.sampleCount),
		gl.DEPTH_COMPONENT24,
		int32(m.config.width),
		int32(m.config.height),
	)
	gl.FramebufferRenderbuffer(gl.FRAMEBUFFER, gl.DEPTH_ATTACHMENT, gl.RENDERBUFFER, m.depthBufferMS)

	if status := gl.CheckFramebufferStatus(gl.FRAMEBUFFER); status != gl.FRAMEBUFFER_COMPLETE {
		m.initErr = graphics_api.NewFramebufferIncomplete("MSAAのフレームバッファが不完全です", nil)
	}

	gl.GenFramebuffers(1, &m.resolveFBO)
	gl.BindFramebuffer(gl.FRAMEBUFFER, m.resolveFBO)

	gl.GenRenderbuffers(1, &m.colorBuffer)
	gl.BindRenderbuffer(gl.RENDERBUFFER, m.colorBuffer)
	gl.RenderbufferStorage(gl.RENDERBUFFER, gl.RGBA8, int32(m.config.width), int32(m.config.height))
	gl.FramebufferRenderbuffer(gl.FRAMEBUFFER, gl.COLOR_ATTACHMENT0, gl.RENDERBUFFER, m.colorBuffer)

	gl.GenRenderbuffers(1, &m.depthBuffer)
	gl.BindRenderbuffer(gl.RENDERBUFFER, m.depthBuffer)
	gl.RenderbufferStorage(gl.RENDERBUFFER, gl.DEPTH_COMPONENT24, int32(m.config.width), int32(m.config.height))
	gl.FramebufferRenderbuffer(gl.FRAMEBUFFER, gl.DEPTH_ATTACHMENT, gl.RENDERBUFFER, m.depthBuffer)

	if status := gl.CheckFramebufferStatus(gl.FRAMEBUFFER); status != gl.FRAMEBUFFER_COMPLETE {
		m.initErr = graphics_api.NewFramebufferIncomplete("MSAA解決用フレームバッファが不完全です", nil)
	}

	gl.BindFramebuffer(gl.FRAMEBUFFER, 0)
}

// initError は初期化時のエラーを返す。
func (m *MsaaBuffer) initError() error {
	return m.initErr
}

// ReadDepthAt は指定座標の深度値を読み取る。
func (m *MsaaBuffer) ReadDepthAt(x, y, width, height int) float32 {
	_ = width
	_ = height
	gl.BindFramebuffer(gl.FRAMEBUFFER, m.resolveFBO)

	var depth float32
	gl.ReadPixels(int32(x), int32(m.config.height-y), 1, 1, gl.DEPTH_COMPONENT, gl.FLOAT, gl.Ptr(&depth))

	m.Unbind()

	return depth
}

// Bind はMSAAの描画先をバインドする。
func (m *MsaaBuffer) Bind() {
	gl.BindFramebuffer(gl.FRAMEBUFFER, m.msFBO)
}

// Unbind は描画先をデフォルトに戻す。
func (m *MsaaBuffer) Unbind() {
	gl.BindFramebuffer(gl.FRAMEBUFFER, 0)
	gl.BindFramebuffer(gl.READ_FRAMEBUFFER, 0)
	gl.BindFramebuffer(gl.DRAW_FRAMEBUFFER, 0)
	gl.BindTexture(gl.TEXTURE_2D, 0)
}

// Resolve はMSAA結果を解決する。
func (m *MsaaBuffer) Resolve() {
	gl.BindFramebuffer(gl.READ_FRAMEBUFFER, m.msFBO)
	gl.BindFramebuffer(gl.DRAW_FRAMEBUFFER, m.resolveFBO)
	gl.BlitFramebuffer(
		0, 0, int32(m.config.width), int32(m.config.height),
		0, 0, int32(m.config.width), int32(m.config.height),
		gl.COLOR_BUFFER_BIT|gl.DEPTH_BUFFER_BIT,
		gl.NEAREST,
	)
	gl.BindFramebuffer(gl.READ_FRAMEBUFFER, m.resolveFBO)
	gl.BindFramebuffer(gl.DRAW_FRAMEBUFFER, 0)
	gl.BlitFramebuffer(
		0, 0, int32(m.config.width), int32(m.config.height),
		0, 0, int32(m.config.width), int32(m.config.height),
		gl.COLOR_BUFFER_BIT, gl.NEAREST,
	)
	gl.BindFramebuffer(gl.READ_FRAMEBUFFER, 0)
}

// Delete はMSAAリソースを解放する。
func (m *MsaaBuffer) Delete() {
	if m.msFBO != 0 {
		gl.DeleteFramebuffers(1, &m.msFBO)
		m.msFBO = 0
	}
	if m.colorBufferMS != 0 {
		gl.DeleteRenderbuffers(1, &m.colorBufferMS)
		m.colorBufferMS = 0
	}
	if m.depthBufferMS != 0 {
		gl.DeleteRenderbuffers(1, &m.depthBufferMS)
		m.depthBufferMS = 0
	}
	if m.resolveFBO != 0 {
		gl.DeleteFramebuffers(1, &m.resolveFBO)
		m.resolveFBO = 0
	}
	if m.colorBuffer != 0 {
		gl.DeleteRenderbuffers(1, &m.colorBuffer)
		m.colorBuffer = 0
	}
	if m.depthBuffer != 0 {
		gl.DeleteRenderbuffers(1, &m.depthBuffer)
		m.depthBuffer = 0
	}
}

// Resize はMSAAバッファサイズを更新する。
func (m *MsaaBuffer) Resize(width, height int) {
	m.config.width = width
	m.config.height = height
	m.initErr = nil

	if m.colorBufferMS != 0 {
		gl.DeleteRenderbuffers(1, &m.colorBufferMS)
	}
	if m.depthBufferMS != 0 {
		gl.DeleteRenderbuffers(1, &m.depthBufferMS)
	}
	if m.colorBuffer != 0 {
		gl.DeleteRenderbuffers(1, &m.colorBuffer)
	}
	if m.depthBuffer != 0 {
		gl.DeleteRenderbuffers(1, &m.depthBuffer)
	}

	gl.BindFramebuffer(gl.FRAMEBUFFER, m.msFBO)
	gl.GenRenderbuffers(1, &m.colorBufferMS)
	gl.BindRenderbuffer(gl.RENDERBUFFER, m.colorBufferMS)
	gl.RenderbufferStorageMultisample(gl.RENDERBUFFER, int32(m.config.sampleCount), gl.RGBA8, int32(width), int32(height))
	gl.FramebufferRenderbuffer(gl.FRAMEBUFFER, gl.COLOR_ATTACHMENT0, gl.RENDERBUFFER, m.colorBufferMS)

	gl.GenRenderbuffers(1, &m.depthBufferMS)
	gl.BindRenderbuffer(gl.RENDERBUFFER, m.depthBufferMS)
	gl.RenderbufferStorageMultisample(gl.RENDERBUFFER, int32(m.config.sampleCount), gl.DEPTH_COMPONENT24, int32(width), int32(height))
	gl.FramebufferRenderbuffer(gl.FRAMEBUFFER, gl.DEPTH_ATTACHMENT, gl.RENDERBUFFER, m.depthBufferMS)

	if status := gl.CheckFramebufferStatus(gl.FRAMEBUFFER); status != gl.FRAMEBUFFER_COMPLETE {
		m.initErr = graphics_api.NewFramebufferIncomplete(
			"MSAAのフレームバッファが不完全です: "+getFrameBufferStatusString(status),
			nil,
		)
	}

	gl.BindFramebuffer(gl.FRAMEBUFFER, m.resolveFBO)
	gl.GenRenderbuffers(1, &m.colorBuffer)
	gl.BindRenderbuffer(gl.RENDERBUFFER, m.colorBuffer)
	gl.RenderbufferStorage(gl.RENDERBUFFER, gl.RGBA8, int32(width), int32(height))
	gl.FramebufferRenderbuffer(gl.FRAMEBUFFER, gl.COLOR_ATTACHMENT0, gl.RENDERBUFFER, m.colorBuffer)

	gl.GenRenderbuffers(1, &m.depthBuffer)
	gl.BindRenderbuffer(gl.RENDERBUFFER, m.depthBuffer)
	gl.RenderbufferStorage(gl.RENDERBUFFER, gl.DEPTH_COMPONENT24, int32(width), int32(height))
	gl.FramebufferRenderbuffer(gl.FRAMEBUFFER, gl.DEPTH_ATTACHMENT, gl.RENDERBUFFER, m.depthBuffer)

	if status := gl.CheckFramebufferStatus(gl.FRAMEBUFFER); status != gl.FRAMEBUFFER_COMPLETE {
		m.initErr = graphics_api.NewFramebufferIncomplete(
			"MSAA解決用フレームバッファが不完全です: "+getFrameBufferStatusString(status),
			nil,
		)
	}

	gl.BindFramebuffer(gl.FRAMEBUFFER, 0)
}
