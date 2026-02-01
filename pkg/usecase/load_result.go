// 指示: miu200521358
package usecase

import (
	"github.com/miu200521358/mlib_go/pkg/domain/model"
	"github.com/miu200521358/mlib_go/pkg/domain/motion"
	portio "github.com/miu200521358/mlib_go/pkg/usecase/port/io"
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
func LoadMotionWithMeta(rep portio.IFileReader, path string) (*MotionLoadResult, error) {
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

// CanLoadPath はリポジトリが指定パスを読み込み可能か判定する。
func CanLoadPath(rep portio.IFileReader, path string) bool {
	if rep == nil || path == "" {
		return false
	}
	return rep.CanLoad(path)
}
