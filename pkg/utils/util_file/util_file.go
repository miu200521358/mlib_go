package util_file

import (
	"fmt"
	"os"
	"path/filepath"

)

// 指定されたパスがファイルとして存在しているか
func ExistsFile(path string) (bool, error) {
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
