package helper

import (
	"github.com/xuri/excelize/v2"
)

func CreateExcelFile(headers []string, fileName string, list [][]interface{}) error {
	f := excelize.NewFile()
	sheet := "Sheet1"
	f.NewSheet(sheet)

	// Menulis header
	for col, header := range headers {
		cell := string('A'+col) + "1"
		f.SetCellValue("Sheet1", cell, header)
	}

	// Menulis data
	// Write data
	for i, row := range list {
		for j, value := range row {
			cell, _ := excelize.CoordinatesToCellName(j+1, i+2) // Start from row 2
			f.SetCellValue(sheet, cell, value)
		}
	}

	if err := f.SaveAs(fileName); err != nil {
		return err
	}

	return nil
}
