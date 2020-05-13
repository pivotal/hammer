/*
Copyright (C) 2019-Present Pivotal Software, Inc. All rights reserved.

This program and the accompanying materials are made available under the terms of the under the Apache License, Version 2.0 (the "License”); you may not use this file except in compliance with the License. You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.
*/

package om_test

import (
	"fmt"
	"net/url"

	"github.com/pivotal/hammer/om"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pivotal/hammer/environment"
	"github.com/pivotal/hammer/scripting/scriptingfakes"
)

var _ = Describe("om runner", func() {
	var (
		err          error
		omRunner     om.Runner
		scriptRunner *scriptingfakes.FakeScriptRunner

		data   environment.Config
		dryRun bool
		omArgs []string
	)

	BeforeEach(func() {
		scriptRunner = new(scriptingfakes.FakeScriptRunner)

		url, _ := url.Parse("https://www.test-url.io")
		data = environment.Config{
			OpsManager: environment.OpsManager{
				URL:          *url,
				Username:     "username",
				Password:     "password",
				ClientID:     "client_id",
				ClientSecret: "client_secret",
			},
		}

		omRunner = om.Runner{
			ScriptRunner: scriptRunner,
		}
	})

	JustBeforeEach(func() {
		err = omRunner.Run(data, dryRun, omArgs...)
	})

	When("no om args are passed to the runner", func() {
		BeforeEach(func() {
			omArgs = []string{}
		})

		It("runs the script with a series of om env var echos and no prerequisites", func() {
			Expect(scriptRunner.RunScriptCallCount()).To(Equal(1))

			lines, prereqs, _ := scriptRunner.RunScriptArgsForCall(0)
			Expect(lines).To(Equal([]string{
				`echo "export OM_TARGET='https://www.test-url.io'"`,
				`echo "export OM_CLIENT_ID='client_id'"`,
				`echo "export OM_CLIENT_SECRET='client_secret'"`,
				`echo "export OM_USERNAME='username'"`,
				`echo "export OM_PASSWORD='password'"`,
			}))
			Expect(prereqs).To(HaveLen(0))
		})

		It("doesn't error", func() {
			Expect(err).NotTo(HaveOccurred())
		})
	})

	When("one or more om args are passed to the runner", func() {
		BeforeEach(func() {
			omArgs = []string{"arg1", "arg2", "arg3"}
		})

		It("runs the script with an om command and specifying om as a prerequisite", func() {
			Expect(scriptRunner.RunScriptCallCount()).To(Equal(1))

			lines, prereqs, _ := scriptRunner.RunScriptArgsForCall(0)
			Expect(lines).To(Equal([]string{
				`OM_CLIENT_ID='client_id' OM_CLIENT_SECRET='client_secret' OM_USERNAME='username' OM_PASSWORD='password' om -t 'https://www.test-url.io' -k 'arg1' 'arg2' 'arg3'`,
			}))
			Expect(prereqs).To(ConsistOf("om"))
		})

		It("doesn't error", func() {
			Expect(err).NotTo(HaveOccurred())
		})
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

	When("running the script errors", func() {
		BeforeEach(func() {
			scriptRunner.RunScriptReturns(fmt.Errorf("run-script-error"))
		})

		It("propagates the error", func() {
			Expect(err).To(MatchError("run-script-error"))
		})
	})
})
