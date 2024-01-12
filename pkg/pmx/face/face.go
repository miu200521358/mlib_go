package face

type Face struct {
	Index       int
	VertexIndex [3]int
}

func NewFace(index, vertexIndex0, vertexIndex1, vertexIndex2 int) *Face {
	return &Face{
		Index:       index,
		VertexIndex: [3]int{vertexIndex0, vertexIndex1, vertexIndex2},
	}
}
