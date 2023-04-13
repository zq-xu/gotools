package utils

import (
	"fmt"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

// go test -v \
// excel.go excel_test.go \
// -test.run TestExportToPath -count=1
func TestExportToPath(t *testing.T) {
	dataKey := make([]map[string]string, 0)
	dataKey = append(dataKey, map[string]string{
		"key":    "id",
		"title":  "索引",
		"width":  "20",
		"is_num": "0",
	})
	dataKey = append(dataKey, map[string]string{
		"key":    "username",
		"title":  "用户名",
		"width":  "20",
		"is_num": "0",
	})
	dataKey = append(dataKey, map[string]string{
		"key":    "remark",
		"title":  "备注",
		"width":  "20",
		"is_num": "0",
	})

	//填充数据
	data := make([]map[string]interface{}, 0)

	for i := 0; i < 10; i++ {
		data = append(data, map[string]interface{}{
			"id":       fmt.Sprintf("%d-id", i),
			"username": fmt.Sprintf("%d-username", i),
			"remark":   fmt.Sprintf("%d-remark", i),
		})
	}

	Convey("TestExportToPath", t, func() {
		Convey("TestExportToPath.Local", func() {
			ee := NewExcelExport("as1")
			_, err := ee.ExportToPath(dataKey, data, "./")
			So(err, ShouldBeNil)
		})
	})
}
