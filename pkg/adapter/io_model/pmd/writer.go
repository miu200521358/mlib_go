// 指示: miu200521358
package pmd

import (
	"io"
	"math"
	"sort"

	"github.com/miu200521358/mlib_go/pkg/adapter/io_common"
	"github.com/miu200521358/mlib_go/pkg/domain/mmath"
	"github.com/miu200521358/mlib_go/pkg/domain/model"
)

const maxUint16 = int(^uint16(0))

// pmdWriter はPMD書き込み処理を表す。
type pmdWriter struct {
	writer *io_common.BinaryWriter
}

type pmdWriteState struct {
	writer             *io_common.BinaryWriter
	model              *model.PmxModel
	includeSystem      bool
	boneMapping        indexMapping
	morphMapping       indexMapping
	rigidMapping       indexMapping
	bonesByIndex       []*model.Bone
	morphsByIndex      []*model.Morph
	rigidsByIndex      []*model.RigidBody
	supportedMorphs    []*model.Morph
	pmdMorphIndexMap   map[int]int
	baseVertexIndexes  []int
	baseVertexIndexMap map[int]int
}

// newPmdWriter はpmdWriterを生成する。
func newPmdWriter(w io.Writer) *pmdWriter {
	return &pmdWriter{writer: io_common.NewBinaryWriter(w)}
}

// Write はPMDを書き込む。
func (p *pmdWriter) Write(modelData *model.PmxModel, opts io_common.SaveOptions) error {
	if modelData == nil {
		return io_common.NewIoEncodeFailed("PMDモデルがnilです", nil)
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
	if err := state.writeMaterials(); err != nil {
		return err
	}
	if err := state.writeBones(); err != nil {
		return err
	}
	if err := state.writeIk(); err != nil {
		return err
	}
	if err := state.writeSkins(); err != nil {
		return err
	}
	if err := state.writeSkinDisplayList(); err != nil {
		return err
	}
	boneSlotNames, boneSlotEnglish, boneSlotIndexes := state.collectBoneDisplaySlots()
	if err := state.writeBoneDisplayNames(boneSlotNames); err != nil {
		return err
	}
	if err := state.writeBoneDisplayList(boneSlotIndexes); err != nil {
		return err
	}
	if err := state.writeExtensions(boneSlotNames, boneSlotEnglish); err != nil {
		return err
	}
	return nil
}

func (p *pmdWriter) prepareState(modelData *model.PmxModel, opts io_common.SaveOptions) (*pmdWriteState, error) {
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

	for _, vertex := range modelData.Vertices.Values() {
		if vertex == nil {
			continue
		}
		if len(vertex.ExtendedUvs) > 0 {
			return nil, io_common.NewIoEncodeFailed("PMDは拡張UVに対応していません", nil)
		}
		if vertex.Deform == nil {
			return nil, io_common.NewIoEncodeFailed("PMDデフォームが空です", nil)
		}
		switch vertex.Deform.DeformType() {
		case model.BDEF1, model.BDEF2:
			// supported
		default:
			return nil, io_common.NewIoEncodeFailed("PMDはBDEF1/BDEF2のみ対応です", nil)
		}
	}

	for _, morph := range morphsByIndex {
		if morph == nil {
			continue
		}
		if !opts.IncludeSystem && morph.IsSystem {
			continue
		}
		switch morph.MorphType {
		case model.MORPH_TYPE_VERTEX, model.MORPH_TYPE_AFTER_VERTEX:
			// supported
		default:
			return nil, io_common.NewIoEncodeFailed("PMDは頂点モーフのみ対応です", nil)
		}
		for _, offset := range morph.Offsets {
			if _, ok := offset.(*model.VertexMorphOffset); !ok {
				return nil, io_common.NewIoEncodeFailed("PMDは頂点モーフのみ対応です", nil)
			}
		}
	}

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

	supportedMorphs := make([]*model.Morph, 0, len(morphMapping.newToOld))
	pmdMorphIndexMap := make(map[int]int)
	for i, oldIndex := range morphMapping.newToOld {
		morph := morphsByIndex[oldIndex]
		if morph == nil {
			continue
		}
		supportedMorphs = append(supportedMorphs, morph)
		pmdMorphIndexMap[oldIndex] = i + 1
	}

	baseVertexSet := make(map[int]struct{})
	for _, morph := range supportedMorphs {
		if morph == nil {
			continue
		}
		for _, offset := range morph.Offsets {
			vtx, ok := offset.(*model.VertexMorphOffset)
			if !ok {
				continue
			}
			baseVertexSet[vtx.VertexIndex] = struct{}{}
		}
	}
	baseVertexIndexes := make([]int, 0, len(baseVertexSet))
	for index := range baseVertexSet {
		baseVertexIndexes = append(baseVertexIndexes, index)
	}
	sort.Ints(baseVertexIndexes)
	baseVertexIndexMap := make(map[int]int, len(baseVertexIndexes))
	for i, idx := range baseVertexIndexes {
		baseVertexIndexMap[idx] = i
	}

	return &pmdWriteState{
		writer:             p.writer,
		model:              modelData,
		includeSystem:      opts.IncludeSystem,
		boneMapping:        boneMapping,
		morphMapping:       morphMapping,
		rigidMapping:       rigidMapping,
		bonesByIndex:       bonesByIndex,
		morphsByIndex:      morphsByIndex,
		rigidsByIndex:      rigidsByIndex,
		supportedMorphs:    supportedMorphs,
		pmdMorphIndexMap:   pmdMorphIndexMap,
		baseVertexIndexes:  baseVertexIndexes,
		baseVertexIndexMap: baseVertexIndexMap,
	}, nil
}

func (s *pmdWriteState) writeHeader() error {
	if err := s.writer.WriteBytes([]byte("Pmd")); err != nil {
		return io_common.NewIoSaveFailed("PMD署名の書き込みに失敗しました", err)
	}
	if err := s.writer.WriteFloat32(1.0, 0, true); err != nil {
		return io_common.NewIoSaveFailed("PMDバージョンの書き込みに失敗しました", err)
	}
	if err := s.writeFixedString(s.model.Name(), 20, "PMDモデル名"); err != nil {
		return err
	}
	if err := s.writeFixedString(s.model.Comment, 256, "PMDコメント"); err != nil {
		return err
	}
	return nil
}

func (s *pmdWriteState) writeVertices() error {
	vertexCount := s.model.Vertices.Len()
	if vertexCount > maxUint16 {
		return io_common.NewIoEncodeFailed("PMD頂点数が上限を超えています", nil)
	}
	if err := s.writer.WriteUint32(uint32(vertexCount)); err != nil {
		return io_common.NewIoSaveFailed("PMD頂点数の書き込みに失敗しました", err)
	}
	for _, vertex := range s.model.Vertices.Values() {
		if vertex == nil {
			continue
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
		bone0, bone1, weight := s.vertexBoneWeights(vertex)
		if err := s.writer.WriteUint16(uint16(bone0)); err != nil {
			return io_common.NewIoSaveFailed("PMD頂点ボーン番号の書き込みに失敗しました", err)
		}
		if err := s.writer.WriteUint16(uint16(bone1)); err != nil {
			return io_common.NewIoSaveFailed("PMD頂点ボーン番号の書き込みに失敗しました", err)
		}
		if err := s.writer.WriteUint8(uint8(weight)); err != nil {
			return io_common.NewIoSaveFailed("PMD頂点ボーンウェイトの書き込みに失敗しました", err)
		}
		edgeFlag := uint8(0)
		if vertex.EdgeFactor <= 0 {
			edgeFlag = 1
		}
		if err := s.writer.WriteUint8(edgeFlag); err != nil {
			return io_common.NewIoSaveFailed("PMD頂点エッジフラグの書き込みに失敗しました", err)
		}
	}
	return nil
}

func (s *pmdWriteState) writeFaces() error {
	faceCount := s.model.Faces.Len()
	faceVertCount := faceCount * 3
	if err := s.writer.WriteUint32(uint32(faceVertCount)); err != nil {
		return io_common.NewIoSaveFailed("PMD面頂点数の書き込みに失敗しました", err)
	}
	for _, face := range s.model.Faces.Values() {
		if face == nil {
			continue
		}
		for _, idx := range face.VertexIndexes {
			if idx < 0 || idx > maxUint16 {
				return io_common.NewIoEncodeFailed("PMD面頂点番号が上限を超えています", nil)
			}
			if err := s.writer.WriteUint16(uint16(idx)); err != nil {
				return io_common.NewIoSaveFailed("PMD面頂点の書き込みに失敗しました", err)
			}
		}
	}
	return nil
}

func (s *pmdWriteState) writeMaterials() error {
	materials := s.model.Materials.Values()
	if err := s.writer.WriteUint32(uint32(len(materials))); err != nil {
		return io_common.NewIoSaveFailed("PMD材質数の書き込みに失敗しました", err)
	}
	for _, material := range materials {
		if material == nil {
			continue
		}
		diffuse := material.Diffuse
		if err := s.writeVec3(diffuse.XYZ(), false); err != nil {
			return err
		}
		if err := s.writer.WriteFloat32(diffuse.W, 1.0, true); err != nil {
			return io_common.NewIoSaveFailed("PMD材質Alphaの書き込みに失敗しました", err)
		}
		if err := s.writer.WriteFloat32(material.Specular.W, 0, true); err != nil {
			return io_common.NewIoSaveFailed("PMD材質Specularityの書き込みに失敗しました", err)
		}
		if err := s.writeVec3(material.Specular.XYZ(), false); err != nil {
			return err
		}
		if err := s.writeVec3(material.Ambient, false); err != nil {
			return err
		}
		toonIndex := uint8(0xFF)
		if material.ToonSharingFlag == model.TOON_SHARING_SHARING && material.ToonTextureIndex >= 0 && material.ToonTextureIndex <= 9 {
			toonIndex = uint8(material.ToonTextureIndex)
		}
		if err := s.writer.WriteUint8(toonIndex); err != nil {
			return io_common.NewIoSaveFailed("PMD材質Toonの書き込みに失敗しました", err)
		}
		edgeFlag := uint8(0)
		if material.DrawFlag&model.DRAW_FLAG_DRAWING_EDGE != 0 {
			edgeFlag = 1
		}
		if err := s.writer.WriteUint8(edgeFlag); err != nil {
			return io_common.NewIoSaveFailed("PMD材質エッジフラグの書き込みに失敗しました", err)
		}
		if err := s.writer.WriteUint32(uint32(material.VerticesCount)); err != nil {
			return io_common.NewIoSaveFailed("PMD材質面頂点数の書き込みに失敗しました", err)
		}
		textureName := s.textureName(material.TextureIndex)
		sphereName := s.textureName(material.SphereTextureIndex)
		combined := buildTextureSpec(textureName, sphereName, material.SphereMode)
		if err := s.writeFixedString(combined, 20, "PMD材質テクスチャ名"); err != nil {
			return err
		}
	}
	return nil
}

func (s *pmdWriteState) writeBones() error {
	boneCount := len(s.boneMapping.newToOld)
	if boneCount > maxUint16 {
		return io_common.NewIoEncodeFailed("PMDボーン数が上限を超えています", nil)
	}
	if err := s.writer.WriteUint16(uint16(boneCount)); err != nil {
		return io_common.NewIoSaveFailed("PMDボーン数の書き込みに失敗しました", err)
	}
	for _, oldIndex := range s.boneMapping.newToOld {
		bone := s.bonesByIndex[oldIndex]
		if bone == nil {
			continue
		}
		if err := s.writeFixedString(bone.Name(), 20, "PMDボーン名"); err != nil {
			return err
		}
		parentIndex := s.boneMapping.mapIndex(bone.ParentIndex)
		parentValue := uint16(0xFFFF)
		if parentIndex >= 0 {
			parentValue = uint16(parentIndex)
		}
		if err := s.writer.WriteUint16(parentValue); err != nil {
			return io_common.NewIoSaveFailed("PMDボーン親番号の書き込みに失敗しました", err)
		}
		tailValue := uint16(0xFFFF)
		if bone.BoneFlag&model.BONE_FLAG_TAIL_IS_BONE != 0 && bone.TailIndex >= 0 {
			mapped := s.boneMapping.mapIndex(bone.TailIndex)
			if mapped >= 0 {
				tailValue = uint16(mapped)
			}
		}
		if err := s.writer.WriteUint16(tailValue); err != nil {
			return io_common.NewIoSaveFailed("PMDボーン接続番号の書き込みに失敗しました", err)
		}
		boneType := boneTypeFromBone(bone)
		if err := s.writer.WriteUint8(boneType); err != nil {
			return io_common.NewIoSaveFailed("PMDボーン種別の書き込みに失敗しました", err)
		}
		ikParentValue := uint16(0)
		if boneType == 4 || boneType == 5 {
			mapped := s.boneMapping.mapIndex(bone.EffectIndex)
			if mapped >= 0 {
				ikParentValue = uint16(mapped)
			}
		}
		if err := s.writer.WriteUint16(ikParentValue); err != nil {
			return io_common.NewIoSaveFailed("PMDボーンIK親の書き込みに失敗しました", err)
		}
		if err := s.writeVec3(bone.Position, false); err != nil {
			return err
		}
	}
	return nil
}

func (s *pmdWriteState) writeIk() error {
	ikBones := make([]*model.Bone, 0)
	ikBoneIndexes := make([]int, 0)
	for _, oldIndex := range s.boneMapping.newToOld {
		bone := s.bonesByIndex[oldIndex]
		if bone == nil || bone.Ik == nil {
			continue
		}
		newIndex := s.boneMapping.mapIndex(oldIndex)
		if newIndex < 0 {
			continue
		}
		ikBones = append(ikBones, bone)
		ikBoneIndexes = append(ikBoneIndexes, newIndex)
	}
	if err := s.writer.WriteUint16(uint16(len(ikBones))); err != nil {
		return io_common.NewIoSaveFailed("PMD IK数の書き込みに失敗しました", err)
	}
	for i, bone := range ikBones {
		ik := bone.Ik
		if ik == nil {
			continue
		}
		if len(ik.Links) > math.MaxUint8 {
			return io_common.NewIoEncodeFailed("PMD IKリンク数が上限を超えています", nil)
		}
		if ik.LoopCount > maxUint16 {
			return io_common.NewIoEncodeFailed("PMD IK反復回数が上限を超えています", nil)
		}
		target := s.boneMapping.mapIndex(ik.BoneIndex)
		if target < 0 {
			return io_common.NewIoEncodeFailed("PMD IKターゲットの参照が不正です", nil)
		}
		if err := s.writer.WriteUint16(uint16(ikBoneIndexes[i])); err != nil {
			return io_common.NewIoSaveFailed("PMD IKボーン番号の書き込みに失敗しました", err)
		}
		if err := s.writer.WriteUint16(uint16(target)); err != nil {
			return io_common.NewIoSaveFailed("PMD IKターゲット番号の書き込みに失敗しました", err)
		}
		if err := s.writer.WriteUint8(uint8(len(ik.Links))); err != nil {
			return io_common.NewIoSaveFailed("PMD IKチェーン長の書き込みに失敗しました", err)
		}
		if err := s.writer.WriteUint16(uint16(ik.LoopCount)); err != nil {
			return io_common.NewIoSaveFailed("PMD IK反復回数の書き込みに失敗しました", err)
		}
		if err := s.writer.WriteFloat32(ik.UnitRotation.X, 0, true); err != nil {
			return io_common.NewIoSaveFailed("PMD IK制限角の書き込みに失敗しました", err)
		}
		for _, link := range ik.Links {
			mapped := s.boneMapping.mapIndex(link.BoneIndex)
			if mapped < 0 {
				return io_common.NewIoEncodeFailed("PMD IKリンクの参照が不正です", nil)
			}
			if err := s.writer.WriteUint16(uint16(mapped)); err != nil {
				return io_common.NewIoSaveFailed("PMD IKリンクの書き込みに失敗しました", err)
			}
		}
	}
	return nil
}

func (s *pmdWriteState) writeSkins() error {
	if len(s.supportedMorphs) == 0 {
		if err := s.writer.WriteUint16(0); err != nil {
			return io_common.NewIoSaveFailed("PMD表情数の書き込みに失敗しました", err)
		}
		return nil
	}
	count := 1 + len(s.supportedMorphs)
	if count > maxUint16 {
		return io_common.NewIoEncodeFailed("PMD表情数が上限を超えています", nil)
	}
	if err := s.writer.WriteUint16(uint16(count)); err != nil {
		return io_common.NewIoSaveFailed("PMD表情数の書き込みに失敗しました", err)
	}
	if err := s.writeFixedString("base", 20, "PMD表情名"); err != nil {
		return err
	}
	if err := s.writer.WriteUint32(uint32(len(s.baseVertexIndexes))); err != nil {
		return io_common.NewIoSaveFailed("PMD表情頂点数の書き込みに失敗しました", err)
	}
	if err := s.writer.WriteUint8(0); err != nil {
		return io_common.NewIoSaveFailed("PMD表情種別の書き込みに失敗しました", err)
	}
	for _, vertexIndex := range s.baseVertexIndexes {
		if vertexIndex < 0 || vertexIndex > math.MaxUint32 {
			return io_common.NewIoEncodeFailed("PMD表情頂点番号が不正です", nil)
		}
		vertex, err := s.model.Vertices.Get(vertexIndex)
		if err != nil || vertex == nil {
			return io_common.NewIoEncodeFailed("PMD表情頂点が存在しません", nil)
		}
		if err := s.writer.WriteUint32(uint32(vertexIndex)); err != nil {
			return io_common.NewIoSaveFailed("PMD表情頂点の書き込みに失敗しました", err)
		}
		if err := s.writeVec3(vertex.Position, false); err != nil {
			return err
		}
	}
	for _, morph := range s.supportedMorphs {
		if morph == nil {
			continue
		}
		if err := s.writeFixedString(morph.Name(), 20, "PMD表情名"); err != nil {
			return err
		}
		if err := s.writer.WriteUint32(uint32(len(morph.Offsets))); err != nil {
			return io_common.NewIoSaveFailed("PMD表情頂点数の書き込みに失敗しました", err)
		}
		if err := s.writer.WriteUint8(skinTypeFromPanel(morph.Panel)); err != nil {
			return io_common.NewIoSaveFailed("PMD表情種別の書き込みに失敗しました", err)
		}
		for _, offset := range morph.Offsets {
			vtx, ok := offset.(*model.VertexMorphOffset)
			if !ok {
				return io_common.NewIoEncodeFailed("PMD表情オフセットが不正です", nil)
			}
			baseIndex, ok := s.baseVertexIndexMap[vtx.VertexIndex]
			if !ok {
				return io_common.NewIoEncodeFailed("PMD表情頂点の参照が不正です", nil)
			}
			if err := s.writer.WriteUint32(uint32(baseIndex)); err != nil {
				return io_common.NewIoSaveFailed("PMD表情頂点の書き込みに失敗しました", err)
			}
			if err := s.writeVec3(vtx.Position, false); err != nil {
				return err
			}
		}
	}
	return nil
}

func (s *pmdWriteState) writeSkinDisplayList() error {
	if len(s.supportedMorphs) == 0 {
		if err := s.writer.WriteUint8(0); err != nil {
			return io_common.NewIoSaveFailed("PMD表情表示数の書き込みに失敗しました", err)
		}
		return nil
	}
	indices := s.collectSkinDisplayList()
	if len(indices) > math.MaxUint8 {
		return io_common.NewIoEncodeFailed("PMD表情表示数が上限を超えています", nil)
	}
	if err := s.writer.WriteUint8(uint8(len(indices))); err != nil {
		return io_common.NewIoSaveFailed("PMD表情表示数の書き込みに失敗しました", err)
	}
	for _, idx := range indices {
		if idx <= 0 || idx > maxUint16 {
			continue
		}
		if err := s.writer.WriteUint16(uint16(idx)); err != nil {
			return io_common.NewIoSaveFailed("PMD表情表示の書き込みに失敗しました", err)
		}
	}
	return nil
}

func (s *pmdWriteState) writeBoneDisplayNames(names []string) error {
	if len(names) > math.MaxUint8 {
		return io_common.NewIoEncodeFailed("PMDボーン枠名数が上限を超えています", nil)
	}
	if err := s.writer.WriteUint8(uint8(len(names))); err != nil {
		return io_common.NewIoSaveFailed("PMDボーン枠名数の書き込みに失敗しました", err)
	}
	for _, name := range names {
		if err := s.writeFixedString(name, 50, "PMDボーン枠名"); err != nil {
			return err
		}
	}
	return nil
}

func (s *pmdWriteState) writeBoneDisplayList(boneSlotIndexes []int) error {
	entries := make([]pmdBoneDisplay, 0)
	rootSlot := s.findRootDisplaySlot()
	if rootSlot != nil {
		for _, ref := range rootSlot.References {
			if ref.DisplayType != model.DISPLAY_TYPE_BONE {
				continue
			}
			mapped := s.boneMapping.mapIndex(ref.DisplayIndex)
			if mapped < 0 {
				continue
			}
			entries = append(entries, pmdBoneDisplay{boneIndex: mapped, frameIndex: 0})
		}
	}
	for i, slotIndex := range boneSlotIndexes {
		slot, err := s.model.DisplaySlots.Get(slotIndex)
		if err != nil || slot == nil {
			continue
		}
		frameIndex := i + 1
		for _, ref := range slot.References {
			if ref.DisplayType != model.DISPLAY_TYPE_BONE {
				continue
			}
			mapped := s.boneMapping.mapIndex(ref.DisplayIndex)
			if mapped < 0 {
				continue
			}
			entries = append(entries, pmdBoneDisplay{boneIndex: mapped, frameIndex: frameIndex})
		}
	}
	if err := s.writer.WriteUint32(uint32(len(entries))); err != nil {
		return io_common.NewIoSaveFailed("PMDボーン表示数の書き込みに失敗しました", err)
	}
	for _, entry := range entries {
		if entry.boneIndex < 0 || entry.boneIndex > maxUint16 {
			return io_common.NewIoEncodeFailed("PMDボーン表示の参照が不正です", nil)
		}
		if err := s.writer.WriteUint16(uint16(entry.boneIndex)); err != nil {
			return io_common.NewIoSaveFailed("PMDボーン表示の書き込みに失敗しました", err)
		}
		if err := s.writer.WriteUint8(uint8(entry.frameIndex)); err != nil {
			return io_common.NewIoSaveFailed("PMDボーン表示の書き込みに失敗しました", err)
		}
	}
	return nil
}

func (s *pmdWriteState) writeExtensions(boneSlotNames, boneSlotEnglish []string) error {
	englishEnabled := s.hasEnglishNames(boneSlotEnglish)
	flag := uint8(0)
	if englishEnabled {
		flag = 1
	}
	if err := s.writer.WriteUint8(flag); err != nil {
		return io_common.NewIoSaveFailed("PMD英名対応フラグの書き込みに失敗しました", err)
	}
	if englishEnabled {
		if err := s.writeFixedString(s.model.EnglishName, 20, "PMD英名モデル名"); err != nil {
			return err
		}
		if err := s.writeFixedString(s.model.EnglishComment, 256, "PMD英名コメント"); err != nil {
			return err
		}
		for _, oldIndex := range s.boneMapping.newToOld {
			bone := s.bonesByIndex[oldIndex]
			name := ""
			if bone != nil {
				name = bone.EnglishName
				if name == "" {
					name = bone.Name()
				}
			}
			if err := s.writeFixedString(name, 20, "PMD英名ボーン名"); err != nil {
				return err
			}
		}
		for _, morph := range s.supportedMorphs {
			name := ""
			if morph != nil {
				name = morph.EnglishName
				if name == "" {
					name = morph.Name()
				}
			}
			if err := s.writeFixedString(name, 20, "PMD英名表情名"); err != nil {
				return err
			}
		}
		for i, name := range boneSlotNames {
			english := name
			if i < len(boneSlotEnglish) && boneSlotEnglish[i] != "" {
				english = boneSlotEnglish[i]
			}
			if err := s.writeFixedString(english, 50, "PMD英名ボーン枠名"); err != nil {
				return err
			}
		}
	}

	toonNames := defaultToonFileNames()
	for _, name := range toonNames {
		if err := s.writeFixedString(name, 100, "PMDトゥーンテクスチャ名"); err != nil {
			return err
		}
	}

	if err := s.writeRigidBodies(); err != nil {
		return err
	}
	if err := s.writeJoints(); err != nil {
		return err
	}
	return nil
}

func (s *pmdWriteState) writeRigidBodies() error {
	count := len(s.rigidMapping.newToOld)
	if err := s.writer.WriteUint32(uint32(count)); err != nil {
		return io_common.NewIoSaveFailed("PMD剛体数の書き込みに失敗しました", err)
	}
	for _, oldIndex := range s.rigidMapping.newToOld {
		rigid := s.rigidsByIndex[oldIndex]
		if rigid == nil {
			continue
		}
		if err := s.writeFixedString(rigid.Name(), 20, "PMD剛体名"); err != nil {
			return err
		}
		boneIndex := s.boneMapping.mapIndex(rigid.BoneIndex)
		boneValue := int16(-1)
		if boneIndex >= 0 {
			boneValue = int16(boneIndex)
		}
		if err := s.writer.WriteInt16(boneValue); err != nil {
			return io_common.NewIoSaveFailed("PMD剛体ボーン番号の書き込みに失敗しました", err)
		}
		if err := s.writer.WriteUint8(rigid.CollisionGroup.Group); err != nil {
			return io_common.NewIoSaveFailed("PMD剛体グループの書き込みに失敗しました", err)
		}
		if err := s.writer.WriteInt16(int16(rigid.CollisionGroup.Mask)); err != nil {
			return io_common.NewIoSaveFailed("PMD剛体グループ対象の書き込みに失敗しました", err)
		}
		if err := s.writer.WriteUint8(uint8(rigid.Shape)); err != nil {
			return io_common.NewIoSaveFailed("PMD剛体形状の書き込みに失敗しました", err)
		}
		if err := s.writeVec3(rigid.Size, true); err != nil {
			return err
		}
		rigidPos := rigid.Position
		if boneIndex >= 0 {
			bone, boneErr := s.model.Bones.Get(rigid.BoneIndex)
			if boneErr == nil {
				// PMDの剛体位置はボーン相対で保存する。
				rigidPos = rigid.Position.Subed(bone.Position)
			}
		}
		if err := s.writeVec3(rigidPos, false); err != nil {
			return err
		}
		if err := s.writeVec3(rigid.Rotation, false); err != nil {
			return err
		}
		if err := s.writer.WriteFloat32(rigid.Param.Mass, 0, true); err != nil {
			return io_common.NewIoSaveFailed("PMD剛体質量の書き込みに失敗しました", err)
		}
		if err := s.writer.WriteFloat32(rigid.Param.LinearDamping, 0, true); err != nil {
			return io_common.NewIoSaveFailed("PMD剛体移動減衰の書き込みに失敗しました", err)
		}
		if err := s.writer.WriteFloat32(rigid.Param.AngularDamping, 0, true); err != nil {
			return io_common.NewIoSaveFailed("PMD剛体回転減衰の書き込みに失敗しました", err)
		}
		if err := s.writer.WriteFloat32(rigid.Param.Restitution, 0, true); err != nil {
			return io_common.NewIoSaveFailed("PMD剛体反発の書き込みに失敗しました", err)
		}
		if err := s.writer.WriteFloat32(rigid.Param.Friction, 0, true); err != nil {
			return io_common.NewIoSaveFailed("PMD剛体摩擦の書き込みに失敗しました", err)
		}
		if err := s.writer.WriteUint8(uint8(rigid.PhysicsType)); err != nil {
			return io_common.NewIoSaveFailed("PMD剛体タイプの書き込みに失敗しました", err)
		}
	}
	return nil
}

func (s *pmdWriteState) writeJoints() error {
	joints := make([]*model.Joint, 0)
	for _, joint := range s.model.Joints.Values() {
		if joint == nil {
			continue
		}
		if s.rigidMapping.mapIndex(joint.RigidBodyIndexA) < 0 || s.rigidMapping.mapIndex(joint.RigidBodyIndexB) < 0 {
			continue
		}
		joints = append(joints, joint)
	}
	if err := s.writer.WriteUint32(uint32(len(joints))); err != nil {
		return io_common.NewIoSaveFailed("PMDジョイント数の書き込みに失敗しました", err)
	}
	for _, joint := range joints {
		if joint == nil {
			continue
		}
		if err := s.writeFixedString(joint.Name(), 20, "PMDジョイント名"); err != nil {
			return err
		}
		mappedA := s.rigidMapping.mapIndex(joint.RigidBodyIndexA)
		mappedB := s.rigidMapping.mapIndex(joint.RigidBodyIndexB)
		if mappedA < 0 || mappedB < 0 {
			return io_common.NewIoEncodeFailed("PMDジョイント剛体参照が不正です", nil)
		}
		if err := s.writer.WriteInt32(int32(mappedA)); err != nil {
			return io_common.NewIoSaveFailed("PMDジョイント剛体Aの書き込みに失敗しました", err)
		}
		if err := s.writer.WriteInt32(int32(mappedB)); err != nil {
			return io_common.NewIoSaveFailed("PMDジョイント剛体Bの書き込みに失敗しました", err)
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

func (s *pmdWriteState) collectSkinDisplayList() []int {
	morphSlot := s.findMorphDisplaySlot()
	indices := make([]int, 0)
	if morphSlot != nil {
		for _, ref := range morphSlot.References {
			if ref.DisplayType != model.DISPLAY_TYPE_MORPH {
				continue
			}
			if idx, ok := s.pmdMorphIndexMap[ref.DisplayIndex]; ok {
				indices = append(indices, idx)
			}
		}
	}
	if len(indices) == 0 {
		for i := range s.supportedMorphs {
			indices = append(indices, i+1)
		}
	}
	return indices
}

func (s *pmdWriteState) collectBoneDisplaySlots() ([]string, []string, []int) {
	names := make([]string, 0)
	english := make([]string, 0)
	indexes := make([]int, 0)
	morphSlot := s.findMorphDisplaySlot()
	rootSlot := s.findRootDisplaySlot()
	for i, slot := range s.model.DisplaySlots.Values() {
		if slot == nil {
			continue
		}
		if slot == rootSlot || slot == morphSlot {
			continue
		}
		hasBone := false
		for _, ref := range slot.References {
			if ref.DisplayType == model.DISPLAY_TYPE_BONE {
				hasBone = true
				break
			}
		}
		if !hasBone {
			continue
		}
		names = append(names, slot.Name())
		english = append(english, slot.EnglishName)
		indexes = append(indexes, i)
	}
	return names, english, indexes
}

func (s *pmdWriteState) findMorphDisplaySlot() *model.DisplaySlot {
	for _, slot := range s.model.DisplaySlots.Values() {
		if slot == nil {
			continue
		}
		if slot.SpecialFlag == model.SPECIAL_FLAG_ON && (slot.Name() == "表情" || slot.EnglishName == "Exp") {
			return slot
		}
	}
	for _, slot := range s.model.DisplaySlots.Values() {
		if slot == nil {
			continue
		}
		for _, ref := range slot.References {
			if ref.DisplayType == model.DISPLAY_TYPE_MORPH {
				return slot
			}
		}
	}
	return nil
}

func (s *pmdWriteState) findRootDisplaySlot() *model.DisplaySlot {
	for _, slot := range s.model.DisplaySlots.Values() {
		if slot == nil {
			continue
		}
		if slot.SpecialFlag == model.SPECIAL_FLAG_ON && (slot.Name() == "Root" || slot.EnglishName == "Root") {
			return slot
		}
	}
	if s.model.DisplaySlots.Len() > 0 {
		slot, _ := s.model.DisplaySlots.Get(0)
		return slot
	}
	return nil
}

func (s *pmdWriteState) hasEnglishNames(boneSlotEnglish []string) bool {
	if s.model.EnglishName != "" || s.model.EnglishComment != "" {
		return true
	}
	for _, bone := range s.bonesByIndex {
		if bone != nil && bone.EnglishName != "" {
			return true
		}
	}
	for _, morph := range s.supportedMorphs {
		if morph != nil && morph.EnglishName != "" {
			return true
		}
	}
	for _, name := range boneSlotEnglish {
		if name != "" {
			return true
		}
	}
	return false
}

func (s *pmdWriteState) textureName(index int) string {
	if index < 0 {
		return ""
	}
	tex, err := s.model.Textures.Get(index)
	if err != nil || tex == nil || !tex.IsValid() {
		return ""
	}
	return tex.Name()
}

func (s *pmdWriteState) vertexBoneWeights(vertex *model.Vertex) (int, int, int) {
	idx0 := 0
	idx1 := 0
	weight := 100
	if vertex == nil || vertex.Deform == nil {
		return idx0, idx1, weight
	}
	switch deform := vertex.Deform.(type) {
	case *model.Bdef1:
		idx0 = s.boneMapping.mapIndex(deform.Indexes()[0])
		if idx0 < 0 {
			idx0 = 0
		}
		idx1 = idx0
		weight = 100
	case *model.Bdef2:
		idx0 = s.boneMapping.mapIndex(deform.Indexes()[0])
		idx1 = s.boneMapping.mapIndex(deform.Indexes()[1])
		w := deform.Weights()[0]
		if idx0 < 0 && idx1 >= 0 {
			idx0 = idx1
			idx1 = idx0
			w = 1.0
		} else if idx0 < 0 {
			idx0 = 0
			idx1 = idx0
			w = 1.0
		} else if idx1 < 0 {
			idx1 = idx0
			w = 1.0
		}
		weight = int(math.Round(w * 100.0))
		if weight < 0 {
			weight = 0
		} else if weight > 100 {
			weight = 100
		}
	default:
		idx0 = 0
		idx1 = 0
		weight = 100
	}
	return idx0, idx1, weight
}

func (s *pmdWriteState) writeFixedString(text string, size int, label string) error {
	encoded, err := io_common.EncodeShiftJISFixed(text, size)
	if err != nil {
		return io_common.NewIoNameEncodeFailed("名称のエンコードに失敗しました: %s", err, label)
	}
	if err := s.writer.WriteBytes(encoded); err != nil {
		return io_common.NewIoSaveFailed("名称の書き込みに失敗しました: %s", err, label)
	}
	return nil
}

func (s *pmdWriteState) writeVec2(vec mmath.Vec2, positiveOnly bool) error {
	if err := s.writer.WriteFloat32(vec.X, 0, positiveOnly); err != nil {
		return io_common.NewIoSaveFailed("PMDVec2の書き込みに失敗しました", err)
	}
	if err := s.writer.WriteFloat32(vec.Y, 0, positiveOnly); err != nil {
		return io_common.NewIoSaveFailed("PMDVec2の書き込みに失敗しました", err)
	}
	return nil
}

func (s *pmdWriteState) writeVec3(vec mmath.Vec3, positiveOnly bool) error {
	if err := s.writer.WriteFloat32(vec.X, 0, positiveOnly); err != nil {
		return io_common.NewIoSaveFailed("PMDVec3の書き込みに失敗しました", err)
	}
	if err := s.writer.WriteFloat32(vec.Y, 0, positiveOnly); err != nil {
		return io_common.NewIoSaveFailed("PMDVec3の書き込みに失敗しました", err)
	}
	if err := s.writer.WriteFloat32(vec.Z, 0, positiveOnly); err != nil {
		return io_common.NewIoSaveFailed("PMDVec3の書き込みに失敗しました", err)
	}
	return nil
}
