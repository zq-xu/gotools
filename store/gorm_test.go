package store

import (
	"testing"
)

// go test -v gorm_test.go gorm.go -test.run TestGorm
func TestGorm(t *testing.T) {
	dbInfo := &DatabaseInfo{
		Address:      "192.168.1.99",
		Port:         "3306",
		Username:     "root",
		Password:     "root",
		DatabaseName: "beluga",
	}
	t.Logf("config is%+v", dbInfo)
}
