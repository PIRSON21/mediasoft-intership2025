package errors

import "errors"

var (
	ErrInventoryAlreadyExists = errors.New("this inventory information already exists")
	ErrForeignKey             = errors.New("wrong product ID or warehouse ID")
)
