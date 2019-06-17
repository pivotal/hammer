package ssh_test

import (
	"fmt"
	"net"
	"net/url"

	"github.com/pivotal/pcf/ssh"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pivotal/pcf/environment"
	"github.com/pivotal/pcf/scripting/scriptingfakes"
)

var _ = Describe("ssh runner", func() {

	var err error
	var sshRunner ssh.Runner
	var scriptRunner *scriptingfakes.FakeScriptRunner

	var data environment.Config
	var dryRun bool

	BeforeEach(func() {
		scriptRunner = new(scriptingfakes.FakeScriptRunner)

		url, _ := url.Parse("www.test-url.io")
		data = environment.Config{
			OpsManager: environment.OpsManager{
				PrivateKey: "private-key-contents",
				IP:         net.ParseIP("10.0.0.6"),
				URL:        *url,
				Username:   "username",
				Password:   "password",
			},
		}
		dryRun = true

		sshRunner = ssh.Runner{
			ScriptRunner: scriptRunner,
		}
	})

	JustBeforeEach(func() {
		err = sshRunner.Run(data, dryRun)
	})

	Context("run", func() {
		It("invokes script runner with an ssh to the opsman vm", func() {
			Expect(scriptRunner.RunScriptCallCount()).To(Equal(1))

			lines, prereqs, dryRun := scriptRunner.RunScriptArgsForCall(0)
			Expect(lines).To(ContainElement(`echo "private-key-contents" >"$ssh_key_path"`))
			Expect(lines).To(ContainElement(`creds="$(om -t www.test-url.io -k -u username -p password curl -s -p /api/v0/deployed/director/credentials/bosh_commandline_credentials)"`))
			Expect(lines).To(ContainElement(`ssh -o IdentitiesOnly=yes -o StrictHostKeyChecking=no -i "${ssh_key_path}" -t ubuntu@"10.0.0.6" "$shell"`))

			Expect(prereqs).To(ConsistOf("ssh", "om"))
			Expect(dryRun).To(Equal(true))
		})

		It("doesn't error", func() {
			Expect(err).NotTo(HaveOccurred())
		})

		When("script runner run script errors", func() {
			BeforeEach(func() {
				scriptRunner.RunScriptReturns(fmt.Errorf("run-script-error"))
			})

			It("propagates the error", func() {
				Expect(err).To(MatchError("run-script-error"))
			})
		})
	})
})
