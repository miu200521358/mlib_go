//go:build windows
// +build windows

package pmx

import (
	"github.com/go-gl/mathgl/mgl32"
	"github.com/miu200521358/mlib_go/pkg/infrastructure/mgl"
)

// SDEF用パラメーターを返す
func (s *Sdef) GetSdefParams() (mgl32.Vec3, mgl32.Vec3, mgl32.Vec3) {
	// CがR0とR1より先にいかないよう、重みに基づいて補正
	copiedSdefR0 := s.SdefR0.Copy()
	copiedSdefR1 := s.SdefR1.Copy()
	copiedSdefCR0 := s.SdefC.Copy()
	copiedSdefCR1 := s.SdefC.Copy()

	weight := copiedSdefR0.MulScalar(s.Weights[0]).Add(copiedSdefR1.MulScalar(1 - s.Weights[0]))
	sdefR0 := copiedSdefCR0.Add(s.SdefR0).Sub(weight)
	sdefR1 := copiedSdefCR1.Add(s.SdefR1).Sub(weight)

	return mgl.NewGlVec3FromMVec3(s.SdefC), mgl.NewGlVec3FromMVec3(sdefR0), mgl.NewGlVec3FromMVec3(sdefR1)
}