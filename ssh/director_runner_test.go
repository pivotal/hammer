/*
Copyright (C) 2019-Present Pivotal Software, Inc. All rights reserved.

This program and the accompanying materials are made available under the terms of the under the Apache License, Version 2.0 (the "License”); you may not use this file except in compliance with the License. You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.
*/

package ssh_test

import (
	"fmt"
	"net"
	"net/url"

	"github.com/pivotal/hammer/ssh"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pivotal/hammer/environment"
	"github.com/pivotal/hammer/scripting/scriptingfakes"
)

var _ = Describe("director ssh runner", func() {
	var (
		err          error
		sshRunner    ssh.DirectorRunner
		scriptRunner *scriptingfakes.FakeScriptRunner

		data   environment.Config
		dryRun bool
	)

	BeforeEach(func() {
		scriptRunner = new(scriptingfakes.FakeScriptRunner)

		url, _ := url.Parse("https://www.test-url.io")
		data = environment.Config{
			OpsManager: environment.OpsManager{
				PrivateKey: "private-key-contents",
				IP:         net.ParseIP("10.0.0.6"),
				URL:        *url,
				Username:   "username",
				Password:   "password",
			},
		}

		sshRunner = ssh.DirectorRunner{
			ScriptRunner: scriptRunner,
		}
	})

	JustBeforeEach(func() {
		err = sshRunner.Run(data, dryRun)
	})

	It("runs the script with an ssh to the director vm", func() {
		Expect(scriptRunner.RunScriptCallCount()).To(Equal(1))

		lines, _, _ := scriptRunner.RunScriptArgsForCall(0)
		Expect(lines).To(Equal([]string{
			`ssh_key_path=$(mktemp temp.XXXXXX)`,
			`echo "private-key-contents" >"$ssh_key_path"`,
			`chmod 0600 "${ssh_key_path}"`,
			`director_ssh_key="$(om -t https://www.test-url.io -k -u username -p password curl -s -p /api/v0/deployed/director/credentials/bbr_ssh_credentials | jq -r .credential.value.private_key_pem)"`,
			`director_ssh_key_path=$(mktemp)`,
			`echo -e "$director_ssh_key" > "$director_ssh_key_path"`,
			`chmod 0600 "${director_ssh_key_path}"`,
			`bosh_env="$(om -t https://www.test-url.io -k -u username -p password curl -s -p /api/v0/deployed/director/credentials/bosh_commandline_credentials | grep -o "BOSH_ENVIRONMENT=\S*" | cut -f2 -d=)"`,
			`trap 'rm -f ${director_ssh_key_path}; rm -f ${ssh_key_path}' EXIT`,
			`jumpbox_cmd="ubuntu@10.0.0.6 -o IdentitiesOnly=yes -o StrictHostKeyChecking=no -i ${ssh_key_path}"`,
			`ssh -o IdentitiesOnly=yes -o StrictHostKeyChecking=no -J "$jumpbox_cmd" "bbr@${bosh_env}" -i "$director_ssh_key_path"`,
		}))
	})

	It("specifies the appropriate prerequisites when running the script", func() {
		Expect(scriptRunner.RunScriptCallCount()).To(Equal(1))

		_, prereqs, _ := scriptRunner.RunScriptArgsForCall(0)

		Expect(prereqs).To(ConsistOf("ssh", "om"))
	})

	When("run with dry run set to false", func() {
		BeforeEach(func() {
			dryRun = false
		})

		It("runs the script in dry run mode", func() {
			Expect(scriptRunner.RunScriptCallCount()).To(Equal(1))

			_, _, dryRun := scriptRunner.RunScriptArgsForCall(0)
			Expect(dryRun).To(Equal(false))
		})
	})

	When("run with dry run set to true", func() {
		BeforeEach(func() {
			dryRun = true
		})

		It("runs the script in dry run mode", func() {
			Expect(scriptRunner.RunScriptCallCount()).To(Equal(1))

			_, _, dryRun := scriptRunner.RunScriptArgsForCall(0)
			Expect(dryRun).To(Equal(true))
		})
	})

	When("running the script succeeds", func() {
		It("doesn't error", func() {
			Expect(err).NotTo(HaveOccurred())
		})
	})

	When("running the script errors", func() {
		BeforeEach(func() {
			scriptRunner.RunScriptReturns(fmt.Errorf("run-script-error"))
		})

		It("propagates the error", func() {
			Expect(err).To(MatchError("run-script-error"))
		})
	})
})
