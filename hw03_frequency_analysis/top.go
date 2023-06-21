package hw03frequencyanalysis

import (
	"sort"
	"strings"
	"unicode"
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

	const frequencyLimit = 10
	frequency := map[string]int{}

	// split
	words := strings.Fields(text)
	for _, word := range words {
		if word == "-" {
			continue
		}

		word = strings.ToLower(word)
		word = strings.TrimFunc(word, func(r rune) bool {
			return unicode.IsPunct(r)
		})
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

	// cut off more than 10
	if len(uniqueWords) > frequencyLimit {
		uniqueWords = uniqueWords[:frequencyLimit]
	}

	return uniqueWords, nil
}

func validateText(textToValidate string) error {
	if !utf8.ValidString(textToValidate) {
		return InvalidUtf8StringError
	}

	return nil
}
