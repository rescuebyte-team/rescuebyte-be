package validator

import (
	"reflect"

	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"

	"github.com/rescuebyte/rescuebyte-be/domain/enums"
	"github.com/rescuebyte/rescuebyte-be/pkg/log"
)

type CustomValidator interface {
	CustomValidate() []CustomValidationError
}

type CustomValidationError struct {
	Scope   string
	Field   string
	Tag     string
	Message string
}

func registerCustomRules(v *validator.Validate) {
	for _, desc := range enums.GetEnumRegistry() {
		err := v.RegisterValidation(
			desc.QueryKey,
			enumValidator(),
		)

		if err != nil {
			log.Error(log.Fields{
				"error": err.Error(),
			}, "[Validator] Failed to register "+desc.QueryKey+" validation")
		}
	}
}

func registerCustomTranslations(
	v *validator.Validate,
	trans ut.Translator,
) {
	for _, desc := range enums.GetEnumRegistry() {
		registerEnumTranslation(v, trans, desc.QueryKey)
	}
}

func registerEnumTranslation(
	v *validator.Validate,
	trans ut.Translator,
	tag string,
) {
	err := v.RegisterTranslation(
		tag,
		trans,
		func(ut ut.Translator) error {
			return ut.Add(
				tag,
				"{0} is invalid",
				true,
			)
		},
		func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T(
				tag,
				fe.Field(),
			)

			return t
		},
	)

	if err != nil {
		log.Error(log.Fields{
			"error": err.Error(),
		}, "[Validator] Failed to register "+tag+" translation")
	}
}

func enumValidator() validator.Func {
	return func(fl validator.FieldLevel) bool {
		field := fl.Field()

		if field.Kind() == reflect.Ptr {
			if field.IsNil() {
				return true
			}

			field = field.Elem()
		}

		method := field.MethodByName("IsValidIdx")
		if !method.IsValid() {
			return true
		}

		return method.Call(nil)[0].Bool()
	}
}
