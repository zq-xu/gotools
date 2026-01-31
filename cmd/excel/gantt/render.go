package gantt

import (
	"fmt"
	"sort"
	"time"

	"github.com/xuri/excelize/v2"
)

// GanttRenderer 负责所有 Excel 绘图逻辑
type GanttRenderer struct {
	f          *excelize.File
	sheet      string
	layout     *TimelineLayout
	taskStyle  int
	monthStyle int
	dateStyle  int
	blackStyle int
}

func NewGanttRenderer(f *excelize.File, sheet string, layout *TimelineLayout) *GanttRenderer {
	return &GanttRenderer{
		f:      f,
		sheet:  sheet,
		layout: layout,
	}
}

// Render 执行完整的渲染流程
func (r *GanttRenderer) Render() error {
	r.prepareStyles()
	r.initSheet()
	r.drawBlackWalls()
	r.drawTimeline()
	r.drawTasks()
	return nil
}

func (r *GanttRenderer) prepareStyles() {
	// 任务条样式（带边框的绿色填充）
	r.taskStyle, _ = r.f.NewStyle(&excelize.Style{
		Fill: excelize.Fill{Type: "pattern", Color: []string{"#C6EFCE"}, Pattern: 1},
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1}, {Type: "top", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1}, {Type: "bottom", Color: "000000", Style: 1},
		},
	})
	// 月份表头样式
	r.monthStyle, _ = r.f.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
		Font:      &excelize.Font{Bold: true, Size: 10},
		Fill:      excelize.Fill{Type: "pattern", Color: []string{"#F2F2F2"}, Pattern: 1},
	})
	// 日期表头样式
	r.dateStyle, _ = r.f.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
		Font:      &excelize.Font{Size: 9},
	})
	// 隔断列全黑样式
	r.blackStyle, _ = r.f.NewStyle(&excelize.Style{
		Fill: excelize.Fill{Type: "pattern", Color: []string{"#000000"}, Pattern: 1},
	})
}

func (r *GanttRenderer) initSheet() {
	index, _ := r.f.NewSheet(r.sheet)
	r.f.SetActiveSheet(index)
	// 冻结第一列(A)和前两行
	_ = r.f.SetPanes(r.sheet, &excelize.Panes{
		Freeze: true, XSplit: 1, YSplit: 2, TopLeftCell: "B3", ActivePane: "bottomRight",
	})
	r.f.SetColWidth(r.sheet, "A", "A", 35)
	r.f.SetCellValue(r.sheet, "A1", "Line & Group ID")
	r.f.SetCellValue(r.sheet, "A2", "Timeline")
}

func (r *GanttRenderer) drawBlackWalls() {
	// 直接从 layout 获取任务数
	maxRow := len(r.layout.Tasks) + 2
	for _, col := range r.layout.BlackCols {
		name, _ := excelize.ColumnNumberToName(col)
		r.f.SetCellStyle(r.sheet, fmt.Sprintf("%s1", name), fmt.Sprintf("%s%d", name, maxRow), r.blackStyle)
		r.f.SetColWidth(r.sheet, name, name, 0.6)
	}
}

func (r *GanttRenderer) drawTimeline() {
	dates := make([]string, 0, len(r.layout.DateToCol))
	for d := range r.layout.DateToCol {
		dates = append(dates, d)
	}
	sort.Strings(dates)

	var lastMonth string
	var monthStartCol int

	for i, dStr := range dates {
		col := r.layout.DateToCol[dStr]
		t, _ := time.Parse(time.DateOnly, dStr)
		monthLabel := t.Format("2006/01")

		colName, _ := excelize.ColumnNumberToName(col)
		r.f.SetColWidth(r.sheet, colName, colName, 3.5)

		cell, _ := excelize.CoordinatesToCellName(col, 2)
		r.f.SetCellValue(r.sheet, cell, t.Day())
		r.f.SetCellStyle(r.sheet, cell, cell, r.dateStyle)

		if lastMonth == "" {
			monthStartCol = col
			lastMonth = monthLabel
		} else if monthLabel != lastMonth {
			r.mergeMonth(monthStartCol, col-2, lastMonth)
			monthStartCol = col
			lastMonth = monthLabel
		}
		if i == len(dates)-1 {
			r.mergeMonth(monthStartCol, col, lastMonth)
		}
	}
}

func (r *GanttRenderer) mergeMonth(start, end int, label string) {
	sCell, _ := excelize.CoordinatesToCellName(start, 1)
	eCell, _ := excelize.CoordinatesToCellName(end, 1)
	r.f.MergeCell(r.sheet, sCell, eCell)
	r.f.SetCellValue(r.sheet, sCell, label)
	r.f.SetCellStyle(r.sheet, sCell, sCell, r.monthStyle)
}

func (r *GanttRenderer) drawTasks() {
	// 直接从持有的 layout.Tasks 遍历
	for i, task := range r.layout.Tasks {
		row := i + 3
		r.f.SetCellValue(r.sheet, fmt.Sprintf("A%d", row), fmt.Sprintf("%s - %s", task.LineName, task.GroupID))

		for d := task.Start; !d.After(task.End); d = d.AddDate(0, 0, 1) {
			if col, ok := r.layout.DateToCol[d.Format(time.DateOnly)]; ok {
				cell, _ := excelize.CoordinatesToCellName(col, row)
				r.f.SetCellStyle(r.sheet, cell, cell, r.taskStyle)
			}
		}
	}
}
