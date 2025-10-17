package gormkit

import (
	"net/http"
	"reflect"
	"time"

	"github.com/zq-xu/gotools/errorx"
	"gorm.io/gorm"
)

// Get all fields of the table
func (g *gormDB) GetDBFields(model interface{}) ([]string, error) {
	return GetDBFields(g.DB, model)
}

// GetDBFields returns all fields of the table
func GetDBFields(db *gorm.DB, model interface{}) ([]string, error) {
	var fields []string

	stmt := &gorm.Statement{DB: db}
	if err := stmt.Parse(model); err != nil {
		return nil, err
	}

	for _, field := range stmt.Schema.Fields {
		if field.IgnoreMigration {
			continue
		}

		if field.FieldType.Kind() == reflect.Struct &&
			field.FieldType != reflect.TypeOf(time.Time{}) {
			continue
		}

		if field.FieldType.Kind() == reflect.Ptr &&
			field.FieldType.Elem().Kind() == reflect.Struct &&
			field.FieldType.Elem() != reflect.TypeOf(time.Time{}) {
			continue
		}

		if field.FieldType.Kind() == reflect.Slice {
			continue
		}

		fields = append(fields, field.DBName)
	}

	return fields, nil
}

func DoGormDBTransaction(db *gorm.DB, fns ...func(db *gorm.DB) errorx.ErrorInfo) errorx.ErrorInfo {
	db = db.Begin()

	for _, fn := range fns {
		ei := fn(db)
		if ei != nil {
			db.Rollback()
			return ei
		}
	}

	err := db.Commit().Error
	if err != nil {
		return errorx.NewError(http.StatusInternalServerError, "transaction commit failed", err)
	}
	return nil
}
