////// BulletCollision/CollisionDispatch/btCollisionWorld.h ----------------

%include "BulletCollision/CollisionDispatch/btCollisionWorld.h"

%{
#include "BulletCollision/CollisionDispatch/btCollisionWorld.h"
%}

%inline %{
// Bullet の ClosestRayResultCallback をそのまま Go から new できる薄いラッパ
class BtClosestRayCallback : public btCollisionWorld::ClosestRayResultCallback {
public:
    BtClosestRayCallback(const btVector3& from, const btVector3& to)
    : btCollisionWorld::ClosestRayResultCallback(from, to) {}

    // Bullet 本家の RayResultCallback::hasHit() をそのまま出す
    bool HasHit() const { return this->hasHit(); }

    // ヒット情報を読み出すための getter を用意
    const btCollisionObject* GetCollisionObject() const { return m_collisionObject; }
    const btVector3&        GetHitPointWorld()   const { return m_hitPointWorld; }
    const btVector3&        GetHitNormalWorld()  const { return m_hitNormalWorld; }
    btScalar                GetHitFraction()     const { return m_closestHitFraction; }
};
%}
