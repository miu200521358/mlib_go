//go:build windows
// +build windows

package mgl

import (
	"sync/atomic"

	"github.com/go-gl/gl/v4.4-core/gl"

	"github.com/miu200521358/mlib_go/pkg/config/mlog"
	"github.com/miu200521358/mlib_go/pkg/domain/mmath"
	"github.com/miu200521358/mlib_go/pkg/domain/rendering"
)

// MShader はOpenGLを使用したシェーダー実装
type MShader struct {
	camera           atomic.Value
	width            int
	height           int
	lightPosition    *mmath.MVec3
	lightDirection   *mmath.MVec3
	msaa             rendering.IMsaa
	floorRenderer    rendering.IFloorRenderer
	overrideRenderer rendering.IOverrideRenderer
	programs         map[rendering.ProgramType]uint32
	boneTextureId    uint32
	sharedTextureId  uint32 // サブウィンドウ共有テクスチャID
	shaderLoader     *ShaderLoader
}

// MShaderFactory はOpenGLシェーダーのファクトリー
type MShaderFactory struct{}

// NewMShaderFactory は新しいOpenGLShaderFactoryを作成
func NewMShaderFactory() *MShaderFactory {
	return &MShaderFactory{}
}

// CreateShader は新しいOpenGLShaderを作成
func (f *MShaderFactory) CreateShader(windowIndex, width, height int) (rendering.IShader, error) {
	cam := rendering.NewDefaultCamera(width, height)

	shader := &MShader{
		width:         width,
		height:        height,
		lightPosition: &mmath.MVec3{X: -0.5, Y: -1.0, Z: 0.5},
		msaa:          NewMsaa(width, height),
		floorRenderer: NewFloorRenderer(),
		programs:      make(map[rendering.ProgramType]uint32),
		shaderLoader:  NewShaderLoader(),
	}

	shader.camera.Store(cam)
	shader.lightDirection = shader.lightPosition.Normalized()

	err := shader.initializePrograms()
	if err != nil {
		return nil, err
	}

	// ウィンドウ番号0（メイン）は isMainWindow = true、それ以外は false
	// initializePrograms() の後に呼ぶ必要がある（ProgramTypeOverrideが必要なため）
	isMainWindow := windowIndex == 0
	shader.overrideRenderer = NewOverrideRenderer(
		width,
		height,
		shader.programs[rendering.ProgramTypeOverride],
		isMainWindow,
		&shader.sharedTextureId,
	)

	return shader, nil
}

// initializePrograms は全シェーダープログラムを初期化
func (s *MShader) initializePrograms() error {
	for programType, config := range ShaderProgramConfigs {
		program, err := s.shaderLoader.CreateProgram(config.VertexShader, config.FragmentShader)
		if err != nil {
			mlog.E("Failed to create %v program: %v", programType, err)
			s.cleanupPrograms()
			return err
		}

		s.programs[programType] = program
		gl.UseProgram(program)
		s.setupProgramUniforms(program)
		gl.UseProgram(0)
	}

	// ボーン行列用テクスチャ生成
	gl.GenTextures(1, &s.boneTextureId)

	return nil
}

// setupProgramUniforms はプログラムのユニフォーム変数を設定
func (s *MShader) setupProgramUniforms(program uint32) {
	gl.UseProgram(program)

	cam := s.Camera()

	// 射影行列
	projection := cam.GetProjectionMatrix(s.width, s.height)
	projectionUniform := gl.GetUniformLocation(program, gl.Str(ShaderProjectionMatrix))
	gl.UniformMatrix4fv(projectionUniform, 1, false, &projection[0])

	// カメラ位置
	cameraPosition := mmath.NewGlVec3(cam.Position)
	cameraPositionUniform := gl.GetUniformLocation(program, gl.Str(ShaderCameraPosition))
	gl.Uniform3fv(cameraPositionUniform, 1, &cameraPosition[0])

	// カメラビュー行列
	camera := cam.GetViewMatrix()
	cameraUniform := gl.GetUniformLocation(program, gl.Str(ShaderViewMatrix))
	gl.UniformMatrix4fv(cameraUniform, 1, false, &camera[0])

	// ライト方向
	lightDirection := mmath.NewGlVec3(s.lightDirection)
	lightDirectionUniform := gl.GetUniformLocation(program, gl.Str(ShaderLightDirection))
	gl.Uniform3fv(lightDirectionUniform, 1, &lightDirection[0])

	gl.UseProgram(0)
}

// IShaderインターフェースの実装

func (s *MShader) Resize(width, height int) {
	if s.width == width && s.height == height {
		return
	}

	s.width = width
	s.height = height

	gl.Viewport(0, 0, int32(width), int32(height))

	if width != 0 && height != 0 {
		// MSAAのリサイズ
		s.msaa.Resize(width, height)

		cam := s.Camera()
		cam.UpdateAspectRatio(width, height)
		s.SetCamera(cam)
	}
}

func (s *MShader) Program(programType rendering.ProgramType) uint32 {
	return s.programs[programType]
}

func (s *MShader) UseProgram(programType rendering.ProgramType) {
	program, exists := s.programs[programType]
	if exists {
		gl.UseProgram(program)
	}
}

func (s *MShader) ResetProgram() {
	gl.UseProgram(0)
}

func (s *MShader) BoneTextureID() uint32 {
	return s.boneTextureId
}

func (s *MShader) OverrideTextureID() uint32 {
	return s.sharedTextureId
}

func (s *MShader) SetCamera(cam *rendering.Camera) {
	s.camera.Store(cam)
}

func (s *MShader) Camera() *rendering.Camera {
	return s.camera.Load().(*rendering.Camera)
}

func (s *MShader) UpdateCamera() {
	cam := s.Camera()

	for _, program := range s.programs {
		gl.UseProgram(program)

		// カメラ位置
		cameraPosition := mmath.NewGlVec3(cam.Position)
		cameraPositionUniform := gl.GetUniformLocation(program, gl.Str(ShaderCameraPosition))
		gl.Uniform3fv(cameraPositionUniform, 1, &cameraPosition[0])

		// カメラビュー行列
		viewMatrix := cam.GetViewMatrix()
		viewMatrixUniform := gl.GetUniformLocation(program, gl.Str(ShaderViewMatrix))
		gl.UniformMatrix4fv(viewMatrixUniform, 1, false, &viewMatrix[0])

		// 射影行列
		projection := cam.GetProjectionMatrix(s.width, s.height)
		projectionUniform := gl.GetUniformLocation(program, gl.Str(ShaderProjectionMatrix))
		gl.UniformMatrix4fv(projectionUniform, 1, false, &projection[0])
	}

	gl.UseProgram(0)
}

func (s *MShader) cleanupPrograms() {
	for _, program := range s.programs {
		if program != 0 {
			gl.DeleteProgram(program)
		}
	}
}

func (s *MShader) Cleanup() {
	s.cleanupPrograms()

	// 床のリソース解放
	if s.floorRenderer != nil {
		s.floorRenderer.Delete()
	}

	// MSAAリソースの解放
	if s.msaa != nil {
		s.msaa.Delete()
	}

	// サブウィンドウのリソース解放
	if s.overrideRenderer != nil {
		s.overrideRenderer.Delete()
	}

	// シェーダープログラム解放
	for _, program := range s.programs {
		gl.DeleteProgram(program)
	}

	// ボーンテクスチャ解放
	if s.boneTextureId != 0 {
		gl.DeleteTextures(1, &s.boneTextureId)
	}
}

func (s *MShader) BoneTextureId() uint32 {
	return s.boneTextureId
}

func (s *MShader) OverrideTextureId() uint32 {
	return s.sharedTextureId
}

func (s *MShader) FloorRenderer() rendering.IFloorRenderer {
	return s.floorRenderer
}

func (s *MShader) OverrideRenderer() rendering.IOverrideRenderer {
	return s.overrideRenderer
}

func (s *MShader) RenderSubWindow() {
	program := s.programs[rendering.ProgramTypeOverride]

	if s.overrideRenderer != nil {
		gl.UseProgram(program)

		s.overrideRenderer.Render()

		gl.UseProgram(0)
	}
}

// Msaa はMSAA機能を取得
func (s *MShader) Msaa() rendering.IMsaa {
	return s.msaa
}
