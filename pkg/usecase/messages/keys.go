// 指示: miu200521358
// Package messages は usecase 層で返却するメッセージキーを提供する。
package messages

// メッセージキー一覧。
const (
	LoadModelRepositoryNotConfigured   = "モデル読み込みリポジトリがありません"
	LoadMotionRepositoryNotConfigured  = "モーション読み込みリポジトリがありません"
	LoadModelFormatNotSupported        = "モデル形式が不正です"
	LoadMotionFormatNotSupported       = "モーション形式が不正です"
	LoadModelOverrideBoneInsertWarning = "不足ボーン補完に失敗したため処理を継続します: %s"
	SaveModelNotLoaded                 = "XまたはPMDファイルが読み込まれていません"
	SavePathInvalid                    = "保存先パスが不正です"
	SaveRepositoryNotConfigured        = "保存リポジトリがありません"
	SavePathServiceNotConfigured       = "保存先判定ができません"
	TextureExistsValidationFailed      = "テクスチャの存在確認に失敗しました: %s"
	TextureImageValidationFailed       = "テクスチャの読込に失敗しました: %s"
)
