package rules

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestLenRule_Validate(t *testing.T) {
	t.Run("incorrect value for rule", func(t *testing.T) {
		var limit ValidationLimit = 5
		rule := LenRule{limit}

		err := rule.Validate(5)

		require.ErrorIs(t, err, ErrCastValueForRule)
	})

	t.Run("incorrect limit for rule", func(t *testing.T) {
		var limit ValidationLimit = "limit"
		rule := LenRule{limit}

		err := rule.Validate("string")

		require.ErrorIs(t, err, ErrCastLimitForRule)
	})

	t.Run("value less than limit", func(t *testing.T) {
		var limit ValidationLimit = "5"
		rule := LenRule{limit}

		err := rule.Validate("str")

		require.ErrorIs(t, err, ErrValidationFailed)

	})

	t.Run("value greater than  limit", func(t *testing.T) {
		var limit ValidationLimit = "5"
		rule := LenRule{limit}

		err := rule.Validate("string")

		require.ErrorIs(t, err, ErrValidationFailed)

	})

	t.Run("value satisfy limit", func(t *testing.T) {
		var limit ValidationLimit = "5"
		rule := LenRule{limit}

		err := rule.Validate("strin")

		require.NoError(t, err)

	})
}
