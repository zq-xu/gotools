package utilsx

import (
	"fmt"
	"strconv"
	"time"
)

func OptBoolPtr(dst, src *bool) {
	if src == nil || dst == nil {
		return
	}

	*dst = *src
}

func OptFloat64Ptr(dst, src *float64) {
	if src == nil || dst == nil {
		return
	}

	*dst = *src
}

func GetInt64PtrDefaultNil(src int64) *int64 {
	if src == 0 {
		return nil
	}

	return &src
}

func GetInt64PtrByStringPtrDefaultNil(src *string) *int64 {
	if src == nil {
		return nil
	}

	dst, _ := strconv.ParseInt(*src, 10, 64)
	if dst == 0 {
		return nil
	}

	return &dst
}

func GetStringFromInt64Ptr(ptr *int64) string {
	if ptr == nil {
		return ""
	}
	return fmt.Sprintf("%d", *ptr)
}

func OptStringPtr(dst, src *string) {
	if src == nil || dst == nil {
		return
	}

	*dst = *src
}

func OptFloat32Ptr(dst, src *float32) {
	if src == nil || dst == nil {
		return
	}

	*dst = *src
}

func OptIntPtr(dst, src *int) {
	if src == nil || dst == nil {
		return
	}

	*dst = *src
}

func OptInt64Ptr(dst, src *int64) {
	if src == nil || dst == nil {
		return
	}

	*dst = *src
}

func OptInt64ByStringPtr(dst *int64, src *string) {
	if src == nil || dst == nil {
		return
	}

	*dst, _ = strconv.ParseInt(*src, 10, 64)
}

func GetStringFromPtr(ptr *string) string {
	if ptr == nil {
		return ""
	}
	return *ptr
}

func GetIntFromPtr(ptr *int) int {
	if ptr == nil {
		return 0
	}
	return *ptr
}

func GetFloat32FromPtr(ptr *float32) float32 {
	if ptr == nil {
		return 0
	}
	return *ptr
}

func PtrBool(b bool) *bool {
	var bt = b
	return &bt
}

func PtrSting(s string) *string {
	var st = s
	return &st
}

func GetStringFromPointer(ptr *string) string {
	if ptr == nil {
		return ""
	}

	return *ptr
}

func PtrInt(i int) *int {
	var it = i
	return &it
}

func PtrInt32(i int32) *int32 {
	it := i
	return &it
}

func PtrInt64(i int64) *int64 {
	it := i
	return &it
}

func PtrIntToInt32(i int) *int32 {
	var it = int32(i)
	return &it
}

func PtrIntToInt64(i int) *int64 {
	var it = int64(i)
	return &it
}

func Int32FromPtr(ptr *int32) int32 {
	if ptr == nil {
		return 0
	}

	return *ptr
}

// ----------------- OPT UnixTime ----------------- //
func OptUnixTimePtr(dst, src *UnixTime) {
	if src == nil || dst == nil {
		return
	}

	*dst = *src
}

func OptTimePtrByUnixTimePtr(dst **time.Time, src *UnixTime) {
	if src == nil || dst == nil {
		return
	}

	*dst = (*time.Time)(src)
}

func OptTimePtrByInt64Ptr(dst **time.Time, src *int64) {
	if src == nil || dst == nil {
		return
	}

	t := time.UnixMilli(*src)
	*dst = &t
}
