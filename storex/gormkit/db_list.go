package gormkit

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/rotisserie/eris"
	"gorm.io/gorm"

	pkgtypesx "github.com/zq-xu/gotools/typesx"
	"github.com/zq-xu/gotools/utilsx"
)

func (g *gormDB) GetCount(t any, listParam *pkgtypesx.ListParams) (int64, error) {
	var count int64
	db := g.DB.Model(t)
	db = GenerateDBForQuery(db, listParam, t)
	db = OptFuzzySearchDB(db, listParam.FuzzySearchColumnList, listParam.FuzzySearchValue)
	result := db.Count(&count)
	return count, result.Error
}

func (g *gormDB) ListWithCount(listParam *pkgtypesx.ListParams, t any, listObj any) (int64, error) {
	count, err := g.GetCount(t, listParam)
	if err != nil {
		return 0, eris.Wrap(err, "get count error")
	}

	err = g.List(listParam, listObj)
	if err != nil {
		return 0, eris.Wrap(err, "list error")
	}

	return count, nil
}

func (g *gormDB) ListAssociationsWithCount(listParam *pkgtypesx.ListParams, t any, listObj any, items ...string) (int64, error) {
	count, err := g.GetCount(t, listParam)
	if err != nil {
		return 0, eris.Wrap(err, "get count error")
	}

	err = g.ListAssociations(listParam, listObj, items...)
	if err != nil {
		return 0, eris.Wrap(err, "list error")
	}

	return count, nil
}

// List
// The value should be initialized slice, or example:
// list := make([]Model, 0)
// List(&list)
func (g *gormDB) List(listParam *pkgtypesx.ListParams, listObj any) error {
	db := optListDB(g.DB, listParam, listObj)
	return db.Find(listObj).Error
}

func optListDB(db *gorm.DB, listParam *pkgtypesx.ListParams, listObj any) *gorm.DB {
	db = GenerateDBForQuery(db, listParam, listObj)
	db = OptFuzzySearchDB(db, listParam.FuzzySearchColumnList, listParam.FuzzySearchValue)
	db = OptPageDB(db, listParam)
	return db
}

func (g *gormDB) ListAssociations(listParam *pkgtypesx.ListParams, listObj any, items ...string) error {
	db := optListDB(g.DB, listParam, listObj)
	return getAssociations(db, items...).Find(listObj).Error
}

func GenerateDBForQuery(db *gorm.DB, listParam *pkgtypesx.ListParams, listObj any) *gorm.DB {
	searchQueries := filterQueriesByStructFields(listObj, listParam.Queries)

	keyList := make([]string, 0)
	valueList := make([]any, 0)
	for k, v := range searchQueries {
		keyList = append(keyList, fmt.Sprintf(" %s = ? ", k))
		valueList = append(valueList, v)
	}

	return db.Where(strings.Join(keyList, " And "), valueList...)
}

func OptFuzzySearchDB(db *gorm.DB, fuzzySearchColumnList []string, value string) *gorm.DB {
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

func OptPageDB(db *gorm.DB, listParam *pkgtypesx.ListParams) *gorm.DB {
	limit := listParam.PageInfo.PageSize
	offset := listParam.PageInfo.PageSize * (listParam.PageInfo.PageNum - 1)
	sortSql := listParam.Sorter.MysqlString()
	return db.Order(sortSql).Limit(limit).Offset(offset)
}

// filterQueriesByStructFields filters query keys that match struct fields of listObj
// listObj should be a pointer to a slice, e.g. &[]Model{}
func filterQueriesByStructFields(listObj any, queries map[string]string) map[string]string {
	typ := reflect.TypeOf(listObj)

	// Handle pointer to slice
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}

	// Handle slice
	if typ.Kind() == reflect.Slice {
		typ = typ.Elem()
	}

	// Handle pointer to struct element
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}

	result := make(map[string]string)
	if typ.Kind() != reflect.Struct {
		return result
	}

	// Build field name map
	fieldMap := make(map[string]struct{})
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		snake := utilsx.ConvertToSnakeCase(field.Name)
		fieldMap[snake] = struct{}{}
	}

	// Filter queries
	for key, val := range queries {
		keyLower := strings.ToLower(key)
		if _, ok := fieldMap[keyLower]; ok {
			result[key] = val
		}
	}

	return result
}
