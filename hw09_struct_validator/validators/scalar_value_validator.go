package validators

import (
	"errors"
	"fmt"

	"github.com/petrenko-alex/otus-golang-hw/hw09_struct_validator/rules"
)

type ScalarValueValidator struct {
	rules rules.ValidationRules
}

func (s ScalarValueValidator) GetValidatorRules() rules.ValidationRules {
	return s.rules
}

func (s ScalarValueValidator) ValidateValue(value interface{}) ([]error, error) {
	validationErrors := make([]error, 0)
	for _, rule := range s.GetValidatorRules() {
		err := rule.Validate(value)
		if err != nil {
			if !errors.Is(err, rule.GetError()) {
				return nil, fmt.Errorf(ErrRuntime.Error()+": %w", err)
			}
			validationErrors = append(validationErrors, err)
		}
	}

	return validationErrors, nil
}
