package hw03frequencyanalysis

import (
	"sort"
	"strings"
	"unicode/utf8"
)

func Top10(text string) ([]string, error) {
	validationError := validateText(text)
	if validationError != nil {
		return nil, validationError
	}

	if len(text) <= 0 {
		return []string{}, nil
	}

	frequency := map[string]int{}

	// split
	words := strings.Fields(text)
	for _, word := range words {
		frequency[word] += 1
	}

	// sort
	uniqueWords := make([]string, 0, len(frequency))
	for word := range frequency {
		uniqueWords = append(uniqueWords, word)
	}

	sort.SliceStable(uniqueWords, func(i, j int) bool {
		if frequency[uniqueWords[i]] == frequency[uniqueWords[j]] {
			return uniqueWords[i] < uniqueWords[j]
		}

		return frequency[uniqueWords[i]] > frequency[uniqueWords[j]]
	})

	return uniqueWords, nil
}

func validateText(textToValidate string) error {
	if !utf8.ValidString(textToValidate) {
		return InvalidUtf8StringError
	}

	return nil
}
