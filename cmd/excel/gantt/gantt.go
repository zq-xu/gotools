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
	logx.Logger.Infof("Loaded excel %s", filePath)

	// 2. 读取原始数据
	tasks, err := readTasksFromSheet1(f, sheet1)
	if err != nil {
		return err
	}
	logx.Logger.Infof("Readed data %s", sheet1)

	// 3. 初始化并构建时间轴布局逻辑 (从第2列开始)
	layout := NewTimelineLayout(2)
	layout.Build(tasks, 2026)

	logx.Logger.Info("Built timeline")

	// 4. 创建渲染器并执行全流程渲染
	// 我们将数据处理和 Excel 写入完全解耦
	err = NewGanttRenderer(f, sheet2, layout).Render(tasks)
	if err != nil {
		return eris.Wrap(err, "failed to render gantt chart")
	}
	logx.Logger.Infof("Rendered gantt in excel sheet %s", sheet2)

	// 5. 保存结果
	if err := f.SaveAs(filePath); err != nil {
		return eris.Wrapf(err, "failed to save file %s", filePath)
	}

	logx.Logger.Infof("Saved gantt in excel sheet %s", sheet2)
	return nil
}
