package commands

import (
	"fmt"

	"github.com/pivotal/pcf/lockfile"
)

type envReader struct{}

func NewEnvReader() envReader {
	return envReader{}
}

func (er *envReader) Read(lockfilePath string) (lockfile.Lockfile, error) {
	if lockfilePath == "" {
		return lockfile.Lockfile{}, fmt.Errorf("You must specify the lockfile path (--lockfile | -l) flag")
	}
	return lockfile.FromFile(lockfilePath)
}
