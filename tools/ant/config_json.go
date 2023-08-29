package antnet

import (
	"encoding/json"
	"os"
)

func ReadConfigFromJson(path string, v interface{}) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return ErrFileRead
	}

	err = json.Unmarshal(data, v)
	if err != nil {
		return ErrJsonUnPack
	}
	return nil
}
