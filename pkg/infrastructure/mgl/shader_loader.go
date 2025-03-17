//go:build windows
// +build windows

package mgl

import (
	"embed"
	"fmt"
	"io/fs"
	"strings"

	"github.com/go-gl/gl/v4.4-core/gl"
)

//go:embed glsl/*
var glslFiles embed.FS

// ShaderLoader はシェーダーファイルのロードと管理を担当
type ShaderLoader struct{}

// NewShaderLoader は新しいShaderLoaderを作成
func NewShaderLoader() *ShaderLoader {
	return &ShaderLoader{}
}

// LoadShaderSource はシェーダーソースを読み込む
func (l *ShaderLoader) LoadShaderSource(filename string) (string, error) {
	bytes, err := fs.ReadFile(glslFiles, filename)
	if err != nil {
		return "", fmt.Errorf("failed to load shader %s: %v", filename, err)
	}
	return string(bytes), nil
}

// CompileShader はシェーダーをコンパイルする
func (l *ShaderLoader) CompileShader(shaderName, source string, shaderType uint32) (uint32, error) {
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

// CreateProgram はシェーダープログラムを作成する
func (l *ShaderLoader) CreateProgram(vertexShaderPath, fragmentShaderPath string) (uint32, error) {
	// バーテックスシェーダーの読み込みとコンパイル
	vertexSource, err := l.LoadShaderSource(vertexShaderPath)
	if err != nil {
		return 0, err
	}
	vertexShader, err := l.CompileShader(vertexShaderPath, vertexSource, gl.VERTEX_SHADER)
	if err != nil {
		return 0, err
	}
	defer gl.DeleteShader(vertexShader)

	// フラグメントシェーダーの読み込みとコンパイル
	fragmentSource, err := l.LoadShaderSource(fragmentShaderPath)
	if err != nil {
		return 0, err
	}
	fragmentShader, err := l.CompileShader(fragmentShaderPath, fragmentSource, gl.FRAGMENT_SHADER)
	if err != nil {
		return 0, err
	}
	defer gl.DeleteShader(fragmentShader)

	// プログラムの作成とリンク
	program := gl.CreateProgram()
	gl.AttachShader(program, vertexShader)
	gl.AttachShader(program, fragmentShader)

	gl.LinkProgram(program)

	// リンク状態の確認
	var status int32
	gl.GetProgramiv(program, gl.LINK_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetProgramiv(program, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetProgramInfoLog(program, logLength, nil, gl.Str(log))

		return 0, fmt.Errorf("failed to link program: %v", log)
	}

	return program, nil
}
