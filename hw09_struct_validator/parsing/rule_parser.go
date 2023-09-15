package parsing

import (
	"github.com/petrenko-alex/otus-golang-hw/hw09_struct_validator/rules"
	"strings"
)

const (
	ValidationRuleSeparator = ":"
	ValidationRulePartCount = 2
)

type RuleParser interface {
	GetRule(string) (rules.ValidationRule, error)
}

type BaseRuleParser struct{}

func (f BaseRuleParser) GetRule(stringRule string) (rules.ValidationRule, error) {
	rulePair := strings.Split(stringRule, ValidationRuleSeparator)
	if len(rulePair) != ValidationRulePartCount {
		return nil, rules.ErrParsingRule
	}

	criteria := rulePair[0]
	limit := rulePair[1]
	if len(criteria) == 0 || len(limit) == 0 {
		return nil, rules.ErrParsingRule
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
		return nil, rules.ErrUnknownRule
	}
}
