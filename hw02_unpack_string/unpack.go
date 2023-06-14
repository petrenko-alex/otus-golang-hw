package hw02unpackstring

import (
	"regexp"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
)

// TODO: linter and styler

type UnpackError string

const (
	ErrStartsWithDigits = UnpackError("Not allowed to start with digits")
	ErrHasNumbers       = UnpackError("Numbers are not allowed, only digits")
	ErrInvalidChars     = UnpackError("Only digits and letters allowed")

	regexNotLettersAndDigits = `[^a-zA-Zа-яА-Я0-9]`
	regexNotStartWithDigit   = `^[^\d]`
	regexNumbers             = `.*\d\d.*`
)

func (e UnpackError) Error() string {
	return string(e)
}

func Unpack(packedString string) (string, error) {
	validationError := validateUnpackedString(packedString)
	if validationError != nil {
		return "", validationError
	}

	packedStringLength := utf8.RuneCountInString(packedString)
	if packedStringLength == 1 {
		return packedString, nil
	}

	var builder strings.Builder
	var unpackedString = buildUnpackedString(packedString, &builder)

	return unpackedString, nil
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

func buildUnpackedString(packedString string, builder *strings.Builder) string {
	runes := []rune(packedString)
	i := 0

	for i < len(runes) {
		currentSymbol := runes[i]

		if unicode.IsLetter(currentSymbol) {
			count := 1
			if i+1 < len(runes) && unicode.IsDigit(runes[i+1]) { // get multiplier
				count, _ = strconv.Atoi(string(runes[i+1])) // error suppressed because symbol was checked before
				i++                                         // skip multiplier
			}

			builder.WriteString(strings.Repeat(string(currentSymbol), count))
			i++
		}
	}

	return builder.String()
}
