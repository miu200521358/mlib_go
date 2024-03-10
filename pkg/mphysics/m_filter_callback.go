package mphysics

import "github.com/miu200521358/mlib_go/pkg/mbt"

type MFilterCallback struct {
	mbt.BtOverlapFilterCallback
	nonFilterProxy []mbt.BtBroadphaseProxy
}

func NewMFilterCallback() *MFilterCallback {
	return &MFilterCallback{
		BtOverlapFilterCallback: mbt.NewBtHashedOverlappingPairCache().GetOverlapFilterCallback(),
	}
}

func (m *MFilterCallback) NeedBroadphaseCollision(proxy0, proxy1 mbt.BtBroadphaseProxy) bool {
	// proxy0とproxy1が同じグループに属している場合には衝突検出を行わない
	collides := (proxy0.GetM_collisionFilterGroup() & proxy1.GetM_collisionFilterMask()) != 0
	collides = collides && ((proxy1.GetM_collisionFilterGroup() & proxy0.GetM_collisionFilterMask()) != 0)
	// Add additional logic here to modify 'collides'
	return collides
}

func (m *MFilterCallback) Swigcptr() uintptr {
	return m.BtOverlapFilterCallback.Swigcptr()
}

func (m *MFilterCallback) SwigIsBtOverlapFilterCallback() {
}
