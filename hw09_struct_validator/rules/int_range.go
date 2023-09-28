package rules

import "strconv"

type IntRangeRule struct{ Limit ValidationLimit }

func (r IntRangeRule) Validate(value interface{}) error {
	valueInt, valueCastOk := value.(int)
	if !valueCastOk {
		return ErrCastValueForRule
	}

	limitSlice, limitCastOk := r.castLimit()
	if !limitCastOk {
		return ErrCastLimitForRule
	}

	leftRange := limitSlice[0]
	rightRange := limitSlice[1]

	if valueInt < leftRange || valueInt > rightRange {
		return r.GetError()
	}

	return nil
}

func (r IntRangeRule) GetLimit() ValidationLimit {
	return r.Limit
}

func (r IntRangeRule) GetError() error {
	return ErrValidationFailed
}

func (r IntRangeRule) castLimit() ([]int, bool) {
	// try cast using type casting
	intLimitSlice, intCastOk := r.GetLimit().([]int)
	if intCastOk {
		return intLimitSlice, true
	}

	// try cast using strconv
	limitSlice, stringCastOk := r.GetLimit().([]string)
	if !stringCastOk {
		return nil, false
	}

	intLimitSlice = make([]int, 0, len(limitSlice))
	for _, val := range limitSlice {
		valInt, err := strconv.Atoi(val)
		if err != nil {
			return nil, false
		}

		intLimitSlice = append(intLimitSlice, valInt)
	}

	return intLimitSlice, true
}
