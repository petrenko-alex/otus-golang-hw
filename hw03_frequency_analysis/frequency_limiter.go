package hw03frequencyanalysis

type FrequencyLimiter interface {
	LimitFrequency(words []string) []string
}

type SimpleFrequencyLimiter struct {
	Limit int
}

func (l SimpleFrequencyLimiter) LimitFrequency(words []string) []string {
	if len(words) > l.Limit {
		words = words[:l.Limit]
	}

	return words
}
