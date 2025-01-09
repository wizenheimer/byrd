package utils

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/creasty/defaults"
	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func InitializeValidator() {
	validate = validator.New()
}

func GetValidator() *validator.Validate {
	// If the validator is not initialized, initialize it
	if validate == nil {
		InitializeValidator()
	}

	return validate
}

func formatValidationError(err error) error {
	var ve validator.ValidationErrors
	if errors.As(err, &ve) {
		var errMsgs []string
		for _, fe := range ve {
			fieldName := fe.Field()
			switch fe.Tag() {
			case "required":
				errMsgs = append(errMsgs, fmt.Sprintf("%s is required", fieldName))
			case "email":
				errMsgs = append(errMsgs, fmt.Sprintf("%s must be a valid email address", fieldName))
			case "oneof":
				errMsgs = append(errMsgs, fmt.Sprintf("%s must be one of [%s]", fieldName, fe.Param()))
			case "min":
				errMsgs = append(errMsgs, fmt.Sprintf("%s must be at least %s", fieldName, fe.Param()))
			case "max":
				errMsgs = append(errMsgs, fmt.Sprintf("%s must not exceed %s", fieldName, fe.Param()))
			default:
				errMsgs = append(errMsgs, fmt.Sprintf("%s failed on validation tag '%s'", fieldName, fe.Tag()))
			}
		}
		return fmt.Errorf("validation failed: %s", strings.Join(errMsgs, "; "))
	}
	return err
}

func validateStruct(s interface{}) error {
	if err := validate.Struct(s); err != nil {
		return formatValidationError(err)
	}
	return nil
}

func setDefaults(s interface{}) error {
	if reflect.ValueOf(s).Kind() != reflect.Ptr {
		return fmt.Errorf("input must be a pointer to a struct")
	}
	return defaults.Set(s)
}

func SetDefaultsAndValidate(s interface{}) error {
	if validate == nil {
		InitializeValidator()
	}

	if err := setDefaults(s); err != nil {
		return fmt.Errorf("failed to set defaults: %w", err)
	}
	return validateStruct(s)
}

func SetDefaultsAndValidateArray(arr interface{}) error {
	if validate == nil {
		InitializeValidator()
	}

	// Handle pointer to slice case
	value := reflect.ValueOf(arr)
	if value.Kind() == reflect.Ptr {
		value = value.Elem()
	}

	if value.Kind() != reflect.Slice && value.Kind() != reflect.Array {
		return fmt.Errorf("input must be a slice or array")
	}

	// Set defaults for each element in the array
	for i := 0; i < value.Len(); i++ {
		// Get a pointer to the actual element in the slice
		elem := value.Index(i).Addr().Interface()

		// Now we're passing a pointer to the actual element
		if err := setDefaults(elem); err != nil {
			return fmt.Errorf("failed to set defaults for element %d: %w", i, err)
		}
	}

	// Create a wrapper struct for validation
	wrapper := struct {
		Array interface{} `validate:"required,dive"`
	}{
		Array: arr,
	}

	return validateStruct(wrapper)
}
