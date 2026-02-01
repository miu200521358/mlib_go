// 指示: miu200521358
package config

import (
	"bytes"
	"embed"
	"errors"
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
	"reflect"
	"testing"
	"testing/fstest"

	"github.com/miu200521358/mlib_go/pkg/shared/base/config"
	"github.com/miu200521358/mlib_go/pkg/shared/base/merr"
)

func removeFile(t *testing.T, path string) {
	t.Helper()
	if err := os.Remove(path); err != nil && !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("remove file failed: %v", err)
	}
}

func assertCommonErrorID(t *testing.T, err error, id string) {
	t.Helper()
	ce, ok := err.(*merr.CommonError)
	if !ok {
		t.Fatalf("error type mismatch: %T", err)
	}
	if ce.ErrorID() != id {
		t.Fatalf("ErrorID mismatch: got=%v want=%v", ce.ErrorID(), id)
	}
}

// TestLoadAppConfigSuccess は設定読み込みの成功を確認する。
func TestLoadAppConfigSuccess(t *testing.T) {
	var buf bytes.Buffer
	img := image.NewRGBA(image.Rect(0, 0, 1, 1))
	img.Set(0, 0, color.RGBA{R: 0, G: 0, B: 0, A: 255})
	if err := png.Encode(&buf, img); err != nil {
		t.Fatalf("encode png failed: %v", err)
	}
	fsys := fstest.MapFS{
		"app/app_config.json": &fstest.MapFile{Data: []byte(`{"Name":"TestApp","Version":"1","Horizontal":true,"CloseConfirm":true}`)},
		"app/app.png":         &fstest.MapFile{Data: buf.Bytes()},
	}
	cfg, err := loadAppConfigFS(fsys)
	if err != nil {
		t.Fatalf("loadAppConfigFS failed: %v", err)
	}
	if cfg.AppName != "TestApp" {
		t.Errorf("AppName: got=%v", cfg.AppName)
	}
	if cfg.IconImagePath != config.AppIconImagePath {
		t.Errorf("IconImagePath: got=%v", cfg.IconImagePath)
	}
}

// TestLoadAppConfigError は読み込み失敗時のerrorを確認する。
func TestLoadAppConfigError(t *testing.T) {
	fsys := fstest.MapFS{
		"app/app_config.json": &fstest.MapFile{Data: []byte("invalid-json")},
		"app/app.png":         &fstest.MapFile{Data: []byte("not-png")},
	}
	if _, err := loadAppConfigFS(fsys); err == nil {
		t.Errorf("loadAppConfigFS expected error")
	} else {
		assertCommonErrorID(t, err, appConfigLoadFailedErrorID)
	}
}

// TestLoadAppConfigMissingFile は設定ファイル欠落を確認する。
func TestLoadAppConfigMissingFile(t *testing.T) {
	fsys := fstest.MapFS{}
	if _, err := loadAppConfigFS(fsys); err == nil {
		t.Errorf("loadAppConfigFS expected read error")
	} else {
		assertCommonErrorID(t, err, appConfigLoadFailedErrorID)
	}
}

// TestLoadAppConfigIconMissing はアイコン欠落を確認する。
func TestLoadAppConfigIconMissing(t *testing.T) {
	fsys := fstest.MapFS{
		"app/app_config.json": &fstest.MapFile{Data: []byte(`{"Name":"TestApp","IconImagePath":"app/missing.png"}`)},
	}
	if _, err := loadAppConfigFS(fsys); err == nil {
		t.Errorf("loadAppConfigFS expected icon read error")
	} else {
		assertCommonErrorID(t, err, appConfigLoadFailedErrorID)
	}
}

// TestLoadAppConfigIconDecode はデコード失敗を確認する。
func TestLoadAppConfigIconDecode(t *testing.T) {
	fsys := fstest.MapFS{
		"app/app_config.json": &fstest.MapFile{Data: []byte(`{"Name":"TestApp","IconImagePath":"app/app.png"}`)},
		"app/app.png":         &fstest.MapFile{Data: []byte("not-png")},
	}
	if _, err := loadAppConfigFS(fsys); err == nil {
		t.Errorf("loadAppConfigFS expected decode error")
	} else {
		assertCommonErrorID(t, err, appConfigLoadFailedErrorID)
	}
}

// TestLoadAppConfig は関数経由の読み込みを確認する。
func TestLoadAppConfig(t *testing.T) {
	var emptyFS embed.FS
	if _, err := LoadAppConfig(emptyFS); err == nil {
		t.Errorf("LoadAppConfig expected error")
	} else {
		assertCommonErrorID(t, err, appConfigLoadFailedErrorID)
	}
}

// TestUserConfigSetStringSlice は保存と重複排除を確認する。
func TestUserConfigSetStringSlice(t *testing.T) {
	store := &UserConfigStore{}
	root, err := AppRootDir()
	if err != nil {
		t.Fatalf("AppRootDir failed: %v", err)
	}
	userPath := filepath.Join(root, config.UserConfigFileName)
	removeFile(t, userPath)
	t.Cleanup(func() { removeFile(t, userPath) })

	if err := os.WriteFile(userPath, []byte(`{"history":["path1","path2"]}`), 0644); err != nil {
		t.Fatalf("seed user config failed: %v", err)
	}

	if err := store.SetStringSlice("history", []string{"path2"}, 3); err != nil {
		t.Fatalf("SetStringSlice failed: %v", err)
	}
	values, err := store.GetStringSlice("history")
	if err != nil {
		t.Fatalf("GetStringSlice failed: %v", err)
	}
	expected := []string{"path2", "path1"}
	if !reflect.DeepEqual(values, expected) {
		t.Errorf("GetStringSlice: got=%v want=%v", values, expected)
	}
}

// TestUserConfigFallback はhistory.jsonのフォールバックを確認する。
func TestUserConfigFallback(t *testing.T) {
	store := &UserConfigStore{}
	root, err := AppRootDir()
	if err != nil {
		t.Fatalf("AppRootDir failed: %v", err)
	}
	userPath := filepath.Join(root, config.UserConfigFileName)
	legacyPath := filepath.Join(root, config.UserConfigLegacyFileName)
	removeFile(t, userPath)
	removeFile(t, legacyPath)
	t.Cleanup(func() { removeFile(t, userPath) })
	t.Cleanup(func() { removeFile(t, legacyPath) })

	if err := os.WriteFile(legacyPath, []byte(`{"lang":["ja"]}`), 0644); err != nil {
		t.Fatalf("seed legacy config failed: %v", err)
	}

	values, err := store.GetStringSlice("lang")
	if err != nil {
		t.Fatalf("GetStringSlice failed: %v", err)
	}
	expected := []string{"ja"}
	if !reflect.DeepEqual(values, expected) {
		t.Errorf("GetStringSlice: got=%v want=%v", values, expected)
	}
}
