package commands

import (
	"fmt"
	"strings"

	"github.com/pivotal/pcf/internal/scripting"
)

type BoshCommand struct {
	Env  EnvReader `group:"environment"`
	File bool      `short:"f" long:"file" description:"write a script file but do not run it"`
}

func (c *BoshCommand) Execute(args []string) error {
	data, err := c.Env.Read()
	if err != nil {
		return err
	}

	lines := []string{
		fmt.Sprintf(`ssh_key_path=$(mktemp)`),
		fmt.Sprintf(`echo "%s" >"$ssh_key_path"`, data.OpsManager.PrivateKey),
		fmt.Sprintf(`chmod 0600 "${ssh_key_path}"`),

		fmt.Sprintf(`bosh_ca_path=$(mktemp)`),
		fmt.Sprintf(`ssh -o IdentitiesOnly=yes -o StrictHostKeyChecking=no -i "${ssh_key_path}" ubuntu@"%s" cat /var/tempest/workspaces/default/root_ca_certificate 1>${bosh_ca_path} 2>/dev/null`, data.OpsManager.IP.String()),
		fmt.Sprintf(`chmod 0600 "${bosh_ca_path}"`),

		fmt.Sprintf(`creds="$(om -t %s -k -u %s -p %s curl -s -p %s)"`, data.OpsManager.URL.String(), data.OpsManager.Username, data.OpsManager.Password, boshCredsPath),
		fmt.Sprintf(`bosh_all="$(echo "$creds" | jq -r .credential | tr ' ' '\n' | grep '=')"`),

		fmt.Sprintf(`bosh_client="$(echo $bosh_all | tr ' ' '\n' | grep 'BOSH_CLIENT=')"`),
		fmt.Sprintf(`bosh_env="$(echo $bosh_all | tr ' ' '\n' | grep 'BOSH_ENVIRONMENT=')"`),
		fmt.Sprintf(`bosh_secret="$(echo $bosh_all | tr ' ' '\n' | grep 'BOSH_CLIENT_SECRET=')"`),
		fmt.Sprintf(`bosh_ca_cert="BOSH_CA_CERT=$bosh_ca_path"`),
		fmt.Sprintf(`bosh_proxy="BOSH_ALL_PROXY=ssh+socks5://ubuntu@%s:22?private-key=${ssh_key_path}"`, data.OpsManager.IP.String()),
	}

	if len(args) > 0 {
		lines = append(
			lines,
			fmt.Sprintf(`trap 'rm -f ${ssh_key_path}' EXIT`),
			fmt.Sprintf(`trap 'rm -f ${bosh_ca_path}' EXIT`),
			fmt.Sprintf(`/usr/bin/env $bosh_client $bosh_env $bosh_secret $bosh_ca_cert $bosh_proxy bosh %s`, strings.Join(args, " ")),
		)
	} else {
		lines = append(
			lines,
			fmt.Sprintf(`echo "export BOSH_ENV_NAME=%s"`, data.Name),
			fmt.Sprintf(`echo "export $bosh_client"`),
			fmt.Sprintf(`echo "export $bosh_env"`),
			fmt.Sprintf(`echo "export $bosh_secret"`),
			fmt.Sprintf(`echo "export $bosh_ca_cert"`),
			fmt.Sprintf(`echo "export $bosh_proxy"`),
			fmt.Sprintf(`echo "export CREDHUB_SERVER=\"\${BOSH_ENVIRONMENT}:8844\""`),
			fmt.Sprintf(`echo "export CREDHUB_PROXY=\"\${BOSH_ALL_PROXY}\""`),
			fmt.Sprintf(`echo "export CREDHUB_CLIENT=\"\${BOSH_CLIENT}\""`),
			fmt.Sprintf(`echo "export CREDHUB_SECRET=\"\${BOSH_CLIENT_SECRET}\""`),
			fmt.Sprintf(`echo "export CREDHUB_CA_CERT=\"\${BOSH_CA_CERT}\""`),
		)
	}

	return scripting.RunScript(lines, []string{"jq", "om", "ssh", "bosh"}, c.File)
}
