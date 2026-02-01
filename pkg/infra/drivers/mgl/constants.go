//go:build windows
// +build windows

// 指示: miu200521358
package mgl

import (
	"github.com/go-gl/mathgl/mgl32"

	"github.com/miu200521358/mlib_go/pkg/adapter/graphics_api"
	"github.com/miu200521358/mlib_go/pkg/domain/mmath"
)

// OpenGL最小バージョン。
const (
	OpenGLMinMajor = 4
	OpenGLMinMinor = 3
	ShaderVersion  = "430 core"
)

// シェーダーユニフォーム定数。
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

// ShaderProgramConfig はプログラム構成を表す。
type ShaderProgramConfig struct {
	VertexShader   string
	FragmentShader string
}

// SHADER_PROGRAM_CONFIGS はプログラム種別ごとのシェーダ設定。
var SHADER_PROGRAM_CONFIGS = map[graphics_api.ProgramType]ShaderProgramConfig{
	graphics_api.ProgramTypeModel:          {"glsl/model.vert", "glsl/model.frag"},
	graphics_api.ProgramTypeBone:           {"glsl/bone.vert", "glsl/bone.frag"},
	graphics_api.ProgramTypeEdge:           {"glsl/edge.vert", "glsl/edge.frag"},
	graphics_api.ProgramTypePhysics:        {"glsl/physics.vert", "glsl/physics.frag"},
	graphics_api.ProgramTypeNormal:         {"glsl/vertex.vert", "glsl/vertex.frag"},
	graphics_api.ProgramTypeFloor:          {"glsl/floor.vert", "glsl/floor.frag"},
	graphics_api.ProgramTypeWire:           {"glsl/vertex.vert", "glsl/vertex.frag"},
	graphics_api.ProgramTypeSelectedVertex: {"glsl/vertex.vert", "glsl/vertex.frag"},
	graphics_api.ProgramTypeOverride:       {"glsl/override.vert", "glsl/override.frag"},
	graphics_api.ProgramTypeCursor:         {"glsl/cursor.vert", "glsl/cursor.frag"},
}

// NewGlVec3 はOpenGL座標系へ変換したVec3を返す。
func NewGlVec3(v *mmath.Vec3) mgl32.Vec3 {
	if v == nil {
		return mgl32.Vec3{}
	}
	return mgl32.Vec3{float32(-v.X), float32(v.Y), float32(v.Z)}
}

// NewGlMat4 はOpenGL座標系へ変換したMat4を返す。
func NewGlMat4(m mmath.Mat4) mgl32.Mat4 {
	return mgl32.Mat4{
		float32(m[0]), float32(-m[1]), float32(-m[2]), float32(m[3]),
		float32(-m[4]), float32(m[5]), float32(m[6]), float32(m[7]),
		float32(-m[8]), float32(m[9]), float32(m[10]), float32(m[11]),
		float32(-m[12]), float32(m[13]), float32(m[14]), float32(m[15]),
	}
}
