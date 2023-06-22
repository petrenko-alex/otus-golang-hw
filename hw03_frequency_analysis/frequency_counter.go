package hw03frequencyanalysis

import (
	"strings"
	"unicode"
)

type FrequencyCounter interface {
	calcWordsFrequency(text string) map[string]int
}

type PunctuationFrequencyCounter struct{}
type NonPunctuationFrequencyCounter struct{}

func (PunctuationFrequencyCounter) calcWordsFrequency(text string) map[string]int {
	frequency := map[string]int{}

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

	return frequency
}

func (NonPunctuationFrequencyCounter) calcWordsFrequency(text string) map[string]int {
	frequency := map[string]int{}

	words := strings.Fields(text)

	for _, word := range words {
		frequency[word] += 1
	}

	return frequency
}
