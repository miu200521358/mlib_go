// 指示: miu200521358
package thread_safety

import "testing"

// TestDefaultThreadSafetyPolicy は既定値を確認する。
func TestDefaultThreadSafetyPolicy(t *testing.T) {
	if !DEFAULT_THREAD_SAFETY_POLICY.ModelMutable {
		t.Errorf("ModelMutable: got=%v", DEFAULT_THREAD_SAFETY_POLICY.ModelMutable)
	}
	if !DEFAULT_THREAD_SAFETY_POLICY.CopyRequiredForConcurrency {
		t.Errorf("CopyRequiredForConcurrency: got=%v", DEFAULT_THREAD_SAFETY_POLICY.CopyRequiredForConcurrency)
	}
	if !DEFAULT_THREAD_SAFETY_POLICY.RepositoryNotThreadSafe {
		t.Errorf("RepositoryNotThreadSafe: got=%v", DEFAULT_THREAD_SAFETY_POLICY.RepositoryNotThreadSafe)
	}
	if !DEFAULT_THREAD_SAFETY_POLICY.AtomicCacheAllowed {
		t.Errorf("AtomicCacheAllowed: got=%v", DEFAULT_THREAD_SAFETY_POLICY.AtomicCacheAllowed)
	}
}
