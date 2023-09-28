package rules

import "golang.org/x/exp/slices"

type StringRangeRule struct{ Limit ValidationLimit }

func (r StringRangeRule) Validate(value interface{}) error {
	valueString, valueCastOk := value.(string)
	if !valueCastOk {
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
	return r.Limit
}

func (r StringRangeRule) GetError() error {
	return ErrValidationFailed
}
