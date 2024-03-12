#include "btMDebugDraw.h"

btMDebugDrawLiner::btMDebugDrawLiner() {}
btMDebugDrawLiner::~btMDebugDrawLiner() {}

// void btMDebugDrawLiner::drawLine(const btVector3& from, const btVector3& to, const btVector3& color) {
//     printf("btMDebugDrawLiner::drawLine: from(%f, %f, %f), to(%f, %f, %f), color(%f, %f, %f)\n", from.x(), from.y(), from.z(), to.x(), to.y(), to.z(), color.x(), color.y(), color.z());
// }

// ---------------------------------------------------------
// コンストラクタ
btMDebugDraw::btMDebugDraw() :
    m_debugMode(0)
    , m_liner(0)
    {}
btMDebugDraw::~btMDebugDraw() {}

// 継承先で実装が必要なメソッド
void btMDebugDraw::drawContactPoint(const btVector3& PointOnB, const btVector3& normalOnB, btScalar distance, int lifeTime, const btVector3& color) {
    // Implement drawContactPoint method here
    // printf("drawContactPoint: PointOnB(%f, %f, %f), normalOnB(%f, %f, %f), distance(%f), lifeTime(%d), color(%f, %f, %f)\n", PointOnB.x(), PointOnB.y(), PointOnB.z(), normalOnB.x(), normalOnB.y(), normalOnB.z(), distance, lifeTime, color.x(), color.y(), color.z());

    drawLine( PointOnB, PointOnB + normalOnB * distance, color );
    btVector3 ncolor( 0, 0, 0 );
    drawLine( PointOnB, PointOnB + normalOnB * 0.01f, ncolor );
}

void btMDebugDraw::drawLine(const btVector3& from, const btVector3& to, const btVector3& color) {
    // Implement drawLine method here
    // printf("drawLine: from(%f, %f, %f), to(%f, %f, %f), color(%f, %f, %f)\n", from.x(), from.y(), from.z(), to.x(), to.y(), to.z(), color.x(), color.y(), color.z());
    getLiner()->drawLine(from, to, color);
}

void btMDebugDraw::reportErrorWarning(const char* warningString) {
    // Implement reportErrorWarning method here
    // printf("reportErrorWarning: %s\n", warningString);
}

void btMDebugDraw::draw3dText(const btVector3& location, const char* textString) {
    // Implement draw3dText method here
    // printf("draw3dText: location(%f, %f, %f), textString(%s)\n", location.x(), location.y(), location.z(), textString);
}

// 独自メソッド
void btMDebugDraw::setLiner(btMDebugDrawLiner* liner) {
    // printf("setLiner: %p\n", liner);
    m_liner = liner;
}

btMDebugDrawLiner* btMDebugDraw::getLiner() {
    // printf("getLiner: %p\n", m_liner);
    return m_liner;
}

// ヘッダーで定義済みのメソッド
void btMDebugDraw::setDebugMode(int debugMode) {
    // printf("setDebugMode: %d\n", debugMode);
    m_debugMode = debugMode;
}

int btMDebugDraw::getDebugMode() const {
    // printf("getDebugMode: %d\n", m_debugMode);
    return m_debugMode;
}

// btIDebugDrawで定義済みのメソッド
void btMDebugDraw::setDefaultColors(const DefaultColors& colors) {
    btIDebugDraw::setDefaultColors(colors);
}

void btMDebugDraw::drawLine(const btVector3& from, const btVector3& to, const btVector3& fromColor, const btVector3& toColor) {
    // printf("drawLine: from(%f, %f, %f), to(%f, %f, %f), fromColor(%f, %f, %f), toColor(%f, %f, %f)\n", from.x(), from.y(), from.z(), to.x(), to.y(), to.z(), fromColor.x(), fromColor.y(), fromColor.z(), toColor.x(), toColor.y(), toColor.z());
    btIDebugDraw::drawLine(from, to, fromColor, toColor);
}

void btMDebugDraw::drawSphere(btScalar radius, const btTransform& transform, const btVector3& color) {
    // printf("drawSphere: radius(%f), transform(%f, %f, %f), color(%f, %f, %f)\n", radius, transform.getOrigin().x(), transform.getOrigin().y(), transform.getOrigin().z(), color.x(), color.y(), color.z());
    btIDebugDraw::drawSphere(radius, transform, color);
}

void btMDebugDraw::drawSphere(const btVector3& p, btScalar radius, const btVector3& color) {
    // printf("drawSphere: p(%f, %f, %f), radius(%f), color(%f, %f, %f)\n", p.x(), p.y(), p.z(), radius, color.x(), color.y(), color.z());
    btIDebugDraw::drawSphere(p, radius, color);
}

void btMDebugDraw::drawTriangle(const btVector3& v0, const btVector3& v1, const btVector3& v2, const btVector3& n0, const btVector3& n1, const btVector3& n2, const btVector3& color, btScalar alpha) {
    // printf("drawTriangle: v0(%f, %f, %f), v1(%f, %f, %f), v2(%f, %f, %f), n0(%f, %f, %f), n1(%f, %f, %f), n2(%f, %f, %f), color(%f, %f, %f), alpha(%f)\n", v0.x(), v0.y(), v0.z(), v1.x(), v1.y(), v1.z(), v2.x(), v2.y(), v2.z(), n0.x(), n0.y(), n0.z(), n1.x(), n1.y(), n1.z(), n2.x(), n2.y(), n2.z(), color.x(), color.y(), color.z(), alpha);
    btIDebugDraw::drawTriangle(v0, v1, v2, n0, n1, n2, color, alpha);
}

void btMDebugDraw::drawTriangle(const btVector3& v0, const btVector3& v1, const btVector3& v2, const btVector3& color, btScalar alpha) {
    // printf("drawTriangle: v0(%f, %f, %f), v1(%f, %f, %f), v2(%f, %f, %f), color(%f, %f, %f), alpha(%f)\n", v0.x(), v0.y(), v0.z(), v1.x(), v1.y(), v1.z(), v2.x(), v2.y(), v2.z(), color.x(), color.y(), color.z(), alpha);
    btIDebugDraw::drawTriangle(v0, v1, v2, color, alpha);
}

void btMDebugDraw::drawAabb(const btVector3& from, const btVector3& to, const btVector3& color) {
    // printf("drawAabb: from(%f, %f, %f), to(%f, %f, %f), color(%f, %f, %f)\n", from.x(), from.y(), from.z(), to.x(), to.y(), to.z(), color.x(), color.y(), color.z());
    btIDebugDraw::drawAabb(from, to, color);
}

void btMDebugDraw::drawTransform(const btTransform& transform, btScalar orthoLen) {
    // printf("drawTransform: transform(%f, %f, %f), orthoLen(%f)\n", transform.getOrigin().x(), transform.getOrigin().y(), transform.getOrigin().z(), orthoLen);
    btIDebugDraw::drawTransform(transform, orthoLen);
}

void btMDebugDraw::drawArc(const btVector3& center, const btVector3& normal, const btVector3& axis, btScalar radiusA, btScalar radiusB, btScalar minAngle, btScalar maxAngle, const btVector3& color, bool drawSect, btScalar stepDegrees) {
    // printf("drawArc: center(%f, %f, %f), normal(%f, %f, %f), axis(%f, %f, %f), radiusA(%f), radiusB(%f), minAngle(%f), maxAngle(%f), color(%f, %f, %f), drawSect(%d), stepDegrees(%f)\n", center.x(), center.y(), center.z(), normal.x(), normal.y(), normal.z(), axis.x(), axis.y(), axis.z(), radiusA, radiusB, minAngle, maxAngle, color.x(), color.y(), color.z(), drawSect, stepDegrees);
    btIDebugDraw::drawArc(center, normal, axis, radiusA, radiusB, minAngle, maxAngle, color, drawSect, stepDegrees);
}

void btMDebugDraw::drawSpherePatch(const btVector3& center, const btVector3& up, const btVector3& axis, btScalar radius, btScalar minTh, btScalar maxTh, btScalar minPs, btScalar maxPs, const btVector3& color, btScalar stepDegrees, bool drawCenter) {
    // printf("drawSpherePatch: center(%f, %f, %f), up(%f, %f, %f), axis(%f, %f, %f), radius(%f), minTh(%f), maxTh(%f), minPs(%f), maxPs(%f), color(%f, %f, %f), stepDegrees(%f), drawCenter(%d)\n", center.x(), center.y(), center.z(), up.x(), up.y(), up.z(), axis.x(), axis.y(), axis.z(), radius, minTh, maxTh, minPs, maxPs, color.x(), color.y(), color.z(), stepDegrees, drawCenter);
    btIDebugDraw::drawSpherePatch(center, up, axis, radius, minTh, maxTh, minPs, maxPs, color, stepDegrees, drawCenter);
}

void btMDebugDraw::drawBox(const btVector3& bbMin, const btVector3& bbMax, const btVector3& color) {
    // printf("drawBox: bbMin(%f, %f, %f), bbMax(%f, %f, %f), color(%f, %f, %f)\n", bbMin.x(), bbMin.y(), bbMin.z(), bbMax.x(), bbMax.y(), bbMax.z(), color.x(), color.y(), color.z());
    btIDebugDraw::drawBox(bbMin, bbMax, color);
}

void btMDebugDraw::drawBox(const btVector3& bbMin, const btVector3& bbMax, const btTransform& trans, const btVector3& color) {
    // printf("drawBox: bbMin(%f, %f, %f), bbMax: (%f, %f, %f), trans: (%f, %f, %f), color(%f, %f, %f)\n", bbMin.x(), bbMin.y(), bbMin.z(), bbMax.x(), bbMax.y(), bbMax.z(), color.x(), color.y(), color.z());
    btIDebugDraw::drawBox(bbMin, bbMax, trans, color);
}

void btMDebugDraw::drawCapsule(btScalar radius, btScalar halfHeight, int upAxis, const btTransform& transform, const btVector3& color) {
    // printf("drawCapsule: radius(%f), halfHeight(%f), upAxis(%d), transform(%f, %f, %f), color(%f, %f, %f)\n", radius, halfHeight, upAxis, transform.getOrigin().x(), transform.getOrigin().y(), transform.getOrigin().z(), color.x(), color.y(), color.z());
    btIDebugDraw::drawCapsule(radius, halfHeight, upAxis, transform, color);
}

void btMDebugDraw::drawCylinder(btScalar radius, btScalar halfHeight, int upAxis, const btTransform& transform, const btVector3& color) {
    // printf("drawCylinder: radius(%f), halfHeight(%f), upAxis(%d), transform(%f, %f, %f), color(%f, %f, %f)\n", radius, halfHeight, upAxis, transform.getOrigin().x(), transform.getOrigin().y(), transform.getOrigin().z(), color.x(), color.y(), color.z());
    btIDebugDraw::drawCylinder(radius, halfHeight, upAxis, transform, color);
}

void btMDebugDraw::drawCone(btScalar radius, btScalar height, int upAxis, const btTransform& transform, const btVector3& color) {
    // printf("drawCone: radius(%f), height(%f), upAxis(%d), transform(%f, %f, %f), color(%f, %f, %f)\n", radius, height, upAxis, transform.getOrigin().x(), transform.getOrigin().y(), transform.getOrigin().z(), color.x(), color.y(), color.z());
    btIDebugDraw::drawCone(radius, height, upAxis, transform, color);
}

void btMDebugDraw::drawPlane(const btVector3& planeNormal, btScalar planeConst, const btTransform& transform, const btVector3& color) {
    // printf("drawPlane: planeNormal(%f, %f, %f), planeConst(%f), transform(%f, %f, %f), color(%f, %f, %f)\n", planeNormal.x(), planeNormal.y(), planeNormal.z(), planeConst, transform.getOrigin().x(), transform.getOrigin().y(), transform.getOrigin().z(), color.x(), color.y(), color.z());
    btIDebugDraw::drawPlane(planeNormal, planeConst, transform, color);
}

void btMDebugDraw::clearLines() {
    // printf("clearLines\n");
    btIDebugDraw::clearLines();
}

void btMDebugDraw::flushLines() {
    // printf("flushLines\n");
    btIDebugDraw::flushLines();
}

