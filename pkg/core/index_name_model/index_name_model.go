package index_name_model

import (
	"sort"
	"sync"
)

type TInterface interface {
	IsValid() bool
	Add(v *TInterface) *T
	Copy() *T
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

func (m *T) IsValid() bool {
	return m.GetIndex() >= 0 && len(m.Name) >= 0
}

func (m *T) Add(v *TInterface) *T {
	// Implement your logic here
	return nil
}

func (m *T) GetIndex() int {
	return m.Index
}

func (m *T) SetIndex(index int) {
	m.Index = index
}

func (b *T) GetName() string {
	return b.Name
}

func (b *T) SetName(name string) {
	b.Name = name
}

// Copy
func (b *T) Copy() *T {
	copied := *b
	return &copied
}

// Cのインタフェース
type CInterface interface {
	GetItem(index int) *T
	Range(start, stop, step int) []*T
	SetItem(index int, v TInterface)
	Append(value TInterface, isSort bool)
}

type C struct {
	Name    string
	data    map[int]*T
	Indexes []int
	names   map[string]int
	mu      sync.Mutex
}

func NewC(name string) *C {
	return &C{
		Name:    name,
		Indexes: make([]int, 0),
		data:    make(map[int]*T),
		names:   make(map[string]int),
	}
}

func (m *C) GetByName(name string) *T {
	index, ok := m.names[name]
	if !ok {
		return nil
	}
	return m.GetByIndex(index)
}

func (m *C) GetByIndex(index int) *T {
	m.mu.Lock()
	defer m.mu.Unlock()

	if index < 0 {
		index = len(m.data) + index
	}

	return m.data[index]
}

func (m *C) Append(value TInterface, isSort, isPositiveIndex bool) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if value.GetIndex() < 0 && isPositiveIndex {
		value.SetIndex(len(m.data))
	}

	if value.GetName() != "" && m.names[value.GetName()] == 0 {
		m.names[value.GetName()] = value.GetIndex()
	}

	m.data[value.GetIndex()] = value.(*T) // Perform type assertion to convert value to *T
	if isSort {
		m.SortIndexes(true)
	} else {
		m.Indexes = append(m.Indexes, value.GetIndex())
	}
}

func (m *C) Remove(value TInterface, isSort bool) map[int]int {
	m.mu.Lock()
	defer m.mu.Unlock()

	replacedMap := make(map[int]int)

	if _, ok := m.data[value.GetIndex()]; !ok {
		return replacedMap
	}

	delete(m.data, value.GetIndex())
	for i := range m.Indexes {
		if i == value.GetIndex() {
			m.Indexes = append(m.Indexes[:i], m.Indexes[i+1:]...)
			break
		}
	}
	for name, idx := range m.names {
		if idx == value.GetIndex() {
			delete(m.names, name)
			break
		}
	}

	replacedMap[value.GetIndex()] = value.GetIndex() - 1

	for i := value.GetIndex() + 1; i <= m.lastIndex(); i++ {
		v := m.data[i]
		replacedMap[v.GetIndex()] = v.GetIndex() - 1
		index := v.GetIndex() - 1
		v.SetIndex(index)
		m.data[v.GetIndex()] = v
		m.names[v.Name] = v.GetIndex()
	}
	for i := value.GetIndex() - 1; i >= 0; i-- {
		v := m.data[i]
		replacedMap[v.GetIndex()] = v.GetIndex()
	}

	delete(m.data, m.lastIndex())

	if isSort {
		m.SortIndexes(true)
	}

	if len(replacedMap) > 0 {
		replacedMap[-1] = -1
	}

	return replacedMap
}

func (m *C) Insert(value TInterface, isSort, isPositiveIndex bool) map[int]int {
	m.mu.Lock()
	defer m.mu.Unlock()

	if value.GetIndex() < 0 && isPositiveIndex {
		value.SetIndex(len(m.data))
	}

	replacedMap := make(map[int]int)
	if _, ok := m.data[value.GetIndex()]; ok {
		for i := m.lastIndex(); i >= value.GetIndex(); i-- {
			v := m.data[i]
			replacedMap[v.GetIndex()] = v.GetIndex() + 1
			index := v.GetIndex() + 1
			v.SetIndex(index)
			m.data[v.GetIndex()] = v
			m.names[v.Name] = v.GetIndex()
		}
		for i := value.GetIndex() - 1; i >= 0; i-- {
			v := m.data[i]
			replacedMap[v.GetIndex()] = v.GetIndex()
		}
	}
	if v, ok := value.(*T); ok {
		m.data[v.GetIndex()] = v
	}
	if value.GetName() != "" && m.names[value.GetName()] == 0 {
		m.names[value.GetName()] = value.GetIndex()
	}

	if isSort {
		m.SortIndexes(false)
	} else {
		m.Indexes = append(m.Indexes, value.GetIndex())
	}

	if len(replacedMap) > 0 {
		replacedMap[-1] = -1
	}

	return replacedMap
}

func (m *C) SortIndexes(isSortName bool) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.Indexes = make([]int, 0, len(m.data))
	for index := range m.data {
		m.Indexes = append(m.Indexes, index)
	}
	sort.Ints(m.Indexes)

	if isSortName {
		m.names = make(map[string]int)
		for _, index := range m.Indexes {
			m.names[m.data[index].Name] = index
		}
	}
}

func (m *C) LastIndex() int {
	m.mu.Lock()
	defer m.mu.Unlock()

	return m.lastIndex()
}

func (m *C) lastIndex() int {
	if len(m.data) == 0 {
		return 0
	}
	maxIndex := 0
	for index := range m.data {
		if index > maxIndex {
			maxIndex = index
		}
	}
	return maxIndex
}

func (m *C) Range(start, stop, step int) []*T {
	m.mu.Lock()
	defer m.mu.Unlock()

	if stop < 0 {
		stop = len(m.data) + stop + 1
	}

	result := make([]*T, 0, (stop-start)/step)
	for i := start; i < stop; i += step {
		result = append(result, m.data[m.Indexes[i]])
	}
	return result
}

func (m *C) Len() int {
	m.mu.Lock()
	defer m.mu.Unlock()

	return len(m.data)
}

func (m *C) Names() []string {
	m.mu.Lock()
	defer m.mu.Unlock()

	names := make([]string, 0, len(m.names))
	for name := range m.names {
		names = append(names, name)
	}
	return names
}

func (m *C) RangeIndexes(index int, offFlg bool, indexes []int) (int, int, int) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if indexes == nil {
		indexes = m.Indexes
	}
	if !offFlg && (len(indexes) == 0 || m.data[index] != nil) {
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
