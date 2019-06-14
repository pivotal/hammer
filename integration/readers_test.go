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
		Entry("ssh", "ssh"),
		Entry("sshuttle", "sshuttle"),
	}

	DescribeTable("failure when the environment path is not found",
		func(subcmd string) {
			session := runPcf([]string{}, subcmd, "-t", "/this/should/not/exist")

			Eventually(session).Should(Exit(1))
			Eventually(session.Err).Should(Say("open /this/should/not/exist: no such file or directory"))
			Eventually(string(session.Out.Contents())).Should(Equal(""))
		},
		readers...,
	)

	DescribeTable("accepting the `-t` flag before the subcommand",
		func(subcmd string) {
			session := runPcf([]string{}, "-t", "/also/should/not/exist", subcmd)

			Eventually(session).Should(Exit(1))
			Eventually(session.Err).Should(Say("open /also/should/not/exist: no such file or directory"))
			Eventually(string(session.Out.Contents())).Should(Equal(""))
		},
		readers...,
	)

	DescribeTable("reading the environment from $TARGET_ENVIRONMENT_CONFIG",
		func(subcmd string) {
			env := []string{"TARGET_ENVIRONMENT_CONFIG=fixtures/claim_manatee_response.json"}
			session := runPcf(env, subcmd, "-f")

			Eventually(session).Should(Exit(0))
			Eventually(string(session.Err.Contents())).Should(Equal(""))
		},
		readers...,
	)

	DescribeTable("failure to specify the `-t` flags",
		func(subcmd string) {
			session := runPcf([]string{}, subcmd)

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
