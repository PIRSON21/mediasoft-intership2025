package errors

import "errors"

var (
	ErrInventoryAlreadyExists     = errors.New("this inventory information already exists")
	ErrForeignKey                 = errors.New("wrong product ID or warehouse ID")
	ErrInventoryNotFound          = errors.New("inventory not found")
	ErrNotEnoughProductCount      = errors.New("there are not enough products at warehouse")
	ErrNotFoundProductAtWarehouse = errors.New("there are not some products at warehouse")
)
