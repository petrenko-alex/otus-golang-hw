package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(packedString string) (string, error) {
	// todo: erase
	// todo: validate input
	// todo: ref split into func?
	// todo: ref class ?

	var builder strings.Builder

	packedStringLength := utf8.RuneCountInString(packedString)
	if packedStringLength == 1 {
		return packedString, nil
	}

	var prevSymbol rune
	for _, symbol := range packedString {

		if unicode.IsLetter(prevSymbol) {
			count := 1
			if unicode.IsDigit(symbol) {
				count, _ = strconv.Atoi(string(symbol)) // error suppressed because symbol was checked before
			}

			builder.WriteString(strings.Repeat(string(prevSymbol), count))
		}

		prevSymbol = symbol
	}

	// process last
	lastSymbol, _ := utf8.DecodeLastRuneInString(packedString)
	if unicode.IsLetter(lastSymbol) {
		builder.WriteRune(lastSymbol)
	}

	return builder.String(), nil
}
