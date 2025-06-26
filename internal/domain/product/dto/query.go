package productdto

import (
	"gobase/internal/pkg/service/crud"
)

type ProductList = crud.PageResult[*Product]
