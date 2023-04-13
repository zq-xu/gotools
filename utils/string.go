package utils

import "strings"

func ContainString(sli []string, elem string) bool {
	for _, e := range sli {
		if elem == e {
			return true
		}
	}
	return false
}

func IsPtrStringNotEmpty(stringPtr *string) bool {
	return stringPtr != nil && strings.TrimSpace(*stringPtr) != ""
}
