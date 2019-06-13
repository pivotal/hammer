package integration

import (
	"io/ioutil"
	"os/exec"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gexec"
)

var _ = Describe("sshuttle", func() {
	When("running shuttle", func() {
		It("generates the correct script", func() {
			command := exec.Command(pathToPcf, "sshuttle", "-l", "fixtures/claim_manatee_response.json", "-f")
			session, err := Start(command, GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())

			Eventually(session).Should(Exit(0))
			Eventually(string(session.Err.Contents())).Should(Equal(""))

			output := strings.TrimSuffix(string(session.Out.Contents()), "\n")
			contents, err := ioutil.ReadFile(output)
			Expect(err).NotTo(HaveOccurred())

			Expect(string(contents)).To(Equal(LoadFixture("sshuttle_script.sh")))
		})
	})
})
