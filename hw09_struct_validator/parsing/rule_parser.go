package parsing

import (
	"errors"
	"strings"

	"github.com/petrenko-alex/otus-golang-hw/hw09_struct_validator/rules"
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
	// GetRule returns one ValidationRule for strings, describing rule (ex. "len:10").
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
		rule, err := rules.NewInRule(limit)
		if err != nil {
			return nil, ErrParsingRule
		}
		return rule, nil
	case "min":
		return rules.MinRule{limit}, nil
	case "max":
		return rules.MaxRule{limit}, nil
	default:
		return nil, ErrUnknownRule
	}
}
