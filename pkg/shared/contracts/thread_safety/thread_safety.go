// 指示: miu200521358
package thread_safety

// ThreadSafetyPolicy はスレッド安全性の方針を表す。
type ThreadSafetyPolicy struct {
	ModelMutable              bool
	CopyRequiredForConcurrency bool
	RepositoryNotThreadSafe   bool
	AtomicCacheAllowed        bool
}

// DEFAULT_THREAD_SAFETY_POLICY は既定の方針。
var DEFAULT_THREAD_SAFETY_POLICY = ThreadSafetyPolicy{
	ModelMutable:              true,
	CopyRequiredForConcurrency: true,
	RepositoryNotThreadSafe:   true,
	AtomicCacheAllowed:        true,
}
