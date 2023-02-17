package validators

import (
	"fiber-api/utils/types"

	"github.com/go-playground/validator/v10"
)

func IsValidParams(params interface{}) []*types.ValidationError {
	var Validator = validator.New()
	var errors []*types.ValidationError
	err := Validator.Struct(params)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var e types.ValidationError
			e.Field = err.Field()
			e.Tag = err.Tag()
			e.Value = err.Param()
			errors = append(errors, &e)
		}
		return errors
	}
	return nil
}
