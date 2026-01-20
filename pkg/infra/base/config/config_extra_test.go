// 指示: miu200521358
package config

import (
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/miu200521358/mlib_go/pkg/shared/base/config"
	"github.com/miu200521358/mlib_go/pkg/shared/base/merr"
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

func withReadFile(t *testing.T, fn func(string) ([]byte, error)) {
	t.Helper()
	prev := readFile
	readFile = fn
	t.Cleanup(func() { readFile = prev })
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
	if val, ok, err := store.Get("name"); err != nil {
		t.Fatalf("Get string failed: %v", err)
	} else if !ok {
		t.Errorf("Get string: ok=%v", ok)
	} else if list, ok := val.([]interface{}); !ok || len(list) != 1 || list[0].(string) != "v1" {
		t.Errorf("Get string: got=%v", val)
	}

	if err := store.Set("flag", true); err != nil {
		t.Fatalf("Set bool failed: %v", err)
	}
	if got, err := store.GetBool("flag", false); err != nil {
		t.Fatalf("GetBool failed: %v", err)
	} else if !got {
		t.Errorf("GetBool expected true")
	}
	if err := store.SetBool("flag2", false); err != nil {
		t.Fatalf("SetBool false failed: %v", err)
	}
	if got, err := store.GetBool("flag2", true); err != nil {
		t.Fatalf("GetBool failed: %v", err)
	} else if got {
		t.Errorf("GetBool expected false")
	}

	if err := store.Set("num", 3); err != nil {
		t.Fatalf("Set int failed: %v", err)
	}
	if got, err := store.GetInt("num", 0); err != nil {
		t.Fatalf("GetInt failed: %v", err)
	} else if got != 3 {
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
	if got, err := store.GetStringSlice("list"); err != nil {
		t.Fatalf("GetStringSlice failed: %v", err)
	} else if !reflect.DeepEqual(got, []string{"a", "b"}) {
		t.Errorf("GetStringSlice: got=%v", got)
	}
	if values, all, err := store.GetAll("list"); err != nil {
		t.Fatalf("GetAll failed: %v", err)
	} else if !reflect.DeepEqual(values, []string{"a", "b"}) || all["list"] == nil {
		t.Errorf("GetAll: values=%v all=%v", values, all)
	}

	if err := store.Set("empty", []string{}); err != nil {
		t.Fatalf("Set empty slice failed: %v", err)
	}
	if _, ok, err := store.Get("empty"); err != nil {
		t.Fatalf("Get empty failed: %v", err)
	} else if ok {
		t.Errorf("empty slice should not be stored")
	}

	if err := store.Set("bad", 1.23); err == nil {
		t.Errorf("Set unsupported type should error")
	} else if ce, ok := err.(*merr.CommonError); !ok || ce.ErrorID() != configValueTypeNotSupportedErrorID {
		t.Errorf("Set unsupported type ErrorID: err=%v", err)
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
	values, _, err := store.GetAll("history")
	if err != nil {
		t.Fatalf("GetAll failed: %v", err)
	}
	if !reflect.DeepEqual(values, []string{"a", "b"}) {
		t.Errorf("GetAll list: got=%v", values)
	}
	if values, _, err := store.GetAll("missing"); err != nil {
		t.Fatalf("GetAll failed: %v", err)
	} else if len(values) != 0 {
		t.Errorf("GetAll missing: got=%v", values)
	}

	if err := os.WriteFile(path, []byte(`{"history":[1]}`), 0644); err != nil {
		t.Fatalf("seed config failed: %v", err)
	}
	_, _, err = store.GetAll("history")
	if err == nil {
		t.Fatalf("GetAll non-string expected error")
	} else if ce, ok := err.(*merr.CommonError); !ok || ce.ErrorID() != configValueTypeNotSupportedErrorID {
		t.Fatalf("GetAll non-string ErrorID: err=%v", err)
	}

	if err := os.WriteFile(path, []byte(`{"history":"x"}`), 0644); err != nil {
		t.Fatalf("seed config failed: %v", err)
	}
	_, _, err = store.GetAll("history")
	if err == nil {
		t.Fatalf("GetAll non-slice expected error")
	} else if ce, ok := err.(*merr.CommonError); !ok || ce.ErrorID() != configValueTypeNotSupportedErrorID {
		t.Fatalf("GetAll non-slice ErrorID: err=%v", err)
	}

	withLoadUserConfig(t, func() (map[string]any, error) {
		return map[string]any{"history": []string{"x", "y"}}, nil
	})
	values, _, err = store.GetAll("history")
	if err != nil {
		t.Fatalf("GetAll failed: %v", err)
	}
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
	values, err := store.GetStringSlice("path")
	if err != nil {
		t.Fatalf("GetStringSlice failed: %v", err)
	}
	if !reflect.DeepEqual(values, []string{"a", "b"}) {
		t.Errorf("GetStringSlice: got=%v", values)
	}
}

// TestUserConfigBoolIntDefaults はデフォルト処理を確認する。
func TestUserConfigBoolIntDefaults(t *testing.T) {
	withTempRoot(t)
	store := &UserConfigStore{}
	if got, err := store.GetBool("missing", true); err != nil {
		t.Fatalf("GetBool failed: %v", err)
	} else if !got {
		t.Errorf("Bool default failed")
	}
	if got, err := store.GetInt("missing", 5); err != nil {
		t.Fatalf("GetInt failed: %v", err)
	} else if got != 5 {
		t.Errorf("Int default failed")
	}

	root := MustAppRootDir()
	path := filepath.Join(root, config.UserConfigFileName)
	if err := os.WriteFile(path, []byte(`{"num":["bad"]}`), 0644); err != nil {
		t.Fatalf("seed config failed: %v", err)
	}
	if got, err := store.GetInt("num", 7); err != nil {
		t.Fatalf("GetInt failed: %v", err)
	} else if got != 7 {
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
	} else if ce, ok := err.(*merr.CommonError); !ok || ce.ErrorID() != userConfigSaveFailedErrorID {
		t.Errorf("SetStringSlice marshal ErrorID: err=%v", err)
	}

	withLoadUserConfig(t, func() (map[string]any, error) {
		return map[string]any{}, nil
	})
	withWriteFile(t, func(string, []byte, os.FileMode) error {
		return errors.New("write error")
	})
	if err := store.SetStringSlice("list", []string{"a"}, 1); err == nil {
		t.Errorf("SetStringSlice write error expected")
	} else if ce, ok := err.(*merr.CommonError); !ok || ce.ErrorID() != userConfigSaveFailedErrorID {
		t.Errorf("SetStringSlice write ErrorID: err=%v", err)
	}

	prev := appRootDirFn
	appRootDirFn = func() (string, error) {
		return "", errors.New("root error")
	}
	t.Cleanup(func() { appRootDirFn = prev })
	if err := store.SetStringSlice("list", []string{"a"}, 1); err == nil {
		t.Errorf("SetStringSlice AppRootDir error expected")
	} else if ce, ok := err.(*merr.CommonError); !ok || ce.ErrorID() != userConfigSaveFailedErrorID {
		t.Errorf("SetStringSlice AppRootDir ErrorID: err=%v", err)
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
	if _, err := store.GetStringSlice("missing"); err == nil {
		t.Fatalf("GetStringSlice expected error")
	} else if ce, ok := err.(*merr.CommonError); !ok || ce.ErrorID() != merr.JsonPackageErrorID {
		t.Fatalf("GetStringSlice ErrorID: err=%v", err)
	}
}

// TestLoadUserConfigInvalidJSONError はパース失敗時のエラーIDを確認する。
func TestLoadUserConfigInvalidJSONError(t *testing.T) {
	withTempRoot(t)
	root := MustAppRootDir()
	path := filepath.Join(root, config.UserConfigFileName)
	if err := os.WriteFile(path, []byte("{invalid"), 0644); err != nil {
		t.Fatalf("seed config failed: %v", err)
	}
	if _, err := loadUserConfig(); err == nil {
		t.Fatalf("loadUserConfig expected error")
	} else if ce, ok := err.(*merr.CommonError); !ok || ce.ErrorID() != merr.JsonPackageErrorID {
		t.Fatalf("loadUserConfig ErrorID: err=%v", err)
	}
}

// TestLoadUserConfigReadError は読込失敗時のエラーIDを確認する。
func TestLoadUserConfigReadError(t *testing.T) {
	withTempRoot(t)
	withReadFile(t, func(string) ([]byte, error) {
		return nil, errors.New("read error")
	})
	if _, err := loadUserConfig(); err == nil {
		t.Fatalf("loadUserConfig expected error")
	} else if ce, ok := err.(*merr.CommonError); !ok || ce.ErrorID() != merr.OsPackageErrorID {
		t.Fatalf("loadUserConfig ErrorID: err=%v", err)
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
	} else if ce, ok := err.(*merr.CommonError); !ok || ce.ErrorID() != appRootDirResolveFailedErrorID {
		t.Errorf("loadUserConfig ErrorID: err=%v", err)
	}
}
