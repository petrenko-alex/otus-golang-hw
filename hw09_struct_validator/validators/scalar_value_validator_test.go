package validators

import (
	"github.com/petrenko-alex/otus-golang-hw/hw09_struct_validator/rules"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestScalarValueValidator_ValidateValue(t *testing.T) {
	t.Run("no validation rules", func(t *testing.T) {
		var validator ValueValidator = ScalarValueValidator{}

		errors := validator.ValidateValue(10)

		require.Len(t, errors, 0)
	})

	t.Run("one rule, no errors", func(t *testing.T) {
		var validator ValueValidator = ScalarValueValidator{
			rules.ValidationRules{
				rules.MaxRule{"20"},
			},
		}

		errors := validator.ValidateValue(10)

		require.Len(t, errors, 0)
	})

	t.Run("one rule, one error", func(t *testing.T) {
		var validator ValueValidator = ScalarValueValidator{
			rules.ValidationRules{
				rules.MaxRule{"20"},
			},
		}

		errors := validator.ValidateValue(25)

		require.Len(t, errors, 1)
	})

	t.Run("multiple rule, no errors", func(t *testing.T) {
		var validator ValueValidator = ScalarValueValidator{
			rules.ValidationRules{
				rules.MaxRule{"20"},
				rules.MinRule{"10"},
			},
		}

		errors := validator.ValidateValue(15)

		require.Len(t, errors, 0)
	})

	t.Run("multiple rule, one error", func(t *testing.T) {
		var validator ValueValidator = ScalarValueValidator{
			rules.ValidationRules{
				rules.MaxRule{"20"},
				rules.MinRule{"10"},
			},
		}

		errors := validator.ValidateValue(5)

		require.Len(t, errors, 1)
	})

	t.Run("multiple rule, multiple errors", func(t *testing.T) {
		inRule, _ := rules.NewInRule("10,30")
		var validator ValueValidator = ScalarValueValidator{
			rules.ValidationRules{
				rules.MaxRule{"40"},
				inRule,
			},
		}

		errors := validator.ValidateValue(50)

		require.Len(t, errors, 2)
	})
}
