package mmath

type MRect struct {
	x      int
	y      int
	width  int
	height int
}

func NewMRect(x, y, width, height int) *MRect {
	return &MRect{
		x:      x,
		y:      y,
		width:  width,
		height: height,
	}
}

func (m *MRect) SetX(x int) {
	m.x = x
}

func (m *MRect) GetX() int {
	return m.x
}

func (m *MRect) SetY(y int) {
	m.y = y
}

func (m *MRect) GetY() int {
	return m.y
}

func (m *MRect) SetWidth(width int) {
	m.width = width
}

func (m *MRect) GetWidth() int {
	return m.width
}

func (m *MRect) SetHeight(height int) {
	m.height = height
}

func (m *MRect) GetHeight() int {
	return m.height
}
