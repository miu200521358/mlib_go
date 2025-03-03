//go:build windows
// +build windows

package mgl

import "github.com/miu200521358/mlib_go/pkg/domain/rendering"

// シェーダーユニフォーム定数
const (
	ShaderBoneMatrixTexture       = "boneMatrixTexture\x00"
	ShaderBoneMatrixTextureWidth  = "boneMatrixWidth\x00"
	ShaderBoneMatrixTextureHeight = "boneMatrixHeight\x00"
	ShaderViewMatrix              = "viewMatrix\x00"
	ShaderProjectionMatrix        = "projectionMatrix\x00"
	ShaderCameraPosition          = "cameraPosition\x00"
	ShaderLightDirection          = "lightDirection\x00"
	ShaderDiffuse                 = "diffuse\x00"
	ShaderAmbient                 = "ambient\x00"
	ShaderSpecular                = "specular\x00"
	ShaderEmissive                = "emissive\x00"
	ShaderLightDiffuse            = "lightDiffuse\x00"
	ShaderLightSpecular           = "lightSpecular\x00"
	ShaderLightAmbient            = "lightAmbient\x00"
	ShaderTextureSampler          = "textureSampler\x00"
	ShaderToonSampler             = "toonSampler\x00"
	ShaderSphereSampler           = "sphereSampler\x00"
	ShaderUseTexture              = "useTexture\x00"
	ShaderUseToon                 = "useToon\x00"
	ShaderUseSphere               = "useSphere\x00"
	ShaderSphereMode              = "sphereMode\x00"
	ShaderMorphTextureFactor      = "textureFactor\x00"
	ShaderMorphToonFactor         = "toonFactor\x00"
	ShaderMorphSphereFactor       = "sphereFactor\x00"
	ShaderColor                   = "color\x00"
	ShaderAlpha                   = "alpha\x00"
	ShaderEdgeColor               = "edgeColor\x00"
	ShaderEdgeSize                = "edgeSize\x00"
	ShaderCursorPositions         = "cursorPositions\x00"
	ShaderCursorThreshold         = "cursorThreshold\x00"
	ShaderLineWidth               = "lineWidth\x00"
)

// シェーダープログラム定義
type ShaderProgramConfig struct {
	VertexShader   string
	FragmentShader string
}

// 各プログラムタイプの設定
var ShaderProgramConfigs = map[rendering.ProgramType]ShaderProgramConfig{
	rendering.ProgramTypeModel:          {"glsl/model.vert", "glsl/model.frag"},
	rendering.ProgramTypeBone:           {"glsl/bone.vert", "glsl/bone.frag"},
	rendering.ProgramTypeEdge:           {"glsl/edge.vert", "glsl/edge.frag"},
	rendering.ProgramTypePhysics:        {"glsl/physics.vert", "glsl/physics.frag"},
	rendering.ProgramTypeNormal:         {"glsl/vertex.vert", "glsl/vertex.frag"},
	rendering.ProgramTypeFloor:          {"glsl/floor.vert", "glsl/floor.frag"},
	rendering.ProgramTypeWire:           {"glsl/vertex.vert", "glsl/vertex.frag"},
	rendering.ProgramTypeSelectedVertex: {"glsl/vertex.vert", "glsl/vertex.frag"},
	rendering.ProgramTypeOverride:       {"glsl/override.vert", "glsl/override.frag"},
	rendering.ProgramTypeCursor:         {"glsl/cursor.vert", "glsl/cursor.frag"},
}
