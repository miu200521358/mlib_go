package mutils

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

// 指定されたパスがファイルとして存在しているか
func ExistsFile(path string) (bool, error) {
	if path == "" {
		return false, fmt.Errorf("path is empty")
	}

	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false, err
	}
	return !info.IsDir(), nil
}

// 指定されたパスをファイルとして開く
func Open(path string) (*os.File, error) {
	isFile, err := ExistsFile(path)
	if err != nil {
		return nil, err
	}
	if !isFile {
		return nil, fmt.Errorf("path not file: %s", path)
	}

	// ファイルを開く
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	return file, nil
}

// 指定されたファイルを閉じる
func Close(file *os.File) {
	defer file.Close()
}

// アプリのルートディレクトリを取得
func GetAppRootDir() string {
	// ファイルのフルパスを取得
	exePath, err := os.Executable()
	if err != nil {
		panic(err)
	}
	return filepath.Dir(exePath)
}

// テキストファイルの全文を読み込んでひとつの文字列で返す
func ReadText(path string) (string, error) {
	isExist, err := ExistsFile(path)
	if err != nil {
		return "", err
	}
	if !isExist {
		return "", fmt.Errorf("path not found: %s", path)
	}

	file, err := Open(path)
	if err != nil {
		return "", err
	}

	content, err := io.ReadAll(file)
	if err != nil {
		return "", err
	}

	defer Close(file)

	return string(content), nil
}

// 保存可能なファイルパスであるか否か
func CanSave(path string) bool {
	if path == "" {
		return false
	}

	// ファイルパスまでのディレクトリ
	dir := filepath.Dir(path)
	if dir == "" || dir == "." {
		return false
	}

	// ディレクトリが存在しない場合は作成不可
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return false
	}

	return true
}

func CreateOutputPath(originalPath, label string) string {
	originalDir, fileName := filepath.Split(originalPath)
	ext := filepath.Ext(fileName)
	if label != "" {
		label = label + "_"
	}
	return filepath.Join(originalDir, fmt.Sprintf("%s_%s%s%s", fileName[:len(fileName)-len(ext)],
		label, time.Now().Format("20060102_150405"), ext))
}
