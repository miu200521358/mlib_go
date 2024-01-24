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

type VsLayout int

const (
	VS_LAYOUT_POSITION_ID VsLayout = iota
	VS_LAYOUT_NORMAL_ID
	VS_LAYOUT_UV_ID
	VS_LAYOUT_EXTEND_UV_ID
	VS_LAYOUT_EDGE_ID
	VS_LAYOUT_BONE_ID
	VS_LAYOUT_WEIGHT_ID
	VS_LAYOUT_MORPH_POS_ID
	VS_LAYOUT_MORPH_UV_ID
	VS_LAYOUT_MORPH_UV1_ID
	VS_LAYOUT_MORPH_AFTER_POS_ID
)

type ProgramType int

const (
	PROGRAM_TYPE_MODEL ProgramType = iota
	PROGRAM_TYPE_EDGE
	PROGRAM_TYPE_BONE
)

const (
	INITIAL_VERTICAL_DEGREES  float64 = 40.0
	INITIAL_CAMERA_POSITION_Y float64 = 11.0
	INITIAL_CAMERA_POSITION_Z float64 = -40.0
	INITIAL_LOOK_AT_CENTER_Y  float64 = 11.0
	INITIAL_CAMERA_POSITION_X float64 = 40.0
	LIGHT_AMBIENT             float64 = 154.0 / 255.0
)

type MShader struct {
	lightAmbient                     *mmath.MVec4
	lightDiffuse                     *mmath.MVec3
	lightSpecular                    *mmath.MVec3
	initialCameraPosition            *mmath.MVec3
	initialCameraOffsetPosition      *mmath.MVec3
	initialLookAtCenterPosition      *mmath.MVec3
	width                            int
	height                           int
	nearPlane                        int
	farPlane                         int
	lightPosition                    *mmath.MVec3
	lightDirection                   *mmath.MVec3
	msaa                             *Msaa
	BoneMatrixTextureUniform         map[ProgramType]int32
	BoneMatrixTextureId              map[ProgramType]uint32
	BoneMatrixTextureWidth           map[ProgramType]int32
	BoneMatrixTextureHeight          map[ProgramType]int32
	ModelViewMatrixUniform           map[ProgramType]int32
	ModelViewProjectionMatrixUniform map[ProgramType]int32
	LightDirectionUniform            map[ProgramType]int32
	CameraVecUniform                 map[ProgramType]int32
	DiffuseColorUniform              map[ProgramType]int32
	DiffuseAlphaUniform              map[ProgramType]int32
	AmbientColorUniform              map[ProgramType]int32
	SpecularColorUniform             map[ProgramType]int32
	SpecularPowerUniform             map[ProgramType]int32
	SelectBoneColorUniform           map[ProgramType]int32
	UnselectBoneColorUniform         map[ProgramType]int32
	EdgeColorUniform                 map[ProgramType]int32
	EdgeAlphaUniform                 map[ProgramType]int32
	EdgeSizeUniform                  map[ProgramType]int32
	UseTextureUniform                map[ProgramType]int32
	TextureUniform                   map[ProgramType]int32
	TextureFactorUniform             map[ProgramType]int32
	UseToonUniform                   map[ProgramType]int32
	ToonUniform                      map[ProgramType]int32
	ToonFactorUniform                map[ProgramType]int32
	UseSphereUniform                 map[ProgramType]int32
	SphereModeUniform                map[ProgramType]int32
	SphereUniform                    map[ProgramType]int32
	SphereFactorUniform              map[ProgramType]int32
	BoneCountUniform                 map[ProgramType]int32
	IsShowBoneWeightUniform          map[ProgramType]int32
	ShowBoneIndexesUniform           map[ProgramType]int32
	ModelProgram                     uint32
	EdgeProgram                      uint32
	BoneProgram                      uint32
	ModelUniform                     int32
}

func NewMShader(width, height int, resourceFiles embed.FS) (*MShader, error) {
	shader := &MShader{
		lightAmbient:                     &mmath.MVec4{LIGHT_AMBIENT, LIGHT_AMBIENT, LIGHT_AMBIENT, 1},
		initialCameraPosition:            &mmath.MVec3{0.0, INITIAL_CAMERA_POSITION_Y, INITIAL_CAMERA_POSITION_Z},
		initialCameraOffsetPosition:      &mmath.MVec3{},
		initialLookAtCenterPosition:      &mmath.MVec3{0.0, INITIAL_LOOK_AT_CENTER_Y, 0.0},
		width:                            width,
		height:                           height,
		nearPlane:                        1,
		farPlane:                         100,
		lightPosition:                    &mmath.MVec3{-20, INITIAL_CAMERA_POSITION_Y * 2, INITIAL_CAMERA_POSITION_Z * 2},
		lightDirection:                   &mmath.MVec3{-1, -1, -1},
		msaa:                             NewMsaa(int32(width), int32(height)),
		BoneMatrixTextureUniform:         make(map[ProgramType]int32, 0),
		BoneMatrixTextureId:              make(map[ProgramType]uint32, 0),
		BoneMatrixTextureWidth:           make(map[ProgramType]int32, 0),
		BoneMatrixTextureHeight:          make(map[ProgramType]int32, 0),
		ModelViewMatrixUniform:           make(map[ProgramType]int32, 0),
		ModelViewProjectionMatrixUniform: make(map[ProgramType]int32, 0),
		LightDirectionUniform:            make(map[ProgramType]int32, 0),
		CameraVecUniform:                 make(map[ProgramType]int32, 0),
		DiffuseColorUniform:              make(map[ProgramType]int32, 0),
		DiffuseAlphaUniform:              make(map[ProgramType]int32, 0),
		AmbientColorUniform:              make(map[ProgramType]int32, 0),
		SpecularColorUniform:             make(map[ProgramType]int32, 0),
		SpecularPowerUniform:             make(map[ProgramType]int32, 0),
		EdgeColorUniform:                 make(map[ProgramType]int32, 0),
		EdgeAlphaUniform:                 make(map[ProgramType]int32, 0),
		EdgeSizeUniform:                  make(map[ProgramType]int32, 0),
		SelectBoneColorUniform:           make(map[ProgramType]int32, 0),
		UnselectBoneColorUniform:         make(map[ProgramType]int32, 0),
		UseTextureUniform:                make(map[ProgramType]int32, 0),
		TextureUniform:                   make(map[ProgramType]int32, 0),
		TextureFactorUniform:             make(map[ProgramType]int32, 0),
		UseToonUniform:                   make(map[ProgramType]int32, 0),
		ToonUniform:                      make(map[ProgramType]int32, 0),
		ToonFactorUniform:                make(map[ProgramType]int32, 0),
		UseSphereUniform:                 make(map[ProgramType]int32, 0),
		SphereModeUniform:                make(map[ProgramType]int32, 0),
		SphereUniform:                    make(map[ProgramType]int32, 0),
		SphereFactorUniform:              make(map[ProgramType]int32, 0),
		BoneCountUniform:                 make(map[ProgramType]int32, 0),
		IsShowBoneWeightUniform:          make(map[ProgramType]int32, 0),
		ShowBoneIndexesUniform:           make(map[ProgramType]int32, 0),
	}
	modelProgram, err := shader.newProgram(
		resourceFiles,
		"resources/glsl/vertex.vert", "resources/glsl/vertex.frag",
		PROGRAM_TYPE_MODEL)
	if err != nil {
		return nil, err
	}
	shader.ModelProgram = modelProgram
	shader.Use(PROGRAM_TYPE_MODEL)
	shader.initialize(shader.ModelProgram, PROGRAM_TYPE_MODEL)
	shader.Unuse()

	return shader, nil
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

	// if programType == PROGRAM_TYPE_MODEL || programType == PROGRAM_TYPE_EDGE {
	// 	vertexShaderSource = fmt.Sprintf(
	// 		string(vertexShaderSource),
	// 		VS_LAYOUT_POSITION_ID,
	// 		VS_LAYOUT_NORMAL_ID,
	// 		VS_LAYOUT_UV_ID,
	// 		VS_LAYOUT_EXTEND_UV_ID,
	// 		VS_LAYOUT_EDGE_ID,
	// 		VS_LAYOUT_BONE_ID,
	// 		VS_LAYOUT_WEIGHT_ID,
	// 		VS_LAYOUT_MORPH_POS_ID,
	// 		VS_LAYOUT_MORPH_UV_ID,
	// 		VS_LAYOUT_MORPH_UV1_ID,
	// 		VS_LAYOUT_MORPH_AFTER_POS_ID,
	// 	)
	// }

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

	// # モデルビュー行列
	projection := mgl32.Perspective(mgl32.DegToRad(45.0), float32(s.width)/float32(s.height), 0.1, 10.0)
	projectionUniform := gl.GetUniformLocation(program, gl.Str("projection\x00"))
	gl.UniformMatrix4fv(projectionUniform, 1, false, &projection[0])

	camera := mgl32.LookAtV(mgl32.Vec3{3, 3, 3}, mgl32.Vec3{0, 0, 0}, mgl32.Vec3{0, 1, 0})
	cameraUniform := gl.GetUniformLocation(program, gl.Str("camera\x00"))
	gl.UniformMatrix4fv(cameraUniform, 1, false, &camera[0])

	model := mgl32.Ident4()
	modelUniform := gl.GetUniformLocation(program, gl.Str("model\x00"))
	gl.UniformMatrix4fv(modelUniform, 1, false, &model[0])
	s.ModelUniform = modelUniform

	textureUniform := gl.GetUniformLocation(program, gl.Str("tex\x00"))
	gl.Uniform1i(textureUniform, 0)

	gl.BindFragDataLocation(program, 0, gl.Str("outputColor\x00"))

	// モデルビュー行列
	s.ModelViewMatrixUniform[programType] = gl.GetUniformLocation(program, gl.Str("modelViewMatrix\x00"))
	// MVP行列
	s.ModelViewProjectionMatrixUniform[programType] = gl.GetUniformLocation(program, gl.Str("modelViewProjectionMatrix\x00"))

	// ボーン変形行列用テクスチャ
	boneMatrixTextureId := s.BoneMatrixTextureId[programType]
	gl.GenTextures(1, &boneMatrixTextureId)
	s.BoneMatrixTextureUniform[programType] = gl.GetUniformLocation(program, gl.Str("boneMatrixTexture\x00"))
	s.BoneMatrixTextureWidth[programType] = gl.GetUniformLocation(program, gl.Str("boneMatrixWidth\x00"))
	s.BoneMatrixTextureHeight[programType] = gl.GetUniformLocation(program, gl.Str("boneMatrixHeight\x00"))

	if programType == PROGRAM_TYPE_EDGE {
		// エッジシェーダーへの割り当て
		// エッジ設定
		s.EdgeColorUniform[programType] = gl.GetUniformLocation(program, gl.Str("edgeColor\x00"))
		s.EdgeAlphaUniform[programType] = gl.GetUniformLocation(program, gl.Str("edgeAlpha\x00"))
		s.EdgeSizeUniform[programType] = gl.GetUniformLocation(program, gl.Str("edgeSize\x00"))
	} else if programType == PROGRAM_TYPE_BONE {
		// 選択ボーン色
		s.SelectBoneColorUniform[programType] = gl.GetUniformLocation(program, gl.Str("selectBoneColor\x00"))
		// 非選択ボーン色
		s.UnselectBoneColorUniform[programType] = gl.GetUniformLocation(program, gl.Str("unselectBoneColor\x00"))
		// ボーン数
		s.BoneCountUniform[programType] = gl.GetUniformLocation(program, gl.Str("boneCount\x00"))
	} else {
		// モデルシェーダーへの割り当て
		s.LightDirectionUniform[programType] = gl.GetUniformLocation(program, gl.Str("lightDirection\x00"))
		gl.Uniform3f(s.LightDirectionUniform[programType],
			float32(s.lightDirection.GetX()), float32(s.lightDirection.GetY()), float32(s.lightDirection.GetZ()))

		// カメラの位置
		s.CameraVecUniform[programType] = gl.GetUniformLocation(program, gl.Str("cameraPos\x00"))

		// マテリアル設定
		s.DiffuseColorUniform[programType] = gl.GetUniformLocation(program, gl.Str("diffuseColor\x00"))
		s.DiffuseAlphaUniform[programType] = gl.GetUniformLocation(program, gl.Str("diffuseAlpha\x00"))
		s.AmbientColorUniform[programType] = gl.GetUniformLocation(program, gl.Str("ambientColor\x00"))
		s.SpecularColorUniform[programType] = gl.GetUniformLocation(program, gl.Str("specularColor\x00"))
		s.SpecularPowerUniform[programType] = gl.GetUniformLocation(program, gl.Str("specularPower\x00"))

		// テクスチャの設定
		s.UseTextureUniform[programType] = gl.GetUniformLocation(program, gl.Str("useTexture\x00"))
		s.TextureUniform[programType] = gl.GetUniformLocation(program, gl.Str("textureSampler\x00"))
		s.TextureFactorUniform[programType] = gl.GetUniformLocation(program, gl.Str("textureFactor\x00"))

		// Toonの設定
		s.UseToonUniform[programType] = gl.GetUniformLocation(program, gl.Str("useToon\x00"))
		s.ToonUniform[programType] = gl.GetUniformLocation(program, gl.Str("toonSampler\x00"))
		s.ToonFactorUniform[programType] = gl.GetUniformLocation(program, gl.Str("toonFactor\x00"))

		// Sphereの設定
		s.UseSphereUniform[programType] = gl.GetUniformLocation(program, gl.Str("useSphere\x00"))
		s.SphereModeUniform[programType] = gl.GetUniformLocation(program, gl.Str("sphereMode\x00"))
		s.SphereUniform[programType] = gl.GetUniformLocation(program, gl.Str("sphereSampler\x00"))
		s.SphereFactorUniform[programType] = gl.GetUniformLocation(program, gl.Str("sphereFactor\x00"))

		// ウェイトの描写
		s.IsShowBoneWeightUniform[programType] = gl.GetUniformLocation(program, gl.Str("isShowBoneWeight\x00"))
		s.ShowBoneIndexesUniform[programType] = gl.GetUniformLocation(program, gl.Str("showBoneIndexes\x00"))
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
	}
}

func (s *MShader) Unuse() {
	gl.UseProgram(0)
}

func (s *MShader) Delete() {
	s.DeleteProgram(s.ModelProgram)
	s.DeleteProgram(s.EdgeProgram)
	s.DeleteProgram(s.BoneProgram)
}
