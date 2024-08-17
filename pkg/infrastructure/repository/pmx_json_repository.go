package repository

import (
	"encoding/json"
	"os"

	"github.com/miu200521358/mlib_go/pkg/domain/core"
	"github.com/miu200521358/mlib_go/pkg/domain/mmath"
	"github.com/miu200521358/mlib_go/pkg/domain/pmx"
	"github.com/miu200521358/mlib_go/pkg/mutils/mlog"
)

type ikLinkJson struct {
	BoneIndex     int          `json:"bone_index"`  // リンクボーンのボーンIndex
	AngleLimit    bool         `json:"angle_limit"` // 角度制限有無
	MinAngleLimit *mmath.MVec3 `json:"min_angle"`   // 下限
	MaxAngleLimit *mmath.MVec3 `json:"max_angle"`   // 上限
}

type ikJson struct {
	BoneIndex    int           `json:"bone_index"`    // IKターゲットボーンのボーンIndex
	LoopCount    int           `json:"loop_count"`    // IKループ回数 (最大255)
	UnitRotation float64       `json:"unit_rotation"` // IKループ計算時の1回あたりの制限角度(ラジアン)
	Links        []*ikLinkJson `json:"links"`         // IKリンクリスト
}

type boneJson struct {
	Index        int          `json:"index"`         // ボーンINDEX
	Name         string       `json:"name"`          // ボーン名
	EnglishName  string       `json:"english_name"`  // ボーン英名
	Position     *mmath.MVec3 `json:"position"`      // 位置
	ParentIndex  int          `json:"parent_index"`  // 親ボーンのボーンIndex(親がない場合は-1)
	Layer        int          `json:"layer"`         // 変形階層
	BoneFlag     int          `json:"bone_flag"`     // ボーンフラグ(16bit) 各bit 0:OFF 1:ON
	TailPosition *mmath.MVec3 `json:"tail_position"` // 接続先:0 の場合 座標オフセット, ボーン位置からの相対分
	TailIndex    int          `json:"tail_index"`    // 接続先:1 の場合 接続先ボーンのボーンIndex
	EffectIndex  int          `json:"effect_index"`  // 回転付与:1 または 移動付与:1 の場合 付与親ボーンのボーンIndex
	EffectFactor float64      `json:"effect_factor"` // 付与率
	FixedAxis    *mmath.MVec3 `json:"fixed_axis"`    // 軸固定:1 の場合 軸の方向ベクトル
	LocalAxisX   *mmath.MVec3 `json:"local_axis_x"`  // ローカル軸:1 の場合 X軸の方向ベクトル
	LocalAxisZ   *mmath.MVec3 `json:"local_axis_z"`  // ローカル軸:1 の場合 Z軸の方向ベクトル
	EffectorKey  int          `json:"effector_key"`  // 外部親変形:1 の場合 Key値
	Ik           *ikJson      `json:"ik"`            // IK:1 の場合 IKデータを格納
}

type pmxJson struct {
	Name  string
	Bones []*boneJson
}

type PmxJsonRepository struct {
	*baseRepository[*pmx.PmxModel]
}

func NewPmxJsonRepository() *PmxJsonRepository {
	return &PmxJsonRepository{
		baseRepository: &baseRepository[*pmx.PmxModel]{
			newFunc: func(path string) *pmx.PmxModel {
				return pmx.NewPmxModel(path)
			},
		},
	}
}

func (rep *PmxJsonRepository) Save(overridePath string, data core.IHashModel, includeSystem bool) error {
	model := data.(*pmx.PmxModel)

	// モデルをJSONに変換
	jsonData := pmxJson{
		Name:  model.Name(),
		Bones: make([]*boneJson, 0),
	}

	for i := range model.Bones.Len() {
		bone := model.Bones.Get(i)
		boneData := boneJson{
			Index:        bone.Index(),
			Name:         bone.Name(),
			EnglishName:  bone.EnglishName(),
			Position:     bone.Position,
			ParentIndex:  bone.ParentIndex,
			Layer:        bone.Layer,
			BoneFlag:     int(bone.BoneFlag),
			TailPosition: bone.TailPosition,
			TailIndex:    bone.TailIndex,
			EffectIndex:  bone.EffectIndex,
			EffectFactor: bone.EffectFactor,
			FixedAxis:    bone.FixedAxis,
			LocalAxisX:   bone.LocalAxisX,
			LocalAxisZ:   bone.LocalAxisZ,
			EffectorKey:  bone.EffectorKey,
		}

		if bone.Ik != nil {
			ikData := ikJson{
				BoneIndex:    bone.Ik.BoneIndex,
				LoopCount:    bone.Ik.LoopCount,
				UnitRotation: bone.Ik.UnitRotation.Radians().X,
				Links:        make([]*ikLinkJson, 0),
			}

			for _, link := range bone.Ik.Links {
				linkData := ikLinkJson{
					BoneIndex:     link.BoneIndex,
					AngleLimit:    link.AngleLimit,
					MinAngleLimit: link.MinAngleLimit.Radians(),
					MaxAngleLimit: link.MaxAngleLimit.Radians(),
				}
				ikData.Links = append(ikData.Links, &linkData)
			}

			boneData.Ik = &ikData
		}

		jsonData.Bones = append(jsonData.Bones, &boneData)
	}

	// JSONに変換
	jsonText, err := json.MarshalIndent(jsonData, "", "  ")
	if err != nil {
		mlog.E("Save.Save error: %v", err)
		return err
	}

	// ファイルに書き込み
	if err := os.WriteFile(overridePath, jsonText, 0666); err != nil {
		mlog.E("Save.Save error: %v", err)
		return err
	}

	return nil
}

// 指定されたパスのファイルからデータを読み込む
func (rep *PmxJsonRepository) Load(path string) (core.IHashModel, error) {
	// モデルを新規作成
	model := rep.newFunc(path)

	// ファイルを開く
	jsonText, err := os.ReadFile(path)
	if err != nil {
		mlog.E("Load.Load error: %v", err)
		return model, err
	}

	// JSON読み込み
	var jsonData pmxJson
	if err := json.Unmarshal(jsonText, &jsonData); err != nil {
		mlog.E("Load.Load error: %v", err)
		return model, err
	}

	model, err = rep.loadModel(model, &jsonData)
	if err != nil {
		mlog.E("Load.readData error: %v", err)
		return model, err
	}

	model.UpdateHash()

	return model, nil
}

func (rep *PmxJsonRepository) LoadName(path string) (string, error) {
	// ファイルを開く
	jsonText, err := os.ReadFile(path)
	if err != nil {
		mlog.E("Load.Load error: %v", err)
		return "", err
	}

	// JSON読み込み
	var jsonData pmxJson
	if err := json.Unmarshal(jsonText, &jsonData); err != nil {
		mlog.E("Load.Load error: %v", err)
		return "", err
	}

	return jsonData.Name, nil
}

func (rep *PmxJsonRepository) loadModel(model *pmx.PmxModel, jsonData *pmxJson) (*pmx.PmxModel, error) {

	for _, boneData := range jsonData.Bones {
		bone := pmx.NewBone()
		bone.SetIndex(boneData.Index)
		bone.SetName(boneData.Name)
		bone.SetEnglishName(boneData.EnglishName)
		bone.Position = boneData.Position
		bone.ParentIndex = boneData.ParentIndex
		bone.Layer = boneData.Layer
		bone.BoneFlag = pmx.BoneFlag(uint16(boneData.BoneFlag))
		bone.TailPosition = boneData.TailPosition
		bone.TailIndex = boneData.TailIndex
		bone.EffectIndex = boneData.EffectIndex
		bone.EffectFactor = boneData.EffectFactor
		bone.FixedAxis = boneData.FixedAxis
		bone.LocalAxisX = boneData.LocalAxisX
		bone.LocalAxisZ = boneData.LocalAxisZ
		bone.EffectorKey = boneData.EffectorKey

		if boneData.Ik != nil {
			ik := pmx.NewIk()
			ik.BoneIndex = boneData.Ik.BoneIndex
			ik.LoopCount = boneData.Ik.LoopCount
			ik.UnitRotation.SetRadians(&mmath.MVec3{X: boneData.Ik.UnitRotation})
			for _, linkData := range boneData.Ik.Links {
				link := pmx.NewIkLink()
				link.BoneIndex = linkData.BoneIndex
				link.AngleLimit = linkData.AngleLimit
				link.MinAngleLimit.SetRadians(linkData.MinAngleLimit)
				link.MaxAngleLimit.SetRadians(linkData.MaxAngleLimit)
				ik.Links = append(ik.Links, link)
			}
			bone.Ik = ik
		}

		model.Bones.Append(bone)
	}

	model.Setup()

	return model, nil
}
