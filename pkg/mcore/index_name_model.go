package mcore

import (
	"sort"
	"sync"
)

type TInterface interface {
	IsValid() bool
	Copy() TInterface
	GetIndex() int
	SetIndex(index int)
	GetName() string
	SetName(name string)
}

type T struct {
	Index       int
	Name        string
	EnglishName string
}

func NewT(index int, name, englishName string) *T {
	return &T{
		Index:       index,
		Name:        name,
		EnglishName: englishName,
	}
}

func (t *T) IsValid() bool {
	return t.GetIndex() >= 0 && len(t.Name) >= 0
}

func (t *T) GetIndex() int {
	return t.Index
}

func (t *T) SetIndex(index int) {
	t.Index = index
}

func (t *T) GetName() string {
	return t.Name
}

func (t *T) SetName(name string) {
	t.Name = name
}

// Copy
func (t *T) Copy() TInterface {
	copied := *t
	return &copied
}

type C struct {
	Name    string
	Data    map[int]TInterface
	Indexes []int
	mu      sync.Mutex
}

func NewC(name string) *C {
	return &C{
		Name:    name,
		Indexes: make([]int, 0),
		Data:    make(map[int]TInterface),
	}
}

func (c *C) GetByName(name string) TInterface {
	names := c.Names()
	index, ok := names[name]
	if !ok {
		return nil
	}
	return c.GetByIndex(index)
}

func (c *C) GetByIndex(index int) TInterface {
	c.mu.Lock()
	defer c.mu.Unlock()

	if index < 0 {
		index = len(c.Data) + index
	}

	return c.Data[index]
}

func (c *C) Append(value TInterface, isSort, isPositiveIndex bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if value.GetIndex() < 0 && isPositiveIndex {
		value.SetIndex(len(c.Data))
	}

	names := c.Names()

	if value.GetName() != "" && names[value.GetName()] == 0 {
		names[value.GetName()] = value.GetIndex()
	}

	c.Data[value.GetIndex()] = value.(*T)
	if isSort {
		c.SortIndexes()
	} else {
		c.Indexes = append(c.Indexes, value.GetIndex())
	}
}

func (c *C) Remove(value T, isSort bool) map[int]int {
	c.mu.Lock()
	defer c.mu.Unlock()

	replacedMap := make(map[int]int)

	if _, ok := c.Data[value.GetIndex()]; !ok {
		return replacedMap
	}

	delete(c.Data, value.GetIndex())
	for i := range c.Indexes {
		if i == value.GetIndex() {
			c.Indexes = append(c.Indexes[:i], c.Indexes[i+1:]...)
			break
		}
	}

	names := c.Names()
	for name, idx := range names {
		if idx == value.GetIndex() {
			delete(names, name)
			break
		}
	}

	replacedMap[value.GetIndex()] = value.GetIndex() - 1

	for i := value.GetIndex() + 1; i <= c.LastIndex(); i++ {
		v := c.Data[i]
		replacedMap[v.GetIndex()] = v.GetIndex() - 1
		index := v.GetIndex() - 1
		v.SetIndex(index)
		c.Data[v.GetIndex()] = v
		names[v.GetName()] = v.GetIndex()
	}
	for i := value.GetIndex() - 1; i >= 0; i-- {
		v := c.Data[i]
		replacedMap[v.GetIndex()] = v.GetIndex()
	}

	delete(c.Data, c.LastIndex())

	if isSort {
		c.SortIndexes()
	}

	if len(replacedMap) > 0 {
		replacedMap[-1] = -1
	}

	return replacedMap
}

func (c *C) Insert(value T, isSort, isPositiveIndex bool) map[int]int {
	c.mu.Lock()
	defer c.mu.Unlock()

	if value.GetIndex() < 0 && isPositiveIndex {
		value.SetIndex(len(c.Data))
	}

	names := c.Names()
	replacedMap := make(map[int]int)
	if _, ok := c.Data[value.GetIndex()]; ok {
		for i := c.LastIndex(); i >= value.GetIndex(); i-- {
			v := c.Data[i]
			replacedMap[v.GetIndex()] = v.GetIndex() + 1
			index := v.GetIndex() + 1
			v.SetIndex(index)
			c.Data[v.GetIndex()] = v
			names[v.GetName()] = v.GetIndex()
		}
		for i := value.GetIndex() - 1; i >= 0; i-- {
			v := c.Data[i]
			replacedMap[v.GetIndex()] = v.GetIndex()
		}
	}

	c.Data[value.GetIndex()] = &value
	if value.GetName() != "" && names[value.GetName()] == 0 {
		names[value.GetName()] = value.GetIndex()
	}

	if isSort {
		c.SortIndexes()
	} else {
		c.Indexes = append(c.Indexes, value.GetIndex())
	}

	if len(replacedMap) > 0 {
		replacedMap[-1] = -1
	}

	return replacedMap
}

func (c *C) SortIndexes() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.Indexes = make([]int, 0, len(c.Data))
	for index := range c.Data {
		c.Indexes = append(c.Indexes, index)
	}
	sort.Ints(c.Indexes)
}

func (c *C) LastIndex() int {
	c.mu.Lock()
	defer c.mu.Unlock()

	if len(c.Data) == 0 {
		return 0
	}
	maxIndex := 0
	for index := range c.Data {
		if index > maxIndex {
			maxIndex = index
		}
	}
	return maxIndex
}

func (c *C) Range(start, stop, step int) []TInterface {
	c.mu.Lock()
	defer c.mu.Unlock()

	if stop < 0 {
		stop = len(c.Data) + stop + 1
	}

	result := make([]TInterface, 0, (stop-start)/step)
	for i := start; i < stop; i += step {
		result = append(result, c.Data[c.Indexes[i]])
	}
	return result
}

func (c *C) Len() int {
	c.mu.Lock()
	defer c.mu.Unlock()

	return len(c.Data)
}

func (c *C) Names() map[string]int {
	c.mu.Lock()
	defer c.mu.Unlock()

	names := make(map[string]int, len(c.Data))
	for k, v := range c.Data {
		names[v.GetName()] = k
	}
	return names
}

func (c *C) RangeIndexes(index int, offFlag bool, indexes []int) (int, int, int) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if indexes == nil {
		indexes = c.Indexes
	}
	if !offFlag && (len(indexes) == 0 || c.Data[index] != nil) {
		return index, index, index
	}

	idx := sort.SearchInts(indexes, index)

	prevIndex := 0
	if idx != 0 {
		prevIndex = indexes[idx-1]
	}

	nextIndex := index
	if idx != len(indexes) {
		nextIndex = indexes[idx]
	}

	return prevIndex, index, nextIndex
}
