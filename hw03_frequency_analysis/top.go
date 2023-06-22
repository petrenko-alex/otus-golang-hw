package hw03frequencyanalysis

import (
	"unicode/utf8"
)

type TextWordFrequency interface {
	Top(text string) ([]string, error)
}

type GeneralTextWordFrequency struct {
	FrequencyCounter FrequencyCounter
	FrequencySorter  FrequencySorter
	FrequencyLimiter FrequencyLimiter
}

func (f GeneralTextWordFrequency) Top(text string) ([]string, error) {
	validationError := validateText(text)
	if validationError != nil {
		return nil, validationError
	}

	if len(text) <= 0 {
		return []string{}, nil
	}

	frequency := f.FrequencyCounter.CalcFrequency(text)
	top := f.FrequencySorter.SortFrequency(frequency)
	top = f.FrequencyLimiter.LimitFrequency(top)

	return top, nil
}

func validateText(textToValidate string) error {
	if !utf8.ValidString(textToValidate) {
		return InvalidUtf8StringError
	}

	return nil
}
