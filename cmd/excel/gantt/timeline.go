package gantt

import (
	"sort"
	"time"
)

// TimelineLayout 管理甘特图的列布局逻辑
type TimelineLayout struct {
	DateToCol  map[string]int // 日期 -> 列索引
	BlackCols  []int          // 隔断列索引
	CurrentCol int            // 当前分配到的最大列
}

func NewTimelineLayout(startCol int) *TimelineLayout {
	return &TimelineLayout{
		DateToCol:  make(map[string]int),
		BlackCols:  []int{},
		CurrentCol: startCol,
	}
}

// Build 根据任务生成布局映射
func (l *TimelineLayout) Build(tasks []Task, year int) {
	dates := l.extractUniqueDates(tasks, year)

	var prevMonth string
	for _, dStr := range dates {
		// 1. 月份切换逻辑：截取 YYYY-MM
		currMonth := dStr[:7]
		if prevMonth != "" && currMonth != prevMonth {
			l.BlackCols = append(l.BlackCols, l.CurrentCol)
			l.CurrentCol++
		}

		// 2. 分配列索引
		l.DateToCol[dStr] = l.CurrentCol
		l.CurrentCol++
		prevMonth = currMonth
	}
}

// extractUniqueDates 内部辅助：收集并排序日期
func (l *TimelineLayout) extractUniqueDates(tasks []Task, year int) []string {
	set := make(map[string]struct{})
	for _, t := range tasks {
		for d := t.Start; !d.After(t.End); d = d.AddDate(0, 0, 1) {
			if d.Year() == year {
				set[d.Format(time.DateOnly)] = struct{}{}
			}
		}
	}

	res := make([]string, 0, len(set))
	for d := range set {
		res = append(res, d)
	}
	sort.Strings(res)
	return res
}
