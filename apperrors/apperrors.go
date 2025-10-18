package apperrors

import "errors"

var ErrOrdersNotFound = errors.New("no orders found matching the criteria")
