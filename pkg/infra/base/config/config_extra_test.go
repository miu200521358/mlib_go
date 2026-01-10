// 指示: miu200521358
package config

import (
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/miu200521358/mlib_go/pkg/shared/base/config"
)

func withTempRoot(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	prev := appRootDirFn
	appRootDirFn = func() (string, error) { return dir, nil }
	t.Cleanup(func() { appRootDirFn = prev })
	return dir
}

func withLoadUserConfig(t *testing.T, fn func() (map[string]any, error)) {
	t.Helper()
	prev := loadUserConfigFn
	loadUserConfigFn = fn
	t.Cleanup(func() { loadUserConfigFn = prev })
}

func withWriteFile(t *testing.T, fn func(string, []byte, os.FileMode) error) {
	t.Helper()
	prev := writeFile
	writeFile = fn
	t.Cleanup(func() { writeFile = prev })
}

func withExecutable(t *testing.T, fn func() (string, error)) {
	t.Helper()
	prev := osExecutable
	osExecutable = fn
	t.Cleanup(func() { osExecutable = prev })
}

// TestConfigStoreAccessors はConfigStoreのgetterを確認する。
func TestConfigStoreAccessors(t *testing.T) {
	appCfg := &config.AppConfig{AppName: "app"}
	userCfg := &UserConfigStore{}
	store := NewConfigStore(appCfg, userCfg)
	if store.AppConfig() != appCfg {
		t.Errorf("AppConfig mismatch")
	}
	if store.UserConfig() != userCfg {
		t.Errorf("UserConfig mismatch")
	}

	var nilStore *ConfigStore
	if nilStore.AppConfig() != nil || nilStore.UserConfig() != nil {
		t.Errorf("nil ConfigStore should return nil")
	}
}

// TestNewUserConfigStoreSingleton はシングルトンを確認する。
func TestNewUserConfigStoreSingleton(t *testing.T) {
	first := NewUserConfigStore()
	second := NewUserConfigStore()
	if first != second {
		t.Errorf("NewUserConfigStore should return singleton")
	}
}

// TestUserConfigSetGet は各型のSet/Getを確認する。
func TestUserConfigSetGet(t *testing.T) {
	withTempRoot(t)
	store := &UserConfigStore{}

	if err := store.Set("name", "v1"); err != nil {
		t.Fatalf("Set string failed: %v", err)
	}
	if val, ok := store.Get("name"); !ok {
		t.Errorf("Get string: ok=%v", ok)
	} else if list, ok := val.([]interface{}); !ok || len(list) != 1 || list[0].(string) != "v1" {
		t.Errorf("Get string: got=%v", val)
	}

	if err := store.Set("flag", true); err != nil {
		t.Fatalf("Set bool failed: %v", err)
	}
	if !store.GetBool("flag", false) {
		t.Errorf("GetBool expected true")
	}
	if err := store.SetBool("flag2", false); err != nil {
		t.Fatalf("SetBool false failed: %v", err)
	}
	if store.GetBool("flag2", true) {
		t.Errorf("GetBool expected false")
	}

	if err := store.Set("num", 3); err != nil {
		t.Fatalf("Set int failed: %v", err)
	}
	if store.GetInt("num", 0) != 3 {
		t.Errorf("GetInt expected 3")
	}
	if err := store.SetInt("num2", 4); err != nil {
		t.Fatalf("SetInt failed: %v", err)
	}
	if root, err := store.AppRootDir(); err != nil || root == "" {
		t.Errorf("AppRootDir: root=%v err=%v", root, err)
	}

	if err := store.Set("list", []string{"a", "b"}); err != nil {
		t.Fatalf("Set slice failed: %v", err)
	}
	if got := store.GetStringSlice("list"); !reflect.DeepEqual(got, []string{"a", "b"}) {
		t.Errorf("GetStringSlice: got=%v", got)
	}
	if values, all := store.GetAll("list"); !reflect.DeepEqual(values, []string{"a", "b"}) || all["list"] == nil {
		t.Errorf("GetAll: values=%v all=%v", values, all)
	}

	if err := store.Set("empty", []string{}); err != nil {
		t.Fatalf("Set empty slice failed: %v", err)
	}
	if _, ok := store.Get("empty"); ok {
		t.Errorf("empty slice should not be stored")
	}

	if err := store.Set("bad", 1.23); err == nil {
		t.Errorf("Set unsupported type should error")
	}
}

// TestUserConfigLoadAllVariants はGetAllの分岐を確認する。
func TestUserConfigLoadAllVariants(t *testing.T) {
	withTempRoot(t)
	store := &UserConfigStore{}

	root := MustAppRootDir()
	path := filepath.Join(root, config.UserConfigFileName)
	if err := os.WriteFile(path, []byte(`{"history":["a","b"]}`), 0644); err != nil {
		t.Fatalf("seed config failed: %v", err)
	}
	values, _ := store.GetAll("history")
	if !reflect.DeepEqual(values, []string{"a", "b"}) {
		t.Errorf("GetAll list: got=%v", values)
	}
	if values, _ := store.GetAll("missing"); len(values) != 0 {
		t.Errorf("GetAll missing: got=%v", values)
	}

	if err := os.WriteFile(path, []byte(`{"history":[1]}`), 0644); err != nil {
		t.Fatalf("seed config failed: %v", err)
	}
	values, _ = store.GetAll("history")
	if len(values) != 0 {
		t.Errorf("GetAll non-string: got=%v", values)
	}

	if err := os.WriteFile(path, []byte(`{"history":"x"}`), 0644); err != nil {
		t.Fatalf("seed config failed: %v", err)
	}
	values, _ = store.GetAll("history")
	if len(values) != 0 {
		t.Errorf("GetAll default: got=%v", values)
	}

	withLoadUserConfig(t, func() (map[string]any, error) {
		return map[string]any{"history": []string{"x", "y"}}, nil
	})
	values, _ = store.GetAll("history")
	if !reflect.DeepEqual(values, []string{"x", "y"}) {
		t.Errorf("GetAll []string: got=%v", values)
	}
}

// TestUserConfigSaveStringSlice は重複排除と制限を確認する。
func TestUserConfigSaveStringSlice(t *testing.T) {
	withTempRoot(t)
	store := &UserConfigStore{}

	if err := store.Set("path", ""); err != nil {
		t.Fatalf("Set empty string failed: %v", err)
	}
	if err := store.SetStringSlice("path", []string{}, 1); err != nil {
		t.Fatalf("SetStringSlice empty failed: %v", err)
	}

	root := MustAppRootDir()
	path := filepath.Join(root, config.UserConfigFileName)
	if err := os.WriteFile(path, []byte(`{"path":["c"]}`), 0644); err != nil {
		t.Fatalf("seed config failed: %v", err)
	}

	if err := store.SetStringSlice("path", []string{"a", "", "b", "a"}, 2); err != nil {
		t.Fatalf("SetStringSlice failed: %v", err)
	}
	values := store.GetStringSlice("path")
	if !reflect.DeepEqual(values, []string{"a", "b"}) {
		t.Errorf("GetStringSlice: got=%v", values)
	}
}

// TestUserConfigBoolIntDefaults はデフォルト処理を確認する。
func TestUserConfigBoolIntDefaults(t *testing.T) {
	withTempRoot(t)
	store := &UserConfigStore{}
	if !store.GetBool("missing", true) {
		t.Errorf("Bool default failed")
	}
	if store.GetInt("missing", 5) != 5 {
		t.Errorf("Int default failed")
	}

	root := MustAppRootDir()
	path := filepath.Join(root, config.UserConfigFileName)
	if err := os.WriteFile(path, []byte(`{"num":["bad"]}`), 0644); err != nil {
		t.Fatalf("seed config failed: %v", err)
	}
	if store.GetInt("num", 7) != 7 {
		t.Errorf("Int parse default failed")
	}
}

// TestSaveStringSliceErrors は保存時のエラー分岐を確認する。
func TestSaveStringSliceErrors(t *testing.T) {
	withTempRoot(t)
	store := &UserConfigStore{}

	withLoadUserConfig(t, func() (map[string]any, error) {
		return map[string]any{"bad": make(chan int)}, nil
	})
	if err := store.SetStringSlice("list", []string{"a"}, 1); err == nil {
		t.Errorf("SetStringSlice marshal error expected")
	}

	withLoadUserConfig(t, func() (map[string]any, error) {
		return map[string]any{}, nil
	})
	withWriteFile(t, func(string, []byte, os.FileMode) error {
		return errors.New("write error")
	})
	if err := store.SetStringSlice("list", []string{"a"}, 1); err == nil {
		t.Errorf("SetStringSlice write error expected")
	}

	prev := appRootDirFn
	appRootDirFn = func() (string, error) {
		return "", errors.New("root error")
	}
	t.Cleanup(func() { appRootDirFn = prev })
	if err := store.SetStringSlice("list", []string{"a"}, 1); err == nil {
		t.Errorf("SetStringSlice AppRootDir error expected")
	}
}

// TestMustAppRootDirPanic はpanic分岐を確認する。
func TestMustAppRootDirPanic(t *testing.T) {
	withExecutable(t, func() (string, error) {
		return "", errors.New("exec error")
	})
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("MustAppRootDir should panic")
		}
	}()
	_ = MustAppRootDir()
}

// TestLoadUserConfigInvalidJSON はパース失敗を確認する。
func TestLoadUserConfigInvalidJSON(t *testing.T) {
	withTempRoot(t)
	root := MustAppRootDir()
	path := filepath.Join(root, config.UserConfigFileName)
	if err := os.WriteFile(path, []byte("{invalid"), 0644); err != nil {
		t.Fatalf("seed config failed: %v", err)
	}
	store := &UserConfigStore{}
	values := store.GetStringSlice("missing")
	if values == nil {
		t.Errorf("GetStringSlice should not be nil")
	}
}

// TestLoadUserConfigAppRootDirError はAppRootDir失敗を確認する。
func TestLoadUserConfigAppRootDirError(t *testing.T) {
	prev := appRootDirFn
	appRootDirFn = func() (string, error) {
		return "", errors.New("root error")
	}
	t.Cleanup(func() { appRootDirFn = prev })

	if _, err := loadUserConfig(); err == nil {
		t.Errorf("loadUserConfig expected error")
	}
}
