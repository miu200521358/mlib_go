package vmd

import (
	"github.com/miu200521358/mlib_go/pkg/mmath"

)

type BoneTree struct {
	BoneName      string
	Frame         float32
	GlobalMatrix  *mmath.MMat4
	LocalMatrix   *mmath.MMat4
	Position      *mmath.MVec3
	FramePosition *mmath.MVec3
	FrameRotation *mmath.MQuaternion
	FrameScale    *mmath.MVec3
}

func NewBoneTree(
	boneName string,
	frame float32,
	globalMatrix, localMatrix *mmath.MMat4,
	framePosition *mmath.MVec3,
	frameRotation *mmath.MQuaternion,
	frameScale *mmath.MVec3,
) *BoneTree {
	p := globalMatrix.Translation()
	return &BoneTree{
		BoneName:      boneName,
		Frame:         frame,
		GlobalMatrix:  globalMatrix,
		LocalMatrix:   localMatrix,
		Position:      p,
		FramePosition: framePosition,
		FrameRotation: frameRotation,
		FrameScale:    frameScale,
	}
}

type BoneNameFrameNo struct {
	BoneName string
	Frame    float32
}

type BoneTrees struct {
	Data map[BoneNameFrameNo]*BoneTree
}

func NewBoneTrees() *BoneTrees {
	return &BoneTrees{
		Data: make(map[BoneNameFrameNo]*BoneTree, 0),
	}
}

func (bts *BoneTrees) GetItem(boneName string, frame float32) *BoneTree {
	return bts.Data[BoneNameFrameNo{boneName, frame}]
}

func (bts *BoneTrees) SetItem(boneName string, frame float32, boneTree *BoneTree) {
	bts.Data[BoneNameFrameNo{boneName, frame}] = boneTree
}

func (bts *BoneTrees) GetBoneNames() []string {
	boneNames := make([]string, 0)
	for key := range bts.Data {
		boneNames = append(boneNames, key.BoneName)
	}
	return boneNames
}

func (bts *BoneTrees) GetFrameNos() []float32 {
	frames := make([]float32, 0)
	for key := range bts.Data {
		frames = append(frames, key.Frame)
	}
	return frames
}

func (bts *BoneTrees) Contains(boneName string, frame float32) bool {
	_, ok := bts.Data[BoneNameFrameNo{boneName, frame}]
	return ok
}

// func (bts *BoneTrees) TransformSDEF(model *pmx.PmxModel, vertex *pmx.Vertex, frame float32) *mgl32.Vec3 {
// 	sdef := vertex.Deform.(*pmx.Sdef)

// 	// R0
// 	w0 := vertex.Deform.GetAllWeights()[0]
// 	bone0 := model.Bones.GetItem(vertex.Deform.GetAllIndexes()[0])
// 	mat0 := bts.GetItem(bone0.Name, frame).LocalMatrix

// 	// R1
// 	// w1 := 1 - w0
// 	bone1 := model.Bones.GetItem(vertex.Deform.GetAllIndexes()[1])
// 	mat1 := bts.GetItem(bone1.Name, frame).LocalMatrix

// 	vecCinB0 := sdef.SdefC.Subed(bone0.Position)
// 	vecCinB1 := sdef.SdefC.Subed(bone1.Position)

// 	vecR0inB0 := sdef.SdefR0.Subed(bone0.Position)
// 	vecR1inB1 := sdef.SdefR1.Subed(bone1.Position)

// 	// R0/R1影響係数算出: 近い方に強く影響される
// 	var r1Bias float64
// 	len0 := vecCinB0.Distance(&vecR0inB0)
// 	len1 := vecCinB1.Distance(&vecR1inB1)
// 	if (len0 > 0.0 && len1 == 0.0) || (len0+len1 <= 0.0) {
// 		r1Bias = 0.0
// 	} else if len0 == 0.0 && len1 > 0.0 {
// 		r1Bias = 1.0
// 	} else {
// 		r1Bias = mmath.ClampFloat(len0/(len0+len1), 0.0, 1.0)
// 	}
// 	r0Bias := 1.0 - r1Bias

// 	// ------------------------
// 	// SDEF回転計算

// 	// // bone0の姿勢（q0）からbone1への相対回転（q1)をウェイト値でslerpをかけて求めている
// 	// // その後q0にq1を加えることで最終的な回転を維持しつつ、slerpのフリップ特性を得る
// 	// q0 := mat0.Quaternion()
// 	// q1 := mat1.Quaternion()
// 	// q0Toq1 := q0.Slerp(&q1, w0*r0Bias)

// 	// // 最終的なSDEF回転行列
// 	// matR := q0Toq1.ToMat4()

// 	// // 回転行列からスケール成分を除去する
// 	// scaleX := mmath.MVec3{matR[0][0], matR[0][1], matR[0][2]}
// 	// scaleY := mmath.MVec3{matR[1][0], matR[1][1], matR[1][2]}
// 	// scaleZ := mmath.MVec3{matR[2][0], matR[2][1], matR[2][2]}

// 	// sx := 1.0 / scaleX.Length()
// 	// sy := 1.0 / scaleY.Length()
// 	// sz := 1.0 / scaleZ.Length()

// 	// matR[0][0] *= sx
// 	// matR[0][1] *= sx
// 	// matR[0][2] *= sx
// 	// matR[1][0] *= sy
// 	// matR[1][1] *= sy
// 	// matR[1][2] *= sy
// 	// matR[2][0] *= sz
// 	// matR[2][1] *= sz
// 	// matR[2][2] *= sz
// 	// // ------------------------

// 	// // 変形後の交点Cの位置姿勢中間値
// 	// matP01 := vecCinB0.ToMat4().Muled(mat0)
// 	// matP02 := matP01.MuledFactor(w0)
// 	// vecP0 := matP02.Translation()

// 	// matP11 := vecCinB1.ToMat4().Muled(mat1)
// 	// matP12 := matP11.MuledFactor(w1)
// 	// vecP1 := matP12.Translation()

// 	// vecMedianC := vecP0.Added(&vecP1)

// 	// // ------------------------

// 	// // 補間点R0/R1をBDEF2移動させて交点Cを補正する
// 	// matCR0 := vecR0inB0.ToMat4().Muled(mat0)
// 	// matCR1 := vecR1inB1.ToMat4().Muled(mat1)

// 	// vecCR0 := matCR0.Translation()
// 	// vecCR1 := matCR1.Translation()

// 	// vecCR01 := vecCR0.MuledScalar(r0Bias)
// 	// vecCR11 := vecCR1.MuledScalar(r1Bias)

// 	// vecCR02 := vecCR01.Added(&vecCR11)
// 	// vecCR12 := vecMedianC.Added(&vecCR02)

// 	// vecFinalC := vecCR12.MuledScalar(0.5)

// 	// // ------------------------

// 	// // 最終計算
// 	// cp0Vec := bone0.NormalizedFixedAxis.ToLocalMatrix4x4().Translation()
// 	// vecCP := cp0Vec.Subed(&vecCinB0)

// 	// // vecCPM := vecCP.MuledScalar(w0)

// 	// scale0X := mmath.MVec3{mat0[0][0], mat0[0][1], mat0[0][2]}
// 	// scale0Y := mmath.MVec3{mat0[1][0], mat0[1][1], mat0[1][2]}
// 	// scale0Z := mmath.MVec3{mat0[2][0], mat0[2][1], mat0[2][2]}
// 	// scale0 := mmath.MVec3{scale0X.Length(), scale0Y.Length(), scale0Z.Length()}

// 	// scale1X := mmath.MVec3{mat1[0][0], mat1[0][1], mat1[0][2]}
// 	// scale1Y := mmath.MVec3{mat1[1][0], mat1[1][1], mat1[1][2]}
// 	// scale1Z := mmath.MVec3{mat1[2][0], mat1[2][1], mat1[2][2]}
// 	// scale1 := mmath.MVec3{scale1X.Length(), scale1Y.Length(), scale1Z.Length()}

// 	// scaleW0 := scale0.MuledScalar(w0)
// 	// scaleW1 := scale1.MuledScalar(w1)
// 	// vecMatS := scaleW0.Added(&scaleW1)

// 	// matS := mmath.NewMMat4()
// 	// matS[0][0] = vecMatS.GetX()
// 	// matS[1][1] = vecMatS.GetY()
// 	// matS[2][2] = vecMatS.GetZ()

// 	// matCPR := vecCP.ToMat4().Muled(&matR)
// 	// matCPRS := matCPR.Muled(matS)
// 	// vecCPRS := matCPRS.Translation()

// 	vertexPosition := matR.MulVec3(vertex.Position)
// 	vertexPositionGL := vertexPosition.GL()

// 	return &vertexPositionGL
// }
