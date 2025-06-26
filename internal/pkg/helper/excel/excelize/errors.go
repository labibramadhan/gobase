package excelize

import "errors"

var (
	// ErrDataSheetIsRequired defined the error message on invalid data sheet missing
	ErrDataSheetIsRequired = errors.New("data sheet is required")
	// ErrHeaderAndColumnOrderNotValid defined the error message on invalid header and column order data
	ErrHeaderAndColumnOrderNotValid = errors.New("header and column order is not valid")
)
