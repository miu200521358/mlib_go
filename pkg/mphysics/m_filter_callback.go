package mphysics

import "github.com/miu200521358/mlib_go/pkg/mbt"

type MFilterCallback struct {
	mbt.BtOverlapFilterCallback
	m_nonFilterProxy []mbt.BtBroadphaseProxy
}

func NewMFilterCallback(cb mbt.BtOverlapFilterCallback) MFilterCallback {
	return MFilterCallback{
		BtOverlapFilterCallback: cb,
	}
}

func (m *MFilterCallback) NeedBroadphaseCollision(proxy0, proxy1 mbt.BtBroadphaseProxy) bool {
	findIt := func() mbt.BtBroadphaseProxy {
		for _, x := range m.m_nonFilterProxy {
			if x == proxy0 || x == proxy1 {
				return x
			}
		}
		return nil
	}()

	if findIt != nil {
		return true
	}

	collides1 := (proxy0.GetM_collisionFilterGroup() & proxy1.GetM_collisionFilterMask()) != 0
	collides2 := (proxy0.GetM_collisionFilterMask() & proxy1.GetM_collisionFilterGroup()) != 0
	return collides1 && collides2
}

func (m *MFilterCallback) Swigcptr() uintptr {
	return m.BtOverlapFilterCallback.Swigcptr()
}

func (m *MFilterCallback) SwigIsBtOverlapFilterCallback() {
}
