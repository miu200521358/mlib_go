// 指示: miu200521358
package pmd

import (
	"path/filepath"
	"strings"

	"github.com/miu200521358/mlib_go/pkg/domain/model"
)

// indexMapping は旧新インデックスの対応を表す。
type indexMapping struct {
	oldToNew []int
	newToOld []int
}

// buildIndexMapping は対象のインデックス対応を生成する。
func buildIndexMapping(total int, include func(index int) bool) indexMapping {
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
	return indexMapping{oldToNew: oldToNew, newToOld: newToOld}
}

// mapIndex は旧インデックスを新インデックスへ変換する。
func (m indexMapping) mapIndex(index int) int {
	if index < 0 || index >= len(m.oldToNew) {
		return -1
	}
	return m.oldToNew[index]
}

type textureSpec struct {
	textureName string
	sphereName  string
	sphereMode  model.SphereMode
}

// parseTextureSpec はPMD材質のテクスチャ指定を分解する。
func parseTextureSpec(raw string) textureSpec {
	name := strings.TrimSpace(raw)
	if name == "" {
		return textureSpec{}
	}
	if strings.Contains(name, "*") {
		parts := strings.SplitN(name, "*", 2)
		base := strings.TrimSpace(parts[0])
		sphere := strings.TrimSpace(parts[1])
		return textureSpec{
			textureName: base,
			sphereName:  sphere,
			sphereMode:  sphereModeFromName(sphere, model.SPHERE_MODE_MULTIPLICATION),
		}
	}
	if strings.Contains(name, "/") {
		parts := strings.SplitN(name, "/", 2)
		base := strings.TrimSpace(parts[0])
		sphere := strings.TrimSpace(parts[1])
		return textureSpec{
			textureName: base,
			sphereName:  sphere,
			sphereMode:  model.SPHERE_MODE_MULTIPLICATION,
		}
	}
	if isSphereTexture(name) {
		return textureSpec{
			sphereName: name,
			sphereMode: sphereModeFromName(name, model.SPHERE_MODE_MULTIPLICATION),
		}
	}
	return textureSpec{textureName: name}
}

// buildTextureSpec はPMD材質のテクスチャ指定を組み立てる。
func buildTextureSpec(textureName, sphereName string, sphereMode model.SphereMode) string {
	base := strings.TrimSpace(textureName)
	sphere := strings.TrimSpace(sphereName)
	if base == "" && sphere == "" {
		return ""
	}
	if base == "" {
		return sphere
	}
	if sphere == "" {
		return base
	}
	return base + "*" + sphere
}

// isSphereTexture はスフィアマップ拡張子か判定する。
func isSphereTexture(name string) bool {
	ext := strings.ToLower(filepath.Ext(name))
	return ext == ".sph" || ext == ".spa"
}

// sphereModeFromName はスフィア名から計算モードを決定する。
func sphereModeFromName(name string, fallback model.SphereMode) model.SphereMode {
	if strings.HasSuffix(strings.ToLower(name), ".spa") {
		return model.SPHERE_MODE_ADDITION
	}
	return fallback
}

// panelFromSkinType はPMDスキン種別からモーフパネルへ変換する。
func panelFromSkinType(skinType byte) model.MorphPanel {
	switch skinType {
	case 1:
		return model.MORPH_PANEL_EYEBROW_LOWER_LEFT
	case 2:
		return model.MORPH_PANEL_EYE_UPPER_LEFT
	case 3:
		return model.MORPH_PANEL_LIP_UPPER_RIGHT
	case 4:
		return model.MORPH_PANEL_OTHER_LOWER_RIGHT
	default:
		return model.MORPH_PANEL_OTHER_LOWER_RIGHT
	}
}

// skinTypeFromPanel はモーフパネルからPMDスキン種別へ変換する。
func skinTypeFromPanel(panel model.MorphPanel) byte {
	switch panel {
	case model.MORPH_PANEL_EYEBROW_LOWER_LEFT:
		return 1
	case model.MORPH_PANEL_EYE_UPPER_LEFT:
		return 2
	case model.MORPH_PANEL_LIP_UPPER_RIGHT:
		return 3
	case model.MORPH_PANEL_OTHER_LOWER_RIGHT:
		return 4
	default:
		return 4
	}
}

// boneFlagsFromType はPMDボーン種別からボーンフラグを組み立てる。
func boneFlagsFromType(boneType byte, hasTail bool) model.BoneFlag {
	flag := model.BONE_FLAG_NONE
	if hasTail {
		flag |= model.BONE_FLAG_TAIL_IS_BONE
	}
	switch boneType {
	case 1:
		flag |= model.BONE_FLAG_CAN_ROTATE | model.BONE_FLAG_CAN_TRANSLATE | model.BONE_FLAG_CAN_MANIPULATE | model.BONE_FLAG_IS_VISIBLE
	case 2:
		flag |= model.BONE_FLAG_CAN_ROTATE | model.BONE_FLAG_CAN_TRANSLATE | model.BONE_FLAG_CAN_MANIPULATE | model.BONE_FLAG_IS_VISIBLE | model.BONE_FLAG_IS_IK
	case 4:
		flag |= model.BONE_FLAG_CAN_ROTATE | model.BONE_FLAG_CAN_MANIPULATE | model.BONE_FLAG_IS_VISIBLE
	case 5:
		flag |= model.BONE_FLAG_CAN_ROTATE | model.BONE_FLAG_CAN_MANIPULATE | model.BONE_FLAG_IS_VISIBLE | model.BONE_FLAG_IS_EXTERNAL_ROTATION
	case 7:
		flag |= model.BONE_FLAG_CAN_ROTATE
	case 8:
		flag |= model.BONE_FLAG_CAN_ROTATE | model.BONE_FLAG_CAN_MANIPULATE | model.BONE_FLAG_IS_VISIBLE
	case 9:
		flag |= model.BONE_FLAG_CAN_ROTATE | model.BONE_FLAG_CAN_MANIPULATE | model.BONE_FLAG_IS_EXTERNAL_ROTATION
	default:
		flag |= model.BONE_FLAG_CAN_ROTATE | model.BONE_FLAG_CAN_MANIPULATE | model.BONE_FLAG_IS_VISIBLE
	}
	return flag
}

// boneTypeFromBone はボーン情報からPMDボーン種別を推定する。
func boneTypeFromBone(bone *model.Bone) byte {
	if bone == nil {
		return 0
	}
	if bone.BoneFlag&model.BONE_FLAG_IS_IK != 0 {
		return 2
	}
	if bone.BoneFlag&model.BONE_FLAG_IS_EXTERNAL_ROTATION != 0 {
		if bone.EffectFactor != 1.0 {
			return 9
		}
		return 5
	}
	if bone.BoneFlag&model.BONE_FLAG_IS_EXTERNAL_TRANSLATION != 0 {
		return 1
	}
	if bone.BoneFlag&model.BONE_FLAG_CAN_TRANSLATE != 0 {
		return 1
	}
	if bone.BoneFlag&model.BONE_FLAG_IS_VISIBLE == 0 {
		return 7
	}
	return 0
}

// defaultToonFileNames は既定のtoon画像ファイル名一覧を返す。
func defaultToonFileNames() []string {
	return []string{
		"toon01.bmp",
		"toon02.bmp",
		"toon03.bmp",
		"toon04.bmp",
		"toon05.bmp",
		"toon06.bmp",
		"toon07.bmp",
		"toon08.bmp",
		"toon09.bmp",
		"toon10.bmp",
	}
}
