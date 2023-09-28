package validators

import (
	"github.com/petrenko-alex/otus-golang-hw/hw09_struct_validator/rules"
)

type MockValidator struct{}

func (m MockValidator) ValidateValue(_ interface{}) []error {
	return nil
}

func (m MockValidator) GetValidatorRules() rules.ValidationRules {
	return nil
}
