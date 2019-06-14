package commands

import (
	"fmt"

	"github.com/pivotal/pcf/lockfile"
)

//go:generate counterfeiter . CFLoginScripter
type CFLoginScripter interface {
	Generate(data lockfile.Lockfile) []string
}

type CFLoginCommand struct {
	Lockfile string `short:"l" long:"lockfile" env:"ENVIRONMENT_LOCK_METADATA" description:"path to a lockfile"`
	File     bool   `short:"f" long:"file" description:"write a script file but do not run it"`

	CFLoginScripter CFLoginScripter
	Env             EnvReader
	ScriptRunner    ScriptRunner
}

func (c *CFLoginCommand) Execute(args []string) error {
	data, err := c.Env.Read(c.Lockfile)
	if err != nil {
		return err
	}

	fmt.Printf("Logging in to: %s\n", data.OpsManager.URL.String())

	cfLoginCommandLines := c.CFLoginScripter.Generate(data)

	return c.ScriptRunner.RunScript(cfLoginCommandLines, []string{"jq", "om", "cf"}, c.File)
}
