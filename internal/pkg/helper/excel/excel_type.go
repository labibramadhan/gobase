package excel

import "github.com/google/uuid"

// HeaderKey is used as a reference for body to place its value into the same header column. Make sure no duplicate key
type HeaderKey string
type BodyRow map[HeaderKey]BodyCell

type Sheet struct {
	Name string
	Data []TableUnit
}

func NewSheet(name string) Sheet {
	return Sheet{
		Name: name,
		Data: []TableUnit{},
	}
}

type TableUnit struct {
	Title            *string
	Header           TableHeader // Requiring ColumnOrder when structuring table header
	ColumnOrder      []HeaderKey
	IsWidthReference bool      // used as reference point for width in one sheet excel
	Body             TableBody // Requiring ColumnOrder when structuring table body
	Footer           TableFooter
}

type TableHeader map[HeaderKey]HeaderCell
type TableBody []map[HeaderKey]BodyCell
type TableFooter []FooterCell

func NewTableUnit() TableUnit {
	return TableUnit{
		Header:      TableHeader{},
		ColumnOrder: []HeaderKey{},
		Body:        TableBody{},
		Footer:      TableFooter{},
	}
}

func (t *TableUnit) SetHeaderColumn(headers []*HeaderCell) {
	for _, h := range headers {
		key := HeaderKey(uuid.NewString())
		if h != nil {
			key = h.Key
		}

		t.ColumnOrder = append(t.ColumnOrder, key)
		t.addHeaderColumn(string(key), h)
	}
}

// func (t *TableUnit) setColumnOrder(keys []string) {
// 	for _, k := range keys {
// 		t.ColumnOrder = append(t.ColumnOrder, HeaderKey(k))
// 	}
// }

func (t *TableUnit) addHeaderColumn(key string, h *HeaderCell) {
	hKey := HeaderKey(key)

	if h == nil {
		t.Header[hKey] = *NewHeaderCell("", string(hKey))
	} else {
		t.Header[hKey] = *h
	}
}

func (t *TableUnit) AddBodyRow(data map[string]BodyCell) {
	row := make(map[HeaderKey]BodyCell)

	for _, hKey := range t.ColumnOrder {
		row[hKey] = data[string(hKey)]
	}

	t.Body = append(t.Body, row)
}

func NewTableFooter() TableFooter {
	return TableFooter{}
}

type CellStyle struct {
	FontSize float64
	Italic   bool
	Bold     bool
}

type HeaderCell struct {
	Name      string
	Key       HeaderKey
	Width     float64
	BodyStyle *CellStyle
	*CellStyle
}

func SetHeaderCellWidth(width float64) func(*HeaderCell) {
	return func(hc *HeaderCell) {
		hc.Width = width
	}
}

func SetHeaderCellStyle(style CellStyle) func(*HeaderCell) {
	return func(hc *HeaderCell) {
		hc.CellStyle = &style
	}
}

func SetHeaderCellBodyStyle(style CellStyle) func(*HeaderCell) {
	return func(hc *HeaderCell) {
		hc.BodyStyle = &style
	}
}

func NewHeaderCell(name, key string, options ...func(*HeaderCell)) *HeaderCell {
	hc := &HeaderCell{
		Name: name,
		Key:  HeaderKey(key),
	}
	// ... (write initializations with default values)...
	for _, option := range options {
		option(hc)
	}
	return hc
}

type BodyCell struct {
	Value interface{}
	Type  string
}

func NewBodyCell(val interface{}, t string) BodyCell {
	return BodyCell{
		Value: val,
		Type:  t,
	}
}

type FooterCell struct {
	Name       string
	Key        string
	MergeNCell int // will merge current cell with N right cell
	*CellStyle
}
