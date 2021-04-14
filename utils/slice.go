package utils

import (
	"encoding/json"
	"fmt"
	"strings"
)

// SliceIntToString this will return slice into string example "[1,2,3,4,5]"
func SliceIntToString(a []int64, delim string) string {
	return strings.Trim(strings.Replace(fmt.Sprint(a), " ", delim, -1), "")
}

// StringToIntSlice returns string into slice of int64
func StringToIntSlice(a string) []int64 {
	var is []int64
	json.Unmarshal([]byte(a), &is)
	return is
}
