package gantt

import (
	"github.com/rotisserie/eris"
	"github.com/xuri/excelize/v2"

	"github.com/zq-xu/gotools/logx"
)

const (
	sheet1 = "Sheet1"
	sheet2 = "Sheet2"
	layout = "01-02-06"
)

func Gantt(filePath string) error {
	logx.Logger.Infof("Generate gantt table for excel %s", filePath)

	// 1. 打开文件
	f, err := excelize.OpenFile(filePath)
	if err != nil {
		return eris.Wrapf(err, "failed to open file %s", filePath)
	}
	defer f.Close()

	// 2. 读取原始数据 (调用 task.go 中的函数)
	tasks, err := ReadTasks(f, sheet1)
	if err != nil {
		return err
	}

	// 3. 构建布局对象 (TimelineLayout 现在拥有 tasks)
	timeline := NewTimelineLayout(tasks, 2)
	timeline.Build(2026)

	// 4. 执行渲染 (GanttRenderer 现在只依赖 timeline)
	renderer := NewGanttRenderer(f, sheet2, timeline)
	if err := renderer.Render(); err != nil {
		return eris.Wrap(err, "failed to render gantt chart")
	}

	// 5. 保存
	if err := f.SaveAs(filePath); err != nil {
		return eris.Wrapf(err, "failed to save file %s", filePath)
	}

	logx.Logger.Infof("Success! Gantt generated in %s", filePath)
	return nil
}
