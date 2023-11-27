package data

import (
	"math"
	"strconv"
	"strings"
)

func ToString(value interface{}, defaultValue string) string {
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
	}

	return defaultValue
}

func ToInt(value interface{}, defaultValue int64) int64 {
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
	}

	return defaultValue
}

func ToFloat(value interface{}, defaultValue float64) float64 {
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
	}

	return defaultValue
}

func ToBool(value interface{}, defaultValue bool) bool {
	switch typedValue := value.(type) {
	case string:
		switch strings.ToLower(typedValue) {
		case "1", "on", "yes", "enable", "enabled", "true":
			return true
		case "0", "off", "no", "disable", "disabled", "false":
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
	}

	return defaultValue
}
