package database

import (
	"context"

	"github.com/zq-xu/gotools/errorx"
	pkgtypesx "github.com/zq-xu/gotools/typesx"
)

//go:generate mockgen -destination=../../../../test/mocks/database_mock.go -package=mocks pkg/store/database Database
type Database interface {
	Context(ctx context.Context) Database

	Create(t any) error
	Update(t any) error
	Delete(t any, id string) error
	DeleteAssociations(t any, id string) error

	Get(t any, id string) error
	GetByName(t any, name string) error
	GetByField(out any, key string, value any) error
	GetByMultiFields(t any, conditions map[string]any) error

	GetAssociations(t any, id string, items ...string) error
	GetAssociationsByName(t any, name string, items ...string) error
	GetAssociationsByField(t any, key string, value any, items ...string) error
	GetAssociationsByMultiFields(t any, conditions map[string]any, items ...string) error

	GetCount(t any, listParam *pkgtypesx.ListParams) (int64, error)
	List(listParam *pkgtypesx.ListParams, obj any) error
	ListWithCount(listParam *pkgtypesx.ListParams, t any, listObj any) (int64, error)

	ListAssociations(listParam *pkgtypesx.ListParams, listObj any, items ...string) error
	ListAssociationsWithCount(listParam *pkgtypesx.ListParams, t any, listObj any, items ...string) (int64, error)

	EnsureExist(obj any, id string) errorx.ErrorInfo
	EnsureNotExistByName(obj any, name string) errorx.ErrorInfo
	EnsureNotExistByField(t any, key string, value any) errorx.ErrorInfo

	DoDBTransaction(fns ...func(db Database) errorx.ErrorInfo) errorx.ErrorInfo

	// GetDBFields return all fields for table
	GetDBFields(model interface{}) ([]string, error)
}
