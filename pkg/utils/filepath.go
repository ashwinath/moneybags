package utils

import (
	"os"
	"path/filepath"
	"strings"
)

func GetLocalRepoLocation() string {
	wd, _ := os.Getwd()
	for !strings.HasSuffix(wd, "moneybags") {
		wd = filepath.Dir(wd)
	}

	return wd
}
