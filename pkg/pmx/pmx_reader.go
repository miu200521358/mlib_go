package pmx

import (
	"fmt"
	"slices"

	"golang.org/x/text/encoding/unicode"

	"github.com/miu200521358/mlib_go/pkg/mcore"
	"github.com/miu200521358/mlib_go/pkg/mmath"
	"github.com/miu200521358/mlib_go/pkg/mutils/mlog"
)

type PmxReader struct {
	mcore.BaseReader[*PmxModel]
}

func (r *PmxReader) createModel(path string) *PmxModel {
	model := NewPmxModel(path)
	return model
}

// 指定されたパスのファイルからデータを読み込む
func (r *PmxReader) ReadByFilepath(path string) (mcore.IHashModel, error) {
	// モデルを新規作成
	model := r.createModel(path)

	hash, err := r.ReadHashByFilePath(path)
	if err != nil {
		mlog.E("ReadByFilepath.ReadHashByFilePath error: %v", err)
		return nil, err
	}
	model.Hash = hash

	// ファイルを開く
	err = r.Open(path)
	if err != nil {
		mlog.E("ReadByFilepath.Open error: %v", err)
		return model, err
	}

	err = r.readHeader(model)
	if err != nil {
		mlog.E("ReadByFilepath.readHeader error: %v", err)
		return model, err
	}

	err = r.readData(model)
	if err != nil {
		mlog.E("ReadByFilepath.readData error: %v", err)
		return model, err
	}

	r.Close()
	model.setup()

	return model, nil
}

func (r *PmxReader) ReadNameByFilepath(path string) (string, error) {
	// モデルを新規作成
	model := r.createModel(path)

	// ファイルを開く
	err := r.Open(path)
	if err != nil {
		mlog.E("ReadNameByFilepath.Open error: %v", err)
		return "", err
	}

	err = r.readHeader(model)
	if err != nil {
		mlog.E("ReadNameByFilepath.readHeader error: %v", err)
		return "", err
	}

	r.Close()

	return model.Name, nil
}

func (r *PmxReader) readHeader(model *PmxModel) error {
	fbytes, err := r.UnpackBytes(4)
	if err != nil {
		mlog.E("readHeader.UnpackBytes error: %v", err)
		return err
	}
	model.Signature = r.DecodeText(unicode.UTF8, fbytes)
	model.Version, err = r.UnpackFloat()

	if err != nil {
		mlog.E("readHeader.Version error: %v", err)
		return err
	}

	if model.Signature[:3] != "PMX" ||
		!slices.Contains([]string{"2.0", "2.1"}, fmt.Sprintf("%.1f", model.Version)) {
		// 整合性チェック
		return fmt.Errorf("PMX2.0/2.1形式外のデータです。signature: %s, version: %.1f ", model.Signature, model.Version)
	}

	// 1 : byte	| 後続するデータ列のバイトサイズ  PMX2.0は 8 で固定
	_, _ = r.UnpackByte()

	// [0] - エンコード方式  | 0:UTF16 1:UTF8
	encodeType, err := r.UnpackByte()
	if err != nil {
		mlog.E("UnpackByte encodeType error: %v", err)
		return err
	}

	switch encodeType {
	case 0:
		r.DefineEncoding(unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM))
	case 1:
		r.DefineEncoding(unicode.UTF8)
	default:
		return fmt.Errorf("未知のエンコードタイプです。encodeType: %d", encodeType)
	}

	// [1] - 追加UV数 	| 0～4 詳細は頂点参照
	extendedUVCount, err := r.UnpackByte()
	if err != nil {
		mlog.E("UnpackByte extendedUVCount error: %v", err)
		return err
	}
	model.ExtendedUVCountType = int(extendedUVCount)
	// [2] - 頂点Indexサイズ | 1,2,4 のいずれか
	vertexCount, err := r.UnpackByte()
	if err != nil {
		mlog.E("UnpackByte vertexCount error: %v", err)
		return err
	}
	model.VertexCountType = int(vertexCount)
	// [3] - テクスチャIndexサイズ | 1,2,4 のいずれか
	textureCount, err := r.UnpackByte()
	if err != nil {
		mlog.E("UnpackByte textureCount error: %v", err)
		return err
	}
	model.TextureCountType = int(textureCount)
	// [4] - 材質Indexサイズ | 1,2,4 のいずれか
	materialCount, err := r.UnpackByte()
	if err != nil {
		mlog.E("UnpackByte materialCount error: %v", err)
		return err
	}
	model.MaterialCountType = int(materialCount)
	// [5] - ボーンIndexサイズ | 1,2,4 のいずれか
	boneCount, err := r.UnpackByte()
	if err != nil {
		mlog.E("UnpackByte boneCount error: %v", err)
		return err
	}
	model.BoneCountType = int(boneCount)
	// [6] - モーフIndexサイズ | 1,2,4 のいずれか
	morphCount, err := r.UnpackByte()
	if err != nil {
		mlog.E("UnpackByte morphCount error: %v", err)
		return err
	}
	model.MorphCountType = int(morphCount)
	// [7] - 剛体Indexサイズ | 1,2,4 のいずれか
	rigidBodyCount, err := r.UnpackByte()
	if err != nil {
		mlog.E("UnpackByte rigidBodyCount error: %v", err)
		return err
	}
	model.RigidBodyCountType = int(rigidBodyCount)

	// 4 + n : TextBuf	| モデル名
	model.Name = r.ReadText()
	// 4 + n : TextBuf	| モデル名英
	model.EnglishName = r.ReadText()
	// 4 + n : TextBuf	| コメント
	model.Comment = r.ReadText()
	// 4 + n : TextBuf	| コメント英
	model.EnglishComment = r.ReadText()

	return nil
}

func (r *PmxReader) readData(model *PmxModel) error {
	err := r.readVertices(model)
	if err != nil {
		mlog.E("readData.readVertices error: %v", err)
		return err
	}

	err = r.readFaces(model)
	if err != nil {
		mlog.E("readData.readFaces error: %v", err)
		return err
	}

	err = r.readTextures(model)
	if err != nil {
		mlog.E("readData.readTextures error: %v", err)
		return err
	}

	err = r.readMaterials(model)
	if err != nil {
		mlog.E("readData.readMaterials error: %v", err)
		return err
	}

	err = r.readBones(model)
	if err != nil {
		mlog.E("readData.readBones error: %v", err)
		return err
	}

	err = r.readMorphs(model)
	if err != nil {
		mlog.E("readData.readMorphs error: %v", err)
		return err
	}

	err = r.readDisplaySlots(model)
	if err != nil {
		mlog.E("readData.readDisplaySlots error: %v", err)
		return err
	}

	err = r.readRigidBodies(model)
	if err != nil {
		mlog.E("readData.readRigidBodies error: %v", err)
		return err
	}

	err = r.readJoints(model)
	if err != nil {
		mlog.E("readData.readJoints error: %v", err)
		return err
	}

	return nil
}

func (r *PmxReader) readVertices(model *PmxModel) error {
	totalVertexCount, err := r.UnpackInt()
	if err != nil {
		mlog.E("readVertices UnpackInt totalVertexCount error: %v", err)
		return err
	}

	for i := 0; i < totalVertexCount; i++ {
		v := Vertex{IndexModel: &mcore.IndexModel{Index: i}}

		// 12 : float3  | 位置(x,y,z)
		pos, err := r.UnpackVec3(true)
		if err != nil {
			mlog.E("[%d] readVertices UnpackFloat Position error: %v", i, err)
			return err
		}
		v.Position = &pos

		// 12 : float3  | 法線(x,y,z)
		normal, err := r.UnpackVec3(true)
		if err != nil {
			mlog.E("[%d] readVertices UnpackFloat Normal[0] error: %v", i, err)
			return err
		}
		v.Normal = &normal

		// 8  : float2  | UV(u,v)
		uv, err := r.UnpackVec2()
		if err != nil {
			mlog.E("[%d] readVertices UnpackFloat UV[0] error: %v", i, err)
			return err
		}
		v.UV = &uv

		// 16 * n : float4[n] | 追加UV(x,y,z,w)  PMXヘッダの追加UV数による
		v.ExtendedUVs = make([]*mmath.MVec4, 0)
		for j := 0; j < model.ExtendedUVCountType; j++ {
			extendedUV, err := r.UnpackVec4(false)
			if err != nil {
				mlog.E("[%d][%d] readVertices UnpackVec4 ExtendedUV error: %v", i, j, err)
				return err
			}
			v.ExtendedUVs = append(v.ExtendedUVs, &extendedUV)
		}

		// 1 : byte    | ウェイト変形方式 0:BDEF1 1:BDEF2 2:BDEF4 3:SDEF
		Type, err := r.UnpackByte()
		if err != nil {
			mlog.E("[%d] readVertices UnpackByte Type error: %v", i, err)
			return err
		}
		v.DeformType = DeformType(Type)

		switch v.DeformType {
		case BDEF1:
			// n : ボーンIndexサイズ  | ウェイト1.0の単一ボーン(参照Index)
			boneIndex, err := r.unpackBoneIndex(model)
			if err != nil {
				mlog.E("[%d] readVertices BDEF1 unpackBoneIndex error: %v", i, err)
				return err
			}
			deform := NewBdef1(boneIndex)
			v.Deform = &deform
		case BDEF2:
			// n : ボーンIndexサイズ  | ボーン1の参照Index
			boneIndex1, err := r.unpackBoneIndex(model)
			if err != nil {
				mlog.E("[%d] readVertices BDEF2 unpackBoneIndex boneIndex1 error: %v", i, err)
				return err
			}
			// n : ボーンIndexサイズ  | ボーン2の参照Index
			boneIndex2, err := r.unpackBoneIndex(model)
			if err != nil {
				mlog.E("[%d] readVertices BDEF2 unpackBoneIndex boneIndex2 error: %v", i, err)
				return err
			}
			// 4 : float              | ボーン1のウェイト値(0～1.0), ボーン2のウェイト値は 1.0-ボーン1ウェイト
			boneWeight, err := r.UnpackFloat()
			if err != nil {
				mlog.E("[%d] readVertices BDEF2 UnpackFloat boneWeight error: %v", i, err)
				return err
			}
			deform := NewBdef2(boneIndex1, boneIndex2, boneWeight)
			v.Deform = &deform
		case BDEF4:
			// n : ボーンIndexサイズ  | ボーン1の参照Index
			boneIndex1, err := r.unpackBoneIndex(model)
			if err != nil {
				mlog.E("[%d] readVertices BDEF4 unpackBoneIndex boneIndex1 error: %v", i, err)
				return err
			}
			// n : ボーンIndexサイズ  | ボーン2の参照Index
			boneIndex2, err := r.unpackBoneIndex(model)
			if err != nil {
				mlog.E("[%d] readVertices BDEF4 unpackBoneIndex boneIndex2 error: %v", i, err)
				return err
			}
			// n : ボーンIndexサイズ  | ボーン3の参照Index
			boneIndex3, err := r.unpackBoneIndex(model)
			if err != nil {
				mlog.E("[%d] readVertices BDEF4 unpackBoneIndex boneIndex3 error: %v", i, err)
				return err
			}
			// n : ボーンIndexサイズ  | ボーン4の参照Index
			boneIndex4, err := r.unpackBoneIndex(model)
			if err != nil {
				mlog.E("[%d] readVertices BDEF4 unpackBoneIndex boneIndex4 error: %v", i, err)
				return err
			}
			// 4 : float              | ボーン1のウェイト値
			boneWeight1, err := r.UnpackFloat()
			if err != nil {
				mlog.E("[%d] readVertices BDEF4 UnpackFloat boneWeight1 error: %v", i, err)
				return err
			}
			// 4 : float              | ボーン2のウェイト値
			boneWeight2, err := r.UnpackFloat()
			if err != nil {
				mlog.E("[%d] readVertices BDEF4 UnpackFloat boneWeight2 error: %v", i, err)
				return err
			}
			// 4 : float              | ボーン3のウェイト値
			boneWeight3, err := r.UnpackFloat()
			if err != nil {
				mlog.E("[%d] readVertices BDEF4 UnpackFloat boneWeight3 error: %v", i, err)
				return err
			}
			// 4 : float              | ボーン4のウェイト値 (ウェイト計1.0の保障はない)
			boneWeight4, err := r.UnpackFloat()
			if err != nil {
				mlog.E("[%d] readVertices BDEF4 UnpackFloat boneWeight4 error: %v", i, err)
				return err
			}
			deform := NewBdef4(boneIndex1, boneIndex2, boneIndex3, boneIndex4,
				boneWeight1, boneWeight2, boneWeight3, boneWeight4)
			v.Deform = &deform
		case SDEF:
			// n : ボーンIndexサイズ  | ボーン1の参照Index
			boneIndex1, err := r.unpackBoneIndex(model)
			if err != nil {
				mlog.E("[%d] readVertices SDEF unpackBoneIndex boneIndex1 error: %v", i, err)
				return err
			}
			// n : ボーンIndexサイズ  | ボーン2の参照Index
			boneIndex2, err := r.unpackBoneIndex(model)
			if err != nil {
				mlog.E("[%d] readVertices SDEF unpackBoneIndex boneIndex2 error: %v", i, err)
				return err
			}
			// 4 : float              | ボーン1のウェイト値(0～1.0), ボーン2のウェイト値は 1.0-ボーン1ウェイト
			boneWeight, err := r.UnpackFloat()
			if err != nil {
				mlog.E("[%d] readVertices SDEF UnpackFloat boneWeight error: %v", i, err)
				return err
			}
			// 12 : float3             | SDEF-C値(x,y,z)
			sdefC, err := r.UnpackVec3(true)
			if err != nil {
				mlog.E("[%d] readVertices SDEF UnpackVec3 sdefC error: %v", i, err)
				return err
			}
			// 12 : float3             | SDEF-R0値(x,y,z)
			sdefR0, err := r.UnpackVec3(true)
			if err != nil {
				mlog.E("[%d] readVertices SDEF UnpackVec3 sdefR0 error: %v", i, err)
				return err
			}
			// 12 : float3             | SDEF-R1値(x,y,z) ※修正値を要計算
			sdefR1, err := r.UnpackVec3(true)
			if err != nil {
				mlog.E("[%d] readVertices SDEF UnpackVec3 sdefR1 error: %v", i, err)
				return err
			}
			deform := NewSdef(boneIndex1, boneIndex2, boneWeight, &sdefC, &sdefR0, &sdefR1)
			v.Deform = &deform
		}

		v.EdgeFactor, err = r.UnpackFloat()
		if err != nil {
			mlog.E("[%d] readVertices UnpackFloat EdgeFactor error: %v", i, err)
			return err
		}

		model.Vertices.Append(&v)
	}

	return nil
}

func (r *PmxReader) readFaces(model *PmxModel) error {
	totalFaceCount, err := r.UnpackInt()
	if err != nil {
		mlog.E("readFaces UnpackInt totalFaceCount error: %v", err)
		return err
	}

	for i := 0; i < totalFaceCount; i += 3 {
		f := Face{
			IndexModel:    &mcore.IndexModel{Index: int(i / 3)},
			VertexIndexes: [3]int{},
		}

		// n : 頂点Indexサイズ     | 頂点の参照Index
		f.VertexIndexes[0], err = r.unpackVertexIndex(model)
		if err != nil {
			mlog.E("[%d] readFaces unpackVertexIndex VertexIndexes[0] error: %v", i, err)
			return err
		}

		// n : 頂点Indexサイズ     | 頂点の参照Index
		f.VertexIndexes[1], err = r.unpackVertexIndex(model)
		if err != nil {
			mlog.E("[%d] readFaces unpackVertexIndex VertexIndexes[1] error: %v", i, err)
			return err
		}

		// n : 頂点Indexサイズ     | 頂点の参照Index
		f.VertexIndexes[2], err = r.unpackVertexIndex(model)
		if err != nil {
			mlog.E("[%d] readFaces unpackVertexIndex VertexIndexes[2] error: %v", i, err)
			return err
		}

		model.Faces.Append(&f)
	}

	return nil
}

func (r *PmxReader) readTextures(model *PmxModel) error {
	totalTextureCount, err := r.UnpackInt()
	if err != nil {
		mlog.E("readTextures UnpackInt totalTextureCount error: %v", err)
		return err
	}

	for i := 0; i < totalTextureCount; i++ {
		t := NewTexture()

		// 4 + n : TextBuf	| テクスチャパス
		t.Name = r.ReadText()

		model.Textures.Append(t)
	}

	return nil
}

func (r *PmxReader) readMaterials(model *PmxModel) error {
	totalMaterialCount, err := r.UnpackInt()
	if err != nil {
		mlog.E("readMaterials UnpackInt totalMaterialCount error: %v", err)
		return err
	}

	for i := 0; i < totalMaterialCount; i++ {
		// 4 + n : TextBuf	| 材質名
		name := r.ReadText()
		// 4 + n : TextBuf	| 材質名英
		englishName := r.ReadText()

		m := Material{
			IndexNameModel: &mcore.IndexNameModel{Index: i, Name: name, EnglishName: englishName},
		}

		// 16 : float4	| Diffuse (R,G,B,A)
		diffuse, err := r.UnpackVec4(false)
		if err != nil {
			mlog.E("[%d] readMaterials UnpackVec4 Diffuse error: %v", i, err)
			return err
		}
		m.Diffuse = &diffuse
		// 12 : float3	| Specular (R,G,B,Specular係数)
		specular, err := r.UnpackVec4(false)
		if err != nil {
			mlog.E("[%d] readMaterials UnpackVec4 Specular error: %v", i, err)
			return err
		}
		m.Specular = &specular
		// 12 : float3	| Ambient (R,G,B)
		ambient, err := r.UnpackVec3(false)
		if err != nil {
			mlog.E("[%d] readMaterials UnpackVec3 Ambient error: %v", i, err)
			return err
		}
		m.Ambient = &ambient
		// 1  : bitFlag  	| 描画フラグ(8bit) - 各bit 0:OFF 1:ON
		drawFlag, err := r.UnpackByte()
		if err != nil {
			mlog.E("[%d] readMaterials UnpackByte DrawFlag error: %v", i, err)
			return err
		}
		m.DrawFlag = DrawFlag(drawFlag)
		// 16 : float4	| エッジ色 (R,G,B,A)
		edge, err := r.UnpackVec4(false)
		if err != nil {
			mlog.E("[%d] readMaterials UnpackVec4 Edge error: %v", i, err)
			return err
		}
		m.Edge = &edge
		// 4  : float	| エッジサイズ
		m.EdgeSize, err = r.UnpackFloat()
		if err != nil {
			mlog.E("[%d] readMaterials UnpackFloat EdgeSize error: %v", i, err)
			return err
		}
		// n  : テクスチャIndexサイズ	| 通常テクスチャ
		m.TextureIndex, err = r.unpackTextureIndex(model)
		if err != nil {
			mlog.E("[%d] readMaterials unpackTextureIndex TextureIndex error: %v", i, err)
			return err
		}
		// n  : テクスチャIndexサイズ	| スフィアテクスチャ
		m.SphereTextureIndex, err = r.unpackTextureIndex(model)
		if err != nil {
			mlog.E("[%d] readMaterials unpackTextureIndex SphereTextureIndex error: %v", i, err)
			return err
		}
		// 1  : byte	| スフィアモード 0:無効 1:乗算(sph) 2:加算(spa) 3:サブテクスチャ(追加UV1のx,yをUV参照して通常テクスチャ描画を行う)
		sphereMode, err := r.UnpackByte()
		if err != nil {
			mlog.E("[%d] readMaterials UnpackByte SphereMode error: %v", i, err)
			return err
		}
		m.SphereMode = SphereMode(sphereMode)
		// 1  : byte	| 共有Toonフラグ 0:継続値は個別Toon 1:継続値は共有Toon
		toonSharingFlag, err := r.UnpackByte()

		if err != nil {
			mlog.E("[%d] readMaterials UnpackByte ToonSharingFlag error: %v", i, err)
			return err
		}
		m.ToonSharingFlag = ToonSharing(toonSharingFlag)

		switch m.ToonSharingFlag {
		case TOON_SHARING_INDIVIDUAL:
			// n  : テクスチャIndexサイズ	| Toonテクスチャ
			m.ToonTextureIndex, err = r.unpackTextureIndex(model)
			if err != nil {
				mlog.E("[%d] readMaterials unpackTextureIndex ToonTextureIndex error: %v", i, err)
				return err
			}
		case TOON_SHARING_SHARING:
			// 1  : byte	| 共有ToonテクスチャIndex 0～9
			toonTextureIndex, err := r.UnpackByte()
			if err != nil {
				mlog.E("[%d] readMaterials UnpackByte ToonTextureIndex error: %v", i, err)
				return err
			}
			m.ToonTextureIndex = int(toonTextureIndex)
		}

		// 4 + n : TextBuf	| メモ
		m.Memo = r.ReadText()

		// 4  : int	| 材質に対応する面(頂点)数 (必ず3の倍数になる)
		m.VerticesCount, err = r.UnpackInt()
		if err != nil {
			mlog.E("[%d] readMaterials UnpackInt VerticesCount error: %v", i, err)
			return err
		}

		model.Materials.Append(&m)
	}

	return nil
}

func (r *PmxReader) readBones(model *PmxModel) error {
	totalBoneCount, err := r.UnpackInt()
	if err != nil {
		mlog.E("readBones UnpackInt totalBoneCount error: %v", err)
		return err
	}

	for i := 0; i < totalBoneCount; i++ {

		// 4 + n : TextBuf	| ボーン名
		name := r.ReadText()
		// 4 + n : TextBuf	| ボーン名英
		englishName := r.ReadText()

		b := Bone{
			IndexNameModel:         &mcore.IndexNameModel{Index: i, Name: name, EnglishName: englishName},
			IkLinkBoneIndexes:      make([]int, 0),
			IkTargetBoneIndexes:    make([]int, 0),
			ParentRelativePosition: mmath.NewMVec3(),
			ChildRelativePosition:  mmath.NewMVec3(),
			NormalizedFixedAxis:    mmath.NewMVec3(),
			TreeBoneIndexes:        make([]int, 0),
			RevertOffsetMatrix:     mmath.NewMMat4(),
			OffsetMatrix:           mmath.NewMMat4(),
			ParentBoneIndexes:      make([]int, 0),
			RelativeBoneIndexes:    make([]int, 0),
			ChildBoneIndexes:       make([]int, 0),
			EffectiveBoneIndexes:   make([]int, 0),
			AngleLimit:             false,
			MinAngleLimit:          mmath.NewRotation(),
			MaxAngleLimit:          mmath.NewRotation(),
			LocalAngleLimit:        false,
			LocalMinAngleLimit:     mmath.NewRotation(),
			LocalMaxAngleLimit:     mmath.NewRotation(),
			AxisSign:               1,
		}

		// 12 : float3	| 位置
		pos, err := r.UnpackVec3(true)
		if err != nil {
			mlog.E("[%d] readBones UnpackVec3 Position error: %v", i, err)
			return err
		}
		b.Position = &pos
		// n  : ボーンIndexサイズ	| 親ボーン
		b.ParentIndex, err = r.unpackBoneIndex(model)
		if err != nil {
			mlog.E("[%d] readBones unpackBoneIndex ParentIndex error: %v", i, err)
			return err
		}
		// 4  : int	| 変形階層
		b.Layer, err = r.UnpackInt()
		if err != nil {
			mlog.E("[%d] UnpackInt Layer error: %v", i, err)
			return err
		}
		// 2  : bitFlag*2	| ボーンフラグ(16bit) 各bit 0:OFF 1:ON
		boneFlag, err := r.UnpackBytes(2)
		if err != nil {
			mlog.E("[%d] readBones UnpackBytes BoneFlag error: %v", i, err)
			return err
		}
		b.BoneFlag = BoneFlag(uint16(boneFlag[0]) | uint16(boneFlag[1])<<8)

		// 0x0001  : 接続先(PMD子ボーン指定)表示方法 -> 0:座標オフセットで指定 1:ボーンで指定
		if b.IsTailBone() {
			// n  : ボーンIndexサイズ  | 接続先ボーンのボーンIndex
			b.TailIndex, err = r.unpackBoneIndex(model)
			if err != nil {
				mlog.E("[%d] readBones unpackBoneIndex TailIndex error: %v", i, err)
				return err
			}
			b.TailPosition = mmath.NewMVec3()
		} else {
			//  12 : float3	| 座標オフセット, ボーン位置からの相対分
			b.TailIndex = -1
			tailPos, err := r.UnpackVec3(true)
			if err != nil {
				mlog.E("[%d] readBones UnpackVec3 TailPosition error: %v", i, err)
				return err
			}
			b.TailPosition = &tailPos
		}

		// 回転付与:1 または 移動付与:1 の場合
		if b.IsEffectorRotation() || b.IsEffectorTranslation() {
			// n  : ボーンIndexサイズ  | 付与親ボーンのボーンIndex
			b.EffectIndex, err = r.unpackBoneIndex(model)
			if err != nil {
				mlog.E("[%d] readBones unpackBoneIndex EffectIndex error: %v", i, err)
				return err
			}
			// 4  : float	| 付与率
			b.EffectFactor, err = r.UnpackFloat()
			if err != nil {
				mlog.E("[%d] readBones UnpackFloat EffectFactor error: %v", i, err)
				return err
			}
		} else {
			b.EffectIndex = -1
			b.EffectFactor = 0
		}

		// 軸固定:1 の場合
		if b.HasFixedAxis() {
			// 12 : float3	| 軸の方向ベクトル
			fixedAxis, err := r.UnpackVec3(true)
			if err != nil {
				mlog.E("[%d] readBones UnpackVec3 FixedAxis error: %v", i, err)
				return err
			}
			b.FixedAxis = &fixedAxis
			b.NormalizeFixedAxis(b.FixedAxis.Normalize())
		} else {
			b.FixedAxis = mmath.NewMVec3()
		}

		// ローカル軸:1 の場合
		if b.HasLocalAxis() {
			// 12 : float3	| X軸の方向ベクトル
			localAxisX, err := r.UnpackVec3(true)
			if err != nil {
				mlog.E("[%d] readBones UnpackVec3 LocalAxisX error: %v", i, err)
				return err
			}
			b.LocalAxisX = &localAxisX
			// 12 : float3	| Z軸の方向ベクトル
			localAxisZ, err := r.UnpackVec3(true)
			if err != nil {
				mlog.E("[%d] readBones UnpackVec3 LocalAxisZ error: %v", i, err)
				return err
			}
			b.LocalAxisZ = &localAxisZ
			b.NormalizeLocalAxis(b.LocalAxisX)
		} else {
			b.LocalAxisX = mmath.NewMVec3()
			b.LocalAxisZ = mmath.NewMVec3()
		}

		// 外部親変形:1 の場合
		if b.IsEffectorParentDeform() {
			// 4  : int	| Key値
			b.EffectorKey, err = r.UnpackInt()
			if err != nil {
				mlog.E("[%d] readBones UnpackInt EffectorKey error: %v", i, err)
				return err
			}
		}

		// IK:1 の場合 IKデータを格納
		if b.IsIK() {
			b.Ik = NewIk()

			// n  : ボーンIndexサイズ  | IKターゲットボーンのボーンIndex
			b.Ik.BoneIndex, err = r.unpackBoneIndex(model)
			if err != nil {
				mlog.E("[%d] readBones unpackBoneIndex Ik.BoneIndex error: %v", i, err)
				return err
			}
			// 4  : int  	| IKループ回数 (PMD及びMMD環境では255回が最大になるようです)
			b.Ik.LoopCount, err = r.UnpackInt()
			if err != nil {
				mlog.E("[%d] readBones UnpackInt Ik.LoopCount error: %v", i, err)
				return err
			}
			// 4  : float	| IKループ計算時の1回あたりの制限角度 -> ラジアン角 | PMDのIK値とは4倍異なるので注意
			unitRot, err := r.UnpackFloat()
			if err != nil {
				mlog.E("[%d] readBones UnpackFloat unitRot error: %v", i, err)
				return err
			}
			b.Ik.UnitRotation.SetRadians(&mmath.MVec3{unitRot, unitRot, unitRot})
			// 4  : int  	| IKリンク数 : 後続の要素数
			ikLinkCount, err := r.UnpackInt()
			if err != nil {
				mlog.E("[%d] readBones UnpackInt ikLinkCount error: %v", i, err)
				return err
			}
			for j := 0; j < ikLinkCount; j++ {
				il := NewIkLink()
				// n  : ボーンIndexサイズ  | リンクボーンのボーンIndex
				il.BoneIndex, err = r.unpackBoneIndex(model)
				if err != nil {
					mlog.E("[%d][%d] readBones unpackBoneIndex IkLink.BoneIndex error: %v", i, j, err)
					return err
				}
				// 1  : byte	| 角度制限 0:OFF 1:ON
				ikLinkAngleLimit, err := r.UnpackByte()
				if err != nil {
					mlog.E("[%d][%d] readBones UnpackByte IkLink.AngleLimit error: %v", i, j, err)
					return err
				}
				il.AngleLimit = ikLinkAngleLimit == 1
				if il.AngleLimit {
					// 12 : float3	| 下限 (x,y,z) -> ラジアン角
					minAngleLimit, err := r.UnpackVec3(false)
					if err != nil {
						mlog.E("[%d][%d] readBones UnpackVec3 IkLink.MinAngleLimit error: %v", i, j, err)
						return err
					}
					il.MinAngleLimit.SetRadians(&minAngleLimit)
					// 12 : float3	| 上限 (x,y,z) -> ラジアン角
					maxAngleLimit, err := r.UnpackVec3(false)
					if err != nil {
						mlog.E("[%d][%d] readBones UnpackVec3 IkLink.MaxAngleLimit error: %v", i, j, err)
						return err
					}
					il.MaxAngleLimit.SetRadians(&maxAngleLimit)
				}
				b.Ik.Links = append(b.Ik.Links, il)
			}
		}

		model.Bones.Append(&b)
	}

	return nil
}

func (r *PmxReader) readMorphs(model *PmxModel) error {
	totalMorphCount, err := r.UnpackInt()
	if err != nil {
		mlog.E("readMorphs UnpackInt totalMorphCount error: %v", err)
		return err
	}

	for i := 0; i < totalMorphCount; i++ {
		// 4 + n : TextBuf	| モーフ名
		name := r.ReadText()
		// 4 + n : TextBuf	| モーフ名英
		englishName := r.ReadText()

		m := Morph{
			IndexNameModel: &mcore.IndexNameModel{Index: i, Name: name, EnglishName: englishName},
		}

		// 1  : byte	| 操作パネル (PMD:カテゴリ) 1:眉(左下) 2:目(左上) 3:口(右上) 4:その他(右下)  | 0:システム予約
		panel, err := r.UnpackByte()
		if err != nil {
			return err
		}
		m.Panel = MorphPanel(panel)
		// 1  : byte	| モーフ種類 - 0:グループ, 1:頂点, 2:ボーン, 3:UV, 4:追加UV1, 5:追加UV2, 6:追加UV3, 7:追加UV4, 8:材質
		morphType, err := r.UnpackByte()
		if err != nil {
			mlog.E("[%d] readMorphs UnpackByte MorphType error: %v", i, err)
			return err
		}
		m.MorphType = MorphType(morphType)

		offsetCount, err := r.UnpackInt()
		if err != nil {
			mlog.E("[%d] readMorphs UnpackInt OffsetCount error: %v", i, err)
			return err
		}
		for j := 0; j < offsetCount; j++ {
			switch m.MorphType {
			case MORPH_TYPE_GROUP:
				// n  : モーフIndexサイズ  | モーフIndex  ※仕様上グループモーフのグループ化は非対応とする
				morphIndex, err := r.unpackMorphIndex(model)
				if err != nil {
					mlog.E("[%d][%d] readMorphs unpackMorphIndex MorphIndex error: %v", i, j, err)
					return err
				}
				// 4  : float	| モーフ率 : グループモーフのモーフ値 * モーフ率 = 対象モーフのモーフ値
				morphFactor, err := r.UnpackFloat()
				if err != nil {
					mlog.E("[%d][%d] readMorphs UnpackFloat MorphFactor error: %v", i, j, err)
					return err
				}
				m.Offsets = append(m.Offsets, NewGroupMorph(morphIndex, morphFactor))
			case MORPH_TYPE_VERTEX:
				// n  : 頂点Indexサイズ  | 頂点Index
				vertexIndex, err := r.unpackVertexIndex(model)
				if err != nil {
					mlog.E("[%d][%d] readMorphs unpackVertexIndex VertexIndex error: %v", i, j, err)
					return err
				}
				// 12 : float3	| 座標オフセット量(x,y,z)
				offset, err := r.UnpackVec3(true)
				if err != nil {
					mlog.E("[%d][%d] readMorphs UnpackVec3 Offset error: %v", i, j, err)
					return err
				}
				m.Offsets = append(m.Offsets, NewVertexMorph(vertexIndex, &offset))
			case MORPH_TYPE_BONE:
				// n  : ボーンIndexサイズ  | ボーンIndex
				boneIndex, err := r.unpackBoneIndex(model)
				if err != nil {
					mlog.E("[%d][%d] readMorphs unpackBoneIndex BoneIndex error: %v", i, j, err)
					return err
				}
				// 12 : float3	| 移動量(x,y,z)
				offset, err := r.UnpackVec3(true)
				if err != nil {
					mlog.E("[%d][%d] readMorphs UnpackVec3 Offset error: %v", i, j, err)
					return err
				}
				// 16 : float4	| 回転量(x,y,z,w)
				qq, err := r.UnpackQuaternion(true)
				if err != nil {
					mlog.E("[%d][%d] readMorphs UnpackQuaternion Quaternion error: %v", i, j, err)
					return err
				}
				m.Offsets = append(m.Offsets, NewBoneMorph(boneIndex, &offset, mmath.NewRotationFromQuaternion(&qq)))
			case MORPH_TYPE_UV, MORPH_TYPE_EXTENDED_UV1, MORPH_TYPE_EXTENDED_UV2, MORPH_TYPE_EXTENDED_UV3, MORPH_TYPE_EXTENDED_UV4:
				// n  : 頂点Indexサイズ  | 頂点Index
				vertexIndex, err := r.unpackVertexIndex(model)
				if err != nil {
					mlog.E("[%d][%d] readMorphs unpackVertexIndex VertexIndex error: %v", i, j, err)
					return err
				}
				// 16 : float4	| UVオフセット量(x,y,z,w) ※通常UVはz,wが不要項目になるがモーフとしてのデータ値は記録しておく
				offset, err := r.UnpackVec4(false)
				if err != nil {
					mlog.E("[%d][%d] readMorphs UnpackVec4 Offset error: %v", i, j, err)
					return err
				}
				m.Offsets = append(m.Offsets, NewUvMorph(vertexIndex, &offset))
			case MORPH_TYPE_MATERIAL:
				// n  : 材質Indexサイズ  | 材質Index -> -1:全材質対象
				materialIndex, err := r.unpackMaterialIndex(model)
				if err != nil {
					mlog.E("[%d][%d] readMorphs unpackMaterialIndex MaterialIndex error: %v", i, j, err)
					return err
				}
				// 1  : オフセット演算形式 | 0:乗算, 1:加算 - 詳細は後述
				calcMode, err := r.UnpackByte()
				if err != nil {
					mlog.E("[%d][%d] readMorphs UnpackByte CalcMode error: %v", i, j, err)
					return err
				}
				// 16 : float4	| Diffuse (R,G,B,A) - 乗算:1.0／加算:0.0 が初期値となる(同以下)
				diffuse, err := r.UnpackVec4(false)
				if err != nil {
					mlog.E("[%d][%d] readMorphs UnpackVec4 Diffuse error: %v", i, j, err)
					return err
				}
				// 12 : float3	| Specular (R,G,B, Specular係数)
				specular, err := r.UnpackVec4(false)
				if err != nil {
					mlog.E("[%d][%d] readMorphs UnpackVec4 Specular error: %v", i, j, err)
					return err
				}
				// 12 : float3	| Ambient (R,G,B)
				ambient, err := r.UnpackVec3(false)
				if err != nil {
					mlog.E("[%d][%d] readMorphs UnpackVec3 Ambient error: %v", i, j, err)
					return err
				}
				// 16 : float4	| エッジ色 (R,G,B,A)
				edge, err := r.UnpackVec4(false)
				if err != nil {
					mlog.E("[%d][%d] readMorphs UnpackVec4 Edge error: %v", i, j, err)
					return err
				}
				// 4  : float	| エッジサイズ
				edgeSize, err := r.UnpackFloat()
				if err != nil {
					mlog.E("[%d][%d] readMorphs UnpackFloat EdgeSize error: %v", i, j, err)
					return err
				}
				// 16 : float4	| テクスチャ係数 (R,G,B,A)
				textureFactor, err := r.UnpackVec4(false)
				if err != nil {
					mlog.E("[%d][%d] readMorphs UnpackVec4 TextureFactor error: %v", i, j, err)
					return err
				}
				// 16 : float4	| スフィアテクスチャ係数 (R,G,B,A)
				sphereTextureFactor, err := r.UnpackVec4(false)
				if err != nil {
					mlog.E("[%d][%d] readMorphs UnpackVec4 SphereTextureFactor error: %v", i, j, err)
					return err
				}
				// 16 : float4	| Toonテクスチャ係数 (R,G,B,A)
				toonTextureFactor, err := r.UnpackVec4(false)
				if err != nil {
					mlog.E("[%d][%d] readMorphs UnpackVec4 ToonTextureFactor error: %v", i, j, err)
					return err
				}
				m.Offsets = append(m.Offsets, NewMaterialMorph(
					materialIndex,
					MaterialMorphCalcMode(calcMode),
					&diffuse,
					&specular,
					&ambient,
					&edge,
					edgeSize,
					&textureFactor,
					&sphereTextureFactor,
					&toonTextureFactor,
				))
			}
		}

		model.Morphs.Append(&m)
	}

	return nil
}

func (r *PmxReader) readDisplaySlots(model *PmxModel) error {
	totalDisplaySlotCount, err := r.UnpackInt()
	if err != nil {
		mlog.E("readDisplaySlots UnpackInt totalDisplaySlotCount error: %v", err)
		return err
	}

	for i := 0; i < totalDisplaySlotCount; i++ {
		// 4 + n : TextBuf	| 枠名
		name := r.ReadText()
		// 4 + n : TextBuf	| 枠名英
		englishName := r.ReadText()

		d := DisplaySlot{
			IndexNameModel: &mcore.IndexNameModel{Index: i, Name: name, EnglishName: englishName},
			References:     make([]Reference, 0),
		}

		// 1  : byte	| 特殊枠フラグ - 0:通常枠 1:特殊枠
		specialFlag, err := r.UnpackByte()
		if err != nil {
			mlog.E("[%d] readDisplaySlots UnpackByte SpecialFlag error: %v", i, err)
			return err
		}
		d.SpecialFlag = SpecialFlag(specialFlag)

		// 4  : int  	| 枠内要素数 : 後続の要素数
		referenceCount, err := r.UnpackInt()
		if err != nil {
			mlog.E("[%d] readDisplaySlots UnpackInt ReferenceCount error: %v", i, err)
			return err
		}

		for j := 0; j < referenceCount; j++ {
			reference := NewDisplaySlotReference()

			// 1  : byte	| 要素種別 - 0:ボーン 1:モーフ
			referenceType, err := r.UnpackByte()
			if err != nil {
				mlog.E("[%d][%d] readDisplaySlots UnpackByte ReferenceType error: %v", i, j, err)
				return err
			}
			reference.DisplayType = DisplayType(referenceType)

			switch reference.DisplayType {
			case DISPLAY_TYPE_BONE:
				// n  : ボーンIndexサイズ  | ボーンIndex
				reference.DisplayIndex, err = r.unpackBoneIndex(model)
				if err != nil {
					mlog.E("[%d][%d] readDisplaySlots unpackBoneIndex DisplayIndex error: %v", i, j, err)
					return err
				}
			case DISPLAY_TYPE_MORPH:
				// n  : モーフIndexサイズ  | モーフIndex
				reference.DisplayIndex, err = r.unpackMorphIndex(model)
				if err != nil {
					mlog.E("[%d][%d] readDisplaySlots unpackMorphIndex DisplayIndex error: %v", i, j, err)
					return err
				}
			}

			d.References = append(d.References, *reference)
		}
		model.DisplaySlots.Append(&d)
	}

	return nil
}

func (r *PmxReader) readRigidBodies(model *PmxModel) error {
	totalRigidBodyCount, err := r.UnpackInt()
	if err != nil {
		mlog.E("readRigidBodies UnpackInt totalRigidBodyCount error: %v", err)
		return err
	}

	for i := 0; i < totalRigidBodyCount; i++ {
		// 4 + n : TextBuf	| 剛体名
		name := r.ReadText()
		// 4 + n : TextBuf	| 剛体名英
		englishName := r.ReadText()

		b := RigidBody{
			IndexNameModel: &mcore.IndexNameModel{Index: i, Name: name, EnglishName: englishName},
			RigidBodyParam: NewRigidBodyParam(),
		}

		// n  : ボーンIndexサイズ  | 関連ボーンIndex - 関連なしの場合は-1
		b.BoneIndex, err = r.unpackBoneIndex(model)
		if err != nil {
			mlog.E("[%d] readRigidBodies unpackBoneIndex BoneIndex error: %v", i, err)
			return err
		}
		// 1  : byte	| グループ
		collisionGroup, err := r.UnpackByte()
		if err != nil {
			mlog.E("[%d] readRigidBodies UnpackByte CollisionGroup error: %v", i, err)
			return err
		}
		b.CollisionGroup = collisionGroup
		// 2  : ushort	| 非衝突グループフラグ
		collisionGroupMask, err := r.UnpackUShort()
		if err != nil {
			mlog.E("[%d] readRigidBodies UnpackUShort CollisionGroupMask error: %v", i, err)
			return err
		}
		b.CollisionGroupMaskValue = int(collisionGroupMask)
		b.CollisionGroupMask.IsCollisions = NewCollisionGroup(collisionGroupMask)
		// 1  : byte	| 形状 - 0:球 1:箱 2:カプセル
		shapeType, err := r.UnpackByte()
		if err != nil {
			mlog.E("[%d] readRigidBodies UnpackByte ShapeType error: %v", i, err)
			return err
		}
		b.ShapeType = Shape(shapeType)
		// 12 : float3	| サイズ(x,y,z)
		size, err := r.UnpackVec3(false)
		if err != nil {
			mlog.E("[%d] readRigidBodies UnpackVec3 Size error: %v", i, err)
			return err
		}
		b.Size = &size
		// 12 : float3	| 位置(x,y,z)
		position, err := r.UnpackVec3(true)
		if err != nil {
			mlog.E("[%d] readRigidBodies UnpackVec3 Position error: %v", i, err)
			return err
		}
		b.Position = &position
		// 12 : float3	| 回転(x,y,z) -> ラジアン角
		rads, err := r.UnpackVec3(false)
		if err != nil {
			mlog.E("[%d] readRigidBodies UnpackVec3 Rotation error: %v", i, err)
			return err
		}
		b.Rotation = mmath.NewRotationFromRadians(&rads)
		// 4  : float	| 質量
		b.RigidBodyParam.Mass, err = r.UnpackFloat()
		if err != nil {
			mlog.E("[%d] readRigidBodies UnpackFloat Mass error: %v", i, err)
			return err
		}
		// 4  : float	| 移動減衰
		b.RigidBodyParam.LinearDamping, err = r.UnpackFloat()
		if err != nil {
			mlog.E("[%d] readRigidBodies UnpackFloat LinearDamping error: %v", i, err)
			return err
		}
		// 4  : float	| 回転減衰
		b.RigidBodyParam.AngularDamping, err = r.UnpackFloat()
		if err != nil {
			mlog.E("[%d] readRigidBodies UnpackFloat AngularDamping error: %v", i, err)
			return err
		}
		// 4  : float	| 反発力
		b.RigidBodyParam.Restitution, err = r.UnpackFloat()
		if err != nil {
			mlog.E("[%d] readRigidBodies UnpackFloat Restitution error: %v", i, err)
			return err
		}
		// 4  : float	| 摩擦力
		b.RigidBodyParam.Friction, err = r.UnpackFloat()
		if err != nil {
			mlog.E("[%d] readRigidBodies UnpackFloat Friction error: %v", i, err)
			return err
		}
		// 1  : byte	| 剛体の物理演算 - 0:ボーン追従(static) 1:物理演算(dynamic) 2:物理演算 + Bone位置合わせ
		physicsType, err := r.UnpackByte()
		if err != nil {
			mlog.E("[%d] readRigidBodies UnpackByte PhysicsType error: %v", i, err)
			return err
		}
		b.PhysicsType = PhysicsType(physicsType)

		model.RigidBodies.Append(&b)
	}

	return nil
}

func (r *PmxReader) readJoints(model *PmxModel) error {
	totalJointCount, err := r.UnpackInt()
	if err != nil {
		mlog.E("readJoints UnpackInt totalJointCount error: %v", err)
		return err
	}

	for i := 0; i < totalJointCount; i++ {
		// 4 + n : TextBuf	| Joint名
		name := r.ReadText()
		// 4 + n : TextBuf	| Joint名英
		englishName := r.ReadText()

		j := Joint{
			IndexNameModel: &mcore.IndexNameModel{Index: i, Name: name, EnglishName: englishName},
			JointParam:     NewJointParam(),
		}

		// 1  : byte	| Joint種類 - 0:スプリング6DOF   | PMX2.0では 0 のみ(拡張用)
		j.JointType, err = r.UnpackByte()
		if err != nil {
			mlog.E("[%d] readJoints UnpackByte JointType error: %v", i, err)
			return err
		}
		// n  : 剛体Indexサイズ  | 関連剛体AのIndex - 関連なしの場合は-1
		j.RigidbodyIndexA, err = r.unpackRigidBodyIndex(model)
		if err != nil {
			mlog.E("[%d] readJoints unpackRigidBodyIndex RigidbodyIndexA error: %v", i, err)
			return err
		}
		// n  : 剛体Indexサイズ  | 関連剛体BのIndex - 関連なしの場合は-1
		j.RigidbodyIndexB, err = r.unpackRigidBodyIndex(model)
		if err != nil {
			mlog.E("[%d] readJoints unpackRigidBodyIndex RigidbodyIndexB error: %v", i, err)
			return err
		}
		// 12 : float3	| 位置(x,y,z)
		position, err := r.UnpackVec3(true)
		if err != nil {
			mlog.E("[%d] readJoints UnpackVec3 Position error: %v", i, err)
			return err
		}
		j.Position = &position
		// 12 : float3	| 回転(x,y,z) -> ラジアン角
		rads, err := r.UnpackVec3(false)
		if err != nil {
			mlog.E("[%d] readJoints UnpackVec3 Rotation error: %v", i, err)
			return err
		}
		j.Rotation = mmath.NewRotationFromRadians(&rads)
		// 12 : float3	| 移動制限-下限(x,y,z)
		translationLimitMin, err := r.UnpackVec3(false)
		if err != nil {
			mlog.E("[%d] readJoints UnpackVec3 TranslationLimitMin error: %v", i, err)
			return err
		}
		j.JointParam.TranslationLimitMin = &translationLimitMin
		// 12 : float3	| 移動制限-上限(x,y,z)
		translationLimitMax, err := r.UnpackVec3(false)
		if err != nil {
			mlog.E("[%d] readJoints UnpackVec3 TranslationLimitMax error: %v", i, err)
			return err
		}
		j.JointParam.TranslationLimitMax = &translationLimitMax
		// 12 : float3	| 回転制限-下限(x,y,z) -> ラジアン角
		rotationLimitMin, err := r.UnpackVec3(false)
		if err != nil {
			mlog.E("[%d] readJoints UnpackVec3 RotationLimitMin error: %v", i, err)
			return err
		}
		j.JointParam.RotationLimitMin = mmath.NewRotationFromRadians(&rotationLimitMin)
		// 12 : float3	| 回転制限-上限(x,y,z) -> ラジアン角
		rotationLimitMax, err := r.UnpackVec3(false)
		if err != nil {
			mlog.E("[%d] readJoints UnpackVec3 RotationLimitMax error: %v", i, err)
			return err
		}
		j.JointParam.RotationLimitMax = mmath.NewRotationFromRadians(&rotationLimitMax)
		// 12 : float3	| バネ定数-移動(x,y,z)
		springConstantTranslation, err := r.UnpackVec3(false)
		if err != nil {
			mlog.E("[%d] readJoints UnpackVec3 SpringConstantTranslation error: %v", i, err)
			return err
		}
		j.JointParam.SpringConstantTranslation = &springConstantTranslation
		// 12 : float3	| バネ定数-回転(x,y,z)
		springConstantRotation, err := r.UnpackVec3(false)
		if err != nil {
			mlog.E("[%d] readJoints UnpackVec3 SpringConstantRotation error: %v", i, err)
			return err
		}
		j.JointParam.SpringConstantRotation = mmath.NewRotationFromDegrees(&springConstantRotation)

		model.Joints.Append(&j)
	}

	return nil
}

// テキストデータを読み取る
func (r *PmxReader) unpackVertexIndex(model *PmxModel) (int, error) {
	switch model.VertexCountType {
	case 1:
		v, err := r.UnpackByte()
		if err != nil {
			mlog.E("unpackVertexIndex.UnpackByte error: %v", err)
			return 0, err
		}
		return int(v), nil
	case 2:
		v, err := r.UnpackUShort()
		if err != nil {
			mlog.E("unpackVertexIndex.UnpackUShort error: %v", err)
			return 0, err
		}
		return int(v), nil
	case 4:
		v, err := r.UnpackInt()
		if err != nil {
			mlog.E("unpackVertexIndex.UnpackInt error: %v", err)
			return 0, err
		}
		return v, nil
	}
	return 0, fmt.Errorf("未知のVertexIndexサイズです。vertexCount: %d", model.VertexCountType)
}

// テクスチャIndexを読み取る
func (r *PmxReader) unpackTextureIndex(model *PmxModel) (int, error) {
	return r.unpackIndex(model.TextureCountType)
}

// 材質Indexを読み取る
func (r *PmxReader) unpackMaterialIndex(model *PmxModel) (int, error) {
	return r.unpackIndex(model.MaterialCountType)
}

// ボーンIndexを読み取る
func (r *PmxReader) unpackBoneIndex(model *PmxModel) (int, error) {
	return r.unpackIndex(model.BoneCountType)
}

// 表情Indexを読み取る
func (r *PmxReader) unpackMorphIndex(model *PmxModel) (int, error) {
	return r.unpackIndex(model.MorphCountType)
}

// 剛体Indexを読み取る
func (r *PmxReader) unpackRigidBodyIndex(model *PmxModel) (int, error) {
	return r.unpackIndex(model.RigidBodyCountType)
}

func (r *PmxReader) unpackIndex(index int) (int, error) {
	switch index {
	case 1:
		v, err := r.UnpackSByte()
		if err != nil {
			mlog.E("unpackIndex.UnpackSByte error: %v", err)
			return 0, err
		}
		return int(v), nil
	case 2:
		v, err := r.UnpackShort()
		if err != nil {
			mlog.E("unpackIndex.UnpackShort error: %v", err)
			return 0, err
		}
		return int(v), nil
	case 4:
		v, err := r.UnpackInt()
		if err != nil {
			mlog.E("unpackIndex.UnpackInt error: %v", err)
			return 0, err
		}
		return v, nil
	}
	return 0, fmt.Errorf("未知のIndexサイズです。index: %d", index)
}
