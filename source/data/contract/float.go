package contract

import (
	"fmt"
	"reflect"
)

type FloatField struct {
	Field
}

func Float(name string) *FloatField {
	var field = &FloatField{
		Field: GenericField(name),
	}

	field.addValidator(&floatTypeValidator{})

	return field
}

type floatTypeValidator struct {
	FieldValidator
}

func (validator *floatTypeValidator) Validate(value interface{}) (err *ValidationError) {
	if reflect.ValueOf(value).Kind() != reflect.Float64 {
		return &ValidationError{
			ErrorCode:    InvalidValueType,
			ErrorMessage: "value must be a float",
		}
	}

	return
}

type floatRangeValidator struct {
	FieldValidator

	min float64
	max float64
}

func (validator *floatRangeValidator) Validate(value interface{}) (err *ValidationError) {
	var intValue = value.(float64)

	if intValue < validator.min || intValue > validator.max {
		return &ValidationError{
			ErrorCode:    InvalidValue,
			ErrorMessage: fmt.Sprintf("value must be between %f and %f", validator.min, validator.max),
		}
	}

	return
}

func (field *FloatField) Range(min float64, max float64) *FloatField {
	field.addValidator(&floatRangeValidator{
		min: min,
		max: max,
	})

	return field
}
