/*
Copyright (C) 2019-Present Pivotal Software, Inc. All rights reserved.

This program and the accompanying materials are made available under the terms of the under the Apache License, Version 2.0 (the "License”); you may not use this file except in compliance with the License. You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.
*/

package integration

import (
	"os/exec"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
	. "github.com/onsi/gomega/gexec"
)

var _ = Describe("Readers", func() {
	readers := []TableEntry{
		Entry("bosh", "bosh"),
		Entry("cf-login", "cf-login"),
		Entry("open", "open"),
		Entry("om", "om"),
		Entry("ssh opsman", "ssh", "opsman"),
		Entry("ssh director", "ssh", "director"),
		Entry("sshuttle", "sshuttle"),
	}

	DescribeTable("failure when the environment path is not found",
		func(subcmds ...string) {
			params := append(subcmds, "-t", "/this/should/not/exist")
			session := runPcf([]string{}, params...)

			Eventually(session).Should(Exit(1))
			Eventually(session.Err).Should(Say("open /this/should/not/exist: no such file or directory"))
			Eventually(string(session.Out.Contents())).Should(Equal(""))
		},
		readers,
	)

	DescribeTable("accepting the `-t` flag before the subcommands",
		func(subcmds ...string) {
			params := append([]string{"-t", "/also/should/not/exist"}, subcmds...)
			session := runPcf([]string{}, params...)

			Eventually(session).Should(Exit(1))
			Eventually(session.Err).Should(Say("open /also/should/not/exist: no such file or directory"))
			Eventually(string(session.Out.Contents())).Should(Equal(""))
		},
		readers,
	)

	DescribeTable("reading the environment from $HAMMER_TARGET_CONFIG",
		func(subcmds ...string) {
			env := []string{"HAMMER_TARGET_CONFIG=fixtures/claim_manatee_response.json"}
			params := append(subcmds, "-f")
			session := runPcf(env, params...)

			Eventually(session).Should(Exit(0))
			Eventually(string(session.Err.Contents())).Should(Equal(""))
		},
		readers,
	)

	DescribeTable("failure to specify the `-t` flags",
		func(subcmds ...string) {
			session := runPcf([]string{}, subcmds...)

			Eventually(session).Should(Exit(1))
			Eventually(string(session.Err.Contents())).Should(Equal("You must specify the target environment config path (--target | -t) flag\n"))
			Eventually(string(session.Out.Contents())).Should(Equal(""))
		},
		readers,
	)

	DescribeTable("specifying a mismatching environment name with a single environments",
		func(subcmds ...string) {
			env := []string{
				"HAMMER_TARGET_CONFIG=fixtures/claim_manatee_response.json",
				"HAMMER_ENVIRONMENT_NAME=environment-that-does-not-match",
			}
			session := runPcf(env, subcmds...)

			Eventually(session).Should(Exit(1))
			Eventually(string(session.Err.Contents())).Should(Equal("Environment name 'environment-that-does-not-match' specified but does not match environment in config\n"))
			Eventually(string(session.Out.Contents())).Should(Equal(""))
		},
		readers,
	)

	DescribeTable("specifying a non-existent environment name with a config containing an environments list",
		func(subcmds ...string) {
			env := []string{
				"HAMMER_TARGET_CONFIG=fixtures/multiple_environment_config.json",
				"HAMMER_ENVIRONMENT_NAME=environment-that-does-not-exist",
			}
			session := runPcf(env, subcmds...)

			Eventually(session).Should(Exit(1))
			Eventually(string(session.Err.Contents())).Should(Equal("Environment name 'environment-that-does-not-exist' specified but does not match environment in config\n"))
			Eventually(string(session.Out.Contents())).Should(Equal(""))
		},
		readers,
	)

	DescribeTable("not specifying environment name with a config containing an environments list",
		func(subcmds ...string) {
			env := []string{
				"HAMMER_TARGET_CONFIG=fixtures/multiple_environment_config.json",
			}
			params := append(subcmds, "-f")
			session := runPcf(env, params...)

			Eventually(session).Should(Exit(0))
			Eventually(string(session.Err.Contents())).Should(Equal(""))
		},
		readers,
	)

	DescribeTable("specifying environment name via an env var",
		func(subcmds ...string) {
			env := []string{
				"HAMMER_TARGET_CONFIG=fixtures/multiple_environment_config.json",
				"HAMMER_ENVIRONMENT_NAME=narwhal",
			}
			params := append(subcmds, "-f")
			session := runPcf(env, params...)

			Eventually(session).Should(Exit(0))
			Eventually(string(session.Err.Contents())).Should(Equal(""))
		},
		readers,
	)

	DescribeTable("specifying environment name via a flag",
		func(subcmds ...string) {
			params := append([]string{"-t", "fixtures/multiple_environment_config.json", "-e", "narwhal"}, subcmds...)
			params = append(params, "-f")
			session := runPcf([]string{}, params...)

			Eventually(session).Should(Exit(0))
			Eventually(string(session.Err.Contents())).Should(Equal(""))
		},
		readers,
	)
})

func runPcf(env []string, params ...string) *Session {
	command := exec.Command(pathToPcf, params...)
	command.Env = env
	session, err := Start(command, GinkgoWriter, GinkgoWriter)
	Expect(err).NotTo(HaveOccurred())
	return session
}
