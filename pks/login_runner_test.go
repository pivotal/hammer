/*
Copyright (C) 2019-Present Pivotal Software, Inc. All rights reserved.

This program and the accompanying materials are made available under the terms of the under the Apache License, Version 2.0 (the "License”); you may not use this file except in compliance with the License. You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.
*/

package pks_test

import (
	"fmt"
	"net/url"

	"github.com/pivotal/hammer/pks"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/pivotal/hammer/environment"
	"github.com/pivotal/hammer/scripting/scriptingfakes"
)

var _ = Describe("pks login runner", func() {
	var (
		err            error
		pksLoginRunner pks.LoginRunner
		scriptRunner   *scriptingfakes.FakeScriptRunner

		data   environment.Config
		dryRun bool
	)

	BeforeEach(func() {
		scriptRunner = new(scriptingfakes.FakeScriptRunner)

		opsmanUrl, _ := url.Parse("https://www.test-url.io")
		pksApiUrl, _ := url.Parse("api.test-url.io")
		data = environment.Config{
			OpsManager: environment.OpsManager{
				Username:     "username",
				Password:     "password",
				ClientID:     "client_id",
				ClientSecret: "client_secret",
				URL:          *opsmanUrl,
			},
			PKSApi: environment.PKSApi{
				URL:      *pksApiUrl,
				Username: "username",
				Password: "password",
			},
		}

		pksLoginRunner = pks.LoginRunner{
			ScriptRunner: scriptRunner,
		}
	})

	JustBeforeEach(func() {
		err = pksLoginRunner.Run(data, dryRun)
	})

	It("runs the script with a pks login", func() {
		Expect(scriptRunner.RunScriptCallCount()).To(Equal(1))

		lines, _, _ := scriptRunner.RunScriptArgsForCall(0)

		Expect(lines).To(Equal([]string{
			`prods="$(OM_CLIENT_ID='client_id' OM_CLIENT_SECRET='client_secret' OM_USERNAME='username' OM_PASSWORD='password' om -t https://www.test-url.io -k curl -s -p /api/v0/staged/products)"`,
			`guid="$(echo "$prods" | jq -r '.[] | select(.type == "pivotal-container-service") | .guid')"`,
			`creds="$(OM_CLIENT_ID='client_id' OM_CLIENT_SECRET='client_secret' OM_USERNAME='username' OM_PASSWORD='password' om -t https://www.test-url.io -k curl -s -p /api/v0/deployed/products/"$guid"/credentials/.properties.uaa_admin_password)"`,
			`pass="$(echo "$creds" | jq -r .credential.value.secret)"`,
			`pks login -a api.test-url.io -u admin -p "$pass" --skip-ssl-validation`,
		}))
	})

	It("specifies the appropriate prerequisites when running the script", func() {
		Expect(scriptRunner.RunScriptCallCount()).To(Equal(1))

		_, prereqs, _ := scriptRunner.RunScriptArgsForCall(0)

		Expect(prereqs).To(ConsistOf("jq", "om", "pks"))
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
		BeforeEach(func() {
			scriptRunner.RunScriptReturns(nil)
		})

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
