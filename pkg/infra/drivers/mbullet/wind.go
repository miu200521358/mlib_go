//go:build windows
// +build windows

// 指示: miu200521358
package mbullet

import (
	"math"

	"github.com/miu200521358/mlib_go/pkg/domain/delta"
	"github.com/miu200521358/mlib_go/pkg/domain/mmath"
	"github.com/miu200521358/mlib_go/pkg/domain/model"
	"github.com/miu200521358/mlib_go/pkg/infra/drivers/mbullet/bt"
)

// EnableWind は風の有効/無効を切り替えます。
func (mp *PhysicsEngine) EnableWind(enable bool) {
	mp.windCfg.Enabled = enable
}

// SetWind は風向き・風速・ランダム性を設定します。
func (mp *PhysicsEngine) SetWind(direction *mmath.Vec3, speed float32, randomness float32) {
	if direction != nil {
		mp.windCfg.Direction = *direction
	}
	mp.windCfg.Speed = speed
	if randomness < 0 {
		randomness = 0
	}
	mp.windCfg.Randomness = randomness
}

// SetWindAdvanced は風の詳細パラメータを設定します。
func (mp *PhysicsEngine) SetWindAdvanced(dragCoeff, liftCoeff, turbulenceFreqHz float32) {
	if dragCoeff >= 0 {
		mp.windCfg.DragCoeff = dragCoeff
	}
	if liftCoeff >= 0 {
		mp.windCfg.LiftCoeff = liftCoeff
	}
	if turbulenceFreqHz >= 0 {
		mp.windCfg.TurbulenceFreqHz = turbulenceFreqHz
	}
}

// applyWindForces は風の力（抵抗 + 簡易揚力）を動的剛体に付与します。
func (mp *PhysicsEngine) applyWindForces(dt float32) {
	if !mp.windCfg.Enabled || mp.windCfg.Speed == 0 {
		return
	}

	mp.simTimeAcc += dt

	r := float64(mmath.Clamped(float64(mp.windCfg.Randomness), 0, 1))
	f := math.Max(0.0001, float64(mp.windCfg.TurbulenceFreqHz))
	t := float64(mp.simTimeAcc)
	gust := 1.0 + r*(0.6*math.Sin(2*math.Pi*f*t)+0.4*math.Sin(2*math.Pi*1.73*f*t+0.9))

	dir := mp.windCfg.Direction.Normalized()
	windSpeed := float64(mp.windCfg.Speed) * gust
	windVecMmd := dir.MuledScalar(windSpeed)
	windVecBt := newBulletFromVec(windVecMmd)
	defer bt.DeleteBtVector3(windVecBt)

	windX := float64(windVecBt.GetX())
	windY := float64(windVecBt.GetY())
	windZ := float64(windVecBt.GetZ())

	for _, bodies := range mp.rigidBodies {
		if bodies == nil {
			continue
		}
		for _, rb := range bodies {
			if rb == nil || rb.BtRigidBody == nil || rb.RigidBody == nil {
				continue
			}
			if rb.RigidBody.PhysicsType == model.PHYSICS_TYPE_STATIC {
				continue
			}

			v := rb.BtRigidBody.GetLinearVelocity()
			vx := float64(v.GetX())
			vy := float64(v.GetY())
			vz := float64(v.GetZ())

			relX := vx - windX
			relY := vy - windY
			relZ := vz - windZ
			speed2 := relX*relX + relY*relY + relZ*relZ
			if speed2 < 1.0e-12 {
				continue
			}

			speed := math.Sqrt(speed2)
			area := mp.approxCrossSectionArea(rb, dir)

			invSpeed := 1.0 / speed
			nrx := relX * invSpeed
			nry := relY * invSpeed
			nrz := relZ * invSpeed

			kd := float64(mp.windCfg.DragCoeff)
			magD := kd * float64(area) * speed2
			fDx := float32(-magD * nrx)
			fDy := float32(-magD * nry)
			fDz := float32(-magD * nrz)

			kl := float64(mp.windCfg.LiftCoeff)
			fLx, fLy, fLz := float32(0), float32(0), float32(0)
			if kl > 0 {
				ux, uy, uz := 0.0, 1.0, 0.0
				dotUn := ux*nrx + uy*nry + uz*nrz
				lx := ux - dotUn*nrx
				ly := uy - dotUn*nry
				lz := uz - dotUn*nrz
				l2 := lx*lx + ly*ly + lz*lz
				if l2 > 1.0e-8 {
					invL := 1.0 / math.Sqrt(l2)
					lx *= invL
					ly *= invL
					lz *= invL
					magL := kl * float64(area) * speed2
					fLx = float32(magL * lx)
					fLy = float32(magL * ly)
					fLz = float32(magL * lz)
				}
			}

			fTx := fDx + fLx
			fTy := fDy + fLy
			fTz := fDz + fLz

			maxA := float64(mp.windCfg.MaxAcceleration)
			if maxA > 0 {
				mass := float64(rb.BtRigidBody.GetMass())
				if mass > 0 {
					f2 := float64(fTx)*float64(fTx) + float64(fTy)*float64(fTy) + float64(fTz)*float64(fTz)
					if f2 > 0 {
						fmag := math.Sqrt(f2)
						a := fmag / mass
						if a > maxA {
							scale := float32((mass * maxA) / fmag)
							fTx *= scale
							fTy *= scale
							fTz *= scale
						}
					}
				}
			}

			fTotal := bt.NewBtVector3(fTx, fTy, fTz)
			rb.BtRigidBody.ApplyCentralForce(fTotal)
			bt.DeleteBtVector3(fTotal)

			rb.BtRigidBody.Activate(true)
		}
	}
}

// approxCrossSectionArea は風向きに対する見かけの断面積を近似計算します。
func (mp *PhysicsEngine) approxCrossSectionArea(r *RigidBodyValue, dir mmath.Vec3) float32 {
	if r == nil || r.RigidBody == nil {
		return 1.0
	}
	if dir.IsZero() {
		return 1.0
	}
	size := r.RigidBody.Size
	if r.HasAppliedParams {
		size = r.AppliedSize
	}

	d := dir.Normalized()
	absX := math.Abs(d.X)
	absY := math.Abs(d.Y)
	absZ := math.Abs(d.Z)

	switch r.RigidBody.Shape {
	case model.SHAPE_SPHERE:
		rads := float64(size.X)
		return float32(math.Pi * rads * rads)
	case model.SHAPE_BOX:
		wx := float64(2.0 * size.X)
		wy := float64(2.0 * size.Y)
		wz := float64(2.0 * size.Z)
		area := absX*(wy*wz) + absY*(wx*wz) + absZ*(wx*wy)
		return float32(area)
	case model.SHAPE_CAPSULE:
		rad := float64(size.X)
		h := float64(size.Y)
		aAxis := math.Pi * rad * rad
		aPerp := 2*rad*h + math.Pi*rad*rad
		w := absY
		return float32(w*aAxis + (1.0-w)*aPerp)
	default:
		rads := float64(size.X)
		return float32(math.Pi * rads * rads)
	}
}

// UpdatePhysicsSelectively は変更が必要な剛体・ジョイントのみを選択的に更新します。
func (mp *PhysicsEngine) UpdatePhysicsSelectively(
	modelIndex int,
	pmxModel *model.PmxModel,
	physicsDeltas *delta.PhysicsDeltas,
) {
	if physicsDeltas == nil {
		return
	}

	if physicsDeltas.RigidBodies != nil {
		mp.UpdateRigidBodiesSelectively(modelIndex, pmxModel, physicsDeltas.RigidBodies)
	}
	if physicsDeltas.Joints != nil {
		mp.UpdateJointsSelectively(modelIndex, pmxModel, physicsDeltas.Joints)
	}
}
