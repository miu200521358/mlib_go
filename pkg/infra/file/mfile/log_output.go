// 指示: miu200521358
package mfile

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/miu200521358/mlib_go/pkg/shared/base/config"
	"github.com/miu200521358/mlib_go/pkg/shared/base/merr"
	"github.com/miu200521358/mlib_go/pkg/shared/base/logging"
)

type verboseFileSink struct {
	mu       sync.Mutex
	file     *os.File
	writeErr error
}

func (s *verboseFileSink) WriteLine(text string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.file == nil {
		s.writeErr = newLogStreamOpenFailed("ログストリームが初期化されていません", nil)
		return
	}
	if s.writeErr != nil {
		return
	}
	if _, err := s.file.WriteString(text + "\n"); err != nil {
		s.writeErr = newLogStreamOpenFailed("ログストリームへの書き込みに失敗しました", merr.NewOsPackageError("file.WriteStringに失敗しました", err))
	}
}

func (s *verboseFileSink) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	var err error
	if s.file != nil {
		if closeErr := s.file.Close(); closeErr != nil {
			err = newLogStreamOpenFailed("ログストリームのクローズに失敗しました", merr.NewOsPackageError("file.Closeに失敗しました", closeErr))
		}
		s.file = nil
	}
	if s.writeErr != nil {
		err = s.writeErr
		s.writeErr = nil
	}
	return err
}

// OpenVerboseLogStream は冗長ログのストリームを生成する。
func OpenVerboseLogStream(userConfig config.IUserConfig, label string) (string, logging.IVerboseSink, error) {
	if userConfig == nil {
		return "", nil, newLogStreamOpenFailed("ユーザー設定が初期化されていません", nil)
	}
	root, err := userConfig.AppRootDir()
	if err != nil {
		return "", nil, newLogStreamOpenFailed("アプリルートの取得に失敗しました", err)
	}
	dir := filepath.Join(root, "logs")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", nil, newLogStreamOpenFailed("ログディレクトリの作成に失敗しました", merr.NewOsPackageError("os.MkdirAllに失敗しました", err))
	}
	path := filepath.Join(dir, logFileName(label))
	file, err := os.Create(path)
	if err != nil {
		return "", nil, newLogStreamOpenFailed("ログファイルの作成に失敗しました", merr.NewOsPackageError("os.Createに失敗しました", err))
	}
	return path, &verboseFileSink{file: file}, nil
}

// SaveConsoleSnapshot はメッセージ欄の全文を保存する。
func SaveConsoleSnapshot(userConfig config.IUserConfig, label string, text string) (string, error) {
	if userConfig == nil {
		return "", newConsoleSnapshotSaveFailed("ユーザー設定が初期化されていません", nil)
	}
	root, err := userConfig.AppRootDir()
	if err != nil {
		return "", newConsoleSnapshotSaveFailed("アプリルートの取得に失敗しました", err)
	}
	dir := filepath.Join(root, "logs")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", newConsoleSnapshotSaveFailed("ログディレクトリの作成に失敗しました", merr.NewOsPackageError("os.MkdirAllに失敗しました", err))
	}
	path := filepath.Join(dir, logFileName(label))
	if err := os.WriteFile(path, []byte(text), 0o644); err != nil {
		return "", newConsoleSnapshotSaveFailed("ログスナップショットの保存に失敗しました", merr.NewOsPackageError("os.WriteFileに失敗しました", err))
	}
	return path, nil
}

func logFileName(label string) string {
	stamp := time.Now().Format("20060102_150405")
	if label == "" {
		return fmt.Sprintf("log_%s.txt", stamp)
	}
	return fmt.Sprintf("log_%s_%s.txt", label, stamp)
}
