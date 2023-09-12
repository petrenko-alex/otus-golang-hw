package hw09structvalidator

import (
	"errors"
	"reflect"
)

const (
	ValidatorTagName = "validate"
)

var (
	ErrInputNotStruct = errors.New("input argument must be a struct")
	ErrValidatorInit  = errors.New("incorrect validate tag value")
)

type (
	ValidationRule = string
)

type Validator interface {
	ValidateValue(value interface{}) ValidationErrors
}

type ValidationError struct {
	Field string
	Err   error
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	panic("implement me")
}

func Validate(v interface{}) error {
	inputType := reflect.TypeOf(v)
	if inputType.Kind() != reflect.Struct {
		return ErrInputNotStruct
	}

	inputValue := reflect.ValueOf(v)

	validationErrors := make(ValidationErrors, 0)
	for i := 0; i < inputType.NumField(); i++ {
		fieldType := inputType.Field(i)
		val := fieldType.Tag.Get(ValidatorTagName)

		validator, err := GetValidator(fieldType.Type.Kind(), val)
		if err != nil {
			return ErrValidatorInit
		}

		fieldValue := inputValue.FieldByName(fieldType.Name)
		validationErrors = append(
			validationErrors,
			validator.ValidateValue(fieldValue.Interface())...,
		)

	}

	return nil
}

func GetValidator(kind reflect.Kind, rule ValidationRule) (Validator, error) {
	return MockValidator{}, nil
}
