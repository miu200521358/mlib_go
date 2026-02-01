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

// TestHashableAccessors はgetter/setterを確認する。
func TestHashableAccessors(t *testing.T) {
	hb := NewHashableBase("n", "p")
	if hb.Name() != "n" || hb.Path() != "p" {
		t.Errorf("initial values: got name=%v path=%v", hb.Name(), hb.Path())
	}
	hb.SetName("n2")
	hb.SetPath("p2")
	hb.SetFileModTime(10)
	hb.SetHash("h")
	if hb.Name() != "n2" || hb.Path() != "p2" || hb.FileModTime() != 10 || hb.Hash() != "h" {
		t.Errorf("updated values: got name=%v path=%v mod=%v hash=%v", hb.Name(), hb.Path(), hb.FileModTime(), hb.Hash())
	}
	if hb.GetHashParts() != "" {
		t.Errorf("GetHashParts default should be empty")
	}
}

// TestHashableUpdateHashNoParts は追加部品なしのハッシュを確認する。
func TestHashableUpdateHashNoParts(t *testing.T) {
	hb := NewHashableBase("name", "path")
	hb.SetFileModTime(456)
	hb.UpdateHash()

	fh := fnv.New32a()
	_, _ = fh.Write([]byte("name"))
	_, _ = fh.Write([]byte("path"))
	_, _ = fh.Write([]byte(strconv.FormatInt(456, 10)))
	expected := fmt.Sprintf("%x", fh.Sum(nil))

	if hb.Hash() != expected {
		t.Errorf("UpdateHash no parts: got=%s want=%s", hb.Hash(), expected)
	}
}
