//go:build windows
// +build windows

package mgl

import (
	"embed"
	"fmt"
	"io/fs"
	"strings"

	"github.com/go-gl/gl/v4.4-core/gl"
	"github.com/go-gl/mathgl/mgl32"

	"github.com/miu200521358/mlib_go/pkg/domain/mmath"
	"github.com/miu200521358/mlib_go/pkg/infrastructure/mgl/buffer"
	"github.com/miu200521358/mlib_go/pkg/mutils/mlog"
)

const (
	SHADER_BONE_MATRIX_TEXTURE          = "boneMatrixTexture\x00"
	SHADER_BONE_MATRIX_TEXTURE_WIDTH    = "boneMatrixWidth\x00"
	SHADER_BONE_MATRIX_TEXTURE_HEIGHT   = "boneMatrixHeight\x00"
	SHADER_MODEL_VIEW_MATRIX            = "modelViewMatrix\x00"
	SHADER_MODEL_VIEW_PROJECTION_MATRIX = "modelViewProjectionMatrix\x00"
	SHADER_CAMERA_POSITION              = "cameraPosition\x00"
	SHADER_LIGHT_DIRECTION              = "lightDirection\x00"
	SHADER_DIFFUSE                      = "diffuse\x00"
	SHADER_AMBIENT                      = "ambient\x00"
	SHADER_SPECULAR                     = "specular\x00"
	SHADER_TEXTURE_SAMPLER              = "textureSampler\x00"
	SHADER_TOON_SAMPLER                 = "toonSampler\x00"
	SHADER_SPHERE_SAMPLER               = "sphereSampler\x00"
	SHADER_USE_TEXTURE                  = "useTexture\x00"
	SHADER_USE_TOON                     = "useToon\x00"
	SHADER_USE_SPHERE                   = "useSphere\x00"
	SHADER_SPHERE_MODE                  = "sphereMode\x00"
	SHADER_MORPH_TEXTURE_FACTOR         = "textureFactor\x00"
	SHADER_MORPH_TOON_FACTOR            = "toonFactor\x00"
	SHADER_MORPH_SPHERE_FACTOR          = "sphereFactor\x00"
	SHADER_COLOR                        = "color\x00"
	SHADER_ALPHA                        = "alpha\x00"
	SHADER_EDGE_COLOR                   = "edgeColor\x00"
	SHADER_EDGE_SIZE                    = "edgeSize\x00"
	SHADER_VERTEX_GL_POSITION           = "gl_Position\x00"
)

type ProgramType int

const (
	PROGRAM_TYPE_MODEL           ProgramType = iota
	PROGRAM_TYPE_EDGE            ProgramType = iota
	PROGRAM_TYPE_BONE            ProgramType = iota
	PROGRAM_TYPE_PHYSICS         ProgramType = iota
	PROGRAM_TYPE_NORMAL          ProgramType = iota
	PROGRAM_TYPE_FLOOR           ProgramType = iota
	PROGRAM_TYPE_WIRE            ProgramType = iota
	PROGRAM_TYPE_SELECTED_VERTEX ProgramType = iota
)

const (
	INITIAL_CAMERA_POSITION_Y float64 = 11.0
	INITIAL_CAMERA_POSITION_Z float64 = -40.0
	INITIAL_LOOK_AT_CENTER_Y  float64 = 11.0
	LIGHT_AMBIENT             float64 = 154.0 / 255.0
	FIELD_OF_VIEW_ANGLE       float32 = 40.0
)

var (
	initialCameraPosition = &mmath.MVec3{X: 0.0, Y: INITIAL_CAMERA_POSITION_Y, Z: INITIAL_CAMERA_POSITION_Z}
	initialLookAtPosition = &mmath.MVec3{X: 0.0, Y: INITIAL_LOOK_AT_CENTER_Y, Z: 0.0}
)

type MShader struct {
	CameraPosition        *mmath.MVec3
	LookAtCenterPosition  *mmath.MVec3
	FieldOfViewAngle      float32
	Width                 int
	Height                int
	NearPlane             float32
	FarPlane              float32
	lightPosition         *mmath.MVec3
	lightDirection        *mmath.MVec3
	Msaa                  *buffer.Msaa
	ModelProgram          uint32
	EdgeProgram           uint32
	BoneProgram           uint32
	PhysicsProgram        uint32
	NormalProgram         uint32
	FloorProgram          uint32
	WireProgram           uint32
	SelectedVertexProgram uint32
	BoneTextureId         uint32
	floor                 *MFloor
}

//go:embed glsl/*
var glslFiles embed.FS

func NewMShader(width, height int) *MShader {
	shader := &MShader{
		CameraPosition:       initialCameraPosition.Copy(),
		LookAtCenterPosition: initialLookAtPosition.Copy(),
		FieldOfViewAngle:     FIELD_OF_VIEW_ANGLE,
		Width:                width,
		Height:               height,
		NearPlane:            0.1,
		FarPlane:             1000.0,
		lightPosition:        &mmath.MVec3{X: -0.5, Y: -1.0, Z: 0.5},
		Msaa:                 buffer.NewMsaa(width, height),
		floor:                newMFloor(),
	}
	shader.lightDirection = shader.lightPosition.Normalized()

	var err error
	{
		shader.ModelProgram, err = shader.newProgram(
			glslFiles, "glsl/model.vert", "glsl/model.frag")
		if err != nil {
			mlog.E("Failed to create model program: %v", err)
			return nil
		}
		gl.UseProgram(shader.ModelProgram)
		shader.initialize(shader.ModelProgram)
		gl.UseProgram(0)
	}

	{
		shader.BoneProgram, err = shader.newProgram(
			glslFiles, "glsl/bone.vert", "glsl/bone.frag")
		if err != nil {
			mlog.E("Failed to create bone program: %v", err)
			return nil
		}
		gl.UseProgram(shader.BoneProgram)
		shader.initialize(shader.BoneProgram)
		gl.UseProgram(0)
	}

	{
		shader.EdgeProgram, err = shader.newProgram(
			glslFiles, "glsl/edge.vert", "glsl/edge.frag")
		if err != nil {
			mlog.E("Failed to create edge program: %v", err)
			return nil
		}
		gl.UseProgram(shader.EdgeProgram)
		shader.initialize(shader.EdgeProgram)
		gl.UseProgram(0)
	}

	{
		shader.PhysicsProgram, err = shader.newProgram(
			glslFiles, "glsl/physics.vert", "glsl/physics.frag")
		if err != nil {
			mlog.E("Failed to create physics program: %v", err)
			return nil
		}
		gl.UseProgram(shader.PhysicsProgram)
		shader.initialize(shader.PhysicsProgram)
		gl.UseProgram(0)
	}

	{
		shader.NormalProgram, err = shader.newProgram(
			glslFiles, "glsl/vertex.vert", "glsl/vertex.frag")
		if err != nil {
			mlog.E("Failed to create normal program: %v", err)
			return nil
		}
		gl.UseProgram(shader.NormalProgram)
		shader.initialize(shader.NormalProgram)
		gl.UseProgram(0)
	}

	{
		shader.FloorProgram, err = shader.newProgram(
			glslFiles, "glsl/floor.vert", "glsl/floor.frag")
		if err != nil {
			mlog.E("Failed to create floor program: %v", err)
			return nil
		}
		gl.UseProgram(shader.FloorProgram)
		shader.initialize(shader.FloorProgram)
		gl.UseProgram(0)
	}

	{
		shader.WireProgram, err = shader.newProgram(
			glslFiles, "glsl/vertex.vert", "glsl/vertex.frag")
		if err != nil {
			mlog.E("Failed to create wire program: %v", err)
			return nil
		}
		gl.UseProgram(shader.WireProgram)
		shader.initialize(shader.WireProgram)
		gl.UseProgram(0)
	}

	{
		shader.SelectedVertexProgram, err = shader.newProgram(
			glslFiles, "glsl/vertex.vert", "glsl/vertex.frag")
		if err != nil {
			mlog.E("Failed to create selected vertex program: %v", err)
			return nil
		}
		gl.UseProgram(shader.SelectedVertexProgram)
		shader.initialize(shader.SelectedVertexProgram)
		gl.UseProgram(0)
	}

	return shader
}

func (s *MShader) Reset() {
	s.CameraPosition = initialCameraPosition.Copy()
	s.LookAtCenterPosition = initialLookAtPosition.Copy()
	s.FieldOfViewAngle = FIELD_OF_VIEW_ANGLE
	s.Resize(int(s.Width), int(s.Height))
}

func (s *MShader) Resize(width, height int) {
	s.Width = width
	s.Height = height
	s.Msaa = buffer.NewMsaa(s.Width, s.Height)
}

func (s *MShader) DeleteProgram(program uint32) {
	if program != 0 {
		gl.DeleteProgram(program)
	}
}

func (s *MShader) compileShader(shaderName, source string, shaderType uint32) (uint32, error) {
	shader := gl.CreateShader(shaderType)

	csources, free := gl.Strs(source + "\x00")
	gl.ShaderSource(shader, 1, csources, nil)
	free()
	gl.CompileShader(shader)

	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status != gl.TRUE {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(log))

		return 0, fmt.Errorf("failed to compile %v: %v", shaderName, log)
	}

	return shader, nil
}

func (s *MShader) newProgram(
	glslFiles embed.FS,
	vertexShaderName, fragmentShaderName string,
) (uint32, error) {
	vertexShaderFile, err := fs.ReadFile(glslFiles, vertexShaderName)
	if err != nil {
		return 0, err
	}

	vertexShaderSource := string(vertexShaderFile)

	fragmentShaderFile, err := fs.ReadFile(glslFiles, fragmentShaderName)
	if err != nil {
		return 0, err
	}

	fragmentShaderSource := string(fragmentShaderFile)

	vertexShader, err := s.compileShader(vertexShaderName, vertexShaderSource, gl.VERTEX_SHADER)
	if err != nil {
		return 0, err
	}

	fragmentShader, err := s.compileShader(fragmentShaderName, fragmentShaderSource, gl.FRAGMENT_SHADER)
	if err != nil {
		return 0, err
	}

	program := gl.CreateProgram()

	gl.AttachShader(program, vertexShader)
	gl.AttachShader(program, fragmentShader)
	gl.LinkProgram(program)

	var status int32
	gl.GetProgramiv(program, gl.LINK_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetProgramiv(program, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetProgramInfoLog(program, logLength, nil, gl.Str(log))

		return 0, fmt.Errorf("failed to link program: %v", log)
	}

	gl.DeleteShader(vertexShader)
	gl.DeleteShader(fragmentShader)

	return program, nil
}

func (s *MShader) initialize(program uint32) {
	gl.UseProgram(program)

	projection := mgl32.Perspective(
		mgl32.DegToRad(s.FieldOfViewAngle), float32(s.Width)/float32(s.Height), s.NearPlane, s.FarPlane)
	projectionUniform := gl.GetUniformLocation(program, gl.Str(SHADER_MODEL_VIEW_PROJECTION_MATRIX))
	gl.UniformMatrix4fv(projectionUniform, 1, false, &projection[0])

	// カメラの位置
	cameraPosition := NewGlVec3(s.CameraPosition)
	cameraPositionUniform := gl.GetUniformLocation(program, gl.Str(SHADER_CAMERA_POSITION))
	gl.Uniform3fv(cameraPositionUniform, 1, &cameraPosition[0])

	// ライト
	lightDirection := NewGlVec3(s.lightDirection)
	lightDirectionUniform := gl.GetUniformLocation(program, gl.Str(SHADER_LIGHT_DIRECTION))
	gl.Uniform3fv(lightDirectionUniform, 1, &lightDirection[0])

	// カメラ中心
	lookAtCenter := NewGlVec3(s.LookAtCenterPosition)
	camera := mgl32.LookAtV(cameraPosition, lookAtCenter, mgl32.Vec3{0, 1, 0})
	cameraUniform := gl.GetUniformLocation(program, gl.Str(SHADER_MODEL_VIEW_MATRIX))
	gl.UniformMatrix4fv(cameraUniform, 1, false, &camera[0])

	// ボーン行列用テクスチャ生成
	gl.GenTextures(1, &s.BoneTextureId)

	s.Fit(int(s.Width), int(s.Height))

	gl.UseProgram(0)
}

func (s *MShader) Fit(width int, height int) {
	s.Width = width
	s.Height = height

	// MSAAも作り直し
	s.Msaa = buffer.NewMsaa(s.Width, s.Height)

	// ビューポートの設定
	gl.Viewport(0, 0, int32(s.Width), int32(s.Height))
}

func (s *MShader) GetProgram(programType ProgramType) uint32 {
	switch programType {
	case PROGRAM_TYPE_MODEL:
		return s.ModelProgram
	case PROGRAM_TYPE_EDGE:
		return s.EdgeProgram
	case PROGRAM_TYPE_BONE:
		return s.BoneProgram
	case PROGRAM_TYPE_PHYSICS:
		return s.PhysicsProgram
	case PROGRAM_TYPE_NORMAL:
		return s.NormalProgram
	case PROGRAM_TYPE_FLOOR:
		return s.FloorProgram
	case PROGRAM_TYPE_WIRE:
		return s.WireProgram
	case PROGRAM_TYPE_SELECTED_VERTEX:
		return s.SelectedVertexProgram
	}
	return 0
}

func (s *MShader) GetPrograms() []uint32 {
	return []uint32{s.ModelProgram, s.EdgeProgram, s.BoneProgram, s.PhysicsProgram, s.NormalProgram, s.FloorProgram, s.WireProgram, s.SelectedVertexProgram}
}

func (s *MShader) Delete() {
	s.DeleteProgram(s.ModelProgram)
	s.DeleteProgram(s.EdgeProgram)
	s.DeleteProgram(s.BoneProgram)
	s.DeleteProgram(s.PhysicsProgram)
	s.DeleteProgram(s.NormalProgram)
	s.DeleteProgram(s.FloorProgram)
	s.DeleteProgram(s.WireProgram)
	s.DeleteProgram(s.SelectedVertexProgram)
	s.Msaa.Delete()
}
