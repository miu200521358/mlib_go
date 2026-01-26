// 指示: miu200521358
package pmd

import (
	"errors"
	"fmt"
	"io"
	"math"

	"github.com/miu200521358/mlib_go/pkg/adapter/io_common"
	"github.com/miu200521358/mlib_go/pkg/domain/mmath"
	"github.com/miu200521358/mlib_go/pkg/domain/model"
	"gonum.org/v1/gonum/spatial/r3"
)

type pmdBoneDisplay struct {
	boneIndex  int
	frameIndex int
}

// pmdReader はPMD読み込み処理を表す。
type pmdReader struct {
	reader                  *io_common.BinaryReader
	textureIndex            map[string]int
	boneCount               int
	skinCount               int
	baseSkinVertexIndexes   []int
	skinIndexMap            []int
	skinDisplayList         []int
	boneDisplayNames        []string
	boneDisplayNamesEnglish []string
	boneDisplayList         []pmdBoneDisplay
}

// newPmdReader はpmdReaderを生成する。
func newPmdReader(r io.Reader) *pmdReader {
	return &pmdReader{
		reader:       io_common.NewBinaryReader(r),
		textureIndex: make(map[string]int),
	}
}

// Read はPMDを読み込む。
func (p *pmdReader) Read(modelData *model.PmxModel) error {
	if modelData == nil {
		return io_common.NewIoParseFailed("PMDモデルがnilです", nil)
	}
	if err := p.readHeader(modelData); err != nil {
		return err
	}
	if err := p.readVertices(modelData); err != nil {
		return err
	}
	if err := p.readFaces(modelData); err != nil {
		return err
	}
	if err := p.readMaterials(modelData); err != nil {
		return err
	}
	if err := p.readBones(modelData); err != nil {
		return err
	}
	if err := p.readIk(modelData); err != nil {
		return err
	}
	if err := p.readSkins(modelData); err != nil {
		return err
	}
	if err := p.readSkinDisplayList(); err != nil {
		return err
	}
	if err := p.readBoneDisplayNames(); err != nil {
		return err
	}
	if err := p.readBoneDisplayList(); err != nil {
		return err
	}
	if err := p.readExtensions(modelData); err != nil {
		return err
	}
	p.buildDisplaySlots(modelData)
	return nil
}

func (p *pmdReader) readHeader(modelData *model.PmxModel) error {
	signature, err := p.reader.ReadBytes(3)
	if err != nil {
		return wrapParseFailed("PMD署名の読み込みに失敗しました", err)
	}
	if string(signature) != "Pmd" {
		return wrapFormatNotSupported("PMD署名が不正です", nil)
	}
	version, err := p.reader.ReadFloat32()
	if err != nil {
		return wrapParseFailed("PMDバージョンの読み込みに失敗しました", err)
	}
	if !nearVersion(version, 1.0) {
		return wrapFormatNotSupported("PMDバージョンが非対応です", nil)
	}
	name, err := p.readFixedString(20, "PMDモデル名")
	if err != nil {
		return err
	}
	comment, err := p.readFixedString(256, "PMDコメント")
	if err != nil {
		return err
	}
	modelData.SetName(name)
	modelData.Comment = comment
	return nil
}

func (p *pmdReader) readVertices(modelData *model.PmxModel) error {
	count, err := p.reader.ReadUint32()
	if err != nil {
		return wrapParseFailed("PMD頂点数の読み込みに失敗しました", err)
	}
	for i := 0; i < int(count); i++ {
		pos, err := p.reader.ReadVec3()
		if err != nil {
			return wrapParseFailed("PMD頂点位置の読み込みに失敗しました", err)
		}
		normal, err := p.reader.ReadVec3()
		if err != nil {
			return wrapParseFailed("PMD頂点法線の読み込みに失敗しました", err)
		}
		uv, err := p.reader.ReadVec2()
		if err != nil {
			return wrapParseFailed("PMD頂点UVの読み込みに失敗しました", err)
		}
		bone0, err := p.reader.ReadUint16()
		if err != nil {
			return wrapParseFailed("PMD頂点ボーン番号の読み込みに失敗しました", err)
		}
		bone1, err := p.reader.ReadUint16()
		if err != nil {
			return wrapParseFailed("PMD頂点ボーン番号の読み込みに失敗しました", err)
		}
		weightRaw, err := p.reader.ReadUint8()
		if err != nil {
			return wrapParseFailed("PMD頂点ボーンウェイトの読み込みに失敗しました", err)
		}
		edgeFlag, err := p.reader.ReadUint8()
		if err != nil {
			return wrapParseFailed("PMD頂点エッジフラグの読み込みに失敗しました", err)
		}
		idx0 := int(bone0)
		idx1 := int(bone1)
		if bone0 == 0xFFFF {
			idx0 = -1
		}
		if bone1 == 0xFFFF {
			idx1 = -1
		}
		weight0 := float64(weightRaw) / 100.0
		if idx0 < 0 && idx1 >= 0 {
			idx0 = idx1
			idx1 = -1
			weight0 = 1.0
		}
		if idx0 < 0 {
			idx0 = 0
		}
		var deform model.IDeform
		if idx1 < 0 || weight0 >= 0.999 {
			deform = model.NewBdef1(idx0)
		} else {
			deform = model.NewBdef2(idx0, idx1, weight0)
		}
		edgeFactor := 1.0
		if edgeFlag != 0 {
			edgeFactor = 0.0
		}

		vertex := &model.Vertex{
			Position:   pos,
			Normal:     normal,
			Uv:         uv,
			DeformType: deform.DeformType(),
			Deform:     deform,
			EdgeFactor: edgeFactor,
		}
		modelData.Vertices.AppendRaw(vertex)
	}
	return nil
}

func (p *pmdReader) readFaces(modelData *model.PmxModel) error {
	count, err := p.reader.ReadUint32()
	if err != nil {
		return wrapParseFailed("PMD面頂点数の読み込みに失敗しました", err)
	}
	if count%3 != 0 {
		return wrapParseFailed("PMD面頂点数が3の倍数ではありません", nil)
	}
	for i := 0; i < int(count); i += 3 {
		v0, err := p.reader.ReadUint16()
		if err != nil {
			return wrapParseFailed("PMD面頂点の読み込みに失敗しました", err)
		}
		v1, err := p.reader.ReadUint16()
		if err != nil {
			return wrapParseFailed("PMD面頂点の読み込みに失敗しました", err)
		}
		v2, err := p.reader.ReadUint16()
		if err != nil {
			return wrapParseFailed("PMD面頂点の読み込みに失敗しました", err)
		}
		face := &model.Face{VertexIndexes: [3]int{int(v0), int(v1), int(v2)}}
		modelData.Faces.AppendRaw(face)
	}
	return nil
}

func (p *pmdReader) readMaterials(modelData *model.PmxModel) error {
	count, err := p.reader.ReadUint32()
	if err != nil {
		return wrapParseFailed("PMD材質数の読み込みに失敗しました", err)
	}
	for i := 0; i < int(count); i++ {
		diffuse, err := p.reader.ReadVec3()
		if err != nil {
			return wrapParseFailed("PMD材質Diffuseの読み込みに失敗しました", err)
		}
		alpha, err := p.reader.ReadFloat32()
		if err != nil {
			return wrapParseFailed("PMD材質Alphaの読み込みに失敗しました", err)
		}
		specularity, err := p.reader.ReadFloat32()
		if err != nil {
			return wrapParseFailed("PMD材質Specularityの読み込みに失敗しました", err)
		}
		specularColor, err := p.reader.ReadVec3()
		if err != nil {
			return wrapParseFailed("PMD材質Specularの読み込みに失敗しました", err)
		}
		ambient, err := p.reader.ReadVec3()
		if err != nil {
			return wrapParseFailed("PMD材質Ambientの読み込みに失敗しました", err)
		}
		toonIndex, err := p.reader.ReadUint8()
		if err != nil {
			return wrapParseFailed("PMD材質Toonの読み込みに失敗しました", err)
		}
		edgeFlag, err := p.reader.ReadUint8()
		if err != nil {
			return wrapParseFailed("PMD材質エッジフラグの読み込みに失敗しました", err)
		}
		faceVertCount, err := p.reader.ReadUint32()
		if err != nil {
			return wrapParseFailed("PMD材質面頂点数の読み込みに失敗しました", err)
		}
		texName, err := p.readFixedString(20, "PMD材質テクスチャ名")
		if err != nil {
			return err
		}
		spec := parseTextureSpec(texName)
		textureIndex := p.findOrAppendTexture(modelData, spec.textureName)
		sphereIndex := p.findOrAppendTexture(modelData, spec.sphereName)
		if spec.sphereName == "" {
			sphereIndex = -1
			spec.sphereMode = model.SPHERE_MODE_INVALID
		}

		material := model.NewMaterial()
		material.SetName(fmt.Sprintf("材質%02d", modelData.Materials.Len()+1))
		material.Diffuse = mmath.Vec4{X: diffuse.X, Y: diffuse.Y, Z: diffuse.Z, W: alpha}
		material.Specular = mmath.Vec4{X: specularColor.X, Y: specularColor.Y, Z: specularColor.Z, W: specularity}
		material.Ambient = ambient
		if edgeFlag == 0 {
			material.DrawFlag = model.DRAW_FLAG_DRAWING_EDGE
			material.EdgeSize = 1.0
		} else {
			material.DrawFlag = model.DRAW_FLAG_NONE
			material.EdgeSize = 0.0
		}
		material.Edge = mmath.UNIT_W_VEC4
		material.TextureIndex = textureIndex
		material.SphereTextureIndex = sphereIndex
		material.SphereMode = spec.sphereMode
		material.ToonSharingFlag = model.TOON_SHARING_SHARING
		if toonIndex <= 9 {
			material.ToonTextureIndex = int(toonIndex)
		} else {
			material.ToonTextureIndex = -1
		}
		material.VerticesCount = int(faceVertCount)
		modelData.Materials.AppendRaw(material)
	}
	return nil
}

func (p *pmdReader) readBones(modelData *model.PmxModel) error {
	count, err := p.reader.ReadUint16()
	if err != nil {
		return wrapParseFailed("PMDボーン数の読み込みに失敗しました", err)
	}
	p.boneCount = int(count)
	for i := 0; i < int(count); i++ {
		name, err := p.readFixedString(20, "PMDボーン名")
		if err != nil {
			return err
		}
		parentRaw, err := p.reader.ReadUint16()
		if err != nil {
			return wrapParseFailed("PMDボーン親番号の読み込みに失敗しました", err)
		}
		tailRaw, err := p.reader.ReadUint16()
		if err != nil {
			return wrapParseFailed("PMDボーン接続番号の読み込みに失敗しました", err)
		}
		boneType, err := p.reader.ReadUint8()
		if err != nil {
			return wrapParseFailed("PMDボーン種別の読み込みに失敗しました", err)
		}
		ikParentRaw, err := p.reader.ReadUint16()
		if err != nil {
			return wrapParseFailed("PMDボーンIK親の読み込みに失敗しました", err)
		}
		pos, err := p.reader.ReadVec3()
		if err != nil {
			return wrapParseFailed("PMDボーン位置の読み込みに失敗しました", err)
		}

		parentIndex := int(parentRaw)
		if parentRaw == 0xFFFF {
			parentIndex = -1
		}
		tailIndex := int(tailRaw)
		if tailRaw == 0xFFFF {
			tailIndex = -1
		}
		bone := &model.Bone{
			Position:    pos,
			ParentIndex: parentIndex,
			TailIndex:   tailIndex,
			Layer:       0,
			BoneFlag:    boneFlagsFromType(boneType, tailIndex >= 0),
		}
		bone.SetName(name)
		if boneType == 4 || boneType == 5 {
			if ikParentRaw != 0xFFFF {
				bone.EffectIndex = int(ikParentRaw)
				bone.EffectFactor = 1.0
			}
		}
		modelData.Bones.AppendRaw(bone)
	}
	return nil
}

func (p *pmdReader) readIk(modelData *model.PmxModel) error {
	count, err := p.reader.ReadUint16()
	if err != nil {
		return wrapParseFailed("PMD IK数の読み込みに失敗しました", err)
	}
	for i := 0; i < int(count); i++ {
		ikBoneRaw, err := p.reader.ReadUint16()
		if err != nil {
			return wrapParseFailed("PMD IKボーン番号の読み込みに失敗しました", err)
		}
		targetRaw, err := p.reader.ReadUint16()
		if err != nil {
			return wrapParseFailed("PMD IKターゲット番号の読み込みに失敗しました", err)
		}
		chainLen, err := p.reader.ReadUint8()
		if err != nil {
			return wrapParseFailed("PMD IKチェーン長の読み込みに失敗しました", err)
		}
		iterations, err := p.reader.ReadUint16()
		if err != nil {
			return wrapParseFailed("PMD IK反復回数の読み込みに失敗しました", err)
		}
		controlWeight, err := p.reader.ReadFloat32()
		if err != nil {
			return wrapParseFailed("PMD IK制限角の読み込みに失敗しました", err)
		}
		links := make([]model.IkLink, 0, chainLen)
		for j := 0; j < int(chainLen); j++ {
			linkRaw, err := p.reader.ReadUint16()
			if err != nil {
				return wrapParseFailed("PMD IKリンクの読み込みに失敗しました", err)
			}
			links = append(links, model.IkLink{BoneIndex: int(linkRaw)})
		}
		ikBoneIndex := int(ikBoneRaw)
		if ikBoneIndex < 0 || ikBoneIndex >= modelData.Bones.Len() {
			continue
		}
		ik := &model.Ik{
			BoneIndex:    int(targetRaw),
			LoopCount:    int(iterations),
			UnitRotation: mmath.Vec3{Vec: r3.Vec{X: controlWeight, Y: controlWeight, Z: controlWeight}},
			Links:        links,
		}
		bone, err := modelData.Bones.Get(ikBoneIndex)
		if err != nil {
			continue
		}
		bone.Ik = ik
		bone.BoneFlag |= model.BONE_FLAG_IS_IK
	}
	return nil
}

func (p *pmdReader) readSkins(modelData *model.PmxModel) error {
	count, err := p.reader.ReadUint16()
	if err != nil {
		return wrapParseFailed("PMD表情数の読み込みに失敗しました", err)
	}
	p.skinCount = int(count)
	p.skinIndexMap = make([]int, p.skinCount)
	for i := range p.skinIndexMap {
		p.skinIndexMap[i] = -1
	}
	for i := 0; i < int(count); i++ {
		name, err := p.readFixedString(20, "PMD表情名")
		if err != nil {
			return err
		}
		vertCount, err := p.reader.ReadUint32()
		if err != nil {
			return wrapParseFailed("PMD表情頂点数の読み込みに失敗しました", err)
		}
		skinType, err := p.reader.ReadUint8()
		if err != nil {
			return wrapParseFailed("PMD表情種別の読み込みに失敗しました", err)
		}

		if skinType == 0 {
			p.baseSkinVertexIndexes = make([]int, 0, vertCount)
			for j := 0; j < int(vertCount); j++ {
				vertIndex, err := p.reader.ReadUint32()
				if err != nil {
					return wrapParseFailed("PMD表情頂点の読み込みに失敗しました", err)
				}
				pos, err := p.reader.ReadVec3()
				if err != nil {
					return wrapParseFailed("PMD表情頂点位置の読み込みに失敗しました", err)
				}
				_ = pos
				p.baseSkinVertexIndexes = append(p.baseSkinVertexIndexes, int(vertIndex))
			}
			continue
		}
		if len(p.baseSkinVertexIndexes) == 0 {
			return wrapParseFailed("PMD表情ベースが存在しません", nil)
		}

		offsets := make([]model.MorphOffset, 0, vertCount)
		for j := 0; j < int(vertCount); j++ {
			baseIndexRaw, err := p.reader.ReadUint32()
			if err != nil {
				return wrapParseFailed("PMD表情頂点の読み込みに失敗しました", err)
			}
			pos, err := p.reader.ReadVec3()
			if err != nil {
				return wrapParseFailed("PMD表情頂点位置の読み込みに失敗しました", err)
			}
			baseIndex := int(baseIndexRaw)
			if baseIndex < 0 || baseIndex >= len(p.baseSkinVertexIndexes) {
				return wrapParseFailed("PMD表情頂点の参照が不正です", nil)
			}
			vertexIndex := p.baseSkinVertexIndexes[baseIndex]
			offsets = append(offsets, &model.VertexMorphOffset{VertexIndex: vertexIndex, Position: pos})
		}

		morph := &model.Morph{
			Panel:     panelFromSkinType(skinType),
			MorphType: model.MORPH_TYPE_VERTEX,
			Offsets:   offsets,
		}
		morph.SetName(name)
		idx := modelData.Morphs.AppendRaw(morph)
		p.skinIndexMap[i] = idx
	}
	return nil
}

func (p *pmdReader) readSkinDisplayList() error {
	count, err := p.reader.ReadUint8()
	if err != nil {
		return wrapParseFailed("PMD表情表示数の読み込みに失敗しました", err)
	}
	p.skinDisplayList = make([]int, 0, count)
	for i := 0; i < int(count); i++ {
		idx, err := p.reader.ReadUint16()
		if err != nil {
			return wrapParseFailed("PMD表情表示の読み込みに失敗しました", err)
		}
		p.skinDisplayList = append(p.skinDisplayList, int(idx))
	}
	return nil
}

func (p *pmdReader) readBoneDisplayNames() error {
	count, err := p.reader.ReadUint8()
	if err != nil {
		return wrapParseFailed("PMDボーン枠名数の読み込みに失敗しました", err)
	}
	p.boneDisplayNames = make([]string, 0, count)
	for i := 0; i < int(count); i++ {
		name, err := p.readFixedString(50, "PMDボーン枠名")
		if err != nil {
			return err
		}
		p.boneDisplayNames = append(p.boneDisplayNames, name)
	}
	return nil
}

func (p *pmdReader) readBoneDisplayList() error {
	count, err := p.reader.ReadUint32()
	if err != nil {
		return wrapParseFailed("PMDボーン表示数の読み込みに失敗しました", err)
	}
	p.boneDisplayList = make([]pmdBoneDisplay, 0, count)
	for i := 0; i < int(count); i++ {
		boneIndex, err := p.reader.ReadUint16()
		if err != nil {
			return wrapParseFailed("PMDボーン表示の読み込みに失敗しました", err)
		}
		frameIndex, err := p.reader.ReadUint8()
		if err != nil {
			return wrapParseFailed("PMDボーン表示の読み込みに失敗しました", err)
		}
		p.boneDisplayList = append(p.boneDisplayList, pmdBoneDisplay{boneIndex: int(boneIndex), frameIndex: int(frameIndex)})
	}
	return nil
}

func (p *pmdReader) readExtensions(modelData *model.PmxModel) error {
	flag, err := p.reader.ReadUint8()
	if err != nil {
		if isEOF(err) {
			return nil
		}
		return wrapParseFailed("PMD英名対応フラグの読み込みに失敗しました", err)
	}
	englishEnabled := flag == 1
	if englishEnabled {
		name, err := p.readFixedString(20, "PMD英名モデル名")
		if err != nil {
			return err
		}
		comment, err := p.readFixedString(256, "PMD英名コメント")
		if err != nil {
			return err
		}
		modelData.EnglishName = name
		modelData.EnglishComment = comment
		for i := 0; i < p.boneCount; i++ {
			englishName, err := p.readFixedString(20, "PMD英名ボーン名")
			if err != nil {
				return err
			}
			bone, err := modelData.Bones.Get(i)
			if err == nil {
				bone.EnglishName = englishName
			}
		}
		for i := 0; i < p.skinCount-1; i++ {
			englishName, err := p.readFixedString(20, "PMD英名表情名")
			if err != nil {
				return err
			}
			pmdIndex := i + 1
			if pmdIndex < len(p.skinIndexMap) {
				morphIndex := p.skinIndexMap[pmdIndex]
				if morphIndex >= 0 {
					morph, err := modelData.Morphs.Get(morphIndex)
					if err == nil {
						morph.EnglishName = englishName
					}
				}
			}
		}
		p.boneDisplayNamesEnglish = make([]string, 0, len(p.boneDisplayNames))
		for i := 0; i < len(p.boneDisplayNames); i++ {
			englishName, err := p.readFixedString(50, "PMD英名ボーン枠名")
			if err != nil {
				return err
			}
			p.boneDisplayNamesEnglish = append(p.boneDisplayNamesEnglish, englishName)
		}
	}

	for i := 0; i < 10; i++ {
		if _, err := p.readFixedString(100, "PMDトゥーンテクスチャ名"); err != nil {
			return err
		}
	}

	rigidCount, err := p.reader.ReadUint32()
	if err != nil {
		if isEOF(err) {
			return nil
		}
		return wrapParseFailed("PMD剛体数の読み込みに失敗しました", err)
	}
	for i := 0; i < int(rigidCount); i++ {
		name, err := p.readFixedString(20, "PMD剛体名")
		if err != nil {
			return err
		}
		boneRaw, err := p.reader.ReadUint16()
		if err != nil {
			return wrapParseFailed("PMD剛体ボーン番号の読み込みに失敗しました", err)
		}
		group, err := p.reader.ReadUint8()
		if err != nil {
			return wrapParseFailed("PMD剛体グループの読み込みに失敗しました", err)
		}
		groupTarget, err := p.reader.ReadUint16()
		if err != nil {
			return wrapParseFailed("PMD剛体グループ対象の読み込みに失敗しました", err)
		}
		shapeType, err := p.reader.ReadUint8()
		if err != nil {
			return wrapParseFailed("PMD剛体形状の読み込みに失敗しました", err)
		}
		shapeW, err := p.reader.ReadFloat32()
		if err != nil {
			return wrapParseFailed("PMD剛体サイズの読み込みに失敗しました", err)
		}
		shapeH, err := p.reader.ReadFloat32()
		if err != nil {
			return wrapParseFailed("PMD剛体サイズの読み込みに失敗しました", err)
		}
		shapeD, err := p.reader.ReadFloat32()
		if err != nil {
			return wrapParseFailed("PMD剛体サイズの読み込みに失敗しました", err)
		}
		pos, err := p.reader.ReadVec3()
		if err != nil {
			return wrapParseFailed("PMD剛体位置の読み込みに失敗しました", err)
		}
		rot, err := p.reader.ReadVec3()
		if err != nil {
			return wrapParseFailed("PMD剛体回転の読み込みに失敗しました", err)
		}
		mass, err := p.reader.ReadFloat32()
		if err != nil {
			return wrapParseFailed("PMD剛体質量の読み込みに失敗しました", err)
		}
		linearDamping, err := p.reader.ReadFloat32()
		if err != nil {
			return wrapParseFailed("PMD剛体移動減衰の読み込みに失敗しました", err)
		}
		angularDamping, err := p.reader.ReadFloat32()
		if err != nil {
			return wrapParseFailed("PMD剛体回転減衰の読み込みに失敗しました", err)
		}
		restitution, err := p.reader.ReadFloat32()
		if err != nil {
			return wrapParseFailed("PMD剛体反発の読み込みに失敗しました", err)
		}
		friction, err := p.reader.ReadFloat32()
		if err != nil {
			return wrapParseFailed("PMD剛体摩擦の読み込みに失敗しました", err)
		}
		physType, err := p.reader.ReadUint8()
		if err != nil {
			return wrapParseFailed("PMD剛体タイプの読み込みに失敗しました", err)
		}

		boneIndex := int(boneRaw)
		if boneRaw == 0xFFFF {
			boneIndex = -1
		}
		rigid := &model.RigidBody{
			EnglishName: "",
			BoneIndex:   boneIndex,
			CollisionGroup: model.CollisionGroup{
				Group: group,
				Mask:  ^groupTarget,
			},
			Shape:       model.Shape(shapeType),
			Size:        mmath.Vec3{Vec: r3.Vec{X: shapeW, Y: shapeH, Z: shapeD}},
			Position:    pos,
			Rotation:    rot,
			PhysicsType: model.PhysicsType(physType),
			Param: model.RigidBodyParam{
				Mass:           mass,
				LinearDamping:  linearDamping,
				AngularDamping: angularDamping,
				Restitution:    restitution,
				Friction:       friction,
			},
		}
		rigid.SetName(name)
		modelData.RigidBodies.AppendRaw(rigid)
	}

	jointCount, err := p.reader.ReadUint32()
	if err != nil {
		return wrapParseFailed("PMDジョイント数の読み込みに失敗しました", err)
	}
	for i := 0; i < int(jointCount); i++ {
		name, err := p.readFixedString(20, "PMDジョイント名")
		if err != nil {
			return err
		}
		rigidA, err := p.reader.ReadUint32()
		if err != nil {
			return wrapParseFailed("PMDジョイント剛体Aの読み込みに失敗しました", err)
		}
		rigidB, err := p.reader.ReadUint32()
		if err != nil {
			return wrapParseFailed("PMDジョイント剛体Bの読み込みに失敗しました", err)
		}
		pos, err := p.reader.ReadVec3()
		if err != nil {
			return wrapParseFailed("PMDジョイント位置の読み込みに失敗しました", err)
		}
		rot, err := p.reader.ReadVec3()
		if err != nil {
			return wrapParseFailed("PMDジョイント回転の読み込みに失敗しました", err)
		}
		limitMin, err := p.reader.ReadVec3()
		if err != nil {
			return wrapParseFailed("PMDジョイント移動下限の読み込みに失敗しました", err)
		}
		limitMax, err := p.reader.ReadVec3()
		if err != nil {
			return wrapParseFailed("PMDジョイント移動上限の読み込みに失敗しました", err)
		}
		rotMin, err := p.reader.ReadVec3()
		if err != nil {
			return wrapParseFailed("PMDジョイント回転下限の読み込みに失敗しました", err)
		}
		rotMax, err := p.reader.ReadVec3()
		if err != nil {
			return wrapParseFailed("PMDジョイント回転上限の読み込みに失敗しました", err)
		}
		springPos, err := p.reader.ReadVec3()
		if err != nil {
			return wrapParseFailed("PMDジョイント移動ばねの読み込みに失敗しました", err)
		}
		springRot, err := p.reader.ReadVec3()
		if err != nil {
			return wrapParseFailed("PMDジョイント回転ばねの読み込みに失敗しました", err)
		}

		joint := &model.Joint{
			RigidBodyIndexA: int(rigidA),
			RigidBodyIndexB: int(rigidB),
			Param: model.JointParam{
				Position:                  pos,
				Rotation:                  rot,
				TranslationLimitMin:       limitMin,
				TranslationLimitMax:       limitMax,
				RotationLimitMin:          rotMin,
				RotationLimitMax:          rotMax,
				SpringConstantTranslation: springPos,
				SpringConstantRotation:    springRot,
			},
		}
		joint.SetName(name)
		modelData.Joints.AppendRaw(joint)
	}
	return nil
}

func (p *pmdReader) buildDisplaySlots(modelData *model.PmxModel) {
	if modelData == nil || modelData.DisplaySlots == nil {
		return
	}
	modelData.CreateDefaultDisplaySlots()
	if modelData.DisplaySlots.Len() < 2 {
		return
	}
	for i, name := range p.boneDisplayNames {
		slot := &model.DisplaySlot{
			SpecialFlag: model.SPECIAL_FLAG_OFF,
			References:  make([]model.Reference, 0),
		}
		slot.SetName(name)
		if i < len(p.boneDisplayNamesEnglish) {
			slot.EnglishName = p.boneDisplayNamesEnglish[i]
		}
		modelData.DisplaySlots.AppendRaw(slot)
	}

	rootSlot, _ := modelData.DisplaySlots.Get(0)
	morphSlot, _ := modelData.DisplaySlots.Get(1)
	for _, entry := range p.boneDisplayList {
		if entry.boneIndex < 0 {
			continue
		}
		var slot *model.DisplaySlot
		if entry.frameIndex == 0 {
			slot = rootSlot
		} else {
			slotIndex := 1 + entry.frameIndex
			if slotIndex < modelData.DisplaySlots.Len() {
				slot, _ = modelData.DisplaySlots.Get(slotIndex)
			}
		}
		if slot == nil {
			continue
		}
		slot.References = append(slot.References, model.Reference{DisplayType: model.DISPLAY_TYPE_BONE, DisplayIndex: entry.boneIndex})
	}

	for _, skinIndex := range p.skinDisplayList {
		if skinIndex < 0 || skinIndex >= len(p.skinIndexMap) {
			continue
		}
		morphIndex := p.skinIndexMap[skinIndex]
		if morphIndex < 0 {
			continue
		}
		if morphSlot != nil {
			morphSlot.References = append(morphSlot.References, model.Reference{DisplayType: model.DISPLAY_TYPE_MORPH, DisplayIndex: morphIndex})
		}
	}
}

func (p *pmdReader) findOrAppendTexture(modelData *model.PmxModel, name string) int {
	if name == "" {
		return -1
	}
	if idx, ok := p.textureIndex[name]; ok {
		return idx
	}
	tex := model.NewTexture()
	tex.SetName(name)
	tex.SetValid(true)
	idx := modelData.Textures.AppendRaw(tex)
	p.textureIndex[name] = idx
	return idx
}

func (p *pmdReader) readFixedString(size int, label string) (string, error) {
	raw, err := p.reader.ReadBytes(size)
	if err != nil {
		return "", wrapParseFailed(label+"の読み込みに失敗しました", err)
	}
	text, err := io_common.DecodeShiftJISFixed(raw)
	if err != nil {
		return "", wrapEncodingUnknown(label+"のデコードに失敗しました", err)
	}
	return text, nil
}

func nearVersion(value, expected float64) bool {
	return math.Abs(value-expected) < 0.01
}

func isEOF(err error) bool {
	return errors.Is(err, io.EOF) || errors.Is(err, io.ErrUnexpectedEOF)
}

func wrapParseFailed(message string, err error) error {
	return io_common.NewIoParseFailed(message, err)
}

func wrapFormatNotSupported(message string, err error) error {
	return io_common.NewIoFormatNotSupported(message, err)
}

func wrapEncodingUnknown(message string, err error) error {
	return io_common.NewIoEncodingUnknown(message, err)
}

func wrapEncodeFailed(message string, err error) error {
	return io_common.NewIoEncodeFailed(message, err)
}
