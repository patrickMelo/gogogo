package config

import (
	"gogogo/data"
	"gogogo/log"
	"strings"
)

type Provider interface {
	Load() (data.GenericMap, error)
}

func Load(providers ...Provider) (err error) {
	var providerParams data.GenericMap
	var allParams = data.NewGenericMap()

	for _, provider := range providers {
		if providerParams, err = provider.Load(); err != nil {
			log.Error(_logTag, err)
		}

		allParams.MergeWith(providerParams)
	}

	_parameters = data.NewGenericMap()

	for key, value := range allParams {
		_parameters[strings.ToLower(key)] = value
	}

	return nil
}

func Get(paramName string, defaultValue interface{}) interface{} {
	return _parameters.Get(strings.ToLower(paramName), defaultValue)
}

func GetString(paramName string, defaultValue string) string {
	return _parameters.GetString(strings.ToLower(paramName), defaultValue)
}

func GetInt(paramName string, defaultValue int64) int64 {
	return _parameters.GetInt(strings.ToLower(paramName), defaultValue)
}

func GetFloat(paramName string, defaultValue float64) float64 {
	return _parameters.GetFloat(strings.ToLower(paramName), defaultValue)
}

func GetBool(paramName string, defaultValue bool) bool {
	return _parameters.GetBool(strings.ToLower(paramName), defaultValue)
}

const (
	_logTag string = "config"
)

var (
	_parameters data.GenericMap = data.NewGenericMap()
)
