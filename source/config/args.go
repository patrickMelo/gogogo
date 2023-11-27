package config

import (
	"fmt"
	"gogogo/data"
	"os"
)

type ArgsProvider struct {
	Provider
}

func Args() *ArgsProvider {
	return &ArgsProvider{}
}

func (provider *ArgsProvider) Load() (params data.GenericMap, err error) {
	var key string

	params = data.NewGenericMap()

	for _, arg := range os.Args[1:] {
		if arg[0] == '-' {
			if key != "" {
				params.Set(key, true)
			}

			key = arg[1:]
			continue
		}

		if key == "" {
			return nil, fmt.Errorf("loose configuration value: %s", arg)
		}

		params.Set(key, arg)
		key = ""
	}

	return
}
