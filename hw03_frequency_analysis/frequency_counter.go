package hw03frequencyanalysis

import (
	"strings"
	"unicode"
)

type FrequencyCounter interface {
	CalcFrequency(text string) Frequency
}

type PunctuationFrequencyCounter struct{}
type NonPunctuationFrequencyCounter struct{}

func (PunctuationFrequencyCounter) CalcFrequency(text string) Frequency {
	frequency := Frequency{}

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

func (NonPunctuationFrequencyCounter) CalcFrequency(text string) Frequency {
	frequency := Frequency{}

	words := strings.Fields(text)

	for _, word := range words {
		frequency[word] += 1
	}

	return frequency
}
