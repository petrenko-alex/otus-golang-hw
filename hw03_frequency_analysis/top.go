package hw03frequencyanalysis

type TextWordFrequency interface {
	Top(text string) ([]string, error)
}

type GeneralTextWordFrequency struct {
	textValidator    TextValidator
	frequencyCounter FrequencyCounter
	frequencySorter  FrequencySorter
	frequencyLimiter FrequencyLimiter
}

func NewPunctuationTextWordFrequency() TextWordFrequency {
	return GeneralTextWordFrequency{
		textValidator:    Utf8Validator{},
		frequencyCounter: PunctuationFrequencyCounter{},
		frequencySorter:  DescendingFrequencySorter{},
		frequencyLimiter: SimpleFrequencyLimiter{Limit: 10},
	}
}

func NewNonPunctuationTextWordFrequency() TextWordFrequency {
	return GeneralTextWordFrequency{
		textValidator:    Utf8Validator{},
		frequencyCounter: NonPunctuationFrequencyCounter{},
		frequencySorter:  DescendingFrequencySorter{},
		frequencyLimiter: SimpleFrequencyLimiter{Limit: 10},
	}
}

func (f GeneralTextWordFrequency) Top(text string) ([]string, error) {
	validationError := f.textValidator.ValidateText(text)
	if validationError != nil {
		return nil, validationError
	}

	if len(text) == 0 {
		return []string{}, nil
	}

	frequency := f.frequencyCounter.CalcFrequency(text)
	top := f.frequencySorter.SortFrequency(frequency)
	top = f.frequencyLimiter.LimitFrequency(top)

	return top, nil
}
