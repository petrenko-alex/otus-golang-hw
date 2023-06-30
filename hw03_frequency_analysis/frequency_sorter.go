package hw03frequencyanalysis

import "sort"

type FrequencySorter interface {
	SortFrequency(frequency Frequency) []string
}

type DescendingFrequencySorter struct{}

func (DescendingFrequencySorter) SortFrequency(frequency Frequency) []string {
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
