package swagger

import (
	"github.com/PIRSON21/mediasoft-intership2025/internal/domain"
	"github.com/PIRSON21/mediasoft-intership2025/internal/dto"
)

// Domain models

// swagger:model Product
type Product domain.Product

// swagger:model Warehouse
type Warehouse domain.Warehouse

// swagger:model Inventory
type Inventory domain.Inventory

// swagger:model Analytics
type Analytics domain.Analytics

// DTO models

// swagger:model ProductAtListResponse
type ProductAtListResponse dto.ProductAtListResponse

// swagger:model ProductRequest
type ProductRequest dto.ProductRequest

// swagger:model WarehouseRequest
type WarehouseRequest dto.WarehouseRequest

// swagger:model WarehouseAtListResponse
type WarehouseAtListResponse dto.WarehouseAtListResponse

// swagger:model InventoryCreateRequest
type InventoryCreateRequest dto.InventoryCreateRequest

// swagger:model ChangeProductCountRequest
type ChangeProductCountRequest dto.ChangeProductCountRequest

// swagger:model DiscountToProductRequest
type DiscountToProductRequest dto.DiscountToProductRequest

// swagger:model Discount
type Discount dto.Discount

// swagger:model ProductFromWarehouseResponse
type ProductFromWarehouseResponse dto.ProductFromWarehouseResponse

// swagger:model CartRequest
type CartRequest dto.CartRequest

// swagger:model ProductInCartRequest
type ProductInCartRequest dto.ProductInCartRequest

// swagger:model CartResponse
type CartResponse dto.CartResponse

// swagger:model ProductInCartResponse
type ProductInCartResponse dto.ProductInCartResponse

// swagger:model Pagination
type Pagination dto.Pagination

// swagger:model ProductsResponse
type ProductsResponse dto.ProductsResponse

// swagger:model ProductAtList
type ProductAtList dto.ProductAtList

// swagger:model WarehouseAnalyticsResponse
type WarehouseAnalyticsResponse dto.WarehouseAnalyticsResponse

// swagger:model ProductAnalytic
type ProductAnalytic dto.ProductAnalytic

// swagger:model WarehouseAnalyticsAtListResponse
type WarehouseAnalyticsAtListResponse dto.WarehouseAnalyticsAtListResponse
