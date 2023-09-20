package validators

import (
	"errors"
	"github.com/petrenko-alex/otus-golang-hw/hw09_struct_validator/rules"
	"reflect"
)

var (
	ErrValueNotIterable  = errors.New("value should be iterable (slice or array)")
	ErrValueNotSupported = errors.New("value type not supported")
)

type SliceValueValidator struct {
	scalarValidator ScalarValueValidator
}

func (v SliceValueValidator) ValidateValue(value interface{}) []error {
	validationErrors := make([]error, 0)
	switch reflect.TypeOf(value).Kind() {
	case reflect.Slice, reflect.Array:
		s := reflect.ValueOf(value)

		for i := 0; i < s.Len(); i++ {
			element := s.Index(i)

			var elementValidationErrors []error
			switch element.Kind() {
			case reflect.String:
				elementValidationErrors = v.scalarValidator.ValidateValue(element.String())
			case reflect.Int:
				elementValidationErrors = v.scalarValidator.ValidateValue(int(element.Int()))
			default:
				return []error{ErrValueNotSupported}
			}

			validationErrors = append(validationErrors, elementValidationErrors...)
		}
	default:
		return []error{ErrValueNotIterable}
	}

	return validationErrors
}

func (v SliceValueValidator) GetValidatorRules() rules.ValidationRules {
	return v.scalarValidator.GetValidatorRules()
}
