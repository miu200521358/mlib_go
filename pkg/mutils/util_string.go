package mutils

import (
	"fmt"
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
