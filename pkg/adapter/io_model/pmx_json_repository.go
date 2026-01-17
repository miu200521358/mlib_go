// 指示: miu200521358
package io_model

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/miu200521358/mlib_go/pkg/adapter/io_common"
	"github.com/miu200521358/mlib_go/pkg/domain/mmath"
	"github.com/miu200521358/mlib_go/pkg/domain/model"
	"github.com/miu200521358/mlib_go/pkg/shared/hashable"
	"gonum.org/v1/gonum/spatial/r3"
)

type ikLinkJSON struct {
	BoneIndex     int        `json:"bone_index"`
	AngleLimit    bool       `json:"angle_limit"`
	MinAngleLimit mmath.Vec3 `json:"min_angle"`
	MaxAngleLimit mmath.Vec3 `json:"max_angle"`
}

type ikJSON struct {
	BoneIndex    int          `json:"bone_index"`
	LoopCount    int          `json:"loop_count"`
	UnitRotation float64      `json:"unit_rotation"`
	Links        []ikLinkJSON `json:"links"`
}

type boneJSON struct {
	Index        int        `json:"index"`
	Name         string     `json:"name"`
	EnglishName  string     `json:"english_name"`
	Position     mmath.Vec3 `json:"position"`
	ParentIndex  int        `json:"parent_index"`
	Layer        int        `json:"layer"`
	BoneFlag     int        `json:"bone_flag"`
	TailPosition mmath.Vec3 `json:"tail_position"`
	TailIndex    int        `json:"tail_index"`
	EffectIndex  int        `json:"effect_index"`
	EffectFactor float64    `json:"effect_factor"`
	FixedAxis    mmath.Vec3 `json:"fixed_axis"`
	LocalAxisX   mmath.Vec3 `json:"local_axis_x"`
	LocalAxisZ   mmath.Vec3 `json:"local_axis_z"`
	EffectorKey  int        `json:"effector_key"`
	Ik           *ikJSON    `json:"ik"`
}

type referenceJSON struct {
	DisplayType  int `json:"display_type"`
	DisplayIndex int `json:"display_index"`
}

type displaySlotJSON struct {
	Index       int             `json:"index"`
	Name        string          `json:"name"`
	EnglishName string          `json:"english_name"`
	SpecialFlag int             `json:"special_flag"`
	References  []referenceJSON `json:"references"`
}

type rigidBodyJSON struct {
	Index              int        `json:"index"`
	Name               string     `json:"name"`
	EnglishName        string     `json:"english_name"`
	BoneIndex          int        `json:"bone_index"`
	CollisionGroup     int        `json:"collision_group"`
	CollisionGroupMask int        `json:"collision_group_mask"`
	ShapeType          int        `json:"shape_type"`
	Size               mmath.Vec3 `json:"size"`
	Position           mmath.Vec3 `json:"position"`
	Rotation           mmath.Vec3 `json:"rotation"`
	Mass               float64    `json:"mass"`
	LinearDamping      float64    `json:"linear_damping"`
	AngularDamping     float64    `json:"angular_damping"`
	Restitution        float64    `json:"restitution"`
	Friction           float64    `json:"friction"`
	PhysicsType        int        `json:"physics_type"`
}

type pmxJSON struct {
	Name         string            `json:"Name"`
	Bones        []boneJSON        `json:"Bones"`
	DisplaySlots []displaySlotJSON `json:"DisplaySlots"`
	RigidBodies  []rigidBodyJSON   `json:"RigidBodies"`
}

// PmxJsonRepository はPMX JSONの入出力を表す。
type PmxJsonRepository struct{}

// NewPmxJsonRepository はPmxJsonRepositoryを生成する。
func NewPmxJsonRepository() *PmxJsonRepository {
	return &PmxJsonRepository{}
}

// CanLoad は拡張子に応じて読み込み可否を判定する。
func (r *PmxJsonRepository) CanLoad(path string) bool {
	return strings.EqualFold(filepath.Ext(path), ".json")
}

// InferName はパスから表示名を推定する。
func (r *PmxJsonRepository) InferName(path string) string {
	base := filepath.Base(path)
	ext := filepath.Ext(base)
	if ext == "" {
		return base
	}
	return strings.TrimSuffix(base, ext)
}

// Load はPMX JSONを読み込む。
func (r *PmxJsonRepository) Load(path string) (hashable.IHashable, error) {
	if !r.CanLoad(path) {
		return nil, io_common.NewIoExtInvalid(path, nil)
	}
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, io_common.NewIoFileNotFound(path, err)
		}
		return nil, io_common.NewIoParseFailed("JSONファイル情報の取得に失敗しました", err)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, io_common.NewIoParseFailed("JSONファイルの読み込みに失敗しました", err)
	}
	var jsonData pmxJSON
	if err := json.Unmarshal(data, &jsonData); err != nil {
		return nil, io_common.NewIoParseFailed("JSONの解析に失敗しました", err)
	}
	modelData := model.NewPmxModel()
	modelData.SetName(jsonData.Name)
	modelData.SetPath(path)
	modelData.SetFileModTime(info.ModTime().UnixNano())
	if err := loadBonesFromJSON(modelData, jsonData.Bones); err != nil {
		return nil, err
	}
	if err := loadDisplaySlotsFromJSON(modelData, jsonData.DisplaySlots); err != nil {
		return nil, err
	}
	if err := loadRigidBodiesFromJSON(modelData, jsonData.RigidBodies); err != nil {
		return nil, err
	}
	modelData.UpdateHash()
	return modelData, nil
}

// Save はPMX JSONを保存する。
func (r *PmxJsonRepository) Save(path string, data hashable.IHashable, opts io_common.SaveOptions) error {
	modelData, ok := data.(*model.PmxModel)
	if !ok {
		return io_common.NewIoEncodeFailed("JSON保存対象が不正です", nil)
	}
	savePath := path
	if savePath == "" {
		savePath = modelData.Path()
	}
	if savePath == "" {
		return io_common.NewIoSaveFailed("保存先パスが空です", nil)
	}

	bonesByIndex := make([]*model.Bone, modelData.Bones.Len())
	for i := 0; i < len(bonesByIndex); i++ {
		bone, err := modelData.Bones.Get(i)
		if err != nil {
			continue
		}
		bonesByIndex[i] = bone
	}
	boneMap := buildIndexMappingJSON(len(bonesByIndex), func(index int) bool {
		if opts.IncludeSystem {
			return true
		}
		bone := bonesByIndex[index]
		return bone != nil && !bone.IsSystem
	})
	layerMap := compressLayersJSON(bonesByIndex, boneMap)

	jsonData := pmxJSON{
		Name:         modelData.Name(),
		Bones:        buildBonesJSON(bonesByIndex, boneMap, layerMap),
		DisplaySlots: buildDisplaySlotsJSON(modelData.DisplaySlots.Values(), boneMap),
		RigidBodies:  buildRigidBodiesJSON(modelData.RigidBodies.Values(), boneMap, opts.IncludeSystem),
	}

	payload, err := json.MarshalIndent(jsonData, "", "  ")
	if err != nil {
		return io_common.NewIoEncodeFailed("JSONのエンコードに失敗しました", err)
	}
	if err := os.WriteFile(savePath, payload, 0o644); err != nil {
		return io_common.NewIoSaveFailed("JSONファイルの保存に失敗しました", err)
	}
	modelData.SetPath(savePath)
	modelData.UpdateHash()
	return nil
}

// loadBonesFromJSON はボーン一覧を読み込む。
func loadBonesFromJSON(modelData *model.PmxModel, bones []boneJSON) error {
	sorted := append([]boneJSON(nil), bones...)
	sort.Slice(sorted, func(i, j int) bool { return sorted[i].Index < sorted[j].Index })
	for _, item := range sorted {
		bone := &model.Bone{
			EnglishName:  item.EnglishName,
			Position:     item.Position,
			ParentIndex:  item.ParentIndex,
			Layer:        item.Layer,
			BoneFlag:     model.BoneFlag(item.BoneFlag),
			TailPosition: item.TailPosition,
			TailIndex:    item.TailIndex,
			EffectIndex:  item.EffectIndex,
			EffectFactor: item.EffectFactor,
			FixedAxis:    item.FixedAxis,
			LocalAxisX:   item.LocalAxisX,
			LocalAxisZ:   item.LocalAxisZ,
			EffectorKey:  item.EffectorKey,
		}
		bone.SetName(item.Name)
		if item.Ik != nil {
			ik := &model.Ik{BoneIndex: item.Ik.BoneIndex, LoopCount: item.Ik.LoopCount}
			ik.UnitRotation = mmath.Vec3{Vec: r3.Vec{X: item.Ik.UnitRotation, Y: item.Ik.UnitRotation, Z: item.Ik.UnitRotation}}
			ik.Links = make([]model.IkLink, 0, len(item.Ik.Links))
			for _, link := range item.Ik.Links {
				ik.Links = append(ik.Links, model.IkLink{
					BoneIndex:     link.BoneIndex,
					AngleLimit:    link.AngleLimit,
					MinAngleLimit: link.MinAngleLimit,
					MaxAngleLimit: link.MaxAngleLimit,
				})
			}
			bone.Ik = ik
		}
		modelData.Bones.AppendRaw(bone)
	}
	return nil
}

// loadDisplaySlotsFromJSON は表示枠一覧を読み込む。
func loadDisplaySlotsFromJSON(modelData *model.PmxModel, slots []displaySlotJSON) error {
	sorted := append([]displaySlotJSON(nil), slots...)
	sort.Slice(sorted, func(i, j int) bool { return sorted[i].Index < sorted[j].Index })
	for _, item := range sorted {
		slot := &model.DisplaySlot{
			EnglishName: item.EnglishName,
			SpecialFlag: model.SpecialFlag(item.SpecialFlag),
			References:  make([]model.Reference, 0, len(item.References)),
		}
		slot.SetName(item.Name)
		for _, ref := range item.References {
			slot.References = append(slot.References, model.Reference{
				DisplayType:  model.DisplayType(ref.DisplayType),
				DisplayIndex: ref.DisplayIndex,
			})
		}
		modelData.DisplaySlots.AppendRaw(slot)
	}
	return nil
}

// loadRigidBodiesFromJSON は剛体一覧を読み込む。
func loadRigidBodiesFromJSON(modelData *model.PmxModel, bodies []rigidBodyJSON) error {
	sorted := append([]rigidBodyJSON(nil), bodies...)
	sort.Slice(sorted, func(i, j int) bool { return sorted[i].Index < sorted[j].Index })
	for _, item := range sorted {
		rigid := &model.RigidBody{
			EnglishName: item.EnglishName,
			BoneIndex:   item.BoneIndex,
			CollisionGroup: model.CollisionGroup{
				Group: byte(item.CollisionGroup),
				Mask:  uint16(item.CollisionGroupMask),
			},
			Shape:       model.Shape(item.ShapeType),
			Size:        item.Size,
			Position:    item.Position,
			Rotation:    item.Rotation,
			Param:       model.RigidBodyParam{Mass: item.Mass, LinearDamping: item.LinearDamping, AngularDamping: item.AngularDamping, Restitution: item.Restitution, Friction: 0},
			PhysicsType: model.PhysicsType(item.PhysicsType),
		}
		rigid.SetName(item.Name)
		modelData.RigidBodies.AppendRaw(rigid)
	}
	return nil
}

// buildIndexMappingJSON はJSON用のインデックス対応を生成する。
func buildIndexMappingJSON(total int, include func(index int) bool) jsonIndexMapping {
	oldToNew := make([]int, total)
	newToOld := make([]int, 0, total)
	for i := 0; i < total; i++ {
		oldToNew[i] = -1
	}
	for i := 0; i < total; i++ {
		if include == nil || include(i) {
			oldToNew[i] = len(newToOld)
			newToOld = append(newToOld, i)
		}
	}
	return jsonIndexMapping{oldToNew: oldToNew, newToOld: newToOld}
}

// jsonIndexMapping はJSON用の対応表を表す。
type jsonIndexMapping struct {
	oldToNew []int
	newToOld []int
}

// mapIndex は旧インデックスを新インデックスへ変換する。
func (m jsonIndexMapping) mapIndex(index int) int {
	if index < 0 || index >= len(m.oldToNew) {
		return -1
	}
	return m.oldToNew[index]
}

// compressLayersJSON はレイヤー値を圧縮する。
func compressLayersJSON(bones []*model.Bone, mapping jsonIndexMapping) map[int]int {
	layerSet := make(map[int]struct{})
	for _, oldIndex := range mapping.newToOld {
		if oldIndex < 0 || oldIndex >= len(bones) {
			continue
		}
		bone := bones[oldIndex]
		if bone == nil {
			continue
		}
		layerSet[bone.Layer] = struct{}{}
	}
	layers := make([]int, 0, len(layerSet))
	for layer := range layerSet {
		layers = append(layers, layer)
	}
	sort.Ints(layers)
	mapped := make(map[int]int, len(layers))
	for i, layer := range layers {
		mapped[layer] = i
	}
	return mapped
}

// buildBonesJSON はボーンJSONを生成する。
func buildBonesJSON(bones []*model.Bone, mapping jsonIndexMapping, layerMap map[int]int) []boneJSON {
	out := make([]boneJSON, 0, len(mapping.newToOld))
	for _, oldIndex := range mapping.newToOld {
		bone := bones[oldIndex]
		if bone == nil {
			continue
		}
		boneFlag := bone.BoneFlag
		tailIsBone := boneFlag&model.BONE_FLAG_TAIL_IS_BONE != 0
		tailIndex := bone.TailIndex
		if tailIsBone {
			mappedTail := mapping.mapIndex(bone.TailIndex)
			if mappedTail < 0 {
				boneFlag &^= model.BONE_FLAG_TAIL_IS_BONE
				tailIsBone = false
				tailIndex = -1
			} else {
				tailIndex = mappedTail
			}
		}
		ik := prepareIkJSON(bone.Ik, mapping)
		if boneFlag&model.BONE_FLAG_IS_IK != 0 && ik == nil {
			boneFlag &^= model.BONE_FLAG_IS_IK
		}
		layer := bone.Layer
		if mapped, ok := layerMap[layer]; ok {
			layer = mapped
		}
		item := boneJSON{
			Index:        mapping.mapIndex(bone.Index()),
			Name:         bone.Name(),
			EnglishName:  bone.EnglishName,
			Position:     bone.Position,
			ParentIndex:  mapping.mapIndex(bone.ParentIndex),
			Layer:        layer,
			BoneFlag:     int(boneFlag),
			TailPosition: bone.TailPosition,
			TailIndex:    tailIndex,
			EffectIndex:  mapping.mapIndex(bone.EffectIndex),
			EffectFactor: bone.EffectFactor,
			FixedAxis:    bone.FixedAxis,
			LocalAxisX:   bone.LocalAxisX,
			LocalAxisZ:   bone.LocalAxisZ,
			EffectorKey:  bone.EffectorKey,
			Ik:           ik,
		}
		out = append(out, item)
	}
	return out
}

// prepareIkJSON はIK情報を再割当して返す。
func prepareIkJSON(ik *model.Ik, mapping jsonIndexMapping) *ikJSON {
	if ik == nil {
		return nil
	}
	target := mapping.mapIndex(ik.BoneIndex)
	if target < 0 {
		return nil
	}
	links := make([]ikLinkJSON, 0, len(ik.Links))
	for _, link := range ik.Links {
		mapped := mapping.mapIndex(link.BoneIndex)
		if mapped < 0 {
			continue
		}
		links = append(links, ikLinkJSON{
			BoneIndex:     mapped,
			AngleLimit:    link.AngleLimit,
			MinAngleLimit: link.MinAngleLimit,
			MaxAngleLimit: link.MaxAngleLimit,
		})
	}
	if len(links) == 0 {
		return nil
	}
	return &ikJSON{
		BoneIndex:    target,
		LoopCount:    ik.LoopCount,
		UnitRotation: ik.UnitRotation.X,
		Links:        links,
	}
}

// buildDisplaySlotsJSON は表示枠JSONを生成する。
func buildDisplaySlotsJSON(slots []*model.DisplaySlot, boneMap jsonIndexMapping) []displaySlotJSON {
	sorted := append([]*model.DisplaySlot(nil), slots...)
	sort.Slice(sorted, func(i, j int) bool { return sorted[i].Index() < sorted[j].Index() })
	out := make([]displaySlotJSON, 0, len(sorted))
	for _, slot := range sorted {
		if slot == nil {
			continue
		}
		refs := make([]referenceJSON, 0, len(slot.References))
		for _, ref := range slot.References {
			if ref.DisplayType == model.DISPLAY_TYPE_BONE {
				mapped := boneMap.mapIndex(ref.DisplayIndex)
				if mapped < 0 {
					continue
				}
				refs = append(refs, referenceJSON{DisplayType: int(ref.DisplayType), DisplayIndex: mapped})
				continue
			}
			refs = append(refs, referenceJSON{DisplayType: int(ref.DisplayType), DisplayIndex: ref.DisplayIndex})
		}
		out = append(out, displaySlotJSON{
			Index:       slot.Index(),
			Name:        slot.Name(),
			EnglishName: slot.EnglishName,
			SpecialFlag: int(slot.SpecialFlag),
			References:  refs,
		})
	}
	return out
}

// buildRigidBodiesJSON は剛体JSONを生成する。
func buildRigidBodiesJSON(bodies []*model.RigidBody, boneMap jsonIndexMapping, includeSystem bool) []rigidBodyJSON {
	sorted := append([]*model.RigidBody(nil), bodies...)
	sort.Slice(sorted, func(i, j int) bool { return sorted[i].Index() < sorted[j].Index() })
	out := make([]rigidBodyJSON, 0, len(sorted))
	for _, rigid := range sorted {
		if rigid == nil {
			continue
		}
		if !includeSystem && rigid.IsSystem {
			continue
		}
		if !strings.HasPrefix(rigid.Name(), "mlib") {
			continue
		}
		out = append(out, rigidBodyJSON{
			Index:              rigid.Index(),
			Name:               rigid.Name(),
			EnglishName:        rigid.EnglishName,
			BoneIndex:          boneMap.mapIndex(rigid.BoneIndex),
			CollisionGroup:     0,
			CollisionGroupMask: 0,
			ShapeType:          int(model.SHAPE_SPHERE),
			Size:               rigid.Size,
			Position:           rigid.Position,
			Rotation:           rigid.Rotation,
			Mass:               0,
			LinearDamping:      0,
			AngularDamping:     0,
			Restitution:        0,
			Friction:           0,
			PhysicsType:        int(model.PHYSICS_TYPE_STATIC),
		})
	}
	return out
}
