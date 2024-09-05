package excel

import (
	"bytes"

	"github.com/tealeg/xlsx"
)

// ExportToExcel 中顺序必须与结构体字段顺序一致，
// 且不能跳过结构体中的字段，只可顺序减少结构体后面的字段
func ExportToExcel(exportFields []string, list [][]interface{}) (*bytes.Buffer, error) {
	// 设置默认字体
	xlsx.SetDefaultFont(12, "宋体")
	file := xlsx.NewFile()
	sheet, err := file.AddSheet("Sheet1")
	if err != nil {
		return nil, err
	}
	// 设置表格头
	row := sheet.AddRow()
	for _, field := range exportFields {
		cell := row.AddCell()
		cell.Value = field
		// 设置字体加粗
		cell.GetStyle().Font.Bold = true
	}

	// 设置数据
	for _, e := range list {
		row := sheet.AddRow()
		row.WriteSlice(&e, len(exportFields))
	}

	var buffer bytes.Buffer
	if err = file.Write(&buffer); err != nil {
		return nil, err
	}
	return &buffer, nil
}
