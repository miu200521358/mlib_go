package pmx

import (
	"math"

	"github.com/go-gl/gl/v4.4-core/gl"

	"github.com/miu200521358/mlib_go/pkg/mbt"
	"github.com/miu200521358/mlib_go/pkg/mcore"
	"github.com/miu200521358/mlib_go/pkg/mgl"
	"github.com/miu200521358/mlib_go/pkg/mmath"
	"github.com/miu200521358/mlib_go/pkg/mphysics"
)

type JointParam struct {
	TranslationLimitMin       *mmath.MVec3     // 移動制限-下限(x,y,z)
	TranslationLimitMax       *mmath.MVec3     // 移動制限-上限(x,y,z)
	RotationLimitMin          *mmath.MRotation // 回転制限-下限
	RotationLimitMax          *mmath.MRotation // 回転制限-上限
	SpringConstantTranslation *mmath.MVec3     // バネ定数-移動(x,y,z)
	SpringConstantRotation    *mmath.MRotation // バネ定数-回転(x,y,z)
}

func NewJointParam() *JointParam {
	return &JointParam{
		TranslationLimitMin:       mmath.NewMVec3(),
		TranslationLimitMax:       mmath.NewMVec3(),
		RotationLimitMin:          mmath.NewRotationModelByDegrees(mmath.NewMVec3()),
		RotationLimitMax:          mmath.NewRotationModelByDegrees(mmath.NewMVec3()),
		SpringConstantTranslation: mmath.NewMVec3(),
		SpringConstantRotation:    mmath.NewRotationModelByDegrees(mmath.NewMVec3()),
	}
}

type Joint struct {
	*mcore.IndexNameModel
	JointType       byte             // Joint種類 - 0:スプリング6DOF   | PMX2.0では 0 のみ(拡張用)
	RigidbodyIndexA int              // 関連剛体AのIndex
	RigidbodyIndexB int              // 関連剛体BのIndex
	Position        *mmath.MVec3     // 位置(x,y,z)
	Rotation        *mmath.MRotation // 回転
	JointParam      *JointParam      // ジョイントパラメーター
	IsSystem        bool
	Constraint      mbt.BtGeneric6DofSpringConstraint // Bulletのジョイント
}

func NewJoint() *Joint {
	return &Joint{
		IndexNameModel: &mcore.IndexNameModel{Index: -1, Name: "", EnglishName: ""},
		// Joint種類 - 0:スプリング6DOF   | PMX2.0では 0 のみ(拡張用)
		JointType:       0,
		RigidbodyIndexA: -1,
		RigidbodyIndexB: -1,
		Position:        mmath.NewMVec3(),
		Rotation:        mmath.NewRotationModelByDegrees(mmath.NewMVec3()),
		JointParam:      NewJointParam(),
		IsSystem:        false,
	}
}

func NewJointByName(name string) *Joint {
	j := NewJoint()
	j.Name = name
	return j
}

func (j *Joint) InitPhysics(modelPhysics *mphysics.MPhysics, rigidBodyA *RigidBody, rigidBodyB *RigidBody) {
	// ジョイントの位置と向き
	jointTransform := mbt.NewBtTransform(j.Rotation.GetQuaternion().Bullet(), j.Position.Bullet())

	// 剛体Aの現在の位置と向きを取得
	worldTransformA := rigidBodyA.BtRigidBody.GetWorldTransform().(mbt.BtTransform)

	// 剛体Aのローカル座標系におけるジョイント
	jointLocalTransformA := mbt.NewBtTransform()
	jointLocalTransformA.SetIdentity()
	jointLocalTransformA.Mult(worldTransformA.Inverse(), jointTransform)

	// 剛体Bの現在の位置と向きを取得
	worldTransformB := rigidBodyB.BtRigidBody.GetWorldTransform().(mbt.BtTransform)

	// 剛体Bのローカル座標系におけるジョイント
	jointLocalTransformB := mbt.NewBtTransform()
	jointLocalTransformB.SetIdentity()
	jointLocalTransformB.Mult(worldTransformB.Inverse(), jointTransform)

	// ジョイント係数
	j.Constraint = mbt.NewBtGeneric6DofSpringConstraint(
		rigidBodyA.BtRigidBody, rigidBodyB.BtRigidBody, jointLocalTransformA, jointLocalTransformB, true)
	// 係数は符号を調整する必要がないため、そのまま設定
	j.Constraint.SetLinearLowerLimit(mbt.NewBtVector3(
		float32(j.JointParam.TranslationLimitMin.GetX()),
		float32(j.JointParam.TranslationLimitMin.GetY()),
		float32(j.JointParam.TranslationLimitMin.GetZ())))
	j.Constraint.SetLinearUpperLimit(mbt.NewBtVector3(
		float32(j.JointParam.TranslationLimitMax.GetX()),
		float32(j.JointParam.TranslationLimitMax.GetY()),
		float32(j.JointParam.TranslationLimitMax.GetZ())))
	j.Constraint.SetAngularLowerLimit(mbt.NewBtVector3(
		float32(j.JointParam.RotationLimitMin.GetRadians().GetX()),
		float32(j.JointParam.RotationLimitMin.GetRadians().GetY()),
		float32(j.JointParam.RotationLimitMin.GetRadians().GetZ())))
	j.Constraint.SetAngularUpperLimit(mbt.NewBtVector3(
		float32(j.JointParam.RotationLimitMax.GetRadians().GetX()),
		float32(j.JointParam.RotationLimitMax.GetRadians().GetY()),
		float32(j.JointParam.RotationLimitMax.GetRadians().GetZ())))
	j.Constraint.EnableSpring(0, true)
	j.Constraint.SetStiffness(0, float32(j.JointParam.SpringConstantTranslation.GetX()))
	j.Constraint.EnableSpring(1, true)
	j.Constraint.SetStiffness(1, float32(j.JointParam.SpringConstantTranslation.GetY()))
	j.Constraint.EnableSpring(2, true)
	j.Constraint.SetStiffness(2, float32(j.JointParam.SpringConstantTranslation.GetZ()))
	j.Constraint.EnableSpring(3, true)
	j.Constraint.SetStiffness(3, float32(j.JointParam.SpringConstantRotation.GetRadians().GetX()))
	j.Constraint.EnableSpring(4, true)
	j.Constraint.SetStiffness(4, float32(j.JointParam.SpringConstantRotation.GetRadians().GetY()))
	j.Constraint.EnableSpring(5, true)
	j.Constraint.SetStiffness(5, float32(j.JointParam.SpringConstantRotation.GetRadians().GetZ()))

	modelPhysics.AddJoint(j.Constraint)
}

func (j *Joint) updateVbo(jointVbo []float32, segments int) []float32 {
	// ジョイントの位置と向き
	// Get the rigid bodies from the constraint
	bodyA := j.Constraint.GetRigidBodyA().(mbt.BtRigidBody)
	bodyB := j.Constraint.GetRigidBodyB().(mbt.BtRigidBody)

	// Get the world transforms of the rigid bodies
	transformA := bodyA.GetWorldTransform().(mbt.BtTransform)
	transformB := bodyB.GetWorldTransform().(mbt.BtTransform)

	// Get the local transforms of the constraint
	localTransformA := j.Constraint.GetFrameOffsetA().(mbt.BtTransform)
	localTransformB := j.Constraint.GetFrameOffsetB().(mbt.BtTransform)

	// Calculate the global transforms by combining the world and local transforms
	globalTransformA := mbt.NewBtTransform()
	globalTransformA.Mult(transformA, localTransformA)
	globalTransformB := mbt.NewBtTransform()
	globalTransformB.Mult(transformB, localTransformB)

	// Get the global positions and rotations from the global transforms
	globalPositionA := globalTransformA.GetOrigin().(mbt.BtVector3)
	globalPositionB := globalTransformB.GetOrigin().(mbt.BtVector3)

	globalRotationA := globalTransformA.GetRotation()
	globalRotationB := globalTransformB.GetRotation()
	globalRotation := globalRotationA.Slerp(globalRotationB, 0.5)

	position := mmath.MVec3{
		float64((globalPositionA.GetX() + globalPositionB.GetX()) / 2),
		float64((globalPositionA.GetY() + globalPositionB.GetY()) / 2),
		float64((globalPositionA.GetZ() + globalPositionB.GetZ()) / 2),
	}

	rotation := mmath.NewMQuaternionByValues(
		float64(globalRotation.GetX()),
		float64(globalRotation.GetY()),
		float64(globalRotation.GetZ()),
		float64(globalRotation.GetW()),
	)

	// Add the global positions and rotations to the VBO
	jointVbo = j.addJointVbo(&position, rotation, jointVbo, segments)

	return jointVbo
}

// 球の頂点を追加
func (j *Joint) addJointVbo(
	position *mmath.MVec3,
	rotation *mmath.MQuaternion,
	jointVbo []float32,
	segments int,
) []float32 {
	// 球体を6x6に分割して線で描くため、球体の頂点を追加

	for m := 0; m <= segments; m++ {
		for n := 0; n <= segments; n++ {
			// Calculate the latitude and longitude
			lat := (float64(m) / float64(segments)) * math.Pi
			lon := (float64(n) / float64(segments)) * 2.0 * math.Pi

			size := 0.2

			// Calculate the x, y, z coordinates
			x := size * math.Sin(lat) * math.Cos(lon)
			y := size * math.Cos(lat)
			z := size * math.Sin(lat) * math.Sin(lon)
			spherePos := &mmath.MVec3{x, y, z}

			// Apply the rigid body's rotation to the vector
			rotatedPos := rotation.RotatedVec3(spherePos).Added(position)

			// Append the coordinates to the VBO
			jointVbo = append(
				jointVbo,
				// 0: typeColor
				float32(0.0),
				float32(0.0),
				float32(1.0),
				float32(0.6),
				// 1: position
				float32(rotatedPos.GetX()),
				float32(rotatedPos.GetY()),
				float32(rotatedPos.GetZ()),
			)
		}
	}

	return jointVbo
}

// ジョイントリスト
type Joints struct {
	*mcore.IndexNameModelCorrection[*Joint]
	vao   *mgl.VAO
	vbo   *mgl.VBO
	ibo   *mgl.IBO
	count int32
}

func NewJoints() *Joints {
	return &Joints{
		IndexNameModelCorrection: mcore.NewIndexNameModelCorrection[*Joint](),
		vao:                      nil,
		vbo:                      nil,
		ibo:                      nil,
		count:                    0,
	}
}

func (j *Joints) InitPhysics(modelPhysics *mphysics.MPhysics, rigidBodies *RigidBodies) {
	// ジョイントを順番に剛体と紐付けていく
	for _, joint := range j.GetSortedData() {
		if joint.RigidbodyIndexA >= 0 && rigidBodies.Contains(joint.RigidbodyIndexA) &&
			joint.RigidbodyIndexB >= 0 && rigidBodies.Contains(joint.RigidbodyIndexB) {
			joint.InitPhysics(
				modelPhysics, rigidBodies.GetItem(joint.RigidbodyIndexA),
				rigidBodies.GetItem(joint.RigidbodyIndexB))
		}
	}

	// 描画準備
	j.prepareDraw()
}

func (j *Joints) prepareDraw() {
	jointIbo := make([]uint32, 0)

	startIdx := 0
	segments := 6

	for range j.Data {
		for m := 0; m <= segments; m++ {
			for n := 0; n <= segments; n++ {
				// Calculate the indices of the current point and the points to the right and below it
				idx := startIdx + m*(segments+1) + n
				idxRight := idx + 1
				idxBelow := idx + segments + 1

				// Add a line from the current point to the point to the right
				if n < segments {
					jointIbo = append(jointIbo, uint32(idx), uint32(idxRight))
				}

				// Add a line from the current point to the point below
				if m < segments {
					jointIbo = append(jointIbo, uint32(idx), uint32(idxBelow))
				}
			}
		}
		startIdx += (segments+1)*(segments+1) + (segments + 1)
	}

	j.vao = mgl.NewVAO()
	j.vao.Bind()
	j.vao.Unbind()

	j.ibo = mgl.NewIBO(gl.Ptr(jointIbo), len(jointIbo))
	j.count = int32(len(jointIbo))
}

func (j *Joints) updateVbo() {
	// ジョイント位置を更新
	startIdx := 0
	segments := 6
	jointVbo := make([]float32, 0)

	for _, joint := range j.GetSortedData() {
		jointVbo = joint.updateVbo(jointVbo, segments)
		startIdx += (segments + 1) * (segments + 1)
	}

	j.vao.Bind()
	j.vbo = mgl.NewVBOForJoint(gl.Ptr(jointVbo), len(jointVbo))
	j.vbo.BindJoint()
	j.vbo.Unbind()
	j.vao.Unbind()
}

func (j *Joints) Draw(
	shader *mgl.MShader,
	windowIndex int,
) {
	j.updateVbo()

	// ボーンをモデルメッシュの前面に描画するために深度テストを無効化
	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.ALWAYS)

	// ブレンディングを有効にする
	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)

	// ブレンディングを有効にする
	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)

	shader.UseJointProgram()

	j.vao.Bind()
	j.vbo.BindJoint()
	j.ibo.Bind()

	gl.DrawElements(gl.LINES, j.count, gl.UNSIGNED_INT, nil)

	j.ibo.Unbind()
	j.vbo.Unbind()
	j.vao.Unbind()

	shader.Unuse()

	gl.Disable(gl.BLEND)
}
