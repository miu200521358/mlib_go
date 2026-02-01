////// LinearMath/btMDebugDraw.h ----------------

%module(directors="1") bt
%feature("director") btMDebugDrawLiner;

%include "LinearMath/btMDebugDraw.h"

%{
#include "LinearMath/btMDebugDraw.h"
%}

// class directorBtMDebugDrawLiner : public btMDebugDrawLiner {
// public:
//     virtual ~directorBtMDebugDrawLiner() {}
//     virtual void drawLine(const btVector3& from, const btVector3& to, const btVector3& color) {}
// };

// void drawLine(directorBtMDebugDrawLiner *director, const btVector3& from, const btVector3& to, const btVector3& color) {
//     director->drawLine(from, to, color);
// }