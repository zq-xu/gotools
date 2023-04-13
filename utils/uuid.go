package utils

import (
	"fmt"

	"github.com/bwmarrin/snowflake"
	"github.com/google/uuid"
	"github.com/sony/sonyflake"
)

var snowflakeNode, _ = snowflake.NewNode(1)
var flake = sonyflake.NewSonyflake(sonyflake.Settings{
	MachineID: func() (uint16, error) { return 128, nil },
})

func GenerateUUID() int64 {
	return GenerateSonyFlakeUUID()
}

func GenerateStringUUID() string {
	return fmt.Sprintf("%d", GenerateSonyFlakeUUID())
}

func GenerateSnowFlakeUUID() int64 {
	return snowflakeNode.Generate().Int64()
}

func GenerateSonyFlakeUUID() int64 {
	id, _ := flake.NextID()
	return int64(id)
}

func GenerateGoogleUUID() uint32 {
	return uuid.New().ID()
}
