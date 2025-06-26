package errors

import "errors"

var (
	ErrProductAlreadyExists = errors.New("product with this name already exists")
	ErrProductNotFound      = errors.New("product not found")
)
