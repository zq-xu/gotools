package storex

import (
	"context"

	"gorm.io/gorm"

	"github.com/zq-xu/gotools/storex/database"
	"github.com/zq-xu/gotools/storex/gormkit"
)

var (
	RegisterTable = gormkit.RegisterTable

	GenerateModel       = gormkit.GenerateModel
	GenerateModelWithID = gormkit.GenerateModelWithID

	GetDBFields = gormkit.GetDBFields

	GenerateDBForQuery = gormkit.GenerateDBForQuery
	OptFuzzySearchDB   = gormkit.OptFuzzySearchDB
	OptPageDB          = gormkit.OptPageDB

	DoGormDBTransaction = gormkit.DoGormDBTransaction

	NewErrorInfoForCreateError = gormkit.NewErrorInfoForCreateError
	NewErrorInfoForUpdateError = gormkit.NewErrorInfoForUpdateError
	NewErrorInfoForDeleteError = gormkit.NewErrorInfoForDeleteError
	NewErrorInfoForGetError    = gormkit.NewErrorInfoForGetError
	NewErrorInfoForListError   = gormkit.NewErrorInfoForListError
	IsNotFoundError            = gormkit.IsNotFoundError
)

type (
	Model = gormkit.Model
)

func DB(ctx context.Context) database.Database {
	return gormkit.GlobalDB.Context(ctx)
}

func GormDB(ctx context.Context) *gorm.DB {
	return gormkit.GormDB.WithContext(ctx)
}

func SetDB(d database.Database) {
	gormkit.GlobalDB = d
}
