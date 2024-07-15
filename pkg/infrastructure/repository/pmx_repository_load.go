package repository

import (
	"fmt"
	"slices"

	"github.com/miu200521358/mlib_go/pkg/domain/core"
	"github.com/miu200521358/mlib_go/pkg/domain/mmath"
	"github.com/miu200521358/mlib_go/pkg/domain/pmx"
	"github.com/miu200521358/mlib_go/pkg/mutils/mlog"
	"golang.org/x/text/encoding/unicode"
)

// 指定されたパスのファイルからデータを読み込む
func (r *PmxRepository) Load(path string) (core.IHashModel, error) {
	// モデルを新規作成
	model := r.newFunc(path)

	hash, err := r.LoadHash(path)
	if err != nil {
		mlog.E("ReadByFilepath.ReadHashByFilePath error: %v", err)
		return nil, err
	}
	model.SetHash(hash)

	// ファイルを開く
	err = r.open(path)
	if err != nil {
		mlog.E("ReadByFilepath.Open error: %v", err)
		return model, err
	}

	err = r.loadHeader(model)
	if err != nil {
		mlog.E("ReadByFilepath.loadHeader error: %v", err)
		return model, err
	}

	err = r.loadModel(model)
	if err != nil {
		mlog.E("ReadByFilepath.loadData error: %v", err)
		return model, err
	}

	r.close()
	model.Setup()

	return model, nil
}

func (r *PmxRepository) LoadName(path string) (string, error) {
	// モデルを新規作成
	model := r.newFunc(path)

	// ファイルを開く
	err := r.open(path)
	if err != nil {
		mlog.E("LoadName.Open error: %v", err)
		return "", err
	}

	err = r.loadHeader(model)
	if err != nil {
		mlog.E("LoadName.loadHeader error: %v", err)
		return "", err
	}

	r.close()

	return model.Name, nil
}

func (r *PmxRepository) loadHeader(model *pmx.PmxModel) error {
	fbytes, err := r.unpackBytes(4)
	if err != nil {
		mlog.E("loadHeader.unpackBytes error: %v", err)
		return err
	}
	model.Signature = r.decodeText(unicode.UTF8, fbytes)
	model.Version, err = r.unpackFloat()

	if err != nil {
		mlog.E("loadHeader.Version error: %v", err)
		return err
	}

	if model.Signature[:3] != "PMX" ||
		!slices.Contains([]string{"2.0", "2.1"}, fmt.Sprintf("%.1f", model.Version)) {
		// 整合性チェック
		return fmt.Errorf("PMX2.0/2.1形式外のデータです。signature: %s, version: %.1f ", model.Signature, model.Version)
	}

	// 1 : byte	| 後続するデータ列のバイトサイズ  PMX2.0は 8 で固定
	_, _ = r.unpackByte()

	// [0] - エンコード方式  | 0:UTF16 1:UTF8
	encodeType, err := r.unpackByte()
	if err != nil {
		mlog.E("UnpackByte encodeType error: %v", err)
		return err
	}

	switch encodeType {
	case 0:
		r.defineEncoding(unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM))
	case 1:
		r.defineEncoding(unicode.UTF8)
	default:
		return fmt.Errorf("未知のエンコードタイプです。encodeType: %d", encodeType)
	}

	// [1] - 追加UV数 	| 0～4 詳細は頂点参照
	extendedUVCount, err := r.unpackByte()
	if err != nil {
		mlog.E("UnpackByte extendedUVCount error: %v", err)
		return err
	}
	model.ExtendedUVCount = int(extendedUVCount)
	// [2] - 頂点Indexサイズ | 1,2,4 のいずれか
	vertexCount, err := r.unpackByte()
	if err != nil {
		mlog.E("UnpackByte vertexCount error: %v", err)
		return err
	}
	model.VertexCountType = int(vertexCount)
	// [3] - テクスチャIndexサイズ | 1,2,4 のいずれか
	textureCount, err := r.unpackByte()
	if err != nil {
		mlog.E("UnpackByte textureCount error: %v", err)
		return err
	}
	model.TextureCountType = int(textureCount)
	// [4] - 材質Indexサイズ | 1,2,4 のいずれか
	materialCount, err := r.unpackByte()
	if err != nil {
		mlog.E("UnpackByte materialCount error: %v", err)
		return err
	}
	model.MaterialCountType = int(materialCount)
	// [5] - ボーンIndexサイズ | 1,2,4 のいずれか
	boneCount, err := r.unpackByte()
	if err != nil {
		mlog.E("UnpackByte boneCount error: %v", err)
		return err
	}
	model.BoneCountType = int(boneCount)
	// [6] - モーフIndexサイズ | 1,2,4 のいずれか
	morphCount, err := r.unpackByte()
	if err != nil {
		mlog.E("UnpackByte morphCount error: %v", err)
		return err
	}
	model.MorphCountType = int(morphCount)
	// [7] - 剛体Indexサイズ | 1,2,4 のいずれか
	rigidBodyCount, err := r.unpackByte()
	if err != nil {
		mlog.E("UnpackByte rigidBodyCount error: %v", err)
		return err
	}
	model.RigidBodyCountType = int(rigidBodyCount)

	// 4 + n : TextBuf	| モデル名
	model.Name = r.readText()
	// 4 + n : TextBuf	| モデル名英
	model.EnglishName = r.readText()
	// 4 + n : TextBuf	| コメント
	model.Comment = r.readText()
	// 4 + n : TextBuf	| コメント英
	model.EnglishComment = r.readText()

	return nil
}

func (r *PmxRepository) loadModel(model *pmx.PmxModel) error {
	err := r.loadVertices(model)
	if err != nil {
		mlog.E("loadData.loadVertices error: %v", err)
		return err
	}

	err = r.loadFaces(model)
	if err != nil {
		mlog.E("loadData.loadFaces error: %v", err)
		return err
	}

	err = r.loadTextures(model)
	if err != nil {
		mlog.E("loadData.loadTextures error: %v", err)
		return err
	}

	err = r.loadMaterials(model)
	if err != nil {
		mlog.E("loadData.loadMaterials error: %v", err)
		return err
	}

	err = r.loadBones(model)
	if err != nil {
		mlog.E("loadData.loadBones error: %v", err)
		return err
	}

	err = r.loadMorphs(model)
	if err != nil {
		mlog.E("loadData.loadMorphs error: %v", err)
		return err
	}

	err = r.loadDisplaySlots(model)
	if err != nil {
		mlog.E("loadData.loadDisplaySlots error: %v", err)
		return err
	}

	err = r.loadRigidBodies(model)
	if err != nil {
		mlog.E("loadData.loadRigidBodies error: %v", err)
		return err
	}

	err = r.loadJoints(model)
	if err != nil {
		mlog.E("loadData.loadJoints error: %v", err)
		return err
	}

	return nil
}

func (r *PmxRepository) loadVertices(model *pmx.PmxModel) error {
	totalVertexCount, err := r.unpackInt()
	if err != nil {
		mlog.E("loadVertices UnpackInt totalVertexCount error: %v", err)
		return err
	}

	vertices := pmx.NewVertices(totalVertexCount)

	for i := 0; i < totalVertexCount; i++ {
		v := &pmx.Vertex{IndexModel: &core.IndexModel{Index: i}}

		// 12 : float3  | 位置(x,y,z)
		pos, err := r.unpackVec3()
		if err != nil {
			mlog.E("[%d] loadVertices UnpackFloat Position error: %v", i, err)
			return err
		}
		v.Position = &pos

		// 12 : float3  | 法線(x,y,z)
		normal, err := r.unpackVec3()
		if err != nil {
			mlog.E("[%d] loadVertices UnpackFloat Normal[0] error: %v", i, err)
			return err
		}
		v.Normal = &normal

		// 8  : float2  | UV(u,v)
		uv, err := r.unpackVec2()
		if err != nil {
			mlog.E("[%d] loadVertices UnpackFloat UV[0] error: %v", i, err)
			return err
		}
		v.Uv = &uv

		// 16 * n : float4[n] | 追加UV(x,y,z,w)  PMXヘッダの追加UV数による
		v.ExtendedUvs = make([]*mmath.MVec4, 0)
		for j := 0; j < model.ExtendedUVCount; j++ {
			extendedUV, err := r.unpackVec4()
			if err != nil {
				mlog.E("[%d][%d] loadVertices UnpackVec4 ExtendedUV error: %v", i, j, err)
				return err
			}
			v.ExtendedUvs = append(v.ExtendedUvs, &extendedUV)
		}

		// 1 : byte    | ウェイト変形方式 0:BDEF1 1:BDEF2 2:BDEF4 3:SDEF
		Type, err := r.unpackByte()
		if err != nil {
			mlog.E("[%d] loadVertices UnpackByte Type error: %v", i, err)
			return err
		}
		v.DeformType = pmx.DeformType(Type)

		switch v.DeformType {
		case pmx.BDEF1:
			// n : ボーンIndexサイズ  | ウェイト1.0の単一ボーン(参照Index)
			boneIndex, err := r.unpackBoneIndex(model)
			if err != nil {
				mlog.E("[%d] loadVertices BDEF1 unpackBoneIndex error: %v", i, err)
				return err
			}
			v.Deform = pmx.NewBdef1(boneIndex)
		case pmx.BDEF2:
			// n : ボーンIndexサイズ  | ボーン1の参照Index
			boneIndex1, err := r.unpackBoneIndex(model)
			if err != nil {
				mlog.E("[%d] loadVertices BDEF2 unpackBoneIndex boneIndex1 error: %v", i, err)
				return err
			}
			// n : ボーンIndexサイズ  | ボーン2の参照Index
			boneIndex2, err := r.unpackBoneIndex(model)
			if err != nil {
				mlog.E("[%d] loadVertices BDEF2 unpackBoneIndex boneIndex2 error: %v", i, err)
				return err
			}
			// 4 : float              | ボーン1のウェイト値(0～1.0), ボーン2のウェイト値は 1.0-ボーン1ウェイト
			boneWeight, err := r.unpackFloat()
			if err != nil {
				mlog.E("[%d] loadVertices BDEF2 UnpackFloat boneWeight error: %v", i, err)
				return err
			}
			v.Deform = pmx.NewBdef2(boneIndex1, boneIndex2, boneWeight)
		case pmx.BDEF4:
			// n : ボーンIndexサイズ  | ボーン1の参照Index
			boneIndex1, err := r.unpackBoneIndex(model)
			if err != nil {
				mlog.E("[%d] loadVertices BDEF4 unpackBoneIndex boneIndex1 error: %v", i, err)
				return err
			}
			// n : ボーンIndexサイズ  | ボーン2の参照Index
			boneIndex2, err := r.unpackBoneIndex(model)
			if err != nil {
				mlog.E("[%d] loadVertices BDEF4 unpackBoneIndex boneIndex2 error: %v", i, err)
				return err
			}
			// n : ボーンIndexサイズ  | ボーン3の参照Index
			boneIndex3, err := r.unpackBoneIndex(model)
			if err != nil {
				mlog.E("[%d] loadVertices BDEF4 unpackBoneIndex boneIndex3 error: %v", i, err)
				return err
			}
			// n : ボーンIndexサイズ  | ボーン4の参照Index
			boneIndex4, err := r.unpackBoneIndex(model)
			if err != nil {
				mlog.E("[%d] loadVertices BDEF4 unpackBoneIndex boneIndex4 error: %v", i, err)
				return err
			}
			// 4 : float              | ボーン1のウェイト値
			boneWeight1, err := r.unpackFloat()
			if err != nil {
				mlog.E("[%d] loadVertices BDEF4 UnpackFloat boneWeight1 error: %v", i, err)
				return err
			}
			// 4 : float              | ボーン2のウェイト値
			boneWeight2, err := r.unpackFloat()
			if err != nil {
				mlog.E("[%d] loadVertices BDEF4 UnpackFloat boneWeight2 error: %v", i, err)
				return err
			}
			// 4 : float              | ボーン3のウェイト値
			boneWeight3, err := r.unpackFloat()
			if err != nil {
				mlog.E("[%d] loadVertices BDEF4 UnpackFloat boneWeight3 error: %v", i, err)
				return err
			}
			// 4 : float              | ボーン4のウェイト値 (ウェイト計1.0の保障はない)
			boneWeight4, err := r.unpackFloat()
			if err != nil {
				mlog.E("[%d] loadVertices BDEF4 UnpackFloat boneWeight4 error: %v", i, err)
				return err
			}
			v.Deform = pmx.NewBdef4(boneIndex1, boneIndex2, boneIndex3, boneIndex4,
				boneWeight1, boneWeight2, boneWeight3, boneWeight4)
		case pmx.SDEF:
			// n : ボーンIndexサイズ  | ボーン1の参照Index
			boneIndex1, err := r.unpackBoneIndex(model)
			if err != nil {
				mlog.E("[%d] loadVertices SDEF unpackBoneIndex boneIndex1 error: %v", i, err)
				return err
			}
			// n : ボーンIndexサイズ  | ボーン2の参照Index
			boneIndex2, err := r.unpackBoneIndex(model)
			if err != nil {
				mlog.E("[%d] loadVertices SDEF unpackBoneIndex boneIndex2 error: %v", i, err)
				return err
			}
			// 4 : float              | ボーン1のウェイト値(0～1.0), ボーン2のウェイト値は 1.0-ボーン1ウェイト
			boneWeight, err := r.unpackFloat()
			if err != nil {
				mlog.E("[%d] loadVertices SDEF UnpackFloat boneWeight error: %v", i, err)
				return err
			}
			// 12 : float3             | SDEF-C値(x,y,z)
			sdefC, err := r.unpackVec3()
			if err != nil {
				mlog.E("[%d] loadVertices SDEF UnpackVec3 sdefC error: %v", i, err)
				return err
			}
			// 12 : float3             | SDEF-R0値(x,y,z)
			sdefR0, err := r.unpackVec3()
			if err != nil {
				mlog.E("[%d] loadVertices SDEF UnpackVec3 sdefR0 error: %v", i, err)
				return err
			}
			// 12 : float3             | SDEF-R1値(x,y,z) ※修正値を要計算
			sdefR1, err := r.unpackVec3()
			if err != nil {
				mlog.E("[%d] loadVertices SDEF UnpackVec3 sdefR1 error: %v", i, err)
				return err
			}
			v.Deform = pmx.NewSdef(boneIndex1, boneIndex2, boneWeight, &sdefC, &sdefR0, &sdefR1)
		}

		v.EdgeFactor, err = r.unpackFloat()
		if err != nil {
			mlog.E("[%d] loadVertices UnpackFloat EdgeFactor error: %v", i, err)
			return err
		}

		vertices.Update(v)
	}

	model.Vertices = vertices

	return nil
}

func (r *PmxRepository) loadFaces(model *pmx.PmxModel) error {
	totalFaceCount, err := r.unpackInt()
	if err != nil {
		mlog.E("loadFaces UnpackInt totalFaceCount error: %v", err)
		return err
	}

	faces := pmx.NewFaces(totalFaceCount / 3)

	for i := 0; i < totalFaceCount; i += 3 {
		f := &pmx.Face{
			IndexModel:    &core.IndexModel{Index: int(i / 3)},
			VertexIndexes: [3]int{},
		}

		// n : 頂点Indexサイズ     | 頂点の参照Index
		f.VertexIndexes[0], err = r.unpackVertexIndex(model)
		if err != nil {
			mlog.E("[%d] loadFaces unpackVertexIndex VertexIndexes[0] error: %v", i, err)
			return err
		}

		// n : 頂点Indexサイズ     | 頂点の参照Index
		f.VertexIndexes[1], err = r.unpackVertexIndex(model)
		if err != nil {
			mlog.E("[%d] loadFaces unpackVertexIndex VertexIndexes[1] error: %v", i, err)
			return err
		}

		// n : 頂点Indexサイズ     | 頂点の参照Index
		f.VertexIndexes[2], err = r.unpackVertexIndex(model)
		if err != nil {
			mlog.E("[%d] loadFaces unpackVertexIndex VertexIndexes[2] error: %v", i, err)
			return err
		}

		faces.Update(f)
	}

	model.Faces = faces

	return nil
}

func (r *PmxRepository) loadTextures(model *pmx.PmxModel) error {
	totalTextureCount, err := r.unpackInt()
	if err != nil {
		mlog.E("loadTextures UnpackInt totalTextureCount error: %v", err)
		return err
	}

	textures := pmx.NewTextures(totalTextureCount)

	for i := 0; i < totalTextureCount; i++ {
		t := &pmx.Texture{IndexModel: &core.IndexModel{Index: i}}

		// 4 + n : TextBuf	| テクスチャパス
		t.Name = r.readText()

		textures.Update(t)
	}

	model.Textures = textures

	return nil
}

func (r *PmxRepository) loadMaterials(model *pmx.PmxModel) error {
	totalMaterialCount, err := r.unpackInt()
	if err != nil {
		mlog.E("loadMaterials UnpackInt totalMaterialCount error: %v", err)
		return err
	}

	materials := pmx.NewMaterials(totalMaterialCount)

	for i := 0; i < totalMaterialCount; i++ {
		// 4 + n : TextBuf	| 材質名
		name := r.readText()
		// 4 + n : TextBuf	| 材質名英
		englishName := r.readText()

		m := &pmx.Material{
			IndexNameModel: &core.IndexNameModel{Index: i, Name: name, EnglishName: englishName},
		}

		// 16 : float4	| Diffuse (R,G,B,A)
		diffuse, err := r.unpackVec4()
		if err != nil {
			mlog.E("[%d] loadMaterials UnpackVec4 Diffuse error: %v", i, err)
			return err
		}
		m.Diffuse = &diffuse
		// 12 : float3	| Specular (R,G,B,Specular係数)
		specular, err := r.unpackVec4()
		if err != nil {
			mlog.E("[%d] loadMaterials UnpackVec4 Specular error: %v", i, err)
			return err
		}
		m.Specular = &specular
		// 12 : float3	| Ambient (R,G,B)
		ambient, err := r.unpackVec3()
		if err != nil {
			mlog.E("[%d] loadMaterials UnpackVec3 Ambient error: %v", i, err)
			return err
		}
		m.Ambient = &ambient
		// 1  : bitFlag  	| 描画フラグ(8bit) - 各bit 0:OFF 1:ON
		drawFlag, err := r.unpackByte()
		if err != nil {
			mlog.E("[%d] loadMaterials UnpackByte DrawFlag error: %v", i, err)
			return err
		}
		m.DrawFlag = pmx.DrawFlag(drawFlag)
		// 16 : float4	| エッジ色 (R,G,B,A)
		edge, err := r.unpackVec4()
		if err != nil {
			mlog.E("[%d] loadMaterials UnpackVec4 Edge error: %v", i, err)
			return err
		}
		m.Edge = &edge
		// 4  : float	| エッジサイズ
		m.EdgeSize, err = r.unpackFloat()
		if err != nil {
			mlog.E("[%d] loadMaterials UnpackFloat EdgeSize error: %v", i, err)
			return err
		}
		// n  : テクスチャIndexサイズ	| 通常テクスチャ
		m.TextureIndex, err = r.unpackTextureIndex(model)
		if err != nil {
			mlog.E("[%d] loadMaterials unpackTextureIndex TextureIndex error: %v", i, err)
			return err
		}
		// n  : テクスチャIndexサイズ	| スフィアテクスチャ
		m.SphereTextureIndex, err = r.unpackTextureIndex(model)
		if err != nil {
			mlog.E("[%d] loadMaterials unpackTextureIndex SphereTextureIndex error: %v", i, err)
			return err
		}
		// 1  : byte	| スフィアモード 0:無効 1:乗算(sph) 2:加算(spa) 3:サブテクスチャ(追加UV1のx,yをUV参照して通常テクスチャ描画を行う)
		sphereMode, err := r.unpackByte()
		if err != nil {
			mlog.E("[%d] loadMaterials UnpackByte SphereMode error: %v", i, err)
			return err
		}
		m.SphereMode = pmx.SphereMode(sphereMode)
		// 1  : byte	| 共有Toonフラグ 0:継続値は個別Toon 1:継続値は共有Toon
		toonSharingFlag, err := r.unpackByte()

		if err != nil {
			mlog.E("[%d] loadMaterials UnpackByte ToonSharingFlag error: %v", i, err)
			return err
		}
		m.ToonSharingFlag = pmx.ToonSharing(toonSharingFlag)

		switch m.ToonSharingFlag {
		case pmx.TOON_SHARING_INDIVIDUAL:
			// n  : テクスチャIndexサイズ	| Toonテクスチャ
			m.ToonTextureIndex, err = r.unpackTextureIndex(model)
			if err != nil {
				mlog.E("[%d] loadMaterials unpackTextureIndex ToonTextureIndex error: %v", i, err)
				return err
			}
		case pmx.TOON_SHARING_SHARING:
			// 1  : byte	| 共有ToonテクスチャIndex 0～9
			toonTextureIndex, err := r.unpackByte()
			if err != nil {
				mlog.E("[%d] loadMaterials UnpackByte ToonTextureIndex error: %v", i, err)
				return err
			}
			m.ToonTextureIndex = int(toonTextureIndex)
		}

		// 4 + n : TextBuf	| メモ
		m.Memo = r.readText()

		// 4  : int	| 材質に対応する面(頂点)数 (必ず3の倍数になる)
		m.VerticesCount, err = r.unpackInt()
		if err != nil {
			mlog.E("[%d] loadMaterials UnpackInt VerticesCount error: %v", i, err)
			return err
		}

		materials.Update(m)

	}

	model.Materials = materials

	return nil
}

func (r *PmxRepository) loadBones(model *pmx.PmxModel) error {
	totalBoneCount, err := r.unpackInt()
	if err != nil {
		mlog.E("loadBones UnpackInt totalBoneCount error: %v", err)
		return err
	}

	bones := pmx.NewBones(totalBoneCount)

	for i := 0; i < totalBoneCount; i++ {

		// 4 + n : TextBuf	| ボーン名
		name := r.readText()
		// 4 + n : TextBuf	| ボーン名英
		englishName := r.readText()

		b := &pmx.Bone{
			IndexNameModel:         &core.IndexNameModel{Index: i, Name: name, EnglishName: englishName},
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
			MinAngleLimit:          mmath.NewMRotation(),
			MaxAngleLimit:          mmath.NewMRotation(),
			LocalAngleLimit:        false,
			LocalMinAngleLimit:     mmath.NewMRotation(),
			LocalMaxAngleLimit:     mmath.NewMRotation(),
			AxisSign:               1,
		}

		// 12 : float3	| 位置
		pos, err := r.unpackVec3()
		if err != nil {
			mlog.E("[%d] loadBones UnpackVec3 Position error: %v", i, err)
			return err
		}
		b.Position = &pos
		// n  : ボーンIndexサイズ	| 親ボーン
		b.ParentIndex, err = r.unpackBoneIndex(model)
		if err != nil {
			mlog.E("[%d] loadBones unpackBoneIndex ParentIndex error: %v", i, err)
			return err
		}
		// 4  : int	| 変形階層
		b.Layer, err = r.unpackInt()
		if err != nil {
			mlog.E("[%d] UnpackInt Layer error: %v", i, err)
			return err
		}
		// 2  : bitFlag*2	| ボーンフラグ(16bit) 各bit 0:OFF 1:ON
		boneFlag, err := r.unpackBytes(2)
		if err != nil {
			mlog.E("[%d] loadBones UnpackBytes BoneFlag error: %v", i, err)
			return err
		}
		b.BoneFlag = pmx.BoneFlag(uint16(boneFlag[0]) | uint16(boneFlag[1])<<8)

		// 0x0001  : 接続先(PMD子ボーン指定)表示方法 -> 0:座標オフセットで指定 1:ボーンで指定
		if b.IsTailBone() {
			// n  : ボーンIndexサイズ  | 接続先ボーンのボーンIndex
			b.TailIndex, err = r.unpackBoneIndex(model)
			if err != nil {
				mlog.E("[%d] loadBones unpackBoneIndex TailIndex error: %v", i, err)
				return err
			}
			b.TailPosition = mmath.NewMVec3()
		} else {
			//  12 : float3	| 座標オフセット, ボーン位置からの相対分
			b.TailIndex = -1
			tailPos, err := r.unpackVec3()
			if err != nil {
				mlog.E("[%d] loadBones UnpackVec3 TailPosition error: %v", i, err)
				return err
			}
			b.TailPosition = &tailPos
		}

		// 回転付与:1 または 移動付与:1 の場合
		if b.IsEffectorRotation() || b.IsEffectorTranslation() {
			// n  : ボーンIndexサイズ  | 付与親ボーンのボーンIndex
			b.EffectIndex, err = r.unpackBoneIndex(model)
			if err != nil {
				mlog.E("[%d] loadBones unpackBoneIndex EffectIndex error: %v", i, err)
				return err
			}
			// 4  : float	| 付与率
			b.EffectFactor, err = r.unpackFloat()
			if err != nil {
				mlog.E("[%d] loadBones UnpackFloat EffectFactor error: %v", i, err)
				return err
			}
		} else {
			b.EffectIndex = -1
			b.EffectFactor = 0
		}

		// 軸固定:1 の場合
		if b.HasFixedAxis() {
			// 12 : float3	| 軸の方向ベクトル
			fixedAxis, err := r.unpackVec3()
			if err != nil {
				mlog.E("[%d] loadBones UnpackVec3 FixedAxis error: %v", i, err)
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
			localAxisX, err := r.unpackVec3()
			if err != nil {
				mlog.E("[%d] loadBones UnpackVec3 LocalAxisX error: %v", i, err)
				return err
			}
			b.LocalAxisX = &localAxisX
			// 12 : float3	| Z軸の方向ベクトル
			localAxisZ, err := r.unpackVec3()
			if err != nil {
				mlog.E("[%d] loadBones UnpackVec3 LocalAxisZ error: %v", i, err)
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
			b.EffectorKey, err = r.unpackInt()
			if err != nil {
				mlog.E("[%d] loadBones UnpackInt EffectorKey error: %v", i, err)
				return err
			}
		}

		// IK:1 の場合 IKデータを格納
		if b.IsIK() {
			b.Ik = pmx.NewIk()

			// n  : ボーンIndexサイズ  | IKターゲットボーンのボーンIndex
			b.Ik.BoneIndex, err = r.unpackBoneIndex(model)
			if err != nil {
				mlog.E("[%d] loadBones unpackBoneIndex Ik.BoneIndex error: %v", i, err)
				return err
			}
			// 4  : int  	| IKループ回数 (PMD及びMMD環境では255回が最大になるようです)
			b.Ik.LoopCount, err = r.unpackInt()
			if err != nil {
				mlog.E("[%d] loadBones UnpackInt Ik.LoopCount error: %v", i, err)
				return err
			}
			// 4  : float	| IKループ計算時の1回あたりの制限角度 -> ラジアン角 | PMDのIK値とは4倍異なるので注意
			unitRot, err := r.unpackFloat()
			if err != nil {
				mlog.E("[%d] loadBones UnpackFloat unitRot error: %v", i, err)
				return err
			}
			b.Ik.UnitRotation.SetRadians(&mmath.MVec3{unitRot, unitRot, unitRot})
			// 4  : int  	| IKリンク数 : 後続の要素数
			ikLinkCount, err := r.unpackInt()
			if err != nil {
				mlog.E("[%d] loadBones UnpackInt ikLinkCount error: %v", i, err)
				return err
			}
			for j := 0; j < ikLinkCount; j++ {
				il := pmx.NewIkLink()
				// n  : ボーンIndexサイズ  | リンクボーンのボーンIndex
				il.BoneIndex, err = r.unpackBoneIndex(model)
				if err != nil {
					mlog.E("[%d][%d] loadBones unpackBoneIndex IkLink.BoneIndex error: %v", i, j, err)
					return err
				}
				// 1  : byte	| 角度制限 0:OFF 1:ON
				ikLinkAngleLimit, err := r.unpackByte()
				if err != nil {
					mlog.E("[%d][%d] loadBones UnpackByte IkLink.AngleLimit error: %v", i, j, err)
					return err
				}
				il.AngleLimit = ikLinkAngleLimit == 1
				if il.AngleLimit {
					// 12 : float3	| 下限 (x,y,z) -> ラジアン角
					minAngleLimit, err := r.unpackVec3()
					if err != nil {
						mlog.E("[%d][%d] loadBones UnpackVec3 IkLink.MinAngleLimit error: %v", i, j, err)
						return err
					}
					il.MinAngleLimit.SetRadians(&minAngleLimit)
					// 12 : float3	| 上限 (x,y,z) -> ラジアン角
					maxAngleLimit, err := r.unpackVec3()
					if err != nil {
						mlog.E("[%d][%d] loadBones UnpackVec3 IkLink.MaxAngleLimit error: %v", i, j, err)
						return err
					}
					il.MaxAngleLimit.SetRadians(&maxAngleLimit)
				}
				b.Ik.Links = append(b.Ik.Links, il)
			}
		}

		bones.Update(b)
	}

	model.Bones = bones

	return nil
}

func (r *PmxRepository) loadMorphs(model *pmx.PmxModel) error {
	totalMorphCount, err := r.unpackInt()
	if err != nil {
		mlog.E("loadMorphs UnpackInt totalMorphCount error: %v", err)
		return err
	}

	morphs := pmx.NewMorphs(totalMorphCount)

	for i := 0; i < totalMorphCount; i++ {
		// 4 + n : TextBuf	| モーフ名
		name := r.readText()
		// 4 + n : TextBuf	| モーフ名英
		englishName := r.readText()

		m := &pmx.Morph{
			IndexNameModel: &core.IndexNameModel{Index: i, Name: name, EnglishName: englishName},
		}

		// 1  : byte	| 操作パネル (PMD:カテゴリ) 1:眉(左下) 2:目(左上) 3:口(右上) 4:その他(右下)  | 0:システム予約
		panel, err := r.unpackByte()
		if err != nil {
			return err
		}
		m.Panel = pmx.MorphPanel(panel)
		// 1  : byte	| モーフ種類 - 0:グループ, 1:頂点, 2:ボーン, 3:UV, 4:追加UV1, 5:追加UV2, 6:追加UV3, 7:追加UV4, 8:材質
		morphType, err := r.unpackByte()
		if err != nil {
			mlog.E("[%d] loadMorphs UnpackByte MorphType error: %v", i, err)
			return err
		}
		m.MorphType = pmx.MorphType(morphType)

		offsetCount, err := r.unpackInt()
		if err != nil {
			mlog.E("[%d] loadMorphs UnpackInt OffsetCount error: %v", i, err)
			return err
		}
		for j := 0; j < offsetCount; j++ {
			switch m.MorphType {
			case pmx.MORPH_TYPE_GROUP:
				// n  : モーフIndexサイズ  | モーフIndex  ※仕様上グループモーフのグループ化は非対応とする
				morphIndex, err := r.unpackMorphIndex(model)
				if err != nil {
					mlog.E("[%d][%d] loadMorphs unpackMorphIndex MorphIndex error: %v", i, j, err)
					return err
				}
				// 4  : float	| モーフ率 : グループモーフのモーフ値 * モーフ率 = 対象モーフのモーフ値
				morphFactor, err := r.unpackFloat()
				if err != nil {
					mlog.E("[%d][%d] loadMorphs UnpackFloat MorphFactor error: %v", i, j, err)
					return err
				}
				m.Offsets = append(m.Offsets, pmx.NewGroupMorphOffset(morphIndex, morphFactor))
			case pmx.MORPH_TYPE_VERTEX:
				// n  : 頂点Indexサイズ  | 頂点Index
				vertexIndex, err := r.unpackVertexIndex(model)
				if err != nil {
					mlog.E("[%d][%d] loadMorphs unpackVertexIndex VertexIndex error: %v", i, j, err)
					return err
				}
				// 12 : float3	| 座標オフセット量(x,y,z)
				offset, err := r.unpackVec3()
				if err != nil {
					mlog.E("[%d][%d] loadMorphs UnpackVec3 Offset error: %v", i, j, err)
					return err
				}
				m.Offsets = append(m.Offsets, pmx.NewVertexMorphOffset(vertexIndex, &offset))
			case pmx.MORPH_TYPE_BONE:
				// n  : ボーンIndexサイズ  | ボーンIndex
				boneIndex, err := r.unpackBoneIndex(model)
				if err != nil {
					mlog.E("[%d][%d] loadMorphs unpackBoneIndex BoneIndex error: %v", i, j, err)
					return err
				}
				// 12 : float3	| 移動量(x,y,z)
				offset, err := r.unpackVec3()
				if err != nil {
					mlog.E("[%d][%d] loadMorphs UnpackVec3 Offset error: %v", i, j, err)
					return err
				}
				// 16 : float4	| 回転量(x,y,z,w)
				qq, err := r.unpackQuaternion()
				if err != nil {
					mlog.E("[%d][%d] loadMorphs UnpackQuaternion Quaternion error: %v", i, j, err)
					return err
				}
				m.Offsets = append(m.Offsets, pmx.NewBoneMorphOffset(boneIndex, &offset, mmath.NewRotationFromQuaternion(&qq)))
			case pmx.MORPH_TYPE_UV, pmx.MORPH_TYPE_EXTENDED_UV1, pmx.MORPH_TYPE_EXTENDED_UV2, pmx.MORPH_TYPE_EXTENDED_UV3, pmx.MORPH_TYPE_EXTENDED_UV4:
				// n  : 頂点Indexサイズ  | 頂点Index
				vertexIndex, err := r.unpackVertexIndex(model)
				if err != nil {
					mlog.E("[%d][%d] loadMorphs unpackVertexIndex VertexIndex error: %v", i, j, err)
					return err
				}
				// 16 : float4	| UVオフセット量(x,y,z,w) ※通常UVはz,wが不要項目になるがモーフとしてのデータ値は記録しておく
				offset, err := r.unpackVec4()
				if err != nil {
					mlog.E("[%d][%d] loadMorphs UnpackVec4 Offset error: %v", i, j, err)
					return err
				}
				m.Offsets = append(m.Offsets, pmx.NewUvMorphOffset(vertexIndex, &offset))
			case pmx.MORPH_TYPE_MATERIAL:
				// n  : 材質Indexサイズ  | 材質Index -> -1:全材質対象
				materialIndex, err := r.unpackMaterialIndex(model)
				if err != nil {
					mlog.E("[%d][%d] loadMorphs unpackMaterialIndex MaterialIndex error: %v", i, j, err)
					return err
				}
				// 1  : オフセット演算形式 | 0:乗算, 1:加算 - 詳細は後述
				calcMode, err := r.unpackByte()
				if err != nil {
					mlog.E("[%d][%d] loadMorphs UnpackByte CalcMode error: %v", i, j, err)
					return err
				}
				// 16 : float4	| Diffuse (R,G,B,A) - 乗算:1.0／加算:0.0 が初期値となる(同以下)
				diffuse, err := r.unpackVec4()
				if err != nil {
					mlog.E("[%d][%d] loadMorphs UnpackVec4 Diffuse error: %v", i, j, err)
					return err
				}
				// 12 : float3	| Specular (R,G,B, Specular係数)
				specular, err := r.unpackVec4()
				if err != nil {
					mlog.E("[%d][%d] loadMorphs UnpackVec4 Specular error: %v", i, j, err)
					return err
				}
				// 12 : float3	| Ambient (R,G,B)
				ambient, err := r.unpackVec3()
				if err != nil {
					mlog.E("[%d][%d] loadMorphs UnpackVec3 Ambient error: %v", i, j, err)
					return err
				}
				// 16 : float4	| エッジ色 (R,G,B,A)
				edge, err := r.unpackVec4()
				if err != nil {
					mlog.E("[%d][%d] loadMorphs UnpackVec4 Edge error: %v", i, j, err)
					return err
				}
				// 4  : float	| エッジサイズ
				edgeSize, err := r.unpackFloat()
				if err != nil {
					mlog.E("[%d][%d] loadMorphs UnpackFloat EdgeSize error: %v", i, j, err)
					return err
				}
				// 16 : float4	| テクスチャ係数 (R,G,B,A)
				textureFactor, err := r.unpackVec4()
				if err != nil {
					mlog.E("[%d][%d] loadMorphs UnpackVec4 TextureFactor error: %v", i, j, err)
					return err
				}
				// 16 : float4	| スフィアテクスチャ係数 (R,G,B,A)
				sphereTextureFactor, err := r.unpackVec4()
				if err != nil {
					mlog.E("[%d][%d] loadMorphs UnpackVec4 SphereTextureFactor error: %v", i, j, err)
					return err
				}
				// 16 : float4	| Toonテクスチャ係数 (R,G,B,A)
				toonTextureFactor, err := r.unpackVec4()
				if err != nil {
					mlog.E("[%d][%d] loadMorphs UnpackVec4 ToonTextureFactor error: %v", i, j, err)
					return err
				}
				m.Offsets = append(m.Offsets, pmx.NewMaterialMorphOffset(
					materialIndex,
					pmx.MaterialMorphCalcMode(calcMode),
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

		morphs.Update(m)
	}

	model.Morphs = morphs

	return nil
}

func (r *PmxRepository) loadDisplaySlots(model *pmx.PmxModel) error {
	totalDisplaySlotCount, err := r.unpackInt()
	if err != nil {
		mlog.E("loadDisplaySlots UnpackInt totalDisplaySlotCount error: %v", err)
		return err
	}

	displaySlots := pmx.NewDisplaySlots(totalDisplaySlotCount)

	for i := 0; i < totalDisplaySlotCount; i++ {
		// 4 + n : TextBuf	| 枠名
		name := r.readText()
		// 4 + n : TextBuf	| 枠名英
		englishName := r.readText()

		d := &pmx.DisplaySlot{
			IndexNameModel: &core.IndexNameModel{Index: i, Name: name, EnglishName: englishName},
			References:     make([]pmx.Reference, 0),
		}

		// 1  : byte	| 特殊枠フラグ - 0:通常枠 1:特殊枠
		specialFlag, err := r.unpackByte()
		if err != nil {
			mlog.E("[%d] loadDisplaySlots UnpackByte SpecialFlag error: %v", i, err)
			return err
		}
		d.SpecialFlag = pmx.SpecialFlag(specialFlag)

		// 4  : int  	| 枠内要素数 : 後続の要素数
		referenceCount, err := r.unpackInt()
		if err != nil {
			mlog.E("[%d] loadDisplaySlots UnpackInt ReferenceCount error: %v", i, err)
			return err
		}

		for j := 0; j < referenceCount; j++ {
			reference := pmx.NewDisplaySlotReference()

			// 1  : byte	| 要素種別 - 0:ボーン 1:モーフ
			referenceType, err := r.unpackByte()
			if err != nil {
				mlog.E("[%d][%d] loadDisplaySlots UnpackByte ReferenceType error: %v", i, j, err)
				return err
			}
			reference.DisplayType = pmx.DisplayType(referenceType)

			switch reference.DisplayType {
			case pmx.DISPLAY_TYPE_BONE:
				// n  : ボーンIndexサイズ  | ボーンIndex
				reference.DisplayIndex, err = r.unpackBoneIndex(model)
				if err != nil {
					mlog.E("[%d][%d] loadDisplaySlots unpackBoneIndex DisplayIndex error: %v", i, j, err)
					return err
				}
			case pmx.DISPLAY_TYPE_MORPH:
				// n  : モーフIndexサイズ  | モーフIndex
				reference.DisplayIndex, err = r.unpackMorphIndex(model)
				if err != nil {
					mlog.E("[%d][%d] loadDisplaySlots unpackMorphIndex DisplayIndex error: %v", i, j, err)
					return err
				}
			}

			d.References = append(d.References, *reference)
		}

		displaySlots.Update(d)
	}

	model.DisplaySlots = displaySlots

	return nil
}

func (r *PmxRepository) loadRigidBodies(model *pmx.PmxModel) error {
	totalRigidBodyCount, err := r.unpackInt()
	if err != nil {
		mlog.E("loadRigidBodies UnpackInt totalRigidBodyCount error: %v", err)
		return err
	}

	rigidBodies := pmx.NewRigidBodies(totalRigidBodyCount)

	for i := 0; i < totalRigidBodyCount; i++ {
		// 4 + n : TextBuf	| 剛体名
		name := r.readText()
		// 4 + n : TextBuf	| 剛体名英
		englishName := r.readText()

		b := &pmx.RigidBody{
			IndexNameModel: &core.IndexNameModel{Index: i, Name: name, EnglishName: englishName},
			RigidBodyParam: pmx.NewRigidBodyParam(),
		}

		// n  : ボーンIndexサイズ  | 関連ボーンIndex - 関連なしの場合は-1
		b.BoneIndex, err = r.unpackBoneIndex(model)
		if err != nil {
			mlog.E("[%d] loadRigidBodies unpackBoneIndex BoneIndex error: %v", i, err)
			return err
		}
		// 1  : byte	| グループ
		collisionGroup, err := r.unpackByte()
		if err != nil {
			mlog.E("[%d] loadRigidBodies UnpackByte CollisionGroup error: %v", i, err)
			return err
		}
		b.CollisionGroup = collisionGroup
		// 2  : ushort	| 非衝突グループフラグ
		collisionGroupMask, err := r.unpackUShort()
		if err != nil {
			mlog.E("[%d] loadRigidBodies UnpackUShort CollisionGroupMask error: %v", i, err)
			return err
		}
		b.CollisionGroupMaskValue = int(collisionGroupMask)
		b.CollisionGroupMask.IsCollisions = pmx.NewCollisionGroup(collisionGroupMask)
		// 1  : byte	| 形状 - 0:球 1:箱 2:カプセル
		shapeType, err := r.unpackByte()
		if err != nil {
			mlog.E("[%d] loadRigidBodies UnpackByte ShapeType error: %v", i, err)
			return err
		}
		b.ShapeType = pmx.Shape(shapeType)
		// 12 : float3	| サイズ(x,y,z)
		size, err := r.unpackVec3()
		if err != nil {
			mlog.E("[%d] loadRigidBodies UnpackVec3 Size error: %v", i, err)
			return err
		}
		b.Size = &size
		// 12 : float3	| 位置(x,y,z)
		position, err := r.unpackVec3()
		if err != nil {
			mlog.E("[%d] loadRigidBodies UnpackVec3 Position error: %v", i, err)
			return err
		}
		b.Position = &position
		// 12 : float3	| 回転(x,y,z) -> ラジアン角
		rads, err := r.unpackVec3()
		if err != nil {
			mlog.E("[%d] loadRigidBodies UnpackVec3 Rotation error: %v", i, err)
			return err
		}
		b.Rotation = mmath.NewRotationFromRadians(&rads)
		// 4  : float	| 質量
		b.RigidBodyParam.Mass, err = r.unpackFloat()
		if err != nil {
			mlog.E("[%d] loadRigidBodies UnpackFloat Mass error: %v", i, err)
			return err
		}
		// 4  : float	| 移動減衰
		b.RigidBodyParam.LinearDamping, err = r.unpackFloat()
		if err != nil {
			mlog.E("[%d] loadRigidBodies UnpackFloat LinearDamping error: %v", i, err)
			return err
		}
		// 4  : float	| 回転減衰
		b.RigidBodyParam.AngularDamping, err = r.unpackFloat()
		if err != nil {
			mlog.E("[%d] loadRigidBodies UnpackFloat AngularDamping error: %v", i, err)
			return err
		}
		// 4  : float	| 反発力
		b.RigidBodyParam.Restitution, err = r.unpackFloat()
		if err != nil {
			mlog.E("[%d] loadRigidBodies UnpackFloat Restitution error: %v", i, err)
			return err
		}
		// 4  : float	| 摩擦力
		b.RigidBodyParam.Friction, err = r.unpackFloat()
		if err != nil {
			mlog.E("[%d] loadRigidBodies UnpackFloat Friction error: %v", i, err)
			return err
		}
		// 1  : byte	| 剛体の物理演算 - 0:ボーン追従(static) 1:物理演算(dynamic) 2:物理演算 + Bone位置合わせ
		physicsType, err := r.unpackByte()
		if err != nil {
			mlog.E("[%d] loadRigidBodies UnpackByte PhysicsType error: %v", i, err)
			return err
		}
		b.PhysicsType = pmx.PhysicsType(physicsType)

		rigidBodies.Update(b)
	}

	model.RigidBodies = rigidBodies

	return nil
}

func (r *PmxRepository) loadJoints(model *pmx.PmxModel) error {
	totalJointCount, err := r.unpackInt()
	if err != nil {
		mlog.E("loadJoints UnpackInt totalJointCount error: %v", err)
		return err
	}

	joints := pmx.NewJoints(totalJointCount)

	for i := 0; i < totalJointCount; i++ {
		// 4 + n : TextBuf	| Joint名
		name := r.readText()
		// 4 + n : TextBuf	| Joint名英
		englishName := r.readText()

		j := &pmx.Joint{
			IndexNameModel: &core.IndexNameModel{Index: i, Name: name, EnglishName: englishName},
			JointParam:     pmx.NewJointParam(),
		}

		// 1  : byte	| Joint種類 - 0:スプリング6DOF   | PMX2.0では 0 のみ(拡張用)
		j.JointType, err = r.unpackByte()
		if err != nil {
			mlog.E("[%d] loadJoints UnpackByte JointType error: %v", i, err)
			return err
		}
		// n  : 剛体Indexサイズ  | 関連剛体AのIndex - 関連なしの場合は-1
		j.RigidbodyIndexA, err = r.unpackRigidBodyIndex(model)
		if err != nil {
			mlog.E("[%d] loadJoints unpackRigidBodyIndex RigidbodyIndexA error: %v", i, err)
			return err
		}
		// n  : 剛体Indexサイズ  | 関連剛体BのIndex - 関連なしの場合は-1
		j.RigidbodyIndexB, err = r.unpackRigidBodyIndex(model)
		if err != nil {
			mlog.E("[%d] loadJoints unpackRigidBodyIndex RigidbodyIndexB error: %v", i, err)
			return err
		}
		// 12 : float3	| 位置(x,y,z)
		position, err := r.unpackVec3()
		if err != nil {
			mlog.E("[%d] loadJoints UnpackVec3 Position error: %v", i, err)
			return err
		}
		j.Position = &position
		// 12 : float3	| 回転(x,y,z) -> ラジアン角
		rads, err := r.unpackVec3()
		if err != nil {
			mlog.E("[%d] loadJoints UnpackVec3 Rotation error: %v", i, err)
			return err
		}
		j.Rotation = mmath.NewRotationFromRadians(&rads)
		// 12 : float3	| 移動制限-下限(x,y,z)
		translationLimitMin, err := r.unpackVec3()
		if err != nil {
			mlog.E("[%d] loadJoints UnpackVec3 TranslationLimitMin error: %v", i, err)
			return err
		}
		j.JointParam.TranslationLimitMin = &translationLimitMin
		// 12 : float3	| 移動制限-上限(x,y,z)
		translationLimitMax, err := r.unpackVec3()
		if err != nil {
			mlog.E("[%d] loadJoints UnpackVec3 TranslationLimitMax error: %v", i, err)
			return err
		}
		j.JointParam.TranslationLimitMax = &translationLimitMax
		// 12 : float3	| 回転制限-下限(x,y,z) -> ラジアン角
		rotationLimitMin, err := r.unpackVec3()
		if err != nil {
			mlog.E("[%d] loadJoints UnpackVec3 RotationLimitMin error: %v", i, err)
			return err
		}
		j.JointParam.RotationLimitMin = mmath.NewRotationFromRadians(&rotationLimitMin)
		// 12 : float3	| 回転制限-上限(x,y,z) -> ラジアン角
		rotationLimitMax, err := r.unpackVec3()
		if err != nil {
			mlog.E("[%d] loadJoints UnpackVec3 RotationLimitMax error: %v", i, err)
			return err
		}
		j.JointParam.RotationLimitMax = mmath.NewRotationFromRadians(&rotationLimitMax)
		// 12 : float3	| バネ定数-移動(x,y,z)
		springConstantTranslation, err := r.unpackVec3()
		if err != nil {
			mlog.E("[%d] loadJoints UnpackVec3 SpringConstantTranslation error: %v", i, err)
			return err
		}
		j.JointParam.SpringConstantTranslation = &springConstantTranslation
		// 12 : float3	| バネ定数-回転(x,y,z)
		springConstantRotation, err := r.unpackVec3()
		if err != nil {
			mlog.E("[%d] loadJoints UnpackVec3 SpringConstantRotation error: %v", i, err)
			return err
		}
		j.JointParam.SpringConstantRotation = &springConstantRotation

		joints.Update(j)
	}

	model.Joints = joints

	return nil
}

// テキストデータを読み取る
func (r *PmxRepository) unpackVertexIndex(model *pmx.PmxModel) (int, error) {
	switch model.VertexCountType {
	case 1:
		v, err := r.unpackByte()
		if err != nil {
			mlog.E("unpackVertexIndex.UnpackByte error: %v", err)
			return 0, err
		}
		return int(v), nil
	case 2:
		v, err := r.unpackUShort()
		if err != nil {
			mlog.E("unpackVertexIndex.UnpackUShort error: %v", err)
			return 0, err
		}
		return int(v), nil
	case 4:
		v, err := r.unpackInt()
		if err != nil {
			mlog.E("unpackVertexIndex.UnpackInt error: %v", err)
			return 0, err
		}
		return v, nil
	}
	return 0, fmt.Errorf("未知のVertexIndexサイズです。vertexCount: %d", model.VertexCountType)
}

// テクスチャIndexを読み取る
func (r *PmxRepository) unpackTextureIndex(model *pmx.PmxModel) (int, error) {
	return r.unpackIndex(model.TextureCountType)
}

// 材質Indexを読み取る
func (r *PmxRepository) unpackMaterialIndex(model *pmx.PmxModel) (int, error) {
	return r.unpackIndex(model.MaterialCountType)
}

// ボーンIndexを読み取る
func (r *PmxRepository) unpackBoneIndex(model *pmx.PmxModel) (int, error) {
	return r.unpackIndex(model.BoneCountType)
}

// 表情Indexを読み取る
func (r *PmxRepository) unpackMorphIndex(model *pmx.PmxModel) (int, error) {
	return r.unpackIndex(model.MorphCountType)
}

// 剛体Indexを読み取る
func (r *PmxRepository) unpackRigidBodyIndex(model *pmx.PmxModel) (int, error) {
	return r.unpackIndex(model.RigidBodyCountType)
}

func (r *PmxRepository) unpackIndex(index int) (int, error) {
	switch index {
	case 1:
		v, err := r.unpackSByte()
		if err != nil {
			mlog.E("unpackIndex.UnpackSByte error: %v", err)
			return 0, err
		}
		return int(v), nil
	case 2:
		v, err := r.unpackShort()
		if err != nil {
			mlog.E("unpackIndex.UnpackShort error: %v", err)
			return 0, err
		}
		return int(v), nil
	case 4:
		v, err := r.unpackInt()
		if err != nil {
			mlog.E("unpackIndex.UnpackInt error: %v", err)
			return 0, err
		}
		return v, nil
	}
	return 0, fmt.Errorf("未知のIndexサイズです。index: %d", index)
}
