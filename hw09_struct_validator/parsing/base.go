package parsing

import "github.com/petrenko-alex/otus-golang-hw/hw09_struct_validator/rules"

type (
	ValidationTag = string
)

type ValidationTagParser interface {
	// GetValidationRules parses ValidationTag and return ValidationRules for it.
	GetValidationRules(tag ValidationTag) (rules.ValidationRules, error)
}
