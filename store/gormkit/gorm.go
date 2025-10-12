package gormkit

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/zq-xu/gotools/apperror"
	"github.com/zq-xu/gotools/store/database"
	pkgTypes "github.com/zq-xu/gotools/types"
)

var (
	GlobalDB database.Database
)

type gormDB struct {
	db *gorm.DB
}

func (g *gormDB) context(ctx context.Context) *gorm.DB {
	if g == nil {
		log.Fatal("nil g")
	}

	if g.db == nil {
		log.Fatal("nil db")
	}
	return g.db.WithContext(ctx)
}

func (g *gormDB) DoDBTransaction(ctx context.Context, fns ...func(db database.Database) apperror.ErrorInfo) apperror.ErrorInfo {
	db := &gormDB{db: g.context(ctx).Begin()}

	for _, fn := range fns {
		ei := fn(db)
		if ei != nil {
			db.db.Rollback()
			return ei
		}
	}

	err := db.db.Commit().Error
	if err != nil {
		return apperror.NewError(http.StatusInternalServerError, "transaction commit failed", err)
	}
	return nil
}

func (g *gormDB) Create(ctx context.Context, t any) error {
	return g.context(ctx).Omit(clause.Associations).Create(t).Error
}

func (g *gormDB) Update(ctx context.Context, t any, asc ...string) error {
	return g.context(ctx).Omit(asc...).Save(t).Error
}

func (g *gormDB) Delete(ctx context.Context, t any, id string) error {
	return g.context(ctx).Delete(t, id).Error
}

func (g *gormDB) DeleteAssociations(ctx context.Context, t any, id string) error {
	return g.context(ctx).Select(clause.Associations).Delete(t, id).Error
}

func (g *gormDB) Get(ctx context.Context, t any, id string) error {
	return g.context(ctx).Where("id = ?", id).First(t).Error
}

func (g *gormDB) GetAssociations(ctx context.Context, t any, id string) error {
	return g.context(ctx).Where("id = ?", id).First(t).Preload(clause.Associations).Error
}

func (g *gormDB) GetByName(ctx context.Context, t any, name string) error {
	return g.GetByField(ctx, t, "name", name)
}

func (g *gormDB) GetByField(ctx context.Context, t any, key string, value any) error {
	return g.context(ctx).Where(map[string]any{key: value}).Preload(clause.Associations).First(t).Error
}

func (g *gormDB) GetCount(ctx context.Context, t any, listParam *pkgTypes.ListParams) (int64, error) {
	var count int64
	db := g.context(ctx).Model(t)
	db = generateDBForQuery(db, listParam.FuzzySearchColumnList, listParam.FuzzySearchValue)
	result := db.Count(&count)
	return count, result.Error
}

func (g *gormDB) EnsureExist(ctx context.Context, t any, id string) apperror.ErrorInfo {
	err := g.Get(ctx, t, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apperror.NewError(http.StatusBadRequest, "not found", err)
		}

		return apperror.NewError(http.StatusBadRequest, "unexpected error", err)
	}

	return nil
}

func (g *gormDB) EnsureNotExistByName(ctx context.Context, t any, name string) apperror.ErrorInfo {
	return g.EnsureNotExistByField(ctx, t, "name", name)
}

func (g *gormDB) EnsureNotExistByField(ctx context.Context, t any, key string, value any) apperror.ErrorInfo {
	err := g.GetByField(ctx, t, key, value)
	if err == nil {
		return apperror.NewError(http.StatusConflict, "already exists", nil)
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil
	}

	return apperror.NewError(http.StatusInternalServerError, "unexpected error", err)
}

// List
// The value should be initialized slice, or example:
// list := make([]Model, 0)
// List(&list)
func (g *gormDB) List(ctx context.Context, listParam *pkgTypes.ListParams, listObj any) error {
	db := g.context(ctx)

	db = generateDBForQuery(db, listParam.FuzzySearchColumnList, listParam.FuzzySearchValue)
	db = optPageDB(db, listParam)

	return db.Find(listObj).Error
}

func generateDBForQuery(db *gorm.DB, fuzzySearchColumnList []string, value string) *gorm.DB {
	value = strings.TrimSpace(value)

	if value == "" || len(fuzzySearchColumnList) == 0 {
		return db
	}

	keyList := make([]string, len(fuzzySearchColumnList))
	valueList := make([]any, len(fuzzySearchColumnList))

	for k, v := range fuzzySearchColumnList {
		keyList[k] = fmt.Sprintf(" %s LIKE ? ", v)
		valueList[k] = fmt.Sprintf("%%%s%%", value)
	}

	return db.Where(strings.Join(keyList, " OR "), valueList...)
}

func optPageDB(db *gorm.DB, listParam *pkgTypes.ListParams) *gorm.DB {
	limit := listParam.PageInfo.PageSize
	offset := listParam.PageInfo.PageSize * (listParam.PageInfo.PageNum - 1)
	sortSql := listParam.Sorter.MysqlString()
	return db.Order(sortSql).Limit(limit).Offset(offset)
}
