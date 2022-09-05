package api

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]

		if name == "-" {
			return ""
		}

		return name
	})

	// This is also other way to get the json tag from field
	// validationErrors := err.(validator.ValidationErrors)
	// validationErr := validationErrors[0]
	// fieldName := validationErr.Field()
	// field, ok := reflect.TypeOf(args).Elem().FieldByName(fieldName)
	// fieldJSONName, _ := field.Tag.Lookup("json")
}

func msgForTag(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "This field is required"

	case "alphanum":
		return "This field can have only alphanumeric characters"

	case "min":
		return fmt.Sprintf("This field must have at least %v characters", fe.Param())

	case "max":
		return fmt.Sprintf("This field can have at most %v characters", fe.Param())

	case "email":
		return "The email address is invalid"

	default:
		return fe.Error()
	}
}

func validateRequest(req interface{}) map[string][]string {
	err := validate.Struct(req)

	if err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			fmt.Println(err)
			return map[string][]string{
				"invalid_validation_error": {"Invalid validation error"},
			}
		}

		validationErrors := err.(validator.ValidationErrors)
		var apiErrors = map[string][]string{}
		for _, fe := range validationErrors {
			apiErrors[fe.Field()] = append(apiErrors[fe.Field()], msgForTag(fe))
		}

		return apiErrors
	}

	return nil
}
