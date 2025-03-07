package rendering

import (
	"github.com/miu200521358/mlib_go/pkg/domain/pmx"
)

// ShowOptions は、ボーンの可視化設定をまとめた構造体。
// 「すべて表示」や「IKのみ表示」といったユーザー設定を想定しています。
type ShowOptions struct {
	ShowAll       bool
	ShowVisible   bool
	ShowIk        bool
	ShowEffector  bool
	ShowFixed     bool
	ShowRotate    bool
	ShowTranslate bool
}

// BoneColor は、Bone の描画色を決定するドメインサービス用インターフェース
type BoneColor interface {
	// GetBoneColor は、指定した Bone を描画する際の色を返します
	GetBoneColor(bone *pmx.Bone, opts ShowOptions) []float32
}

// boneColor は BoneColor の実装
type boneColor struct{}

// 各種ボーンのデバッグ表示用カラー。ここで定義する値は純粋な色情報であり、OpenGL呼び出しは行いません。
var (
	boneColorsIK             = []float32{1.0, 0.38, 0.0, 1.0}
	boneColorsIKLink         = []float32{1.0, 0.83, 0.49, 1.0}
	boneColorsIKTarget       = []float32{1.0, 0.57, 0.61, 1.0}
	boneColorsFixed          = []float32{0.72, 0.32, 1.0, 1.0}
	boneColorsEffect         = []float32{0.68, 0.64, 1.0, 1.0}
	boneColorsEffectEffector = []float32{0.88, 0.84, 1.0, 0.7}
	boneColorsTranslate      = []float32{0.70, 1.0, 0.54, 1.0}
	boneColorsRotate         = []float32{0.56, 0.78, 1.0, 1.0}
	boneColorsInvisible      = []float32{0.82, 0.82, 0.82, 1.0}
)

// NewBoneColor は新しい BoneColor を生成して返します。
func NewBoneColor() BoneColor {
	return &boneColor{}
}

// GetBoneColor はボーンの状態と表示オプションから適切な表示カラーを返すメソッドです。
// 元々の getBoneDebugColor() 相当のロジックをドメインサービスに切り出しています。
func (s *boneColor) GetBoneColor(bone *pmx.Bone, opts ShowOptions) []float32 {
	// IK
	if (opts.ShowAll || opts.ShowVisible || opts.ShowIk) && bone.IsIK() {
		return boneColorsIK
	} else if (opts.ShowAll || opts.ShowVisible || opts.ShowIk) &&
		len(bone.IkLinkBoneIndexes) > 0 {
		return boneColorsIKLink
	} else if (opts.ShowAll || opts.ShowVisible || opts.ShowIk) &&
		len(bone.IkTargetBoneIndexes) > 0 {
		return boneColorsIKTarget
	}

	// 付与（エフェクタ）
	if (opts.ShowAll || opts.ShowVisible || opts.ShowEffector) &&
		(bone.IsEffectorRotation() || bone.IsEffectorTranslation()) {
		return boneColorsEffect
	} else if (opts.ShowAll || opts.ShowVisible || opts.ShowEffector) &&
		len(bone.EffectiveBoneIndexes) > 0 {
		return boneColorsEffectEffector
	}

	// 軸固定
	if (opts.ShowAll || opts.ShowVisible || opts.ShowFixed) && bone.HasFixedAxis() {
		return boneColorsFixed
	}

	// 移動
	if (opts.ShowAll || opts.ShowVisible || opts.ShowTranslate) && bone.CanTranslate() {
		return boneColorsTranslate
	}

	// 回転
	if (opts.ShowAll || opts.ShowVisible || opts.ShowRotate) && bone.CanRotate() {
		return boneColorsRotate
	}

	// 非表示
	if opts.ShowAll && !bone.IsVisible() {
		return boneColorsInvisible
	}

	// 何にも該当しない場合は無色（アルファ0）
	return []float32{0.0, 0.0, 0.0, 0.0}
}
