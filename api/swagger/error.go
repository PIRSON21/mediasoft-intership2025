package swagger

// ErrorResponse represents an error message.
// swagger:response ErrorResponse
type ErrorResponse struct {
	// in: body
	Body struct {
		// Error description
		Error string `json:"error"`
	}
}
