package ssh

import (
	"fmt"

	"github.com/pivotal/pcf/environment"
	"github.com/pivotal/pcf/scripting"
)

type Runner struct {
	ScriptRunner scripting.ScriptRunner
}

func (b Runner) Run(data environment.Config, dryRun bool, args ...string) error {
	sshCommand := fmt.Sprintf(`ssh -o IdentitiesOnly=yes -o StrictHostKeyChecking=no -i "${ssh_key_path}" -t ubuntu@"%s"`, data.OpsManager.IP.String())

	sshCommandLines := []string{
		fmt.Sprintf(`ssh_key_path=$(mktemp)`),
		fmt.Sprintf(`echo "%s" >"$ssh_key_path"`, data.OpsManager.PrivateKey),
		fmt.Sprintf(`trap 'rm -f ${ssh_key_path}' EXIT`),
		fmt.Sprintf(`chmod 0600 "${ssh_key_path}"`),
		fmt.Sprintf(`creds="$(om -t %s -k -u %s -p %s curl -s -p %s)"`, data.OpsManager.URL.String(), data.OpsManager.Username, data.OpsManager.Password, "/api/v0/deployed/director/credentials/bosh_commandline_credentials"),
		fmt.Sprintf(`bosh="$(echo "$creds" | jq -r .credential | tr ' ' '\n' | grep '=')"`),
		fmt.Sprintf(`echo "$bosh"`),
		fmt.Sprintf(`shell="/usr/bin/env $(echo $bosh | tr '\n' ' ') bash -l"`),
		fmt.Sprintf(`%s "$shell"`, sshCommand),
	}

	return b.ScriptRunner.RunScript(sshCommandLines, []string{"ssh", "om"}, dryRun)
}
