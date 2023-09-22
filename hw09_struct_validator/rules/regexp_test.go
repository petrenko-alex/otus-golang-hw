package rules

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRegexpRule_Validate(t *testing.T) {
	t.Run("incorrect value for rule", func(t *testing.T) {
		var limit ValidationLimit = "\\d+"
		rule := RegexpRule{limit}

		err := rule.Validate(5)

		require.ErrorIs(t, err, ErrCastValueForRule)
	})

	t.Run("incorrect limit for rule", func(t *testing.T) {
		var limit ValidationLimit = struct{}{}
		rule := RegexpRule{limit}

		err := rule.Validate("5")

		require.ErrorIs(t, err, ErrCastLimitForRule)
	})

	t.Run("incorrect regexp", func(t *testing.T) {
		var limit ValidationLimit = "[d+"
		rule := RegexpRule{limit}

		err := rule.Validate("5")

		require.ErrorIs(t, err, ErrCastLimitForRule)

	})
	t.Run("value satisfy limit", func(t *testing.T) {
		var limit ValidationLimit = "\\d+"
		rule := RegexpRule{limit}

		err := rule.Validate("5")

		require.NoError(t, err)

	})
	t.Run("value does not satisfy limit", func(t *testing.T) {
		var limit ValidationLimit = "\\d+"
		rule := RegexpRule{limit}

		err := rule.Validate("not a number")

		require.ErrorIs(t, err, ErrValidationFailed)
	})
}
