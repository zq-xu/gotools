package types

import (
	"fmt"
	"net/http"
	"slices"
	"strings"

	"zq-xu/gotools/apperror"
	"zq-xu/gotools/utils"
)

const (
	sortByQuery    = "sortBy"
	defaultSortStr = "updated_at,desc"
	sortSplit      = ":::"

	asc  = "asc"
	desc = "desc"
)

var (
	orderList      = []string{"asc", "ascending", "desc", "descending"}
	baseConditions = []string{"updated_at", "created_at", "id"}
)

type sorter struct {
	conditions []sortCondition
}

type sortCondition struct {
	condition string
	order     string
}

// NewSorter analyse the sort string, like: name,asc:::age,desc
func NewSorter(sortQuery string) *sorter {
	cs := strings.Split(sortQuery, sortSplit)
	res := make([]sortCondition, 0, len(cs))

	for _, v := range cs {
		sc := generateCondition(v)
		if sc.condition == "" {
			continue
		}
		res = append(res, sc)
	}

	s := &sorter{conditions: res}
	return s
}

// return likes "name asc,alias desc"
func (sh *sorter) MysqlString() string {
	res := make([]string, 0, len(sh.conditions))
	for _, sc := range sh.conditions {
		res = append(res, fmt.Sprintf("%s %s", sc.condition, sc.order))
	}
	return strings.Join(res, ",")
}

func generateCondition(str string) sortCondition {
	tmp := strings.Split(str, ",")
	condition, order := "", ""

	switch len(tmp) {
	case 1:
		condition = tmp[0]
	case 2:
		condition = tmp[0]
		order = tmp[1]
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
	if slices.Contains(orderList, str) {
		return strings.TrimSuffix(str, "ending")
	}
	return asc
}

func (s *sorter) validateSortQuery(obj interface{}) apperror.ErrorInfo {
	for _, sc := range s.conditions {
		if !slices.Contains(baseConditions, sc.condition) && !utils.IsStructHasField(obj, sc.condition) {
			return apperror.Errorf(nil, http.StatusBadRequest, "invalid sort key %s", sc.condition)
		}
	}
	return nil
}
