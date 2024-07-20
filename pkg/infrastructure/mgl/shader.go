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

const (
	INITIAL_CAMERA_POSITION_Y float64 = 11.0
	INITIAL_CAMERA_POSITION_Z float64 = -40.0
	INITIAL_LOOK_AT_CENTER_Y  float64 = 11.0
	FIELD_OF_VIEW_ANGLE       float32 = 40.0
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

var (
	initialCameraPosition = &mmath.MVec3{X: 0.0, Y: INITIAL_CAMERA_POSITION_Y, Z: INITIAL_CAMERA_POSITION_Z}
	initialLookAtPosition = &mmath.MVec3{X: 0.0, Y: INITIAL_LOOK_AT_CENTER_Y, Z: 0.0}
)

type IShader interface {
	GetProgram(programType ProgramType) uint32
	BoneTextureId() uint32
}

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
	modelProgram          uint32
	edgeProgram           uint32
	boneProgram           uint32
	physicsProgram        uint32
	normalProgram         uint32
	floorProgram          uint32
	wireProgram           uint32
	selectedVertexProgram uint32
	boneTextureId         uint32
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
		shader.modelProgram, err = shader.newProgram(
			glslFiles, "glsl/model.vert", "glsl/model.frag")
		if err != nil {
			mlog.E("Failed to create model program: %v", err)
			return nil
		}
		gl.UseProgram(shader.modelProgram)
		shader.initialize(shader.modelProgram)
		gl.UseProgram(0)
	}

	{
		shader.boneProgram, err = shader.newProgram(
			glslFiles, "glsl/bone.vert", "glsl/bone.frag")
		if err != nil {
			mlog.E("Failed to create bone program: %v", err)
			return nil
		}
		gl.UseProgram(shader.boneProgram)
		shader.initialize(shader.boneProgram)
		gl.UseProgram(0)
	}

	{
		shader.edgeProgram, err = shader.newProgram(
			glslFiles, "glsl/edge.vert", "glsl/edge.frag")
		if err != nil {
			mlog.E("Failed to create edge program: %v", err)
			return nil
		}
		gl.UseProgram(shader.edgeProgram)
		shader.initialize(shader.edgeProgram)
		gl.UseProgram(0)
	}

	{
		shader.physicsProgram, err = shader.newProgram(
			glslFiles, "glsl/physics.vert", "glsl/physics.frag")
		if err != nil {
			mlog.E("Failed to create physics program: %v", err)
			return nil
		}
		gl.UseProgram(shader.physicsProgram)
		shader.initialize(shader.physicsProgram)
		gl.UseProgram(0)
	}

	{
		shader.normalProgram, err = shader.newProgram(
			glslFiles, "glsl/vertex.vert", "glsl/vertex.frag")
		if err != nil {
			mlog.E("Failed to create normal program: %v", err)
			return nil
		}
		gl.UseProgram(shader.normalProgram)
		shader.initialize(shader.normalProgram)
		gl.UseProgram(0)
	}

	{
		shader.floorProgram, err = shader.newProgram(
			glslFiles, "glsl/floor.vert", "glsl/floor.frag")
		if err != nil {
			mlog.E("Failed to create floor program: %v", err)
			return nil
		}
		gl.UseProgram(shader.floorProgram)
		shader.initialize(shader.floorProgram)
		gl.UseProgram(0)
	}

	{
		shader.wireProgram, err = shader.newProgram(
			glslFiles, "glsl/vertex.vert", "glsl/vertex.frag")
		if err != nil {
			mlog.E("Failed to create wire program: %v", err)
			return nil
		}
		gl.UseProgram(shader.wireProgram)
		shader.initialize(shader.wireProgram)
		gl.UseProgram(0)
	}

	{
		shader.selectedVertexProgram, err = shader.newProgram(
			glslFiles, "glsl/vertex.vert", "glsl/vertex.frag")
		if err != nil {
			mlog.E("Failed to create selected vertex program: %v", err)
			return nil
		}
		gl.UseProgram(shader.selectedVertexProgram)
		shader.initialize(shader.selectedVertexProgram)
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
	gl.GenTextures(1, &s.boneTextureId)

	s.Fit(int(s.Width), int(s.Height))

	gl.UseProgram(0)
}

func (s *MShader) BoneTextureId() uint32 {
	return s.boneTextureId
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
		return s.modelProgram
	case PROGRAM_TYPE_EDGE:
		return s.edgeProgram
	case PROGRAM_TYPE_BONE:
		return s.boneProgram
	case PROGRAM_TYPE_PHYSICS:
		return s.physicsProgram
	case PROGRAM_TYPE_NORMAL:
		return s.normalProgram
	case PROGRAM_TYPE_FLOOR:
		return s.floorProgram
	case PROGRAM_TYPE_WIRE:
		return s.wireProgram
	case PROGRAM_TYPE_SELECTED_VERTEX:
		return s.selectedVertexProgram
	}
	return 0
}

func (s *MShader) GetPrograms() []uint32 {
	return []uint32{s.modelProgram, s.edgeProgram, s.boneProgram, s.physicsProgram, s.normalProgram, s.floorProgram, s.wireProgram, s.selectedVertexProgram}
}

func (s *MShader) Delete() {
	s.DeleteProgram(s.modelProgram)
	s.DeleteProgram(s.edgeProgram)
	s.DeleteProgram(s.boneProgram)
	s.DeleteProgram(s.physicsProgram)
	s.DeleteProgram(s.normalProgram)
	s.DeleteProgram(s.floorProgram)
	s.DeleteProgram(s.wireProgram)
	s.DeleteProgram(s.selectedVertexProgram)
	s.Msaa.Delete()
}
