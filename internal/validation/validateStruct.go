package validation

import (
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
)

func ValidateStruct(s interface{}) []string {
	validate := validator.New()
	en := en.New()
	uni := ut.New(en, en)
	trans, _ := uni.GetTranslator("en")
	en_translations.RegisterDefaultTranslations(validate, trans)
	err := validate.Struct(s)
	if err != nil {
		errs := err.(validator.ValidationErrors)
		var errorMessages []string
		for _, err := range errs {
			errorMessages = append(errorMessages, err.Translate(trans))
		}
		return errorMessages
	}
	return nil
}
