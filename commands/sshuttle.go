package commands

type SshuttleCommand struct {
	Lockfile string `short:"l" long:"lockfile" env:"ENVIRONMENT_LOCK_METADATA" description:"path to a lockfile"`
	File     bool   `short:"f" long:"file" description:"write a script file but do not run it"`

	Env            EnvReader
	SshuttleRunner ToolRunner
}

func (c *SshuttleCommand) Execute(args []string) error {
	data, err := c.Env.Read(c.Lockfile)
	if err != nil {
		return err
	}

	return c.SshuttleRunner.Run(data, c.File)
}
