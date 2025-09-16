//go:build windows
// +build windows

package mbt

import (
	"math"

	"github.com/miu200521358/mlib_go/pkg/domain/delta"
	"github.com/miu200521358/mlib_go/pkg/domain/mmath"
	"github.com/miu200521358/mlib_go/pkg/domain/physics"
	"github.com/miu200521358/mlib_go/pkg/domain/pmx"
	"github.com/miu200521358/mlib_go/pkg/infrastructure/bt"
)

// MPhysics 物理エンジンの実装
type MPhysics struct {
	world       bt.BtDiscreteDynamicsWorld // ワールド
	drawer      bt.BtMDebugDraw            // デバッグビューワー
	liner       *mDebugDrawLiner           // ライナー
	config      physics.PhysicsConfig      // 設定パラメータ
	DeformSpf   float32                    // デフォームspf
	PhysicsSpf  float32                    // 物理spf
	joints      map[int][]*jointValue      // ジョイント
	rigidBodies map[int][]*rigidBodyValue  // 剛体

	// 風の設定
	windCfg    physics.WindConfig
	simTimeAcc float32 // 経過時間[秒]
}

// NewMPhysics は物理エンジンのインスタンスを生成します
func NewMPhysics(gravity *mmath.MVec3) physics.IPhysics {
	world := createWorld(gravity)

	// デフォルト設定
	physics := &MPhysics{
		world: world,
		config: physics.PhysicsConfig{
			FixedTimeStep: 1 / 60.0,
		},
		rigidBodies: make(map[int][]*rigidBodyValue),
		joints:      make(map[int][]*jointValue),

		// 風のデフォルト設定（無効）
		windCfg: physics.WindConfig{
			Enabled:          false,
			Direction:        &mmath.MVec3{X: 1, Y: 0, Z: 0},
			Speed:            0,
			Randomness:       0,
			TurbulenceFreqHz: 0.5,
			DragCoeff:        0.8,
			LiftCoeff:        0.2,
			MaxAcceleration:  80.0,
		},
		simTimeAcc: 0,
	}

	// デバッグビューワーの初期化
	physics.initDebugDrawer()

	return physics
}

// initDebugDrawer はデバッグ描画機能を初期化します
func (mp *MPhysics) initDebugDrawer() {
	liner := newMDebugDrawLiner()
	drawer := bt.NewBtMDebugDraw()
	drawer.SetLiner(liner)
	drawer.SetMDefaultColors(newConstBtMDefaultColors())
	mp.world.SetDebugDrawer(drawer)

	mp.drawer = drawer
	mp.liner = liner
}

// ResetWorld はワールドをリセットします
func (mp *MPhysics) ResetWorld(gravity *mmath.MVec3) {
	// ワールド削除
	bt.DeleteBtDynamicsWorld(mp.world)
	// ワールド作成
	world := createWorld(gravity)
	world.SetDebugDrawer(mp.drawer)
	mp.world = world
}

// AddModel はモデルを物理エンジンに追加します
func (mp *MPhysics) AddModel(modelIndex int, model *pmx.PmxModel) {
	// 根元から追加していく
	mp.initRigidBodies(modelIndex, model.RigidBodies)
	mp.initJoints(modelIndex, model.RigidBodies, model.Joints)
}

// AddModelByDeltas はボーンデルタ情報を使用してモデルを物理エンジンに追加します
func (mp *MPhysics) AddModelByDeltas(modelIndex int, model *pmx.PmxModel, boneDeltas *delta.BoneDeltas, physicsDeltas *delta.PhysicsDeltas) {
	// 根元から追加していく
	var rigidBodyDeltas *delta.RigidBodyDeltas
	// var jointDeltas *delta.JointDeltas
	if physicsDeltas != nil {
		rigidBodyDeltas = physicsDeltas.RigidBodies
		// jointDeltas = physicsDeltas.Joints
	}

	mp.initRigidBodiesByBoneDeltas(modelIndex, model.RigidBodies, boneDeltas, rigidBodyDeltas)
	mp.initJointsByBoneDeltas(modelIndex, model.RigidBodies, model.Joints, boneDeltas, nil)
}

// DeleteModel はモデルを物理エンジンから削除します
func (mp *MPhysics) DeleteModel(modelIndex int) {
	// 末端から削除していく
	mp.deleteJoints(modelIndex)
	mp.deleteRigidBodies(modelIndex)
}

// StepSimulation は物理シミュレーションを1ステップ進めます
func (mp *MPhysics) StepSimulation(timeStep float32, maxSubSteps int, fixedTimeStep float32) {
	// 風力の適用（物理更新の直前）
	mp.applyWindForces(timeStep)
	mp.world.StepSimulation(timeStep, maxSubSteps, fixedTimeStep)
}

// EnableWind 風の有効/無効を切り替える
func (mp *MPhysics) EnableWind(enable bool) {
	mp.windCfg.Enabled = enable
}

// SetWind 風向き・風速・ランダム性を設定する
//
//	direction: MMD座標系の風向ベクトル（正規化されていなくてもOK）
//	speed: 風速の基本値（単位/秒）
//	randomness: 0..1 程度の乱れの強さ
func (mp *MPhysics) SetWind(direction *mmath.MVec3, speed float32, randomness float32) {
	if direction != nil {
		mp.windCfg.Direction = direction.Copy()
	}
	mp.windCfg.Speed = speed
	if randomness < 0 {
		randomness = 0
	}
	mp.windCfg.Randomness = randomness
}

// SetWindAdvanced 風の詳細パラメータを設定する
//
//	dragCoeff, liftCoeff は 0.5*rho*Cd*A, 0.5*rho*Cl*A を吸収した係数として扱う
//	turbulenceFreqHz はガストの周波数[Hz]
func (mp *MPhysics) SetWindAdvanced(dragCoeff, liftCoeff, turbulenceFreqHz float32) {
	if dragCoeff >= 0 {
		mp.windCfg.DragCoeff = dragCoeff
	}
	if liftCoeff >= 0 {
		mp.windCfg.LiftCoeff = liftCoeff
	}
	if turbulenceFreqHz > 0 {
		mp.windCfg.TurbulenceFreqHz = turbulenceFreqHz
	}
}

// applyWindForces 風の力（抵抗 + 簡易揚力）を動的剛体に付与する
func (mp *MPhysics) applyWindForces(dt float32) {
	if !mp.windCfg.Enabled {
		return
	}
	if mp.windCfg.Speed == 0 {
		return
	}

	mp.simTimeAcc += dt

	// 乱流係数: 1 + r*sin + r*0.5*sin
	r := float64(mmath.Clamped(float64(mp.windCfg.Randomness), 0, 1))
	f := math.Max(0.0001, float64(mp.windCfg.TurbulenceFreqHz))
	t := float64(mp.simTimeAcc)
	gust := 1.0 + r*(0.6*math.Sin(2*math.Pi*f*t)+0.4*math.Sin(2*math.Pi*1.73*f*t+0.9))

	// 風速（MMD座標系）→ Bullet座標系ベクトル
	dir := mp.windCfg.Direction.Copy().Normalized()
	windSpeed := float64(mp.windCfg.Speed) * gust
	windVecMmd := dir.MuledScalar(windSpeed)
	windVecBt := newBulletFromVec(windVecMmd)

	// Bullet座標の風の各成分
	windX := float64(windVecBt.GetX())
	windY := float64(windVecBt.GetY())
	windZ := float64(windVecBt.GetZ())

	// 各モデル・各剛体に適用
	for _, bodies := range mp.rigidBodies {
		if bodies == nil {
			continue
		}
		for _, rb := range bodies {
			if rb == nil || rb.btRigidBody == nil || rb.pmxRigidBody == nil {
				continue
			}
			// 静的剛体はスキップ
			if rb.pmxRigidBody.PhysicsType == pmx.PHYSICS_TYPE_STATIC {
				continue
			}

			// 相対速度 v_rel = v_body - v_wind
			v := rb.btRigidBody.GetLinearVelocity()
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

			// 断面積の近似（MMD座標系でOK: 絶対値なので X 反転の影響なし）
			area := mp.approxCrossSectionArea(rb.pmxRigidBody, dir)

			// 相対速度の単位ベクトル
			invSpeed := 1.0 / speed
			nrx := relX * invSpeed
			nry := relY * invSpeed
			nrz := relZ * invSpeed

			// 抵抗（風に合わせる方向: -v_rel = v_wind - v_body）
			// Fd = k_d * A * |v_rel|^2
			kd := float64(mp.windCfg.DragCoeff)
			magD := kd * float64(area) * speed2
			fDx := float32(-magD * nrx)
			fDy := float32(-magD * nry)
			fDz := float32(-magD * nrz)

			// 簡易揚力: up を v_rel に直交な平面へ正射影した方向
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

			// 合力
			fTx := fDx + fLx
			fTy := fDy + fLy
			fTz := fDz + fLz

			// 安定化: 最大加速度でクランプ
			maxA := float64(mp.windCfg.MaxAcceleration)
			if maxA > 0 {
				mass := float64(rb.btRigidBody.GetMass())
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
			rb.btRigidBody.ApplyCentralForce(fTotal)
			bt.DeleteBtVector3(fTotal)

			// 剛体をアクティブ化
			rb.btRigidBody.Activate(true)
		}
	}
}

// approxCrossSectionArea 風向きに対する見かけの断面積を近似計算する
// shape の回転は無視し，ボックスは各面の面積を風向の各軸の寄与で線形補間
func (mp *MPhysics) approxCrossSectionArea(r *pmx.RigidBody, dir *mmath.MVec3) float32 {
	if r == nil || dir == nil {
		return 1.0
	}
	d := dir.Normalized()
	absX := math.Abs(d.X)
	absY := math.Abs(d.Y)
	absZ := math.Abs(d.Z)

	switch r.ShapeType {
	case pmx.SHAPE_SPHERE:
		// size.X を半径として使用
		rads := float64(r.Size.X)
		return float32(math.Pi * rads * rads)
	case pmx.SHAPE_BOX:
		// Bullet では半径系。フルサイズ = 2*size
		wx := float64(2.0 * r.Size.X)
		wy := float64(2.0 * r.Size.Y)
		wz := float64(2.0 * r.Size.Z)
		area := absX*(wy*wz) + absY*(wx*wz) + absZ*(wx*wy)
		return float32(area)
	case pmx.SHAPE_CAPSULE:
		// Y軸キャップスル形状を想定
		rad := float64(r.Size.X)
		h := float64(r.Size.Y)
		// 軸方向と垂直方向の断面積を補間（スタジアム断面を近似）
		aAxis := math.Pi * rad * rad
		aPerp := 2*rad*h + math.Pi*rad*rad
		// 寄与: Y 成分が大きいほど軸方向の断面に近い
		w := absY
		return float32(w*aAxis + (1.0-w)*aPerp)
	default:
		// デフォルトは球と同様に扱う
		rads := float64(r.Size.X)
		return float32(math.Pi * rads * rads)
	}
}

func createWorld(gravity *mmath.MVec3) bt.BtDiscreteDynamicsWorld {
	broadphase := bt.NewBtDbvtBroadphase()
	collisionConfiguration := bt.NewBtDefaultCollisionConfiguration()
	dispatcher := bt.NewBtCollisionDispatcher(collisionConfiguration)
	solver := bt.NewBtSequentialImpulseConstraintSolver()
	// solver.GetM_analyticsData().SetM_numIterationsUsed(200)
	world := bt.NewBtDiscreteDynamicsWorld(dispatcher, broadphase, solver, collisionConfiguration)
	world.SetGravity(bt.NewBtVector3(float32(gravity.X), float32(gravity.Y*10), float32(gravity.Z)))
	// world.GetSolverInfo().(bt.BtContactSolverInfo).SetM_numIterations(100)
	// world.GetSolverInfo().(bt.BtContactSolverInfo).SetM_splitImpulse(1)

	groundShape := bt.NewBtStaticPlaneShape(bt.NewBtVector3(float32(0), float32(1), float32(0)), float32(0))
	groundTransform := bt.NewBtTransform()
	groundTransform.SetIdentity()
	groundTransform.SetOrigin(bt.NewBtVector3(float32(0), float32(0), float32(0)))
	groundMotionState := bt.NewBtDefaultMotionState(groundTransform)
	groundRigidBody := bt.NewBtRigidBody(float32(0), groundMotionState, groundShape)

	world.AddRigidBody(groundRigidBody, 1<<15, 0xFFFF)

	return world
}

// UpdatePhysicsSelectively は変更が必要な剛体・ジョイントのみを選択的に更新します
func (mp *MPhysics) UpdatePhysicsSelectively(
	modelIndex int,
	model *pmx.PmxModel,
	physicsDeltas *delta.PhysicsDeltas,
) {
	if physicsDeltas == nil {
		return
	}

	// 剛体の選択的更新
	if physicsDeltas.RigidBodies != nil {
		mp.UpdateRigidBodiesSelectively(modelIndex, model, physicsDeltas.RigidBodies)
	}

	// ジョイントの選択的更新
	if physicsDeltas.Joints != nil {
		mp.UpdateJointsSelectively(modelIndex, model, physicsDeltas.Joints)
	}
}
