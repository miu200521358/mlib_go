// 指示: miu200521358
package vrm

import "github.com/miu200521358/mlib_go/pkg/domain/mmath"

// Vrm1SpringBone は VRM1 VRMC_springBone 拡張情報を表す。
type Vrm1SpringBone struct {
	SpecVersion    string
	Colliders      []Vrm1SpringCollider
	ColliderGroups []Vrm1SpringColliderGroup
	Springs        []Vrm1Spring
}

// NewVrm1SpringBone は Vrm1SpringBone を既定値で生成する。
func NewVrm1SpringBone() *Vrm1SpringBone {
	return &Vrm1SpringBone{
		Colliders:      make([]Vrm1SpringCollider, 0),
		ColliderGroups: make([]Vrm1SpringColliderGroup, 0),
		Springs:        make([]Vrm1Spring, 0),
	}
}

// Vrm1SpringCollider は VRM1 springBone collider を表す。
type Vrm1SpringCollider struct {
	Node     int
	Shape    Vrm1SpringColliderShape
	Extended *Vrm1SpringExtendedCollider
}

// Vrm1SpringColliderShape は VRM1 springBone collider shape を表す。
type Vrm1SpringColliderShape struct {
	Sphere  *Vrm1SpringColliderSphere
	Capsule *Vrm1SpringColliderCapsule
}

// Vrm1SpringColliderSphere は VRM1 springBone sphere collider を表す。
type Vrm1SpringColliderSphere struct {
	Offset mmath.Vec3
	Radius float64
}

// Vrm1SpringColliderCapsule は VRM1 springBone capsule collider を表す。
type Vrm1SpringColliderCapsule struct {
	Offset mmath.Vec3
	Radius float64
	Tail   mmath.Vec3
}

// Vrm1SpringExtendedCollider は VRM1 springBone extended collider を表す。
type Vrm1SpringExtendedCollider struct {
	SpecVersion string
	Shape       Vrm1SpringExtendedColliderShape
}

// Vrm1SpringExtendedColliderShape は VRM1 springBone extended collider shape を表す。
type Vrm1SpringExtendedColliderShape struct {
	Sphere  *Vrm1SpringExtendedSphereCollider
	Capsule *Vrm1SpringExtendedCapsuleCollider
	Plane   *Vrm1SpringExtendedPlaneCollider
}

// Vrm1SpringExtendedSphereCollider は VRM1 springBone inside sphere collider を表す。
type Vrm1SpringExtendedSphereCollider struct {
	Offset mmath.Vec3
	Radius float64
	Inside bool
}

// Vrm1SpringExtendedCapsuleCollider は VRM1 springBone inside capsule collider を表す。
type Vrm1SpringExtendedCapsuleCollider struct {
	Offset mmath.Vec3
	Radius float64
	Tail   mmath.Vec3
	Inside bool
}

// Vrm1SpringExtendedPlaneCollider は VRM1 springBone plane collider を表す。
type Vrm1SpringExtendedPlaneCollider struct {
	Offset mmath.Vec3
	Normal mmath.Vec3
}

// Vrm1SpringColliderGroup は VRM1 springBone colliderGroup を表す。
type Vrm1SpringColliderGroup struct {
	Name      string
	Colliders []int
}

// Vrm1Spring は VRM1 springBone spring を表す。
type Vrm1Spring struct {
	Name           string
	Joints         []Vrm1SpringJoint
	ColliderGroups []int
	Center         *int
}

// Vrm1SpringJoint は VRM1 springBone joint を表す。
type Vrm1SpringJoint struct {
	Node         int
	HitRadius    float64
	Stiffness    float64
	GravityPower float64
	GravityDir   mmath.Vec3
	DragForce    float64
}
