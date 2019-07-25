/*
Copyright (C) 2019-Present Pivotal Software, Inc. All rights reserved.

This program and the accompanying materials are made available under the terms of the under the Apache License, Version 2.0 (the "License”); you may not use this file except in compliance with the License. You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.
*/

package integration

import (
	"os/exec"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
	. "github.com/onsi/gomega/gexec"
)

var _ = Describe("CLI", func() {
	When("no command or flag is passed", func() {
		It("displays help and exits zero", func() {
			command := exec.Command(pathToPcf)
			session, err := Start(command, GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())

			Eventually(session).Should(Exit(0))

			Eventually(session).Should(SatisfyAll(
				Say(`Usage:`),
				Say(`hammer \[OPTIONS\] <command>`),
				Say(`Application Options:`),
				Say(`-t, --target= path to the target environment config`),
				Say(`Help Options:`),
				Say(`-h, --help    Show this help message`),
				Say(`Available commands:`),
				Say(`bosh        display BOSH credentials, or run a BOSH command`),
				Say(`cf-login    log in to the cf for the environment`),
				Say(`completion  command completion script`),
				Say(`om          run the 'om' command with credentials for this environment`),
				Say(`open        open a browser to this environment`),
				Say(`ssh         open an ssh connection to the ops manager or director of this environment`),
				Say(`sshuttle    sshuttle to this environment`),
				Say(`version     version of command \(aliases: ver\)`),
			))
		})
	})

	When("--help is passed", func() {
		It("displays help and exits zero", func() {
			command := exec.Command(pathToPcf, "--help")
			session, err := Start(command, GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())

			Eventually(session).Should(Exit(0))

			Eventually(session).Should(SatisfyAll(
				Say(`Usage:`),
				Say(`hammer \[OPTIONS\] <command>`),
				Say(`Application Options:`),
				Say(`-t, --target= path to the target environment config`),
				Say(`Help Options:`),
				Say(`-h, --help    Show this help message`),
				Say(`Available commands:`),
				Say(`bosh        display BOSH credentials, or run a BOSH command`),
				Say(`cf-login    log in to the cf for the environment`),
				Say(`completion  command completion script`),
				Say(`om          run the 'om' command with credentials for this environment`),
				Say(`open        open a browser to this environment`),
				Say(`ssh         open an ssh connection to the ops manager or director of this environment`),
				Say(`sshuttle    sshuttle to this environment`),
				Say(`version     version of command \(aliases: ver\)`),
			))
		})
	})
})
