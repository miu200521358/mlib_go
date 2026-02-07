// 指示: miu200521358
package base

import (
	"errors"
	"testing"
)

// TestRunWithSetupTeardownOrder は setup->run->teardown の順で呼ばれることを確認する。
func TestRunWithSetupTeardownOrder(t *testing.T) {
	order := make([]string, 0, 3)
	err := RunWithSetupTeardown(
		func() {
			order = append(order, "setup")
		},
		func() {
			order = append(order, "teardown")
		},
		func() error {
			order = append(order, "run")
			return nil
		},
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(order) != 3 || order[0] != "setup" || order[1] != "run" || order[2] != "teardown" {
		t.Fatalf("unexpected order: %v", order)
	}
}

// TestRunWithSetupTeardownError は run のエラーを返しつつ teardown が実行されることを確認する。
func TestRunWithSetupTeardownError(t *testing.T) {
	runErr := errors.New("run failed")
	calledTeardown := false
	err := RunWithSetupTeardown(
		nil,
		func() {
			calledTeardown = true
		},
		func() error {
			return runErr
		},
	)
	if !errors.Is(err, runErr) {
		t.Fatalf("expected run error, got: %v", err)
	}
	if !calledTeardown {
		t.Fatalf("teardown should be called")
	}
}

// TestRunWithBoolState は一時値設定と復元が行われることを確認する。
func TestRunWithBoolState(t *testing.T) {
	state := false
	seenInRun := false
	err := RunWithBoolState(
		func(v bool) {
			state = v
		},
		true,
		false,
		func() error {
			seenInRun = state
			return nil
		},
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !seenInRun {
		t.Fatalf("state should be true in run")
	}
	if state {
		t.Fatalf("state should be restored to false")
	}
}
