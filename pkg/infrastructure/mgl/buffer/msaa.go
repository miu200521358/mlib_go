//go:build windows
// +build windows

package buffer

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

	if gl.CheckFramebufferStatus(gl.FRAMEBUFFER) != gl.FRAMEBUFFER_COMPLETE {
		mlog.F("解決フレームバッファ生成失敗")
	}

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
}

func (msaa *Msaa) Delete() {
	gl.DeleteFramebuffers(1, &msaa.msFBO)
	gl.DeleteFramebuffers(1, &msaa.resolveFBO)
	gl.DeleteRenderbuffers(1, &msaa.colorBuffer)
	gl.DeleteRenderbuffers(1, &msaa.depthBuffer)
	gl.DeleteRenderbuffers(1, &msaa.colorBufferMS)
	gl.DeleteRenderbuffers(1, &msaa.depthBufferMS)
}

func (msaa *Msaa) saveImage(filename string) error {
	w := msaa.width
	h := msaa.height

	pixels := make([]byte, msaa.width*msaa.height*4) // RGBA形式で4バイト/ピクセル
	gl.ReadPixels(0, 0, msaa.width, msaa.height, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(pixels))

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
