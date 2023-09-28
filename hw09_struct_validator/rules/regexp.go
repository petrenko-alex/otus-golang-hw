package rules

import "regexp"

type RegexpRule struct{ Limit ValidationLimit }

func (r RegexpRule) Validate(value interface{}) error {
	valueStr, valueCastOk := value.(string)
	if !valueCastOk {
		return ErrCastValueForRule
	}

	limitStr, limitCastOk := r.GetLimit().(string)
	if !limitCastOk {
		return ErrCastLimitForRule
	}

	match, regexpErr := regexp.MatchString(limitStr, valueStr)
	if regexpErr != nil {
		return ErrCastLimitForRule
	}

	if !match {
		return r.GetError()
	}

	return nil
}

func (r RegexpRule) GetLimit() ValidationLimit {
	return r.Limit
}

func (r RegexpRule) GetError() error {
	return ErrValidationFailed
}
