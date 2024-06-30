//go:build windows
// +build windows

package mview

import (
	"embed"
	"fmt"
	"io/fs"
	"strings"

	"github.com/go-gl/gl/v4.4-core/gl"
	"github.com/go-gl/mathgl/mgl32"

	"github.com/miu200521358/mlib_go/pkg/mmath"
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
	SHADER_TEXTURE_MUL_FACTOR           = "textureMulFactor\x00"
	SHADER_TEXTURE_ADD_FACTOR           = "textureAddFactor\x00"
	SHADER_SPHERE_MUL_FACTOR            = "sphereMulFactor\x00"
	SHADER_SPHERE_ADD_FACTOR            = "sphereAddFactor\x00"
	SHADER_TOON_MUL_FACTOR              = "toonMulFactor\x00"
	SHADER_TOON_ADD_FACTOR              = "toonAddFactor\x00"
	SHADER_VERTEX_GL_POSITION           = "gl_Position\x00"
)

type ProgramType int

const (
	PROGRAM_TYPE_MODEL   ProgramType = iota
	PROGRAM_TYPE_EDGE    ProgramType = iota
	PROGRAM_TYPE_BONE    ProgramType = iota
	PROGRAM_TYPE_PHYSICS ProgramType = iota
	PROGRAM_TYPE_NORMAL  ProgramType = iota
	PROGRAM_TYPE_FLOOR   ProgramType = iota
	PROGRAM_TYPE_WIRE    ProgramType = iota
)

const (
	INITIAL_CAMERA_POSITION_Y float64 = 11.0
	INITIAL_CAMERA_POSITION_Z float64 = -40.0
	INITIAL_LOOK_AT_CENTER_Y  float64 = 11.0
	LIGHT_AMBIENT             float64 = 154.0 / 255.0
	FIELD_OF_VIEW_ANGLE       float32 = 40.0
)

type MShader struct {
	lightAmbient         *mmath.MVec4
	CameraPosition       *mmath.MVec3
	LookAtCenterPosition *mmath.MVec3
	FieldOfViewAngle     float32
	Width                int32
	Height               int32
	NearPlane            float32
	FarPlane             float32
	lightPosition        *mmath.MVec3
	lightDirection       *mmath.MVec3
	Msaa                 *Msaa
	ModelProgram         uint32
	EdgeProgram          uint32
	BoneProgram          uint32
	PhysicsProgram       uint32
	NormalProgram        uint32
	FloorProgram         uint32
	WireProgram          uint32
	BoneTextureId        uint32
}

func NewMShader(width, height int, resourceFiles embed.FS) (*MShader, error) {
	shader := &MShader{
		lightAmbient:         &mmath.MVec4{LIGHT_AMBIENT, LIGHT_AMBIENT, LIGHT_AMBIENT, 1},
		CameraPosition:       &mmath.MVec3{0.0, INITIAL_CAMERA_POSITION_Y, INITIAL_CAMERA_POSITION_Z},
		LookAtCenterPosition: &mmath.MVec3{0.0, INITIAL_LOOK_AT_CENTER_Y, 0.0},
		FieldOfViewAngle:     FIELD_OF_VIEW_ANGLE,
		Width:                int32(width),
		Height:               int32(height),
		NearPlane:            0.1,
		FarPlane:             1000.0,
		lightPosition:        &mmath.MVec3{-0.5, -1.0, 0.5},
		Msaa:                 NewMsaa(int32(width), int32(height)),
	}
	shader.lightDirection = shader.lightPosition.Normalized()

	{
		modelProgram, err := shader.newProgram(
			resourceFiles, "resources/glsl/model.vert", "resources/glsl/model.frag")
		if err != nil {
			return nil, err
		}
		shader.ModelProgram = modelProgram
		shader.Use(PROGRAM_TYPE_MODEL)
		shader.initialize(shader.ModelProgram)
		shader.Unuse()
	}

	{
		boneProgram, err := shader.newProgram(
			resourceFiles, "resources/glsl/bone.vert", "resources/glsl/bone.frag")
		if err != nil {
			return nil, err
		}
		shader.BoneProgram = boneProgram
		shader.Use(PROGRAM_TYPE_BONE)
		shader.initialize(shader.BoneProgram)
		shader.Unuse()
	}

	{
		edgeProgram, err := shader.newProgram(
			resourceFiles, "resources/glsl/edge.vert", "resources/glsl/edge.frag")
		if err != nil {
			return nil, err
		}
		shader.EdgeProgram = edgeProgram
		shader.Use(PROGRAM_TYPE_EDGE)
		shader.initialize(shader.EdgeProgram)
		shader.Unuse()
	}

	{
		physicsProgram, err := shader.newProgram(
			resourceFiles, "resources/glsl/physics.vert", "resources/glsl/physics.frag")
		if err != nil {
			return nil, err
		}
		shader.PhysicsProgram = physicsProgram
		shader.Use(PROGRAM_TYPE_PHYSICS)
		shader.initialize(shader.PhysicsProgram)
		shader.Unuse()
	}

	{
		normalProgram, err := shader.newProgram(
			resourceFiles, "resources/glsl/normal.vert", "resources/glsl/normal.frag")
		if err != nil {
			return nil, err
		}
		shader.NormalProgram = normalProgram
		shader.Use(PROGRAM_TYPE_NORMAL)
		shader.initialize(shader.NormalProgram)
		shader.Unuse()
	}

	{
		floorProgram, err := shader.newProgram(
			resourceFiles, "resources/glsl/floor.vert", "resources/glsl/floor.frag")
		if err != nil {
			return nil, err
		}
		shader.FloorProgram = floorProgram
		shader.Use(PROGRAM_TYPE_FLOOR)
		shader.initialize(shader.FloorProgram)
		shader.Unuse()
	}

	{
		wireProgram, err := shader.newProgram(
			resourceFiles, "resources/glsl/wire.vert", "resources/glsl/wire.frag")
		if err != nil {
			return nil, err
		}
		shader.WireProgram = wireProgram
		shader.Use(PROGRAM_TYPE_WIRE)
		shader.initialize(shader.WireProgram)
		shader.Unuse()
	}

	return shader, nil
}

func (s *MShader) Reset() {
	s.CameraPosition = &mmath.MVec3{0.0, INITIAL_CAMERA_POSITION_Y, INITIAL_CAMERA_POSITION_Z}
	s.LookAtCenterPosition = &mmath.MVec3{0.0, INITIAL_LOOK_AT_CENTER_Y, 0.0}
	s.FieldOfViewAngle = FIELD_OF_VIEW_ANGLE
	s.Resize(int(s.Width), int(s.Height))
}

func (s *MShader) Resize(width, height int) {
	s.Width = int32(width)
	s.Height = int32(height)
	s.Msaa = NewMsaa(s.Width, s.Height)
	s.updateCamera()
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
	resourceFiles embed.FS,
	vertexShaderName, fragmentShaderName string,
) (uint32, error) {
	vertexShaderFile, err := fs.ReadFile(resourceFiles, vertexShaderName)
	if err != nil {
		return 0, err
	}

	vertexShaderSource := string(vertexShaderFile)

	fragmentShaderFile, err := fs.ReadFile(resourceFiles, fragmentShaderName)
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
	// light color
	// MMD Light Diffuse は必ず0
	// MMDの照明色そのまま: light_diffuse == MMDのambient
	s.lightAmbient = &mmath.MVec4{LIGHT_AMBIENT, LIGHT_AMBIENT, LIGHT_AMBIENT, 1}

	gl.UseProgram(program)

	projection := mgl32.Perspective(
		mgl32.DegToRad(s.FieldOfViewAngle), float32(s.Width)/float32(s.Height), s.NearPlane, s.FarPlane)
	projectionUniform := gl.GetUniformLocation(program, gl.Str(SHADER_MODEL_VIEW_PROJECTION_MATRIX))
	gl.UniformMatrix4fv(projectionUniform, 1, false, &projection[0])

	// カメラの位置
	cameraPosition := s.CameraPosition.GL()
	cameraPositionUniform := gl.GetUniformLocation(program, gl.Str(SHADER_CAMERA_POSITION))
	gl.Uniform3fv(cameraPositionUniform, 1, &cameraPosition[0])

	// ライト
	lightDirection := s.lightDirection.GL()
	lightDirectionUniform := gl.GetUniformLocation(program, gl.Str(SHADER_LIGHT_DIRECTION))
	gl.Uniform3fv(lightDirectionUniform, 1, &lightDirection[0])

	// カメラ中心
	lookAtCenter := s.LookAtCenterPosition.GL()
	camera := mgl32.LookAtV(cameraPosition, lookAtCenter, mgl32.Vec3{0, 1, 0})
	cameraUniform := gl.GetUniformLocation(program, gl.Str(SHADER_MODEL_VIEW_MATRIX))
	gl.UniformMatrix4fv(cameraUniform, 1, false, &camera[0])

	// ボーン行列用テクスチャ生成
	gl.GenTextures(1, &s.BoneTextureId)

	s.Fit(int(s.Width), int(s.Height))

	gl.UseProgram(0)
}

func (s *MShader) Fit(
	width int,
	height int,
) {
	s.Width = int32(width)
	s.Height = int32(height)

	// MSAAも作り直し
	s.Msaa = NewMsaa(s.Width, s.Height)

	// ビューポートの設定
	gl.Viewport(0, 0, s.Width, s.Height)

	s.updateCamera()
}

func (s *MShader) updateCamera() {
	for _, p := range s.GetPrograms() {
		s.Use(ProgramType(p))

		s.Unuse()
	}
}

func (s *MShader) Use(programType ProgramType) {
	switch programType {
	case PROGRAM_TYPE_MODEL:
		gl.UseProgram(s.ModelProgram)
	case PROGRAM_TYPE_EDGE:
		gl.UseProgram(s.EdgeProgram)
	case PROGRAM_TYPE_BONE:
		gl.UseProgram(s.BoneProgram)
	case PROGRAM_TYPE_PHYSICS:
		gl.UseProgram(s.PhysicsProgram)
	case PROGRAM_TYPE_NORMAL:
		gl.UseProgram(s.NormalProgram)
	case PROGRAM_TYPE_FLOOR:
		gl.UseProgram(s.FloorProgram)
	case PROGRAM_TYPE_WIRE:
		gl.UseProgram(s.WireProgram)
	}
}

func (s *MShader) GetPrograms() []uint32 {
	return []uint32{s.ModelProgram, s.EdgeProgram, s.BoneProgram, s.PhysicsProgram, s.NormalProgram, s.FloorProgram, s.WireProgram}
}

func (s *MShader) Unuse() {
	gl.UseProgram(0)
}

func (s *MShader) Delete() {
	s.DeleteProgram(s.ModelProgram)
	s.DeleteProgram(s.EdgeProgram)
	s.DeleteProgram(s.BoneProgram)
	s.DeleteProgram(s.PhysicsProgram)
	s.DeleteProgram(s.NormalProgram)
	s.DeleteProgram(s.FloorProgram)
	s.DeleteProgram(s.WireProgram)
	s.Msaa.Delete()
}
