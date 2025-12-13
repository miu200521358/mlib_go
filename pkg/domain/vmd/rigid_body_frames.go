package vmd

import (
	"math"
	"slices"
	"sync/atomic"

	"github.com/miu200521358/mlib_go/pkg/domain/mmath"
	"github.com/tiendc/go-deepcopy"
)

type RigidBodyFrames struct {
	names  atomic.Value
	values []*RigidBodyNameFrames
}

func NewRigidBodyFrames() *RigidBodyFrames {
	return &RigidBodyFrames{
		names:  atomic.Value{},
		values: make([]*RigidBodyNameFrames, 0),
	}
}

func (rigidBodyFrames *RigidBodyFrames) getNames() []string {
	names := rigidBodyFrames.names.Load()
	if names == nil {
		return make([]string, 0)
	}
	return names.([]string)
}

func (rigidBodyFrames *RigidBodyFrames) setNames(names []string) {
	rigidBodyFrames.names.Store(names)
}

func (rigidBodyFrames *RigidBodyFrames) Contains(rigidBodyName string) bool {
	if slices.Contains(rigidBodyFrames.getNames(), rigidBodyName) {
		return true
	}

	return false
}

func (rigidBodyFrames *RigidBodyFrames) Update(rigidBodyNameFrames *RigidBodyNameFrames) {
	index := slices.Index(rigidBodyFrames.getNames(), rigidBodyNameFrames.Name)
	if index < 0 {
		rigidBodyFrames.setNames(append(rigidBodyFrames.getNames(), rigidBodyNameFrames.Name))
		rigidBodyFrames.values = append(rigidBodyFrames.values, rigidBodyNameFrames)
	} else {
		rigidBodyFrames.getNames()[index] = rigidBodyNameFrames.Name
		rigidBodyFrames.values[index] = rigidBodyNameFrames
	}
}

func (rigidBodyFrames *RigidBodyFrames) Delete(rigidBodyName string) {
	index := slices.Index(rigidBodyFrames.getNames(), rigidBodyName)
	if index < 0 {
		return
	}
	rigidBodyFrames.setNames(append(rigidBodyFrames.getNames()[:index], rigidBodyFrames.getNames()[index+1:]...))
	rigidBodyFrames.values = append(rigidBodyFrames.values[:index], rigidBodyFrames.values[index+1:]...)
}

func (rigidBodyFrames *RigidBodyFrames) Get(rigidBodyName string) *RigidBodyNameFrames {
	index := slices.Index(rigidBodyFrames.getNames(), rigidBodyName)

	if index < 0 || index >= len(rigidBodyFrames.values) {
		rigidBodyNameFrames := NewRigidBodyNameFrames(rigidBodyName)
		rigidBodyFrames.setNames(append(rigidBodyFrames.getNames(), rigidBodyNameFrames.Name))
		rigidBodyFrames.values = append(rigidBodyFrames.values, rigidBodyNameFrames)
		return rigidBodyNameFrames
	}

	return rigidBodyFrames.values[index]
}

func (rigidBodyFrames *RigidBodyFrames) Names() []string {
	return rigidBodyFrames.getNames()
}

func (rigidBodyFrames *RigidBodyFrames) Indexes() []int {
	indexes := make([]int, 0)
	for _, rigidBodyFrames := range rigidBodyFrames.values {
		rigidBodyFrames.values.ForEach(func(index float32) bool {
			indexes = append(indexes, int(index))
			return true
		})
	}
	mmath.Unique(indexes)
	mmath.Sort(indexes)
	return indexes
}

func (rigidBodyFrames *RigidBodyFrames) IndexesByNames(names []string) []int {
	indexes := make([]int, 0)
	for _, rigidBodyName := range rigidBodyFrames.getNames() {
		if !slices.Contains(names, rigidBodyName) {
			continue
		}
		rigidBodyFrames.Get(rigidBodyName).values.ForEach(func(index float32) bool {
			indexes = append(indexes, int(index))
			return true
		})
	}
	indexes = mmath.Unique(indexes)
	mmath.Sort(indexes)
	return indexes
}

func (rigidBodyFrames *RigidBodyFrames) Length() int {
	count := 0
	for _, rigidBodyFrames := range rigidBodyFrames.values {
		count += rigidBodyFrames.values.Len()
	}
	return count
}

func (rigidBodyFrames *RigidBodyFrames) MaxFrame() float32 {
	maxFno := float32(0)
	for _, rigidBodyFrames := range rigidBodyFrames.values {
		fno := float32(rigidBodyFrames.MaxFrame())
		if fno > maxFno {
			maxFno = fno
		}
	}
	return maxFno
}

func (rigidBodyFrames *RigidBodyFrames) MinFrame() float32 {
	minFno := float32(math.MaxFloat32)
	for _, rigidBodyFrames := range rigidBodyFrames.values {
		fno := float32(rigidBodyFrames.MinFrame())
		if fno < minFno {
			minFno = fno
		}
	}
	if minFno == float32(math.MaxFloat32) {
		return 0
	}
	return minFno
}

func (rigidBodyFrames *RigidBodyFrames) ForEach(fn func(rigidBodyName string, rigidBodyNameFrames *RigidBodyNameFrames)) {
	for _, rigidBodyName := range rigidBodyFrames.getNames() {
		fn(rigidBodyName, rigidBodyFrames.Get(rigidBodyName))
	}
}

func (rigidBodyFrames *RigidBodyFrames) Copy() (*RigidBodyFrames, error) {
	copied := new(RigidBodyFrames)
	err := deepcopy.Copy(copied, rigidBodyFrames)
	return copied, err
}
