package integration

import (
	"io/ioutil"
	"os/exec"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
	. "github.com/onsi/gomega/gexec"
	"github.com/pivotal/pcf/internal/testhelpers"
)

var _ = Describe("Open", func() {
	It("generates the correct script", func() {
		command := exec.Command(pathToPcf, "open", "-l", "fixtures/claim_manatee_response.json", "-f")
		session, err := Start(command, GinkgoWriter, GinkgoWriter)
		Expect(err).NotTo(HaveOccurred())

		Eventually(session).Should(Exit(0))
		Eventually(string(session.Err.Contents())).Should(Equal(""))
		Eventually(session.Out).Should(Say("Username is: pivotalcf"))
		Eventually(session.Out).Should(Say("Password is in the clipboard"))

		output := strings.TrimSuffix(string(session.Out.Contents()), "\n")
		lines := strings.Split(output, "\n")
		pathToFile := lines[len(lines)-1]
		contents, err := ioutil.ReadFile(pathToFile)
		Expect(err).NotTo(HaveOccurred())

		Expect(string(contents)).To(Equal(testhelpers.LoadFixture("open_script.sh")))
	})

	When("run with the --show flag", func() {
		It("prints the credentials to the screen", func() {
			command := exec.Command(pathToPcf, "open", "-l", "fixtures/claim_manatee_response.json", "--show")
			session, err := Start(command, GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())

			Eventually(session).Should(Exit(0))
			Eventually(string(session.Err.Contents())).Should(Equal(""))
			Eventually(session.Out).Should(Say("https://pcf.manatee.cf-app.com"))
			Eventually(session.Out).Should(Say("username: pivotalcf"))
			Eventually(session.Out).Should(Say("password: fakePassword"))
		})
	})
})
