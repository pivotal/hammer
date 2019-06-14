package commands

import (
	"fmt"

	"github.com/pivotal/pcf/lockfile"
)

//go:generate counterfeiter . OMScripter
type OMScripter interface {
	Generate(data lockfile.Lockfile, omArgs []string) []string
}

type OMCommand struct {
	Lockfile string `short:"l" long:"lockfile" env:"ENVIRONMENT_LOCK_METADATA" description:"path to a lockfile"`
	File     bool   `short:"f" long:"file" description:"write a script file but do not run it"`

	OMScripter   OMScripter
	Env          EnvReader
	ScriptRunner ScriptRunner
}

func (c *OMCommand) Execute(args []string) error {
	data, err := c.Env.Read(c.Lockfile)
	if err != nil {
		return err
	}

	if len(args) > 0 {
		omCommandLines := c.OMScripter.Generate(data, args)
		return c.ScriptRunner.RunScript(omCommandLines, []string{"om"}, c.File)
	}

	fmt.Printf("export OM_TARGET=%s\n", data.OpsManager.URL.String())
	fmt.Printf("export OM_USERNAME=%s\n", data.OpsManager.Username)
	fmt.Printf("export OM_PASSWORD=%s\n", data.OpsManager.Password)
	return nil
}
