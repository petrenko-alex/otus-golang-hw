package rules

import "strconv"

type LenRule struct{ Limit ValidationLimit }

func (r LenRule) Validate(value interface{}) error {
	valueStr, valueCastOk := value.(string)
	if valueCastOk != true {
		return ErrCastValueForRule
	}

	limitInt, limitCastOk := strconv.Atoi(r.GetLimit().(string))
	if limitCastOk != nil {
		return ErrCastLimitForRule
	}

	if len(valueStr) != limitInt {
		return r.GetError()
	}

	return nil
}

func (r LenRule) GetLimit() ValidationLimit {
	return r.Limit
}

func (r LenRule) GetError() error {
	return ErrValidationFailed
}
