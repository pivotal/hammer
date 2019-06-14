package commands

import (
	"fmt"
)

type SSHCommand struct {
	Lockfile string `short:"l" long:"lockfile" env:"ENVIRONMENT_LOCK_METADATA" description:"path to a lockfile"`
	File     bool   `short:"f" long:"file" description:"write a script file but do not run it"`

	Env       EnvReader
	SSHRunner ToolRunner
}

func (c *SSHCommand) Execute(args []string) error {
	data, err := c.Env.Read(c.Lockfile)
	if err != nil {
		return err
	}

	fmt.Printf("Connecting to: %s\n", data.Name)

	return c.SSHRunner.Run(data, c.File)
}
