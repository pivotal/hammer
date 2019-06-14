package commands

import (
	"github.com/pivotal/pcf/lockfile"
)

//go:generate counterfeiter . SshuttleScripter
type SshuttleScripter interface {
	Generate(data lockfile.Lockfile) []string
}

type SshuttleCommand struct {
	Lockfile string `short:"l" long:"lockfile" env:"ENVIRONMENT_LOCK_METADATA" description:"path to a lockfile"`
	File     bool   `short:"f" long:"file" description:"write a script file but do not run it"`

	SshuttleScripter SshuttleScripter
	Env              EnvReader
	ScriptRunner     ScriptRunner
}

func (c *SshuttleCommand) Execute(args []string) error {
	data, err := c.Env.Read(c.Lockfile)
	if err != nil {
		return err
	}

	lines := c.SshuttleScripter.Generate(data)

	return c.ScriptRunner.RunScript(lines, []string{"jq", "om", "sshuttle"}, c.File)
}
