package rules

import "errors"

var (
	ErrCastValueForRule = errors.New("can't cast value for validation rule")
	ErrCastLimitForRule = errors.New("can't cast limit for validation rule")
	ErrValidationFailed = errors.New("value does not satisfy limit")
)

type (
	ValidationLimit = interface{}
	ValidationRules = []ValidationRule
)

// ValidationRule validates value against ValidationLimit returning error if not satisfy.
type ValidationRule interface {
	Validate(value interface{}) error
	GetLimit() ValidationLimit
	GetError() error
}
