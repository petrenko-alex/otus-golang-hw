package validators

import (
	"testing"

	"github.com/petrenko-alex/otus-golang-hw/hw09_struct_validator/rules"
	"github.com/stretchr/testify/require"
)

func TestSliceValueValidator_ValidateValue_Errors(t *testing.T) {
	t.Run("unsupported slice", func(t *testing.T) {
		var validator ValueValidator = SliceValueValidator{
			ScalarValueValidator{
				rules.ValidationRules{
					rules.LenRule{Limit: "3"},
				},
			},
		}

		validationErrors, runtimeErr := validator.ValidateValue([]bool{true, false})

		require.Nil(t, validationErrors)
		require.ErrorIs(t, runtimeErr, ErrValueNotSupported)
	})

	t.Run("incorrect validator rule", func(t *testing.T) {
		var validator ValueValidator = SliceValueValidator{
			ScalarValueValidator{
				rules.ValidationRules{
					rules.LenRule{Limit: "3"},
				},
			},
		}

		validationErrors, runtimeErr := validator.ValidateValue([]int{1, 2, 3})

		require.Nil(t, validationErrors)
		require.ErrorIs(t, runtimeErr, rules.ErrCastValueForRule)
	})

	t.Run("slice of interfaces", func(t *testing.T) {
		var validator ValueValidator = SliceValueValidator{
			ScalarValueValidator{
				rules.ValidationRules{
					rules.LenRule{Limit: "3"},
				},
			},
		}

		validationErrors, runtimeErr := validator.ValidateValue([]interface{}{1, "string"})

		require.Nil(t, validationErrors)
		require.ErrorIs(t, runtimeErr, ErrValueNotSupported)
	})

	t.Run("validate non slice value", func(t *testing.T) {
		var validator ValueValidator = SliceValueValidator{
			ScalarValueValidator{
				rules.ValidationRules{
					rules.LenRule{Limit: "3"},
				},
			},
		}

		validationErrors, runtimeErr := validator.ValidateValue(1)

		require.Nil(t, validationErrors)
		require.ErrorIs(t, runtimeErr, ErrValueNotIterable)
	})
}

func TestSliceValueValidator_ValidateValue(t *testing.T) {
	t.Run("no validation rules", func(t *testing.T) {
		var validator ValueValidator = SliceValueValidator{
			ScalarValueValidator{},
		}

		validationErrors, runtimeErr := validator.ValidateValue([]int{10, 20, 30})

		require.Len(t, validationErrors, 0)
		require.Nil(t, runtimeErr)
	})

	t.Run("slice of int", func(t *testing.T) {
		var validator ValueValidator = SliceValueValidator{
			ScalarValueValidator{
				rules.ValidationRules{
					rules.MinRule{Limit: "15"},
				},
			},
		}

		validationErrors, runtimeErr := validator.ValidateValue([]int{10, 20, 30})

		require.Len(t, validationErrors, 1)
		require.Nil(t, runtimeErr)
	})

	t.Run("slice of string", func(t *testing.T) {
		var validator ValueValidator = SliceValueValidator{
			ScalarValueValidator{
				rules.ValidationRules{
					rules.LenRule{Limit: "3"},
				},
			},
		}

		validationErrors, runtimeErr := validator.ValidateValue([]string{"foo", "bar", "goo1", "gooo"})

		require.Len(t, validationErrors, 2)
		require.Nil(t, runtimeErr)
	})

	t.Run("one element slice, satisfy", func(t *testing.T) {
		var validator ValueValidator = SliceValueValidator{
			ScalarValueValidator{
				rules.ValidationRules{
					rules.LenRule{Limit: "3"},
				},
			},
		}

		validationErrors, runtimeErr := validator.ValidateValue([]string{"foo"})

		require.Len(t, validationErrors, 0)
		require.Nil(t, runtimeErr)
	})

	t.Run("one element slice, failed", func(t *testing.T) {
		var validator ValueValidator = SliceValueValidator{
			ScalarValueValidator{
				rules.ValidationRules{
					rules.LenRule{Limit: "4"},
				},
			},
		}

		validationErrors, runtimeErr := validator.ValidateValue([]string{"foo"})

		require.Len(t, validationErrors, 1)
		require.Nil(t, runtimeErr)
	})

	t.Run("all elements failed", func(t *testing.T) {
		var validator ValueValidator = SliceValueValidator{
			ScalarValueValidator{
				rules.ValidationRules{
					rules.LenRule{Limit: "4"},
				},
			},
		}

		validationErrors, runtimeErr := validator.ValidateValue([]string{"foo", "bar", "goo"})

		require.Len(t, validationErrors, 3)
		require.Nil(t, runtimeErr)
	})

	t.Run("all elements satisfy", func(t *testing.T) {
		var validator ValueValidator = SliceValueValidator{
			ScalarValueValidator{
				rules.ValidationRules{
					rules.LenRule{Limit: "3"},
				},
			},
		}

		validationErrors, runtimeErr := validator.ValidateValue([]string{"foo", "bar", "goo"})

		require.Len(t, validationErrors, 0)
		require.Nil(t, runtimeErr)
	})
}
