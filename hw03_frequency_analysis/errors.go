package hw03frequencyanalysis

type Top10Error string

func (e Top10Error) Error() string {
	return string(e)
}

const (
	InvalidUtf8StringError = Top10Error("String should be valid utf8")
)
