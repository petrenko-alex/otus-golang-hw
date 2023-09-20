package validators

import (
	"github.com/petrenko-alex/otus-golang-hw/hw09_struct_validator/rules"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestSliceValueValidator_ValidateValue(t *testing.T) {
	t.Run("no validation rules", func(t *testing.T) {
		var validator ValueValidator = SliceValueValidator{
			ScalarValueValidator{},
		}

		errors := validator.ValidateValue([]int{10, 20, 30})

		require.Len(t, errors, 0)
	})

	t.Run("slice of int", func(t *testing.T) {
		var validator ValueValidator = SliceValueValidator{
			ScalarValueValidator{
				rules.ValidationRules{
					rules.MinRule{"15"},
				},
			},
		}

		errors := validator.ValidateValue([]int{10, 20, 30})

		require.Len(t, errors, 1)
	})

	t.Run("slice of string", func(t *testing.T) {
		var validator ValueValidator = SliceValueValidator{
			ScalarValueValidator{
				rules.ValidationRules{
					rules.LenRule{"3"},
				},
			},
		}

		errors := validator.ValidateValue([]string{"foo", "bar", "goo1", "gooo"})

		require.Len(t, errors, 2)
	})

	t.Run("one element slice, satisfy", func(t *testing.T) {
		var validator ValueValidator = SliceValueValidator{
			ScalarValueValidator{
				rules.ValidationRules{
					rules.LenRule{"3"},
				},
			},
		}

		errors := validator.ValidateValue([]string{"foo"})

		require.Len(t, errors, 0)
	})

	t.Run("one element slice, failed", func(t *testing.T) {
		var validator ValueValidator = SliceValueValidator{
			ScalarValueValidator{
				rules.ValidationRules{
					rules.LenRule{"4"},
				},
			},
		}

		errors := validator.ValidateValue([]string{"foo"})

		require.Len(t, errors, 1)
	})

	t.Run("all elements failed", func(t *testing.T) {
		var validator ValueValidator = SliceValueValidator{
			ScalarValueValidator{
				rules.ValidationRules{
					rules.LenRule{"4"},
				},
			},
		}

		errors := validator.ValidateValue([]string{"foo", "bar", "goo"})

		require.Len(t, errors, 3)
	})

	t.Run("all elements satisfy", func(t *testing.T) {
		var validator ValueValidator = SliceValueValidator{
			ScalarValueValidator{
				rules.ValidationRules{
					rules.LenRule{"3"},
				},
			},
		}

		errors := validator.ValidateValue([]string{"foo", "bar", "goo"})

		require.Len(t, errors, 0)
	})

	t.Run("unsupported slice", func(t *testing.T) {
		var validator ValueValidator = SliceValueValidator{
			ScalarValueValidator{
				rules.ValidationRules{
					rules.LenRule{"3"},
				},
			},
		}

		errors := validator.ValidateValue([]bool{true, false})

		require.Len(t, errors, 1)
		require.ErrorIs(t, errors[0], ErrValueNotSupported)
	})

	t.Run("incorrect validator rule", func(t *testing.T) {
		var validator ValueValidator = SliceValueValidator{
			ScalarValueValidator{
				rules.ValidationRules{
					rules.LenRule{"3"},
				},
			},
		}

		errors := validator.ValidateValue([]int{1, 2, 3})

		require.Len(t, errors, 3)
		require.ErrorIs(t, errors[0], rules.ErrCastValueForRule)
	})

	t.Run("slice of interfaces", func(t *testing.T) {
		var validator ValueValidator = SliceValueValidator{
			ScalarValueValidator{
				rules.ValidationRules{
					rules.LenRule{"3"},
				},
			},
		}

		errors := validator.ValidateValue([]interface{}{1, "string"})

		require.Len(t, errors, 1)
		require.ErrorIs(t, errors[0], ErrValueNotSupported)
	})

	t.Run("validate non slice value", func(t *testing.T) {
		var validator ValueValidator = SliceValueValidator{
			ScalarValueValidator{
				rules.ValidationRules{
					rules.LenRule{"3"},
				},
			},
		}

		errors := validator.ValidateValue(1)

		require.Len(t, errors, 1)
		require.ErrorIs(t, errors[0], ErrValueNotIterable)
	})

}
