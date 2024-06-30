//go:build windows
// +build windows

package mview

import (
	"image"
	"image/color"
	"image/png"
	"os"

	"github.com/go-gl/gl/v4.4-core/gl"
	"github.com/miu200521358/mlib_go/pkg/mutils/mlog"
)

type Msaa struct {
	width         int32
	height        int32
	msaa_samples  int32
	msFBO         uint32
	resolveFBO    uint32
	colorBuffer   uint32
	depthBuffer   uint32
	colorBufferMS uint32
	depthBufferMS uint32
}

func NewMsaa(width int32, height int32) *Msaa {
	msaa := &Msaa{
		width:        width,
		height:       height,
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
		mlog.F("Multisample Framebuffer is not complete")
	}

	// アンバインド
	gl.BindFramebuffer(gl.FRAMEBUFFER, 0)

	// シングルサンプルフレームバッファの作成
	gl.GenFramebuffers(1, &msaa.resolveFBO)
	gl.BindFramebuffer(gl.FRAMEBUFFER, msaa.resolveFBO)

	// シングルサンプルカラーおよび深度バッファの作成
	gl.GenTextures(1, &msaa.colorBuffer)
	gl.BindTexture(gl.TEXTURE_2D, msaa.colorBuffer)
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA8, int32(msaa.width), int32(msaa.height), 0, gl.RGBA, gl.UNSIGNED_BYTE, nil)
	gl.FramebufferTexture2D(gl.FRAMEBUFFER, gl.COLOR_ATTACHMENT0, gl.TEXTURE_2D, msaa.colorBuffer, 0)

	gl.GenTextures(1, &msaa.depthBuffer)
	gl.BindTexture(gl.TEXTURE_2D, msaa.depthBuffer)
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.DEPTH_COMPONENT24, int32(msaa.width), int32(msaa.height), 0, gl.DEPTH_COMPONENT, gl.UNSIGNED_BYTE, nil)
	gl.FramebufferTexture2D(gl.FRAMEBUFFER, gl.DEPTH_ATTACHMENT, gl.TEXTURE_2D, msaa.depthBuffer, 0)

	if gl.CheckFramebufferStatus(gl.FRAMEBUFFER) != gl.FRAMEBUFFER_COMPLETE {
		mlog.F("Single-sample Framebuffer is not complete")
	}

	// アンバインド
	gl.BindFramebuffer(gl.FRAMEBUFFER, 0)

	return msaa
}

func (m *Msaa) ReadDepthAt(x, y, width, height int) float32 {
	// マウスクリック時に解像度をダウンしてFBOを解決
	gl.BindFramebuffer(gl.READ_FRAMEBUFFER, m.msFBO)
	gl.BindFramebuffer(gl.DRAW_FRAMEBUFFER, m.resolveFBO)
	gl.BlitFramebuffer(0, 0, int32(m.width), int32(m.height), 0, 0, int32(m.width), int32(m.height),
		gl.COLOR_BUFFER_BIT|gl.DEPTH_BUFFER_BIT, gl.NEAREST)

	// エラーが発生していないかチェック
	if err := gl.GetError(); err != gl.NO_ERROR {
		mlog.E("OpenGL Error after blit: %v", err)
	}

	// シングルサンプルFBOから読み取る
	gl.BindFramebuffer(gl.FRAMEBUFFER, m.resolveFBO)
	if status := gl.CheckFramebufferStatus(gl.FRAMEBUFFER); status != gl.FRAMEBUFFER_COMPLETE {
		mlog.E("Framebuffer is not complete: %v", status)
	}

	var depth float32
	// yは下が0なので、上下反転
	gl.ReadPixels(int32(x), int32(height-y), 1, 1, gl.DEPTH_COMPONENT, gl.FLOAT, gl.Ptr(&depth))

	// エラーが発生していないかチェック
	if err := gl.GetError(); err != gl.NO_ERROR {
		mlog.E("OpenGL Error after ReadPixels: %v", err)
	}

	mlog.D("Depth at (%d, %d): %f", x, y, depth)

	// フレームバッファの内容を画像ファイルとして保存
	pixels := readFramebuffer(int32(m.width), int32(m.height))
	err := saveImage("framebuffer_output.png", int32(m.width), int32(m.height), pixels)
	if err != nil {
		mlog.E("Failed to save framebuffer image: %v", err)
	}

	// フレームバッファをアンバインド
	gl.BindFramebuffer(gl.FRAMEBUFFER, 0)
	gl.BindFramebuffer(gl.READ_FRAMEBUFFER, 0)
	gl.BindFramebuffer(gl.DRAW_FRAMEBUFFER, 0)

	return depth
}

func readFramebuffer(w, h int32) []byte {
	pixels := make([]byte, w*h*4) // RGBA形式で4バイト/ピクセル
	gl.ReadPixels(0, 0, w, h, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(pixels))
	return pixels
}

func saveImage(filename string, w, h int32, pixels []byte) error {
	img := image.NewRGBA(image.Rect(0, 0, int(w), int(h)))
	for y := int32(0); y < h; y++ {
		for x := int32(0); x < w; x++ {
			i := (y*w + x) * 4
			r := pixels[i]
			g := pixels[i+1]
			b := pixels[i+2]
			a := pixels[i+3]
			img.SetRGBA(int(x), int(h-y-1), color.RGBA{r, g, b, a}) // 画像の上下を反転
		}
	}
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	return png.Encode(file, img)
}

func (m *Msaa) Bind() {
	gl.BindFramebuffer(gl.FRAMEBUFFER, m.msFBO)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
}

func (m *Msaa) Unbind() {
	gl.BindFramebuffer(gl.FRAMEBUFFER, 0)
}

func (m *Msaa) Delete() {
	gl.DeleteFramebuffers(1, &m.msFBO)
	gl.DeleteFramebuffers(1, &m.resolveFBO)
	gl.DeleteRenderbuffers(1, &m.colorBuffer)
	gl.DeleteRenderbuffers(1, &m.depthBuffer)
	gl.DeleteRenderbuffers(1, &m.colorBufferMS)
	gl.DeleteRenderbuffers(1, &m.depthBufferMS)
}
