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

		// complex
		{desc: "Complex 1", input: "a4bc2d5e", expected: "aaaabccddddde"},
		{desc: "Complex 2", input: "a4bc2d5e0f2", expected: "aaaabccdddddff"},
		{desc: "Cyrillic", input: "а3п2орпz3", expected: "аааппорпzzz"},

		// extra
		// {input: `qwe\4\5`, expected: `qwe45`},
		// {input: `qwe\45`, expected: `qwe44444`},
		// {input: `qwe\\5`, expected: `qwe\\\\\`},
		// {input: `qwe\\\3`, expected: `qwe\3`},
		// {input: `qwe\\\n`, expected: `qwe\3`},
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
	}

	for i := range testCases {
		testCase := testCases[i]
		t.Run(testCase.desc, func(t *testing.T) {
			_, err := hw02unpackstring.Unpack(testCase.input)

			require.Truef(t, errors.Is(err, testCase.expectedError), "actual error %q", err)
		})
	}
}
