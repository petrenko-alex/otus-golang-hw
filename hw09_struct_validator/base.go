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

func (v ValidationErrors) GetFields() []string {
	fields := make([]string, 0, len(v))

	for _, fieldName := range v {
		fields = append(fields, fieldName.Field)
	}

	return fields
}
