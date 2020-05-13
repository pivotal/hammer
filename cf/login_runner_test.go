/*
Copyright (C) 2019-Present Pivotal Software, Inc. All rights reserved.

This program and the accompanying materials are made available under the terms of the under the Apache License, Version 2.0 (the "License”); you may not use this file except in compliance with the License. You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.
*/

package cf_test

import (
	"fmt"
	"net/url"

	"github.com/pivotal/hammer/cf"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pivotal/hammer/environment"
	"github.com/pivotal/hammer/scripting/scriptingfakes"
)

var _ = Describe("cf login runner", func() {
	var (
		err           error
		cfLoginRunner cf.LoginRunner
		scriptRunner  *scriptingfakes.FakeScriptRunner

		data   environment.Config
		dryRun bool
	)

	BeforeEach(func() {
		scriptRunner = new(scriptingfakes.FakeScriptRunner)

		url, _ := url.Parse("https://www.test-url.io")
		data = environment.Config{
			CFDomain: "sys.test-url.io",
			OpsManager: environment.OpsManager{
				URL:          *url,
				Username:     "username",
				Password:     "password",
				ClientID:     "client_id",
				ClientSecret: "client_secret",
			},
		}

		cfLoginRunner = cf.LoginRunner{
			ScriptRunner: scriptRunner,
		}
	})

	JustBeforeEach(func() {
		err = cfLoginRunner.Run(data, dryRun)
	})

	It("runs the script with a cf login", func() {
		Expect(scriptRunner.RunScriptCallCount()).To(Equal(1))

		lines, _, _ := scriptRunner.RunScriptArgsForCall(0)

		Expect(lines).To(Equal([]string{
			`prods="$(OM_CLIENT_ID='client_id' OM_CLIENT_SECRET='client_secret' OM_USERNAME='username' OM_PASSWORD='password' om -t https://www.test-url.io -k curl -s -p /api/v0/staged/products)"`,
			`guid="$(echo "$prods" | jq -r '.[] | select(.type == "cf") | .guid')"`,
			`creds="$(OM_CLIENT_ID='client_id' OM_CLIENT_SECRET='client_secret' OM_USERNAME='username' OM_PASSWORD='password' om -t https://www.test-url.io -k curl -s -p /api/v0/deployed/products/"$guid"/credentials/.uaa.admin_credentials)"`, `user="$(echo "$creds" | jq -r .credential.value.identity)"`,
			`pass="$(echo "$creds" | jq -r .credential.value.password)"`,
			`cf login -a "api.sys.test-url.io" -u "$user" -p "$pass" --skip-ssl-validation`,
		}))
	})

	It("specifies the appropriate prerequisites when running the script", func() {
		Expect(scriptRunner.RunScriptCallCount()).To(Equal(1))

		_, prereqs, _ := scriptRunner.RunScriptArgsForCall(0)

		Expect(prereqs).To(ConsistOf("jq", "om", "cf"))
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
