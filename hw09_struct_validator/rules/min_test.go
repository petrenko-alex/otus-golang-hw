package rules

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestMinRule_Validate(t *testing.T) {
	t.Run("incorrect value for rule", func(t *testing.T) {
		var limit ValidationLimit = "5"
		rule := MinRule{limit}

		err := rule.Validate("5")

		require.ErrorIs(t, err, ErrCastValueForRule)
	})

	t.Run("incorrect limit for rule", func(t *testing.T) {
		var limit ValidationLimit = "limit"
		rule := MinRule{limit}

		err := rule.Validate("5")

		require.ErrorIs(t, err, ErrCastValueForRule)
	})

	t.Run("value less than limit", func(t *testing.T) {
		var limit ValidationLimit = "10"
		rule := MinRule{limit}

		err := rule.Validate(3)

		require.ErrorIs(t, err, ErrValidationFailed)
	})

	t.Run("value equals limit", func(t *testing.T) {
		var limit ValidationLimit = "10"
		rule := MinRule{limit}

		err := rule.Validate(10)

		require.NoError(t, err)
	})

	t.Run("value greater than limit", func(t *testing.T) {
		var limit ValidationLimit = "10"
		rule := MinRule{limit}

		err := rule.Validate(15)

		require.NoError(t, err)
	})

	t.Run("limit is negative", func(t *testing.T) {
		var limit ValidationLimit = "-10"
		rule := MinRule{limit}

		err := rule.Validate(-5)

		require.NoError(t, err)
	})

	t.Run("limit is zero", func(t *testing.T) {
		var limit ValidationLimit = "0"
		rule := MinRule{limit}

		err := rule.Validate(1)

		require.NoError(t, err)
	})
}
