package commands

import (
	"fmt"

	"github.com/pivotal/pcf/lockfile"
)

//go:generate counterfeiter . SSHScripter
type SSHScripter interface {
	Generate(data lockfile.Lockfile) []string
}

type SSHCommand struct {
	Lockfile string `short:"l" long:"lockfile" env:"ENVIRONMENT_LOCK_METADATA" description:"path to a lockfile"`
	File     bool   `short:"f" long:"file" description:"write a script file but do not run it"`

	SSHScripter  SSHScripter
	Env          EnvReader
	ScriptRunner ScriptRunner
}

func (c *SSHCommand) Execute(args []string) error {
	data, err := c.Env.Read(c.Lockfile)
	if err != nil {
		return err
	}

	fmt.Printf("Connecting to: %s\n", data.Name)

	lines := c.SSHScripter.Generate(data)

	dependencies := []string{"ssh", "om"}

	return c.ScriptRunner.RunScript(lines, dependencies, c.File)
}
