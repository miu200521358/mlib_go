// 指示: miu200521358
package deform

import (
	"math"

	"github.com/miu200521358/mlib_go/pkg/domain/delta"
	"github.com/miu200521358/mlib_go/pkg/domain/mmath"
	"github.com/miu200521358/mlib_go/pkg/domain/model"
	"github.com/miu200521358/mlib_go/pkg/domain/model/collection"
)

// ApplySkinning はスキニングを適用し頂点/法線を更新する。
func ApplySkinning(
	vertices *collection.IndexedCollection[*model.Vertex],
	boneDeltas *delta.BoneDeltas,
	morphDeltas *delta.MorphDeltas,
) {
	if vertices == nil || boneDeltas == nil {
		return
	}
	for _, vertex := range vertices.Values() {
		if vertex == nil || vertex.Deform == nil {
			continue
		}
		mat := skinningMatrix(vertex.Deform, boneDeltas)
		morphDelta := vertexMorphDelta(vertex, morphDeltas)

		pos := vertex.Position
		if morphDelta != nil && morphDelta.Position != nil {
			pos = pos.Added(*morphDelta.Position)
		}
		pos = mat.MulVec3(pos)
		if morphDelta != nil && morphDelta.AfterPosition != nil {
			pos = pos.Added(*morphDelta.AfterPosition)
		}
		vertex.Position = pos

		normal := mat.MulVec3(vertex.Normal)
		vertex.Normal = normal.Normalized()

		if sdef, ok := vertex.Deform.(*model.Sdef); ok {
			bone0 := boneDeltas.Get(sdef.Indexes()[0])
			bone1 := boneDeltas.Get(sdef.Indexes()[1])
			if bone0 == nil || bone1 == nil {
				continue
			}
			c, r0, r1 := RecomputeSdef(bone0.FilledGlobalPosition(), bone1.FilledGlobalPosition(), vertex.Position)
			sdef.SdefC = c
			sdef.SdefR0 = r0
			sdef.SdefR1 = r1
		}
	}
}

// RecomputeSdef はSDEFの再計算結果を返す。
func RecomputeSdef(bone0Global, bone1Global, vertexPos mmath.Vec3) (mmath.Vec3, mmath.Vec3, mmath.Vec3) {
	if isInvalidVec3(bone0Global) {
		bone0Global = vertexPos
	}
	if isInvalidVec3(bone1Global) {
		bone1Global = vertexPos
	}
	c := mmath.IntersectLinePoint(bone0Global, bone1Global, vertexPos)
	if isInvalidVec3(c) {
		c = vertexPos
	}
	r0 := bone0Global.Added(c).MuledScalar(0.5)
	r1 := bone1Global.Added(c).MuledScalar(0.5)
	return c, r0, r1
}

// skinningMatrix はウェイト合成行列を返す。
func skinningMatrix(deform model.IDeform, boneDeltas *delta.BoneDeltas) mmath.Mat4 {
	if deform == nil || boneDeltas == nil {
		return mmath.ZERO_MAT4
	}
	indexes := deform.Indexes()
	weights := deform.Weights()
	mat := mmath.ZERO_MAT4
	for i, boneIndex := range indexes {
		if i >= len(weights) {
			break
		}
		weight := weights[i]
		if weight == 0 {
			continue
		}
		boneDelta := boneDeltas.Get(boneIndex)
		if boneDelta == nil {
			// 差分が無いボーンは単位行列として扱う。
			mat = mat.Added(mmath.IDENT_MAT4.MuledScalar(weight))
			continue
		}
		boneMat := boneDelta.FilledLocalMatrix().MuledScalar(weight)
		mat = mat.Added(boneMat)
	}
	return mat
}

// vertexMorphDelta は頂点モーフ差分を返す。
func vertexMorphDelta(vertex *model.Vertex, morphDeltas *delta.MorphDeltas) *delta.VertexMorphDelta {
	if vertex == nil || morphDeltas == nil || morphDeltas.Vertices() == nil {
		return nil
	}
	return morphDeltas.Vertices().Get(vertex.Index())
}

func isInvalidFloat(f float64) bool {
	return math.IsNaN(f) || math.IsInf(f, 0)
}

// isInvalidVec3 はNaN/Inf判定を返す。
func isInvalidVec3(v mmath.Vec3) bool {
	return math.IsNaN(v.X) || math.IsNaN(v.Y) || math.IsNaN(v.Z) ||
		math.IsInf(v.X, 0) || math.IsInf(v.Y, 0) || math.IsInf(v.Z, 0)
}
