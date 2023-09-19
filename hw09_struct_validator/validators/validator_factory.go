package validators

import (
	"errors"
	"fmt"
	"github.com/petrenko-alex/otus-golang-hw/hw09_struct_validator/parsing"
	"reflect"
)

var (
	ErrFieldTypeArg         = errors.New("can't get field type")
	ErrValidationTagArg     = errors.New("can't parse validation tag")
	ErrUnsupportedFieldType = errors.New("can't validate field type")
)

type ValidatorFactory interface {
	// GetValidator return validator for value and validation tag.
	GetValidator(value interface{}, tag parsing.ValidationTag) (ValueValidator, error)
}

type FieldTypeValidatorFactory struct{}

func (f FieldTypeValidatorFactory) GetValidator(fieldType interface{}, validationTag string) (ValueValidator, error) {
	if _, ok := fieldType.(reflect.Kind); !ok {
		return nil, ErrFieldTypeArg
	}

	var tagParser parsing.ValidationTagParser = parsing.TagParser{
		RuleParser: parsing.BaseRuleParser{},
	}

	rules, err := tagParser.GetValidationRules(validationTag)
	if err != nil {
		return nil, fmt.Errorf(ErrValidationTagArg.Error()+": %w", err)
	}

	switch fieldType {
	case reflect.String, reflect.Int:
		return ScalarValueValidator{rules}, nil
	case reflect.Slice, reflect.Array:
		return SliceValueValidator{}, nil
	}

	return nil, ErrUnsupportedFieldType
}
