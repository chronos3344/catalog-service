package entity

import "errors"

var (
	ErrProductAlreadyExists  = errors.New("product already exists")
	ErrCategoryAlreadyExists = errors.New("category already exists")

	ErrNotFound = errors.New("object not found")
	//ErrConflict = errors.New("object already exists")
)
