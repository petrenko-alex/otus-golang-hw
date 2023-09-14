package hw09structvalidator

import (
	"errors"
	"github.com/petrenko-alex/otus-golang-hw/hw09_struct_validator/validators"
	"reflect"
)

// todo: словарь терминов
// todo: rename to struct validator?
// todo: add constructors for types

// validation (raw)rule - min:10, len:32
// validation limit - 10, 32
// validation criteria - min, len
// validation tag - ?

const (
	ValidatorTagName = "validate"
)

var (
	ErrInputNotStruct = errors.New("input argument must be a struct")
	ErrValidatorInit  = errors.New("incorrect validate tag value")
)

type ValidationError struct {
	Field string
	Err   error
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	panic("implement me")
}

// Задача StructValidator - проходить по поням структуры и каждое поле валидировать
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

		validator, err := validators.GetValidator(fieldType.Type.Kind(), val)
		if err != nil {
			return ErrValidatorInit
		}

		fieldValue := inputValue.FieldByName(fieldType.Name)
		fieldErrors := validator.ValidateValue(fieldValue.Interface())
		if fieldErrors != nil && len(fieldErrors) > 0 {
			for _, fieldErr := range fieldErrors {
				validationErrors = append(validationErrors, ValidationError{
					Field: fieldType.Name,
					Err:   fieldErr,
				})
			}
		}
	}

	return validationErrors
}
