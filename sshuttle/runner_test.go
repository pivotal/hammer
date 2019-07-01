package sshuttle_test

import (
	"fmt"
	"net"

	"github.com/pivotal/hammer/sshuttle"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pivotal/hammer/environment"
	"github.com/pivotal/hammer/scripting/scriptingfakes"
)

var _ = Describe("sshuttle runner", func() {
	var (
		err            error
		sshuttleRunner sshuttle.Runner
		scriptRunner   *scriptingfakes.FakeScriptRunner

		data   environment.Config
		dryRun bool
	)

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

		sshuttleRunner = sshuttle.Runner{
			ScriptRunner: scriptRunner,
		}
	})

	JustBeforeEach(func() {
		err = sshuttleRunner.Run(data, dryRun)
	})

	It("runs the script with a sshuttle to the opsman vm", func() {
		Expect(scriptRunner.RunScriptCallCount()).To(Equal(1))

		lines, _, _ := scriptRunner.RunScriptArgsForCall(0)
		Expect(lines).To(Equal([]string{
			`ssh_key_path=$(mktemp)`,
			`echo "private-key-contents" >"$ssh_key_path"`,
			`trap 'rm -f ${ssh_key_path}' EXIT`,
			`chmod 0600 "${ssh_key_path}"`,
			`sshuttle --ssh-cmd "ssh -o IdentitiesOnly=yes -i ${ssh_key_path}" -r ubuntu@"10.0.0.6" 10.0.0.0/24 10.0.4.0/24 10.0.8.0/24`,
		}))
	})

	It("specifies the appropriate prerequisites when running the script", func() {
		Expect(scriptRunner.RunScriptCallCount()).To(Equal(1))

		_, prereqs, _ := scriptRunner.RunScriptArgsForCall(0)

		Expect(prereqs).To(ConsistOf("jq", "om", "sshuttle"))
	})

	When("run with dry run set to false", func() {
		BeforeEach(func() {
			dryRun = false
		})

		It("runs the script in dry run mode", func() {
			Expect(scriptRunner.RunScriptCallCount()).To(Equal(1))

			_, _, dryRun := scriptRunner.RunScriptArgsForCall(0)
			Expect(dryRun).To(Equal(false))
		})
	})

	When("run with dry run set to true", func() {
		BeforeEach(func() {
			dryRun = true
		})

		It("runs the script in dry run mode", func() {
			Expect(scriptRunner.RunScriptCallCount()).To(Equal(1))

			_, _, dryRun := scriptRunner.RunScriptArgsForCall(0)
			Expect(dryRun).To(Equal(true))
		})
	})

	When("running the script succeeds", func() {
		It("doesn't error", func() {
			Expect(err).NotTo(HaveOccurred())
		})
	})

	When("running the script errors", func() {
		BeforeEach(func() {
			scriptRunner.RunScriptReturns(fmt.Errorf("run-script-error"))
		})

		It("propagates the error", func() {
			Expect(err).To(MatchError("run-script-error"))
		})
	})
})
