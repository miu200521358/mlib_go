// 指示: miu200521358
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/miu200521358/mlib_go/pkg/domain/deform"
	"github.com/miu200521358/mlib_go/pkg/domain/delta"
	"github.com/miu200521358/mlib_go/pkg/domain/mmath"
	"github.com/miu200521358/mlib_go/pkg/infra/file/mfile"
	"github.com/miu200521358/mlib_go/pkg/shared/contracts/mtime"

	"github.com/miu200521358/mlib_go/pkg/adapter/io_common"
	"github.com/miu200521358/mlib_go/pkg/adapter/io_model/pmx"
	"github.com/miu200521358/mlib_go/pkg/adapter/io_motion/vmd"
	"github.com/miu200521358/mlib_go/pkg/domain/model"
	"github.com/miu200521358/mlib_go/pkg/domain/motion"
)

// main は腕IKモーションをFKボーンへ焼き込んで保存する。
func main() {
	vmdPath := flag.String("vmd", "", "VMDパス")
	pmxPath := flag.String("pmx", "", "変形元PMXパス")
	dstPmxPath := flag.String("dst-pmx", "", "焼き込み先PMXパス(省略時は-pmxと同じ)")
	flag.Parse()
	if *vmdPath == "" || *pmxPath == "" {
		log.Fatal("-vmd と -pmx は必須です")
	}
	if *dstPmxPath == "" {
		*dstPmxPath = *pmxPath
	}
	outputPath := mfile.CreateOutputPath(*vmdPath, "baked")

	motionData, err := loadVmd(*vmdPath)
	if err != nil {
		log.Fatalf("VMD読み込みに失敗しました: %v", err)
	}
	sourceModelData, err := loadPmx(*pmxPath)
	if err != nil {
		log.Fatalf("PMX読み込みに失敗しました: %v", err)
	}
	targetModelData, err := loadPmx(*dstPmxPath)
	if err != nil {
		log.Fatalf("焼き込み先PMX読み込みに失敗しました: %v", err)
	}

	copiedMotion, err := motionData.Copy()
	if err != nil {
		log.Fatalf("VMD複製に失敗しました: %v", err)
	}
	bakedMotion := &copiedMotion

	armChains := [][]string{
		{
			model.ARM.Left(),
			model.ARM_TWIST.Left(),
			model.ELBOW.Left(),
			model.WRIST_TWIST.Left(),
			model.WRIST.Left(),
		},
		{
			model.ARM.Right(),
			model.ARM_TWIST.Right(),
			model.ELBOW.Right(),
			model.WRIST_TWIST.Right(),
			model.WRIST.Right(),
		},
	}

	for i := 0; i <= int(motionData.MaxFrame()); i++ {
		if i%100 == 0 {
			log.Printf("frame: %d", i)
		}

		f := mtime.Frame(i)

		bakedBoneDeltas := delta.NewBoneDeltas(targetModelData.Bones)

		sourceBoneDeltas, indexes := deform.ComputeBoneDeltas(
			sourceModelData,
			motionData,
			f,
			nil,
			true,
			false,
			false,
			nil,
		)
		deform.ApplyBoneMatricesWithIndexes(sourceModelData, sourceBoneDeltas, indexes)

		targetBoneDeltas, _ := deform.ComputeBoneDeltas(
			targetModelData,
			bakedMotion,
			f,
			nil,
			true,
			false,
			false,
			nil,
		)
		for _, chain := range armChains {
			for _, boneName := range chain[:len(chain)-1] {
				if !hasBone(sourceModelData, boneName) || !hasBone(targetModelData, boneName) {
					continue
				}
				d := deform.BakeBoneFrameByGlobalMatrix(
					targetModelData,
					targetBoneDeltas,
					bakedMotion,
					boneName,
					f,
					sourceBoneDeltas.GetByName(boneName).FilledGlobalMatrix(),
				)
				bakedBoneDeltas.Update(d)
			}
		}

		for _, chain := range armChains {
			for _, boneName := range chain[:len(chain)-1] {
				if !hasBone(sourceModelData, boneName) || !hasBone(targetModelData, boneName) {
					continue
				}
				bf := bakedMotion.BoneFrames.Get(boneName).Get(f)
				v := mmath.NewVec3()
				r := bakedBoneDeltas.GetByName(boneName).FilledFrameRotation()
				bf.Position = &v
				bf.Rotation = &r
				bakedMotion.InsertBoneFrame(boneName, bf)
			}
		}

		targetBoneDeltas, _ = deform.ComputeBoneDeltas(
			targetModelData,
			bakedMotion,
			f,
			nil,
			true,
			false,
			false,
			nil,
		)
		for _, chain := range armChains {
			wristBoneName := chain[len(chain)-1]
			if !hasBone(sourceModelData, wristBoneName) || !hasBone(targetModelData, wristBoneName) {
				continue
			}
			d := deform.BakeBoneFrameByGlobalMatrix(
				targetModelData,
				targetBoneDeltas,
				bakedMotion,
				wristBoneName,
				f,
				sourceBoneDeltas.GetByName(wristBoneName).FilledGlobalMatrix(),
			)
			bf := bakedMotion.BoneFrames.Get(wristBoneName).Get(f)
			v := mmath.NewVec3()
			r := d.FilledFrameRotation()
			bf.Position = &v
			bf.Rotation = &r
			bakedMotion.InsertBoneFrame(wristBoneName, bf)
		}
	}

	if err := saveMotion(outputPath, bakedMotion); err != nil {
		log.Fatalf("VMD保存に失敗しました: %v", err)
	}
	_, _ = fmt.Fprintf(os.Stdout, "VMD保存完了: %s\n", outputPath)
}

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

// saveMotion はVMDファイルとしてモデルを保存する。
func saveMotion(outPath string, m *motion.VmdMotion) error {
	if err := os.MkdirAll(filepath.Dir(outPath), 0o755); err != nil {
		return fmt.Errorf("保存先ディレクトリ作成に失敗: %w", err)
	}
	r := vmd.NewVmdRepository()
	if err := r.Save(outPath, m, io_common.SaveOptions{IncludeSystem: true}); err != nil {
		return fmt.Errorf("VMD保存に失敗: %w", err)
	}
	return nil
}
