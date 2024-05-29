package main

import (
	"log"

	"github.com/pkg/profile"

	"github.com/miu200521358/mlib_go/pkg/pmx"
	"github.com/miu200521358/mlib_go/pkg/vmd"
)

func main() {
	defer profile.Start(profile.CPUProfile, profile.ProfilePath(".")).Stop()
	// defer profile.Start(profile.MemProfile, profile.ProfilePath(".")).Stop()

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

	for i := 0; i < 500; i++ {
		if i%100 == 0 {
			log.Printf("i: %d", i)
		}
		motion.BoneFrames.Deform(i, model, nil, true, nil)
	}
}
