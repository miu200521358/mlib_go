package mphysics

import (
	"github.com/go-gl/gl/v4.4-core/gl"

	"github.com/miu200521358/mlib_go/pkg/mbt"
	"github.com/miu200521358/mlib_go/pkg/mgl"
)

func NewConstBtMDefaultColors() mbt.BtMDefaultColors {
	return mbt.NewBtMDefaultColors(
		mbt.NewBtVector3(float32(1.0), float32(0.0), float32(0.0)), // activeObject	(物理剛体)
		mbt.NewBtVector3(float32(0.5), float32(0.5), float32(0.0)), // deactivatedObject
		mbt.NewBtVector3(float32(0.5), float32(0.0), float32(0.5)), // wantsDeactivationObject
		mbt.NewBtVector3(float32(0.0), float32(0.5), float32(0.5)), // disabledDeactivationObject
		mbt.NewBtVector3(float32(0.0), float32(1.0), float32(0.0)), // disabledSimulationObject	(ボーン追従剛体)
		mbt.NewBtVector3(float32(1.0), float32(1.0), float32(0.0)), // aabb
		mbt.NewBtVector3(float32(0.0), float32(0.0), float32(1.0)), // contactPoint
	)
}

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

func (ddl MDebugDrawLiner) DrawLine(from mbt.BtVector3, to mbt.BtVector3, color mbt.BtVector3) {
	// fmt.Println("MDebugDrawLiner.DrawLine")
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

	// fmt.Printf("ddl.color: %.5f, %.5f, %.5f\n", color.GetX(), color.GetY(), color.GetZ())

	// 色を設定
	colorUniform := gl.GetUniformLocation(ddl.shader.PhysicsProgram, gl.Str(mgl.SHADER_COLOR))
	gl.Uniform3f(colorUniform, color.GetX(), color.GetY(), color.GetZ())

	alphaUniform := gl.GetUniformLocation(ddl.shader.PhysicsProgram, gl.Str(mgl.SHADER_ALPHA))
	gl.Uniform1f(alphaUniform, 0.6)

	// 描画
	gl.BindVertexArray(vao)
	gl.DrawArrays(gl.LINES, 0, int32(len(vertices)/3))

	// バッファとVAOをアンバインド
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	gl.BindVertexArray(0)

	ddl.shader.Unuse()
}
