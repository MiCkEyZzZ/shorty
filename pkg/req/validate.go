package req

import (
	"github.com/go-playground/validator/v10"
)

// IsValidate проверяет валидность структуры payload.
func IsValidate[T any](payload T) error {
	validate := validator.New()
	err := validate.Struct(payload)
	return err
}
