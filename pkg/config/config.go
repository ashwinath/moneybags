package config

import (
	"github.com/ashwinath/moneybags/pbgo/configpb"
	"github.com/ashwinath/moneybags/pkg/utils"
)

func New(path string) (*configpb.Config, error) {
	c := &configpb.Config{}
	if err := utils.UnmarshalYAML(path, c); err != nil {
		return nil, err
	}

	return c, nil
}
