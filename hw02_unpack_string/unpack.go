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
	ErrInvalidEscaping  = UnpackError("Only digits and slash symbol can be escaped")

	regexNotAllowedSymbols = `[^a-zA-Zа-яА-Я0-9\\]`
	regexNotStartWithDigit = `^[^\d]`
	regexNumbers           = `[^\\]+\d\d.*`
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

	return buildUnpackedString(packedString), nil
}

func validateUnpackedString(input string) error {
	if len(input) == 0 {
		return nil
	}

	matched, err := regexp.MatchString(regexNotAllowedSymbols, input)
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

	return validateEscaping(input)
}

func validateEscaping(input string) error {
	runes := []rune(input)
	for i := 0; i < len(runes); i++ {
		if runes[i] != '\\' {
			continue
		}

		if i+1 >= len(runes) {
			return ErrInvalidEscaping
		}

		if !isValidEscapedSymbol(runes[i+1]) {
			return ErrInvalidEscaping
		}

		i++ // skip escaped symbol
	}

	return nil
}

func isValidEscapedSymbol(symbol rune) bool {
	return unicode.IsDigit(symbol) || symbol == '\\'
}

func buildUnpackedString(packedString string) string {
	var builder strings.Builder

	i := 0
	runes := []rune(packedString)
	for i < len(runes) {
		multiplier := 1
		currentSymbol := runes[i]

		if unicode.IsLetter(currentSymbol) {
			// get multiplier
			if i+1 < len(runes) && unicode.IsDigit(runes[i+1]) {
				multiplier = AtoiRune(runes[i+1])
				i++
			}

			builder.WriteString(strings.Repeat(string(currentSymbol), multiplier))
			i++
		} else if currentSymbol == '\\' {
			// start escaping
			if i+1 >= len(runes) {
				break
			}

			// get escaped symbol
			escapedSymbol := runes[i+1]

			// try to find multiplier
			if i+2 < len(runes) && unicode.IsDigit(runes[i+2]) {
				multiplier = AtoiRune(runes[i+2])
				i++
			}

			builder.WriteString(strings.Repeat(string(escapedSymbol), multiplier))
			i += 2
		}
	}

	return builder.String()
}

func AtoiRune(symbol rune) int {
	num, _ := strconv.Atoi(string(symbol)) // error suppressed because we suppose symbol was checked before
	return num
}
