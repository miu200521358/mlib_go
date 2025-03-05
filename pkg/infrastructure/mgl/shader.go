//go:build windows
// +build windows

package mgl

import (
	"github.com/go-gl/gl/v4.4-core/gl"
	"github.com/go-gl/mathgl/mgl32"

	"github.com/miu200521358/mlib_go/pkg/config/mlog"
	"github.com/miu200521358/mlib_go/pkg/domain/mmath"
	"github.com/miu200521358/mlib_go/pkg/domain/rendering"
)

// MShader はOpenGLを使用したシェーダー実装
type MShader struct {
	camera         *rendering.Camera
	width          int
	height         int
	lightPosition  *mmath.MVec3
	lightDirection *mmath.MVec3
	msaa           rendering.IMsaa
	floor          *FloorRenderer // 床レンダラー
	programs       map[rendering.ProgramType]uint32
	boneTextureId  uint32
	shaderLoader   *ShaderLoader
}

// MShaderFactory はOpenGLシェーダーのファクトリー
type MShaderFactory struct{}

// NewMShaderFactory は新しいOpenGLShaderFactoryを作成
func NewMShaderFactory() *MShaderFactory {
	return &MShaderFactory{}
}

// CreateShader は新しいOpenGLShaderを作成
func (f *MShaderFactory) CreateShader(width, height int) (rendering.IShader, error) {
	cam := rendering.NewDefaultCamera(width, height)

	// MSAAの作成
	msaaFactory := NewMsaaFactory()
	msaa := msaaFactory.CreateMsaa(width, height)

	shader := &MShader{
		camera:        cam,
		width:         width,
		height:        height,
		lightPosition: &mmath.MVec3{X: -0.5, Y: -1.0, Z: 0.5},
		msaa:          msaa,
		floor:         NewFloorRenderer(),
		programs:      make(map[rendering.ProgramType]uint32),
		shaderLoader:  NewShaderLoader(),
	}
	shader.lightDirection = shader.lightPosition.Normalized()

	err := shader.initializePrograms()
	if err != nil {
		return nil, err
	}

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

	// 射影行列
	projection := mgl32.Perspective(
		mgl32.DegToRad(s.camera.FieldOfView),
		s.camera.AspectRatio,
		s.camera.NearPlane,
		s.camera.FarPlane,
	)
	projectionUniform := gl.GetUniformLocation(program, gl.Str(ShaderProjectionMatrix))
	gl.UniformMatrix4fv(projectionUniform, 1, false, &projection[0])

	// カメラ位置
	cameraPosition := NewGlVec3(s.camera.Position)
	cameraPositionUniform := gl.GetUniformLocation(program, gl.Str(ShaderCameraPosition))
	gl.Uniform3fv(cameraPositionUniform, 1, &cameraPosition[0])

	// ライト方向
	lightDirection := NewGlVec3(s.lightDirection)
	lightDirectionUniform := gl.GetUniformLocation(program, gl.Str(ShaderLightDirection))
	gl.Uniform3fv(lightDirectionUniform, 1, &lightDirection[0])

	// カメラビュー行列
	cameraUp := NewGlVec3(s.camera.Up)
	lookAtCenter := NewGlVec3(s.camera.LookAtCenter)
	camera := mgl32.LookAtV(cameraPosition, lookAtCenter, cameraUp)
	cameraUniform := gl.GetUniformLocation(program, gl.Str(ShaderViewMatrix))
	gl.UniformMatrix4fv(cameraUniform, 1, false, &camera[0])

	gl.UseProgram(0)
}

// IShaderインターフェースの実装

func (s *MShader) Resize(width, height int) {
	s.width = width
	s.height = height

	gl.Viewport(0, 0, int32(width), int32(height))

	if width != 0 && height != 0 {
		// MSAAのリサイズ
		s.msaa.Resize(width, height)

		s.camera.UpdateAspectRatio(width, height)

		// 全プログラムを更新
		for _, program := range s.programs {
			s.updateProgramProjection(program)
		}
	}
}

func (s *MShader) updateProgramProjection(program uint32) {
	gl.UseProgram(program)
	projection := mgl32.Perspective(
		mgl32.DegToRad(s.camera.FieldOfView),
		s.camera.AspectRatio,
		s.camera.NearPlane,
		s.camera.FarPlane,
	)
	projectionUniform := gl.GetUniformLocation(program, gl.Str(ShaderProjectionMatrix))
	gl.UniformMatrix4fv(projectionUniform, 1, false, &projection[0])
	gl.UseProgram(0)
}

func (s *MShader) GetProgram(programType rendering.ProgramType) uint32 {
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

func (s *MShader) GetBoneTextureID() uint32 {
	return s.boneTextureId
}

func (s *MShader) GetOverrideTextureID() uint32 {
	return 0
	// return s.msaa.OverrideTextureId()
}

func (s *MShader) UpdateCameraSettings(cam *rendering.Camera) {
	s.camera = cam

	for _, program := range s.programs {
		gl.UseProgram(program)

		// カメラ位置
		cameraPosition := NewGlVec3(s.camera.Position)
		cameraPositionUniform := gl.GetUniformLocation(program, gl.Str(ShaderCameraPosition))
		gl.Uniform3fv(cameraPositionUniform, 1, &cameraPosition[0])

		// カメラビュー行列
		cameraUp := NewGlVec3(s.camera.Up)
		lookAtCenter := NewGlVec3(s.camera.LookAtCenter)
		viewMatrix := mgl32.LookAtV(cameraPosition, lookAtCenter, cameraUp)
		viewMatrixUniform := gl.GetUniformLocation(program, gl.Str(ShaderViewMatrix))
		gl.UniformMatrix4fv(viewMatrixUniform, 1, false, &viewMatrix[0])

		// 射影行列
		projection := mgl32.Perspective(
			mgl32.DegToRad(s.camera.FieldOfView),
			s.camera.AspectRatio,
			s.camera.NearPlane,
			s.camera.FarPlane,
		)
		projectionUniform := gl.GetUniformLocation(program, gl.Str(ShaderProjectionMatrix))
		gl.UniformMatrix4fv(projectionUniform, 1, false, &projection[0])
	}

	gl.UseProgram(0)
}

func (s *MShader) GetFieldOfView() float32 {
	return s.camera.FieldOfView
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
	if s.floor != nil {
		s.floor.Delete()
	}

	// MSAAリソースの解放
	if s.msaa != nil {
		s.msaa.Delete()
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
	return 0
	// return s.msaa.OverrideTextureId()
}

// DrawFloor は床を描画する
func (s *MShader) DrawFloor() {
	// 床描画用のプログラム取得
	program := s.programs[rendering.ProgramTypeFloor]

	// 床レンダラーを使用して描画
	if s.floor != nil {
		s.floor.Render(program)
	}
}

// GetMsaa はMSAA機能を取得
func (s *MShader) GetMsaa() rendering.IMsaa {
	return s.msaa
}
