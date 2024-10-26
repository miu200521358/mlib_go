//go:build windows
// +build windows

package mgl

import (
	"github.com/go-gl/gl/v4.4-core/gl"
	"github.com/miu200521358/mlib_go/pkg/mutils/mlog"
)

type Msaa struct {
	width                 int32
	height                int32
	msaa_samples          int32
	msFBO                 uint32
	resolveFBO            uint32
	colorBuffer           uint32
	depthBuffer           uint32
	colorBufferMS         uint32
	depthBufferMS         uint32
	overrideTexture       uint32
	overrideTargetTexture uint32
	overrideVao           *VAO // オーバーライドVAO
	overrideVbo           *VBO // オーバーライドVBO
}

func NewMsaa(width int, height int) *Msaa {
	msaa := &Msaa{
		width:        int32(width),
		height:       int32(height),
		msaa_samples: 4,
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
	gl.RenderbufferStorageMultisample(gl.RENDERBUFFER, msaa.msaa_samples, gl.RGBA8, int32(msaa.width), int32(msaa.height))
	gl.FramebufferRenderbuffer(gl.FRAMEBUFFER, gl.COLOR_ATTACHMENT0, gl.RENDERBUFFER, msaa.colorBufferMS)

	gl.GenRenderbuffers(1, &msaa.depthBufferMS)
	gl.BindRenderbuffer(gl.RENDERBUFFER, msaa.depthBufferMS)
	gl.RenderbufferStorageMultisample(gl.RENDERBUFFER, msaa.msaa_samples, gl.DEPTH_COMPONENT24, int32(msaa.width), int32(msaa.height))
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

	msaa.overrideVao = NewVAO()
	msaa.overrideVao.Bind()
	msaa.overrideVbo = NewVBOForOverride(gl.Ptr(overrideVertices), len(overrideVertices))
	msaa.overrideVbo.Unbind()
	msaa.overrideVao.Unbind()

	// アンバインド
	msaa.Unbind()

	return msaa
}

func (msaa *Msaa) ReadDepthAt(x, y int) float32 {
	// シングルサンプルFBOから読み取る
	gl.BindFramebuffer(gl.FRAMEBUFFER, msaa.resolveFBO)
	if status := gl.CheckFramebufferStatus(gl.FRAMEBUFFER); status != gl.FRAMEBUFFER_COMPLETE {
		mlog.E("Framebuffer is not complete: %v", status)
	}

	var depth float32
	// yは下が0なので、上下反転
	gl.ReadPixels(int32(x), msaa.height-int32(y), 1, 1, gl.DEPTH_COMPONENT, gl.FLOAT, gl.Ptr(&depth))

	// エラーが発生していないかチェック
	if err := gl.GetError(); err != gl.NO_ERROR {
		mlog.E("OpenGL Error after ReadPixels: %v", err)
		return -1
	}

	// // フレームバッファの内容を画像ファイルとして保存
	// err := msaa.saveImage("framebuffer_output.png")
	// if err != nil {
	// 	mlog.E("Failed to save framebuffer image: %v", err)
	// }

	gl.BindFramebuffer(gl.FRAMEBUFFER, 0)

	mlog.V("ReadDepthAt: %v, %v, depth: %.5f", x, y, depth)

	return depth
}

func (msaa *Msaa) Bind() {
	gl.BindFramebuffer(gl.FRAMEBUFFER, msaa.msFBO)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
}

func (msaa *Msaa) Unbind() {
	gl.BindFramebuffer(gl.FRAMEBUFFER, 0)
	gl.BindFramebuffer(gl.READ_FRAMEBUFFER, 0)
	gl.BindFramebuffer(gl.DRAW_FRAMEBUFFER, 0)
	gl.BindTexture(gl.TEXTURE_2D, 0)
}

func (msaa *Msaa) Resolve() {
	// マルチサンプルフレームバッファの内容を解決フレームバッファにコピー
	gl.BindFramebuffer(gl.READ_FRAMEBUFFER, msaa.msFBO)
	gl.BindFramebuffer(gl.DRAW_FRAMEBUFFER, msaa.resolveFBO)
	gl.BlitFramebuffer(0, 0, int32(msaa.width), int32(msaa.height), 0, 0, int32(msaa.width), int32(msaa.height),
		gl.COLOR_BUFFER_BIT|gl.DEPTH_BUFFER_BIT, gl.NEAREST)

	// 解決フレームバッファの内容をウィンドウのデフォルトフレームバッファにコピー
	gl.BindFramebuffer(gl.READ_FRAMEBUFFER, msaa.resolveFBO)
	gl.BindFramebuffer(gl.DRAW_FRAMEBUFFER, 0)
	gl.BlitFramebuffer(0, 0, msaa.width, msaa.height, 0, 0, msaa.width, msaa.height, gl.COLOR_BUFFER_BIT, gl.NEAREST)

	if msaa.overrideTargetTexture != 0 {
		gl.BindFramebuffer(gl.READ_FRAMEBUFFER, msaa.resolveFBO)
		gl.BindTexture(gl.TEXTURE_2D, msaa.overrideTargetTexture)
		gl.CopyTexImage2D(gl.TEXTURE_2D, 0, gl.RGBA, 0, 0, msaa.width, msaa.height, 0)
	}
}

func (msaa *Msaa) Delete() {
	gl.DeleteFramebuffers(1, &msaa.msFBO)
	gl.DeleteFramebuffers(1, &msaa.resolveFBO)
	gl.DeleteRenderbuffers(1, &msaa.colorBuffer)
	gl.DeleteRenderbuffers(1, &msaa.depthBuffer)
	gl.DeleteRenderbuffers(1, &msaa.colorBufferMS)
	gl.DeleteRenderbuffers(1, &msaa.depthBufferMS)
	gl.DeleteTextures(1, &msaa.overrideTexture)
}

func (msaa *Msaa) SetOverrideTargetTexture(texture uint32) {
	msaa.overrideTargetTexture = texture
}

func (msaa *Msaa) OverrideTargetTexture() uint32 {
	return msaa.overrideTargetTexture
}

func (msaa *Msaa) OverrideTextureId() uint32 {
	return msaa.overrideTexture
}

func (msaa *Msaa) BindOverrideTexture(
	windowIndex int, program uint32,
) {
	msaa.overrideVao.Bind()
	msaa.overrideVbo.BindOverride()

	// テクスチャをアクティブにする
	switch windowIndex {
	case 0:
		gl.ActiveTexture(gl.TEXTURE23)
	case 1:
		gl.ActiveTexture(gl.TEXTURE24)
	case 2:
		gl.ActiveTexture(gl.TEXTURE25)
	}

	// テクスチャをバインドする
	gl.BindTexture(gl.TEXTURE_2D, msaa.overrideTexture)

	// テクスチャのパラメーターの設定
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAX_LEVEL, 0)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)

	boneUniform := gl.GetUniformLocation(program, gl.Str("overrideTexture\x00"))
	switch windowIndex {
	case 0:
		gl.Uniform1i(boneUniform, 23)
	case 1:
		gl.Uniform1i(boneUniform, 24)
	case 2:
		gl.Uniform1i(boneUniform, 25)
	}
}

func (msaa *Msaa) UnbindOverrideTexture() {
	gl.BindTexture(gl.TEXTURE_2D, 0)

	msaa.overrideVbo.Unbind()
	msaa.overrideVao.Unbind()
}
