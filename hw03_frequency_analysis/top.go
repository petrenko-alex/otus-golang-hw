package hw03frequencyanalysis

import "unicode/utf8"

func Top10(text string) ([]string, error) {
	validationError := validateText(text)
	if validationError != nil {
		return nil, validationError
	}

	return nil, nil
}

func validateText(textToValidate string) error {
	if !utf8.ValidString(textToValidate) {
		return InvalidUtf8StringError
	}

	return nil
}
