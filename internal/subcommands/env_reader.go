package subcommands

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/pivotal/pcf/lockfile"
)

type EnvReader struct {
	Lockfile string `short:"l" long:"lockfile" env:"ENVIRONMENT_LOCK_METADATA" description:"path to a lockfile"`
}

func (er *EnvReader) ReadRaw(result interface{}) error {
	if err := er.processParameters(); err != nil {
		return err
	}

	contents, err := ioutil.ReadFile(er.Lockfile)
	if err != nil {
		return err
	}
	return json.Unmarshal(contents, result)
}

func (er *EnvReader) Read() (lockfile.Lockfile, error) {
	if err := er.processParameters(); err != nil {
		return lockfile.Lockfile{}, err
	}

	return lockfile.FromFile(er.Lockfile)
}

func (er *EnvReader) processParameters() error {
	if er.Lockfile == "" {
		return fmt.Errorf("You must specify the lockfile path (--lockfile | -l) flag")
	}

	return nil
}
