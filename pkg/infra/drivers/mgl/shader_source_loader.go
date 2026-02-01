//go:build windows
// +build windows

// 指示: miu200521358
package mgl

import (
	"embed"
	"io/fs"
	"strings"

	"github.com/go-gl/gl/v4.3-core/gl"
	"github.com/miu200521358/mlib_go/pkg/adapter/graphics_api"
)

//go:embed glsl/*
var glslFiles embed.FS

// ShaderSourceLoader はシェーダーソースの読み込みとコンパイルを担当する。
type ShaderSourceLoader struct{}

// NewShaderSourceLoader はShaderSourceLoaderを生成する。
func NewShaderSourceLoader() *ShaderSourceLoader {
	return &ShaderSourceLoader{}
}

// LoadSource はシェーダーソースを読み込む。
func (l *ShaderSourceLoader) LoadSource(path string) (string, error) {
	bytes, err := fs.ReadFile(glslFiles, path)
	if err != nil {
		return "", graphics_api.NewShaderSourceLoadFailed("シェーダーソースの読み込みに失敗しました: %s", err, path)
	}
	return string(bytes), nil
}

// Compile はシェーダーをコンパイルする。
func (l *ShaderSourceLoader) Compile(name, source string, shaderType uint32) (uint32, error) {
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
		gl.DeleteShader(shader)

		return 0, graphics_api.NewShaderCompileFailed(
			"シェーダーコンパイルに失敗しました: %s (%s)",
			nil,
			name,
			strings.TrimRight(log, "\x00"),
		)
	}

	return shader, nil
}

// CreateProgram はシェーダープログラムを作成する。
func (l *ShaderSourceLoader) CreateProgram(vertexPath, fragmentPath string) (uint32, error) {
	vertexSource, err := l.LoadSource(vertexPath)
	if err != nil {
		return 0, err
	}
	vertexShader, err := l.Compile(vertexPath, vertexSource, gl.VERTEX_SHADER)
	if err != nil {
		return 0, err
	}
	defer gl.DeleteShader(vertexShader)

	fragmentSource, err := l.LoadSource(fragmentPath)
	if err != nil {
		return 0, err
	}
	fragmentShader, err := l.Compile(fragmentPath, fragmentSource, gl.FRAGMENT_SHADER)
	if err != nil {
		return 0, err
	}
	defer gl.DeleteShader(fragmentShader)

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
		gl.DeleteProgram(program)

		return 0, graphics_api.NewShaderLinkFailed(
			"シェーダーリンクに失敗しました: %s",
			nil,
			strings.TrimRight(log, "\x00"),
		)
	}

	return program, nil
}
