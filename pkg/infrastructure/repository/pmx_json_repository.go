package repository

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/miu200521358/mlib_go/pkg/config/mi18n"
	"github.com/miu200521358/mlib_go/pkg/config/mlog"
	"github.com/miu200521358/mlib_go/pkg/config/mproc"
	"github.com/miu200521358/mlib_go/pkg/domain/core"
	"github.com/miu200521358/mlib_go/pkg/domain/mmath"
	"github.com/miu200521358/mlib_go/pkg/domain/pmx"
	"github.com/miu200521358/mlib_go/pkg/infrastructure/mfile"
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

type referenceJson struct {
	DisplayType  int `json:"display_type"`  // 要素対象 0:ボーン 1:モーフ
	DisplayIndex int `json:"display_index"` // ボーンIndex or モーフIndex
}

type displaySlotJson struct {
	Index       int              `json:"index"`        // 表示枠INDEX
	Name        string           `json:"name"`         // 表示枠名
	EnglishName string           `json:"english_name"` // 表示枠英名
	SpecialFlag int              `json:"special_flag"` // 特殊枠フラグ - 0:通常枠 1:特殊枠
	References  []*referenceJson `json:"references"`   // 表示枠要素
}

type rigidBodyJson struct {
	Index              int          `json:"index"`                // 剛体INDEX
	Name               string       `json:"name"`                 // 剛体名
	EnglishName        string       `json:"english_name"`         // 剛体英名
	BoneIndex          int          `json:"bone_index"`           // 関連ボーンIndex
	CollisionGroup     int          `json:"collision_group"`      // グループ
	CollisionGroupMask int          `json:"collision_group_mask"` // 非衝突グループフラグ
	ShapeType          int          `json:"shape_type"`           // 形状
	Size               *mmath.MVec3 `json:"size"`                 // サイズ(x,y,z)
	Position           *mmath.MVec3 `json:"position"`             // 位置(x,y,z)
	Rotation           *mmath.MVec3 `json:"rotation"`             // 回転(x,y,z) -> ラジアン角
	Mass               float64      `json:"mass"`                 // 質量
	LinearDamping      float64      `json:"linear_damping"`       // 移動減衰
	AngularDamping     float64      `json:"angular_damping"`      // 回転減衰
	Restitution        float64      `json:"restitution"`          // 反発力
	Friction           float64      `json:"friction"`             // 摩擦力
	PhysicsType        int          `json:"physics_type"`         // 剛体の物理演算
}

type pmxJson struct {
	Name         string
	Bones        []*boneJson
	DisplaySlots []*displaySlotJson
	RigidBodies  []*rigidBodyJson
	// TODO モーフ
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

func (rep *PmxJsonRepository) CanLoad(path string) (bool, error) {
	if isExist, err := mfile.ExistsFile(path); err != nil || !isExist {
		return false, fmt.Errorf("%s", mi18n.T("ファイル存在エラー", map[string]interface{}{"Path": path}))
	}

	_, _, ext := mfile.SplitPath(path)
	if strings.ToLower(ext) != ".json" && strings.ToLower(ext) != ".pmx" {
		return false, fmt.Errorf("%s", mi18n.T("拡張子エラー", map[string]interface{}{"Path": path, "Ext": ".json, .pmx"}))
	}

	return true, nil
}

func (rep *PmxJsonRepository) Save(overridePath string, data core.IHashModel, includeSystem bool) error {
	mproc.SetMaxProcess(true)
	defer mproc.SetMaxProcess(false)

	mlog.IL("%s", mi18n.T("保存開始", map[string]interface{}{"Type": "Json", "Path": overridePath}))
	defer mlog.I("%s", mi18n.T("保存終了", map[string]interface{}{"Type": "Json"}))

	model := data.(*pmx.PmxModel)

	// モデルをJSONに変換
	jsonData := pmxJson{
		Name:         model.Name(),
		Bones:        make([]*boneJson, 0),
		DisplaySlots: make([]*displaySlotJson, 0),
	}

	// 頂点をボーンINDEX別に纏める
	allBoneVertices := model.Vertices.GetMapByBoneIndex(0.0)

	model.Bones.ForEach(func(index int, bone *pmx.Bone) bool {
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
				UnitRotation: bone.Ik.UnitRotation.X,
				Links:        make([]*ikLinkJson, 0),
			}

			for _, link := range bone.Ik.Links {
				linkData := ikLinkJson{
					BoneIndex:     link.BoneIndex,
					AngleLimit:    link.AngleLimit,
					MinAngleLimit: link.MinAngleLimit,
					MaxAngleLimit: link.MaxAngleLimit,
				}
				ikData.Links = append(ikData.Links, &linkData)
			}

			boneData.Ik = &ikData
		}

		jsonData.Bones = append(jsonData.Bones, &boneData)

		if bone.Config() == nil && bone.HasDynamicPhysics() {
			// 準標準ボーンではなく、物理剛体に紐付いていない場合、親の中で最も近い準標準ボーンに頂点INDEXリストを載せ替える
			for _, parentIndex := range bone.ParentBoneIndexes {
				parentBone, _ := model.Bones.Get(parentIndex)
				if parentBone.Config() != nil && parentBone.Name() == pmx.HEAD.String() {
					// 親が準標準ボーンの場合、頂点INDEXリストを載せ替える
					allBoneVertices[parentIndex] = append(allBoneVertices[parentIndex], allBoneVertices[bone.Index()]...)
					break
				}
			}

		}

		return true
	})

	// 表示枠をJSONに変換
	model.DisplaySlots.ForEach(func(index int, displaySlot *pmx.DisplaySlot) bool {
		displaySlotData := displaySlotJson{
			Index:       displaySlot.Index(),
			Name:        displaySlot.Name(),
			EnglishName: displaySlot.EnglishName(),
			SpecialFlag: int(displaySlot.SpecialFlag),
			References:  make([]*referenceJson, 0),
		}

		for _, reference := range displaySlot.References {
			referenceData := referenceJson{
				DisplayType:  int(reference.DisplayType),
				DisplayIndex: reference.DisplayIndex,
			}
			displaySlotData.References = append(displaySlotData.References, &referenceData)
		}

		jsonData.DisplaySlots = append(jsonData.DisplaySlots, &displaySlotData)
		return true
	})

	// システム用剛体をJSONに変換
	model.RigidBodies.ForEach(func(index int, rb *pmx.RigidBody) bool {
		if !strings.Contains(rb.Name(), pmx.MLIB_PREFIX) {
			return true
		}

		rigidBody := &rigidBodyJson{
			Index:              len(jsonData.RigidBodies),
			Name:               rb.Name(),
			EnglishName:        rb.EnglishName(),
			BoneIndex:          rb.Index(),
			CollisionGroup:     0,
			CollisionGroupMask: 0,
			Mass:               0.0,
			LinearDamping:      0.0,
			AngularDamping:     0.0,
			Restitution:        0.0,
			Friction:           0.0,
			PhysicsType:        int(pmx.PHYSICS_TYPE_STATIC),
			ShapeType:          int(pmx.SHAPE_SPHERE),
			Size:               rb.Size,
			Position:           rb.Position,
			Rotation:           rb.Rotation,
		}

		jsonData.RigidBodies = append(jsonData.RigidBodies, rigidBody)

		return true
	})

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
	mproc.SetMaxProcess(true)
	defer mproc.SetMaxProcess(false)

	mlog.IL("%s", mi18n.T("読み込み開始", map[string]interface{}{"Type": "Json", "Path": path}))
	defer mlog.I("%s", mi18n.T("読み込み終了", map[string]interface{}{"Type": "Json"}))

	// モデルを新規作成
	model := rep.newFunc(path)

	// ファイルを開く
	jsonText, err := os.ReadFile(path)
	if err != nil {
		mlog.E("Load.ReadFile error: %v", err)
		return model, err
	}

	// JSON読み込み
	var jsonData pmxJson
	if err := json.Unmarshal(jsonText, &jsonData); err != nil {
		mlog.E("Load.Unmarshal error: %v", err)
		return model, err
	}

	model, err = rep.loadModel(model, &jsonData)
	if err != nil {
		mlog.E("Load.loadModel error: %v", err)
		return model, err
	}

	model.UpdateHash()

	return model, nil
}

func (rep *PmxJsonRepository) LoadName(path string) string {
	if ok, err := rep.CanLoad(path); !ok || err != nil {
		return mi18n.T("読み込み失敗")
	}

	// ファイルを開く
	jsonText, err := os.ReadFile(path)
	if err != nil {
		// mlog.E("Load.Load error: %v", err)
		return mi18n.T("読み込み失敗")
	}

	// JSON読み込み
	var jsonData pmxJson
	if err := json.Unmarshal(jsonText, &jsonData); err != nil {
		// mlog.E("Load.Load error: %v", err)
		return mi18n.T("読み込み失敗")
	}

	return jsonData.Name
}

func (rep *PmxJsonRepository) loadModel(model *pmx.PmxModel, jsonData *pmxJson) (*pmx.PmxModel, error) {

	for _, boneData := range jsonData.Bones {
		bone := &pmx.Bone{}
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
		bone.OriginalLayer = bone.Layer

		if boneData.Ik != nil {
			ik := pmx.NewIk()
			ik.BoneIndex = boneData.Ik.BoneIndex
			ik.LoopCount = boneData.Ik.LoopCount
			ik.UnitRotation = &mmath.MVec3{X: boneData.Ik.UnitRotation}
			for _, linkData := range boneData.Ik.Links {
				link := pmx.NewIkLink()
				link.BoneIndex = linkData.BoneIndex
				link.AngleLimit = linkData.AngleLimit
				link.MinAngleLimit = linkData.MinAngleLimit
				link.MaxAngleLimit = linkData.MaxAngleLimit
				ik.Links = append(ik.Links, link)
			}
			bone.Ik = ik
		}

		model.Bones.Append(bone)
	}

	for _, displaySlotData := range jsonData.DisplaySlots {
		displaySlot := &pmx.DisplaySlot{}
		displaySlot.SetIndex(displaySlotData.Index)
		displaySlot.SetName(displaySlotData.Name)
		displaySlot.SetEnglishName(displaySlotData.EnglishName)
		displaySlot.SpecialFlag = pmx.SpecialFlag(displaySlotData.SpecialFlag)
		displaySlot.References = make([]*pmx.Reference, 0)
		for _, referenceData := range displaySlotData.References {
			reference := pmx.NewDisplaySlotReferenceByValues(
				pmx.DisplayType(referenceData.DisplayType), referenceData.DisplayIndex)
			displaySlot.References = append(displaySlot.References, reference)
		}
		model.DisplaySlots.Append(displaySlot)
	}

	for _, rigidBodyData := range jsonData.RigidBodies {
		rigidBody := &pmx.RigidBody{}
		rigidBody.SetIndex(rigidBodyData.Index)
		rigidBody.SetName(rigidBodyData.Name)
		rigidBody.SetEnglishName(rigidBodyData.EnglishName)
		rigidBody.BoneIndex = rigidBodyData.BoneIndex
		rigidBody.CollisionGroup = byte(rigidBodyData.CollisionGroup)
		rigidBody.CollisionGroupMaskValue = int(rigidBodyData.CollisionGroupMask)
		rigidBody.CollisionGroupMask.IsCollisions = pmx.NewCollisionGroup(uint16(rigidBodyData.CollisionGroupMask))
		rigidBody.ShapeType = pmx.Shape(rigidBodyData.ShapeType)
		rigidBody.Size = rigidBodyData.Size
		rigidBody.Position = rigidBodyData.Position
		rigidBody.Rotation = rigidBodyData.Rotation
		rigidBody.RigidBodyParam = pmx.NewRigidBodyParam()
		rigidBody.RigidBodyParam.Mass = rigidBodyData.Mass
		rigidBody.RigidBodyParam.LinearDamping = rigidBodyData.LinearDamping
		rigidBody.RigidBodyParam.AngularDamping = rigidBodyData.AngularDamping
		rigidBody.RigidBodyParam.Restitution = rigidBodyData.Restitution

		model.RigidBodies.Append(rigidBody)
	}

	model.Setup()

	return model, nil
}
