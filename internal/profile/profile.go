// 指示: miu200521358
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	pmxio "github.com/miu200521358/mlib_go/pkg/adapter/io_model/pmx"
	vmdio "github.com/miu200521358/mlib_go/pkg/adapter/io_motion/vmd"
	"github.com/miu200521358/mlib_go/pkg/domain/deform"
	"github.com/miu200521358/mlib_go/pkg/domain/model"
	"github.com/miu200521358/mlib_go/pkg/domain/motion"
	"github.com/pkg/profile"
)

const (
	defaultVmdPath = "D:/MMD/MikuMikuDance_v926x64/UserFile/Motion/ダンス_1人/[A]ddiction_モーション hino/[A]ddiction_Lat式.vmd"
	defaultPmxPath = "D:/MMD/MikuMikuDance_v926x64/UserFile/Model/VOCALOID/初音ミク/ISAO式ミク/I_ミクv4チャイナ/Miku_V4_チャイナ.pmx"

	// defaultVmdPath = "D:/MMD/MikuMikuDance_v926x64/UserFile/Motion/ダンス_1人/CH4NGE mobiusP/CH4NGE.vmd"
	// defaultPmxPath = "D:/MMD/MikuMikuDance_v926x64/UserFile/Model/_VMDサイジング/mlibkiller/mlibkiller.pmx"
)

// main はCPUプロファイルを取得しながらボーン変形を実行する。
func main() {
	vmdPath := flag.String("vmd", "", "VMDパス")
	pmxPath := flag.String("pmx", "", "PMXパス")
	frameCount := flag.Int("frames", 500, "計測するフレーム数")
	includeIk := flag.Bool("ik", true, "IKを含める")
	removeTwist := flag.Bool("remove-twist", false, "捩り除去を行う")
	flag.Parse()

	resolvedVmdPath := convertWindowsPathToWsl(resolveProfilePath(*vmdPath, "MLIB_PROFILE_VMD", defaultVmdPath))
	resolvedPmxPath := convertWindowsPathToWsl(resolveProfilePath(*pmxPath, "MLIB_PROFILE_PMX", defaultPmxPath))

	profileDir := filepath.Join("internal", "profile", fmt.Sprintf("cpu_%s", time.Now().Format("20060102_150405")))
	defer profile.Start(profile.CPUProfile, profile.ProfilePath(profileDir)).Stop()

	log.Printf("プロファイル開始: vmd=%s pmx=%s frames=%d", filepath.Base(resolvedVmdPath), filepath.Base(resolvedPmxPath), *frameCount)

	motionData, err := loadVmd(resolvedVmdPath)
	if err != nil {
		log.Fatalf("VMD読み込みに失敗しました: %v", err)
	}
	modelData, err := loadPmx(resolvedPmxPath)
	if err != nil {
		log.Fatalf("PMX読み込みに失敗しました: %v", err)
	}

	for i := 0; i <= *frameCount; i++ {
		if i%100 == 0 {
			log.Printf("frame: %d", i)
		}
		boneDeltas, _ := deform.ComputeBoneDeltas(modelData, motionData, motion.Frame(i), nil, *includeIk, false, *removeTwist)
		deform.ApplyBoneMatrices(modelData, boneDeltas)
	}
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

// convertWindowsPathToWsl はLinux環境でWindowsパスをWSLパスに変換する。
func convertWindowsPathToWsl(path string) string {
	if runtime.GOOS != "linux" {
		return path
	}
	if len(path) < 2 || path[1] != ':' {
		return path
	}
	drive := strings.ToLower(path[:1])
	rest := strings.ReplaceAll(path[2:], "\\", "/")
	if rest == "" {
		return "/mnt/" + drive
	}
	if !strings.HasPrefix(rest, "/") {
		rest = "/" + rest
	}
	return "/mnt/" + drive + rest
}

// loadVmd はVMDを読み込んで返す。
func loadVmd(path string) (*motion.VmdMotion, error) {
	repo := vmdio.NewVmdRepository()
	data, err := repo.Load(path)
	if err != nil {
		return nil, err
	}
	motionData, ok := data.(*motion.VmdMotion)
	if !ok || motionData == nil {
		return nil, fmt.Errorf("VMD読み込み結果が不正です")
	}
	return motionData, nil
}

// loadPmx はPMXを読み込んで返す。
func loadPmx(path string) (*model.PmxModel, error) {
	repo := pmxio.NewPmxRepository()
	data, err := repo.Load(path)
	if err != nil {
		return nil, err
	}
	modelData, ok := data.(*model.PmxModel)
	if !ok || modelData == nil {
		return nil, fmt.Errorf("PMX読み込み結果が不正です")
	}
	return modelData, nil
}
