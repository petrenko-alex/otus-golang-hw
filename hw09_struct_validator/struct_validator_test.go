package hw09structvalidator

import (
	"encoding/json"
	"github.com/petrenko-alex/otus-golang-hw/hw09_struct_validator/validators"
	"github.com/stretchr/testify/require"
	"testing"
)

type UserRole string

// Test the function on different structures and other types.
type (
	User struct {
		ID     string `json:"id" validate:"len:36"`
		Name   string
		Age    int             `validate:"min:18|max:50"`
		Email  string          `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role   UserRole        `validate:"in:admin,stuff"`
		Phones []string        `validate:"len:11"`
		meta   json.RawMessage //nolint:unused
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
		Code int    `validate:"in:200,404,500"`
		Body string `json:"omitempty"`
	}
)

// todo: test cases for each validator
// todo: how to work with unexported fields

func TestValidate(t *testing.T) {
	testCases := []struct {
		name        string
		input       interface{}
		expectedErr error
	}{
		{
			name:        "only struct is available",
			input:       "simple string",
			expectedErr: ErrInputNotStruct,
		},
		{
			name:        "empty struct",
			input:       struct{}{},
			expectedErr: nil,
		},
		{
			name:        "one field, no tag",
			input:       struct{ Fld string }{Fld: "value"},
			expectedErr: nil,
		},
		{
			name: "one field, alien tag",
			input: struct {
				Fld string `json:"name"`
			}{Fld: "value"},
			expectedErr: nil,
		},
		{
			name: "one field, one tag, corrupted",
			input: struct {
				Fld string `validate:`
			}{Fld: "val"},
			expectedErr: ErrValidatorInit,
		},
		{
			name: "one field, one tag, empty",
			input: struct {
				Fld string `validate:""`
			}{Fld: "val"},
			expectedErr: ErrValidatorInit,
		},
		{
			name: "one field, one tag, satisfy",
			input: struct {
				Fld string `validate:"len:5"`
			}{Fld: "val"},
			expectedErr: nil,
		},
		// name: "one field, many tags"
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			testCase := testCase
			t.Parallel()

			var validatorFactory validators.ValidatorFactory = validators.FieldTypeValidatorFactory{}
			var structValidator Validator = StructValidator{validatorFactory}

			err := structValidator.Validate(testCase.input)
			require.ErrorIs(t, err, testCase.expectedErr)
		})
	}
}
