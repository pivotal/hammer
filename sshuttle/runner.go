package sshuttle

import (
	"fmt"

	"github.com/pivotal/pcf/lockfile"
	"github.com/pivotal/pcf/scripting"
)

type Runner struct {
	ScriptRunner scripting.ScriptRunner
}

func (b Runner) Run(data lockfile.Lockfile, dryRun bool, args ...string) error {
	sshuttleCommandLines := []string{
		fmt.Sprintf(`ssh_key_path=$(mktemp)`),
		fmt.Sprintf(`echo "%s" >"$ssh_key_path"`, data.OpsManager.PrivateKey),
		fmt.Sprintf(`trap 'rm -f ${ssh_key_path}' EXIT`),
		fmt.Sprintf(`chmod 0600 "${ssh_key_path}"`),
		fmt.Sprintf(`sshuttle --ssh-cmd "ssh -o IdentitiesOnly=yes -i ${ssh_key_path}" -r ubuntu@"%s" %s %s %s`, data.OpsManager.IP.String(), data.OpsManager.CIDR.String(), data.PasCIDR.String(), data.ServicesCIDR.String()),
	}

	return b.ScriptRunner.RunScript(sshuttleCommandLines, []string{"jq", "om", "sshuttle"}, dryRun)
}
