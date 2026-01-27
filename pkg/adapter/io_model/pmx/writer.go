// 指示: miu200521358
package pmx

import (
	"io"

	"github.com/miu200521358/mlib_go/pkg/adapter/io_common"
	"github.com/miu200521358/mlib_go/pkg/domain/mmath"
	"github.com/miu200521358/mlib_go/pkg/domain/model"
	"github.com/miu200521358/mlib_go/pkg/shared/base/logging"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/unicode"
)

// pmxWriter はPMX書き込み処理を表す。
type pmxWriter struct {
	writer *io_common.BinaryWriter
}

// pmxWriteState は書き込みに必要な状態を表す。
type pmxWriteState struct {
	writer             *io_common.BinaryWriter
	model              *model.PmxModel
	encoding           encoding.Encoding
	encodeType         byte
	extendedUVCount    int
	vertexIndexSize    byte
	textureIndexSize   byte
	materialIndexSize  byte
	boneIndexSize      byte
	morphIndexSize     byte
	rigidBodyIndexSize byte
	boneMapping        indexMapping
	morphMapping       indexMapping
	rigidMapping       indexMapping
	layerMapping       map[int]int
	includeSystem      bool
	bonesByIndex       []*model.Bone
	morphsByIndex      []*model.Morph
	rigidsByIndex      []*model.RigidBody
}

// newPmxWriter はpmxWriterを生成する。
func newPmxWriter(w io.Writer) *pmxWriter {
	return &pmxWriter{writer: io_common.NewBinaryWriter(w)}
}

// Write はPMXを書き込む。
func (p *pmxWriter) Write(modelData *model.PmxModel, opts io_common.SaveOptions) error {
	if modelData == nil {
		return io_common.NewIoEncodeFailed("PMXモデルがnilです", nil)
	}
	state, err := p.prepareState(modelData, opts)
	if err != nil {
		return err
	}
	if err := state.writeHeader(); err != nil {
		return err
	}
	if err := state.writeVertices(); err != nil {
		return err
	}
	if err := state.writeFaces(); err != nil {
		return err
	}
	if err := state.writeTextures(); err != nil {
		return err
	}
	if err := state.writeMaterials(); err != nil {
		return err
	}
	if err := state.writeBones(); err != nil {
		return err
	}
	if err := state.writeMorphs(); err != nil {
		return err
	}
	if err := state.writeDisplaySlots(); err != nil {
		return err
	}
	if err := state.writeRigidBodies(); err != nil {
		return err
	}
	if err := state.writeJoints(); err != nil {
		return err
	}
	return nil
}

// prepareState は書き込み状態を生成する。
func (p *pmxWriter) prepareState(modelData *model.PmxModel, opts io_common.SaveOptions) (*pmxWriteState, error) {
	bonesByIndex := make([]*model.Bone, modelData.Bones.Len())
	for i := 0; i < len(bonesByIndex); i++ {
		bone, err := modelData.Bones.Get(i)
		if err != nil {
			continue
		}
		bonesByIndex[i] = bone
	}
	morphsByIndex := modelData.Morphs.Values()
	rigidsByIndex := modelData.RigidBodies.Values()

	boneMapping := buildIndexMapping(len(bonesByIndex), func(index int) bool {
		if opts.IncludeSystem {
			return true
		}
		bone := bonesByIndex[index]
		return bone != nil && !bone.IsSystem
	})
	morphMapping := buildIndexMapping(len(morphsByIndex), func(index int) bool {
		if opts.IncludeSystem {
			return true
		}
		morph := morphsByIndex[index]
		return morph != nil && !morph.IsSystem
	})
	rigidMapping := buildIndexMapping(len(rigidsByIndex), func(index int) bool {
		if opts.IncludeSystem {
			return true
		}
		rigid := rigidsByIndex[index]
		return rigid != nil && !rigid.IsSystem
	})

	extendedUVCount := calcExtendedUVCount(modelData.Vertices.Values())
	if extendedUVCount > 4 {
		return nil, io_common.NewIoEncodeFailed("拡張UV数が上限を超えています", nil)
	}

	encodeType := byte(0)
	encoding := unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM)

	vertexIndexSize := defineVertexIndexSize(modelData.Vertices.Len())
	textureIndexSize := defineOtherIndexSize(modelData.Textures.Len())
	materialIndexSize := defineOtherIndexSize(modelData.Materials.Len())
	boneIndexSize := defineOtherIndexSize(len(boneMapping.newToOld))
	morphIndexSize := defineOtherIndexSize(len(morphMapping.newToOld))
	rigidIndexSize := defineOtherIndexSize(len(rigidMapping.newToOld))

	layerMapping := compressLayers(bonesByIndex, boneMapping)

	return &pmxWriteState{
		writer:             p.writer,
		model:              modelData,
		encoding:           encoding,
		encodeType:         encodeType,
		extendedUVCount:    extendedUVCount,
		vertexIndexSize:    vertexIndexSize,
		textureIndexSize:   textureIndexSize,
		materialIndexSize:  materialIndexSize,
		boneIndexSize:      boneIndexSize,
		morphIndexSize:     morphIndexSize,
		rigidBodyIndexSize: rigidIndexSize,
		boneMapping:        boneMapping,
		morphMapping:       morphMapping,
		rigidMapping:       rigidMapping,
		layerMapping:       layerMapping,
		includeSystem:      opts.IncludeSystem,
		bonesByIndex:       bonesByIndex,
		morphsByIndex:      morphsByIndex,
		rigidsByIndex:      rigidsByIndex,
	}, nil
}

// writeHeader はヘッダを書き込む。
func (s *pmxWriteState) writeHeader() error {
	if err := s.writer.WriteBytes([]byte("PMX ")); err != nil {
		return io_common.NewIoSaveFailed("PMX署名の書き込みに失敗しました", err)
	}
	if err := s.writer.WriteFloat32(2.0, 0, true); err != nil {
		return io_common.NewIoSaveFailed("PMXバージョンの書き込みに失敗しました", err)
	}
	if err := s.writer.WriteUint8(8); err != nil {
		return io_common.NewIoSaveFailed("PMXヘッダサイズの書き込みに失敗しました", err)
	}
	if err := s.writer.WriteUint8(s.encodeType); err != nil {
		return io_common.NewIoSaveFailed("PMXエンコード種別の書き込みに失敗しました", err)
	}
	if err := s.writer.WriteUint8(uint8(s.extendedUVCount)); err != nil {
		return io_common.NewIoSaveFailed("PMX拡張UV数の書き込みに失敗しました", err)
	}
	if err := s.writer.WriteUint8(s.vertexIndexSize); err != nil {
		return io_common.NewIoSaveFailed("PMX頂点インデックスサイズの書き込みに失敗しました", err)
	}
	if err := s.writer.WriteUint8(s.textureIndexSize); err != nil {
		return io_common.NewIoSaveFailed("PMXテクスチャインデックスサイズの書き込みに失敗しました", err)
	}
	if err := s.writer.WriteUint8(s.materialIndexSize); err != nil {
		return io_common.NewIoSaveFailed("PMX材質インデックスサイズの書き込みに失敗しました", err)
	}
	if err := s.writer.WriteUint8(s.boneIndexSize); err != nil {
		return io_common.NewIoSaveFailed("PMXボーンインデックスサイズの書き込みに失敗しました", err)
	}
	if err := s.writer.WriteUint8(s.morphIndexSize); err != nil {
		return io_common.NewIoSaveFailed("PMXモーフインデックスサイズの書き込みに失敗しました", err)
	}
	if err := s.writer.WriteUint8(s.rigidBodyIndexSize); err != nil {
		return io_common.NewIoSaveFailed("PMX剛体インデックスサイズの書き込みに失敗しました", err)
	}

	if err := s.writeText(s.model.Name()); err != nil {
		return err
	}
	if err := s.writeText(s.model.EnglishName); err != nil {
		return err
	}
	if err := s.writeText(s.model.Comment); err != nil {
		return err
	}
	if err := s.writeText(s.model.EnglishComment); err != nil {
		return err
	}
	return nil
}

// writeText はTextBufを書き込む。
func (s *pmxWriteState) writeText(text string) error {
	bytes, err := encodeText(s.encoding, text)
	if err != nil {
		return io_common.NewIoEncodeFailed("PMX文字列のエンコードに失敗しました", err)
	}
	if err := s.writer.WriteInt32(int32(len(bytes))); err != nil {
		return io_common.NewIoSaveFailed("PMX文字列長の書き込みに失敗しました", err)
	}
	if err := s.writer.WriteBytes(bytes); err != nil {
		return io_common.NewIoSaveFailed("PMX文字列の書き込みに失敗しました", err)
	}
	return nil
}

// writeVertices は頂点セクションを書き込む。
func (s *pmxWriteState) writeVertices() error {
	if err := s.writer.WriteInt32(int32(s.model.Vertices.Len())); err != nil {
		return io_common.NewIoSaveFailed("PMX頂点数の書き込みに失敗しました", err)
	}
	for _, vertex := range s.model.Vertices.Values() {
		if vertex == nil {
			return io_common.NewIoEncodeFailed("PMX頂点がnilです", nil)
		}
		if err := s.writeVec3(vertex.Position, false); err != nil {
			return err
		}
		if err := s.writeVec3(vertex.Normal, false); err != nil {
			return err
		}
		if err := s.writeVec2(vertex.Uv, false); err != nil {
			return err
		}
		uvs := vertex.ExtendedUvs
		for i := 0; i < s.extendedUVCount; i++ {
			var uv mmath.Vec4
			if i < len(uvs) {
				uv = uvs[i]
			}
			if err := s.writeVec4(uv, false); err != nil {
				return err
			}
		}

		deformType := vertex.Deform.DeformType()
		if err := s.writer.WriteUint8(byte(deformType)); err != nil {
			return io_common.NewIoSaveFailed("PMXデフォーム種別の書き込みに失敗しました", err)
		}
		if err := s.writeDeform(vertex); err != nil {
			return err
		}
		if err := s.writer.WriteFloat32(vertex.EdgeFactor, 0, true); err != nil {
			return io_common.NewIoSaveFailed("PMXエッジ倍率の書き込みに失敗しました", err)
		}
	}
	return nil
}

// writeDeform はデフォーム情報を書き込む。
func (s *pmxWriteState) writeDeform(vertex *model.Vertex) error {
	if vertex == nil || vertex.Deform == nil {
		return io_common.NewIoEncodeFailed("PMXデフォームが空です", nil)
	}
	switch deform := vertex.Deform.(type) {
	case *model.Bdef1:
		idx := s.boneMapping.mapIndex(deform.Indexes()[0])
		if err := writeSignedIndex(s.writer, s.boneIndexSize, idx); err != nil {
			return io_common.NewIoSaveFailed("PMXデフォーム(BDEF1)の書き込みに失敗しました", err)
		}
	case *model.Bdef2:
		idx0 := s.boneMapping.mapIndex(deform.Indexes()[0])
		idx1 := s.boneMapping.mapIndex(deform.Indexes()[1])
		weight0 := deform.Weights()[0]
		if idx0 < 0 && idx1 < 0 {
			idx0 = -1
			idx1 = -1
			weight0 = 1.0
		} else if idx0 < 0 {
			idx0 = idx1
			idx1 = -1
			weight0 = 1.0
		} else if idx1 < 0 {
			idx1 = -1
			weight0 = 1.0
		}
		if err := writeSignedIndex(s.writer, s.boneIndexSize, idx0); err != nil {
			return io_common.NewIoSaveFailed("PMXデフォーム(BDEF2)の書き込みに失敗しました", err)
		}
		if err := writeSignedIndex(s.writer, s.boneIndexSize, idx1); err != nil {
			return io_common.NewIoSaveFailed("PMXデフォーム(BDEF2)の書き込みに失敗しました", err)
		}
		if err := s.writer.WriteFloat32(weight0, 0, true); err != nil {
			return io_common.NewIoSaveFailed("PMXデフォーム(BDEF2)の書き込みに失敗しました", err)
		}
	case *model.Bdef4:
		idxs := deform.Indexes()
		weights := deform.Weights()
		mappedIdxs := make([]int, 4)
		mappedWeights := make([]float64, 4)
		var sum float64
		for i := 0; i < 4; i++ {
			mappedIdxs[i] = s.boneMapping.mapIndex(idxs[i])
			if mappedIdxs[i] >= 0 {
				mappedWeights[i] = weights[i]
				sum += weights[i]
			}
		}
		if sum == 0 {
			for i := 0; i < 4; i++ {
				mappedIdxs[i] = -1
				mappedWeights[i] = 0
			}
		} else {
			for i := 0; i < 4; i++ {
				if mappedIdxs[i] < 0 {
					mappedWeights[i] = 0
					continue
				}
				mappedWeights[i] /= sum
			}
		}
		for i := 0; i < 4; i++ {
			if err := writeSignedIndex(s.writer, s.boneIndexSize, mappedIdxs[i]); err != nil {
				return io_common.NewIoSaveFailed("PMXデフォーム(BDEF4)の書き込みに失敗しました", err)
			}
		}
		for i := 0; i < 4; i++ {
			if err := s.writer.WriteFloat32(mappedWeights[i], 0, true); err != nil {
				return io_common.NewIoSaveFailed("PMXデフォーム(BDEF4)の書き込みに失敗しました", err)
			}
		}
	case *model.Sdef:
		idx0 := s.boneMapping.mapIndex(deform.Indexes()[0])
		idx1 := s.boneMapping.mapIndex(deform.Indexes()[1])
		weight0 := deform.Weights()[0]
		if idx0 < 0 && idx1 < 0 {
			idx0 = -1
			idx1 = -1
			weight0 = 1.0
		} else if idx0 < 0 {
			idx0 = idx1
			idx1 = -1
			weight0 = 1.0
		} else if idx1 < 0 {
			idx1 = -1
			weight0 = 1.0
		}
		if err := writeSignedIndex(s.writer, s.boneIndexSize, idx0); err != nil {
			return io_common.NewIoSaveFailed("PMXデフォーム(SDEF)の書き込みに失敗しました", err)
		}
		if err := writeSignedIndex(s.writer, s.boneIndexSize, idx1); err != nil {
			return io_common.NewIoSaveFailed("PMXデフォーム(SDEF)の書き込みに失敗しました", err)
		}
		if err := s.writer.WriteFloat32(weight0, 0, true); err != nil {
			return io_common.NewIoSaveFailed("PMXデフォーム(SDEF)の書き込みに失敗しました", err)
		}
		if err := s.writeVec3(deform.SdefC, false); err != nil {
			return err
		}
		if err := s.writeVec3(deform.SdefR0, false); err != nil {
			return err
		}
		if err := s.writeVec3(deform.SdefR1, false); err != nil {
			return err
		}
	default:
		return io_common.NewIoEncodeFailed("PMXデフォーム種別が未対応です", nil)
	}
	return nil
}

// writeFaces は面セクションを書き込む。
func (s *pmxWriteState) writeFaces() error {
	faceCount := s.model.Faces.Len() * 3
	if err := s.writer.WriteInt32(int32(faceCount)); err != nil {
		return io_common.NewIoSaveFailed("PMX面数の書き込みに失敗しました", err)
	}
	for _, face := range s.model.Faces.Values() {
		if face == nil {
			return io_common.NewIoEncodeFailed("PMX面がnilです", nil)
		}
		for i := 0; i < 3; i++ {
			if err := writeVertexIndex(s.writer, s.vertexIndexSize, face.VertexIndexes[i]); err != nil {
				return io_common.NewIoSaveFailed("PMX面インデックスの書き込みに失敗しました", err)
			}
		}
	}
	return nil
}

// writeTextures はテクスチャセクションを書き込む。
func (s *pmxWriteState) writeTextures() error {
	if err := s.writer.WriteInt32(int32(s.model.Textures.Len())); err != nil {
		return io_common.NewIoSaveFailed("PMXテクスチャ数の書き込みに失敗しました", err)
	}
	for _, tex := range s.model.Textures.Values() {
		if tex == nil {
			return io_common.NewIoEncodeFailed("PMXテクスチャがnilです", nil)
		}
		if err := s.writeText(tex.Name()); err != nil {
			return err
		}
	}
	return nil
}

// writeMaterials は材質セクションを書き込む。
func (s *pmxWriteState) writeMaterials() error {
	if err := s.writer.WriteInt32(int32(s.model.Materials.Len())); err != nil {
		return io_common.NewIoSaveFailed("PMX材質数の書き込みに失敗しました", err)
	}
	for _, mat := range s.model.Materials.Values() {
		if mat == nil {
			return io_common.NewIoEncodeFailed("PMX材質がnilです", nil)
		}
		if err := s.writeText(mat.Name()); err != nil {
			return err
		}
		if err := s.writeText(mat.EnglishName); err != nil {
			return err
		}
		if err := s.writeVec4(mat.Diffuse, false); err != nil {
			return err
		}
		if err := s.writeVec3(mat.Specular.XYZ(), false); err != nil {
			return err
		}
		if err := s.writer.WriteFloat32(mat.Specular.W, 0, true); err != nil {
			return io_common.NewIoSaveFailed("PMX材質SpecularPowerの書き込みに失敗しました", err)
		}
		if err := s.writeVec3(mat.Ambient, false); err != nil {
			return err
		}
		if err := s.writer.WriteUint8(uint8(mat.DrawFlag)); err != nil {
			return io_common.NewIoSaveFailed("PMX材質描画フラグの書き込みに失敗しました", err)
		}
		if err := s.writeVec4(mat.Edge, false); err != nil {
			return err
		}
		if err := s.writer.WriteFloat32(mat.EdgeSize, 0, true); err != nil {
			return io_common.NewIoSaveFailed("PMX材質エッジサイズの書き込みに失敗しました", err)
		}
		if err := writeSignedIndex(s.writer, s.textureIndexSize, mat.TextureIndex); err != nil {
			return io_common.NewIoSaveFailed("PMX材質テクスチャ参照の書き込みに失敗しました", err)
		}
		if err := writeSignedIndex(s.writer, s.textureIndexSize, mat.SphereTextureIndex); err != nil {
			return io_common.NewIoSaveFailed("PMX材質スフィア参照の書き込みに失敗しました", err)
		}
		if err := s.writer.WriteUint8(uint8(mat.SphereMode)); err != nil {
			return io_common.NewIoSaveFailed("PMX材質スフィアモードの書き込みに失敗しました", err)
		}
		if err := s.writer.WriteUint8(uint8(mat.ToonSharingFlag)); err != nil {
			return io_common.NewIoSaveFailed("PMX材質トゥーン共有フラグの書き込みに失敗しました", err)
		}
		if mat.ToonSharingFlag == model.TOON_SHARING_INDIVIDUAL {
			if err := writeSignedIndex(s.writer, s.textureIndexSize, mat.ToonTextureIndex); err != nil {
				return io_common.NewIoSaveFailed("PMX材質トゥーン参照の書き込みに失敗しました", err)
			}
		} else {
			if err := s.writer.WriteUint8(uint8(mat.ToonTextureIndex)); err != nil {
				return io_common.NewIoSaveFailed("PMX材質共有トゥーン参照の書き込みに失敗しました", err)
			}
		}
		if err := s.writeText(mat.Memo); err != nil {
			return err
		}
		if err := s.writer.WriteInt32(int32(mat.VerticesCount)); err != nil {
			return io_common.NewIoSaveFailed("PMX材質頂点数の書き込みに失敗しました", err)
		}
	}
	return nil
}

// writeBones はボーンセクションを書き込む。
func (s *pmxWriteState) writeBones() error {
	if err := s.writer.WriteInt32(int32(len(s.boneMapping.newToOld))); err != nil {
		return io_common.NewIoSaveFailed("PMXボーン数の書き込みに失敗しました", err)
	}
	for newIndex, oldIndex := range s.boneMapping.newToOld {
		bone := s.bonesByIndex[oldIndex]
		if bone == nil {
			return io_common.NewIoEncodeFailed("PMXボーンがnilです", nil)
		}
		if err := s.writeText(bone.Name()); err != nil {
			return err
		}
		if err := s.writeText(bone.EnglishName); err != nil {
			return err
		}
		if err := s.writeVec3(bone.Position, false); err != nil {
			return err
		}
		parentIndex := s.boneMapping.mapIndex(bone.ParentIndex)
		if err := writeSignedIndex(s.writer, s.boneIndexSize, parentIndex); err != nil {
			return io_common.NewIoSaveFailed("PMXボーン親インデックスの書き込みに失敗しました", err)
		}
		layer := bone.Layer
		if mapped, ok := s.layerMapping[layer]; ok {
			layer = mapped
		}
		if err := s.writer.WriteInt32(int32(layer)); err != nil {
			return io_common.NewIoSaveFailed("PMXボーンレイヤーの書き込みに失敗しました", err)
		}

		boneFlag := bone.BoneFlag
		tailIsBone := boneFlag&model.BONE_FLAG_TAIL_IS_BONE != 0
		mappedTail := bone.TailIndex
		if tailIsBone {
			mappedTail = s.boneMapping.mapIndex(bone.TailIndex)
			if mappedTail < 0 {
				mappedTail = -1
			}
		}
		ik := s.prepareIk(bone.Ik)
		if boneFlag&model.BONE_FLAG_IS_IK != 0 && ik == nil {
			boneFlag &^= model.BONE_FLAG_IS_IK
		}

		if err := s.writer.WriteUint16(uint16(boneFlag)); err != nil {
			return io_common.NewIoSaveFailed("PMXボーンフラグの書き込みに失敗しました", err)
		}
		if tailIsBone {
			if err := writeSignedIndex(s.writer, s.boneIndexSize, mappedTail); err != nil {
				return io_common.NewIoSaveFailed("PMXボーン接続インデックスの書き込みに失敗しました", err)
			}
		} else {
			if err := s.writeVec3(bone.TailPosition, false); err != nil {
				return err
			}
		}

		if boneFlag&(model.BONE_FLAG_IS_EXTERNAL_ROTATION|model.BONE_FLAG_IS_EXTERNAL_TRANSLATION) != 0 {
			effectIndex := s.boneMapping.mapIndex(bone.EffectIndex)
			if err := writeSignedIndex(s.writer, s.boneIndexSize, effectIndex); err != nil {
				return io_common.NewIoSaveFailed("PMXボーン付与親インデックスの書き込みに失敗しました", err)
			}
			if err := s.writer.WriteFloat32(bone.EffectFactor, 0, false); err != nil {
				return io_common.NewIoSaveFailed("PMXボーン付与率の書き込みに失敗しました", err)
			}
		}

		if boneFlag&model.BONE_FLAG_HAS_FIXED_AXIS != 0 {
			if err := s.writeVec3(bone.FixedAxis, false); err != nil {
				return err
			}
		}
		if boneFlag&model.BONE_FLAG_HAS_LOCAL_AXIS != 0 {
			if err := s.writeVec3(bone.LocalAxisX, false); err != nil {
				return err
			}
			if err := s.writeVec3(bone.LocalAxisZ, false); err != nil {
				return err
			}
		}
		if boneFlag&model.BONE_FLAG_IS_EXTERNAL_PARENT_DEFORM != 0 {
			if err := s.writer.WriteInt32(int32(bone.EffectorKey)); err != nil {
				return io_common.NewIoSaveFailed("PMXボーン外部親変形キーの書き込みに失敗しました", err)
			}
		}
		if boneFlag&model.BONE_FLAG_IS_IK != 0 {
			logger := logging.DefaultLogger()
			if logger.IsVerboseEnabled(logging.VERBOSE_INDEX_IK) && ik != nil {
				logger.Verbose(
					logging.VERBOSE_INDEX_IK,
					"PMX IK単位角(出力直前): bone=%s oldIndex=%d newIndex=%d raw=%.6f deg(rad換算)=%.6f",
					bone.Name(),
					oldIndex,
					newIndex,
					ik.UnitRotation.X,
					mmath.RadToDeg(ik.UnitRotation.X),
				)
			}
			if err := s.writeIk(ik); err != nil {
				return err
			}
		}
	}
	return nil
}

// writeIk はIK情報を書き込む。
func (s *pmxWriteState) writeIk(ik *model.Ik) error {
	if ik == nil {
		return io_common.NewIoEncodeFailed("PMX IKがnilです", nil)
	}
	if err := writeSignedIndex(s.writer, s.boneIndexSize, ik.BoneIndex); err != nil {
		return io_common.NewIoSaveFailed("PMX IKターゲットの書き込みに失敗しました", err)
	}
	if err := s.writer.WriteInt32(int32(ik.LoopCount)); err != nil {
		return io_common.NewIoSaveFailed("PMX IKループ数の書き込みに失敗しました", err)
	}
	if err := s.writer.WriteFloat32(ik.UnitRotation.X, 0, true); err != nil {
		return io_common.NewIoSaveFailed("PMX IK単位回転の書き込みに失敗しました", err)
	}

	if err := s.writer.WriteInt32(int32(len(ik.Links))); err != nil {
		return io_common.NewIoSaveFailed("PMX IKリンク数の書き込みに失敗しました", err)
	}
	for _, link := range ik.Links {
		if err := writeSignedIndex(s.writer, s.boneIndexSize, link.BoneIndex); err != nil {
			return io_common.NewIoSaveFailed("PMX IKリンクの書き込みに失敗しました", err)
		}
		if err := s.writer.WriteUint8(boolToByte(link.AngleLimit)); err != nil {
			return io_common.NewIoSaveFailed("PMX IK角度制限の書き込みに失敗しました", err)
		}
		if link.AngleLimit {
			if err := s.writeVec3(link.MinAngleLimit, false); err != nil {
				return err
			}
			if err := s.writeVec3(link.MaxAngleLimit, false); err != nil {
				return err
			}
		}
	}
	return nil
}

// prepareIk はIKのインデックスを再割当して返す。
func (s *pmxWriteState) prepareIk(ik *model.Ik) *model.Ik {
	if ik == nil {
		return nil
	}
	target := s.boneMapping.mapIndex(ik.BoneIndex)
	if target < 0 {
		return nil
	}
	links := make([]model.IkLink, 0, len(ik.Links))
	for _, link := range ik.Links {
		mapped := s.boneMapping.mapIndex(link.BoneIndex)
		if mapped < 0 {
			continue
		}
		link.BoneIndex = mapped
		links = append(links, link)
	}
	if len(links) == 0 {
		return nil
	}
	copied := *ik
	copied.BoneIndex = target
	copied.Links = links
	return &copied
}

// writeMorphs はモーフセクションを書き込む。
func (s *pmxWriteState) writeMorphs() error {
	if err := s.writer.WriteInt32(int32(len(s.morphMapping.newToOld))); err != nil {
		return io_common.NewIoSaveFailed("PMXモーフ数の書き込みに失敗しました", err)
	}
	for _, oldIndex := range s.morphMapping.newToOld {
		morph := s.morphsByIndex[oldIndex]
		if morph == nil {
			return io_common.NewIoEncodeFailed("PMXモーフがnilです", nil)
		}
		if err := s.writeText(morph.Name()); err != nil {
			return err
		}
		if err := s.writeText(morph.EnglishName); err != nil {
			return err
		}
		if err := s.writer.WriteUint8(byte(morph.Panel)); err != nil {
			return io_common.NewIoSaveFailed("PMXモーフパネルの書き込みに失敗しました", err)
		}
		if err := s.writer.WriteUint8(byte(morph.MorphType)); err != nil {
			return io_common.NewIoSaveFailed("PMXモーフ種別の書き込みに失敗しました", err)
		}
		filtered := s.filterMorphOffsets(morph)
		if err := s.writer.WriteInt32(int32(len(filtered))); err != nil {
			return io_common.NewIoSaveFailed("PMXモーフオフセット数の書き込みに失敗しました", err)
		}
		for _, offset := range filtered {
			if err := s.writeMorphOffset(morph.MorphType, offset); err != nil {
				return err
			}
		}
	}
	return nil
}

// filterMorphOffsets は参照更新済みのオフセット一覧を返す。
func (s *pmxWriteState) filterMorphOffsets(morph *model.Morph) []model.IMorphOffset {
	if morph == nil {
		return nil
	}
	out := make([]model.IMorphOffset, 0, len(morph.Offsets))
	for _, offset := range morph.Offsets {
		switch o := offset.(type) {
		case *model.GroupMorphOffset:
			mapped := s.morphMapping.mapIndex(o.MorphIndex)
			if mapped < 0 {
				continue
			}
			out = append(out, &model.GroupMorphOffset{MorphIndex: mapped, MorphFactor: o.MorphFactor})
		case *model.BoneMorphOffset:
			mapped := s.boneMapping.mapIndex(o.BoneIndex)
			if mapped < 0 {
				continue
			}
			out = append(out, &model.BoneMorphOffset{BoneIndex: mapped, Position: o.Position, Rotation: o.Rotation})
		default:
			out = append(out, offset)
		}
	}
	return out
}

// writeMorphOffset はモーフオフセットを書き込む。
func (s *pmxWriteState) writeMorphOffset(morphType model.MorphType, offset model.IMorphOffset) error {
	switch morphType {
	case model.MORPH_TYPE_GROUP:
		group, ok := offset.(*model.GroupMorphOffset)
		if !ok {
			return io_common.NewIoEncodeFailed("PMXグループモーフが不正です", nil)
		}
		if err := writeSignedIndex(s.writer, s.morphIndexSize, group.MorphIndex); err != nil {
			return io_common.NewIoSaveFailed("PMXグループモーフの書き込みに失敗しました", err)
		}
		if err := s.writer.WriteFloat32(group.MorphFactor, 0, true); err != nil {
			return io_common.NewIoSaveFailed("PMXグループモーフ率の書き込みに失敗しました", err)
		}
	case model.MORPH_TYPE_VERTEX, model.MORPH_TYPE_AFTER_VERTEX:
		vertex, ok := offset.(*model.VertexMorphOffset)
		if !ok {
			return io_common.NewIoEncodeFailed("PMX頂点モーフが不正です", nil)
		}
		if err := writeVertexIndex(s.writer, s.vertexIndexSize, vertex.VertexIndex); err != nil {
			return io_common.NewIoSaveFailed("PMX頂点モーフの書き込みに失敗しました", err)
		}
		if err := s.writeVec3(vertex.Position, false); err != nil {
			return err
		}
	case model.MORPH_TYPE_BONE:
		bone, ok := offset.(*model.BoneMorphOffset)
		if !ok {
			return io_common.NewIoEncodeFailed("PMXボーンモーフが不正です", nil)
		}
		if err := writeSignedIndex(s.writer, s.boneIndexSize, bone.BoneIndex); err != nil {
			return io_common.NewIoSaveFailed("PMXボーンモーフの書き込みに失敗しました", err)
		}
		if err := s.writeVec3(bone.Position, false); err != nil {
			return err
		}
		rot := bone.Rotation.Vec4()
		if err := s.writeVec4(rot, false); err != nil {
			return err
		}
	case model.MORPH_TYPE_UV,
		model.MORPH_TYPE_EXTENDED_UV1,
		model.MORPH_TYPE_EXTENDED_UV2,
		model.MORPH_TYPE_EXTENDED_UV3,
		model.MORPH_TYPE_EXTENDED_UV4:
		uv, ok := offset.(*model.UvMorphOffset)
		if !ok {
			return io_common.NewIoEncodeFailed("PMX UVモーフが不正です", nil)
		}
		if err := writeVertexIndex(s.writer, s.vertexIndexSize, uv.VertexIndex); err != nil {
			return io_common.NewIoSaveFailed("PMX UVモーフの書き込みに失敗しました", err)
		}
		if err := s.writeVec4(uv.Uv, false); err != nil {
			return err
		}
	case model.MORPH_TYPE_MATERIAL:
		mat, ok := offset.(*model.MaterialMorphOffset)
		if !ok {
			return io_common.NewIoEncodeFailed("PMX材質モーフが不正です", nil)
		}
		if err := writeSignedIndex(s.writer, s.materialIndexSize, mat.MaterialIndex); err != nil {
			return io_common.NewIoSaveFailed("PMX材質モーフの書き込みに失敗しました", err)
		}
		if err := s.writer.WriteUint8(uint8(mat.CalcMode)); err != nil {
			return io_common.NewIoSaveFailed("PMX材質モーフ計算モードの書き込みに失敗しました", err)
		}
		if err := s.writeVec4(mat.Diffuse, false); err != nil {
			return err
		}
		if err := s.writeVec3(mat.Specular.XYZ(), false); err != nil {
			return err
		}
		if err := s.writer.WriteFloat32(mat.Specular.W, 0, true); err != nil {
			return io_common.NewIoSaveFailed("PMX材質モーフSpecularPowerの書き込みに失敗しました", err)
		}
		if err := s.writeVec3(mat.Ambient, false); err != nil {
			return err
		}
		if err := s.writeVec4(mat.Edge, false); err != nil {
			return err
		}
		if err := s.writer.WriteFloat32(mat.EdgeSize, 0, true); err != nil {
			return io_common.NewIoSaveFailed("PMX材質モーフエッジサイズの書き込みに失敗しました", err)
		}
		if err := s.writeVec4(mat.TextureFactor, false); err != nil {
			return err
		}
		if err := s.writeVec4(mat.SphereTextureFactor, false); err != nil {
			return err
		}
		if err := s.writeVec4(mat.ToonTextureFactor, false); err != nil {
			return err
		}
	default:
		return io_common.NewIoEncodeFailed("PMXモーフ種別が未対応です", nil)
	}
	return nil
}

// writeDisplaySlots は表示枠セクションを書き込む。
func (s *pmxWriteState) writeDisplaySlots() error {
	if err := s.writer.WriteInt32(int32(s.model.DisplaySlots.Len())); err != nil {
		return io_common.NewIoSaveFailed("PMX表示枠数の書き込みに失敗しました", err)
	}
	for _, slot := range s.model.DisplaySlots.Values() {
		if slot == nil {
			return io_common.NewIoEncodeFailed("PMX表示枠がnilです", nil)
		}
		if err := s.writeText(slot.Name()); err != nil {
			return err
		}
		if err := s.writeText(slot.EnglishName); err != nil {
			return err
		}
		if err := s.writer.WriteUint8(uint8(slot.SpecialFlag)); err != nil {
			return io_common.NewIoSaveFailed("PMX表示枠フラグの書き込みに失敗しました", err)
		}
		references := s.filterDisplayReferences(slot.References)
		if err := s.writer.WriteInt32(int32(len(references))); err != nil {
			return io_common.NewIoSaveFailed("PMX表示枠参照数の書き込みに失敗しました", err)
		}
		for _, ref := range references {
			if err := s.writer.WriteUint8(uint8(ref.DisplayType)); err != nil {
				return io_common.NewIoSaveFailed("PMX表示枠参照種別の書き込みに失敗しました", err)
			}
			switch ref.DisplayType {
			case model.DISPLAY_TYPE_BONE:
				if err := writeSignedIndex(s.writer, s.boneIndexSize, ref.DisplayIndex); err != nil {
					return io_common.NewIoSaveFailed("PMX表示枠ボーン参照の書き込みに失敗しました", err)
				}
			case model.DISPLAY_TYPE_MORPH:
				if err := writeSignedIndex(s.writer, s.morphIndexSize, ref.DisplayIndex); err != nil {
					return io_common.NewIoSaveFailed("PMX表示枠モーフ参照の書き込みに失敗しました", err)
				}
			default:
				return io_common.NewIoEncodeFailed("PMX表示枠参照種別が未対応です", nil)
			}
		}
	}
	return nil
}

// filterDisplayReferences は参照更新済みの表示枠参照を返す。
func (s *pmxWriteState) filterDisplayReferences(refs []model.Reference) []model.Reference {
	out := make([]model.Reference, 0, len(refs))
	for _, ref := range refs {
		switch ref.DisplayType {
		case model.DISPLAY_TYPE_BONE:
			mapped := s.boneMapping.mapIndex(ref.DisplayIndex)
			if mapped < 0 {
				continue
			}
			ref.DisplayIndex = mapped
		case model.DISPLAY_TYPE_MORPH:
			mapped := s.morphMapping.mapIndex(ref.DisplayIndex)
			if mapped < 0 {
				continue
			}
			ref.DisplayIndex = mapped
		}
		out = append(out, ref)
	}
	return out
}

// writeRigidBodies は剛体セクションを書き込む。
func (s *pmxWriteState) writeRigidBodies() error {
	if err := s.writer.WriteInt32(int32(len(s.rigidMapping.newToOld))); err != nil {
		return io_common.NewIoSaveFailed("PMX剛体数の書き込みに失敗しました", err)
	}
	for _, oldIndex := range s.rigidMapping.newToOld {
		rigid := s.rigidsByIndex[oldIndex]
		if rigid == nil {
			return io_common.NewIoEncodeFailed("PMX剛体がnilです", nil)
		}
		if err := s.writeText(rigid.Name()); err != nil {
			return err
		}
		if err := s.writeText(rigid.EnglishName); err != nil {
			return err
		}
		boneIndex := s.boneMapping.mapIndex(rigid.BoneIndex)
		if err := writeSignedIndex(s.writer, s.boneIndexSize, boneIndex); err != nil {
			return io_common.NewIoSaveFailed("PMX剛体ボーン参照の書き込みに失敗しました", err)
		}
		if err := s.writer.WriteUint8(rigid.CollisionGroup.Group); err != nil {
			return io_common.NewIoSaveFailed("PMX剛体衝突グループの書き込みに失敗しました", err)
		}
		if err := s.writer.WriteUint16(rigid.CollisionGroup.Mask); err != nil {
			return io_common.NewIoSaveFailed("PMX剛体衝突マスクの書き込みに失敗しました", err)
		}
		if err := s.writer.WriteUint8(uint8(rigid.Shape)); err != nil {
			return io_common.NewIoSaveFailed("PMX剛体形状の書き込みに失敗しました", err)
		}
		if err := s.writeVec3(rigid.Size, true); err != nil {
			return err
		}
		if err := s.writeVec3(rigid.Position, false); err != nil {
			return err
		}
		if err := s.writeVec3(rigid.Rotation, false); err != nil {
			return err
		}
		if err := s.writer.WriteFloat32(rigid.Param.Mass, 0, true); err != nil {
			return io_common.NewIoSaveFailed("PMX剛体質量の書き込みに失敗しました", err)
		}
		if err := s.writer.WriteFloat32(rigid.Param.LinearDamping, 0, true); err != nil {
			return io_common.NewIoSaveFailed("PMX剛体移動減衰の書き込みに失敗しました", err)
		}
		if err := s.writer.WriteFloat32(rigid.Param.AngularDamping, 0, true); err != nil {
			return io_common.NewIoSaveFailed("PMX剛体回転減衰の書き込みに失敗しました", err)
		}
		if err := s.writer.WriteFloat32(rigid.Param.Restitution, 0, true); err != nil {
			return io_common.NewIoSaveFailed("PMX剛体反発力の書き込みに失敗しました", err)
		}
		if err := s.writer.WriteFloat32(rigid.Param.Friction, 0, true); err != nil {
			return io_common.NewIoSaveFailed("PMX剛体摩擦の書き込みに失敗しました", err)
		}
		if err := s.writer.WriteUint8(uint8(rigid.PhysicsType)); err != nil {
			return io_common.NewIoSaveFailed("PMX剛体物理種別の書き込みに失敗しました", err)
		}
	}
	return nil
}

// writeJoints はジョイントセクションを書き込む。
func (s *pmxWriteState) writeJoints() error {
	if err := s.writer.WriteInt32(int32(s.model.Joints.Len())); err != nil {
		return io_common.NewIoSaveFailed("PMXジョイント数の書き込みに失敗しました", err)
	}
	for _, joint := range s.model.Joints.Values() {
		if joint == nil {
			return io_common.NewIoEncodeFailed("PMXジョイントがnilです", nil)
		}
		if err := s.writeText(joint.Name()); err != nil {
			return err
		}
		if err := s.writeText(joint.EnglishName); err != nil {
			return err
		}
		if err := s.writer.WriteUint8(0); err != nil {
			return io_common.NewIoSaveFailed("PMXジョイント種別の書き込みに失敗しました", err)
		}
		idxA := s.rigidMapping.mapIndex(joint.RigidBodyIndexA)
		idxB := s.rigidMapping.mapIndex(joint.RigidBodyIndexB)
		if err := writeSignedIndex(s.writer, s.rigidBodyIndexSize, idxA); err != nil {
			return io_common.NewIoSaveFailed("PMXジョイント剛体A参照の書き込みに失敗しました", err)
		}
		if err := writeSignedIndex(s.writer, s.rigidBodyIndexSize, idxB); err != nil {
			return io_common.NewIoSaveFailed("PMXジョイント剛体B参照の書き込みに失敗しました", err)
		}
		if err := s.writeVec3(joint.Param.Position, false); err != nil {
			return err
		}
		if err := s.writeVec3(joint.Param.Rotation, false); err != nil {
			return err
		}
		if err := s.writeVec3(joint.Param.TranslationLimitMin, false); err != nil {
			return err
		}
		if err := s.writeVec3(joint.Param.TranslationLimitMax, false); err != nil {
			return err
		}
		if err := s.writeVec3(joint.Param.RotationLimitMin, false); err != nil {
			return err
		}
		if err := s.writeVec3(joint.Param.RotationLimitMax, false); err != nil {
			return err
		}
		if err := s.writeVec3(joint.Param.SpringConstantTranslation, false); err != nil {
			return err
		}
		if err := s.writeVec3(joint.Param.SpringConstantRotation, false); err != nil {
			return err
		}
	}
	return nil
}

// writeVec2 はVec2を書き込む。
func (s *pmxWriteState) writeVec2(vec mmath.Vec2, positiveOnly bool) error {
	if err := s.writer.WriteFloat32(vec.X, 0, positiveOnly); err != nil {
		return io_common.NewIoSaveFailed("PMXVec2の書き込みに失敗しました", err)
	}
	if err := s.writer.WriteFloat32(vec.Y, 0, positiveOnly); err != nil {
		return io_common.NewIoSaveFailed("PMXVec2の書き込みに失敗しました", err)
	}
	return nil
}

// writeVec3 はVec3を書き込む。
func (s *pmxWriteState) writeVec3(vec mmath.Vec3, positiveOnly bool) error {
	if err := s.writer.WriteFloat32(vec.X, 0, positiveOnly); err != nil {
		return io_common.NewIoSaveFailed("PMXVec3の書き込みに失敗しました", err)
	}
	if err := s.writer.WriteFloat32(vec.Y, 0, positiveOnly); err != nil {
		return io_common.NewIoSaveFailed("PMXVec3の書き込みに失敗しました", err)
	}
	if err := s.writer.WriteFloat32(vec.Z, 0, positiveOnly); err != nil {
		return io_common.NewIoSaveFailed("PMXVec3の書き込みに失敗しました", err)
	}
	return nil
}

// writeVec4 はVec4を書き込む。
func (s *pmxWriteState) writeVec4(vec mmath.Vec4, positiveOnly bool) error {
	if err := s.writer.WriteFloat32(vec.X, 0, positiveOnly); err != nil {
		return io_common.NewIoSaveFailed("PMXVec4の書き込みに失敗しました", err)
	}
	if err := s.writer.WriteFloat32(vec.Y, 0, positiveOnly); err != nil {
		return io_common.NewIoSaveFailed("PMXVec4の書き込みに失敗しました", err)
	}
	if err := s.writer.WriteFloat32(vec.Z, 0, positiveOnly); err != nil {
		return io_common.NewIoSaveFailed("PMXVec4の書き込みに失敗しました", err)
	}
	if err := s.writer.WriteFloat32(vec.W, 0, positiveOnly); err != nil {
		return io_common.NewIoSaveFailed("PMXVec4の書き込みに失敗しました", err)
	}
	return nil
}

// encodeText は文字列を指定エンコードで変換する。
func encodeText(enc encoding.Encoding, text string) ([]byte, error) {
	if text == "" {
		return []byte{}, nil
	}
	if enc == nil {
		return nil, io_common.NewIoEncodeFailed("エンコード設定が未指定です", nil)
	}
	return enc.NewEncoder().Bytes([]byte(text))
}

// calcExtendedUVCount は拡張UV数を算出する。
func calcExtendedUVCount(vertices []*model.Vertex) int {
	max := 0
	for _, v := range vertices {
		if v == nil {
			continue
		}
		if len(v.ExtendedUvs) > max {
			max = len(v.ExtendedUvs)
		}
	}
	return max
}

// boolToByte はboolを0/1に変換する。
func boolToByte(value bool) uint8 {
	if value {
		return 1
	}
	return 0
}
