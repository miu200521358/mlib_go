package main

import (
	"log"

	"github.com/miu200521358/mlib_go/pkg/domain/vmd"
	"github.com/miu200521358/mlib_go/pkg/infrastructure/repository"
)

func main() {
	vr := repository.NewVmdRepository(true)
	motionData, err := vr.Load("E:/downloads/M00-基本の動き_落ちる方.vmd")

	if err != nil {
		log.Fatalf("Expected error to be nil, got %q", err)
	}

	motion := motionData.(*vmd.VmdMotion)

	if err := vr.Save("E:/downloads/M00-基本の動き_落ちる方_copy_ikあり.vmd", motion, false); err != nil {
		log.Fatalf("Expected error to be nil, got %q", err)
	}

	motion.IkFrames = vmd.NewIkFrames()

	if err := vr.Save("E:/downloads/M00-基本の動き_落ちる方_copy_ikなし.vmd", motion, false); err != nil {
		log.Fatalf("Expected error to be nil, got %q", err)
	}
}
