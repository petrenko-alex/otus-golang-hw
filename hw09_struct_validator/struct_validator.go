package hw09structvalidator

import (
	"errors"
	"reflect"

	"github.com/petrenko-alex/otus-golang-hw/hw09_struct_validator/validators" //nolint:depguard
)

const (
	ValidatorTagName = "validate"
)

var (
	ErrInputNotStruct = errors.New("input argument must be a struct")
	ErrValidatorInit  = errors.New("incorrect validate tag value")
)

type StructValidator struct {
	factory validators.ValidatorFactory
}

func (v StructValidator) Validate(value interface{}) error {
	inputType := reflect.TypeOf(value)
	if inputType.Kind() != reflect.Struct {
		return ErrInputNotStruct
	}

	inputValue := reflect.ValueOf(value)

	validationErrors := make(ValidationErrors, 0)
	for i := 0; i < inputType.NumField(); i++ {
		fieldType := inputType.Field(i)
		val := fieldType.Tag.Get(ValidatorTagName)
		if len(val) == 0 {
			continue
		}

		validator, err := v.factory.GetValidator(fieldType.Type.Kind(), val)
		if err != nil {
			return ErrValidatorInit
		}

		fieldValue := inputValue.FieldByName(fieldType.Name)
		if !fieldValue.CanInterface() {
			continue // skip validation of unexported fields
		}

		fieldErrors, runtimeError := validator.ValidateValue(fieldValue.Interface())
		if runtimeError != nil {
			return runtimeError
		}

		if len(fieldErrors) > 0 {
			for _, fieldErr := range fieldErrors {
				validationErrors = append(validationErrors, ValidationError{
					Field: fieldType.Name,
					Err:   fieldErr,
				})
			}
		}
	}

	if len(validationErrors) == 0 {
		return nil
	}

	return validationErrors
}
