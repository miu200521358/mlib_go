//go:build windows
// +build windows

package mgl

import (
	"github.com/go-gl/gl/v4.4-core/gl"
	"github.com/miu200521358/mlib_go/pkg/config/mlog"
	"github.com/miu200521358/mlib_go/pkg/domain/rendering"
)

// MsaaFactory はMSAAオブジェクトのファクトリー
type MsaaFactory struct{}

// NewMsaaFactory は新しいMSAAファクトリーを作成する
func NewMsaaFactory() *MsaaFactory {
	return &MsaaFactory{}
}

// CreateMsaa は新しいMSAAオブジェクトを作成する
func (f *MsaaFactory) CreateMsaa(width, height int) rendering.IMsaa {
	config := rendering.MSAAConfig{
		Width:       width,
		Height:      height,
		SampleCount: 4, // デフォルトのサンプル数
	}
	return NewMsaa(config)
}

// CreateMsaaWithConfig は指定された設定でMSAAオブジェクトを作成する
func (f *MsaaFactory) CreateMsaaWithConfig(config rendering.MSAAConfig) rendering.IMsaa {
	return NewMsaa(config)
}

// --------------------------------------------------

// Msaa はOpenGLを使用したMSAA実装
type Msaa struct {
	width                 int32
	height                int32
	sampleCount           int32
	msFBO                 uint32
	resolveFBO            uint32
	colorBuffer           uint32
	depthBuffer           uint32
	colorBufferMS         uint32
	depthBufferMS         uint32
	overrideTexture       uint32
	overrideTargetTexture uint32
	overrideBuffer        *VertexBufferHandle
}

// NewMsaa は新しいMSAAオブジェクトを作成する
func NewMsaa(config rendering.MSAAConfig) *Msaa {
	msaa := &Msaa{
		width:       int32(config.Width),
		height:      int32(config.Height),
		sampleCount: int32(config.SampleCount),
	}

	// 深度テストの有効化
	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LEQUAL)

	// マルチサンプルフレームバッファの作成
	gl.GenFramebuffers(1, &msaa.msFBO)
	gl.BindFramebuffer(gl.FRAMEBUFFER, msaa.msFBO)

	// マルチサンプルカラーおよび深度バッファの作成
	gl.GenRenderbuffers(1, &msaa.colorBufferMS)
	gl.BindRenderbuffer(gl.RENDERBUFFER, msaa.colorBufferMS)
	gl.RenderbufferStorageMultisample(gl.RENDERBUFFER, msaa.sampleCount, gl.RGBA8, msaa.width, msaa.height)
	gl.FramebufferRenderbuffer(gl.FRAMEBUFFER, gl.COLOR_ATTACHMENT0, gl.RENDERBUFFER, msaa.colorBufferMS)

	gl.GenRenderbuffers(1, &msaa.depthBufferMS)
	gl.BindRenderbuffer(gl.RENDERBUFFER, msaa.depthBufferMS)
	gl.RenderbufferStorageMultisample(gl.RENDERBUFFER, msaa.sampleCount, gl.DEPTH_COMPONENT24, msaa.width, msaa.height)
	gl.FramebufferRenderbuffer(gl.FRAMEBUFFER, gl.DEPTH_ATTACHMENT, gl.RENDERBUFFER, msaa.depthBufferMS)

	if gl.CheckFramebufferStatus(gl.FRAMEBUFFER) != gl.FRAMEBUFFER_COMPLETE {
		mlog.F("MSAA生成失敗")
	}

	// 解決フレームバッファの作成
	gl.GenFramebuffers(1, &msaa.resolveFBO)
	gl.BindFramebuffer(gl.FRAMEBUFFER, msaa.resolveFBO)

	// 解決フレームバッファのカラーおよび深度バッファの作成
	gl.GenRenderbuffers(1, &msaa.colorBuffer)
	gl.BindRenderbuffer(gl.RENDERBUFFER, msaa.colorBuffer)
	gl.RenderbufferStorage(gl.RENDERBUFFER, gl.RGBA8, msaa.width, msaa.height)
	gl.FramebufferRenderbuffer(gl.FRAMEBUFFER, gl.COLOR_ATTACHMENT0, gl.RENDERBUFFER, msaa.colorBuffer)

	gl.GenRenderbuffers(1, &msaa.depthBuffer)
	gl.BindRenderbuffer(gl.RENDERBUFFER, msaa.depthBuffer)
	gl.RenderbufferStorage(gl.RENDERBUFFER, gl.DEPTH_COMPONENT24, msaa.width, msaa.height)
	gl.FramebufferRenderbuffer(gl.FRAMEBUFFER, gl.DEPTH_ATTACHMENT, gl.RENDERBUFFER, msaa.depthBuffer)

	gl.GenTextures(1, &msaa.overrideTexture)
	gl.BindTexture(gl.TEXTURE_2D, msaa.overrideTexture)
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA, msaa.width, msaa.height, 0, gl.RGBA, gl.UNSIGNED_BYTE, nil)

	if gl.CheckFramebufferStatus(gl.FRAMEBUFFER) != gl.FRAMEBUFFER_COMPLETE {
		mlog.F("解決フレームバッファ生成失敗")
	}

	// 透過描画用に四角形の頂点とテクスチャ座標を定義
	overrideVertices := []float32{
		// positions   // texCoords
		1.0, 1.0, 0.0, 1.0, 1.0,
		1.0, -1.0, 0.0, 1.0, 0.0,
		-1.0, -1.0, 0.0, 0.0, 0.0,
		-1.0, -1.0, 0.0, 0.0, 0.0,
		-1.0, 1.0, 0.0, 0.0, 1.0,
		1.0, 1.0, 0.0, 1.0, 1.0,
	}

	// 新しいバッファ抽象化を使用する
	factory := NewBufferFactory()
	msaa.overrideBuffer = factory.CreateOverrideBuffer(gl.Ptr(overrideVertices), len(overrideVertices), 5)

	// アンバインド
	msaa.Unbind()

	return msaa
}

// ReadDepthAt は指定座標の深度値を読み取る
func (m *Msaa) ReadDepthAt(x, y int) float32 {
	// シングルサンプルFBOから読み取る
	gl.BindFramebuffer(gl.FRAMEBUFFER, m.resolveFBO)
	if status := gl.CheckFramebufferStatus(gl.FRAMEBUFFER); status != gl.FRAMEBUFFER_COMPLETE {
		mlog.E("Framebuffer is not complete: %v", status)
	}

	var depth float32
	// yは下が0なので、上下反転
	gl.ReadPixels(int32(x), m.height-int32(y), 1, 1, gl.DEPTH_COMPONENT, gl.FLOAT, gl.Ptr(&depth))

	// エラーが発生していないかチェック
	if err := gl.GetError(); err != gl.NO_ERROR {
		mlog.E("OpenGL Error after ReadPixels: %v", err)
		return -1
	}

	gl.BindFramebuffer(gl.FRAMEBUFFER, 0)

	mlog.V("ReadDepthAt: %v, %v, depth: %.5f", x, y, depth)

	return depth
}

// Bind はMSAAフレームバッファをバインドする
func (m *Msaa) Bind() {
	gl.BindFramebuffer(gl.FRAMEBUFFER, m.msFBO)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
}

// Unbind はMSAAフレームバッファをアンバインドする
func (m *Msaa) Unbind() {
	gl.BindFramebuffer(gl.FRAMEBUFFER, 0)
	gl.BindFramebuffer(gl.READ_FRAMEBUFFER, 0)
	gl.BindFramebuffer(gl.DRAW_FRAMEBUFFER, 0)
	gl.BindTexture(gl.TEXTURE_2D, 0)
}

// Resolve はMSAAの結果をデフォルトフレームバッファに解決する
func (m *Msaa) Resolve() {
	// マルチサンプルフレームバッファの内容を解決フレームバッファにコピー
	gl.BindFramebuffer(gl.READ_FRAMEBUFFER, m.msFBO)
	gl.BindFramebuffer(gl.DRAW_FRAMEBUFFER, m.resolveFBO)
	gl.BlitFramebuffer(0, 0, m.width, m.height, 0, 0, m.width, m.height,
		gl.COLOR_BUFFER_BIT|gl.DEPTH_BUFFER_BIT, gl.NEAREST)

	// 解決フレームバッファの内容をウィンドウのデフォルトフレームバッファにコピー
	gl.BindFramebuffer(gl.READ_FRAMEBUFFER, m.resolveFBO)
	gl.BindFramebuffer(gl.DRAW_FRAMEBUFFER, 0)
	gl.BlitFramebuffer(0, 0, m.width, m.height, 0, 0, m.width, m.height, gl.COLOR_BUFFER_BIT, gl.NEAREST)

	if m.overrideTargetTexture != 0 {
		gl.BindFramebuffer(gl.READ_FRAMEBUFFER, m.resolveFBO)
		gl.BindTexture(gl.TEXTURE_2D, m.overrideTargetTexture)
		gl.CopyTexImage2D(gl.TEXTURE_2D, 0, gl.RGBA, 0, 0, m.width, m.height, 0)
	}
}

// Delete はMSAAリソースを解放する
func (m *Msaa) Delete() {
	gl.DeleteFramebuffers(1, &m.msFBO)
	gl.DeleteFramebuffers(1, &m.resolveFBO)
	gl.DeleteRenderbuffers(1, &m.colorBuffer)
	gl.DeleteRenderbuffers(1, &m.depthBuffer)
	gl.DeleteRenderbuffers(1, &m.colorBufferMS)
	gl.DeleteRenderbuffers(1, &m.depthBufferMS)
	gl.DeleteTextures(1, &m.overrideTexture)

	if m.overrideBuffer != nil {
		m.overrideBuffer.Delete()
	}
}

// SetOverrideTargetTexture はオーバーライドターゲットテクスチャを設定する
func (m *Msaa) SetOverrideTargetTexture(texture uint32) {
	m.overrideTargetTexture = texture
}

// OverrideTargetTexture はオーバーライドターゲットテクスチャのIDを取得する
func (m *Msaa) OverrideTargetTexture() uint32 {
	return m.overrideTargetTexture
}

// BindOverrideTexture はオーバーライドテクスチャをバインドする
func (m *Msaa) BindOverrideTexture(windowIndex int, program uint32) {
	m.overrideBuffer.Bind()

	// テクスチャユニットの選択
	var textureUnit uint32
	switch windowIndex {
	case 0:
		textureUnit = gl.TEXTURE23
	case 1:
		textureUnit = gl.TEXTURE24
	case 2:
		textureUnit = gl.TEXTURE25
	default:
		textureUnit = gl.TEXTURE23
	}

	// テクスチャをアクティブにする
	gl.ActiveTexture(textureUnit)

	// テクスチャをバインドする
	gl.BindTexture(gl.TEXTURE_2D, m.overrideTexture)

	// テクスチャのパラメーターの設定
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAX_LEVEL, 0)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)

	// シェーダーにテクスチャユニットを設定
	uniform := gl.GetUniformLocation(program, gl.Str("overrideTexture\x00"))
	gl.Uniform1i(uniform, int32(windowIndex+23))
}

// UnbindOverrideTexture はオーバーライドテクスチャをアンバインドする
func (m *Msaa) UnbindOverrideTexture() {
	gl.BindTexture(gl.TEXTURE_2D, 0)
	m.overrideBuffer.Unbind()
}

// Resize はMSAAバッファのサイズを変更する
func (m *Msaa) Resize(width, height int) {
	// 既存のリソースを削除
	m.Delete()

	// 新しいサイズで再作成
	newConfig := rendering.MSAAConfig{
		Width:       width,
		Height:      height,
		SampleCount: int(m.sampleCount),
	}

	newMsaa := NewMsaa(newConfig)

	// フィールドをコピー
	*m = *newMsaa
}
