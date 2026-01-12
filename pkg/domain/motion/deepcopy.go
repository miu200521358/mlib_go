// 指示: miu200521358
package motion

import (
	baseerr "github.com/miu200521358/mlib_go/pkg/shared/base/err"
	"github.com/tiendc/go-deepcopy"
)

const deepcopyErrorMessage = "go-deepcopyパッケージでエラーが発生しました"

// deepCopy はgo-deepcopyを使って複製する。
func deepCopy[T any](src T) (T, error) {
	var dst T
	if err := deepcopy.Copy(&dst, src); err != nil {
		return dst, baseerr.NewDeepcopyPackageError(deepcopyErrorMessage, err)
	}
	return dst, nil
}
