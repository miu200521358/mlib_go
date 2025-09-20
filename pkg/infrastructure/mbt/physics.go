//go:build windows
// +build windows

package mbt

import (
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
