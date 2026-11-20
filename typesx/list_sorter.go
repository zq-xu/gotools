package typesx

import (
	"fmt"
	"net/http"
	"regexp"
	"slices"
	"strings"

	"github.com/zq-xu/gotools/errorx"
	"github.com/zq-xu/gotools/utilsx"
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

	// sortKeyPattern 限制排序键只能是普通标识符，杜绝 "sortBy" 里塞进 SQL 片段。
	sortKeyPattern = regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_]*$`)
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

	// 全部非法（或没给）时回退到默认排序，保证分页结果稳定，且 ORDER BY 不会为空。
	if len(res) == 0 {
		res = append(res, generateCondition(defaultSortStr))
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

	// 非法键直接丢弃，避免未经校验的字符串进入 ORDER BY。
	if !sortKeyPattern.MatchString(condition) {
		return sortCondition{}
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

func (s *sorter) validateSortQuery(obj interface{}) errorx.ErrorInfo {
	for _, sc := range s.conditions {
		if !slices.Contains(baseConditions, sc.condition) && !utilsx.IsStructHasField(obj, sc.condition) {
			return errorx.Errorf(nil, http.StatusBadRequest, "invalid sort key %s", sc.condition)
		}
	}
	return nil
}
