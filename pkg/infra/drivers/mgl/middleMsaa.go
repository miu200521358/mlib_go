//go:build windows
// +build windows

// 指示: miu200521358
package mgl

import (
	"github.com/go-gl/gl/v4.3-core/gl"
	"github.com/miu200521358/mlib_go/pkg/adapter/graphics_api"
	"github.com/miu200521358/mlib_go/pkg/shared/base/logging"
)

// IntermediateMsaaBuffer は中間FBOを使ったMSAA実装。
type IntermediateMsaaBuffer struct {
	config                   msaaConfig
	msaaFBO                  uint32
	msaaColorBuffer          uint32
	msaaDepthBuffer          uint32
	intermediateFBO          uint32
	intermediateColorBuffer  uint32
	intermediateDepthTexture uint32
	initErr                  error
}

// NewIntermediateMsaaBuffer は中間FBOを使用するMSAAバッファを生成する。
func NewIntermediateMsaaBuffer(width, height int) graphics_api.IMsaa {
	config := msaaConfig{width: width, height: height, sampleCount: 4}
	buf := &IntermediateMsaaBuffer{config: config}
	buf.init()
	return buf
}

// init は中間FBOを初期化する。
func (m *IntermediateMsaaBuffer) init() {
	gl.GenFramebuffers(1, &m.msaaFBO)
	gl.BindFramebuffer(gl.FRAMEBUFFER, m.msaaFBO)

	gl.GenRenderbuffers(1, &m.msaaColorBuffer)
	gl.BindRenderbuffer(gl.RENDERBUFFER, m.msaaColorBuffer)
	gl.RenderbufferStorageMultisample(gl.RENDERBUFFER, int32(m.config.sampleCount), gl.RGBA8, int32(m.config.width), int32(m.config.height))
	gl.FramebufferRenderbuffer(gl.FRAMEBUFFER, gl.COLOR_ATTACHMENT0, gl.RENDERBUFFER, m.msaaColorBuffer)

	gl.GenRenderbuffers(1, &m.msaaDepthBuffer)
	gl.BindRenderbuffer(gl.RENDERBUFFER, m.msaaDepthBuffer)
	gl.RenderbufferStorageMultisample(gl.RENDERBUFFER, int32(m.config.sampleCount), gl.DEPTH24_STENCIL8, int32(m.config.width), int32(m.config.height))
	gl.FramebufferRenderbuffer(gl.FRAMEBUFFER, gl.DEPTH_STENCIL_ATTACHMENT, gl.RENDERBUFFER, m.msaaDepthBuffer)

	if gl.CheckFramebufferStatus(gl.FRAMEBUFFER) != gl.FRAMEBUFFER_COMPLETE {
		m.initErr = graphics_api.NewFramebufferIncomplete("中間MSAAのマルチサンプルFBOが不完全です", nil)
		logging.DefaultLogger().Warn("中間MSAAのマルチサンプルFBOが不完全です")
	}
	gl.BindFramebuffer(gl.FRAMEBUFFER, 0)

	gl.GenFramebuffers(1, &m.intermediateFBO)
	gl.BindFramebuffer(gl.FRAMEBUFFER, m.intermediateFBO)

	gl.GenRenderbuffers(1, &m.intermediateColorBuffer)
	gl.BindRenderbuffer(gl.RENDERBUFFER, m.intermediateColorBuffer)
	gl.RenderbufferStorage(gl.RENDERBUFFER, gl.RGBA8, int32(m.config.width), int32(m.config.height))
	gl.FramebufferRenderbuffer(gl.FRAMEBUFFER, gl.COLOR_ATTACHMENT0, gl.RENDERBUFFER, m.intermediateColorBuffer)

	gl.GenTextures(1, &m.intermediateDepthTexture)
	gl.BindTexture(gl.TEXTURE_2D, m.intermediateDepthTexture)
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.DEPTH_COMPONENT32F, int32(m.config.width), int32(m.config.height), 0, gl.DEPTH_COMPONENT, gl.FLOAT, nil)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
	gl.FramebufferTexture2D(gl.FRAMEBUFFER, gl.DEPTH_ATTACHMENT, gl.TEXTURE_2D, m.intermediateDepthTexture, 0)

	if gl.CheckFramebufferStatus(gl.FRAMEBUFFER) != gl.FRAMEBUFFER_COMPLETE {
		m.initErr = graphics_api.NewFramebufferIncomplete("中間MSAAの中間FBOが不完全です", nil)
		logging.DefaultLogger().Warn("中間MSAAの中間FBOが不完全です")
	}
	gl.BindFramebuffer(gl.FRAMEBUFFER, 0)
}

// initError は初期化時のエラーを返す。
func (m *IntermediateMsaaBuffer) initError() error {
	return m.initErr
}

// ReadDepthAt は指定座標の深度値を読み取る。
func (m *IntermediateMsaaBuffer) ReadDepthAt(x, y, width, height int) float32 {
	var readFBO int32
	var drawFBO int32
	gl.GetIntegerv(gl.READ_FRAMEBUFFER_BINDING, &readFBO)
	gl.GetIntegerv(gl.DRAW_FRAMEBUFFER_BINDING, &drawFBO)

	gl.BindFramebuffer(gl.READ_FRAMEBUFFER, m.msaaFBO)
	gl.BindFramebuffer(gl.DRAW_FRAMEBUFFER, m.intermediateFBO)
	gl.BlitFramebuffer(
		0, 0, int32(m.config.width), int32(m.config.height),
		0, 0, int32(m.config.width), int32(m.config.height),
		gl.DEPTH_BUFFER_BIT, gl.NEAREST,
	)

	gl.BindFramebuffer(gl.FRAMEBUFFER, m.intermediateFBO)
	if status := gl.CheckFramebufferStatus(gl.FRAMEBUFFER); status != gl.FRAMEBUFFER_COMPLETE {
		logging.DefaultLogger().Warn("Framebuffer is not complete: %v", status)
	}

	var depth float32
	gl.ReadPixels(int32(x), int32(height-y), 1, 1, gl.DEPTH_COMPONENT, gl.FLOAT, gl.Ptr(&depth))

	gl.BindFramebuffer(gl.READ_FRAMEBUFFER, uint32(readFBO))
	gl.BindFramebuffer(gl.DRAW_FRAMEBUFFER, uint32(drawFBO))

	return depth
}

// ReadDepthRegion は指定矩形の深度値を読み取る。
func (m *IntermediateMsaaBuffer) ReadDepthRegion(x, y, width, height, framebufferHeight int) []float32 {
	if width <= 0 || height <= 0 || framebufferHeight <= 0 {
		return nil
	}
	var readFBO int32
	var drawFBO int32
	gl.GetIntegerv(gl.READ_FRAMEBUFFER_BINDING, &readFBO)
	gl.GetIntegerv(gl.DRAW_FRAMEBUFFER_BINDING, &drawFBO)

	gl.BindFramebuffer(gl.READ_FRAMEBUFFER, m.msaaFBO)
	gl.BindFramebuffer(gl.DRAW_FRAMEBUFFER, m.intermediateFBO)
	gl.BlitFramebuffer(
		0, 0, int32(m.config.width), int32(m.config.height),
		0, 0, int32(m.config.width), int32(m.config.height),
		gl.DEPTH_BUFFER_BIT, gl.NEAREST,
	)

	gl.BindFramebuffer(gl.FRAMEBUFFER, m.intermediateFBO)
	var depth []float32
	readY := framebufferHeight - (y + height)
	if readY < 0 {
		readY = 0
	}
	if width > 0 && height > 0 {
		depth = make([]float32, width*height)
		gl.ReadPixels(int32(x), int32(readY), int32(width), int32(height), gl.DEPTH_COMPONENT, gl.FLOAT, gl.Ptr(&depth[0]))
	}

	gl.BindFramebuffer(gl.READ_FRAMEBUFFER, uint32(readFBO))
	gl.BindFramebuffer(gl.DRAW_FRAMEBUFFER, uint32(drawFBO))

	return depth
}

// Bind は描画先をMSAAに切り替える。
func (m *IntermediateMsaaBuffer) Bind() {
	gl.BindFramebuffer(gl.FRAMEBUFFER, m.msaaFBO)
}

// Unbind は描画先をデフォルトに戻す。
func (m *IntermediateMsaaBuffer) Unbind() {
	gl.BindFramebuffer(gl.FRAMEBUFFER, 0)
}

// Resolve はMSAA結果を解決する。
func (m *IntermediateMsaaBuffer) Resolve() {
	gl.BindFramebuffer(gl.READ_FRAMEBUFFER, m.msaaFBO)
	gl.BindFramebuffer(gl.DRAW_FRAMEBUFFER, m.intermediateFBO)
	gl.BlitFramebuffer(
		0, 0, int32(m.config.width), int32(m.config.height),
		0, 0, int32(m.config.width), int32(m.config.height),
		gl.COLOR_BUFFER_BIT, gl.NEAREST,
	)
	gl.BindFramebuffer(gl.READ_FRAMEBUFFER, m.intermediateFBO)
	gl.BindFramebuffer(gl.DRAW_FRAMEBUFFER, 0)
	gl.BlitFramebuffer(
		0, 0, int32(m.config.width), int32(m.config.height),
		0, 0, int32(m.config.width), int32(m.config.height),
		gl.COLOR_BUFFER_BIT, gl.NEAREST,
	)
	gl.BindFramebuffer(gl.READ_FRAMEBUFFER, 0)
}

// ResolveDepth は深度のみを解決して読み出し用FBOへ反映する。
func (m *IntermediateMsaaBuffer) ResolveDepth() {
	var readFBO int32
	var drawFBO int32
	gl.GetIntegerv(gl.READ_FRAMEBUFFER_BINDING, &readFBO)
	gl.GetIntegerv(gl.DRAW_FRAMEBUFFER_BINDING, &drawFBO)

	gl.BindFramebuffer(gl.READ_FRAMEBUFFER, m.msaaFBO)
	gl.BindFramebuffer(gl.DRAW_FRAMEBUFFER, m.intermediateFBO)
	gl.BlitFramebuffer(
		0, 0, int32(m.config.width), int32(m.config.height),
		0, 0, int32(m.config.width), int32(m.config.height),
		gl.DEPTH_BUFFER_BIT, gl.NEAREST,
	)
	gl.BindFramebuffer(gl.READ_FRAMEBUFFER, uint32(readFBO))
	gl.BindFramebuffer(gl.DRAW_FRAMEBUFFER, uint32(drawFBO))
}

// Delete はMSAAリソースを解放する。
func (m *IntermediateMsaaBuffer) Delete() {
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
	if m.intermediateDepthTexture != 0 {
		gl.DeleteTextures(1, &m.intermediateDepthTexture)
		m.intermediateDepthTexture = 0
	}
}

// Resize はMSAAバッファサイズを更新する。
func (m *IntermediateMsaaBuffer) Resize(width, height int) {
	m.config.width = width
	m.config.height = height

	if m.msaaColorBuffer != 0 {
		gl.DeleteRenderbuffers(1, &m.msaaColorBuffer)
	}
	if m.msaaDepthBuffer != 0 {
		gl.DeleteRenderbuffers(1, &m.msaaDepthBuffer)
	}
	gl.BindFramebuffer(gl.FRAMEBUFFER, m.msaaFBO)
	gl.GenRenderbuffers(1, &m.msaaColorBuffer)
	gl.BindRenderbuffer(gl.RENDERBUFFER, m.msaaColorBuffer)
	gl.RenderbufferStorageMultisample(gl.RENDERBUFFER, int32(m.config.sampleCount), gl.RGBA8, int32(width), int32(height))
	gl.FramebufferRenderbuffer(gl.FRAMEBUFFER, gl.COLOR_ATTACHMENT0, gl.RENDERBUFFER, m.msaaColorBuffer)

	gl.GenRenderbuffers(1, &m.msaaDepthBuffer)
	gl.BindRenderbuffer(gl.RENDERBUFFER, m.msaaDepthBuffer)
	gl.RenderbufferStorageMultisample(gl.RENDERBUFFER, int32(m.config.sampleCount), gl.DEPTH24_STENCIL8, int32(width), int32(height))
	gl.FramebufferRenderbuffer(gl.FRAMEBUFFER, gl.DEPTH_STENCIL_ATTACHMENT, gl.RENDERBUFFER, m.msaaDepthBuffer)
	gl.BindFramebuffer(gl.FRAMEBUFFER, 0)

	if m.intermediateColorBuffer != 0 {
		gl.DeleteRenderbuffers(1, &m.intermediateColorBuffer)
	}
	if m.intermediateDepthTexture != 0 {
		gl.DeleteTextures(1, &m.intermediateDepthTexture)
	}
	gl.BindFramebuffer(gl.FRAMEBUFFER, m.intermediateFBO)
	gl.GenRenderbuffers(1, &m.intermediateColorBuffer)
	gl.BindRenderbuffer(gl.RENDERBUFFER, m.intermediateColorBuffer)
	gl.RenderbufferStorage(gl.RENDERBUFFER, gl.RGBA8, int32(width), int32(height))
	gl.FramebufferRenderbuffer(gl.FRAMEBUFFER, gl.COLOR_ATTACHMENT0, gl.RENDERBUFFER, m.intermediateColorBuffer)

	gl.GenTextures(1, &m.intermediateDepthTexture)
	gl.BindTexture(gl.TEXTURE_2D, m.intermediateDepthTexture)
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.DEPTH_COMPONENT32F, int32(width), int32(height), 0, gl.DEPTH_COMPONENT, gl.FLOAT, nil)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
	gl.FramebufferTexture2D(gl.FRAMEBUFFER, gl.DEPTH_ATTACHMENT, gl.TEXTURE_2D, m.intermediateDepthTexture, 0)

	gl.BindFramebuffer(gl.FRAMEBUFFER, 0)
}
