package swagger

// This file contains swagger:route annotations for API endpoints

// swagger:route GET /api/health health healthCheck
// Health check endpoint
//
// responses:
//   200: none

// swagger:route GET /warehouses warehouses listWarehouses
// Returns list of warehouses
//
// responses:
//   200: WarehousesResponse
//   500: ErrorResponse

// swagger:route POST /warehouses warehouses createWarehouse
// Creates a warehouse
//
// responses:
//   201: none
//   409: ErrorResponse
//   422: ErrorResponse

// swagger:route GET /products products getProducts
// Returns list of products
//
// responses:
//   200: ProductResponse
//   500: ErrorResponse

// swagger:route POST /products products addProduct
// Adds a product
//
// responses:
//   201: none
//   409: ErrorResponse
//   422: ErrorResponse

// swagger:route PUT /product/{id} products updateProduct
// Update product information
//
// responses:
//   204: none
//   400: ErrorResponse
//   500: ErrorResponse

// swagger:route PATCH /product/{id} products patchProduct
// Partially update product information
//
// responses:
//   204: none
//   400: ErrorResponse
//   500: ErrorResponse

// swagger:route POST /inventory inventory createInventory
// Create inventory record
//
// responses:
//   201: none
//   400: ErrorResponse
//   409: ErrorResponse
//   500: ErrorResponse

// swagger:route POST /inventory/change_count inventory changeProductCount
// Change product count in warehouse
//
// responses:
//   204: none
//   400: ErrorResponse
//   404: ErrorResponse
//   500: ErrorResponse

// swagger:route POST /inventory/add_discount inventory addDiscount
// Add discount to products
//
// responses:
//   204: none
//   400: ErrorResponse
//   404: ErrorResponse
//   500: ErrorResponse

// swagger:route GET /warehouse/{id} inventory getWarehouseProducts
// Returns products at warehouse or one product if query provided
//
// responses:
//   200: ProductsResponse
//   404: ErrorResponse
//   500: ErrorResponse

// swagger:route POST /inventory/check_cart inventory checkCart
// Calculate cart
//
// responses:
//   200: CartResponse
//   400: ErrorResponse
//   500: ErrorResponse

// swagger:route POST /inventory/buy inventory buyProducts
// Buy products
//
// responses:
//   200: CartResponse
//   400: ErrorResponse
//   500: ErrorResponse

// swagger:route GET /analytics/{id} analytics getWarehouseAnalytics
// Get analytics for warehouse
//
// responses:
//   200: WarehouseAnalyticsResponse
//   500: ErrorResponse

// swagger:route GET /analytics/top_warehouses analytics getTopWarehouses
// Get top warehouses
//
// responses:
//   200: WarehouseAnalyticsAtListResponse
//   500: ErrorResponse
