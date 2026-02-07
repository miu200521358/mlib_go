// 指示: miu200521358
package base

// RunWithSetupTeardown は setup 実行後に run を呼び、最後に teardown を必ず実行する。
func RunWithSetupTeardown(setup func(), teardown func(), run func() error) error {
	if setup != nil {
		setup()
	}
	if teardown != nil {
		defer teardown()
	}
	if run == nil {
		return nil
	}
	return run()
}

// RunWithBoolState は bool setter を一時値に切り替えて処理後に復元する。
func RunWithBoolState(setter func(bool), temporary bool, restore bool, run func() error) error {
	setup := func() {
		if setter != nil {
			setter(temporary)
		}
	}
	teardown := func() {
		if setter != nil {
			setter(restore)
		}
	}
	return RunWithSetupTeardown(setup, teardown, run)
}
