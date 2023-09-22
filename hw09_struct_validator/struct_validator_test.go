package hw09structvalidator

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/petrenko-alex/otus-golang-hw/hw09_struct_validator/validators" //nolint:depguard
	"github.com/stretchr/testify/require"                                      //nolint:depguard
)

type UserRole string

// Test the function on different structures and other types.
type (
	User struct {
		ID     string `json:"id" validate:"len:36"`
		Name   string
		Age    int      `validate:"min:18|max:50"`
		Email  string   `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role   UserRole `validate:"in:admin,stuff"`
		Phones []string `validate:"len:11"`
		meta   json.RawMessage
	}

	App struct {
		Version string `validate:"len:5"`
	}

	Token struct {
		Header    []byte
		Payload   []byte
		Signature []byte
	}

	Response struct {
		Code string `validate:"in:200,404,500"`
		Body string `json:"omitempty"`
	}

	Unexported struct {
		fieldOne string `validate:"len:11"`
		fieldTwo int    `validate:"min:18"`
	}

	UnsupportedFieldType struct {
		Field bool `validate:"in:true"`
	}

	ValidateRuntimeError struct {
		Field string `validate:"min:10"`
	}
)

func TestStructValidator_Validate_Errors(t *testing.T) {
	var validatorFactory validators.ValidatorFactory = validators.FieldTypeValidatorFactory{}
	var structValidator Validator = StructValidator{validatorFactory}

	t.Run("only struct is available", func(t *testing.T) {
		err := structValidator.Validate("simple string")

		require.ErrorIs(t, err, ErrInputNotStruct)
	})

	t.Run("not supported field type", func(t *testing.T) {
		err := structValidator.Validate(UnsupportedFieldType{Field: false})

		require.ErrorIs(t, err, ErrValidatorInit)
	})

	t.Run("validate runtime error", func(t *testing.T) {
		err := structValidator.Validate(ValidateRuntimeError{Field: "value"})

		require.ErrorContains(t, err, validators.ErrRuntime.Error())
	})
}

func TestStructValidator_Validate(t *testing.T) {
	var validatorFactory validators.ValidatorFactory = validators.FieldTypeValidatorFactory{}
	var structValidator Validator = StructValidator{validatorFactory}

	t.Run("empty struct", func(t *testing.T) {
		input := struct{}{}

		err := structValidator.Validate(input)

		require.Nil(t, err)
	})

	t.Run("User struct", func(t *testing.T) {
		input := User{
			ID:     "510",
			Name:   "Alex",
			Age:    17,
			Email:  "test$test.ru",
			Role:   UserRole("client"),
			Phones: []string{"88007006", "8100200302"},
			meta:   json.RawMessage("{error: null}"),
		}

		err := structValidator.Validate(input)

		var validationErr ValidationErrors
		if !errors.As(err, &validationErr) {
			require.Fail(t, "unexpected return type")
		}

		require.Len(t, validationErr, 6)
		require.Equal(t, []string{"ID", "Age", "Email", "Role", "Phones", "Phones"}, validationErr.GetFields())
	})

	t.Run("App struct", func(t *testing.T) {
		input := App{Version: "1.0"}

		err := structValidator.Validate(input)

		var validationErr ValidationErrors
		if !errors.As(err, &validationErr) {
			require.Fail(t, "unexpected return type")
		}

		require.Len(t, validationErr, 1)
		require.IsType(t, ValidationError{}, validationErr[0])
		require.Equal(t, "Version", validationErr[0].Field)
	})

	t.Run("Token struct", func(t *testing.T) {
		input := Token{
			Header:    []byte("some header"),
			Payload:   []byte("some payload"),
			Signature: []byte("some signature"),
		}

		err := structValidator.Validate(input)

		require.Nil(t, err)
	})

	t.Run("Response struct", func(t *testing.T) {
		input := Response{
			Code: "201",
			Body: "<h1>Hello World</h1>",
		}

		err := structValidator.Validate(input)

		var validationErr ValidationErrors
		if !errors.As(err, &validationErr) {
			require.Fail(t, "unexpected return type")
		}

		require.Len(t, validationErr, 1)
		require.IsType(t, ValidationError{}, validationErr[0])
		require.Equal(t, "Code", validationErr[0].Field)
	})

	t.Run("Unexported struct", func(t *testing.T) {
		input := Unexported{
			fieldOne: "value",
			fieldTwo: 10,
		}

		err := structValidator.Validate(input)

		require.Nil(t, err)
	})
}
