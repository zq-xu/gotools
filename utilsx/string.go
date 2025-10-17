package utilsx

import (
	"strings"
	"unicode"
)

func IsPtrStringNotEmpty(stringPtr *string) bool {
	return stringPtr != nil && strings.TrimSpace(*stringPtr) != ""
}

// ConvertToSnakeCase converts CamelCase to snake_case
// Specially handles consecutive uppercase (like ID -> id)
func ConvertToSnakeCase(s string) string {
	var result []rune
	runes := []rune(s)
	length := len(runes)

	for i := 0; i < length; i++ {
		r := runes[i]

		// Insert underscore if:
		// 1. current is uppercase, AND
		// 2. previous exists, AND
		// 3. (next is lowercase OR previous is lowercase)
		if unicode.IsUpper(r) {
			if i > 0 && (unicode.IsLower(runes[i-1]) || (i+1 < length && unicode.IsLower(runes[i+1]))) {
				result = append(result, '_')
			}
			r = unicode.ToLower(r)
		}
		result = append(result, r)
	}

	return string(result)
}
