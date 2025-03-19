package rendering

type IFloorRenderer interface {
	Bind()

	Unbind()

	Render()

	Delete()
}
