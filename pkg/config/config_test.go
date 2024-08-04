package config

import (
	"os"
	"path"
	"testing"

	"github.com/ashwinath/moneybags/pbgo/configpb"
	"github.com/stretchr/testify/assert"
)

func TestNewConfig(t *testing.T) {
	p, err := os.Getwd()
	assert.Nil(t, err)

	p = path.Join(p, "./testdata/config.yaml")

	c, err := New(p)
	assert.Nil(t, err)

	assert.Equal(t, &configpb.Config{
		PostgresDb: &configpb.PostgresDB{
			Host:     "127.0.0.1",
			User:     "postgres",
			Password: "very_secure",
			DbName:   "postgres",
			Port:     5432,
		},
	}, c)
}
