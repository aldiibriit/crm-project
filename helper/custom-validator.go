package helper

import (
	"github.com/go-playground/validator/v10"
	"github.com/iancoleman/strcase"
)

type ApiError struct {
	Field string
	Msg   string
}

func CustomValidator(ve validator.ValidationErrors) interface{} {
	out := make([]ApiError, len(ve))
	for i, fe := range ve {
		out[i] = ApiError{strcase.ToLowerCamel(fe.Field()), msgForTag(fe.Tag())}
	}
	return out
}

func msgForTag(tag string) string {
	switch tag {
	case "required":
		return "This field is required"
	case "email":
		return "Invalid email"
	}
	return ""
}
