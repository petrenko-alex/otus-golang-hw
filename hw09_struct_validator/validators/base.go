package validators

import (
	"errors"

	"github.com/petrenko-alex/otus-golang-hw/hw09_struct_validator/rules"
)

var ErrRuntime = errors.New("unexpected error during validation")

type ValueValidator interface {
	// ValidateValue validates value against ValidatorRules using various rule combining logic.
	ValidateValue(value interface{}) ([]error, error)

	// GetValidatorRules returns rules assigned to validator.
	GetValidatorRules() rules.ValidationRules
}
