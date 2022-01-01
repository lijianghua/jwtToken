package userService

import "fmt"

const (
	Unknown = iota
	InvalidArgument
	NotFound
)

type UserError struct {
	Type    uint   `json:"-"`
	Code    string `json:"code,omitempty"`
	Message string `json:"message"`
}

func (e *UserError) Error() string {
	return fmt.Sprintf("code=%s message=%s", e.Code, e.Message)
}

type FieldError struct {
	UserError
	Field string `json:"field"`
}

func NewMissingFieldError(field string) *FieldError {
	return &FieldError{
		UserError: UserError{
			Type:    InvalidArgument,
			Code:    "missing_field",
			Message: "Field was not provided",
		},
		Field: field,
	}
}

func (e *FieldError) Error() string {
	return fmt.Sprintf("code=%s field=%s message=%s", e.Code, e.Field, e.Message)
}

var ErrNotFound = &UserError{Type: NotFound, Code: "not_found", Message: "user not found"}
