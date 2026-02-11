// 指示: miu200521358
package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/miu200521358/mlib_go/pkg/adapter/io_common"
	"github.com/miu200521358/mlib_go/pkg/adapter/io_model/pmx"
	"github.com/miu200521358/mlib_go/pkg/domain/mmath"
	"github.com/miu200521358/mlib_go/pkg/domain/model"
	"gonum.org/v1/gonum/spatial/r3"
)

const (
	positionMin  = -50
	positionMax  = 50
	positionStep = 10
)

// main はダミーボーン組み合わせPMXを生成して保存する。
func main() {
	outPath := parseArgs()

	m := buildDummyBonesModel()
	if err := saveModel(outPath, m); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "PMX保存に失敗しました: %v\n", err)
		os.Exit(1)
	}

	_, _ = fmt.Fprintf(os.Stdout, "PMX保存完了: %s (ボーン数=%d)\n", outPath, m.Bones.Len())
}

// parseArgs はCLI引数を解析して保存先パスを返す。
func parseArgs() string {
	outPath := flag.String("output", "dummy_bones_by_x_column.pmx", "保存先PMXパス")
	flag.Parse()
	return *outPath
}

// buildDummyBonesModel はX列単位で操作可能なダミーボーンPMXモデルを生成する。
func buildDummyBonesModel() *model.PmxModel {
	m := model.NewPmxModel()
	m.SetName("dummy_bones_by_x_column")
	m.EnglishName = "dummy_bones_by_x_column"
	m.Comment = "X列単位の親ボーンとXYZ(-100〜100, step=10)ダミーボーン検証用モデル"
	m.EnglishComment = "Dummy bone model with X-column parent bones and XYZ combinations (range: -100 to 100, step: 10)"

	centerIndex := appendCenterBone(m)
	appendRootDisplaySlot(m)
	appendCenterDisplaySlot(m, centerIndex)
	appendXColumnBonesAndDisplaySlots(m, centerIndex)

	return m
}

// appendCenterBone は要件のセンターボーンを追加し、追加 index を返す。
func appendCenterBone(m *model.PmxModel) int {
	center := model.NewBoneByName("センター")
	center.EnglishName = "Center"
	center.Position = mmath.ZERO_VEC3
	center.ParentIndex = -1
	center.Layer = 0
	center.BoneFlag = model.BONE_FLAG_CAN_ROTATE |
		model.BONE_FLAG_CAN_TRANSLATE |
		model.BONE_FLAG_IS_VISIBLE |
		model.BONE_FLAG_CAN_MANIPULATE
	center.TailPosition = mmath.UNIT_Y_VEC3
	center.TailIndex = -1
	return m.Bones.AppendRaw(center)
}

// appendRootDisplaySlot は Root 表示枠を追加する。
func appendRootDisplaySlot(m *model.PmxModel) {
	root := model.NewRootDisplaySlot()
	m.DisplaySlots.AppendRaw(root)
}

// appendCenterDisplaySlot はセンターボーン表示枠を追加する。
func appendCenterDisplaySlot(m *model.PmxModel, centerIndex int) {
	slot := &model.DisplaySlot{
		SpecialFlag: model.SPECIAL_FLAG_OFF,
		References: []model.Reference{
			{
				DisplayType:  model.DISPLAY_TYPE_BONE,
				DisplayIndex: centerIndex,
			},
		},
	}
	slot.SetName("センター")
	slot.EnglishName = "Center"
	m.DisplaySlots.AppendRaw(slot)
}

// appendXColumnBonesAndDisplaySlots はX列ごとの親ボーンと表示枠を追加する。
func appendXColumnBonesAndDisplaySlots(m *model.PmxModel, centerIndex int) {
	columnBoneCount := axisValueCount() * axisValueCount()
	for x := positionMin; x <= positionMax; x += positionStep {
		parentBoneIndex := m.Bones.AppendRaw(newXColumnParentBone(centerIndex, x))
		columnBoneIndices := make([]int, 0, columnBoneCount+1)
		columnBoneIndices = append(columnBoneIndices, parentBoneIndex)

		for y := 0; y <= positionMax-positionMin; y += positionStep {
			for z := positionMin; z <= positionMax; z += positionStep {
				boneIndex := m.Bones.AppendRaw(newDummyBone(parentBoneIndex, x, y, z))
				columnBoneIndices = append(columnBoneIndices, boneIndex)
			}
		}

		m.DisplaySlots.AppendRaw(newXColumnDisplaySlot(x, columnBoneIndices))
	}
}

// axisValueCount はX/Y/Z軸で生成する値の件数を返す。
func axisValueCount() int {
	return ((positionMax - positionMin) / positionStep) + 1
}

// newXColumnParentBone は同一X列を一括操作するための親ボーンを生成する。
func newXColumnParentBone(centerIndex, x int) *model.Bone {
	name := formatXColumnParentBoneName(x)
	bone := model.NewBoneByName(name)
	bone.EnglishName = name
	bone.Position = mmath.Vec3{Vec: r3.Vec{X: float64(x), Y: 0.0, Z: 0.0}}
	bone.ParentIndex = centerIndex
	bone.Layer = 0
	bone.BoneFlag = model.BONE_FLAG_CAN_ROTATE |
		model.BONE_FLAG_CAN_TRANSLATE |
		model.BONE_FLAG_IS_VISIBLE |
		model.BONE_FLAG_CAN_MANIPULATE
	bone.TailPosition = mmath.ZERO_VEC3
	bone.TailIndex = -1
	return bone
}

// newDummyBone は指定した親indexと座標を持つダミーボーンを生成する。
func newDummyBone(parentIndex, x, y, z int) *model.Bone {
	name := formatDummyBoneName(x, y, z)
	bone := model.NewBoneByName(name)
	bone.EnglishName = name
	bone.Position = mmath.Vec3{Vec: r3.Vec{X: float64(x), Y: float64(y), Z: float64(z)}}
	bone.ParentIndex = parentIndex
	bone.Layer = 0
	bone.BoneFlag = model.BONE_FLAG_CAN_ROTATE |
		model.BONE_FLAG_CAN_TRANSLATE |
		model.BONE_FLAG_IS_VISIBLE |
		model.BONE_FLAG_CAN_MANIPULATE
	bone.TailPosition = mmath.ZERO_VEC3
	bone.TailIndex = -1
	return bone
}

// newXColumnDisplaySlot は同一X列のボーン一覧を持つ表示枠を生成する。
func newXColumnDisplaySlot(x int, boneIndices []int) *model.DisplaySlot {
	references := make([]model.Reference, 0, len(boneIndices))
	for _, boneIndex := range boneIndices {
		references = append(references, model.Reference{
			DisplayType:  model.DISPLAY_TYPE_BONE,
			DisplayIndex: boneIndex,
		})
	}

	slot := &model.DisplaySlot{
		SpecialFlag: model.SPECIAL_FLAG_OFF,
		References:  references,
	}
	slotName := formatXColumnDisplaySlotName(x)
	slot.SetName(slotName)
	slot.EnglishName = slotName
	return slot
}

// formatXColumnParentBoneName はX列親ボーン名を返す。
func formatXColumnParentBoneName(x int) string {
	return fmt.Sprintf("X%+02d", x)
}

// formatDummyBoneName はダミーボーン名を返す。
func formatDummyBoneName(x, y, z int) string {
	return fmt.Sprintf("X%+02dY%+02dZ%+02d", x, y, z)
}

// formatXColumnDisplaySlotName はX列表示枠名を返す。
func formatXColumnDisplaySlotName(x int) string {
	return fmt.Sprintf("X%+02d", x)
}

// saveModel はPMXファイルとしてモデルを保存する。
func saveModel(outPath string, m *model.PmxModel) error {
	if err := os.MkdirAll(filepath.Dir(outPath), 0o755); err != nil {
		return fmt.Errorf("保存先ディレクトリ作成に失敗: %w", err)
	}
	r := pmx.NewPmxRepository()
	if err := r.Save(outPath, m, io_common.SaveOptions{IncludeSystem: true}); err != nil {
		return fmt.Errorf("PMX保存に失敗: %w", err)
	}
	return nil
}
