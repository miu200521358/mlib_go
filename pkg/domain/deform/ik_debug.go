// 指示: miu200521358
package deform

import (
	"fmt"

	"github.com/miu200521358/mlib_go/pkg/domain/mmath"
	"github.com/miu200521358/mlib_go/pkg/domain/model"
	"github.com/miu200521358/mlib_go/pkg/domain/motion"
)

// IIkDebugFactory はIKデバッグ出力のセッション生成を行うI/F。
type IIkDebugFactory interface {
	// NewIkDebugSession はIKデバッグ用セッションを生成する。
	// 生成に失敗した場合はnilを返す。
	NewIkDebugSession(input IkDebugSessionInput) IIkDebugSession
}

// IIkDebugSession はIKデバッグ出力を行うI/F。
type IIkDebugSession interface {
	// AppendIkRotation はIKデバッグ用の回転フレームを追加する。
	AppendIkRotation(frameIndex int, boneName string, rotation mmath.Quaternion)
	// AppendGlobalPosition はIKデバッグ用のグローバル位置フレームを追加する。
	AppendGlobalPosition(frameIndex int, boneName string, position mmath.Vec3)
	// Logf はIKデバッグ用のログを出力する。
	Logf(frameIndex int, format string, params ...any)
	// Close はIKデバッグ出力を終了する。
	Close()
}

// IkDebugSessionInput はIKデバッグセッション生成の入力を表す。
type IkDebugSessionInput struct {
	ModelPath   string
	MotionPath  string
	Frame       motion.Frame
	IkBoneName  string
	IkBoneIndex int
	OrderIndex  int
}

// IkDebugContext はIK計算中のデバッグ情報を表す。
type IkDebugContext struct {
	Session      IIkDebugSession
	FrameIndex   int
	Frame        motion.Frame
	LinkBoneName string
}

// ikDebugPosition はIKデバッグ用の位置情報を表す。
type ikDebugPosition struct {
	Name     string
	Position mmath.Vec3
}

// newIkDebugContext はIKデバッグ用のコンテキストを生成する。
func newIkDebugContext(factory IIkDebugFactory, modelData *model.PmxModel, motionData *motion.VmdMotion, frame motion.Frame, ikBone *model.Bone, orderIndex int) *IkDebugContext {
	if factory == nil || modelData == nil || ikBone == nil {
		return nil
	}
	motionPath := ""
	if motionData != nil {
		motionPath = motionData.Path()
	}
	session := factory.NewIkDebugSession(IkDebugSessionInput{
		ModelPath:   modelData.Path(),
		MotionPath:  motionPath,
		Frame:       frame,
		IkBoneName:  ikBone.Name(),
		IkBoneIndex: ikBone.Index(),
		OrderIndex:  orderIndex,
	})
	if session == nil {
		return nil
	}
	return &IkDebugContext{
		Session:    session,
		FrameIndex: 1,
		Frame:      frame,
	}
}

// closeIkDebugContext はIKデバッグ用コンテキストをクローズする。
func closeIkDebugContext(ctx *IkDebugContext) {
	if ctx == nil || ctx.Session == nil {
		return
	}
	ctx.Session.Close()
}

// appendIkRotation はIKデバッグ用の回転フレームを追加する。
func appendIkRotation(ctx *IkDebugContext, boneName string, rotation mmath.Quaternion) int {
	if ctx == nil || ctx.Session == nil {
		return 0
	}
	frameIndex := ctx.FrameIndex
	ctx.Session.AppendIkRotation(frameIndex, boneName, rotation)
	ctx.FrameIndex++
	return frameIndex
}

// appendGlobalPositions はIKデバッグ用のグローバル位置フレームを追加する。
func appendGlobalPositions(ctx *IkDebugContext, positions []ikDebugPosition) int {
	if ctx == nil || ctx.Session == nil {
		return 0
	}
	frameIndex := ctx.FrameIndex
	for _, pos := range positions {
		ctx.Session.AppendGlobalPosition(frameIndex, pos.Name, pos.Position)
	}
	ctx.FrameIndex++
	return frameIndex
}

// ikDebugStepIndex はIKデバッグ用の直近ステップ番号を返す。
func ikDebugStepIndex(ctx *IkDebugContext) int {
	if ctx == nil || ctx.FrameIndex <= 0 {
		return 0
	}
	return ctx.FrameIndex - 1
}

// logIkDebugf はIKデバッグ用ログを出力する。
func logIkDebugf(ctx *IkDebugContext, format string, params ...any) {
	if ctx == nil || ctx.Session == nil {
		return
	}
	step := ikDebugStepIndex(ctx)
	prefix := fmt.Sprintf("frame=%v step=%d", ctx.Frame, step)
	if ctx.LinkBoneName != "" {
		prefix = fmt.Sprintf("%s link=%s", prefix, ctx.LinkBoneName)
	}
	message := fmt.Sprintf(format, params...)
	ctx.Session.Logf(step, "%s %s", prefix, message)
}
