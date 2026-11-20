package gantt

import (
	"time"

	"github.com/rotisserie/eris"
	"github.com/xuri/excelize/v2"
)

type Task struct {
	LineName string
	GroupID  string
	Start    time.Time
	End      time.Time
}

// ReadTasks 从指定 Sheet 解析任务数据
func ReadTasks(f *excelize.File, sheet string) ([]Task, error) {
	rows, err := f.GetRows(sheet)
	if err != nil {
		return nil, eris.Wrapf(err, "failed to get rows from sheet: %s", sheet)
	}

	var tasks []Task
	for i, row := range rows {
		// 跳过表头或列数不足的行
		if i == 0 || len(row) < 4 {
			continue
		}

		// 假设 layout 在 gantt.go 中定义
		start, errS := time.Parse(layout, row[2])
		end, errE := time.Parse(layout, row[3])
		if errS != nil || errE != nil {
			continue // 可以在这里记录日志，跳过日期格式错误的行
		}

		tasks = append(tasks, Task{
			LineName: row[0],
			GroupID:  row[1],
			Start:    start,
			End:      end,
		})
	}
	return tasks, nil
}
