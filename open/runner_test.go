/*
Copyright (C) 2019-Present Pivotal Software, Inc. All rights reserved.

This program and the accompanying materials are made available under the terms of the under the Apache License, Version 2.0 (the "License”); you may not use this file except in compliance with the License. You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.
*/

package open_test

import (
	"fmt"
	"net/url"

	"github.com/pivotal/hammer/open"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/pivotal/hammer/environment"
	"github.com/pivotal/hammer/scripting/scriptingfakes"
)

var _ = Describe("open runner", func() {
	var (
		err          error
		openRunner   open.Runner
		scriptRunner *scriptingfakes.FakeScriptRunner

		data   environment.Config
		dryRun bool
	)

	BeforeEach(func() {
		scriptRunner = new(scriptingfakes.FakeScriptRunner)

		url, _ := url.Parse("https://www.test-url.io")
		data = environment.Config{
			OpsManager: environment.OpsManager{
				URL:      *url,
				Password: "password",
			},
		}

		openRunner = open.Runner{
			ScriptRunner: scriptRunner,
		}
	})

	JustBeforeEach(func() {
		err = openRunner.Run(data, dryRun)
	})

	It("runs the script with opsman url open and copying of the password into the clipboard", func() {
		Expect(scriptRunner.RunScriptCallCount()).To(Equal(1))

		lines, _, _ := scriptRunner.RunScriptArgsForCall(0)
		Expect(lines).To(Equal([]string{
			`open "https://www.test-url.io"`,
			`echo "password" | pbcopy`,
		}))
	})

	When("client credentials are specified", func() {
		BeforeEach(func() {
			url, _ := url.Parse("https://www.test-url.io")
			data = environment.Config{
				OpsManager: environment.OpsManager{
					URL:          *url,
					Password:     "password",
					ClientID:     "client_id",
					ClientSecret: "client_secret",
				},
			}
		})

		It("runs the script with opsman url open and copying of the client secret into the clipboard", func() {
			Expect(scriptRunner.RunScriptCallCount()).To(Equal(1))

			lines, _, _ := scriptRunner.RunScriptArgsForCall(0)
			Expect(lines).To(Equal([]string{
				`open "https://www.test-url.io"`,
				`echo "client_secret" | pbcopy`,
			}))
		})
	})

	It("specifies the appropriate prerequisites when running the script", func() {
		Expect(scriptRunner.RunScriptCallCount()).To(Equal(1))

		_, prereqs, _ := scriptRunner.RunScriptArgsForCall(0)

		Expect(prereqs).To(ConsistOf("open", "pbcopy"))
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
