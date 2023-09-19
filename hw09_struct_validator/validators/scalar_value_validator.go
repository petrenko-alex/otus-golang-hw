package validators

import "github.com/petrenko-alex/otus-golang-hw/hw09_struct_validator/rules"

type ScalarValueValidator struct {
	rules rules.ValidationRules
}

func (s ScalarValueValidator) GetValidatorRules() rules.ValidationRules {
	return s.rules
}

func (s ScalarValueValidator) ValidateValue(value interface{}) []error {
	validationErrors := make([]error, 0)
	for _, rule := range s.GetValidatorRules() {
		err := rule.Validate(value)
		if err != nil {
			validationErrors = append(validationErrors, err)
		}
	}

	return validationErrors
}
