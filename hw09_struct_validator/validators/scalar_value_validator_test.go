package validators

import (
	"github.com/petrenko-alex/otus-golang-hw/hw09_struct_validator/rules"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestScalarValueValidator_ValidateValue(t *testing.T) {
	t.Run("no validation rules", func(t *testing.T) {
		var validator ValueValidator = ScalarValueValidator{}

		validationErrors, runtimeErr := validator.ValidateValue(10)

		require.Len(t, validationErrors, 0)
		require.Nil(t, runtimeErr)
	})

	t.Run("one rule, no validationErrors", func(t *testing.T) {
		var validator ValueValidator = ScalarValueValidator{
			rules.ValidationRules{
				rules.MaxRule{"20"},
			},
		}

		validationErrors, runtimeErr := validator.ValidateValue(10)

		require.Len(t, validationErrors, 0)
		require.Nil(t, runtimeErr)
	})

	t.Run("one rule, one error", func(t *testing.T) {
		var validator ValueValidator = ScalarValueValidator{
			rules.ValidationRules{
				rules.MaxRule{"20"},
			},
		}

		validationErrors, runtimeErr := validator.ValidateValue(25)

		require.Len(t, validationErrors, 1)
		require.Nil(t, runtimeErr)
	})

	t.Run("multiple rule, no validationErrors", func(t *testing.T) {
		var validator ValueValidator = ScalarValueValidator{
			rules.ValidationRules{
				rules.MaxRule{"20"},
				rules.MinRule{"10"},
			},
		}

		validationErrors, runtimeErr := validator.ValidateValue(15)

		require.Len(t, validationErrors, 0)
		require.Nil(t, runtimeErr)
	})

	t.Run("multiple rule, one error", func(t *testing.T) {
		var validator ValueValidator = ScalarValueValidator{
			rules.ValidationRules{
				rules.MaxRule{"20"},
				rules.MinRule{"10"},
			},
		}

		validationErrors, runtimeErr := validator.ValidateValue(5)

		require.Len(t, validationErrors, 1)
		require.Nil(t, runtimeErr)
	})

	t.Run("multiple rule, multiple validationErrors", func(t *testing.T) {
		inRule, _ := rules.NewInRule("10,30")
		var validator ValueValidator = ScalarValueValidator{
			rules.ValidationRules{
				rules.MaxRule{"40"},
				inRule,
			},
		}

		validationErrors, runtimeErr := validator.ValidateValue(50)

		require.Len(t, validationErrors, 2)
		require.Nil(t, runtimeErr)
	})

	t.Run("runtime error", func(t *testing.T) {
		var validator ValueValidator = ScalarValueValidator{
			rules.ValidationRules{
				rules.MaxRule{"20"},
			},
		}

		validationErrors, runtimeErr := validator.ValidateValue("25")

		require.Nil(t, validationErrors)
		require.ErrorContains(t, runtimeErr, RuntimeError.Error())
	})
}
