// 指示: miu200521358
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/miu200521358/mlib_go/pkg/adapter/io_common"
	"github.com/miu200521358/mlib_go/pkg/adapter/io_model/pmx"
	"github.com/miu200521358/mlib_go/pkg/adapter/io_motion/vmd"
	"github.com/miu200521358/mlib_go/pkg/domain/deform"
	"github.com/miu200521358/mlib_go/pkg/domain/model"
	"github.com/miu200521358/mlib_go/pkg/domain/motion"
	"github.com/miu200521358/mlib_go/pkg/infra/file/mfile"
	"github.com/miu200521358/mlib_go/pkg/shared/contracts/mtime"
)

// main はFKで記録された足首の姿勢を足IKボーンへ焼き込んで保存する。
func main() {
	vmdPath := flag.String("vmd", "", "VMDパス")
	pmxPath := flag.String("pmx", "", "変形元PMXパス")
	flag.Parse()
	if *vmdPath == "" || *pmxPath == "" {
		log.Fatal("-vmd と -pmx は必須です")
	}

	outputPath := mfile.CreateOutputPath(*vmdPath, "legik")
	motionData, err := loadVmd(*vmdPath)
	if err != nil {
		log.Fatalf("VMD読み込みに失敗しました: %v", err)
	}
	modelData, err := loadPmx(*pmxPath)
	if err != nil {
		log.Fatalf("PMX読み込みに失敗しました: %v", err)
	}

	copiedMotion, err := motionData.Copy()
	if err != nil {
		log.Fatalf("VMD複製に失敗しました: %v", err)
	}
	outMotion := &copiedMotion

	legIKNames := []string{model.LEG_IK.Left(), model.LEG_IK.Right()}
	ankleNames := []string{model.ANKLE.Left(), model.ANKLE.Right()}
	for i := 0; i <= int(motionData.MaxFrame()); i++ {
		if i%100 == 0 {
			log.Printf("frame: %d", i)
		}

		frame := mtime.Frame(i)
		// 元モーションはIKを無効にして、FK足首の目標グローバル行列を確定する。
		fkDeltas, fkIndexes := deform.ComputeBoneDeltas(
			modelData,
			motionData,
			frame,
			nil,
			false,
			false,
			false,
			nil,
		)
		deform.ApplyBoneMatricesWithIndexes(modelData, fkDeltas, fkIndexes)

		// 出力側もIK無効で親ボーン行列を確定し、足IKの親相対分解に使う。
		outDeltas, outIndexes := deform.ComputeBoneDeltas(
			modelData,
			outMotion,
			frame,
			nil,
			false,
			false,
			false,
			nil,
		)
		deform.ApplyBoneMatricesWithIndexes(modelData, outDeltas, outIndexes)

		for side := range legIKNames {
			legIKName := legIKNames[side]
			ankleName := ankleNames[side]
			if !hasBone(modelData, legIKName) || !hasBone(modelData, ankleName) {
				continue
			}

			ankleDelta := fkDeltas.GetByName(ankleName)
			if ankleDelta == nil {
				continue
			}
			bakedDelta := deform.BakeBoneFrameByGlobalMatrix(
				modelData,
				outDeltas,
				outMotion,
				legIKName,
				frame,
				ankleDelta.FilledGlobalMatrix(),
			)
			if bakedDelta == nil {
				continue
			}
			// Bake関数内のGetは未登録フレームを格納しないため、値を明示して挿入する。
			boneFrame := outMotion.BoneFrames.Get(legIKName).Get(frame)
			position := bakedDelta.FilledFramePosition()
			rotation := bakedDelta.FilledFrameRotation()
			boneFrame.Position = &position
			boneFrame.Rotation = &rotation
			outMotion.InsertBoneFrame(legIKName, boneFrame)
		}
	}

	if err := saveMotion(outputPath, outMotion); err != nil {
		log.Fatalf("VMD保存に失敗しました: %v", err)
	}
	_, _ = fmt.Fprintf(os.Stdout, "VMD保存完了: %s\n", outputPath)
}

// hasBone はモデルに指定名のボーンが存在するかを返す。
func hasBone(modelData *model.PmxModel, boneName string) bool {
	if modelData == nil {
		return false
	}
	bone, err := modelData.Bones.GetByName(boneName)
	return err == nil && bone != nil
}

// loadVmd はVMDを読み込んで返す。
func loadVmd(path string) (*motion.VmdMotion, error) {
	repo := vmd.NewVmdRepository()
	data, err := repo.Load(path)
	if err != nil {
		return nil, err
	}
	motionData, ok := data.(*motion.VmdMotion)
	if !ok || motionData == nil {
		return nil, io_common.NewIoParseFailed("VMD読み込み結果が不正です", nil)
	}
	return motionData, nil
}

// loadPmx はPMXを読み込んで返す。
func loadPmx(path string) (*model.PmxModel, error) {
	repo := pmx.NewPmxRepository()
	data, err := repo.Load(path)
	if err != nil {
		return nil, err
	}
	modelData, ok := data.(*model.PmxModel)
	if !ok || modelData == nil {
		return nil, io_common.NewIoParseFailed("PMX読み込み結果が不正です", nil)
	}
	return modelData, nil
}

// saveMotion はVMDファイルとしてモーションを保存する。
func saveMotion(outPath string, motionData *motion.VmdMotion) error {
	if err := os.MkdirAll(filepath.Dir(outPath), 0o755); err != nil {
		return fmt.Errorf("保存先ディレクトリ作成に失敗: %w", err)
	}
	repo := vmd.NewVmdRepository()
	if err := repo.Save(outPath, motionData, io_common.SaveOptions{IncludeSystem: true}); err != nil {
		return fmt.Errorf("VMD保存に失敗: %w", err)
	}
	return nil
}
