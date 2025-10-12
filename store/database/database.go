package database

import (
	"context"

	"github.com/zq-xu/gotools/apperror"
	pkgTypes "github.com/zq-xu/gotools/types"
)

//go:generate mockgen -destination=../../../../test/mocks/database_mock.go -package=mocks pkg/store/database Database

type Database interface {
	Create(ctx context.Context, t any) error
	Update(ctx context.Context, t any, asc ...string) error
	Delete(ctx context.Context, t any, id string) error
	DeleteAssociations(ctx context.Context, t any, id string) error

	Get(ctx context.Context, t any, id string) error
	GetAssociations(ctx context.Context, t any, id string) error
	GetByName(ctx context.Context, t any, name string) error
	GetByField(ctx context.Context, out any, key string, value any) error
	GetCount(ctx context.Context, t any, listParam *pkgTypes.ListParams) (int64, error)
	List(ctx context.Context, listParam *pkgTypes.ListParams, obj any) error

	DoDBTransaction(ctx context.Context, fns ...func(db Database) apperror.ErrorInfo) apperror.ErrorInfo
	EnsureExist(ctx context.Context, obj any, id string) apperror.ErrorInfo
	EnsureNotExistByName(ctx context.Context, obj any, name string) apperror.ErrorInfo
	EnsureNotExistByField(ctx context.Context, t any, key string, value any) apperror.ErrorInfo
}
