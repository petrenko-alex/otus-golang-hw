package validators

import (
	"errors"
	"github.com/petrenko-alex/otus-golang-hw/hw09_struct_validator/parsing"
	"reflect"
)

type ValidatorFactory interface {
	// GetValidator return validator for value and validation tag.
	GetValidator(value interface{}, tag parsing.ValidationTag) (ValueValidator, error)
}

type FieldTypeValidatorFactory struct{}

func (f FieldTypeValidatorFactory) GetValidator(fieldType interface{}, validationTag string) (ValueValidator, error) {
	if _, ok := fieldType.(reflect.Kind); !ok {
		// todo: return err
		return nil, errors.New("tmp err")
	}

	var tagParser parsing.ValidationTagParser = parsing.TagParser{
		Factory: parsing.BaseRuleParser{},
	}

	rules, err := tagParser.GetValidationRules(validationTag)
	if err != nil {
		// todo: wrap err
		return nil, err
	}

	switch fieldType {
	case reflect.String:
	case reflect.Int:
		return ScalarValueValidator{rules}, nil
	case reflect.Slice:
	case reflect.Array:
		return SliceValueValidator{}, nil
	}

	return nil, errors.New("tmp err") // todo: return err not supported (why not default?)
}
