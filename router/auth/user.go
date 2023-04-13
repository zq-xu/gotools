package auth

import (
	"zq-xu/helper/store"
)

const (
	UserTableName = "user"
)

// description:"default the column ID is the primaryKey
type User struct {
	store.Model

	Name     string
	Alias    string
	Password string
	Comment  string
	Role     int
	Status   int
}

func (u *User) TableName() string {
	return UserTableName
}

func init() {
	store.RegisterTable(&User{})
}
