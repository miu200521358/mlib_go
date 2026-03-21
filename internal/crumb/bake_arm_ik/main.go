// 指示: miu200521358
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/miu200521358/mlib_go/pkg/domain/deform"
	"github.com/miu200521358/mlib_go/pkg/domain/mmath"
	"github.com/miu200521358/mlib_go/pkg/infra/file/mfile"
	"github.com/miu200521358/mlib_go/pkg/shared/contracts/mtime"

	"github.com/miu200521358/mlib_go/pkg/adapter/io_common"
	"github.com/miu200521358/mlib_go/pkg/adapter/io_model/pmx"
	"github.com/miu200521358/mlib_go/pkg/adapter/io_motion/vmd"
	"github.com/miu200521358/mlib_go/pkg/domain/model"
	"github.com/miu200521358/mlib_go/pkg/domain/motion"
)

// main はダミーボーン組み合わせPMXを生成して保存する。
func main() {
	vmdPath := flag.String("vmd", "", "VMDパス")
	pmxPath := flag.String("pmx", "", "PMXパス")
	flag.Parse()
	outputPath := mfile.CreateOutputPath(*vmdPath, "baked")

	motionData, err := loadVmd(*vmdPath)
	if err != nil {
		log.Fatalf("VMD読み込みに失敗しました: %v", err)
	}
	modelData, err := loadPmx(*pmxPath)
	if err != nil {
		log.Fatalf("PMX読み込みに失敗しました: %v", err)
	}

	// 手首の角度がないやつ
	wristInitialMotionData, err := loadVmd(*vmdPath)
	if err != nil {
		log.Fatalf("VMD読み込みに失敗しました: %v", err)
	}
	wristInitialMotionData.BoneFrames.Delete("左手首")
	wristInitialMotionData.BoneFrames.Get("左腕ＩＫ").ForEach(func(frame motion.Frame, value *motion.BoneFrame) bool {
		q := mmath.NewQuaternion()
		value.Rotation = &q
		return true
	})

	bakedMotion := motion.NewVmdMotion(outputPath)
	for i := 0; i <= int(motionData.MaxFrame()); i++ {
		if i%100 == 0 {
			log.Printf("frame: %d", i)
		}

		f := mtime.Frame(i)

		boneDeltas, indexes := deform.ComputeBoneDeltas(
			modelData,
			motionData,
			f,
			[]string{"左腕", "左腕捩", "左ひじ", "左手捩", "左手首"},
			true,
			false,
			false,
			nil,
		)
		deform.ApplyBoneMatricesWithIndexes(modelData, boneDeltas, indexes)

		for _, boneName := range []string{"腕", "腕捩", "ひじ", "手捩", "手首"} {
			for _, direction := range []string{"左"} {
				directionBoneName := direction + boneName
				// directionBone, _ := modelData.Bones.GetByName(directionBoneName)

				d := boneDeltas.GetByName(directionBoneName)
				bf := motion.NewBoneFrame(f)
				q := d.FilledFrameRotation()
				// if boneName == "手首" {
				// 	// wristTailPosition, _ := directionBone.TailPosition.Copy()
				// 	// if directionBone.TailIndex > 0 {
				// 	// 	tailBone, _ := modelData.Bones.Get(directionBone.TailIndex)
				// 	// 	wristTailPosition = tailBone.Position.Subed(wristTailPosition)
				// 	// }

				// 	// wristInitialGlobalPosition := wristInitialboneDeltas.GetByName("左手首").FilledGlobalPosition()
				// 	// index1InitialGlobalPosition := wristInitialboneDeltas.GetByName("左人指１").FilledGlobalPosition()
				// 	// index1InitialDiff := index1InitialGlobalPosition.Subed(wristInitialGlobalPosition).Normalized()

				// 	// wristGlobalPosition := d.FilledGlobalPosition()
				// 	// index1GlobalPosition := boneDeltas.GetByName("左人指１").FilledGlobalPosition()
				// 	// index1Diff := index1GlobalPosition.Subed(wristGlobalPosition).Normalized()
				// 	// // index1DiffCross := index1InitialDiff.Cross(index1Diff)

				// 	// q = mmath.NewQuaternionFromDirection(index1InitialDiff, index1Diff)

				// 	// ikDelta := boneDeltas.GetByName("左腕ＩＫ")
				// 	// wristParentDelta := boneDeltas.Get(directionBone.ParentIndex)

				// 	// q = (wristParentDelta.FilledGlobalMatrix().Inverted().Muled(ikDelta.FilledGlobalMatrix())).Inverted().Muled(d.FilledGlobalMatrix()).Quaternion()

				// 	// index1Bone, _ := modelData.Bones.GetByName()
				// 	// diffDelta := delta.NewBoneDeltaByGlobalMatrix(index1Bone, f, index1GlobalMat, d)
				// 	// q = diffDelta.UnitMatrix.Quaternion()
				// 	ikDelta := boneDeltas.GetByName("左腕ＩＫ")
				// 	twistDelta := boneDeltas.GetByName("左手捩")
				// 	// q = twistDelta.FilledFrameRotation().Muled(ikDelta.FilledFrameRotation()).Muled(q)
				// 	q = q.Muled(ikDelta.FilledFrameRotation().Inverted()).Muled(twistDelta.FilledFrameRotation().Inverted())
				// }
				bf.Rotation = &q
				bakedMotion.InsertBoneFrame(directionBoneName, bf)
			}
		}
	}

	saveMotion(outputPath, bakedMotion)
	_, _ = fmt.Fprintf(os.Stdout, "VMD保存完了: %s\n", outputPath)
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
