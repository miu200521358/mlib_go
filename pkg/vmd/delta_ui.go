//go:build windows
// +build windows

package vmd

func (md *VertexMorphDelta) GL() []float32 {
	var p0, p1, p2 float32
	if md.Position != nil {
		p := md.Position.GL()
		p0, p1, p2 = p[0], p[1], p[2]
	}
	var ap0, ap1, ap2 float32
	if md.AfterPosition != nil {
		ap := md.AfterPosition.GL()
		ap0, ap1, ap2 = ap[0], ap[1], ap[2]
	}
	// UVは符号関係ないのでそのまま取得する
	var u0x, u0y, u1x, u1y float32
	if md.Uv != nil {
		u0x = float32(md.Uv.GetX())
		u0y = float32(md.Uv.GetY())
	}
	if md.Uv1 != nil {
		u1x = float32(md.Uv1.GetX())
		u1y = float32(md.Uv1.GetY())
	}
	return []float32{
		p0, p1, p2,
		u0x, u0y, 0, 0,
		u1x, u1y, 0, 0,
		ap0, ap1, ap2,
	}
}
