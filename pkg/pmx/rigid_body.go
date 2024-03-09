package pmx

import (
	"math"

	"github.com/go-gl/gl/v4.4-core/gl"
	"github.com/go-gl/mathgl/mgl32"

	"github.com/miu200521358/mlib_go/pkg/mbt"
	"github.com/miu200521358/mlib_go/pkg/mcore"
	"github.com/miu200521358/mlib_go/pkg/mgl"
	"github.com/miu200521358/mlib_go/pkg/mmath"
	"github.com/miu200521358/mlib_go/pkg/mphysics"
)

type RigidBodyParam struct {
	Mass           float64 // 質量
	LinearDamping  float64 // 移動減衰
	AngularDamping float64 // 回転減衰
	Restitution    float64 // 反発力
	Friction       float64 // 摩擦力
}

func NewRigidBodyParam() *RigidBodyParam {
	return &RigidBodyParam{
		Mass:           0,
		LinearDamping:  0,
		AngularDamping: 0,
		Restitution:    0,
		Friction:       0,
	}
}

// 剛体の形状
type Shape int

const (
	SHAPE_SPHERE  Shape = 0 // 球
	SHAPE_BOX     Shape = 1 // 箱
	SHAPE_CAPSULE Shape = 2 // カプセル
)

// 剛体物理の計算モード
type PhysicsType int

const (
	PHYSICS_TYPE_STATIC       PhysicsType = 0 // ボーン追従(static)
	PHYSICS_TYPE_DYNAMIC      PhysicsType = 1 // 物理演算(dynamic)
	PHYSICS_TYPE_DYNAMIC_BONE PhysicsType = 2 // 物理演算 + Bone位置合わせ
)

type CollisionGroup struct {
	IsCollisions []uint16
}

var CollisionGroupFlags = []uint16{
	0x0001, // 0:グループ1
	0x0002, // 1:グループ2
	0x0004, // 2:グループ3
	0x0008, // 3:グループ4
	0x0010, // 4:グループ5
	0x0020, // 5:グループ6
	0x0040, // 6:グループ7
	0x0080, // 7:グループ8
	0x0100, // 8:グループ9
	0x0200, // 9:グループ10
	0x0400, // 10:グループ11
	0x0800, // 11:グループ12
	0x1000, // 12:グループ13
	0x2000, // 13:グループ14
	0x4000, // 14:グループ15
	0x8000, // 15:グループ16
}

func NewCollisionGroupFromSlice(collisionGroup []uint16) CollisionGroup {
	groups := CollisionGroup{}
	collisionGroupMask := uint16(0)
	for i, v := range collisionGroup {
		if v == 1 {
			collisionGroupMask |= CollisionGroupFlags[i]
		}
	}
	groups.IsCollisions = NewCollisionGroup(collisionGroupMask)

	return groups
}

func NewCollisionGroup(collisionGroupMask uint16) []uint16 {
	collisionGroup := []uint16{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	for i, v := range CollisionGroupFlags {
		if collisionGroupMask&v == v {
			collisionGroup[i] = 0
		} else {
			collisionGroup[i] = 1
		}
	}
	return collisionGroup
}

type RigidBody struct {
	*mcore.IndexNameModel
	BoneIndex                    int              // 関連ボーンIndex
	CollisionGroup               byte             // グループ
	CollisionGroupMask           CollisionGroup   // 非衝突グループフラグ
	CollisionGroupMaskValue      int              // 非衝突グループフラグ値
	ShapeType                    Shape            // 形状
	Size                         *mmath.MVec3     // サイズ(x,y,z)
	Position                     *mmath.MVec3     // 位置(x,y,z)
	Rotation                     *mmath.MRotation // 回転(x,y,z) -> ラジアン角
	RigidBodyParam               *RigidBodyParam  // 剛体パラ
	PhysicsType                  PhysicsType      // 剛体の物理演算
	XDirection                   *mmath.MVec3     // X軸方向
	YDirection                   *mmath.MVec3     // Y軸方向
	ZDirection                   *mmath.MVec3     // Z軸方向
	IsSystem                     bool             // システムで追加した剛体か
	Matrix                       *mmath.MMat4     // 剛体の行列
	BtRigidBody                  mbt.BtRigidBody  // 物理剛体
	BtRigidBodyTransform         mbt.BtTransform  // 剛体の初期位置情報
	BtRigidBodyPositionTransform mbt.BtTransform  // 剛体の初期逆位置情報
}

// NewRigidBody creates a new rigid body.
func NewRigidBody() *RigidBody {
	return &RigidBody{
		IndexNameModel:               &mcore.IndexNameModel{Index: -1, Name: "", EnglishName: ""},
		BoneIndex:                    -1,
		CollisionGroup:               0,
		CollisionGroupMask:           NewCollisionGroupFromSlice([]uint16{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}),
		CollisionGroupMaskValue:      0,
		ShapeType:                    SHAPE_BOX,
		Size:                         mmath.NewMVec3(),
		Position:                     mmath.NewMVec3(),
		Rotation:                     mmath.NewRotationModelByDegrees(mmath.NewMVec3()),
		RigidBodyParam:               NewRigidBodyParam(),
		PhysicsType:                  PHYSICS_TYPE_STATIC,
		XDirection:                   mmath.NewMVec3(),
		YDirection:                   mmath.NewMVec3(),
		ZDirection:                   mmath.NewMVec3(),
		IsSystem:                     false,
		Matrix:                       mmath.NewMMat4(),
		BtRigidBody:                  nil,
		BtRigidBodyTransform:         nil,
		BtRigidBodyPositionTransform: nil,
	}
}

func (r *RigidBody) InitPhysics(modelPhysics *mphysics.MPhysics, bone *Bone) {
	var btCollisionShape mbt.BtCollisionShape
	switch r.ShapeType {
	case SHAPE_SPHERE:
		btCollisionShape = mbt.NewBtSphereShape(float32(r.Size.GetX()))
	case SHAPE_BOX:
		btCollisionShape = mbt.NewBtBoxShape(
			mbt.NewBtVector3(float32(r.Size.GetX()), float32(r.Size.GetY()), float32(r.Size.GetZ())))
	case SHAPE_CAPSULE:
		btCollisionShape = mbt.NewBtCapsuleShape(float32(r.Size.GetX()), float32(r.Size.GetY()))
	}

	// 質量
	mass := float32(0.0)
	localInertia := mbt.NewBtVector3(float32(0.0), float32(0.0), float32(0.0))
	if r.PhysicsType != PHYSICS_TYPE_STATIC {
		// ボーン追従ではない場合そのまま設定
		mass = float32(r.RigidBodyParam.Mass)
	}
	if mass != 0 {
		// 質量が設定されている場合、慣性を計算
		btCollisionShape.CalculateLocalInertia(mass, localInertia)
	}

	// // ボーンのローカル位置
	// boneTransform := mbt.NewBtTransform()
	// boneTransform.SetIdentity()
	// boneTransform.SetOrigin(boneLocalPosition.Bullet())

	// 剛体の初期位置と回転
	r.BtRigidBodyTransform = mbt.NewBtTransform(
		r.Rotation.GetQuaternion().Bullet(), r.Position.Bullet())
	// 剛体の初期位置(回は加味しない)
	r.BtRigidBodyPositionTransform = mbt.NewBtTransform()
	r.BtRigidBodyPositionTransform.SetIdentity()
	r.BtRigidBodyPositionTransform.SetOrigin(r.Position.Bullet())

	// {
	// 	fmt.Println("---------------------------------")
	// }
	// {
	// 	mat := mgl32.Mat4{}
	// 	r.BtRigidBodyTransform.GetOpenGLMatrix(&mat[0])
	// 	fmt.Printf("1. [%s] BtRigidBodyTransform: \n%v\n", r.Name, mat)
	// }

	// 剛体のグローバル位置と回転
	motionState := mbt.NewBtDefaultMotionState(r.BtRigidBodyTransform)

	r.BtRigidBody = mbt.NewBtRigidBody(mass, motionState, btCollisionShape, localInertia)
	r.BtRigidBody.SetDamping(float32(r.RigidBodyParam.LinearDamping), float32(r.RigidBodyParam.AngularDamping))
	r.BtRigidBody.SetRestitution(float32(r.RigidBodyParam.Restitution))
	r.BtRigidBody.SetFriction(float32(r.RigidBodyParam.Friction))
	// btRigidBody.SetUserIndex(mbt.)
	r.BtRigidBody.SetSleepingThresholds(0.01, (180.0 * 0.1 / math.Pi))

	if r.PhysicsType == PHYSICS_TYPE_STATIC {
		// 剛体の位置更新に物理演算を使わない。
		// MotionState::getWorldTransformが毎ステップコールされるようになるのでここで剛体位置を更新する。
		r.BtRigidBody.SetCollisionFlags(
			r.BtRigidBody.GetCollisionFlags() | int(mbt.BtCollisionObjectCF_KINEMATIC_OBJECT))
		// 毎ステップの剛体位置通知を無効にする
		// MotionState::setWorldTransformの毎ステップ呼び出しが無効になる(剛体位置は判っているので不要)
		r.BtRigidBody.SetActivationState(mbt.DISABLE_SIMULATION)
	} else {
		// 物理演算・物理+ボーン位置合わせの場合
		// 剛体の位置更新に物理演算を使う。
		// MotionState::getWorldTransformが毎ステップコールされるようになるのでここで剛体位置を更新する。
		r.BtRigidBody.SetCollisionFlags(
			r.BtRigidBody.GetCollisionFlags() & int(^mbt.BtCollisionObjectCF_KINEMATIC_OBJECT))
		// 毎ステップの剛体位置通知を有効にする
		// MotionState::setWorldTransformの毎ステップ呼び出しが有効になる(剛体位置が変わるので必要)
		r.BtRigidBody.SetActivationState(mbt.ACTIVE_TAG)
	}

	modelPhysics.AddRigidBody(r.BtRigidBody, int(r.CollisionGroup), r.CollisionGroupMaskValue)
}

func (r *RigidBody) UpdateTransform(
	boneMatrixes []*mgl32.Mat4,
	boneTransforms []*mbt.BtTransform,
) {
	if r.BtRigidBody == nil || r.BtRigidBody.GetMotionState() == nil ||
		r.BoneIndex < 0 || r.BoneIndex >= len(boneTransforms) || r.PhysicsType != PHYSICS_TYPE_STATIC {
		return
	}

	// {
	// 	fmt.Println("----------")
	// }

	// 剛体のグローバル位置を確定
	rigidBodyTransform := mbt.NewBtTransform()
	rigidBodyTransform.Mult(*boneTransforms[r.BoneIndex], r.BtRigidBodyTransform)

	// {
	// 	mat := mgl32.Mat4{}
	// 	(*boneTransforms[r.BoneIndex]).GetOpenGLMatrix(&mat[0])
	// 	fmt.Printf("2. [%d] boneTransform: \n%v\n", r.BoneIndex, mat)
	// }
	// {
	// 	mat := mgl32.Mat4{}
	// 	rigidBodyTransform.GetOpenGLMatrix(&mat[0])
	// 	fmt.Printf("2. [%s] rigidBodyTransform: \n%v\n", r.Name, mat)
	// }

	motionState := r.BtRigidBody.GetMotionState().(mbt.BtMotionState)
	motionState.SetWorldTransform(rigidBodyTransform)
}

func (r *RigidBody) UpdateMatrix(
	boneMatrixes []*mgl32.Mat4,
	boneTransforms []*mbt.BtTransform,
) {
	if r.BtRigidBody == nil || r.BtRigidBody.GetMotionState() == nil ||
		r.BoneIndex < 0 || r.BoneIndex >= len(boneMatrixes) || r.PhysicsType == PHYSICS_TYPE_STATIC {
		return
	}

	// {
	// 	fmt.Println("----------")
	// }

	motionState := r.BtRigidBody.GetMotionState().(mbt.BtMotionState)

	rigidBodyTransform := mbt.NewBtTransform()
	motionState.GetWorldTransform(rigidBodyTransform)

	if r.PhysicsType == PHYSICS_TYPE_DYNAMIC_BONE {
		// {
		// 	mat := mgl32.Mat4{}
		// 	rigidBodyTransform.GetOpenGLMatrix(&mat[0])
		// 	fmt.Printf("3. [%s] rigidBodyTransform Before: \n%v\n", r.Name, mat)
		// }

		// 物理+ボーン追従はボーン移動成分を現在のボーン位置にする(回転は加味しない)
		boneGlobalTransform := mbt.NewBtTransform()
		boneGlobalTransform.Mult(*boneTransforms[r.BoneIndex], r.BtRigidBodyPositionTransform)

		rigidBodyTransform.SetOrigin(boneGlobalTransform.GetOrigin().(mbt.BtVector3))

		// {
		// 	mat := mgl32.Mat4{}
		// 	rigidBodyTransform.GetOpenGLMatrix(&mat[0])
		// 	fmt.Printf("3. [%s] rigidBodyTransform After: \n%v\n", r.Name, mat)
		// }
	}

	boneLocalTransform := mbt.NewBtTransform()
	boneLocalTransform.Mult(rigidBodyTransform, r.BtRigidBodyTransform.Inverse())

	physicsBoneMatrix := mgl32.Mat4{}
	boneLocalTransform.GetOpenGLMatrix(&physicsBoneMatrix[0])

	// {
	// 	mat := mgl32.Mat4{}
	// 	(*boneTransforms[r.BoneIndex]).GetOpenGLMatrix(&mat[0])
	// 	fmt.Printf("3. [%d] boneTransform: \n%v\n", r.BoneIndex, mat)
	// }
	// {
	// 	mat := mgl32.Mat4{}
	// 	rigidBodyTransform.GetOpenGLMatrix(&mat[0])
	// 	fmt.Printf("3. [%s] rigidBodyTransform: \n%v\n", r.Name, mat)
	// }
	// {
	// 	fmt.Printf("3. [%s] physicsBoneMatrix: \n%v\n", r.Name, physicsBoneMatrix)
	// }

	boneMatrixes[r.BoneIndex] = &physicsBoneMatrix
}

func (r *RigidBody) Bullet() []float32 {
	rb := make([]float32, 0)

	// ボーンINDEX
	rb = append(rb, float32(r.BoneIndex))

	// 剛体の色
	if r.PhysicsType == PHYSICS_TYPE_STATIC {
		// ボーン追従剛体：黄緑色
		rb = append(rb, float32(0.0))
		rb = append(rb, float32(1.0))
		rb = append(rb, float32(0.0))
	} else if r.PhysicsType == PHYSICS_TYPE_DYNAMIC {
		// 物理剛体: 赤色
		rb = append(rb, float32(1.0))
		rb = append(rb, float32(0.0))
		rb = append(rb, float32(0.0))
	} else {
		// ボーン追従 + 物理剛体: 黄色
		rb = append(rb, float32(1.0))
		rb = append(rb, float32(1.0))
		rb = append(rb, float32(0.0))
	}
	rb = append(rb, float32(0.6))

	return rb
}

// 剛体リスト
type RigidBodies struct {
	*mcore.IndexNameModelCorrection[*RigidBody]
	vao         *mgl.VAO
	vbo         *mgl.VBO
	ibo         *mgl.IBO
	count       int32
	sphereVao   *mgl.VAO
	sphereVbo   *mgl.VBO
	shapeCount  int32
	sphereSizes []float32
}

func NewRigidBodies() *RigidBodies {
	return &RigidBodies{
		IndexNameModelCorrection: mcore.NewIndexNameModelCorrection[*RigidBody](),
		vao:                      nil,
		vbo:                      nil,
		ibo:                      nil,
		count:                    0,
		sphereVao:                nil,
		sphereVbo:                nil,
		shapeCount:               0,
		sphereSizes:              make([]float32, 0),
	}
}

func (r *RigidBodies) addCornerBox(
	rigidBody *RigidBody,
	rigidBodyVbo []float32,
	rigidBodyIbo []uint32,
	startIdx int,
) ([]float32, []uint32) {
	// 直方体の角を設定していく
	size := rigidBody.Size

	// 傾き加味
	rotation := rigidBody.Rotation.GetQuaternion().Copy()
	rotation.SetX(-rotation.GetX())
	rotation.SetY(-rotation.GetY())
	rotation.Normalize()

	pos := rigidBody.Position.Copy()
	// pos.SetZ(-pos.GetZ())

	mat := mmath.NewMMat4()
	mat.Translate(pos)
	mat.Mul(rotation.ToMat4())

	{
		// 手前左上
		posGl := mat.MulVec3(&mmath.MVec3{-size.GetX(), -size.GetY(), -size.GetZ()}).GL()

		rigidBodyVbo = append(rigidBodyVbo, rigidBody.Bullet()...)
		rigidBodyVbo = append(rigidBodyVbo, posGl[0])
		rigidBodyVbo = append(rigidBodyVbo, posGl[1])
		rigidBodyVbo = append(rigidBodyVbo, posGl[2])
	}
	{
		// 手前右上
		posGl := mat.MulVec3(&mmath.MVec3{size.GetX(), -size.GetY(), -size.GetZ()}).GL()

		rigidBodyVbo = append(rigidBodyVbo, rigidBody.Bullet()...)
		rigidBodyVbo = append(rigidBodyVbo, posGl[0])
		rigidBodyVbo = append(rigidBodyVbo, posGl[1])
		rigidBodyVbo = append(rigidBodyVbo, posGl[2])
	}
	{
		// 手前右下
		posGl := mat.MulVec3(&mmath.MVec3{size.GetX(), size.GetY(), -size.GetZ()}).GL()

		rigidBodyVbo = append(rigidBodyVbo, rigidBody.Bullet()...)
		rigidBodyVbo = append(rigidBodyVbo, posGl[0])
		rigidBodyVbo = append(rigidBodyVbo, posGl[1])
		rigidBodyVbo = append(rigidBodyVbo, posGl[2])
	}
	{
		// 手前左下
		posGl := mat.MulVec3(&mmath.MVec3{-size.GetX(), size.GetY(), -size.GetZ()}).GL()

		rigidBodyVbo = append(rigidBodyVbo, rigidBody.Bullet()...)
		rigidBodyVbo = append(rigidBodyVbo, posGl[0])
		rigidBodyVbo = append(rigidBodyVbo, posGl[1])
		rigidBodyVbo = append(rigidBodyVbo, posGl[2])
	}
	{
		// 奥左上
		posGl := mat.MulVec3(&mmath.MVec3{-size.GetX(), -size.GetY(), size.GetZ()}).GL()

		rigidBodyVbo = append(rigidBodyVbo, rigidBody.Bullet()...)
		rigidBodyVbo = append(rigidBodyVbo, posGl[0])
		rigidBodyVbo = append(rigidBodyVbo, posGl[1])
		rigidBodyVbo = append(rigidBodyVbo, posGl[2])
	}
	{
		// 奥右上
		posGl := mat.MulVec3(&mmath.MVec3{size.GetX(), -size.GetY(), size.GetZ()}).GL()

		rigidBodyVbo = append(rigidBodyVbo, rigidBody.Bullet()...)
		rigidBodyVbo = append(rigidBodyVbo, posGl[0])
		rigidBodyVbo = append(rigidBodyVbo, posGl[1])
		rigidBodyVbo = append(rigidBodyVbo, posGl[2])
	}
	{
		// 奥右下
		posGl := mat.MulVec3(&mmath.MVec3{size.GetX(), size.GetY(), size.GetZ()}).GL()

		rigidBodyVbo = append(rigidBodyVbo, rigidBody.Bullet()...)
		rigidBodyVbo = append(rigidBodyVbo, posGl[0])
		rigidBodyVbo = append(rigidBodyVbo, posGl[1])
		rigidBodyVbo = append(rigidBodyVbo, posGl[2])
	}
	{
		// 奥左下
		posGl := mat.MulVec3(&mmath.MVec3{-size.GetX(), size.GetY(), size.GetZ()}).GL()

		rigidBodyVbo = append(rigidBodyVbo, rigidBody.Bullet()...)
		rigidBodyVbo = append(rigidBodyVbo, posGl[0])
		rigidBodyVbo = append(rigidBodyVbo, posGl[1])
		rigidBodyVbo = append(rigidBodyVbo, posGl[2])
	}
	{
		// 手前左上-手前右上
		rigidBodyIbo = append(rigidBodyIbo, uint32(startIdx+0))
		rigidBodyIbo = append(rigidBodyIbo, uint32(startIdx+1))
	}
	{
		// 手前右上-手前右下
		rigidBodyIbo = append(rigidBodyIbo, uint32(startIdx+1))
		rigidBodyIbo = append(rigidBodyIbo, uint32(startIdx+2))
	}
	{
		// 手前右下-手前左下
		rigidBodyIbo = append(rigidBodyIbo, uint32(startIdx+2))
		rigidBodyIbo = append(rigidBodyIbo, uint32(startIdx+3))
	}
	{
		// 手前左下-手前左上
		rigidBodyIbo = append(rigidBodyIbo, uint32(startIdx+3))
		rigidBodyIbo = append(rigidBodyIbo, uint32(startIdx+0))
	}
	{
		// 奥左上-奥右上
		rigidBodyIbo = append(rigidBodyIbo, uint32(startIdx+4))
		rigidBodyIbo = append(rigidBodyIbo, uint32(startIdx+5))
	}
	{
		// 奥右上-奥右下
		rigidBodyIbo = append(rigidBodyIbo, uint32(startIdx+5))
		rigidBodyIbo = append(rigidBodyIbo, uint32(startIdx+6))
	}
	{
		// 奥右下-奥左下
		rigidBodyIbo = append(rigidBodyIbo, uint32(startIdx+6))
		rigidBodyIbo = append(rigidBodyIbo, uint32(startIdx+7))
	}
	{
		// 奥左下-奥左上
		rigidBodyIbo = append(rigidBodyIbo, uint32(startIdx+7))
		rigidBodyIbo = append(rigidBodyIbo, uint32(startIdx+4))
	}
	{
		// 手前左上-奥左上
		rigidBodyIbo = append(rigidBodyIbo, uint32(startIdx+0))
		rigidBodyIbo = append(rigidBodyIbo, uint32(startIdx+4))
	}
	{
		// 手前右上-奥右上
		rigidBodyIbo = append(rigidBodyIbo, uint32(startIdx+1))
		rigidBodyIbo = append(rigidBodyIbo, uint32(startIdx+5))
	}
	{
		// 手前右下-奥右下
		rigidBodyIbo = append(rigidBodyIbo, uint32(startIdx+2))
		rigidBodyIbo = append(rigidBodyIbo, uint32(startIdx+6))
	}
	{
		// 手前左下-奥左下
		rigidBodyIbo = append(rigidBodyIbo, uint32(startIdx+3))
		rigidBodyIbo = append(rigidBodyIbo, uint32(startIdx+7))
	}

	return rigidBodyVbo, rigidBodyIbo
}

// 半球の頂点を追加
func (r *RigidBodies) addCornerSphereHalf(
	rigidBody *RigidBody,
	rigidBodyVbo []float32,
	rigidBodyIbo []uint32,
	startIdx int,
	pos *mmath.MVec3,
	segments int,
	isUpper bool,
) ([]float32, []uint32) {
	// 球体を6x6に分割して線で描くため、球体の頂点を追加

	// isUpperがtrueの場合は上半球、falseの場合は下半球
	var mStart int
	var mEnd int
	if isUpper {
		mStart = 0
		mEnd = int(segments / 2)
	} else {
		mStart = int(segments / 2)
		mEnd = segments
	}

	for m := mStart; m <= mEnd; m++ {
		for n := 0; n <= segments; n++ {
			// Calculate the latitude and longitude
			lat := (float64(m) / float64(segments)) * math.Pi
			lon := (float64(n) / float64(segments)) * 2.0 * math.Pi

			// Calculate the x, y, z coordinates
			x := rigidBody.Size.GetX() * math.Sin(lat) * math.Cos(lon)
			y := rigidBody.Size.GetX() * math.Cos(lat)
			z := rigidBody.Size.GetX() * math.Sin(lat) * math.Sin(lon)
			spherePos := &mmath.MVec3{x, y, z}

			// Apply the rigid body's rotation to the vector
			rotation := rigidBody.Rotation.GetQuaternion().Copy()
			rotation.SetX(-rotation.GetX())
			rotation.SetY(-rotation.GetY())
			rotation.Normalize()
			rotatedPosGl := rotation.RotatedVec3(spherePos).Added(pos).GL()

			// Append the coordinates to the VBO
			rigidBodyVbo = append(rigidBodyVbo, rigidBody.Bullet()...)
			rigidBodyVbo = append(rigidBodyVbo, rotatedPosGl[0], rotatedPosGl[1], rotatedPosGl[2])
		}
	}

	for m := mStart; m <= mEnd; m++ {
		for n := 0; n <= segments; n++ {
			// Calculate the indices of the current point and the points to the right and below it
			idx := startIdx + m*(segments+1) + n
			idxRight := idx + 1
			idxBelow := idx + segments + 1

			// Add a line from the current point to the point to the right
			if n < segments {
				rigidBodyIbo = append(rigidBodyIbo, uint32(idx), uint32(idxRight))
			}

			// Add a line from the current point to the point below
			if m < segments {
				rigidBodyIbo = append(rigidBodyIbo, uint32(idx), uint32(idxBelow))
			}
		}
	}

	return rigidBodyVbo, rigidBodyIbo
}

func (r *RigidBodies) addCornerSphere(
	rigidBody *RigidBody,
	rigidBodyVbo []float32,
	rigidBodyIbo []uint32,
	startIdx int,
	segments int,
) ([]float32, []uint32) {
	// 上半球
	upperPos := rigidBody.Position.Copy()
	rigidBodyVbo, rigidBodyIbo = r.addCornerSphereHalf(
		rigidBody, rigidBodyVbo, rigidBodyIbo, startIdx, upperPos, segments, true)
	// 下半球
	lowerPos := rigidBody.Position.Copy()
	rigidBodyVbo, rigidBodyIbo = r.addCornerSphereHalf(
		rigidBody, rigidBodyVbo, rigidBodyIbo, startIdx, lowerPos, segments, false)

	return rigidBodyVbo, rigidBodyIbo
}

func (r *RigidBodies) addCornerCapsule(
	rigidBody *RigidBody,
	rigidBodyVbo []float32,
	rigidBodyIbo []uint32,
	startIdx int,
	segments int,
) ([]float32, []uint32) {
	// カプセルを6x6に分割して線で描くため、カプセルの頂点を追加
	rotation := rigidBody.Rotation.GetQuaternion().Copy()
	rotation.SetX(-rotation.GetX())
	rotation.SetY(-rotation.GetY())
	rotation.Normalize()

	// 上半球
	upperMat := mmath.NewMMat4()
	upperMat.Translate(rigidBody.Position)
	upperMat.Mul(rotation.ToMat4())
	upperPos := upperMat.MulVec3(&mmath.MVec3{0, rigidBody.Size.GetY() / 2, 0})
	rigidBodyVbo, rigidBodyIbo = r.addCornerSphereHalf(
		rigidBody, rigidBodyVbo, rigidBodyIbo, startIdx, upperPos, segments, true)
	// 下半球
	lowerMat := mmath.NewMMat4()
	lowerMat.Translate(rigidBody.Position)
	lowerMat.Mul(rotation.ToMat4())
	lowerPos := lowerMat.MulVec3(&mmath.MVec3{0, -rigidBody.Size.GetY() / 2, 0})
	rigidBodyVbo, rigidBodyIbo = r.addCornerSphereHalf(
		rigidBody, rigidBodyVbo, rigidBodyIbo, startIdx, lowerPos, segments, false)

	return rigidBodyVbo, rigidBodyIbo
}

func (r *RigidBodies) prepareDraw() {
	rigidBodyVbo := make([]float32, 0, len(r.Data))
	rigidBodyIbo := make([]uint32, 0, len(r.Data))

	startIdx := 0
	segments := 6

	for _, rigidBody := range r.GetSortedData() {
		if rigidBody.ShapeType == SHAPE_BOX {
			// 箱剛体
			rigidBodyVbo, rigidBodyIbo = r.addCornerBox(
				rigidBody, rigidBodyVbo, rigidBodyIbo, startIdx)
			startIdx += 8
		} else if rigidBody.ShapeType == SHAPE_SPHERE {
			// 球剛体
			rigidBodyVbo, rigidBodyIbo = r.addCornerSphere(
				rigidBody, rigidBodyVbo, rigidBodyIbo, startIdx, segments)
			startIdx += (segments+1)*(segments+1) + (segments + 1)
		} else if rigidBody.ShapeType == SHAPE_CAPSULE {
			// カプセル剛体
			rigidBodyVbo, rigidBodyIbo = r.addCornerCapsule(
				rigidBody, rigidBodyVbo, rigidBodyIbo, startIdx, segments)
			startIdx += (segments+1)*(segments+1) + (segments + 1)
		}
	}

	r.vao = mgl.NewVAO()
	r.vao.Bind()
	r.vbo = mgl.NewVBOForRigidBody(gl.Ptr(rigidBodyVbo), len(rigidBodyVbo))
	r.vbo.BindRigidBody()
	r.vbo.Unbind()
	r.vao.Unbind()

	r.ibo = mgl.NewIBO(gl.Ptr(rigidBodyIbo), len(rigidBodyIbo))
	r.count = int32(len(rigidBodyIbo))
}

func (r *RigidBodies) InitPhysics(physics *mphysics.MPhysics, bones *Bones) {
	// 剛体を順番にボーンと紐付けていく
	for _, rigidBody := range r.GetSortedData() {
		// 物理設定の初期化
		if rigidBody.BoneIndex >= 0 && bones.Contains(rigidBody.BoneIndex) {
			rigidBody.InitPhysics(physics, bones.GetItem(rigidBody.BoneIndex))
		} else {
			rigidBody.InitPhysics(physics, nil)
		}
	}

	// 剛体の描画準備
	r.prepareDraw()
}

func (r *RigidBodies) Draw(
	shader *mgl.MShader,
	boneMatrixes []*mgl32.Mat4,
	windowIndex int,
) {
	// ボーンをモデルメッシュの前面に描画するために深度テストを無効化
	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.ALWAYS)

	// ブレンディングを有効にする
	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)

	// ブレンディングを有効にする
	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)

	shader.UseRigidBodyProgram()

	// ボーンデフォームテクスチャ設定
	BindBoneMatrixes(boneMatrixes, shader, shader.RigidBodyProgram, windowIndex)

	// 箱剛体の描画

	r.vao.Bind()
	r.vbo.BindRigidBody()
	r.ibo.Bind()

	gl.DrawElements(gl.LINES, r.count, gl.UNSIGNED_INT, nil)

	r.ibo.Unbind()
	r.vbo.Unbind()
	r.vao.Unbind()

	UnbindBoneMatrixes()

	shader.Unuse()

	gl.Disable(gl.BLEND)
}
