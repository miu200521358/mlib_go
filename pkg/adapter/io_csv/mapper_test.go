// 指示: miu200521358
package io_csv

import (
	"reflect"
	"testing"

	"github.com/miu200521358/mlib_go/pkg/shared/base/merr"
)

type sampleCsvRow struct {
	Key     string `csv:"キー"`
	Value   int    `csv:"値"`
	Enabled bool   `csv:"有効"`
}

func TestMarshal(t *testing.T) {
	rows := []sampleCsvRow{
		{Key: "A", Value: 10, Enabled: true},
		{Key: "B", Value: 20, Enabled: false},
	}

	model, err := Marshal(rows)
	if err != nil {
		t.Fatalf("expected marshal to succeed, got %v", err)
	}

	expected := [][]string{
		{"キー", "値", "有効"},
		{"A", "10", "true"},
		{"B", "20", "false"},
	}
	if !reflect.DeepEqual(model.Records(), expected) {
		t.Fatalf("expected records %v, got %v", expected, model.Records())
	}
}

func TestMarshalInvalidInput(t *testing.T) {
	_, err := Marshal(map[string]string{"a": "b"})
	if err == nil {
		t.Fatalf("expected error")
	}
	if merr.ExtractErrorID(err) != "14106" {
		t.Fatalf("expected error id 14106, got %s", merr.ExtractErrorID(err))
	}
}

func TestUnmarshal(t *testing.T) {
	model := NewCsvModel([][]string{
		{"キー", "値", "有効"},
		{"A", "10", "true"},
		{"B", "20", "false"},
	})

	var rows []sampleCsvRow
	if err := Unmarshal(model, &rows); err != nil {
		t.Fatalf("expected unmarshal to succeed, got %v", err)
	}
	if len(rows) != 2 {
		t.Fatalf("expected 2 rows, got %d", len(rows))
	}
	if rows[0].Key != "A" || rows[0].Value != 10 || !rows[0].Enabled {
		t.Fatalf("unexpected first row: %+v", rows[0])
	}
	if rows[1].Key != "B" || rows[1].Value != 20 || rows[1].Enabled {
		t.Fatalf("unexpected second row: %+v", rows[1])
	}
}

func TestUnmarshalPointerSlice(t *testing.T) {
	model := NewCsvModel([][]string{
		{"キー", "値", "有効"},
		{"A", "10", "true"},
	})

	var rows []*sampleCsvRow
	if err := Unmarshal(model, &rows); err != nil {
		t.Fatalf("expected unmarshal to succeed, got %v", err)
	}
	if len(rows) != 1 || rows[0] == nil {
		t.Fatalf("unexpected rows: %v", rows)
	}
	if rows[0].Key != "A" || rows[0].Value != 10 || !rows[0].Enabled {
		t.Fatalf("unexpected row: %+v", rows[0])
	}
}

func TestUnmarshalParseError(t *testing.T) {
	model := NewCsvModel([][]string{
		{"キー", "値", "有効"},
		{"A", "not-number", "true"},
	})

	var rows []sampleCsvRow
	err := Unmarshal(model, &rows)
	if err == nil {
		t.Fatalf("expected error")
	}
	if merr.ExtractErrorID(err) != "14105" {
		t.Fatalf("expected error id 14105, got %s", merr.ExtractErrorID(err))
	}
}
