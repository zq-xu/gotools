package gantt

import (
	"sort"
	"time"
)

// TimelineLayout 管理甘特图的列布局逻辑
type TimelineLayout struct {
	Tasks      []Task         // 持有原始任务数据
	DateToCol  map[string]int // 日期 -> 列索引
	BlackCols  []int          // 隔断列索引
	CurrentCol int            // 当前分配到的最大列
}

// NewTimelineLayout 初始化布局器
func NewTimelineLayout(tasks []Task, startCol int) *TimelineLayout {
	return &TimelineLayout{
		Tasks:      tasks,
		DateToCol:  make(map[string]int),
		BlackCols:  []int{},
		CurrentCol: startCol,
	}
}

// Build 自动识别任务中的年份范围并生成布局
func (l *TimelineLayout) Build() {
	// 1. 提取所有涉及的日期（自动跨年）
	dates := l.extractUniqueDates()

	// 2. 分配列索引，处理跨月黑墙
	var prevMonth string
	for _, dStr := range dates {
		// 截取 YYYY-MM 进行跨月判断
		currMonth := dStr[:7]
		if prevMonth != "" && currMonth != prevMonth {
			l.BlackCols = append(l.BlackCols, l.CurrentCol)
			l.CurrentCol++
		}

		l.DateToCol[dStr] = l.CurrentCol
		l.CurrentCol++
		prevMonth = currMonth
	}
}

// extractUniqueDates 内部辅助：自动从 Tasks 收集所有年份的日期并排序
func (l *TimelineLayout) extractUniqueDates() []string {
	set := make(map[string]struct{})

	for _, t := range l.Tasks {
		// 遍历任务从 Start 到 End 的每一天
		for d := t.Start; !d.After(t.End); d = d.AddDate(0, 0, 1) {
			// 直接格式化存储，不再过滤年份
			set[d.Format(time.DateOnly)] = struct{}{}
		}
	}

	res := make([]string, 0, len(set))
	for d := range set {
		res = append(res, d)
	}

	// 排序后，跨年的日期也会按先后顺序排好（如 2025-12-31, 2026-01-01）
	sort.Strings(res)
	return res
}
