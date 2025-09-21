////// BulletCollision/CollisionDispatch/btCollisionWorld.h ----------------

%include "BulletCollision/CollisionDispatch/btCollisionWorld.h"

%{
#include "BulletCollision/CollisionDispatch/btCollisionWorld.h"
%}

%inline %{
// ClosestHit の公開ラッパ
class BtClosestRayCallback : public btCollisionWorld::ClosestRayResultCallback {
public:
    BtClosestRayCallback(const btVector3& from, const btVector3& to)
    : btCollisionWorld::ClosestRayResultCallback(from, to) {}

    bool HasHit() const { return this->hasHit(); }
    const btCollisionObject* GetCollisionObject() const { return m_collisionObject; }
    const btVector3& GetHitPointWorld() const { return m_hitPointWorld; }
    const btVector3& GetHitNormalWorld() const { return m_hitNormalWorld; }
    btScalar GetHitFraction() const { return m_closestHitFraction; }

    void SetCollisionFilterGroup(int g) { m_collisionFilterGroup = g; }
    void SetCollisionFilterMask(int m)  { m_collisionFilterMask  = m; }
    void SetFlags(unsigned int f)       { m_flags = f; }
};

// AllHits の公開ラッパ
class BtAllHitsRayCallback : public btCollisionWorld::AllHitsRayResultCallback {
public:
    BtAllHitsRayCallback(const btVector3& from, const btVector3& to)
    : btCollisionWorld::AllHitsRayResultCallback(from, to) {}

    bool HasHit() const { return this->hasHit(); }

    // 結果配列にアクセスするための簡易 getter
    int  GetNumHits() const { return m_collisionObjects.size(); }
    const btCollisionObject* GetCollisionObject(int i) const { return m_collisionObjects[i]; }
    btScalar GetHitFraction(int i) const { return m_hitFractions[i]; }
    const btVector3& GetHitPointWorld(int i) const { return m_hitPointWorld[i]; }
    const btVector3& GetHitNormalWorld(int i) const { return m_hitNormalWorld[i]; }

    void SetCollisionFilterGroup(int g) { m_collisionFilterGroup = g; }
    void SetCollisionFilterMask(int m)  { m_collisionFilterMask  = m; }
    void SetFlags(unsigned int f)       { m_flags = f; }
};
%}
