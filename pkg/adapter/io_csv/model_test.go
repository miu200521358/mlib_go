// 指示: miu200521358
package io_csv

import (
	"reflect"
	"testing"
)

func TestCsvModelRecords(t *testing.T) {
	model := NewCsvModel([][]string{
		{"h1", "h2"},
		{"a", "b"},
	})

	records := model.Records()
	records[1][0] = "x"

	if model.Records()[1][0] != "a" {
		t.Fatalf("expected model records to be immutable copy")
	}
}

func TestCsvModelSetRecords(t *testing.T) {
	model := NewCsvModel(nil)
	source := [][]string{{"A", "1"}}

	model.SetRecords(source)
	source[0][0] = "B"

	expected := [][]string{{"A", "1"}}
	if !reflect.DeepEqual(model.Records(), expected) {
		t.Fatalf("expected %v, got %v", expected, model.Records())
	}
}
