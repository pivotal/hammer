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
			command := exec.Command(pathToPcf, "bosh", "-l", "fixtures/claim_manatee_response.json", "-f")
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
			command := exec.Command(pathToPcf, "bosh", "-l", "fixtures/claim_manatee_response.json", "-f", "deployments")
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
