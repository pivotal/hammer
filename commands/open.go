package commands

import (
	"fmt"

	"github.com/pivotal/pcf/scripting"
)

type OpenCommand struct {
	Env  EnvReader `group:"environment"`
	File bool      `short:"f" long:"file" description:"write a script file but do not run it"`
	Show bool      `short:"s" long:"show" description:"only show the credentials"`
}

func (c *OpenCommand) Execute(args []string) error {
	data, err := c.Env.Read()
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

	lines := []string{
		fmt.Sprintf(`open "%s"`, data.OpsManager.URL.String()),
		fmt.Sprintf(`echo "%s" | pbcopy`, data.OpsManager.Password),
	}

	return scripting.RunScript(lines, []string{"open", "pbcopy"}, c.File)
}
