//go:build windows
// +build windows

package deform

func (md *VertexMorphDelta) GL() []float32 {
	p := md.Position.GL()
	ap := md.AfterPosition.GL()
	// UVは符号関係ないのでそのまま取得する
	return []float32{
		p[0], p[1], p[2],
		float32(md.Uv.GetX()), float32(md.Uv.GetY()),
		float32(md.Uv1.GetX()), float32(md.Uv1.GetY()),
		ap[0], ap[1], ap[2],
	}
}
