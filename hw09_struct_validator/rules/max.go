package rules

import "strconv"

type MaxRule struct{ ValidationLimit }

func (r MaxRule) Validate(value interface{}) error {
	valueInt, valueCastOk := value.(int)
	if valueCastOk != true {
		return ErrCastValueForRule
	}

	limitInt, limitCastOk := strconv.Atoi(r.GetLimit().(string))
	if limitCastOk != nil {
		return ErrCastLimitForRule
	}

	if valueInt > limitInt {
		return r.GetError()
	}

	return nil
}

func (r MaxRule) GetLimit() ValidationLimit {
	return r.ValidationLimit
}

func (r MaxRule) GetError() error {
	return ErrValidationFailed
}
