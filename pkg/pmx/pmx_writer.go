package pmx

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math"
	"os"
	"unicode/utf16"

	"github.com/miu200521358/mlib_go/pkg/mmath"
	"github.com/miu200521358/mlib_go/pkg/mutils"
	"github.com/miu200521358/mlib_go/pkg/mutils/mlog"
)

type pmxBinaryType string

const (
	pmxBinaryType_float         pmxBinaryType = "<f"
	pmxBinaryType_byte          pmxBinaryType = "<b"
	pmxBinaryType_unsignedByte  pmxBinaryType = "<B"
	pmxBinaryType_short         pmxBinaryType = "<h"
	pmxBinaryType_unsignedShort pmxBinaryType = "<H"
	pmxBinaryType_int           pmxBinaryType = "<i"
	pmxBinaryType_unsignedInt   pmxBinaryType = "<I"
	pmxBinaryType_long          pmxBinaryType = "<l"
	pmxBinaryType_unsignedLong  pmxBinaryType = "<L"
)

func (model *PmxModel) Save(includeSystem bool, overridePath string) error {
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

	filteredBones := []*Bone{}
	for i := range model.Bones.Len() {
		bone := model.Bones.Get(i)
		if (!includeSystem && !bone.IsSystem) || includeSystem {
			filteredBones = append(filteredBones, bone)
		}
	}

	filteredMorphs := []*Morph{}
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

	err = writeNumber(fout, pmxBinaryType_float, 2.0, 0.0, true)
	if err != nil {
		return err
	}

	err = writeByte(fout, 8, true)
	if err != nil {
		return err
	}

	err = writeByte(fout, 0, true)
	if err != nil {
		return err
	}

	err = writeByte(fout, model.ExtendedUVCount, true)
	if err != nil {
		return err
	}

	vertexIdxSize, vertexIdxType := defineWriteIndexForVertex(model.Vertices.Len())
	err = writeByte(fout, vertexIdxSize, true)
	if err != nil {
		return err
	}

	textureIdxSize, textureIdxType := defineWriteIndexForOthers(model.Textures.Len())
	err = writeByte(fout, textureIdxSize, true)
	if err != nil {
		return err
	}

	materialIdxSize, materialIdxType := defineWriteIndexForOthers(model.Materials.Len())
	err = writeByte(fout, materialIdxSize, true)
	if err != nil {
		return err
	}

	boneIdxSize, boneIdxType := defineWriteIndexForOthers(len(filteredBones))
	err = writeByte(fout, boneIdxSize, true)
	if err != nil {
		return err
	}

	morphIdxSize, morphIdxType := defineWriteIndexForOthers(len(filteredMorphs))
	err = writeByte(fout, morphIdxSize, true)
	if err != nil {
		return err
	}

	rigidbodyIdxSize, rigidbodyIdxType := defineWriteIndexForOthers(model.RigidBodies.Len())
	err = writeByte(fout, rigidbodyIdxSize, true)
	if err != nil {
		return err
	}

	err = writeText(fout, model.Name, "Pmx Model")
	if err != nil {
		return err
	}

	err = writeText(fout, model.EnglishName, "Pmx Model")
	if err != nil {
		return err
	}

	err = writeText(fout, model.Comment, "")
	if err != nil {
		return err
	}

	err = writeText(fout, model.EnglishComment, "")
	if err != nil {
		return err
	}

	err = writeVertices(fout, model, boneIdxType)
	if err != nil {
		return err
	}

	err = writeFaces(fout, model, vertexIdxType)
	if err != nil {
		return err
	}

	err = writeTextures(fout, model)
	if err != nil {
		return err
	}

	err = writeMaterials(fout, model, textureIdxType)
	if err != nil {
		return err
	}

	err = writeBones(fout, filteredBones, boneIdxType)
	if err != nil {
		return err
	}

	err = writeMorphs(fout, filteredMorphs, vertexIdxType, boneIdxType, materialIdxType, morphIdxType)
	if err != nil {
		return err
	}

	err = writeDisplaySlots(fout, model, boneIdxType, morphIdxType)
	if err != nil {
		return err
	}

	err = writeRigidBodies(fout, model, boneIdxType)
	if err != nil {
		return err
	}

	err = writeJoints(fout, model, rigidbodyIdxType)
	if err != nil {
		return err
	}

	return nil
}

// writeVertices 頂点データの書き込み
func writeVertices(fout *os.File, model *PmxModel, boneIdxType pmxBinaryType) error {
	err := writeNumber(fout, pmxBinaryType_int, float64(model.Vertices.Len()), 0.0, true)
	if err != nil {
		return err
	}

	for i := range model.Vertices.Len() {
		vertex := model.Vertices.Get(i)

		writeNumber(fout, pmxBinaryType_float, vertex.Position.GetX(), 0.0, false)
		writeNumber(fout, pmxBinaryType_float, vertex.Position.GetY(), 0.0, false)
		writeNumber(fout, pmxBinaryType_float, vertex.Position.GetZ(), 0.0, false)
		writeNumber(fout, pmxBinaryType_float, vertex.Normal.GetX(), 0.0, false)
		writeNumber(fout, pmxBinaryType_float, vertex.Normal.GetY(), 0.0, false)
		writeNumber(fout, pmxBinaryType_float, vertex.Normal.GetZ(), 0.0, false)
		writeNumber(fout, pmxBinaryType_float, vertex.Uv.GetX(), 0.0, false)
		writeNumber(fout, pmxBinaryType_float, vertex.Uv.GetY(), 0.0, false)

		for _, uv := range vertex.ExtendedUvs {
			writeNumber(fout, pmxBinaryType_float, uv.GetX(), 0.0, false)
			writeNumber(fout, pmxBinaryType_float, uv.GetY(), 0.0, false)
			writeNumber(fout, pmxBinaryType_float, uv.GetZ(), 0.0, false)
			writeNumber(fout, pmxBinaryType_float, uv.GetW(), 0.0, false)
		}

		for j := len(vertex.ExtendedUvs); j < model.ExtendedUVCount; j++ {
			writeNumber(fout, pmxBinaryType_float, 0.0, 0.0, false)
			writeNumber(fout, pmxBinaryType_float, 0.0, 0.0, false)
			writeNumber(fout, pmxBinaryType_float, 0.0, 0.0, false)
			writeNumber(fout, pmxBinaryType_float, 0.0, 0.0, false)
		}

		switch v := vertex.Deform.(type) {
		case *Bdef1:
			writeByte(fout, 0, false)
			writeNumber(fout, boneIdxType, float64(v.Indexes[0]), 0.0, false)
		case *Bdef2:
			writeByte(fout, 1, false)
			writeNumber(fout, boneIdxType, float64(v.Indexes[0]), 0.0, false)
			writeNumber(fout, boneIdxType, float64(v.Indexes[1]), 0.0, false)
			writeNumber(fout, pmxBinaryType_float, v.Weights[0], 0.0, true)
		case *Bdef4:
			writeByte(fout, 2, false)
			writeNumber(fout, boneIdxType, float64(v.Indexes[0]), 0.0, false)
			writeNumber(fout, boneIdxType, float64(v.Indexes[1]), 0.0, false)
			writeNumber(fout, boneIdxType, float64(v.Indexes[2]), 0.0, false)
			writeNumber(fout, boneIdxType, float64(v.Indexes[3]), 0.0, false)
			writeNumber(fout, pmxBinaryType_float, v.Weights[0], 0.0, true)
			writeNumber(fout, pmxBinaryType_float, v.Weights[1], 0.0, true)
			writeNumber(fout, pmxBinaryType_float, v.Weights[2], 0.0, true)
			writeNumber(fout, pmxBinaryType_float, v.Weights[3], 0.0, true)
		case *Sdef:
			writeByte(fout, 3, false)
			writeNumber(fout, boneIdxType, float64(v.Indexes[0]), 0.0, false)
			writeNumber(fout, boneIdxType, float64(v.Indexes[1]), 0.0, false)
			writeNumber(fout, pmxBinaryType_float, v.Weights[0], 0.0, true)
			writeNumber(fout, pmxBinaryType_float, v.SdefC.GetX(), 0.0, false)
			writeNumber(fout, pmxBinaryType_float, v.SdefC.GetY(), 0.0, false)
			writeNumber(fout, pmxBinaryType_float, v.SdefC.GetZ(), 0.0, false)
			writeNumber(fout, pmxBinaryType_float, v.SdefR0.GetX(), 0.0, false)
			writeNumber(fout, pmxBinaryType_float, v.SdefR0.GetY(), 0.0, false)
			writeNumber(fout, pmxBinaryType_float, v.SdefR0.GetZ(), 0.0, false)
			writeNumber(fout, pmxBinaryType_float, v.SdefR1.GetX(), 0.0, false)
			writeNumber(fout, pmxBinaryType_float, v.SdefR1.GetY(), 0.0, false)
			writeNumber(fout, pmxBinaryType_float, v.SdefR1.GetZ(), 0.0, false)
		default:
			mlog.W("頂点deformなし: %v\n", vertex)
		}

		writeNumber(fout, pmxBinaryType_float, vertex.EdgeFactor, 0.0, true)
	}

	return nil
}

// writeFaces 面データの書き込み
func writeFaces(fout *os.File, model *PmxModel, vertexIdxType pmxBinaryType) error {
	err := writeNumber(fout, pmxBinaryType_int, float64(model.Faces.Len()*3), 0.0, true)
	if err != nil {
		return err
	}

	for i := range model.Faces.Len() {
		face := model.Faces.Get(i)

		for _, vidx := range face.VertexIndexes {
			err = writeNumber(fout, vertexIdxType, float64(vidx), 0.0, true)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// writeTextures テクスチャデータの書き込み
func writeTextures(fout *os.File, model *PmxModel) error {
	err := writeNumber(fout, pmxBinaryType_int, float64(model.Textures.Len()), 0.0, true)
	if err != nil {
		return err
	}

	for i := range model.Textures.Len() {
		texture := model.Textures.Get(i)
		err = writeText(fout, texture.Name, "")
		if err != nil {
			return err
		}
	}

	return nil
}

// writeMaterials 材質データの書き込み
func writeMaterials(fout *os.File, model *PmxModel, textureIdxType pmxBinaryType) error {
	err := writeNumber(fout, pmxBinaryType_int, float64(model.Materials.Len()), 0.0, true)
	if err != nil {
		return err
	}

	for i := range model.Materials.Len() {
		material := model.Materials.Get(i)
		err = writeText(fout, material.Name, fmt.Sprintf("Material %d", i))
		if err != nil {
			return err
		}
		err = writeText(fout, material.EnglishName, fmt.Sprintf("Material %d", i))
		if err != nil {
			return err
		}
		err = writeNumber(fout, pmxBinaryType_float, material.Diffuse.GetX(), 0.0, true)
		if err != nil {
			return err
		}
		err = writeNumber(fout, pmxBinaryType_float, material.Diffuse.GetY(), 0.0, true)
		if err != nil {
			return err
		}
		err = writeNumber(fout, pmxBinaryType_float, material.Diffuse.GetZ(), 0.0, true)
		if err != nil {
			return err
		}
		err = writeNumber(fout, pmxBinaryType_float, material.Diffuse.GetW(), 0.0, true)
		if err != nil {
			return err
		}
		err = writeNumber(fout, pmxBinaryType_float, material.Specular.GetX(), 0.0, true)
		if err != nil {
			return err
		}
		err = writeNumber(fout, pmxBinaryType_float, material.Specular.GetY(), 0.0, true)
		if err != nil {
			return err
		}
		err = writeNumber(fout, pmxBinaryType_float, material.Specular.GetZ(), 0.0, true)
		if err != nil {
			return err
		}
		err = writeNumber(fout, pmxBinaryType_float, material.Specular.GetW(), 0.0, true)
		if err != nil {
			return err
		}
		err = writeNumber(fout, pmxBinaryType_float, material.Ambient.GetX(), 0.0, true)
		if err != nil {
			return err
		}
		err = writeNumber(fout, pmxBinaryType_float, material.Ambient.GetY(), 0.0, true)
		if err != nil {
			return err
		}
		err = writeNumber(fout, pmxBinaryType_float, material.Ambient.GetZ(), 0.0, true)
		if err != nil {
			return err
		}
		err = writeByte(fout, int(material.DrawFlag), true)
		if err != nil {
			return err
		}
		err = writeNumber(fout, pmxBinaryType_float, material.Edge.GetX(), 0.0, true)
		if err != nil {
			return err
		}
		err = writeNumber(fout, pmxBinaryType_float, material.Edge.GetY(), 0.0, true)
		if err != nil {
			return err
		}
		err = writeNumber(fout, pmxBinaryType_float, material.Edge.GetZ(), 0.0, true)
		if err != nil {
			return err
		}
		err = writeNumber(fout, pmxBinaryType_float, material.Edge.GetW(), 0.0, true)
		if err != nil {
			return err
		}
		err = writeNumber(fout, pmxBinaryType_float, material.EdgeSize, 0.0, true)
		if err != nil {
			return err
		}
		err = writeNumber(fout, textureIdxType, float64(material.TextureIndex), 0.0, true)
		if err != nil {
			return err
		}
		err = writeNumber(fout, textureIdxType, float64(material.SphereTextureIndex), 0.0, true)
		if err != nil {
			return err
		}
		err = writeByte(fout, int(material.SphereMode), true)
		if err != nil {
			return err
		}
		err = writeByte(fout, int(material.ToonSharingFlag), true)
		if err != nil {
			return err
		}
		if material.ToonSharingFlag == TOON_SHARING_SHARING {
			err = writeNumber(fout, textureIdxType, float64(material.ToonTextureIndex), 0.0, true)
		} else {
			err = writeByte(fout, int(material.ToonTextureIndex), true)
		}
		if err != nil {
			return err
		}
		err = writeText(fout, material.Memo, "")
		if err != nil {
			return err
		}
		err = writeNumber(fout, pmxBinaryType_int, float64(material.VerticesCount), 0.0, true)
		if err != nil {
			return err
		}
	}

	return nil
}

// writeBones ボーンデータの書き込み
func writeBones(fout *os.File, targetBones []*Bone, boneIdxType pmxBinaryType) error {
	err := writeNumber(fout, pmxBinaryType_int, float64(len(targetBones)), 0.0, true)
	if err != nil {
		return err
	}

	for i, bone := range targetBones {
		err = writeText(fout, bone.Name, fmt.Sprintf("Bone %d", i))
		if err != nil {
			return err
		}
		err = writeText(fout, bone.EnglishName, fmt.Sprintf("Bone %d", i))
		if err != nil {
			return err
		}
		err = writeNumber(fout, pmxBinaryType_float, bone.Position.GetX(), 0.0, false)
		if err != nil {
			return err
		}
		err = writeNumber(fout, pmxBinaryType_float, bone.Position.GetY(), 0.0, false)
		if err != nil {
			return err
		}
		err = writeNumber(fout, pmxBinaryType_float, bone.Position.GetZ(), 0.0, false)
		if err != nil {
			return err
		}
		err = writeNumber(fout, boneIdxType, float64(bone.ParentIndex), 0.0, false)
		if err != nil {
			return err
		}
		err = writeNumber(fout, pmxBinaryType_int, float64(bone.Layer), 0.0, true)
		if err != nil {
			return err
		}
		err = writeShort(fout, uint16(bone.BoneFlag))
		if err != nil {
			return err
		}

		if bone.IsTailBone() {
			err = writeNumber(fout, boneIdxType, float64(bone.TailIndex), 0.0, false)
			if err != nil {
				return err
			}
		} else {
			err = writeNumber(fout, pmxBinaryType_float, bone.TailPosition.GetX(), 0.0, false)
			if err != nil {
				return err
			}
			err = writeNumber(fout, pmxBinaryType_float, bone.TailPosition.GetY(), 0.0, false)
			if err != nil {
				return err
			}
			err = writeNumber(fout, pmxBinaryType_float, bone.TailPosition.GetZ(), 0.0, false)
			if err != nil {
				return err
			}
		}
		if bone.IsEffectorTranslation() || bone.IsEffectorRotation() {
			err = writeNumber(fout, boneIdxType, float64(bone.EffectIndex), 0.0, false)
			if err != nil {
				return err
			}
			err = writeNumber(fout, pmxBinaryType_float, bone.EffectFactor, 0.0, false)
			if err != nil {
				return err
			}
		}
		if bone.HasFixedAxis() {
			err = writeNumber(fout, pmxBinaryType_float, bone.FixedAxis.GetX(), 0.0, false)
			if err != nil {
				return err
			}
			err = writeNumber(fout, pmxBinaryType_float, bone.FixedAxis.GetY(), 0.0, false)
			if err != nil {
				return err
			}
			err = writeNumber(fout, pmxBinaryType_float, bone.FixedAxis.GetZ(), 0.0, false)
			if err != nil {
				return err
			}
		}
		if bone.HasLocalAxis() {
			err = writeNumber(fout, pmxBinaryType_float, bone.LocalAxisX.GetX(), 0.0, false)
			if err != nil {
				return err
			}
			err = writeNumber(fout, pmxBinaryType_float, bone.LocalAxisX.GetY(), 0.0, false)
			if err != nil {
				return err
			}
			err = writeNumber(fout, pmxBinaryType_float, bone.LocalAxisX.GetZ(), 0.0, false)
			if err != nil {
				return err
			}
			err = writeNumber(fout, pmxBinaryType_float, bone.LocalAxisZ.GetX(), 0.0, false)
			if err != nil {
				return err
			}
			err = writeNumber(fout, pmxBinaryType_float, bone.LocalAxisZ.GetY(), 0.0, false)
			if err != nil {
				return err
			}
			err = writeNumber(fout, pmxBinaryType_float, bone.LocalAxisZ.GetZ(), 0.0, false)
			if err != nil {
				return err
			}
		}
		if bone.IsEffectorParentDeform() {
			err = writeNumber(fout, pmxBinaryType_int, float64(bone.EffectorKey), 0.0, true)
			if err != nil {
				return err
			}
		}
		if bone.IsIK() {
			err = writeNumber(fout, boneIdxType, float64(bone.Ik.BoneIndex), 0.0, false)
			if err != nil {
				return err
			}
			err = writeNumber(fout, pmxBinaryType_int, float64(bone.Ik.LoopCount), 0.0, true)
			if err != nil {
				return err
			}
			err = writeNumber(fout, pmxBinaryType_float, bone.Ik.UnitRotation.GetRadians().GetX(), 0.0, false)
			if err != nil {
				return err
			}
			err = writeNumber(fout, pmxBinaryType_int, float64(len(bone.Ik.Links)), 0.0, true)
			if err != nil {
				return err
			}

			for _, link := range bone.Ik.Links {
				err = writeNumber(fout, boneIdxType, float64(link.BoneIndex), 0.0, false)
				if err != nil {
					return err
				}
				err = writeByte(fout, int(mmath.BoolToInt(link.AngleLimit)), true)
				if err != nil {
					return err
				}
				if link.AngleLimit {
					err = writeNumber(fout, pmxBinaryType_float, link.MinAngleLimit.GetRadians().GetX(), 0.0, false)
					if err != nil {
						return err
					}
					err = writeNumber(fout, pmxBinaryType_float, link.MinAngleLimit.GetRadians().GetY(), 0.0, false)
					if err != nil {
						return err
					}
					err = writeNumber(fout, pmxBinaryType_float, link.MinAngleLimit.GetRadians().GetZ(), 0.0, false)
					if err != nil {
						return err
					}
					err = writeNumber(fout, pmxBinaryType_float, link.MaxAngleLimit.GetRadians().GetX(), 0.0, false)
					if err != nil {
						return err
					}
					err = writeNumber(fout, pmxBinaryType_float, link.MaxAngleLimit.GetRadians().GetY(), 0.0, false)
					if err != nil {
						return err
					}
					err = writeNumber(fout, pmxBinaryType_float, link.MaxAngleLimit.GetRadians().GetZ(), 0.0, false)
					if err != nil {
						return err
					}
				}
			}
		}
	}

	return nil
}

// writeMorphs モーフデータの書き込み
func writeMorphs(
	fout *os.File, targetMorphs []*Morph, vertexIdxType, boneIdxType, materialIdxType, morphIdxType pmxBinaryType,
) error {
	err := writeNumber(fout, pmxBinaryType_int, float64(len(targetMorphs)), 0.0, true)
	if err != nil {
		return err
	}

	for i, morph := range targetMorphs {
		err = writeText(fout, morph.Name, fmt.Sprintf("Morph %d", i))
		if err != nil {
			return err
		}
		err = writeText(fout, morph.EnglishName, fmt.Sprintf("Morph %d", i))
		if err != nil {
			return err
		}
		err = writeByte(fout, int(morph.Panel), true)
		if err != nil {
			return err
		}
		err = writeByte(fout, int(morph.MorphType), true)
		if err != nil {
			return err
		}
		err = writeNumber(fout, pmxBinaryType_int, float64(len(morph.Offsets)), 0.0, true)
		if err != nil {
			return err
		}

		for _, offset := range morph.Offsets {
			switch off := offset.(type) {
			case *VertexMorphOffset:
				err = writeNumber(fout, vertexIdxType, float64(off.VertexIndex), 0.0, true)
				if err != nil {
					return err
				}
				err = writeNumber(fout, pmxBinaryType_float, off.Position.GetX(), 0.0, false)
				if err != nil {
					return err
				}
				err = writeNumber(fout, pmxBinaryType_float, off.Position.GetY(), 0.0, false)
				if err != nil {
					return err
				}
				err = writeNumber(fout, pmxBinaryType_float, off.Position.GetZ(), 0.0, false)
				if err != nil {
					return err
				}
			case *UvMorphOffset:
				err = writeNumber(fout, vertexIdxType, float64(off.VertexIndex), 0.0, true)
				if err != nil {
					return err
				}
				err = writeNumber(fout, pmxBinaryType_float, off.Uv.GetX(), 0.0, false)
				if err != nil {
					return err
				}
				err = writeNumber(fout, pmxBinaryType_float, off.Uv.GetY(), 0.0, false)
				if err != nil {
					return err
				}
				err = writeNumber(fout, pmxBinaryType_float, off.Uv.GetZ(), 0.0, false)
				if err != nil {
					return err
				}
				err = writeNumber(fout, pmxBinaryType_float, off.Uv.GetZ(), 0.0, false)
				if err != nil {
					return err
				}
			case *BoneMorphOffset:
				err = writeNumber(fout, boneIdxType, float64(off.BoneIndex), 0.0, true)
				if err != nil {
					return err
				}
				err = writeNumber(fout, pmxBinaryType_float, off.Position.GetX(), 0.0, false)
				if err != nil {
					return err
				}
				err = writeNumber(fout, pmxBinaryType_float, off.Position.GetY(), 0.0, false)
				if err != nil {
					return err
				}
				err = writeNumber(fout, pmxBinaryType_float, off.Position.GetZ(), 0.0, false)
				if err != nil {
					return err
				}
				err = writeNumber(fout, pmxBinaryType_float, off.Rotation.GetQuaternion().GetX(), 0.0, false)
				if err != nil {
					return err
				}
				err = writeNumber(fout, pmxBinaryType_float, off.Rotation.GetQuaternion().GetY(), 0.0, false)
				if err != nil {
					return err
				}
				err = writeNumber(fout, pmxBinaryType_float, off.Rotation.GetQuaternion().GetZ(), 0.0, false)
				if err != nil {
					return err
				}
				err = writeNumber(fout, pmxBinaryType_float, off.Rotation.GetQuaternion().GetW(), 0.0, false)
				if err != nil {
					return err
				}
			case *MaterialMorphOffset:
				err = writeNumber(fout, materialIdxType, float64(off.MaterialIndex), 0.0, true)
				if err != nil {
					return err
				}
				err = writeByte(fout, int(off.CalcMode), true)
				if err != nil {
					return err
				}
				err = writeNumber(fout, pmxBinaryType_float, off.Diffuse.GetX(), 0.0, false)
				if err != nil {
					return err
				}
				err = writeNumber(fout, pmxBinaryType_float, off.Diffuse.GetY(), 0.0, false)
				if err != nil {
					return err
				}
				err = writeNumber(fout, pmxBinaryType_float, off.Diffuse.GetZ(), 0.0, false)
				if err != nil {
					return err
				}
				err = writeNumber(fout, pmxBinaryType_float, off.Diffuse.GetW(), 0.0, false)
				if err != nil {
					return err
				}
				err = writeNumber(fout, pmxBinaryType_float, off.Specular.GetX(), 0.0, false)
				if err != nil {
					return err
				}
				err = writeNumber(fout, pmxBinaryType_float, off.Specular.GetY(), 0.0, false)
				if err != nil {
					return err
				}
				err = writeNumber(fout, pmxBinaryType_float, off.Specular.GetZ(), 0.0, false)
				if err != nil {
					return err
				}
				err = writeNumber(fout, pmxBinaryType_float, off.Specular.GetW(), 0.0, false)
				if err != nil {
					return err
				}
				err = writeNumber(fout, pmxBinaryType_float, off.Ambient.GetX(), 0.0, false)
				if err != nil {
					return err
				}
				err = writeNumber(fout, pmxBinaryType_float, off.Ambient.GetY(), 0.0, false)
				if err != nil {
					return err
				}
				err = writeNumber(fout, pmxBinaryType_float, off.Ambient.GetZ(), 0.0, false)
				if err != nil {
					return err
				}
				err = writeNumber(fout, pmxBinaryType_float, off.Edge.GetX(), 0.0, false)
				if err != nil {
					return err
				}
				err = writeNumber(fout, pmxBinaryType_float, off.Edge.GetY(), 0.0, false)
				if err != nil {
					return err
				}
				err = writeNumber(fout, pmxBinaryType_float, off.Edge.GetZ(), 0.0, false)
				if err != nil {
					return err
				}
				err = writeNumber(fout, pmxBinaryType_float, off.Edge.GetW(), 0.0, false)
				if err != nil {
					return err
				}
				err = writeNumber(fout, pmxBinaryType_float, off.EdgeSize, 0.0, false)
				if err != nil {
					return err
				}
				err = writeNumber(fout, pmxBinaryType_float, off.TextureFactor.GetX(), 0.0, false)
				if err != nil {
					return err
				}
				err = writeNumber(fout, pmxBinaryType_float, off.TextureFactor.GetY(), 0.0, false)
				if err != nil {
					return err
				}
				err = writeNumber(fout, pmxBinaryType_float, off.TextureFactor.GetZ(), 0.0, false)
				if err != nil {
					return err
				}
				err = writeNumber(fout, pmxBinaryType_float, off.TextureFactor.GetW(), 0.0, false)
				if err != nil {
					return err
				}
				err = writeNumber(fout, pmxBinaryType_float, off.SphereTextureFactor.GetX(), 0.0, false)
				if err != nil {
					return err
				}
				err = writeNumber(fout, pmxBinaryType_float, off.SphereTextureFactor.GetY(), 0.0, false)
				if err != nil {
					return err
				}
				err = writeNumber(fout, pmxBinaryType_float, off.SphereTextureFactor.GetZ(), 0.0, false)
				if err != nil {
					return err
				}
				err = writeNumber(fout, pmxBinaryType_float, off.SphereTextureFactor.GetW(), 0.0, false)
				if err != nil {
					return err
				}
				err = writeNumber(fout, pmxBinaryType_float, off.ToonTextureFactor.GetX(), 0.0, false)
				if err != nil {
					return err
				}
				err = writeNumber(fout, pmxBinaryType_float, off.ToonTextureFactor.GetY(), 0.0, false)
				if err != nil {
					return err
				}
				err = writeNumber(fout, pmxBinaryType_float, off.ToonTextureFactor.GetZ(), 0.0, false)
				if err != nil {
					return err
				}
				err = writeNumber(fout, pmxBinaryType_float, off.ToonTextureFactor.GetW(), 0.0, false)
				if err != nil {
					return err
				}
			case *GroupMorphOffset:
				err = writeNumber(fout, morphIdxType, float64(off.MorphIndex), 0.0, true)
				if err != nil {
					return err
				}
				err = writeNumber(fout, pmxBinaryType_float, off.MorphFactor, 0.0, false)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

// writeDisplaySlots 表示枠データの書き込み
func writeDisplaySlots(fout *os.File, model *PmxModel, boneIdxType, morphIdxType pmxBinaryType) error {
	err := writeNumber(fout, pmxBinaryType_int, float64(model.DisplaySlots.Len()), 0.0, true)
	if err != nil {
		return err
	}

	for i := range model.DisplaySlots.Len() {
		displaySlot := model.DisplaySlots.Get(i)

		err = writeText(fout, displaySlot.Name, fmt.Sprintf("Display %d", i))
		if err != nil {
			return err
		}
		err = writeText(fout, displaySlot.EnglishName, fmt.Sprintf("Display %d", i))
		if err != nil {
			return err
		}
		err = writeByte(fout, int(displaySlot.SpecialFlag), true)
		if err != nil {
			return err
		}
		err = writeNumber(fout, pmxBinaryType_int, float64(len(displaySlot.References)), 0.0, true)
		if err != nil {
			return err
		}

		for _, reference := range displaySlot.References {
			err = writeByte(fout, int(reference.DisplayType), true)
			if err != nil {
				return err
			}
			if reference.DisplayType == DISPLAY_TYPE_BONE {
				err = writeNumber(fout, boneIdxType, float64(reference.DisplayIndex), 0.0, true)
			} else {
				err = writeNumber(fout, morphIdxType, float64(reference.DisplayIndex), 0.0, true)
			}
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// writeRigidBodies 剛体データの書き込み
func writeRigidBodies(fout *os.File, model *PmxModel, boneIdxType pmxBinaryType) error {
	err := writeNumber(fout, pmxBinaryType_int, float64(model.RigidBodies.Len()), 0.0, true)
	if err != nil {
		return err
	}

	for i := range model.RigidBodies.Len() {
		rigidbody := model.RigidBodies.Get(i)
		err = writeText(fout, rigidbody.Name, fmt.Sprintf("Rigidbody %d", i))
		if err != nil {
			return err
		}
		err = writeText(fout, rigidbody.EnglishName, fmt.Sprintf("Rigidbody %d", i))
		if err != nil {
			return err
		}
		err = writeNumber(fout, boneIdxType, float64(rigidbody.BoneIndex), 0.0, false)
		if err != nil {
			return err
		}
		err = writeByte(fout, int(rigidbody.CollisionGroup), true)
		if err != nil {
			return err
		}
		err = writeShort(fout, uint16(rigidbody.CollisionGroupMaskValue))
		if err != nil {
			return err
		}
		err = writeByte(fout, int(rigidbody.ShapeType), true)
		if err != nil {
			return err
		}
		err = writeNumber(fout, pmxBinaryType_float, rigidbody.Size.GetX(), 0.0, true)
		if err != nil {
			return err
		}
		err = writeNumber(fout, pmxBinaryType_float, rigidbody.Size.GetY(), 0.0, true)
		if err != nil {
			return err
		}
		err = writeNumber(fout, pmxBinaryType_float, rigidbody.Size.GetZ(), 0.0, true)
		if err != nil {
			return err
		}
		err = writeNumber(fout, pmxBinaryType_float, rigidbody.Position.GetX(), 0.0, false)
		if err != nil {
			return err
		}
		err = writeNumber(fout, pmxBinaryType_float, rigidbody.Position.GetY(), 0.0, false)
		if err != nil {
			return err
		}
		err = writeNumber(fout, pmxBinaryType_float, rigidbody.Position.GetZ(), 0.0, false)
		if err != nil {
			return err
		}
		err = writeNumber(fout, pmxBinaryType_float, rigidbody.Rotation.GetRadians().GetX(), 0.0, false)
		if err != nil {
			return err
		}
		err = writeNumber(fout, pmxBinaryType_float, rigidbody.Rotation.GetRadians().GetY(), 0.0, false)
		if err != nil {
			return err
		}
		err = writeNumber(fout, pmxBinaryType_float, rigidbody.Rotation.GetRadians().GetZ(), 0.0, false)
		if err != nil {
			return err
		}
		err = writeNumber(fout, pmxBinaryType_float, rigidbody.RigidBodyParam.Mass, 0.0, true)
		if err != nil {
			return err
		}
		err = writeNumber(fout, pmxBinaryType_float, rigidbody.RigidBodyParam.LinearDamping, 0.0, true)
		if err != nil {
			return err
		}
		err = writeNumber(fout, pmxBinaryType_float, rigidbody.RigidBodyParam.AngularDamping, 0.0, true)
		if err != nil {
			return err
		}
		err = writeNumber(fout, pmxBinaryType_float, rigidbody.RigidBodyParam.Restitution, 0.0, true)
		if err != nil {
			return err
		}
		err = writeNumber(fout, pmxBinaryType_float, rigidbody.RigidBodyParam.Friction, 0.0, true)
		if err != nil {
			return err
		}
		err = writeByte(fout, int(rigidbody.PhysicsType), true)
		if err != nil {
			return err
		}
	}

	return nil
}

// writeJoints ジョイントデータの書き込み
func writeJoints(fout *os.File, model *PmxModel, rigidbodyIdxType pmxBinaryType) error {
	err := writeNumber(fout, pmxBinaryType_int, float64(model.Joints.Len()), 0.0, true)
	if err != nil {
		return err
	}

	for i := range model.Joints.Len() {
		joint := model.Joints.Get(i)

		err = writeText(fout, joint.Name, fmt.Sprintf("Joint %d", i))
		if err != nil {
			return err
		}
		err = writeText(fout, joint.EnglishName, fmt.Sprintf("Joint %d", i))
		if err != nil {
			return err
		}
		err = writeByte(fout, int(joint.JointType), true)
		if err != nil {
			return err
		}
		err = writeNumber(fout, rigidbodyIdxType, float64(joint.RigidbodyIndexA), 0.0, false)
		if err != nil {
			return err
		}
		err = writeNumber(fout, rigidbodyIdxType, float64(joint.RigidbodyIndexB), 0.0, false)
		if err != nil {
			return err
		}
		err = writeNumber(fout, pmxBinaryType_float, joint.Position.GetX(), 0.0, false)
		if err != nil {
			return err
		}
		err = writeNumber(fout, pmxBinaryType_float, joint.Position.GetY(), 0.0, false)
		if err != nil {
			return err
		}
		err = writeNumber(fout, pmxBinaryType_float, joint.Position.GetZ(), 0.0, false)
		if err != nil {
			return err
		}
		err = writeNumber(fout, pmxBinaryType_float, joint.Rotation.GetRadians().GetX(), 0.0, false)
		if err != nil {
			return err
		}
		err = writeNumber(fout, pmxBinaryType_float, joint.Rotation.GetRadians().GetY(), 0.0, false)
		if err != nil {
			return err
		}
		err = writeNumber(fout, pmxBinaryType_float, joint.Rotation.GetRadians().GetZ(), 0.0, false)
		if err != nil {
			return err
		}
		err = writeNumber(fout, pmxBinaryType_float, joint.JointParam.TranslationLimitMin.GetX(), 0.0, false)
		if err != nil {
			return err
		}
		err = writeNumber(fout, pmxBinaryType_float, joint.JointParam.TranslationLimitMin.GetY(), 0.0, false)
		if err != nil {
			return err
		}
		err = writeNumber(fout, pmxBinaryType_float, joint.JointParam.TranslationLimitMin.GetZ(), 0.0, false)
		if err != nil {
			return err
		}
		err = writeNumber(fout, pmxBinaryType_float, joint.JointParam.TranslationLimitMax.GetX(), 0.0, false)
		if err != nil {
			return err
		}
		err = writeNumber(fout, pmxBinaryType_float, joint.JointParam.TranslationLimitMax.GetY(), 0.0, false)
		if err != nil {
			return err
		}
		err = writeNumber(fout, pmxBinaryType_float, joint.JointParam.TranslationLimitMax.GetZ(), 0.0, false)
		if err != nil {
			return err
		}
		err = writeNumber(fout, pmxBinaryType_float, joint.JointParam.RotationLimitMin.GetRadians().GetX(), 0.0, false)
		if err != nil {
			return err
		}
		err = writeNumber(fout, pmxBinaryType_float, joint.JointParam.RotationLimitMin.GetRadians().GetY(), 0.0, false)
		if err != nil {
			return err
		}
		err = writeNumber(fout, pmxBinaryType_float, joint.JointParam.RotationLimitMin.GetRadians().GetZ(), 0.0, false)
		if err != nil {
			return err
		}
		err = writeNumber(fout, pmxBinaryType_float, joint.JointParam.RotationLimitMax.GetRadians().GetX(), 0.0, false)
		if err != nil {
			return err
		}
		err = writeNumber(fout, pmxBinaryType_float, joint.JointParam.RotationLimitMax.GetRadians().GetY(), 0.0, false)
		if err != nil {
			return err
		}
		err = writeNumber(fout, pmxBinaryType_float, joint.JointParam.RotationLimitMax.GetRadians().GetZ(), 0.0, false)
		if err != nil {
			return err
		}
		err = writeNumber(fout, pmxBinaryType_float, joint.JointParam.SpringConstantTranslation.GetX(), 0.0, false)
		if err != nil {
			return err
		}
		err = writeNumber(fout, pmxBinaryType_float, joint.JointParam.SpringConstantTranslation.GetY(), 0.0, false)
		if err != nil {
			return err
		}
		err = writeNumber(fout, pmxBinaryType_float, joint.JointParam.SpringConstantTranslation.GetZ(), 0.0, false)
		if err != nil {
			return err
		}
		err = writeNumber(fout, pmxBinaryType_float, joint.JointParam.SpringConstantRotation.GetRadians().GetX(), 0.0, false)
		if err != nil {
			return err
		}
		err = writeNumber(fout, pmxBinaryType_float, joint.JointParam.SpringConstantRotation.GetRadians().GetY(), 0.0, false)
		if err != nil {
			return err
		}
		err = writeNumber(fout, pmxBinaryType_float, joint.JointParam.SpringConstantRotation.GetRadians().GetZ(), 0.0, false)
		if err != nil {
			return err
		}
	}

	return nil
}

// ------------------------------
func writeText(fout *os.File, text string, defaultText string) error {
	var binaryTxt []byte
	var err error

	// エンコードの試行
	binaryTxt, err = encodeUTF16LE(text)
	if err != nil {
		binaryTxt, _ = encodeUTF16LE(defaultText)
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

func encodeUTF16LE(s string) ([]byte, error) {
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

func defineWriteIndexForVertex(size int) (int, pmxBinaryType) {
	if size < 256 {
		return 1, pmxBinaryType_unsignedByte
	} else if size <= 65535 {
		return 2, pmxBinaryType_unsignedShort
	}
	return 4, pmxBinaryType_int
}

func defineWriteIndexForOthers(size int) (int, pmxBinaryType) {
	if size < 128 {
		return 1, pmxBinaryType_byte
	} else if size <= 32767 {
		return 2, pmxBinaryType_short
	}
	return 4, pmxBinaryType_int
}

func writeNumber(fout *os.File, valType pmxBinaryType, val float64, defaultValue float64, isPositiveOnly bool) error {
	// 値の検証と修正
	if math.IsNaN(val) || math.IsInf(val, 0) {
		val = defaultValue
	}
	if isPositiveOnly && val < 0 {
		val = 0
	}

	// バイナリデータの作成
	var buf bytes.Buffer
	var err error
	switch valType {
	case pmxBinaryType_float:
		err = binary.Write(&buf, binary.LittleEndian, float32(val))
	case pmxBinaryType_unsignedInt:
		err = binary.Write(&buf, binary.LittleEndian, uint32(val))
	case pmxBinaryType_unsignedByte:
		err = binary.Write(&buf, binary.LittleEndian, uint8(val))
	case pmxBinaryType_unsignedShort:
		err = binary.Write(&buf, binary.LittleEndian, uint16(val))
	case pmxBinaryType_byte:
		err = binary.Write(&buf, binary.LittleEndian, int8(val))
	case pmxBinaryType_short:
		err = binary.Write(&buf, binary.LittleEndian, int16(val))
	default:
		err = binary.Write(&buf, binary.LittleEndian, int32(val))
	}
	if err != nil {
		return writeDefaultNumber(fout, valType, defaultValue)
	}

	// ファイルへの書き込み
	_, err = fout.Write(buf.Bytes())
	if err != nil {
		return writeDefaultNumber(fout, valType, defaultValue)
	}
	return nil
}

func writeDefaultNumber(fout *os.File, valType pmxBinaryType, defaultValue float64) error {
	var buf bytes.Buffer
	var err error
	switch valType {
	case pmxBinaryType_float:
		err = binary.Write(&buf, binary.LittleEndian, float32(defaultValue))
	default:
		err = binary.Write(&buf, binary.LittleEndian, int32(defaultValue))
	}
	if err != nil {
		return err
	}
	_, err = fout.Write(buf.Bytes())
	return err
}

func writeByte(fout *os.File, val int, isUnsigned bool) error {
	var buf bytes.Buffer
	var err error

	if isUnsigned {
		err = binary.Write(&buf, binary.LittleEndian, uint8(val))
	} else {
		err = binary.Write(&buf, binary.LittleEndian, int8(val))
	}

	if err != nil {
		return err
	}

	_, err = fout.Write(buf.Bytes())
	return err
}

func writeShort(fout *os.File, val uint16) error {
	var buf bytes.Buffer
	err := binary.Write(&buf, binary.LittleEndian, val)
	if err != nil {
		return err
	}
	_, err = fout.Write(buf.Bytes())
	return err
}
