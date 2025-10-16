package gormkit

import (
	"context"
	"errors"
	"net/http"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/zq-xu/gotools/apperror"
	"github.com/zq-xu/gotools/store/database"
)

var (
	GormDB   *gorm.DB
	GlobalDB database.Database
)

type gormDB struct {
	*gorm.DB
}

func (g *gormDB) Context(ctx context.Context) database.Database {
	g.DB = g.DB.WithContext(ctx)
	return g
}

func (g *gormDB) DoDBTransaction(fns ...func(db database.Database) apperror.ErrorInfo) apperror.ErrorInfo {
	db := &gormDB{DB: g.DB.Begin()}

	for _, fn := range fns {
		ei := fn(db)
		if ei != nil {
			db.DB.Rollback()
			return ei
		}
	}

	err := db.DB.Commit().Error
	if err != nil {
		return apperror.NewError(http.StatusInternalServerError, "transaction commit failed", err)
	}
	return nil
}

func (g *gormDB) Create(t any) error {
	return g.DB.Omit(clause.Associations).Create(t).Error
}

func (g *gormDB) Update(t any) error {
	return g.DB.Omit(clause.Associations).Save(t).Error
}

func (g *gormDB) Delete(t any, id string) error {
	return g.DB.Delete(t, id).Error
}

func (g *gormDB) DeleteAssociations(t any, id string) error {
	return g.DB.Select(clause.Associations).Delete(t, id).Error
}

func (g *gormDB) Get(t any, id string) error {
	return g.DB.Where("id = ?", id).First(t).Error
}

func (g *gormDB) GetByName(t any, name string) error {
	return g.GetByField(t, "name", name)
}

func (g *gormDB) GetByField(t any, key string, value any) error {
	return g.DB.Where(map[string]any{key: value}).First(t).Error
}

func (g *gormDB) GetByMultiFields(t any, conditions map[string]any) error {
	return g.DB.Where(conditions).First(t).Error
}

func getAssociations(db *gorm.DB, items ...string) *gorm.DB {
	if len(items) == 0 {
		return db.Preload(clause.Associations)
	}

	for _, v := range items {
		db = db.Preload(v)
	}

	return db
}

func (g *gormDB) GetAssociations(t any, id string, items ...string) error {
	return getAssociations(g.DB, items...).Where("id = ?", id).First(t).Error
}

func (g *gormDB) GetAssociationsByName(t any, name string, items ...string) error {
	return g.GetAssociationsByField(t, "name", name, items...)
}

func (g *gormDB) GetAssociationsByField(t any, key string, value any, items ...string) error {
	return getAssociations(g.DB, items...).Where(map[string]any{key: value}).First(t).Error
}

func (g *gormDB) GetAssociationsByMultiFields(t any, conditions map[string]any, items ...string) error {
	return getAssociations(g.DB, items...).Where(conditions).First(t).Error
}

func (g *gormDB) EnsureExist(t any, id string) apperror.ErrorInfo {
	err := g.Get(t, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apperror.NewError(http.StatusBadRequest, "not found", err)
		}

		return apperror.NewError(http.StatusBadRequest, "unexpected error", err)
	}

	return nil
}

func (g *gormDB) EnsureNotExistByName(t any, name string) apperror.ErrorInfo {
	return g.EnsureNotExistByField(t, "name", name)
}

func (g *gormDB) EnsureNotExistByField(t any, key string, value any) apperror.ErrorInfo {
	err := g.GetByField(t, key, value)
	if err == nil {
		return apperror.NewError(http.StatusConflict, "already exists", nil)
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil
	}

	return apperror.NewError(http.StatusInternalServerError, "unexpected error", err)
}
