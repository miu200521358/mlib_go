package mgl

import (
	"fmt"

	"github.com/go-gl/gl/v4.4-core/gl"

)

func checkGLError() error {
	if errCode := gl.GetError(); errCode != gl.NO_ERROR {
		return fmt.Errorf("OpenGL error: %v", errCode)
	}
	return nil
}
