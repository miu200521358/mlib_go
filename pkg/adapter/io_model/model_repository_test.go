// 指示: miu200521358
package io_model

import (
	"path/filepath"
	"testing"

	"github.com/miu200521358/mlib_go/pkg/shared/base/merr"
)

func TestModelRepositoryCanLoadSupportedExt(t *testing.T) {
	repository := NewModelRepository()
	if !repository.CanLoad("sample.pmx") {
		t.Fatalf("expected sample.pmx to be loadable")
	}
	if !repository.CanLoad("sample.pmd") {
		t.Fatalf("expected sample.pmd to be loadable")
	}
	if !repository.CanLoad("sample.x") {
		t.Fatalf("expected sample.x to be loadable")
	}
}

func TestModelRepositoryCanLoadRejectsVrm(t *testing.T) {
	repository := NewModelRepository()
	if repository.CanLoad("sample.vrm") {
		t.Fatalf("expected sample.vrm to be not loadable")
	}
}

func TestModelRepositoryLoadInvalidExt(t *testing.T) {
	repository := NewModelRepository()
	_, err := repository.Load(filepath.Join(t.TempDir(), "sample.vrm"))
	if err == nil {
		t.Fatalf("expected error")
	}
	if merr.ExtractErrorID(err) != "14102" {
		t.Fatalf("expected error id 14102, got %s", merr.ExtractErrorID(err))
	}
}
