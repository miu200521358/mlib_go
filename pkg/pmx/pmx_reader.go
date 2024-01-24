package pmx

import (
	"fmt"

	"golang.org/x/text/encoding/unicode"

	"github.com/miu200521358/mlib_go/pkg/mcore"
	"github.com/miu200521358/mlib_go/pkg/mmath"
	"github.com/miu200521358/mlib_go/pkg/mutils"
)

type PmxReader struct {
	mcore.BaseReader[*PmxModel]
}

func (r *PmxReader) createModel(path string) *PmxModel {
	model := NewPmxModel(path)
	return model
}

// 指定されたパスのファイルからデータを読み込む
func (r *PmxReader) ReadByFilepath(path string) (mcore.HashModelInterface, error) {
	// モデルを新規作成
	model := r.createModel(path)

	// ファイルを開く
	err := r.Open(path)
	if err != nil {
		return model, err
	}

	err = r.readHeader(model)
	if err != nil {
		return model, err
	}

	err = r.readData(model)
	if err != nil {
		return model, err
	}

	r.Close()

	return model, nil
}

func (r *PmxReader) ReadNameByFilepath(path string) (string, error) {
	// モデルを新規作成
	model := r.createModel(path)

	// ファイルを開く
	err := r.Open(path)
	if err != nil {
		return "", err
	}

	err = r.readHeader(model)
	if err != nil {
		return "", err
	}

	r.Close()

	return model.Name, nil
}

func (r *PmxReader) readHeader(model *PmxModel) error {
	fbytes, err := r.UnpackBytes(4)
	if err != nil {
		return err
	}
	model.Signature = r.DecodeText(unicode.UTF8, fbytes)
	model.Version, err = r.UnpackFloat()

	if err != nil {
		return err
	}

	if model.Signature[:3] != "PMX" ||
		!mutils.Contains([]string{"2.0", "2.1"}, fmt.Sprintf("%.1f", model.Version)) {
		// 整合性チェック
		return fmt.Errorf("PMX2.0/2.1形式外のデータです。signature: %s, version: %.1f ", model.Signature, model.Version)
	}

	// 1 : byte	| 後続するデータ列のバイトサイズ  PMX2.0は 8 で固定
	_, _ = r.UnpackByte()

	// [0] - エンコード方式  | 0:UTF16 1:UTF8
	encodeType, err := r.UnpackByte()
	if err != nil {
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
		return err
	}
	model.ExtendedUVCount = int(extendedUVCount)
	// [2] - 頂点Indexサイズ | 1,2,4 のいずれか
	vertexCount, err := r.UnpackByte()
	if err != nil {
		return err
	}
	model.VertexCount = int(vertexCount)
	// [3] - テクスチャIndexサイズ | 1,2,4 のいずれか
	textureCount, err := r.UnpackByte()
	if err != nil {
		return err
	}
	model.TextureCount = int(textureCount)
	// [4] - 材質Indexサイズ | 1,2,4 のいずれか
	materialCount, err := r.UnpackByte()
	if err != nil {
		return err
	}
	model.MaterialCount = int(materialCount)
	// [5] - ボーンIndexサイズ | 1,2,4 のいずれか
	boneCount, err := r.UnpackByte()
	if err != nil {
		return err
	}
	model.BoneCount = int(boneCount)
	// [6] - モーフIndexサイズ | 1,2,4 のいずれか
	morphCount, err := r.UnpackByte()
	if err != nil {
		return err
	}
	model.MorphCount = int(morphCount)
	// [7] - 剛体Indexサイズ | 1,2,4 のいずれか
	rigidBodyCount, err := r.UnpackByte()
	if err != nil {
		return err
	}
	model.RigidBodyCount = int(rigidBodyCount)

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
		return err
	}

	err = r.readFaces(model)
	if err != nil {
		return err
	}

	err = r.readTextures(model)
	if err != nil {
		return err
	}

	err = r.readMaterials(model)
	if err != nil {
		return err
	}

	err = r.readBones(model)
	if err != nil {
		return err
	}

	err = r.readMorphs(model)
	if err != nil {
		return err
	}

	err = r.readDisplaySlots(model)
	if err != nil {
		return err
	}

	err = r.readRigidBodies(model)
	if err != nil {
		return err
	}

	err = r.readJoints(model)
	if err != nil {
		return err
	}

	return nil
}

func (r *PmxReader) readVertices(model *PmxModel) error {
	totalVertexCount, err := r.UnpackInt()
	if err != nil {
		return err
	}

	for i := 0; i < totalVertexCount; i++ {
		v := NewVertex()

		// 12 : float3  | 位置(x,y,z)
		v.Position[0], err = r.UnpackFloat()
		if err != nil {
			return err
		}
		v.Position[1], err = r.UnpackFloat()
		if err != nil {
			return err
		}
		v.Position[2], err = r.UnpackFloat()
		if err != nil {
			return err
		}
		// 12 : float3  | 法線(x,y,z)
		v.Normal[0], err = r.UnpackFloat()
		if err != nil {
			return err
		}
		v.Normal[1], err = r.UnpackFloat()
		if err != nil {
			return err
		}
		v.Normal[2], err = r.UnpackFloat()
		if err != nil {
			return err
		}
		// 8  : float2  | UV(u,v)
		v.UV[0], err = r.UnpackFloat()
		if err != nil {
			return err
		}
		v.UV[1], err = r.UnpackFloat()
		if err != nil {
			return err
		}

		// 16 * n : float4[n] | 追加UV(x,y,z,w)  PMXヘッダの追加UV数による
		v.ExtendedUVs = make([]mmath.MVec4, 0)
		for j := 0; j < model.ExtendedUVCount; j++ {
			extendedUV, err := r.UnpackVec4()
			if err != nil {
				return err
			}
			v.ExtendedUVs = append(v.ExtendedUVs, extendedUV)
		}

		// 1 : byte    | ウェイト変形方式 0:BDEF1 1:BDEF2 2:BDEF4 3:SDEF
		Type, err := r.UnpackByte()
		if err != nil {
			return err
		}
		v.DeformType = DeformType(Type)

		switch v.DeformType {
		case BDEF1:
			// n : ボーンIndexサイズ  | ウェイト1.0の単一ボーン(参照Index)
			boneIndex, err := r.unpackBoneIndex(model)
			if err != nil {
				return err
			}
			v.Deform = NewBdef1(boneIndex)
		case BDEF2:
			// n : ボーンIndexサイズ  | ボーン1の参照Index
			boneIndex1, err := r.unpackBoneIndex(model)
			if err != nil {
				return err
			}
			// n : ボーンIndexサイズ  | ボーン2の参照Index
			boneIndex2, err := r.unpackBoneIndex(model)
			if err != nil {
				return err
			}
			// 4 : float              | ボーン1のウェイト値(0～1.0), ボーン2のウェイト値は 1.0-ボーン1ウェイト
			boneWeight, err := r.UnpackFloat()
			if err != nil {
				return err
			}
			v.Deform = NewBdef2(boneIndex1, boneIndex2, boneWeight)
		case BDEF4:
			// n : ボーンIndexサイズ  | ボーン1の参照Index
			boneIndex1, err := r.unpackBoneIndex(model)
			if err != nil {
				return err
			}
			// n : ボーンIndexサイズ  | ボーン2の参照Index
			boneIndex2, err := r.unpackBoneIndex(model)
			if err != nil {
				return err
			}
			// n : ボーンIndexサイズ  | ボーン3の参照Index
			boneIndex3, err := r.unpackBoneIndex(model)
			if err != nil {
				return err
			}
			// n : ボーンIndexサイズ  | ボーン4の参照Index
			boneIndex4, err := r.unpackBoneIndex(model)
			if err != nil {
				return err
			}
			// 4 : float              | ボーン1のウェイト値
			boneWeight1, err := r.UnpackFloat()
			if err != nil {
				return err
			}
			// 4 : float              | ボーン2のウェイト値
			boneWeight2, err := r.UnpackFloat()
			if err != nil {
				return err
			}
			// 4 : float              | ボーン3のウェイト値
			boneWeight3, err := r.UnpackFloat()
			if err != nil {
				return err
			}
			// 4 : float              | ボーン4のウェイト値 (ウェイト計1.0の保障はない)
			boneWeight4, err := r.UnpackFloat()
			if err != nil {
				return err
			}
			v.Deform = NewBdef4(boneIndex1, boneIndex2, boneIndex3, boneIndex4, boneWeight1, boneWeight2, boneWeight3, boneWeight4)
		case SDEF:
			// n : ボーンIndexサイズ  | ボーン1の参照Index
			boneIndex1, err := r.unpackBoneIndex(model)
			if err != nil {
				return err
			}
			// n : ボーンIndexサイズ  | ボーン2の参照Index
			boneIndex2, err := r.unpackBoneIndex(model)
			if err != nil {
				return err
			}
			// 4 : float              | ボーン1のウェイト値(0～1.0), ボーン2のウェイト値は 1.0-ボーン1ウェイト
			boneWeight, err := r.UnpackFloat()
			if err != nil {
				return err
			}
			// 12 : float3             | SDEF-C値(x,y,z)
			sdefC, err := r.UnpackVec3()
			if err != nil {
				return err
			}
			// 12 : float3             | SDEF-R0値(x,y,z)
			sdefR0, err := r.UnpackVec3()
			if err != nil {
				return err
			}
			// 12 : float3             | SDEF-R1値(x,y,z) ※修正値を要計算
			sdefR1, err := r.UnpackVec3()
			if err != nil {
				return err
			}
			v.Deform = NewSdef(boneIndex1, boneIndex2, boneWeight, sdefC, sdefR0, sdefR1)
		}

		v.EdgeFactor, err = r.UnpackFloat()
		if err != nil {
			return err
		}

		model.Vertices.Append(v, false)
	}

	return nil
}

func (r *PmxReader) readFaces(model *PmxModel) error {
	totalFaceCount, err := r.UnpackInt()
	if err != nil {
		return err
	}

	for i := 0; i < totalFaceCount; i += 3 {
		f := NewFace()

		// n : 頂点Indexサイズ     | 頂点の参照Index
		f.VertexIndexes[0], err = r.unpackVertexIndex(model)
		if err != nil {
			return err
		}

		// n : 頂点Indexサイズ     | 頂点の参照Index
		f.VertexIndexes[1], err = r.unpackVertexIndex(model)
		if err != nil {
			return err
		}

		// n : 頂点Indexサイズ     | 頂点の参照Index
		f.VertexIndexes[2], err = r.unpackVertexIndex(model)
		if err != nil {
			return err
		}

		model.Faces.Append(f, false)
	}

	return nil
}

func (r *PmxReader) readTextures(model *PmxModel) error {
	totalTextureCount, err := r.UnpackInt()
	if err != nil {
		return err
	}

	for i := 0; i < totalTextureCount; i++ {
		t := NewTexture()

		// 4 + n : TextBuf	| テクスチャパス
		t.Name = r.ReadText()

		model.Textures.Append(t, false)
	}

	return nil
}

func (r *PmxReader) readMaterials(model *PmxModel) error {
	totalMaterialCount, err := r.UnpackInt()
	if err != nil {
		return err
	}

	for i := 0; i < totalMaterialCount; i++ {
		m := NewMaterial()

		// 4 + n : TextBuf	| 材質名
		m.Name = r.ReadText()
		// 4 + n : TextBuf	| 材質名英
		m.EnglishName = r.ReadText()

		// 16 : float4	| Diffuse (R,G,B,A)
		m.DiffuseColor, err = r.UnpackVec3()
		if err != nil {
			return err
		}
		m.DiffuseAlpha, err = r.UnpackFloat()
		if err != nil {
			return err
		}
		// 12 : float3	| Specular (R,G,B)
		m.SpecularColor, err = r.UnpackVec3()
		if err != nil {
			return err
		}
		// 4  : float	| Specular係数
		m.SpecularPower, err = r.UnpackFloat()
		if err != nil {
			return err
		}
		// 12 : float3	| Ambient (R,G,B)
		m.AmbientColor, err = r.UnpackVec3()
		if err != nil {
			return err
		}
		// 1  : bitFlag  	| 描画フラグ(8bit) - 各bit 0:OFF 1:ON
		drawFlag, err := r.UnpackByte()
		if err != nil {
			return err
		}
		m.DrawFlag = DrawFlag(drawFlag)
		// 16 : float4	| エッジ色 (R,G,B,A)
		m.EdgeColor, err = r.UnpackVec3()
		if err != nil {
			return err
		}
		m.EdgeAlpha, err = r.UnpackFloat()
		if err != nil {
			return err
		}
		// 4  : float	| エッジサイズ
		m.EdgeSize, err = r.UnpackFloat()
		if err != nil {
			return err
		}
		// n  : テクスチャIndexサイズ	| 通常テクスチャ
		m.TextureIndex, err = r.unpackTextureIndex(model)
		if err != nil {
			return err
		}
		// n  : テクスチャIndexサイズ	| スフィアテクスチャ
		m.SphereTextureIndex, err = r.unpackTextureIndex(model)
		if err != nil {
			return err
		}
		// 1  : byte	| スフィアモード 0:無効 1:乗算(sph) 2:加算(spa) 3:サブテクスチャ(追加UV1のx,yをUV参照して通常テクスチャ描画を行う)
		sphereMode, err := r.UnpackByte()
		if err != nil {
			return err
		}
		m.SphereMode = SphereMode(sphereMode)
		// 1  : byte	| 共有Toonフラグ 0:継続値は個別Toon 1:継続値は共有Toon
		toonSharingFlag, err := r.UnpackByte()

		if err != nil {
			return err
		}
		m.ToonSharingFlag = ToonSharing(toonSharingFlag)

		switch m.ToonSharingFlag {
		case TOON_SHARING_INDIVIDUAL:
			// n  : テクスチャIndexサイズ	| Toonテクスチャ
			m.ToonTextureIndex, err = r.unpackTextureIndex(model)
			if err != nil {
				return err
			}
		case TOON_SHARING_SHARING:
			// 1  : byte	| 共有ToonテクスチャIndex 0～9
			toonTextureIndex, err := r.UnpackByte()
			if err != nil {
				return err
			}
			m.ToonTextureIndex = int(toonTextureIndex)
		}

		// 4 + n : TextBuf	| メモ
		m.Memo = r.ReadText()

		// 4  : int	| 材質に対応する面(頂点)数 (必ず3の倍数になる)
		m.VerticesCount, err = r.UnpackInt()
		if err != nil {
			return err
		}

		model.Materials.Append(m, false)
	}

	return nil
}

func (r *PmxReader) readBones(model *PmxModel) error {
	totalBoneCount, err := r.UnpackInt()
	if err != nil {
		return err
	}

	for i := 0; i < totalBoneCount; i++ {
		b := NewBone()

		// 4 + n : TextBuf	| ボーン名
		b.Name = r.ReadText()
		// 4 + n : TextBuf	| ボーン名英
		b.EnglishName = r.ReadText()
		// 12 : float3	| 位置
		b.Position, err = r.UnpackVec3()
		if err != nil {
			return err
		}
		// n  : ボーンIndexサイズ	| 親ボーン
		b.ParentIndex, err = r.unpackBoneIndex(model)
		if err != nil {
			return err
		}
		// 4  : int	| 変形階層
		b.Layer, err = r.UnpackInt()
		if err != nil {
			return err
		}
		// 2  : bitFlag*2	| ボーンフラグ(16bit) 各bit 0:OFF 1:ON
		boneFlag, err := r.UnpackBytes(2)
		if err != nil {
			return err
		}
		b.BoneFlag = BoneFlag(uint16(boneFlag[0]) | uint16(boneFlag[1])<<8)

		// 0x0001  : 接続先(PMD子ボーン指定)表示方法 -> 0:座標オフセットで指定 1:ボーンで指定
		if b.IsTailBone() {
			// n  : ボーンIndexサイズ  | 接続先ボーンのボーンIndex
			b.TailIndex, err = r.unpackBoneIndex(model)
			if err != nil {
				return err
			}
		} else {
			//  12 : float3	| 座標オフセット, ボーン位置からの相対分
			b.TailPosition, err = r.UnpackVec3()
			if err != nil {
				return err
			}
		}

		// 回転付与:1 または 移動付与:1 の場合
		if b.IsExternalRotation() || b.IsExternalTranslation() {
			// n  : ボーンIndexサイズ  | 付与親ボーンのボーンIndex
			b.EffectIndex, err = r.unpackBoneIndex(model)
			if err != nil {
				return err
			}
			// 4  : float	| 付与率
			b.EffectFactor, err = r.UnpackFloat()
			if err != nil {
				return err
			}
		}

		// 軸固定:1 の場合
		if b.HasFixedAxis() {
			// 12 : float3	| 軸の方向ベクトル
			b.FixedAxis, err = r.UnpackVec3()
			if err != nil {
				return err
			}
			b.NormalizeFixedAxis(b.FixedAxis.Normalized())
		}

		// ローカル軸:1 の場合
		if b.HasLocalAxis() {
			// 12 : float3	| X軸の方向ベクトル
			b.LocalAxisX, err = r.UnpackVec3()
			if err != nil {
				return err
			}
			// 12 : float3	| Z軸の方向ベクトル
			b.LocalAxisZ, err = r.UnpackVec3()
			if err != nil {
				return err
			}
			b.NormalizeLocalAxis(b.LocalAxisX.Normalized())
		}

		// 外部親変形:1 の場合
		if b.IsExternalParentDeform() {
			// 4  : int	| Key値
			b.ExternalKey, err = r.UnpackInt()
			if err != nil {
				return err
			}
		}

		// IK:1 の場合 IKデータを格納
		if b.IsIK() {
			// n  : ボーンIndexサイズ  | IKターゲットボーンのボーンIndex
			b.Ik.BoneIndex, err = r.unpackBoneIndex(model)
			if err != nil {
				return err
			}
			// 4  : int  	| IKループ回数 (PMD及びMMD環境では255回が最大になるようです)
			b.Ik.LoopCount, err = r.UnpackInt()
			if err != nil {
				return err
			}
			// 4  : float	| IKループ計算時の1回あたりの制限角度 -> ラジアン角 | PMDのIK値とは4倍異なるので注意
			unitRot, err := r.UnpackFloat()
			if err != nil {
				return err
			}
			b.Ik.UnitRotation.SetRadians(mmath.MVec3{unitRot, unitRot, unitRot})
			// 4  : int  	| IKリンク数 : 後続の要素数
			ikLinkCount, err := r.UnpackInt()
			if err != nil {
				return err
			}
			for j := 0; j < ikLinkCount; j++ {
				il := NewIkLink()
				// n  : ボーンIndexサイズ  | リンクボーンのボーンIndex
				il.BoneIndex, err = r.unpackBoneIndex(model)
				if err != nil {
					return err
				}
				// 1  : byte	| 角度制限 0:OFF 1:ON
				ikLinkAngleLimit, err := r.UnpackByte()
				if err != nil {
					return err
				}
				il.AngleLimit = ikLinkAngleLimit == 1
				if il.AngleLimit {
					// 12 : float3	| 下限 (x,y,z) -> ラジアン角
					minAngleLimit, err := r.UnpackVec3()
					if err != nil {
						return err
					}
					il.MinAngleLimit.SetRadians(minAngleLimit)
					// 12 : float3	| 上限 (x,y,z) -> ラジアン角
					maxAngleLimit, err := r.UnpackVec3()
					if err != nil {
						return err
					}
					il.MaxAngleLimit.SetRadians(maxAngleLimit)
				}
				b.Ik.Links = append(b.Ik.Links, *il)
			}
		}

		model.Bones.Append(b, false)
	}

	return nil
}

func (r *PmxReader) readMorphs(model *PmxModel) error {
	totalMorphCount, err := r.UnpackInt()
	if err != nil {
		return err
	}

	for i := 0; i < totalMorphCount; i++ {
		m := NewMorph()

		// 4 + n : TextBuf	| モーフ名
		m.Name = r.ReadText()
		// 4 + n : TextBuf	| モーフ名英
		m.EnglishName = r.ReadText()
		// 1  : byte	| 操作パネル (PMD:カテゴリ) 1:眉(左下) 2:目(左上) 3:口(右上) 4:その他(右下)  | 0:システム予約
		panel, err := r.UnpackByte()
		if err != nil {
			return err
		}
		m.Panel = MorphPanel(panel)
		// 1  : byte	| モーフ種類 - 0:グループ, 1:頂点, 2:ボーン, 3:UV, 4:追加UV1, 5:追加UV2, 6:追加UV3, 7:追加UV4, 8:材質
		morphType, err := r.UnpackByte()
		if err != nil {
			return err
		}
		m.MorphType = MorphType(morphType)

		offsetCount, err := r.UnpackInt()
		if err != nil {
			return err
		}
		for j := 0; j < offsetCount; j++ {
			switch m.MorphType {
			case MORPH_TYPE_GROUP:
				// n  : モーフIndexサイズ  | モーフIndex  ※仕様上グループモーフのグループ化は非対応とする
				morphIndex, err := r.unpackMorphIndex(model)
				if err != nil {
					return err
				}
				// 4  : float	| モーフ率 : グループモーフのモーフ値 * モーフ率 = 対象モーフのモーフ値
				morphFactor, err := r.UnpackFloat()
				if err != nil {
					return err
				}
				m.Offsets = append(m.Offsets, NewGroupMorph(morphIndex, morphFactor))
			case MORPH_TYPE_VERTEX:
				// n  : 頂点Indexサイズ  | 頂点Index
				vertexIndex, err := r.unpackVertexIndex(model)
				if err != nil {
					return err
				}
				// 12 : float3	| 座標オフセット量(x,y,z)
				offset, err := r.UnpackVec3()
				if err != nil {
					return err
				}
				m.Offsets = append(m.Offsets, NewVertexMorph(vertexIndex, offset))
			case MORPH_TYPE_BONE:
				// n  : ボーンIndexサイズ  | ボーンIndex
				boneIndex, err := r.unpackBoneIndex(model)
				if err != nil {
					return err
				}
				// 12 : float3	| 移動量(x,y,z)
				offset, err := r.UnpackVec3()
				if err != nil {
					return err
				}
				// 16 : float4	| 回転量(x,y,z,w)
				qq, err := r.UnpackQuaternion()
				if err != nil {
					return err
				}
				rotation := mmath.MRotation{}
				rotation.SetQuaternion(qq)
				m.Offsets = append(m.Offsets, NewBoneMorph(boneIndex, offset, rotation))
			case MORPH_TYPE_UV, MORPH_TYPE_EXTENDED_UV1, MORPH_TYPE_EXTENDED_UV2, MORPH_TYPE_EXTENDED_UV3, MORPH_TYPE_EXTENDED_UV4:
				// n  : 頂点Indexサイズ  | 頂点Index
				vertexIndex, err := r.unpackVertexIndex(model)
				if err != nil {
					return err
				}
				// 16 : float4	| UVオフセット量(x,y,z,w) ※通常UVはz,wが不要項目になるがモーフとしてのデータ値は記録しておく
				offset, err := r.UnpackVec4()
				if err != nil {
					return err
				}
				m.Offsets = append(m.Offsets, NewUvMorph(vertexIndex, offset))
			case MORPH_TYPE_MATERIAL:
				// n  : 材質Indexサイズ  | 材質Index -> -1:全材質対象
				materialIndex, err := r.unpackMaterialIndex(model)
				if err != nil {
					return err
				}
				// 1  : オフセット演算形式 | 0:乗算, 1:加算 - 詳細は後述
				calcMode, err := r.UnpackByte()
				if err != nil {
					return err
				}
				// 16 : float4	| Diffuse (R,G,B,A) - 乗算:1.0／加算:0.0 が初期値となる(同以下)
				diffuseColor, err := r.UnpackVec3()
				if err != nil {
					return err
				}
				diffuseAlpha, err := r.UnpackFloat()
				if err != nil {
					return err
				}
				// 12 : float3	| Specular (R,G,B)
				specularColor, err := r.UnpackVec3()
				if err != nil {
					return err
				}
				// 4  : float	| Specular係数
				specularPower, err := r.UnpackFloat()
				if err != nil {
					return err
				}
				// 12 : float3	| Ambient (R,G,B)
				ambientColor, err := r.UnpackVec3()
				if err != nil {
					return err
				}
				// 16 : float4	| エッジ色 (R,G,B,A)
				edgeColor, err := r.UnpackVec3()
				if err != nil {
					return err
				}
				edgeAlpha, err := r.UnpackFloat()
				if err != nil {
					return err
				}
				// 4  : float	| エッジサイズ
				edgeSize, err := r.UnpackFloat()
				if err != nil {
					return err
				}
				// 16 : float4	| テクスチャ係数 (R,G,B,A)
				textureCoefficient, err := r.UnpackVec3()
				if err != nil {
					return err
				}
				textureAlpha, err := r.UnpackFloat()
				if err != nil {
					return err
				}
				// 16 : float4	| スフィアテクスチャ係数 (R,G,B,A)
				sphereTextureCoefficient, err := r.UnpackVec3()
				if err != nil {
					return err
				}
				sphereTextureAlpha, err := r.UnpackFloat()
				if err != nil {
					return err
				}
				// 16 : float4	| Toonテクスチャ係数 (R,G,B,A)
				toonTextureCoefficient, err := r.UnpackVec3()
				if err != nil {
					return err
				}
				toonTextureAlpha, err := r.UnpackFloat()
				if err != nil {
					return err
				}
				m.Offsets = append(m.Offsets, NewMaterialMorph(materialIndex, MaterialMorphCalcMode(calcMode), diffuseColor, diffuseAlpha, specularColor, specularPower, ambientColor, edgeColor, edgeAlpha, edgeSize, textureCoefficient, textureAlpha, sphereTextureCoefficient, sphereTextureAlpha, toonTextureCoefficient, toonTextureAlpha))
			}
		}

		model.Morphs.Append(m, false)
	}

	return nil
}

func (r *PmxReader) readDisplaySlots(model *PmxModel) error {
	totalDisplaySlotCount, err := r.UnpackInt()
	if err != nil {
		return err
	}

	for i := 0; i < totalDisplaySlotCount; i++ {
		d := NewDisplaySlot()

		// 4 + n : TextBuf	| 枠名
		d.Name = r.ReadText()
		// 4 + n : TextBuf	| 枠名英
		d.EnglishName = r.ReadText()

		// 1  : byte	| 特殊枠フラグ - 0:通常枠 1:特殊枠
		specialFlag, err := r.UnpackByte()
		if err != nil {
			return err
		}
		d.SpecialFlag = SpecialFlag(specialFlag)

		// 4  : int  	| 枠内要素数 : 後続の要素数
		referenceCount, err := r.UnpackInt()
		if err != nil {
			return err
		}

		for j := 0; j < referenceCount; j++ {
			reference := NewDisplaySlotReference()

			// 1  : byte	| 要素種別 - 0:ボーン 1:モーフ
			referenceType, err := r.UnpackByte()
			if err != nil {
				return err
			}
			reference.DisplayType = DisplayType(referenceType)

			switch reference.DisplayType {
			case DISPLAY_TYPE_BONE:
				// n  : ボーンIndexサイズ  | ボーンIndex
				reference.DisplayIndex, err = r.unpackBoneIndex(model)
				if err != nil {
					return err
				}
			case DISPLAY_TYPE_MORPH:
				// n  : モーフIndexサイズ  | モーフIndex
				reference.DisplayIndex, err = r.unpackMorphIndex(model)
				if err != nil {
					return err
				}
			}

			d.References = append(d.References, *reference)
		}
		model.DisplaySlots.Append(d, false)
	}

	return nil
}

func (r *PmxReader) readRigidBodies(model *PmxModel) error {
	totalRigidBodyCount, err := r.UnpackInt()
	if err != nil {
		return err
	}

	for i := 0; i < totalRigidBodyCount; i++ {
		b := NewRigidBody()

		// 4 + n : TextBuf	| 剛体名
		b.Name = r.ReadText()
		// 4 + n : TextBuf	| 剛体名英
		b.EnglishName = r.ReadText()
		// n  : ボーンIndexサイズ  | 関連ボーンIndex - 関連なしの場合は-1
		b.BoneIndex, err = r.unpackBoneIndex(model)
		if err != nil {
			return err
		}
		// 1  : byte	| グループ
		collisionGroup, err := r.UnpackByte()
		if err != nil {
			return err
		}
		b.CollisionGroup = collisionGroup
		// 2  : ushort	| 非衝突グループフラグ
		collisionGroupMask, err := r.UnpackUShort()
		if err != nil {
			return err
		}
		b.CollisionGroupMask.IsCollisions = NewCollisionGroup(collisionGroupMask)
		// 1  : byte	| 形状 - 0:球 1:箱 2:カプセル
		shapeType, err := r.UnpackByte()
		if err != nil {
			return err
		}
		b.ShapeType = Shape(shapeType)
		// 12 : float3	| サイズ(x,y,z)
		b.Size, err = r.UnpackVec3()
		if err != nil {
			return err
		}
		// 12 : float3	| 位置(x,y,z)
		b.Position, err = r.UnpackVec3()
		if err != nil {
			return err
		}
		// 12 : float3	| 回転(x,y,z) -> ラジアン角
		rads, err := r.UnpackVec3()
		if err != nil {
			return err
		}
		b.Rotation.SetRadians(rads)
		// 4  : float	| 質量
		b.RigidBodyParam.Mass, err = r.UnpackFloat()
		if err != nil {
			return err
		}
		// 4  : float	| 移動減衰
		b.RigidBodyParam.LinearDamping, err = r.UnpackFloat()
		if err != nil {
			return err
		}
		// 4  : float	| 回転減衰
		b.RigidBodyParam.AngularDamping, err = r.UnpackFloat()
		if err != nil {
			return err
		}
		// 4  : float	| 反発力
		b.RigidBodyParam.Restitution, err = r.UnpackFloat()
		if err != nil {
			return err
		}
		// 4  : float	| 摩擦力
		b.RigidBodyParam.Friction, err = r.UnpackFloat()
		if err != nil {
			return err
		}
		// 1  : byte	| 剛体の物理演算 - 0:ボーン追従(static) 1:物理演算(dynamic) 2:物理演算 + Bone位置合わせ
		physicsType, err := r.UnpackByte()
		if err != nil {
			return err
		}
		b.PhysicsType = PhysicsType(physicsType)

		model.RigidBodies.Append(b, false)
	}

	return nil
}

func (r *PmxReader) readJoints(model *PmxModel) error {
	totalJointCount, err := r.UnpackInt()
	if err != nil {
		return err
	}

	for i := 0; i < totalJointCount; i++ {
		j := NewJoint()

		// 4 + n : TextBuf	| Joint名
		j.Name = r.ReadText()
		// 4 + n : TextBuf	| Joint名英
		j.EnglishName = r.ReadText()
		// 1  : byte	| Joint種類 - 0:スプリング6DOF   | PMX2.0では 0 のみ(拡張用)
		j.JointType, err = r.UnpackByte()
		if err != nil {
			return err
		}
		// n  : 剛体Indexサイズ  | 関連剛体AのIndex - 関連なしの場合は-1
		j.RigidbodyIndexA, err = r.unpackRigidBodyIndex(model)
		if err != nil {
			return err
		}
		// n  : 剛体Indexサイズ  | 関連剛体BのIndex - 関連なしの場合は-1
		j.RigidbodyIndexB, err = r.unpackRigidBodyIndex(model)
		if err != nil {
			return err
		}
		// 12 : float3	| 位置(x,y,z)
		j.Position, err = r.UnpackVec3()
		if err != nil {
			return err
		}
		// 12 : float3	| 回転(x,y,z) -> ラジアン角
		rads, err := r.UnpackVec3()
		if err != nil {
			return err
		}
		j.Rotation.SetRadians(rads)
		// 12 : float3	| 移動制限-下限(x,y,z)
		j.JointParam.TranslationLimitMin, err = r.UnpackVec3()
		if err != nil {
			return err
		}
		// 12 : float3	| 移動制限-上限(x,y,z)
		j.JointParam.TranslationLimitMax, err = r.UnpackVec3()
		if err != nil {
			return err
		}
		// 12 : float3	| 回転制限-下限(x,y,z) -> ラジアン角
		rotationLimitMin, err := r.UnpackVec3()
		if err != nil {
			return err
		}
		j.JointParam.RotationLimitMin.SetRadians(rotationLimitMin)
		// 12 : float3	| 回転制限-上限(x,y,z) -> ラジアン角
		rotationLimitMax, err := r.UnpackVec3()
		if err != nil {
			return err
		}
		j.JointParam.RotationLimitMax.SetRadians(rotationLimitMax)
		// 12 : float3	| バネ定数-移動(x,y,z)
		j.JointParam.SpringConstantTranslation, err = r.UnpackVec3()
		if err != nil {
			return err
		}
		// 12 : float3	| バネ定数-回転(x,y,z)
		springConstantRotation, err := r.UnpackVec3()
		if err != nil {
			return err
		}
		j.JointParam.SpringConstantRotation.SetDegrees(springConstantRotation)

		model.Joints.Append(j, false)
	}

	return nil
}

// テキストデータを読み取る
func (r *PmxReader) unpackVertexIndex(model *PmxModel) (int, error) {
	switch model.VertexCount {
	case 1:
		v, err := r.UnpackByte()
		if err != nil {
			return 0, err
		}
		return int(v), nil
	case 2:
		v, err := r.UnpackUShort()
		if err != nil {
			return 0, err
		}
		return int(v), nil
	case 4:
		v, err := r.UnpackInt()
		if err != nil {
			return 0, err
		}
		return v, nil
	}
	return 0, fmt.Errorf("未知のVertexIndexサイズです。vertexCount: %d", model.VertexCount)
}

// テクスチャIndexを読み取る
func (r *PmxReader) unpackTextureIndex(model *PmxModel) (int, error) {
	return r.unpackIndex(model, model.TextureCount)
}

// 材質Indexを読み取る
func (r *PmxReader) unpackMaterialIndex(model *PmxModel) (int, error) {
	return r.unpackIndex(model, model.MaterialCount)
}

// ボーンIndexを読み取る
func (r *PmxReader) unpackBoneIndex(model *PmxModel) (int, error) {
	return r.unpackIndex(model, model.BoneCount)
}

// 表情Indexを読み取る
func (r *PmxReader) unpackMorphIndex(model *PmxModel) (int, error) {
	return r.unpackIndex(model, model.MorphCount)
}

// 剛体Indexを読み取る
func (r *PmxReader) unpackRigidBodyIndex(model *PmxModel) (int, error) {
	return r.unpackIndex(model, model.RigidBodyCount)
}

func (r *PmxReader) unpackIndex(model *PmxModel, index int) (int, error) {
	switch index {
	case 1:
		v, err := r.UnpackSByte()
		if err != nil {
			return 0, err
		}
		return int(v), nil
	case 2:
		v, err := r.UnpackShort()
		if err != nil {
			return 0, err
		}
		return int(v), nil
	case 4:
		v, err := r.UnpackInt()
		if err != nil {
			return 0, err
		}
		return v, nil
	}
	return 0, fmt.Errorf("未知のIndexサイズです。index: %d", index)
}
