/*
Copyright (C) 2019-Present Pivotal Software, Inc. All rights reserved.

This program and the accompanying materials are made available under the terms of the under the Apache License, Version 2.0 (the "License”); you may not use this file except in compliance with the License. You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.
*/

package integration

import (
	"io/ioutil"
	"os/exec"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
	. "github.com/onsi/gomega/gexec"
)

var _ = Describe("Open", func() {
	It("generates the correct script", func() {
		command := exec.Command(pathToPcf, "open", "-t", "fixtures/claim_manatee_response.json", "-f")
		session, err := Start(command, GinkgoWriter, GinkgoWriter)
		Expect(err).NotTo(HaveOccurred())

		Eventually(session).Should(Exit(0))
		Eventually(string(session.Err.Contents())).Should(Equal(""))
		Eventually(session.Out).Should(Say("Username is: pivotalcf"))
		Eventually(session.Out).Should(Say("Password is in the clipboard"))

		output := strings.TrimSuffix(string(session.Out.Contents()), "\n")
		pathToFile := LastLine(output)
		contents, err := ioutil.ReadFile(pathToFile)
		Expect(err).NotTo(HaveOccurred())

		Expect(string(contents)).To(Equal(LoadFixture("open_script.sh")))
	})

	When("run with the --show flag", func() {
		It("prints the credentials to the screen", func() {
			command := exec.Command(pathToPcf, "open", "-t", "fixtures/claim_manatee_response.json", "--show")
			session, err := Start(command, GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())

			Eventually(session).Should(Exit(0))
			Eventually(string(session.Err.Contents())).Should(Equal(""))
			Eventually(session.Out).Should(Say("https://pcf.manatee.cf-app.com"))
			Eventually(session.Out).Should(Say("username: pivotalcf"))
			Eventually(session.Out).Should(Say("password: fakePassword"))
		})
	})
})
