package config

import (
	"gogogo/data"
	"os"
	"strings"
)

type EnvironmentProvider struct {
	Provider

	prefix string
}

func Environment(prefix string) *EnvironmentProvider {
	return &EnvironmentProvider{
		prefix: prefix + "_",
	}
}

func (provider *EnvironmentProvider) Load() (params data.GenericMap, err error) {
	var varData []string
	var key string
	var prefixLength = len(provider.prefix)

	params = data.NewGenericMap()

	for _, envVar := range os.Environ() {
		varData = strings.SplitN(envVar, "=", 2)

		if strings.Index(varData[0], provider.prefix) == 0 {
			key = strings.ToLower(strings.ReplaceAll(varData[0], "_", "."))[prefixLength:]
			params.Set(key, varData[1])
		}
	}

	return
}
