package main

import (
	"log"

	"github.com/miu200521358/mlib_go/pkg/mutils/miter"
	"github.com/miu200521358/mlib_go/pkg/pmx"
	"github.com/miu200521358/mlib_go/pkg/vmd"
	"github.com/pkg/profile"
)

func main() {
	// // CPUプロファイル用のファイルを作成
	// {
	// 	f, err := os.Create("cpu.pprof")
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// 	defer f.Close()

	// 	// CPUプロファイリングを開始
	// 	if err := pprof.StartCPUProfile(f); err != nil {
	// 		panic(err)
	// 	}
	// 	defer pprof.StopCPUProfile()
	// }

	// {
	// 	// メモリプロファイル用のファイルを作成
	// 	f, err := os.Create("memory.pprof")
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// 	defer f.Close()
	// 	runtime.GC()

	// 	// ヒーププロファイリングを開始
	// 	if err := pprof.WriteHeapProfile(f); err != nil {
	// 		panic(err)
	// 	}
	// }

	// // CPUプロファイル
	// defer profile.Start(profile.CPUProfile, profile.ProfilePath(".")).Stop()

	// メモリプロファイル
	defer profile.Start(profile.MemProfile, profile.ProfilePath(".")).Stop()

	// --------------------------------------------

	vr := &vmd.VmdMotionReader{}
	motionData, err := vr.ReadByFilepath("D:/MMD/MikuMikuDance_v926x64/UserFile/Motion/ダンス_1人/[A]ddiction_モーション hino/[A]ddiction_Lat式.vmd")
	// motionData, err := vr.ReadByFilepath("D:/MMD/MikuMikuDance_v926x64/UserFile/Motion/ダンス_1人/CH4NGE mobiusP/CH4NGE.vmd")

	if err != nil {
		log.Fatalf("Expected error to be nil, got %q", err)
	}

	motion := motionData.(*vmd.VmdMotion)

	pr := &pmx.PmxReader{}
	modelData, err := pr.ReadByFilepath("D:/MMD/MikuMikuDance_v926x64/UserFile/Model/VOCALOID/初音ミク/ISAO式ミク/I_ミクv4チャイナ/Miku_V4_チャイナ.pmx")

	if err != nil {
		log.Fatalf("Expected error to be nil, got %q", err)
	}

	model := modelData.(*pmx.PmxModel)
	model.SetUp()

	miter.SetBlockSize(12)

	for i := 0; i < 500; i++ {
		if i%100 == 0 {
			log.Printf("i: %d", i)
		}
		motion.AnimateBone(i, model, nil, true)
	}
}
