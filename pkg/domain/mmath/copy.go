package mmath

import "github.com/tiendc/go-deepcopy"

func deepCopy[T any](src T) (T, error) {
	var dst T
	if err := deepcopy.Copy(&dst, src); err != nil {
		return dst, err
	}
	return dst, nil
}
