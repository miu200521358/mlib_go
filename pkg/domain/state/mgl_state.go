//go:build windows
// +build windows

package state

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
)

type IShader interface {
	GetProgram(programType ProgramType) uint32
	BoneTextureId() uint32
}
