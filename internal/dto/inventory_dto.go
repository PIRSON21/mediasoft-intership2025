package dto

// InventoryCreateRequest представляет запрос на создание инвентаризации.
type InventoryCreateRequest struct {
	WarehouseID string   `json:"warehouse_id"`
	ProductID   string   `json:"product_id"`
	Count       *int     `json:"product_count"`
	Price       *float64 `json:"product_price"`
}

// ChangeProductCountRequest представляет запрос на изменение количества продукта на складе.
type ChangeProductCountRequest struct {
	WarehouseID string `json:"warehouse_id"`
	ProductID   string `json:"product_id"`
	Count       *int   `json:"product_count"`
}

// DiscountToProductRequest представляет запрос на применение скидок к продуктам на складе.
type DiscountToProductRequest struct {
	WarehouseID string      `json:"warehouse_id"`
	Discounts   []*Discount `json:"discounts"`
}

// Discount представляет скидку на продукт.
type Discount struct {
	ProductID     string `json:"product_id"`
	DiscountValue *int   `json:"discount"`
}

// ProductFromWarehouseResponse представляет продукт на складе с его деталями.
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

// CartRequest представляет запрос на корзину товаров.
type CartRequest struct {
	WarehouseID string                  `json:"warehouse_id"`
	Products    []*ProductInCartRequest `json:"products"`
}

// ProductInCartRequest представляет продукт в корзине с его количеством.
type ProductInCartRequest struct {
	ProductID string `json:"product_id"`
	Count     *int   `json:"product_count"`
}

// CartResponse представляет ответ на запрос корзины товаров.
type CartResponse struct {
	Products                      []*ProductInCartResponse `json:"products"`
	TotalProductPrice             float64                  `json:"total_price"`
	TotalProductPriceWithDiscount float64                  `json:"total_price_with_discount"`
}

// ProductInCartResponse представляет продукт в корзине с его деталями.
type ProductInCartResponse struct {
	ProductID         string  `json:"product_id"`
	Count             int     `json:"product_count"`
	FullPrice         float64 `json:"product_price"`
	PriceWithDiscount float64 `json:"product_price_with_discount"`
}

// Pagination представляет параметры пагинации для запросов.
type Pagination struct {
	Page   int
	Offset int
	Limit  int
}

// ProductsResponse представляет ответ на запрос списка продуктов.
type ProductsResponse struct {
	Page     int              `json:"page"`
	Limit    int              `json:"limit"`
	Products []*ProductAtList `json:"products"`
}

// ProductAtList представляет продукт в списке с его деталями.
type ProductAtList struct {
	ProductID                string  `json:"product_id"`
	ProductName              string  `json:"product_name"`
	ProductPrice             float64 `json:"product_price"`
	ProductPriceWithDiscount float64 `json:"product_discount_price"`
}
