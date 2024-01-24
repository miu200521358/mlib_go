package mutils

import (
	"fmt"
	"strings"

)

func Contains(arr []string, str string) bool {
	for _, a := range arr {
		if a == str {
			return true
		}
	}
	return false
}

func ConvertFloat32ToInterfaceSlice(slice []float32) []interface{} {
	interfaceSlice := make([]interface{}, len(slice))
	for i, v := range slice {
		interfaceSlice[i] = v
	}
	return interfaceSlice
}

func ConvertUint32ToInterfaceSlice(slice []uint32) []interface{} {
	interfaceSlice := make([]interface{}, len(slice))
	for i, v := range slice {
		interfaceSlice[i] = v
	}
	return interfaceSlice
}

func JoinSlice(slice []interface{}) string {
	// Convert each item in the slice to a string.
	strSlice := make([]string, len(slice))
	for i, v := range slice {
		strSlice[i] = fmt.Sprint(v)
	}

	// Join the string slice into a single string with commas in between.
	return "[" + strings.Join(strSlice, ", ") + "]"
}
