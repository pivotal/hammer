package commands

import (
	"github.com/pivotal/pcf/lockfile"
)

//go:generate counterfeiter . BoshScripter
type BoshScripter interface {
	Generate(data lockfile.Lockfile, boshArgs []string) []string
}

type BoshCommand struct {
	Lockfile string `short:"l" long:"lockfile" env:"ENVIRONMENT_LOCK_METADATA" description:"path to a lockfile"`
	File     bool   `short:"f" long:"file" description:"write a script file but do not run it"`

	BoshScripter BoshScripter
	Env          EnvReader
	ScriptRunner ScriptRunner
}

func (c *BoshCommand) Execute(args []string) error {
	data, err := c.Env.Read(c.Lockfile)
	if err != nil {
		return err
	}

	boshCommandLines := c.BoshScripter.Generate(data, args)

	return c.ScriptRunner.RunScript(boshCommandLines, []string{"jq", "om", "ssh", "bosh"}, c.File)
}
