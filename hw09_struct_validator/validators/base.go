package validators

import "github.com/petrenko-alex/otus-golang-hw/hw09_struct_validator/rules"

type ValueValidator interface {
	// ValidateValue validates value against ValidatorRules using various rule combining logic.
	ValidateValue(value interface{}) []error

	// GetValidatorRules returns rules assigned to validator.
	GetValidatorRules() rules.ValidationRules
}
