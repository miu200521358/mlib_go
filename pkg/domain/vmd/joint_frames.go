package vmd

import (
	"math"
	"slices"
	"sync/atomic"

	"github.com/miu200521358/mlib_go/pkg/domain/mmath"
	"github.com/tiendc/go-deepcopy"
)

type JointFrames struct {
	names  atomic.Value
	values []*JointNameFrames
}

func NewJointFrames() *JointFrames {
	return &JointFrames{
		names:  atomic.Value{},
		values: make([]*JointNameFrames, 0),
	}
}

func (jointFrames *JointFrames) getNames() []string {
	names := jointFrames.names.Load()
	if names == nil {
		return make([]string, 0)
	}
	return names.([]string)
}

func (jointFrames *JointFrames) setNames(names []string) {
	jointFrames.names.Store(names)
}

func (jointFrames *JointFrames) Contains(jointName string) bool {
	if slices.Contains(jointFrames.getNames(), jointName) {
		return true
	}

	return false
}

func (jointFrames *JointFrames) Update(jointNameFrames *JointNameFrames) {
	index := slices.Index(jointFrames.getNames(), jointNameFrames.Name)
	if index < 0 {
		jointFrames.setNames(append(jointFrames.getNames(), jointNameFrames.Name))
		jointFrames.values = append(jointFrames.values, jointNameFrames)
	} else {
		jointFrames.getNames()[index] = jointNameFrames.Name
		jointFrames.values[index] = jointNameFrames
	}
}

func (jointFrames *JointFrames) Delete(jointName string) {
	index := slices.Index(jointFrames.getNames(), jointName)
	if index < 0 {
		return
	}
	jointFrames.setNames(append(jointFrames.getNames()[:index], jointFrames.getNames()[index+1:]...))
	jointFrames.values = append(jointFrames.values[:index], jointFrames.values[index+1:]...)
}

func (jointFrames *JointFrames) Get(jointName string) *JointNameFrames {
	index := slices.Index(jointFrames.getNames(), jointName)

	if index < 0 || index >= len(jointFrames.values) {
		jointNameFrames := NewJointNameFrames(jointName)
		jointFrames.setNames(append(jointFrames.getNames(), jointNameFrames.Name))
		jointFrames.values = append(jointFrames.values, jointNameFrames)
		return jointNameFrames
	}

	return jointFrames.values[index]
}

func (jointFrames *JointFrames) Names() []string {
	return jointFrames.getNames()
}

func (jointFrames *JointFrames) Indexes() []int {
	indexes := make([]int, 0)
	for _, jointFrames := range jointFrames.values {
		jointFrames.values.ForEach(func(index float32) bool {
			indexes = append(indexes, int(index))
			return true
		})
	}
	mmath.Unique(indexes)
	mmath.Sort(indexes)
	return indexes
}

func (jointFrames *JointFrames) IndexesByNames(names []string) []int {
	indexes := make([]int, 0)
	for _, jointName := range jointFrames.getNames() {
		if !slices.Contains(names, jointName) {
			continue
		}
		jointFrames.Get(jointName).values.ForEach(func(index float32) bool {
			indexes = append(indexes, int(index))
			return true
		})
	}
	indexes = mmath.Unique(indexes)
	mmath.Sort(indexes)
	return indexes
}

func (jointFrames *JointFrames) Length() int {
	count := 0
	for _, jointFrames := range jointFrames.values {
		count += jointFrames.values.Len()
	}
	return count
}

func (jointFrames *JointFrames) MaxFrame() float32 {
	maxFno := float32(0)
	for _, jointFrames := range jointFrames.values {
		fno := float32(jointFrames.MaxFrame())
		if fno > maxFno {
			maxFno = fno
		}
	}
	return maxFno
}

func (jointFrames *JointFrames) MinFrame() float32 {
	minFno := float32(math.MaxFloat32)
	for _, jointFrames := range jointFrames.values {
		fno := float32(jointFrames.MinFrame())
		if fno < minFno {
			minFno = fno
		}
	}
	if minFno == float32(math.MaxFloat32) {
		return 0
	}
	return minFno
}

func (jointFrames *JointFrames) ForEach(fn func(jointName string, jointNameFrames *JointNameFrames)) {
	for _, jointName := range jointFrames.getNames() {
		fn(jointName, jointFrames.Get(jointName))
	}
}

func (jointFrames *JointFrames) Copy() (*JointFrames, error) {
	copied := new(JointFrames)
	err := deepcopy.Copy(copied, jointFrames)
	return copied, err
}
