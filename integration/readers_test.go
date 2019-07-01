package integration

import (
	"os/exec"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
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
		readers...,
	)

	DescribeTable("accepting the `-t` flag before the subcommands",
		func(subcmds ...string) {
			params := append([]string{"-t", "/also/should/not/exist"}, subcmds...)
			session := runPcf([]string{}, params...)

			Eventually(session).Should(Exit(1))
			Eventually(session.Err).Should(Say("open /also/should/not/exist: no such file or directory"))
			Eventually(string(session.Out.Contents())).Should(Equal(""))
		},
		readers...,
	)

	DescribeTable("reading the environment from $TARGET_ENVIRONMENT_CONFIG",
		func(subcmds ...string) {
			env := []string{"TARGET_ENVIRONMENT_CONFIG=fixtures/claim_manatee_response.json"}
			params := append(subcmds, "-f")
			session := runPcf(env, params...)

			Eventually(session).Should(Exit(0))
			Eventually(string(session.Err.Contents())).Should(Equal(""))
		},
		readers...,
	)

	DescribeTable("failure to specify the `-t` flags",
		func(subcmds ...string) {
			session := runPcf([]string{}, subcmds...)

			Eventually(session).Should(Exit(1))
			Eventually(string(session.Err.Contents())).Should(Equal("You must specify the target environment config path (--target | -t) flag\n"))
			Eventually(string(session.Out.Contents())).Should(Equal(""))
		},
		readers...,
	)
})

func runPcf(env []string, params ...string) *Session {
	command := exec.Command(pathToPcf, params...)
	command.Env = env
	session, err := Start(command, GinkgoWriter, GinkgoWriter)
	Expect(err).NotTo(HaveOccurred())
	return session
}
