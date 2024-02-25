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

%ignore Operator_Sm__Se_;

%include "LinearMath/btScalar.h.i"
%include "LinearMath/btMinMax.h.i"
%include "LinearMath/btAlignedAllocator.h.i"
%include "LinearMath/btAlignedAllocator.cpp.i"
%include "LinearMath/btVector3.h.i"
%include "LinearMath/btQuadWord.h.i"
%include "LinearMath/btQuaternion.h.i"
%include "LinearMath/btMatrix3x3.h.i"
%include "LinearMath/btTransform.h.i"
%include "LinearMath/btMotionState.h.i"
%include "LinearMath/btDefaultMotionState.h.i"

%include "LinearMath/btAlignedObjectArray.h.i"
%include "LinearMath/btHashMap.h.i"
%include "LinearMath/btSerializer.h.i"
%include "LinearMath/btSerializer.cpp.i"
%include "LinearMath/btSerializer64.cpp.i"
%include "LinearMath/btAabbUtil2.h.i"
%include "LinearMath/btConvexHullComputer.h.i"
%include "LinearMath/btConvexHullComputer.cpp.i"

%include "BulletCollision/BroadphaseCollision/btBroadphaseProxy.h.i"
%include "BulletCollision/CollisionShapes/btCollisionMargin.h.i"
%include "BulletCollision/CollisionShapes/btCollisionShape.h.i"
%include "BulletCollision/CollisionShapes/btCollisionShape.cpp.i"
%include "BulletCollision/CollisionShapes/btConvexShape.h.i"
%include "BulletCollision/CollisionShapes/btConvexShape.cpp.i"
%include "BulletCollision/CollisionShapes/btConvexInternalShape.h.i"
%include "BulletCollision/CollisionShapes/btConvexInternalShape.cpp.i"
%include "BulletCollision/CollisionShapes/btPolyhedralConvexShape.h.i"
%include "BulletCollision/CollisionShapes/btPolyhedralConvexShape.cpp.i"
%include "BulletCollision/CollisionShapes/btConvexPolyhedron.h.i"
%include "BulletCollision/CollisionShapes/btConvexPolyhedron.cpp.i"

// %include "BulletCollision/CollisionShapes/btTriangleShape.h.i"
// %include "BulletCollision/CollisionShapes/btBoxShape.h.i"

%include "LinearMath/btGeometryUtil.h.i"
%include "LinearMath/btGeometryUtil.cpp.i"

// %include "BulletCollision/CollisionShapes/btBoxShape.cpp.i"

// %include "BulletCollision/CollisionShapes/btSphereShape.h.i"
// %include "BulletCollision/CollisionShapes/btConvexPolyhedron.h.i"
// %include "BulletCollision/CollisionShapes/btConvexHullShape.h.i"

// %include "BulletCollision/CollisionShapes/btCylinderShape.h.i"
// %include "BulletCollision/CollisionShapes/btCylinderShapeX.h.i"
// %include "BulletCollision/CollisionShapes/btCylinderShapeZ.h.i"
// %include "BulletCollision/CollisionShapes/btConeShape.h.i"
// %include "BulletCollision/CollisionShapes/btConeShapeX.h.i"
// %include "BulletCollision/CollisionShapes/btConeShapeZ.h.i"
// %include "BulletCollision/CollisionShapes/btCapsuleShape.h.i"
// %include "BulletCollision/CollisionShapes/btCapsuleShapeX.h.i"
// %include "BulletCollision/CollisionShapes/btCapsuleShapeZ.h.i"
// %include "BulletCollision/CollisionShapes/btConvexPointCloudShape.h.i"
// %include "BulletCollision/CollisionShapes/btConvexShape.cpp.i"
// %include "BulletDynamics/Dynamics/btDiscreteDynamicsWorld.h"

