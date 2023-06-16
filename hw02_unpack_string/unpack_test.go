package hw02unpackstring_test

import (
	"errors"
	"testing"

	hw02unpackstring "github.com/petrenko-alex/otus-golang-hw/hw02_unpack_string"
	"github.com/stretchr/testify/require"
)

func TestUnpack(t *testing.T) {
	testCases := []struct {
		desc     string
		input    string
		expected string
	}{
		// nothing to do
		{desc: "Empty string", input: "", expected: ""},
		{desc: "String with only one letter", input: "a", expected: "a"},
		{desc: "String w/o digits", input: "abccd", expected: "abccd"},

		// simple
		{desc: "One letter, multiplier 1", input: "a1", expected: "a"},
		{desc: "One letter, multiplier 2", input: "a2", expected: "aa"},
		{desc: "Many letter, different multiplier", input: "a2b3c4", expected: "aabbbcccc"},

		// erase
		{desc: "Erase to empty string", input: "a0", expected: ""},
		{desc: "Erase to one letter", input: "ab0", expected: "a"},
		{desc: "Normal erase", input: "aaa0b", expected: "aab"},

		// escaping
		{desc: "Escape one digit", input: `\4`, expected: `4`},
		{desc: "Escape one slash", input: `\\`, expected: `\`},
		{desc: "Escape one digit and erase", input: `\40`, expected: ``},
		{desc: "Escape one slash and erase", input: `\\0`, expected: ``},
		{desc: "Escape digits", input: `qwe\4\5`, expected: `qwe45`},
		{desc: "Escape digits with multiplier", input: `qwe\45`, expected: `qwe44444`},
		{desc: "Escape slash with multiplier", input: `qwe\\5`, expected: `qwe\\\\\`},
		{desc: "Escape slash and digits", input: `qwe\\\3`, expected: `qwe\3`},
		{desc: "Escape with erase", input: `qwe\\2\30`, expected: `qwe\\`},

		// complex
		{desc: "Complex 1", input: "a4bc2d5e", expected: "aaaabccddddde"},
		{desc: "Complex 2", input: "a4bc2d5e0f2", expected: "aaaabccdddddff"},
		{desc: "Complex with cyrillic", input: "а3п2орпz3", expected: "аааппорпzzz"},
		{desc: "Complex with escaping", input: `а3п2орпz3\\\3n2\\3\32\90`, expected: `аааппорпzzz\3nn\\\33`},
	}

	for i := range testCases {
		testCase := testCases[i]
		t.Run(testCase.desc, func(t *testing.T) {
			result, err := hw02unpackstring.Unpack(testCase.input)

			require.NoError(t, err)
			require.Equal(t, testCase.expected, result)
		})
	}
}

func TestUnpackInvalidString(t *testing.T) {
	testCases := []struct {
		desc          string
		input         string
		expectedError error
	}{
		{desc: "String starts with digit", input: "3abc", expectedError: hw02unpackstring.ErrStartsWithDigits},
		{desc: "String has only digit", input: "4", expectedError: hw02unpackstring.ErrStartsWithDigits},
		{desc: "String with number instead of digit", input: "aaa10b", expectedError: hw02unpackstring.ErrHasNumbers},
		{desc: "String with non symbols", input: "a4d&", expectedError: hw02unpackstring.ErrInvalidChars},
		{desc: "String emojis", input: "😀", expectedError: hw02unpackstring.ErrInvalidChars},
		{desc: "Invalid escaped symbol", input: `qwe\\\n`, expectedError: hw02unpackstring.ErrInvalidEscaping},
		{desc: "Invalid escaping 1", input: `\`, expectedError: hw02unpackstring.ErrInvalidEscaping},
		{desc: "Invalid escaping 2", input: `\3\\\`, expectedError: hw02unpackstring.ErrInvalidEscaping},
	}

	for i := range testCases {
		testCase := testCases[i]
		t.Run(testCase.desc, func(t *testing.T) {
			_, err := hw02unpackstring.Unpack(testCase.input)

			require.Truef(t, errors.Is(err, testCase.expectedError), "actual error %q", err)
		})
	}
}
