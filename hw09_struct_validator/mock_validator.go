package hw09structvalidator

type MockValidator struct{}

func (m MockValidator) ValidateValue(value interface{}) ValidationErrors {
	return ValidationErrors{}
}
