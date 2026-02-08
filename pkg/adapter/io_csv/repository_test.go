// 指示: miu200521358
package io_csv

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/miu200521358/mlib_go/pkg/adapter/io_common"
	"github.com/miu200521358/mlib_go/pkg/shared/base/merr"
	"github.com/miu200521358/mlib_go/pkg/shared/hashable"
)

func TestCsvRepositoryCanLoad(t *testing.T) {
	repository := NewCsvRepository()

	if !repository.CanLoad("sample.csv") {
		t.Fatalf("expected sample.csv to be loadable")
	}
	if !repository.CanLoad("sample.CSV") {
		t.Fatalf("expected sample.CSV to be loadable")
	}
	if repository.CanLoad("sample.txt") {
		t.Fatalf("expected sample.txt to be not loadable")
	}
}

func TestCsvRepositoryLoadNotFound(t *testing.T) {
	repository := NewCsvRepository()

	_, err := repository.Load(filepath.Join(t.TempDir(), "missing.csv"))
	if err == nil {
		t.Fatalf("expected error")
	}
	if merr.ExtractErrorID(err) != "14101" {
		t.Fatalf("expected error id 14101, got %s", merr.ExtractErrorID(err))
	}
}

func TestCsvRepositoryLoadInvalidExt(t *testing.T) {
	repository := NewCsvRepository()

	_, err := repository.Load(filepath.Join(t.TempDir(), "sample.txt"))
	if err == nil {
		t.Fatalf("expected error")
	}
	if merr.ExtractErrorID(err) != "14102" {
		t.Fatalf("expected error id 14102, got %s", merr.ExtractErrorID(err))
	}
}

func TestCsvRepositoryLoadSaveRoundTrip(t *testing.T) {
	repository := NewCsvRepository()
	path := filepath.Join(t.TempDir(), "sample.csv")

	source := NewCsvModel([][]string{
		{"列1", "列2"},
		{"A", "1"},
		{"B", "2"},
	})
	if err := repository.Save(path, source, io_common.SaveOptions{}); err != nil {
		t.Fatalf("expected save to succeed, got %v", err)
	}

	loadedData, err := repository.Load(path)
	if err != nil {
		t.Fatalf("expected load to succeed, got %v", err)
	}
	loaded, ok := loadedData.(*CsvModel)
	if !ok {
		t.Fatalf("expected loaded type to be *CsvModel, got %T", loadedData)
	}
	if loaded.Path() != path {
		t.Fatalf("expected path %q, got %q", path, loaded.Path())
	}
	if loaded.Name() != "sample" {
		t.Fatalf("expected name sample, got %q", loaded.Name())
	}
	if !reflect.DeepEqual(loaded.Records(), source.Records()) {
		t.Fatalf("expected records %v, got %v", source.Records(), loaded.Records())
	}
}

func TestCsvRepositoryLoadWithCRLF(t *testing.T) {
	repository := NewCsvRepository()
	path := filepath.Join(t.TempDir(), "sample.csv")
	if err := os.WriteFile(path, []byte("キー,値\r\nA,1\r\n"), 0o644); err != nil {
		t.Fatalf("failed to write csv: %v", err)
	}

	data, err := repository.Load(path)
	if err != nil {
		t.Fatalf("expected load to succeed, got %v", err)
	}
	model, ok := data.(*CsvModel)
	if !ok {
		t.Fatalf("expected loaded type to be *CsvModel, got %T", data)
	}
	expected := [][]string{{"キー", "値"}, {"A", "1"}}
	if !reflect.DeepEqual(model.Records(), expected) {
		t.Fatalf("expected records %v, got %v", expected, model.Records())
	}
}

func TestCsvRepositoryLoadProfileExactColumns(t *testing.T) {
	repository := NewCsvRepository()
	repository.SetProfile(CsvProfile{
		HasHeader:    true,
		ExactColumns: 2,
	})
	path := filepath.Join(t.TempDir(), "sample.csv")
	if err := os.WriteFile(path, []byte("キー,値\nA\n"), 0o644); err != nil {
		t.Fatalf("failed to write csv: %v", err)
	}

	_, err := repository.Load(path)
	if err == nil {
		t.Fatalf("expected error")
	}
	if merr.ExtractErrorID(err) != "14105" {
		t.Fatalf("expected error id 14105, got %s", merr.ExtractErrorID(err))
	}
}

func TestCsvRepositoryLoadProfileHeaderMismatch(t *testing.T) {
	repository := NewCsvRepositoryWithProfile(CsvProfile{
		HasHeader:    true,
		ExactColumns: 2,
		Header:       []string{"キー", "値"},
	})
	path := filepath.Join(t.TempDir(), "sample.csv")
	if err := os.WriteFile(path, []byte("name,value\nA,1\n"), 0o644); err != nil {
		t.Fatalf("failed to write csv: %v", err)
	}

	_, err := repository.Load(path)
	if err == nil {
		t.Fatalf("expected error")
	}
	if merr.ExtractErrorID(err) != "14105" {
		t.Fatalf("expected error id 14105, got %s", merr.ExtractErrorID(err))
	}
}

func TestCsvRepositoryLoadParseFailed(t *testing.T) {
	repository := NewCsvRepository()
	path := filepath.Join(t.TempDir(), "sample.csv")
	if err := os.WriteFile(path, []byte("キー,値\n\"A,1\n"), 0o644); err != nil {
		t.Fatalf("failed to write csv: %v", err)
	}

	_, err := repository.Load(path)
	if err == nil {
		t.Fatalf("expected error")
	}
	if merr.ExtractErrorID(err) != "14105" {
		t.Fatalf("expected error id 14105, got %s", merr.ExtractErrorID(err))
	}
}

func TestCsvRepositorySaveInvalidModel(t *testing.T) {
	repository := NewCsvRepository()
	path := filepath.Join(t.TempDir(), "sample.csv")

	err := repository.Save(path, hashable.NewHashableBase("", ""), io_common.SaveOptions{})
	if err == nil {
		t.Fatalf("expected error")
	}
	if merr.ExtractErrorID(err) != "14106" {
		t.Fatalf("expected error id 14106, got %s", merr.ExtractErrorID(err))
	}
}

func TestCsvRepositorySaveCreateFailed(t *testing.T) {
	repository := NewCsvRepository()
	path := filepath.Join(t.TempDir(), "missing_dir", "sample.csv")
	data := NewCsvModel([][]string{{"キー", "値"}})

	err := repository.Save(path, data, io_common.SaveOptions{})
	if err == nil {
		t.Fatalf("expected error")
	}
	if merr.ExtractErrorID(err) != "14108" {
		t.Fatalf("expected error id 14108, got %s", merr.ExtractErrorID(err))
	}
}
