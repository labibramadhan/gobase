package excel

type Excel interface {
	Export(filename string, data Sheet) error
	ExportMultipleSheet(filename string, sheets []Sheet) error
}
