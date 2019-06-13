package cmdtests_test

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

	DescribeTable("failure when the lockfile path is not found",
		func(subcmd string) {
			session := runPcf([]string{}, subcmd, "-l", "/this/should/not/exist")

			Eventually(session).Should(Exit(1))
			Eventually(session.Err).Should(Say("open /this/should/not/exist: no such file or directory"))
			Eventually(string(session.Out.Contents())).Should(Equal(""))
		},
		readers...,
	)

	DescribeTable("accepting the `-l` flag before the subcommand",
		func(subcmd string) {
			session := runPcf([]string{}, "-l", "/also/should/not/exist", subcmd)

			Eventually(session).Should(Exit(1))
			Eventually(session.Err).Should(Say("open /also/should/not/exist: no such file or directory"))
			Eventually(string(session.Out.Contents())).Should(Equal(""))
		},
		readers...,
	)

	DescribeTable("reading the lockfile from $ENVIRONMENT_LOCK_METADATA",
		func(subcmd string) {
			env := []string{"ENVIRONMENT_LOCK_METADATA=fixtures/claim_manatee_response.json"}
			session := runPcf(env, subcmd, "-f")

			Eventually(session).Should(Exit(0))
			Eventually(string(session.Err.Contents())).Should(Equal(""))
		},
		readers...,
	)

	DescribeTable("failure to specify the `-l` flags",
		func(subcmd string) {
			session := runPcf([]string{}, subcmd)

			Eventually(session).Should(Exit(1))
			Eventually(string(session.Err.Contents())).Should(Equal("You must specify the lockfile path (--lockfile | -l) flag\n"))
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
