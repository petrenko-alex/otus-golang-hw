package parsing

import (
	"github.com/petrenko-alex/otus-golang-hw/hw09_struct_validator/rules"
	"strings"
)

const (
	ValidationTagSeparator = "|"
)

type TagParser struct {
	Factory RuleParser // todo: rename from Factory
}

func (t TagParser) GetValidationRules(tag ValidationTag) (rules.ValidationRules, error) {
	// todo: common regex for validationtag

	ruleStrings := strings.Split(tag, ValidationTagSeparator)
	validationRules := make(rules.ValidationRules, 0, len(ruleStrings))

	for _, ruleString := range ruleStrings {

		rule, err := t.Factory.GetRule(ruleString)
		if err != nil {
			// todo: wrap ?
			return nil, err
		}

		validationRules = append(validationRules, rule)
	}

	return validationRules, nil
}
