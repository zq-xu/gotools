package utils

import (
	"fmt"
	"net/url"
	"path"
	"reflect"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
)

const (
	defaultSheetName = "Sheet1"
	defaultHeight    = 25.0 //default row height
	startColumn      = 'A'
)

var (
	centerAlignment = &excelize.Alignment{
		Horizontal: "center",
		Vertical:   "center",
	}

	defaultTopStyle = &excelize.Style{
		Font: &excelize.Font{
			Bold: true,
			Size: 12,
		},
		Alignment: centerAlignment,
	}

	defaultLineStyle = &excelize.Style{Alignment: centerAlignment}

	defaultCellStyle = &excelize.Style{
		Font: &excelize.Font{
			//Color:  "#666666",
			Bold: false,
			Size: 10,
			//Family: "arial",
		},
		Alignment: centerAlignment,
	}
)

type excelExport struct {
	file      *excelize.File
	sheetName string
	endCol    string
}

func NewExcelExport(sheetName string) *excelExport {
	if sheetName == "" {
		sheetName = defaultSheetName
	}

	return &excelExport{file: createFile(sheetName), sheetName: sheetName}
}

// save to local path
func (l *excelExport) ExportToPath(params []map[string]string, data []map[string]interface{}, filePath string) (string, error) {
	l.export(params, data)
	err := l.file.SaveAs(filePath)
	return filePath, err
}

func (l *excelExport) ExportToDir(params []map[string]string, data []map[string]interface{}, dir string) (string, error) {
	filePath := path.Join(dir, createFileName())
	return l.ExportToPath(params, data, filePath)
}

func (l *excelExport) ExportToWeb(params []map[string]string, data []map[string]interface{}, c *gin.Context) {
	l.export(params, data)
	buffer, _ := l.file.WriteToBuffer()
	c.Header("Content-Type", "application/vnd.ms-excel;charset=utf8")
	c.Header("Content-Disposition", "attachment; filename="+url.QueryEscape(createFileName()))
	_, _ = c.Writer.Write(buffer.Bytes())
}

// set the top line
func (l *excelExport) writeTop(params []map[string]string) {
	topStyle, _ := l.file.NewStyle(defaultTopStyle)
	word := startColumn

	for _, conf := range params {
		// set the content of the top cell
		title := conf["title"]
		width, _ := strconv.ParseFloat(conf["width"], 64)
		line := fmt.Sprintf("%c1", word)

		_ = l.file.SetCellValue(l.sheetName, line, title)
		_ = l.file.SetCellStyle(l.sheetName, line, line, topStyle)
		_ = l.file.SetColWidth(l.sheetName, fmt.Sprintf("%c", word), fmt.Sprintf("%c", word), width)

		word++
	}
}

func (l *excelExport) writeData(params []map[string]string, data []map[string]interface{}) {
	lineStyle, _ := l.file.NewStyle(defaultLineStyle)
	// write from the second line
	var j = 2

	for i, val := range data {
		_ = l.file.SetRowHeight(l.sheetName, i+1, defaultHeight)

		// write column by column
		word := startColumn
		for _, conf := range params {
			valKey := conf["key"]
			line := fmt.Sprintf("%c%v", word, j)

			_ = l.file.SetCellValue(l.sheetName, line, fmt.Sprintf("%v", val[valKey]))
			_ = l.file.SetCellStyle(l.sheetName, line, line, lineStyle)
			word++
		}
		j++
	}

	_ = l.file.SetRowHeight(l.sheetName, len(data)+1, defaultHeight)
}

func (l *excelExport) export(params []map[string]string, data []map[string]interface{}) {
	l.writeTop(params)
	l.writeData(params, data)
}

func (l *excelExport) WriteExcelByStruct(sheetName string,
	titleList []string, data []interface{}) {
	if sheetName == "" {
		sheetName = l.sheetName
	}
	l.initializeColumn(sheetName, titleList)
	l.writeDataRow(sheetName, data)
}

func (l *excelExport) initializeColumn(sheetName string, titleList []string) {
	_ = l.file.SetSheetRow(sheetName, "A1", &titleList)
	_ = l.file.SetRowHeight("Sheet1", 1, 30)

	headStyle := Letter(len(titleList))
	l.endCol = "A"
	if len(titleList) > 0 {
		l.endCol = headStyle[len(titleList)-1]
	}

	_ = l.file.SetColWidth(sheetName, "A", l.endCol, 30)
	topStyle, _ := l.file.NewStyle(defaultTopStyle)
	_ = l.file.SetCellStyle(l.sheetName, "A1", fmt.Sprintf("%s1", l.endCol), topStyle)
}

func (l *excelExport) writeDataRow(sheetName string, data []interface{}) {
	rowNum := 1
	rowStyleID, _ := l.file.NewStyle(defaultCellStyle)

	for _, v := range data {
		t := reflect.TypeOf(v)
		value := reflect.ValueOf(v)
		row := make([]interface {
		}, 0)
		for l := 0; l < t.NumField(); l++ {
			val := value.Field(l).Interface()
			row = append(row, val)
		}

		rowNum++
		_ = l.file.SetSheetRow(sheetName, fmt.Sprintf("A%d", rowNum), &row)
		_ = l.file.SetCellStyle(sheetName, fmt.Sprintf("A%d", rowNum), fmt.Sprintf("%s%d", l.endCol, rowNum), rowStyleID)
	}
}

// ExportExcelByStruct the data should be the struct list
func (l *excelExport) ExportExcelToGin(fileName string, ctx *gin.Context) error {
	ctx.Writer.Header().Set("Content-Type", "application/octet-stream")
	ctx.Writer.Header().Set("Content-Transfer-Encoding", "binary")
	ctx.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Disposition")
	ctx.Writer.Header().Set("Content-Disposition",
		fmt.Sprintf("attachment; filename=%s.xlsx", url.QueryEscape(fileName)))

	return l.file.Write(ctx.Writer)
}

// Letter return a-z
func Letter(length int) []string {
	var str []string
	for i := 0; i < length; i++ {
		str = append(str, string(rune('A'+i)))
	}
	return str
}

func createFile(sheetName string) *excelize.File {
	f := excelize.NewFile()
	index, _ := f.NewSheet(sheetName)
	// set default active sheet
	f.SetActiveSheet(index)
	return f
}

func createFileName() string {
	return fmt.Sprintf("excel-%s.xlsx", GenerateStringUUID())
}
