// 指示: miu200521358
package merr

import (
	"errors"
	"io/fs"
	"strings"
	"testing"
)

// TestLoadDefaultRegistry は埋め込みCSVの読込を確認する。
func TestLoadDefaultRegistry(t *testing.T) {
	recs, err := LoadDefaultRegistry()
	if err != nil {
		t.Fatalf("LoadDefaultRegistry failed: %v", err)
	}
	if len(recs) == 0 {
		t.Fatalf("LoadDefaultRegistry returned empty")
	}
	if recs[0].ID == "" {
		t.Errorf("first record ID empty")
	}
}

// TestLoadDefaultRegistryError はOpen失敗を確認する。
func TestLoadDefaultRegistryError(t *testing.T) {
	prev := openRegistryFile
	openRegistryFile = func(string) (fs.File, error) {
		return nil, errors.New("open error")
	}
	t.Cleanup(func() { openRegistryFile = prev })

	if _, err := LoadDefaultRegistry(); err == nil {
		t.Errorf("LoadDefaultRegistry expected error")
	}
}

// TestLoadRegistryEmpty は空CSVの扱いを確認する。
func TestLoadRegistryEmpty(t *testing.T) {
	recs, err := LoadRegistry(strings.NewReader(""))
	if err != nil {
		t.Fatalf("LoadRegistry failed: %v", err)
	}
	if len(recs) != 0 {
		t.Errorf("LoadRegistry empty: got=%v", len(recs))
	}
}

// TestLoadRegistryHeader はヘッダとパス分割を確認する。
func TestLoadRegistryHeader(t *testing.T) {
	csv := "ID,Kind,Layer,Module,ErrorName,Summary,Remedy,SourcePaths\n" +
		"100,Internal,shared,01,ErrA,summary,remedy, a.go ; b.go \n" +
		"short,only,five,cols\n"
	recs, err := LoadRegistry(strings.NewReader(csv))
	if err != nil {
		t.Fatalf("LoadRegistry failed: %v", err)
	}
	if len(recs) != 1 {
		t.Fatalf("LoadRegistry length: got=%v", len(recs))
	}
	if len(recs[0].SourcePaths) != 2 || recs[0].SourcePaths[0] != "a.go" || recs[0].SourcePaths[1] != "b.go" {
		t.Errorf("SourcePaths: got=%v", recs[0].SourcePaths)
	}
}

// TestLoadRegistryNoHeader はヘッダなしと空パスを確認する。
func TestLoadRegistryNoHeader(t *testing.T) {
	csv := "200,External,infra,02,ErrB,summary,remedy,-\n"
	recs, err := LoadRegistry(strings.NewReader(csv))
	if err != nil {
		t.Fatalf("LoadRegistry failed: %v", err)
	}
	if len(recs) != 1 {
		t.Fatalf("LoadRegistry length: got=%v", len(recs))
	}
	if recs[0].SourcePaths != nil {
		t.Errorf("SourcePaths: expected nil got=%v", recs[0].SourcePaths)
	}
}

// TestLoadRegistryInvalid はCSVエラーを確認する。
func TestLoadRegistryInvalid(t *testing.T) {
	if _, err := LoadRegistry(strings.NewReader("\"bad")); err == nil {
		t.Errorf("LoadRegistry expected error")
	}
}

// TestSplitPaths はパス分割の分岐を確認する。
func TestSplitPaths(t *testing.T) {
	if splitPaths("") != nil || splitPaths("-") != nil {
		t.Errorf("splitPaths empty should be nil")
	}
	out := splitPaths(" a ; ; b ")
	if len(out) != 2 || out[0] != "a" || out[1] != "b" {
		t.Errorf("splitPaths: got=%v", out)
	}
}

// TestCommonErrorNil はnil受信時の挙動を確認する。
func TestCommonErrorNil(t *testing.T) {
	var ce *CommonError
	if ce.Error() != "" {
		t.Errorf("nil Error should be empty")
	}
	if ce.Unwrap() != nil {
		t.Errorf("nil Unwrap should be nil")
	}
	if ce.ErrorID() != "" {
		t.Errorf("nil ErrorID should be empty")
	}
	if ce.ErrorKind() != "" {
		t.Errorf("nil ErrorKind should be empty")
	}
}

// TestCommonErrorValues は値とメッセージの組み立てを確認する。
func TestCommonErrorValues(t *testing.T) {
	cause := errors.New("cause")
	ce := NewCommonError("id", ErrorKindExternal, "msg", cause)
	if ce.Error() != "msg: cause" {
		t.Errorf("Error: got=%v", ce.Error())
	}
	if ce.Unwrap() != cause {
		t.Errorf("Unwrap mismatch")
	}
	if ce.ErrorID() != "id" {
		t.Errorf("ErrorID: got=%v", ce.ErrorID())
	}
	if ce.ErrorKind() != ErrorKindExternal {
		t.Errorf("ErrorKind: got=%v", ce.ErrorKind())
	}
}

// TestCommonErrorMessageOnly はメッセージのみの分岐を確認する。
func TestCommonErrorMessageOnly(t *testing.T) {
	ce := NewCommonError("id", ErrorKindInternal, "only", nil)
	if ce.Error() != "only" {
		t.Errorf("Error: got=%v", ce.Error())
	}
}

// TestCommonErrorCauseOnly は原因のみの分岐を確認する。
func TestCommonErrorCauseOnly(t *testing.T) {
	cause := errors.New("only-cause")
	ce := NewCommonError("id", ErrorKindInternal, "", cause)
	if ce.Error() != "only-cause" {
		t.Errorf("Error: got=%v", ce.Error())
	}
}

// TestNewPackageErrors は共通委譲エラーのID/種別を確認する。
func TestNewPackageErrors(t *testing.T) {
	cause := errors.New("cause")
	tests := []struct {
		name string
		err  *CommonError
		id   string
	}{
		{name: "os", err: NewOsPackageError("msg", cause), id: OsPackageErrorID},
		{name: "json", err: NewJsonPackageError("msg", cause), id: JsonPackageErrorID},
		{name: "image", err: NewImagePackageError("msg", cause), id: ImagePackageErrorID},
		{name: "fs", err: NewFsPackageError("msg", cause), id: FsPackageErrorID},
		{name: "deepcopy", err: NewDeepcopyPackageError("msg", cause), id: DeepcopyPackageErrorID},
	}
	for _, tt := range tests {
		if tt.err.ErrorID() != tt.id {
			t.Errorf("%s ErrorID: got=%v", tt.name, tt.err.ErrorID())
		}
		if tt.err.ErrorKind() != ErrorKindExternal {
			t.Errorf("%s ErrorKind: got=%v", tt.name, tt.err.ErrorKind())
		}
	}
}
