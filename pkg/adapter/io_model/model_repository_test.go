// 指示: miu200521358
package io_model

import (
	"path/filepath"
	"testing"

	"github.com/miu200521358/mlib_go/pkg/shared/base/merr"
)

func TestModelRepositoryCanLoadVrm(t *testing.T) {
	repository := NewModelRepository()
	if !repository.CanLoad("sample.vrm") {
		t.Fatalf("expected sample.vrm to be loadable")
	}
}

func TestModelRepositoryLoadVrmRouted(t *testing.T) {
	repository := NewModelRepository()
	_, err := repository.Load(filepath.Join(t.TempDir(), "missing.vrm"))
	if err == nil {
		t.Fatalf("expected error")
	}
	if merr.ExtractErrorID(err) != "14101" {
		t.Fatalf("expected error id 14101, got %s", merr.ExtractErrorID(err))
	}
}
