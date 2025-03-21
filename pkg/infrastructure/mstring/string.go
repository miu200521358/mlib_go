package mstring

import (
	"fmt"
	"runtime"
	"slices"
	"strconv"
	"strings"
)

func SplitAll(input string, separators []string) []string {
	results := []string{input}

	for _, sep := range separators {
		var tempResult []string
		for _, str := range results {
			tempResult = append(tempResult, strings.Split(str, sep)...)
		}
		results = tempResult
	}

	return results
}

func JoinSlice[T any](slice []T) string {
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

func DeepCopyIntSlice(original []int) []int {
	newSlice := make([]int, len(original))
	copy(newSlice, original)
	return newSlice
}

func DeepCopyStringSlice(original []string) []string {
	newSlice := make([]string, len(original))
	copy(newSlice, original)
	return newSlice
}

func GetStackTrace() string {
	buf := make([]byte, 1024)
	n := runtime.Stack(buf, true)
	return string(buf[:n])
}
