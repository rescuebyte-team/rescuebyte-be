package validator

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enTranslations "github.com/go-playground/validator/v10/translations/en"
	"github.com/goccy/go-json"

	"github.com/rescuebyte/rescuebyte-be/pkg/log"
)

type Service interface {
	Validate(data interface{}) MultiError
}

type validatorImpl struct {
	validator *validator.Validate
	trans     ut.Translator
}

var Validator = getValidator()

func getValidator() Service {
	enLocale := en.New()
	uni := ut.New(enLocale, enLocale)

	trans, found := uni.GetTranslator("en")
	if !found {
		log.Error(log.Fields{
			"error": "translator not found",
		}, "[Validator] Translator not found")
	}

	v := validator.New()

	err := enTranslations.RegisterDefaultTranslations(v, trans)
	if err != nil {
		log.Error(log.Fields{
			"error": err.Error(),
		}, "[Validator] Failed to register default translations")
	}

	return &validatorImpl{
		validator: v,
		trans:     trans,
	}
}

func (vi *validatorImpl) Validate(data interface{}) MultiError {
	res := MultiError{}

	err := vi.validator.Struct(data)
	if err != nil {
		var valErrs validator.ValidationErrors

		if errors.As(err, &valErrs) {
			dataType := reflect.TypeOf(data)

			if dataType.Kind() == reflect.Ptr {
				dataType = dataType.Elem()
			}

			for _, valErr := range valErrs {
				field, _ := dataType.FieldByName(valErr.StructField())
				tag, fieldName := getTagAndFieldName(field)

				res.Add(CustomValidationError{
					Scope:   getScope(tag),
					Field:   fieldName,
					Tag:     valErr.Tag(),
					Message: valErr.Translate(vi.trans),
				})
			}
		} else {
			log.Error(log.Fields{
				"error": err.Error(),
			}, "[Validator] Failed to validate data")
		}
	}

	if len(res) > 0 {
		return res
	}

	return nil
}

type CustomValidationError struct {
	Scope   string
	Field   string
	Tag     string
	Message string
}

type FieldError struct {
	Tag     string `json:"tag"`
	Message string `json:"message"`
}

type ValidationError struct {
	Message string                `json:"message"`
	Fields  map[string]FieldError `json:"fields"`
}

type MultiError map[string]ValidationError

func (v MultiError) Error() string {
	j, err := json.Marshal(v)
	if err != nil {
		return ""
	}
	return string(j)
}

func (v MultiError) Serialize() any {
	return v
}

func (v MultiError) Add(err CustomValidationError) {
	if _, exists := v[err.Scope]; !exists {
		v[err.Scope] = ValidationError{
			Fields: make(map[string]FieldError),
		}
	}

	validationErr := v[err.Scope]
	validationErr.Fields[err.Field] = FieldError{
		Tag:     err.Tag,
		Message: err.Message,
	}
	validationErr.Message = fmt.Sprintf(
		"%d validation error(s) in %s",
		len(validationErr.Fields),
		err.Scope,
	)
	v[err.Scope] = validationErr
}

func getTagAndFieldName(field reflect.StructField) (string, string) {
	checkTags := []string{"json", "query", "param"}
	for _, tag := range checkTags {
		fieldName, ok := field.Tag.Lookup(tag)
		if ok {
			return tag, fieldName
		}
	}
	return "", field.Name
}

func getScope(tag string) string {
	switch tag {
	case "json":
		return "body"
	case "query":
		return "query"
	case "param":
		return "param"
	default:
		return "others"
	}
}
