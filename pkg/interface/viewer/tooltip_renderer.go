package viewer

import (
	"image"
	"image/color"
	"image/draw"
	"unsafe"

	"github.com/go-gl/gl/v4.4-core/gl"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"

	"github.com/miu200521358/mlib_go/pkg/infrastructure/mgl"
)

type tooltipRenderer struct {
	program  uint32
	vao      uint32
	vbo      uint32
	texture  uint32
	width    int
	height   int
	lastText string
}

func newTooltipRenderer() (*tooltipRenderer, error) {
	loader := mgl.NewShaderLoader()
	program, err := loader.CreateProgram("glsl/tooltip.vert", "glsl/tooltip.frag")
	if err != nil {
		return nil, err
	}

	renderer := &tooltipRenderer{program: program}

	gl.GenVertexArrays(1, &renderer.vao)
	gl.GenBuffers(1, &renderer.vbo)

	gl.BindVertexArray(renderer.vao)
	gl.BindBuffer(gl.ARRAY_BUFFER, renderer.vbo)
	gl.BufferData(gl.ARRAY_BUFFER, 16*4, nil, gl.DYNAMIC_DRAW)

	stride := int32(4 * 4)
	gl.EnableVertexAttribArray(0)
	gl.VertexAttribPointer(0, 2, gl.FLOAT, false, stride, unsafe.Pointer(uintptr(0)))
	gl.EnableVertexAttribArray(1)
	gl.VertexAttribPointer(1, 2, gl.FLOAT, false, stride, unsafe.Pointer(uintptr(8)))

	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	gl.BindVertexArray(0)

	return renderer, nil
}

func (r *tooltipRenderer) Render(text string, cursorX, cursorY float32, winWidth, winHeight int) {
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

func (r *tooltipRenderer) uploadTextTexture(text string) {
	padding := 6
	face := basicfont.Face7x13
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

func (r *tooltipRenderer) buildQuad(cursorX, cursorY float32, winWidth, winHeight int) []float32 {
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
