package utils

import "strings"

func IsPtrStringNotEmpty(stringPtr *string) bool {
	return stringPtr != nil && strings.TrimSpace(*stringPtr) != ""
}
