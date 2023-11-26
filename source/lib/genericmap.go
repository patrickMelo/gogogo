package lib

import (
	"math"
	"reflect"
	"strconv"
	"strings"
)

type GenericMap map[string]interface{}

// Creates a new, empty generic map.
func NewGenericMap() GenericMap {
	return make(GenericMap)
}

// Merges the map with another one (adding or replacing map keys with otherMap keys)
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
		switch typedValue := value.(type) {
		case string:
			return typedValue
		case int:
			return strconv.FormatInt(int64(typedValue), 10)
		case int32:
			return strconv.FormatInt(int64(typedValue), 10)
		case int64:
			return strconv.FormatInt(typedValue, 10)
		case uint:
			return strconv.FormatUint(uint64(typedValue), 10)
		case uint32:
			return strconv.FormatUint(uint64(typedValue), 10)
		case uint64:
			return strconv.FormatUint(typedValue, 10)
		case float32:
			return strconv.FormatFloat(float64(typedValue), 'f', -1, 32)
		case float64:
			return strconv.FormatFloat(typedValue, 'f', -1, 64)
		case bool:
			return strconv.FormatBool(typedValue)
		default:
			return defaultValue
		}
	}

	return defaultValue
}

// Returns a value from the map as an integer value. If the key does not exists the defaultValue is returned.
func (m GenericMap) GetInt(key string, defaultValue int64) int64 {
	if value, keyExists := m[key]; keyExists {
		switch typedValue := value.(type) {
		case string:
			if intValue, err := strconv.ParseInt(typedValue, 10, 64); err == nil {
				return intValue
			}

		case int:
			return int64(typedValue)
		case int32:
			return int64(typedValue)
		case int64:
			return typedValue
		case uint:
			return int64(typedValue)
		case uint32:
			return int64(typedValue)
		case uint64:
			return int64(typedValue)
		case float32:
			return int64(math.Round(float64(typedValue)))
		case float64:
			return int64(math.Round(typedValue))
		case bool:
			if typedValue {
				return 1
			} else {
				return 0
			}
		default:
			return defaultValue
		}
	}

	return defaultValue
}

// Returns a value from the map a float value. If the key does not exists the defaultValue is returned.
func (m GenericMap) GetFloat(key string, defaultValue float64) float64 {
	if value, keyExists := m[key]; keyExists {
		switch typedValue := value.(type) {
		case string:
			if floatValue, err := strconv.ParseFloat(typedValue, 64); err == nil {
				return floatValue
			}

		case int:
			return float64(typedValue)
		case int32:
			return float64(typedValue)
		case int64:
			return float64(typedValue)
		case uint:
			return float64(typedValue)
		case uint32:
			return float64(typedValue)
		case uint64:
			return float64(typedValue)
		case float32:
			return float64(typedValue)
		case float64:
			return typedValue
		case bool:
			if typedValue {
				return 1.0
			} else {
				return 0.0
			}
		default:
			return defaultValue
		}
	}

	return defaultValue
}

// Returns a value from the map as a boolean value. If the key does not exists the defaultValue is returned.
func (m GenericMap) GetBool(key string, defaultValue bool) bool {
	if value, keyExists := m[key]; keyExists {
		switch typedValue := value.(type) {
		case string:
			switch strings.ToLower(typedValue) {
			case "1", "on", "yes", "enable", "enabled":
				return true
			case "0", "off", "no", "disable", "disabled":
				return false
			}

		case int:
			return typedValue == 1
		case int32:
			return typedValue == 1
		case int64:
			return typedValue == 1
		case uint:
			return typedValue == 1
		case uint32:
			return typedValue == 1
		case uint64:
			return typedValue == 1
		case float32:
			return typedValue == 1.0
		case float64:
			return typedValue == 1.0
		case bool:
			return typedValue
		default:
			return defaultValue
		}
	}

	return defaultValue
}
