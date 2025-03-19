package mgl

import (
	"fmt"
	"unsafe"

	"github.com/go-gl/gl/v4.4-core/gl"
	"github.com/miu200521358/mlib_go/pkg/domain/rendering"
)

// MOverrideRenderer は、IOverrideRenderer を実装し、
// サブウィンドウの描画内容をテクスチャに書き込み、
// メインウィンドウで半透明合成して描画するためのレンダラーです。
type MOverrideRenderer struct {
	width           int
	height          int
	program         uint32
	isMainWindow    bool
	sharedTextureID *uint32

	fbo     uint32
	texture uint32

	// VertexBufferHandle を活用したフルスクリーンクアッド用バッファ
	quadBuffer *VertexBufferHandle
}

// NewOverrideRenderer は新しい MOverrideRenderer を作成します。
// isMainWindow が true の場合、メインウィンドウ用として動作し、
// そうでない場合はサブウィンドウ用として共有テクスチャを設定します。
func NewOverrideRenderer(width, height int, program uint32, isMainWindow bool) rendering.IOverrideRenderer {
	renderer := &MOverrideRenderer{
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

func (m *MOverrideRenderer) SetSharedTextureID(sharedTextureID *uint32) {
	m.sharedTextureID = sharedTextureID
}

func (m *MOverrideRenderer) SharedTextureIDPtr() *uint32 {
	return m.sharedTextureID
}

func (m *MOverrideRenderer) TextureIDPtr() *uint32 {
	return &m.texture
}

// initFBOAndTexture はレンダリング結果を受け取るための FBO とテクスチャを初期化します。
func (m *MOverrideRenderer) initFBOAndTexture() {
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
	var depthRenderbuffer uint32
	gl.GenRenderbuffers(1, &depthRenderbuffer)
	gl.BindRenderbuffer(gl.RENDERBUFFER, depthRenderbuffer)
	gl.RenderbufferStorage(gl.RENDERBUFFER, gl.DEPTH_COMPONENT24, int32(m.width), int32(m.height))
	gl.FramebufferRenderbuffer(gl.FRAMEBUFFER, gl.DEPTH_ATTACHMENT, gl.RENDERBUFFER, depthRenderbuffer)
	// レンダーバッファのバインド解除
	gl.BindRenderbuffer(gl.RENDERBUFFER, 0)

	// FBO の状態チェック
	status := gl.CheckFramebufferStatus(gl.FRAMEBUFFER)
	if status != gl.FRAMEBUFFER_COMPLETE {
		panic(fmt.Sprintf("Error: Framebuffer is not complete: %s", getFrameBufferStatusString(status)))
	}
	gl.BindFramebuffer(gl.FRAMEBUFFER, 0)
}

// initScreenQuad は、VertexBufferBuilder を利用してフルスクリーンクアッド用の頂点バッファを生成します。
// ここでは、2次元の位置（2要素）とテクスチャ座標（2要素）の頂点データを設定しています。
func (m *MOverrideRenderer) initScreenQuad() {
	quadVertices := []float32{
		// positions   // texCoords
		1.0, 1.0, 0.0, 1.0, 1.0,
		1.0, -1.0, 0.0, 1.0, 0.0,
		-1.0, -1.0, 0.0, 0.0, 0.0,
		-1.0, -1.0, 0.0, 0.0, 0.0,
		-1.0, 1.0, 0.0, 0.0, 1.0,
		1.0, 1.0, 0.0, 1.0, 1.0,
	}
	// VertexBufferBuilder を使用して、2次元位置（2要素）とテクスチャ座標（2要素）の属性を設定
	builder := NewVertexBufferBuilder().
		AddOverrideAttributes()
	builder.SetData(unsafe.Pointer(&quadVertices[0]), len(quadVertices))
	m.quadBuffer = builder.Build()
}

// Bind はレンダリング先をオフスクリーン FBO に設定します。
func (m *MOverrideRenderer) Bind() {
	gl.UseProgram(m.program)
	gl.BindFramebuffer(gl.FRAMEBUFFER, m.fbo)
	gl.Viewport(0, 0, int32(m.width), int32(m.height))
	// カラーと深度の両方をクリアする
	gl.ClearColor(0, 0, 0, 0)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
}

// Unbind は FBO のバインドを解除し、デフォルトのフレームバッファに戻します。
func (m *MOverrideRenderer) Unbind() {
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

// Resolve は、メインウィンドウ側で、subwindow で更新された共有テクスチャを元に、
// 半透明合成（フルスクリーンクアッド描画）を default フレームバッファ上に行い、
// その最終結果をファイルに保存します。
func (m *MOverrideRenderer) Resolve() {
	// 合成描画
	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	// 深度テストを無効化
	gl.Disable(gl.DEPTH_TEST)

	gl.UseProgram(m.program)
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, *m.sharedTextureID)
	location := gl.GetUniformLocation(m.program, gl.Str("overrideTexture\x00"))
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

// Resize は、レンダリング対象のサイズ変更に伴い FBO とテクスチャを再生成します。
func (m *MOverrideRenderer) Resize(width, height int) {
	m.width = width
	m.height = height

	if m.fbo != 0 {
		gl.DeleteFramebuffers(1, &m.fbo)
	}
	if m.texture != 0 {
		gl.DeleteTextures(1, &m.texture)
	}
	m.initFBOAndTexture()

	// サブウィンドウの場合は共有テクスチャも更新
	if !m.isMainWindow {
		*m.sharedTextureID = m.texture
	}
}

// Delete は、FBO、テクスチャ、VertexBufferHandle などのリソースを解放します。
func (m *MOverrideRenderer) Delete() {
	if m.fbo != 0 {
		gl.DeleteFramebuffers(1, &m.fbo)
	}
	if m.texture != 0 {
		gl.DeleteTextures(1, &m.texture)
	}
	if m.quadBuffer != nil {
		m.quadBuffer.Delete()
	}
}
