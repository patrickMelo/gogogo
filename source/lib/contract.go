package lib

import (
	"math"
	"regexp"
	"unicode/utf8"
)

const ContractValidationErrorsFieldName = "errors"

type ContractField interface {
	GetName() string
	IsRequired() bool
	Validate(value interface{}) ContractValidationErrorCode
}

type Contract struct {
	fields map[string]ContractField
}

func NewContract(fields []ContractField) (contract *Contract) {
	contract = &Contract{
		fields: make(map[string]ContractField),
	}

	for _, field := range fields {
		contract.fields[field.GetName()] = field
	}

	return
}

type ContractValidationErrorCode int

func (code ContractValidationErrorCode) String() string {
	switch code {
	case ContractFieldOK:
		return "OK"
	case ContractInvalidLength:
		return "InvalidLength"
	case ContractInvalidValue:
		return "InvalidValue"
	case ContractInvalidValueType:
		return "InvalidValueType"
	case ContractMissingRequiredField:
		return "MissingRequiredField"
	case ContractUnknownField:
		return "UnknownField"
	case ContractUnknownFieldType:
		return "UnknownFieldType"
	case ContractValueOutOfRange:
		return "ValueOutOfRange"
	}

	return "?"
}

const (
	ContractFieldOK ContractValidationErrorCode = iota
	ContractInvalidLength
	ContractInvalidValue
	ContractInvalidValueType
	ContractMissingRequiredField
	ContractUnknownField
	ContractUnknownFieldType
	ContractValueOutOfRange
)

type ContractStringField struct {
	ContractField

	name           string
	isRequired     bool
	minLength      int64
	maxLength      int64
	regex          *regexp.Regexp
	acceptedValues []string
}

func ContractString(name string) *ContractStringField {
	return &ContractStringField{
		name:           name,
		isRequired:     false,
		minLength:      -1,
		maxLength:      -1,
		regex:          nil,
		acceptedValues: nil,
	}
}

func (field *ContractStringField) GetName() string {
	return field.name
}

func (field *ContractStringField) IsRequired() bool {
	return field.isRequired
}

func (field *ContractStringField) Optional() *ContractStringField {
	field.isRequired = false
	return field
}

func (field *ContractStringField) Required() *ContractStringField {
	field.isRequired = true
	return field
}

func (field *ContractStringField) Length(min int64, max int64) *ContractStringField {
	field.minLength = min
	field.maxLength = max
	return field
}

func (field *ContractStringField) Validate(value interface{}) ContractValidationErrorCode {
	var stringValue, typeOK = value.(string)

	if !typeOK {
		return ContractInvalidValueType
	}

	if field.acceptedValues != nil {
		var valueFound = false

		for _, acceptedValue := range field.acceptedValues {
			if acceptedValue == stringValue {
				valueFound = true
				break
			}
		}

		if !valueFound {
			return ContractInvalidValue
		}
	}

	if (field.minLength > 0 && utf8.RuneCountInString(stringValue) < int(field.minLength)) || (field.maxLength > 0 && utf8.RuneCountInString(stringValue) > int(field.maxLength)) {
		return ContractInvalidLength
	}

	if field.regex != nil && !field.regex.MatchString(stringValue) {
		return ContractInvalidValue
	}

	return ContractFieldOK
}

func (field *ContractStringField) Regex(regex *regexp.Regexp) *ContractStringField {
	field.regex = regex
	return field
}

func (field *ContractStringField) Accept(values ...string) *ContractStringField {
	field.acceptedValues = values
	return field
}

type ContractIntegerField struct {
	ContractField

	name       string
	isRequired bool
	minValue   int64
	maxValue   int64
}

func ContractInteger(name string) *ContractIntegerField {
	return &ContractIntegerField{
		name:       name,
		isRequired: false,
		minValue:   math.MinInt64,
		maxValue:   math.MaxInt64,
	}
}

func (field *ContractIntegerField) GetName() string {
	return field.name
}

func (field *ContractIntegerField) IsRequired() bool {
	return field.isRequired
}

func (field *ContractIntegerField) Optional() *ContractIntegerField {
	field.isRequired = false
	return field
}

func (field *ContractIntegerField) Required() *ContractIntegerField {
	field.isRequired = true
	return field
}

func (field *ContractIntegerField) Min(min int64) *ContractIntegerField {
	field.minValue = min
	return field
}

func (field *ContractIntegerField) Max(max int64) *ContractIntegerField {
	field.maxValue = max
	return field
}

func (field *ContractIntegerField) Range(min int64, max int64) *ContractIntegerField {
	field.minValue = min
	field.maxValue = max
	return field
}

func (field *ContractIntegerField) Validate(value interface{}) ContractValidationErrorCode {
	var intValue, typeOK = value.(int64)

	if !typeOK {
		return ContractInvalidValueType
	}

	if (intValue < field.minValue) || (intValue > field.maxValue) {
		return ContractValueOutOfRange
	}

	return ContractFieldOK
}

type ContractFloatField struct {
	ContractField

	name       string
	isRequired bool
	minValue   float64
	maxValue   float64
}

func ContractFloat(name string) *ContractFloatField {
	return &ContractFloatField{
		name:       name,
		isRequired: false,
		minValue:   math.SmallestNonzeroFloat64,
		maxValue:   math.MaxFloat64,
	}
}

func (field *ContractFloatField) GetName() string {
	return field.name
}

func (field *ContractFloatField) IsRequired() bool {
	return field.isRequired
}

func (field *ContractFloatField) Optional() *ContractFloatField {
	field.isRequired = false
	return field
}

func (field *ContractFloatField) Required() *ContractFloatField {
	field.isRequired = true
	return field
}

func (field *ContractFloatField) Min(min float64) *ContractFloatField {
	field.minValue = min
	return field
}

func (field *ContractFloatField) Max(max float64) *ContractFloatField {
	field.maxValue = max
	return field
}

func (field *ContractFloatField) Range(min float64, max float64) *ContractFloatField {
	field.minValue = min
	field.maxValue = max
	return field
}

func (field *ContractFloatField) Validate(value interface{}) ContractValidationErrorCode {
	var floatValue, typeOK = value.(float64)

	if !typeOK {
		return ContractInvalidValueType
	}

	if (floatValue < field.minValue) || (floatValue > field.maxValue) {
		return ContractValueOutOfRange
	}

	return ContractFieldOK
}

type ContractValidationError struct {
	FieldName string                      `json:"fieldName"`
	ErrorCode ContractValidationErrorCode `json:"errorCode"`
}

func (contract *Contract) Validate(data GenericMap) (errors []ContractValidationError) {
	var errorCode ContractValidationErrorCode
	errors = make([]ContractValidationError, 0)

	for fieldName, field := range contract.fields {
		if !data.Has(fieldName) {
			if field.IsRequired() {
				errors = append(errors, ContractValidationError{
					FieldName: fieldName,
					ErrorCode: ContractMissingRequiredField,
				})
			}

			continue
		}

		errorCode = ContractUnknownFieldType

		switch typedField := field.(type) {
		case *ContractStringField:
			errorCode = typedField.Validate(data.GetString(fieldName, ""))

		case *ContractIntegerField:
			errorCode = typedField.Validate(data.GetInt(fieldName, math.MinInt64))

		case *ContractFloatField:
			errorCode = typedField.Validate(data.GetFloat(fieldName, math.SmallestNonzeroFloat64))
		}

		if errorCode != ContractFieldOK {
			errors = append(errors, ContractValidationError{
				FieldName: fieldName,
				ErrorCode: errorCode,
			})
		}
	}

	for fieldName := range data {
		if _, nameExists := contract.fields[fieldName]; !nameExists {
			errors = append(errors, ContractValidationError{
				FieldName: fieldName,
				ErrorCode: ContractUnknownField,
			})
		}
	}

	return
}
