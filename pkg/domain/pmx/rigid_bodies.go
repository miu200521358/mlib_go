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
