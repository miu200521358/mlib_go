package face

type T struct {
	Index       int
	VertexIndex [3]int
}

func NewFace(index, vertexIndex0, vertexIndex1, vertexIndex2 int) *T {
	return &T{
		Index:       index,
		VertexIndex: [3]int{vertexIndex0, vertexIndex1, vertexIndex2},
	}
}

// Copy
func (v *T) Copy() *T {
	copied := *v
	return &copied
}
