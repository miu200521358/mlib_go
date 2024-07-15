//go:build windows
// +build windows

package renderer

import (
	"github.com/go-gl/mathgl/mgl32"
	"github.com/miu200521358/mlib_go/pkg/domain/pmx"
)

// SDEF用パラメーターを返す
func GetSdefParams(d pmx.IDeform) (mgl32.Vec3, mgl32.Vec3, mgl32.Vec3) {
	// pmx.Sdefにキャストできる場合
	if _, ok := d.(*pmx.Sdef); ok {
		s := d.(*pmx.Sdef)
		// CがR0とR1より先にいかないよう、重みに基づいて補正
		copiedSdefR0 := s.SdefR0.Copy()
		copiedSdefR1 := s.SdefR1.Copy()
		copiedSdefCR0 := s.SdefC.Copy()
		copiedSdefCR1 := s.SdefC.Copy()

		weight := copiedSdefR0.MulScalar(s.Weights[0]).Add(copiedSdefR1.MulScalar(1 - s.Weights[0]))
		sdefR0 := copiedSdefCR0.Add(s.SdefR0).Sub(weight)
		sdefR1 := copiedSdefCR1.Add(s.SdefR1).Sub(weight)

		return s.SdefC.GL(), sdefR0.GL(), sdefR1.GL()
	}
	return mgl32.Vec3{}, mgl32.Vec3{}, mgl32.Vec3{}
}
