// 指示: miu200521358
package deform

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"os"
	"strings"
	"testing"

	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/encoding/unicode"

	"github.com/miu200521358/mlib_go/pkg/domain/delta"
	"github.com/miu200521358/mlib_go/pkg/domain/mmath"
	"github.com/miu200521358/mlib_go/pkg/domain/model"
	motionpkg "github.com/miu200521358/mlib_go/pkg/domain/motion"
)

// standardBoneName はボーン名のテンプレートを表す。
type standardBoneName string

// String はボーン名を返す。
func (s standardBoneName) String() string {
	return string(s)
}

// Left は左方向のボーン名を返す。
func (s standardBoneName) Left() string {
	return replaceBoneDirection(string(s), "左")
}

// Right は右方向のボーン名を返す。
func (s standardBoneName) Right() string {
	return replaceBoneDirection(string(s), "右")
}

// StringFromDirection は方向指定のボーン名を返す。
func (s standardBoneName) StringFromDirection(direction string) string {
	return replaceBoneDirection(string(s), direction)
}

// replaceBoneDirection は {d} を方向文字列で置換する。
func replaceBoneDirection(name, direction string) string {
	return strings.ReplaceAll(name, "{d}", direction)
}

var pmx = struct {
	ROOT        standardBoneName
	CENTER      standardBoneName
	GROOVE      standardBoneName
	WAIST       standardBoneName
	LOWER       standardBoneName
	UPPER       standardBoneName
	UPPER2      standardBoneName
	NECK        standardBoneName
	HEAD        standardBoneName
	SHOULDER_P  standardBoneName
	SHOULDER    standardBoneName
	ARM         standardBoneName
	ARM_TWIST   standardBoneName
	ELBOW       standardBoneName
	WRIST_TWIST standardBoneName
	WRIST       standardBoneName
	INDEX1      standardBoneName
	INDEX2      standardBoneName
	INDEX3      standardBoneName
	LEG         standardBoneName
	LEG_D       standardBoneName
	KNEE        standardBoneName
	KNEE_D      standardBoneName
	ANKLE       standardBoneName
	ANKLE_D     standardBoneName
	HEEL        standardBoneName
	TOE_EX      standardBoneName
	LEG_IK      standardBoneName
	TOE_IK      standardBoneName
}{
	ROOT:        "全ての親",
	CENTER:      "センター",
	GROOVE:      "グルーブ",
	WAIST:       "腰",
	LOWER:       "下半身",
	UPPER:       "上半身",
	UPPER2:      "上半身2",
	NECK:        "首",
	HEAD:        "頭",
	SHOULDER_P:  "{d}肩P",
	SHOULDER:    "{d}肩",
	ARM:         "{d}腕",
	ARM_TWIST:   "{d}腕捩",
	ELBOW:       "{d}ひじ",
	WRIST_TWIST: "{d}手捩",
	WRIST:       "{d}手首",
	INDEX1:      "{d}人指１",
	INDEX2:      "{d}人指２",
	INDEX3:      "{d}人指３",
	LEG:         "{d}足",
	LEG_D:       "{d}足D",
	KNEE:        "{d}ひざ",
	KNEE_D:      "{d}ひざD",
	ANKLE:       "{d}足首",
	ANKLE_D:     "{d}足首D",
	HEEL:        "{d}かかと",
	TOE_EX:      "{d}足先EX",
	LEG_IK:      "{d}足ＩＫ",
	TOE_IK:      "{d}つま先ＩＫ",
}

func TestVmdMotion_Deform_Exists(t *testing.T) {
	motion := loadVmd(t, "../../../test_resources/サンプルモーション.vmd")

	model := loadPmx(t, "../../../test_resources/サンプルモデル.pmx")

	{

		boneDeltas, _ := computeBoneDeltas(model, motion, motionpkg.Frame(10), []string{pmx.INDEX3.Left()}, false, false, false)
		{
			expectedPosition := vec3(0.0, 0.0, 0.0)
			if !boneDeltas.GetByName(pmx.ROOT.String()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.ROOT.String()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.ROOT.String()).FilledGlobalPosition()))
			}
		}
		{
			expectedPosition := vec3(0.044920, 8.218059, 0.069347)
			if !boneDeltas.GetByName(pmx.CENTER.String()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.CENTER.String()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.CENTER.String()).FilledGlobalPosition()))
			}
		}
		{
			expectedPosition := vec3(0.044920, 9.392067, 0.064877)
			if !boneDeltas.GetByName(pmx.GROOVE.String()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.GROOVE.String()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.GROOVE.String()).FilledGlobalPosition()))
			}
		}
		{
			expectedPosition := vec3(0.044920, 11.740084, 0.055937)
			if !boneDeltas.GetByName(pmx.WAIST.String()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.WAIST.String()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.WAIST.String()).FilledGlobalPosition()))
			}
		}
		{
			expectedPosition := vec3(0.044920, 12.390969, -0.100531)
			if !boneDeltas.GetByName(pmx.UPPER.String()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.UPPER.String()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.UPPER.String()).FilledGlobalPosition()))
			}
		}
		{
			expectedPosition := vec3(0.044920, 13.803633, -0.138654)
			if !boneDeltas.GetByName(pmx.UPPER2.String()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.UPPER2.String()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.UPPER2.String()).FilledGlobalPosition()))
			}
		}
		{
			expectedPosition := vec3(0.044920, 15.149180, 0.044429)
			if !boneDeltas.GetByName("上半身3").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("上半身3").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("上半身3").FilledGlobalPosition()))
			}
		}
		{
			expectedPosition := vec3(0.324862, 16.470263, 0.419041)
			if !boneDeltas.GetByName(pmx.SHOULDER_P.Left()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.SHOULDER_P.Left()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.SHOULDER_P.Left()).FilledGlobalPosition()))
			}
		}
		{
			expectedPosition := vec3(0.324862, 16.470263, 0.419041)
			if !boneDeltas.GetByName(pmx.SHOULDER.Left()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.SHOULDER.Left()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.SHOULDER.Left()).FilledGlobalPosition()))
			}
		}
		{
			expectedPosition := vec3(1.369838, 16.312170, 0.676838)
			if !boneDeltas.GetByName(pmx.ARM.Left()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.ARM.Left()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.ARM.Left()).FilledGlobalPosition()))
			}
		}
		{
			expectedPosition := vec3(1.845001, 15.024807, 0.747681)
			if !boneDeltas.GetByName(pmx.ARM_TWIST.Left()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.ARM_TWIST.Left()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.ARM_TWIST.Left()).FilledGlobalPosition()))
			}
		}
		{
			expectedPosition := vec3(2.320162, 13.737446, 0.818525)
			if !boneDeltas.GetByName(pmx.ELBOW.Left()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.ELBOW.Left()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.ELBOW.Left()).FilledGlobalPosition()))
			}
		}
		{
			expectedPosition := vec3(2.526190, 12.502445, 0.336127)
			if !boneDeltas.GetByName(pmx.WRIST_TWIST.Left()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.WRIST_TWIST.Left()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.WRIST_TWIST.Left()).FilledGlobalPosition()))
			}
		}
		{
			expectedPosition := vec3(2.732219, 11.267447, -0.146273)
			if !boneDeltas.GetByName(pmx.WRIST.Left()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.WRIST.Left()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.WRIST.Left()).FilledGlobalPosition()))
			}
		}
		{
			expectedPosition := vec3(2.649188, 10.546797, -0.607412)
			if !boneDeltas.GetByName(pmx.INDEX1.Left()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.INDEX1.Left()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.INDEX1.Left()).FilledGlobalPosition()))
			}
		}
		{
			expectedPosition := vec3(2.408238, 10.209290, -0.576288)
			if !boneDeltas.GetByName(pmx.INDEX2.Left()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.INDEX2.Left()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.INDEX2.Left()).FilledGlobalPosition()))
			}
		}
		{
			expectedPosition := vec3(2.360455, 10.422402, -0.442668)
			if !boneDeltas.GetByName(pmx.INDEX3.Left()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.INDEX3.Left()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.INDEX3.Left()).FilledGlobalPosition()))
			}
		}
	}
}

func TestVmdMotion_Deform_Lerp(t *testing.T) {
	motion := loadVmd(t, "../../../test_resources/サンプルモーション.vmd")

	model := loadPmx(t, "../../../test_resources/サンプルモデル.pmx")

	{
		boneDeltas, _ := computeBoneDeltas(model, motion, motionpkg.Frame(999), []string{pmx.INDEX3.Left()}, true, false, false)
		{
			expectedPosition := vec3(0.0, 0.0, 0.0)
			if !boneDeltas.GetByName(pmx.ROOT.String()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.ROOT.String()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.ROOT.String()).FilledGlobalPosition()))
			}
		}
		{
			expectedPosition := vec3(-0.508560, 8.218059, 0.791827)
			if !boneDeltas.GetByName(pmx.CENTER.String()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.CENTER.String()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.CENTER.String()).FilledGlobalPosition()))
			}
		}
		{
			expectedPosition := vec3(-0.508560, 9.182008, 0.787357)
			if !boneDeltas.GetByName(pmx.GROOVE.String()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.GROOVE.String()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.GROOVE.String()).FilledGlobalPosition()))
			}
		}
		{
			expectedPosition := vec3(-0.508560, 11.530025, 0.778416)
			if !boneDeltas.GetByName(pmx.WAIST.String()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.WAIST.String()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.WAIST.String()).FilledGlobalPosition()))
			}
		}
		{
			expectedPosition := vec3(-0.508560, 12.180910, 0.621949)
			if !boneDeltas.GetByName(pmx.UPPER.String()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.UPPER.String()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.UPPER.String()).FilledGlobalPosition()))
			}
		}
		{
			expectedPosition := vec3(-0.437343, 13.588836, 0.523215)
			if !boneDeltas.GetByName(pmx.UPPER2.String()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.UPPER2.String()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.UPPER2.String()).FilledGlobalPosition()))
			}
		}
		{
			expectedPosition := vec3(-0.552491, 14.941880, 0.528703)
			if !boneDeltas.GetByName("上半身3").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("上半身3").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("上半身3").FilledGlobalPosition()))
			}
		}
		{
			expectedPosition := vec3(-0.590927, 16.312325, 0.819156)
			if !boneDeltas.GetByName(pmx.SHOULDER_P.Left()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.SHOULDER_P.Left()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.SHOULDER_P.Left()).FilledGlobalPosition()))
			}
		}
		{
			expectedPosition := vec3(-0.590927, 16.312325, 0.819156)
			if !boneDeltas.GetByName(pmx.SHOULDER.Left()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.SHOULDER.Left()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.SHOULDER.Left()).FilledGlobalPosition()))
			}
		}
		{
			expectedPosition := vec3(0.072990, 16.156742, 1.666761)
			if !boneDeltas.GetByName(pmx.ARM.Left()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.ARM.Left()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.ARM.Left()).FilledGlobalPosition()))
			}
		}
		{
			expectedPosition := vec3(0.043336, 15.182318, 2.635117)
			if !boneDeltas.GetByName(pmx.ARM_TWIST.Left()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.ARM_TWIST.Left()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.ARM_TWIST.Left()).FilledGlobalPosition()))
			}
		}
		{
			expectedPosition := vec3(0.013682, 14.207894, 3.603473)
			if !boneDeltas.GetByName(pmx.ELBOW.Left()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.ELBOW.Left()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.ELBOW.Left()).FilledGlobalPosition()))
			}
		}
		{
			expectedPosition := vec3(1.222444, 13.711100, 3.299384)
			if !boneDeltas.GetByName(pmx.WRIST_TWIST.Left()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.WRIST_TWIST.Left()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.WRIST_TWIST.Left()).FilledGlobalPosition()))
			}
		}
		{
			expectedPosition := vec3(2.431205, 13.214306, 2.995294)
			if !boneDeltas.GetByName(pmx.WRIST.Left()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.WRIST.Left()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.WRIST.Left()).FilledGlobalPosition()))
			}
		}
		{
			expectedPosition := vec3(3.283628, 13.209089, 2.884702)
			if !boneDeltas.GetByName(pmx.INDEX1.Left()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.INDEX1.Left()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.INDEX1.Left()).FilledGlobalPosition()))
			}
		}
		{
			expectedPosition := vec3(3.665809, 13.070156, 2.797680)
			if !boneDeltas.GetByName(pmx.INDEX2.Left()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.INDEX2.Left()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.INDEX2.Left()).FilledGlobalPosition()))
			}
		}
		{
			expectedPosition := vec3(3.886795, 12.968100, 2.718276)
			if !boneDeltas.GetByName(pmx.INDEX3.Left()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.INDEX3.Left()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.INDEX3.Left()).FilledGlobalPosition()))
			}
		}
	}

}

func TestVmdMotion_DeformLegIk1_Matsu(t *testing.T) {
	motion := loadVmd(t, "../../../test_resources/サンプルモーション.vmd")

	model := loadPmx(t, "../../../test_resources/サンプルモデル.pmx")

	{
		boneDeltas, _ := computeBoneDeltas(model, motion, motionpkg.Frame(29), []string{"左つま先", pmx.HEEL.Left()}, true, false, false)
		{
			expectedPosition := vec3(-0.781335, 11.717622, 1.557067)
			if !boneDeltas.GetByName(pmx.LOWER.String()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.LOWER.String()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.LOWER.String()).FilledGlobalPosition()))
			}
		}
		{
			expectedPosition := vec3(-0.368843, 10.614175, 2.532657)
			if !boneDeltas.GetByName(pmx.LEG.Left()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.LEG.Left()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.LEG.Left()).FilledGlobalPosition()))
			}
		}
		{
			expectedPosition := vec3(0.983212, 6.945313, 0.487476)
			if !boneDeltas.GetByName(pmx.KNEE.Left()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.KNEE.Left()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.KNEE.Left()).FilledGlobalPosition()))
			}
		}
		{
			expectedPosition := vec3(-0.345842, 2.211842, 2.182894)
			if !boneDeltas.GetByName(pmx.ANKLE.Left()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.ANKLE.Left()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.ANKLE.Left()).FilledGlobalPosition()))
			}
		}
		{
			expectedPosition := vec3(-0.109262, -0.025810, 1.147780)
			if !boneDeltas.GetByName("左つま先").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左つま先").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("左つま先").FilledGlobalPosition()))
			}
		}
		{
			expectedPosition := vec3(-0.923587, 0.733788, 2.624565)
			if !boneDeltas.GetByName(pmx.HEEL.Left()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.HEEL.Left()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.HEEL.Left()).FilledGlobalPosition()))
			}
		}
	}
}

func TestVmdMotion_DeformLegIk2_Matsu(t *testing.T) {
	motion := loadVmd(t, "../../../test_resources/サンプルモーション.vmd")

	model := loadPmx(t, "../../../test_resources/サンプルモデル.pmx")

	{
		boneDeltas, _ := computeBoneDeltas(model, motion, motionpkg.Frame(3152), []string{"左つま先", pmx.HEEL.Left()}, true, false, false)
		{
			expectedPosition := vec3(7.928583, 11.713336, 1.998830)
			if !boneDeltas.GetByName(pmx.LOWER.String()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.LOWER.String()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.LOWER.String()).FilledGlobalPosition()))
			}
		}
		{
			expectedPosition := vec3(7.370017, 10.665785, 2.963280)
			if !boneDeltas.GetByName(pmx.LEG.Left()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.LEG.Left()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.LEG.Left()).FilledGlobalPosition()))
			}
		}
		{
			expectedPosition := vec3(9.282883, 6.689319, 2.96825)
			if !boneDeltas.GetByName(pmx.KNEE.Left()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.KNEE.Left()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.KNEE.Left()).FilledGlobalPosition()))
			}
		}
		{
			expectedPosition := vec3(4.115521, 7.276527, 2.980609)
			if !boneDeltas.GetByName(pmx.ANKLE.Left()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.ANKLE.Left()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.ANKLE.Left()).FilledGlobalPosition()))
			}
		}
		{
			expectedPosition := vec3(1.931355, 6.108739, 2.994883)
			if !boneDeltas.GetByName("左つま先").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左つま先").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("左つま先").FilledGlobalPosition()))
			}
		}
		{
			expectedPosition := vec3(2.569512, 7.844740, 3.002920)
			if !boneDeltas.GetByName(pmx.HEEL.Left()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.HEEL.Left()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.HEEL.Left()).FilledGlobalPosition()))
			}
		}
	}
}

func TestVmdMotion_DeformLegIk3_Matsu(t *testing.T) {
	// mlog.SetLevel(mlog.IK_VERBOSE)
	motion := loadVmd(t, "../../../test_resources/腰元.vmd")

	model := loadPmx(t, "../../../test_resources/サンプルモデル.pmx")

	{
		boneDeltas, _ := computeBoneDeltas(model, motion, motionpkg.Frame(60), nil, true, false, false)
		{
			expectedPosition := vec3(1.931959, 11.695199, -1.411883)
			if !boneDeltas.GetByName(pmx.LOWER.String()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.LOWER.String()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.LOWER.String()).FilledGlobalPosition()))
			}
		}
		{
			expectedPosition := vec3(2.927524, 10.550287, -1.218106)
			if !boneDeltas.GetByName(pmx.LEG.Left()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.LEG.Left()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.LEG.Left()).FilledGlobalPosition()))
			}
		}
		{
			expectedPosition := vec3(2.263363, 7.061642, -3.837192)
			if !boneDeltas.GetByName(pmx.KNEE.Left()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.KNEE.Left()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.KNEE.Left()).FilledGlobalPosition()))
			}
		}
		{
			expectedPosition := vec3(2.747242, 2.529942, -1.331971)
			if !boneDeltas.GetByName(pmx.ANKLE.Left()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.ANKLE.Left()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.ANKLE.Left()).FilledGlobalPosition()))
			}
		}
		{
			expectedPosition := vec3(2.263363, 7.061642, -3.837192)
			if !boneDeltas.GetByName(pmx.KNEE_D.Left()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.KNEE_D.Left()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.KNEE_D.Left()).FilledGlobalPosition()))
			}
		}
		{
			expectedPosition := vec3(1.916109, 1.177077, -1.452845)
			if !boneDeltas.GetByName(pmx.TOE_EX.Left()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.TOE_EX.Left()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.TOE_EX.Left()).FilledGlobalPosition()))
			}
		}
		{
			expectedPosition := vec3(1.809291, 0.242514, -1.182168)
			if !boneDeltas.GetByName("左つま先").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左つま先").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("左つま先").FilledGlobalPosition()))
			}
		}
		{
			expectedPosition := vec3(3.311764, 1.159233, -0.613653)
			if !boneDeltas.GetByName(pmx.HEEL.Left()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.HEEL.Left()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.HEEL.Left()).FilledGlobalPosition()))
			}
		}
	}
}

func TestVmdMotion_DeformLegIk4_Snow(t *testing.T) {
	// mlog.SetLevel(mlog.IK_VERBOSE)

	motion := loadVmd(t, "../../../test_resources/好き雪_2794.vmd")

	model := loadPmx(t, "../../../test_resources/サンプルモデル.pmx")

	boneDeltas, _ := computeBoneDeltas(model, motion, motionpkg.Frame(0), []string{"右つま先", pmx.HEEL.Right()}, true, false, false)
	{
		expectedPosition := vec3(1.316121, 11.687257, 2.263307)
		if !boneDeltas.GetByName(pmx.LOWER.String()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.LOWER.String()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.LOWER.String()).FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(0.175478, 10.780540, 2.728409)
		if !boneDeltas.GetByName(pmx.LEG.Right()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.LEG.Right()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.LEG.Right()).FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(0.950410, 11.256771, -1.589462)
		if !boneDeltas.GetByName(pmx.KNEE.Right()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.KNEE.Right()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.KNEE.Right()).FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(-1.025194, 7.871110, 1.828258)
		if !boneDeltas.GetByName(pmx.ANKLE.Right()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.ANKLE.Right()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.ANKLE.Right()).FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(-1.701147, 6.066556, 3.384271)
		if !boneDeltas.GetByName("右つま先").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("右つま先").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("右つま先").FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(-1.379169, 7.887148, 3.436968)
		if !boneDeltas.GetByName(pmx.HEEL.Right()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.HEEL.Right()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.HEEL.Right()).FilledGlobalPosition()))
		}
	}
}

func TestVmdMotion_DeformLegIk5_Koshi(t *testing.T) {
	motion := loadVmd(t, "../../../test_resources/腰元.vmd")

	model := loadPmx(t, "../../../test_resources/サンプルモデル.pmx")

	boneDeltas, _ := computeBoneDeltas(model, motion, motionpkg.Frame(7409), []string{"右つま先", pmx.HEEL.Right()}, true, false, false)
	{
		expectedPosition := vec3(-7.652257, 11.990970, -4.511993)
		if !boneDeltas.GetByName(pmx.LOWER.String()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.LOWER.String()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.LOWER.String()).FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(-8.637265, 10.835548, -4.326830)
		if !boneDeltas.GetByName(pmx.LEG.Right()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.LEG.Right()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.LEG.Right()).FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(-8.693436, 7.595280, -7.321638)
		if !boneDeltas.GetByName(pmx.KNEE.Right()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.KNEE.Right()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.KNEE.Right()).FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(-7.521027, 2.827226, -9.035607)
		if !boneDeltas.GetByName(pmx.ANKLE.Right()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.ANKLE.Right()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.ANKLE.Right()).FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(-7.453236, 0.356456, -8.876783)
		if !boneDeltas.GetByName("右つま先").FilledGlobalPosition().NearEquals(expectedPosition, 0.04) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("右つま先").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("右つま先").FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(-7.030497, 1.820072, -7.827912)
		if !boneDeltas.GetByName(pmx.HEEL.Right()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.HEEL.Right()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.HEEL.Right()).FilledGlobalPosition()))
		}
	}
}

func TestVmdMotion_DeformLegIk6_KoshiOff(t *testing.T) {
	motion := loadVmd(t, "../../../test_resources/腰元.vmd")

	model := loadPmx(t, "../../../test_resources/サンプルモデル.pmx")

	// IK OFF
	boneDeltas, _ := computeBoneDeltas(model, motion, motionpkg.Frame(0), nil, false, false, false)
	{
		expectedPosition := vec3(1.622245, 6.632885, 0.713205)
		if !boneDeltas.GetByName(pmx.KNEE.Left()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.KNEE.Left()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.KNEE.Left()).FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(1.003185, 1.474691, 0.475763)
		if !boneDeltas.GetByName(pmx.ANKLE.Left()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.ANKLE.Left()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.ANKLE.Left()).FilledGlobalPosition()))
		}
	}

}

func TestVmdMotion_DeformLegIk6_KoshiOn(t *testing.T) {
	motion := loadVmd(t, "../../../test_resources/腰元.vmd")

	model := loadPmx(t, "../../../test_resources/サンプルモデル.pmx")

	// IK ON
	boneDeltas, _ := computeBoneDeltas(model, motion, motionpkg.Frame(0), nil, true, false, false)
	{
		expectedPosition := vec3(2.143878, 6.558880, 1.121747)
		if !boneDeltas.GetByName(pmx.KNEE.Left()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.KNEE.Left()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.KNEE.Left()).FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(2.214143, 1.689811, 2.947619)
		if !boneDeltas.GetByName(pmx.ANKLE.Left()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.ANKLE.Left()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.ANKLE.Left()).FilledGlobalPosition()))
		}
	}

}

func TestVmdMotion_DeformLegIk6_KoshiIkOn(t *testing.T) {
	motion := loadVmd(t, "../../../test_resources/腰元.vmd")

	model := loadPmx(t, "../../../test_resources/サンプルモデル.pmx")

	// IK ON
	fno := int(0)

	ikEnabledFrame := motionpkg.NewIkEnabledFrame(motionpkg.Frame(fno), pmx.LEG_IK.Left())
	ikEnabledFrame.Enabled = true

	ikFrame := motionpkg.NewIkFrame(motionpkg.Frame(fno))
	ikFrame.IkList = append(ikFrame.IkList, ikEnabledFrame)

	boneDeltas, _ := computeBoneDeltas(model, motion, motionpkg.Frame(0), nil, true, false, false)

	{
		expectedPosition := vec3(2.143878, 6.558880, 1.121747)
		if !boneDeltas.GetByName(pmx.KNEE.Left()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.KNEE.Left()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.KNEE.Left()).FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(2.214143, 1.689811, 2.947619)
		if !boneDeltas.GetByName(pmx.ANKLE.Left()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.ANKLE.Left()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.ANKLE.Left()).FilledGlobalPosition()))
		}
	}

}

func TestVmdMotion_DeformLegIk6_KoshiIkOff(t *testing.T) {
	motion := loadVmd(t, "../../../test_resources/腰元.vmd")

	model := loadPmx(t, "../../../test_resources/サンプルモデル.pmx")

	// IK OFF

	fno := int(0)

	ikEnabledFrame := motionpkg.NewIkEnabledFrame(motionpkg.Frame(fno), pmx.LEG_IK.Left())
	ikEnabledFrame.Enabled = false

	ikFrame := motionpkg.NewIkFrame(motionpkg.Frame(fno))
	ikFrame.IkList = append(ikFrame.IkList, ikEnabledFrame)

	boneDeltas, _ := computeBoneDeltas(model, motion, motionpkg.Frame(0), nil, false, false, false)
	{
		expectedPosition := vec3(1.622245, 6.632885, 0.713205)
		if !boneDeltas.GetByName(pmx.KNEE.Left()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.KNEE.Left()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.KNEE.Left()).FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(1.003185, 1.474691, 0.475763)
		if !boneDeltas.GetByName(pmx.ANKLE.Left()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.ANKLE.Left()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.ANKLE.Left()).FilledGlobalPosition()))
		}
	}
}

func TestVmdMotion_DeformLegIk7_Syou_ISAO(t *testing.T) {
	// mlog.SetLevel(mlog.IK_VERBOSE)

	motion := loadVmd(t, "C:/MMD/mmd_base/tests/resources/唱(ダンスのみ)_0274F.vmd")

	model := loadPmx(t, "D:/MMD/MikuMikuDance_v926x64/UserFile/Model/VOCALOID/初音ミク/ISAO式ミク/I_ミクv4/Miku_V4_準標準.pmx")

	boneDeltas, _ := computeBoneDeltas(model, motion, motionpkg.Frame(0), nil, true, false, false)

	{
		expectedPosition := vec3(0.04952335, 9.0, 1.72378033)
		if !boneDeltas.GetByName(pmx.CENTER.String()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.CENTER.String()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.CENTER.String()).FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(0.04952335, 7.97980869, 1.72378033)
		if !boneDeltas.GetByName(pmx.GROOVE.String()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.GROOVE.String()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.GROOVE.String()).FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(0.04952335, 11.02838314, 2.29172656)
		if !boneDeltas.GetByName("腰").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("腰").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("腰").FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(0.04952335, 11.9671191, 1.06765032)
		if !boneDeltas.GetByName("下半身").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("下半身").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("下半身").FilledGlobalPosition()))
		}
	}
	// FIXME: 物理後なので求められない
	// {
	// 	expectedPosition := vec3(-0.24102019, 9.79926074, 1.08498769)
	// 	if !boneDeltas.GetByName("下半身先").GetGlobalPosition().NearEquals(expectedPosition, 0.03) {
	// 		t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("下半身先").GetGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("下半身先").GetGlobalPosition()))
	// 	}
	// }
	{
		expectedPosition := vec3(0.90331914, 10.27362702, 1.009499759)
		if !boneDeltas.GetByName("腰キャンセル左").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("腰キャンセル左").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("腰キャンセル左").FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(0.90331914, 10.27362702, 1.00949975)
		if !boneDeltas.GetByName("左足").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左足").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("左足").FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(0.08276818, 5.59348757, -1.24981795)
		if !boneDeltas.GetByName("左ひざ").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左ひざ").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("左ひざ").FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(5.63290634e-01, -2.12439821e-04, -3.87768478e-01)
		if !boneDeltas.GetByName("左つま先").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左つま先").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("左つま先").FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(0.90331914, 10.27362702, 1.00949975)
		if !boneDeltas.GetByName("左足D").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左足D").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("左足D").FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(0.23453057, 5.6736954, -0.76228439)
		if !boneDeltas.GetByName("左ひざ2").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左ひざ2").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("左ひざ2").FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(0.12060311, 4.95396153, -1.23761938)
		if !boneDeltas.GetByName("左ひざ2先").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左ひざ2先").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("左ひざ2先").FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(0.90331914, 10.27362702, 1.00949975)
		if !boneDeltas.GetByName("左足y+").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左足y+").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("左足y+").FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(0.74736036, 9.38409308, 0.58008117)
		if !boneDeltas.GetByName("左足yTgt").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左足yTgt").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("左足yTgt").FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(0.74736036, 9.38409308, 0.58008117)
		if !boneDeltas.GetByName("左足yIK").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左足yIK").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("左足yIK").FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(-0.03018836, 10.40081089, 1.26859617)
		if !boneDeltas.GetByName("左尻").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左尻").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("左尻").FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(0.08276818, 5.59348757, -1.24981795)
		if !boneDeltas.GetByName("左ひざsub").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左ひざsub").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("左ひざsub").FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(-0.09359026, 5.54494997, -1.80895985)
		if !boneDeltas.GetByName("左ひざsub先").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左ひざsub先").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("左ひざsub先").FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(1.23779916, 1.28891465, 1.65257835)
		if !boneDeltas.GetByName("左ひざD2").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左ひざD2").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("左ひざD2").FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(0.1106881, 4.98643066, -1.26321915)
		if !boneDeltas.GetByName("左ひざD2先").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左ひざD2先").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("左ひざD2先").FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(0.12060311, 4.95396153, -1.23761938)
		if !boneDeltas.GetByName("左ひざD2IK").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左ひざD2IK").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("左ひざD2IK").FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(0.88590917, 0.38407067, 0.56801614)
		if !boneDeltas.GetByName("左足ゆび").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左足ゆび").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("左足ゆび").FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(5.63290634e-01, -2.12439821e-04, -3.87768478e-01)
		if !boneDeltas.GetByName("左つま先D").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左つま先D").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("左つま先D").FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(0.90331914, 10.27362702, 1.00949975)
		if !boneDeltas.GetByName("左足D").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左足D").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("左足D").FilledGlobalPosition()))
		}
	}
	// {
	// 	expectedPosition := vec3(0.08276818, 5.59348757, -1.24981795)
	// 	if !boneDeltas.GetByName("左ひざD").GetGlobalPosition().NearEquals(expectedPosition, 0.03) {
	// 		t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左ひざD").GetGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("左ひざD").GetGlobalPosition()))
	// 	}
	// }
	// {
	// 	expectedPosition := vec3(1.23779916, 1.28891465, 1.65257835)
	// 	if !boneDeltas.GetByName("左足首D").GetGlobalPosition().NearEquals(expectedPosition, 0.03) {
	// 		t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左足首D").GetGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("左足首D").GetGlobalPosition()))
	// 	}
	// }
}

func TestVmdMotion_DeformLegIk7_Syou(t *testing.T) {
	// mlog.SetLevel(mlog.IK_VERBOSE)

	motion := loadVmd(t, "../../../test_resources/唱(ダンスのみ)_0278F.vmd")

	model := loadPmx(t, "../../../test_resources/サンプルモデル.pmx")

	// 残存回転判定用
	boneDeltas, _ := computeBoneDeltas(model, motion, motionpkg.Frame(0), []string{"右つま先", pmx.HEEL.Right()}, true, false, false)
	{
		expectedPosition := vec3(0.721499, 11.767294, 1.638818)
		if !boneDeltas.GetByName(pmx.LOWER.String()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.LOWER.String()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.LOWER.String()).FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(-0.133304, 10.693992, 2.314730)
		if !boneDeltas.GetByName(pmx.LEG.Right()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.LEG.Right()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.LEG.Right()).FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(-2.833401, 8.174604, -0.100545)
		if !boneDeltas.GetByName(pmx.KNEE.Right()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.KNEE.Right()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.KNEE.Right()).FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(-0.409387, 5.341005, 3.524572)
		if !boneDeltas.GetByName(pmx.ANKLE.Right()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.ANKLE.Right()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.ANKLE.Right()).FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(-0.578271, 2.874233, 3.669599)
		if !boneDeltas.GetByName("右つま先").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("右つま先").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("右つま先").FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(0.322606, 4.249237, 4.517416)
		if !boneDeltas.GetByName(pmx.HEEL.Right()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.HEEL.Right()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.HEEL.Right()).FilledGlobalPosition()))
		}
	}
}

func TestVmdMotion_DeformLegIk8_Syou(t *testing.T) {
	motion := loadVmd(t, "../../../test_resources/唱(ダンスのみ)_0-300F.vmd")

	model := loadPmx(t, "../../../test_resources/サンプルモデル.pmx")

	boneDeltas, _ := computeBoneDeltas(model, motion, motionpkg.Frame(278), []string{"右つま先", pmx.HEEL.Right()}, true, false, false)
	{
		expectedPosition := vec3(0.721499, 11.767294, 1.638818)
		if !boneDeltas.GetByName(pmx.LOWER.String()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.LOWER.String()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.LOWER.String()).FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(-0.133304, 10.693992, 2.314730)
		if !boneDeltas.GetByName(pmx.LEG.Right()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.LEG.Right()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.LEG.Right()).FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(-2.833401, 8.174604, -0.100545)
		if !boneDeltas.GetByName(pmx.KNEE.Right()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.KNEE.Right()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.KNEE.Right()).FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(-0.409387, 5.341005, 3.524572)
		if !boneDeltas.GetByName(pmx.ANKLE.Right()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.ANKLE.Right()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.ANKLE.Right()).FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(-0.578271, 2.874233, 3.669599)
		if !boneDeltas.GetByName("右つま先").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("右つま先").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("右つま先").FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(0.322606, 4.249237, 4.517416)
		if !boneDeltas.GetByName(pmx.HEEL.Right()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.HEEL.Right()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.HEEL.Right()).FilledGlobalPosition()))
		}
	}
}

func TestVmdMotion_DeformLegIk10_Syou1(t *testing.T) {
	motion := loadVmd(t, "../../../test_resources/唱(ダンスのみ)_0-300F.vmd")

	model := loadPmx(t, "../../../test_resources/サンプルモデル.pmx")

	boneDeltas, _ := computeBoneDeltas(model, motion, motionpkg.Frame(100), []string{"右つま先", pmx.HEEL.Right()}, true, false, false)
	{
		expectedPosition := vec3(0.365000, 11.411437, 1.963828)
		if !boneDeltas.GetByName(pmx.LOWER.String()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.LOWER.String()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.LOWER.String()).FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(-0.513678, 10.280550, 2.500991)
		if !boneDeltas.GetByName(pmx.LEG.Right()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.LEG.Right()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.LEG.Right()).FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(-2.891708, 8.162312, -0.553409)
		if !boneDeltas.GetByName(pmx.KNEE.Right()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.KNEE.Right()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.KNEE.Right()).FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(-0.826174, 4.330670, 2.292396)
		if !boneDeltas.GetByName(pmx.ANKLE.Right()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.ANKLE.Right()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.ANKLE.Right()).FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(-1.063101, 1.865613, 2.335564)
		if !boneDeltas.GetByName("右つま先").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("右つま先").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("右つま先").FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(-0.178356, 3.184965, 3.282950)
		if !boneDeltas.GetByName(pmx.HEEL.Right()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.HEEL.Right()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.HEEL.Right()).FilledGlobalPosition()))
		}
	}
}

func TestVmdMotion_DeformLegIk10_Syou2(t *testing.T) {
	// mlog.SetLevel(mlog.IK_VERBOSE)
	motion := loadVmd(t, "../../../test_resources/唱(ダンスのみ)_0-300F.vmd")

	model := loadPmx(t, "../../../test_resources/サンプルモデル.pmx")

	boneDeltas, _ := computeBoneDeltas(model, motion, motionpkg.Frame(107), []string{"右つま先", pmx.HEEL.Right()}, true, false, false)
	{
		expectedPosition := vec3(0.365000, 12.042871, 2.034023)
		if !boneDeltas.GetByName(pmx.LOWER.String()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.LOWER.String()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.LOWER.String()).FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(-0.488466, 10.920292, 2.626419)
		if !boneDeltas.GetByName(pmx.LEG.Right()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.LEG.Right()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.LEG.Right()).FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(-1.607765, 6.763937, 1.653586)
		if !boneDeltas.GetByName(pmx.KNEE.Right()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.KNEE.Right()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.KNEE.Right()).FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(-1.110289, 1.718307, 2.809817)
		if !boneDeltas.GetByName(pmx.ANKLE.Right()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.ANKLE.Right()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.ANKLE.Right()).FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(-1.753089, -0.026766, 1.173958)
		if !boneDeltas.GetByName("右つま先").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("右つま先").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("右つま先").FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(-0.952785, 0.078826, 2.838099)
		if !boneDeltas.GetByName(pmx.HEEL.Right()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.HEEL.Right()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.HEEL.Right()).FilledGlobalPosition()))
		}
	}

}

func TestVmdMotion_DeformLegIk10_Syou3(t *testing.T) {
	// mlog.SetLevel(mlog.IK_VERBOSE)
	motion := loadVmd(t, "../../../test_resources/唱(ダンスのみ)_0-300F.vmd")

	model := loadPmx(t, "../../../test_resources/サンプルモデル.pmx")

	boneDeltas, _ := computeBoneDeltas(model, motion, motionpkg.Frame(272), []string{"右つま先", pmx.HEEL.Right()}, true, false, false)
	{
		expectedPosition := vec3(-0.330117, 10.811301, 1.914508)
		if !boneDeltas.GetByName(pmx.LOWER.String()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.LOWER.String()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.LOWER.String()).FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(-1.325985, 9.797281, 2.479780)
		if !boneDeltas.GetByName(pmx.LEG.Right()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.LEG.Right()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.LEG.Right()).FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(-1.394679, 6.299243, -0.209150)
		if !boneDeltas.GetByName(pmx.KNEE.Right()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.KNEE.Right()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.KNEE.Right()).FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(-0.865021, 1.642431, 2.044760)
		if !boneDeltas.GetByName(pmx.ANKLE.Right()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.ANKLE.Right()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.ANKLE.Right()).FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(-1.191817, -0.000789, 0.220605)
		if !boneDeltas.GetByName("右つま先").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("右つま先").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("右つま先").FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(-0.958608, -0.002146, 2.055439)
		if !boneDeltas.GetByName(pmx.HEEL.Right()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.HEEL.Right()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.HEEL.Right()).FilledGlobalPosition()))
		}
	}

}

func TestVmdMotion_DeformLegIk10_Syou4(t *testing.T) {
	motion := loadVmd(t, "../../../test_resources/唱(ダンスのみ)_0-300F.vmd")

	model := loadPmx(t, "../../../test_resources/サンプルモデル.pmx")

	boneDeltas, _ := computeBoneDeltas(model, motion, motionpkg.Frame(273), []string{"右つま先", pmx.HEEL.Right()}, true, false, false)
	{
		expectedPosition := vec3(-0.154848, 10.862784, 1.868560)
		if !boneDeltas.GetByName(pmx.LOWER.String()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.LOWER.String()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.LOWER.String()).FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(-1.153633, 9.846655, 2.436846)
		if !boneDeltas.GetByName(pmx.LEG.Right()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.LEG.Right()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.LEG.Right()).FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(-1.498977, 6.380789, -0.272370)
		if !boneDeltas.GetByName(pmx.KNEE.Right()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.KNEE.Right()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.KNEE.Right()).FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(-0.845777, 1.802650, 2.106815)
		if !boneDeltas.GetByName(pmx.ANKLE.Right()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.ANKLE.Right()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.ANKLE.Right()).FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(-1.239674, 0.026274, 0.426385)
		if !boneDeltas.GetByName("右つま先").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("右つま先").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("右つま先").FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(-0.797867, 0.159797, 2.217469)
		if !boneDeltas.GetByName(pmx.HEEL.Right()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.HEEL.Right()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.HEEL.Right()).FilledGlobalPosition()))
		}
	}

}

func TestVmdMotion_DeformLegIk10_Syou5(t *testing.T) {
	motion := loadVmd(t, "../../../test_resources/唱(ダンスのみ)_0-300F.vmd")

	model := loadPmx(t, "../../../test_resources/サンプルモデル.pmx")

	boneDeltas, _ := computeBoneDeltas(model, motion, motionpkg.Frame(274), []string{"右つま先", pmx.HEEL.Right()}, true, false, false)
	{
		expectedPosition := vec3(0.049523, 10.960778, 1.822612)
		if !boneDeltas.GetByName(pmx.LOWER.String()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.LOWER.String()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.LOWER.String()).FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(-0.930675, 9.938401, 2.400088)
		if !boneDeltas.GetByName(pmx.LEG.Right()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.LEG.Right()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.LEG.Right()).FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(-1.710987, 6.669293, -0.459177)
		if !boneDeltas.GetByName(pmx.KNEE.Right()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.KNEE.Right()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.KNEE.Right()).FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(-0.773748, 2.387820, 2.340310)
		if !boneDeltas.GetByName(pmx.ANKLE.Right()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.ANKLE.Right()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.ANKLE.Right()).FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(-1.256876, 0.365575, 0.994345)
		if !boneDeltas.GetByName("右つま先").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("右つま先").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("右つま先").FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(-0.556038, 0.785363, 2.653745)
		if !boneDeltas.GetByName(pmx.HEEL.Right()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.HEEL.Right()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.HEEL.Right()).FilledGlobalPosition()))
		}
	}

}

func TestVmdMotion_DeformLegIk10_Syou6(t *testing.T) {
	motion := loadVmd(t, "../../../test_resources/唱(ダンスのみ)_0-300F.vmd")

	model := loadPmx(t, "../../../test_resources/サンプルモデル.pmx")

	boneDeltas, _ := computeBoneDeltas(model, motion, motionpkg.Frame(278), []string{"右つま先", pmx.HEEL.Right()}, true, false, false)
	{
		expectedPosition := vec3(0.721499, 11.767294, 1.638818)
		if !boneDeltas.GetByName(pmx.LOWER.String()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.LOWER.String()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.LOWER.String()).FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(-0.133304, 10.693992, 2.314730)
		if !boneDeltas.GetByName(pmx.LEG.Right()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.LEG.Right()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.LEG.Right()).FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(-2.833401, 8.174604, -0.100545)
		if !boneDeltas.GetByName(pmx.KNEE.Right()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.KNEE.Right()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.KNEE.Right()).FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(-0.409387, 5.341005, 3.524572)
		if !boneDeltas.GetByName(pmx.ANKLE.Right()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.ANKLE.Right()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.ANKLE.Right()).FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(-0.578271, 2.874233, 3.669599)
		if !boneDeltas.GetByName("右つま先").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("右つま先").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("右つま先").FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(0.322606, 4.249237, 4.517416)
		if !boneDeltas.GetByName(pmx.HEEL.Right()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.HEEL.Right()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.HEEL.Right()).FilledGlobalPosition()))
		}
	}

}

func TestVmdMotion_DeformLegIk11_Shining_Miku(t *testing.T) {
	// mlog.SetLevel(mlog.IK_VERBOSE)
	motion := loadVmd(t, "../../../test_resources/シャイニングミラクル_50F.vmd")

	model := loadPmx(t, "D:/MMD/MikuMikuDance_v926x64/UserFile/Model/VOCALOID/初音ミク/Tda式初音ミク_盗賊つばき流Ｍトレースモデル配布 v1.07/Tda式初音ミク_盗賊つばき流Mトレースモデルv1.07_かかと.pmx")

	boneDeltas, _ := computeBoneDeltas(model, motion, motionpkg.Frame(0), []string{pmx.LEG_IK.Right(), pmx.HEEL.Right(), "足首_R_"}, true, false, false)
	{
		expectedPosition := vec3(-1.869911, 2.074591, -0.911531)
		if !boneDeltas.GetByName(pmx.LEG_IK.Right()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.LEG_IK.Right()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.LEG_IK.Right()).FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(0.0, 0.002071, 0.0)
		if !boneDeltas.GetByName(pmx.ROOT.String()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.ROOT.String()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.ROOT.String()).FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(0.0, 8.404771, -0.850001)
		if !boneDeltas.GetByName(pmx.CENTER.String()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.CENTER.String()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.CENTER.String()).FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(0.0, 5.593470, -0.850001)
		if !boneDeltas.GetByName(pmx.GROOVE.String()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.GROOVE.String()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.GROOVE.String()).FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(0.0, 9.311928, -0.586922)
		if !boneDeltas.GetByName(pmx.WAIST.String()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.WAIST.String()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.WAIST.String()).FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(0.0, 10.142656, -1.362172)
		if !boneDeltas.GetByName(pmx.LOWER.String()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.LOWER.String()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.LOWER.String()).FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(-0.843381, 8.895412, -0.666409)
		if !boneDeltas.GetByName(pmx.LEG.Right()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.LEG.Right()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.LEG.Right()).FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(-0.274925, 5.679991, -4.384042)
		if !boneDeltas.GetByName(pmx.KNEE.Right()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.KNEE.Right()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.KNEE.Right()).FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(-1.870632, 2.072767, -0.910016)
		if !boneDeltas.GetByName(pmx.ANKLE.Right()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.ANKLE.Right()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.ANKLE.Right()).FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(-1.485913, -0.300011, -1.310446)
		if !boneDeltas.GetByName("足首_R_").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("足首_R_").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("足首_R_").FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(-1.894769, 0.790468, 0.087442)
		if !boneDeltas.GetByName(pmx.HEEL.Right()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.HEEL.Right()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.HEEL.Right()).FilledGlobalPosition()))
		}
	}
}

func TestVmdMotion_DeformLegIk11_Shining_Vroid(t *testing.T) {
	// mlog.SetLevel(mlog.IK_VERBOSE)

	motion := loadVmd(t, "../../../test_resources/シャイニングミラクル_50F.vmd")

	model := loadPmx(t, "../../../test_resources/サンプルモデル.pmx")

	boneDeltas, _ := computeBoneDeltas(model, motion, motionpkg.Frame(0), []string{"右つま先", pmx.HEEL.Right()}, true, false, false)
	{
		expectedPosition := vec3(0.0, 9.379668, -1.051170)
		if !boneDeltas.GetByName(pmx.LOWER.String()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.LOWER.String()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.LOWER.String()).FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(-0.919751, 8.397145, -0.324375)
		if !boneDeltas.GetByName(pmx.LEG.Right()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.LEG.Right()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.LEG.Right()).FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(-0.422861, 6.169319, -4.100779)
		if !boneDeltas.GetByName(pmx.KNEE.Right()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.KNEE.Right()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.KNEE.Right()).FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(-1.821804, 2.095607, -1.186269)
		if !boneDeltas.GetByName(pmx.ANKLE.Right()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.ANKLE.Right()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.ANKLE.Right()).FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(-1.390510, -0.316872, -1.544655)
		if !boneDeltas.GetByName("右つま先").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("右つま先").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("右つま先").FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(-1.852786, 0.811991, -0.154341)
		if !boneDeltas.GetByName(pmx.HEEL.Right()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.HEEL.Right()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.HEEL.Right()).FilledGlobalPosition()))
		}
	}

}

func TestVmdMotion_DeformLegIk12_Down_Miku(t *testing.T) {
	// mlog.SetLevel(mlog.IK_VERBOSE)
	motion := loadVmd(t, "../../../test_resources/しゃがむ.vmd")

	model := loadPmx(t, "D:/MMD/MikuMikuDance_v926x64/UserFile/Model/VOCALOID/初音ミク/Tda式初音ミク_盗賊つばき流Ｍトレースモデル配布 v1.07/Tda式初音ミク_盗賊つばき流Mトレースモデルv1.07_かかと.pmx")

	boneDeltas, _ := computeBoneDeltas(model, motion, motionpkg.Frame(0), []string{pmx.LEG_IK.Right(), pmx.HEEL.Right(), "足首_R_"}, true, false, false)
	{
		expectedPosition := vec3(-1.012964, 1.623157, 0.680305)
		if !boneDeltas.GetByName(pmx.LEG_IK.Right()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.LEG_IK.Right()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.LEG_IK.Right()).FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(0.0, 5.953951, -0.512170)
		if !boneDeltas.GetByName(pmx.LOWER.String()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.LOWER.String()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.LOWER.String()).FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(-0.896440, 4.569404, -0.337760)
		if !boneDeltas.GetByName(pmx.LEG.Right()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.LEG.Right()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.LEG.Right()).FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(-0.691207, 1.986888, -4.553376)
		if !boneDeltas.GetByName(pmx.KNEE.Right()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.KNEE.Right()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.KNEE.Right()).FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(-1.012964, 1.623157, 0.680305)
		if !boneDeltas.GetByName(pmx.ANKLE.Right()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.ANKLE.Right()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.ANKLE.Right()).FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(-1.013000, 0.002578, -1.146909)
		if !boneDeltas.GetByName("足首_R_").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("足首_R_").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("足首_R_").FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(-1.056216, -0.001008, 0.676086)
		if !boneDeltas.GetByName(pmx.HEEL.Right()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.HEEL.Right()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.HEEL.Right()).FilledGlobalPosition()))
		}
	}
}

func TestVmdMotion_DeformLegIk13_Lamb(t *testing.T) {
	motion := loadVmd(t, "../../../test_resources/Lamb_2689F.vmd")

	model := loadPmx(t, "D:/MMD/MikuMikuDance_v926x64/UserFile/Model/ゲーム/戦国BASARA/幸村 たぬき式 ver.1.24/真田幸村没第二衣装1.24軽量版.pmx")
	boneDeltas, _ := computeBoneDeltas(model, motion, motionpkg.Frame(0), []string{pmx.LEG_IK.Right(), "右つま先", pmx.LEG_IK.Left(), "左つま先", pmx.HEEL.Left()}, true, false, false)

	{

		{
			expectedPosition := vec3(-1.216134, 1.887670, -10.78867)
			if !boneDeltas.GetByName(pmx.LEG_IK.Right()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.LEG_IK.Right()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.LEG_IK.Right()).FilledGlobalPosition()))
			}
		}
		{
			expectedPosition := vec3(0.803149, 6.056844, -10.232766)
			if !boneDeltas.GetByName(pmx.LOWER.String()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.LOWER.String()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.LOWER.String()).FilledGlobalPosition()))
			}
		}
		{
			expectedPosition := vec3(0.728442, 4.560226, -11.571869)
			if !boneDeltas.GetByName(pmx.LEG.Right()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.LEG.Right()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.LEG.Right()).FilledGlobalPosition()))
			}
		}
		{
			expectedPosition := vec3(4.173470, 0.361388, -11.217197)
			if !boneDeltas.GetByName(pmx.KNEE.Right()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.KNEE.Right()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.KNEE.Right()).FilledGlobalPosition()))
			}
		}
		{
			expectedPosition := vec3(-1.217569, 1.885731, -10.788104)
			if !boneDeltas.GetByName(pmx.ANKLE.Right()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.ANKLE.Right()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.ANKLE.Right()).FilledGlobalPosition()))
			}
		}
		{
			expectedPosition := vec3(-0.922247, -1.163554, -10.794323)
			if !boneDeltas.GetByName("右つま先").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("右つま先").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("右つま先").FilledGlobalPosition()))
			}
		}
	}
	{

		{
			expectedPosition := vec3(2.322227, 1.150214, -9.644499)
			if !boneDeltas.GetByName(pmx.LEG_IK.Left()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.LEG_IK.Left()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.LEG_IK.Left()).FilledGlobalPosition()))
			}
		}
		{
			expectedPosition := vec3(0.803149, 6.056844, -10.232766)
			if !boneDeltas.GetByName(pmx.LOWER.String()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.LOWER.String()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.LOWER.String()).FilledGlobalPosition()))
			}
		}
		{
			expectedPosition := vec3(0.720821, 4.639688, -8.810255)
			if !boneDeltas.GetByName(pmx.LEG.Left()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.LEG.Left()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.LEG.Left()).FilledGlobalPosition()))
			}
		}
		{
			expectedPosition := vec3(6.126388, 5.074682, -8.346903)
			if !boneDeltas.GetByName(pmx.KNEE.Left()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.KNEE.Left()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.KNEE.Left()).FilledGlobalPosition()))
			}
		}
		{
			expectedPosition := vec3(2.323599, 1.147291, -9.645196)
			if !boneDeltas.GetByName(pmx.ANKLE.Left()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.ANKLE.Left()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.ANKLE.Left()).FilledGlobalPosition()))
			}
		}
		{
			expectedPosition := vec3(5.163002, -0.000894, -9.714369)
			if !boneDeltas.GetByName("左つま先").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左つま先").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("左つま先").FilledGlobalPosition()))
			}
		}
	}
}

func TestVmdMotion_DeformLegIk14_Ballet(t *testing.T) {
	motion := loadVmd(t, "../../../test_resources/ミク用バレリーコ_1069.vmd")

	model := loadPmx(t, "D:/MMD/MikuMikuDance_v926x64/UserFile/Model/_あにまさ式/初音ミク_準標準.pmx")

	boneDeltas, _ := computeBoneDeltas(model, motion, motionpkg.Frame(0), []string{pmx.LEG_IK.Right(), "右つま先", pmx.HEEL.Right()}, true, false, false)
	{
		expectedPosition := vec3(11.324574, 10.920002, -7.150005)
		if !boneDeltas.GetByName(pmx.LEG_IK.Right()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.LEG_IK.Right()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.LEG_IK.Right()).FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(2.433170, 13.740387, 0.992719)
		if !boneDeltas.GetByName(pmx.LOWER.String()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.LOWER.String()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.LOWER.String()).FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(1.982654, 11.188538, 0.602013)
		if !boneDeltas.GetByName(pmx.LEG.Right()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.LEG.Right()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.LEG.Right()).FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(5.661557, 11.008962, -2.259013)
		if !boneDeltas.GetByName(pmx.KNEE.Right()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.KNEE.Right()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.KNEE.Right()).FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(9.224476, 10.979847, -5.407887)
		if !boneDeltas.GetByName(pmx.ANKLE.Right()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.ANKLE.Right()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.ANKLE.Right()).FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(11.345482, 10.263426, -7.003638)
		if !boneDeltas.GetByName("右つま先").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("右つま先").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("右つま先").FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(9.406674, 9.687277, -5.710646)
		if !boneDeltas.GetByName(pmx.HEEL.Right()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.HEEL.Right()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.HEEL.Right()).FilledGlobalPosition()))
		}
	}
}

func TestVmdMotion_DeformLegIk15_Bottom(t *testing.T) {
	// mlog.SetLevel(mlog.IK_VERBOSE)

	motion := loadVmd(t, "../../../test_resources/●ボトム_0-300.vmd")

	model := loadPmx(t, "D:/MMD/MikuMikuDance_v926x64/UserFile/Model/VOCALOID/初音ミク/Tda式初音ミク_盗賊つばき流Ｍトレースモデル配布 v1.07/Tda式初音ミク_盗賊つばき流Mトレースモデルv1.07_かかと.pmx")

	boneDeltas, _ := computeBoneDeltas(model, motion, motionpkg.Frame(218), []string{pmx.LEG_IK.Right(), pmx.HEEL.Right(), "足首_R_"}, true, false, false)
	{
		expectedPosition := vec3(-1.358434, 1.913062, 0.611182)
		if !boneDeltas.GetByName(pmx.LEG_IK.Right()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.LEG_IK.Right()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.LEG_IK.Right()).FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(0.150000, 4.253955, 0.237829)
		if !boneDeltas.GetByName(pmx.LOWER.String()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.LOWER.String()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.LOWER.String()).FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(-0.906292, 2.996784, 0.471846)
		if !boneDeltas.GetByName(pmx.LEG.Right()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.LEG.Right()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.LEG.Right()).FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(-2.533418, 3.889916, -4.114837)
		if !boneDeltas.GetByName(pmx.KNEE.Right()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.KNEE.Right()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.KNEE.Right()).FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(-1.358807, 1.912181, 0.611265)
		if !boneDeltas.GetByName(pmx.ANKLE.Right()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.ANKLE.Right()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.ANKLE.Right()).FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(-2.040872, -0.188916, -0.430442)
		if !boneDeltas.GetByName("足首_R_").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("足首_R_").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("足首_R_").FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(-1.292688, 0.375211, 1.133899)
		if !boneDeltas.GetByName(pmx.HEEL.Right()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.HEEL.Right()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.HEEL.Right()).FilledGlobalPosition()))
		}
	}
}

func TestVmdMotion_DeformLegIk16_Lamb(t *testing.T) {
	motion := loadVmd(t, "../../../test_resources/Lamb_2689F.vmd")

	model := loadPmx(t, "D:/MMD/MikuMikuDance_v926x64/UserFile/Model/ゲーム/戦国BASARA/幸村 たぬき式 ver.1.24/真田幸村没第二衣装1.24軽量版.pmx")

	boneDeltas, _ := computeBoneDeltas(model, motion, motionpkg.Frame(0), []string{pmx.LEG_IK.Right(), "右つま先", pmx.HEEL.Right()}, true, false, false)

	{
		expectedPosition := vec3(-1.216134, 1.887670, -10.78867)
		if !boneDeltas.GetByName(pmx.LEG_IK.Right()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.LEG_IK.Right()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.LEG_IK.Right()).FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(0.803149, 6.056844, -10.232766)
		if !boneDeltas.GetByName(pmx.LOWER.String()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.LOWER.String()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.LOWER.String()).FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(0.728442, 4.560226, -11.571869)
		if !boneDeltas.GetByName(pmx.LEG.Right()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.LEG.Right()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.LEG.Right()).FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(4.173470, 0.361388, -11.217197)
		if !boneDeltas.GetByName(pmx.KNEE.Right()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.KNEE.Right()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.KNEE.Right()).FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(-1.217569, 1.885731, -10.788104)
		if !boneDeltas.GetByName(pmx.ANKLE.Right()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.ANKLE.Right()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.ANKLE.Right()).FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(-0.922247, -1.163554, -10.794323)
		if !boneDeltas.GetByName("右つま先").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("右つま先").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("右つま先").FilledGlobalPosition()))
		}
	}
}

func TestVmdMotion_DeformLegIk17_Snow(t *testing.T) {
	// mlog.SetLevel(mlog.IK_VERBOSE)

	motion := loadVmd(t, "../../../test_resources/好き雪_1075.vmd")

	model := loadPmx(t, "D:/MMD/MikuMikuDance_v926x64/UserFile/Model/VOCALOID/初音ミク/Lat式ミクVer2.31/Lat式ミクVer2.31_White_準標準.pmx")

	boneDeltas, _ := computeBoneDeltas(model, motion, motionpkg.Frame(0), []string{"右つま先", pmx.HEEL.Right()}, true, false, false)

	{
		expectedPosition := vec3(2.049998, 12.957623, 1.477440)
		if !boneDeltas.GetByName(pmx.LOWER.String()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.LOWER.String()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.LOWER.String()).FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(1.201382, 11.353215, 2.266898)
		if !boneDeltas.GetByName(pmx.LEG.Right()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.LEG.Right()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.LEG.Right()).FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(0.443043, 7.640018, -1.308741)
		if !boneDeltas.GetByName(pmx.KNEE.Right()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.KNEE.Right()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.KNEE.Right()).FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(-0.574753, 7.943915, 3.279809)
		if !boneDeltas.GetByName(pmx.ANKLE.Right()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.ANKLE.Right()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.ANKLE.Right()).FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(-1.443098, 6.324932, 4.837177)
		if !boneDeltas.GetByName("右つま先").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("右つま先").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("右つま先").FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(-0.701516, 8.181108, 4.687274)
		if !boneDeltas.GetByName(pmx.HEEL.Right()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.HEEL.Right()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.HEEL.Right()).FilledGlobalPosition()))
		}
	}
}

func TestVmdMotion_DeformLegIk18_Syou(t *testing.T) {
	motion := loadVmd(t, "../../../test_resources/唱(ダンスのみ)_0-300F.vmd")

	model := loadPmx(t, "../../../test_resources/サンプルモデル.pmx")

	boneDeltas, _ := computeBoneDeltas(model, motion, motionpkg.Frame(107), []string{"右つま先", pmx.HEEL.Right()}, true, false, false)

	{
		expectedPosition := vec3(0.365000, 12.042871, 2.034023)
		if !boneDeltas.GetByName(pmx.LOWER.String()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.LOWER.String()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.LOWER.String()).FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(-0.488466, 10.920292, 2.626419)
		if !boneDeltas.GetByName(pmx.LEG.Right()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.LEG.Right()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.LEG.Right()).FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(-1.607765, 6.763937, 1.653586)
		if !boneDeltas.GetByName(pmx.KNEE.Right()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.KNEE.Right()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.KNEE.Right()).FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(-1.110289, 1.718307, 2.809817)
		if !boneDeltas.GetByName(pmx.ANKLE.Right()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.ANKLE.Right()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.ANKLE.Right()).FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(-1.753089, -0.026766, 1.173958)
		if !boneDeltas.GetByName("右つま先").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("右つま先").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("右つま先").FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(-0.952785, 0.078826, 2.838099)
		if !boneDeltas.GetByName(pmx.HEEL.Right()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.HEEL.Right()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.HEEL.Right()).FilledGlobalPosition()))
		}
	}
}

func TestVmdMotion_DeformLegIk19_Wa(t *testing.T) {
	motion := loadVmd(t, "../../../test_resources/129cm_001_10F.vmd")

	model := loadPmx(t, "D:/MMD/MikuMikuDance_v926x64/UserFile/Model/_VMDサイジング/wa_129cm 20231028/wa_129cm_bone-structure.pmx")

	boneDeltas, _ := computeBoneDeltas(model, motion, motionpkg.Frame(0), []string{"右つま先", pmx.HEEL.Right()}, true, false, false)
	{
		expectedPosition := vec3(0.000000, 9.900000, 0.000000)
		if !boneDeltas.GetByName(pmx.LOWER.String()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.LOWER.String()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.LOWER.String()).FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(-0.599319, 8.639606, 0.369618)
		if !boneDeltas.GetByName(pmx.LEG.Right()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.LEG.Right()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.LEG.Right()).FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(-0.486516, 6.323577, -2.217865)
		if !boneDeltas.GetByName(pmx.KNEE.Right()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.KNEE.Right()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.KNEE.Right()).FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(-2.501665, 2.859252, -1.902513)
		if !boneDeltas.GetByName(pmx.ANKLE.Right()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.ANKLE.Right()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.ANKLE.Right()).FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(-3.071062, 0.841962, -2.077063)
		if !boneDeltas.GetByName("右つま先").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("右つま先").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("右つま先").FilledGlobalPosition()))
		}
	}
}

func TestVmdMotion_DeformLegIk20_Syou(t *testing.T) {
	motion := loadVmd(t, "../../../test_resources/唱(ダンスのみ)_0-300F.vmd")

	model := loadPmx(t, "../../../test_resources/サンプルモデル.pmx")

	boneDeltas, _ := computeBoneDeltas(model, motion, motionpkg.Frame(107), []string{"右つま先", pmx.HEEL.Right()}, true, false, false)
	{
		expectedPosition := vec3(0.365000, 12.042871, 2.034023)
		if !boneDeltas.GetByName(pmx.LOWER.String()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.LOWER.String()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.LOWER.String()).FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(-0.488466, 10.920292, 2.626419)
		if !boneDeltas.GetByName(pmx.LEG.Right()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.LEG.Right()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.LEG.Right()).FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(-1.607765, 6.763937, 1.653586)
		if !boneDeltas.GetByName(pmx.KNEE.Right()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.KNEE.Right()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.KNEE.Right()).FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(-1.110289, 1.718307, 2.809817)
		if !boneDeltas.GetByName(pmx.ANKLE.Right()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.ANKLE.Right()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.ANKLE.Right()).FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(-1.753089, -0.026766, 1.173958)
		if !boneDeltas.GetByName("右つま先").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("右つま先").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("右つま先").FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(-0.952785, 0.078826, 2.838099)
		if !boneDeltas.GetByName(pmx.HEEL.Right()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.HEEL.Right()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.HEEL.Right()).FilledGlobalPosition()))
		}
	}
}

func TestVmdMotion_DeformLegIk21_FK(t *testing.T) {
	motion := loadVmd(t, "../../../test_resources/足FK.vmd")

	model := loadPmx(t, "../../../test_resources/サンプルモデル_ひざ制限なし.pmx")

	boneDeltas, _ := computeBoneDeltas(model, motion, motionpkg.Frame(0), []string{"右つま先", pmx.HEEL.Right()}, false, false, false)
	{
		expectedPosition := vec3(-0.133305, 10.693993, 2.314730)
		if !boneDeltas.GetByName(pmx.LEG.Right()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.LEG.Right()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.LEG.Right()).FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(2.708069, 9.216356, -0.720822)
		if !boneDeltas.GetByName(pmx.KNEE.Right()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.KNEE.Right()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.KNEE.Right()).FilledGlobalPosition()))
		}
	}
}

func TestVmdMotion_DeformLegIk22_Bake(t *testing.T) {
	motion := loadVmd(t, "../../../test_resources/足FK焼き込み.vmd")

	model := loadPmx(t, "../../../test_resources/サンプルモデル_ひざ制限なし.pmx")

	boneDeltas, _ := computeBoneDeltas(model, motion, motionpkg.Frame(0), []string{"右つま先", pmx.HEEL.Right()}, true, false, false)
	{
		expectedPosition := vec3(-0.133306, 10.693994, 2.314731)
		if !boneDeltas.GetByName(pmx.LEG.Right()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.LEG.Right()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.LEG.Right()).FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(-3.753989, 8.506582, 1.058842)
		if !boneDeltas.GetByName(pmx.KNEE.Right()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.KNEE.Right()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.KNEE.Right()).FilledGlobalPosition()))
		}
	}
}

func TestVmdMotion_DeformLegIk22_NoLimit(t *testing.T) {
	motion := loadVmd(t, "../../../test_resources/足FK.vmd")

	model := loadPmx(t, "../../../test_resources/サンプルモデル_ひざ制限なし.pmx")

	boneDeltas, _ := computeBoneDeltas(model, motion, motionpkg.Frame(0), []string{"右つま先", pmx.HEEL.Right()}, true, false, false)
	{
		expectedPosition := vec3(-0.133305, 10.693993, 2.314730)
		if !boneDeltas.GetByName(pmx.LEG.Right()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.LEG.Right()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.LEG.Right()).FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(2.081436, 7.884178, -0.268146)
		if !boneDeltas.GetByName(pmx.KNEE.Right()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.KNEE.Right()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.KNEE.Right()).FilledGlobalPosition()))
		}
	}
}

func TestVmdMotion_DeformLegIk23_Addiction(t *testing.T) {
	motion := loadVmd(t, "../../../test_resources/[A]ddiction_Lat式_0171F.vmd")

	model := loadPmx(t, "D:/MMD/MikuMikuDance_v926x64/UserFile/Model/VOCALOID/初音ミク/Tda式ミクワンピース/Tda式ミクワンピースRSP.pmx")

	boneDeltas, _ := computeBoneDeltas(model, motion, motionpkg.Frame(0), []string{pmx.TOE_IK.Right(), "右つま先"}, true, false, false)

	{
		expectedPosition := vec3(0, 0.2593031, 0)
		if !boneDeltas.GetByName(pmx.ROOT.String()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.ROOT.String()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.ROOT.String()).FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(-1.528317, 5.033707, 3.125487)
		if !boneDeltas.GetByName(pmx.LEG_IK.Right()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.LEG_IK.Right()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.LEG_IK.Right()).FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(0.609285, 12.001350, 1.666402)
		if !boneDeltas.GetByName(pmx.LOWER.String()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.LOWER.String()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.LOWER.String()).FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(-0.129098, 10.550634, 1.348259)
		if !boneDeltas.GetByName(pmx.LEG.Right()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.LEG.Right()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.LEG.Right()).FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(1.661012, 6.604201, -1.196993)
		if !boneDeltas.GetByName(pmx.KNEE.Right()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.KNEE.Right()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.KNEE.Right()).FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(-1.529553, 5.033699, 3.127081)
		if !boneDeltas.GetByName(pmx.ANKLE.Right()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.ANKLE.Right()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.ANKLE.Right()).FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(-2.044619, 3.204468, 2.877363)
		if !boneDeltas.GetByName("右つま先").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("右つま先").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("右つま先").FilledGlobalPosition()))
		}
	}
}

func TestVmdMotion_DeformLegIk24_Positive(t *testing.T) {
	// mlog.SetLevel(mlog.IK_VERBOSE)

	motion := loadVmd(t, "../../../test_resources/ポジティブパレード_0526.vmd")

	model := loadPmx(t, "D:/MMD/MikuMikuDance_v926x64/UserFile/Model/_VMDサイジング/wa_129cm 20231028/wa_129cm_20240406.pmx")
	boneDeltas, _ := computeBoneDeltas(model, motion, motionpkg.Frame(0), nil, true, false, false)
	{
		expectedPosition := vec3(0, 0, 0)
		if !boneDeltas.GetByName(pmx.ROOT.String()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.ROOT.String()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.ROOT.String()).FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(-3.312041, 6.310613, -1.134230)
		if !boneDeltas.GetByName(pmx.LEG_IK.Left()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.LEG_IK.Right()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.LEG_IK.Right()).FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(-2.754258, 7.935882, -2.298871)
		if !boneDeltas.GetByName(pmx.LOWER.String()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.LOWER.String()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.LOWER.String()).FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(-2.455364, 6.571013, -1.935295)
		if !boneDeltas.GetByName(pmx.LEG.Left()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.LEG.Right()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.LEG.Right()).FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(-2.695464, 4.323516, -4.574024)
		if !boneDeltas.GetByName(pmx.KNEE.Left()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.KNEE.Right()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.KNEE.Right()).FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(-3.322137, 6.302598, -1.131305)
		if !boneDeltas.GetByName("左脛骨").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左脛骨").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("左脛骨").FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(-2.575414, 5.447266, -3.254661)
		if !boneDeltas.GetByName("左足捩").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左足捩").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("左足捩").FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(-2.229677, 5.626327, -3.481028)
		if !boneDeltas.GetByName("左足捩先").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左足捩先").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("左足捩先").FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(-2.455364, 6.571013, -1.935295)
		if !boneDeltas.GetByName("左足向検A").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左足向検A").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("左足向検A").FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(-2.695177, 4.324148, -4.574588)
		if !boneDeltas.GetByName("左足向検A先").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左足向検A先").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("左足向検A先").FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(-2.695177, 4.324148, -4.574588)
		if !boneDeltas.GetByName("左足捩検B").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左足捩検B").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("左足捩検B").FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(-0.002697, 5.869486, -6.134800)
		if !boneDeltas.GetByName("左足捩検B先").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左足捩検B先").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("左足捩検B先").FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(-2.877639, 4.4450495, -4.164494)
		if !boneDeltas.GetByName("左膝補").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左膝補").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("左膝補").FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(-3.523895, 4.135535, -3.716305)
		if !boneDeltas.GetByName("左膝補先").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左膝補先").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("左膝補先").FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(-2.118768, 6.263350, -2.402574)
		if !boneDeltas.GetByName("左足w").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左足w").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("左足w").FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(-2.480717, 3.120446, -5.602753)
		if !boneDeltas.GetByName("左足w先").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左足w先").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("左足w先").FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(-2.455364, 6.571013, -1.935294)
		if !boneDeltas.GetByName("左足向-").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左足向-").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("左足向-").FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(-3.322137, 6.302598, -1.131305)
		if !boneDeltas.GetByName("左脛骨D").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左脛骨D").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("左脛骨D").FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(-3.199167, 3.952319, -4.391296)
		if !boneDeltas.GetByName("左脛骨D先").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左脛骨D先").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("左脛骨D先").FilledGlobalPosition()))
		}
	}
}

func TestVmdMotion_DeformArmIk(t *testing.T) {
	motion := loadVmd(t, "../../../test_resources/サンプルモーション.vmd")

	model := loadPmx(t, "../../../test_resources/ボーンツリーテストモデル.pmx")

	boneDeltas, _ := computeBoneDeltas(model, motion, motionpkg.Frame(3182), nil, true, false, false)
	{
		expectedPosition := vec3(0, 0, 0)
		if !boneDeltas.GetByName(pmx.ROOT.String()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.ROOT.String()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.ROOT.String()).FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(12.400011, 9.000000, 1.885650)
		if !boneDeltas.GetByName(pmx.CENTER.String()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.CENTER.String()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.CENTER.String()).FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(12.400011, 8.580067, 1.885650)
		if !boneDeltas.GetByName(pmx.GROOVE.String()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.GROOVE.String()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.GROOVE.String()).FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(12.400011, 11.628636, 2.453597)
		if !boneDeltas.GetByName(pmx.WAIST.String()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.WAIST.String()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.WAIST.String()).FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(12.400011, 12.567377, 1.229520)
		if !boneDeltas.GetByName(pmx.UPPER.String()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.UPPER.String()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.UPPER.String()).FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(12.344202, 13.782951, 1.178849)
		if !boneDeltas.GetByName(pmx.UPPER2.String()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.UPPER2.String()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.UPPER2.String()).FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(12.425960, 15.893852, 1.481421)
		if !boneDeltas.GetByName(pmx.SHOULDER_P.Left()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.SHOULDER_P.Left()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.SHOULDER_P.Left()).FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(12.425960, 15.893852, 1.481421)
		if !boneDeltas.GetByName(pmx.SHOULDER.Left()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.SHOULDER.Left()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.SHOULDER.Left()).FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(13.348320, 15.767927, 1.802947)
		if !boneDeltas.GetByName(pmx.ARM.Left()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.ARM.Left()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.ARM.Left()).FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(13.564770, 14.998386, 1.289923)
		if !boneDeltas.GetByName(pmx.ARM_TWIST.Left()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.ARM_TWIST.Left()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.ARM_TWIST.Left()).FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(14.043257, 13.297290, 0.155864)
		if !boneDeltas.GetByName(pmx.ELBOW.Left()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.ELBOW.Left()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.ELBOW.Left()).FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(13.811955, 13.552182, -0.388005)
		if !boneDeltas.GetByName(pmx.WRIST_TWIST.Left()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.WRIST_TWIST.Left()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.WRIST_TWIST.Left()).FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(13.144803, 14.287374, -1.956703)
		if !boneDeltas.GetByName(pmx.WRIST.Left()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.WRIST.Left()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.WRIST.Left()).FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(12.813587, 14.873419, -2.570278)
		if !boneDeltas.GetByName(pmx.INDEX1.Left()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.INDEX1.Left()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.INDEX1.Left()).FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(12.541822, 15.029200, -2.709604)
		if !boneDeltas.GetByName(pmx.INDEX2.Left()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.INDEX2.Left()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.INDEX2.Left()).FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(12.476499, 14.950351, -2.502167)
		if !boneDeltas.GetByName(pmx.INDEX3.Left()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.INDEX3.Left()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.INDEX3.Left()).FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(12.620306, 14.795185, -2.295859)
		if !boneDeltas.GetByName("左人指先").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左人指先").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("左人指先").FilledGlobalPosition()))
		}
	}
}

func TestVmdMotion_DeformArmIk3(t *testing.T) {
	motion := loadVmd(t, "C:/MMD/mlib_go/test_resources/Addiction_0F.vmd")

	model := loadPmx(t, "D:/MMD/MikuMikuDance_v926x64/UserFile/Model/VOCALOID/初音ミク/Sour式初音ミクVer.1.02/Black_全表示.pmx")

	boneDeltas, _ := computeBoneDeltas(model, motion, motionpkg.Frame(0), nil, true, false, false)
	{
		expectedPosition := vec3(1.018832, 15.840092, 0.532239)
		if !boneDeltas.GetByName("左腕").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左腕").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("左腕").FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(1.186002, 14.510550, 0.099023)
		if !boneDeltas.GetByName("左腕捩").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左腕捩").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("左腕捩").FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(1.353175, 13.181011, -0.334196)
		if !boneDeltas.GetByName("左ひじ").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左ひじ").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("左ひじ").FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(1.018832, 15.840092, 0.532239)
		if !boneDeltas.GetByName("左腕W").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左腕W").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("左腕W").FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(1.353175, 13.181011, -0.334196)
		if !boneDeltas.GetByName("左腕W先").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左腕W先").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("左腕W先").FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(1.353175, 13.181011, -0.334196)
		if !boneDeltas.GetByName("左腕WIK").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左腕WIK").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("左腕WIK").FilledGlobalPosition()))
		}
	}
}

func TestVmdMotion_DeformLegIk25_Ballet(t *testing.T) {
	// mlog.SetLevel(mlog.IK_VERBOSE)

	motion := loadVmd(t, "../../../test_resources/青江バレリーコ_1543F.vmd")

	model := loadPmx(t, "D:/MMD/MikuMikuDance_v926x64/UserFile/Model/刀剣乱舞/019_にっかり青江/にっかり青江 帽子屋式 ver2.1/帽子屋式にっかり青江（戦装束）_表示枠.pmx")

	boneDeltas, _ := computeBoneDeltas(model, motion, motionpkg.Frame(0), []string{"左つま先", pmx.HEEL.Left(), pmx.TOE_EX.Left()}, true, false, false)

	{
		expectedPosition := vec3(-4.374956, 13.203792, 1.554190)
		if !boneDeltas.GetByName(pmx.LOWER.String()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.LOWER.String()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.LOWER.String()).FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(-3.481956, 11.214747, 1.127255)
		if !boneDeltas.GetByName(pmx.LEG.Left()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.LEG.Left()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.LEG.Left()).FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(-7.173243, 7.787793, 0.013533)
		if !boneDeltas.GetByName(pmx.KNEE.Left()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.KNEE.Left()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.KNEE.Left()).FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(-11.529483, 3.689184, -1.119154)
		if !boneDeltas.GetByName(pmx.ANKLE.Left()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.ANKLE.Left()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.ANKLE.Left()).FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(-13.408189, 1.877100, -2.183821)
		if !boneDeltas.GetByName("左つま先").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左つま先").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("左つま先").FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(-12.545708, 4.008257, -0.932670)
		if !boneDeltas.GetByName(pmx.HEEL.Left()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.HEEL.Left()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.HEEL.Left()).FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(-3.481956, 11.214747, 1.127255)
		if !boneDeltas.GetByName(pmx.LEG_D.Left()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.LEG_D.Left()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.LEG_D.Left()).FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(-7.173243, 7.787793, 0.013533)
		if !boneDeltas.GetByName(pmx.KNEE_D.Left()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.KNEE_D.Left()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.KNEE_D.Left()).FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(-11.529483, 3.689184, -1.119154)
		if !boneDeltas.GetByName(pmx.ANKLE_D.Left()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.ANKLE_D.Left()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.ANKLE_D.Left()).FilledGlobalPosition()))
		}
	}
	// {
	// 	expectedPosition := vec3(-12.845280, 2.816309, -2.136874)
	// 	if !boneDeltas.GetByName(pmx.TOE_EX.Left()).GetGlobalPosition().NearEquals(expectedPosition, 0.03) {
	// 		t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.TOE_EX.Left()).GetGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.TOE_EX.Left()).GetGlobalPosition()))
	// 	}
	// }
}

func TestVmdMotion_DeformLegIk26_Far(t *testing.T) {
	// mlog.SetLevel(mlog.IK_VERBOSE)

	motion := loadVmd(t, "../../../test_resources/足IK乖離.vmd")

	model := loadPmx(t, "D:/MMD/MikuMikuDance_v926x64/UserFile/Model/_あにまさ式ミク準標準見せパン/初音ミクVer2 準標準 見せパン 3.pmx")
	boneDeltas, _ := computeBoneDeltas(model, motion, motionpkg.Frame(0), []string{"右つま先", pmx.TOE_EX.Right(), pmx.HEEL.Right()}, true, false, false)

	{
		expectedPosition := vec3(-0.796811, 10.752734, -0.072743)
		if !boneDeltas.GetByName(pmx.LEG.Right()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.LEG.Right()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.LEG.Right()).FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(-2.202487, 10.921064, -4.695134)
		if !boneDeltas.GetByName(pmx.KNEE.Right()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.KNEE.Right()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.KNEE.Right()).FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(-4.193142, 11.026311, -8.844866)
		if !boneDeltas.GetByName(pmx.ANKLE.Right()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.ANKLE.Right()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.ANKLE.Right()).FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(-5.108798, 10.935530, -11.494570)
		if !boneDeltas.GetByName("右つま先").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("右つま先").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("右つま先").FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(-4.800813, 10.964218, -10.612234)
		if !boneDeltas.GetByName(pmx.TOE_EX.Right()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.TOE_EX.Right()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.TOE_EX.Right()).FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(-4.331888, 12.178923, -9.514071)
		if !boneDeltas.GetByName(pmx.HEEL.Right()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.HEEL.Right()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.HEEL.Right()).FilledGlobalPosition()))
		}
	}
}

func TestVmdMotion_DeformLegIk27_Addiction_Shoes(t *testing.T) {
	// mlog.SetLevel(mlog.IK_VERBOSE)

	motion := loadVmd(t, "../../../test_resources/[A]ddiction_和洋_1074-1078F.vmd")

	model := loadPmx(t, "D:/MMD/MikuMikuDance_v926x64/UserFile/Model/_VMDサイジング/wa_129cm 20231028/wa_129cm_20240406.pmx")

	boneDeltas, _ := computeBoneDeltas(model, motion, motionpkg.Frame(2), nil, true, false, false)
	{
		expectedPosition := vec3(0, 0, 0)
		if !boneDeltas.GetByName(pmx.ROOT.String()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.ROOT.String()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.ROOT.String()).FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(1.406722, 1.841236, 0.277818)
		if !boneDeltas.GetByName(pmx.LEG_IK.Left()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.LEG_IK.Right()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.LEG_IK.Right()).FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(0.510231, 9.009953, 0.592482)
		if !boneDeltas.GetByName(pmx.LOWER.String()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.LOWER.String()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.LOWER.String()).FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(1.355914, 7.853320, 0.415251)
		if !boneDeltas.GetByName(pmx.LEG.Left()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.LEG.Right()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.LEG.Right()).FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(-0.327781, 5.203806, -1.073718)
		if !boneDeltas.GetByName(pmx.KNEE.Left()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.KNEE.Right()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.KNEE.Right()).FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(1.407848, 1.839228, 0.278700)
		if !boneDeltas.GetByName("左脛骨").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左脛骨").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("左脛骨").FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(1.407848, 1.839228, 0.278700)
		if !boneDeltas.GetByName("左脛骨D").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左脛骨D").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("左脛骨D").FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(-0.498054, 5.045506, -1.221016)
		if !boneDeltas.GetByName("左脛骨D先").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左脛骨D先").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("左脛骨D先").FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(1.462306, 7.684025, 0.087026)
		if !boneDeltas.GetByName("左足Dw").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左足Dw").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("左足Dw").FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(1.593721, 0.784840, -0.054141)
		if !boneDeltas.GetByName("左足先EX").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左足先EX").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("左足先EX").FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(1.551940, 1.045847, 0.034003)
		if !boneDeltas.GetByName("左素足先A").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左素足先A").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("左素足先A").FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(1.453982, 0.305976, -0.510022)
		if !boneDeltas.GetByName("左素足先A先").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左素足先A先").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("左素足先A先").FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(1.453982, 0.305976, -0.510022)
		if !boneDeltas.GetByName("左素足先AIK").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左素足先AIK").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("左素足先AIK").FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(0.941880, 2.132958, 0.020403)
		if !boneDeltas.GetByName("左素足先B").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左素足先B").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("左素足先B").FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(1.359364, 0.974298, -0.226041)
		if !boneDeltas.GetByName("左素足先B先").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左素足先B先").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("左素足先B先").FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(1.460890, 0.692527, -0.285973)
		if !boneDeltas.GetByName("左素足先BIK").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左素足先BIK").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("左素足先BIK").FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(1.173929, 2.066327, 0.182685)
		if !boneDeltas.GetByName("左靴調節").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左靴調節").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("左靴調節").FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(1.739235, 1.171441, 0.485052)
		if !boneDeltas.GetByName("左靴追従").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左靴追従").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("左靴追従").FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(1.186359, 2.046771, 0.189367)
		if !boneDeltas.GetByName("左靴追従先").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左靴追従先").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("左靴追従先").FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(1.173929, 2.066327, 0.182685)
		if !boneDeltas.GetByName("左靴追従IK").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左靴追従IK").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("左靴追従IK").FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(1.574899, 6.873434, 0.342768)
		if !boneDeltas.GetByName("左足補D").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左足補D").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("左足補D").FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(0.150401, 5.170907, -0.712416)
		if !boneDeltas.GetByName("左足補D先").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左足補D先").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("左足補D先").FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(0.150401, 5.170907, -0.712416)
		if !boneDeltas.GetByName("左足補DIK").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左足補DIK").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("左足補DIK").FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(1.355915, 7.853319, 0.415251)
		if !boneDeltas.GetByName("左足向検A").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左足向検A").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("左足向検A").FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(-0.327781, 5.203805, -1.073719)
		if !boneDeltas.GetByName("左足向検A先").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左足向検A先").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("左足向検A先").FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(-0.327781, 5.203805, -1.073719)
		if !boneDeltas.GetByName("左足向検AIK").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左足向検AIK").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("左足向検AIK").FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(1.355914, 7.853319, 0.415251)
		if !boneDeltas.GetByName("左足向-").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左足向-").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("左足向-").FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(1.264808, 7.561551, -0.161703)
		if !boneDeltas.GetByName("左足w").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左足w").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("左足w").FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(-0.714029, 3.930234, -1.935889)
		if !boneDeltas.GetByName("左足w先").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左足w先").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("左足w先").FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(0.016770, 5.319929, -0.781771)
		if !boneDeltas.GetByName("左膝補").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左膝補").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("左膝補").FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(-0.164672, 4.511360, -0.957886)
		if !boneDeltas.GetByName("左膝補先").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左膝補先").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("左膝補先").FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(-0.099887, 4.800064, -0.895003)
		if !boneDeltas.GetByName("左膝補IK").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左膝補IK").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("左膝補IK").FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(-0.327781, 5.203806, -1.073718)
		if !boneDeltas.GetByName("左足捩検B").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左足捩検B").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("左足捩検B").FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(-2.392915, 7.450026, -2.735495)
		if !boneDeltas.GetByName("左足捩検B先").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左足捩検B先").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("左足捩検B先").FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(-2.392915, 7.450026, -2.735495)
		if !boneDeltas.GetByName("左足捩検BIK").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左足捩検BIK").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("左足捩検BIK").FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(0.514067, 6.528563, -0.329234)
		if !boneDeltas.GetByName("左足捩").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左足捩").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("左足捩").FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(0.231636, 6.794109, -0.557747)
		if !boneDeltas.GetByName("左足捩先").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左足捩先").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("左足捩先").FilledGlobalPosition()))
		}
	}
}

func TestVmdMotion_DeformLegIk28_Gimme_Mitsu(t *testing.T) {
	// mlog.SetLevel(mlog.IK_VERBOSE)

	motion := loadVmd(t, "../../../test_resources/ぎみぎみ_498F.vmd")

	model := loadPmx(t, "D:/MMD/MikuMikuDance_v926x64/UserFile/Model/オリジナル/折岸みつ つみだんご/折岸みつ.pmx")

	boneDeltas, _ := computeBoneDeltas(model, motion, motionpkg.Frame(0), []string{"右つま先", pmx.HEEL.Right(), pmx.TOE_EX.Right()}, true, false, false)

	{
		expectedPosition := vec3(0.704942, 9.193451, 0.070969)
		if !boneDeltas.GetByName(pmx.LOWER.String()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.LOWER.String()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.LOWER.String()).FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(-0.954316, 7.572014, 1.019005)
		if !boneDeltas.GetByName(pmx.LEG.Right()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.LEG.Right()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.LEG.Right()).FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(0.545182, 5.180062, -2.267060)
		if !boneDeltas.GetByName(pmx.KNEE.Right()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.KNEE.Right()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.KNEE.Right()).FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(0.863384, 1.755991, 0.945758)
		if !boneDeltas.GetByName(pmx.ANKLE.Right()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.ANKLE.Right()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.ANKLE.Right()).FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(0.374198, 0.001257, -0.396838)
		if !boneDeltas.GetByName("右つま先").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("右つま先").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("右つま先").FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(0.472255, 0.627241, -0.116600)
		if !boneDeltas.GetByName(pmx.TOE_EX.Right()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.TOE_EX.Right()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.TOE_EX.Right()).FilledGlobalPosition()))
		}
	}
}

func TestVmdMotion_DeformLegIk28_Gimme_Mitsu_loop3(t *testing.T) {
	// mlog.SetLevel(mlog.IK_VERBOSE)

	motion := loadVmd(t, "../../../test_resources/ぎみぎみ_498F.vmd")

	model := loadPmx(t, "D:/MMD/MikuMikuDance_v926x64/UserFile/Model/オリジナル/折岸みつ つみだんご/折岸みつ_つま先IKループ3.pmx")

	boneDeltas, _ := computeBoneDeltas(model, motion, motionpkg.Frame(0), []string{"右つま先", pmx.HEEL.Right(), pmx.TOE_EX.Right()}, true, false, false)

	{
		expectedPosition := vec3(0.704942, 9.193451, 0.070969)
		if !boneDeltas.GetByName(pmx.LOWER.String()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.LOWER.String()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.LOWER.String()).FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(-0.954316, 7.572014, 1.019005)
		if !boneDeltas.GetByName(pmx.LEG.Right()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.LEG.Right()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.LEG.Right()).FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(0.545182, 5.180062, -2.267060)
		if !boneDeltas.GetByName(pmx.KNEE.Right()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.KNEE.Right()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.KNEE.Right()).FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(0.863384, 1.755991, 0.945758)
		if !boneDeltas.GetByName(pmx.ANKLE.Right()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.ANKLE.Right()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.ANKLE.Right()).FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(0.374198, 0.001257, -0.396838)
		if !boneDeltas.GetByName("右つま先").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("右つま先").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("右つま先").FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(0.472255, 0.627241, -0.116600)
		if !boneDeltas.GetByName(pmx.TOE_EX.Right()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.TOE_EX.Right()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.TOE_EX.Right()).FilledGlobalPosition()))
		}
	}
}

func TestVmdMotion_DeformLegIk28_Gimme_Mitsu_toe_order(t *testing.T) {
	// mlog.SetLevel(mlog.IK_VERBOSE)

	motion := loadVmd(t, "../../../test_resources/ぎみぎみ_498F.vmd")

	model := loadPmx(t, "D:/MMD/MikuMikuDance_v926x64/UserFile/Model/オリジナル/折岸みつ つみだんご/折岸みつ_つま先計算順前.pmx")

	boneDeltas, _ := computeBoneDeltas(model, motion, motionpkg.Frame(0), []string{"右つま先", pmx.HEEL.Right(), pmx.TOE_EX.Right()}, true, false, false)

	{
		expectedPosition := vec3(0.704942, 9.193451, 0.070969)
		if !boneDeltas.GetByName(pmx.LOWER.String()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.LOWER.String()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.LOWER.String()).FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(-0.954316, 7.572014, 1.019005)
		if !boneDeltas.GetByName(pmx.LEG.Right()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.LEG.Right()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.LEG.Right()).FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(0.545182, 5.180062, -2.267060)
		if !boneDeltas.GetByName(pmx.KNEE.Right()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.KNEE.Right()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.KNEE.Right()).FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(0.863384, 1.755991, 0.945758)
		if !boneDeltas.GetByName(pmx.ANKLE.Right()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.ANKLE.Right()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.ANKLE.Right()).FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(0.374198, 0.001257, -0.396838)
		if !boneDeltas.GetByName("右つま先").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("右つま先").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("右つま先").FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(0.596727, 0.597577, -0.123183)
		if !boneDeltas.GetByName(pmx.TOE_EX.Right()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.TOE_EX.Right()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.TOE_EX.Right()).FilledGlobalPosition()))
		}
	}
}

func TestVmdMotion_DeformLegIk28_Gimme_Miku(t *testing.T) {
	// mlog.SetLevel(mlog.IK_VERBOSE)

	motion := loadVmd(t, "../../../test_resources/ぎみぎみ_498F.vmd")

	model := loadPmx(t, "D:/MMD/MikuMikuDance_v926x64/UserFile/Model/_あにまさ式ミク準標準見せパン/初音ミクVer2 準標準 見せパン 3.pmx")

	boneDeltas, _ := computeBoneDeltas(model, motion, motionpkg.Frame(0), []string{"右つま先", pmx.HEEL.Right(), pmx.TOE_EX.Right()}, true, false, false)

	{
		expectedPosition := vec3(0.704942, 9.442580, 0.420454)
		if !boneDeltas.GetByName(pmx.LOWER.String()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.LOWER.String()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.LOWER.String()).FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(-0.784250, 7.438829, 1.217240)
		if !boneDeltas.GetByName(pmx.LEG.Right()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.LEG.Right()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.LEG.Right()).FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(1.200370, 4.815614, -2.325471)
		if !boneDeltas.GetByName(pmx.KNEE.Right()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.KNEE.Right()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.KNEE.Right()).FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(0.932342, 1.342473, 0.684434)
		if !boneDeltas.GetByName(pmx.ANKLE.Right()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.ANKLE.Right()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.ANKLE.Right()).FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(0.039125, -0.000434, -1.610423)
		if !boneDeltas.GetByName("右つま先").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("右つま先").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("右つま先").FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(0.219215, 0.218346, 0.837373)
		if !boneDeltas.GetByName(pmx.HEEL.Right()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.HEEL.Right()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.HEEL.Right()).FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(0.336862, 0.447203, -0.845470)
		if !boneDeltas.GetByName(pmx.TOE_EX.Right()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.TOE_EX.Right()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.TOE_EX.Right()).FilledGlobalPosition()))
		}
	}
}

func TestVmdMotion_DeformLegIk28_Gimme_Miku_toe_order(t *testing.T) {
	// mlog.SetLevel(mlog.IK_VERBOSE)

	motion := loadVmd(t, "../../../test_resources/ぎみぎみ_498F.vmd")

	model := loadPmx(t, "D:/MMD/MikuMikuDance_v926x64/UserFile/Model/_あにまさ式ミク準標準見せパン/初音ミクVer2 準標準 見せパン 3_つま先計算順後.pmx")

	boneDeltas, _ := computeBoneDeltas(model, motion, motionpkg.Frame(0), []string{"右つま先", pmx.HEEL.Right(), pmx.TOE_EX.Right()}, true, false, false)

	{
		expectedPosition := vec3(0.704942, 9.442580, 0.420454)
		if !boneDeltas.GetByName(pmx.LOWER.String()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.LOWER.String()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.LOWER.String()).FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(-0.784250, 7.438829, 1.217240)
		if !boneDeltas.GetByName(pmx.LEG.Right()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.LEG.Right()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.LEG.Right()).FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(1.200370, 4.815614, -2.325471)
		if !boneDeltas.GetByName(pmx.KNEE.Right()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.KNEE.Right()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.KNEE.Right()).FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(0.932342, 1.342473, 0.684434)
		if !boneDeltas.GetByName(pmx.ANKLE.Right()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.ANKLE.Right()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.ANKLE.Right()).FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(0.039126, -0.000434, -1.610423)
		if !boneDeltas.GetByName("右つま先").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("右つま先").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("右つま先").FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(1.146487, 0.023015, 0.590759)
		if !boneDeltas.GetByName(pmx.HEEL.Right()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.HEEL.Right()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.HEEL.Right()).FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(0.336863, 0.447203, -0.845470)
		if !boneDeltas.GetByName(pmx.TOE_EX.Right()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.TOE_EX.Right()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.TOE_EX.Right()).FilledGlobalPosition()))
		}
	}
}

func TestVmdMotion_DeformLegIk28_Gimme_Tda(t *testing.T) {
	// mlog.SetLevel(mlog.IK_VERBOSE)

	motion := loadVmd(t, "../../../test_resources/ぎみぎみ_498F.vmd")

	model := loadPmx(t, "D:/MMD/MikuMikuDance_v926x64/UserFile/Model/VOCALOID/初音ミク/Tda式初音ミク_盗賊つばき流Ｍトレースモデル配布 v1.07/Tda式初音ミク_盗賊つばき流Mトレースモデルv1.07_かかと.pmx")

	toeName := "足首_R_"
	boneDeltas, _ := computeBoneDeltas(model, motion, motionpkg.Frame(0), []string{"右つま先", pmx.HEEL.Right(), pmx.TOE_EX.Right(), toeName}, true, false, false)

	{
		expectedPosition := vec3(0.704941, 9.353957, 0.163552)
		if !boneDeltas.GetByName(pmx.LOWER.String()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.LOWER.String()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.LOWER.String()).FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(-0.561875, 8.374916, 0.596736)
		if !boneDeltas.GetByName(pmx.LEG.Right()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.LEG.Right()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.LEG.Right()).FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(1.184561, 5.035730, -2.609931)
		if !boneDeltas.GetByName(pmx.KNEE.Right()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.KNEE.Right()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.KNEE.Right()).FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(0.794708, 1.622819, 1.368421)
		if !boneDeltas.GetByName(pmx.ANKLE.Right()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.ANKLE.Right()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.ANKLE.Right()).FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(0.169612, 0.002951, -0.349215)
		if !boneDeltas.GetByName(toeName).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(toeName).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(toeName).FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(0.380133, 0.546320, 0.223149)
		if !boneDeltas.GetByName(pmx.TOE_EX.Right()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.TOE_EX.Right()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.TOE_EX.Right()).FilledGlobalPosition()))
		}
	}
}

func TestVmdMotion_DeformLegIk28_Gimme_Wa(t *testing.T) {
	// mlog.SetLevel(mlog.IK_VERBOSE)

	motion := loadVmd(t, "../../../test_resources/ぎみぎみ_498F.vmd")

	model := loadPmx(t, "D:/MMD/MikuMikuDance_v926x64/UserFile/Model/_VMDサイジング/wa_129cm 20231028/wa_129cm_20240406.pmx")

	boneDeltas, _ := computeBoneDeltas(model, motion, motionpkg.Frame(0), []string{"右つま先", pmx.HEEL.Right(), pmx.TOE_EX.Right()}, true, false, false)

	{
		expectedPosition := vec3(0.704942, 6.099999, 0.675723)
		if !boneDeltas.GetByName(pmx.LOWER.String()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.LOWER.String()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.LOWER.String()).FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(-0.258664, 5.228151, 1.304799)
		if !boneDeltas.GetByName(pmx.LEG.Right()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.LEG.Right()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.LEG.Right()).FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(1.135185, 4.090482, -1.667520)
		if !boneDeltas.GetByName(pmx.KNEE.Right()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.KNEE.Right()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.KNEE.Right()).FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(1.208500, 1.161467, 1.085173)
		if !boneDeltas.GetByName(pmx.ANKLE.Right()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.ANKLE.Right()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.ANKLE.Right()).FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(0.608805, -0.000103, -0.592631)
		if !boneDeltas.GetByName("右つま先").FilledGlobalPosition().NearEquals(expectedPosition, 0.031) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("右つま先").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("右つま先").FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(0.857188, 0.452834, 0.290514)
		if !boneDeltas.GetByName(pmx.TOE_EX.Right()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.TOE_EX.Right()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.TOE_EX.Right()).FilledGlobalPosition()))
		}
	}
}

func TestVmdMotion_DeformLegIk28_Gimme_Rin(t *testing.T) {
	// mlog.SetLevel(mlog.IK_VERBOSE)

	motion := loadVmd(t, "../../../test_resources/ぎみぎみ_498F.vmd")

	model := loadPmx(t, "D:/MMD/MikuMikuDance_v926x64/UserFile/Model/VOCALOID/鏡音リン/つみ式鏡音リン/つみ式鏡音リン.pmx")

	boneDeltas, _ := computeBoneDeltas(model, motion, motionpkg.Frame(0), []string{"右つま先", pmx.HEEL.Right(), pmx.TOE_EX.Right()}, true, false, false)

	{
		expectedPosition := vec3(0.704942, 8.031305, 0.156231)
		if !boneDeltas.GetByName(pmx.LOWER.String()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.LOWER.String()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.LOWER.String()).FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(-0.745842, 6.428913, 0.657299)
		if !boneDeltas.GetByName(pmx.LEG.Right()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.LEG.Right()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.LEG.Right()).FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(0.288399, 4.507010, -2.389454)
		if !boneDeltas.GetByName(pmx.KNEE.Right()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.KNEE.Right()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.KNEE.Right()).FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(0.899115, 1.789909, 0.907105)
		if !boneDeltas.GetByName(pmx.ANKLE.Right()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.ANKLE.Right()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.ANKLE.Right()).FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(0.469594, -0.000650, -0.271858)
		if !boneDeltas.GetByName("右つま先").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("右つま先").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("右つま先").FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(0.622100, 0.470730, 0.124509)
		if !boneDeltas.GetByName(pmx.TOE_EX.Right()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.TOE_EX.Right()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.TOE_EX.Right()).FilledGlobalPosition()))
		}
	}
}

func TestVmdMotion_DeformIk28_Simple(t *testing.T) {
	// mlog.SetLevel(mlog.IK_VERBOSE)

	motion := loadVmd(t, "../../../test_resources/IKの挙動を見たい_020.vmd")

	model := loadPmx(t, "../../../test_resources/IKの挙動を見たい.pmx")
	boneDeltas, _ := computeBoneDeltas(model, motion, motionpkg.Frame(0), nil, true, false, false)

	{
		expectedPosition := vec3(-9.433129, 1.363848, 1.867427)
		if !boneDeltas.GetByName("A+tail").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("A+tail").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("A+tail").FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(-9.433129, 1.363847, 1.867427)
		if !boneDeltas.GetByName("A+IK").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("A+IK").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("A+IK").FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(5.0, 4.517528, 2.142881)
		if !boneDeltas.GetByName("B+tail").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("B+tail").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("B+tail").FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(0.566871, 1.363847, 1.867427)
		if !boneDeltas.GetByName("B+IK").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("B+IK").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("B+IK").FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(10.0, 3.020634, 3.984441)
		if !boneDeltas.GetByName("C+tail").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("C+tail").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("C+tail").FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(5.566871, 1.363848, 1.867427)
		if !boneDeltas.GetByName("C+IK").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("C+IK").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("C+IK").FilledGlobalPosition()))
		}
	}
}

func TestVmdMotion_DeformIk29_Simple(t *testing.T) {
	// mlog.SetLevel(mlog.IK_VERBOSE)

	motion := loadVmd(t, "../../../test_resources/IKの挙動を見たい2_040.vmd")

	model := loadPmx(t, "../../../test_resources/IKの挙動を見たい2.pmx")
	boneDeltas, _ := computeBoneDeltas(model, motion, motionpkg.Frame(0), nil, true, false, false)
	{
		boneName := "A+2"
		expectedPosition := vec3(-5.440584, 2.324726, 0.816799)
		if !boneDeltas.GetByName(boneName).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(boneName).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(boneName).FilledGlobalPosition()))
		}
	}
	{
		boneName := "A+2tail"
		expectedPosition := vec3(-4.671312, 3.980981, -0.895119)
		if !boneDeltas.GetByName(boneName).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(boneName).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(boneName).FilledGlobalPosition()))
		}
	}
	{
		boneName := "B+2"
		expectedPosition := vec3(4.559244, 2.324562, 0.817174)
		if !boneDeltas.GetByName(boneName).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(boneName).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(boneName).FilledGlobalPosition()))
		}
	}
	{
		boneName := "B+2tail"
		expectedPosition := vec3(5.328533, 3.980770, -0.894783)
		if !boneDeltas.GetByName(boneName).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(boneName).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(boneName).FilledGlobalPosition()))
		}
	}
	{
		boneName := "C+2"
		expectedPosition := vec3(8.753987, 2.042284, -0.736314)
		if !boneDeltas.GetByName(boneName).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(boneName).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(boneName).FilledGlobalPosition()))
		}
	}
	{
		boneName := "C+2tail"
		expectedPosition := vec3(10.328943, 3.981413, -0.894101)
		if !boneDeltas.GetByName(boneName).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(boneName).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(boneName).FilledGlobalPosition()))
		}
	}
}

func TestVmdMotion_DeformArmIk2(t *testing.T) {
	// mlog.SetLevel(mlog.IK_VERBOSE)

	motion := loadVmd(t, "C:/MMD/mmd_base/tests/resources/唱(ダンスのみ)_0274F.vmd")

	model := loadPmx(t, "D:/MMD/MikuMikuDance_v926x64/UserFile/Model/VOCALOID/初音ミク/ISAO式ミク/I_ミクv4/Miku_V4_準標準.pmx")
	boneDeltas, _ := computeBoneDeltas(model, motion, motionpkg.Frame(0), nil, true, false, false)
	{
		expectedPosition := vec3(0.04952335, 9.0, 1.72378033)
		if !boneDeltas.GetByName(pmx.CENTER.String()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.CENTER.String()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.CENTER.String()).FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(0.04952335, 7.97980869, 1.72378033)
		if !boneDeltas.GetByName(pmx.GROOVE.String()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.GROOVE.String()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.GROOVE.String()).FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(0.04952335, 11.02838314, 2.29172656)
		if !boneDeltas.GetByName("腰").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("腰").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("腰").FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(0.04952335, 11.9671191, 1.06765032)
		if !boneDeltas.GetByName("上半身").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("上半身").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("上半身").FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(0.26284261, 13.14576297, 0.84720008)
		if !boneDeltas.GetByName("上半身2").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("上半身2").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("上半身2").FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(0.33636433, 15.27729547, 0.77435588)
		if !boneDeltas.GetByName("右肩").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("右肩").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("右肩").FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(-0.63104276, 15.44542768, 0.8507726)
		if !boneDeltas.GetByName("右肩C").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("右肩C").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("右肩C").FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(-0.63104276, 15.44542768, 0.8507726)
		if !boneDeltas.GetByName("右腕").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("右腕").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("右腕").FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(-0.90326269, 14.53727204, 0.7925801)
		if !boneDeltas.GetByName("右腕捩").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("右腕捩").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("右腕捩").FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(-1.50502977, 12.52976106, 0.66393998)
		if !boneDeltas.GetByName("右ひじ").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("右ひじ").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("右ひじ").FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(-1.46843236, 12.88476121, 0.12831076)
		if !boneDeltas.GetByName("右手捩").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("右手捩").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("右手捩").FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(-1.36287259, 13.90869981, -1.41662258)
		if !boneDeltas.GetByName("右手首").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("右手首").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("右手首").FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(-1.81521586, 14.00661535, -1.55616424)
		if !boneDeltas.GetByName("右手先").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("右手先").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("右手先").FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(-0.63104276, 15.44542768, 0.8507726)
		if !boneDeltas.GetByName("右腕YZ").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("右腕YZ").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("右腕YZ").FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(-0.72589296, 15.12898892, 0.83049645)
		if !boneDeltas.GetByName("右腕YZ先").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("右腕YZ先").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("右腕YZ先").FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(-0.72589374, 15.12898632, 0.83049628)
		if !boneDeltas.GetByName("右腕YZIK").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("右腕YZIK").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("右腕YZIK").FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(-0.63104276, 15.44542768, 0.8507726)
		if !boneDeltas.GetByName("右腕X").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("右腕X").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("右腕X").FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(-1.125321, 15.600293, 0.746130)
		if !boneDeltas.GetByName("右腕X先").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("右腕X先").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("右腕X先").FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(-1.1253241, 15.60029489, 0.7461294)
		if !boneDeltas.GetByName("右腕XIK").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("右腕XIK").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("右腕XIK").FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(-0.90325538, 14.53727326, 0.79258165)
		if !boneDeltas.GetByName("右腕捩YZ").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("右腕捩YZ").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("右腕捩YZ").FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(-1.01247534, 14.17289417, 0.76923367)
		if !boneDeltas.GetByName("右腕捩YZTgt").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("右腕捩YZTgt").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("右腕捩YZTgt").FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(-1.01248754, 14.17289597, 0.76923112)
		if !boneDeltas.GetByName("右腕捩YZIK").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("右腕捩YZIK").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("右腕捩YZIK").FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(-0.90325538, 14.53727326, 0.79258165)
		if !boneDeltas.GetByName("右腕捩X").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("右腕捩X").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("右腕捩X").FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(-1.40656426, 14.68386802, 0.85919594)
		if !boneDeltas.GetByName("右腕捩XTgt").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("右腕捩XTgt").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("右腕捩XTgt").FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(-1.40657579, 14.68387899, 0.8591982)
		if !boneDeltas.GetByName("右腕捩XIK").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("右腕捩XIK").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("右腕捩XIK").FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(-1.50499623, 12.52974836, 0.66394738)
		if !boneDeltas.GetByName("右ひじYZ").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("右ひじYZ").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("右ひじYZ").FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(-1.48334366, 12.74011791, 0.34655051)
		if !boneDeltas.GetByName("右ひじYZ先").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("右ひじYZ先").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("右ひじYZ先").FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(-1.48334297, 12.74012453, 0.34654052)
		if !boneDeltas.GetByName("右ひじYZIK").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("右ひじYZIK").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("右ひじYZIK").FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(-1.50499623, 12.52974836, 0.66394738)
		if !boneDeltas.GetByName("右ひじX").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("右ひじX").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("右ひじX").FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(-2.01179616, 12.66809052, 0.72106658)
		if !boneDeltas.GetByName("右ひじX先").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("右ひじX先").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("右ひじX先").FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(-2.00760407, 12.67958516, 0.7289003)
		if !boneDeltas.GetByName("右ひじXIK").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("右ひじXIK").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("右ひじXIK").FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(-1.50499623, 12.52974836, 0.66394738)
		if !boneDeltas.GetByName("右ひじY").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("右ひじY").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("右ひじY").FilledGlobalPosition()))
		}
	}
	{

		expectedPosition := vec3(-1.485519, 12.740760, 0.346835)
		if !boneDeltas.GetByName("右ひじY先").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("右ひじY先").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("右ひじY先").FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(-1.48334297, 12.74012453, 0.34654052)
		if !boneDeltas.GetByName("右ひじYIK").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("右ひじYIK").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("右ひじYIK").FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(-1.46845628, 12.88475892, 0.12832214)
		if !boneDeltas.GetByName("右手捩YZ").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("右手捩YZ").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("右手捩YZ").FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(-1.41168478, 13.4363328, -0.7038697)
		if !boneDeltas.GetByName("右手捩YZTgt").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("右手捩YZTgt").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("右手捩YZTgt").FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(-1.41156715, 13.43632015, -0.70389025)
		if !boneDeltas.GetByName("右手捩YZIK").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("右手捩YZIK").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("右手捩YZIK").FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(-1.46845628, 12.88475892, 0.12832214)
		if !boneDeltas.GetByName("右手捩X").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("右手捩X").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("右手捩X").FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(-1.5965686, 12.06213832, -0.42564769)
		if !boneDeltas.GetByName("右手捩XTgt").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("右手捩XTgt").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("右手捩XTgt").FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(-1.5965684, 12.06214091, -0.42565404)
		if !boneDeltas.GetByName("右手捩XIK").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("右手捩XIK").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("右手捩XIK").FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(-1.7198605, 13.98597326, -1.5267472)
		if !boneDeltas.GetByName("右手YZ先").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("右手YZ先").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("右手YZ先").FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(-1.71969424, 13.98593727, -1.52669587)
		if !boneDeltas.GetByName("右手YZIK").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("右手YZIK").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("右手YZIK").FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(-1.36306295, 13.90872698, -1.41659848)
		if !boneDeltas.GetByName("右手X").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("右手X").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("右手X").FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(-1.54727182, 13.56147176, -1.06342964)
		if !boneDeltas.GetByName("右手X先").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("右手X先").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("右手X先").FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(-1.54700171, 13.5614545, -1.0633896)
		if !boneDeltas.GetByName("右手XIK").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("右手XIK").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("右手XIK").FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(-0.90581859, 14.5370842, 0.80752276)
		if !boneDeltas.GetByName("右腕捩1").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("右腕捩1").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("右腕捩1").FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(-0.99954005, 14.2243783, 0.78748743)
		if !boneDeltas.GetByName("右腕捩2").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("右腕捩2").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("右腕捩2").FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(-1.10880907, 13.85976329, 0.76412793)
		if !boneDeltas.GetByName("右腕捩3").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("右腕捩3").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("右腕捩3").FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(-1.21298069, 13.51216081, 0.74185819)
		if !boneDeltas.GetByName("右腕捩4").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("右腕捩4").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("右腕捩4").FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(-1.5074743, 12.52953348, 0.67889319)
		if !boneDeltas.GetByName("右ひじsub").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("右ひじsub").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("右ひじsub").FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(-1.617075, 12.131149, 0.786797)
		if !boneDeltas.GetByName("右ひじsub先").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("右ひじsub先").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("右ひじsub先").FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(-1.472866, 12.872813, 0.120103)
		if !boneDeltas.GetByName("右手捩1").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("右手捩1").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("右手捩1").FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(-1.458749, 13.009759, -0.086526)
		if !boneDeltas.GetByName("右手捩2").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("右手捩2").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("右手捩2").FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(-1.440727, 13.184620, -0.350361)
		if !boneDeltas.GetByName("右手捩3").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("右手捩3").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("右手捩3").FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(-1.42368773, 13.34980879, -0.59962077)
		if !boneDeltas.GetByName("右手捩4").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("右手捩4").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("右手捩4").FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(-1.40457204, 13.511055, -0.84384039)
		if !boneDeltas.GetByName("右手捩5").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("右手捩5").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("右手捩5").FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(-1.39275926, 13.62582429, -1.01699954)
		if !boneDeltas.GetByName("右手捩6").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("右手捩6").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("右手捩6").FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(-1.36500465, 13.89623575, -1.42501008)
		if !boneDeltas.GetByName("右手首R").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("右手首R").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("右手首R").FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(-1.36500465, 13.89623575, -1.42501008)
		if !boneDeltas.GetByName("右手首1").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("右手首1").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("右手首1").FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(-1.472418, 13.917203, -1.529887)
		if !boneDeltas.GetByName("右手首2").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("右手首2").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("右手首2").FilledGlobalPosition()))
		}
	}
}

func TestVmdMotion_DeformLegIk25_Addiction_Wa_Right(t *testing.T) {
	// mlog.SetLevel(mlog.IK_VERBOSE)

	motion := loadVmd(t, "../../../test_resources/[A]ddiction_和洋_0126F.vmd")

	model := loadPmx(t, "D:/MMD/MikuMikuDance_v926x64/UserFile/Model/_VMDサイジング/wa_129cm 20231028/wa_129cm_20240406.pmx")
	boneDeltas, _ := computeBoneDeltas(model, motion, motionpkg.Frame(0), []string{"右襟先"}, true, false, false)

	{
		expectedPosition := vec3(-0.225006, 9.705784, 2.033072)
		if !boneDeltas.GetByName(pmx.UPPER.String()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.UPPER.String()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.UPPER.String()).FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(-0.237383, 10.769137, 2.039952)
		if !boneDeltas.GetByName(pmx.UPPER2.String()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.UPPER2.String()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.UPPER2.String()).FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(-0.630130, 13.306682, 2.752505)
		if !boneDeltas.GetByName(pmx.SHOULDER_P.Right()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.SHOULDER_P.Right()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.SHOULDER_P.Right()).FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(-0.630131, 13.306683, 2.742505)
		if !boneDeltas.GetByName(pmx.SHOULDER.Right()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.SHOULDER.Right()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.SHOULDER.Right()).FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(-0.948004, 13.753115, 2.690539)
		if !boneDeltas.GetByName(pmx.ARM.Right()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.ARM.Right()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.ARM.Right()).FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(-0.611438, 12.394744, 2.353463)
		if !boneDeltas.GetByName("右上半身C-A").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("右上半身C-A").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("右上半身C-A").FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(-0.664344, 13.835273, 2.403165)
		if !boneDeltas.GetByName("右鎖骨IK").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("右鎖骨IK").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("右鎖骨IK").FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(-0.270636, 13.350624, 2.258960)
		if !boneDeltas.GetByName("右鎖骨").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("右鎖骨").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("右鎖骨").FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(-0.963317, 14.098928, 2.497183)
		if !boneDeltas.GetByName("右鎖骨先").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("右鎖骨先").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("右鎖骨先").FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(-0.235138, 13.300934, 2.666039)
		if !boneDeltas.GetByName("右肩Rz検").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("右肩Rz検").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("右肩Rz検").FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(-0.847069, 13.997178, 2.886786)
		if !boneDeltas.GetByName("右肩Rz検先").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("右肩Rz検先").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("右肩Rz検先").FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(-0.235138, 13.300934, 2.666039)
		if !boneDeltas.GetByName("右肩Ry検").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("右肩Ry検").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("右肩Ry検").FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(-1.172100, 13.315790, 2.838742)
		if !boneDeltas.GetByName("右肩Ry検先").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("右肩Ry検先").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("右肩Ry検先").FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(-0.591152, 12.674325, 2.391185)
		if !boneDeltas.GetByName("右上半身C-B").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("右上半身C-B").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("右上半身C-B").FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(-0.588046, 12.954157, 2.432232)
		if !boneDeltas.GetByName("右上半身C-C").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("右上半身C-C").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("右上半身C-C").FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(-0.672292, 10.939227, 2.148515)
		if !boneDeltas.GetByName("右上半身2").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("右上半身2").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("右上半身2").FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(-0.520068, 14.089510, 2.812157)
		if !boneDeltas.GetByName("右襟").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("右襟").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("右襟").FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(-0.491354, 14.225309, 2.502640)
		if !boneDeltas.GetByName("右襟先").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("右襟先").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("右襟先").FilledGlobalPosition()))
		}
	}
}

func TestVmdMotion_DeformIk_Down(t *testing.T) {
	// mlog.SetLevel(mlog.IK_VERBOSE)

	motion := loadVmd(t, "../../../test_resources/センター下げる.vmd")

	model := loadPmx(t, "D:/MMD/MikuMikuDance_v926x64/UserFile/Model/_あにまさ式/MEIKO準標準_400.pmx")
	_, _ = computeBoneDeltas(model, motion, motionpkg.Frame(0), nil, true, false, false)
}

func TestVmdMotion_DeformArmIk4_DMF(t *testing.T) {
	// mlog.SetLevel(mlog.IK_VERBOSE)

	motion := loadVmd(t, "../../../test_resources/nac_dmf_601.vmd")

	model := loadPmx(t, "D:/MMD/MikuMikuDance_v926x64/UserFile/Model/VOCALOID/初音ミク/ISAO式ミク/I_ミクv4/Miku_V4.pmx")

	boneDeltas, _ := computeBoneDeltas(model, motion, motionpkg.Frame(0), nil, true, false, false)
	{
		expectedPosition := vec3(6.210230, 8.439670, 0.496305)
		if !boneDeltas.GetByName(pmx.CENTER.String()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.CENTER.String()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.CENTER.String()).FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(6.210230, 8.849669, 0.496305)
		if !boneDeltas.GetByName(pmx.GROOVE.String()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.GROOVE.String()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.GROOVE.String()).FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(6.210230, 12.836980, -0.159825)
		if !boneDeltas.GetByName("上半身").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("上半身").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("上半身").FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(6.261481, 13.968025, 0.288966)
		if !boneDeltas.GetByName("上半身2").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("上半身2").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("上半身2").FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(6.541666, 15.754716, 1.421828)
		if !boneDeltas.GetByName("左肩").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左肩").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("左肩").FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(7.451898, 16.031992, 1.675949)
		if !boneDeltas.GetByName("左腕").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左腕").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("左腕").FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(8.135534, 15.373729, 1.715530)
		if !boneDeltas.GetByName("左腕捩").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左腕捩").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("左腕捩").FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(9.646749, 13.918620, 1.803021)
		if !boneDeltas.GetByName("左ひじ").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左ひじ").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("左ひじ").FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(9.164164, 13.503792, 1.706635)
		if !boneDeltas.GetByName("左手捩").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左手捩").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("左手捩").FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(7.772219, 12.307291, 1.428628)
		if !boneDeltas.GetByName("左手首").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左手首").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("左手首").FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(7.390504, 12.011601, 1.405503)
		if !boneDeltas.GetByName("左手先").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左手先").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("左手先").FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(7.451900, 16.031990, 1.675949)
		if !boneDeltas.GetByName("左腕YZ").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左腕YZ").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("左腕YZ").FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(7.690105, 15.802624, 1.689741)
		if !boneDeltas.GetByName("左腕YZ先").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左腕YZ先").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("左腕YZ先").FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(7.690105, 15.802622, 1.689740)
		if !boneDeltas.GetByName("左腕YZIK").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左腕YZIK").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("左腕YZIK").FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(7.451899, 16.031988, 1.675950)
		if !boneDeltas.GetByName("左腕X").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左腕X").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("左腕X").FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(7.816861, 16.406412, 1.599419)
		if !boneDeltas.GetByName("左腕X先").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左腕X先").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("左腕X先").FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(7.816858, 16.406418, 1.599418)
		if !boneDeltas.GetByName("左腕XIK").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左腕XIK").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("左腕XIK").FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(8.135530, 15.373726, 1.715530)
		if !boneDeltas.GetByName("左腕捩YZ").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左腕捩YZ").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("左腕捩YZ").FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(8.409824, 15.109610, 1.731412)
		if !boneDeltas.GetByName("左腕捩YZTgt").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左腕捩YZTgt").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("左腕捩YZTgt").FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(8.409830, 15.109617, 1.731411)
		if !boneDeltas.GetByName("左腕捩YZIK").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左腕捩YZIK").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("左腕捩YZIK").FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(8.135530, 15.373725, 1.715531)
		if !boneDeltas.GetByName("左腕捩X").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左腕捩X").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("左腕捩X").FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(8.500528, 15.748149, 1.639511)
		if !boneDeltas.GetByName("左腕捩XTgt").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左腕捩XTgt").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("左腕捩XTgt").FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(8.500531, 15.748233, 1.639508)
		if !boneDeltas.GetByName("左腕捩XIK").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左腕捩XIK").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("左腕捩XIK").FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(9.646743, 13.918595, 1.803029)
		if !boneDeltas.GetByName("左ひじYZ").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左ひじYZ").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("左ひじYZ").FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(9.360763, 13.672787, 1.745903)
		if !boneDeltas.GetByName("左ひじYZ先").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左ひじYZ先").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("左ひじYZ先").FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(9.360781, 13.672805, 1.745905)
		if !boneDeltas.GetByName("左ひじYZIK").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左ひじYZIK").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("左ひじYZIK").FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(9.646734, 13.918593, 1.803028)
		if !boneDeltas.GetByName("左ひじX").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左ひじX").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("左ひじX").FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(9.944283, 13.652989, 1.456379)
		if !boneDeltas.GetByName("左ひじX先").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左ひじX先").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("左ひじX先").FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(9.944304, 13.653007, 1.456381)
		if !boneDeltas.GetByName("左ひじXIK").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左ひじXIK").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("左ひじXIK").FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(9.646734, 13.918596, 1.803028)
		if !boneDeltas.GetByName("左ひじY").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左ひじY").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("左ひじY").FilledGlobalPosition()))
		}
	}
	{
		// FIXME
		expectedPosition := vec3(9.560862, 13.926876, 1.431514)
		if !boneDeltas.GetByName("左ひじY先").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左ひじY先").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("左ひじY先").FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(9.360781, 13.672805, 1.745905)
		if !boneDeltas.GetByName("左ひじYIK").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左ひじYIK").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("左ひじYIK").FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(9.164141, 13.503780, 1.706625)
		if !boneDeltas.GetByName("左手捩YZ").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左手捩YZ").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("左手捩YZ").FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(8.414344, 12.859288, 1.556843)
		if !boneDeltas.GetByName("左手捩YZTgt").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左手捩YZTgt").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("左手捩YZTgt").FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(8.414370, 12.859282, 1.556885)
		if !boneDeltas.GetByName("左手捩YZIK").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左手捩YZIK").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("左手捩YZIK").FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(9.164142, 13.503780, 1.706624)
		if !boneDeltas.GetByName("左手捩X").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左手捩X").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("左手捩X").FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(9.511073, 12.928087, 2.447041)
		if !boneDeltas.GetByName("左手捩XTgt").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左手捩XTgt").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("左手捩XTgt").FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(9.511120, 12.928122, 2.447057)
		if !boneDeltas.GetByName("左手捩XIK").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左手捩XIK").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("左手捩XIK").FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(7.471097, 12.074032, 1.410383)
		if !boneDeltas.GetByName("左手YZ先").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左手YZ先").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("左手YZ先").FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(7.471111, 12.074042, 1.410384)
		if !boneDeltas.GetByName("左手YZIK").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左手YZIK").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("左手YZIK").FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(7.772183, 12.307314, 1.428564)
		if !boneDeltas.GetByName("左手X").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左手X").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("左手X").FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(7.802912, 12.308764, 0.901022)
		if !boneDeltas.GetByName("左手X先").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左手X先").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("左手X先").FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(7.802991, 12.308830, 0.901079)
		if !boneDeltas.GetByName("左手XIK").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左手XIK").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("左手XIK").FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(8.130125, 15.368912, 1.728851)
		if !boneDeltas.GetByName("左腕捩1").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左腕捩1").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("左腕捩1").FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(8.365511, 15.142246, 1.742475)
		if !boneDeltas.GetByName("左腕捩2").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左腕捩2").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("左腕捩2").FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(8.639965, 14.877952, 1.758356)
		if !boneDeltas.GetByName("左腕捩3").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左腕捩3").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("左腕捩3").FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(8.901615, 14.625986, 1.773497)
		if !boneDeltas.GetByName("左腕捩4").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左腕捩4").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("左腕捩4").FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(9.641270, 13.913721, 1.816324)
		if !boneDeltas.GetByName("左ひじsub").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左ひじsub").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("左ひじsub").FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(9.907782, 13.661371, 2.034630)
		if !boneDeltas.GetByName("左ひじsub先").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左ひじsub先").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("左ひじsub先").FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(9.165060, 13.499348, 1.721094)
		if !boneDeltas.GetByName("左手捩1").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左手捩1").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("左手捩1").FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(8.978877, 13.339340, 1.683909)
		if !boneDeltas.GetByName("左手捩2").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左手捩2").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("左手捩2").FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(8.741154, 13.135028, 1.636428)
		if !boneDeltas.GetByName("左手捩3").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左手捩3").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("左手捩3").FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(8.516553, 12.942023, 1.591578)
		if !boneDeltas.GetByName("左手捩4").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左手捩4").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("左手捩4").FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(8.301016, 12.748707, 1.544439)
		if !boneDeltas.GetByName("左手捩5").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左手捩5").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("左手捩5").FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(8.145000, 12.614601, 1.513277)
		if !boneDeltas.GetByName("左手捩6").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左手捩6").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("左手捩6").FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(7.777408, 12.298634, 1.439762)
		if !boneDeltas.GetByName("左手首R").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左手首R").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("左手首R").FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(7.777408, 12.298635, 1.439762)
		if !boneDeltas.GetByName("左手首1").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左手首1").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("左手首1").FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(7.670320, 12.202144, 1.486689)
		if !boneDeltas.GetByName("左手首2").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左手首2").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("左手首2").FilledGlobalPosition()))
		}
	}
}

func TestVmdMotion_DeformLegIk25_Addiction_Wa_Left(t *testing.T) {
	// mlog.SetLevel(mlog.IK_VERBOSE)

	motion := loadVmd(t, "../../../test_resources/[A]ddiction_和洋_0126F.vmd")

	model := loadPmx(t, "D:/MMD/MikuMikuDance_v926x64/UserFile/Model/_VMDサイジング/wa_129cm 20231028/wa_129cm_20240406.pmx")

	boneDeltas, _ := computeBoneDeltas(model, motion, motionpkg.Frame(0), []string{"左襟先"}, true, false, false)

	{
		expectedPosition := vec3(-0.225006, 9.705784, 2.033072)
		if !boneDeltas.GetByName(pmx.UPPER.String()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.UPPER.String()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.UPPER.String()).FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(-0.237383, 10.769137, 2.039952)
		if !boneDeltas.GetByName(pmx.UPPER2.String()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.UPPER2.String()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.UPPER2.String()).FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(0.460140, 13.290816, 2.531440)
		if !boneDeltas.GetByName(pmx.SHOULDER_P.Left()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.SHOULDER_P.Left()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.SHOULDER_P.Left()).FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(0.460140, 13.290816, 2.531440)
		if !boneDeltas.GetByName(pmx.SHOULDER.Left()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.SHOULDER.Left()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.SHOULDER.Left()).FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(0.784452, 13.728909, 2.608527)
		if !boneDeltas.GetByName(pmx.ARM.Left()).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.ARM.Left()).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(pmx.ARM.Left()).FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(0.272067, 12.381887, 2.182425)
		if !boneDeltas.GetByName("左上半身C-A").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左上半身C-A").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("左上半身C-A").FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(0.406217, 13.797803, 2.460243)
		if !boneDeltas.GetByName("左鎖骨IK").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左鎖骨IK").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("左鎖骨IK").FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(-0.052427, 13.347448, 2.216718)
		if !boneDeltas.GetByName("左鎖骨").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左鎖骨").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("左鎖骨").FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(0.659554, 14.017852, 2.591099)
		if !boneDeltas.GetByName("左鎖骨先").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左鎖骨先").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("左鎖骨先").FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(0.065147, 13.296564, 2.607907)
		if !boneDeltas.GetByName("左肩Rz検").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左肩Rz検").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("左肩Rz検").FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(0.517776, 14.134196, 2.645912)
		if !boneDeltas.GetByName("左肩Rz検先").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左肩Rz検先").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("左肩Rz検先").FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(0.065148, 13.296564, 2.607907)
		if !boneDeltas.GetByName("左肩Ry検").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左肩Ry検").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("左肩Ry検").FilledGlobalPosition()))
		}
	}
	{
		// FIXME
		expectedPosition := vec3(0.860159, 13.190875, 3.122428)
		if !boneDeltas.GetByName("左肩Ry検先").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左肩Ry検先").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("左肩Ry検先").FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(0.195053, 12.648546, 2.236849)
		if !boneDeltas.GetByName("左上半身C-B").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左上半身C-B").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("左上半身C-B").FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(0.294257, 12.912640, 2.257159)
		if !boneDeltas.GetByName("左上半身C-C").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左上半身C-C").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("左上半身C-C").FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(0.210011, 10.897711, 1.973442)
		if !boneDeltas.GetByName("左上半身2").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左上半身2").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("左上半身2").FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(0.320589, 14.049745, 2.637018)
		if !boneDeltas.GetByName("左襟").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左襟").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("左襟").FilledGlobalPosition()))
		}
	}
	{
		expectedPosition := vec3(0.297636, 14.263302, 2.374467)
		if !boneDeltas.GetByName("左襟先").FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左襟先").FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName("左襟先").FilledGlobalPosition()))
		}
	}
}

func TestVmdMotion_DeformArmIk_Mahoujin_02(t *testing.T) {
	// mlog.SetLevel(mlog.IK_VERBOSE)

	motion := loadVmd(t, "../../../test_resources/arm_ik_mahoujin_006F.vmd")

	model := loadPmx(t, "D:/MMD/MikuMikuDance_v926x64/UserFile/Model/刀剣乱舞/107_髭切/髭切mkmk009c 刀剣乱舞/髭切mkmk009c/髭切上着無mkmk009b_腕ＩＫ2.pmx")

	boneDeltas, _ := computeBoneDeltas(model, motion, motionpkg.Frame(0), nil, true, false, false)
	{
		boneName := pmx.ARM.Right()
		expectedPosition := vec3(-1.801768, 18.555544, 0.482812)
		if !boneDeltas.GetByName(boneName).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(boneName).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(boneName).FilledGlobalPosition()))
		}
	}
	{
		boneName := pmx.ELBOW.Right()
		expectedPosition := vec3(-3.273916, 17.405672, -2.046059)
		if !boneDeltas.GetByName(boneName).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(boneName).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(boneName).FilledGlobalPosition()))
		}
	}
	{
		boneName := pmx.WRIST.Right()
		expectedPosition := vec3(-1.240410, 18.910606, -4.062796)
		if !boneDeltas.GetByName(boneName).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(boneName).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(boneName).FilledGlobalPosition()))
		}
	}
	{
		boneName := pmx.INDEX3.Right()
		expectedPosition := vec3(-0.614190, 19.042362, -5.691705)
		if !boneDeltas.GetByName(boneName).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(boneName).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(boneName).FilledGlobalPosition()))
		}
	}
}

func TestVmdMotion_DeformArmIk_Mahoujin_03(t *testing.T) {
	// mlog.SetLevel(mlog.IK_VERBOSE)

	motion := loadVmd(t, "../../../test_resources/arm_ik_mahoujin_060F.vmd")

	model := loadPmx(t, "D:/MMD/MikuMikuDance_v926x64/UserFile/Model/刀剣乱舞/107_髭切/髭切mkmk009c 刀剣乱舞/髭切mkmk009c/髭切上着無mkmk009b_腕ＩＫ2.pmx")

	boneDeltas, _ := computeBoneDeltas(model, motion, motionpkg.Frame(0), nil, true, false, false)
	{
		boneName := pmx.ARM.Left()
		expectedPosition := vec3(1.801768, 18.555544, 0.457727)
		if !boneDeltas.GetByName(boneName).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(boneName).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(boneName).FilledGlobalPosition()))
		}
	}
	{
		boneName := pmx.ELBOW.Left()
		expectedPosition := vec3(4.422032, 18.073154, -1.174010)
		if !boneDeltas.GetByName(boneName).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(boneName).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(boneName).FilledGlobalPosition()))
		}
	}
	{
		boneName := pmx.WRIST.Left()
		expectedPosition := vec3(2.107284, 16.968552, -3.176913)
		if !boneDeltas.GetByName(boneName).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(boneName).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(boneName).FilledGlobalPosition()))
		}
	}
	{
		boneName := pmx.INDEX3.Left()
		expectedPosition := vec3(1.581160, 17.498112, -4.760089)
		if !boneDeltas.GetByName(boneName).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(boneName).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(boneName).FilledGlobalPosition()))
		}
	}
}

func TestVmdMotion_DeformArmIk_Choco_01(t *testing.T) {
	// mlog.SetLevel(mlog.IK_VERBOSE)

	motion := loadVmd(t, "../../../test_resources/ビタチョコ_0676F.vmd")

	model := loadPmx(t, "D:/MMD/MikuMikuDance_v926x64/UserFile/Model/ゲーム/Fate/眞白式ロマニ・アーキマン ver.1.01/眞白式ロマニ・アーキマン_ビタチョコ.pmx")

	boneDeltas, _ := computeBoneDeltas(model, motion, motionpkg.Frame(0), nil, true, false, false)
	{
		boneName := pmx.ARM.Left()
		expectedPosition := vec3(2.260640, 12.404558, -1.519635)
		if !boneDeltas.GetByName(boneName).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(boneName).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(boneName).FilledGlobalPosition()))
		}
	}
	{
		boneName := pmx.ELBOW.Left()
		expectedPosition := vec3(1.121608, 11.217656, -4.486015)
		if !boneDeltas.GetByName(boneName).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(boneName).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(boneName).FilledGlobalPosition()))
		}
	}
	{
		boneName := pmx.WRIST.Left()
		expectedPosition := vec3(0.717674, 13.924381, -3.561227)
		if !boneDeltas.GetByName(boneName).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(boneName).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(boneName).FilledGlobalPosition()))
		}
	}
	{
		boneName := pmx.INDEX3.Left()
		expectedPosition := vec3(1.002670, 15.652058, -3.506799)
		if !boneDeltas.GetByName(boneName).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(boneName).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(boneName).FilledGlobalPosition()))
		}
	}
	{
		boneName := pmx.ARM.Right()
		expectedPosition := vec3(-2.412614, 12.565295, -1.774290)
		if !boneDeltas.GetByName(boneName).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(boneName).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(boneName).FilledGlobalPosition()))
		}
	}
	{
		boneName := pmx.ELBOW.Right()
		expectedPosition := vec3(-1.009609, 11.296631, -4.589892)
		if !boneDeltas.GetByName(boneName).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(boneName).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(boneName).FilledGlobalPosition()))
		}
	}
	{
		boneName := pmx.WRIST.Right()
		expectedPosition := vec3(-0.137049, 14.029240, -4.235312)
		if !boneDeltas.GetByName(boneName).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(boneName).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(boneName).FilledGlobalPosition()))
		}
	}
	{
		boneName := pmx.INDEX3.Right()
		expectedPosition := vec3(-0.395239, 15.750233, -3.984484)
		if !boneDeltas.GetByName(boneName).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(boneName).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(boneName).FilledGlobalPosition()))
		}
	}
}

func TestVmdMotion_AdjustBones(t *testing.T) {
	// mlog.SetLevel(mlog.IK_VERBOSE)

	motion := loadVmd(t, "../../../test_resources/調整用ボーン移動.vmd")

	model := loadPmx(t, "D:/MMD/MikuMikuDance_v926x64/UserFile/Model/_あにまさ式ミク準標準見せパン/初音ミクVer2 準標準 見せパン 3_調整用ボーン追加.pmx")

	boneDeltas, _ := computeBoneDeltas(model, motion, motionpkg.Frame(0), nil, true, false, false)
	{
		boneName := pmx.CENTER.String()
		expectedPosition := vec3(1.84999, 8.0, -2.2)
		if !boneDeltas.GetByName(boneName).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(boneName).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(boneName).FilledGlobalPosition()))
		}
	}
	{
		boneName := pmx.GROOVE.String()
		expectedPosition := vec3(1.84999, 4.5, -2.2)
		if !boneDeltas.GetByName(boneName).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(boneName).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(boneName).FilledGlobalPosition()))
		}
	}
	{
		boneName := pmx.LOWER.String()
		expectedPosition := vec3(1.84999, 9.542581, -2.455269)
		if !boneDeltas.GetByName(boneName).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(boneName).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(boneName).FilledGlobalPosition()))
		}
	}
}

func TestVmdMotion_Neck(t *testing.T) {
	// mlog.SetLevel(mlog.IK_VERBOSE)

	motion := loadVmd(t, "../../../test_resources/くるりん_150F.vmd")

	model := loadPmx(t, "D:/MMD/MikuMikuDance_v926x64/UserFile/Model/VOCALOID/初音ミク/ISAO式ミク/I_ミクv4/Miku_V4.pmx")

	boneDeltas, _ := computeBoneDeltas(model, motion, motionpkg.Frame(0), []string{"頭"}, true, false, false)
	{
		boneName := pmx.NECK.String()
		expectedPosition := vec3(0.883310, 17.340812, -1.313977)
		if !boneDeltas.GetByName(boneName).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(boneName).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(boneName).FilledGlobalPosition()))
		}
	}
	{
		boneName := pmx.HEAD.String()
		expectedPosition := vec3(0.812887, 18.080100, -1.292382)
		if !boneDeltas.GetByName(boneName).FilledGlobalPosition().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(boneName).FilledGlobalPosition(), expectedPosition.Distance(boneDeltas.GetByName(boneName).FilledGlobalPosition()))
		}
	}
}

// loadVmd はVMDファイルを読み込んでモーションを返す。
func loadVmd(t *testing.T, path string) *motionpkg.VmdMotion {
	t.Helper()

	file, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			t.Skip("テスト用VMDが見つからないためスキップ")
		}
		t.Fatalf("VMDファイルの読み込み準備に失敗しました")
	}
	defer file.Close()

	motionData := motionpkg.NewVmdMotion(path)
	if info, err := file.Stat(); err == nil {
		motionData.SetFileModTime(info.ModTime().Unix())
	}

	reader := newBinaryReader(file)
	vmd := newVmdReader(reader)
	if err := vmd.readHeader(motionData); err != nil {
		t.Fatalf("VMDヘッダの読み込みに失敗しました: %v", err)
	}
	if err := vmd.readMotion(motionData); err != nil {
		t.Fatalf("VMD本体の読み込みに失敗しました: %v", err)
	}
	motionData.UpdateHash()

	return motionData
}

// loadPmx はPMXファイルを読み込んでモデルを返す。
func loadPmx(t *testing.T, path string) *model.PmxModel {
	t.Helper()

	file, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			t.Skip("テスト用PMXが見つからないためスキップ")
		}
		t.Fatalf("PMXファイルの読み込み準備に失敗しました")
	}
	defer file.Close()

	modelData := model.NewPmxModel()
	modelData.SetPath(path)
	if info, err := file.Stat(); err == nil {
		modelData.SetFileModTime(info.ModTime().Unix())
	}

	reader := newBinaryReader(file)
	pmx := newPmxReader(reader)
	if err := pmx.readHeader(modelData); err != nil {
		t.Fatalf("PMXヘッダの読み込みに失敗しました: %v", err)
	}
	if err := pmx.readModel(modelData); err != nil {
		t.Fatalf("PMX本体の読み込みに失敗しました: %v", err)
	}
	modelData.UpdateHash()

	return modelData
}

// computeBoneDeltas は行列更新まで実施したボーン差分を返す。
func computeBoneDeltas(
	modelData *model.PmxModel,
	motionData *motionpkg.VmdMotion,
	frame motionpkg.Frame,
	boneNames []string,
	includeIk bool,
	afterPhysics bool,
	removeTwist bool,
) (*delta.BoneDeltas, []int) {
	boneDeltas, indexes := ComputeBoneDeltas(
		modelData,
		motionData,
		frame,
		boneNames,
		includeIk,
		afterPhysics,
		removeTwist,
	)
	ApplyBoneMatrices(modelData, boneDeltas)
	return boneDeltas, indexes
}

// binaryReader はバイナリ読み込みを補助する。
type binaryReader struct {
	reader *bufio.Reader
}

// newBinaryReader はバイナリ読み込み用のリーダーを生成する。
func newBinaryReader(file *os.File) *binaryReader {
	return &binaryReader{reader: bufio.NewReader(file)}
}

// readBytes は指定サイズのバイト列を読み込む。
func (b *binaryReader) readBytes(size int) ([]byte, error) {
	if size <= 0 {
		return []byte{}, nil
	}
	buf := make([]byte, size)
	if _, err := io.ReadFull(b.reader, buf); err != nil {
		return nil, err
	}
	return buf, nil
}

// readUint8 はuint8を読み込む。
func (b *binaryReader) readUint8() (uint8, error) {
	buf, err := b.readBytes(1)
	if err != nil {
		return 0, err
	}
	return buf[0], nil
}

// readInt8 はint8を読み込む。
func (b *binaryReader) readInt8() (int8, error) {
	buf, err := b.readBytes(1)
	if err != nil {
		return 0, err
	}
	return int8(buf[0]), nil
}

// readUint16 はuint16を読み込む。
func (b *binaryReader) readUint16() (uint16, error) {
	buf, err := b.readBytes(2)
	if err != nil {
		return 0, err
	}
	return binary.LittleEndian.Uint16(buf), nil
}

// readInt16 はint16を読み込む。
func (b *binaryReader) readInt16() (int16, error) {
	buf, err := b.readBytes(2)
	if err != nil {
		return 0, err
	}
	return int16(binary.LittleEndian.Uint16(buf)), nil
}

// readUint32 はuint32を読み込む。
func (b *binaryReader) readUint32() (uint32, error) {
	buf, err := b.readBytes(4)
	if err != nil {
		return 0, err
	}
	return binary.LittleEndian.Uint32(buf), nil
}

// readInt32 はint32を読み込む。
func (b *binaryReader) readInt32() (int32, error) {
	buf, err := b.readBytes(4)
	if err != nil {
		return 0, err
	}
	return int32(binary.LittleEndian.Uint32(buf)), nil
}

// readFloat32 はfloat32を読み込んでfloat64で返す。
func (b *binaryReader) readFloat32() (float64, error) {
	value, err := b.readUint32()
	if err != nil {
		return 0, err
	}
	return float64(math.Float32frombits(value)), nil
}

// readFloat32s はfloat32配列を読み込む。
func (b *binaryReader) readFloat32s(values []float64) ([]float64, error) {
	for i := range values {
		v, err := b.readFloat32()
		if err != nil {
			return nil, err
		}
		values[i] = v
	}
	return values, nil
}

// readVec2 はVec2を読み込む。
func (b *binaryReader) readVec2() (mmath.Vec2, error) {
	x, err := b.readFloat32()
	if err != nil {
		return mmath.Vec2{}, err
	}
	y, err := b.readFloat32()
	if err != nil {
		return mmath.Vec2{}, err
	}
	return mmath.Vec2{X: x, Y: y}, nil
}

// readVec3 はVec3を読み込む。
func (b *binaryReader) readVec3() (mmath.Vec3, error) {
	x, err := b.readFloat32()
	if err != nil {
		return mmath.Vec3{}, err
	}
	y, err := b.readFloat32()
	if err != nil {
		return mmath.Vec3{}, err
	}
	z, err := b.readFloat32()
	if err != nil {
		return mmath.Vec3{}, err
	}
	return vec3(x, y, z), nil
}

// readVec4 はVec4を読み込む。
func (b *binaryReader) readVec4() (mmath.Vec4, error) {
	x, err := b.readFloat32()
	if err != nil {
		return mmath.Vec4{}, err
	}
	y, err := b.readFloat32()
	if err != nil {
		return mmath.Vec4{}, err
	}
	z, err := b.readFloat32()
	if err != nil {
		return mmath.Vec4{}, err
	}
	w, err := b.readFloat32()
	if err != nil {
		return mmath.Vec4{}, err
	}
	return mmath.Vec4{X: x, Y: y, Z: z, W: w}, nil
}

// readIndex はサイズ指定のインデックスを読み込む。
func (b *binaryReader) readIndex(size int) (int, error) {
	switch size {
	case 1:
		v, err := b.readInt8()
		return int(v), err
	case 2:
		v, err := b.readInt16()
		return int(v), err
	case 4:
		v, err := b.readInt32()
		return int(v), err
	default:
		return 0, fmt.Errorf("未対応のインデックスサイズ: %d", size)
	}
}

// decodeShiftJIS はShift-JISのバイト列をデコードする。
func decodeShiftJIS(raw []byte) (string, error) {
	decoded, err := japanese.ShiftJIS.NewDecoder().Bytes(raw)
	if err != nil {
		return "", err
	}
	decoded = bytes.TrimRight(decoded, "\xfd")
	decoded = bytes.TrimRight(decoded, "\x00")
	decoded = bytes.ReplaceAll(decoded, []byte("\x00"), []byte(" "))
	return string(decoded), nil
}

// decodeText は指定エンコーディングでテキストをデコードする。
func decodeText(enc encoding.Encoding, raw []byte) (string, error) {
	if enc == nil {
		return string(raw), nil
	}
	decoded, err := enc.NewDecoder().Bytes(raw)
	if err != nil {
		return "", err
	}
	return string(decoded), nil
}

// vmdReader はVMD読み込みの状態を保持する。
type vmdReader struct {
	br *binaryReader
}

// newVmdReader はVMDリーダーを生成する。
func newVmdReader(br *binaryReader) *vmdReader {
	return &vmdReader{br: br}
}

// readShiftJISText はShift-JIS固定長テキストを読み込む。
func (v *vmdReader) readShiftJISText(size int) (string, error) {
	raw, err := v.br.readBytes(size)
	if err != nil {
		return "", err
	}
	return decodeShiftJIS(raw)
}

// readHeader はVMDヘッダを読み込む。
func (v *vmdReader) readHeader(motionData *motionpkg.VmdMotion) error {
	signature, err := v.readShiftJISText(30)
	if err != nil {
		return err
	}
	name, err := v.readShiftJISText(20)
	if err != nil {
		return err
	}
	motionData.Signature = signature
	motionData.SetName(name)
	return nil
}

// readMotion はVMD各セクションを読み込む。
func (v *vmdReader) readMotion(motionData *motionpkg.VmdMotion) error {
	if err := v.readBones(motionData); err != nil {
		return err
	}
	if err := v.readMorphs(motionData); err != nil {
		return err
	}
	if err := v.readCameras(motionData); err != nil {
		return err
	}
	if err := v.readLights(motionData); err != nil {
		return err
	}
	if err := v.readShadows(motionData); err != nil {
		return err
	}
	if err := v.readIks(motionData); err != nil {
		return err
	}
	return nil
}

// readBones はボーンフレームを読み込む。
func (v *vmdReader) readBones(motionData *motionpkg.VmdMotion) error {
	total, err := v.br.readUint32()
	if err != nil {
		return err
	}
	curves := make([]byte, 64)
	values := make([]float64, 7)
	for i := 0; i < int(total); i++ {
		boneName, err := v.readShiftJISText(15)
		if err != nil {
			return err
		}
		index, err := v.br.readUint32()
		if err != nil {
			return err
		}
		bf := motionpkg.NewBoneFrame(motionpkg.Frame(index))
		bf.Read = true

		values, err = v.br.readFloat32s(values)
		if err != nil {
			return err
		}
		pos := vec3(values[0], values[1], values[2])
		bf.Position = &pos
		rot := mmath.NewQuaternionByValues(values[3], values[4], values[5], values[6])
		bf.Rotation = &rot

		rawCurves, err := v.br.readBytes(len(curves))
		if err != nil {
			return err
		}
		bf.Curves = motionpkg.NewBoneCurvesByValues(rawCurves)

		motionData.AppendBoneFrame(boneName, bf)
	}
	return nil
}

// readMorphs はモーフフレームを読み込む。
func (v *vmdReader) readMorphs(motionData *motionpkg.VmdMotion) error {
	total, err := v.br.readUint32()
	if err != nil {
		return err
	}
	for i := 0; i < int(total); i++ {
		morphName, err := v.readShiftJISText(15)
		if err != nil {
			return err
		}
		index, err := v.br.readUint32()
		if err != nil {
			return err
		}
		ratio, err := v.br.readFloat32()
		if err != nil {
			return err
		}
		mf := motionpkg.NewMorphFrame(motionpkg.Frame(index))
		mf.Read = true
		mf.Ratio = ratio
		motionData.AppendMorphFrame(morphName, mf)
	}
	return nil
}

// readCameras はカメラフレームを読み込む。
func (v *vmdReader) readCameras(motionData *motionpkg.VmdMotion) error {
	total, err := v.br.readUint32()
	if err != nil {
		if err == io.EOF || err == io.ErrUnexpectedEOF {
			return nil
		}
		return err
	}
	for i := 0; i < int(total); i++ {
		index, err := v.br.readUint32()
		if err != nil {
			return err
		}
		cf := motionpkg.NewCameraFrame(motionpkg.Frame(index))
		cf.Read = true

		distance, err := v.br.readFloat32()
		if err != nil {
			return err
		}
		cf.Distance = distance

		pos, err := v.br.readVec3()
		if err != nil {
			return err
		}
		cf.Position = &pos

		deg, err := v.br.readVec3()
		if err != nil {
			return err
		}
		cf.Degrees = &deg

		curves, err := v.br.readBytes(24)
		if err != nil {
			return err
		}
		cf.Curves = motionpkg.NewCameraCurvesByValues(curves)

		viewOfAngle, err := v.br.readUint32()
		if err != nil {
			return err
		}
		cf.ViewOfAngle = int(viewOfAngle)

		perspective, err := v.br.readUint8()
		if err != nil {
			return err
		}
		cf.IsPerspectiveOff = perspective == 1

		motionData.AppendCameraFrame(cf)
	}
	return nil
}

// readLights はライトフレームを読み込む。
func (v *vmdReader) readLights(motionData *motionpkg.VmdMotion) error {
	total, err := v.br.readUint32()
	if err != nil {
		if err == io.EOF || err == io.ErrUnexpectedEOF {
			return nil
		}
		return err
	}
	for i := 0; i < int(total); i++ {
		index, err := v.br.readUint32()
		if err != nil {
			return err
		}
		lf := motionpkg.NewLightFrame(motionpkg.Frame(index))
		lf.Read = true

		color, err := v.br.readVec3()
		if err != nil {
			return err
		}
		lf.Color = color

		pos, err := v.br.readVec3()
		if err != nil {
			return err
		}
		lf.Position = pos

		motionData.AppendLightFrame(lf)
	}
	return nil
}

// readShadows は影フレームを読み込む。
func (v *vmdReader) readShadows(motionData *motionpkg.VmdMotion) error {
	total, err := v.br.readUint32()
	if err != nil {
		if err == io.EOF || err == io.ErrUnexpectedEOF {
			return nil
		}
		return err
	}
	for i := 0; i < int(total); i++ {
		index, err := v.br.readUint32()
		if err != nil {
			return err
		}
		sf := motionpkg.NewShadowFrame(motionpkg.Frame(index))
		sf.Read = true

		mode, err := v.br.readUint8()
		if err != nil {
			return err
		}
		sf.ShadowMode = int(mode)

		distance, err := v.br.readFloat32()
		if err != nil {
			return err
		}
		sf.Distance = distance

		motionData.AppendShadowFrame(sf)
	}
	return nil
}

// readIks はIKフレームを読み込む。
func (v *vmdReader) readIks(motionData *motionpkg.VmdMotion) error {
	total, err := v.br.readUint32()
	if err != nil {
		if err == io.EOF || err == io.ErrUnexpectedEOF {
			return nil
		}
		return err
	}
	for i := 0; i < int(total); i++ {
		index, err := v.br.readUint32()
		if err != nil {
			return err
		}
		ikf := motionpkg.NewIkFrame(motionpkg.Frame(index))
		ikf.Read = true

		visible, err := v.br.readUint8()
		if err != nil {
			return err
		}
		ikf.Visible = visible == 1

		ikCount, err := v.br.readUint32()
		if err != nil {
			return err
		}
		for j := 0; j < int(ikCount); j++ {
			boneName, err := v.readShiftJISText(20)
			if err != nil {
				return err
			}
			enabled, err := v.br.readUint8()
			if err != nil {
				return err
			}
			ik := motionpkg.NewIkEnabledFrame(ikf.Index(), boneName)
			ik.Enabled = enabled == 1
			ikf.IkList = append(ikf.IkList, ik)
		}
		motionData.AppendIkFrame(ikf)
	}
	return nil
}

// pmxReader はPMX読み込みの状態を保持する。
type pmxReader struct {
	br                 *binaryReader
	textEncoding       encoding.Encoding
	extendedUVCount    int
	vertexIndexSize    int
	textureIndexSize   int
	materialIndexSize  int
	boneIndexSize      int
	morphIndexSize     int
	rigidBodyIndexSize int
}

// newPmxReader はPMXリーダーを生成する。
func newPmxReader(br *binaryReader) *pmxReader {
	return &pmxReader{br: br}
}

// readHeader はPMXヘッダを読み込む。
func (p *pmxReader) readHeader(modelData *model.PmxModel) error {
	signatureBytes, err := p.br.readBytes(4)
	if err != nil {
		return err
	}
	signature := string(signatureBytes)
	version, err := p.br.readFloat32()
	if err != nil {
		return err
	}
	if len(signature) < 3 || signature[:3] != "PMX" {
		return fmt.Errorf("PMX署名が不正です: %s", signature)
	}
	versionTag := fmt.Sprintf("%.1f", version)
	if versionTag != "2.0" && versionTag != "2.1" {
		return fmt.Errorf("PMXバージョンが不正です: %s", versionTag)
	}
	_, err = p.br.readUint8()
	if err != nil {
		return err
	}
	encodeType, err := p.br.readUint8()
	if err != nil {
		return err
	}
	switch encodeType {
	case 0:
		p.textEncoding = unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM)
	case 1:
		p.textEncoding = unicode.UTF8
	default:
		return fmt.Errorf("未対応エンコード種別です: %d", encodeType)
	}
	extendedUVCount, err := p.br.readUint8()
	if err != nil {
		return err
	}
	p.extendedUVCount = int(extendedUVCount)
	vertexIndexSize, err := p.br.readUint8()
	if err != nil {
		return err
	}
	p.vertexIndexSize = int(vertexIndexSize)
	textureIndexSize, err := p.br.readUint8()
	if err != nil {
		return err
	}
	p.textureIndexSize = int(textureIndexSize)
	materialIndexSize, err := p.br.readUint8()
	if err != nil {
		return err
	}
	p.materialIndexSize = int(materialIndexSize)
	boneIndexSize, err := p.br.readUint8()
	if err != nil {
		return err
	}
	p.boneIndexSize = int(boneIndexSize)
	morphIndexSize, err := p.br.readUint8()
	if err != nil {
		return err
	}
	p.morphIndexSize = int(morphIndexSize)
	rigidBodyIndexSize, err := p.br.readUint8()
	if err != nil {
		return err
	}
	p.rigidBodyIndexSize = int(rigidBodyIndexSize)

	name, err := p.readText()
	if err != nil {
		return err
	}
	modelData.SetName(name)
	englishName, err := p.readText()
	if err != nil {
		return err
	}
	modelData.EnglishName = englishName
	comment, err := p.readText()
	if err != nil {
		return err
	}
	modelData.Comment = comment
	englishComment, err := p.readText()
	if err != nil {
		return err
	}
	modelData.EnglishComment = englishComment
	return nil
}

// readText は可変長テキストを読み込む。
func (p *pmxReader) readText() (string, error) {
	size, err := p.br.readInt32()
	if err != nil {
		return "", err
	}
	if size == 0 {
		return "", nil
	}
	raw, err := p.br.readBytes(int(size))
	if err != nil {
		return "", err
	}
	return decodeText(p.textEncoding, raw)
}

// readVertexIndex は頂点インデックスを読み込む。
func (p *pmxReader) readVertexIndex() (int, error) {
	return p.br.readIndex(p.vertexIndexSize)
}

// readTextureIndex はテクスチャインデックスを読み込む。
func (p *pmxReader) readTextureIndex() (int, error) {
	return p.br.readIndex(p.textureIndexSize)
}

// readMaterialIndex は材質インデックスを読み込む。
func (p *pmxReader) readMaterialIndex() (int, error) {
	return p.br.readIndex(p.materialIndexSize)
}

// readBoneIndex はボーンインデックスを読み込む。
func (p *pmxReader) readBoneIndex() (int, error) {
	return p.br.readIndex(p.boneIndexSize)
}

// readMorphIndex はモーフインデックスを読み込む。
func (p *pmxReader) readMorphIndex() (int, error) {
	return p.br.readIndex(p.morphIndexSize)
}

// readRigidBodyIndex は剛体インデックスを読み込む。
func (p *pmxReader) readRigidBodyIndex() (int, error) {
	return p.br.readIndex(p.rigidBodyIndexSize)
}

// readModel はPMX本体を読み込む。
func (p *pmxReader) readModel(modelData *model.PmxModel) error {
	if err := p.readVertices(modelData); err != nil {
		return err
	}
	if err := p.readFaces(modelData); err != nil {
		return err
	}
	if err := p.readTextures(modelData); err != nil {
		return err
	}
	if err := p.readMaterials(modelData); err != nil {
		return err
	}
	if err := p.readBones(modelData); err != nil {
		return err
	}
	if err := p.readMorphs(modelData); err != nil {
		return err
	}
	if err := p.readDisplaySlots(modelData); err != nil {
		return err
	}
	if err := p.readRigidBodies(modelData); err != nil {
		return err
	}
	if err := p.readJoints(modelData); err != nil {
		return err
	}
	return nil
}

// readVertices は頂点を読み込む。
func (p *pmxReader) readVertices(modelData *model.PmxModel) error {
	count, err := p.br.readInt32()
	if err != nil {
		return err
	}
	for i := 0; i < int(count); i++ {
		vertex := &model.Vertex{}
		pos, err := p.br.readVec3()
		if err != nil {
			return err
		}
		vertex.Position = pos

		normal, err := p.br.readVec3()
		if err != nil {
			return err
		}
		vertex.Normal = normal

		uv, err := p.br.readVec2()
		if err != nil {
			return err
		}
		vertex.Uv = uv

		if p.extendedUVCount > 0 {
			vertex.ExtendedUvs = make([]mmath.Vec4, p.extendedUVCount)
			for j := 0; j < p.extendedUVCount; j++ {
				ext, err := p.br.readVec4()
				if err != nil {
					return err
				}
				vertex.ExtendedUvs[j] = ext
			}
		}

		weightType, err := p.br.readUint8()
		if err != nil {
			return err
		}

		switch weightType {
		case 0:
			boneIndex, err := p.readBoneIndex()
			if err != nil {
				return err
			}
			vertex.DeformType = model.BDEF1
			vertex.Deform = model.NewBdef1(boneIndex)
		case 1:
			boneIndex0, err := p.readBoneIndex()
			if err != nil {
				return err
			}
			boneIndex1, err := p.readBoneIndex()
			if err != nil {
				return err
			}
			weight0, err := p.br.readFloat32()
			if err != nil {
				return err
			}
			vertex.DeformType = model.BDEF2
			vertex.Deform = model.NewBdef2(boneIndex0, boneIndex1, weight0)
		case 2, 4:
			var indexes [4]int
			var weights [4]float64
			for j := 0; j < 4; j++ {
				idx, err := p.readBoneIndex()
				if err != nil {
					return err
				}
				indexes[j] = idx
			}
			for j := 0; j < 4; j++ {
				weight, err := p.br.readFloat32()
				if err != nil {
					return err
				}
				weights[j] = weight
			}
			vertex.DeformType = model.BDEF4
			vertex.Deform = model.NewBdef4(indexes, weights)
		case 3:
			boneIndex0, err := p.readBoneIndex()
			if err != nil {
				return err
			}
			boneIndex1, err := p.readBoneIndex()
			if err != nil {
				return err
			}
			weight0, err := p.br.readFloat32()
			if err != nil {
				return err
			}
			sdef := model.NewSdef(boneIndex0, boneIndex1, weight0)
			sdefC, err := p.br.readVec3()
			if err != nil {
				return err
			}
			sdefR0, err := p.br.readVec3()
			if err != nil {
				return err
			}
			sdefR1, err := p.br.readVec3()
			if err != nil {
				return err
			}
			sdef.SdefC = sdefC
			sdef.SdefR0 = sdefR0
			sdef.SdefR1 = sdefR1
			vertex.DeformType = model.SDEF
			vertex.Deform = sdef
		default:
			return fmt.Errorf("未対応ウェイト種別です: %d", weightType)
		}

		edgeFactor, err := p.br.readFloat32()
		if err != nil {
			return err
		}
		vertex.EdgeFactor = edgeFactor

		modelData.Vertices.Append(vertex)
	}
	return nil
}

// readFaces は面情報を読み込む。
func (p *pmxReader) readFaces(modelData *model.PmxModel) error {
	count, err := p.br.readInt32()
	if err != nil {
		return err
	}
	faceCount := int(count) / 3
	for i := 0; i < faceCount; i++ {
		idx0, err := p.readVertexIndex()
		if err != nil {
			return err
		}
		idx1, err := p.readVertexIndex()
		if err != nil {
			return err
		}
		idx2, err := p.readVertexIndex()
		if err != nil {
			return err
		}
		face := &model.Face{VertexIndexes: [3]int{idx0, idx1, idx2}}
		modelData.Faces.Append(face)
	}
	return nil
}

// readTextures はテクスチャ情報を読み込む。
func (p *pmxReader) readTextures(modelData *model.PmxModel) error {
	count, err := p.br.readInt32()
	if err != nil {
		return err
	}
	for i := 0; i < int(count); i++ {
		name, err := p.readText()
		if err != nil {
			return err
		}
		texture := model.NewTexture()
		texture.SetName(name)
		texture.SetValid(true)
		modelData.Textures.Append(texture)
	}
	return nil
}

// readMaterials は材質情報を読み込む。
func (p *pmxReader) readMaterials(modelData *model.PmxModel) error {
	count, err := p.br.readInt32()
	if err != nil {
		return err
	}
	for i := 0; i < int(count); i++ {
		mat := model.NewMaterial()
		name, err := p.readText()
		if err != nil {
			return err
		}
		mat.SetName(name)
		englishName, err := p.readText()
		if err != nil {
			return err
		}
		mat.EnglishName = englishName

		diffuse, err := p.br.readVec4()
		if err != nil {
			return err
		}
		mat.Diffuse = diffuse

		specular, err := p.br.readVec3()
		if err != nil {
			return err
		}
		specularPower, err := p.br.readFloat32()
		if err != nil {
			return err
		}
		mat.Specular = mmath.Vec4{X: specular.X, Y: specular.Y, Z: specular.Z, W: specularPower}

		ambient, err := p.br.readVec3()
		if err != nil {
			return err
		}
		mat.Ambient = ambient

		drawFlag, err := p.br.readUint8()
		if err != nil {
			return err
		}
		mat.DrawFlag = model.DrawFlag(drawFlag)

		edge, err := p.br.readVec4()
		if err != nil {
			return err
		}
		mat.Edge = edge

		edgeSize, err := p.br.readFloat32()
		if err != nil {
			return err
		}
		mat.EdgeSize = edgeSize

		textureIndex, err := p.readTextureIndex()
		if err != nil {
			return err
		}
		mat.TextureIndex = textureIndex

		sphereTextureIndex, err := p.readTextureIndex()
		if err != nil {
			return err
		}
		mat.SphereTextureIndex = sphereTextureIndex

		sphereMode, err := p.br.readUint8()
		if err != nil {
			return err
		}
		mat.SphereMode = model.SphereMode(sphereMode)

		toonSharing, err := p.br.readUint8()
		if err != nil {
			return err
		}
		mat.ToonSharingFlag = model.ToonSharingFlag(toonSharing)
		switch mat.ToonSharingFlag {
		case model.TOON_SHARING_INDIVIDUAL:
			toonTextureIndex, err := p.readTextureIndex()
			if err != nil {
				return err
			}
			mat.ToonTextureIndex = toonTextureIndex
		case model.TOON_SHARING_SHARING:
			toonTextureIndex, err := p.br.readUint8()
			if err != nil {
				return err
			}
			mat.ToonTextureIndex = int(toonTextureIndex)
		}

		memo, err := p.readText()
		if err != nil {
			return err
		}
		mat.Memo = memo

		verticesCount, err := p.br.readInt32()
		if err != nil {
			return err
		}
		mat.VerticesCount = int(verticesCount)

		modelData.Materials.Append(mat)
	}
	return nil
}

// readBones はボーン情報を読み込む。
func (p *pmxReader) readBones(modelData *model.PmxModel) error {
	count, err := p.br.readInt32()
	if err != nil {
		return err
	}
	for i := 0; i < int(count); i++ {
		bone := &model.Bone{DisplaySlotIndex: -1}
		name, err := p.readText()
		if err != nil {
			return err
		}
		bone.SetName(name)
		englishName, err := p.readText()
		if err != nil {
			return err
		}
		bone.EnglishName = englishName

		pos, err := p.br.readVec3()
		if err != nil {
			return err
		}
		bone.Position = pos

		parentIndex, err := p.readBoneIndex()
		if err != nil {
			return err
		}
		bone.ParentIndex = parentIndex

		layer, err := p.br.readInt32()
		if err != nil {
			return err
		}
		bone.Layer = int(layer)

		flagBytes, err := p.br.readBytes(2)
		if err != nil {
			return err
		}
		bone.BoneFlag = model.BoneFlag(uint16(flagBytes[0]) | uint16(flagBytes[1])<<8)

		if bone.BoneFlag&model.BONE_FLAG_TAIL_IS_BONE != 0 {
			tailIndex, err := p.readBoneIndex()
			if err != nil {
				return err
			}
			bone.TailIndex = tailIndex
		} else {
			offset, err := p.br.readVec3()
			if err != nil {
				return err
			}
			bone.TailIndex = -1
			bone.TailPosition = offset
		}

		if bone.BoneFlag&(model.BONE_FLAG_IS_EXTERNAL_ROTATION|model.BONE_FLAG_IS_EXTERNAL_TRANSLATION) != 0 {
			effectIndex, err := p.readBoneIndex()
			if err != nil {
				return err
			}
			effectFactor, err := p.br.readFloat32()
			if err != nil {
				return err
			}
			bone.EffectIndex = effectIndex
			bone.EffectFactor = effectFactor
		} else {
			bone.EffectIndex = -1
			bone.EffectFactor = 0
		}

		if bone.BoneFlag&model.BONE_FLAG_HAS_FIXED_AXIS != 0 {
			fixedAxis, err := p.br.readVec3()
			if err != nil {
				return err
			}
			bone.FixedAxis = fixedAxis
		}

		if bone.BoneFlag&model.BONE_FLAG_HAS_LOCAL_AXIS != 0 {
			localAxisX, err := p.br.readVec3()
			if err != nil {
				return err
			}
			localAxisZ, err := p.br.readVec3()
			if err != nil {
				return err
			}
			bone.LocalAxisX = localAxisX
			bone.LocalAxisZ = localAxisZ
		}

		if bone.BoneFlag&model.BONE_FLAG_IS_EXTERNAL_PARENT_DEFORM != 0 {
			effectorKey, err := p.br.readInt32()
			if err != nil {
				return err
			}
			bone.EffectorKey = int(effectorKey)
		}

		if bone.BoneFlag&model.BONE_FLAG_IS_IK != 0 {
			ik := &model.Ik{}
			targetIndex, err := p.readBoneIndex()
			if err != nil {
				return err
			}
			ik.BoneIndex = targetIndex
			loopCount, err := p.br.readInt32()
			if err != nil {
				return err
			}
			ik.LoopCount = int(loopCount)
			unitRot, err := p.br.readFloat32()
			if err != nil {
				return err
			}
			ik.UnitRotation = vec3(unitRot, unitRot, unitRot)
			linkCount, err := p.br.readInt32()
			if err != nil {
				return err
			}
			ik.Links = make([]model.IkLink, 0, int(linkCount))
			for j := 0; j < int(linkCount); j++ {
				link := model.IkLink{}
				linkIndex, err := p.readBoneIndex()
				if err != nil {
					return err
				}
				link.BoneIndex = linkIndex
				angleLimit, err := p.br.readUint8()
				if err != nil {
					return err
				}
				link.AngleLimit = angleLimit == 1
				if link.AngleLimit {
					minLimit, err := p.br.readVec3()
					if err != nil {
						return err
					}
					maxLimit, err := p.br.readVec3()
					if err != nil {
						return err
					}
					link.MinAngleLimit = minLimit
					link.MaxAngleLimit = maxLimit
				}
				ik.Links = append(ik.Links, link)
			}
			bone.Ik = ik
		}

		modelData.Bones.Append(bone)
	}
	return nil
}

// readMorphs はモーフ情報を読み込む。
func (p *pmxReader) readMorphs(modelData *model.PmxModel) error {
	count, err := p.br.readInt32()
	if err != nil {
		return err
	}
	boneOffsetValues := make([]float64, 7)
	materialOffsetValues := make([]float64, 28)
	for i := 0; i < int(count); i++ {
		morph := &model.Morph{}
		name, err := p.readText()
		if err != nil {
			return err
		}
		morph.SetName(name)
		englishName, err := p.readText()
		if err != nil {
			return err
		}
		morph.EnglishName = englishName

		panel, err := p.br.readUint8()
		if err != nil {
			return err
		}
		morph.Panel = model.MorphPanel(panel)

		morphType, err := p.br.readUint8()
		if err != nil {
			return err
		}
		morph.MorphType = model.MorphType(morphType)

		offsetCount, err := p.br.readInt32()
		if err != nil {
			return err
		}
		morph.Offsets = make([]model.MorphOffset, 0, int(offsetCount))
		for j := 0; j < int(offsetCount); j++ {
			switch morph.MorphType {
			case model.MORPH_TYPE_GROUP:
				morphIndex, err := p.readMorphIndex()
				if err != nil {
					return err
				}
				morphFactor, err := p.br.readFloat32()
				if err != nil {
					return err
				}
				morph.Offsets = append(morph.Offsets, &model.GroupMorphOffset{
					MorphIndex:  morphIndex,
					MorphFactor: morphFactor,
				})
			case model.MORPH_TYPE_VERTEX, model.MORPH_TYPE_AFTER_VERTEX:
				vertexIndex, err := p.readVertexIndex()
				if err != nil {
					return err
				}
				offset, err := p.br.readVec3()
				if err != nil {
					return err
				}
				morph.Offsets = append(morph.Offsets, &model.VertexMorphOffset{
					VertexIndex: vertexIndex,
					Position:    offset,
				})
			case model.MORPH_TYPE_BONE:
				boneIndex, err := p.readBoneIndex()
				if err != nil {
					return err
				}
				boneOffsetValues, err = p.br.readFloat32s(boneOffsetValues)
				if err != nil {
					return err
				}
				rot := mmath.NewQuaternionByValues(
					boneOffsetValues[3], boneOffsetValues[4], boneOffsetValues[5], boneOffsetValues[6],
				)
				morph.Offsets = append(morph.Offsets, &model.BoneMorphOffset{
					BoneIndex: boneIndex,
					Position:  vec3(boneOffsetValues[0], boneOffsetValues[1], boneOffsetValues[2]),
					Rotation:  rot,
				})
			case model.MORPH_TYPE_UV, model.MORPH_TYPE_EXTENDED_UV1, model.MORPH_TYPE_EXTENDED_UV2, model.MORPH_TYPE_EXTENDED_UV3, model.MORPH_TYPE_EXTENDED_UV4:
				vertexIndex, err := p.readVertexIndex()
				if err != nil {
					return err
				}
				uv, err := p.br.readVec4()
				if err != nil {
					return err
				}
				morph.Offsets = append(morph.Offsets, &model.UvMorphOffset{
					VertexIndex: vertexIndex,
					Uv:          uv,
					UvType:      morph.MorphType,
				})
			case model.MORPH_TYPE_MATERIAL:
				materialIndex, err := p.readMaterialIndex()
				if err != nil {
					return err
				}
				calcMode, err := p.br.readUint8()
				if err != nil {
					return err
				}
				materialOffsetValues, err = p.br.readFloat32s(materialOffsetValues)
				if err != nil {
					return err
				}
				morph.Offsets = append(morph.Offsets, &model.MaterialMorphOffset{
					MaterialIndex:       materialIndex,
					CalcMode:            model.MaterialMorphCalcMode(calcMode),
					Diffuse:             mmath.Vec4{X: materialOffsetValues[0], Y: materialOffsetValues[1], Z: materialOffsetValues[2], W: materialOffsetValues[3]},
					Specular:            mmath.Vec4{X: materialOffsetValues[4], Y: materialOffsetValues[5], Z: materialOffsetValues[6], W: materialOffsetValues[7]},
					Ambient:             vec3(materialOffsetValues[8], materialOffsetValues[9], materialOffsetValues[10]),
					Edge:                mmath.Vec4{X: materialOffsetValues[11], Y: materialOffsetValues[12], Z: materialOffsetValues[13], W: materialOffsetValues[14]},
					EdgeSize:            materialOffsetValues[15],
					TextureFactor:       mmath.Vec4{X: materialOffsetValues[16], Y: materialOffsetValues[17], Z: materialOffsetValues[18], W: materialOffsetValues[19]},
					SphereTextureFactor: mmath.Vec4{X: materialOffsetValues[20], Y: materialOffsetValues[21], Z: materialOffsetValues[22], W: materialOffsetValues[23]},
					ToonTextureFactor:   mmath.Vec4{X: materialOffsetValues[24], Y: materialOffsetValues[25], Z: materialOffsetValues[26], W: materialOffsetValues[27]},
				})
			default:
				return fmt.Errorf("未対応モーフ種別です: %d", morph.MorphType)
			}
		}
		modelData.Morphs.Append(morph)
	}
	return nil
}

// readDisplaySlots は表示枠情報を読み込む。
func (p *pmxReader) readDisplaySlots(modelData *model.PmxModel) error {
	count, err := p.br.readInt32()
	if err != nil {
		return err
	}
	for i := 0; i < int(count); i++ {
		slot := &model.DisplaySlot{}
		name, err := p.readText()
		if err != nil {
			return err
		}
		slot.SetName(name)
		englishName, err := p.readText()
		if err != nil {
			return err
		}
		slot.EnglishName = englishName

		specialFlag, err := p.br.readUint8()
		if err != nil {
			return err
		}
		slot.SpecialFlag = model.SpecialFlag(specialFlag)

		refCount, err := p.br.readInt32()
		if err != nil {
			return err
		}
		slot.References = make([]model.Reference, 0, int(refCount))
		for j := 0; j < int(refCount); j++ {
			refType, err := p.br.readUint8()
			if err != nil {
				return err
			}
			ref := model.Reference{DisplayType: model.DisplayType(refType)}
			switch ref.DisplayType {
			case model.DISPLAY_TYPE_BONE:
				index, err := p.readBoneIndex()
				if err != nil {
					return err
				}
				ref.DisplayIndex = index
			case model.DISPLAY_TYPE_MORPH:
				index, err := p.readMorphIndex()
				if err != nil {
					return err
				}
				ref.DisplayIndex = index
			default:
				return fmt.Errorf("未対応表示枠種別です: %d", ref.DisplayType)
			}
			slot.References = append(slot.References, ref)
		}
		modelData.DisplaySlots.Append(slot)
	}
	return nil
}

// readRigidBodies は剛体情報を読み込む。
func (p *pmxReader) readRigidBodies(modelData *model.PmxModel) error {
	count, err := p.br.readInt32()
	if err != nil {
		return err
	}
	for i := 0; i < int(count); i++ {
		body := &model.RigidBody{}
		name, err := p.readText()
		if err != nil {
			return err
		}
		body.SetName(name)
		englishName, err := p.readText()
		if err != nil {
			return err
		}
		body.EnglishName = englishName

		boneIndex, err := p.readBoneIndex()
		if err != nil {
			return err
		}
		body.BoneIndex = boneIndex

		group, err := p.br.readUint8()
		if err != nil {
			return err
		}
		mask, err := p.br.readUint16()
		if err != nil {
			return err
		}
		body.CollisionGroup = model.CollisionGroup{Group: group, Mask: mask}

		shape, err := p.br.readUint8()
		if err != nil {
			return err
		}
		body.Shape = model.Shape(shape)

		size, err := p.br.readVec3()
		if err != nil {
			return err
		}
		body.Size = size

		pos, err := p.br.readVec3()
		if err != nil {
			return err
		}
		body.Position = pos

		rot, err := p.br.readVec3()
		if err != nil {
			return err
		}
		body.Rotation = rot

		mass, err := p.br.readFloat32()
		if err != nil {
			return err
		}
		linearDamping, err := p.br.readFloat32()
		if err != nil {
			return err
		}
		angularDamping, err := p.br.readFloat32()
		if err != nil {
			return err
		}
		restitution, err := p.br.readFloat32()
		if err != nil {
			return err
		}
		friction, err := p.br.readFloat32()
		if err != nil {
			return err
		}
		body.Param = model.RigidBodyParam{
			Mass:           mass,
			LinearDamping:  linearDamping,
			AngularDamping: angularDamping,
			Restitution:    restitution,
			Friction:       friction,
		}

		physicsType, err := p.br.readUint8()
		if err != nil {
			return err
		}
		body.PhysicsType = model.PhysicsType(physicsType)

		modelData.RigidBodies.Append(body)
	}
	return nil
}

// readJoints はジョイント情報を読み込む。
func (p *pmxReader) readJoints(modelData *model.PmxModel) error {
	count, err := p.br.readInt32()
	if err != nil {
		return err
	}
	for i := 0; i < int(count); i++ {
		joint := &model.Joint{}
		name, err := p.readText()
		if err != nil {
			return err
		}
		joint.SetName(name)
		englishName, err := p.readText()
		if err != nil {
			return err
		}
		joint.EnglishName = englishName

		_, err = p.br.readUint8()
		if err != nil {
			return err
		}

		rigidBodyIndexA, err := p.readRigidBodyIndex()
		if err != nil {
			return err
		}
		rigidBodyIndexB, err := p.readRigidBodyIndex()
		if err != nil {
			return err
		}
		joint.RigidBodyIndexA = rigidBodyIndexA
		joint.RigidBodyIndexB = rigidBodyIndexB

		pos, err := p.br.readVec3()
		if err != nil {
			return err
		}
		rot, err := p.br.readVec3()
		if err != nil {
			return err
		}
		transMin, err := p.br.readVec3()
		if err != nil {
			return err
		}
		transMax, err := p.br.readVec3()
		if err != nil {
			return err
		}
		rotMin, err := p.br.readVec3()
		if err != nil {
			return err
		}
		rotMax, err := p.br.readVec3()
		if err != nil {
			return err
		}
		springTrans, err := p.br.readVec3()
		if err != nil {
			return err
		}
		springRot, err := p.br.readVec3()
		if err != nil {
			return err
		}
		joint.Param = model.JointParam{
			Position:                  pos,
			Rotation:                  rot,
			TranslationLimitMin:       transMin,
			TranslationLimitMax:       transMax,
			RotationLimitMin:          rotMin,
			RotationLimitMax:          rotMax,
			SpringConstantTranslation: springTrans,
			SpringConstantRotation:    springRot,
		}

		modelData.Joints.Append(joint)
	}
	return nil
}
