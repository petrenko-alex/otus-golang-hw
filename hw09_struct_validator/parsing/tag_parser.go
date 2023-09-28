package parsing

import (
	"errors"
	"fmt"
	"strings"

	"github.com/petrenko-alex/otus-golang-hw/hw09_struct_validator/rules"
)

const (
	ValidationTagSeparator = "|"
)

var (
	ErrValidationTagEmpty   = errors.New("no validation rules provided")
	ErrParsingValidationTag = errors.New("validation tag corrupted")
)

type TagParser struct {
	RuleParser RuleParser
}

func (t TagParser) GetValidationRules(tag ValidationTag) (rules.ValidationRules, error) {
	if len(tag) == 0 {
		return nil, ErrValidationTagEmpty
	}

	ruleStrings := strings.Split(tag, ValidationTagSeparator)
	validationRules := make(rules.ValidationRules, 0, len(ruleStrings))

	for _, ruleString := range ruleStrings {
		rule, err := t.RuleParser.GetRule(ruleString)
		if err != nil {
			return nil, fmt.Errorf(ErrParsingValidationTag.Error()+": %w", err)
		}

		validationRules = append(validationRules, rule)
	}

	return validationRules, nil
}
