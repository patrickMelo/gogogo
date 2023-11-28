package contract

import (
	"gogogo/data"
)

type ValidationErrorCode int

const (
	FieldOK ValidationErrorCode = iota
	InvalidLength
	InvalidValue
	InvalidValueType
	MissingRequiredField
	UnknownField
	ValueOutOfRange
)

func (code ValidationErrorCode) String() string {
	switch code {
	case FieldOK:
		return "OK"
	case InvalidLength:
		return "InvalidLength"
	case InvalidValue:
		return "InvalidValue"
	case InvalidValueType:
		return "InvalidValueType"
	case MissingRequiredField:
		return "MissingRequiredField"
	case UnknownField:
		return "UnknownField"
	case ValueOutOfRange:
		return "ValueOutOfRange"
	}

	return "?"
}

type ValidationError struct {
	ErrorCode    ValidationErrorCode `json:"errorCode"`
	ErrorMessage string              `json:"errorMessage"`
}

type FieldValidator interface {
	Validate(value interface{}) *ValidationError
}

type Field struct {
	name       string
	isRequired bool
	validators []FieldValidator
}

func GenericField(name string) Field {
	return Field{
		name:       name,
		isRequired: false,
		validators: make([]FieldValidator, 0),
	}
}

func (field *Field) Name() string {
	return field.name
}

func (field *Field) Validate(value interface{}) *ValidationError {
	for _, validator := range field.validators {
		if validatorError := validator.Validate(value); validatorError != nil {
			return validatorError
		}
	}

	return nil
}

func (field *Field) Optional() *Field {
	field.isRequired = false
	return field
}

func (field *Field) Required() *Field {
	field.isRequired = true
	return field
}

func (field *Field) IsRequired() bool {
	return field.isRequired
}

func (field *Field) addValidator(validator FieldValidator) {
	field.validators = append(field.validators, validator)
}

type Contract struct {
	fields map[string]*Field
}

func New(fields ...*Field) (contract *Contract) {
	contract = &Contract{
		fields: make(map[string]*Field),
	}

	for _, field := range fields {
		contract.fields[field.Name()] = field
	}

	return
}

func (contract *Contract) Validate(data data.GenericMap) (errors map[string]*ValidationError) {
	errors = make(map[string]*ValidationError, 0)

	for fieldName, field := range contract.fields {
		if !data.Has(fieldName) {
			if field.IsRequired() {
				errors[fieldName] = &ValidationError{
					ErrorCode:    MissingRequiredField,
					ErrorMessage: "missing required field",
				}
			}

			continue
		}

		var fieldError = field.Validate(data.Get(fieldName, nil))

		if fieldError != nil {
			errors[fieldName] = fieldError
		}
	}

	for fieldName := range data {
		if _, nameExists := contract.fields[fieldName]; !nameExists {
			errors[fieldName] = &ValidationError{
				ErrorCode:    UnknownField,
				ErrorMessage: "unknown field",
			}
		}
	}

	return
}
