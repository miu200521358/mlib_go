// 指示: miu200521358
package err

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
