// go run crumb/profile.go
// go tool pprof -http=:8080
package main

import (
	"fmt"
	"log"
	"time"

	"github.com/miu200521358/mlib_go/pkg/domain/pmx"
	"github.com/miu200521358/mlib_go/pkg/domain/vmd"
	"github.com/miu200521358/mlib_go/pkg/infrastructure/repository"
	"github.com/miu200521358/mlib_go/pkg/usecase/deform"
	"github.com/pkg/profile"
)

func main() {
	defer profile.Start(profile.CPUProfile, profile.ProfilePath(fmt.Sprintf("crumb/cpu_%s", time.Now().Format("20060102_150405")))).Stop()

	// defer profile.Start(profile.CPUProfile, profile.ProfilePath(".")).Stop()
	// defer profile.Start(profile.MemProfile, profile.ProfilePath(".")).Stop()

	// --------------------------------------------

	vr := repository.NewVmdRepository()
	motionData, err := vr.Load("D:/MMD/MikuMikuDance_v926x64/UserFile/Motion/ダンス_1人/[A]ddiction_モーション hino/[A]ddiction_Lat式.vmd")
	// motionData, err := vr.Load("D:/MMD/MikuMikuDance_v926x64/UserFile/Motion/ダンス_1人/CH4NGE mobiusP/CH4NGE.vmd")

	if err != nil {
		log.Fatalf("Expected error to be nil, got %q", err)
	}

	motion := motionData.(*vmd.VmdMotion)

	pr := repository.NewPmxRepository()
	// modelData, err := pr.Load("D:/MMD/MikuMikuDance_v926x64/UserFile/Model/_VMDサイジング/mlibkiller/mlibkiller.pmx")
	modelData, err := pr.Load("D:/MMD/MikuMikuDance_v926x64/UserFile/Model/VOCALOID/初音ミク/ISAO式ミク/I_ミクv4チャイナ/Miku_V4_チャイナ.pmx")

	if err != nil {
		log.Fatalf("Expected error to be nil, got %q", err)
	}

	model := modelData.(*pmx.PmxModel)

	for i := 0; i <= 500; i++ {
		if i%100 == 0 {
			log.Printf("i: %d", i)
		}
		deform.DeformBone(model, motion, true, i, nil)
	}
}
