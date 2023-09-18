package rules

import (
	"strings"
)

type InRule struct{ ValidationLimit }

func NewInRule(limit ValidationLimit) (*InRule, error) {
	strLimit, strCastOk := limit.(string)
	if !strCastOk {
		return nil, ErrCastLimitForRule
	}

	return &InRule{strings.Split(strLimit, ",")}, nil
}

func (r InRule) Validate(value interface{}) error {
	valueInt, intCastOk := value.(int)
	if intCastOk {
		return IntRangeRule{r.GetLimit()}.Validate(valueInt)
	}

	valueStr, strCastOk := value.(string)
	if strCastOk {
		return StringRangeRule{r.GetLimit()}.Validate(valueStr)
	}

	return ErrCastValueForRule
}

func (r InRule) GetLimit() ValidationLimit {
	return r.ValidationLimit
}

func (r InRule) GetError() error {
	return ErrValidationFailed
}
