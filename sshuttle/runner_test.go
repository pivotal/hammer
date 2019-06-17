package sshuttle_test

import (
	"fmt"
	"net"

	"github.com/pivotal/pcf/sshuttle"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pivotal/pcf/environment"
	"github.com/pivotal/pcf/scripting/scriptingfakes"
)

var _ = Describe("sshuttle runner", func() {

	var err error
	var sshuttleRunner sshuttle.Runner
	var scriptRunner *scriptingfakes.FakeScriptRunner

	var data environment.Config
	var dryRun bool

	BeforeEach(func() {
		scriptRunner = new(scriptingfakes.FakeScriptRunner)

		_, omCIDR, _ := net.ParseCIDR("10.0.0.0/24")
		_, pasCIDR, _ := net.ParseCIDR("10.0.4.0/24")
		_, servicesCIDR, _ := net.ParseCIDR("10.0.8.0/24")

		data = environment.Config{
			OpsManager: environment.OpsManager{
				PrivateKey: "private-key-contents",
				IP:         net.ParseIP("10.0.0.6"),
				CIDR:       *omCIDR,
			},
			PasCIDR:      *pasCIDR,
			ServicesCIDR: *servicesCIDR,
		}
		dryRun = false

		sshuttleRunner = sshuttle.Runner{
			ScriptRunner: scriptRunner,
		}
	})

	JustBeforeEach(func() {
		err = sshuttleRunner.Run(data, dryRun)
	})

	Context("run", func() {
		It("invokes script runner with a sshuttle to the opsman vm", func() {
			Expect(scriptRunner.RunScriptCallCount()).To(Equal(1))

			lines, prereqs, dryRun := scriptRunner.RunScriptArgsForCall(0)
			Expect(lines).To(ContainElement(`echo "private-key-contents" >"$ssh_key_path"`))
			Expect(lines).To(ContainElement(`sshuttle --ssh-cmd "ssh -o IdentitiesOnly=yes -i ${ssh_key_path}" -r ubuntu@"10.0.0.6" 10.0.0.0/24 10.0.4.0/24 10.0.8.0/24`))

			Expect(prereqs).To(ConsistOf("jq", "om", "sshuttle"))
			Expect(dryRun).To(Equal(false))
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
