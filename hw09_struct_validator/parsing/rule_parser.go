package parsing

import (
	"errors"
	"github.com/petrenko-alex/otus-golang-hw/hw09_struct_validator/rules"
	"strings"
)

const (
	ValidationRuleSeparator = ":"
	ValidationRulePartCount = 2
)

var (
	ErrParsingRule = errors.New("incorrect rule string")
	ErrUnknownRule = errors.New("unknown validation rule")
)

type RuleParser interface {
	GetRule(string) (rules.ValidationRule, error)
}

type BaseRuleParser struct{}

func (f BaseRuleParser) GetRule(stringRule string) (rules.ValidationRule, error) {
	rulePair := strings.Split(stringRule, ValidationRuleSeparator)
	if len(rulePair) != ValidationRulePartCount {
		return nil, ErrParsingRule
	}

	criteria := rulePair[0]
	limit := rulePair[1]
	if len(criteria) == 0 || len(limit) == 0 {
		return nil, ErrParsingRule
	}

	switch criteria {
	case "len":
		return rules.LenRule{limit}, nil
	case "regexp":
		return rules.RegexpRule{limit}, nil
	case "in":
		return rules.InRule{limit}, nil
	case "min":
		return rules.MinRule{limit}, nil
	case "max":
		return rules.MaxRule{limit}, nil
	default:
		return nil, ErrUnknownRule
	}
}
