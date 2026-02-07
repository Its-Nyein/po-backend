package validation

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

var Validate = validator.New(validator.WithRequiredStructEnabled())

func ValidateStruct(data interface{}) error {
	return Validate.Struct(data)
}

func FormatValidationErrors(err error) map[string]string {
	formatted := make(map[string]string)
	if errs, ok := err.(validator.ValidationErrors); ok {
		for _, e := range errs {
			field := e.Field()
			switch e.Tag() {
			case "required":
				formatted[field] = fmt.Sprintf("%s is required", field)
			case "min":
				formatted[field] = fmt.Sprintf("%s must be at least %s characters", field, e.Param())
			default:
				formatted[field] = fmt.Sprintf("%s is not valid", field)
			}
		}
	}
	return formatted
}
