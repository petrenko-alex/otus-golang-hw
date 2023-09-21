package rules

import (
	"reflect"
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
	typeValue := reflect.ValueOf(value)
	if typeValue.Kind() == reflect.Int {
		return IntRangeRule{r.GetLimit()}.Validate(int(typeValue.Int()))
	}

	if typeValue.Kind() == reflect.String {
		return StringRangeRule{r.GetLimit()}.Validate(typeValue.String())
	}

	return ErrCastValueForRule
}

func (r InRule) GetLimit() ValidationLimit {
	return r.ValidationLimit
}

func (r InRule) GetError() error {
	return ErrValidationFailed
}
