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
	. "github.com/onsi/gomega/gexec"
)

var _ = Describe("BOSH", func() {
	When("getting the BOSH credentials", func() {
		It("generates the correct script", func() {
			command := exec.Command(pathToPcf, "bosh", "-t", "fixtures/claim_manatee_response.json", "-f")
			session, err := Start(command, GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())

			Eventually(session).Should(Exit(0))
			Eventually(string(session.Err.Contents())).Should(Equal(""))

			output := strings.TrimSuffix(string(session.Out.Contents()), "\n")
			contents, err := ioutil.ReadFile(output)
			Expect(err).NotTo(HaveOccurred())

			Expect(string(contents)).To(Equal(LoadFixture("bosh_creds_script.sh")))
		})
	})

	When("running a BOSH command", func() {
		It("generates the correct script", func() {
			command := exec.Command(pathToPcf, "bosh", "-t", "fixtures/claim_manatee_response.json", "-f", "deployments")
			session, err := Start(command, GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())

			Eventually(session).Should(Exit(0))
			Eventually(string(session.Err.Contents())).Should(Equal(""))

			output := strings.TrimSuffix(string(session.Out.Contents()), "\n")
			contents, err := ioutil.ReadFile(output)
			Expect(err).NotTo(HaveOccurred())

			Expect(string(contents)).To(Equal(LoadFixture("bosh_cmd_script.sh")))
		})
	})
})
