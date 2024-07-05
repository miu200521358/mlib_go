package mutils

import (
	"fmt"
	"slices"
	"strconv"
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

func JoinIntsWithComma(ints []int) string {
	var strList []string
	for _, num := range ints {
		strList = append(strList, strconv.Itoa(num))
	}
	return strings.Join(strList, ", ")
}

func SplitCommaSeparatedInts(s string) ([]int, error) {
	var ints []int
	strList := strings.Split(s, ", ")
	for _, str := range strList {
		num, err := strconv.Atoi(str)
		if err != nil {
			return nil, err
		}
		ints = append(ints, num)
	}
	return ints, nil
}
