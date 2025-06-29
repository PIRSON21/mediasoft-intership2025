package dto

type InventoryCreateRequest struct {
	WarehouseID string   `json:"warehouse_id"`
	ProductID   string   `json:"product_id"`
	Count       *int     `json:"product_count"`
	Price       *float64 `json:"product_price"`
}

type ChangeProductCountRequest struct {
	WarehouseID string `json:"warehouse_id"`
	ProductID   string `json:"product_id"`
	Count       *int   `json:"product_count"`
}

type DiscountToProductRequest struct {
	WarehouseID string      `json:"warehouse_id"`
	Discounts   []*Discount `json:"discounts"`
}

type Discount struct {
	ProductID     string `json:"product_id"`
	DiscountValue *int   `json:"discount"`
}

type ProductFromWarehouseResponse struct {
	ProductID            string         `json:"product_id"`
	ProductName          string         `json:"product_name"`
	ProductDescription   string         `json:"product_description"`
	ProductWeight        float64        `json:"product_weight"`
	ProductParams        map[string]any `json:"product_params,omitempty"`
	ProductBarcode       string         `json:"product_barcode"`
	ProductCount         int            `json:"product_count"`
	ProductPrice         float64        `json:"product_price"`
	ProductPriceWithSale float64        `json:"product_sale"`
}

type CartRequest struct {
	WarehouseID string                  `json:"warehouse_id"`
	Products    []*ProductInCartRequest `json:"products"`
}

type ProductInCartRequest struct {
	ProductID string `json:"product_id"`
	Count     *int   `json:"product_count"`
}

type CartResponse struct {
	Products                      []*ProductInCartResponse `json:"products"`
	TotalProductPrice             float64                  `json:"total_price"`
	TotalProductPriceWithDiscount float64                  `json:"total_price_with_discount"`
}

type ProductInCartResponse struct {
	ProductID         string  `json:"product_id"`
	Count             int     `json:"product_count"`
	FullPrice         float64 `json:"product_price"`
	PriceWithDiscount float64 `json:"product_price_with_discount"`
}
