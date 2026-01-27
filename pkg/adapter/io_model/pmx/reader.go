// 指示: miu200521358
package pmx

import (
	"io"
	"math"

	"github.com/miu200521358/mlib_go/pkg/adapter/io_common"
	"github.com/miu200521358/mlib_go/pkg/domain/mmath"
	"github.com/miu200521358/mlib_go/pkg/domain/model"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/unicode"
	"gonum.org/v1/gonum/spatial/r3"
)

// pmxHeader はPMXヘッダ情報を表す。
type pmxHeader struct {
	version            float64
	encodeType         byte
	encoding           encoding.Encoding
	extendedUVCount    int
	vertexIndexSize    byte
	textureIndexSize   byte
	materialIndexSize  byte
	boneIndexSize      byte
	morphIndexSize     byte
	rigidBodyIndexSize byte
}

// pmxReader はPMX読み込み処理を表す。
type pmxReader struct {
	reader *io_common.BinaryReader
	header *pmxHeader
}

// newPmxReader はpmxReaderを生成する。
func newPmxReader(r io.Reader) *pmxReader {
	return &pmxReader{reader: io_common.NewBinaryReader(r)}
}

// Read はPMXを読み込む。
func (p *pmxReader) Read(modelData *model.PmxModel) error {
	if modelData == nil {
		return io_common.NewIoParseFailed("PMXモデルがnilです", nil)
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
	if err := p.readTextures(modelData); err != nil {
		return err
	}
	if err := p.readMaterials(modelData); err != nil {
		return err
	}
	if err := p.readBones(modelData); err != nil {
		return err
	}
	if err := p.readMorphs(modelData); err != nil {
		return err
	}
	if err := p.readDisplaySlots(modelData); err != nil {
		return err
	}
	if err := p.readRigidBodies(modelData); err != nil {
		return err
	}
	if err := p.readJoints(modelData); err != nil {
		return err
	}
	if err := p.skipSoftBodies(); err != nil {
		return err
	}
	return nil
}

// readHeader はヘッダとモデル情報を読み込む。
func (p *pmxReader) readHeader(modelData *model.PmxModel) error {
	signature, err := p.reader.ReadBytes(4)
	if err != nil {
		return wrapParseFailed("PMX署名の読み込みに失敗しました", err)
	}
	if len(signature) < 3 || string(signature[:3]) != "PMX" {
		return wrapFormatNotSupported("PMX署名が不正です", nil)
	}
	version, err := p.reader.ReadFloat32()
	if err != nil {
		return wrapParseFailed("PMXバージョンの読み込みに失敗しました", err)
	}
	if !nearVersion(version, 2.0) && !nearVersion(version, 2.1) {
		return wrapFormatNotSupported("PMXバージョンが非対応です", nil)
	}
	_, err = p.reader.ReadUint8()
	if err != nil {
		return wrapParseFailed("PMXヘッダサイズの読み込みに失敗しました", err)
	}
	encodeType, err := p.reader.ReadUint8()
	if err != nil {
		return wrapParseFailed("PMX文字コード種別の読み込みに失敗しました", err)
	}
	enc, err := resolveEncoding(encodeType)
	if err != nil {
		return wrapEncodingUnknown("PMX文字コード種別が不正です", err)
	}
	extendedUVCount, err := p.reader.ReadUint8()
	if err != nil {
		return wrapParseFailed("PMX拡張UV数の読み込みに失敗しました", err)
	}
	if extendedUVCount > 4 {
		return wrapFormatNotSupported("PMX拡張UV数が不正です", nil)
	}
	vertexIndexSize, err := p.reader.ReadUint8()
	if err != nil {
		return wrapParseFailed("PMX頂点インデックスサイズの読み込みに失敗しました", err)
	}
	textureIndexSize, err := p.reader.ReadUint8()
	if err != nil {
		return wrapParseFailed("PMXテクスチャインデックスサイズの読み込みに失敗しました", err)
	}
	materialIndexSize, err := p.reader.ReadUint8()
	if err != nil {
		return wrapParseFailed("PMX材質インデックスサイズの読み込みに失敗しました", err)
	}
	boneIndexSize, err := p.reader.ReadUint8()
	if err != nil {
		return wrapParseFailed("PMXボーンインデックスサイズの読み込みに失敗しました", err)
	}
	morphIndexSize, err := p.reader.ReadUint8()
	if err != nil {
		return wrapParseFailed("PMXモーフインデックスサイズの読み込みに失敗しました", err)
	}
	rigidBodyIndexSize, err := p.reader.ReadUint8()
	if err != nil {
		return wrapParseFailed("PMX剛体インデックスサイズの読み込みに失敗しました", err)
	}
	if !isValidIndexSize(vertexIndexSize) ||
		!isValidIndexSize(textureIndexSize) ||
		!isValidIndexSize(materialIndexSize) ||
		!isValidIndexSize(boneIndexSize) ||
		!isValidIndexSize(morphIndexSize) ||
		!isValidIndexSize(rigidBodyIndexSize) {
		return wrapFormatNotSupported("PMXインデックスサイズが不正です", nil)
	}

	p.header = &pmxHeader{
		version:            version,
		encodeType:         encodeType,
		encoding:           enc,
		extendedUVCount:    int(extendedUVCount),
		vertexIndexSize:    vertexIndexSize,
		textureIndexSize:   textureIndexSize,
		materialIndexSize:  materialIndexSize,
		boneIndexSize:      boneIndexSize,
		morphIndexSize:     morphIndexSize,
		rigidBodyIndexSize: rigidBodyIndexSize,
	}

	name, err := p.readText()
	if err != nil {
		return err
	}
	englishName, err := p.readText()
	if err != nil {
		return err
	}
	comment, err := p.readText()
	if err != nil {
		return err
	}
	englishComment, err := p.readText()
	if err != nil {
		return err
	}

	modelData.SetName(name)
	modelData.EnglishName = englishName
	modelData.Comment = comment
	modelData.EnglishComment = englishComment
	return nil
}

// readText はTextBufを読み込む。
func (p *pmxReader) readText() (string, error) {
	length, err := p.reader.ReadInt32()
	if err != nil {
		return "", wrapParseFailed("PMX文字列長の読み込みに失敗しました", err)
	}
	if length < 0 {
		return "", wrapParseFailed("PMX文字列長が不正です", nil)
	}
	if length == 0 {
		return "", nil
	}
	bytes, err := p.reader.ReadBytes(int(length))
	if err != nil {
		return "", wrapParseFailed("PMX文字列の読み込みに失敗しました", err)
	}
	text, err := io_common.DecodePmxText(p.header.encoding, bytes)
	if err != nil {
		return "", wrapEncodingUnknown("PMX文字列のデコードに失敗しました", err)
	}
	return text, nil
}

// readVertices は頂点セクションを読み込む。
func (p *pmxReader) readVertices(modelData *model.PmxModel) error {
	count, err := p.reader.ReadInt32()
	if err != nil {
		return wrapParseFailed("PMX頂点数の読み込みに失敗しました", err)
	}
	if count < 0 {
		return wrapParseFailed("PMX頂点数が不正です", nil)
	}
	for i := 0; i < int(count); i++ {
		pos, err := p.reader.ReadVec3()
		if err != nil {
			return wrapParseFailed("PMX頂点位置の読み込みに失敗しました", err)
		}
		normal, err := p.reader.ReadVec3()
		if err != nil {
			return wrapParseFailed("PMX頂点法線の読み込みに失敗しました", err)
		}
		uv, err := p.reader.ReadVec2()
		if err != nil {
			return wrapParseFailed("PMX頂点UVの読み込みに失敗しました", err)
		}
		extended := make([]mmath.Vec4, 0, p.header.extendedUVCount)
		for j := 0; j < p.header.extendedUVCount; j++ {
			uv4, err := p.reader.ReadVec4()
			if err != nil {
				return wrapParseFailed("PMX頂点拡張UVの読み込みに失敗しました", err)
			}
			extended = append(extended, uv4)
		}

		deformType, err := p.reader.ReadUint8()
		if err != nil {
			return wrapParseFailed("PMXデフォーム種別の読み込みに失敗しました", err)
		}
		deform, err := p.readDeform(model.DeformType(deformType))
		if err != nil {
			return err
		}
		edgeFactor, err := p.reader.ReadFloat32()
		if err != nil {
			return wrapParseFailed("PMXエッジ倍率の読み込みに失敗しました", err)
		}

		vertex := &model.Vertex{
			Position:    pos,
			Normal:      normal,
			Uv:          uv,
			ExtendedUvs: extended,
			DeformType:  model.DeformType(deformType),
			Deform:      deform,
			EdgeFactor:  edgeFactor,
		}
		modelData.Vertices.AppendRaw(vertex)
	}
	return nil
}

// readDeform はデフォーム情報を読み込む。
func (p *pmxReader) readDeform(deformType model.DeformType) (model.IDeform, error) {
	switch deformType {
	case model.BDEF1:
		boneIndex, err := readSignedIndex(p.reader, p.header.boneIndexSize)
		if err != nil {
			return nil, wrapParseFailed("PMXデフォーム(BDEF1)の読み込みに失敗しました", err)
		}
		return model.NewBdef1(boneIndex), nil
	case model.BDEF2:
		boneIndex0, err := readSignedIndex(p.reader, p.header.boneIndexSize)
		if err != nil {
			return nil, wrapParseFailed("PMXデフォーム(BDEF2)の読み込みに失敗しました", err)
		}
		boneIndex1, err := readSignedIndex(p.reader, p.header.boneIndexSize)
		if err != nil {
			return nil, wrapParseFailed("PMXデフォーム(BDEF2)の読み込みに失敗しました", err)
		}
		weight0, err := p.reader.ReadFloat32()
		if err != nil {
			return nil, wrapParseFailed("PMXデフォーム(BDEF2)の読み込みに失敗しました", err)
		}
		return model.NewBdef2(boneIndex0, boneIndex1, weight0), nil
	case model.BDEF4:
		indexes := [4]int{}
		weights := [4]float64{}
		for i := 0; i < 4; i++ {
			idx, err := readSignedIndex(p.reader, p.header.boneIndexSize)
			if err != nil {
				return nil, wrapParseFailed("PMXデフォーム(BDEF4)の読み込みに失敗しました", err)
			}
			indexes[i] = idx
		}
		for i := 0; i < 4; i++ {
			weight, err := p.reader.ReadFloat32()
			if err != nil {
				return nil, wrapParseFailed("PMXデフォーム(BDEF4)の読み込みに失敗しました", err)
			}
			weights[i] = weight
		}
		return model.NewBdef4(indexes, weights), nil
	case model.SDEF:
		boneIndex0, err := readSignedIndex(p.reader, p.header.boneIndexSize)
		if err != nil {
			return nil, wrapParseFailed("PMXデフォーム(SDEF)の読み込みに失敗しました", err)
		}
		boneIndex1, err := readSignedIndex(p.reader, p.header.boneIndexSize)
		if err != nil {
			return nil, wrapParseFailed("PMXデフォーム(SDEF)の読み込みに失敗しました", err)
		}
		weight0, err := p.reader.ReadFloat32()
		if err != nil {
			return nil, wrapParseFailed("PMXデフォーム(SDEF)の読み込みに失敗しました", err)
		}
		sdef := model.NewSdef(boneIndex0, boneIndex1, weight0)
		c, err := p.reader.ReadVec3()
		if err != nil {
			return nil, wrapParseFailed("PMXデフォーム(SDEF)の読み込みに失敗しました", err)
		}
		r0, err := p.reader.ReadVec3()
		if err != nil {
			return nil, wrapParseFailed("PMXデフォーム(SDEF)の読み込みに失敗しました", err)
		}
		r1, err := p.reader.ReadVec3()
		if err != nil {
			return nil, wrapParseFailed("PMXデフォーム(SDEF)の読み込みに失敗しました", err)
		}
		sdef.SdefC = c
		sdef.SdefR0 = r0
		sdef.SdefR1 = r1
		return sdef, nil
	default:
		return nil, wrapFormatNotSupported("PMXデフォーム種別が未対応です", nil)
	}
}

// readFaces は面セクションを読み込む。
func (p *pmxReader) readFaces(modelData *model.PmxModel) error {
	count, err := p.reader.ReadInt32()
	if err != nil {
		return wrapParseFailed("PMX面インデックス数の読み込みに失敗しました", err)
	}
	if count < 0 {
		return wrapParseFailed("PMX面インデックス数が不正です", nil)
	}
	if count%3 != 0 {
		return wrapParseFailed("PMX面インデックス数が3の倍数ではありません", nil)
	}
	faces := int(count / 3)
	for i := 0; i < faces; i++ {
		idx0, err := readVertexIndex(p.reader, p.header.vertexIndexSize)
		if err != nil {
			return wrapParseFailed("PMX面インデックスの読み込みに失敗しました", err)
		}
		idx1, err := readVertexIndex(p.reader, p.header.vertexIndexSize)
		if err != nil {
			return wrapParseFailed("PMX面インデックスの読み込みに失敗しました", err)
		}
		idx2, err := readVertexIndex(p.reader, p.header.vertexIndexSize)
		if err != nil {
			return wrapParseFailed("PMX面インデックスの読み込みに失敗しました", err)
		}
		face := &model.Face{VertexIndexes: [3]int{idx0, idx1, idx2}}
		modelData.Faces.AppendRaw(face)
	}
	return nil
}

// readTextures はテクスチャセクションを読み込む。
func (p *pmxReader) readTextures(modelData *model.PmxModel) error {
	count, err := p.reader.ReadInt32()
	if err != nil {
		return wrapParseFailed("PMXテクスチャ数の読み込みに失敗しました", err)
	}
	if count < 0 {
		return wrapParseFailed("PMXテクスチャ数が不正です", nil)
	}
	for i := 0; i < int(count); i++ {
		name, err := p.readText()
		if err != nil {
			return err
		}
		tex := model.NewTexture()
		tex.SetName(name)
		tex.SetValid(true)
		modelData.Textures.AppendRaw(tex)
	}
	return nil
}

// readMaterials は材質セクションを読み込む。
func (p *pmxReader) readMaterials(modelData *model.PmxModel) error {
	count, err := p.reader.ReadInt32()
	if err != nil {
		return wrapParseFailed("PMX材質数の読み込みに失敗しました", err)
	}
	if count < 0 {
		return wrapParseFailed("PMX材質数が不正です", nil)
	}
	for i := 0; i < int(count); i++ {
		name, err := p.readText()
		if err != nil {
			return err
		}
		englishName, err := p.readText()
		if err != nil {
			return err
		}
		diffuse, err := p.reader.ReadVec4()
		if err != nil {
			return wrapParseFailed("PMX材質Diffuseの読み込みに失敗しました", err)
		}
		specular, err := p.reader.ReadVec3()
		if err != nil {
			return wrapParseFailed("PMX材質Specularの読み込みに失敗しました", err)
		}
		specularPower, err := p.reader.ReadFloat32()
		if err != nil {
			return wrapParseFailed("PMX材質SpecularPowerの読み込みに失敗しました", err)
		}
		ambient, err := p.reader.ReadVec3()
		if err != nil {
			return wrapParseFailed("PMX材質Ambientの読み込みに失敗しました", err)
		}
		drawFlag, err := p.reader.ReadUint8()
		if err != nil {
			return wrapParseFailed("PMX材質描画フラグの読み込みに失敗しました", err)
		}
		edgeColor, err := p.reader.ReadVec4()
		if err != nil {
			return wrapParseFailed("PMX材質エッジ色の読み込みに失敗しました", err)
		}
		edgeSize, err := p.reader.ReadFloat32()
		if err != nil {
			return wrapParseFailed("PMX材質エッジサイズの読み込みに失敗しました", err)
		}
		textureIndex, err := readSignedIndex(p.reader, p.header.textureIndexSize)
		if err != nil {
			return wrapParseFailed("PMX材質テクスチャ参照の読み込みに失敗しました", err)
		}
		sphereTextureIndex, err := readSignedIndex(p.reader, p.header.textureIndexSize)
		if err != nil {
			return wrapParseFailed("PMX材質スフィア参照の読み込みに失敗しました", err)
		}
		sphereMode, err := p.reader.ReadUint8()
		if err != nil {
			return wrapParseFailed("PMX材質スフィアモードの読み込みに失敗しました", err)
		}
		toonSharingFlag, err := p.reader.ReadUint8()
		if err != nil {
			return wrapParseFailed("PMX材質トゥーン共有フラグの読み込みに失敗しました", err)
		}
		toonTextureIndex := -1
		if toonSharingFlag == 0 {
			idx, err := readSignedIndex(p.reader, p.header.textureIndexSize)
			if err != nil {
				return wrapParseFailed("PMX材質トゥーン参照の読み込みに失敗しました", err)
			}
			toonTextureIndex = idx
		} else {
			idx, err := p.reader.ReadUint8()
			if err != nil {
				return wrapParseFailed("PMX材質共有トゥーン参照の読み込みに失敗しました", err)
			}
			toonTextureIndex = int(idx)
		}
		memo, err := p.readText()
		if err != nil {
			return err
		}
		verticesCount, err := p.reader.ReadInt32()
		if err != nil {
			return wrapParseFailed("PMX材質頂点数の読み込みに失敗しました", err)
		}

		material := model.NewMaterial()
		material.SetName(name)
		material.EnglishName = englishName
		material.Diffuse = diffuse
		material.Specular = mmath.Vec4{X: specular.X, Y: specular.Y, Z: specular.Z, W: specularPower}
		material.Ambient = ambient
		material.DrawFlag = model.DrawFlag(drawFlag)
		material.Edge = edgeColor
		material.EdgeSize = edgeSize
		material.TextureIndex = textureIndex
		material.SphereTextureIndex = sphereTextureIndex
		material.SphereMode = model.SphereMode(sphereMode)
		material.ToonSharingFlag = model.ToonSharingFlag(toonSharingFlag)
		material.ToonTextureIndex = toonTextureIndex
		material.Memo = memo
		material.VerticesCount = int(verticesCount)
		modelData.Materials.AppendRaw(material)
	}
	return nil
}

// readBones はボーンセクションを読み込む。
func (p *pmxReader) readBones(modelData *model.PmxModel) error {
	count, err := p.reader.ReadInt32()
	if err != nil {
		return wrapParseFailed("PMXボーン数の読み込みに失敗しました", err)
	}
	if count < 0 {
		return wrapParseFailed("PMXボーン数が不正です", nil)
	}
	for i := 0; i < int(count); i++ {
		name, err := p.readText()
		if err != nil {
			return err
		}
		englishName, err := p.readText()
		if err != nil {
			return err
		}
		pos, err := p.reader.ReadVec3()
		if err != nil {
			return wrapParseFailed("PMXボーン位置の読み込みに失敗しました", err)
		}
		parentIndex, err := readSignedIndex(p.reader, p.header.boneIndexSize)
		if err != nil {
			return wrapParseFailed("PMXボーン親インデックスの読み込みに失敗しました", err)
		}
		layer, err := p.reader.ReadInt32()
		if err != nil {
			return wrapParseFailed("PMXボーンレイヤーの読み込みに失敗しました", err)
		}
		boneFlagRaw, err := p.reader.ReadUint16()
		if err != nil {
			return wrapParseFailed("PMXボーンフラグの読み込みに失敗しました", err)
		}
		boneFlag := model.BoneFlag(boneFlagRaw)

		bone := &model.Bone{
			EnglishName: englishName,
			Position:    pos,
			ParentIndex: parentIndex,
			Layer:       int(layer),
			BoneFlag:    boneFlag,
		}
		bone.SetName(name)

		if boneFlag&model.BONE_FLAG_TAIL_IS_BONE != 0 {
			idx, err := readSignedIndex(p.reader, p.header.boneIndexSize)
			if err != nil {
				return wrapParseFailed("PMXボーン接続インデックスの読み込みに失敗しました", err)
			}
			bone.TailIndex = idx
		} else {
			pos, err := p.reader.ReadVec3()
			if err != nil {
				return wrapParseFailed("PMXボーン接続座標の読み込みに失敗しました", err)
			}
			bone.TailPosition = pos
			bone.TailIndex = -1
		}

		if boneFlag&(model.BONE_FLAG_IS_EXTERNAL_ROTATION|model.BONE_FLAG_IS_EXTERNAL_TRANSLATION) != 0 {
			idx, err := readSignedIndex(p.reader, p.header.boneIndexSize)
			if err != nil {
				return wrapParseFailed("PMXボーン付与親インデックスの読み込みに失敗しました", err)
			}
			rate, err := p.reader.ReadFloat32()
			if err != nil {
				return wrapParseFailed("PMXボーン付与率の読み込みに失敗しました", err)
			}
			bone.EffectIndex = idx
			bone.EffectFactor = rate
		}

		if boneFlag&model.BONE_FLAG_HAS_FIXED_AXIS != 0 {
			axis, err := p.reader.ReadVec3()
			if err != nil {
				return wrapParseFailed("PMXボーン軸固定の読み込みに失敗しました", err)
			}
			bone.FixedAxis = axis
		}

		if boneFlag&model.BONE_FLAG_HAS_LOCAL_AXIS != 0 {
			axisX, err := p.reader.ReadVec3()
			if err != nil {
				return wrapParseFailed("PMXボーンローカル軸Xの読み込みに失敗しました", err)
			}
			axisZ, err := p.reader.ReadVec3()
			if err != nil {
				return wrapParseFailed("PMXボーンローカル軸Zの読み込みに失敗しました", err)
			}
			bone.LocalAxisX = axisX
			bone.LocalAxisZ = axisZ
		}

		if boneFlag&model.BONE_FLAG_IS_EXTERNAL_PARENT_DEFORM != 0 {
			key, err := p.reader.ReadInt32()
			if err != nil {
				return wrapParseFailed("PMXボーン外部親変形キーの読み込みに失敗しました", err)
			}
			bone.EffectorKey = int(key)
		}

		if boneFlag&model.BONE_FLAG_IS_IK != 0 {
			ik, err := p.readIk()
			if err != nil {
				return err
			}
			bone.Ik = ik
		}

		modelData.Bones.AppendRaw(bone)
	}
	return nil
}

// readIk はIK情報を読み込む。
func (p *pmxReader) readIk() (*model.Ik, error) {
	ik := &model.Ik{}
	idx, err := readSignedIndex(p.reader, p.header.boneIndexSize)
	if err != nil {
		return nil, wrapParseFailed("PMX IKターゲットの読み込みに失敗しました", err)
	}
	ik.BoneIndex = idx
	loopCount, err := p.reader.ReadInt32()
	if err != nil {
		return nil, wrapParseFailed("PMX IKループ数の読み込みに失敗しました", err)
	}
	ik.LoopCount = int(loopCount)
	unitRot, err := p.reader.ReadFloat32()
	if err != nil {
		return nil, wrapParseFailed("PMX IK単位回転の読み込みに失敗しました", err)
	}
	ik.UnitRotation = mmath.Vec3{Vec: r3.Vec{X: unitRot, Y: unitRot, Z: unitRot}}
	linkCount, err := p.reader.ReadInt32()
	if err != nil {
		return nil, wrapParseFailed("PMX IKリンク数の読み込みに失敗しました", err)
	}
	ik.Links = make([]model.IkLink, 0, linkCount)
	for i := 0; i < int(linkCount); i++ {
		linkIndex, err := readSignedIndex(p.reader, p.header.boneIndexSize)
		if err != nil {
			return nil, wrapParseFailed("PMX IKリンクの読み込みに失敗しました", err)
		}
		limitRaw, err := p.reader.ReadUint8()
		if err != nil {
			return nil, wrapParseFailed("PMX IK角度制限の読み込みに失敗しました", err)
		}
		link := model.IkLink{BoneIndex: linkIndex, AngleLimit: limitRaw == 1}
		if link.AngleLimit {
			minAngle, err := p.reader.ReadVec3()
			if err != nil {
				return nil, wrapParseFailed("PMX IK角度下限の読み込みに失敗しました", err)
			}
			maxAngle, err := p.reader.ReadVec3()
			if err != nil {
				return nil, wrapParseFailed("PMX IK角度上限の読み込みに失敗しました", err)
			}
			link.MinAngleLimit = minAngle
			link.MaxAngleLimit = maxAngle
		}
		ik.Links = append(ik.Links, link)
	}
	return ik, nil
}

// readMorphs はモーフセクションを読み込む。
func (p *pmxReader) readMorphs(modelData *model.PmxModel) error {
	count, err := p.reader.ReadInt32()
	if err != nil {
		return wrapParseFailed("PMXモーフ数の読み込みに失敗しました", err)
	}
	if count < 0 {
		return wrapParseFailed("PMXモーフ数が不正です", nil)
	}
	for i := 0; i < int(count); i++ {
		name, err := p.readText()
		if err != nil {
			return err
		}
		englishName, err := p.readText()
		if err != nil {
			return err
		}
		panel, err := p.reader.ReadUint8()
		if err != nil {
			return wrapParseFailed("PMXモーフパネルの読み込みに失敗しました", err)
		}
		morphType, err := p.reader.ReadUint8()
		if err != nil {
			return wrapParseFailed("PMXモーフ種別の読み込みに失敗しました", err)
		}
		offsetCount, err := p.reader.ReadInt32()
		if err != nil {
			return wrapParseFailed("PMXモーフオフセット数の読み込みに失敗しました", err)
		}

		morph := &model.Morph{
			EnglishName: englishName,
			Panel:       model.MorphPanel(panel),
			MorphType:   model.MorphType(morphType),
			Offsets:     make([]model.IMorphOffset, 0, offsetCount),
		}
		morph.SetName(name)

		for j := 0; j < int(offsetCount); j++ {
			offset, err := p.readMorphOffset(model.MorphType(morphType))
			if err != nil {
				return err
			}
			morph.Offsets = append(morph.Offsets, offset)
		}
		modelData.Morphs.AppendRaw(morph)
	}
	return nil
}

// readMorphOffset はモーフオフセットを読み込む。
func (p *pmxReader) readMorphOffset(morphType model.MorphType) (model.IMorphOffset, error) {
	switch morphType {
	case model.MORPH_TYPE_GROUP:
		idx, err := readSignedIndex(p.reader, p.header.morphIndexSize)
		if err != nil {
			return nil, wrapParseFailed("PMXグループモーフの読み込みに失敗しました", err)
		}
		rate, err := p.reader.ReadFloat32()
		if err != nil {
			return nil, wrapParseFailed("PMXグループモーフ率の読み込みに失敗しました", err)
		}
		return &model.GroupMorphOffset{MorphIndex: idx, MorphFactor: rate}, nil
	case model.MORPH_TYPE_VERTEX, model.MORPH_TYPE_AFTER_VERTEX:
		idx, err := readVertexIndex(p.reader, p.header.vertexIndexSize)
		if err != nil {
			return nil, wrapParseFailed("PMX頂点モーフの読み込みに失敗しました", err)
		}
		pos, err := p.reader.ReadVec3()
		if err != nil {
			return nil, wrapParseFailed("PMX頂点モーフ座標の読み込みに失敗しました", err)
		}
		return &model.VertexMorphOffset{VertexIndex: idx, Position: pos}, nil
	case model.MORPH_TYPE_BONE:
		idx, err := readSignedIndex(p.reader, p.header.boneIndexSize)
		if err != nil {
			return nil, wrapParseFailed("PMXボーンモーフの読み込みに失敗しました", err)
		}
		pos, err := p.reader.ReadVec3()
		if err != nil {
			return nil, wrapParseFailed("PMXボーンモーフ移動の読み込みに失敗しました", err)
		}
		rot, err := p.reader.ReadVec4()
		if err != nil {
			return nil, wrapParseFailed("PMXボーンモーフ回転の読み込みに失敗しました", err)
		}
		q := mmath.NewQuaternionByValues(rot.X, rot.Y, rot.Z, rot.W)
		return &model.BoneMorphOffset{BoneIndex: idx, Position: pos, Rotation: q}, nil
	case model.MORPH_TYPE_UV,
		model.MORPH_TYPE_EXTENDED_UV1,
		model.MORPH_TYPE_EXTENDED_UV2,
		model.MORPH_TYPE_EXTENDED_UV3,
		model.MORPH_TYPE_EXTENDED_UV4:
		idx, err := readVertexIndex(p.reader, p.header.vertexIndexSize)
		if err != nil {
			return nil, wrapParseFailed("PMX UVモーフの読み込みに失敗しました", err)
		}
		uv, err := p.reader.ReadVec4()
		if err != nil {
			return nil, wrapParseFailed("PMX UVモーフの読み込みに失敗しました", err)
		}
		return &model.UvMorphOffset{VertexIndex: idx, Uv: uv, UvType: morphType}, nil
	case model.MORPH_TYPE_MATERIAL:
		idx, err := readSignedIndex(p.reader, p.header.materialIndexSize)
		if err != nil {
			return nil, wrapParseFailed("PMX材質モーフの読み込みに失敗しました", err)
		}
		modeRaw, err := p.reader.ReadUint8()
		if err != nil {
			return nil, wrapParseFailed("PMX材質モーフ計算モードの読み込みに失敗しました", err)
		}
		diffuse, err := p.reader.ReadVec4()
		if err != nil {
			return nil, wrapParseFailed("PMX材質モーフDiffuseの読み込みに失敗しました", err)
		}
		specular, err := p.reader.ReadVec3()
		if err != nil {
			return nil, wrapParseFailed("PMX材質モーフSpecularの読み込みに失敗しました", err)
		}
		specularPower, err := p.reader.ReadFloat32()
		if err != nil {
			return nil, wrapParseFailed("PMX材質モーフSpecularPowerの読み込みに失敗しました", err)
		}
		ambient, err := p.reader.ReadVec3()
		if err != nil {
			return nil, wrapParseFailed("PMX材質モーフAmbientの読み込みに失敗しました", err)
		}
		edge, err := p.reader.ReadVec4()
		if err != nil {
			return nil, wrapParseFailed("PMX材質モーフエッジ色の読み込みに失敗しました", err)
		}
		edgeSize, err := p.reader.ReadFloat32()
		if err != nil {
			return nil, wrapParseFailed("PMX材質モーフエッジサイズの読み込みに失敗しました", err)
		}
		textureFactor, err := p.reader.ReadVec4()
		if err != nil {
			return nil, wrapParseFailed("PMX材質モーフテクスチャ係数の読み込みに失敗しました", err)
		}
		sphereFactor, err := p.reader.ReadVec4()
		if err != nil {
			return nil, wrapParseFailed("PMX材質モーフスフィア係数の読み込みに失敗しました", err)
		}
		toonFactor, err := p.reader.ReadVec4()
		if err != nil {
			return nil, wrapParseFailed("PMX材質モーフトゥーン係数の読み込みに失敗しました", err)
		}
		return &model.MaterialMorphOffset{
			MaterialIndex:       idx,
			CalcMode:            model.MaterialMorphCalcMode(modeRaw),
			Diffuse:             diffuse,
			Specular:            mmath.Vec4{X: specular.X, Y: specular.Y, Z: specular.Z, W: specularPower},
			Ambient:             ambient,
			Edge:                edge,
			EdgeSize:            edgeSize,
			TextureFactor:       textureFactor,
			SphereTextureFactor: sphereFactor,
			ToonTextureFactor:   toonFactor,
		}, nil
	default:
		return nil, wrapFormatNotSupported("PMXモーフ種別が未対応です", nil)
	}
}

// readDisplaySlots は表示枠セクションを読み込む。
func (p *pmxReader) readDisplaySlots(modelData *model.PmxModel) error {
	count, err := p.reader.ReadInt32()
	if err != nil {
		return wrapParseFailed("PMX表示枠数の読み込みに失敗しました", err)
	}
	if count < 0 {
		return wrapParseFailed("PMX表示枠数が不正です", nil)
	}
	for i := 0; i < int(count); i++ {
		name, err := p.readText()
		if err != nil {
			return err
		}
		englishName, err := p.readText()
		if err != nil {
			return err
		}
		specialFlag, err := p.reader.ReadUint8()
		if err != nil {
			return wrapParseFailed("PMX表示枠フラグの読み込みに失敗しました", err)
		}
		referenceCount, err := p.reader.ReadInt32()
		if err != nil {
			return wrapParseFailed("PMX表示枠参照数の読み込みに失敗しました", err)
		}
		slot := &model.DisplaySlot{
			EnglishName: englishName,
			SpecialFlag: model.SpecialFlag(specialFlag),
			References:  make([]model.Reference, 0, referenceCount),
		}
		slot.SetName(name)
		for j := 0; j < int(referenceCount); j++ {
			displayType, err := p.reader.ReadUint8()
			if err != nil {
				return wrapParseFailed("PMX表示枠参照種別の読み込みに失敗しました", err)
			}
			ref := model.Reference{DisplayType: model.DisplayType(displayType)}
			switch ref.DisplayType {
			case model.DISPLAY_TYPE_BONE:
				idx, err := readSignedIndex(p.reader, p.header.boneIndexSize)
				if err != nil {
					return wrapParseFailed("PMX表示枠ボーン参照の読み込みに失敗しました", err)
				}
				ref.DisplayIndex = idx
			case model.DISPLAY_TYPE_MORPH:
				idx, err := readSignedIndex(p.reader, p.header.morphIndexSize)
				if err != nil {
					return wrapParseFailed("PMX表示枠モーフ参照の読み込みに失敗しました", err)
				}
				ref.DisplayIndex = idx
			default:
				return wrapFormatNotSupported("PMX表示枠参照種別が未対応です", nil)
			}
			slot.References = append(slot.References, ref)
		}
		modelData.DisplaySlots.AppendRaw(slot)
	}
	return nil
}

// readRigidBodies は剛体セクションを読み込む。
func (p *pmxReader) readRigidBodies(modelData *model.PmxModel) error {
	count, err := p.reader.ReadInt32()
	if err != nil {
		return wrapParseFailed("PMX剛体数の読み込みに失敗しました", err)
	}
	if count < 0 {
		return wrapParseFailed("PMX剛体数が不正です", nil)
	}
	for i := 0; i < int(count); i++ {
		name, err := p.readText()
		if err != nil {
			return err
		}
		englishName, err := p.readText()
		if err != nil {
			return err
		}
		boneIndex, err := readSignedIndex(p.reader, p.header.boneIndexSize)
		if err != nil {
			return wrapParseFailed("PMX剛体ボーン参照の読み込みに失敗しました", err)
		}
		group, err := p.reader.ReadUint8()
		if err != nil {
			return wrapParseFailed("PMX剛体衝突グループの読み込みに失敗しました", err)
		}
		mask, err := p.reader.ReadUint16()
		if err != nil {
			return wrapParseFailed("PMX剛体衝突マスクの読み込みに失敗しました", err)
		}
		shape, err := p.reader.ReadUint8()
		if err != nil {
			return wrapParseFailed("PMX剛体形状の読み込みに失敗しました", err)
		}
		size, err := p.reader.ReadVec3()
		if err != nil {
			return wrapParseFailed("PMX剛体サイズの読み込みに失敗しました", err)
		}
		pos, err := p.reader.ReadVec3()
		if err != nil {
			return wrapParseFailed("PMX剛体位置の読み込みに失敗しました", err)
		}
		rot, err := p.reader.ReadVec3()
		if err != nil {
			return wrapParseFailed("PMX剛体回転の読み込みに失敗しました", err)
		}
		mass, err := p.reader.ReadFloat32()
		if err != nil {
			return wrapParseFailed("PMX剛体質量の読み込みに失敗しました", err)
		}
		linearDamping, err := p.reader.ReadFloat32()
		if err != nil {
			return wrapParseFailed("PMX剛体移動減衰の読み込みに失敗しました", err)
		}
		angularDamping, err := p.reader.ReadFloat32()
		if err != nil {
			return wrapParseFailed("PMX剛体回転減衰の読み込みに失敗しました", err)
		}
		restitution, err := p.reader.ReadFloat32()
		if err != nil {
			return wrapParseFailed("PMX剛体反発力の読み込みに失敗しました", err)
		}
		friction, err := p.reader.ReadFloat32()
		if err != nil {
			return wrapParseFailed("PMX剛体摩擦の読み込みに失敗しました", err)
		}
		physicsType, err := p.reader.ReadUint8()
		if err != nil {
			return wrapParseFailed("PMX剛体物理種別の読み込みに失敗しました", err)
		}

		rigid := &model.RigidBody{
			EnglishName: englishName,
			BoneIndex:   boneIndex,
			CollisionGroup: model.CollisionGroup{
				Group: group,
				Mask:  mask,
			},
			Shape:       model.Shape(shape),
			Size:        size,
			Position:    pos,
			Rotation:    rot,
			Param:       model.RigidBodyParam{Mass: mass, LinearDamping: linearDamping, AngularDamping: angularDamping, Restitution: restitution, Friction: friction},
			PhysicsType: model.PhysicsType(physicsType),
		}
		rigid.SetName(name)
		modelData.RigidBodies.AppendRaw(rigid)
	}
	return nil
}

// readJoints はジョイントセクションを読み込む。
func (p *pmxReader) readJoints(modelData *model.PmxModel) error {
	count, err := p.reader.ReadInt32()
	if err != nil {
		return wrapParseFailed("PMXジョイント数の読み込みに失敗しました", err)
	}
	if count < 0 {
		return wrapParseFailed("PMXジョイント数が不正です", nil)
	}
	for i := 0; i < int(count); i++ {
		name, err := p.readText()
		if err != nil {
			return err
		}
		englishName, err := p.readText()
		if err != nil {
			return err
		}
		jointType, err := p.reader.ReadUint8()
		if err != nil {
			return wrapParseFailed("PMXジョイント種別の読み込みに失敗しました", err)
		}
		if jointType != 0 {
			return wrapFormatNotSupported("PMXジョイント種別が未対応です", nil)
		}
		idxA, err := readSignedIndex(p.reader, p.header.rigidBodyIndexSize)
		if err != nil {
			return wrapParseFailed("PMXジョイント剛体A参照の読み込みに失敗しました", err)
		}
		idxB, err := readSignedIndex(p.reader, p.header.rigidBodyIndexSize)
		if err != nil {
			return wrapParseFailed("PMXジョイント剛体B参照の読み込みに失敗しました", err)
		}
		pos, err := p.reader.ReadVec3()
		if err != nil {
			return wrapParseFailed("PMXジョイント位置の読み込みに失敗しました", err)
		}
		rot, err := p.reader.ReadVec3()
		if err != nil {
			return wrapParseFailed("PMXジョイント回転の読み込みに失敗しました", err)
		}
		transMin, err := p.reader.ReadVec3()
		if err != nil {
			return wrapParseFailed("PMXジョイント移動下限の読み込みに失敗しました", err)
		}
		transMax, err := p.reader.ReadVec3()
		if err != nil {
			return wrapParseFailed("PMXジョイント移動上限の読み込みに失敗しました", err)
		}
		rotMin, err := p.reader.ReadVec3()
		if err != nil {
			return wrapParseFailed("PMXジョイント回転下限の読み込みに失敗しました", err)
		}
		rotMax, err := p.reader.ReadVec3()
		if err != nil {
			return wrapParseFailed("PMXジョイント回転上限の読み込みに失敗しました", err)
		}
		constTrans, err := p.reader.ReadVec3()
		if err != nil {
			return wrapParseFailed("PMXジョイント移動ばねの読み込みに失敗しました", err)
		}
		constRot, err := p.reader.ReadVec3()
		if err != nil {
			return wrapParseFailed("PMXジョイント回転ばねの読み込みに失敗しました", err)
		}

		joint := &model.Joint{
			EnglishName:     englishName,
			RigidBodyIndexA: idxA,
			RigidBodyIndexB: idxB,
			Param: model.JointParam{
				Position:                  pos,
				Rotation:                  rot,
				TranslationLimitMin:       transMin,
				TranslationLimitMax:       transMax,
				RotationLimitMin:          rotMin,
				RotationLimitMax:          rotMax,
				SpringConstantTranslation: constTrans,
				SpringConstantRotation:    constRot,
			},
		}
		joint.SetName(name)
		modelData.Joints.AppendRaw(joint)
	}
	return nil
}

// skipSoftBodies はSoftBodyセクションを読み飛ばす。
func (p *pmxReader) skipSoftBodies() error {
	if p.header == nil {
		return nil
	}
	if !nearVersion(p.header.version, 2.1) {
		return nil
	}
	count, err := p.reader.ReadInt32()
	if err != nil {
		return wrapParseFailed("PMXソフトボディ数の読み込みに失敗しました", err)
	}
	if count < 0 {
		return wrapParseFailed("PMXソフトボディ数が不正です", nil)
	}
	if count == 0 {
		return nil
	}
	if err := p.reader.DiscardAll(); err != nil {
		return wrapParseFailed("PMXソフトボディの読み飛ばしに失敗しました", err)
	}
	return nil
}

// resolveEncoding はPMXのエンコード方式を解決する。
func resolveEncoding(encodeType byte) (encoding.Encoding, error) {
	switch encodeType {
	case 0:
		return unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM), nil
	case 1:
		return unicode.UTF8, nil
	default:
		return nil, io_common.NewIoFormatNotSupported("エンコード種別が未対応です: %d", nil, encodeType)
	}
}

// nearVersion はバージョンの近似一致を判定する。
func nearVersion(value, expected float64) bool {
	return math.Abs(value-expected) < 0.01
}

// isValidIndexSize はインデックスサイズの妥当性を判定する。
func isValidIndexSize(size byte) bool {
	return size == 1 || size == 2 || size == 4
}

// wrapParseFailed は解析失敗エラーを生成する。
func wrapParseFailed(message string, err error) error {
	return io_common.NewIoParseFailed(message, err)
}

// wrapFormatNotSupported は形式未対応エラーを生成する。
func wrapFormatNotSupported(message string, err error) error {
	return io_common.NewIoFormatNotSupported(message, err)
}

// wrapEncodingUnknown は未知エンコードエラーを生成する。
func wrapEncodingUnknown(message string, err error) error {
	return io_common.NewIoEncodingUnknown(message, err)
}
