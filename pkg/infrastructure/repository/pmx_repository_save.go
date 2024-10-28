package repository

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
	"runtime"
	"unicode/utf16"

	"github.com/miu200521358/mlib_go/pkg/domain/core"
	"github.com/miu200521358/mlib_go/pkg/domain/mmath"
	"github.com/miu200521358/mlib_go/pkg/domain/pmx"
	"github.com/miu200521358/mlib_go/pkg/mutils"
	"github.com/miu200521358/mlib_go/pkg/mutils/mi18n"
	"github.com/miu200521358/mlib_go/pkg/mutils/mlog"
)

func (rep *PmxRepository) Save(overridePath string, data core.IHashModel, includeSystem bool) error {
	runtime.GOMAXPROCS(int(runtime.NumCPU()))
	defer runtime.GOMAXPROCS(max(1, int(runtime.NumCPU()/4)))

	mlog.IL(mi18n.T("保存開始", map[string]interface{}{"Type": "Pmx", "Path": overridePath}))
	defer mlog.I(mi18n.T("保存終了", map[string]interface{}{"Type": "Pmx"}))

	model := data.(*pmx.PmxModel)

	path := model.Path()
	// 保存可能なパスである場合、上書き
	if mutils.CanSave(overridePath) {
		path = overridePath
	}

	// Open the output file
	fout, err := os.Create(path)
	if err != nil {
		return err
	}
	defer fout.Close()

	filteredBones := []*pmx.Bone{}
	for i := range model.Bones.Len() {
		bone := model.Bones.Get(i)
		if (!includeSystem && !bone.IsSystem) || includeSystem {
			filteredBones = append(filteredBones, bone)
		}
	}

	filteredMorphs := []*pmx.Morph{}
	for i := range model.Morphs.Len() {
		morph := model.Morphs.Get(i)
		if (!includeSystem && !morph.IsSystem) || includeSystem {
			filteredMorphs = append(filteredMorphs, morph)
		}
	}

	_, err = fout.Write([]byte("PMX "))
	if err != nil {
		return fmt.Errorf("failed to write PMX signature: %v", err)
	}

	err = rep.writeNumber(fout, binaryType_float, 2.0, 0.0, true)
	if err != nil {
		return err
	}

	err = rep.writeByte(fout, 8, true)
	if err != nil {
		return err
	}

	err = rep.writeByte(fout, 0, true)
	if err != nil {
		return err
	}

	err = rep.writeByte(fout, model.ExtendedUVCount, true)
	if err != nil {
		return err
	}

	vertexIdxSize, vertexIdxType := rep.defineWriteIndexForVertex(model.Vertices.Len())
	err = rep.writeByte(fout, vertexIdxSize, true)
	if err != nil {
		return err
	}

	textureIdxSize, textureIdxType := rep.defineWriteIndexForOthers(model.Textures.Len())
	err = rep.writeByte(fout, textureIdxSize, true)
	if err != nil {
		return err
	}

	materialIdxSize, materialIdxType := rep.defineWriteIndexForOthers(model.Materials.Len())
	err = rep.writeByte(fout, materialIdxSize, true)
	if err != nil {
		return err
	}

	boneIdxSize, boneIdxType := rep.defineWriteIndexForOthers(len(filteredBones))
	err = rep.writeByte(fout, boneIdxSize, true)
	if err != nil {
		return err
	}

	morphIdxSize, morphIdxType := rep.defineWriteIndexForOthers(len(filteredMorphs))
	err = rep.writeByte(fout, morphIdxSize, true)
	if err != nil {
		return err
	}

	rigidbodyIdxSize, rigidbodyIdxType := rep.defineWriteIndexForOthers(model.RigidBodies.Len())
	err = rep.writeByte(fout, rigidbodyIdxSize, true)
	if err != nil {
		return err
	}

	err = rep.writeText(fout, model.Name(), "Pmx Model")
	if err != nil {
		return err
	}

	err = rep.writeText(fout, model.EnglishName(), "Pmx Model")
	if err != nil {
		return err
	}

	err = rep.writeText(fout, model.Comment, "")
	if err != nil {
		return err
	}

	err = rep.writeText(fout, model.EnglishComment, "")
	if err != nil {
		return err
	}

	err = rep.saveVertices(fout, model, boneIdxType)
	if err != nil {
		return err
	}

	err = rep.saveFaces(fout, model, vertexIdxType)
	if err != nil {
		return err
	}

	err = rep.saveTextures(fout, model)
	if err != nil {
		return err
	}

	err = rep.saveMaterials(fout, model, textureIdxType)
	if err != nil {
		return err
	}

	err = rep.saveBones(fout, filteredBones, boneIdxType)
	if err != nil {
		return err
	}

	err = rep.saveMorphs(fout, filteredMorphs, vertexIdxType, boneIdxType, materialIdxType, morphIdxType)
	if err != nil {
		return err
	}

	err = rep.saveDisplaySlots(fout, model, boneIdxType, morphIdxType)
	if err != nil {
		return err
	}

	err = rep.saveRigidBodies(fout, model, boneIdxType)
	if err != nil {
		return err
	}

	err = rep.saveJoints(fout, model, rigidbodyIdxType)
	if err != nil {
		return err
	}

	return nil
}

// saveVertices 頂点データの書き込み
func (rep *PmxRepository) saveVertices(fout *os.File, model *pmx.PmxModel, boneIdxType binaryType) error {
	defer mlog.I(mi18n.T("保存途中完了", map[string]interface{}{"Type": mi18n.T("頂点")}))

	err := rep.writeNumber(fout, binaryType_int, float64(model.Vertices.Len()), 0.0, true)
	if err != nil {
		return err
	}

	for i := range model.Vertices.Len() {
		vertex := model.Vertices.Get(i)

		rep.writeNumber(fout, binaryType_float, vertex.Position.X, 0.0, false)
		rep.writeNumber(fout, binaryType_float, vertex.Position.Y, 0.0, false)
		rep.writeNumber(fout, binaryType_float, vertex.Position.Z, 0.0, false)
		rep.writeNumber(fout, binaryType_float, vertex.Normal.X, 0.0, false)
		rep.writeNumber(fout, binaryType_float, vertex.Normal.Y, 0.0, false)
		rep.writeNumber(fout, binaryType_float, vertex.Normal.Z, 0.0, false)
		rep.writeNumber(fout, binaryType_float, vertex.Uv.X, 0.0, false)
		rep.writeNumber(fout, binaryType_float, vertex.Uv.Y, 0.0, false)

		for _, uv := range vertex.ExtendedUvs {
			rep.writeNumber(fout, binaryType_float, uv.X, 0.0, false)
			rep.writeNumber(fout, binaryType_float, uv.Y, 0.0, false)
			rep.writeNumber(fout, binaryType_float, uv.Z, 0.0, false)
			rep.writeNumber(fout, binaryType_float, uv.W, 0.0, false)
		}

		for j := len(vertex.ExtendedUvs); j < model.ExtendedUVCount; j++ {
			rep.writeNumber(fout, binaryType_float, 0.0, 0.0, false)
			rep.writeNumber(fout, binaryType_float, 0.0, 0.0, false)
			rep.writeNumber(fout, binaryType_float, 0.0, 0.0, false)
			rep.writeNumber(fout, binaryType_float, 0.0, 0.0, false)
		}

		switch v := vertex.Deform.(type) {
		case *pmx.Bdef1:
			rep.writeByte(fout, 0, false)
			rep.writeNumber(fout, boneIdxType, float64(v.AllIndexes()[0]), 0.0, false)
		case *pmx.Bdef2:
			rep.writeByte(fout, 1, false)
			rep.writeNumber(fout, boneIdxType, float64(v.AllIndexes()[0]), 0.0, false)
			rep.writeNumber(fout, boneIdxType, float64(v.AllIndexes()[1]), 0.0, false)
			rep.writeNumber(fout, binaryType_float, v.AllWeights()[0], 0.0, true)
		case *pmx.Bdef4:
			rep.writeByte(fout, 2, false)
			rep.writeNumber(fout, boneIdxType, float64(v.AllIndexes()[0]), 0.0, false)
			rep.writeNumber(fout, boneIdxType, float64(v.AllIndexes()[1]), 0.0, false)
			rep.writeNumber(fout, boneIdxType, float64(v.AllIndexes()[2]), 0.0, false)
			rep.writeNumber(fout, boneIdxType, float64(v.AllIndexes()[3]), 0.0, false)
			rep.writeNumber(fout, binaryType_float, v.AllWeights()[0], 0.0, true)
			rep.writeNumber(fout, binaryType_float, v.AllWeights()[1], 0.0, true)
			rep.writeNumber(fout, binaryType_float, v.AllWeights()[2], 0.0, true)
			rep.writeNumber(fout, binaryType_float, v.AllWeights()[3], 0.0, true)
		case *pmx.Sdef:
			rep.writeByte(fout, 3, false)
			rep.writeNumber(fout, boneIdxType, float64(v.AllIndexes()[0]), 0.0, false)
			rep.writeNumber(fout, boneIdxType, float64(v.AllIndexes()[1]), 0.0, false)
			rep.writeNumber(fout, binaryType_float, v.AllWeights()[0], 0.0, true)
			rep.writeNumber(fout, binaryType_float, v.SdefC.X, 0.0, false)
			rep.writeNumber(fout, binaryType_float, v.SdefC.Y, 0.0, false)
			rep.writeNumber(fout, binaryType_float, v.SdefC.Z, 0.0, false)
			rep.writeNumber(fout, binaryType_float, v.SdefR0.X, 0.0, false)
			rep.writeNumber(fout, binaryType_float, v.SdefR0.Y, 0.0, false)
			rep.writeNumber(fout, binaryType_float, v.SdefR0.Z, 0.0, false)
			rep.writeNumber(fout, binaryType_float, v.SdefR1.X, 0.0, false)
			rep.writeNumber(fout, binaryType_float, v.SdefR1.Y, 0.0, false)
			rep.writeNumber(fout, binaryType_float, v.SdefR1.Z, 0.0, false)
		default:
			mlog.W("頂点deformなし: %v\n", vertex)
		}

		rep.writeNumber(fout, binaryType_float, vertex.EdgeFactor, 0.0, true)
	}

	return nil
}

// saveFaces 面データの書き込み
func (rep *PmxRepository) saveFaces(fout *os.File, model *pmx.PmxModel, vertexIdxType binaryType) error {
	defer mlog.I(mi18n.T("保存途中完了", map[string]interface{}{"Type": mi18n.T("面")}))

	err := rep.writeNumber(fout, binaryType_int, float64(model.Faces.Len()*3), 0.0, true)
	if err != nil {
		return err
	}

	for i := range model.Faces.Len() {
		face := model.Faces.Get(i)

		for _, vidx := range face.VertexIndexes {
			err = rep.writeNumber(fout, vertexIdxType, float64(vidx), 0.0, true)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// saveTextures テクスチャデータの書き込み
func (rep *PmxRepository) saveTextures(fout *os.File, model *pmx.PmxModel) error {
	err := rep.writeNumber(fout, binaryType_int, float64(model.Textures.Len()), 0.0, true)
	if err != nil {
		return err
	}

	for i := range model.Textures.Len() {
		texture := model.Textures.Get(i)
		err = rep.writeText(fout, texture.Name(), "")
		if err != nil {
			return err
		}
	}

	return nil
}

// saveMaterials 材質データの書き込み
func (rep *PmxRepository) saveMaterials(fout *os.File, model *pmx.PmxModel, textureIdxType binaryType) error {
	defer mlog.I(mi18n.T("保存途中完了", map[string]interface{}{"Type": mi18n.T("材質")}))

	err := rep.writeNumber(fout, binaryType_int, float64(model.Materials.Len()), 0.0, true)
	if err != nil {
		return err
	}

	for i := range model.Materials.Len() {
		material := model.Materials.Get(i)
		err = rep.writeText(fout, material.Name(), fmt.Sprintf("Material %d", i))
		if err != nil {
			return err
		}
		err = rep.writeText(fout, material.EnglishName(), fmt.Sprintf("Material %d", i))
		if err != nil {
			return err
		}
		err = rep.writeNumber(fout, binaryType_float, material.Diffuse.X, 0.0, true)
		if err != nil {
			return err
		}
		err = rep.writeNumber(fout, binaryType_float, material.Diffuse.Y, 0.0, true)
		if err != nil {
			return err
		}
		err = rep.writeNumber(fout, binaryType_float, material.Diffuse.Z, 0.0, true)
		if err != nil {
			return err
		}
		err = rep.writeNumber(fout, binaryType_float, material.Diffuse.W, 0.0, true)
		if err != nil {
			return err
		}
		err = rep.writeNumber(fout, binaryType_float, material.Specular.X, 0.0, true)
		if err != nil {
			return err
		}
		err = rep.writeNumber(fout, binaryType_float, material.Specular.Y, 0.0, true)
		if err != nil {
			return err
		}
		err = rep.writeNumber(fout, binaryType_float, material.Specular.Z, 0.0, true)
		if err != nil {
			return err
		}
		err = rep.writeNumber(fout, binaryType_float, material.Specular.W, 0.0, true)
		if err != nil {
			return err
		}
		err = rep.writeNumber(fout, binaryType_float, material.Ambient.X, 0.0, true)
		if err != nil {
			return err
		}
		err = rep.writeNumber(fout, binaryType_float, material.Ambient.Y, 0.0, true)
		if err != nil {
			return err
		}
		err = rep.writeNumber(fout, binaryType_float, material.Ambient.Z, 0.0, true)
		if err != nil {
			return err
		}
		err = rep.writeByte(fout, int(material.DrawFlag), true)
		if err != nil {
			return err
		}
		err = rep.writeNumber(fout, binaryType_float, material.Edge.X, 0.0, true)
		if err != nil {
			return err
		}
		err = rep.writeNumber(fout, binaryType_float, material.Edge.Y, 0.0, true)
		if err != nil {
			return err
		}
		err = rep.writeNumber(fout, binaryType_float, material.Edge.Z, 0.0, true)
		if err != nil {
			return err
		}
		err = rep.writeNumber(fout, binaryType_float, material.Edge.W, 0.0, true)
		if err != nil {
			return err
		}
		err = rep.writeNumber(fout, binaryType_float, material.EdgeSize, 0.0, true)
		if err != nil {
			return err
		}
		err = rep.writeNumber(fout, textureIdxType, float64(material.TextureIndex), 0.0, false)
		if err != nil {
			return err
		}
		err = rep.writeNumber(fout, textureIdxType, float64(material.SphereTextureIndex), 0.0, false)
		if err != nil {
			return err
		}
		err = rep.writeByte(fout, int(material.SphereMode), true)
		if err != nil {
			return err
		}
		err = rep.writeByte(fout, int(material.ToonSharingFlag), true)
		if err != nil {
			return err
		}
		if material.ToonSharingFlag == pmx.TOON_SHARING_SHARING {
			err = rep.writeNumber(fout, textureIdxType, float64(material.ToonTextureIndex), 0.0, false)
		} else {
			err = rep.writeByte(fout, int(material.ToonTextureIndex), true)
		}
		if err != nil {
			return err
		}
		err = rep.writeText(fout, material.Memo, "")
		if err != nil {
			return err
		}
		err = rep.writeNumber(fout, binaryType_int, float64(material.VerticesCount), 0.0, true)
		if err != nil {
			return err
		}
	}

	return nil
}

// saveBones ボーンデータの書き込み
func (rep *PmxRepository) saveBones(fout *os.File, targetBones []*pmx.Bone, boneIdxType binaryType) error {
	defer mlog.I(mi18n.T("保存途中完了", map[string]interface{}{"Type": mi18n.T("ボーン")}))

	err := rep.writeNumber(fout, binaryType_int, float64(len(targetBones)), 0.0, true)
	if err != nil {
		return err
	}

	for i, bone := range targetBones {
		err = rep.writeText(fout, bone.Name(), fmt.Sprintf("Bone %d", i))
		if err != nil {
			return err
		}
		err = rep.writeText(fout, bone.EnglishName(), fmt.Sprintf("Bone %d", i))
		if err != nil {
			return err
		}
		err = rep.writeNumber(fout, binaryType_float, bone.Position.X, 0.0, false)
		if err != nil {
			return err
		}
		err = rep.writeNumber(fout, binaryType_float, bone.Position.Y, 0.0, false)
		if err != nil {
			return err
		}
		err = rep.writeNumber(fout, binaryType_float, bone.Position.Z, 0.0, false)
		if err != nil {
			return err
		}
		err = rep.writeNumber(fout, boneIdxType, float64(bone.ParentIndex), 0.0, false)
		if err != nil {
			return err
		}
		err = rep.writeNumber(fout, binaryType_int, float64(bone.Layer), 0.0, true)
		if err != nil {
			return err
		}
		err = rep.writeShort(fout, uint16(bone.BoneFlag))
		if err != nil {
			return err
		}

		if bone.IsTailBone() {
			err = rep.writeNumber(fout, boneIdxType, float64(bone.TailIndex), 0.0, false)
			if err != nil {
				return err
			}
		} else {
			err = rep.writeNumber(fout, binaryType_float, bone.TailPosition.X, 0.0, false)
			if err != nil {
				return err
			}
			err = rep.writeNumber(fout, binaryType_float, bone.TailPosition.Y, 0.0, false)
			if err != nil {
				return err
			}
			err = rep.writeNumber(fout, binaryType_float, bone.TailPosition.Z, 0.0, false)
			if err != nil {
				return err
			}
		}
		if bone.IsEffectorTranslation() || bone.IsEffectorRotation() {
			err = rep.writeNumber(fout, boneIdxType, float64(bone.EffectIndex), 0.0, false)
			if err != nil {
				return err
			}
			err = rep.writeNumber(fout, binaryType_float, bone.EffectFactor, 0.0, false)
			if err != nil {
				return err
			}
		}
		if bone.HasFixedAxis() {
			err = rep.writeNumber(fout, binaryType_float, bone.FixedAxis.X, 0.0, false)
			if err != nil {
				return err
			}
			err = rep.writeNumber(fout, binaryType_float, bone.FixedAxis.Y, 0.0, false)
			if err != nil {
				return err
			}
			err = rep.writeNumber(fout, binaryType_float, bone.FixedAxis.Z, 0.0, false)
			if err != nil {
				return err
			}
		}
		if bone.HasLocalAxis() {
			err = rep.writeNumber(fout, binaryType_float, bone.LocalAxisX.X, 0.0, false)
			if err != nil {
				return err
			}
			err = rep.writeNumber(fout, binaryType_float, bone.LocalAxisX.Y, 0.0, false)
			if err != nil {
				return err
			}
			err = rep.writeNumber(fout, binaryType_float, bone.LocalAxisX.Z, 0.0, false)
			if err != nil {
				return err
			}
			err = rep.writeNumber(fout, binaryType_float, bone.LocalAxisZ.X, 0.0, false)
			if err != nil {
				return err
			}
			err = rep.writeNumber(fout, binaryType_float, bone.LocalAxisZ.Y, 0.0, false)
			if err != nil {
				return err
			}
			err = rep.writeNumber(fout, binaryType_float, bone.LocalAxisZ.Z, 0.0, false)
			if err != nil {
				return err
			}
		}
		if bone.IsEffectorParentDeform() {
			err = rep.writeNumber(fout, binaryType_int, float64(bone.EffectorKey), 0.0, true)
			if err != nil {
				return err
			}
		}
		if bone.IsIK() {
			err = rep.writeNumber(fout, boneIdxType, float64(bone.Ik.BoneIndex), 0.0, false)
			if err != nil {
				return err
			}
			err = rep.writeNumber(fout, binaryType_int, float64(bone.Ik.LoopCount), 0.0, true)
			if err != nil {
				return err
			}
			err = rep.writeNumber(fout, binaryType_float, bone.Ik.UnitRotation.Radians().X, 0.0, false)
			if err != nil {
				return err
			}
			err = rep.writeNumber(fout, binaryType_int, float64(len(bone.Ik.Links)), 0.0, true)
			if err != nil {
				return err
			}

			for _, link := range bone.Ik.Links {
				err = rep.writeNumber(fout, boneIdxType, float64(link.BoneIndex), 0.0, false)
				if err != nil {
					return err
				}
				err = rep.writeByte(fout, int(mmath.BoolToInt(link.AngleLimit)), true)
				if err != nil {
					return err
				}
				if link.AngleLimit {
					err = rep.writeNumber(fout, binaryType_float, link.MinAngleLimit.Radians().X, 0.0, false)
					if err != nil {
						return err
					}
					err = rep.writeNumber(fout, binaryType_float, link.MinAngleLimit.Radians().Y, 0.0, false)
					if err != nil {
						return err
					}
					err = rep.writeNumber(fout, binaryType_float, link.MinAngleLimit.Radians().Z, 0.0, false)
					if err != nil {
						return err
					}
					err = rep.writeNumber(fout, binaryType_float, link.MaxAngleLimit.Radians().X, 0.0, false)
					if err != nil {
						return err
					}
					err = rep.writeNumber(fout, binaryType_float, link.MaxAngleLimit.Radians().Y, 0.0, false)
					if err != nil {
						return err
					}
					err = rep.writeNumber(fout, binaryType_float, link.MaxAngleLimit.Radians().Z, 0.0, false)
					if err != nil {
						return err
					}
				}
			}
		}
	}

	return nil
}

// saveMorphs モーフデータの書き込み
func (rep *PmxRepository) saveMorphs(
	fout *os.File, targetMorphs []*pmx.Morph, vertexIdxType, boneIdxType, materialIdxType, morphIdxType binaryType,
) error {
	defer mlog.I(mi18n.T("保存途中完了", map[string]interface{}{"Type": mi18n.T("モーフ")}))

	err := rep.writeNumber(fout, binaryType_int, float64(len(targetMorphs)), 0.0, true)
	if err != nil {
		return err
	}

	for i, morph := range targetMorphs {
		err = rep.writeText(fout, morph.Name(), fmt.Sprintf("Morph %d", i))
		if err != nil {
			return err
		}
		err = rep.writeText(fout, morph.EnglishName(), fmt.Sprintf("Morph %d", i))
		if err != nil {
			return err
		}
		err = rep.writeByte(fout, int(morph.Panel), true)
		if err != nil {
			return err
		}
		err = rep.writeByte(fout, int(morph.MorphType), true)
		if err != nil {
			return err
		}
		err = rep.writeNumber(fout, binaryType_int, float64(len(morph.Offsets)), 0.0, true)
		if err != nil {
			return err
		}

		for _, offset := range morph.Offsets {
			switch off := offset.(type) {
			case *pmx.VertexMorphOffset:
				err = rep.writeNumber(fout, vertexIdxType, float64(off.VertexIndex), 0.0, true)
				if err != nil {
					return err
				}
				err = rep.writeNumber(fout, binaryType_float, off.Position.X, 0.0, false)
				if err != nil {
					return err
				}
				err = rep.writeNumber(fout, binaryType_float, off.Position.Y, 0.0, false)
				if err != nil {
					return err
				}
				err = rep.writeNumber(fout, binaryType_float, off.Position.Z, 0.0, false)
				if err != nil {
					return err
				}
			case *pmx.UvMorphOffset:
				err = rep.writeNumber(fout, vertexIdxType, float64(off.VertexIndex), 0.0, true)
				if err != nil {
					return err
				}
				err = rep.writeNumber(fout, binaryType_float, off.Uv.X, 0.0, false)
				if err != nil {
					return err
				}
				err = rep.writeNumber(fout, binaryType_float, off.Uv.Y, 0.0, false)
				if err != nil {
					return err
				}
				err = rep.writeNumber(fout, binaryType_float, off.Uv.Z, 0.0, false)
				if err != nil {
					return err
				}
				err = rep.writeNumber(fout, binaryType_float, off.Uv.Z, 0.0, false)
				if err != nil {
					return err
				}
			case *pmx.BoneMorphOffset:
				err = rep.writeNumber(fout, boneIdxType, float64(off.BoneIndex), 0.0, false)
				if err != nil {
					return err
				}
				err = rep.writeNumber(fout, binaryType_float, off.Position.X, 0.0, false)
				if err != nil {
					return err
				}
				err = rep.writeNumber(fout, binaryType_float, off.Position.Y, 0.0, false)
				if err != nil {
					return err
				}
				err = rep.writeNumber(fout, binaryType_float, off.Position.Z, 0.0, false)
				if err != nil {
					return err
				}
				err = rep.writeNumber(fout, binaryType_float, off.Rotation.Vec3().X, 0.0, false)
				if err != nil {
					return err
				}
				err = rep.writeNumber(fout, binaryType_float, off.Rotation.Vec3().Y, 0.0, false)
				if err != nil {
					return err
				}
				err = rep.writeNumber(fout, binaryType_float, off.Rotation.Vec3().Z, 0.0, false)
				if err != nil {
					return err
				}
				err = rep.writeNumber(fout, binaryType_float, off.Rotation.W, 0.0, false)
				if err != nil {
					return err
				}
			case *pmx.MaterialMorphOffset:
				err = rep.writeNumber(fout, materialIdxType, float64(off.MaterialIndex), 0.0, false)
				if err != nil {
					return err
				}
				err = rep.writeByte(fout, int(off.CalcMode), true)
				if err != nil {
					return err
				}
				err = rep.writeNumber(fout, binaryType_float, off.Diffuse.X, 0.0, false)
				if err != nil {
					return err
				}
				err = rep.writeNumber(fout, binaryType_float, off.Diffuse.Y, 0.0, false)
				if err != nil {
					return err
				}
				err = rep.writeNumber(fout, binaryType_float, off.Diffuse.Z, 0.0, false)
				if err != nil {
					return err
				}
				err = rep.writeNumber(fout, binaryType_float, off.Diffuse.W, 0.0, false)
				if err != nil {
					return err
				}
				err = rep.writeNumber(fout, binaryType_float, off.Specular.X, 0.0, false)
				if err != nil {
					return err
				}
				err = rep.writeNumber(fout, binaryType_float, off.Specular.Y, 0.0, false)
				if err != nil {
					return err
				}
				err = rep.writeNumber(fout, binaryType_float, off.Specular.Z, 0.0, false)
				if err != nil {
					return err
				}
				err = rep.writeNumber(fout, binaryType_float, off.Specular.W, 0.0, false)
				if err != nil {
					return err
				}
				err = rep.writeNumber(fout, binaryType_float, off.Ambient.X, 0.0, false)
				if err != nil {
					return err
				}
				err = rep.writeNumber(fout, binaryType_float, off.Ambient.Y, 0.0, false)
				if err != nil {
					return err
				}
				err = rep.writeNumber(fout, binaryType_float, off.Ambient.Z, 0.0, false)
				if err != nil {
					return err
				}
				err = rep.writeNumber(fout, binaryType_float, off.Edge.X, 0.0, false)
				if err != nil {
					return err
				}
				err = rep.writeNumber(fout, binaryType_float, off.Edge.Y, 0.0, false)
				if err != nil {
					return err
				}
				err = rep.writeNumber(fout, binaryType_float, off.Edge.Z, 0.0, false)
				if err != nil {
					return err
				}
				err = rep.writeNumber(fout, binaryType_float, off.Edge.W, 0.0, false)
				if err != nil {
					return err
				}
				err = rep.writeNumber(fout, binaryType_float, off.EdgeSize, 0.0, false)
				if err != nil {
					return err
				}
				err = rep.writeNumber(fout, binaryType_float, off.TextureFactor.X, 0.0, false)
				if err != nil {
					return err
				}
				err = rep.writeNumber(fout, binaryType_float, off.TextureFactor.Y, 0.0, false)
				if err != nil {
					return err
				}
				err = rep.writeNumber(fout, binaryType_float, off.TextureFactor.Z, 0.0, false)
				if err != nil {
					return err
				}
				err = rep.writeNumber(fout, binaryType_float, off.TextureFactor.W, 0.0, false)
				if err != nil {
					return err
				}
				err = rep.writeNumber(fout, binaryType_float, off.SphereTextureFactor.X, 0.0, false)
				if err != nil {
					return err
				}
				err = rep.writeNumber(fout, binaryType_float, off.SphereTextureFactor.Y, 0.0, false)
				if err != nil {
					return err
				}
				err = rep.writeNumber(fout, binaryType_float, off.SphereTextureFactor.Z, 0.0, false)
				if err != nil {
					return err
				}
				err = rep.writeNumber(fout, binaryType_float, off.SphereTextureFactor.W, 0.0, false)
				if err != nil {
					return err
				}
				err = rep.writeNumber(fout, binaryType_float, off.ToonTextureFactor.X, 0.0, false)
				if err != nil {
					return err
				}
				err = rep.writeNumber(fout, binaryType_float, off.ToonTextureFactor.Y, 0.0, false)
				if err != nil {
					return err
				}
				err = rep.writeNumber(fout, binaryType_float, off.ToonTextureFactor.Z, 0.0, false)
				if err != nil {
					return err
				}
				err = rep.writeNumber(fout, binaryType_float, off.ToonTextureFactor.W, 0.0, false)
				if err != nil {
					return err
				}
			case *pmx.GroupMorphOffset:
				err = rep.writeNumber(fout, morphIdxType, float64(off.MorphIndex), 0.0, false)
				if err != nil {
					return err
				}
				err = rep.writeNumber(fout, binaryType_float, off.MorphFactor, 0.0, false)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

// saveDisplaySlots 表示枠データの書き込み
func (rep *PmxRepository) saveDisplaySlots(fout *os.File, model *pmx.PmxModel, boneIdxType, morphIdxType binaryType) error {
	defer mlog.I(mi18n.T("保存途中完了", map[string]interface{}{"Type": mi18n.T("表示枠")}))

	err := rep.writeNumber(fout, binaryType_int, float64(model.DisplaySlots.Len()), 0.0, true)
	if err != nil {
		return err
	}

	for i := range model.DisplaySlots.Len() {
		displaySlot := model.DisplaySlots.Get(i)

		err = rep.writeText(fout, displaySlot.Name(), fmt.Sprintf("Display %d", i))
		if err != nil {
			return err
		}
		err = rep.writeText(fout, displaySlot.EnglishName(), fmt.Sprintf("Display %d", i))
		if err != nil {
			return err
		}
		err = rep.writeByte(fout, int(displaySlot.SpecialFlag), true)
		if err != nil {
			return err
		}
		err = rep.writeNumber(fout, binaryType_int, float64(len(displaySlot.References)), 0.0, true)
		if err != nil {
			return err
		}

		for _, reference := range displaySlot.References {
			err = rep.writeByte(fout, int(reference.DisplayType), true)
			if err != nil {
				return err
			}
			if reference.DisplayType == pmx.DISPLAY_TYPE_BONE {
				err = rep.writeNumber(fout, boneIdxType, float64(reference.DisplayIndex), 0.0, false)
			} else {
				err = rep.writeNumber(fout, morphIdxType, float64(reference.DisplayIndex), 0.0, false)
			}
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// saveRigidBodies 剛体データの書き込み
func (rep *PmxRepository) saveRigidBodies(fout *os.File, model *pmx.PmxModel, boneIdxType binaryType) error {
	defer mlog.I(mi18n.T("保存途中完了", map[string]interface{}{"Type": mi18n.T("剛体")}))

	err := rep.writeNumber(fout, binaryType_int, float64(model.RigidBodies.Len()), 0.0, true)
	if err != nil {
		return err
	}

	for i := range model.RigidBodies.Len() {
		rigidbody := model.RigidBodies.Get(i)
		err = rep.writeText(fout, rigidbody.Name(), fmt.Sprintf("Rigidbody %d", i))
		if err != nil {
			return err
		}
		err = rep.writeText(fout, rigidbody.EnglishName(), fmt.Sprintf("Rigidbody %d", i))
		if err != nil {
			return err
		}
		err = rep.writeNumber(fout, boneIdxType, float64(rigidbody.BoneIndex), 0.0, false)
		if err != nil {
			return err
		}
		err = rep.writeByte(fout, int(rigidbody.CollisionGroup), true)
		if err != nil {
			return err
		}
		err = rep.writeShort(fout, uint16(rigidbody.CollisionGroupMaskValue))
		if err != nil {
			return err
		}
		err = rep.writeByte(fout, int(rigidbody.ShapeType), true)
		if err != nil {
			return err
		}
		err = rep.writeNumber(fout, binaryType_float, rigidbody.Size.X, 0.0, true)
		if err != nil {
			return err
		}
		err = rep.writeNumber(fout, binaryType_float, rigidbody.Size.Y, 0.0, true)
		if err != nil {
			return err
		}
		err = rep.writeNumber(fout, binaryType_float, rigidbody.Size.Z, 0.0, true)
		if err != nil {
			return err
		}
		err = rep.writeNumber(fout, binaryType_float, rigidbody.Position.X, 0.0, false)
		if err != nil {
			return err
		}
		err = rep.writeNumber(fout, binaryType_float, rigidbody.Position.Y, 0.0, false)
		if err != nil {
			return err
		}
		err = rep.writeNumber(fout, binaryType_float, rigidbody.Position.Z, 0.0, false)
		if err != nil {
			return err
		}
		err = rep.writeNumber(fout, binaryType_float, rigidbody.Rotation.Radians().X, 0.0, false)
		if err != nil {
			return err
		}
		err = rep.writeNumber(fout, binaryType_float, rigidbody.Rotation.Radians().Y, 0.0, false)
		if err != nil {
			return err
		}
		err = rep.writeNumber(fout, binaryType_float, rigidbody.Rotation.Radians().Z, 0.0, false)
		if err != nil {
			return err
		}
		err = rep.writeNumber(fout, binaryType_float, rigidbody.RigidBodyParam.Mass, 0.0, true)
		if err != nil {
			return err
		}
		err = rep.writeNumber(fout, binaryType_float, rigidbody.RigidBodyParam.LinearDamping, 0.0, true)
		if err != nil {
			return err
		}
		err = rep.writeNumber(fout, binaryType_float, rigidbody.RigidBodyParam.AngularDamping, 0.0, true)
		if err != nil {
			return err
		}
		err = rep.writeNumber(fout, binaryType_float, rigidbody.RigidBodyParam.Restitution, 0.0, true)
		if err != nil {
			return err
		}
		err = rep.writeNumber(fout, binaryType_float, rigidbody.RigidBodyParam.Friction, 0.0, true)
		if err != nil {
			return err
		}
		err = rep.writeByte(fout, int(rigidbody.PhysicsType), true)
		if err != nil {
			return err
		}
	}

	return nil
}

// saveJoints ジョイントデータの書き込み
func (rep *PmxRepository) saveJoints(fout *os.File, model *pmx.PmxModel, rigidbodyIdxType binaryType) error {
	defer mlog.I(mi18n.T("保存途中完了", map[string]interface{}{"Type": mi18n.T("ジョイント")}))

	err := rep.writeNumber(fout, binaryType_int, float64(model.Joints.Len()), 0.0, true)
	if err != nil {
		return err
	}

	for i := range model.Joints.Len() {
		joint := model.Joints.Get(i)

		err = rep.writeText(fout, joint.Name(), fmt.Sprintf("Joint %d", i))
		if err != nil {
			return err
		}
		err = rep.writeText(fout, joint.EnglishName(), fmt.Sprintf("Joint %d", i))
		if err != nil {
			return err
		}
		err = rep.writeByte(fout, int(joint.JointType), true)
		if err != nil {
			return err
		}
		err = rep.writeNumber(fout, rigidbodyIdxType, float64(joint.RigidbodyIndexA), 0.0, false)
		if err != nil {
			return err
		}
		err = rep.writeNumber(fout, rigidbodyIdxType, float64(joint.RigidbodyIndexB), 0.0, false)
		if err != nil {
			return err
		}
		err = rep.writeNumber(fout, binaryType_float, joint.Position.X, 0.0, false)
		if err != nil {
			return err
		}
		err = rep.writeNumber(fout, binaryType_float, joint.Position.Y, 0.0, false)
		if err != nil {
			return err
		}
		err = rep.writeNumber(fout, binaryType_float, joint.Position.Z, 0.0, false)
		if err != nil {
			return err
		}
		err = rep.writeNumber(fout, binaryType_float, joint.Rotation.Radians().X, 0.0, false)
		if err != nil {
			return err
		}
		err = rep.writeNumber(fout, binaryType_float, joint.Rotation.Radians().Y, 0.0, false)
		if err != nil {
			return err
		}
		err = rep.writeNumber(fout, binaryType_float, joint.Rotation.Radians().Z, 0.0, false)
		if err != nil {
			return err
		}
		err = rep.writeNumber(fout, binaryType_float, joint.JointParam.TranslationLimitMin.X, 0.0, false)
		if err != nil {
			return err
		}
		err = rep.writeNumber(fout, binaryType_float, joint.JointParam.TranslationLimitMin.Y, 0.0, false)
		if err != nil {
			return err
		}
		err = rep.writeNumber(fout, binaryType_float, joint.JointParam.TranslationLimitMin.Z, 0.0, false)
		if err != nil {
			return err
		}
		err = rep.writeNumber(fout, binaryType_float, joint.JointParam.TranslationLimitMax.X, 0.0, false)
		if err != nil {
			return err
		}
		err = rep.writeNumber(fout, binaryType_float, joint.JointParam.TranslationLimitMax.Y, 0.0, false)
		if err != nil {
			return err
		}
		err = rep.writeNumber(fout, binaryType_float, joint.JointParam.TranslationLimitMax.Z, 0.0, false)
		if err != nil {
			return err
		}
		err = rep.writeNumber(fout, binaryType_float, joint.JointParam.RotationLimitMin.Radians().X, 0.0, false)
		if err != nil {
			return err
		}
		err = rep.writeNumber(fout, binaryType_float, joint.JointParam.RotationLimitMin.Radians().Y, 0.0, false)
		if err != nil {
			return err
		}
		err = rep.writeNumber(fout, binaryType_float, joint.JointParam.RotationLimitMin.Radians().Z, 0.0, false)
		if err != nil {
			return err
		}
		err = rep.writeNumber(fout, binaryType_float, joint.JointParam.RotationLimitMax.Radians().X, 0.0, false)
		if err != nil {
			return err
		}
		err = rep.writeNumber(fout, binaryType_float, joint.JointParam.RotationLimitMax.Radians().Y, 0.0, false)
		if err != nil {
			return err
		}
		err = rep.writeNumber(fout, binaryType_float, joint.JointParam.RotationLimitMax.Radians().Z, 0.0, false)
		if err != nil {
			return err
		}
		err = rep.writeNumber(fout, binaryType_float, joint.JointParam.SpringConstantTranslation.X, 0.0, false)
		if err != nil {
			return err
		}
		err = rep.writeNumber(fout, binaryType_float, joint.JointParam.SpringConstantTranslation.Y, 0.0, false)
		if err != nil {
			return err
		}
		err = rep.writeNumber(fout, binaryType_float, joint.JointParam.SpringConstantTranslation.Z, 0.0, false)
		if err != nil {
			return err
		}
		err = rep.writeNumber(fout, binaryType_float, joint.JointParam.SpringConstantRotation.X, 0.0, false)
		if err != nil {
			return err
		}
		err = rep.writeNumber(fout, binaryType_float, joint.JointParam.SpringConstantRotation.Y, 0.0, false)
		if err != nil {
			return err
		}
		err = rep.writeNumber(fout, binaryType_float, joint.JointParam.SpringConstantRotation.Z, 0.0, false)
		if err != nil {
			return err
		}
	}

	return nil
}

// ------------------------------
func (rep *PmxRepository) writeText(fout *os.File, text string, defaultText string) error {
	var binaryTxt []byte
	var err error

	// エンコードの試行
	binaryTxt, err = rep.encodeUTF16LE(text)
	if err != nil {
		binaryTxt, _ = rep.encodeUTF16LE(defaultText)
	}

	// バイナリサイズの書き込み
	var buf bytes.Buffer
	err = binary.Write(&buf, binary.LittleEndian, int32(len(binaryTxt)))
	if err != nil {
		return err
	}
	_, err = fout.Write(buf.Bytes())
	if err != nil {
		return err
	}

	// 文字列の書き込み
	_, err = fout.Write(binaryTxt)
	return err
}

func (rep *PmxRepository) encodeUTF16LE(s string) ([]byte, error) {
	runes := []rune(s)
	encoded := utf16.Encode(runes)
	buf := new(bytes.Buffer)
	for _, r := range encoded {
		err := binary.Write(buf, binary.LittleEndian, r)
		if err != nil {
			return nil, err
		}
	}
	return buf.Bytes(), nil
}

func (rep *PmxRepository) defineWriteIndexForVertex(size int) (int, binaryType) {
	if size < 256 {
		return 1, binaryType_unsignedByte
	} else if size <= 65535 {
		return 2, binaryType_unsignedShort
	}
	return 4, binaryType_int
}

func (rep *PmxRepository) defineWriteIndexForOthers(size int) (int, binaryType) {
	if size < 128 {
		return 1, binaryType_byte
	} else if size <= 32767 {
		return 2, binaryType_short
	}
	return 4, binaryType_int
}
