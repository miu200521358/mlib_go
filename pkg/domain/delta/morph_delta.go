// 指示: miu200521358
package delta

import (
	"github.com/miu200521358/mlib_go/pkg/domain/mmath"
	"github.com/miu200521358/mlib_go/pkg/domain/model"
)

// VertexMorphDelta は頂点モーフ差分を表す。
type VertexMorphDelta struct {
	Index         int
	Position      *mmath.Vec3
	Uv            *mmath.Vec2
	Uv1           *mmath.Vec2
	AfterPosition *mmath.Vec3
}

// NewVertexMorphDelta はVertexMorphDeltaを生成する。
func NewVertexMorphDelta(index int) *VertexMorphDelta {
	return &VertexMorphDelta{Index: index}
}

// IsZero はゼロ差分か判定する。
func (d *VertexMorphDelta) IsZero() bool {
	return d == nil ||
		((d.Position == nil || d.Position.NearEquals(mmath.ZERO_VEC3, 1e-4)) &&
			(d.Uv == nil || d.Uv.NearEquals(mmath.ZERO_VEC2, 1e-4)) &&
			(d.Uv1 == nil || d.Uv1.NearEquals(mmath.ZERO_VEC2, 1e-4)) &&
			(d.AfterPosition == nil || d.AfterPosition.NearEquals(mmath.ZERO_VEC3, 1e-4)))
}

// BoneMorphDelta はボーンモーフ差分を表す。
type BoneMorphDelta struct {
	BoneIndex               int
	FramePosition           *mmath.Vec3
	FrameCancelablePosition *mmath.Vec3
	FrameRotation           *mmath.Quaternion
	FrameCancelableRotation *mmath.Quaternion
	FrameScale              *mmath.Vec3
	FrameCancelableScale    *mmath.Vec3
	FrameLocalMat           *mmath.Mat4
}

// NewBoneMorphDelta はBoneMorphDeltaを生成する。
func NewBoneMorphDelta(boneIndex int) *BoneMorphDelta {
	return &BoneMorphDelta{BoneIndex: boneIndex}
}

// FilledMorphPosition はモーフ位置を返す。
func (d *BoneMorphDelta) FilledMorphPosition() mmath.Vec3 {
	if d == nil || d.FramePosition == nil {
		return mmath.NewVec3()
	}
	return *d.FramePosition
}

// FilledMorphCancelablePosition はモーフキャンセル位置を返す。
func (d *BoneMorphDelta) FilledMorphCancelablePosition() mmath.Vec3 {
	if d == nil || d.FrameCancelablePosition == nil {
		return mmath.NewVec3()
	}
	return *d.FrameCancelablePosition
}

// FilledMorphRotation はモーフ回転を返す。
func (d *BoneMorphDelta) FilledMorphRotation() mmath.Quaternion {
	if d == nil || d.FrameRotation == nil {
		return mmath.NewQuaternion()
	}
	return *d.FrameRotation
}

// FilledMorphCancelableRotation はモーフキャンセル回転を返す。
func (d *BoneMorphDelta) FilledMorphCancelableRotation() mmath.Quaternion {
	if d == nil || d.FrameCancelableRotation == nil {
		return mmath.NewQuaternion()
	}
	return *d.FrameCancelableRotation
}

// FilledMorphScale はモーフスケールを返す。
func (d *BoneMorphDelta) FilledMorphScale() mmath.Vec3 {
	if d == nil || d.FrameScale == nil {
		return mmath.ONE_VEC3
	}
	return *d.FrameScale
}

// FilledMorphCancelableScale はモーフキャンセルスケールを返す。
func (d *BoneMorphDelta) FilledMorphCancelableScale() mmath.Vec3 {
	if d == nil || d.FrameCancelableScale == nil {
		return mmath.ONE_VEC3
	}
	return *d.FrameCancelableScale
}

// FilledMorphLocalMat はモーフローカル行列を返す。
func (d *BoneMorphDelta) FilledMorphLocalMat() mmath.Mat4 {
	if d == nil || d.FrameLocalMat == nil {
		return mmath.NewMat4()
	}
	return *d.FrameLocalMat
}

// Copy は差分を複製する。
func (d *BoneMorphDelta) Copy() (BoneMorphDelta, error) {
	if d == nil {
		return BoneMorphDelta{}, nil
	}
	copyDelta := BoneMorphDelta{BoneIndex: d.BoneIndex}
	{
		pos := d.FilledMorphPosition()
		copyDelta.FramePosition = &pos
	}
	{
		pos := d.FilledMorphCancelablePosition()
		copyDelta.FrameCancelablePosition = &pos
	}
	{
		rot := d.FilledMorphRotation()
		copyDelta.FrameRotation = &rot
	}
	{
		rot := d.FilledMorphCancelableRotation()
		copyDelta.FrameCancelableRotation = &rot
	}
	{
		scale := d.FilledMorphScale()
		copyDelta.FrameScale = &scale
	}
	{
		scale := d.FilledMorphCancelableScale()
		copyDelta.FrameCancelableScale = &scale
	}
	{
		mat := d.FilledMorphLocalMat()
		copyDelta.FrameLocalMat = &mat
	}
	return copyDelta, nil
}

// MaterialMorphDelta は材質モーフ差分を表す。
type MaterialMorphDelta struct {
	Material    model.Material
	AddMaterial model.Material
	MulMaterial model.Material
}

// NewMaterialMorphDelta はMaterialMorphDeltaを生成する。
func NewMaterialMorphDelta(material *model.Material) *MaterialMorphDelta {
	base := model.NewMaterial()
	if material != nil {
		base = copyMaterial(material)
	}
	add := model.NewMaterial()
	mul := model.NewMaterial()
	mul.Diffuse = mmath.Vec4{X: 1, Y: 1, Z: 1, W: 1}
	mul.Specular = mmath.Vec4{X: 1, Y: 1, Z: 1, W: 1}
	mul.Ambient = mmath.ONE_VEC3
	mul.Edge = mmath.Vec4{X: 1, Y: 1, Z: 1, W: 1}
	mul.EdgeSize = 1
	mul.TextureFactor = mmath.ONE_VEC4
	mul.SphereTextureFactor = mmath.ONE_VEC4
	mul.ToonTextureFactor = mmath.ONE_VEC4
	return &MaterialMorphDelta{
		Material:    *base,
		AddMaterial: *add,
		MulMaterial: *mul,
	}
}

// Add は加算差分を反映する。
func (d *MaterialMorphDelta) Add(offset *model.MaterialMorphOffset, ratio float64) {
	if d == nil || offset == nil {
		return
	}
	d.AddMaterial.Diffuse = d.AddMaterial.Diffuse.Added(offset.Diffuse.MuledScalar(ratio))
	d.AddMaterial.Specular = d.AddMaterial.Specular.Added(offset.Specular.MuledScalar(ratio))
	d.AddMaterial.Ambient = d.AddMaterial.Ambient.Added(offset.Ambient.MuledScalar(ratio))
	d.AddMaterial.Edge = d.AddMaterial.Edge.Added(offset.Edge.MuledScalar(ratio))
	d.AddMaterial.EdgeSize += offset.EdgeSize * ratio
	d.AddMaterial.TextureFactor = d.AddMaterial.TextureFactor.Added(offset.TextureFactor.MuledScalar(ratio))
	d.AddMaterial.SphereTextureFactor = d.AddMaterial.SphereTextureFactor.Added(offset.SphereTextureFactor.MuledScalar(ratio))
	d.AddMaterial.ToonTextureFactor = d.AddMaterial.ToonTextureFactor.Added(offset.ToonTextureFactor.MuledScalar(ratio))
}

// Mul は乗算差分を反映する。
func (d *MaterialMorphDelta) Mul(offset *model.MaterialMorphOffset, ratio float64) {
	if d == nil || offset == nil {
		return
	}
	d.MulMaterial.Diffuse = d.MulMaterial.Diffuse.Muled(lerpVec4(offset.Diffuse, ratio))
	d.MulMaterial.Specular = d.MulMaterial.Specular.Muled(lerpVec4(offset.Specular, ratio))
	d.MulMaterial.Ambient = d.MulMaterial.Ambient.Muled(lerpVec3(offset.Ambient, ratio))
	d.MulMaterial.Edge = d.MulMaterial.Edge.Muled(lerpVec4(offset.Edge, ratio))
	d.MulMaterial.EdgeSize *= mmath.Lerp(1, offset.EdgeSize, ratio)
	d.MulMaterial.TextureFactor = d.MulMaterial.TextureFactor.Muled(lerpVec4(offset.TextureFactor, ratio))
	d.MulMaterial.SphereTextureFactor = d.MulMaterial.SphereTextureFactor.Muled(lerpVec4(offset.SphereTextureFactor, ratio))
	d.MulMaterial.ToonTextureFactor = d.MulMaterial.ToonTextureFactor.Muled(lerpVec4(offset.ToonTextureFactor, ratio))
}

// copyMaterial は材質を複製する。
func copyMaterial(src *model.Material) *model.Material {
	dst := model.NewMaterial()
	if src == nil {
		return dst
	}
	// index/name はメソッド経由で複製する。
	dst.SetIndex(src.Index())
	dst.SetName(src.Name())
	dst.EnglishName = src.EnglishName
	dst.Memo = src.Memo
	dst.Diffuse = src.Diffuse
	dst.Specular = src.Specular
	dst.Ambient = src.Ambient
	dst.DrawFlag = src.DrawFlag
	dst.Edge = src.Edge
	dst.EdgeSize = src.EdgeSize
	dst.TextureFactor = src.TextureFactor
	dst.SphereTextureFactor = src.SphereTextureFactor
	dst.ToonTextureFactor = src.ToonTextureFactor
	dst.TextureIndex = src.TextureIndex
	dst.SphereTextureIndex = src.SphereTextureIndex
	dst.SphereMode = src.SphereMode
	dst.ToonSharingFlag = src.ToonSharingFlag
	dst.ToonTextureIndex = src.ToonTextureIndex
	dst.VerticesCount = src.VerticesCount
	return dst
}

// lerpVec3 は乗算用の補間ベクトルを返す。
func lerpVec3(v mmath.Vec3, ratio float64) mmath.Vec3 {
	out := mmath.NewVec3()
	out.X = mmath.Lerp(1, v.X, ratio)
	out.Y = mmath.Lerp(1, v.Y, ratio)
	out.Z = mmath.Lerp(1, v.Z, ratio)
	return out
}

// lerpVec4 は乗算用の補間ベクトルを返す。
func lerpVec4(v mmath.Vec4, ratio float64) mmath.Vec4 {
	return mmath.Vec4{
		X: mmath.Lerp(1, v.X, ratio),
		Y: mmath.Lerp(1, v.Y, ratio),
		Z: mmath.Lerp(1, v.Z, ratio),
		W: mmath.Lerp(1, v.W, ratio),
	}
}
