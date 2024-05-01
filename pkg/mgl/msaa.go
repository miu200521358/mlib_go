//go:build windows
// +build windows

package mgl

import (
	"github.com/go-gl/gl/v4.4-core/gl"
)

type Msaa struct {
	width             int32
	height            int32
	msaa_samples      int32
	msaa_buffer       uint32
	msaa_color_buffer uint32
	msaa_depth_buffer uint32
}

func NewMsaa(width int32, height int32) *Msaa {
	msaa := &Msaa{
		width:        width,
		height:       height,
		msaa_samples: 4,
	}

	gl.GenFramebuffers(1, &msaa.msaa_buffer)
	gl.BindFramebuffer(gl.FRAMEBUFFER, msaa.msaa_buffer)

	gl.GenRenderbuffers(1, &msaa.msaa_color_buffer)
	gl.BindRenderbuffer(gl.RENDERBUFFER, msaa.msaa_color_buffer)
	gl.RenderbufferStorageMultisample(gl.RENDERBUFFER, msaa.msaa_samples, gl.RGBA, msaa.width, msaa.height)
	gl.FramebufferRenderbuffer(gl.FRAMEBUFFER, gl.COLOR_ATTACHMENT0, gl.RENDERBUFFER, msaa.msaa_color_buffer)

	gl.GenRenderbuffers(1, &msaa.msaa_depth_buffer)
	gl.BindRenderbuffer(gl.RENDERBUFFER, msaa.msaa_depth_buffer)
	gl.RenderbufferStorageMultisample(gl.RENDERBUFFER, msaa.msaa_samples, gl.DEPTH_COMPONENT, msaa.width, msaa.height)
	gl.FramebufferRenderbuffer(gl.FRAMEBUFFER, gl.DEPTH_ATTACHMENT, gl.RENDERBUFFER, msaa.msaa_depth_buffer)

	gl.BindFramebuffer(gl.FRAMEBUFFER, 0)

	return msaa
}

func (m *Msaa) Bind() {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	gl.BindFramebuffer(gl.FRAMEBUFFER, m.msaa_buffer)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
}

func (m *Msaa) Unbind() {
	gl.BindFramebuffer(gl.READ_FRAMEBUFFER, m.msaa_buffer)
	gl.BindFramebuffer(gl.DRAW_FRAMEBUFFER, 0)
	gl.BlitFramebuffer(0, 0, m.width, m.height, 0, 0, m.width, m.height, gl.COLOR_BUFFER_BIT, gl.NEAREST)
}
