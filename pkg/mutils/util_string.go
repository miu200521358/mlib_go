package mutils

import (
	"fmt"
	"slices"
	"strings"

)

func JoinSlice(slice []interface{}) string {
	// Convert each item in the slice to a string.
	strSlice := make([]string, len(slice))
	for i, v := range slice {
		strSlice[i] = fmt.Sprint(v)
	}

	// Join the string slice into a single string with commas in between.
	return "[" + strings.Join(strSlice, ", ") + "]"
}

func RemoveFromSlice[S ~[]E, E comparable](slice S, value E) []E {
	// value が含まれている index を探す
	index := slices.Index(slice, value)
	if index == -1 {
		return slice
	}
	return append(slice[:index], slice[index+1:]...)
}
