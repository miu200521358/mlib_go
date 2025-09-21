package mbt

import (
	"math"

	"github.com/go-gl/gl/v4.4-core/gl"
	"github.com/go-gl/mathgl/mgl32"

	"github.com/miu200521358/mlib_go/pkg/domain/mmath"
	"github.com/miu200521358/mlib_go/pkg/domain/physics"
	"github.com/miu200521358/mlib_go/pkg/domain/rendering"
	"github.com/miu200521358/mlib_go/pkg/infrastructure/bt"
	"github.com/miu200521358/mlib_go/pkg/infrastructure/mgl"
)

const (
	highlightColorR = 1.0
	highlightColorG = 1.0
	highlightColorB = 0.0
	highlightColorA = 0.35
	epsilon         = 1e-6
)

func (mp *MPhysics) UpdateDebugHover(rayFrom, rayTo *mmath.MVec3, enable bool) {
	if !enable || rayFrom == nil || rayTo == nil {
		mp.clearDebugHover()
		return
	}

	btRayFrom := NewBulletFromVec(rayFrom)
	btRayTo := NewBulletFromVec(rayTo)
	defer bt.DeleteBtVector3(btRayFrom)
	defer bt.DeleteBtVector3(btRayTo)

	bestDistance := math.MaxFloat64
	var bestRigid *rigidBodyValue
	var bestHit *mmath.MVec3

	for _, bodies := range mp.rigidBodies {
		for _, rb := range bodies {
			if rb == nil || rb.btRigidBody == nil {
				continue
			}

			hitPoint, ok := mp.intersectRigidBody(rayFrom, rayTo, rb, btRayFrom, btRayTo)
			if !ok || hitPoint == nil {
				continue
			}

			distance := rayFrom.Distance(hitPoint)
			if distance < bestDistance {
				bestDistance = distance
				bestRigid = rb
				bestHit = hitPoint
			}
		}
	}

	if bestRigid == nil {
		mp.clearDebugHover()
		return
	}

	mp.debugHover = &physics.DebugRigidBodyHover{
		RigidBody: bestRigid.pmxRigidBody,
		HitPoint:  bestHit,
	}
	mp.debugHoverRigid = bestRigid
	mp.debugHoverDistance = bestDistance
	mp.rebuildHighlightVertices(bestRigid)
}

func (mp *MPhysics) DebugHoverInfo() *physics.DebugRigidBodyHover {
	return mp.debugHover
}

func (mp *MPhysics) DrawDebugHighlight(shader rendering.IShader, isDrawRigidBodyFront bool) {
	if len(mp.highlightVertices) == 0 || mp.debugHoverRigid == nil {
		return
	}

	program := shader.Program(rendering.ProgramTypePhysics)
	gl.UseProgram(program)

	mp.liner.configureDepthTest(isDrawRigidBodyFront)

	if mp.highlightBuffer == nil {
		mp.highlightBuffer = mgl.NewBufferFactory().CreateDebugBuffer(gl.Ptr(&mp.highlightVertices[0]), len(mp.highlightVertices))
	}

	mp.highlightBuffer.Bind()
	mp.highlightBuffer.UpdateDebugBuffer(mp.highlightVertices)

	gl.DrawArrays(gl.TRIANGLES, 0, int32(len(mp.highlightVertices)/7))

	mp.highlightBuffer.Unbind()
	mp.liner.restoreDepthTest(isDrawRigidBodyFront)
	gl.UseProgram(0)
}

func (mp *MPhysics) clearDebugHover() {
	mp.debugHover = nil
	mp.debugHoverRigid = nil
	mp.debugHoverDistance = 0
	if mp.highlightVertices != nil {
		mp.highlightVertices = mp.highlightVertices[:0]
	}
}

func (mp *MPhysics) intersectRigidBody(rayFrom, rayTo *mmath.MVec3, rb *rigidBodyValue, btRayFrom, btRayTo bt.BtVector3) (*mmath.MVec3, bool) {
	worldTransformIface := rb.btRigidBody.GetWorldTransform()
	transform, ok := worldTransformIface.(bt.BtTransform)
	if !ok {
		return nil, false
	}

	localFromBt := transform.InvXform(btRayFrom)
	defer bt.DeleteBtVector3(localFromBt)
	localToBt := transform.InvXform(btRayTo)
	defer bt.DeleteBtVector3(localToBt)

	localFrom := vec3FromBt(localFromBt)
	localTo := vec3FromBt(localToBt)
	dir := localTo.sub(localFrom)

	collisionShapeIface := rb.btRigidBody.GetCollisionShape()
	shape, ok := collisionShapeIface.(bt.BtCollisionShape)
	if !ok {
		return nil, false
	}

	var t float64
	var hit bool

	switch s := shape.(type) {
	case bt.BtBoxShape:
		halfExtentsBt := s.GetHalfExtentsWithMargin()
		hx := math.Abs(float64(halfExtentsBt.GetX()))
		hy := math.Abs(float64(halfExtentsBt.GetY()))
		hz := math.Abs(float64(halfExtentsBt.GetZ()))
		bt.DeleteBtVector3(halfExtentsBt)
		min := vec3{-hx, -hy, -hz}
		max := vec3{hx, hy, hz}
		t, hit = rayAabbIntersection(localFrom, dir, min, max)
	case bt.BtSphereShape:
		radius := float64(s.GetRadius())
		t, hit = raySphereIntersection(localFrom, dir, radius)
	case bt.BtCapsuleShape:
		radius := float64(s.GetRadius())
		halfHeight := float64(s.GetHalfHeight())
		t, hit = rayCapsuleIntersection(localFrom, dir, halfHeight, radius)
	default:
		return nil, false
	}

	if !hit || t < 0 || t > 1 {
		return nil, false
	}

	hitPoint := rayFrom.Added(rayTo.Subed(rayFrom).MuledScalar(t))
	return hitPoint, true
}

func (mp *MPhysics) rebuildHighlightVertices(rb *rigidBodyValue) {
	mp.highlightVertices = mp.highlightVertices[:0]
	if rb == nil || rb.btRigidBody == nil {
		return
	}

	transformIface := rb.btRigidBody.GetWorldTransform()
	transform, ok := transformIface.(bt.BtTransform)
	if !ok {
		return
	}

	mat := mgl32.Mat4{}
	transform.GetOpenGLMatrix(&mat[0])
	worldMat := NewMMat4ByMgl(&mat)

	shapeIface := rb.btRigidBody.GetCollisionShape()
	shape, ok := shapeIface.(bt.BtCollisionShape)
	if !ok {
		return
	}

	switch s := shape.(type) {
	case bt.BtBoxShape:
		mp.appendBoxHighlight(worldMat, s)
	case bt.BtSphereShape:
		mp.appendSphereHighlight(worldMat, float64(s.GetRadius()))
	case bt.BtCapsuleShape:
		mp.appendCapsuleHighlight(worldMat, float64(s.GetRadius()), float64(s.GetHalfHeight()))
	}
}

func (mp *MPhysics) appendBoxHighlight(world *mmath.MMat4, shape bt.BtBoxShape) {
	halfExtentsBt := shape.GetHalfExtentsWithMargin()
	hx := math.Abs(float64(halfExtentsBt.GetX()))
	hy := math.Abs(float64(halfExtentsBt.GetY()))
	hz := math.Abs(float64(halfExtentsBt.GetZ()))
	bt.DeleteBtVector3(halfExtentsBt)

	corners := []vec3{
		{-hx, -hy, -hz},
		{hx, -hy, -hz},
		{hx, hy, -hz},
		{-hx, hy, -hz},
		{-hx, -hy, hz},
		{hx, -hy, hz},
		{hx, hy, hz},
		{-hx, hy, hz},
	}

	indices := []int{
		0, 1, 2, 0, 2, 3,
		4, 5, 6, 4, 6, 7,
		0, 1, 5, 0, 5, 4,
		2, 3, 7, 2, 7, 6,
		1, 2, 6, 1, 6, 5,
		0, 3, 7, 0, 7, 4,
	}

	mp.appendVertices(world, corners, indices)
}

func (mp *MPhysics) appendSphereHighlight(world *mmath.MMat4, radius float64) {
	latSegments := 12
	lonSegments := 24

	var vertices []vec3
	var indices []int

	for lat := 0; lat < latSegments; lat++ {
		theta1 := float64(lat) / float64(latSegments) * math.Pi
		theta2 := float64(lat+1) / float64(latSegments) * math.Pi

		y1 := radius * math.Cos(theta1)
		y2 := radius * math.Cos(theta2)
		r1 := radius * math.Sin(theta1)
		r2 := radius * math.Sin(theta2)

		for lon := 0; lon < lonSegments; lon++ {
			phi1 := float64(lon) / float64(lonSegments) * 2 * math.Pi
			phi2 := float64(lon+1) / float64(lonSegments) * 2 * math.Pi

			p1 := vec3{r1 * math.Cos(phi1), y1, r1 * math.Sin(phi1)}
			p2 := vec3{r2 * math.Cos(phi1), y2, r2 * math.Sin(phi1)}
			p3 := vec3{r2 * math.Cos(phi2), y2, r2 * math.Sin(phi2)}
			p4 := vec3{r1 * math.Cos(phi2), y1, r1 * math.Sin(phi2)}

			base := len(vertices)
			vertices = append(vertices, p1, p2, p3, p4)
			indices = append(indices,
				base, base+1, base+2,
				base, base+2, base+3,
			)
		}
	}

	mp.appendVertices(world, vertices, indices)
}

func (mp *MPhysics) appendCapsuleHighlight(world *mmath.MMat4, radius, halfHeight float64) {
	segments := 16
	capSegments := segments / 2

	var vertices []vec3
	var indices []int

	// Cylinder
	for i := 0; i < segments; i++ {
		phi1 := float64(i) / float64(segments) * 2 * math.Pi
		phi2 := float64(i+1) / float64(segments) * 2 * math.Pi

		top1 := vec3{radius * math.Cos(phi1), halfHeight, radius * math.Sin(phi1)}
		top2 := vec3{radius * math.Cos(phi2), halfHeight, radius * math.Sin(phi2)}
		bottom1 := vec3{radius * math.Cos(phi1), -halfHeight, radius * math.Sin(phi1)}
		bottom2 := vec3{radius * math.Cos(phi2), -halfHeight, radius * math.Sin(phi2)}

		base := len(vertices)
		vertices = append(vertices, top1, bottom1, top2, bottom2)
		indices = append(indices,
			base, base+1, base+2,
			base+1, base+3, base+2,
		)
	}

	// Top hemisphere
	centerTop := halfHeight
	for lat := 0; lat < capSegments; lat++ {
		theta1 := float64(lat) / float64(capSegments) * (math.Pi / 2)
		theta2 := float64(lat+1) / float64(capSegments) * (math.Pi / 2)

		y1 := centerTop + radius*math.Cos(theta1)
		y2 := centerTop + radius*math.Cos(theta2)
		r1 := radius * math.Sin(theta1)
		r2 := radius * math.Sin(theta2)

		for lon := 0; lon < segments; lon++ {
			phi1 := float64(lon) / float64(segments) * 2 * math.Pi
			phi2 := float64(lon+1) / float64(segments) * 2 * math.Pi

			p1 := vec3{r1 * math.Cos(phi1), y1, r1 * math.Sin(phi1)}
			p2 := vec3{r2 * math.Cos(phi1), y2, r2 * math.Sin(phi1)}
			p3 := vec3{r2 * math.Cos(phi2), y2, r2 * math.Sin(phi2)}
			p4 := vec3{r1 * math.Cos(phi2), y1, r1 * math.Sin(phi2)}

			base := len(vertices)
			vertices = append(vertices, p1, p2, p3, p4)
			indices = append(indices,
				base, base+1, base+2,
				base, base+2, base+3,
			)
		}
	}

	// Bottom hemisphere
	centerBottom := -halfHeight
	for lat := 0; lat < capSegments; lat++ {
		theta1 := float64(lat)/float64(capSegments)*(math.Pi/2) + math.Pi/2
		theta2 := float64(lat+1)/float64(capSegments)*(math.Pi/2) + math.Pi/2

		y1 := centerBottom + radius*math.Cos(theta1)
		y2 := centerBottom + radius*math.Cos(theta2)
		r1 := radius * math.Sin(theta1)
		r2 := radius * math.Sin(theta2)

		for lon := 0; lon < segments; lon++ {
			phi1 := float64(lon) / float64(segments) * 2 * math.Pi
			phi2 := float64(lon+1) / float64(segments) * 2 * math.Pi

			p1 := vec3{r1 * math.Cos(phi1), y1, r1 * math.Sin(phi1)}
			p2 := vec3{r2 * math.Cos(phi1), y2, r2 * math.Sin(phi1)}
			p3 := vec3{r2 * math.Cos(phi2), y2, r2 * math.Sin(phi2)}
			p4 := vec3{r1 * math.Cos(phi2), y1, r1 * math.Sin(phi2)}

			base := len(vertices)
			vertices = append(vertices, p1, p2, p3, p4)
			indices = append(indices,
				base, base+1, base+2,
				base, base+2, base+3,
			)
		}
	}

	mp.appendVertices(world, vertices, indices)
}

func (mp *MPhysics) appendVertices(world *mmath.MMat4, vertices []vec3, indices []int) {
	color := []float32{float32(highlightColorR), float32(highlightColorG), float32(highlightColorB), float32(highlightColorA)}
	for _, idx := range indices {
		if idx < 0 || idx >= len(vertices) {
			continue
		}
		local := vertices[idx]
		worldPos := world.MulVec3(&mmath.MVec3{X: local.x, Y: local.y, Z: local.z})
		mp.highlightVertices = append(mp.highlightVertices,
			float32(worldPos.X), float32(worldPos.Y), float32(worldPos.Z),
			color[0], color[1], color[2], color[3],
		)
	}
}

func raySphereIntersection(origin, dir vec3, radius float64) (float64, bool) {
	a := dir.dot(dir)
	if math.Abs(a) < epsilon {
		return 0, false
	}
	b := 2 * origin.dot(dir)
	c := origin.dot(origin) - radius*radius
	discriminant := b*b - 4*a*c
	if discriminant < 0 {
		return 0, false
	}
	sqrtDisc := math.Sqrt(discriminant)
	invDenom := 1 / (2 * a)
	t0 := (-b - sqrtDisc) * invDenom
	t1 := (-b + sqrtDisc) * invDenom
	if t0 >= 0 && t0 <= 1 {
		return t0, true
	}
	if t1 >= 0 && t1 <= 1 {
		return t1, true
	}
	return 0, false
}

func rayAabbIntersection(origin, dir, min, max vec3) (float64, bool) {
	tMin := 0.0
	tMax := 1.0

	for axis := 0; axis < 3; axis++ {
		var o, d, minVal, maxVal float64
		switch axis {
		case 0:
			o, d, minVal, maxVal = origin.x, dir.x, min.x, max.x
		case 1:
			o, d, minVal, maxVal = origin.y, dir.y, min.y, max.y
		case 2:
			o, d, minVal, maxVal = origin.z, dir.z, min.z, max.z
		}

		if math.Abs(d) < epsilon {
			if o < minVal || o > maxVal {
				return 0, false
			}
			continue
		}

		invD := 1 / d
		t1 := (minVal - o) * invD
		t2 := (maxVal - o) * invD
		if t1 > t2 {
			t1, t2 = t2, t1
		}

		if t1 > tMin {
			tMin = t1
		}
		if t2 < tMax {
			tMax = t2
		}
		if tMax < tMin {
			return 0, false
		}
	}

	if tMin >= 0 && tMin <= 1 {
		return tMin, true
	}
	if tMax >= 0 && tMax <= 1 {
		return tMax, true
	}
	return 0, false
}

func rayCapsuleIntersection(origin, dir vec3, halfHeight, radius float64) (float64, bool) {
	minT := math.MaxFloat64
	hit := false

	a := dir.x*dir.x + dir.z*dir.z
	if math.Abs(a) > epsilon {
		b := 2 * (origin.x*dir.x + origin.z*dir.z)
		c := origin.x*origin.x + origin.z*origin.z - radius*radius
		discriminant := b*b - 4*a*c
		if discriminant >= 0 {
			sqrtDisc := math.Sqrt(discriminant)
			invDenom := 1 / (2 * a)
			t0 := (-b - sqrtDisc) * invDenom
			t1 := (-b + sqrtDisc) * invDenom
			if t0 > t1 {
				t0, t1 = t1, t0
			}
			if t0 >= 0 {
				y := origin.y + dir.y*t0
				if t0 <= 1 && y >= -halfHeight && y <= halfHeight {
					minT = t0
					hit = true
				}
			}
			if t1 >= 0 && (!hit || t1 < minT) {
				y := origin.y + dir.y*t1
				if t1 <= 1 && y >= -halfHeight && y <= halfHeight {
					minT = t1
					hit = true
				}
			}
		}
	} else if origin.x*origin.x+origin.z*origin.z <= radius*radius {
		if dir.y > epsilon {
			t := (-halfHeight - origin.y) / dir.y
			if t >= 0 && t <= 1 {
				minT = t
				hit = true
			}
		} else if dir.y < -epsilon {
			t := (halfHeight - origin.y) / dir.y
			if t >= 0 && t <= 1 {
				minT = t
				hit = true
			}
		}
	}

	if t, ok := raySphereIntersection(origin.sub(vec3{0, halfHeight, 0}), dir, radius); ok {
		if t >= 0 && t <= 1 && (!hit || t < minT) {
			minT = t
			hit = true
		}
	}
	if t, ok := raySphereIntersection(origin.sub(vec3{0, -halfHeight, 0}), dir, radius); ok {
		if t >= 0 && t <= 1 && (!hit || t < minT) {
			minT = t
			hit = true
		}
	}

	return minT, hit
}

type vec3 struct {
	x, y, z float64
}

func vec3FromBt(v bt.BtVector3) vec3 {
	return vec3{float64(v.GetX()), float64(v.GetY()), float64(v.GetZ())}
}

func (v vec3) sub(o vec3) vec3 {
	return vec3{v.x - o.x, v.y - o.y, v.z - o.z}
}

func (v vec3) dot(o vec3) float64 {
	return v.x*o.x + v.y*o.y + v.z*o.z
}
