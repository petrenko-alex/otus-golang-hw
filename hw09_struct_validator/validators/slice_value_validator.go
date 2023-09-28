package validators

import (
	"errors"
	"reflect"

	"github.com/petrenko-alex/otus-golang-hw/hw09_struct_validator/rules"
)

var (
	ErrValueNotIterable  = errors.New("value should be iterable (slice or array)")
	ErrValueNotSupported = errors.New("value type not supported")
)

type SliceValueValidator struct {
	scalarValidator ScalarValueValidator
}

func (v SliceValueValidator) ValidateValue(value interface{}) ([]error, error) {
	validationErrors := make([]error, 0)
	switch reflect.TypeOf(value).Kind() {
	case reflect.Slice, reflect.Array:
		s := reflect.ValueOf(value)

		for i := 0; i < s.Len(); i++ {
			element := s.Index(i)

			var elementValidationErrors []error
			var runtimeError error
			switch element.Kind() {
			case reflect.String:
				elementValidationErrors, runtimeError = v.scalarValidator.ValidateValue(element.String())
			case reflect.Int:
				elementValidationErrors, runtimeError = v.scalarValidator.ValidateValue(int(element.Int()))
			default:
				return nil, ErrValueNotSupported
			}

			if runtimeError != nil {
				return nil, runtimeError
			}

			validationErrors = append(validationErrors, elementValidationErrors...)
		}
	default:
		return nil, ErrValueNotIterable
	}

	return validationErrors, nil
}

func (v SliceValueValidator) GetValidatorRules() rules.ValidationRules {
	return v.scalarValidator.GetValidatorRules()
}
