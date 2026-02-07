// 指示: miu200521358
package usecase

import (
	"github.com/miu200521358/mlib_go/pkg/domain/model"
	"github.com/miu200521358/mlib_go/pkg/domain/motion"
	"github.com/miu200521358/mlib_go/pkg/usecase/port/io"
)

// ModelLoadResult はモデル読み込み結果を表す。
// テクスチャ検証の実行順序はツール固有のため、Validation は呼び出し側で設定する。
type ModelLoadResult struct {
	Model      *model.PmxModel
	Validation *TextureValidationResult
}

// MotionLoadResult はモーション読み込み結果を表す。
type MotionLoadResult struct {
	Motion   *motion.VmdMotion
	MaxFrame motion.Frame
}

// LoadMotionWithMeta はモーションを読み込み、最大フレームを返す。
func LoadMotionWithMeta(rep io.IFileReader, path string) (*MotionLoadResult, error) {
	motionData, err := LoadMotion(rep, path)
	if err != nil {
		return nil, err
	}
	result := &MotionLoadResult{Motion: motionData}
	if motionData == nil {
		return result, nil
	}
	result.MaxFrame = motionData.MaxFrame()
	return result, nil
}

// LoadCameraMotionWithMeta はカメラフレームのみを採用してモーションを読み込み、最大フレームを返す。
func LoadCameraMotionWithMeta(rep io.IFileReader, path string) (*MotionLoadResult, error) {
	motionData, err := LoadMotion(rep, path)
	if err != nil {
		return nil, err
	}
	cameraMotion, err := extractCameraMotion(motionData)
	if err != nil {
		return nil, err
	}
	result := &MotionLoadResult{Motion: cameraMotion}
	if cameraMotion == nil {
		return result, nil
	}
	result.MaxFrame = cameraMotion.MaxFrame()
	return result, nil
}

// extractCameraMotion は入力モーションからカメラフレームのみを複製して返す。
func extractCameraMotion(source *motion.VmdMotion) (*motion.VmdMotion, error) {
	if source == nil {
		return nil, nil
	}
	cameraMotion := motion.NewVmdMotion(source.Path())
	cameraMotion.SetName(source.Name())
	cameraMotion.SetFileModTime(source.FileModTime())
	cameraMotion.Signature = source.Signature

	if source.CameraFrames != nil {
		cameraFrames, err := source.CameraFrames.Copy()
		if err != nil {
			return nil, err
		}
		cameraMotion.CameraFrames = &cameraFrames
	}

	cameraMotion.UpdateHash()
	return cameraMotion, nil
}

// CanLoadPath はリポジトリが指定パスを読み込み可能か判定する。
func CanLoadPath(rep io.IFileReader, path string) bool {
	if rep == nil || path == "" {
		return false
	}
	return rep.CanLoad(path)
}
