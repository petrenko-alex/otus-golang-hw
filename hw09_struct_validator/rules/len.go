package rules

type LenRule struct{ ValidationLimit }

func (r LenRule) Validate(value interface{}) error {
	valueStr, valueCastOk := value.(string)
	if valueCastOk != true {
		return ErrCastValueForRule
	}

	limitInt, limitCastOk := r.GetLimit().(int)
	if limitCastOk != true {
		return ErrCastLimitForRule
	}

	if len(valueStr) != limitInt {
		return r.GetError()
	}

	return nil
}

func (r LenRule) GetLimit() ValidationLimit {
	return r.ValidationLimit
}

func (r LenRule) GetError() error {
	return ErrValidationFailed
}
