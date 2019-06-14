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

var _ = Describe("SSH", func() {
	It("generates the correct scriptp", func() {
		command := exec.Command(pathToPcf, "ssh", "-t", "fixtures/claim_manatee_response.json", "-f")
		session, err := Start(command, GinkgoWriter, GinkgoWriter)
		Expect(err).NotTo(HaveOccurred())

		Eventually(session).Should(Exit(0))
		Eventually(string(session.Err.Contents())).Should(Equal(""))
		Eventually(session.Out).Should(Say("Connecting to: manatee"))

		output := strings.TrimSuffix(string(session.Out.Contents()), "\n")
		lines := strings.Split(output, "\n")
		pathToFile := lines[len(lines)-1]
		contents, err := ioutil.ReadFile(pathToFile)
		Expect(err).NotTo(HaveOccurred())

		Expect(string(contents)).To(Equal(LoadFixture("ssh_script.sh")))
	})
})
