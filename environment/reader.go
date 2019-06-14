package environment

import (
	"fmt"
)

type Reader struct{}

func (er *Reader) Read(targetConfigPath string) (Config, error) {
	if targetConfigPath == "" {
		return Config{}, fmt.Errorf("You must specify the target environment config path (--target | -t) flag")
	}
	return FromFile(targetConfigPath)
}
