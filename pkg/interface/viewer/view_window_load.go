package viewer

import (
	"github.com/miu200521358/mlib_go/pkg/domain/delta"
	"github.com/miu200521358/mlib_go/pkg/domain/state"
	"github.com/miu200521358/mlib_go/pkg/infrastructure/render"
)

func (vw *ViewWindow) loadModelRenderers(shared *state.SharedState) {
	for i := range shared.ModelCount(vw.windowIndex) {
		for i >= len(vw.modelRenderers) {
			vw.modelRenderers = append(vw.modelRenderers, nil)
		}
		model := shared.LoadModel(vw.windowIndex, i)
		if vw.modelRenderers[i] != nil && vw.modelRenderers[i].Hash() != model.Hash() {
			vw.physics.DeleteModel(i)
			vw.modelRenderers[i].Delete()
			vw.modelRenderers[i] = nil
		}
		if vw.modelRenderers[i] == nil {
			vw.modelRenderers[i] = render.NewModelRenderer(vw.windowIndex, model)
			vw.physics.AddModel(i, model)
		}
	}

	// モデルが新しく読み込まれた場合、vmdDeltasの初期化が必要
	if len(vw.vmdDeltas) != len(vw.modelRenderers) {
		oldDeltas := vw.vmdDeltas
		vw.vmdDeltas = make([]*delta.VmdDeltas, len(vw.modelRenderers))
		// 既存のデータを保持
		for i := 0; i < len(oldDeltas) && i < len(vw.vmdDeltas); i++ {
			vw.vmdDeltas[i] = oldDeltas[i]
		}

		// 初期フレームで事前計算が必要
		vw.initializeFrameCalculation(shared.Frame())
	}
}

func (vw *ViewWindow) loadMotions(shared *state.SharedState) {
	motionChanged := false

	for i := range shared.MotionCount(vw.windowIndex) {
		for i >= len(vw.motions) {
			vw.motions = append(vw.motions, nil)
			vw.vmdDeltas = append(vw.vmdDeltas, nil)
			motionChanged = true
		}

		motion := shared.LoadMotion(vw.windowIndex, i)

		if vw.motions[i] != nil && vw.motions[i].Hash() != motion.Hash() {
			// モーションのハッシュが変わった場合
			vw.motions[i] = nil
			vw.vmdDeltas[i] = nil // デルタも初期化
			motionChanged = true
		}

		if vw.motions[i] == nil {
			// 新しいモーションが設定された場合
			vw.motions[i] = motion
			vw.vmdDeltas[i] = nil // デルタも初期化
			motionChanged = true
		}
	}

	// モーションが変更された場合、初期計算が必要
	if motionChanged && len(vw.motions) > 0 && len(vw.modelRenderers) > 0 {
		vw.initializeFrameCalculation(shared.Frame())
	}
}
