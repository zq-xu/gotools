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

func readTasksFromSheet1(f *excelize.File, sheet string) ([]Task, error) {
	rows, err := f.GetRows(sheet)
	if err != nil {
		return nil, eris.Wrap(err, "failed to get rows in sheet1")
	}

	var tasks []Task
	for i, row := range rows {
		if i == 0 || len(row) < 4 {
			continue
		}
		start, _ := time.Parse(layout, row[2])
		end, _ := time.Parse(layout, row[3])

		tasks = append(tasks, Task{
			LineName: row[0],
			GroupID:  row[1],
			Start:    start,
			End:      end,
		})
	}
	return tasks, nil
}
