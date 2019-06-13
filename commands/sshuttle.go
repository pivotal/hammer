package commands

import (
	"fmt"

	"github.com/pivotal/pcf/scripting"
)

type SshuttleCommand struct {
	Env  EnvReader `group:"environment"`
	File bool      `short:"f" long:"file" description:"write a script file but do not run it"`
}

func (c *SshuttleCommand) Execute(args []string) error {
	data, err := c.Env.Read()
	if err != nil {
		return err
	}
	lines := []string{
		fmt.Sprintf(`ssh_key_path=$(mktemp)`),
		fmt.Sprintf(`echo "%s" >"$ssh_key_path"`, data.OpsManager.PrivateKey),
		fmt.Sprintf(`trap 'rm -f ${ssh_key_path}' EXIT`),
		fmt.Sprintf(`chmod 0600 "${ssh_key_path}"`),
		fmt.Sprintf(`sshuttle --ssh-cmd "ssh -o IdentitiesOnly=yes -i ${ssh_key_path}" -r ubuntu@"%s" %s %s %s`, data.OpsManager.IP.String(), data.OpsManager.CIDR.String(), data.PasCIDR.String(), data.ServicesCIDR.String()),
	}

	return scripting.RunScript(lines, []string{"jq", "om", "sshuttle"}, c.File)
}
