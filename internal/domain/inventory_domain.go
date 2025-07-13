package domain

// Inventory представляет информацию о продукте на складе.
type Inventory struct {
	Product      *Product
	Warehouse    *Warehouse
	ProductCount int
	ProductPrice float64
	ProductSale  int
}
