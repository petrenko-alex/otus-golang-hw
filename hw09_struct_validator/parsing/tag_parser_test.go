package parsing

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTagParser_GetValidationRules(t *testing.T) {
	tests := []struct {
		name, input string
		rulesCount  int
	}{
		{
			name:       "one rule",
			input:      "len:5",
			rulesCount: 1,
		},
		{
			name:       "multiple rules",
			input:      "len:5|min:10|in:256,1024|regexp:\\d+",
			rulesCount: 4,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var ruleParser RuleParser = BaseRuleParser{}
			var tagParser ValidationTagParser = TagParser{ruleParser}

			rules, err := tagParser.GetValidationRules(test.input)

			require.NoError(t, err)
			require.Len(t, rules, test.rulesCount)
		})
	}
}

func TestTagParser_GetValidationRules_Errors(t *testing.T) {
	tests := []struct {
		name, input string
		output      error
	}{
		{
			name:   "empty tag",
			input:  "",
			output: ErrValidationTagEmpty,
		},
		{
			name:   "corrupted tag #1",
			input:  "len:5|min:",
			output: ErrParsingValidationTag,
		},
		{
			name:   "corrupted tag #2",
			input:  "min:10|",
			output: ErrParsingValidationTag,
		},
		{
			name:   "multiple rules w/ unknown",
			input:  "len:5|min:10|in:256,1024|magic:rule",
			output: ErrParsingValidationTag,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var ruleParser RuleParser = BaseRuleParser{}
			var tagParser ValidationTagParser = TagParser{ruleParser}

			_, err := tagParser.GetValidationRules(test.input)

			require.ErrorContains(t, err, test.output.Error())
		})
	}
}
