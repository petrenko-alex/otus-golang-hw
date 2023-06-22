package hw03frequencyanalysis

import "unicode/utf8"

const InvalidUtf8TextError = TextValidationError("Text should be valid utf8")

type (
	Utf8Validator       struct{}
	TextValidationError string
)

type TextValidator interface {
	ValidateText(text string) error
}

func (e TextValidationError) Error() string {
	return string(e)
}

func (Utf8Validator) ValidateText(text string) error {
	if !utf8.ValidString(text) {
		return InvalidUtf8TextError
	}

	return nil
}
