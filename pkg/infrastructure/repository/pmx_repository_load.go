package repository

import (
	"bufio"
	"fmt"
	"io/fs"
	"slices"
	"strings"

	"github.com/miu200521358/mlib_go/pkg/config/mi18n"
	"github.com/miu200521358/mlib_go/pkg/config/mlog"
	"github.com/miu200521358/mlib_go/pkg/config/mproc"
	"github.com/miu200521358/mlib_go/pkg/domain/core"
	"github.com/miu200521358/mlib_go/pkg/domain/mmath"
	"github.com/miu200521358/mlib_go/pkg/domain/pmx"
	"github.com/miu200521358/mlib_go/pkg/infrastructure/mfile"
	"github.com/miu200521358/mlib_go/pkg/infrastructure/mstring"
	"golang.org/x/text/encoding/unicode"
)

func (rep *PmxRepository) CanLoad(path string) (bool, error) {
	if isExist, err := mfile.ExistsFile(path); err != nil || !isExist {
		return false, fmt.Errorf("%s", mi18n.T("ファイル存在エラー", map[string]interface{}{"Path": path}))
	}

	_, _, ext := mfile.SplitPath(path)
	if strings.ToLower(ext) != ".pmx" {
		return false, fmt.Errorf("%s", mi18n.T("拡張子エラー", map[string]interface{}{"Path": path, "Ext": ".pmx"}))
	}

	return true, nil
}

// 指定されたパスのファイルからデータを読み込む
func (rep *PmxRepository) Load(path string) (core.IHashModel, error) {
	mproc.SetMaxProcess(true)
	defer mproc.SetMaxProcess(false)

	if rep.isLog {
		mlog.IL("%s", mi18n.T("読み込み開始", map[string]interface{}{"Type": "Pmx", "Path": path}))
		defer mlog.I("%s", mi18n.T("読み込み終了", map[string]interface{}{"Type": "Pmx"}))
	}

	// モデルを新規作成
	model := rep.newFunc(path)

	// ファイルを開く
	err := rep.open(path)
	if err != nil {
		mlog.E("Load.Open error: %v", err)
		return model, err
	}

	err = rep.loadHeader(model)
	if err != nil {
		mlog.E("Load.loadHeader error: %v", err)
		return model, err
	}

	err = rep.loadModel(model)
	if err != nil {
		mlog.E("Load.loadModel error: %v", err)
		return model, err
	}

	rep.close()
	model.Setup()

	return model, nil
}

// 指定されたファイルオブジェクトからデータを読み込む
func (rep *PmxRepository) LoadByFile(file fs.File) (core.IHashModel, error) {
	// モデルを新規作成
	model := rep.newFunc("")

	// ファイルを開く
	rep.file = file
	rep.reader = bufio.NewReader(rep.file)

	err := rep.loadHeader(model)
	if err != nil {
		mlog.E("ReadByFilepath.loadHeader error: %v", err)
		return model, err
	}

	err = rep.loadModel(model)
	if err != nil {
		mlog.E("ReadByFilepath.loadData error: %v", err)
		return model, err
	}

	rep.close()
	model.Setup()

	return model, nil
}

func (rep *PmxRepository) LoadName(path string) string {
	if ok, err := rep.CanLoad(path); !ok || err != nil {
		return mi18n.T("読み込み失敗")
	}

	// モデルを新規作成
	model := rep.newFunc(path)

	// ファイルを開く
	err := rep.open(path)
	if err != nil {
		return mi18n.T("読み込み失敗")
	}

	err = rep.loadHeader(model)
	if err != nil {
		return mi18n.T("読み込み失敗")
	}

	rep.close()

	return model.Name()
}

func (rep *PmxRepository) loadHeader(model *pmx.PmxModel) error {
	fbytes, err := rep.unpackBytes(4)
	if err != nil {
		return err
	}
	model.Signature = rep.decodeText(unicode.UTF8, fbytes)
	model.Version, err = rep.unpackFloat()

	if err != nil {
		return err
	}

	if model.Signature[:3] != "PMX" ||
		!slices.Contains([]string{"2.0", "2.1"}, fmt.Sprintf("%.1f", model.Version)) {
		// 整合性チェック
		return fmt.Errorf("PMX2.0/2.1形式外のデータです。signature: %s, version: %.1f ", model.Signature, model.Version)
	}

	// 1 : byte	| 後続するデータ列のバイトサイズ  PMX2.0は 8 で固定
	_, _ = rep.unpackByte()

	// [0] - エンコード方式  | 0:UTF16 1:UTF8
	encodeType, err := rep.unpackByte()
	if err != nil {
		return err
	}

	switch encodeType {
	case 0:
		rep.defineEncoding(unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM))
	case 1:
		rep.defineEncoding(unicode.UTF8)
	default:
		return fmt.Errorf("未知のエンコードタイプです。encodeType: %d", encodeType)
	}

	// [1] - 追加UV数 	| 0～4 詳細は頂点参照
	extendedUVCount, err := rep.unpackByte()
	if err != nil {
		return err
	}
	model.ExtendedUVCount = int(extendedUVCount)
	// [2] - 頂点Indexサイズ | 1,2,4 のいずれか
	vertexCount, err := rep.unpackByte()
	if err != nil {
		return err
	}
	model.VertexCountType = int(vertexCount)
	// [3] - テクスチャIndexサイズ | 1,2,4 のいずれか
	textureCount, err := rep.unpackByte()
	if err != nil {
		return err
	}
	model.TextureCountType = int(textureCount)
	// [4] - 材質Indexサイズ | 1,2,4 のいずれか
	materialCount, err := rep.unpackByte()
	if err != nil {
		return err
	}
	model.MaterialCountType = int(materialCount)
	// [5] - ボーンIndexサイズ | 1,2,4 のいずれか
	boneCount, err := rep.unpackByte()
	if err != nil {
		return err
	}
	model.BoneCountType = int(boneCount)
	// [6] - モーフIndexサイズ | 1,2,4 のいずれか
	morphCount, err := rep.unpackByte()
	if err != nil {
		return err
	}
	model.MorphCountType = int(morphCount)
	// [7] - 剛体Indexサイズ | 1,2,4 のいずれか
	rigidBodyCount, err := rep.unpackByte()
	if err != nil {
		return err
	}
	model.RigidBodyCountType = int(rigidBodyCount)

	// 4 + n : TextBuf	| モデル名
	model.SetName(rep.readText())
	// 4 + n : TextBuf	| モデル名英
	model.SetEnglishName(rep.readText())
	// 4 + n : TextBuf	| コメント
	model.Comment = rep.readText()
	// 4 + n : TextBuf	| コメント英
	model.EnglishComment = rep.readText()

	return nil
}

func (rep *PmxRepository) loadModel(model *pmx.PmxModel) error {
	err := rep.loadVertices(model)
	if err != nil {
		return err
	}

	err = rep.loadFaces(model)
	if err != nil {
		return err
	}

	err = rep.loadTextures(model)
	if err != nil {
		return err
	}

	err = rep.loadMaterials(model)
	if err != nil {
		return err
	}

	err = rep.loadBones(model)
	if err != nil {
		return err
	}

	err = rep.loadMorphs(model)
	if err != nil {
		return err
	}

	err = rep.loadDisplaySlots(model)
	if err != nil {
		return err
	}

	err = rep.loadRigidBodies(model)
	if err != nil {
		return err
	}

	err = rep.loadJoints(model)
	if err != nil {
		return err
	}

	return nil
}

func (rep *PmxRepository) loadVertices(model *pmx.PmxModel) error {
	if rep.isLog {
		defer mlog.I("%s", mi18n.T("読み込み途中完了", map[string]interface{}{"Type": mi18n.T("頂点")}))
	}

	totalVertexCount, err := rep.unpackInt()
	if err != nil {
		return err
	}

	vertices := pmx.NewVertices(totalVertexCount)

	vertexValues := make([]float64, 8)
	extendedUvs := make([]float64, model.ExtendedUVCount*4)
	bdef4Weights := make([]float64, 4)
	sdefWeights := make([]float64, 10)
	for i := range totalVertexCount {
		if i%20000 == 0 && i > 0 && rep.isLog {
			mlog.I("%s", mi18n.T("読み込み途中", map[string]interface{}{"Type": mi18n.T("頂点"), "Index": i, "Total": totalVertexCount}))
		}

		// 12 : float3  | 位置(x,y,z)
		// 12 : float3  | 法線(x,y,z)
		// 8  : float2  | UV(u,v)
		vertexValues, err = rep.unpackFloats(vertexValues, 8)
		if err != nil {
			return err
		}

		vertex := &pmx.Vertex{
			Position:    &mmath.MVec3{X: vertexValues[0], Y: vertexValues[1], Z: vertexValues[2]},
			Normal:      &mmath.MVec3{X: vertexValues[3], Y: vertexValues[4], Z: vertexValues[5]},
			Uv:          &mmath.MVec2{X: vertexValues[6], Y: vertexValues[7]},
			ExtendedUvs: make([]*mmath.MVec4, model.ExtendedUVCount),
		}
		vertex.SetIndex(i)

		// 16 * n : float4[n] | 追加UV(x,y,z,w)  PMXヘッダの追加UV数による
		if model.ExtendedUVCount > 0 {
			extendedUvs, err = rep.unpackFloats(extendedUvs, model.ExtendedUVCount*4)
			if err != nil {
				return err
			}

			for j := 0; j < model.ExtendedUVCount; j++ {
				vertex.ExtendedUvs[j] = &mmath.MVec4{X: extendedUvs[j*4], Y: extendedUvs[j*4+1],
					Z: extendedUvs[j*4+2], W: extendedUvs[j*4+3]}
			}
		}

		// 1 : byte    | ウェイト変形方式 0:BDEF1 1:BDEF2 2:BDEF4 3:SDEF
		Type, err := rep.unpackByte()
		if err != nil {
			return err
		}
		vertex.DeformType = pmx.DeformType(Type)

		switch vertex.DeformType {
		case pmx.BDEF1:
			// n : ボーンIndexサイズ  | ウェイト1.0の単一ボーン(参照Index)
			boneIndex, err := rep.unpackBoneIndex(model)
			if err != nil {
				return err
			}
			vertex.Deform = pmx.NewBdef1(boneIndex)
		case pmx.BDEF2:
			// n : ボーンIndexサイズ  | ボーン1の参照Index
			boneIndex1, err := rep.unpackBoneIndex(model)
			if err != nil {
				return err
			}
			// n : ボーンIndexサイズ  | ボーン2の参照Index
			boneIndex2, err := rep.unpackBoneIndex(model)
			if err != nil {
				return err
			}
			// 4 : float              | ボーン1のウェイト値(0～1.0), ボーン2のウェイト値は 1.0-ボーン1ウェイト
			boneWeight, err := rep.unpackFloat()
			if err != nil {
				return err
			}
			vertex.Deform = pmx.NewBdef2(boneIndex1, boneIndex2, boneWeight)
		case pmx.BDEF4:
			// n : ボーンIndexサイズ  | ボーン1の参照Index
			boneIndex1, err := rep.unpackBoneIndex(model)
			if err != nil {
				return err
			}
			// n : ボーンIndexサイズ  | ボーン2の参照Index
			boneIndex2, err := rep.unpackBoneIndex(model)
			if err != nil {
				return err
			}
			// n : ボーンIndexサイズ  | ボーン3の参照Index
			boneIndex3, err := rep.unpackBoneIndex(model)
			if err != nil {
				return err
			}
			// n : ボーンIndexサイズ  | ボーン4の参照Index
			boneIndex4, err := rep.unpackBoneIndex(model)
			if err != nil {
				return err
			}

			// 4 : float              | ボーン1のウェイト値
			// 4 : float              | ボーン2のウェイト値
			// 4 : float              | ボーン3のウェイト値
			// 4 : float              | ボーン4のウェイト値 (ウェイト計1.0の保障はない)
			bdef4Weights, err = rep.unpackFloats(bdef4Weights, 4)
			if err != nil {
				return err
			}

			vertex.Deform = pmx.NewBdef4(boneIndex1, boneIndex2, boneIndex3, boneIndex4,
				bdef4Weights[0], bdef4Weights[1], bdef4Weights[2], bdef4Weights[3])
		case pmx.SDEF:
			// n : ボーンIndexサイズ  | ボーン1の参照Index
			boneIndex1, err := rep.unpackBoneIndex(model)
			if err != nil {
				return err
			}
			// n : ボーンIndexサイズ  | ボーン2の参照Index
			boneIndex2, err := rep.unpackBoneIndex(model)
			if err != nil {
				return err
			}

			sdefWeights, err = rep.unpackFloats(sdefWeights, 10)
			if err != nil {
				return err
			}

			// 4 : float              | ボーン1のウェイト値(0～1.0), ボーン2のウェイト値は 1.0-ボーン1ウェイト
			// 12 : float3             | SDEF-C値(x,y,z)
			// 12 : float3             | SDEF-R0値(x,y,z)
			// 12 : float3             | SDEF-R1値(x,y,z) ※修正値を要計算
			vertex.Deform = pmx.NewSdef(boneIndex1, boneIndex2, sdefWeights[0],
				&mmath.MVec3{X: sdefWeights[1], Y: sdefWeights[2], Z: sdefWeights[3]},
				&mmath.MVec3{X: sdefWeights[4], Y: sdefWeights[5], Z: sdefWeights[6]},
				&mmath.MVec3{X: sdefWeights[7], Y: sdefWeights[8], Z: sdefWeights[9]})
		}

		vertex.EdgeFactor, err = rep.unpackFloat()
		if err != nil {
			return err
		}

		vertices.Update(vertex)
	}

	model.Vertices = vertices

	return nil
}

func (rep *PmxRepository) loadFaces(model *pmx.PmxModel) error {
	if rep.isLog {
		defer mlog.I("%s", mi18n.T("読み込み途中完了", map[string]interface{}{"Type": mi18n.T("面")}))
	}

	totalFaceCount, err := rep.unpackInt()
	if err != nil {
		return err
	}

	faces := pmx.NewFaces(totalFaceCount / 3)

	vertexIndexes, err := rep.unpackVertexIndexes(model, totalFaceCount)
	if err != nil {
		return err
	}

	for i := 0; i < totalFaceCount; i += 3 {
		if i%60000 == 0 && i > 0 && rep.isLog {
			mlog.I("%s", mi18n.T("読み込み途中", map[string]interface{}{"Type": mi18n.T("面"), "Index": i, "Total": totalFaceCount}))
		}

		// n : 頂点Indexサイズ     | 頂点の参照Index
		// n : 頂点Indexサイズ     | 頂点の参照Index
		// n : 頂点Indexサイズ     | 頂点の参照Index
		face := &pmx.Face{
			VertexIndexes: [3]int{
				vertexIndexes[i],
				vertexIndexes[i+1],
				vertexIndexes[i+2],
			},
		}
		face.SetIndex(int(i / 3))
		faces.Update(face)
	}

	model.Faces = faces

	return nil
}

func (rep *PmxRepository) loadTextures(model *pmx.PmxModel) error {
	totalTextureCount, err := rep.unpackInt()
	if err != nil {
		return err
	}

	textures := pmx.NewTextures(totalTextureCount)

	for i := range totalTextureCount {
		// 4 + n : TextBuf	| テクスチャパス
		tex := &pmx.Texture{}
		tex.SetIndex(i)
		tex.SetName(rep.readText())

		textures.Update(tex)
	}

	model.Textures = textures

	return nil
}

func (rep *PmxRepository) loadMaterials(model *pmx.PmxModel) error {
	if rep.isLog {
		defer mlog.I("%s", mi18n.T("読み込み途中完了", map[string]interface{}{"Type": mi18n.T("材質")}))
	}

	totalMaterialCount, err := rep.unpackInt()
	if err != nil {
		return err
	}

	materials := pmx.NewMaterials(totalMaterialCount)

	materialValues := make([]float64, 11)
	edgeValues := make([]float64, 5)
	for i := range totalMaterialCount {
		material := &pmx.Material{}
		material.SetIndex(i)
		// 4 + n : TextBuf	| 材質名
		material.SetName(rep.readText())
		// 4 + n : TextBuf	| 材質名英
		material.SetEnglishName(rep.readText())

		// 16 : float4	| Diffuse (R,G,B,A)
		// 12 : float3	| Specular (R,G,B)
		// 4: float | Specular係数
		// 12 : float3	| Ambient (R,G,B)
		materialValues, err = rep.unpackFloats(materialValues, 11)
		if err != nil {
			return err
		}

		material.Diffuse = &mmath.MVec4{X: materialValues[0], Y: materialValues[1],
			Z: materialValues[2], W: materialValues[3]}
		material.Specular = &mmath.MVec4{X: materialValues[4], Y: materialValues[5],
			Z: materialValues[6], W: materialValues[7]}
		material.Ambient = &mmath.MVec3{X: materialValues[8], Y: materialValues[9],
			Z: materialValues[10]}

		// 1  : bitFlag  	| 描画フラグ(8bit) - 各bit 0:OFF 1:ON
		drawFlag, err := rep.unpackByte()
		if err != nil {
			return err
		}
		material.DrawFlag = pmx.DrawFlag(drawFlag)

		// 16 : float4	| エッジ色 (R,G,B,A)
		// 4  : float	| エッジサイズ
		edgeValues, err = rep.unpackFloats(edgeValues, 5)
		if err != nil {
			return err
		}

		material.Edge = &mmath.MVec4{X: edgeValues[0], Y: edgeValues[1],
			Z: edgeValues[2], W: edgeValues[3]}
		material.EdgeSize = edgeValues[4]

		// n  : テクスチャIndexサイズ	| 通常テクスチャ
		material.TextureIndex, err = rep.unpackTextureIndex(model)
		if err != nil {
			return err
		}
		// n  : テクスチャIndexサイズ	| スフィアテクスチャ
		material.SphereTextureIndex, err = rep.unpackTextureIndex(model)
		if err != nil {
			return err
		}
		// 1  : byte	| スフィアモード 0:無効 1:乗算(sph) 2:加算(spa) 3:サブテクスチャ(追加UV1のx,yをUV参照して通常テクスチャ描画を行う)
		sphereMode, err := rep.unpackByte()
		if err != nil {
			return err
		}
		material.SphereMode = pmx.SphereMode(sphereMode)
		// 1  : byte	| 共有Toonフラグ 0:継続値は個別Toon 1:継続値は共有Toon
		toonSharingFlag, err := rep.unpackByte()

		if err != nil {
			return err
		}
		material.ToonSharingFlag = pmx.ToonSharing(toonSharingFlag)

		switch material.ToonSharingFlag {
		case pmx.TOON_SHARING_INDIVIDUAL:
			// n  : テクスチャIndexサイズ	| Toonテクスチャ
			material.ToonTextureIndex, err = rep.unpackTextureIndex(model)
			if err != nil {
				return err
			}
		case pmx.TOON_SHARING_SHARING:
			// 1  : byte	| 共有ToonテクスチャIndex 0～9
			toonTextureIndex, err := rep.unpackByte()
			if err != nil {
				return err
			}
			material.ToonTextureIndex = int(toonTextureIndex)
		}

		// 4 + n : TextBuf	| メモ
		material.Memo = rep.readText()

		// 4  : int	| 材質に対応する面(頂点)数 (必ず3の倍数になる)
		material.VerticesCount, err = rep.unpackInt()
		if err != nil {
			return err
		}

		materials.Update(material)

	}

	model.Materials = materials

	return nil
}

func (rep *PmxRepository) loadBones(model *pmx.PmxModel) error {
	if rep.isLog {
		defer mlog.I("%s", mi18n.T("読み込み途中完了", map[string]interface{}{"Type": mi18n.T("ボーン")}))
	}

	totalBoneCount, err := rep.unpackInt()
	if err != nil {
		return err
	}

	bones := pmx.NewBones(totalBoneCount)

	for i := range totalBoneCount {
		bone := &pmx.Bone{}
		bone.SetIndex(i)
		// 4 + n : TextBuf	| ボーン名
		bone.SetName(rep.readText())
		// 4 + n : TextBuf	| ボーン名英
		bone.SetEnglishName(rep.readText())

		// 12 : float3	| 位置
		bone.Position, err = rep.unpackVec3()
		if err != nil {
			return err
		}
		// n  : ボーンIndexサイズ	| 親ボーン
		bone.ParentIndex, err = rep.unpackBoneIndex(model)
		if err != nil {
			return err
		}
		// 4  : int	| 変形階層
		bone.Layer, err = rep.unpackInt()
		if err != nil {
			return err
		}
		bone.OriginalLayer = bone.Layer
		// 2  : bitFlag*2	| ボーンフラグ(16bit) 各bit 0:OFF 1:ON
		boneFlag, err := rep.unpackBytes(2)
		if err != nil {
			return err
		}
		bone.BoneFlag = pmx.BoneFlag(uint16(boneFlag[0]) | uint16(boneFlag[1])<<8)

		// 0x0001  : 接続先(PMD子ボーン指定)表示方法 -> 0:座標オフセットで指定 1:ボーンで指定
		if bone.IsTailBone() {
			// n  : ボーンIndexサイズ  | 接続先ボーンのボーンIndex
			bone.TailIndex, err = rep.unpackBoneIndex(model)
			if err != nil {
				return err
			}
			bone.TailPosition = mmath.NewMVec3()
		} else {
			//  12 : float3	| 座標オフセット, ボーン位置からの相対分
			bone.TailIndex = -1
			bone.TailPosition, err = rep.unpackVec3()
			if err != nil {
				return err
			}
		}

		// 回転付与:1 または 移動付与:1 の場合
		if bone.IsEffectorRotation() || bone.IsEffectorTranslation() {
			// n  : ボーンIndexサイズ  | 付与親ボーンのボーンIndex
			bone.EffectIndex, err = rep.unpackBoneIndex(model)
			if err != nil {
				return err
			}
			// 4  : float	| 付与率
			bone.EffectFactor, err = rep.unpackFloat()
			if err != nil {
				return err
			}
		} else {
			bone.EffectIndex = -1
			bone.EffectFactor = 0
		}

		// 軸固定:1 の場合
		if bone.HasFixedAxis() {
			// 12 : float3	| 軸の方向ベクトル
			bone.FixedAxis, err = rep.unpackVec3()
			if err != nil {
				return err
			}
			bone.NormalizeFixedAxis(bone.FixedAxis.Normalize())
		} else {
			bone.FixedAxis = mmath.NewMVec3()
		}

		// ローカル軸:1 の場合
		if bone.HasLocalAxis() {
			// 12 : float3	| X軸の方向ベクトル
			bone.LocalAxisX, err = rep.unpackVec3()
			if err != nil {
				return err
			}
			// 12 : float3	| Z軸の方向ベクトル
			bone.LocalAxisZ, err = rep.unpackVec3()
			if err != nil {
				return err
			}
			bone.NormalizeLocalAxis(bone.LocalAxisX)
		} else {
			bone.LocalAxisX = mmath.NewMVec3()
			bone.LocalAxisZ = mmath.NewMVec3()
		}

		// 外部親変形:1 の場合
		if bone.IsEffectorParentDeform() {
			// 4  : int	| Key値
			bone.EffectorKey, err = rep.unpackInt()
			if err != nil {
				return err
			}
		}

		// IK:1 の場合 IKデータを格納
		if bone.IsIK() {
			bone.Ik = pmx.NewIk()

			// n  : ボーンIndexサイズ  | IKターゲットボーンのボーンIndex
			bone.Ik.BoneIndex, err = rep.unpackBoneIndex(model)
			if err != nil {
				return err
			}
			// 4  : int  	| IKループ回数 (PMD及びMMD環境では255回が最大になるようです)
			bone.Ik.LoopCount, err = rep.unpackInt()
			if err != nil {
				return err
			}
			// 4  : float	| IKループ計算時の1回あたりの制限角度 -> ラジアン角 | PMDのIK値とは4倍異なるので注意
			unitRot, err := rep.unpackFloat()
			if err != nil {
				return err
			}
			bone.Ik.UnitRotation = &mmath.MVec3{X: unitRot, Y: unitRot, Z: unitRot}
			// 4  : int  	| IKリンク数 : 後続の要素数
			ikLinkCount, err := rep.unpackInt()
			if err != nil {
				return err
			}
			for j := 0; j < ikLinkCount; j++ {
				il := pmx.NewIkLink()
				// n  : ボーンIndexサイズ  | リンクボーンのボーンIndex
				il.BoneIndex, err = rep.unpackBoneIndex(model)
				if err != nil {
					return err
				}
				// 1  : byte	| 角度制限 0:OFF 1:ON
				ikLinkAngleLimit, err := rep.unpackByte()
				if err != nil {
					return err
				}
				il.AngleLimit = ikLinkAngleLimit == 1
				if il.AngleLimit {
					// 12 : float3	| 下限 (x,y,z) -> ラジアン角
					minAngleLimit, err := rep.unpackVec3()
					if err != nil {
						return err
					}
					il.MinAngleLimit = minAngleLimit
					// 12 : float3	| 上限 (x,y,z) -> ラジアン角
					maxAngleLimit, err := rep.unpackVec3()
					if err != nil {
						return err
					}
					il.MaxAngleLimit = maxAngleLimit
				}
				bone.Ik.Links = append(bone.Ik.Links, il)
			}
		}

		bones.Update(bone)
	}

	model.Bones = bones

	return nil
}

func (rep *PmxRepository) loadMorphs(model *pmx.PmxModel) error {
	if rep.isLog {
		defer mlog.I("%s", mi18n.T("読み込み途中完了", map[string]interface{}{"Type": mi18n.T("モーフ")}))
	}

	totalMorphCount, err := rep.unpackInt()
	if err != nil {
		return err
	}

	morphs := pmx.NewMorphs(totalMorphCount)

	boneOffsetValues := make([]float64, 7)
	materialOffsetValues := make([]float64, 28)
	for i := range totalMorphCount {
		morph := &pmx.Morph{}
		morph.SetIndex(i)
		// 4 + n : TextBuf	| モーフ名
		morph.SetName(rep.readText())
		// 4 + n : TextBuf	| モーフ名英
		morph.SetEnglishName(rep.readText())

		// 1  : byte	| 操作パネル (PMD:カテゴリ) 1:眉(左下) 2:目(左上) 3:口(右上) 4:その他(右下)  | 0:システム予約
		panel, err := rep.unpackByte()
		if err != nil {
			return err
		}
		morph.Panel = pmx.MorphPanel(panel)
		// 1  : byte	| モーフ種類 - 0:グループ, 1:頂点, 2:ボーン, 3:UV, 4:追加UV1, 5:追加UV2, 6:追加UV3, 7:追加UV4, 8:材質
		morphType, err := rep.unpackByte()
		if err != nil {
			return err
		}
		morph.MorphType = pmx.MorphType(morphType)

		offsetCount, err := rep.unpackInt()
		if err != nil {
			return err
		}
		for j := 0; j < offsetCount; j++ {
			switch morph.MorphType {
			case pmx.MORPH_TYPE_GROUP:
				// n  : モーフIndexサイズ  | モーフIndex  ※仕様上グループモーフのグループ化は非対応とする
				morphIndex, err := rep.unpackMorphIndex(model)
				if err != nil {
					return err
				}
				// 4  : float	| モーフ率 : グループモーフのモーフ値 * モーフ率 = 対象モーフのモーフ値
				morphFactor, err := rep.unpackFloat()
				if err != nil {
					return err
				}
				morph.Offsets = append(morph.Offsets, pmx.NewGroupMorphOffset(morphIndex, morphFactor))
			case pmx.MORPH_TYPE_VERTEX:
				// n  : 頂点Indexサイズ  | 頂点Index
				vertexIndex, err := rep.unpackVertexIndex(model)
				if err != nil {
					return err
				}
				// 12 : float3	| 座標オフセット量(x,y,z)
				offset, err := rep.unpackVec3()
				if err != nil {
					return err
				}
				morph.Offsets = append(morph.Offsets, pmx.NewVertexMorphOffset(vertexIndex, offset))
			case pmx.MORPH_TYPE_BONE:
				// n  : ボーンIndexサイズ  | ボーンIndex
				boneIndex, err := rep.unpackBoneIndex(model)
				if err != nil {
					return err
				}
				// 12 : float3	| 移動量(x,y,z)
				// 16 : float4	| 回転量(x,y,z,w)
				boneOffsetValues, err = rep.unpackFloats(boneOffsetValues, 7)
				if err != nil {
					return err
				}

				offset := pmx.NewBoneMorphOffset(boneIndex)
				offset.Position = &mmath.MVec3{X: boneOffsetValues[0], Y: boneOffsetValues[1], Z: boneOffsetValues[2]}
				offset.Rotation = mmath.NewMQuaternionByValues(
					boneOffsetValues[3], boneOffsetValues[4], boneOffsetValues[5], boneOffsetValues[6])
				morph.Offsets = append(morph.Offsets, offset)
			case pmx.MORPH_TYPE_UV, pmx.MORPH_TYPE_EXTENDED_UV1, pmx.MORPH_TYPE_EXTENDED_UV2, pmx.MORPH_TYPE_EXTENDED_UV3, pmx.MORPH_TYPE_EXTENDED_UV4:
				// n  : 頂点Indexサイズ  | 頂点Index
				vertexIndex, err := rep.unpackVertexIndex(model)
				if err != nil {
					return err
				}
				// 16 : float4	| UVオフセット量(x,y,z,w) ※通常UVはz,wが不要項目になるがモーフとしてのデータ値は記録しておく
				offset, err := rep.unpackVec4()
				if err != nil {
					return err
				}
				morph.Offsets = append(morph.Offsets, pmx.NewUvMorphOffset(vertexIndex, offset))
			case pmx.MORPH_TYPE_MATERIAL:
				// n  : 材質Indexサイズ  | 材質Index -> -1:全材質対象
				materialIndex, err := rep.unpackMaterialIndex(model)
				if err != nil {
					return err
				}
				// 1  : オフセット演算形式 | 0:乗算, 1:加算 - 詳細は後述
				calcMode, err := rep.unpackByte()
				if err != nil {
					return err
				}

				// 16 : float4	| Diffuse (R,G,B,A) - 乗算:1.0／加算:0.0 が初期値となる(同以下)
				// 16 : float4	| Specular (R,G,B, Specular係数)
				// 12 : float3	| Ambient (R,G,B)
				// 16 : float4	| エッジ色 (R,G,B,A)
				// 4  : float	| エッジサイズ
				// 16 : float4	| テクスチャ係数 (R,G,B,A)
				// 16 : float4	| スフィアテクスチャ係数 (R,G,B,A)
				// 16 : float4	| Toonテクスチャ係数 (R,G,B,A)
				materialOffsetValues, err = rep.unpackFloats(materialOffsetValues, 28)
				if err != nil {
					return err
				}
				morph.Offsets = append(morph.Offsets, pmx.NewMaterialMorphOffset(
					materialIndex,
					pmx.MaterialMorphCalcMode(calcMode),
					&mmath.MVec4{X: materialOffsetValues[0], Y: materialOffsetValues[1], Z: materialOffsetValues[2], W: materialOffsetValues[3]},
					&mmath.MVec4{X: materialOffsetValues[4], Y: materialOffsetValues[5], Z: materialOffsetValues[6], W: materialOffsetValues[7]},
					&mmath.MVec3{X: materialOffsetValues[8], Y: materialOffsetValues[9], Z: materialOffsetValues[10]},
					&mmath.MVec4{X: materialOffsetValues[11], Y: materialOffsetValues[12], Z: materialOffsetValues[13], W: materialOffsetValues[14]},
					materialOffsetValues[15],
					&mmath.MVec4{X: materialOffsetValues[16], Y: materialOffsetValues[17], Z: materialOffsetValues[18], W: materialOffsetValues[19]},
					&mmath.MVec4{X: materialOffsetValues[20], Y: materialOffsetValues[21], Z: materialOffsetValues[22], W: materialOffsetValues[23]},
					&mmath.MVec4{X: materialOffsetValues[24], Y: materialOffsetValues[25], Z: materialOffsetValues[26], W: materialOffsetValues[27]},
				))
			}
		}

		morphs.Update(morph)
	}

	model.Morphs = morphs

	return nil
}

func (rep *PmxRepository) loadDisplaySlots(model *pmx.PmxModel) error {
	if rep.isLog {
		defer mlog.I("%s", mi18n.T("読み込み途中完了", map[string]interface{}{"Type": mi18n.T("表示枠")}))
	}

	totalDisplaySlotCount, err := rep.unpackInt()
	if err != nil {
		return err
	}

	displaySlots := pmx.NewDisplaySlots(totalDisplaySlotCount)

	for i := range totalDisplaySlotCount {
		displaySlot := &pmx.DisplaySlot{
			References: make([]*pmx.Reference, 0),
		}
		displaySlot.SetIndex(i)
		// 4 + n : TextBuf	| 枠名
		displaySlot.SetName(rep.readText())
		// 4 + n : TextBuf	| 枠名英
		displaySlot.SetEnglishName(rep.readText())

		// 1  : byte	| 特殊枠フラグ - 0:通常枠 1:特殊枠
		specialFlag, err := rep.unpackByte()
		if err != nil {
			return err
		}
		displaySlot.SpecialFlag = pmx.SpecialFlag(specialFlag)

		// 4  : int  	| 枠内要素数 : 後続の要素数
		referenceCount, err := rep.unpackInt()
		if err != nil {
			return err
		}

		for j := 0; j < referenceCount; j++ {
			reference := pmx.NewDisplaySlotReference()

			// 1  : byte	| 要素種別 - 0:ボーン 1:モーフ
			referenceType, err := rep.unpackByte()
			if err != nil {
				return err
			}
			reference.DisplayType = pmx.DisplayType(referenceType)

			switch reference.DisplayType {
			case pmx.DISPLAY_TYPE_BONE:
				// n  : ボーンIndexサイズ  | ボーンIndex
				reference.DisplayIndex, err = rep.unpackBoneIndex(model)
				if err != nil {
					return err
				}
			case pmx.DISPLAY_TYPE_MORPH:
				// n  : モーフIndexサイズ  | モーフIndex
				reference.DisplayIndex, err = rep.unpackMorphIndex(model)
				if err != nil {
					return err
				}
			}

			displaySlot.References = append(displaySlot.References, reference)
		}

		displaySlots.Update(displaySlot)
	}

	model.DisplaySlots = displaySlots

	return nil
}

func (rep *PmxRepository) loadRigidBodies(model *pmx.PmxModel) error {
	if rep.isLog {
		defer mlog.I("%s", mi18n.T("読み込み途中完了", map[string]interface{}{"Type": mi18n.T("剛体")}))
	}

	totalRigidBodyCount, err := rep.unpackInt()
	if err != nil {
		return err
	}

	rigidBodies := pmx.NewRigidBodies(totalRigidBodyCount)

	rigidBodyValues := make([]float64, 14)
	for i := range totalRigidBodyCount {
		rigidBody := &pmx.RigidBody{
			RigidBodyParam: pmx.NewRigidBodyParam(),
		}
		rigidBody.SetIndex(i)
		// 4 + n : TextBuf	| 剛体名
		rigidBody.SetName(rep.readText())
		// 4 + n : TextBuf	| 剛体名英
		rigidBody.SetEnglishName(rep.readText())

		// n  : ボーンIndexサイズ  | 関連ボーンIndex - 関連なしの場合は-1
		rigidBody.BoneIndex, err = rep.unpackBoneIndex(model)
		if err != nil {
			return err
		}
		// 1  : byte	| グループ
		collisionGroup, err := rep.unpackByte()
		if err != nil {
			return err
		}
		rigidBody.CollisionGroup = collisionGroup
		// 2  : ushort	| 非衝突グループフラグ
		collisionGroupMask, err := rep.unpackUShort()
		if err != nil {
			return err
		}
		rigidBody.CollisionGroupMaskValue = int(collisionGroupMask)
		rigidBody.CollisionGroupMask.IsCollisions = pmx.NewCollisionGroup(collisionGroupMask)
		// 1  : byte	| 形状 - 0:球 1:箱 2:カプセル
		shapeType, err := rep.unpackByte()
		if err != nil {
			return err
		}
		rigidBody.ShapeType = pmx.Shape(shapeType)

		// 12 : float3	| サイズ(x,y,z)
		// 12 : float3	| 位置(x,y,z)
		// 12 : float3	| 回転(x,y,z) -> ラジアン角
		// 4  : float	| 質量
		// 4  : float	| 移動減衰
		// 4  : float	| 回転減衰
		// 4  : float	| 反発力
		// 4  : float	| 摩擦力
		rigidBodyValues, err = rep.unpackFloats(rigidBodyValues, 14)
		if err != nil {
			return err
		}

		rigidBody.Size = &mmath.MVec3{X: rigidBodyValues[0], Y: rigidBodyValues[1], Z: rigidBodyValues[2]}
		rigidBody.Position = &mmath.MVec3{X: rigidBodyValues[3], Y: rigidBodyValues[4], Z: rigidBodyValues[5]}
		rigidBody.Rotation = &mmath.MVec3{X: rigidBodyValues[6], Y: rigidBodyValues[7], Z: rigidBodyValues[8]}
		rigidBody.RigidBodyParam.Mass = rigidBodyValues[9]
		rigidBody.RigidBodyParam.LinearDamping = rigidBodyValues[10]
		rigidBody.RigidBodyParam.AngularDamping = rigidBodyValues[11]
		rigidBody.RigidBodyParam.Restitution = rigidBodyValues[12]
		rigidBody.RigidBodyParam.Friction = rigidBodyValues[13]

		// 1  : byte	| 剛体の物理演算 - 0:ボーン追従(static) 1:物理演算(dynamic) 2:物理演算 + Bone位置合わせ
		physicsType, err := rep.unpackByte()
		if err != nil {
			return err
		}
		rigidBody.PhysicsType = pmx.PhysicsType(physicsType)

		rigidBodies.Update(rigidBody)
	}

	model.RigidBodies = rigidBodies

	return nil
}

func (rep *PmxRepository) loadJoints(model *pmx.PmxModel) error {
	if rep.isLog {
		defer mlog.I("%s", mi18n.T("読み込み途中完了", map[string]interface{}{"Type": mi18n.T("ジョイント")}))
	}

	totalJointCount, err := rep.unpackInt()
	if err != nil {
		return err
	}

	joints := pmx.NewJoints(totalJointCount)

	jointValues := make([]float64, 24)
	for i := range totalJointCount {
		joint := &pmx.Joint{
			JointParam: &pmx.JointParam{},
		}
		joint.SetIndex(i)
		// 4 + n : TextBuf	| Joint名
		joint.SetName(rep.readText())
		// 4 + n : TextBuf	| Joint名英
		joint.SetEnglishName(rep.readText())

		// 1  : byte	| Joint種類 - 0:スプリング6DOF   | PMX2.0では 0 のみ(拡張用)
		joint.JointType, err = rep.unpackByte()
		if err != nil {
			return err
		}
		// n  : 剛体Indexサイズ  | 関連剛体AのIndex - 関連なしの場合は-1
		joint.RigidbodyIndexA, err = rep.unpackRigidBodyIndex(model)
		if err != nil {
			return err
		}
		// n  : 剛体Indexサイズ  | 関連剛体BのIndex - 関連なしの場合は-1
		joint.RigidbodyIndexB, err = rep.unpackRigidBodyIndex(model)
		if err != nil {
			return err
		}

		// 12 : float3	| 位置(x,y,z)
		// 12 : float3	| 回転(x,y,z) -> ラジアン角
		// 12 : float3	| 移動制限-下限(x,y,z)
		// 12 : float3	| 移動制限-上限(x,y,z)
		// 12 : float3	| 回転制限-下限(x,y,z) -> ラジアン角
		// 12 : float3	| 回転制限-上限(x,y,z) -> ラジアン角
		// 12 : float3	| バネ定数-移動(x,y,z)
		// 12 : float3	| バネ定数-回転(x,y,z)
		jointValues, err := rep.unpackFloats(jointValues, 24)
		if err != nil {
			return err
		}

		joint.Position = &mmath.MVec3{X: jointValues[0], Y: jointValues[1], Z: jointValues[2]}
		joint.Rotation = &mmath.MVec3{X: jointValues[3], Y: jointValues[4], Z: jointValues[5]}
		joint.JointParam.TranslationLimitMin = &mmath.MVec3{X: jointValues[6], Y: jointValues[7], Z: jointValues[8]}
		joint.JointParam.TranslationLimitMax = &mmath.MVec3{X: jointValues[9], Y: jointValues[10], Z: jointValues[11]}
		joint.JointParam.RotationLimitMin = &mmath.MVec3{X: jointValues[12], Y: jointValues[13], Z: jointValues[14]}
		joint.JointParam.RotationLimitMax = &mmath.MVec3{X: jointValues[15], Y: jointValues[16], Z: jointValues[17]}
		joint.JointParam.SpringConstantTranslation =
			&mmath.MVec3{X: jointValues[18], Y: jointValues[19], Z: jointValues[20]}
		joint.JointParam.SpringConstantRotation =
			&mmath.MVec3{X: jointValues[21], Y: jointValues[22], Z: jointValues[23]}

		joints.Update(joint)
	}

	model.Joints = joints

	return nil
}

// テキストデータを読み取る
func (rep *PmxRepository) unpackVertexIndex(model *pmx.PmxModel) (int, error) {
	switch model.VertexCountType {
	case 1:
		v, err := rep.unpackByte()
		if err != nil {
			return 0, err
		}
		return int(v), nil
	case 2:
		v, err := rep.unpackUShort()
		if err != nil {
			return 0, err
		}
		return int(v), nil
	case 4:
		v, err := rep.unpackInt()
		if err != nil {
			return 0, err
		}
		return v, nil
	}
	return 0, fmt.Errorf("未知のVertexIndexサイズです。vertexCount: %d", model.VertexCountType)
}

// テキストデータを読み取る
func (rep *PmxRepository) unpackVertexIndexes(model *pmx.PmxModel, count int) ([]int, error) {
	indexes := make([]int, count)

	switch model.VertexCountType {
	case 1:
		values, err := rep.unpackBytes(count)
		if err != nil {
			return indexes, err
		}

		for i, b := range values {
			indexes[i] = int(b)
		}
		return indexes, nil
	case 2:
		values, err := rep.unpackUShorts(count)
		if err != nil {
			return indexes, err
		}

		for i, s := range values {
			indexes[i] = int(s)
		}

		return indexes, nil
	case 4:
		indexes, err := rep.unpackInts(count)
		if err != nil {
			return indexes, err
		}
		return indexes, nil
	}
	return indexes, fmt.Errorf("未知のVertexIndexサイズです。vertexCount: %d", model.VertexCountType)
}

// テクスチャIndexを読み取る
func (rep *PmxRepository) unpackTextureIndex(model *pmx.PmxModel) (int, error) {
	return rep.unpackIndex(model.TextureCountType)
}

// 材質Indexを読み取る
func (rep *PmxRepository) unpackMaterialIndex(model *pmx.PmxModel) (int, error) {
	return rep.unpackIndex(model.MaterialCountType)
}

// ボーンIndexを読み取る
func (rep *PmxRepository) unpackBoneIndex(model *pmx.PmxModel) (int, error) {
	return rep.unpackIndex(model.BoneCountType)
}

// 表情Indexを読み取る
func (rep *PmxRepository) unpackMorphIndex(model *pmx.PmxModel) (int, error) {
	return rep.unpackIndex(model.MorphCountType)
}

// 剛体Indexを読み取る
func (rep *PmxRepository) unpackRigidBodyIndex(model *pmx.PmxModel) (int, error) {
	return rep.unpackIndex(model.RigidBodyCountType)
}

func (rep *PmxRepository) unpackIndex(index int) (int, error) {
	switch index {
	case 1:
		v, err := rep.unpackSByte()
		if err != nil {
			return 0, err
		}
		return int(v), nil
	case 2:
		v, err := rep.unpackShort()
		if err != nil {
			return 0, err
		}
		return int(v), nil
	case 4:
		v, err := rep.unpackInt()
		if err != nil {
			return 0, err
		}
		return v, nil
	}
	return 0, fmt.Errorf("未知のIndexサイズです。index: %d\n\n%v", index, mstring.GetStackTrace())
}
