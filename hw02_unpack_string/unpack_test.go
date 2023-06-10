package hw02unpackstring_test

import (
	"errors"
	hw02unpackstring "github.com/petrenko-alex/otus-golang-hw/hw02_unpack_string"
	"testing"

	"github.com/stretchr/testify/require"
)

// TODO: Uncomment extra cases
// TODO: run coverage

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

		// extra
		// {input: `qwe\4\5`, expected: `qwe45`},
		// {input: `qwe\45`, expected: `qwe44444`},
		// {input: `qwe\\5`, expected: `qwe\\\\\`},
		// {input: `qwe\\\3`, expected: `qwe\3`},
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
		desc  string
		input string
	}{
		{desc: "String starts with digit", input: "3abc"},
		{desc: "String without letters", input: "45"},
		{desc: "String with number instead of digit", input: "aaa10b"},
		{desc: "String with non symbols", input: "a4d&"},
		{desc: "String emojis", input: "😀"},
	}

	for i := range testCases {
		testCase := testCases[i]
		t.Run(testCase.desc, func(t *testing.T) {
			_, err := hw02unpackstring.Unpack(testCase.input)

			require.Truef(t, errors.Is(err, hw02unpackstring.ErrInvalidString), "actual error %q", err)
		})
	}
}
