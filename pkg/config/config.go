package config

import (
	"encoding/json"
	"os"

	"github.com/ashwinath/moneybags/pbgo/configpb"
	"sigs.k8s.io/yaml"
)

func New(path string) (*configpb.Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	j, err := yaml.YAMLToJSON(data)
	if err != nil {
		return nil, err
	}

	c := &configpb.Config{}
	if err = json.Unmarshal(j, c); err != nil {
		return nil, err
	}

	return c, nil
}
