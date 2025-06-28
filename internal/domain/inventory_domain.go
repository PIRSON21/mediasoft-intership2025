package domain

type Inventory struct {
	Product      *Product
	Warehouse    *Warehouse
	ProductCount int
	ProductPrice float64
	ProductSale  int
}
