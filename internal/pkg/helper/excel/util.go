package excel

type CellCoordinate struct {
	Col int
	Row int
}

func NewCellCoordinate() CellCoordinate {
	return CellCoordinate{
		Col: 1,
		Row: 1,
	}
}

func (c *CellCoordinate) IncRow() {
	c.Row++
}

func (c *CellCoordinate) IncCol() {
	c.Col++
}

func (c *CellCoordinate) ResetCol(newCol int) {
	if newCol <= 1 {
		c.Col = 1
	} else {
		c.Col = newCol
	}
}
