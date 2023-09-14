package validators

import "github.com/petrenko-alex/otus-golang-hw/hw09_struct_validator/rules"

type SliceValueValidator struct {
	rules rules.ValidationRules
	// todo: inject scalarValidator?
}

func (s SliceValueValidator) ValidateValue(value interface{}) []error {
	//TODO implement me
	panic("implement me")
}

func (s SliceValueValidator) GetValidatorRules() rules.ValidationRules {
	return s.rules
}
