// 指示: miu200521358
package hashable

import (
	"fmt"
	"hash/fnv"
	"math/rand"
	"strconv"
	"time"
)

// IHashable はI/O共通のハッシュ契約。
type IHashable interface {
	Name() string
	SetName(name string)
	Path() string
	SetPath(path string)
	FileModTime() int64
	SetFileModTime(modTime int64)
	Hash() string
	SetHash(hash string)
	GetHashParts() string
	UpdateHash()
	UpdateRandomHash()
}

// HashableBase はIHashableの基本実装。
type HashableBase struct {
	name        string
	path        string
	fileModTime int64
	hash        string
	partsFunc   func() string
}

// init は乱数シードを初期化する。
func init() {
	rand.Seed(time.Now().UnixNano())
}

// NewHashableBase はHashableBaseを生成する。
func NewHashableBase(name, path string) *HashableBase {
	return &HashableBase{
		name: name,
		path: path,
	}
}

// SetHashPartsFunc はGetHashPartsの代替関数を設定する。
func (hb *HashableBase) SetHashPartsFunc(fn func() string) {
	hb.partsFunc = fn
}

// Name は名称を返す。
func (hb *HashableBase) Name() string {
	return hb.name
}

// SetName は名称を設定する。
func (hb *HashableBase) SetName(name string) {
	hb.name = name
}

// Path はパスを返す。
func (hb *HashableBase) Path() string {
	return hb.path
}

// SetPath はパスを設定する。
func (hb *HashableBase) SetPath(path string) {
	hb.path = path
}

// FileModTime は更新時刻を返す。
func (hb *HashableBase) FileModTime() int64 {
	return hb.fileModTime
}

// SetFileModTime は更新時刻を設定する。
func (hb *HashableBase) SetFileModTime(modTime int64) {
	hb.fileModTime = modTime
}

// Hash はハッシュ文字列を返す。
func (hb *HashableBase) Hash() string {
	return hb.hash
}

// SetHash はハッシュ文字列を設定する。
func (hb *HashableBase) SetHash(hash string) {
	hb.hash = hash
}

// GetHashParts はハッシュ用の追加要素を返す。
func (hb *HashableBase) GetHashParts() string {
	if hb.partsFunc == nil {
		return ""
	}
	return hb.partsFunc()
}

// UpdateHash はName/Path/FileModTime/GetHashPartsを連結して更新する。
func (hb *HashableBase) UpdateHash() {
	h := fnv.New32a()
	_, _ = h.Write([]byte(hb.name))
	_, _ = h.Write([]byte(hb.path))
	_, _ = h.Write([]byte(strconv.FormatInt(hb.fileModTime, 10)))
	_, _ = h.Write([]byte(hb.GetHashParts()))
	hb.hash = fmt.Sprintf("%x", h.Sum(nil))
}

// UpdateRandomHash は再読込用にランダム値を設定する。
func (hb *HashableBase) UpdateRandomHash() {
	hb.hash = strconv.FormatInt(rand.Int63(), 10)
}
