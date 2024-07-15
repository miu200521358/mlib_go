package repository

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
	"unicode/utf16"

	"github.com/miu200521358/mlib_go/pkg/domain/core"
	"github.com/miu200521358/mlib_go/pkg/domain/mmath"
	"github.com/miu200521358/mlib_go/pkg/domain/pmx"
	"github.com/miu200521358/mlib_go/pkg/mutils"
	"github.com/miu200521358/mlib_go/pkg/mutils/mlog"
)

func (r *PmxRepository) Save(overridePath string, data core.IHashModel, includeSystem bool) error {
	model := data.(*pmx.PmxModel)

	path := model.GetPath()
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

	err = r.writeNumber(fout, binaryType_float, 2.0, 0.0, true)
	if err != nil {
		return err
	}

	err = r.writeByte(fout, 8, true)
	if err != nil {
		return err
	}

	err = r.writeByte(fout, 0, true)
	if err != nil {
		return err
	}

	err = r.writeByte(fout, model.ExtendedUVCount, true)
	if err != nil {
		return err
	}

	vertexIdxSize, vertexIdxType := r.defineWriteIndexForVertex(model.Vertices.Len())
	err = r.writeByte(fout, vertexIdxSize, true)
	if err != nil {
		return err
	}

	textureIdxSize, textureIdxType := r.defineWriteIndexForOthers(model.Textures.Len())
	err = r.writeByte(fout, textureIdxSize, true)
	if err != nil {
		return err
	}

	materialIdxSize, materialIdxType := r.defineWriteIndexForOthers(model.Materials.Len())
	err = r.writeByte(fout, materialIdxSize, true)
	if err != nil {
		return err
	}

	boneIdxSize, boneIdxType := r.defineWriteIndexForOthers(len(filteredBones))
	err = r.writeByte(fout, boneIdxSize, true)
	if err != nil {
		return err
	}

	morphIdxSize, morphIdxType := r.defineWriteIndexForOthers(len(filteredMorphs))
	err = r.writeByte(fout, morphIdxSize, true)
	if err != nil {
		return err
	}

	rigidbodyIdxSize, rigidbodyIdxType := r.defineWriteIndexForOthers(model.RigidBodies.Len())
	err = r.writeByte(fout, rigidbodyIdxSize, true)
	if err != nil {
		return err
	}

	err = r.writeText(fout, model.Name, "Pmx Model")
	if err != nil {
		return err
	}

	err = r.writeText(fout, model.EnglishName, "Pmx Model")
	if err != nil {
		return err
	}

	err = r.writeText(fout, model.Comment, "")
	if err != nil {
		return err
	}

	err = r.writeText(fout, model.EnglishComment, "")
	if err != nil {
		return err
	}

	err = r.saveVertices(fout, model, boneIdxType)
	if err != nil {
		return err
	}

	err = r.saveFaces(fout, model, vertexIdxType)
	if err != nil {
		return err
	}

	err = r.saveTextures(fout, model)
	if err != nil {
		return err
	}

	err = r.saveMaterials(fout, model, textureIdxType)
	if err != nil {
		return err
	}

	err = r.saveBones(fout, filteredBones, boneIdxType)
	if err != nil {
		return err
	}

	err = r.saveMorphs(fout, filteredMorphs, vertexIdxType, boneIdxType, materialIdxType, morphIdxType)
	if err != nil {
		return err
	}

	err = r.saveDisplaySlots(fout, model, boneIdxType, morphIdxType)
	if err != nil {
		return err
	}

	err = r.saveRigidBodies(fout, model, boneIdxType)
	if err != nil {
		return err
	}

	err = r.saveJoints(fout, model, rigidbodyIdxType)
	if err != nil {
		return err
	}

	return nil
}

// saveVertices 頂点データの書き込み
func (r *PmxRepository) saveVertices(fout *os.File, model *pmx.PmxModel, boneIdxType binaryType) error {
	err := r.writeNumber(fout, binaryType_int, float64(model.Vertices.Len()), 0.0, true)
	if err != nil {
		return err
	}

	for i := range model.Vertices.Len() {
		vertex := model.Vertices.Get(i)

		r.writeNumber(fout, binaryType_float, vertex.Position.GetX(), 0.0, false)
		r.writeNumber(fout, binaryType_float, vertex.Position.GetY(), 0.0, false)
		r.writeNumber(fout, binaryType_float, vertex.Position.GetZ(), 0.0, false)
		r.writeNumber(fout, binaryType_float, vertex.Normal.GetX(), 0.0, false)
		r.writeNumber(fout, binaryType_float, vertex.Normal.GetY(), 0.0, false)
		r.writeNumber(fout, binaryType_float, vertex.Normal.GetZ(), 0.0, false)
		r.writeNumber(fout, binaryType_float, vertex.Uv.X, 0.0, false)
		r.writeNumber(fout, binaryType_float, vertex.Uv.Y, 0.0, false)

		for _, uv := range vertex.ExtendedUvs {
			r.writeNumber(fout, binaryType_float, uv.GetX(), 0.0, false)
			r.writeNumber(fout, binaryType_float, uv.GetY(), 0.0, false)
			r.writeNumber(fout, binaryType_float, uv.GetZ(), 0.0, false)
			r.writeNumber(fout, binaryType_float, uv.GetW(), 0.0, false)
		}

		for j := len(vertex.ExtendedUvs); j < model.ExtendedUVCount; j++ {
			r.writeNumber(fout, binaryType_float, 0.0, 0.0, false)
			r.writeNumber(fout, binaryType_float, 0.0, 0.0, false)
			r.writeNumber(fout, binaryType_float, 0.0, 0.0, false)
			r.writeNumber(fout, binaryType_float, 0.0, 0.0, false)
		}

		switch v := vertex.Deform.(type) {
		case *pmx.Bdef1:
			r.writeByte(fout, 0, false)
			r.writeNumber(fout, boneIdxType, float64(v.Indexes[0]), 0.0, false)
		case *pmx.Bdef2:
			r.writeByte(fout, 1, false)
			r.writeNumber(fout, boneIdxType, float64(v.Indexes[0]), 0.0, false)
			r.writeNumber(fout, boneIdxType, float64(v.Indexes[1]), 0.0, false)
			r.writeNumber(fout, binaryType_float, v.Weights[0], 0.0, true)
		case *pmx.Bdef4:
			r.writeByte(fout, 2, false)
			r.writeNumber(fout, boneIdxType, float64(v.Indexes[0]), 0.0, false)
			r.writeNumber(fout, boneIdxType, float64(v.Indexes[1]), 0.0, false)
			r.writeNumber(fout, boneIdxType, float64(v.Indexes[2]), 0.0, false)
			r.writeNumber(fout, boneIdxType, float64(v.Indexes[3]), 0.0, false)
			r.writeNumber(fout, binaryType_float, v.Weights[0], 0.0, true)
			r.writeNumber(fout, binaryType_float, v.Weights[1], 0.0, true)
			r.writeNumber(fout, binaryType_float, v.Weights[2], 0.0, true)
			r.writeNumber(fout, binaryType_float, v.Weights[3], 0.0, true)
		case *pmx.Sdef:
			r.writeByte(fout, 3, false)
			r.writeNumber(fout, boneIdxType, float64(v.Indexes[0]), 0.0, false)
			r.writeNumber(fout, boneIdxType, float64(v.Indexes[1]), 0.0, false)
			r.writeNumber(fout, binaryType_float, v.Weights[0], 0.0, true)
			r.writeNumber(fout, binaryType_float, v.SdefC.GetX(), 0.0, false)
			r.writeNumber(fout, binaryType_float, v.SdefC.GetY(), 0.0, false)
			r.writeNumber(fout, binaryType_float, v.SdefC.GetZ(), 0.0, false)
			r.writeNumber(fout, binaryType_float, v.SdefR0.GetX(), 0.0, false)
			r.writeNumber(fout, binaryType_float, v.SdefR0.GetY(), 0.0, false)
			r.writeNumber(fout, binaryType_float, v.SdefR0.GetZ(), 0.0, false)
			r.writeNumber(fout, binaryType_float, v.SdefR1.GetX(), 0.0, false)
			r.writeNumber(fout, binaryType_float, v.SdefR1.GetY(), 0.0, false)
			r.writeNumber(fout, binaryType_float, v.SdefR1.GetZ(), 0.0, false)
		default:
			mlog.W("頂点deformなし: %v\n", vertex)
		}

		r.writeNumber(fout, binaryType_float, vertex.EdgeFactor, 0.0, true)
	}

	return nil
}

// saveFaces 面データの書き込み
func (r *PmxRepository) saveFaces(fout *os.File, model *pmx.PmxModel, vertexIdxType binaryType) error {
	err := r.writeNumber(fout, binaryType_int, float64(model.Faces.Len()*3), 0.0, true)
	if err != nil {
		return err
	}

	for i := range model.Faces.Len() {
		face := model.Faces.Get(i)

		for _, vidx := range face.VertexIndexes {
			err = r.writeNumber(fout, vertexIdxType, float64(vidx), 0.0, true)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// saveTextures テクスチャデータの書き込み
func (r *PmxRepository) saveTextures(fout *os.File, model *pmx.PmxModel) error {
	err := r.writeNumber(fout, binaryType_int, float64(model.Textures.Len()), 0.0, true)
	if err != nil {
		return err
	}

	for i := range model.Textures.Len() {
		texture := model.Textures.Get(i)
		err = r.writeText(fout, texture.Name, "")
		if err != nil {
			return err
		}
	}

	return nil
}

// saveMaterials 材質データの書き込み
func (r *PmxRepository) saveMaterials(fout *os.File, model *pmx.PmxModel, textureIdxType binaryType) error {
	err := r.writeNumber(fout, binaryType_int, float64(model.Materials.Len()), 0.0, true)
	if err != nil {
		return err
	}

	for i := range model.Materials.Len() {
		material := model.Materials.Get(i)
		err = r.writeText(fout, material.Name, fmt.Sprintf("Material %d", i))
		if err != nil {
			return err
		}
		err = r.writeText(fout, material.EnglishName, fmt.Sprintf("Material %d", i))
		if err != nil {
			return err
		}
		err = r.writeNumber(fout, binaryType_float, material.Diffuse.GetX(), 0.0, true)
		if err != nil {
			return err
		}
		err = r.writeNumber(fout, binaryType_float, material.Diffuse.GetY(), 0.0, true)
		if err != nil {
			return err
		}
		err = r.writeNumber(fout, binaryType_float, material.Diffuse.GetZ(), 0.0, true)
		if err != nil {
			return err
		}
		err = r.writeNumber(fout, binaryType_float, material.Diffuse.GetW(), 0.0, true)
		if err != nil {
			return err
		}
		err = r.writeNumber(fout, binaryType_float, material.Specular.GetX(), 0.0, true)
		if err != nil {
			return err
		}
		err = r.writeNumber(fout, binaryType_float, material.Specular.GetY(), 0.0, true)
		if err != nil {
			return err
		}
		err = r.writeNumber(fout, binaryType_float, material.Specular.GetZ(), 0.0, true)
		if err != nil {
			return err
		}
		err = r.writeNumber(fout, binaryType_float, material.Specular.GetW(), 0.0, true)
		if err != nil {
			return err
		}
		err = r.writeNumber(fout, binaryType_float, material.Ambient.GetX(), 0.0, true)
		if err != nil {
			return err
		}
		err = r.writeNumber(fout, binaryType_float, material.Ambient.GetY(), 0.0, true)
		if err != nil {
			return err
		}
		err = r.writeNumber(fout, binaryType_float, material.Ambient.GetZ(), 0.0, true)
		if err != nil {
			return err
		}
		err = r.writeByte(fout, int(material.DrawFlag), true)
		if err != nil {
			return err
		}
		err = r.writeNumber(fout, binaryType_float, material.Edge.GetX(), 0.0, true)
		if err != nil {
			return err
		}
		err = r.writeNumber(fout, binaryType_float, material.Edge.GetY(), 0.0, true)
		if err != nil {
			return err
		}
		err = r.writeNumber(fout, binaryType_float, material.Edge.GetZ(), 0.0, true)
		if err != nil {
			return err
		}
		err = r.writeNumber(fout, binaryType_float, material.Edge.GetW(), 0.0, true)
		if err != nil {
			return err
		}
		err = r.writeNumber(fout, binaryType_float, material.EdgeSize, 0.0, true)
		if err != nil {
			return err
		}
		err = r.writeNumber(fout, textureIdxType, float64(material.TextureIndex), 0.0, true)
		if err != nil {
			return err
		}
		err = r.writeNumber(fout, textureIdxType, float64(material.SphereTextureIndex), 0.0, true)
		if err != nil {
			return err
		}
		err = r.writeByte(fout, int(material.SphereMode), true)
		if err != nil {
			return err
		}
		err = r.writeByte(fout, int(material.ToonSharingFlag), true)
		if err != nil {
			return err
		}
		if material.ToonSharingFlag == pmx.TOON_SHARING_SHARING {
			err = r.writeNumber(fout, textureIdxType, float64(material.ToonTextureIndex), 0.0, true)
		} else {
			err = r.writeByte(fout, int(material.ToonTextureIndex), true)
		}
		if err != nil {
			return err
		}
		err = r.writeText(fout, material.Memo, "")
		if err != nil {
			return err
		}
		err = r.writeNumber(fout, binaryType_int, float64(material.VerticesCount), 0.0, true)
		if err != nil {
			return err
		}
	}

	return nil
}

// saveBones ボーンデータの書き込み
func (r *PmxRepository) saveBones(fout *os.File, targetBones []*pmx.Bone, boneIdxType binaryType) error {
	err := r.writeNumber(fout, binaryType_int, float64(len(targetBones)), 0.0, true)
	if err != nil {
		return err
	}

	for i, bone := range targetBones {
		err = r.writeText(fout, bone.Name, fmt.Sprintf("Bone %d", i))
		if err != nil {
			return err
		}
		err = r.writeText(fout, bone.EnglishName, fmt.Sprintf("Bone %d", i))
		if err != nil {
			return err
		}
		err = r.writeNumber(fout, binaryType_float, bone.Position.GetX(), 0.0, false)
		if err != nil {
			return err
		}
		err = r.writeNumber(fout, binaryType_float, bone.Position.GetY(), 0.0, false)
		if err != nil {
			return err
		}
		err = r.writeNumber(fout, binaryType_float, bone.Position.GetZ(), 0.0, false)
		if err != nil {
			return err
		}
		err = r.writeNumber(fout, boneIdxType, float64(bone.ParentIndex), 0.0, false)
		if err != nil {
			return err
		}
		err = r.writeNumber(fout, binaryType_int, float64(bone.Layer), 0.0, true)
		if err != nil {
			return err
		}
		err = r.writeShort(fout, uint16(bone.BoneFlag))
		if err != nil {
			return err
		}

		if bone.IsTailBone() {
			err = r.writeNumber(fout, boneIdxType, float64(bone.TailIndex), 0.0, false)
			if err != nil {
				return err
			}
		} else {
			err = r.writeNumber(fout, binaryType_float, bone.TailPosition.GetX(), 0.0, false)
			if err != nil {
				return err
			}
			err = r.writeNumber(fout, binaryType_float, bone.TailPosition.GetY(), 0.0, false)
			if err != nil {
				return err
			}
			err = r.writeNumber(fout, binaryType_float, bone.TailPosition.GetZ(), 0.0, false)
			if err != nil {
				return err
			}
		}
		if bone.IsEffectorTranslation() || bone.IsEffectorRotation() {
			err = r.writeNumber(fout, boneIdxType, float64(bone.EffectIndex), 0.0, false)
			if err != nil {
				return err
			}
			err = r.writeNumber(fout, binaryType_float, bone.EffectFactor, 0.0, false)
			if err != nil {
				return err
			}
		}
		if bone.HasFixedAxis() {
			err = r.writeNumber(fout, binaryType_float, bone.FixedAxis.GetX(), 0.0, false)
			if err != nil {
				return err
			}
			err = r.writeNumber(fout, binaryType_float, bone.FixedAxis.GetY(), 0.0, false)
			if err != nil {
				return err
			}
			err = r.writeNumber(fout, binaryType_float, bone.FixedAxis.GetZ(), 0.0, false)
			if err != nil {
				return err
			}
		}
		if bone.HasLocalAxis() {
			err = r.writeNumber(fout, binaryType_float, bone.LocalAxisX.GetX(), 0.0, false)
			if err != nil {
				return err
			}
			err = r.writeNumber(fout, binaryType_float, bone.LocalAxisX.GetY(), 0.0, false)
			if err != nil {
				return err
			}
			err = r.writeNumber(fout, binaryType_float, bone.LocalAxisX.GetZ(), 0.0, false)
			if err != nil {
				return err
			}
			err = r.writeNumber(fout, binaryType_float, bone.LocalAxisZ.GetX(), 0.0, false)
			if err != nil {
				return err
			}
			err = r.writeNumber(fout, binaryType_float, bone.LocalAxisZ.GetY(), 0.0, false)
			if err != nil {
				return err
			}
			err = r.writeNumber(fout, binaryType_float, bone.LocalAxisZ.GetZ(), 0.0, false)
			if err != nil {
				return err
			}
		}
		if bone.IsEffectorParentDeform() {
			err = r.writeNumber(fout, binaryType_int, float64(bone.EffectorKey), 0.0, true)
			if err != nil {
				return err
			}
		}
		if bone.IsIK() {
			err = r.writeNumber(fout, boneIdxType, float64(bone.Ik.BoneIndex), 0.0, false)
			if err != nil {
				return err
			}
			err = r.writeNumber(fout, binaryType_int, float64(bone.Ik.LoopCount), 0.0, true)
			if err != nil {
				return err
			}
			err = r.writeNumber(fout, binaryType_float, bone.Ik.UnitRotation.GetRadians().GetX(), 0.0, false)
			if err != nil {
				return err
			}
			err = r.writeNumber(fout, binaryType_int, float64(len(bone.Ik.Links)), 0.0, true)
			if err != nil {
				return err
			}

			for _, link := range bone.Ik.Links {
				err = r.writeNumber(fout, boneIdxType, float64(link.BoneIndex), 0.0, false)
				if err != nil {
					return err
				}
				err = r.writeByte(fout, int(mmath.BoolToInt(link.AngleLimit)), true)
				if err != nil {
					return err
				}
				if link.AngleLimit {
					err = r.writeNumber(fout, binaryType_float, link.MinAngleLimit.GetRadians().GetX(), 0.0, false)
					if err != nil {
						return err
					}
					err = r.writeNumber(fout, binaryType_float, link.MinAngleLimit.GetRadians().GetY(), 0.0, false)
					if err != nil {
						return err
					}
					err = r.writeNumber(fout, binaryType_float, link.MinAngleLimit.GetRadians().GetZ(), 0.0, false)
					if err != nil {
						return err
					}
					err = r.writeNumber(fout, binaryType_float, link.MaxAngleLimit.GetRadians().GetX(), 0.0, false)
					if err != nil {
						return err
					}
					err = r.writeNumber(fout, binaryType_float, link.MaxAngleLimit.GetRadians().GetY(), 0.0, false)
					if err != nil {
						return err
					}
					err = r.writeNumber(fout, binaryType_float, link.MaxAngleLimit.GetRadians().GetZ(), 0.0, false)
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
func (r *PmxRepository) saveMorphs(
	fout *os.File, targetMorphs []*pmx.Morph, vertexIdxType, boneIdxType, materialIdxType, morphIdxType binaryType,
) error {
	err := r.writeNumber(fout, binaryType_int, float64(len(targetMorphs)), 0.0, true)
	if err != nil {
		return err
	}

	for i, morph := range targetMorphs {
		err = r.writeText(fout, morph.Name, fmt.Sprintf("Morph %d", i))
		if err != nil {
			return err
		}
		err = r.writeText(fout, morph.EnglishName, fmt.Sprintf("Morph %d", i))
		if err != nil {
			return err
		}
		err = r.writeByte(fout, int(morph.Panel), true)
		if err != nil {
			return err
		}
		err = r.writeByte(fout, int(morph.MorphType), true)
		if err != nil {
			return err
		}
		err = r.writeNumber(fout, binaryType_int, float64(len(morph.Offsets)), 0.0, true)
		if err != nil {
			return err
		}

		for _, offset := range morph.Offsets {
			switch off := offset.(type) {
			case *pmx.VertexMorphOffset:
				err = r.writeNumber(fout, vertexIdxType, float64(off.VertexIndex), 0.0, true)
				if err != nil {
					return err
				}
				err = r.writeNumber(fout, binaryType_float, off.Position.GetX(), 0.0, false)
				if err != nil {
					return err
				}
				err = r.writeNumber(fout, binaryType_float, off.Position.GetY(), 0.0, false)
				if err != nil {
					return err
				}
				err = r.writeNumber(fout, binaryType_float, off.Position.GetZ(), 0.0, false)
				if err != nil {
					return err
				}
			case *pmx.UvMorphOffset:
				err = r.writeNumber(fout, vertexIdxType, float64(off.VertexIndex), 0.0, true)
				if err != nil {
					return err
				}
				err = r.writeNumber(fout, binaryType_float, off.Uv.GetX(), 0.0, false)
				if err != nil {
					return err
				}
				err = r.writeNumber(fout, binaryType_float, off.Uv.GetY(), 0.0, false)
				if err != nil {
					return err
				}
				err = r.writeNumber(fout, binaryType_float, off.Uv.GetZ(), 0.0, false)
				if err != nil {
					return err
				}
				err = r.writeNumber(fout, binaryType_float, off.Uv.GetZ(), 0.0, false)
				if err != nil {
					return err
				}
			case *pmx.BoneMorphOffset:
				err = r.writeNumber(fout, boneIdxType, float64(off.BoneIndex), 0.0, true)
				if err != nil {
					return err
				}
				err = r.writeNumber(fout, binaryType_float, off.Position.GetX(), 0.0, false)
				if err != nil {
					return err
				}
				err = r.writeNumber(fout, binaryType_float, off.Position.GetY(), 0.0, false)
				if err != nil {
					return err
				}
				err = r.writeNumber(fout, binaryType_float, off.Position.GetZ(), 0.0, false)
				if err != nil {
					return err
				}
				err = r.writeNumber(fout, binaryType_float, off.Rotation.GetQuaternion().GetX(), 0.0, false)
				if err != nil {
					return err
				}
				err = r.writeNumber(fout, binaryType_float, off.Rotation.GetQuaternion().GetY(), 0.0, false)
				if err != nil {
					return err
				}
				err = r.writeNumber(fout, binaryType_float, off.Rotation.GetQuaternion().GetZ(), 0.0, false)
				if err != nil {
					return err
				}
				err = r.writeNumber(fout, binaryType_float, off.Rotation.GetQuaternion().GetW(), 0.0, false)
				if err != nil {
					return err
				}
			case *pmx.MaterialMorphOffset:
				err = r.writeNumber(fout, materialIdxType, float64(off.MaterialIndex), 0.0, true)
				if err != nil {
					return err
				}
				err = r.writeByte(fout, int(off.CalcMode), true)
				if err != nil {
					return err
				}
				err = r.writeNumber(fout, binaryType_float, off.Diffuse.GetX(), 0.0, false)
				if err != nil {
					return err
				}
				err = r.writeNumber(fout, binaryType_float, off.Diffuse.GetY(), 0.0, false)
				if err != nil {
					return err
				}
				err = r.writeNumber(fout, binaryType_float, off.Diffuse.GetZ(), 0.0, false)
				if err != nil {
					return err
				}
				err = r.writeNumber(fout, binaryType_float, off.Diffuse.GetW(), 0.0, false)
				if err != nil {
					return err
				}
				err = r.writeNumber(fout, binaryType_float, off.Specular.GetX(), 0.0, false)
				if err != nil {
					return err
				}
				err = r.writeNumber(fout, binaryType_float, off.Specular.GetY(), 0.0, false)
				if err != nil {
					return err
				}
				err = r.writeNumber(fout, binaryType_float, off.Specular.GetZ(), 0.0, false)
				if err != nil {
					return err
				}
				err = r.writeNumber(fout, binaryType_float, off.Specular.GetW(), 0.0, false)
				if err != nil {
					return err
				}
				err = r.writeNumber(fout, binaryType_float, off.Ambient.GetX(), 0.0, false)
				if err != nil {
					return err
				}
				err = r.writeNumber(fout, binaryType_float, off.Ambient.GetY(), 0.0, false)
				if err != nil {
					return err
				}
				err = r.writeNumber(fout, binaryType_float, off.Ambient.GetZ(), 0.0, false)
				if err != nil {
					return err
				}
				err = r.writeNumber(fout, binaryType_float, off.Edge.GetX(), 0.0, false)
				if err != nil {
					return err
				}
				err = r.writeNumber(fout, binaryType_float, off.Edge.GetY(), 0.0, false)
				if err != nil {
					return err
				}
				err = r.writeNumber(fout, binaryType_float, off.Edge.GetZ(), 0.0, false)
				if err != nil {
					return err
				}
				err = r.writeNumber(fout, binaryType_float, off.Edge.GetW(), 0.0, false)
				if err != nil {
					return err
				}
				err = r.writeNumber(fout, binaryType_float, off.EdgeSize, 0.0, false)
				if err != nil {
					return err
				}
				err = r.writeNumber(fout, binaryType_float, off.TextureFactor.GetX(), 0.0, false)
				if err != nil {
					return err
				}
				err = r.writeNumber(fout, binaryType_float, off.TextureFactor.GetY(), 0.0, false)
				if err != nil {
					return err
				}
				err = r.writeNumber(fout, binaryType_float, off.TextureFactor.GetZ(), 0.0, false)
				if err != nil {
					return err
				}
				err = r.writeNumber(fout, binaryType_float, off.TextureFactor.GetW(), 0.0, false)
				if err != nil {
					return err
				}
				err = r.writeNumber(fout, binaryType_float, off.SphereTextureFactor.GetX(), 0.0, false)
				if err != nil {
					return err
				}
				err = r.writeNumber(fout, binaryType_float, off.SphereTextureFactor.GetY(), 0.0, false)
				if err != nil {
					return err
				}
				err = r.writeNumber(fout, binaryType_float, off.SphereTextureFactor.GetZ(), 0.0, false)
				if err != nil {
					return err
				}
				err = r.writeNumber(fout, binaryType_float, off.SphereTextureFactor.GetW(), 0.0, false)
				if err != nil {
					return err
				}
				err = r.writeNumber(fout, binaryType_float, off.ToonTextureFactor.GetX(), 0.0, false)
				if err != nil {
					return err
				}
				err = r.writeNumber(fout, binaryType_float, off.ToonTextureFactor.GetY(), 0.0, false)
				if err != nil {
					return err
				}
				err = r.writeNumber(fout, binaryType_float, off.ToonTextureFactor.GetZ(), 0.0, false)
				if err != nil {
					return err
				}
				err = r.writeNumber(fout, binaryType_float, off.ToonTextureFactor.GetW(), 0.0, false)
				if err != nil {
					return err
				}
			case *pmx.GroupMorphOffset:
				err = r.writeNumber(fout, morphIdxType, float64(off.MorphIndex), 0.0, true)
				if err != nil {
					return err
				}
				err = r.writeNumber(fout, binaryType_float, off.MorphFactor, 0.0, false)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

// saveDisplaySlots 表示枠データの書き込み
func (r *PmxRepository) saveDisplaySlots(fout *os.File, model *pmx.PmxModel, boneIdxType, morphIdxType binaryType) error {
	err := r.writeNumber(fout, binaryType_int, float64(model.DisplaySlots.Len()), 0.0, true)
	if err != nil {
		return err
	}

	for i := range model.DisplaySlots.Len() {
		displaySlot := model.DisplaySlots.Get(i)

		err = r.writeText(fout, displaySlot.Name, fmt.Sprintf("Display %d", i))
		if err != nil {
			return err
		}
		err = r.writeText(fout, displaySlot.EnglishName, fmt.Sprintf("Display %d", i))
		if err != nil {
			return err
		}
		err = r.writeByte(fout, int(displaySlot.SpecialFlag), true)
		if err != nil {
			return err
		}
		err = r.writeNumber(fout, binaryType_int, float64(len(displaySlot.References)), 0.0, true)
		if err != nil {
			return err
		}

		for _, reference := range displaySlot.References {
			err = r.writeByte(fout, int(reference.DisplayType), true)
			if err != nil {
				return err
			}
			if reference.DisplayType == pmx.DISPLAY_TYPE_BONE {
				err = r.writeNumber(fout, boneIdxType, float64(reference.DisplayIndex), 0.0, true)
			} else {
				err = r.writeNumber(fout, morphIdxType, float64(reference.DisplayIndex), 0.0, true)
			}
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// saveRigidBodies 剛体データの書き込み
func (r *PmxRepository) saveRigidBodies(fout *os.File, model *pmx.PmxModel, boneIdxType binaryType) error {
	err := r.writeNumber(fout, binaryType_int, float64(model.RigidBodies.Len()), 0.0, true)
	if err != nil {
		return err
	}

	for i := range model.RigidBodies.Len() {
		rigidbody := model.RigidBodies.Get(i)
		err = r.writeText(fout, rigidbody.Name, fmt.Sprintf("Rigidbody %d", i))
		if err != nil {
			return err
		}
		err = r.writeText(fout, rigidbody.EnglishName, fmt.Sprintf("Rigidbody %d", i))
		if err != nil {
			return err
		}
		err = r.writeNumber(fout, boneIdxType, float64(rigidbody.BoneIndex), 0.0, false)
		if err != nil {
			return err
		}
		err = r.writeByte(fout, int(rigidbody.CollisionGroup), true)
		if err != nil {
			return err
		}
		err = r.writeShort(fout, uint16(rigidbody.CollisionGroupMaskValue))
		if err != nil {
			return err
		}
		err = r.writeByte(fout, int(rigidbody.ShapeType), true)
		if err != nil {
			return err
		}
		err = r.writeNumber(fout, binaryType_float, rigidbody.Size.GetX(), 0.0, true)
		if err != nil {
			return err
		}
		err = r.writeNumber(fout, binaryType_float, rigidbody.Size.GetY(), 0.0, true)
		if err != nil {
			return err
		}
		err = r.writeNumber(fout, binaryType_float, rigidbody.Size.GetZ(), 0.0, true)
		if err != nil {
			return err
		}
		err = r.writeNumber(fout, binaryType_float, rigidbody.Position.GetX(), 0.0, false)
		if err != nil {
			return err
		}
		err = r.writeNumber(fout, binaryType_float, rigidbody.Position.GetY(), 0.0, false)
		if err != nil {
			return err
		}
		err = r.writeNumber(fout, binaryType_float, rigidbody.Position.GetZ(), 0.0, false)
		if err != nil {
			return err
		}
		err = r.writeNumber(fout, binaryType_float, rigidbody.Rotation.GetRadians().GetX(), 0.0, false)
		if err != nil {
			return err
		}
		err = r.writeNumber(fout, binaryType_float, rigidbody.Rotation.GetRadians().GetY(), 0.0, false)
		if err != nil {
			return err
		}
		err = r.writeNumber(fout, binaryType_float, rigidbody.Rotation.GetRadians().GetZ(), 0.0, false)
		if err != nil {
			return err
		}
		err = r.writeNumber(fout, binaryType_float, rigidbody.RigidBodyParam.Mass, 0.0, true)
		if err != nil {
			return err
		}
		err = r.writeNumber(fout, binaryType_float, rigidbody.RigidBodyParam.LinearDamping, 0.0, true)
		if err != nil {
			return err
		}
		err = r.writeNumber(fout, binaryType_float, rigidbody.RigidBodyParam.AngularDamping, 0.0, true)
		if err != nil {
			return err
		}
		err = r.writeNumber(fout, binaryType_float, rigidbody.RigidBodyParam.Restitution, 0.0, true)
		if err != nil {
			return err
		}
		err = r.writeNumber(fout, binaryType_float, rigidbody.RigidBodyParam.Friction, 0.0, true)
		if err != nil {
			return err
		}
		err = r.writeByte(fout, int(rigidbody.PhysicsType), true)
		if err != nil {
			return err
		}
	}

	return nil
}

// saveJoints ジョイントデータの書き込み
func (r *PmxRepository) saveJoints(fout *os.File, model *pmx.PmxModel, rigidbodyIdxType binaryType) error {
	err := r.writeNumber(fout, binaryType_int, float64(model.Joints.Len()), 0.0, true)
	if err != nil {
		return err
	}

	for i := range model.Joints.Len() {
		joint := model.Joints.Get(i)

		err = r.writeText(fout, joint.Name, fmt.Sprintf("Joint %d", i))
		if err != nil {
			return err
		}
		err = r.writeText(fout, joint.EnglishName, fmt.Sprintf("Joint %d", i))
		if err != nil {
			return err
		}
		err = r.writeByte(fout, int(joint.JointType), true)
		if err != nil {
			return err
		}
		err = r.writeNumber(fout, rigidbodyIdxType, float64(joint.RigidbodyIndexA), 0.0, false)
		if err != nil {
			return err
		}
		err = r.writeNumber(fout, rigidbodyIdxType, float64(joint.RigidbodyIndexB), 0.0, false)
		if err != nil {
			return err
		}
		err = r.writeNumber(fout, binaryType_float, joint.Position.GetX(), 0.0, false)
		if err != nil {
			return err
		}
		err = r.writeNumber(fout, binaryType_float, joint.Position.GetY(), 0.0, false)
		if err != nil {
			return err
		}
		err = r.writeNumber(fout, binaryType_float, joint.Position.GetZ(), 0.0, false)
		if err != nil {
			return err
		}
		err = r.writeNumber(fout, binaryType_float, joint.Rotation.GetRadians().GetX(), 0.0, false)
		if err != nil {
			return err
		}
		err = r.writeNumber(fout, binaryType_float, joint.Rotation.GetRadians().GetY(), 0.0, false)
		if err != nil {
			return err
		}
		err = r.writeNumber(fout, binaryType_float, joint.Rotation.GetRadians().GetZ(), 0.0, false)
		if err != nil {
			return err
		}
		err = r.writeNumber(fout, binaryType_float, joint.JointParam.TranslationLimitMin.GetX(), 0.0, false)
		if err != nil {
			return err
		}
		err = r.writeNumber(fout, binaryType_float, joint.JointParam.TranslationLimitMin.GetY(), 0.0, false)
		if err != nil {
			return err
		}
		err = r.writeNumber(fout, binaryType_float, joint.JointParam.TranslationLimitMin.GetZ(), 0.0, false)
		if err != nil {
			return err
		}
		err = r.writeNumber(fout, binaryType_float, joint.JointParam.TranslationLimitMax.GetX(), 0.0, false)
		if err != nil {
			return err
		}
		err = r.writeNumber(fout, binaryType_float, joint.JointParam.TranslationLimitMax.GetY(), 0.0, false)
		if err != nil {
			return err
		}
		err = r.writeNumber(fout, binaryType_float, joint.JointParam.TranslationLimitMax.GetZ(), 0.0, false)
		if err != nil {
			return err
		}
		err = r.writeNumber(fout, binaryType_float, joint.JointParam.RotationLimitMin.GetRadians().GetX(), 0.0, false)
		if err != nil {
			return err
		}
		err = r.writeNumber(fout, binaryType_float, joint.JointParam.RotationLimitMin.GetRadians().GetY(), 0.0, false)
		if err != nil {
			return err
		}
		err = r.writeNumber(fout, binaryType_float, joint.JointParam.RotationLimitMin.GetRadians().GetZ(), 0.0, false)
		if err != nil {
			return err
		}
		err = r.writeNumber(fout, binaryType_float, joint.JointParam.RotationLimitMax.GetRadians().GetX(), 0.0, false)
		if err != nil {
			return err
		}
		err = r.writeNumber(fout, binaryType_float, joint.JointParam.RotationLimitMax.GetRadians().GetY(), 0.0, false)
		if err != nil {
			return err
		}
		err = r.writeNumber(fout, binaryType_float, joint.JointParam.RotationLimitMax.GetRadians().GetZ(), 0.0, false)
		if err != nil {
			return err
		}
		err = r.writeNumber(fout, binaryType_float, joint.JointParam.SpringConstantTranslation.GetX(), 0.0, false)
		if err != nil {
			return err
		}
		err = r.writeNumber(fout, binaryType_float, joint.JointParam.SpringConstantTranslation.GetY(), 0.0, false)
		if err != nil {
			return err
		}
		err = r.writeNumber(fout, binaryType_float, joint.JointParam.SpringConstantTranslation.GetZ(), 0.0, false)
		if err != nil {
			return err
		}
		err = r.writeNumber(fout, binaryType_float, joint.JointParam.SpringConstantRotation.GetX(), 0.0, false)
		if err != nil {
			return err
		}
		err = r.writeNumber(fout, binaryType_float, joint.JointParam.SpringConstantRotation.GetY(), 0.0, false)
		if err != nil {
			return err
		}
		err = r.writeNumber(fout, binaryType_float, joint.JointParam.SpringConstantRotation.GetZ(), 0.0, false)
		if err != nil {
			return err
		}
	}

	return nil
}

// ------------------------------
func (r *PmxRepository) writeText(fout *os.File, text string, defaultText string) error {
	var binaryTxt []byte
	var err error

	// エンコードの試行
	binaryTxt, err = r.encodeUTF16LE(text)
	if err != nil {
		binaryTxt, _ = r.encodeUTF16LE(defaultText)
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

func (r *PmxRepository) encodeUTF16LE(s string) ([]byte, error) {
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

func (r *PmxRepository) defineWriteIndexForVertex(size int) (int, binaryType) {
	if size < 256 {
		return 1, binaryType_unsignedByte
	} else if size <= 65535 {
		return 2, binaryType_unsignedShort
	}
	return 4, binaryType_int
}

func (r *PmxRepository) defineWriteIndexForOthers(size int) (int, binaryType) {
	if size < 128 {
		return 1, binaryType_byte
	} else if size <= 32767 {
		return 2, binaryType_short
	}
	return 4, binaryType_int
}
