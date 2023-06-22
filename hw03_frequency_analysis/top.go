package hw03frequencyanalysis

type TextWordFrequency interface {
	Top(text string) ([]string, error)
}

type GeneralTextWordFrequency struct {
	TextValidator    TextValidator
	FrequencyCounter FrequencyCounter
	FrequencySorter  FrequencySorter
	FrequencyLimiter FrequencyLimiter
}

func (f GeneralTextWordFrequency) Top(text string) ([]string, error) {
	validationError := f.TextValidator.ValidateText(text)
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
