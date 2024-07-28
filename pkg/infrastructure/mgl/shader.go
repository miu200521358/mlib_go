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
	SHADER_BONE_MATRIX_TEXTURE        = "boneMatrixTexture\x00"
	SHADER_BONE_MATRIX_TEXTURE_WIDTH  = "boneMatrixWidth\x00"
	SHADER_BONE_MATRIX_TEXTURE_HEIGHT = "boneMatrixHeight\x00"
	SHADER_VIEW_MATRIX                = "viewMatrix\x00"
	SHADER_PROJECTION_MATRIX          = "projectionMatrix\x00"
	SHADER_CAMERA_POSITION            = "cameraPosition\x00"
	SHADER_LIGHT_DIRECTION            = "lightDirection\x00"
	SHADER_DIFFUSE                    = "diffuse\x00"
	SHADER_AMBIENT                    = "ambient\x00"
	SHADER_SPECULAR                   = "specular\x00"
	SHADER_TEXTURE_SAMPLER            = "textureSampler\x00"
	SHADER_TOON_SAMPLER               = "toonSampler\x00"
	SHADER_SPHERE_SAMPLER             = "sphereSampler\x00"
	SHADER_USE_TEXTURE                = "useTexture\x00"
	SHADER_USE_TOON                   = "useToon\x00"
	SHADER_USE_SPHERE                 = "useSphere\x00"
	SHADER_SPHERE_MODE                = "sphereMode\x00"
	SHADER_MORPH_TEXTURE_FACTOR       = "textureFactor\x00"
	SHADER_MORPH_TOON_FACTOR          = "toonFactor\x00"
	SHADER_MORPH_SPHERE_FACTOR        = "sphereFactor\x00"
	SHADER_COLOR                      = "color\x00"
	SHADER_ALPHA                      = "alpha\x00"
	SHADER_EDGE_COLOR                 = "edgeColor\x00"
	SHADER_EDGE_SIZE                  = "edgeSize\x00"
	SHADER_CURSOR_POSITIONS           = "cursorPositions\x00"
	SHADER_LINE_WIDTH                 = "lineWidth\x00"
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
	PROGRAM_TYPE_OVERRIDE        ProgramType = iota
	PROGRAM_TYPE_CURSOR          ProgramType = iota
)

var (
	initialCameraPosition = &mmath.MVec3{X: 0.0, Y: INITIAL_CAMERA_POSITION_Y, Z: INITIAL_CAMERA_POSITION_Z}
	initialLookAtPosition = &mmath.MVec3{X: 0.0, Y: INITIAL_LOOK_AT_CENTER_Y, Z: 0.0}
	initialCameraUp       = &mmath.MVec3{X: 0.0, Y: 1.0, Z: 0.0}
)

type IShader interface {
	Program(programType ProgramType) uint32
	BoneTextureId() uint32
	OverrideTextureId() uint32
	Resize(width, height int)
}

type MShader struct {
	CameraPosition        *mmath.MVec3
	LookAtCenterPosition  *mmath.MVec3
	CameraUp              *mmath.MVec3
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
	overrideProgram       uint32
	cursorProgram         uint32
	boneTextureId         uint32
	floor                 *MFloor
}

//go:embed glsl/*
var glslFiles embed.FS

func NewMShader(width, height int) *MShader {
	shader := &MShader{
		CameraPosition:       initialCameraPosition.Copy(),
		LookAtCenterPosition: initialLookAtPosition.Copy(),
		CameraUp:             initialCameraUp.Copy(),
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
		shader.modelProgram, err = shader.newProgram("glsl/model.vert", "glsl/model.frag", "")
		if err != nil {
			mlog.E("Failed to create model program: %v", err)
			return nil
		}
		gl.UseProgram(shader.modelProgram)
		shader.initialize(shader.modelProgram)
		gl.UseProgram(0)
	}

	{
		shader.boneProgram, err = shader.newProgram("glsl/bone.vert", "glsl/bone.frag", "")
		if err != nil {
			mlog.E("Failed to create bone program: %v", err)
			return nil
		}
		gl.UseProgram(shader.boneProgram)
		shader.initialize(shader.boneProgram)
		gl.UseProgram(0)
	}

	{
		shader.edgeProgram, err = shader.newProgram("glsl/edge.vert", "glsl/edge.frag", "")
		if err != nil {
			mlog.E("Failed to create edge program: %v", err)
			return nil
		}
		gl.UseProgram(shader.edgeProgram)
		shader.initialize(shader.edgeProgram)
		gl.UseProgram(0)
	}

	{
		shader.physicsProgram, err = shader.newProgram("glsl/physics.vert", "glsl/physics.frag", "")
		if err != nil {
			mlog.E("Failed to create physics program: %v", err)
			return nil
		}
		gl.UseProgram(shader.physicsProgram)
		shader.initialize(shader.physicsProgram)
		gl.UseProgram(0)
	}

	{
		shader.normalProgram, err = shader.newProgram("glsl/vertex.vert", "glsl/vertex.frag", "")
		if err != nil {
			mlog.E("Failed to create normal program: %v", err)
			return nil
		}
		gl.UseProgram(shader.normalProgram)
		shader.initialize(shader.normalProgram)
		gl.UseProgram(0)
	}

	{
		shader.floorProgram, err = shader.newProgram("glsl/floor.vert", "glsl/floor.frag", "")
		if err != nil {
			mlog.E("Failed to create floor program: %v", err)
			return nil
		}
		gl.UseProgram(shader.floorProgram)
		shader.initialize(shader.floorProgram)
		gl.UseProgram(0)
	}

	{
		shader.wireProgram, err = shader.newProgram("glsl/vertex.vert", "glsl/vertex.frag", "")
		if err != nil {
			mlog.E("Failed to create wire program: %v", err)
			return nil
		}
		gl.UseProgram(shader.wireProgram)
		shader.initialize(shader.wireProgram)
		gl.UseProgram(0)
	}

	{
		shader.selectedVertexProgram, err = shader.newProgram("glsl/vertex.vert", "glsl/vertex.frag", "")
		if err != nil {
			mlog.E("Failed to create selected vertex program: %v", err)
			return nil
		}
		gl.UseProgram(shader.selectedVertexProgram)
		shader.initialize(shader.selectedVertexProgram)
		gl.UseProgram(0)
	}

	{
		shader.overrideProgram, err = shader.newProgram("glsl/override.vert", "glsl/override.frag", "")
		if err != nil {
			mlog.E("Failed to create selected override program: %v", err)
			return nil
		}
		gl.UseProgram(shader.overrideProgram)
		shader.initialize(shader.overrideProgram)
		gl.UseProgram(0)
	}

	{
		shader.cursorProgram, err = shader.newProgram("glsl/cursor.vert", "glsl/cursor.frag", "")
		if err != nil {
			mlog.E("Failed to create selected cursor program: %v", err)
			return nil
		}
		gl.UseProgram(shader.cursorProgram)
		shader.initialize(shader.cursorProgram)
		gl.UseProgram(0)
	}

	return shader
}

func (shader *MShader) Reset(isResize bool) {
	shader.CameraPosition = initialCameraPosition.Copy()
	shader.LookAtCenterPosition = initialLookAtPosition.Copy()
	shader.CameraUp = initialCameraUp.Copy()
	shader.FieldOfViewAngle = FIELD_OF_VIEW_ANGLE
	if isResize {
		shader.Resize(int(shader.Width), int(shader.Height))
	}
}

func (shader *MShader) DeleteProgram(program uint32) {
	if program != 0 {
		gl.DeleteProgram(program)
	}
}

func (shader *MShader) compileShader(shaderName, source string, shaderType uint32) (uint32, error) {
	s := gl.CreateShader(shaderType)

	csources, free := gl.Strs(source + "\x00")
	gl.ShaderSource(s, 1, csources, nil)
	free()
	gl.CompileShader(s)

	var status int32
	gl.GetShaderiv(s, gl.COMPILE_STATUS, &status)
	if status != gl.TRUE {
		var logLength int32
		gl.GetShaderiv(s, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(s, logLength, nil, gl.Str(log))

		return 0, fmt.Errorf("failed to compile %v: %v", shaderName, log)
	}

	return s, nil
}

func (shader *MShader) newProgram(
	vertexShaderName, fragmentShaderName, geometryShaderName string,
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

	vertexShader, err := shader.compileShader(vertexShaderName, vertexShaderSource, gl.VERTEX_SHADER)
	if err != nil {
		return 0, err
	}

	fragmentShader, err := shader.compileShader(fragmentShaderName, fragmentShaderSource, gl.FRAGMENT_SHADER)
	if err != nil {
		return 0, err
	}

	program := gl.CreateProgram()

	var geometryShader uint32
	if geometryShaderName != "" {
		geometryShaderFile, err := fs.ReadFile(glslFiles, geometryShaderName)
		if err != nil {
			return 0, err
		}

		geometryShaderSource := string(geometryShaderFile)

		geometryShader, err = shader.compileShader(geometryShaderName, geometryShaderSource, gl.GEOMETRY_SHADER)
		if err != nil {
			return 0, err
		}
	}

	gl.AttachShader(program, vertexShader)
	gl.AttachShader(program, fragmentShader)
	if geometryShaderName != "" && geometryShader != 0 {
		gl.AttachShader(program, geometryShader)
	}
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

func (shader *MShader) initialize(program uint32) {
	gl.UseProgram(program)

	projection := mgl32.Perspective(
		mgl32.DegToRad(shader.FieldOfViewAngle), float32(shader.Width)/float32(shader.Height),
		shader.NearPlane, shader.FarPlane)
	projectionUniform := gl.GetUniformLocation(program, gl.Str(SHADER_PROJECTION_MATRIX))
	gl.UniformMatrix4fv(projectionUniform, 1, false, &projection[0])

	// カメラの位置
	cameraPosition := NewGlVec3(shader.CameraPosition)
	cameraPositionUniform := gl.GetUniformLocation(program, gl.Str(SHADER_CAMERA_POSITION))
	gl.Uniform3fv(cameraPositionUniform, 1, &cameraPosition[0])

	// ライト
	lightDirection := NewGlVec3(shader.lightDirection)
	lightDirectionUniform := gl.GetUniformLocation(program, gl.Str(SHADER_LIGHT_DIRECTION))
	gl.Uniform3fv(lightDirectionUniform, 1, &lightDirection[0])

	// カメラUP
	cameraUp := NewGlVec3(shader.CameraUp)

	// カメラ中心
	lookAtCenter := NewGlVec3(shader.LookAtCenterPosition)
	camera := mgl32.LookAtV(cameraPosition, lookAtCenter, cameraUp)
	cameraUniform := gl.GetUniformLocation(program, gl.Str(SHADER_VIEW_MATRIX))
	gl.UniformMatrix4fv(cameraUniform, 1, false, &camera[0])

	// ボーン行列用テクスチャ生成
	gl.GenTextures(1, &shader.boneTextureId)

	shader.Resize(int(shader.Width), int(shader.Height))

	gl.UseProgram(0)
}

func (shader *MShader) BoneTextureId() uint32 {
	return shader.boneTextureId
}

func (shader *MShader) OverrideTextureId() uint32 {
	return shader.Msaa.OverrideTextureId()
}

func (shader *MShader) Resize(width int, height int) {
	shader.Width = width
	shader.Height = height

	// ビューポートの設定
	gl.Viewport(0, 0, int32(shader.Width), int32(shader.Height))

	if width != 0 && height != 0 {
		// MSAAも作り直し
		shader.Msaa.Delete()
		shader.Msaa = buffer.NewMsaa(shader.Width, shader.Height)
	}
}

func (shader *MShader) Program(programType ProgramType) uint32 {
	switch programType {
	case PROGRAM_TYPE_MODEL:
		return shader.modelProgram
	case PROGRAM_TYPE_EDGE:
		return shader.edgeProgram
	case PROGRAM_TYPE_BONE:
		return shader.boneProgram
	case PROGRAM_TYPE_PHYSICS:
		return shader.physicsProgram
	case PROGRAM_TYPE_NORMAL:
		return shader.normalProgram
	case PROGRAM_TYPE_FLOOR:
		return shader.floorProgram
	case PROGRAM_TYPE_WIRE:
		return shader.wireProgram
	case PROGRAM_TYPE_SELECTED_VERTEX:
		return shader.selectedVertexProgram
	case PROGRAM_TYPE_OVERRIDE:
		return shader.overrideProgram
	case PROGRAM_TYPE_CURSOR:
		return shader.cursorProgram
	}
	return 0
}

func (shader *MShader) Programs() []uint32 {
	return []uint32{shader.modelProgram, shader.edgeProgram, shader.boneProgram, shader.physicsProgram,
		shader.normalProgram, shader.floorProgram, shader.wireProgram, shader.selectedVertexProgram,
		shader.overrideProgram, shader.cursorProgram}
}

func (shader *MShader) Delete() {
	shader.DeleteProgram(shader.modelProgram)
	shader.DeleteProgram(shader.edgeProgram)
	shader.DeleteProgram(shader.boneProgram)
	shader.DeleteProgram(shader.physicsProgram)
	shader.DeleteProgram(shader.normalProgram)
	shader.DeleteProgram(shader.floorProgram)
	shader.DeleteProgram(shader.wireProgram)
	shader.DeleteProgram(shader.selectedVertexProgram)
	shader.DeleteProgram(shader.overrideProgram)
	shader.DeleteProgram(shader.cursorProgram)
	shader.Msaa.Delete()
}
