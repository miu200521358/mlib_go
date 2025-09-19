//go:build windows
// +build windows

package mbt

import (
	"math"

	"github.com/miu200521358/mlib_go/pkg/config/mlog"
	"github.com/miu200521358/mlib_go/pkg/domain/delta"
	"github.com/miu200521358/mlib_go/pkg/domain/mmath"
	"github.com/miu200521358/mlib_go/pkg/domain/physics"
	"github.com/miu200521358/mlib_go/pkg/domain/pmx"
	"github.com/miu200521358/mlib_go/pkg/domain/rendering"
	"github.com/miu200521358/mlib_go/pkg/infrastructure/bt"
)

// RigidBodyHit はレイキャストでヒットした剛体の情報
type RigidBodyHit struct {
	ModelIndex     int            // モデルのインデックス
	RigidBodyIndex int            // 剛体のインデックス
	RigidBody      *pmx.RigidBody // 剛体の情報
	Distance       float32        // カメラからの距離
	HitPoint       *mmath.MVec3   // ヒットした座標（ワールド座標）
}

// MPhysics 物理エンジンの実装
type MPhysics struct {
	world       bt.BtDiscreteDynamicsWorld // ワールド
	drawer      bt.BtMDebugDraw            // デバッグビューワー
	liner       *mDebugDrawLiner           // ライナー
	highlighter *mDebugDrawHighlighter     // ハイライト描画器
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
	highlighter := newMDebugDrawHighlighter()
	drawer := bt.NewBtMDebugDraw()
	drawer.SetLiner(liner)
	drawer.SetMDefaultColors(newConstBtMDefaultColors())
	mp.world.SetDebugDrawer(drawer)

	mp.drawer = drawer
	mp.liner = liner
	mp.highlighter = highlighter
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

// RaycastRigidBody は画面座標からレイキャストを行い、最前面の剛体を取得します
func (mp *MPhysics) RaycastRigidBody(screenX, screenY float64, camera *rendering.Camera, width, height int) (*RigidBodyHit, error) {
	// 簡易実装：カメラ情報を使わずに距離ベースの選択を行う
	// スクリーン座標からワールド座標のレイを生成
	rayStart, rayEnd := mp.screenToWorldRay(screenX, screenY, camera, width, height)
	if rayStart == nil || rayEnd == nil {
		// デバッグログ：座標変換エラー
		// mlog.W("Screen to world coordinate conversion failed")
		return nil, nil
	}

	// デバッグログ：レイキャスト実行
	// mlog.I("Raycast: screen(%f, %f) -> world(%f,%f,%f) to (%f,%f,%f)",
	//   screenX, screenY, rayStart.X, rayStart.Y, rayStart.Z, rayEnd.X, rayEnd.Y, rayEnd.Z)

	// Bullet物理エンジンでのレイキャスト
	btRayStart := newBulletFromVec(rayStart)
	btRayEnd := newBulletFromVec(rayEnd)
	defer bt.DeleteBtVector3(btRayStart)
	defer bt.DeleteBtVector3(btRayEnd)

	// レイキャスト結果を格納するコールバック（全ての交差を取得）
	// TODO: Bulletの正しいAPI名を調査
	// 現在は簡易実装として距離ベースの選択を行う
	return mp.performRaycastSelection(rayStart, rayEnd, btRayStart, btRayEnd)
}

// screenToWorldRay はスクリーン座標からワールド座標のレイを生成します
func (mp *MPhysics) screenToWorldRay(screenX, screenY float64, camera *rendering.Camera, width, height int) (*mmath.MVec3, *mmath.MVec3) {
	// NDC座標に変換 (-1 to 1)
	ndcX := (2.0 * screenX / float64(width)) - 1.0
	ndcY := 1.0 - (2.0 * screenY / float64(height)) // Y軸反転

	// デバッグログ：入力値とNDC変換結果
	// mlog.I("スクリーン座標: (%f, %f), ウィンドウサイズ: (%d, %d), NDC座標: (%f, %f)",
	//   screenX, screenY, width, height, ndcX, ndcY)

	// 実際のカメラ情報を使用（Camera構造体に直接アクセス）
	return mp.createRayFromCameraInfo(ndcX, ndcY, camera, width, height)
}

// createRayFromCameraInfo は実際のカメラ情報からレイを生成します
func (mp *MPhysics) createRayFromCameraInfo(ndcX, ndcY float64, camera *rendering.Camera, width, height int) (*mmath.MVec3, *mmath.MVec3) {
	// 実際のカメラの前方ベクトルを計算
	forward := camera.LookAtCenter.Subed(camera.Position).Normalized()

	// カメラの右方向と上方向ベクトルを計算
	up := camera.Up.Copy().Normalized()
	right := forward.Cross(up).Normalized()
	up = right.Cross(forward).Normalized()

	// 実際のFOV（視野角）を使用
	fov := float64(camera.FieldOfView) * math.Pi / 180.0
	aspect := float64(width) / float64(height)

	// レイの方向を計算
	h := math.Tan(fov / 2.0)
	w := h * aspect

	// スクリーン座標に対応するワールド方向ベクトル
	rayDirection := forward.Copy()
	rayDirection.Add(right.MulScalar(ndcX * w))
	rayDirection.Add(up.MulScalar(ndcY * h))
	rayDirection = rayDirection.Normalized()

	// レイの開始点（実際のカメラ位置）
	rayStart := camera.Position.Copy()

	// レイの終了点（遠い距離）
	farDistance := 1000.0
	rayEnd := rayStart.Copy()
	rayEnd.Add(rayDirection.MulScalar(farDistance))

	// デバッグログ：実際のカメラ情報
	mlog.I("Camera pos: (%f,%f,%f), lookAt: (%f,%f,%f), FOV: %f",
		camera.Position.X, camera.Position.Y, camera.Position.Z,
		camera.LookAtCenter.X, camera.LookAtCenter.Y, camera.LookAtCenter.Z, camera.FieldOfView)

	return rayStart, rayEnd
}

// performRaycastSelection は実際のレイキャストを実行し、最適な剛体を選択します
func (mp *MPhysics) performRaycastSelection(rayStart, rayEnd *mmath.MVec3, btRayStart, btRayEnd bt.BtVector3) (*RigidBodyHit, error) {
	// 現在は簡易実装：全ての剛体をチェックし、レイとの距離が最も近いものを選択
	var closestHit *RigidBodyHit
	minDistance := float64(math.MaxFloat64)

	for modelIndex, bodies := range mp.rigidBodies {
		if len(bodies) == 0 {
			continue
		}

		for rigidBodyIndex, rb := range bodies {
			if rb == nil || rb.pmxRigidBody == nil || rb.btRigidBody == nil {
				continue
			}

			// Bullet物理世界での実際の剛体位置を取得（現在は一時的にPMX位置を使用）
			// TODO: Bulletの型アサーションを修正後に有効化
			// transform := rb.btRigidBody.GetWorldTransform()
			// 現在はPMX初期位置を使用
			rigidBodyPos := rb.pmxRigidBody.Position

			// デバッグログ：使用中の位置情報
			// mlog.I("剛体位置（PMX初期値使用）: %s (%f,%f,%f)",
			//   rb.pmxRigidBody.Name(), rigidBodyPos.X, rigidBodyPos.Y, rigidBodyPos.Z)

			// レイと剛体の中心点との距離を計算（簡易的な判定）
			distance := mp.calculateRayPointDistance(rayStart, rayEnd, rigidBodyPos)

			// 剛体のサイズを考慮した判定（簡易的）
			rigidBodyRadius := float64(math.Max(rb.pmxRigidBody.Size.X, math.Max(rb.pmxRigidBody.Size.Y, rb.pmxRigidBody.Size.Z)))

			// デバッグログ：詳細な剛体情報
			// mlog.I("剛体チェック: %s 位置:(%f,%f,%f) 距離:%f 半径:%f",
			//   rb.pmxRigidBody.Name(), rigidBodyPos.X, rigidBodyPos.Y, rigidBodyPos.Z, distance, rigidBodyRadius)

			if distance <= rigidBodyRadius && distance < minDistance {
				minDistance = distance
				closestHit = &RigidBodyHit{
					ModelIndex:     modelIndex,
					RigidBodyIndex: rigidBodyIndex,
					RigidBody:      rb.pmxRigidBody,
					Distance:       float32(distance),
					HitPoint:       rigidBodyPos.Copy(),
				}
				// デバッグログ：候補剛体
				// mlog.I("候補剛体更新: %s (距離: %f)", rb.pmxRigidBody.Name(), distance)
			}
		}
	}

	if closestHit != nil {
		// デバッグログ：最終選択された剛体
		mlog.I("Raycast hit: %s (distance: %f)", closestHit.RigidBody.Name(), closestHit.Distance)
	} else {
		// デバッグログ：剛体なし
		mlog.I("No rigid body hit")
	}

	return closestHit, nil
}

// calculateRayPointDistance はレイと点との最短距離を計算します
func (mp *MPhysics) calculateRayPointDistance(rayStart, rayEnd, point *mmath.MVec3) float64 {
	// レイのベクトル
	rayVec := rayEnd.Subed(rayStart)

	// レイの開始点から対象点へのベクトル
	toPoint := point.Subed(rayStart)

	// レイ上の最近点を計算
	rayLength := rayVec.Length()
	if rayLength < 1e-10 {
		return rayStart.Distance(point)
	}

	rayDir := rayVec.Normalized()
	projLength := toPoint.Dot(rayDir)

	// プロジェクションをレイの範囲内にクランプ
	if projLength < 0 {
		projLength = 0
	} else if projLength > rayLength {
		projLength = rayLength
	}

	// レイ上の最近点
	closestPointOnRay := rayStart.Copy()
	closestPointOnRay.Add(rayDir.MulScalar(projLength))

	// 点とレイ上の最近点との距離
	return closestPointOnRay.Distance(point)
}

// SetSelectedRigidBody はハイライト表示する剛体を設定します
func (mp *MPhysics) SetSelectedRigidBody(hit *RigidBodyHit) {
	if mp.highlighter != nil {
		mp.highlighter.SetSelectedRigidBody(hit)
	}
}

// ClearSelectedRigidBody はハイライト表示を解除します
func (mp *MPhysics) ClearSelectedRigidBody() {
	if mp.highlighter != nil {
		mp.highlighter.ClearSelection()
	}
}

// DrawRigidBodyHighlight はハイライトした剛体を描画します
func (mp *MPhysics) DrawRigidBodyHighlight(shader rendering.IShader, isDrawRigidBodyFront bool) {
	if mp.highlighter != nil {
		// TODO: 型エラー回避のため一旦コメントアウト
		// 基本的なマウスホバー処理の動作確認後に修正
		// mp.highlighter.drawHighlight(shader, isDrawRigidBodyFront)
	}
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
