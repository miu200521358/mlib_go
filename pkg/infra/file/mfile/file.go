// 指示: miu200521358
package mfile

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	baseerr "github.com/miu200521358/mlib_go/pkg/shared/base/err"
)

// ExistsFile はファイルが存在するか判定する。
func ExistsFile(path string) (bool, error) {
	if path == "" {
		return false, newFileReadFailed("パスが空です", nil)
	}
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, newFileReadFailed("ファイル確認に失敗しました", baseerr.NewOsPackageError("os.Statに失敗しました", err))
	}
	return info != nil && !info.IsDir(), nil
}

// ReadText はファイル全体を文字列で返す。
func ReadText(path string) (string, error) {
	if path == "" {
		return "", newFileNotFound("パスが空です", nil)
	}
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return "", newFileNotFound("ファイルが存在しません: "+path, nil)
		}
		return "", newFileReadFailed("ファイル確認に失敗しました", baseerr.NewOsPackageError("os.Statに失敗しました", err))
	}
	if info.IsDir() {
		return "", newFileNotFound("ファイルが存在しません: "+path, nil)
	}
	file, err := os.Open(path)
	if err != nil {
		return "", newFileReadFailed("ファイルを開けません: "+path, baseerr.NewOsPackageError("os.Openに失敗しました", err))
	}
	defer file.Close()
	content, err := io.ReadAll(file)
	if err != nil {
		return "", newFileReadFailed("ファイル読込に失敗しました: "+path, baseerr.NewOsPackageError("io.ReadAllに失敗しました", err))
	}
	return string(content), nil
}

// CanSave は保存可能なパスか判定する。
func CanSave(path string) bool {
	if path == "" {
		return false
	}
	dir := filepath.Dir(path)
	if dir == "" || dir == "." {
		return false
	}
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return false
	}
	return true
}

// CreateOutputPath はラベル付きの出力パスを生成する。
func CreateOutputPath(originalPath, label string) string {
	dir, name, ext := SplitPath(originalPath)
	stamp := time.Now().Format("20060102_150405")
	if name == "" {
		if label == "" {
			return filepath.Join(dir, fmt.Sprintf("%s%s", stamp, ext))
		}
		return filepath.Join(dir, fmt.Sprintf("%s_%s%s", label, stamp, ext))
	}
	if label == "" {
		return filepath.Join(dir, fmt.Sprintf("%s_%s%s", name, stamp, ext))
	}
	return filepath.Join(dir, fmt.Sprintf("%s_%s_%s%s", name, label, stamp, ext))
}

// SplitPath はパスを dir/name/ext に分割する。
func SplitPath(path string) (dir, name, ext string) {
	if path == "" {
		return "", "", ""
	}
	dir, base := filepath.Split(path)
	ext = filepath.Ext(base)
	if ext == "" {
		return dir, base, ""
	}
	return dir, base[:len(base)-len(ext)], ext
}
