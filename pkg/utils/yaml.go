package utils

import (
	"encoding/json"
	"os"

	"sigs.k8s.io/yaml"
)

func UnmarshalYAML(path string, obj interface{}) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	j, err := yaml.YAMLToJSON(data)
	if err != nil {
		return err
	}

	if err = json.Unmarshal(j, obj); err != nil {
		return err
	}

	return nil
}
