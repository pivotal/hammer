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
				Say(`cf-login    log in to cf on the environment`),
				Say(`completion  command completion script`),
				Say(`om          run the 'om' command with credentials for this environment`),
				Say(`open        open a browser to this environment`),
				Say(`ssh         open an ssh connection to ops manager on this environment`),
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
				Say(`cf-login    log in to cf on the environment`),
				Say(`completion  command completion script`),
				Say(`om          run the 'om' command with credentials for this environment`),
				Say(`open        open a browser to this environment`),
				Say(`ssh         open an ssh connection to ops manager on this environment`),
				Say(`sshuttle    sshuttle to this environment`),
				Say(`version     version of command \(aliases: ver\)`),
			))
		})
	})
})
