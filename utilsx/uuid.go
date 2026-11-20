package utilsx

import (
	"fmt"
	"hash/fnv"
	"net"
	"os"
	"strconv"

	"github.com/sony/sonyflake"
)

// machineID 计算顺序：环境变量 GOTOLS_MACHINE_ID → 私网 IPv4 低 16 位 → hostname 哈希。
// 用于多实例部署时错开 sonyflake 的 ID 段，避免不同机器生成重复 ID。
func machineID() (uint16, error) {
	if v := os.Getenv("GOTOLS_MACHINE_ID"); v != "" {
		if n, err := strconv.ParseUint(v, 10, 16); err == nil {
			return uint16(n), nil
		}
	}

	if ip := privateIPv4(); ip != nil {
		return uint16(uint32(ip[2])<<8 | uint32(ip[3])), nil
	}

	host, err := os.Hostname()
	if err == nil {
		h := fnv.New32a()
		_, _ = h.Write([]byte(host))
		return uint16(h.Sum32() & 0xffff), nil
	}

	return 128, nil
}

func privateIPv4() net.IP {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return nil
	}
	for _, addr := range addrs {
		ipnet, ok := addr.(*net.IPNet)
		if !ok || ipnet.IP.IsLoopback() {
			continue
		}
		if ip := ipnet.IP.To4(); ip != nil {
			return ip
		}
	}
	return nil
}

var flake = newSonyflake()

func newSonyflake() *sonyflake.Sonyflake {
	f := sonyflake.NewSonyflake(sonyflake.Settings{
		MachineID: machineID,
	})
	if f == nil {
		// machineID 初始化失败时退回历史默认值，保证 GenerateUUID 仍可用
		f = sonyflake.NewSonyflake(sonyflake.Settings{
			MachineID: func() (uint16, error) { return 128, nil },
		})
	}
	return f
}

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
