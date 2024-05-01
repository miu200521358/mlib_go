//go:build windows
// +build windows

package pmx

import "github.com/go-gl/mathgl/mgl32"

// SDEF用パラメーターを返す
func (s *Sdef) GetSdefParams() (*mgl32.Vec3, *mgl32.Vec3, *mgl32.Vec3) {
	// CがR0とR1より先にいかないよう、重みに基づいて補正
	weight := s.SdefR0.MuledScalar(s.Weights[0]).Added(s.SdefR1.MuledScalar(1 - s.Weights[0]))
	sdefR0 := s.SdefC.Added(s.SdefR0).Subed(weight)
	sdefR1 := s.SdefC.Added(s.SdefR1).Subed(weight)

	return s.SdefC.GL(), sdefR0.GL(), sdefR1.GL()
}
