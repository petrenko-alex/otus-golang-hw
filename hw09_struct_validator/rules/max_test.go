package rules

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestMaxRule_Validate(t *testing.T) {
	t.Run("incorrect value for rule", func(t *testing.T) {
		var limit ValidationLimit = "5"
		rule := MaxRule{limit}

		err := rule.Validate("5")

		require.ErrorIs(t, err, ErrCastValueForRule)
	})

	t.Run("incorrect limit for rule", func(t *testing.T) {
		var limit ValidationLimit = "limit"
		rule := MaxRule{limit}

		err := rule.Validate(5)

		require.ErrorIs(t, err, ErrCastLimitForRule)
	})

	t.Run("value less than limit", func(t *testing.T) {
		var limit ValidationLimit = "20"
		rule := MaxRule{limit}

		err := rule.Validate(3)

		require.NoError(t, err)
	})

	t.Run("value equals limit", func(t *testing.T) {
		var limit ValidationLimit = "20"
		rule := MaxRule{limit}

		err := rule.Validate(20)

		require.NoError(t, err)
	})

	t.Run("value greater than limit", func(t *testing.T) {
		var limit ValidationLimit = "20"
		rule := MaxRule{limit}

		err := rule.Validate(22)

		require.ErrorIs(t, err, ErrValidationFailed)
	})

	t.Run("limit is negative", func(t *testing.T) {
		var limit ValidationLimit = "-5"
		rule := MaxRule{limit}

		err := rule.Validate(-10)

		require.NoError(t, err)
	})

	t.Run("limit is zero", func(t *testing.T) {
		var limit ValidationLimit = "0"
		rule := MaxRule{limit}

		err := rule.Validate(-2)

		require.NoError(t, err)
	})
}
