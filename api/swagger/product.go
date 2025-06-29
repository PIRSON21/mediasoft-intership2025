package swagger

import "github.com/PIRSON21/mediasoft-go/internal/dto"

// ProductResponse represents a response containing product details.
// swagger:response ProductResponse
type ProductResponse struct {
	// A list of products
	// in: body
	Body []dto.ProductAtListResponse
}
