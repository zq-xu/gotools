package utilsx

import (
	"strconv"
	"time"

	"github.com/jinzhu/copier"
)

// Copy
func Copy(dst interface{}, src interface{}) error {
	return copier.CopyWithOption(dst, src, copier.Option{
		DeepCopy: true,
		Converters: []copier.TypeConverter{
			{
				SrcType: *new(int64),
				DstType: *new(string),
				Fn:      int64ToStringConverter,
			},
			{
				SrcType: time.Time{},
				DstType: UnixTime{},
				Fn:      timeToUnixTime,
			},
			{
				SrcType: *new(int64),
				DstType: new(time.Time),
				Fn:      int64ToTimePointer,
			},
		},
	})
}

func int64ToStringConverter(src interface{}) (interface{}, error) {
	if v, ok := src.(int64); ok {
		return strconv.FormatInt(v, 10), nil
	}
	return src, nil
}

func timeToUnixTime(src interface{}) (interface{}, error) {
	if v, ok := src.(time.Time); ok {
		return UnixTime(v), nil
	}
	return src, nil
}

func int64ToTimePointer(src interface{}) (interface{}, error) {
	if v, ok := src.(int64); ok {
		t := time.UnixMilli(v)
		return &t, nil
	}

	return src, nil
}
