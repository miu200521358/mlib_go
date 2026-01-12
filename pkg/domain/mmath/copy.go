// 指示: miu200521358
package mmath

import "github.com/tiendc/go-deepcopy"

// deepCopy は汎用の深いコピーを行う。
func deepCopy[T any](src T) (T, error) {
	var dst T
	if err := deepcopy.Copy(&dst, src); err != nil {
		return dst, err
	}
	return dst, nil
}


