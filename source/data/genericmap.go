package data

import (
	"reflect"
)

type GenericMap map[string]interface{}

// Creates a new, empty generic map.
func NewGenericMap() GenericMap {
	return make(GenericMap)
}

// Merges a map with another one (adding or replacing map keys with otherMap keys)
func (m GenericMap) MergeWith(otherMap GenericMap) GenericMap {
	for key, value := range otherMap {
		m[key] = value
	}

	return m
}

// Checkes whether the map matches another one (if the map contains keys with the same value from otherMap)
func (m GenericMap) Matches(otherMap GenericMap) bool {
	for key, value := range otherMap {
		if localValue, keyExists := m[key]; keyExists {
			if localValue != value {
				return false
			}
		} else {
			return false
		}
	}

	return true
}

// Sets the map entries from a struct fields values.
func (m GenericMap) FromStruct(s interface{}) GenericMap {
	var value = reflect.ValueOf(s)

	if value.Kind() == reflect.Pointer {
		value = value.Elem()
	}

	if !value.IsValid() {
		return m
	}

	var valueType = value.Type()

	if valueType.Kind() != reflect.Struct {
		return m
	}

	for fieldIndex := 0; fieldIndex < valueType.NumField(); fieldIndex++ {
		var fieldTags = valueType.Field(fieldIndex).Tag
		var fieldKey = fieldTags.Get("key")

		if fieldKey == "" {
			fieldKey = valueType.Field(fieldIndex).Name
		}

		var fieldValue = value.Field(fieldIndex)

		switch fieldValue.Type().Kind() {
		case reflect.String:
			m.Set(fieldKey, value.Field(fieldIndex).String())

		case reflect.Bool:
			m.Set(fieldKey, value.Field(fieldIndex).Bool())

		case reflect.Int,
			reflect.Int8,
			reflect.Int16,
			reflect.Int32,
			reflect.Int64:
			m.Set(fieldKey, value.Field(fieldIndex).Int())

		case reflect.Uint,
			reflect.Uint8,
			reflect.Uint16,
			reflect.Uint32,
			reflect.Uint64:
			m.Set(fieldKey, value.Field(fieldIndex).Uint())

		case reflect.Float32,
			reflect.Float64:
			m.Set(fieldKey, value.Field(fieldIndex).Float())
		}
	}

	return m
}

// Sets a struct fields values from the map entries.
func (m GenericMap) ToStruct(s interface{}) GenericMap {
	var value = reflect.ValueOf(s)

	if value.Kind() == reflect.Pointer {
		value = value.Elem()
	}

	if !value.IsValid() {
		return m
	}

	var valueType = value.Type()

	if valueType.Kind() != reflect.Struct {
		return m
	}

	for fieldIndex := 0; fieldIndex < valueType.NumField(); fieldIndex++ {
		if !value.CanSet() {
			continue
		}

		var fieldTags = valueType.Field(fieldIndex).Tag
		var fieldKey = fieldTags.Get("key")

		if fieldKey == "" {
			fieldKey = valueType.Field(fieldIndex).Name
		}

		if !m.Has(fieldKey) {
			continue
		}

		var fieldValue = value.Field(fieldIndex)

		switch fieldValue.Type().Kind() {
		case reflect.String:
			fieldValue.SetString(m.GetString(fieldKey, ""))

		case reflect.Bool:
			fieldValue.SetBool(m.GetBool(fieldKey, false))

		case reflect.Int,
			reflect.Int8,
			reflect.Int16,
			reflect.Int32,
			reflect.Int64:
			fieldValue.SetInt(m.GetInt(fieldKey, 0))

		case reflect.Uint,
			reflect.Uint8,
			reflect.Uint16,
			reflect.Uint32,
			reflect.Uint64:
			fieldValue.SetUint(uint64(m.GetInt(fieldKey, 0)))

		case reflect.Float32,
			reflect.Float64:
			fieldValue.SetFloat(m.GetFloat(fieldKey, 0.0))
		}
	}

	return m
}

// Sets a value for an entry (the entry is added if it does not exists).
func (m GenericMap) Set(key string, value interface{}) GenericMap {
	m[key] = value
	return m
}

// Removes an entry from the map.
func (m GenericMap) Unset(key string) GenericMap {
	delete(m, key)
	return m
}

// Checks if the map has the specified entry.
func (m GenericMap) Has(key string) bool {
	var _, keyExists = m[key]
	return keyExists
}

// Returns a value from the map. If the key does not exists the defaultValue is returned.
func (m GenericMap) Get(key string, defaultValue interface{}) interface{} {
	if value, keyExists := m[key]; keyExists {
		return value
	}

	return defaultValue
}

// Returns a value from the map as a string value. If the key does not exists the defaultValue is returned.
func (m GenericMap) GetString(key string, defaultValue string) string {
	if value, keyExists := m[key]; keyExists {
		return ToString(value, defaultValue)
	}

	return defaultValue
}

// Returns a value from the map as an integer value. If the key does not exists the defaultValue is returned.
func (m GenericMap) GetInt(key string, defaultValue int64) int64 {
	if value, keyExists := m[key]; keyExists {
		return ToInt(value, defaultValue)
	}

	return defaultValue
}

// Returns a value from the map a float value. If the key does not exists the defaultValue is returned.
func (m GenericMap) GetFloat(key string, defaultValue float64) float64 {
	if value, keyExists := m[key]; keyExists {
		return ToFloat(value, defaultValue)
	}

	return defaultValue
}

// Returns a value from the map as a boolean value. If the key does not exists the defaultValue is returned.
func (m GenericMap) GetBool(key string, defaultValue bool) bool {
	if value, keyExists := m[key]; keyExists {
		return ToBool(value, defaultValue)
	}

	return defaultValue
}

func flattenValues(path string, separator string, values GenericMap) GenericMap {
	var flatValues = NewGenericMap()

	for key, value := range values {
		switch typedValue := value.(type) {
		case map[string]interface{}:
			flatValues.MergeWith(flattenValues(path+key+separator, separator, typedValue))

		default:
			flatValues.Set(path+key, value)
		}
	}

	return flatValues
}

func (m GenericMap) Flatten(separator string) GenericMap {
	return flattenValues("", ".", m)
}
