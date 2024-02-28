%module mbt

%include <cmath>
%include <string>

%{
    #include <cmath>
    #include <string>
%}

// %include <float.h> ----------------------------------------------

%{
#define FLT_EPSILON      1.192092896e-07F        // smallest such that 1.0+FLT_EPSILON != 1.0
#define FLT_MAX          3.402823466e+38F        // max value
%}

// %include <math.h> ----------------------------------------------
%{

/* 7.12.4 Trigonometric functions: Double in C89 */
  extern float __cdecl sinf(float _X);
  extern long double __cdecl sinl(long double);

  extern float __cdecl cosf(float _X);
  extern long double __cdecl cosl(long double);

  extern float __cdecl tanf(float _X);
  extern long double __cdecl tanl(long double);
  extern float __cdecl asinf(float _X);
  extern long double __cdecl asinl(long double);

  extern float __cdecl acosf (float);
  extern long double __cdecl acosl (long double);

  extern float __cdecl atanf (float);
  extern long double __cdecl atanl (long double);

  extern float __cdecl atan2f (float, float);
  extern long double __cdecl atan2l (long double, long double);

/* 7.12.6.1 Double in C89 */
  extern float __cdecl expf(float _X);

/* 7.12.6.7 Double in C89 */
  extern float __cdecl logf (float);

/* 7.12.7.4 The pow functions. Double in C89 */
  extern float __cdecl powf(float _X,float _Y);

/* 7.12.7.5 The sqrt functions. Double in C89. */
  extern float __cdecl sqrtf (float);

/* 7.12.7.2 The fabs functions: Double in C89 */
  extern  float __cdecl fabsf (float x);

/* 7.12.10.1 Double in C89 */
  extern float __cdecl fmodf (float, float);

%}


////// included headers [LinearMath/btDefaultMotionState.h] ----------------------------------
%include "LinearMath/btScalar.h.i"
%include "LinearMath/btMinMax.h.i"
%include "LinearMath/btAlignedAllocator.cpp.i"
%include "LinearMath/btAlignedAllocator.h.i"
%include "LinearMath/btVector3.h.i"
%include "LinearMath/btVector3.h.i"
%include "LinearMath/btQuadWord.h.i"
%include "LinearMath/btQuaternion.h.i"
%include "LinearMath/btMatrix3x3.h.i"
%include "LinearMath/btTransform.h.i"
%include "LinearMath/btMotionState.h.i"
%include "LinearMath/btDefaultMotionState.h.i"


////// included headers [BulletCollision/CollisionShapes/btSphereShape.h] ----------------------------------
%include "BulletCollision/BroadphaseCollision/btBroadphaseProxy.h.i"
%include "BulletCollision/CollisionShapes/btCollisionShape.cpp.i"
%include "BulletCollision/CollisionShapes/btCollisionShape.h.i"
%include "BulletCollision/CollisionShapes/btCollisionMargin.h.i"
%include "BulletCollision/CollisionShapes/btConvexShape.cpp.i"
%include "BulletCollision/CollisionShapes/btConvexShape.h.i"
%include "LinearMath/btAabbUtil2.h.i"
%include "BulletCollision/CollisionShapes/btConvexInternalShape.cpp.i"
%include "BulletCollision/CollisionShapes/btConvexInternalShape.h.i"
%include "BulletCollision/CollisionShapes/btSphereShape.cpp.i"
%include "BulletCollision/CollisionShapes/btSphereShape.h.i"

// ////// included headers [BulletCollision/CollisionShapes/btBoxShape.h] ----------------------------------
// %include "BulletCollision/BroadphaseCollision/btBroadphaseProxy.h.i"
// %include "BulletCollision/CollisionShapes/btCollisionShape.h.i"
// %include "BulletCollision/CollisionShapes/btCollisionMargin.h.i"
// %include "BulletCollision/CollisionShapes/btConvexShape.h.i"
// %include "LinearMath/btAabbUtil2.h.i"
// %include "BulletCollision/CollisionShapes/btConvexInternalShape.h.i"
// %include "BulletCollision/CollisionShapes/btPolyhedralConvexShape.h.i"
// %include "BulletCollision/CollisionShapes/btPolyhedralConvexShape.cpp"
// %include "BulletCollision/CollisionShapes/btBoxShape.h.i"
// %include "BulletCollision/CollisionShapes/btBoxShape.cpp"

// %include "LinearMath/btConvexHullComputer.h.i"
// %include "LinearMath/btConvexHullComputer.cpp.i"
// %include "LinearMath/btAlignedObjectArray.h.i"
// %include "LinearMath/btGeometryUtil.h.i"
// %include "LinearMath/btGeometryUtil.cpp.i"
// %include "BulletCollision/CollisionShapes/btConvexPolyhedron.h.i"
// %include "LinearMath/btHashMap.h.i"
// %include "BulletCollision/CollisionShapes/btConvexPolyhedron.cpp.i"
// %include "BulletCollision/CollisionShapes/btPolyhedralConvexShape.h.i"
// %include "BulletCollision/CollisionShapes/btPolyhedralConvexShape.cpp.i"
// %include "BulletCollision/CollisionShapes/btBoxShape.h.i"
// %include "BulletCollision/CollisionShapes/btBoxShape.cpp.i"


// ////// included headers [BulletCollision/CollisionShapes/btCapsuleShape.h] ----------------------------------
// %include "BulletCollision/CollisionShapes/btCapsuleShape.h.i"
// %include "BulletCollision/CollisionShapes/btCapsuleShape.cpp.i"


// ////// included headers [BulletDynamics/Dynamics/btRigidBody.h] ----------------------------------
// %include "BulletCollision/CollisionDispatch/btCollisionObject.h.i"
// %include "BulletCollision/CollisionDispatch/btCollisionObject.cpp.i"
// %include "BulletDynamics/Dynamics/btRigidBody.h.i"
// %include "BulletDynamics/Dynamics/btRigidBody.cpp.i"


// ////// included headers [BulletDynamics/ConstraintSolver/btTypedConstraint.h] ----------------------------------
// %include "BulletDynamics/ConstraintSolver/btJacobianEntry.h.i"
// %include "LinearMath/btTransformUtil.h.i"
// %include "BulletDynamics/ConstraintSolver/btSolverBody.h.i"
// %include "BulletDynamics/ConstraintSolver/btSolverConstraint.h.i"
// %include "BulletDynamics/ConstraintSolver/btTypedConstraint.h.i"
// %include "BulletDynamics/ConstraintSolver/btTypedConstraint.cpp.i"


// ////// included headers [BulletDynamics/Dynamics/btDiscreteDynamicsWorld.h] ----------------------------------
// %include "BulletCollision/CollisionShapes/btMinkowskiSumShape.h.i"
// %include "BulletCollision/CollisionShapes/btMinkowskiSumShape.cpp.i"
// %include "BulletCollision/CollisionDispatch/btCollisionObjectWrapper.h.i"
// %include "BulletCollision/CollisionShapes/btTriangleCallback.h.i"
// %include "BulletCollision/CollisionShapes/btConcaveShape.h.i"
// %include "BulletCollision/CollisionShapes/btStridingMeshInterface.h.i"
// %include "LinearMath/btSerializer.h.i"
// %include "BulletCollision/CollisionShapes/btStridingMeshInterface.cpp.i"
// %include "BulletDynamics/ConstraintSolver/btContactSolverInfo.h.i"
// %include "BulletCollision/NarrowPhaseCollision/btManifoldPoint.h.i"
// %include "BulletDynamics/ConstraintSolver/btConstraintSolver.h.i"
// %include "BulletDynamics/ConstraintSolver/btSequentialImpulseConstraintSolver.h.i"
// %include "BulletCollision/NarrowPhaseCollision/btPersistentManifold.h.i"
// %include "LinearMath/btIDebugDraw.h.i"
// %include "LinearMath/btCpuFeatureUtility.h.i"
// %include "LinearMath/btStackAlloc.h.i"
// %include "BulletDynamics/ConstraintSolver/btSequentialImpulseConstraintSolver.cpp.i"
// %include "BulletDynamics/ConstraintSolver/btContactConstraint.h.i"
// %include "BulletDynamics/ConstraintSolver/btContactConstraint.cpp.i"
// %include "BulletCollision/NarrowPhaseCollision/btPersistentManifold.cpp.i"
// %include "BulletCollision/NarrowPhaseCollision/btDiscreteCollisionDetectorInterface.h.i"
// %include "BulletCollision/CollisionDispatch/btManifoldResult.h.i"
// %include "BulletCollision/CollisionDispatch/btManifoldResult.cpp.i"
// %include "BulletCollision/BroadphaseCollision/btDispatcher.h.i"
// %include "BulletCollision/CollisionDispatch/btCollisionCreateFunc.h.i"
// %include "BulletCollision/CollisionDispatch/btCollisionDispatcher.h.i"
// %include "BulletCollision/BroadphaseCollision/btCollisionAlgorithm.h.i"
// %include "BulletCollision/BroadphaseCollision/btBroadphaseInterface.h.i"
// %include "BulletCollision/BroadphaseCollision/btOverlappingPairCallback.h.i"
// %include "BulletCollision/BroadphaseCollision/btOverlappingPairCache.h.i"
// %include "LinearMath/btThreads.h.i"
// %include "LinearMath/btPoolAllocator.h.i"
// %include "BulletCollision/CollisionDispatch/btCollisionConfiguration.h.i"
// %include "BulletCollision/CollisionDispatch/btCollisionDispatcher.cpp.i"
// %include "BulletCollision/BroadphaseCollision/btOverlappingPairCache.cpp.i"
// %include "BulletCollision/CollisionDispatch/btCollisionWorld.h.i"
// %include "BulletCollision/NarrowPhaseCollision/btVoronoiSimplexSolver.h.i"
// %include "BulletCollision/NarrowPhaseCollision/btSimplexSolverInterface.h.i"
// %include "BulletCollision/NarrowPhaseCollision/btConvexPenetrationDepthSolver.h.i"
// %include "BulletCollision/NarrowPhaseCollision/btGjkEpaPenetrationDepthSolver.h.i"
// %include "BulletCollision/CollisionShapes/btTriangleMeshShape.h.i"
// %include "BulletCollision/BroadphaseCollision/btQuantizedBvh.h.i"
// %include "BulletCollision/CollisionShapes/btOptimizedBvh.h.i"
// %include "BulletCollision/CollisionShapes/btTriangleInfoMap.h.i"
// %include "BulletCollision/CollisionShapes/btBvhTriangleMeshShape.h.i"
// %include "BulletCollision/CollisionShapes/btScaledBvhTriangleMeshShape.h.i"
// %include "BulletCollision/CollisionShapes/btHeightfieldTerrainShape.h.i"
// %include "BulletCollision/NarrowPhaseCollision/btRaycastCallback.h.i"
// %include "BulletCollision/CollisionShapes/btCompoundShape.h.i"
// %include "BulletCollision/NarrowPhaseCollision/btConvexCast.h.i"
// %include "BulletCollision/NarrowPhaseCollision/btSubSimplexConvexCast.h.i"
// %include "BulletCollision/NarrowPhaseCollision/btGjkConvexCast.h.i"
// %include "BulletCollision/NarrowPhaseCollision/btContinuousConvexCollision.h.i"
// %include "BulletCollision/BroadphaseCollision/btDbvt.h.i"
// %include "BulletCollision/BroadphaseCollision/btSimpleBroadphase.h.i"
// %include "BulletCollision/CollisionShapes/btConeShape.h.i"
// %include "BulletCollision/CollisionShapes/btConvexTriangleMeshShape.h.i"
// %include "BulletCollision/CollisionShapes/btCylinderShape.h.i"
// %include "BulletCollision/CollisionShapes/btMultiSphereShape.h.i"
// %include "BulletCollision/CollisionShapes/btStaticPlaneShape.h.i"
// %include "BulletCollision/CollisionDispatch/btCollisionWorld.cpp.i"
// %include "BulletDynamics/Dynamics/btDynamicsWorld.h.i"
// %include "BulletDynamics/Dynamics/btDiscreteDynamicsWorld.h.i"
// %include "BulletCollision/CollisionDispatch/btUnionFind.h.i"
// %include "BulletCollision/CollisionDispatch/btSimulationIslandManager.h.i"
// %include "BulletDynamics/ConstraintSolver/btPoint2PointConstraint.h.i"
// %include "BulletDynamics/ConstraintSolver/btHingeConstraint.h.i"
// %include "BulletDynamics/ConstraintSolver/btConeTwistConstraint.h.i"
// %include "BulletDynamics/ConstraintSolver/btGeneric6DofConstraint.h.i"
// %include "BulletDynamics/ConstraintSolver/btGeneric6DofSpring2Constraint.h.i"
// %include "BulletDynamics/ConstraintSolver/btSliderConstraint.h.i"
// %include "BulletDynamics/Dynamics/btActionInterface.h.i"
// %include "BulletDynamics/Dynamics/btDiscreteDynamicsWorld.cpp.i"
// %include "LinearMath/btThreads.cpp.i"
// %include "BulletCollision/CollisionShapes/btTriangleCallback.cpp.i"
// %include "BulletCollision/BroadphaseCollision/btDispatcher.cpp.i"
// %include "BulletCollision/CollisionShapes/btScaledBvhTriangleMeshShape.cpp.i"
// %include "BulletCollision/NarrowPhaseCollision/btPointCollector.h.i"
// %include "BulletCollision/NarrowPhaseCollision/btSubSimplexConvexCast.cpp.i"
// %include "BulletCollision/NarrowPhaseCollision/btGjkPairDetector.h.i"
// %include "BulletCollision/NarrowPhaseCollision/btGjkPairDetector.cpp.i"
// %include "BulletCollision/CollisionShapes/btCylinderShape.cpp.i"
// %include "BulletCollision/CollisionShapes/btConvexTriangleMeshShape.cpp.i"
// %include "BulletCollision/NarrowPhaseCollision/btGjkEpa2.h.i"
// %include "BulletCollision/NarrowPhaseCollision/btGjkEpa2.cpp.i"
// %include "BulletCollision/CollisionShapes/btMultiSphereShape.cpp.i"
// %include "BulletDynamics/ConstraintSolver/btPoint2PointConstraint.cpp.i"
// %include "BulletDynamics/ConstraintSolver/btHingeConstraint.cpp.i"
// %include "BulletCollision/CollisionShapes/btStaticPlaneShape.cpp.i"
// %include "BulletDynamics/ConstraintSolver/btGeneric6DofSpring2Constraint.cpp.i"
// %include "BulletDynamics/ConstraintSolver/btSliderConstraint.cpp.i"
// %include "BulletCollision/BroadphaseCollision/btQuantizedBvh.cpp.i"
// %include "BulletCollision/CollisionShapes/btOptimizedBvh.cpp.i"
// %include "LinearMath/btSerializer.cpp.i"
// %include "LinearMath/btSerializer64.cpp.i"
// %include "BulletCollision/CollisionShapes/btBvhTriangleMeshShape.cpp.i"
// %include "BulletCollision/CollisionDispatch/btUnionFind.cpp.i"
// %include "BulletCollision/CollisionDispatch/btSimulationIslandManager.cpp.i"
// %include "BulletDynamics/ConstraintSolver/btGeneric6DofConstraint.cpp.i"
// %include "BulletCollision/NarrowPhaseCollision/btConvexCast.cpp.i"
// %include "BulletCollision/CollisionShapes/btConcaveShape.cpp.i"
// %include "BulletCollision/CollisionShapes/btHeightfieldTerrainShape.cpp.i"
// %include "BulletDynamics/ConstraintSolver/btConeTwistConstraint.cpp.i"
// %include "BulletCollision/NarrowPhaseCollision/btVoronoiSimplexSolver.cpp.i"
// %include "BulletCollision/NarrowPhaseCollision/btGjkConvexCast.cpp.i"
// %include "BulletCollision/NarrowPhaseCollision/btGjkEpaPenetrationDepthSolver.cpp.i"
// %include "BulletCollision/CollisionShapes/btConeShape.cpp.i"
// %include "BulletCollision/BroadphaseCollision/btDbvt.cpp.i"
// %include "BulletCollision/CollisionShapes/btCompoundShape.cpp.i"
// %include "BulletCollision/CollisionShapes/btTriangleShape.h.i"
// %include "BulletCollision/NarrowPhaseCollision/btRaycastCallback.cpp.i"
// %include "BulletCollision/NarrowPhaseCollision/btContinuousConvexCollision.cpp.i"
// %include "BulletCollision/CollisionShapes/btTriangleMeshShape.cpp.i"
// %include "BulletCollision/BroadphaseCollision/btSimpleBroadphase.cpp.i"
