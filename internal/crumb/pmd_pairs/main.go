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
)

const (
	startOpacityMilli = 900
	endOpacityMilli   = 1000
	opacityStepMilli  = 5
	flagPatternCount  = 32
)

// main は材質組み合わせPMXを生成して保存する。
func main() {
	outPath := parseArgs()

	m := buildPairsModel()
	if err := saveModel(outPath, m); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "PMX保存に失敗しました: %v\n", err)
		os.Exit(1)
	}

	_, _ = fmt.Fprintf(os.Stdout, "PMX保存完了: %s (材質数=%d)\n", outPath, m.Materials.Len())
}

// parseArgs はCLI引数を解析して保存先パスを返す。
func parseArgs() string {
	defaultOut := filepath.Join("internal", "crumb", "pmd_pairs", "pmd_pairs_materials.pmx")
	outPath := flag.String("output", defaultOut, "保存先PMXパス")
	flag.Parse()
	return *outPath
}

// buildPairsModel は指定条件の材質組み合わせを持つPMXモデルを生成する。
func buildPairsModel() *model.PmxModel {
	m := model.NewPmxModel()
	m.SetName("pmd_pairs_materials")
	m.EnglishName = "pmd_pairs_materials"
	m.Comment = "非透過度×drawFlag全組み合わせの材質検証用モデル"
	m.EnglishComment = "Material pair model for opacity and draw flag combinations"

	appendCenterBone(m)
	appendDisplaySlots(m)
	appendMaterials(m)

	return m
}

// appendCenterBone は要件のセンターボーンを追加する。
func appendCenterBone(m *model.PmxModel) {
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
	m.Bones.AppendRaw(center)
}

// appendDisplaySlots は Root と表示枠を追加し、表示枠にセンターを紐付ける。
func appendDisplaySlots(m *model.PmxModel) {
	root := model.NewRootDisplaySlot()
	m.DisplaySlots.AppendRaw(root)

	slot := &model.DisplaySlot{
		SpecialFlag: model.SPECIAL_FLAG_OFF,
		References: []model.Reference{
			{
				DisplayType:  model.DISPLAY_TYPE_BONE,
				DisplayIndex: 0,
			},
		},
	}
	slot.SetName("表示枠")
	slot.EnglishName = "Display"
	m.DisplaySlots.AppendRaw(slot)
}

// appendMaterials は非透過度と描画フラグ全組み合わせの材質を追加する。
func appendMaterials(m *model.PmxModel) {
	for opacityMilli := startOpacityMilli; opacityMilli <= endOpacityMilli; opacityMilli += opacityStepMilli {
		opacity := float64(opacityMilli) / 1000.0
		for pattern := 0; pattern < flagPatternCount; pattern++ {
			mat := model.NewMaterial()
			name := formatMaterialName(opacity, pattern)
			mat.SetName(name)
			mat.EnglishName = name
			mat.Diffuse = mmath.Vec4{X: 1.0, Y: 1.0, Z: 1.0, W: opacity}
			mat.Specular = mmath.ZERO_VEC4
			mat.Ambient = mmath.ZERO_VEC3
			mat.DrawFlag = toDrawFlag(pattern)
			mat.Edge = mmath.Vec4{X: 0.0, Y: 0.0, Z: 0.0, W: 1.0}
			mat.EdgeSize = 1.0
			mat.TextureFactor = mmath.ONE_VEC4
			mat.SphereTextureFactor = mmath.ONE_VEC4
			mat.ToonTextureFactor = mmath.ONE_VEC4
			mat.TextureIndex = -1
			mat.SphereTextureIndex = -1
			mat.SphereMode = model.SPHERE_MODE_INVALID
			mat.ToonSharingFlag = model.TOON_SHARING_INDIVIDUAL
			mat.ToonTextureIndex = -1
			mat.VerticesCount = 0
			m.Materials.AppendRaw(mat)
		}
	}
}

// formatMaterialName は仕様形式の材質名を返す。
func formatMaterialName(opacity float64, pattern int) string {
	doubleSided := (pattern >> 0) & 1
	groundShadow := (pattern >> 1) & 1
	onSelfShadowMap := (pattern >> 2) & 1
	drawSelfShadow := (pattern >> 3) & 1
	drawEdge := (pattern >> 4) & 1
	return fmt.Sprintf("%.3f_%d%d%d%d%d", opacity, doubleSided, groundShadow, onSelfShadowMap, drawSelfShadow, drawEdge)
}

// toDrawFlag はビット列を DrawFlag に変換する。
func toDrawFlag(pattern int) model.DrawFlag {
	var f model.DrawFlag
	if (pattern & (1 << 0)) != 0 {
		f |= model.DRAW_FLAG_DOUBLE_SIDED_DRAWING
	}
	if (pattern & (1 << 1)) != 0 {
		f |= model.DRAW_FLAG_GROUND_SHADOW
	}
	if (pattern & (1 << 2)) != 0 {
		f |= model.DRAW_FLAG_DRAWING_ON_SELF_SHADOW_MAPS
	}
	if (pattern & (1 << 3)) != 0 {
		f |= model.DRAW_FLAG_DRAWING_SELF_SHADOWS
	}
	if (pattern & (1 << 4)) != 0 {
		f |= model.DRAW_FLAG_DRAWING_EDGE
	}
	return f
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
