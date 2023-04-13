package utils

import (
	"fmt"
	"strconv"
	"time"
)

type UnixTime time.Time

// MarshalJSON implements json.Marshaler.
func (t UnixTime) MarshalJSON() ([]byte, error) {
	stamp := fmt.Sprintf("%d", time.Time(t).UnixMilli())
	return []byte(stamp), nil
}

// MarshalJSON implements json.Marshaler.
func (t *UnixTime) UnmarshalJSON(data []byte) error {
	// Ignore null, like in the main JSON package.
	s := string(data)
	if s == "null" {
		return nil
	}
	// Fractional seconds are handled implicitly by Parse.
	dataUnix, err := strconv.ParseInt(string(data), 10, 64)
	*t = UnixTime(time.UnixMilli(dataUnix))
	return err
}

func (t UnixTime) String() string {
	return (time.Time)(t).String()
}
