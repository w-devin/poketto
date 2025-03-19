package array

import "strings"

func In(a string, arr []string) bool {
	for _, b := range arr {
		if a == b {
			return true
		}
	}

	return false
}

func InFold(a string, arr []string) bool {
	for _, b := range arr {
		if strings.EqualFold(a, b) {
			return true
		}
	}

	return false
}

func IndexOfFold(a string, arr []string) int {
	for index, b := range arr {
		if strings.EqualFold(a, b) {
			return index
		}
	}

	return -1
}

func Distinct(arr []string) []string {
	ret := make([]string, 0, len(arr))
	for _, v := range arr {
		if !In(v, ret) {
			ret = append(ret, v)
		}
	}

	return ret
}

func DistinctFold(arr []string) []string {
	ret := make([]string, 0, len(arr))
	for _, v := range arr {
		if !InFold(v, ret) {
			ret = append(ret, v)
		}
	}

	return ret
}
