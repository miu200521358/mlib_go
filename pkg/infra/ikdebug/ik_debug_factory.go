// 指示: miu200521358
package ikdebug

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/miu200521358/mlib_go/pkg/adapter/io_common"
	"github.com/miu200521358/mlib_go/pkg/adapter/io_motion/vmd"
	"github.com/miu200521358/mlib_go/pkg/domain/deform"
	"github.com/miu200521358/mlib_go/pkg/domain/mmath"
	"github.com/miu200521358/mlib_go/pkg/domain/motion"
	"github.com/miu200521358/mlib_go/pkg/infra/file/mfile"
	"github.com/miu200521358/mlib_go/pkg/shared/base/logging"
)

// Factory はIKデバッグ出力の生成を担う。
type Factory struct {
	logger logging.ILogger
}

// NewFactory はIKデバッグ用ファクトリを生成する。
func NewFactory() *Factory {
	return NewFactoryWithLogger(logging.DefaultLogger())
}

// NewFactoryWithLogger はログ出力先を指定してIKデバッグ用ファクトリを生成する。
func NewFactoryWithLogger(logger logging.ILogger) *Factory {
	if logger == nil {
		logger = logging.DefaultLogger()
	}
	return &Factory{logger: logger}
}

// NewIkDebugSession はIKデバッグセッションを生成する。
func (f *Factory) NewIkDebugSession(input deform.IkDebugSessionInput) deform.IIkDebugSession {
	if f == nil {
		return nil
	}
	logger := f.logger
	if logger == nil {
		logger = logging.DefaultLogger()
	}
	if input.ModelPath == "" {
		if logger != nil {
			logger.Verbose(logging.VERBOSE_INDEX_IK, "IK冗長ログ: model.Pathが空のため出力をスキップ: bone=%s frame=%v", input.IkBoneName, input.Frame)
		}
		return nil
	}
	outputDir := filepath.Join(filepath.Dir(input.ModelPath), "IK_step")
	if err := os.MkdirAll(outputDir, 0o755); err != nil {
		if logger != nil {
			logger.Error("IK冗長ログ: 出力先作成に失敗しました: dir=%s err=%s", filepath.Base(outputDir), pathErrorText(err))
		}
		return nil
	}
	modelName := sanitizeLabel(nameFromPath(input.ModelPath))
	motionName := sanitizeLabel(nameFromPath(input.MotionPath))
	if motionName == "" {
		motionName = modelName
	}
	if motionName == "" {
		motionName = "motion"
	}
	boneName := sanitizeLabel(input.IkBoneName)
	if boneName == "" {
		boneName = "IK"
	}
	stamp := time.Now().Format("20060102_150405")
	prefix := fmt.Sprintf("%02d_%s_F%05d_%s", input.OrderIndex, motionName, int(input.Frame), stamp)
	logPath := filepath.Join(outputDir, fmt.Sprintf("%s_%s.log", prefix, boneName))
	ikMotionPath := filepath.Join(outputDir, fmt.Sprintf("%s_%s.vmd", prefix, boneName))
	globalMotionPath := filepath.Join(outputDir, fmt.Sprintf("%s_%s_global.vmd", prefix, boneName))
	logFile, err := os.OpenFile(logPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o644)
	if err != nil {
		if logger != nil {
			logger.Error("IK冗長ログ: ログファイルの作成に失敗しました: file=%s err=%s", filepath.Base(logPath), pathErrorText(err))
		}
		return nil
	}
	sess := &session{
		logger:           logger,
		frame:            input.Frame,
		ikBoneName:       input.IkBoneName,
		orderIndex:       input.OrderIndex,
		logFile:          logFile,
		logWriter:        bufio.NewWriter(logFile),
		ikMotion:         motion.NewVmdMotion(ikMotionPath),
		globalMotion:     motion.NewVmdMotion(globalMotionPath),
		logFileName:      filepath.Base(logPath),
		ikMotionFileName: filepath.Base(ikMotionPath),
		globalFileName:   filepath.Base(globalMotionPath),
		modelFileName:    baseNameOrEmpty(input.ModelPath),
		motionFileName:   baseNameOrEmpty(input.MotionPath),
	}
	sess.writeHeader(modelName, motionName)
	return sess
}

// session はIKデバッグ出力のセッションを表す。
type session struct {
	logger           logging.ILogger
	frame            motion.Frame
	ikBoneName       string
	orderIndex       int
	logFile          *os.File
	logWriter        *bufio.Writer
	ikMotion         *motion.VmdMotion
	globalMotion     *motion.VmdMotion
	logFileName      string
	ikMotionFileName string
	globalFileName   string
	modelFileName    string
	motionFileName   string
}

// AppendIkRotation はIKデバッグ用の回転フレームを追加する。
func (s *session) AppendIkRotation(frameIndex int, boneName string, rotation mmath.Quaternion) {
	if s == nil || s.ikMotion == nil {
		return
	}
	bf := motion.NewBoneFrame(motion.Frame(frameIndex))
	rot := rotation
	bf.Rotation = &rot
	s.ikMotion.BoneFrames.Get(boneName).Append(bf)
}

// AppendGlobalPosition はIKデバッグ用のグローバル位置フレームを追加する。
func (s *session) AppendGlobalPosition(frameIndex int, boneName string, position mmath.Vec3) {
	if s == nil || s.globalMotion == nil {
		return
	}
	bf := motion.NewBoneFrame(motion.Frame(frameIndex))
	pos := position
	bf.Position = &pos
	s.globalMotion.BoneFrames.Get(boneName).Append(bf)
}

// Logf はIKデバッグ用ログを出力する。
func (s *session) Logf(frameIndex int, format string, params ...any) {
	if s == nil || s.logWriter == nil {
		return
	}
	message := fmt.Sprintf(format, params...)
	line := fmt.Sprintf("[frame=%v][step=%05d] %s\n", s.frame, frameIndex, message)
	if _, err := s.logWriter.WriteString(line); err != nil && s.logger != nil {
		s.logger.Error("IK冗長ログ: 書き込みに失敗しました: %v", err)
	}
}

// Close はIKデバッグ出力を終了する。
func (s *session) Close() {
	if s == nil {
		return
	}
	if s.logWriter != nil {
		if err := s.logWriter.Flush(); err != nil && s.logger != nil {
			s.logger.Error("IK冗長ログ: フラッシュに失敗しました: %v", err)
		}
	}
	if s.logFile != nil {
		if err := s.logFile.Close(); err != nil && s.logger != nil {
			s.logger.Error("IK冗長ログ: ファイルクローズに失敗しました: %v", err)
		}
	}
	repo := vmd.NewVmdRepository()
	if s.ikMotion != nil {
		if err := repo.Save("", s.ikMotion, io_common.SaveOptions{}); err != nil && s.logger != nil {
			s.logger.Error("IK冗長ログ: IKモーション保存に失敗しました: file=%s err=%v", s.ikMotionFileName, err)
		}
	}
	if s.globalMotion != nil {
		if err := repo.Save("", s.globalMotion, io_common.SaveOptions{}); err != nil && s.logger != nil {
			s.logger.Error("IK冗長ログ: グローバルモーション保存に失敗しました: file=%s err=%v", s.globalFileName, err)
		}
	}
}

// writeHeader はログファイルのヘッダを出力する。
func (s *session) writeHeader(modelName string, motionName string) {
	if s == nil || s.logWriter == nil {
		return
	}
	lines := []string{
		"----------------------------------------",
		fmt.Sprintf("IKデバッグログ: bone=%s frame=%v", s.ikBoneName, s.frame),
		fmt.Sprintf("model=%s motion=%s", modelName, motionName),
		fmt.Sprintf("orderIndex=%02d", s.orderIndex),
	}
	if s.modelFileName != "" {
		lines = append(lines, fmt.Sprintf("modelFile=%s", s.modelFileName))
	}
	if s.motionFileName != "" {
		lines = append(lines, fmt.Sprintf("motionFile=%s", s.motionFileName))
	}
	lines = append(lines,
		fmt.Sprintf("ikMotion=%s", s.ikMotionFileName),
		fmt.Sprintf("globalMotion=%s", s.globalFileName),
		fmt.Sprintf("log=%s", s.logFileName),
		"----------------------------------------",
	)
	for _, line := range lines {
		if _, err := s.logWriter.WriteString(line + "\n"); err != nil && s.logger != nil {
			s.logger.Error("IK冗長ログ: ヘッダ出力に失敗しました: %v", err)
			return
		}
	}
}

// nameFromPath はパスから拡張子を除いた名前を返す。
func nameFromPath(path string) string {
	_, name, _ := mfile.SplitPath(path)
	return name
}

// baseNameOrEmpty はパスの末尾要素を返す。
func baseNameOrEmpty(path string) string {
	if path == "" {
		return ""
	}
	return filepath.Base(path)
}

// sanitizeLabel はファイル名に使えない文字を置換する。
func sanitizeLabel(value string) string {
	label := strings.TrimSpace(value)
	replacer := strings.NewReplacer("/", "_", "\\", "_", ":", "_", "\t", "_")
	label = replacer.Replace(label)
	label = strings.ReplaceAll(label, " ", "_")
	return label
}

// pathErrorText はパス情報を含まないエラー文字列を返す。
func pathErrorText(err error) string {
	if err == nil {
		return ""
	}
	if pathErr, ok := err.(*os.PathError); ok {
		if pathErr.Err != nil {
			return pathErr.Err.Error()
		}
	}
	return err.Error()
}
