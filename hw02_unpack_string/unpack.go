package hw02unpackstring

import (
	"regexp"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
)

type UnpackError string

const (
	ErrStartsWithDigits = UnpackError("Not allowed to start with digits")
	ErrHasNumbers       = UnpackError("Numbers are not allowed, only digits")
	ErrInvalidChars     = UnpackError("Only digits and letters allowed")

	regexNotLettersAndDigits = `[^a-zA-Zа-яА-Я0-9]`
	regexNotStartWithDigit   = `^[^\d]`
	regexNumbers             = `.*\d\d.*`
)

func Unpack(packedString string) (string, error) {
	validationError := validateUnpackedString(packedString)
	if validationError != nil {
		return "", validationError
	}

	packedStringLength := utf8.RuneCountInString(packedString)
	if packedStringLength == 1 {
		return packedString, nil
	}

	var prevSymbol rune
	var builder strings.Builder
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

func validateUnpackedString(input string) error {
	if len(input) == 0 {
		return nil
	}

	matched, err := regexp.MatchString(regexNotLettersAndDigits, input)
	if matched || err != nil {
		return ErrInvalidChars
	}

	matched, err = regexp.MatchString(regexNotStartWithDigit, input)
	if !matched || err != nil {
		return ErrStartsWithDigits
	}

	matched, err = regexp.MatchString(regexNumbers, input)
	if matched || err != nil {
		return ErrHasNumbers
	}

	return nil
}

func (e UnpackError) Error() string {
	return string(e)
}
