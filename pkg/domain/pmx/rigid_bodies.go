package pmx

import (
	"github.com/miu200521358/mlib_go/pkg/domain/core"
)

// 剛体リスト
type RigidBodies struct {
	*core.IndexNameModels[*RigidBody]
}

func NewRigidBodies(capacity int) *RigidBodies {
	return &RigidBodies{
		IndexNameModels: core.NewIndexNameModels[*RigidBody](capacity),
	}
}

func (rigidBodies *RigidBodies) Setup(bones *Bones) {
	// ボーンごとの剛体のリストを初期化
	bones.ForEach(func(index int, bone *Bone) bool {
		bone.RigidBodies = make([]*RigidBody, 0)
		return true
	})

	// 剛体
	rigidBodies.ForEach(func(index int, rb *RigidBody) bool {
		if rb.BoneIndex >= 0 && bones.Contains(rb.BoneIndex) {
			// 剛体に関連付けられたボーンが存在する場合、剛体とボーンを関連付ける
			if bone, err := bones.Get(rb.BoneIndex); err == nil {
				bone.RigidBodies = append(bone.RigidBodies, rb)
				rb.Bone = bone
			}
		}
		return true
	})
}
