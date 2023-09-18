package rules

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestIntRangeRule_Validate(t *testing.T) {
	t.Run("incorrect value for rule", func(t *testing.T) {
		var limit ValidationLimit = []int{256, 1024}
		rule := IntRangeRule{limit}

		err := rule.Validate("5")

		require.ErrorIs(t, err, ErrCastValueForRule)
	})

	t.Run("incorrect limit for rule", func(t *testing.T) {
		var limit ValidationLimit = 256
		rule := IntRangeRule{limit}

		err := rule.Validate(5)

		require.ErrorIs(t, err, ErrCastLimitForRule)
	})

	t.Run("value inside range", func(t *testing.T) {
		var limit ValidationLimit = []int{256, 1024}
		rule := IntRangeRule{limit}

		err := rule.Validate(512)

		require.NoError(t, err)
	})

	t.Run("value greater than range", func(t *testing.T) {
		var limit ValidationLimit = []int{256, 1024}
		rule := IntRangeRule{limit}

		err := rule.Validate(2048)

		require.ErrorIs(t, err, ErrValidationFailed)
	})

	t.Run("value lower than range", func(t *testing.T) {
		var limit ValidationLimit = []int{256, 1024}
		rule := IntRangeRule{limit}

		err := rule.Validate(128)

		require.ErrorIs(t, err, ErrValidationFailed)
	})

	t.Run("value starts range", func(t *testing.T) {
		var limit ValidationLimit = []int{256, 1024}
		rule := IntRangeRule{limit}

		err := rule.Validate(256)

		require.NoError(t, err)
	})

	t.Run("value ends range", func(t *testing.T) {
		var limit ValidationLimit = []int{256, 1024}
		rule := IntRangeRule{limit}

		err := rule.Validate(1024)

		require.NoError(t, err)
	})
}
