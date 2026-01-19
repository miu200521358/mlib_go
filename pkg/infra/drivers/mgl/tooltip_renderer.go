//go:build windows
// +build windows

// 指示: miu200521358
package mgl

import (
	_ "embed"
	"image"
	"image/color"
	"image/draw"
	"unsafe"

	"github.com/go-gl/gl/v4.3-core/gl"
	"github.com/miu200521358/mlib_go/pkg/infra/base/logging"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/math/fixed"
)

//go:embed font/NotoSansJP-Light.ttf
var notoSansJP []byte

// TooltipRenderer はツールチップ描画を担当する。
type TooltipRenderer struct {
	program      uint32
	vao          uint32
	vbo          uint32
	texture      uint32
	width        int
	height       int
	lastText     string
	face         font.Face
	fallbackFace font.Face
}

// NewTooltipRenderer はツールチップ描画を初期化する。
func NewTooltipRenderer() (*TooltipRenderer, error) {
	loader := NewShaderSourceLoader()
	program, err := loader.CreateProgram("glsl/tooltip.vert", "glsl/tooltip.frag")
	if err != nil {
		return nil, err
	}

	var face font.Face
	var fallbackFace font.Face = basicfont.Face7x13

	// NotoSansJP-Light.ttfを埋め込みフォントから読み込み
	tt, err := opentype.Parse(notoSansJP)
	if err != nil {
		logging.DefaultLogger().Warn("NotoSansJPフォントの読み込みに失敗しました: %v", err)
		face = fallbackFace
	} else {
		f, err := opentype.NewFace(tt, &opentype.FaceOptions{
			Size:    12,
			DPI:     72,
			Hinting: font.HintingFull,
		})
		if err != nil {
			logging.DefaultLogger().Warn("NotoSansJPフォントの初期化に失敗しました: %v", err)
			face = fallbackFace
		} else {
			face = f
		}
	}

	renderer := &TooltipRenderer{
		program:      program,
		face:         face,
		fallbackFace: fallbackFace,
	}

	gl.GenVertexArrays(1, &renderer.vao)
	gl.GenBuffers(1, &renderer.vbo)

	gl.BindVertexArray(renderer.vao)
	gl.BindBuffer(gl.ARRAY_BUFFER, renderer.vbo)
	gl.BufferData(gl.ARRAY_BUFFER, 16*4, nil, gl.DYNAMIC_DRAW)

	stride := int32(4 * 4)
	gl.EnableVertexAttribArray(0)
	gl.VertexAttribPointer(0, 2, gl.FLOAT, false, stride, unsafe.Add(unsafe.Pointer(nil), 0))
	gl.EnableVertexAttribArray(1)
	gl.VertexAttribPointer(1, 2, gl.FLOAT, false, stride, unsafe.Add(unsafe.Pointer(nil), 8))

	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	gl.BindVertexArray(0)

	return renderer, nil
}

// Render は指定テキストを描画する。
func (r *TooltipRenderer) Render(text string, cursorX, cursorY float32, winWidth, winHeight int) {
	if text == "" || r.program == 0 {
		return
	}

	if text != r.lastText {
		r.uploadTextTexture(text)
	}

	if r.texture == 0 || r.width == 0 || r.height == 0 {
		return
	}

	vertices := r.buildQuad(cursorX, cursorY, winWidth, winHeight)
	if len(vertices) == 0 {
		return
	}

	gl.BindBuffer(gl.ARRAY_BUFFER, r.vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4, gl.Ptr(&vertices[0]), gl.DYNAMIC_DRAW)
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)

	gl.Disable(gl.DEPTH_TEST)
	gl.UseProgram(r.program)

	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, r.texture)

	uniform := gl.GetUniformLocation(r.program, gl.Str("tooltipTexture\x00"))
	gl.Uniform1i(uniform, 0)

	gl.BindVertexArray(r.vao)
	gl.BindBuffer(gl.ARRAY_BUFFER, r.vbo)
	gl.DrawArrays(gl.TRIANGLE_STRIP, 0, 4)
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	gl.BindVertexArray(0)

	gl.BindTexture(gl.TEXTURE_2D, 0)
	gl.UseProgram(0)
	gl.Enable(gl.DEPTH_TEST)
}

// uploadTextTexture はテキスト用テクスチャを更新する。
func (r *TooltipRenderer) uploadTextTexture(text string) {
	padding := 6

	// NotoSansJPフォントを使用、失敗時はフォールバック
	face := r.face
	if face == nil {
		face = r.fallbackFace
	}

	metrics := face.Metrics()
	ascent := metrics.Ascent.Round()
	lineHeight := (metrics.Ascent + metrics.Descent).Round()

	drawer := font.Drawer{Face: face, Src: image.Black}
	textWidth := drawer.MeasureString(text).Round()
	width := textWidth + padding*2
	height := lineHeight + padding*2

	if width <= 0 {
		width = padding * 2
	}
	if height <= 0 {
		height = padding * 2
	}

	img := image.NewRGBA(image.Rect(0, 0, width, height))
	draw.Draw(img, img.Bounds(), &image.Uniform{C: color.RGBA{255, 255, 255, 200}}, image.Point{}, draw.Src)

	drawer.Dst = img
	drawer.Dot = fixed.Point26_6{X: fixed.I(padding), Y: fixed.I(padding + ascent)}

	// 文字描画
	drawer.DrawString(text)

	if r.texture == 0 {
		gl.GenTextures(1, &r.texture)
		gl.BindTexture(gl.TEXTURE_2D, r.texture)
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	} else {
		gl.BindTexture(gl.TEXTURE_2D, r.texture)
	}

	if len(img.Pix) == 0 {
		gl.BindTexture(gl.TEXTURE_2D, 0)
		return
	}
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA, int32(width), int32(height), 0, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(&img.Pix[0]))
	gl.BindTexture(gl.TEXTURE_2D, 0)

	r.width = width
	r.height = height
	r.lastText = text
}

// buildQuad はツールチップ用の四角形頂点を生成する。
func (r *TooltipRenderer) buildQuad(cursorX, cursorY float32, winWidth, winHeight int) []float32 {
	if winWidth == 0 || winHeight == 0 {
		return nil
	}

	offset := float32(12)
	w := float32(r.width)
	h := float32(r.height)
	if w <= 0 || h <= 0 {
		return nil
	}

	x0 := cursorX + offset
	y0 := cursorY + offset

	if x0 < 4 {
		x0 = 4
	}
	if y0 < 4 {
		y0 = 4
	}

	if x0+w > float32(winWidth)-4 {
		x0 = float32(winWidth) - w - 4
		if x0 < 4 {
			x0 = 4
			w = float32(winWidth) - 8
			if w < 1 {
				w = float32(winWidth)
			}
		}
	}
	if y0+h > float32(winHeight)-4 {
		y0 = float32(winHeight) - h - 4
		if y0 < 4 {
			y0 = 4
			h = float32(winHeight) - 8
			if h < 1 {
				h = float32(winHeight)
			}
		}
	}

	x1 := x0 + w
	y1 := y0 + h

	toClipX := func(v float32) float32 {
		return (v/float32(winWidth))*2 - 1
	}
	toClipY := func(v float32) float32 {
		return 1 - (v/float32(winHeight))*2
	}

	return []float32{
		toClipX(x0), toClipY(y0), 0, 0,
		toClipX(x1), toClipY(y0), 1, 0,
		toClipX(x0), toClipY(y1), 0, 1,
		toClipX(x1), toClipY(y1), 1, 1,
	}
}
