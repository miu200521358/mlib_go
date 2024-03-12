package mphysics

import (
	"fmt"

	"github.com/go-gl/gl/v4.4-core/gl"

	"github.com/miu200521358/mlib_go/pkg/mbt"
	"github.com/miu200521358/mlib_go/pkg/mgl"
)

// type MDefaultColors struct {
// }

// func NewMDefaultColors() *MDefaultColors {
// 	return &MDefaultColors{}
// }

// func (mdc *MDefaultColors) Swigcptr() uintptr {
// 	return uintptr(unsafe.Pointer(mdc))
// }

type MDebugDrawLiner struct {
	mbt.BtMDebugDrawLiner
	shader *mgl.MShader
}

func NewMDebugDrawLiner(shader *mgl.MShader) *MDebugDrawLiner {
	liner := &MDebugDrawLiner{
		shader: shader,
	}
	liner.BtMDebugDrawLiner = mbt.NewDirectorBtMDebugDrawLiner(liner)
	return liner
}

// func (ddl MDebugDrawLiner) Swigcptr() uintptr {
// 	return uintptr(ddl.SwigcptrBtMDebugDrawLiner)
// }

// func (ddl MDebugDrawLiner) SwigIsBtMDebugDrawLiner() {
// 	fmt.Println("MDebugDrawLiner.SwigIsBtMDebugDrawLiner")
// }

func (ddl MDebugDrawLiner) DrawLine(from mbt.BtVector3, to mbt.BtVector3, color mbt.BtVector3) {
	fmt.Println("MDebugDrawLiner.DrawLine")
	ddl.shader.UsePhysicsProgram()

	// 頂点データを準備
	vertices := []float32{
		from.GetX(), from.GetY(), from.GetZ(),
		to.GetX(), to.GetY(), to.GetZ(),
	}

	// VBOを生成
	var vbo uint32
	gl.GenBuffers(1, &vbo)

	// VBOをバインド
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)

	// 頂点データをVBOに送信
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4, gl.Ptr(vertices), gl.STATIC_DRAW)

	// VAOを生成
	var vao uint32
	gl.GenVertexArrays(1, &vao)

	// VAOをバインド
	gl.BindVertexArray(vao)

	// 頂点属性を有効化
	gl.EnableVertexAttribArray(0)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 0, nil)

	// 色を設定
	colorUniform := gl.GetUniformLocation(ddl.shader.PhysicsProgram, gl.Str(mgl.SHADER_PHYSICS_COLOR))
	gl.Uniform3f(colorUniform, color.GetX(), color.GetY(), color.GetZ())

	alphaUniform := gl.GetUniformLocation(ddl.shader.PhysicsProgram, gl.Str(mgl.SHADER_PHYSICS_ALPHA))
	gl.Uniform1f(alphaUniform, 0.6)

	// 描画
	gl.BindVertexArray(vao)
	gl.DrawArrays(gl.LINES, 0, int32(len(vertices)/3))

	// バッファとVAOをアンバインド
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	gl.BindVertexArray(0)

	ddl.shader.Unuse()
}

// type MDebugView struct {
// 	mbt.SwigcptrBtMDebugDraw
// 	liner *MDebugDrawLiner
// }

// func NewMDebugView(shader *mgl.MShader) *MDebugView {
// 	return &MDebugView{
// 		liner: NewMDebugDrawLiner(shader),
// 	}
// }

// func (mdv *MDebugView) GetLiner() *MDebugDrawLiner {
// 	return mdv.liner
// }

// func NewMDebugView(shader *mgl.MShader) *MDebugView {
// 	viewer := &MDebugView{
// 		BtMDebugDraw: mbt.NewBtMDebugDraw(),
// 		shader:       shader,
// 	}
// 	return viewer
// }

// // func (mdv *MDebugView) Swigcptr() uintptr {
// // 	// mdvのポインターを返す
// // 	fmt.Println("MDebugView.Swigcptr")
// // 	fmt.Printf("%+v\n", mdv)
// // 	ptr := mdv.BtMDebugDraw.Swigcptr()
// // 	fmt.Printf("%v\n", ptr)
// // 	return ptr
// // }

// // func (mdv *MDebugView) SwigIsBtIDebugDraw() {
// // 	// implementation goes here
// // 	fmt.Println("MDebugView.SwigIsBtIDebugDraw")
// // }

// // func (mdv *MDebugView) GetDefaultColors() mbt.BtIDebugDraw_DefaultColors {
// // 	fmt.Println("MDebugView.GetDefaultColors")
// // 	return mdv.defaultColors
// // }

// // func (mdv *MDebugView) SetDefaultColors(arg2 mbt.BtIDebugDraw_DefaultColors) {
// // 	fmt.Println("MDebugView.SetDefaultColors")
// // 	mdv.defaultColors = arg2
// // }

// func (mdv *MDebugView) DrawLine(a ...interface{}) {
// 	fmt.Println("MDebugView.DrawLine")

// 	argc := len(a)
// 	if argc == 3 {
// 		from := a[0].(mbt.BtVector3)
// 		to := a[1].(mbt.BtVector3)
// 		color := a[2].(mbt.BtVector3)

// 	}
// }

// // func (mdv *MDebugView) DrawSphere(a ...interface{}) {
// // 	// implementation goes here
// // 	fmt.Println("MDebugView.DrawSphere")
// // }

// // func (mdv *MDebugView) DrawTriangle(a ...interface{}) {
// // 	// implementation goes here
// // 	fmt.Println("MDebugView.DrawTriangle")
// // }

// // func (mdv *MDebugView) DrawContactPoint(arg2 mbt.BtVector3, arg3 mbt.BtVector3, arg4 float32, arg5 int, arg6 mbt.BtVector3) {
// // 	fmt.Println("MDebugView.DrawContactPoint")

// // 	// Calculate the second point for drawing the contact point
// // 	point1 := mbt.NewBtVector3(
// // 		arg2.GetX()+arg3.GetX()*arg4,
// // 		arg2.GetY()+arg3.GetY()*arg4,
// // 		arg2.GetZ()+arg3.GetZ()*arg4,
// // 	)

// // 	// Draw the contact point
// // 	mdv.DrawLine(arg2, point1, arg6)

// // 	color := mbt.NewBtVector3()

// // 	point2 := mbt.NewBtVector3(
// // 		arg2.GetX()+arg3.GetX()*0.01,
// // 		arg2.GetY()+arg3.GetY()*0.01,
// // 		arg2.GetZ()+arg3.GetZ()*0.01,
// // 	)

// // 	mdv.DrawLine(arg2, point2, color)
// // }

// // func (mdv *MDebugView) ReportErrorWarning(arg2 string) {
// // 	fmt.Println("MDebugView.ReportErrorWarning")
// // 	// implementation goes here
// // }

// // func (mdv *MDebugView) Draw3dText(arg2 mbt.BtVector3, arg3 string) {
// // 	fmt.Println("MDebugView.Draw3dText")
// // 	// implementation goes here
// // }

// // func (mdv *MDebugView) SetDebugMode(arg2 int) {
// // 	fmt.Println("MDebugView.SetDebugMode")
// // 	mdv.debugMode = arg2
// // }

// // func (mdv *MDebugView) GetDebugMode() int {
// // 	fmt.Println("MDebugView.GetDebugMode")
// // 	return mdv.debugMode
// // }

// // func (mdv *MDebugView) DrawAabb(arg2 mbt.BtVector3, arg3 mbt.BtVector3, arg4 mbt.BtVector3) {
// // 	// implementation goes here
// // 	fmt.Println("MDebugView.DrawAabb")
// // }

// // func (mdv *MDebugView) DrawTransform(arg2 mbt.BtTransform, arg3 float32) {
// // 	// implementation goes here
// // 	fmt.Println("MDebugView.DrawTransform")
// // }

// // func (mdv *MDebugView) DrawArc(a ...interface{}) {
// // 	// implementation goes here
// // 	fmt.Println("MDebugView.DrawArc")
// // }

// // func (mdv *MDebugView) DrawSpherePatch(a ...interface{}) {
// // 	// implementation goes here
// // 	fmt.Println("MDebugView.DrawSpherePatch")
// // }

// // func (mdv *MDebugView) DrawBox(a ...interface{}) {
// // 	// implementation goes here
// // 	fmt.Println("MDebugView.DrawBox")
// // }

// // func (mdv *MDebugView) DrawCapsule(arg2 float32, arg3 float32, arg4 int, arg5 mbt.BtTransform, arg6 mbt.BtVector3) {
// // 	// implementation goes here
// // 	fmt.Println("MDebugView.DrawCapsule")
// // }

// // func (mdv *MDebugView) DrawCylinder(arg2 float32, arg3 float32, arg4 int, arg5 mbt.BtTransform, arg6 mbt.BtVector3) {
// // 	// implementation goes here
// // 	fmt.Println("MDebugView.DrawCylinder")
// // }

// // func (mdv *MDebugView) DrawCone(arg2 float32, arg3 float32, arg4 int, arg5 mbt.BtTransform, arg6 mbt.BtVector3) {
// // 	// implementation goes here
// // 	fmt.Println("MDebugView.DrawCone")
// // }

// // func (mdv *MDebugView) DrawPlane(arg2 mbt.BtVector3, arg3 float32, arg4 mbt.BtTransform, arg5 mbt.BtVector3) {
// // 	// implementation goes here
// // 	fmt.Println("MDebugView.DrawPlane")
// // }

// // func (mdv *MDebugView) ClearLines() {
// // 	// implementation goes here
// // 	fmt.Println("MDebugView.ClearLines")
// // }

// // func (mdv *MDebugView) FlushLines() {
// // 	// implementation goes here
// // 	fmt.Println("MDebugView.FlushLines")
// // }
