//go:build windows
// +build windows

// 指示: miu200521358
package mbullet

import (
	"github.com/miu200521358/mlib_go/pkg/adapter/physics_api"
	"github.com/miu200521358/mlib_go/pkg/domain/mmath"
	"github.com/miu200521358/mlib_go/pkg/infra/drivers/mbullet/bt"
)

// RayTest はレイキャスト結果を取得する。
func (mp *PhysicsEngine) RayTest(
	rayFrom, rayTo *mmath.Vec3,
	filter *physics_api.RaycastFilter,
) (raycastHit *physics_api.RaycastHit) {
	if mp == nil || rayFrom == nil || rayTo == nil {
		return nil
	}
	defer func() {
		if recover() != nil {
			raycastHit = nil
		}
	}()

	// MMD座標系のレイを Bullet 座標系へ変換する。
	btFrom := newBulletFromVec(*rayFrom)
	btTo := newBulletFromVec(*rayTo)
	defer bt.DeleteBtVector3(btFrom)
	defer bt.DeleteBtVector3(btTo)

	cb := bt.NewBtClosestRayCallback(btFrom, btTo)
	defer bt.DeleteBtClosestRayCallback(cb)

	if filter != nil {
		cb.SetCollisionFilterGroup(filter.Group)
		cb.SetCollisionFilterMask(filter.Mask)
	}

	mp.world.RayTest(btFrom, btTo, cb)

	if !cb.HasHit() {
		return nil
	}

	modelIndex, rigidBodyIndex, ok := mp.findRigidBodyByCollisionObject(cb.GetCollisionObject(), true)
	if !ok {
		return nil
	}

	raycastHit = &physics_api.RaycastHit{
		RigidBodyRef: physics_api.RigidBodyRef{
			ModelIndex:     modelIndex,
			RigidBodyIndex: rigidBodyIndex,
		},
		HitFraction: cb.GetHitFraction(),
	}
	return raycastHit
}
