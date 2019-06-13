package cmdtests_test

import (
	"io/ioutil"
	"os/exec"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gexec"
	"github.com/pivotal/pcf/internal/testhelpers"
)

var _ = Describe("OM", func() {
	When("run with parameters", func() {
		It("generates the correct script", func() {
			command := exec.Command(pathToPcf, "om", "-l", "fixtures/claim_manatee_response.json", "-f", "--", "foo")
			session, err := Start(command, GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())

			Eventually(session).Should(Exit(0))
			Eventually(string(session.Err.Contents())).Should(Equal(""))

			pathToFile := strings.TrimSuffix(string(session.Out.Contents()), "\n")
			contents, err := ioutil.ReadFile(pathToFile)
			Expect(err).NotTo(HaveOccurred())

			Expect(string(contents)).To(Equal(testhelpers.LoadFixture("om_script.sh")))
		})
	})

	When("run with parameters including JSON", func() {
		It("generates the correct script", func() {
			command := exec.Command(pathToPcf, "om", "-l", "fixtures/claim_manatee_response.json", "-f", "--", "configure-product",
				"--product-name", "p-rabbitmq", "--product-properties",
				`{".rabbitmq-server.server_admin_credentials":{"value":{"identity":"admin","password":"admin"}}}`)
			session, err := Start(command, GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())

			Eventually(session).Should(Exit(0))
			Eventually(string(session.Err.Contents())).Should(Equal(""))

			pathToFile := strings.TrimSuffix(string(session.Out.Contents()), "\n")
			contents, err := ioutil.ReadFile(pathToFile)
			Expect(err).NotTo(HaveOccurred())

			Expect(string(contents)).To(Equal(testhelpers.LoadFixture("om_script_json.sh")))
		})
	})

	When("run with no parameters", func() {
		It("generates the correct script", func() {
			command := exec.Command(pathToPcf, "om", "-l", "fixtures/claim_manatee_response.json")
			session, err := Start(command, GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())

			Eventually(session).Should(Exit(0))
			Eventually(string(session.Err.Contents())).Should(Equal(""))
			Eventually(string(session.Out.Contents())).Should(Equal("export OM_TARGET=https://pcf.manatee.cf-app.com\nexport OM_USERNAME=pivotalcf\nexport OM_PASSWORD=fakePassword\n"))
		})
	})
})
