//go:build windows
// +build windows

// 指示: miu200521358
package mgl

import (
	"sync/atomic"

	"github.com/go-gl/gl/v4.3-core/gl"

	"github.com/miu200521358/mlib_go/pkg/adapter/graphics_api"
	"github.com/miu200521358/mlib_go/pkg/domain/mmath"
	"github.com/miu200521358/mlib_go/pkg/infra/base/logging"
	"gonum.org/v1/gonum/spatial/r3"
)

// Shader はOpenGLを使用したシェーダー実装。
type Shader struct {
	camera           atomic.Value
	width            int
	height           int
	lightPosition    *mmath.Vec3
	lightDirection   *mmath.Vec3
	msaa             graphics_api.IMsaa
	floorRenderer    graphics_api.IFloorRenderer
	overrideRenderer graphics_api.IOverrideRenderer
	programs         map[graphics_api.ProgramType]uint32
	boneTextureId    uint32
	sharedTextureId  uint32
	shaderLoader     *ShaderSourceLoader
}

// ShaderFactory はOpenGLシェーダーのファクトリー。
type ShaderFactory struct {
	windowIndex int
}

// NewShaderFactory はシェーダーファクトリーを生成する。
func NewShaderFactory(windowIndex int) *ShaderFactory {
	return &ShaderFactory{windowIndex: windowIndex}
}

// CreateShader はOpenGLシェーダーを生成する。
func (f *ShaderFactory) CreateShader(width, height int) (graphics_api.IShader, error) {
	cam := graphics_api.NewDefaultCamera(width, height)
	lightPos := mmath.Vec3{r3.Vec{X: -0.5, Y: -1.0, Z: 0.5}}
	lightDir := lightPos.Normalized()

	msaa := NewMsaaBuffer(width, height)
	if err := msaaInitError(msaa); err != nil {
		return nil, err
	}

	shader := &Shader{
		width:         width,
		height:        height,
		lightPosition: &lightPos,
		lightDirection: &lightDir,
		msaa:          msaa,
		floorRenderer: NewFloorRenderer(),
		programs:      make(map[graphics_api.ProgramType]uint32),
		shaderLoader:  NewShaderSourceLoader(),
	}

	shader.camera.Store(cam)

	if err := shader.initializePrograms(); err != nil {
		return nil, err
	}

	isMainWindow := f.windowIndex == 0
	shader.overrideRenderer = NewOverrideRenderer(
		width,
		height,
		shader.programs[graphics_api.ProgramTypeOverride],
		isMainWindow,
	)

	return shader, nil
}

// initializePrograms はシェーダープログラムを初期化する。
func (s *Shader) initializePrograms() error {
	for programType, config := range SHADER_PROGRAM_CONFIGS {
		program, err := s.shaderLoader.CreateProgram(config.VertexShader, config.FragmentShader)
		if err != nil {
			logging.DefaultLogger().Error("シェーダープログラム生成に失敗: %v (%v)", programType, err)
			s.cleanupPrograms()
			return err
		}

		s.programs[programType] = program
		gl.UseProgram(program)
		s.setupProgramUniforms(program)
		gl.UseProgram(0)
	}

	gl.GenTextures(1, &s.boneTextureId)
	return nil
}

// setupProgramUniforms はユニフォーム初期値を設定する。
func (s *Shader) setupProgramUniforms(program uint32) {
	gl.UseProgram(program)

	cam := s.Camera()

	projection := NewGlMat4(cam.GetProjectionMatrix(s.width, s.height))
	projectionUniform := gl.GetUniformLocation(program, gl.Str(ShaderProjectionMatrix))
	gl.UniformMatrix4fv(projectionUniform, 1, false, &projection[0])

	cameraPosition := NewGlVec3(cam.Position)
	cameraPositionUniform := gl.GetUniformLocation(program, gl.Str(ShaderCameraPosition))
	gl.Uniform3fv(cameraPositionUniform, 1, &cameraPosition[0])

	camera := NewGlMat4(cam.GetViewMatrix())
	cameraUniform := gl.GetUniformLocation(program, gl.Str(ShaderViewMatrix))
	gl.UniformMatrix4fv(cameraUniform, 1, false, &camera[0])

	lightDirection := NewGlVec3(s.lightDirection)
	lightDirectionUniform := gl.GetUniformLocation(program, gl.Str(ShaderLightDirection))
	gl.Uniform3fv(lightDirectionUniform, 1, &lightDirection[0])

	gl.UseProgram(0)
}

// Resize は描画サイズを更新する。
func (s *Shader) Resize(width, height int) {
	if s.width == width && s.height == height {
		return
	}

	s.width = width
	s.height = height

	gl.Viewport(0, 0, int32(width), int32(height))

	if width != 0 && height != 0 {
		s.msaa.Resize(width, height)
		s.overrideRenderer.Resize(width, height)

		cam := s.Camera()
		cam.UpdateAspectRatio(width, height)
		s.SetCamera(cam)
	}
}

// Program はプログラムIDを返す。
func (s *Shader) Program(programType graphics_api.ProgramType) uint32 {
	return s.programs[programType]
}

// UseProgram は指定プログラムを使用する。
func (s *Shader) UseProgram(programType graphics_api.ProgramType) {
	program, exists := s.programs[programType]
	if exists {
		gl.UseProgram(program)
	}
}

// ResetProgram はプログラム利用を解除する。
func (s *Shader) ResetProgram() {
	gl.UseProgram(0)
}

// BoneTextureID はボーン行列テクスチャIDを返す。
func (s *Shader) BoneTextureID() uint32 {
	return s.boneTextureId
}

// OverrideTextureID は共有テクスチャIDを返す。
func (s *Shader) OverrideTextureID() uint32 {
	return s.sharedTextureId
}

// SetCamera はカメラを設定する。
func (s *Shader) SetCamera(cam *graphics_api.Camera) {
	s.camera.Store(cam)
}

// Camera はカメラを返す。
func (s *Shader) Camera() *graphics_api.Camera {
	return s.camera.Load().(*graphics_api.Camera)
}

// UpdateCamera はカメラ情報を反映する。
func (s *Shader) UpdateCamera() {
	cam := s.Camera()

	for _, program := range s.programs {
		gl.UseProgram(program)

		cameraPosition := NewGlVec3(cam.Position)
		cameraPositionUniform := gl.GetUniformLocation(program, gl.Str(ShaderCameraPosition))
		gl.Uniform3fv(cameraPositionUniform, 1, &cameraPosition[0])

		viewMatrix := NewGlMat4(cam.GetViewMatrix())
		viewMatrixUniform := gl.GetUniformLocation(program, gl.Str(ShaderViewMatrix))
		gl.UniformMatrix4fv(viewMatrixUniform, 1, false, &viewMatrix[0])

		projection := NewGlMat4(cam.GetProjectionMatrix(s.width, s.height))
		projectionUniform := gl.GetUniformLocation(program, gl.Str(ShaderProjectionMatrix))
		gl.UniformMatrix4fv(projectionUniform, 1, false, &projection[0])
	}

	gl.UseProgram(0)
}

// Msaa はMSAA実装を返す。
func (s *Shader) Msaa() graphics_api.IMsaa {
	return s.msaa
}

// SetMsaa はMSAA実装を設定する。
func (s *Shader) SetMsaa(msaa graphics_api.IMsaa) {
	s.msaa = msaa
}

// FloorRenderer は床描画を返す。
func (s *Shader) FloorRenderer() graphics_api.IFloorRenderer {
	return s.floorRenderer
}

// OverrideRenderer はオーバーライド描画を返す。
func (s *Shader) OverrideRenderer() graphics_api.IOverrideRenderer {
	return s.overrideRenderer
}

// cleanupPrograms はプログラムを削除する。
func (s *Shader) cleanupPrograms() {
	for _, program := range s.programs {
		if program != 0 {
			gl.DeleteProgram(program)
		}
	}
}

// Cleanup は関連リソースを解放する。
func (s *Shader) Cleanup() {
	s.cleanupPrograms()

	if s.floorRenderer != nil {
		s.floorRenderer.Delete()
	}

	if s.msaa != nil {
		s.msaa.Delete()
	}

	if s.overrideRenderer != nil {
		s.overrideRenderer.Delete()
	}

	if s.boneTextureId != 0 {
		gl.DeleteTextures(1, &s.boneTextureId)
	}
}

// msaaInitError はMSAA初期化エラーを取得する。
func msaaInitError(msaa graphics_api.IMsaa) error {
	if msaa == nil {
		return nil
	}
	typed, ok := msaa.(interface{ initError() error })
	if !ok {
		return nil
	}
	return typed.initError()
}
