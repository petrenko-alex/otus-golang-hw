package rules

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestStringRangeRule_Validate(t *testing.T) {
	t.Run("incorrect value for rule", func(t *testing.T) {
		var limit ValidationLimit = []string{"foo", "bar"}
		rule := StringRangeRule{limit}

		err := rule.Validate(5)

		require.ErrorIs(t, err, ErrCastValueForRule)
	})

	t.Run("incorrect limit for rule", func(t *testing.T) {
		var limit ValidationLimit = true
		rule := StringRangeRule{limit}

		err := rule.Validate("foo")

		require.ErrorIs(t, err, ErrCastLimitForRule)
	})

	t.Run("value inside range", func(t *testing.T) {
		var limit ValidationLimit = []string{"foo", "bar"}
		rule := StringRangeRule{limit}

		err := rule.Validate("foo")

		require.NoError(t, err)
	})

	t.Run("value outside range", func(t *testing.T) {
		var limit ValidationLimit = []string{"foo", "bar"}
		rule := StringRangeRule{limit}

		err := rule.Validate("goo")

		require.ErrorIs(t, err, ErrValidationFailed)
	})
}
