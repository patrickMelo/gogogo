package config

import (
	"encoding/json"
	"gogogo/data"
	"os"
)

type FileProvider struct {
	Provider
}

func File() *FileProvider {
	return &FileProvider{}
}

func (provider *FileProvider) Load() (params data.GenericMap, err error) {
	var fileName = os.Args[0] + ".json"
	var fileData []byte

	if fileData, err = os.ReadFile(fileName); err != nil {
		return
	}

	params = data.NewGenericMap()
	if err = json.Unmarshal(fileData, &params); err != nil {
		return
	}

	params = params.Flatten(".")
	return
}
