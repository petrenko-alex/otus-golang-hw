package rules

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestInRule_Validate(t *testing.T) {
	t.Run("incorrect value for rule", func(t *testing.T) {
		var limit ValidationLimit = []string{"256", "1024"}
		rule := InRule{limit}

		err := rule.Validate(struct{}{})

		require.ErrorIs(t, err, ErrCastValueForRule)
	})

	t.Run("incorrect limit for rule", func(t *testing.T) {
		var limit []struct{}
		rule := InRule{limit}

		err := rule.Validate(5)

		require.ErrorIs(t, err, ErrCastLimitForRule)
	})

	t.Run("work with int", func(t *testing.T) {
		var limit ValidationLimit = []string{"256", "1024"}
		rule := InRule{limit}

		err := rule.Validate(512)

		require.NoError(t, err)
	})

	t.Run("work with string", func(t *testing.T) {
		var limit ValidationLimit = []string{"foo", "bar"}
		rule := InRule{limit}

		err := rule.Validate("foo")

		require.NoError(t, err)
	})
}
