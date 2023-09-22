package validators

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFieldTypeValidatorFactory_GetValidator(t *testing.T) {
	var factory ValidatorFactory = FieldTypeValidatorFactory{}
	const validationTag = "min:0|max:10"

	t.Run("get validator for scalar values", func(t *testing.T) {
		fieldType := reflect.Int

		validator, err := factory.GetValidator(fieldType, validationTag)

		require.NoError(t, err)
		require.IsType(t, ScalarValueValidator{}, validator)
	})

	t.Run("get validator for slice values", func(t *testing.T) {
		fieldType := reflect.Array

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

		_, err := factory.GetValidator(fieldType, validationTag)

		require.ErrorIs(t, err, ErrUnsupportedFieldType)
	})

	t.Run("incorrect fieldType arg", func(t *testing.T) {
		fieldType := "bool"

		_, err := factory.GetValidator(fieldType, validationTag)

		require.ErrorIs(t, err, ErrFieldTypeArg)
	})
}
