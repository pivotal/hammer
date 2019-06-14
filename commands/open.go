package commands

import (
	"fmt"
)

type OpenCommand struct {
	Lockfile string `short:"l" long:"lockfile" env:"ENVIRONMENT_LOCK_METADATA" description:"path to a lockfile"`
	File     bool   `short:"f" long:"file" description:"write a script file but do not run it"`
	Show     bool   `short:"s" long:"show" description:"only show the credentials"`

	Env        EnvReader
	OpenRunner ToolRunner
}

func (c *OpenCommand) Execute(args []string) error {
	data, err := c.Env.Read(c.Lockfile)
	if err != nil {
		return err
	}

	if c.Show {
		fmt.Printf("%s\n", data.OpsManager.URL.String())
		fmt.Printf("username: %s\n", data.OpsManager.Username)
		fmt.Printf("password: %s\n", data.OpsManager.Password)
		return nil
	}

	fmt.Printf("Opening: %s\n", data.OpsManager.URL.String())
	fmt.Printf("Username is: %s\n", data.OpsManager.Username)

	fmt.Println("Password is in the clipboard")

	return c.OpenRunner.Run(data, c.File)
}
