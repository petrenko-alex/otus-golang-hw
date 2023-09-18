package rules

import "strconv"

type MinRule struct{ ValidationLimit }

func (r MinRule) Validate(value interface{}) error {
	valueInt, valueCastOk := value.(int)
	if valueCastOk != true {
		return ErrCastValueForRule
	}

	limitInt, limitCastOk := strconv.Atoi(r.GetLimit().(string))
	if limitCastOk != nil {
		return ErrCastLimitForRule
	}

	if valueInt < limitInt {
		return r.GetError()
	}

	return nil
}

func (r MinRule) GetLimit() ValidationLimit {
	return r.ValidationLimit
}

func (r MinRule) GetError() error {
	return ErrValidationFailed
}
