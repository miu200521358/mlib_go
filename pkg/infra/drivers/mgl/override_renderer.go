//go:build windows
// +build windows

// 指示: miu200521358
package mgl

import (
	"fmt"
	"unsafe"

	"github.com/go-gl/gl/v4.3-core/gl"
	"github.com/miu200521358/mlib_go/pkg/adapter/graphics_api"
	"github.com/miu200521358/mlib_go/pkg/infra/base/i18n"
	"github.com/miu200521358/mlib_go/pkg/shared/base/logging"
)

// OverrideRenderer は、オーバーライド描画用のレンダラーを表す。
type OverrideRenderer struct {
	width           int
	height          int
	program         uint32
	isMainWindow    bool
	sharedTextureID *uint32

	fbo               uint32
	texture           uint32
	depthRenderbuffer uint32

	// VertexBufferHandle を活用したフルスクリーンクアッド用バッファ
	quadBuffer *VertexBufferHandle

	initErr error

	warnedMissingSharedTexture bool
}

// NewOverrideRenderer は新しい OverrideRenderer を作成する。
// isMainWindow が true の場合はメインウィンドウ用として動作する。
func NewOverrideRenderer(width, height int, program uint32, isMainWindow bool) graphics_api.IOverrideRenderer {
	renderer := &OverrideRenderer{
		width:        width,
		height:       height,
		program:      program,
		isMainWindow: isMainWindow,
	}

	// オフスクリーンレンダリング用の FBO とテクスチャを初期化
	renderer.initFBOAndTexture()
	// フルスクリーンクアッドのバッファを VertexBufferHandle を使って生成
	renderer.initScreenQuad()

	return renderer
}

// SetSharedTextureID は共有テクスチャIDを設定する。
func (m *OverrideRenderer) SetSharedTextureID(sharedTextureID *uint32) {
	m.sharedTextureID = sharedTextureID
	m.warnedMissingSharedTexture = false
}

// SharedTextureIDPtr は共有テクスチャIDの参照を返す。
func (m *OverrideRenderer) SharedTextureIDPtr() *uint32 {
	return m.sharedTextureID
}

// TextureIDPtr はこのレンダラーのテクスチャID参照を返す。
func (m *OverrideRenderer) TextureIDPtr() *uint32 {
	return &m.texture
}

// initError は初期化エラーを返す。
func (m *OverrideRenderer) initError() error {
	if m == nil {
		return nil
	}
	return m.initErr
}

// initFBOAndTexture はFBOとテクスチャを初期化する。
func (m *OverrideRenderer) initFBOAndTexture() {
	m.initErr = nil

	// カラー用テクスチャ生成
	gl.GenTextures(1, &m.texture)
	gl.BindTexture(gl.TEXTURE_2D, m.texture)
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA8, int32(m.width), int32(m.height), 0, gl.RGBA, gl.UNSIGNED_BYTE, nil)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.BindTexture(gl.TEXTURE_2D, 0)

	// FBO 生成
	gl.GenFramebuffers(1, &m.fbo)
	gl.BindFramebuffer(gl.FRAMEBUFFER, m.fbo)
	gl.FramebufferTexture2D(gl.FRAMEBUFFER, gl.COLOR_ATTACHMENT0, gl.TEXTURE_2D, m.texture, 0)

	// 深度レンダーバッファの作成と添付
	gl.GenRenderbuffers(1, &m.depthRenderbuffer)
	gl.BindRenderbuffer(gl.RENDERBUFFER, m.depthRenderbuffer)
	gl.RenderbufferStorage(gl.RENDERBUFFER, gl.DEPTH_COMPONENT24, int32(m.width), int32(m.height))
	gl.FramebufferRenderbuffer(gl.FRAMEBUFFER, gl.DEPTH_ATTACHMENT, gl.RENDERBUFFER, m.depthRenderbuffer)
	// レンダーバッファのバインド解除
	gl.BindRenderbuffer(gl.RENDERBUFFER, 0)

	// FBO の状態チェック
	status := gl.CheckFramebufferStatus(gl.FRAMEBUFFER)
	if status != gl.FRAMEBUFFER_COMPLETE {
		m.initErr = graphics_api.NewFramebufferIncomplete(
			fmt.Sprintf(i18n.T("オーバーライド用フレームバッファが不完全です: %s"), getFrameBufferStatusString(status)),
			nil,
		)
		logging.DefaultLogger().Warn(i18n.T("オーバーライドFBOの初期化に失敗しました: %v"), m.initErr)
	}
	gl.BindFramebuffer(gl.FRAMEBUFFER, 0)
}

// initScreenQuad はフルスクリーンクアッド用バッファを生成する。
func (m *OverrideRenderer) initScreenQuad() {
	quadVertices := []float32{
		// positions   // texCoords
		1.0, 1.0, 0.0, 1.0, 1.0,
		1.0, -1.0, 0.0, 1.0, 0.0,
		-1.0, -1.0, 0.0, 0.0, 0.0,
		-1.0, -1.0, 0.0, 0.0, 0.0,
		-1.0, 1.0, 0.0, 0.0, 1.0,
		1.0, 1.0, 0.0, 1.0, 1.0,
	}
	// VertexBufferBuilder を使用して、3次元位置（3要素）とテクスチャ座標（2要素）の属性を設定
	builder := NewVertexBufferBuilder().
		AddOverrideAttributes()
	builder.SetData(unsafe.Pointer(&quadVertices[0]), len(quadVertices))
	m.quadBuffer = builder.Build()
}

// Bind は描画先をFBOに切り替える。
func (m *OverrideRenderer) Bind() {
	gl.UseProgram(m.program)
	gl.BindFramebuffer(gl.FRAMEBUFFER, m.fbo)
	gl.Viewport(0, 0, int32(m.width), int32(m.height))
	// カラーと深度の両方をクリアする
	gl.ClearColor(0, 0, 0, 0)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
}

// Unbind はFBOのバインドを解除する。
func (m *OverrideRenderer) Unbind() {
	// FBOからデフォルトフレームバッファへ内容をブリット（サブウィンドウ表示用）
	if !m.isMainWindow {
		gl.BindFramebuffer(gl.READ_FRAMEBUFFER, m.fbo)
		gl.BindFramebuffer(gl.DRAW_FRAMEBUFFER, 0)
		gl.BlitFramebuffer(
			0, 0, int32(m.width), int32(m.height),
			0, 0, int32(m.width), int32(m.height),
			gl.COLOR_BUFFER_BIT, gl.NEAREST,
		)
		gl.BindFramebuffer(gl.READ_FRAMEBUFFER, 0)
	}

	// FBOのバインドを解除
	gl.BindFramebuffer(gl.FRAMEBUFFER, 0)
	gl.UseProgram(0)
}

// Resolve は共有テクスチャを用いて半透明合成を行う。
func (m *OverrideRenderer) Resolve() {
	if m.sharedTextureID == nil || *m.sharedTextureID == 0 {
		if !m.warnedMissingSharedTexture {
			logging.DefaultLogger().Warn(i18n.T("共有テクスチャが未設定のためオーバーライド合成をスキップします"))
			m.warnedMissingSharedTexture = true
		}
		return
	}
	// 合成描画
	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	// 深度テストを無効化
	gl.Disable(gl.DEPTH_TEST)

	gl.UseProgram(m.program)
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, *m.sharedTextureID)
	location := GetUniformLocation(m.program, "overrideTexture\x00")
	gl.Uniform1i(location, 0)

	m.quadBuffer.Bind()
	gl.DrawArrays(gl.TRIANGLES, 0, 6)
	m.quadBuffer.Unbind()

	// テクスチャのバインド解除
	gl.BindTexture(gl.TEXTURE_2D, 0)
	gl.UseProgram(0)

	// 深度テストを有効化
	gl.Enable(gl.DEPTH_TEST)
	gl.Disable(gl.BLEND)
}

// Resize はレンダリング対象のサイズ変更に伴いFBOとテクスチャを再生成する。
func (m *OverrideRenderer) Resize(width, height int) {
	m.width = width
	m.height = height

	if m.fbo != 0 {
		gl.DeleteFramebuffers(1, &m.fbo)
	}
	if m.texture != 0 {
		gl.DeleteTextures(1, &m.texture)
	}
	if m.depthRenderbuffer != 0 {
		gl.DeleteRenderbuffers(1, &m.depthRenderbuffer)
		m.depthRenderbuffer = 0
	}
	m.initFBOAndTexture()

	// サブウィンドウの場合は共有テクスチャも更新
	if !m.isMainWindow && m.sharedTextureID != nil {
		*m.sharedTextureID = m.texture
	}
}

// Delete はリソースを解放する。
func (m *OverrideRenderer) Delete() {
	if m.fbo != 0 {
		gl.DeleteFramebuffers(1, &m.fbo)
	}
	if m.texture != 0 {
		gl.DeleteTextures(1, &m.texture)
	}
	if m.depthRenderbuffer != 0 {
		gl.DeleteRenderbuffers(1, &m.depthRenderbuffer)
		m.depthRenderbuffer = 0
	}
	if m.quadBuffer != nil {
		m.quadBuffer.Delete()
	}
}
