package ssh_test

import (
	"fmt"
	"net"
	"net/url"

	"github.com/pivotal/hammer/ssh"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pivotal/hammer/environment"
	"github.com/pivotal/hammer/scripting/scriptingfakes"
)

var _ = Describe("ssh runner", func() {
	var (
		err          error
		sshRunner    ssh.Runner
		scriptRunner *scriptingfakes.FakeScriptRunner

		data   environment.Config
		dryRun bool
	)

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

		sshRunner = ssh.Runner{
			ScriptRunner: scriptRunner,
		}
	})

	JustBeforeEach(func() {
		err = sshRunner.Run(data, dryRun)
	})

	It("runs the script with an ssh to the opsman vm", func() {
		Expect(scriptRunner.RunScriptCallCount()).To(Equal(1))

		lines, _, _ := scriptRunner.RunScriptArgsForCall(0)
		Expect(lines).To(Equal([]string{
			`ssh_key_path=$(mktemp)`,
			`echo "private-key-contents" >"$ssh_key_path"`,
			`trap 'rm -f ${ssh_key_path}' EXIT`,
			`chmod 0600 "${ssh_key_path}"`,
			`creds="$(om -t www.test-url.io -k -u username -p password curl -s -p /api/v0/deployed/director/credentials/bosh_commandline_credentials)"`,
			`bosh="$(echo "$creds" | jq -r .credential | tr ' ' '\n' | grep '=')"`,
			`echo "$bosh"`,
			`shell="/usr/bin/env $(echo $bosh | tr '\n' ' ') bash -l"`,
			`ssh -o IdentitiesOnly=yes -o StrictHostKeyChecking=no -i "${ssh_key_path}" -t ubuntu@"10.0.0.6" "$shell"`,
		}))
	})

	It("specifies the appropriate prerequisites when running the script", func() {
		Expect(scriptRunner.RunScriptCallCount()).To(Equal(1))

		_, prereqs, _ := scriptRunner.RunScriptArgsForCall(0)

		Expect(prereqs).To(ConsistOf("ssh", "om"))
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
