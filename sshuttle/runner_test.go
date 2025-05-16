/*
Copyright (C) 2019-Present Pivotal Software, Inc. All rights reserved.

This program and the accompanying materials are made available under the terms of the under the Apache License, Version 2.0 (the "License”); you may not use this file except in compliance with the License. You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.
*/

package sshuttle_test

import (
	"fmt"
	"net"
	"net/url"

	"github.com/pivotal/hammer/sshuttle"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/pivotal/hammer/environment"
	"github.com/pivotal/hammer/scripting/scriptingfakes"
)

var _ = Describe("sshuttle runner", func() {
	var (
		err            error
		sshuttleRunner sshuttle.Runner
		scriptRunner   *scriptingfakes.FakeScriptRunner

		data   environment.Config
		dryRun bool
	)

	BeforeEach(func() {
		scriptRunner = new(scriptingfakes.FakeScriptRunner)

		url, _ := url.Parse("https://www.test-url.io")
		data = environment.Config{
			OpsManager: environment.OpsManager{
				SshUser:      "ssh_user",
				PrivateKey:   "private-key-contents",
				IP:           net.ParseIP("10.0.0.6"),
				URL:          *url,
				Username:     "username",
				Password:     "password",
				ClientID:     "client_id",
				ClientSecret: "client_secret",
			},
		}

		sshuttleRunner = sshuttle.Runner{
			ScriptRunner: scriptRunner,
		}
	})

	JustBeforeEach(func() {
		err = sshuttleRunner.Run(data, dryRun)
	})

	It("runs the script with a sshuttle to the opsman vm", func() {
		Expect(scriptRunner.RunScriptCallCount()).To(Equal(1))

		lines, _, _ := scriptRunner.RunScriptArgsForCall(0)
		Expect(lines).To(Equal([]string{
			`ssh_key_path=$(mktemp)`,
			`echo "private-key-contents" >"$ssh_key_path"`,
			`trap 'rm -f ${ssh_key_path}' EXIT`,
			`chmod 0600 "${ssh_key_path}"`,
			`cidrs="$(OM_CLIENT_ID='client_id' OM_CLIENT_SECRET='client_secret' OM_USERNAME='username' OM_PASSWORD='password' om -t https://www.test-url.io -k curl -s -p /api/v0/staged/director/networks | jq -r .networks[].subnets[].cidr | xargs echo)"`,
			`sshuttle --ssh-cmd "ssh -o IdentitiesOnly=yes -i ${ssh_key_path}" -r ssh_user@"10.0.0.6" "${cidrs}"`,
		}))
	})

	It("specifies the appropriate prerequisites when running the script", func() {
		Expect(scriptRunner.RunScriptCallCount()).To(Equal(1))

		_, prereqs, _ := scriptRunner.RunScriptArgsForCall(0)

		Expect(prereqs).To(ConsistOf("jq", "om", "sshuttle"))
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
