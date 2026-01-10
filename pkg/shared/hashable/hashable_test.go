// 指示: miu200521358
package hashable

import (
	"fmt"
	"hash/fnv"
	"strconv"
	"testing"
)

type testHashable struct {
	*HashableBase
	parts string
}

// GetHashParts はテスト用のハッシュ部品を返す。
func (t *testHashable) GetHashParts() string {
	return t.parts
}

// TestHashableUpdateHash は更新ロジックを確認する。
func TestHashableUpdateHash(t *testing.T) {
	h := &testHashable{HashableBase: NewHashableBase("name", "path"), parts: "00000002"}
	h.SetFileModTime(123)
	h.SetHashPartsFunc(h.GetHashParts)

	h.UpdateHash()

	fh := fnv.New32a()
	_, _ = fh.Write([]byte("name"))
	_, _ = fh.Write([]byte("path"))
	_, _ = fh.Write([]byte(strconv.FormatInt(123, 10)))
	_, _ = fh.Write([]byte("00000002"))
	expected := fmt.Sprintf("%x", fh.Sum(nil))

	if h.Hash() != expected {
		t.Errorf("UpdateHash: got=%s want=%s", h.Hash(), expected)
	}
}

// TestHashableUpdateRandomHash はランダム更新を確認する。
func TestHashableUpdateRandomHash(t *testing.T) {
	h := &testHashable{HashableBase: NewHashableBase("name", "path"), parts: ""}
	h.SetHashPartsFunc(h.GetHashParts)

	h.UpdateRandomHash()
	if h.Hash() == "" {
		t.Errorf("UpdateRandomHash: empty hash")
	}
}
