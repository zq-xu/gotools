package store

import (
	"fmt"
	"reflect"
	"strings"
)

const (
	asc  = "asc"
	desc = "desc"

	sortSplit = ":::"
)

type sorter struct {
	conditions []sortCondition
	data       interface{}
}

type sortCondition struct {
	condition string
	order     string
}

func NewSorter(s string) *sorter {
	cs := strings.Split(s, sortSplit)
	res := make([]sortCondition, 0, len(cs))

	for _, v := range cs {
		sc := generateCondition(v)
		if sc.condition == "" {
			continue
		}

		res = append(res, sc)
	}

	sh := &sorter{conditions: res}
	return sh
}

// return likes "name asc"
func (sc sortCondition) SQLString() string {
	return fmt.Sprintf("%s %s", sc.condition, sc.order)
}

// return likes "name asc,alias desc"
func (sh *sorter) SQLString() string {
	res := make([]string, 0, len(sh.conditions))
	for _, sc := range sh.conditions {
		res = append(res, fmt.Sprintf("%s %s", sc.condition, sc.order))
	}
	return strings.Join(res, ",")
}

func (sh *sorter) Purge(st interface{}) {
	res := make([]sortCondition, 0, len(sh.conditions))
	for _, sc := range sh.conditions {
		if StructHasField(st, sc.condition) {
			res = append(res, sc)
		}
	}
	sh.conditions = res
}

func generateCondition(str string) sortCondition {
	tmp := strings.Split(str, ",")
	condition, order := "", ""

	switch len(tmp) {
	case 2:
		order = tmp[1]
		fallthrough
	case 1:
		condition = tmp[0]
	}

	if order == "" {
		order = asc
	}

	return sortCondition{
		condition: condition,
		order:     generateOrder(order),
	}
}

func generateOrder(str string) string {
	switch str {
	case desc:
		return desc
	case asc:
		return asc
	default:
		return asc
	}
}

// check if the property exists
// Attention: only for struct, return false for pointer
func StructHasField(s interface{}, f string) bool {
	t := reflect.TypeOf(s)
	if t.Kind() == reflect.Struct {
		_, ok := t.FieldByNameFunc(func(n string) bool {
			return strings.ToLower(n) == strings.ToLower(f)
		})
		return ok
	}

	return false
}
