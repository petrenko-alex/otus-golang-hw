package validators

import (
	"github.com/petrenko-alex/otus-golang-hw/hw09_struct_validator/parsing"
	"reflect"
)

// TODO: WIP
// todo: interface
// todo: make more abstract, pass field?

func GetValidator(fieldType reflect.Kind, validationTag string) (ValueValidator, error) {
	var tagParser parsing.ValidationTagParser = parsing.TagParser{
		Factory: parsing.BaseRuleParser{},
	}

	rules, err := tagParser.GetValidationRules(validationTag)
	if err != nil {
		// todo: return err, wrap
	}

	switch fieldType {
	case reflect.String:
	case reflect.Int:
		return ScalarValueValidator{rules}, nil
	case reflect.Slice:
	case reflect.Array:
		return SliceValueValidator{}, nil
	default:
		// todo: err field not supported
	}

	return nil, nil // todo: return err not supported (why not default?)
}
