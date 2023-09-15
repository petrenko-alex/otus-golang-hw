package hw09structvalidator

type Validator interface {
	Validate(value interface{}) error
}

type ValidationError struct {
	Field string
	Err   error
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	panic("implement me")
}
