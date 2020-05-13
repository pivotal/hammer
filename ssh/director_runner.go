/*
Copyright (C) 2019-Present Pivotal Software, Inc. All rights reserved.

This program and the accompanying materials are made available under the terms of the under the Apache License, Version 2.0 (the "License”); you may not use this file except in compliance with the License. You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.
*/

package ssh

import (
	"fmt"

	"github.com/pivotal/hammer/environment"
	"github.com/pivotal/hammer/scripting"
)

type DirectorRunner struct {
	ScriptRunner scripting.ScriptRunner
}

func (b DirectorRunner) Run(data environment.Config, dryRun bool, args ...string) error {
	boshCreds := "/api/v0/deployed/director/credentials/bosh_commandline_credentials"
	bbrCredsPath := "/api/v0/deployed/director/credentials/bbr_ssh_credentials"
	privateKeyPath := ".credential.value.private_key_pem"

	sshCommandLines := []string{
		`ssh_key_path=$(mktemp temp.XXXXXX)`,
		fmt.Sprintf(`echo "%s" >"$ssh_key_path"`, data.OpsManager.PrivateKey),
		`chmod 0600 "${ssh_key_path}"`,

		fmt.Sprintf(`director_ssh_key="$(OM_CLIENT_ID='%s' OM_CLIENT_SECRET='%s' OM_USERNAME='%s' OM_PASSWORD='%s' om -t %s -k curl -s -p %s | jq -r %s)"`,
			data.OpsManager.ClientID,
			data.OpsManager.ClientSecret,
			data.OpsManager.Username,
			data.OpsManager.Password,
			data.OpsManager.URL.String(),
			bbrCredsPath,
			privateKeyPath),
		`director_ssh_key_path=$(mktemp)`,
		`echo -e "$director_ssh_key" > "$director_ssh_key_path"`,
		`chmod 0600 "${director_ssh_key_path}"`,

		fmt.Sprintf(`bosh_env="$(OM_CLIENT_ID='%s' OM_CLIENT_SECRET='%s' OM_USERNAME='%s' OM_PASSWORD='%s' om -t %s -k curl -s -p %s | grep -o "BOSH_ENVIRONMENT=\S*" | cut -f2 -d=)"`,
			data.OpsManager.ClientID,
			data.OpsManager.ClientSecret,
			data.OpsManager.Username,
			data.OpsManager.Password,
			data.OpsManager.URL.String(),
			boshCreds),

		`trap 'rm -f ${director_ssh_key_path}; rm -f ${ssh_key_path}' EXIT`,

		fmt.Sprintf(`jumpbox_cmd="ubuntu@%s -o IdentitiesOnly=yes -o StrictHostKeyChecking=no -i ${ssh_key_path}"`, data.OpsManager.IP.String()),
		`ssh -o IdentitiesOnly=yes -o StrictHostKeyChecking=no -J "$jumpbox_cmd" "bbr@${bosh_env}" -i "$director_ssh_key_path"`,
	}

	return b.ScriptRunner.RunScript(sshCommandLines, []string{"ssh", "om"}, dryRun)
}
