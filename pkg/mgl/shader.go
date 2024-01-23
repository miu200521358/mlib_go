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
	boneMatrixTextureUniform         map[ProgramType]interface{}
	boneMatrixTextureId              map[ProgramType]interface{}
	boneMatrixTextureWidth           map[ProgramType]interface{}
	boneMatrixTextureHeight          map[ProgramType]interface{}
	modelViewMatrixUniform           map[ProgramType]interface{}
	modelViewProjectionMatrixUniform map[ProgramType]interface{}
	lightDirectionUniform            map[ProgramType]interface{}
	cameraVecUniform                 map[ProgramType]interface{}
	diffuseUniform                   map[ProgramType]interface{}
	ambientUniform                   map[ProgramType]interface{}
	specularUniform                  map[ProgramType]interface{}
	edgeColorUniform                 map[ProgramType]interface{}
	selectBoneColorUniform           map[ProgramType]interface{}
	unselectBoneColorUniform         map[ProgramType]interface{}
	edgeSizeUniform                  map[ProgramType]interface{}
	useTextureUniform                map[ProgramType]interface{}
	textureUniform                   map[ProgramType]interface{}
	textureFactorUniform             map[ProgramType]interface{}
	useToonUniform                   map[ProgramType]interface{}
	toonUniform                      map[ProgramType]interface{}
	toonFactorUniform                map[ProgramType]interface{}
	useSphereUniform                 map[ProgramType]interface{}
	sphereModeUniform                map[ProgramType]interface{}
	sphereUniform                    map[ProgramType]interface{}
	sphereFactorUniform              map[ProgramType]interface{}
	boneCountUniform                 map[ProgramType]interface{}
	isShowBoneWeightUniform          map[ProgramType]interface{}
	showBoneIndexesUniform           map[ProgramType]interface{}
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
		boneMatrixTextureUniform:         make(map[ProgramType]interface{}),
		boneMatrixTextureId:              make(map[ProgramType]interface{}),
		boneMatrixTextureWidth:           make(map[ProgramType]interface{}),
		boneMatrixTextureHeight:          make(map[ProgramType]interface{}),
		modelViewMatrixUniform:           make(map[ProgramType]interface{}),
		modelViewProjectionMatrixUniform: make(map[ProgramType]interface{}),
		lightDirectionUniform:            make(map[ProgramType]interface{}),
		cameraVecUniform:                 make(map[ProgramType]interface{}),
		diffuseUniform:                   make(map[ProgramType]interface{}),
		ambientUniform:                   make(map[ProgramType]interface{}),
		specularUniform:                  make(map[ProgramType]interface{}),
		edgeColorUniform:                 make(map[ProgramType]interface{}),
		selectBoneColorUniform:           make(map[ProgramType]interface{}),
		unselectBoneColorUniform:         make(map[ProgramType]interface{}),
		edgeSizeUniform:                  make(map[ProgramType]interface{}),
		useTextureUniform:                make(map[ProgramType]interface{}),
		textureUniform:                   make(map[ProgramType]interface{}),
		textureFactorUniform:             make(map[ProgramType]interface{}),
		useToonUniform:                   make(map[ProgramType]interface{}),
		toonUniform:                      make(map[ProgramType]interface{}),
		toonFactorUniform:                make(map[ProgramType]interface{}),
		useSphereUniform:                 make(map[ProgramType]interface{}),
		sphereModeUniform:                make(map[ProgramType]interface{}),
		sphereUniform:                    make(map[ProgramType]interface{}),
		sphereFactorUniform:              make(map[ProgramType]interface{}),
		boneCountUniform:                 make(map[ProgramType]interface{}),
		isShowBoneWeightUniform:          make(map[ProgramType]interface{}),
		showBoneIndexesUniform:           make(map[ProgramType]interface{}),
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

	// shader.newProgram("edge.vert", "edge.frag", PROGRAM_TYPE_EDGE)
	// shader.use(PROGRAM_TYPE_EDGE)
	// shader.initialize(shader.edgeProgram, PROGRAM_TYPE_EDGE)
	// shader.unuse()

	// shader.newProgram("bone.vert", "bone.frag", PROGRAM_TYPE_BONE)
	// shader.use(PROGRAM_TYPE_BONE)
	// shader.initialize(shader.boneProgram, PROGRAM_TYPE_BONE)
	// shader.unuse()

	// shader.fit(
	// 	shader.width,
	// 	shader.height,
	// 	camera_position,
	// 	camera_offset_position,
	// 	camera_degrees,
	// 	look_at_center,
	// 	vertical_degrees,
	// 	aspect_ratio,
	// )

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
