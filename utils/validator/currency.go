package validators

import (
	"github.com/go-playground/validator/v10"
	"github.com/terajari/bank-api/utils"
)

var ValidCurrency validator.Func = func(fl validator.FieldLevel) bool {
	if currency, ok := fl.Field().Interface().(string); ok {
		return utils.IsSupportedCurrency(currency)
	}
	return false
}
