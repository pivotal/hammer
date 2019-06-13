package subcommands

import (
	"fmt"

	"github.com/pivotal/pcf/internal/scripting"
)

type SSHCommand struct {
	Env  EnvReader `group:"environment"`
	File bool      `short:"f" long:"file" description:"write a script file but do not run it"`
}

func (c *SSHCommand) Execute(args []string) error {
	data, err := c.Env.Read()
	if err != nil {
		return err
	}

	fmt.Printf("Connecting to: %s\n", data.Name)

	sshCommand := fmt.Sprintf(`ssh -o IdentitiesOnly=yes -o StrictHostKeyChecking=no -i "${ssh_key_path}" -t ubuntu@"%s"`, data.OpsManager.IP.String())

	lines := []string{
		fmt.Sprintf(`ssh_key_path=$(mktemp)`),
		fmt.Sprintf(`echo "%s" >"$ssh_key_path"`, data.OpsManager.PrivateKey),
		fmt.Sprintf(`trap 'rm -f ${ssh_key_path}' EXIT`),
		fmt.Sprintf(`chmod 0600 "${ssh_key_path}"`),
		fmt.Sprintf(`creds="$(om -t %s -k -u %s -p %s curl -s -p %s)"`, data.OpsManager.URL.String(), data.OpsManager.Username, data.OpsManager.Password, boshCredsPath),
		fmt.Sprintf(`bosh="$(echo "$creds" | jq -r .credential | tr ' ' '\n' | grep '=')"`),
		fmt.Sprintf(`echo "$bosh"`),
		fmt.Sprintf(`shell="/usr/bin/env $(echo $bosh | tr '\n' ' ') bash -l"`),
		fmt.Sprintf(`%s "$shell"`, sshCommand),
	}

	dependencies := []string{"ssh", "om"}

	return scripting.RunScript(lines, dependencies, c.File)
}
