// 指示: miu200521358
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/pkg/profile"
)

const (
	defaultVmdPath = "internal/test_resources/サンプルモーション.vmd"
	defaultPmxPath = "internal/test_resources/サンプルモデル.pmx"
)

// main はCPUプロファイルを取得しながらボーン変形を実行する。
func main() {
	vmdPath := flag.String("vmd", "", "VMDパス")
	pmxPath := flag.String("pmx", "", "PMXパス")
	frameCount := flag.Int("frames", 500, "計測するフレーム数")
	// includeIk := flag.Bool("ik", true, "IKを含める")
	// removeTwist := flag.Bool("remove-twist", false, "捩り除去を行う")
	flag.Parse()

	resolvedVmdPath := resolveProfilePath(*vmdPath, "MLIB_PROFILE_VMD", defaultVmdPath)
	resolvedPmxPath := resolveProfilePath(*pmxPath, "MLIB_PROFILE_PMX", defaultPmxPath)

	profileDir := filepath.Join("internal", "profile", fmt.Sprintf("cpu_%s", time.Now().Format("20060102_150405")))
	defer profile.Start(profile.CPUProfile, profile.ProfilePath(profileDir)).Stop()

	log.Printf("プロファイル開始: vmd=%s pmx=%s frames=%d", filepath.Base(resolvedVmdPath), filepath.Base(resolvedPmxPath), *frameCount)

	// motionData, err := loadVmd(resolvedVmdPath)
	// if err != nil {
	// 	log.Fatalf("VMD読み込みに失敗しました: %v", err)
	// }
	// modelData, err := loadPmx(resolvedPmxPath)
	// if err != nil {
	// 	log.Fatalf("PMX読み込みに失敗しました: %v", err)
	// }

	// for i := 0; i <= *frameCount; i++ {
	// 	if i%100 == 0 {
	// 		log.Printf("frame: %d", i)
	// 	}
	// 	boneDeltas, _ := deform.ComputeBoneDeltas(modelData, motionData, motion.Frame(i), nil, *includeIk, false, *removeTwist)
	// 	deform.ApplyBoneMatrices(modelData, boneDeltas)
	// }
}

// resolveProfilePath はプロファイル用パスを決定する。
func resolveProfilePath(flagValue, envKey, fallback string) string {
	if flagValue != "" {
		return flagValue
	}
	if envValue := os.Getenv(envKey); envValue != "" {
		return envValue
	}
	return fallback
}
