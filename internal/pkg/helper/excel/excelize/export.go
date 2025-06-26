package excelize

import (
	"github.com/xuri/excelize/v2"

	"gobase/internal/pkg/helper/excel"
)

func (h Excel) ExportMultipleSheet(filename string, sheets []excel.Sheet) error {

	f := excelize.NewFile()

	if len(sheets) <= 0 {
		return ErrDataSheetIsRequired
	}
	for _, sheet := range sheets {
		if len(sheet.Data) <= 0 {
			return ErrDataSheetIsRequired
		}
	}

	for _, sheet := range sheets {
		err := generateSheet(f, sheet)
		if err != nil {
			return err
		}
	}

	err := f.DeleteSheet("Sheet1")
	if err != nil {
		return err
	}
	err = f.SaveAs(filename)
	if err != nil {
		return err
	}

	return nil
}

func (h Excel) Export(filename string, sheet excel.Sheet) error {
	f := excelize.NewFile()

	if len(sheet.Data) <= 0 {
		return ErrDataSheetIsRequired
	}

	err := generateSheet(f, sheet)
	if err != nil {
		return err
	}

	err = f.DeleteSheet("Sheet1")
	if err != nil {
		return err
	}
	err = f.SaveAs(filename)
	if err != nil {
		return err
	}

	return nil
}

func generateSheet(f *excelize.File, sheet excel.Sheet) error {
	_, err := f.NewSheet(sheet.Name)
	if err != nil {
		return err
	}
	coor := excel.NewCellCoordinate()

	for _, data := range sheet.Data {

		if data.Title != nil {
			err := setCellStr(f, sheet.Name, &coor, *data.Title)
			if err != nil {
				return err
			}

			coor.IncRow()
		}

		err := prepareTableHeader(f, sheet.Name, &coor, data)
		if err != nil {
			return err
		}

		err = prepareTableBody(f, sheet.Name, &coor, data)
		if err != nil {
			return err
		}

		// TODO: handle footer
	}

	return nil
}

func prepareTableHeader(f *excelize.File, sheetName string, c *excel.CellCoordinate, sheetData excel.TableUnit) error {
	if len(sheetData.Header) <= 0 {
		return nil
	}

	if len(sheetData.Header) != len(sheetData.ColumnOrder) {
		return ErrHeaderAndColumnOrderNotValid
	}

	for _, hKey := range sheetData.ColumnOrder {
		h := sheetData.Header[hKey]

		err := prepareColumnHeader(f, sheetName, c, h)
		if err != nil {
			return err
		}

		c.IncCol()
	}

	c.ResetCol(1)
	c.IncRow()

	return nil
}

func prepareColumnHeader(f *excelize.File, sheetName string, c *excel.CellCoordinate, h excel.HeaderCell) (err error) {

	err = setColWidth(f, sheetName, c, h.Width)
	if err != nil {
		return err
	}

	err = setCellStyle(f, sheetName, c, h.CellStyle)
	if err != nil {
		return err
	}

	err = setCellStr(f, sheetName, c, h.Name)
	if err != nil {
		return err
	}

	return err
}

func prepareTableBody(f *excelize.File, sheetName string, c *excel.CellCoordinate, sheetData excel.TableUnit) error {
	if len(sheetData.Body) <= 0 {
		return nil
	}

	for _, row := range sheetData.Body {

		for _, hKey := range sheetData.ColumnOrder {
			data := row[hKey]

			err := setCellValue(f, sheetName, c, data.Value, data.Type)
			if err != nil {
				return err
			}

			c.IncCol()
		}

		c.ResetCol(1)
		c.IncRow()
	}

	return nil
}

func coordinatesToCellNames(c *excel.CellCoordinate) (string, error) {
	return excelize.CoordinatesToCellName(c.Col, c.Row)
}

func columnNumberToName(c *excel.CellCoordinate) (string, error) {
	return excelize.ColumnNumberToName(c.Col)
}

func setCellValue(f *excelize.File, sheetName string, c *excel.CellCoordinate, value interface{}, t string) error {
	if t == "string" {
		return setCellStr(f, sheetName, c, value.(string))
	}

	return setCellUnknownType(f, sheetName, c, value)
}

func setCellUnknownType(f *excelize.File, sheetName string, c *excel.CellCoordinate, value interface{}) error {
	axis, err := coordinatesToCellNames(c)
	if err != nil {
		return err
	}

	return f.SetCellValue(sheetName, axis, value)
}

func setCellStr(f *excelize.File, sheetName string, c *excel.CellCoordinate, value string) error {
	axis, err := coordinatesToCellNames(c)
	if err != nil {
		return err
	}

	return f.SetCellStr(sheetName, axis, value)
}

func setCellStyle(f *excelize.File, sheetName string, c *excel.CellCoordinate, cs *excel.CellStyle) error {
	if cs == nil {
		return nil
	}

	fontSize := float64(10)

	if cs.FontSize > 0 {
		fontSize = cs.FontSize
	}

	// s := fmt.Sprintf(`{
	// 	"font": {
	// 		"bold": %t,
	// 		"italic": %t,
	// 		"size": %d,
	// 	}
	// }`, cs.Bold, cs.Italic, fontSize)
	s := excelize.Style{
		Font: &excelize.Font{
			Bold:   cs.Bold,
			Italic: cs.Italic,
			Size:   fontSize,
		},
	}

	style, err := f.NewStyle(&s)
	if err != nil {
		return err
	}

	axis, err := coordinatesToCellNames(c)
	if err != nil {
		return err
	}

	err = f.SetCellStyle(sheetName, axis, axis, style)
	if err != nil {
		return err
	}

	return nil
}

func setColWidth(f *excelize.File, sheetName string, c *excel.CellCoordinate, width float64) error {
	colName, err := columnNumberToName(c)
	if err != nil {
		return err
	}

	w := float64(10)
	if width > 0 {
		w = width
	}

	return f.SetColWidth(sheetName, colName, colName, w)
}
