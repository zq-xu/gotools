package utilsx

import "unsafe"

type InterfaceStructure struct {
	pt uintptr // value type pointer
	pv uintptr // value content pointer
}

// AsInterfaceStructure transfers the interface to InterfaceStructure
func AsInterfaceStructure(i interface{}) InterfaceStructure {
	return *(*InterfaceStructure)(unsafe.Pointer(&i))
}

// IsInterfaceValueNil checks if the interface value is nil
func IsInterfaceValueNil(i interface{}) bool {
	is := AsInterfaceStructure(i)
	return is.pv == 0
}
