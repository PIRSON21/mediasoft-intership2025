package domain

// Analytics представляет аналитику по складу и продуктам.
type Analytics struct {
	Warehouse    *Warehouse
	Product      *Product
	ProductCount int
	ProductPrice float64
}
