// 指示: miu200521358
package io_audio

import (
	"path/filepath"
	"strings"

	"github.com/miu200521358/mlib_go/pkg/adapter/io_common"
	"github.com/miu200521358/mlib_go/pkg/shared/base/i18n"
	"github.com/miu200521358/mlib_go/pkg/shared/hashable"
)

var audioExtensions = map[string]struct{}{
	".wav": {},
	".mp3": {},
}

// AudioRepository は音声ファイル読み込みのI/Fを表す。
type AudioRepository struct {
	translator i18n.II18n
}

// NewAudioRepository はAudioRepositoryを生成する。
func NewAudioRepository(translator i18n.II18n) *AudioRepository {
	return &AudioRepository{translator: translator}
}

// CanLoad は読み込み可能か判定する。
func (r *AudioRepository) CanLoad(path string) bool {
	ext := strings.ToLower(filepath.Ext(path))
	_, ok := audioExtensions[ext]
	return ok
}

// Load は音声ファイルを読み込む。
func (r *AudioRepository) Load(path string) (hashable.IHashable, error) {
	return nil, io_common.NewIoFormatNotSupported(r.t("音楽ファイルの読み込みは未実装です"), nil)
}

// InferName はパスから表示名を推定する。
func (r *AudioRepository) InferName(path string) string {
	base := filepath.Base(path)
	ext := filepath.Ext(base)
	return strings.TrimSuffix(base, ext)
}

// t は翻訳済み文言を返す。
func (r *AudioRepository) t(key string) string {
	if r == nil || r.translator == nil || !r.translator.IsReady() {
		return "●●" + key + "●●"
	}
	return r.translator.T(key)
}
