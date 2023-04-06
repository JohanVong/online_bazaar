package tools

import "github.com/go-playground/validator"

// CustomValidator - кастомный валидатор для приложения
type CustomValidator struct {
	Validator *validator.Validate
}

// Validate - валидирует входящие данные в контроллере по тэгам `validate:""`
func (cv *CustomValidator) Validate(i interface{}) error {
	if err := cv.Validator.Struct(i); err != nil {
		return err
	}

	return nil
}
