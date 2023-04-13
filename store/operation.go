package store

import (
	"strings"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func Create(db *gorm.DB, t interface{}) error {
	return db.Omit(clause.Associations).Create(t).Error
}

func CreateOmit(db *gorm.DB, t interface{}, asc ...string) error {
	return db.Omit(asc...).Create(t).Error
}

func Update(db *gorm.DB, t interface{}) error {
	return db.Save(t).Error
}

func Delete(db *gorm.DB, t interface{}) error {
	return db.Delete(t).Error
}

func DeleteByID(db *gorm.DB, t interface{}, id string) error {
	return db.Delete(t, id).Error
}

func DeleteAssociationsByID(db *gorm.DB, t interface{}, id string) error {
	return db.Select(clause.Associations).Delete(t, id).Error
}

func GetByID(db *gorm.DB, t interface{}, id string) error {
	return db.Where("id = ?", id).First(t).Error
}

func GetAssociationsByID(db *gorm.DB, t interface{}, id string) error {
	return db.Preload(clause.Associations).Where("id = ?", id).First(t).Error
}

func GetByName(db *gorm.DB, t interface{}, name string) error {
	return db.Where("name = ?", name).Preload(clause.Associations).First(t).Error
}

func GetCount(db *gorm.DB, t interface{}) (int64, error) {
	var count int64
	result := db.Model(t).Count(&count)
	return count, result.Error
}

// The value should be initialized slice, or example:
// list := make([]Model, 0)
// List(&list)
func List(db *gorm.DB, value interface{}) error {
	return db.Preload(clause.Associations).Find(value).Error
}

func GenerateDBForQuery(db *gorm.DB, fuzzySearchColumnList []string, value string) *gorm.DB {
	if value != "" && len(fuzzySearchColumnList) > 0 {
		keyList, valueList := make([]string, len(fuzzySearchColumnList)), make([]interface{}, len(fuzzySearchColumnList))
		for k, v := range fuzzySearchColumnList {
			keyList[k] = v + " LIKE ? "
			valueList[k] = "%" + value + "%"
		}
		db = db.Where(strings.Join(keyList, "OR"), valueList...)
	}
	return db
}

func OptPageDB(db *gorm.DB, pageSize, pageNum int, sortQuery string, obj interface{}) *gorm.DB {
	limit := pageSize
	offset := pageSize * (pageNum - 1)
	sortSql := GenerateSortSql(sortQuery, obj)
	return db.Order(sortSql).Limit(limit).Offset(offset)
}

func GenerateSortSql(sortQuery string, obj interface{}) string {
	sort := NewSorter(sortQuery)
	sort.Purge(obj)
	return sort.SQLString()
}
