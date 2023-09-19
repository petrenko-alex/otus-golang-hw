package validators

import (
	"github.com/stretchr/testify/require"
	"reflect"
	"testing"
)

func TestFieldTypeValidatorFactory_GetValidator(t *testing.T) {
	var factory ValidatorFactory = FieldTypeValidatorFactory{}

	t.Run("get validator for scalar values", func(t *testing.T) {
		fieldType := reflect.Int
		validationTag := "min:0|max:10"

		validator, err := factory.GetValidator(fieldType, validationTag)

		require.NoError(t, err)
		require.IsType(t, ScalarValueValidator{}, validator)
	})

	t.Run("get validator for slice values", func(t *testing.T) {
		fieldType := reflect.Array
		validationTag := "min:0|max:10"

		validator, err := factory.GetValidator(fieldType, validationTag)

		require.NoError(t, err)
		require.IsType(t, SliceValueValidator{}, validator)
	})

	t.Run("incorrect validation tag", func(t *testing.T) {
		fieldType := reflect.Int
		validationTag := "min:"

		_, err := factory.GetValidator(fieldType, validationTag)

		require.ErrorContains(t, err, ErrValidationTagArg.Error())
	})

	t.Run("unsupported field type", func(t *testing.T) {
		fieldType := reflect.Bool
		validationTag := "min:0|max:10"

		_, err := factory.GetValidator(fieldType, validationTag)

		require.ErrorIs(t, err, ErrUnsupportedFieldType)
	})

	t.Run("incorrect fieldType arg", func(t *testing.T) {
		fieldType := "bool"
		validationTag := "min:0|max:10"

		_, err := factory.GetValidator(fieldType, validationTag)

		require.ErrorIs(t, err, ErrFieldTypeArg)
	})
}
