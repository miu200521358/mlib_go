package mgl

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
	SHADER_BONE_TRANSFORM_MATRIX        = "boneTransformMatrix\x00"
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
)

type ProgramType int

const (
	PROGRAM_TYPE_MODEL ProgramType = iota
	PROGRAM_TYPE_EDGE
	PROGRAM_TYPE_BONE
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
	lightDiffuse         *mmath.MVec3
	lightSpecular        *mmath.MVec3
	CameraPosition       *mmath.MVec3
	LookAtCenterPosition *mmath.MVec3
	FieldOfViewAngle     float32
	Width                int32
	Height               int32
	NearPlane            float32
	FarPlane             float32
	lightPosition        *mmath.MVec3
	lightDirection       *mmath.MVec3
	msaa                 *Msaa
	ModelProgram         uint32
	EdgeProgram          uint32
	BoneProgram          uint32
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
		lightPosition:        &mmath.MVec3{-20, INITIAL_CAMERA_POSITION_Y * 2, INITIAL_CAMERA_POSITION_Z * 2},
		lightDirection:       &mmath.MVec3{-1, -1, -1},
		msaa:                 NewMsaa(int32(width), int32(height)),
	}
	modelProgram, err := shader.newProgram(
		resourceFiles,
		"resources/glsl/model.vert", "resources/glsl/model.frag",
		PROGRAM_TYPE_MODEL)
	if err != nil {
		return nil, err
	}
	shader.ModelProgram = modelProgram
	shader.UseModelProgram()
	shader.initialize(shader.ModelProgram, PROGRAM_TYPE_MODEL)
	shader.Unuse()

	return shader, nil
}

func (s *MShader) Reset() {
	s.CameraPosition = &mmath.MVec3{0.0, INITIAL_CAMERA_POSITION_Y, INITIAL_CAMERA_POSITION_Z}
	s.LookAtCenterPosition = &mmath.MVec3{0.0, INITIAL_LOOK_AT_CENTER_Y, 0.0}
	s.FieldOfViewAngle = FIELD_OF_VIEW_ANGLE
	s.msaa = NewMsaa(s.Width, s.Height)
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
	programType ProgramType,
) (uint32, error) {
	vertexShaderFile, err := fs.ReadFile(resourceFiles, vertexShaderName)
	if err != nil {
		return 0, err
	}

	vertexShaderSource := string(vertexShaderFile)
	if err != nil {
		return 0, err
	}

	fragmentShaderFile, err := fs.ReadFile(resourceFiles, fragmentShaderName)
	if err != nil {
		return 0, err
	}

	fragmentShaderSource := string(fragmentShaderFile)
	if err != nil {
		return 0, err
	}

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

func (s *MShader) initialize(program uint32, programType ProgramType) {
	// light color
	// MMD Light Diffuse は必ず0
	s.lightDiffuse = &mmath.MVec3Zero
	// MMDの照明色そのまま
	s.lightSpecular = &mmath.MVec3{LIGHT_AMBIENT, LIGHT_AMBIENT, LIGHT_AMBIENT}
	// light_diffuse == MMDのambient
	s.lightAmbient = &mmath.MVec4{LIGHT_AMBIENT, LIGHT_AMBIENT, LIGHT_AMBIENT, 1}

	gl.UseProgram(program)

	projection := mgl32.Perspective(
		mgl32.DegToRad(s.FieldOfViewAngle), float32(s.Width)/float32(s.Height), s.NearPlane, s.FarPlane)
	projectionUniform := gl.GetUniformLocation(program, gl.Str(SHADER_MODEL_VIEW_PROJECTION_MATRIX))
	gl.UniformMatrix4fv(projectionUniform, 1, false, &projection[0])

	// カメラの位置
	cameraPosition := s.CameraPosition.Mgl()
	cameraPositionUniform := gl.GetUniformLocation(program, gl.Str(SHADER_CAMERA_POSITION))
	gl.Uniform3fv(cameraPositionUniform, 1, &cameraPosition[0])

	// ライト
	lightDirection := s.lightDirection.Mgl()
	lightDirectionUniform := gl.GetUniformLocation(program, gl.Str(SHADER_LIGHT_DIRECTION))
	gl.Uniform3fv(lightDirectionUniform, 1, &lightDirection[0])

	lookAtCenter := s.LookAtCenterPosition.Mgl()
	camera := mgl32.LookAtV(cameraPosition, lookAtCenter, mgl32.Vec3{0, 1, 0})
	cameraUniform := gl.GetUniformLocation(program, gl.Str(SHADER_MODEL_VIEW_MATRIX))
	gl.UniformMatrix4fv(cameraUniform, 1, false, &camera[0])

	model := mgl32.Ident4()
	modelUniform := gl.GetUniformLocation(program, gl.Str(SHADER_BONE_TRANSFORM_MATRIX))
	gl.UniformMatrix4fv(modelUniform, 1, false, &model[0])

	// // マテリアル設定
	// diffuseUniform := gl.GetUniformLocation(program, gl.Str(SHADER_DIFFUSE))
	// s.DiffuseUniform[programType] = diffuseUniform

	// ambientUniform := gl.GetUniformLocation(program, gl.Str(SHADER_AMBIENT))
	// s.AmbientUniform[programType] = ambientUniform

	// specularUniform := gl.GetUniformLocation(program, gl.Str(SHADER_SPECULAR))
	// s.SpecularUniform[programType] = specularUniform

	// // テクスチャの設定
	// s.UseTextureUniform[programType] = gl.GetUniformLocation(program, gl.Str(SHADER_USE_TEXTURE))
	// // s.TextureUniform[programType] = gl.GetUniformLocation(program, gl.Str("textureSampler\x00"))
	// // s.TextureFactorUniform[programType] = gl.GetUniformLocation(program, gl.Str("textureFactor\x00"))

	// // Toonの設定
	// s.UseToonUniform[programType] = gl.GetUniformLocation(program, gl.Str(SHADER_USE_TOON))
	// // s.ToonUniform[programType] = gl.GetUniformLocation(program, gl.Str("toonSampler\x00"))
	// // s.ToonFactorUniform[programType] = gl.GetUniformLocation(program, gl.Str("toonFactor\x00"))

	// // Sphereの設定
	// s.UseSphereUniform[programType] = gl.GetUniformLocation(program, gl.Str(SHADER_USE_SPHERE))
	// s.SphereModeUniform[programType] = gl.GetUniformLocation(program, gl.Str(SHADER_SPHERE_MODE))
	// // s.SphereUniform[programType] = gl.GetUniformLocation(program, gl.Str("sphereSampler\x00"))
	// // s.SphereFactorUniform[programType] = gl.GetUniformLocation(program, gl.Str("sphereFactor\x00"))

	// // ボーン変形行列用テクスチャ
	// var boneMatrixTextureId uint32
	// gl.GenTextures(1, &boneMatrixTextureId)
	// s.BoneMatrixTextureId[programType] = boneMatrixTextureId
	// s.BoneMatrixTextureUniform[programType] = gl.GetUniformLocation(program, gl.Str("boneMatrixTexture\x00"))
	// s.BoneMatrixTextureWidth[programType] = gl.GetUniformLocation(program, gl.Str("boneMatrixWidth\x00"))
	// s.BoneMatrixTextureHeight[programType] = gl.GetUniformLocation(program, gl.Str("boneMatrixHeight\x00"))

	// if programType == PROGRAM_TYPE_EDGE {
	// 	// エッジシェーダーへの割り当て
	// 	// エッジ設定
	// 	s.EdgeUniform[programType] = gl.GetUniformLocation(program, gl.Str("edge\x00"))
	// 	s.EdgeSizeUniform[programType] = gl.GetUniformLocation(program, gl.Str("edgeSize\x00"))
	// } else if programType == PROGRAM_TYPE_BONE {
	// 	// 選択ボーン色
	// 	s.SelectBoneColorUniform[programType] = gl.GetUniformLocation(program, gl.Str("selectBoneColor\x00"))
	// 	// 非選択ボーン色
	// 	s.UnselectBoneColorUniform[programType] = gl.GetUniformLocation(program, gl.Str("unselectBoneColor\x00"))
	// 	// ボーン数
	// 	s.BoneCountUniform[programType] = gl.GetUniformLocation(program, gl.Str("boneCount\x00"))
	// } else {

	// 	// // テクスチャの設定
	// 	// s.UseTextureUniform[programType] = gl.GetUniformLocation(program, gl.Str("useTexture\x00"))
	// 	// s.TextureUniform[programType] = gl.GetUniformLocation(program, gl.Str("textureSampler\x00"))
	// 	// s.TextureFactorUniform[programType] = gl.GetUniformLocation(program, gl.Str("textureFactor\x00"))

	// 	// // Toonの設定
	// 	// s.UseToonUniform[programType] = gl.GetUniformLocation(program, gl.Str("useToon\x00"))
	// 	// s.ToonUniform[programType] = gl.GetUniformLocation(program, gl.Str("toonSampler\x00"))
	// 	// s.ToonFactorUniform[programType] = gl.GetUniformLocation(program, gl.Str("toonFactor\x00"))

	// 	// // Sphereの設定
	// 	// s.UseSphereUniform[programType] = gl.GetUniformLocation(program, gl.Str("useSphere\x00"))
	// 	// s.SphereModeUniform[programType] = gl.GetUniformLocation(program, gl.Str("sphereMode\x00"))
	// 	// s.SphereUniform[programType] = gl.GetUniformLocation(program, gl.Str("sphereSampler\x00"))
	// 	// s.SphereFactorUniform[programType] = gl.GetUniformLocation(program, gl.Str("sphereFactor\x00"))

	// 	// // ウェイトの描写
	// 	// s.IsShowBoneWeightUniform[programType] = gl.GetUniformLocation(program, gl.Str("isShowBoneWeight\x00"))
	// 	// s.ShowBoneIndexesUniform[programType] = gl.GetUniformLocation(program, gl.Str("showBoneIndexes\x00"))
	// }

	s.Fit(int(s.Width), int(s.Height))
}

func (s *MShader) Fit(
	width int,
	height int,
) {
	s.Width = int32(width)
	s.Height = int32(height)

	// MSAAも作り直し
	s.msaa = NewMsaa(s.Width, s.Height)

	// ビューポートの設定
	gl.Viewport(0, 0, s.Width, s.Height)

	s.updateCamera()
}

func (s *MShader) updateCamera() {
	for p := range []ProgramType{PROGRAM_TYPE_MODEL, PROGRAM_TYPE_EDGE, PROGRAM_TYPE_BONE} {
		programType := ProgramType(p)
		s.Use(programType)

		s.Unuse()
	}
}

func (s *MShader) UseModelProgram() {
	s.Use(PROGRAM_TYPE_MODEL)
}

func (s *MShader) UseEdgeProgram() {
	s.Use(PROGRAM_TYPE_EDGE)
}

func (s *MShader) UseBoneProgram() {
	s.Use(PROGRAM_TYPE_BONE)
}

func (s *MShader) Use(programType ProgramType) {
	switch programType {
	case PROGRAM_TYPE_MODEL:
		gl.UseProgram(s.ModelProgram)
	case PROGRAM_TYPE_EDGE:
		gl.UseProgram(s.EdgeProgram)
	case PROGRAM_TYPE_BONE:
		gl.UseProgram(s.BoneProgram)
	}
}

func (s *MShader) GetPrograms() []uint32 {
	return []uint32{s.ModelProgram}
}

func (s *MShader) Unuse() {
	gl.UseProgram(0)
}

func (s *MShader) Delete() {
	s.DeleteProgram(s.ModelProgram)
	s.DeleteProgram(s.EdgeProgram)
	s.DeleteProgram(s.BoneProgram)
}