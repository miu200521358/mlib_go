// 指示: miu200521358
package io_csv

import (
	"reflect"
	"testing"
)

func TestNewFixedColumn4Profile(t *testing.T) {
	profile := NewFixedColumn4Profile()

	if !profile.HasHeader {
		t.Fatalf("expected HasHeader true")
	}
	if profile.ExactColumns != 4 {
		t.Fatalf("expected ExactColumns 4, got %d", profile.ExactColumns)
	}
	if profile.AllowExtraColumns {
		t.Fatalf("expected AllowExtraColumns false")
	}
}

func TestNewSimpleKeyValueProfile(t *testing.T) {
	profile := NewSimpleKeyValueProfile()

	expectedHeader := []string{"キー", "値"}
	if !profile.HasHeader {
		t.Fatalf("expected HasHeader true")
	}
	if profile.ExactColumns != 2 {
		t.Fatalf("expected ExactColumns 2, got %d", profile.ExactColumns)
	}
	if !reflect.DeepEqual(profile.Header, expectedHeader) {
		t.Fatalf("expected header %v, got %v", expectedHeader, profile.Header)
	}
}

func TestNewFreeTableProfile(t *testing.T) {
	profile := NewFreeTableProfile()

	if !profile.HasHeader {
		t.Fatalf("expected HasHeader true")
	}
	if profile.MinColumns != 1 {
		t.Fatalf("expected MinColumns 1, got %d", profile.MinColumns)
	}
	if !profile.AllowExtraColumns {
		t.Fatalf("expected AllowExtraColumns true")
	}
}
