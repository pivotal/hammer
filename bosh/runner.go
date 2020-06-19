/*
Copyright (C) 2019-Present Pivotal Software, Inc. All rights reserved.

This program and the accompanying materials are made available under the terms of the under the Apache License, Version 2.0 (the "License”); you may not use this file except in compliance with the License. You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.
*/

package bosh

import (
	"fmt"
	"strings"

	"github.com/pivotal/hammer/scripting"

	"github.com/pivotal/hammer/environment"
)

type Runner struct {
	ScriptRunner scripting.ScriptRunner
}

func (r Runner) Run(data environment.Config, dryRun bool, boshArgs ...string) error {
	lines := []string{
		`ssh_key_path=$(mktemp)`,
		fmt.Sprintf(`echo "%s" >"$ssh_key_path"`, data.OpsManager.PrivateKey),
		`chmod 0600 "${ssh_key_path}"`,

		`bosh_ca_path=$(mktemp)`,
		fmt.Sprintf(`ssh -o IdentitiesOnly=yes -o StrictHostKeyChecking=no -i "${ssh_key_path}" %s@"%s" cat /var/tempest/workspaces/default/root_ca_certificate 1>"${bosh_ca_path}" 2>/dev/null`, data.OpsManager.SshUser, data.OpsManager.IP.String()),
		`chmod 0600 "${bosh_ca_path}"`,

		fmt.Sprintf(`creds="$(OM_CLIENT_ID='%s' OM_CLIENT_SECRET='%s' OM_USERNAME='%s' OM_PASSWORD='%s' om -t %s -k curl -s -p /api/v0/deployed/director/credentials/bosh_commandline_credentials)"`,
			data.OpsManager.ClientID,
			data.OpsManager.ClientSecret,
			data.OpsManager.Username,
			data.OpsManager.Password,
			data.OpsManager.URL.String()),
		`bosh_all="$(echo "$creds" | jq -r .credential | tr ' ' '\n' | grep '=')"`,

		`bosh_client="$(echo "$bosh_all" | tr ' ' '\n' | grep 'BOSH_CLIENT=')"`,
		`bosh_env="$(echo "$bosh_all" | tr ' ' '\n' | grep 'BOSH_ENVIRONMENT=')"`,
		`bosh_secret="$(echo "$bosh_all" | tr ' ' '\n' | grep 'BOSH_CLIENT_SECRET=')"`,
		`bosh_ca_cert="BOSH_CA_CERT=$bosh_ca_path"`,
		fmt.Sprintf(`bosh_proxy="BOSH_ALL_PROXY=ssh+socks5://%s@%s:22?private-key=${ssh_key_path}"`, data.OpsManager.SshUser, data.OpsManager.IP.String()),
		fmt.Sprintf(`bosh_gw_host="BOSH_GW_HOST=%s"`, data.OpsManager.IP.String()),
		fmt.Sprintf(`bosh_gw_user="BOSH_GW_USER=%s"`, data.OpsManager.SshUser),
		`bosh_gw_private_key="BOSH_GW_PRIVATE_KEY=${ssh_key_path}"`,
	}

	prereqs := []string{"jq", "om", "ssh"}

	if len(boshArgs) > 0 {
		lines = append(
			lines,
			`trap 'rm -f ${ssh_key_path} ${bosh_ca_path}' EXIT`,
			fmt.Sprintf(`/usr/bin/env "$bosh_client" "$bosh_env" "$bosh_secret" "$bosh_ca_cert" "$bosh_proxy" "$bosh_gw_host" "$bosh_gw_user" "$bosh_gw_private_key" bosh %s`, strings.Join(boshArgs, " ")),
		)
		prereqs = append(prereqs, "bosh")
	} else {
		lines = append(
			lines,
			fmt.Sprintf(`echo "export BOSH_ENV_NAME=%s"`, data.Name),
			`echo "export $bosh_client"`,
			`echo "export $bosh_env"`,
			`echo "export $bosh_secret"`,
			`echo "export $bosh_ca_cert"`,
			`echo "export $bosh_proxy"`,
			`echo "export $bosh_gw_host"`,
			`echo "export $bosh_gw_user"`,
			`echo "export $bosh_gw_private_key"`,
			`echo "export CREDHUB_SERVER=\"\${BOSH_ENVIRONMENT}:8844\""`,
			`echo "export CREDHUB_PROXY=\"\${BOSH_ALL_PROXY}\""`,
			`echo "export CREDHUB_CLIENT=\"\${BOSH_CLIENT}\""`,
			`echo "export CREDHUB_SECRET=\"\${BOSH_CLIENT_SECRET}\""`,
			`echo "export CREDHUB_CA_CERT=\"\${BOSH_CA_CERT}\""`,
		)
	}

	return r.ScriptRunner.RunScript(lines, prereqs, dryRun)
}
