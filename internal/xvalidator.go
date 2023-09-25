package internal

import "github.com/go-playground/validator/v10"

type XValidator struct {
	validate *validator.Validate
}

type ErrorResponse struct {
	Message     string `json:"message"`
	FailedField string `json:"field"`
}

func NewXValidator() *XValidator {
	return &XValidator{
		validate: validator.New(),
	}
}

func (v *XValidator) Validate(data interface{}) []ErrorResponse {
	var validationErrors []ErrorResponse

	errs := v.validate.Struct(data)
	if errs != nil {
		for _, err := range errs.(validator.ValidationErrors) {
			var elem ErrorResponse
			elem.FailedField = err.Tag()
			elem.Message = err.Error()
			validationErrors = append(validationErrors, elem)
		}
	}
	return validationErrors
}
