package parsing

import (
	"github.com/petrenko-alex/otus-golang-hw/hw09_struct_validator/rules"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestBaseRuleParser_GetRule(t *testing.T) {
	tests := []struct {
		name, input string
		output      rules.ValidationRule
	}{
		{
			name:   "len rule",
			input:  "len:5",
			output: rules.LenRule{},
		},
		{
			name:   "regexp rule",
			input:  "regexp:\\d+",
			output: rules.RegexpRule{},
		},
		{
			name:   "in rule",
			input:  "in:122,322",
			output: rules.InRule{},
		},
		{
			name:   "min rule",
			input:  "min:10",
			output: rules.MinRule{},
		},
		{
			name:   "max rule",
			input:  "max:20",
			output: rules.MaxRule{},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var ruleParser RuleParser = BaseRuleParser{}

			rule, err := ruleParser.GetRule(test.input)

			require.NoError(t, err)
			require.IsType(t, test.output, rule)
		})
	}
}

func TestBaseRuleParser_GetRule_Errors(t *testing.T) {
	tests := []struct {
		name, input string
		output      error
	}{
		{
			name:   "empty string rule",
			input:  "",
			output: ErrParsingRule,
		},
		{
			name:   "incorrect string rule structure #1",
			input:  "no delimiter",
			output: ErrParsingRule,
		},
		{
			name:   "incorrect string rule structure #2",
			input:  "len:",
			output: ErrParsingRule,
		},
		{
			name:   "incorrect string rule structure #3",
			input:  ":12",
			output: ErrParsingRule,
		},
		{
			name:   "unknown rule",
			input:  "magic:rule",
			output: ErrUnknownRule,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var ruleParser RuleParser = BaseRuleParser{}

			_, err := ruleParser.GetRule(test.input)

			require.ErrorIs(t, err, test.output)
		})
	}
}
