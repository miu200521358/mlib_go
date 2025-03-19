//go:build windows
// +build windows

package mgl

import (
	"fmt"

	"github.com/go-gl/gl/v4.4-core/gl"
	"github.com/miu200521358/mlib_go/pkg/domain/rendering"
)

// msaaImpl は IMsaa インターフェースの実装
type msaaImpl struct {
	config      rendering.MSAAConfig
	fbo         uint32 // マルチサンプル用のフレームバッファオブジェクト
	colorBuffer uint32 // カラー用マルチサンプルレンダーバッファ
	depthBuffer uint32 // 深度/ステンシル用マルチサンプルレンダーバッファ
}

// NewMsaa は指定されたサイズで MSAA を初期化し、IMsaa を返す関数です。
// sampleCount はここでは 4 をデフォルトとしていますが、必要に応じて変更可能です。
func NewMsaa(width, height int) rendering.IMsaa {
	config := rendering.MSAAConfig{
		Width:       width,
		Height:      height,
		SampleCount: 4, // サンプル数（例: 4xMSAA）
	}
	msaa := &msaaImpl{
		config: config,
	}
	msaa.init()
	return msaa
}

// init は MSAA 用の FBO とレンダーバッファを初期化します。
func (m *msaaImpl) init() {
	// FBO の生成
	gl.GenFramebuffers(1, &m.fbo)
	gl.BindFramebuffer(gl.FRAMEBUFFER, m.fbo)

	// カラー用マルチサンプルレンダーバッファの生成と設定
	gl.GenRenderbuffers(1, &m.colorBuffer)
	gl.BindRenderbuffer(gl.RENDERBUFFER, m.colorBuffer)
	gl.RenderbufferStorageMultisample(gl.RENDERBUFFER, int32(m.config.SampleCount), gl.RGBA8, int32(m.config.Width), int32(m.config.Height))
	gl.FramebufferRenderbuffer(gl.FRAMEBUFFER, gl.COLOR_ATTACHMENT0, gl.RENDERBUFFER, m.colorBuffer)

	// 深度／ステンシル用マルチサンプルレンダーバッファの生成と設定
	gl.GenRenderbuffers(1, &m.depthBuffer)
	gl.BindRenderbuffer(gl.RENDERBUFFER, m.depthBuffer)
	gl.RenderbufferStorageMultisample(gl.RENDERBUFFER, int32(m.config.SampleCount), gl.DEPTH24_STENCIL8, int32(m.config.Width), int32(m.config.Height))
	gl.FramebufferRenderbuffer(gl.FRAMEBUFFER, gl.DEPTH_STENCIL_ATTACHMENT, gl.RENDERBUFFER, m.depthBuffer)

	// FBO の状態確認
	if gl.CheckFramebufferStatus(gl.FRAMEBUFFER) != gl.FRAMEBUFFER_COMPLETE {
		fmt.Println("MSAA framebuffer is not complete")
	}

	gl.BindFramebuffer(gl.FRAMEBUFFER, 0)
}

// Bind は MSAA FBO をバインドし、描画先を MSAA に切り替えます。
func (m *msaaImpl) Bind() {
	gl.BindFramebuffer(gl.FRAMEBUFFER, m.fbo)
}

// Unbind は FBO のバインドを解除し、デフォルトのフレームバッファに戻します。
func (m *msaaImpl) Unbind() {
	gl.BindFramebuffer(gl.FRAMEBUFFER, 0)
}

// Resolve はマルチサンプルの結果をデフォルトフレームバッファに解決（ブリット）します。
func (m *msaaImpl) Resolve() {
	// MSAA FBO からデフォルト FBO にカラー情報をコピー
	gl.BindFramebuffer(gl.READ_FRAMEBUFFER, m.fbo)
	gl.BindFramebuffer(gl.DRAW_FRAMEBUFFER, 0)
	gl.BlitFramebuffer(
		0, 0, int32(m.config.Width), int32(m.config.Height),
		0, 0, int32(m.config.Width), int32(m.config.Height),
		gl.COLOR_BUFFER_BIT, gl.NEAREST,
	)
	// 解除
	gl.BindFramebuffer(gl.READ_FRAMEBUFFER, 0)
}

// Delete は MSAA 関連のリソース（FBO とレンダーバッファ）を解放します。
func (m *msaaImpl) Delete() {
	if m.fbo != 0 {
		gl.DeleteFramebuffers(1, &m.fbo)
		m.fbo = 0
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

// Resize は MSAA のレンダーバッファサイズを更新します。
// ウィンドウサイズの変更に合わせ、レンダーバッファを再生成します。
func (m *msaaImpl) Resize(width, height int) {
	m.config.Width = width
	m.config.Height = height

	// 既存のレンダーバッファを削除
	if m.colorBuffer != 0 {
		gl.DeleteRenderbuffers(1, &m.colorBuffer)
	}
	if m.depthBuffer != 0 {
		gl.DeleteRenderbuffers(1, &m.depthBuffer)
	}

	// 再度 FBO にバインド
	gl.BindFramebuffer(gl.FRAMEBUFFER, m.fbo)

	// 新しいカラー用マルチサンプルレンダーバッファの作成
	gl.GenRenderbuffers(1, &m.colorBuffer)
	gl.BindRenderbuffer(gl.RENDERBUFFER, m.colorBuffer)
	gl.RenderbufferStorageMultisample(gl.RENDERBUFFER, int32(m.config.SampleCount), gl.RGBA8, int32(width), int32(height))
	gl.FramebufferRenderbuffer(gl.FRAMEBUFFER, gl.COLOR_ATTACHMENT0, gl.RENDERBUFFER, m.colorBuffer)

	// 新しい深度／ステンシル用レンダーバッファの作成
	gl.GenRenderbuffers(1, &m.depthBuffer)
	gl.BindRenderbuffer(gl.RENDERBUFFER, m.depthBuffer)
	gl.RenderbufferStorageMultisample(gl.RENDERBUFFER, int32(m.config.SampleCount), gl.DEPTH24_STENCIL8, int32(width), int32(height))
	gl.FramebufferRenderbuffer(gl.FRAMEBUFFER, gl.DEPTH_STENCIL_ATTACHMENT, gl.RENDERBUFFER, m.depthBuffer)

	if gl.CheckFramebufferStatus(gl.FRAMEBUFFER) != gl.FRAMEBUFFER_COMPLETE {
		fmt.Println("Resized MSAA framebuffer is not complete")
	}

	gl.BindFramebuffer(gl.FRAMEBUFFER, 0)
}
