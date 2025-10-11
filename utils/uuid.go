package utils

import (
	"fmt"

	"github.com/sony/sonyflake"
)

var flake = sonyflake.NewSonyflake(sonyflake.Settings{
	MachineID: func() (uint16, error) { return 128, nil },
})

func GenerateUUID() int64 {
	return GenerateSonyFlakeUUID()
}

func GenerateStringUUID() string {
	return fmt.Sprintf("%d", GenerateSonyFlakeUUID())
}

func GenerateSonyFlakeUUID() int64 {
	id, _ := flake.NextID()
	return int64(id)
}
