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
	GetRule(rawRule string) (rules.ValidationRule, error)
}

type BaseRuleParser struct{}

func (f BaseRuleParser) GetRule(rawRule string) (rules.ValidationRule, error) {
	// todo: rawRule to type?

	rulePair := strings.Split(rawRule, ValidationRuleSeparator)
	if len(rulePair) != ValidationRulePartCount {
		return nil, rules.ErrParsingRule
	}

	switch rulePair[0] {
	case "len":
		return rules.LenRule{rulePair[1]}, nil
	case "regexp":
		return rules.RegexpRule{rulePair[1]}, nil
	case "in":
		return rules.InRule{rulePair[1]}, nil
	default:
		return nil, rules.ErrUnknownRule
	}
}
