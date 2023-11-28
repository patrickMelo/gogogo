package contract

import (
	"fmt"
	"reflect"
)

type IntegerField struct {
	Field
}

func Integer(name string) *IntegerField {
	var field = &IntegerField{
		Field: GenericField(name),
	}

	field.addValidator(&integerTypeValidator{})

	return field
}

type integerTypeValidator struct {
	FieldValidator
}

func (validator *integerTypeValidator) Validate(value interface{}) (err *ValidationError) {
	if reflect.ValueOf(value).Kind() != reflect.Int64 {
		return &ValidationError{
			ErrorCode:    InvalidValueType,
			ErrorMessage: "value must be an integer",
		}
	}

	return
}

type integerRangeValidator struct {
	FieldValidator

	min int64
	max int64
}

func (validator *integerRangeValidator) Validate(value interface{}) (err *ValidationError) {
	var intValue = value.(int64)

	if intValue < validator.min || intValue > validator.max {
		return &ValidationError{
			ErrorCode:    InvalidValue,
			ErrorMessage: fmt.Sprintf("value must be between %d and %d", validator.min, validator.max),
		}
	}

	return
}

func (field *IntegerField) Range(min int64, max int64) *IntegerField {
	field.addValidator(&integerRangeValidator{
		min: min,
		max: max,
	})

	return field
}
