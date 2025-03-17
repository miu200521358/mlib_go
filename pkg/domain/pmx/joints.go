package pmx

import (
	"github.com/miu200521358/mlib_go/pkg/domain/core"
)

// ジョイントリスト
type Joints struct {
	*core.IndexNameModels[*Joint]
}

func NewJoints(capacity int) *Joints {
	return &Joints{
		IndexNameModels: core.NewIndexNameModels[*Joint](capacity),
	}
}
