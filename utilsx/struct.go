package utilsx

import (
	"reflect"
	"strings"
)

// IsStructHasField checks if the roperty exists
// Attention: only for struct, return false for pointer
func IsStructHasField(s interface{}, f string) bool {
	t := reflect.TypeOf(s)
	if t.Kind() == reflect.Ptr {
		t = t.Elem() // Get the struct for the pointer
	}

	if t.Kind() == reflect.Struct {
		_, ok := t.FieldByNameFunc(func(n string) bool {
			return strings.ToLower(n) == strings.ReplaceAll(strings.ToLower(f), "_", "")
		})
		return ok
	}

	return false
}
