package rules

import "golang.org/x/exp/slices"

type StringRangeRule struct{ ValidationLimit }

func (r StringRangeRule) Validate(value interface{}) error {
	valueString, valueCastOk := value.(string)
	if valueCastOk != true {
		return ErrCastValueForRule
	}

	limitSlice, limitCastOk := r.GetLimit().([]string)
	if !limitCastOk {
		return ErrCastLimitForRule
	}

	if !slices.Contains(limitSlice, valueString) {
		return r.GetError()
	}

	return nil
}

func (r StringRangeRule) GetLimit() ValidationLimit {
	return r.ValidationLimit
}

func (r StringRangeRule) GetError() error {
	return ErrValidationFailed
}
