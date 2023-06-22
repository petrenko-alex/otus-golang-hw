package hw03frequencyanalysis

import (
	"sort"
	"unicode/utf8"
)

func Top10(text string, freqCounter FrequencyCounter) ([]string, error) {
	validationError := validateText(text)
	if validationError != nil {
		return nil, validationError
	}

	if len(text) <= 0 {
		return []string{}, nil
	}

	frequency := freqCounter.calcWordsFrequency(text)
	uniqueWords := calcMostFrequentWords(frequency)
	uniqueWords = getTopFrequentWords(uniqueWords, 10)

	return uniqueWords, nil
}

func validateText(textToValidate string) error {
	if !utf8.ValidString(textToValidate) {
		return InvalidUtf8StringError
	}

	return nil
}

func calcMostFrequentWords(frequency map[string]int) []string {
	uniqueWords := make([]string, 0, len(frequency))
	for word := range frequency {
		uniqueWords = append(uniqueWords, word)
	}

	sort.SliceStable(uniqueWords, func(i, j int) bool {
		wordI, wordJ := uniqueWords[i], uniqueWords[j]

		if frequency[wordI] == frequency[wordJ] {
			return wordI < wordJ
		}

		return frequency[wordI] > frequency[wordJ]
	})

	return uniqueWords
}

func getTopFrequentWords(words []string, limit int) []string {
	if len(words) > limit {
		words = words[:limit]
	}

	return words
}
