// 指示: miu200521358
package mstring

import (
	"fmt"
	"runtime"
	"slices"
	"strconv"
	"strings"
)

// SplitAll は複数セパレータで分割する。
func SplitAll(input string, separators []string) []string {
	results := []string{input}
	for _, sep := range separators {
		var temp []string
		for _, str := range results {
			temp = append(temp, strings.Split(str, sep)...)
		}
		results = temp
	}
	return results
}

// JoinSlice はスライスを文字列表現にまとめる。
func JoinSlice[T any](slice []T) string {
	strSlice := make([]string, len(slice))
	for i, v := range slice {
		strSlice[i] = fmt.Sprint(v)
	}
	return "[" + strings.Join(strSlice, ", ") + "]"
}

// RemoveFromSlice は指定値を削除したスライスを返す。
func RemoveFromSlice[S ~[]E, E comparable](slice S, value E) []E {
	index := slices.Index(slice, value)
	if index == -1 {
		return slice
	}
	return append(slice[:index], slice[index+1:]...)
}

// JoinIntsWithComma はintスライスをカンマ区切り文字列に変換する。
func JoinIntsWithComma(ints []int) string {
	strList := make([]string, 0, len(ints))
	for _, num := range ints {
		strList = append(strList, strconv.Itoa(num))
	}
	return strings.Join(strList, ", ")
}

// SplitCommaSeparatedInts はカンマ区切り文字列をintスライスに変換する。
func SplitCommaSeparatedInts(s string) ([]int, error) {
	trimmed := strings.TrimSpace(s)
	if trimmed == "" {
		return []int{}, nil
	}
	parts := strings.Split(trimmed, ",")
	ints := make([]int, 0, len(parts))
	for _, part := range parts {
		value := strings.TrimSpace(part)
		if value == "" {
			continue
		}
		num, err := strconv.Atoi(value)
		if err != nil {
			return nil, newIntParseFailed("数値パースに失敗しました", err)
		}
		ints = append(ints, num)
	}
	return ints, nil
}

// DeepCopyIntSlice はintスライスを複製する。
func DeepCopyIntSlice(original []int) []int {
	newSlice := make([]int, len(original))
	copy(newSlice, original)
	return newSlice
}

// DeepCopyStringSlice はstringスライスを複製する。
func DeepCopyStringSlice(original []string) []string {
	newSlice := make([]string, len(original))
	copy(newSlice, original)
	return newSlice
}

// GetStackTrace はスタックトレースを返す。
func GetStackTrace() string {
	buf := make([]byte, 1024)
	n := runtime.Stack(buf, true)
	return string(buf[:n])
}
