package contract

import (
	"fmt"
	"reflect"
	"regexp"
)

type StringField struct {
	Field
}

func String(name string) *StringField {
	var field = &StringField{
		Field: GenericField(name),
	}

	field.addValidator(&stringTypeValidator{})

	return field
}

type stringTypeValidator struct {
	FieldValidator
}

func (validator *stringTypeValidator) Validate(value interface{}) (err *ValidationError) {
	if reflect.ValueOf(value).Kind() != reflect.String {
		return &ValidationError{
			ErrorCode:    InvalidValueType,
			ErrorMessage: "value must be a string",
		}
	}

	return
}

type stringLengthValidator struct {
	FieldValidator

	min int64
	max int64
}

func (validator *stringLengthValidator) Validate(value interface{}) (err *ValidationError) {
	var stringLength = int64(len(value.(string)))

	if validator.min > 0 && stringLength < validator.min {
		return &ValidationError{
			ErrorCode:    InvalidLength,
			ErrorMessage: fmt.Sprintf("length must be greater or equals to %d", validator.min),
		}
	}

	if validator.max > 0 && stringLength > validator.max {
		return &ValidationError{
			ErrorCode:    InvalidLength,
			ErrorMessage: fmt.Sprintf("length must be less or equals to %d", validator.max),
		}
	}

	return
}

func (field *StringField) Length(min int64, max int64) *StringField {
	field.addValidator(&stringLengthValidator{
		min: min,
		max: max,
	})

	return field
}

type stringRegexValidator struct {
	FieldValidator

	regex *regexp.Regexp
}

func (validator *stringRegexValidator) Validate(value interface{}) (err *ValidationError) {
	if !validator.regex.MatchString(value.(string)) {
		return &ValidationError{
			ErrorCode:    InvalidValue,
			ErrorMessage: "value does not match regex",
		}
	}

	return
}

func (field *StringField) Regex(regex *regexp.Regexp) *StringField {
	field.addValidator(&stringRegexValidator{
		regex: regex,
	})

	return field
}

type stringAcceptValidator struct {
	FieldValidator

	acceptedValues []string
}

func (validator *stringAcceptValidator) Validate(value interface{}) (err *ValidationError) {
	var valueFound = false
	var stringValue = value.(string)

	for _, acceptedValue := range validator.acceptedValues {
		if acceptedValue == stringValue {
			valueFound = true
			break
		}
	}

	if !valueFound {
		return &ValidationError{
			ErrorCode:    InvalidValue,
			ErrorMessage: "value does not match accepted values",
		}
	}

	return
}

func (field *StringField) Accept(values ...string) *StringField {
	field.addValidator(&stringAcceptValidator{
		acceptedValues: values,
	})

	return field
}
